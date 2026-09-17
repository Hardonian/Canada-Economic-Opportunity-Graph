package lakehouse

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// DataType represents supported columnar primitive types.
type DataType string

const (
	TypeString  DataType = "STRING"
	TypeInt64   DataType = "INT64"
	TypeFloat64 DataType = "FLOAT64"
	TypeBool    DataType = "BOOL"
	TypeTime    DataType = "TIMESTAMP"
)

// FieldSchema defines a single column in an Iceberg/Arrow table.
type FieldSchema struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Type     DataType `json:"type"`
	Required bool     `json:"required"`
}

// TableSchema defines the full schema of a lakehouse table.
type TableSchema struct {
	SchemaID int           `json:"schema_id"`
	Fields   []FieldSchema `json:"fields"`
}

// PartitionSpec defines partitioning keys (e.g., province, sector).
type PartitionSpec struct {
	SpecID int      `json:"spec_id"`
	Fields []string `json:"fields"`
}

// Snapshot represents an immutable point-in-time state of an Iceberg table.
type Snapshot struct {
	SnapshotID   int64     `json:"snapshot_id"`
	ParentID     int64     `json:"parent_snapshot_id,omitempty"`
	TimestampMs  int64     `json:"timestamp_ms"`
	ManifestList string    `json:"manifest_list"`
	RecordCount  int64     `json:"record_count"`
	ContentHash  string    `json:"content_hash"`
}

// IcebergMetadata stores table metadata according to Apache Iceberg format v2.
type IcebergMetadata struct {
	TableUUID        string        `json:"table_uuid"`
	FormatVersion    int           `json:"format_version"`
	Location         string        `json:"location"`
	CurrentSchemaID  int           `json:"current_schema_id"`
	Schemas          []TableSchema `json:"schemas"`
	PartitionSpec    PartitionSpec `json:"partition_spec"`
	CurrentSnapshot  int64         `json:"current_snapshot_id"`
	Snapshots        []Snapshot    `json:"snapshots"`
}

// RecordBatch represents an in-memory columnar slice (Arrow-compatible).
type RecordBatch struct {
	Schema  TableSchema
	Columns map[string][]interface{}
	Length  int
}

// NewRecordBatch initializes an empty columnar record batch.
func NewRecordBatch(schema TableSchema) *RecordBatch {
	cols := make(map[string][]interface{})
	for _, f := range schema.Fields {
		cols[f.Name] = make([]interface{}, 0)
	}
	return &RecordBatch{
		Schema:  schema,
		Columns: cols,
		Length:  0,
	}
}

// AppendRow appends a row mapping column name to value.
func (rb *RecordBatch) AppendRow(row map[string]interface{}) error {
	for _, f := range rb.Schema.Fields {
		val, exists := row[f.Name]
		if !exists && f.Required {
			return fmt.Errorf("missing required column %s", f.Name)
		}
		rb.Columns[f.Name] = append(rb.Columns[f.Name], val)
	}
	rb.Length++
	return nil
}

// CDCOperation represents row-level change capture operations.
type CDCOperation string

const (
	CDCOpInsert CDCOperation = "INSERT"
	CDCOpUpdate CDCOperation = "UPDATE"
	CDCOpDelete CDCOperation = "DELETE"
)

// CDCEvent represents a streaming Change Data Capture record.
type CDCEvent struct {
	EventID      string                 `json:"event_id"`
	SourceSystem string                 `json:"source_system"`
	Table        string                 `json:"table"`
	Operation    CDCOperation           `json:"operation"`
	Before       map[string]interface{} `json:"before,omitempty"`
	After        map[string]interface{} `json:"after,omitempty"`
	EmittedAt    time.Time              `json:"emitted_at"`
	TxID         string                 `json:"tx_id"`
}

// ComputeHash generates SHA-256 for CDC event verification.
func (e *CDCEvent) ComputeHash() string {
	payload := fmt.Sprintf("%s|%s|%s|%s|%d|%s",
		e.EventID, e.SourceSystem, e.Table, e.Operation, e.EmittedAt.UnixNano(), e.TxID)
	h := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(h[:])
}

// QuarantineRecord represents a row that violated data quality SLAs.
type QuarantineRecord struct {
	RecordID    string                 `json:"record_id"`
	SourceTable string                 `json:"source_table"`
	RawData     map[string]interface{} `json:"raw_data"`
	Violations  []string               `json:"violations"`
	Severity    string                 `json:"severity"` // "WARNING", "CRITICAL"
	DetectedAt  time.Time              `json:"detected_at"`
}
