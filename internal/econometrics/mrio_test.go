package econometrics

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestEconometricsMRIOAndShock(t *testing.T) {
	engine := NewEngine()

	project := &domain.Project{
		ID:       "p-darlington",
		Name:     "Darlington SMR Project",
		Sector:   domain.SectorNuclearEnergy,
		CapexCAD: 7_700_000_000, // $7.7B CAD
	}

	t.Run("CalculateMRIO generates direct indirect induced breakdown", func(t *testing.T) {
		mrio := engine.CalculateMRIO(project)
		if mrio.TotalGDPCAD <= mrio.DirectGDPCAD {
			t.Fatalf("Total GDP ($%d) must exceed direct GDP ($%d)", mrio.TotalGDPCAD, mrio.DirectGDPCAD)
		}
		if mrio.TotalMultipler < 1.5 {
			t.Fatalf("Expected nuclear multiplier > 1.5, got %.2f", mrio.TotalMultipler)
		}
		if mrio.PersonYearsJobs <= 0 {
			t.Fatalf("Expected positive person-years jobs, got %d", mrio.PersonYearsJobs)
		}
		if mrio.TotalFiscalReturn <= 0 {
			t.Fatalf("Expected positive total fiscal return, got %d", mrio.TotalFiscalReturn)
		}
	})

	t.Run("EvaluateShock captures interest rate and FX inflation", func(t *testing.T) {
		shock := SensitivityShockRequest{
			InterestRateDeltaBps: 150,
			FXDepreciationCAD:    -0.05,
			CommodityPriceDelta:  -0.10,
			ImportedCapexShare:   0.40,
		}
		res := engine.EvaluateShock(project, shock)
		if res.AdjustedCapexCAD <= res.OriginalCapexCAD {
			t.Fatalf("Adjusted capex must increase due to FX depreciation on imported items")
		}
		if res.DebtServiceDeltaCAD <= 0 {
			t.Fatalf("Debt service delta must be positive with +150bps interest rate increase")
		}
		if res.ProjectViabilityShift >= 0 {
			t.Fatalf("Viability shift must be negative under adverse macro conditions")
		}
	})
}
