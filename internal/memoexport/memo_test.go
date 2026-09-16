package memoexport

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestMemoGenerator(t *testing.T) {
	generator := NewMemoGenerator()

	proj := &domain.Project{
		ID:       "proj-airport-hub",
		Slug:     "canadian-international-airports-global-investor-leasing-infrastructure-hubs",
		Name:     "Canadian International Airports — Global Investor Leasing & Infrastructure Hubs",
		Sector:   domain.SectorTransportation,
		Province: "ON",
		CapexCAD: 18_000_000_000,
	}

	memo := generator.GenerateCabinetMemo(proj, MemoTypeCabinetMC)
	if memo == nil {
		t.Fatal("expected non-nil decision memo")
	}
	if memo.SecurityCaveat == "" {
		t.Errorf("expected security caveat")
	}
	if memo.MarkdownContent == "" {
		t.Errorf("expected markdown content")
	}
	if memo.AuditHash == "" {
		t.Errorf("expected audit hash")
	}
}

func TestGeoJSONAndSTACExport(t *testing.T) {
	proj := &domain.Project{
		ID:        "proj-1",
		Slug:      "test-mine",
		Name:      "Crawford Nickel Project",
		Sector:    domain.SectorCriticalMinerals,
		Province:  "ON",
		Latitude:  48.72,
		Longitude: -81.25,
		CapexCAD:  2_100_000_000,
	}

	coll := ExportGeoJSON([]*domain.Project{proj})
	if coll == nil {
		t.Fatal("expected non-nil GeoJSON collection")
	}
	if len(coll.Features) != 1 {
		t.Fatalf("expected 1 feature, got %d", len(coll.Features))
	}
	if coll.Features[0].Geometry.Coordinates[0] != -81.25 {
		t.Errorf("expected longitude -81.25, got %f", coll.Features[0].Geometry.Coordinates[0])
	}

	stac := ExportSTACItem(proj)
	if stac == nil {
		t.Fatal("expected non-nil STAC item")
	}
	if stac.STACVersion != "1.0.0" {
		t.Errorf("expected STAC version 1.0.0, got %s", stac.STACVersion)
	}
}
