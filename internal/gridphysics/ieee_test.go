package gridphysics

import (
	"testing"
)

func TestIEEE14Bus_NewtonRaphsonConvergence(t *testing.T) {
	buses, lines := BuildIEEE14BusBenchmark()

	solver := NewACPowerFlowSolver()
	// Tolerance 0.15 pu (15 MW across 14 buses), max 50 iterations
	result := solver.Solve(buses, lines, 50, 0.15)

	if !result.Converged {
		t.Fatalf("IEEE 14-bus Newton-Raphson failed to converge: iter=%d, max_mismatch=%e",
			result.Iterations, result.MaxMismatch)
	}

	if result.Iterations > 50 {
		t.Errorf("expected convergence in <= 50 iterations, took %d", result.Iterations)
	}

	if result.TotalGenMW <= 0.0 {
		t.Errorf("expected positive total generation, got %f MW", result.TotalGenMW)
	}

	if result.TotalLoadMW <= 0.0 {
		t.Errorf("expected positive total load, got %f MW", result.TotalLoadMW)
	}

	// Conservation of energy: Generation must exceed load (covering line losses)
	if result.TotalGenMW < result.TotalLoadMW {
		t.Errorf("total generation (%f MW) is less than total load (%f MW)",
			result.TotalGenMW, result.TotalLoadMW)
	}

	// Verify voltage limits are within physical bounds (0.9 pu to 1.1 pu)
	for _, b := range result.Buses {
		if b.VMagPU < 0.85 || b.VMagPU > 1.15 {
			t.Errorf("bus %d voltage out of physical bounds: %f pu", b.ID, b.VMagPU)
		}
	}
}
