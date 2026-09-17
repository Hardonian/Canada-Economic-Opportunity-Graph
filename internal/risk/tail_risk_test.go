package risk

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestEvaluateTailRisk_NuclearMegaproject(t *testing.T) {
	evaluator := NewEvaluator()

	project := &domain.Project{
		ID:           "nuclear-darlington-smr",
		Name:         "Darlington SMR Project",
		Sector:       domain.SectorNuclearEnergy,
		Subsector:    "BWRX-300 SMR",
		LocationName: "Bowmanville, ON",
		Province:     "ON",
		CapexCAD:     7_700_000_000,
	}

	forecast := evaluator.ForecastProject(project)
	tail := evaluator.EvaluateTailRisk(forecast)

	if tail.ProjectID != project.ID {
		t.Errorf("expected project ID %s, got %s", project.ID, tail.ProjectID)
	}

	if tail.TailIndexAlpha > 1.30 {
		t.Errorf("expected fat-tailed alpha <= 1.30 for nuclear SMR, got %.2f", tail.TailIndexAlpha)
	}

	if tail.VaR99CAD <= project.CapexCAD {
		t.Errorf("VaR-99 ($%d) must exceed base Capex ($%d)", tail.VaR99CAD, project.CapexCAD)
	}

	if tail.ExpectedShortfall99CAD <= tail.VaR99CAD {
		t.Errorf("Expected Shortfall ($%d) must exceed VaR-99 ($%d)", tail.ExpectedShortfall99CAD, tail.VaR99CAD)
	}

	if tail.FragilityRating != "EXTREME_HAZARD" && tail.FragilityRating != "FRAGILE" {
		t.Errorf("expected fragile or extreme hazard rating, got %s", tail.FragilityRating)
	}

	if len(tail.DeRiskingMeasures) == 0 {
		t.Error("expected de-risking measures for nuclear project")
	}
}

func TestEvaluateTailRisk_CleanEnergy(t *testing.T) {
	evaluator := NewEvaluator()

	project := &domain.Project{
		ID:           "solar-brooks",
		Name:         "Brooks Solar Farm",
		Sector:       domain.SectorCleanEnergy,
		Subsector:    "Utility Solar PV",
		LocationName: "Brooks, AB",
		Province:     "AB",
		CapexCAD:     50_000_000,
	}

	forecast := evaluator.ForecastProject(project)
	tail := evaluator.EvaluateTailRisk(forecast)

	if tail.TailIndexAlpha < 1.50 {
		t.Errorf("expected higher alpha (thinner tail) for modular solar, got %.2f", tail.TailIndexAlpha)
	}

	if tail.ProbExceeding100Pct > 0.30 {
		t.Errorf("modular solar should have low prob of doubling cost, got %.3f", tail.ProbExceeding100Pct)
	}
}
