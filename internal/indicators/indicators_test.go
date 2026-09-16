package indicators_test

import (
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/indicators"
)

func TestCanonicalRegistry(t *testing.T) {
	registry := indicators.CanonicalRegistry()
	if len(registry) < 32 {
		t.Fatalf("expected at least 32 registered KPIs, got %d", len(registry))
	}

	codeMap := indicators.RegistryByCode()
	if len(codeMap) != len(registry) {
		t.Fatalf("duplicate KPI code found in registry: map len %d vs list len %d", len(codeMap), len(registry))
	}

	categoriesFound := make(map[indicators.KPICategory]int)
	for _, kpi := range registry {
		if kpi.Code == "" {
			t.Errorf("KPI has empty code: %+v", kpi)
		}
		if kpi.Name == "" {
			t.Errorf("KPI %s has empty name", kpi.Code)
		}
		if kpi.Unit == "" {
			t.Errorf("KPI %s has empty unit", kpi.Code)
		}
		if kpi.StatutoryBasis == "" {
			t.Errorf("KPI %s missing statutory basis", kpi.Code)
		}
		categoriesFound[kpi.Category]++
	}

	expectedCategories := []indicators.KPICategory{
		indicators.CategoryESGDecarbonization,
		indicators.CategoryIndigenousEquity,
		indicators.CategoryCapitalVelocity,
		indicators.CategorySupplyChainContent,
		indicators.CategoryGridPhysics,
		indicators.CategoryRegulatorySpeed,
		indicators.CategoryLaborSkills,
		indicators.CategoryCommodityMacro,
	}

	for _, cat := range expectedCategories {
		if categoriesFound[cat] == 0 {
			t.Errorf("category %s has no registered KPIs", cat)
		}
	}
}

func TestLiveFeedEngine(t *testing.T) {
	engine := indicators.NewLiveFeedEngine()

	// Verify initial ticks seeded
	ticks := engine.GetLatestTicks()
	if len(ticks) < 8 {
		t.Fatalf("expected at least 8 seeded live feeds, got %d", len(ticks))
	}

	// Verify recent ticks limit
	recent := engine.GetRecentTicks(5)
	if len(recent) != 5 {
		t.Errorf("expected 5 recent ticks, got %d", len(recent))
	}

	// Record manual tick
	now := time.Now().UTC()
	tick := engine.RecordTick(
		indicators.KPIWCSWTIDifferential,
		"WCS Differential",
		indicators.CategoryCommodityMacro,
		13.10,
		"USD/bbl",
		0.25,
		1.94,
		"Test Exchange",
		now,
	)
	if tick.SequenceID <= 0 {
		t.Errorf("expected positive sequence ID, got %d", tick.SequenceID)
	}
	if tick.Direction != "UP" {
		t.Errorf("expected UP direction for positive change, got %s", tick.Direction)
	}
	if tick.Hash == "" {
		t.Error("expected non-empty tick hash")
	}

	// Test simulated drift tick
	simTick := engine.GenerateSimulatedTick(indicators.KPIWCSWTIDifferential)
	if simTick == nil {
		t.Fatal("expected simulated tick to generate successfully")
	}
	if simTick.Value <= 0 {
		t.Errorf("expected positive simulated tick value, got %f", simTick.Value)
	}

	// Test observations
	nationalObs := engine.GetObservations("NATIONAL")
	if len(nationalObs) == 0 {
		t.Error("expected non-empty national observations")
	}

	// Test Macro Summary
	summary := engine.BuildMacroSummary()
	if summary == nil {
		t.Fatal("expected macro summary to be built")
	}
	if summary.AverageIndigenousEquityPct <= 0 {
		t.Errorf("expected positive average indigenous equity, got %f", summary.AverageIndigenousEquityPct)
	}

	// Test Snapshot Generation
	snapshot, err := engine.GenerateSnapshot()
	if err != nil {
		t.Fatalf("failed to generate snapshot: %v", err)
	}
	if len(snapshot.Definitions) < 32 {
		t.Errorf("snapshot missing definitions: %d", len(snapshot.Definitions))
	}
	if len(snapshot.Observations) == 0 {
		t.Error("snapshot missing observations")
	}
	if snapshot.AuditHash == "" {
		t.Error("snapshot missing audit hash")
	}
}

func TestProjectEvaluator(t *testing.T) {
	evaluator := indicators.NewProjectEvaluator()

	testProjects := []*domain.Project{
		{
			ID:           "proj-darlington-smr",
			Slug:         "darlington-new-nuclear-project-unit-1",
			Name:         "Darlington New Nuclear Project (Unit 1)",
			Sector:       domain.SectorNuclearEnergy,
			Province:     "ON",
			CurrentStage: domain.StageConstruction,
			CapexCAD:     3_400_000_000,
		},
		{
			ID:           "proj-crawford-nickel",
			Slug:         "crawford-nickel-cobalt-project",
			Name:         "Crawford Nickel Sulphide Project",
			Sector:       domain.SectorCriticalMinerals,
			Province:     "ON",
			CurrentStage: domain.StagePermitting,
			CapexCAD:     1_900_000_000,
		},
		{
			ID:           "proj-oneida-battery",
			Slug:         "oneida-energy-storage-facility",
			Name:         "Oneida Energy Storage Facility",
			Sector:       domain.SectorCleanEnergy,
			Province:     "ON",
			CurrentStage: domain.StageOperating,
			CapexCAD:     600_000_000,
		},
	}

	for _, p := range testProjects {
		scorecard := evaluator.Evaluate(p)
		if scorecard == nil {
			t.Fatalf("evaluator returned nil for project %s", p.Name)
		}
		if scorecard.OverallKPIRating <= 0 || scorecard.OverallKPIRating > 100 {
			t.Errorf("invalid overall rating %f for project %s", scorecard.OverallKPIRating, p.Name)
		}
		if len(scorecard.Pillars) != 8 {
			t.Errorf("expected 8 pillars, got %d for project %s", len(scorecard.Pillars), p.Name)
		}
		if scorecard.AuditHash == "" {
			t.Errorf("missing audit hash for project %s", p.Name)
		}
		for _, pillar := range scorecard.Pillars {
			if len(pillar.Metrics) == 0 {
				t.Errorf("pillar %s has no metrics for project %s", pillar.Title, p.Name)
			}
			if pillar.Score < 0 || pillar.Score > 100 {
				t.Errorf("pillar %s has out-of-range score %f for project %s", pillar.Title, pillar.Score, p.Name)
			}
		}
	}
}
