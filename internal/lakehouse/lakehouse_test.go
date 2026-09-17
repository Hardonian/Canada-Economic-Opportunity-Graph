package lakehouse

import (
	"testing"
)

func TestVectorQueryEngine(t *testing.T) {
	schema := TableSchema{
		SchemaID: 1,
		Fields: []FieldSchema{
			{ID: 1, Name: "id", Type: TypeString, Required: true},
			{ID: 2, Name: "province", Type: TypeString, Required: true},
			{ID: 3, Name: "sector", Type: TypeString, Required: true},
			{ID: 4, Name: "capex_cad", Type: TypeFloat64, Required: true},
		},
	}

	rb := NewRecordBatch(schema)
	rows := []map[string]interface{}{
		{"id": "p1", "province": "ON", "sector": "Nuclear", "capex_cad": 3400000000.0},
		{"id": "p2", "province": "ON", "sector": "Nuclear", "capex_cad": 12500000000.0},
		{"id": "p3", "province": "QC", "sector": "Mining", "capex_cad": 1800000000.0},
		{"id": "p4", "province": "SK", "sector": "Mining", "capex_cad": 650000000.0},
		{"id": "p5", "province": "ON", "sector": "Clean Energy", "capex_cad": 750000000.0},
	}

	for _, r := range rows {
		if err := rb.AppendRow(r); err != nil {
			t.Fatalf("failed to append row: %v", err)
		}
	}

	engine := NewVectorQueryEngine()

	// Query: Group by province, sum capex_cad, count projects, order by sum desc
	req := QueryRequest{
		GroupBy: []string{"province"},
		Aggregations: []AggregationSpec{
			{Column: "capex_cad", Type: AggSum, Alias: "total_capex"},
			{Column: "id", Type: AggCount, Alias: "project_count"},
		},
		OrderBy:   "total_capex",
		OrderDesc: true,
	}

	res, err := engine.Execute(rb, req)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if len(res.Rows) != 3 {
		t.Fatalf("expected 3 province groups, got %d", len(res.Rows))
	}

	// First row should be ON (16.65B)
	topRow := res.Rows[0]
	if topRow["province"] != "ON" {
		t.Errorf("expected ON as top capex province, got %v", topRow["province"])
	}
	expectedONCapex := 3400000000.0 + 12500000000.0 + 750000000.0
	if topRow["total_capex"] != expectedONCapex {
		t.Errorf("expected ON capex %v, got %v", expectedONCapex, topRow["total_capex"])
	}
	if topRow["project_count"] != 3 {
		t.Errorf("expected 3 projects in ON, got %v", topRow["project_count"])
	}
}

func TestIcebergMetadata(t *testing.T) {
	itm := NewIcebergTableManager()
	schema := TableSchema{
		SchemaID: 1,
		Fields: []FieldSchema{
			{ID: 1, Name: "id", Type: TypeString, Required: true},
			{ID: 2, Name: "name", Type: TypeString, Required: true},
		},
	}
	partSpec := PartitionSpec{SpecID: 1, Fields: []string{"province"}}

	meta, err := itm.CreateTable("sovereign_projects", "s3://lakehouse/projects", schema, partSpec)
	if err != nil {
		t.Fatalf("create table err: %v", err)
	}
	if meta.FormatVersion != 2 {
		t.Errorf("expected format version 2, got %d", meta.FormatVersion)
	}

	snap, err := itm.CommitSnapshot("sovereign_projects", 299, "s3://lakehouse/manifests/snap1.avro")
	if err != nil {
		t.Fatalf("commit snapshot err: %v", err)
	}
	if snap.RecordCount != 299 {
		t.Errorf("expected 299 records, got %d", snap.RecordCount)
	}

	metaUpdated, _ := itm.GetMetadata("sovereign_projects")
	if metaUpdated.CurrentSnapshot != snap.SnapshotID {
		t.Errorf("expected current snapshot to match committed snapshot ID")
	}
}

func TestQuarantineEngine(t *testing.T) {
	qe := NewQuarantineEngine()

	// Valid project
	valid := map[string]interface{}{
		"id":        "proj-1",
		"name":      "Darlington SMR",
		"capex_cad": 3400000000.0,
		"latitude":  43.8688,
	}
	ok, _ := qe.Validate("projects", valid)
	if !ok {
		t.Errorf("expected valid row to pass quarantine check")
	}

	// Invalid project (negative capex, invalid latitude)
	invalid := map[string]interface{}{
		"id":        "proj-2",
		"name":      "Broken Proj",
		"capex_cad": -500.0,
		"latitude":  120.0, // Invalid latitude
	}
	ok2, qr := qe.Validate("projects", invalid)
	if ok2 {
		t.Errorf("expected invalid row to be quarantined")
	}
	if len(qr.Violations) < 2 {
		t.Errorf("expected at least 2 violations, got %d", len(qr.Violations))
	}

	records := qe.GetQuarantinedRecords()
	if len(records) != 1 {
		t.Errorf("expected 1 quarantined record, got %d", len(records))
	}
}
