package nationalplanning

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestNationalPlanningSuite(t *testing.T) {
	projects := []*domain.Project{
		{
			ID:           "p-darlington",
			Name:         "Darlington SMR Deployment",
			Sector:       domain.SectorNuclearEnergy,
			Province:     "ON",
			CapexCAD:     7_700_000_000,
			Scores:       map[string]float64{"strategicity": 95.0, "buildability": 88.0},
		},
		{
			ID:           "p-crawford",
			Name:         "Crawford Nickel Project",
			Sector:       domain.SectorCriticalMinerals,
			Province:     "ON",
			CapexCAD:     3_500_000_000,
			Scores:       map[string]float64{"strategicity": 90.0, "buildability": 75.0},
		},
		{
			ID:           "p-chisasibi",
			Name:         "Chisasibi AI Data Centre",
			Sector:       domain.SectorAICompute,
			Province:     "QC",
			CapexCAD:     2_800_000_000,
			Scores:       map[string]float64{"strategicity": 92.0, "buildability": 80.0},
		},
	}

	t.Run("Optimizer allocates public capital respecting envelopes", func(t *testing.T) {
		opt := NewOptimizer()
		req := OptimizationRequest{
			Objective: ObjectiveBalancedStrategy,
			Envelopes: BudgetEnvelopes{
				CIBConcessionaryCAD: 2_000_000_000,
				SIFGrantsCAD:        1_000_000_000,
				ITCTaxCreditsCAD:    1_500_000_000,
				IndigenousLoansCAD:  500_000_000,
			},
		}
		res := opt.Optimize(projects, req)

		if len(res.AllocatedProjects) == 0 {
			t.Fatal("Expected at least one allocated project")
		}
		if res.TotalPublicInvestedCAD <= 0 {
			t.Fatal("Expected positive public investment")
		}
		if res.CrowdingInMultiplier <= 1.0 {
			t.Fatalf("Expected crowding-in multiplier > 1.0x, got %.2fx", res.CrowdingInMultiplier)
		}
	})

	t.Run("WarGameEngine simulates USMCA 25% tariff shock", func(t *testing.T) {
		wge := NewWarGameEngine()
		req := WarGameRequest{
			Scenario: ShockUSMCATariffs,
		}
		res := wge.SimulateScenario(projects, req)

		if res.TotalAssetsStalledCount == 0 {
			t.Fatal("Expected critical minerals/manufacturing assets to be affected")
		}
		if res.TotalFrozenCapexCAD <= 0 {
			t.Fatal("Expected positive frozen capex")
		}
		if len(res.SovereignMitigations) == 0 {
			t.Fatal("Expected sovereign mitigations")
		}
	})

	t.Run("LaborAggregator detects regional trades pinch-points", func(t *testing.T) {
		la := NewLaborAggregator()
		rep := la.AnalyzeProvince("ON", projects)

		if rep.TotalLaborDemandFTE <= 0 {
			t.Fatal("Expected positive total labor demand")
		}
		if len(rep.Trades) == 0 {
			t.Fatal("Expected trade breakdowns")
		}
	})

	t.Run("IndigenousModeler models equity co-investment returns", func(t *testing.T) {
		im := NewIndigenousModeler()
		p := projects[0]
		m := im.ModelEquityParticipation(p, 25.0)

		if m.IndigenousEquityCAD != int64(float64(p.CapexCAD)*0.25) {
			t.Fatalf("Expected 25%% equity, got $%d", m.IndigenousEquityCAD)
		}
		if m.AnnualNetCommunityCashYield <= 0 {
			t.Fatalf("Expected positive net community cash yield, got $%d", m.AnnualNetCommunityCashYield)
		}
		if m.ThirtyYearCumulativeWealth <= 0 {
			t.Fatalf("Expected positive 30-year wealth, got $%d", m.ThirtyYearCumulativeWealth)
		}
	})
}
