package sedar

import (
	"context"
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestSEDARAdapter_FetchAndParse(t *testing.T) {
	adapter := NewSEDARAdapter(nil)

	if adapter.Name() != "sedar_plus_disclosures" {
		t.Fatalf("expected adapter name sedar_plus_disclosures, got %s", adapter.Name())
	}
	if adapter.Tier() != domain.SourceTier1 {
		t.Fatalf("expected tier 1, got %v", adapter.Tier())
	}

	data, err := adapter.Fetch(context.Background())
	if err != nil {
		t.Fatalf("fetch failed: %v", err)
	}

	res, err := adapter.Parse(data)
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if len(res.Projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(res.Projects))
	}
	if len(res.Entities) != 2 {
		t.Fatalf("expected 2 entities, got %d", len(res.Entities))
	}
	if len(res.CapitalItems) != 2 {
		t.Fatalf("expected 2 capital items, got %d", len(res.CapitalItems))
	}
	if len(res.Evidence) != 2 {
		t.Fatalf("expected 2 evidence records, got %d", len(res.Evidence))
	}

	// Verify CNC project details
	p1 := res.Projects[0]
	if p1.CapexCAD != 2_500_000_000 {
		t.Fatalf("expected $2.5B capex, got %d", p1.CapexCAD)
	}
	if p1.Sector != domain.SectorCriticalMinerals {
		t.Fatalf("expected SectorCriticalMinerals, got %v", p1.Sector)
	}

	// Verify health
	health := adapter.Health()
	if health.Status != "HEALTHY" {
		t.Fatalf("expected HEALTHY health status, got %s", health.Status)
	}
	if health.DocumentsSeen != 2 {
		t.Fatalf("expected 2 documents seen, got %d", health.DocumentsSeen)
	}
}

func TestExtractCapexFromText(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"Project requires $1.2 billion for phase 1", 1_200_000_000},
		{"Total estimated cost is CAD 450 million", 450_000_000},
		{"C$ 85M capital budget allocated", 85_000_000},
		{"No financial numbers mentioned", 0},
	}

	for _, tt := range tests {
		got := extractCapexFromText(tt.input)
		if got != tt.want {
			t.Errorf("extractCapexFromText(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
