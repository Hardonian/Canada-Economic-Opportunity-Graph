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

	t.Run("AC Newton-Raphson Power Flow Solver", func(t *testing.T) {
		solver := NewACPowerFlowSolver()
		buses := []Bus{
			{ID: 1, Name: "Darlington Gen", Type: BusSlack, PGenMW: 300.0, QGenMVAR: 50.0, VMagPU: 1.0, VAngRad: 0.0},
			{ID: 2, Name: "Durham Substation", Type: BusPQ, PLoadMW: 120.0, QLoadMVAR: 30.0, VMagPU: 1.0, VAngRad: 0.0},
			{ID: 3, Name: "AI Compute Park", Type: BusPQ, PLoadMW: 150.0, QLoadMVAR: 40.0, VMagPU: 1.0, VAngRad: 0.0},
		}
		lines := []TransmissionLine{
			{ID: "L1-2", FromBus: 1, ToBus: 2, Resistance: 0.01, Reactance: 0.05, MVARating: 400.0},
			{ID: "L2-3", FromBus: 2, ToBus: 3, Resistance: 0.01, Reactance: 0.04, MVARating: 250.0},
		}

		res := solver.Solve(buses, lines, 50, 1e-2)
		if !res.Converged {
			t.Fatalf("expected AC power flow to converge, mismatch: %f", res.MaxMismatch)
		}
		if res.TotalGenMW != 300.0 {
			t.Errorf("expected 300 MW generation, got %f", res.TotalGenMW)
		}
		if res.TotalLoadMW != 270.0 {
			t.Errorf("expected 270 MW load, got %f", res.TotalLoadMW)
		}
		if len(res.Buses) != 3 {
			t.Errorf("expected 3 buses in result, got %d", len(res.Buses))
		}
	})

	t.Run("N-1 Grid Contingency Analysis", func(t *testing.T) {
		ca := NewContingencyAnalyzer()
		buses := []Bus{
			{ID: 1, Name: "Gen Bus", Type: BusSlack, PGenMW: 200.0},
			{ID: 2, Name: "Load Bus", Type: BusPQ, PLoadMW: 180.0},
		}
		lines := []TransmissionLine{
			{ID: "Circuit-1", FromBus: 1, ToBus: 2, Resistance: 0.02, Reactance: 0.08, MVARating: 250.0},
			{ID: "Circuit-2", FromBus: 1, ToBus: 2, Resistance: 0.02, Reactance: 0.08, MVARating: 250.0},
		}

		contingencies := ca.AnalyzeN1(buses, lines)
		if len(contingencies) != 2 {
			t.Fatalf("expected 2 N-1 contingency evaluations, got %d", len(contingencies))
		}
		for _, c := range contingencies {
			if !c.IsSystemSecure {
				t.Errorf("parallel circuit should maintain system security upon single line trip")
			}
		}
	})
}
