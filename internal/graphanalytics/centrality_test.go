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

	t.Run("ISO GQL Parser and Execution", func(t *testing.T) {
		q, err := ParseGQL("MATCH (p:Project)-[r:intertie]->(t:Project) WHERE p.province = 'ON'")
		if err != nil {
			t.Fatalf("failed to parse GQL: %v", err)
		}
		if q.RelType != "INTERTIE" {
			t.Errorf("expected rel type INTERTIE, got %s", q.RelType)
		}

		res := ExecuteGQL(q, projects, relationships)
		if res.ExecutionMs < 0 {
			t.Errorf("invalid execution time")
		}
	})

	t.Run("Graph Traversal Engine", func(t *testing.T) {
		adj := map[string][]string{
			"p1": {"p2"},
			"p2": {"p3"},
			"p3": {"p4"},
		}
		gte := NewGraphTraversalEngine(adj)

		path, err := gte.FindShortestPath("p1", "p4")
		if err != nil {
			t.Fatalf("shortest path err: %v", err)
		}
		if len(path) != 4 || path[0] != "p1" || path[3] != "p4" {
			t.Errorf("unexpected path: %v", path)
		}

		neighbors := gte.GetNeighborsK("p1", 2)
		if len(neighbors) != 2 {
			t.Errorf("expected 2 neighbors within 2 hops, got %d", len(neighbors))
		}
	})

	t.Run("GNN Link Prediction", func(t *testing.T) {
		predictor := NewGNNLinkPredictor()
		testProjs := []*domain.Project{
			{ID: "smr-1", Name: "Darlington SMR", Sector: domain.SectorCleanEnergy, Province: "ON", CapexCAD: 3400000000},
			{ID: "ai-1", Name: "Bruce AI Compute", Sector: domain.SectorAICompute, Province: "ON", CapexCAD: 1500000000},
		}

		links := predictor.PredictPartnerships(testProjs, 0.5)
		if len(links) == 0 {
			t.Fatalf("expected predicted link between SMR and AI compute, got 0")
		}
		if links[0].PredictedRelType != "CLEAN_POWER_PPA" {
			t.Errorf("expected CLEAN_POWER_PPA, got %s", links[0].PredictedRelType)
		}
	})

	t.Run("Louvain Community Detection", func(t *testing.T) {
		clusterer := NewLouvainClusterer()
		nodes := []string{"p1", "p2", "p3", "p4", "p5"}
		edges := map[string][]string{
			"p1": {"p2"},
			"p2": {"p1", "p3"},
			"p3": {"p2"},
			"p4": {"p5"},
			"p5": {"p4"},
		}

		clusters := clusterer.DetectCommunities(nodes, edges)
		if len(clusters) != 2 {
			t.Fatalf("expected 2 communities, got %d", len(clusters))
		}
	})

	t.Run("HyperGraph Management", func(t *testing.T) {
		hgm := NewHyperGraphManager()
		edge := &HyperEdge{
			ID:      "he-syndicate-1",
			Name:    "Cedar LNG First Nations Syndicate",
			Type:    "MULTI_NATION_EQUITY",
			NodeIDs: []string{"proj-cedar-lng", "haisla-nation", "pembina-pipeline", "cib-guarantor"},
		}

		if err := hgm.AddHyperEdge(edge); err != nil {
			t.Fatalf("failed to add hyper edge: %v", err)
		}

		edges := hgm.GetHyperEdgesForNode("haisla-nation")
		if len(edges) != 1 || edges[0].ID != "he-syndicate-1" {
			t.Errorf("expected node to participate in he-syndicate-1")
		}
	})
}
