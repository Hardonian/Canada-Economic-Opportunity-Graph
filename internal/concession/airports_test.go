package concession

import (
	"testing"
)

func TestDefaultCanonicalAirports(t *testing.T) {
	airports := DefaultCanonicalAirports()
	if len(airports) != 5 {
		t.Fatalf("expected 5 canonical airports, got %d", len(airports))
	}

	var totalCapex int64
	expectedCodes := map[AirportHubCode]bool{
		HubYYZ: true,
		HubYVR: true,
		HubYUL: true,
		HubYYC: true,
		HubYEG: true,
	}

	for _, a := range airports {
		if !expectedCodes[a.Code] {
			t.Errorf("unexpected airport hub code: %s", a.Code)
		}
		if a.BaselineRevenueCAD <= 0 {
			t.Errorf("expected positive baseline revenue for %s, got %d", a.Code, a.BaselineRevenueCAD)
		}
		if a.BaselineEBITDACAD <= 0 {
			t.Errorf("expected positive baseline EBITDA for %s, got %d", a.Code, a.BaselineEBITDACAD)
		}
		if len(a.KeyProjects) == 0 {
			t.Errorf("expected key projects for %s", a.Code)
		}
		if len(a.DecarbonizationInitiatives) == 0 {
			t.Errorf("expected decarbonization initiatives for %s", a.Code)
		}
		if len(a.TargetPensionInvestors) == 0 {
			t.Errorf("expected target pension investors for %s", a.Code)
		}
		totalCapex += a.TargetModernizationCapexCAD
	}

	// Must match Mark Carney's $18B CAD NAS announcement exactly
	if totalCapex != 18_000_000_000 {
		t.Fatalf("expected total NAS capex to be 18,000,000,000 CAD, got %d", totalCapex)
	}
}

func TestSimulateConcessions_DefaultParams(t *testing.T) {
	params := DefaultSimulationParams()
	sim := SimulateConcessions(params)

	if sim == nil {
		t.Fatalf("expected simulation result, got nil")
	}

	if sim.TotalNASModernizationCapexCAD != 18_000_000_000 {
		t.Errorf("expected total capex 18B, got %d", sim.TotalNASModernizationCapexCAD)
	}

	// 50% pension equity -> 9B CAD
	expectedPensionEquity := int64(9_000_000_000)
	if sim.TotalPensionEquityCAD != expectedPensionEquity {
		t.Errorf("expected 9B pension equity, got %d", sim.TotalPensionEquityCAD)
	}

	// 45% commercial debt -> 8.1B CAD
	expectedDebt := int64(8_100_000_000)
	if sim.TotalCommercialDebtCAD != expectedDebt {
		t.Errorf("expected 8.1B debt, got %d", sim.TotalCommercialDebtCAD)
	}

	if sim.TotalUpfrontFederalProceedsCAD <= 0 {
		t.Errorf("expected positive upfront federal proceeds, got %d", sim.TotalUpfrontFederalProceedsCAD)
	}

	if sim.TotalCumulativeRoyaltiesCAD <= 0 {
		t.Errorf("expected positive cumulative royalties, got %d", sim.TotalCumulativeRoyaltiesCAD)
	}

	if sim.CrowdingInMultiplier <= 3.0 {
		t.Errorf("expected crowding in multiplier > 3.0x, got %.1f", sim.CrowdingInMultiplier)
	}

	if sim.PortfolioWeightedIRRPercent < 8.0 || sim.PortfolioWeightedIRRPercent > 14.0 {
		t.Errorf("expected portfolio IRR between 8%% and 14%%, got %.1f%%", sim.PortfolioWeightedIRRPercent)
	}

	if sim.PortfolioAverageDSCR < 1.3 {
		t.Errorf("expected portfolio average DSCR >= 1.3, got %.2f", sim.PortfolioAverageDSCR)
	}

	if len(sim.Maple8AllocationsCAD) == 0 {
		t.Errorf("expected Maple 8 allocations")
	}

	var maple8Total int64
	for _, val := range sim.Maple8AllocationsCAD {
		maple8Total += val
	}
	if maple8Total <= 0 {
		t.Errorf("expected positive Maple 8 total allocation, got %d", maple8Total)
	}

	if sim.AuditHash == "" {
		t.Errorf("expected non-empty audit hash")
	}
}

func TestSimulateConcessions_CustomOverrides(t *testing.T) {
	params := SimulationParams{
		ConcessionHorizonYears:      50,
		FederalRoyaltyRatePercent:   12.0,
		PensionEquitySharePercent:   60.0,
		CommercialDebtSharePercent:  35.0,
		FederalSubordinatedSharePct: 5.0,
		PassengerCAGRPercent:        3.2,
		InflationPercent:            2.2,
		DiscountRatePercent:         6.0,
		CustomCapexCAD: map[AirportHubCode]int64{
			HubYYZ: 8_500_000_000,
		},
	}

	sim := SimulateConcessions(params)
	if sim == nil {
		t.Fatalf("expected simulation result, got nil")
	}

	// 18B base - 7.5B YYZ + 8.5B YYZ = 19B total
	if sim.TotalNASModernizationCapexCAD != 19_000_000_000 {
		t.Errorf("expected total capex 19B with override, got %d", sim.TotalNASModernizationCapexCAD)
	}

	if sim.Params.ConcessionHorizonYears != 50 {
		t.Errorf("expected 50 year horizon, got %d", sim.Params.ConcessionHorizonYears)
	}
}
