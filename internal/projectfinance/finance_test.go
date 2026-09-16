package projectfinance

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestFinanceSimulatorRunSimulation(t *testing.T) {
	simulator := NewFinanceSimulator()

	proj := &domain.Project{
		ID:       "test-project-finance",
		Name:     "Ontario Clean Hyperscale Data Campus",
		Sector:   domain.SectorAICompute,
		CapexCAD: 1_200_000_000,
	}

	result := simulator.RunSimulation(proj, 500)
	if result == nil {
		t.Fatal("expected non-nil simulation result")
	}
	if result.IterationsRun != 500 {
		t.Errorf("expected 500 iterations, got %d", result.IterationsRun)
	}
	if result.ProjectIRRPercent.P50 <= 0 {
		t.Errorf("expected positive median project IRR, got %f", result.ProjectIRRPercent.P50)
	}
	if result.AvgDSCR.P50 <= 0 {
		t.Errorf("expected positive median DSCR, got %f", result.AvgDSCR.P50)
	}
	if result.SyntheticCreditRating == "" {
		t.Errorf("expected synthetic credit rating")
	}
	if result.AuditHash == "" {
		t.Errorf("expected non-empty audit hash")
	}
}

func TestTaxCreditCalculator(t *testing.T) {
	calc := NewTaxCreditCalculator()

	proj := &domain.Project{
		ID:        "test-clean-itc",
		Name:      "BWRX-300 SMR Unit",
		Sector:    domain.SectorNuclearEnergy,
		Subsector: "Nuclear Power",
		CapexCAD:  2_000_000_000,
	}

	profile := calc.CalculateCredits(proj)
	if profile == nil {
		t.Fatal("expected non-nil tax credit profile")
	}
	if profile.ApplicableITC != ITCCleanElectricity {
		t.Errorf("expected Clean Electricity ITC, got %s", profile.ApplicableITC)
	}
	if profile.EffectiveCreditRatePercent != 15.0 {
		t.Errorf("expected 15 percent credit rate, got %f", profile.EffectiveCreditRatePercent)
	}
	if profile.TotalTaxCreditYieldCAD <= 0 {
		t.Errorf("expected positive tax credit yield, got %d", profile.TotalTaxCreditYieldCAD)
	}
	if profile.AuditHash == "" {
		t.Errorf("expected non-empty audit hash")
	}
}
