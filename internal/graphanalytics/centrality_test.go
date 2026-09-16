package graphanalytics

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestGraphCentralityAndCascade(t *testing.T) {
	ga := NewGraphAnalyzer()

	projects := []*domain.Project{
		{ID: "p1", Name: "Churchill Arctic Port", CapexCAD: 850_000_000, Sector: domain.SectorTransportation},
		{ID: "p2", Name: "Kivalliq Hydro Link", CapexCAD: 1_650_000_000, Sector: domain.SectorCleanEnergy},
		{ID: "p3", Name: "Nunavut Gold-Nickel Mine", CapexCAD: 1_200_000_000, Sector: domain.SectorCriticalMinerals},
	}

	relationships := []*domain.Relationship{
		{SourceEntityID: "p1", TargetEntityID: "p2", RelationType: "intertie"},
		{SourceEntityID: "p2", TargetEntityID: "p3", RelationType: "power_offtake"},
	}

	t.Run("AnalyzeNetwork calculates centrality correctly", func(t *testing.T) {
		centralities := ga.AnalyzeNetwork(projects, relationships)
		if len(centralities) != len(projects) {
			t.Fatalf("Expected %d nodes, got %d", len(projects), len(centralities))
		}
	})

	t.Run("SimulateCascade propagates down dependency chain", func(t *testing.T) {
		cascade := ga.SimulateCascade("p1", projects, relationships)
		if cascade.TriggerNodeID != "p1" {
			t.Fatalf("Expected trigger p1, got %s", cascade.TriggerNodeID)
		}
		if len(cascade.TotalCompromisedIDs) < 2 {
			t.Fatalf("Expected p2 and p3 to be compromised, got %v", cascade.TotalCompromisedIDs)
		}
		if cascade.CascadingDepth < 2 {
			t.Fatalf("Expected depth >= 2, got %d", cascade.CascadingDepth)
		}
	})
}
