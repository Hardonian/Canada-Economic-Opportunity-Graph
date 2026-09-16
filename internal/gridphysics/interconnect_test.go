package gridphysics

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestGridInterconnectAssessment(t *testing.T) {
	engine := NewEngine()

	project := &domain.Project{
		ID:       "p-chisasibi-ai",
		Name:     "Chisasibi Sovereign AI Compute",
		Sector:   domain.SectorAICompute,
		Province: "QC",
		CapexCAD: 2_800_000_000,
	}

	assessment := engine.AssessProject(project)

	if assessment.Operator != OperatorHydroQuebec {
		t.Fatalf("Expected Hydro-Quebec, got %s", assessment.Operator)
	}

	if assessment.EstimatedLoadOrGenMW <= 0 {
		t.Fatalf("Expected positive load MW, got %.1f", assessment.EstimatedLoadOrGenMW)
	}

	if !assessment.DedicatedSubstationNeeded {
		t.Fatal("150MW AI data centre must require dedicated substation")
	}

	if assessment.ReinforcementCostCAD <= 0 {
		t.Fatalf("Expected positive reinforcement cost, got $%d", assessment.ReinforcementCostCAD)
	}

	if assessment.CleanPowerPurityPct < 90.0 {
		t.Fatalf("Expected Hydro-Quebec purity > 90%%, got %.1f%%", assessment.CleanPowerPurityPct)
	}
}
