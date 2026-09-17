package gridphysics

import (
	"fmt"
	"math"
)

// BusType defines the nodal power flow classification.
type BusType string

const (
	BusSlack BusType = "SLACK"
	BusPV    BusType = "PV" // Voltage-controlled generator
	BusPQ    BusType = "PQ" // Load bus
)

// Bus represents an electrical substation node in the AC network.
type Bus struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Type      BusType `json:"type"`
	PGenMW    float64 `json:"p_gen_mw"`
	QGenMVAR  float64 `json:"q_gen_mvar"`
	PLoadMW   float64 `json:"p_load_mw"`
	QLoadMVAR float64 `json:"q_load_mvar"`
	VMagPU    float64 `json:"v_mag_pu"` // Voltage magnitude in per-unit
	VAngRad   float64 `json:"v_ang_rad"` // Voltage phase angle in radians
}

// TransmissionLine represents an AC branch between two substations.
type TransmissionLine struct {
	ID          string  `json:"id"`
	FromBus     int     `json:"from_bus"`
	ToBus       int     `json:"to_bus"`
	Resistance  float64 `json:"resistance_pu"`  // R (p.u.)
	Reactance   float64 `json:"reactance_pu"`   // X (p.u.)
	MVARating   float64 `json:"mva_rating"`     // Thermal ampacity rating
	CurrentFlow float64 `json:"current_flow_mw"`
}

// ACSolverResult contains the converged Newton-Raphson power flow state.
type ACSolverResult struct {
	Converged        bool                `json:"converged"`
	Iterations       int                 `json:"iterations"`
	MaxMismatch      float64             `json:"max_mismatch"`
	TotalGenMW       float64             `json:"total_gen_mw"`
	TotalLoadMW      float64             `json:"total_load_mw"`
	TotalLossesMW    float64             `json:"total_losses_mw"`
	Buses            []Bus               `json:"buses"`
	Lines            []TransmissionLine  `json:"lines"`
	OverloadedLines  []string            `json:"overloaded_lines"`
	VoltageViolations []string           `json:"voltage_violations"`
}

// ACPowerFlowSolver solves the non-linear AC Newton-Raphson power flow equations.
type ACPowerFlowSolver struct{}

// NewACPowerFlowSolver creates an AC power flow solver.
func NewACPowerFlowSolver() *ACPowerFlowSolver {
	return &ACPowerFlowSolver{}
}

// Solve executes the Newton-Raphson algorithm over the provided network topology.
func (s *ACPowerFlowSolver) Solve(buses []Bus, lines []TransmissionLine, maxIter int, tol float64) *ACSolverResult {
	n := len(buses)
	if n == 0 {
		return &ACSolverResult{Converged: false}
	}

	// Make local working copies of buses
	workingBuses := make([]Bus, n)
	copy(workingBuses, buses)
	for i := range workingBuses {
		if workingBuses[i].VMagPU <= 0 {
			workingBuses[i].VMagPU = 1.0 // Flat start
		}
	}

	// Admittance matrix Y_bus = G + jB
	G := make([][]float64, n)
	B := make([][]float64, n)
	for i := 0; i < n; i++ {
		G[i] = make([]float64, n)
		B[i] = make([]float64, n)
	}

	for _, line := range lines {
		f := line.FromBus - 1
		t := line.ToBus - 1
		if f < 0 || f >= n || t < 0 || t >= n {
			continue
		}
		denom := line.Resistance*line.Resistance + line.Reactance*line.Reactance
		if denom <= 0 {
			denom = 1e-5
		}
		g := line.Resistance / denom
		b := -line.Reactance / denom

		G[f][t] -= g
		B[f][t] -= b
		G[t][f] -= g
		B[t][f] -= b

		G[f][f] += g
		B[f][f] += b
		G[t][t] += g
		B[t][t] += b
	}

	// Newton-Raphson iterative loop
	converged := false
	iter := 0
	maxMismatch := 0.0

	for iter < maxIter {
		iter++
		maxMismatch = 0.0

		// Calculate scheduled mismatches Delta P and Delta Q
		for i := 0; i < n; i++ {
			if workingBuses[i].Type == BusSlack {
				continue
			}

			// P_calc and Q_calc
			pCalc := 0.0
			qCalc := 0.0
			vi := workingBuses[i].VMagPU
			ai := workingBuses[i].VAngRad

			for j := 0; j < n; j++ {
				vj := workingBuses[j].VMagPU
				aj := workingBuses[j].VAngRad
				angDiff := ai - aj
				pCalc += vi * vj * (G[i][j]*math.Cos(angDiff) + B[i][j]*math.Sin(angDiff))
				qCalc += vi * vj * (G[i][j]*math.Sin(angDiff) - B[i][j]*math.Cos(angDiff))
			}

			pNet := (workingBuses[i].PGenMW - workingBuses[i].PLoadMW) / 100.0 // 100 MVA base
			delP := pNet - pCalc
			if math.Abs(delP) > maxMismatch {
				maxMismatch = math.Abs(delP)
			}

			// Angle update step (decoupled Newton-Raphson)
			bii := B[i][i]
			if math.Abs(bii) > 1e-4 {
				workingBuses[i].VAngRad -= delP / bii
			}

			if workingBuses[i].Type == BusPQ {
				qNet := (workingBuses[i].QGenMVAR - workingBuses[i].QLoadMVAR) / 100.0
				delQ := qNet - qCalc
				if math.Abs(delQ) > maxMismatch {
					maxMismatch = math.Abs(delQ)
				}
				// Voltage magnitude update
				if math.Abs(bii) > 1e-4 {
					workingBuses[i].VMagPU -= (delQ / bii) * 0.5
					workingBuses[i].VMagPU = math.Max(0.85, math.Min(1.15, workingBuses[i].VMagPU))
				}
			}
		}

		if maxMismatch < tol {
			converged = true
			break
		}
	}

	// Compute branch power flows and line overloads
	workingLines := make([]TransmissionLine, len(lines))
	copy(workingLines, lines)
	var overloads []string
	var voltageViols []string

	totalGen := 0.0
	totalLoad := 0.0

	for i := range workingBuses {
		totalGen += workingBuses[i].PGenMW
		totalLoad += workingBuses[i].PLoadMW

		if workingBuses[i].VMagPU < 0.95 || workingBuses[i].VMagPU > 1.05 {
			voltageViols = append(voltageViols, fmt.Sprintf("Bus %d (%s) voltage %.3f p.u. violates [0.95, 1.05]",
				workingBuses[i].ID, workingBuses[i].Name, workingBuses[i].VMagPU))
		}
	}

	for i := range workingLines {
		f := workingLines[i].FromBus - 1
		t := workingLines[i].ToBus - 1
		if f >= 0 && f < n && t >= 0 && t < n {
			angDiff := workingBuses[f].VAngRad - workingBuses[t].VAngRad
			x := workingLines[i].Reactance
			if math.Abs(x) < 1e-5 {
				x = 0.01
			}
			flowMW := (workingBuses[f].VMagPU * workingBuses[t].VMagPU * math.Sin(angDiff) / x) * 100.0
			workingLines[i].CurrentFlow = math.Round(math.Abs(flowMW)*10) / 10

			if workingLines[i].MVARating > 0 && workingLines[i].CurrentFlow > workingLines[i].MVARating {
				overloads = append(overloads, fmt.Sprintf("Line %s (%d->%d) flow %.1f MW exceeds rating %.1f MW (%.1f%%)",
					workingLines[i].ID, workingLines[i].FromBus, workingLines[i].ToBus,
					workingLines[i].CurrentFlow, workingLines[i].MVARating,
					(workingLines[i].CurrentFlow/workingLines[i].MVARating)*100.0))
			}
		}
	}

	losses := math.Max(0.0, totalGen-totalLoad)

	return &ACSolverResult{
		Converged:         converged,
		Iterations:        iter,
		MaxMismatch:       math.Round(maxMismatch*1e5) / 1e5,
		TotalGenMW:        totalGen,
		TotalLoadMW:       totalLoad,
		TotalLossesMW:     math.Round(losses*10) / 10,
		Buses:             workingBuses,
		Lines:             workingLines,
		OverloadedLines:   overloads,
		VoltageViolations: voltageViols,
	}
}
