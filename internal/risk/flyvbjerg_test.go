package risk

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestFlyvbjergRiskForecast(t *testing.T) {
	evaluator := NewEvaluator()

	project := &domain.Project{
		ID:           "p-darlington-smr",
		Name:         "Darlington SMR",
		Sector:       domain.SectorNuclearEnergy,
		Subsector:    "BWRX-300 SMR",
		LocationName: "Bowmanville, ON",
		Province:     "ON",
		CapexCAD:     7_700_000_000,
	}

	forecast := evaluator.ForecastProject(project)

	if forecast.ExpectedCostOverrunPct <= 0 {
		t.Fatalf("Expected positive cost overrun, got %.2f%%", forecast.ExpectedCostOverrunPct)
	}

	if len(forecast.Percentiles) != 4 {
		t.Fatalf("Expected 4 percentile points, got %d", len(forecast.Percentiles))
	}

	p10 := forecast.Percentiles[0]
	p90 := forecast.Percentiles[3]

	if p90.ForecastCapexCAD <= p10.ForecastCapexCAD {
		t.Fatalf("P90 capex ($%d) must exceed P10 capex ($%d)", p90.ForecastCapexCAD, p10.ForecastCapexCAD)
	}

	if p90.ScheduleDelayMonths <= p10.ScheduleDelayMonths {
		t.Fatalf("P90 delay (%d mos) must exceed P10 delay (%d mos)", p90.ScheduleDelayMonths, p10.ScheduleDelayMonths)
	}

	if forecast.TechNoveltyPenalty <= 0 {
		t.Fatalf("Expected tech novelty penalty for SMR")
	}
}
