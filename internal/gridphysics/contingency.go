package gridphysics

import (
	"math"
)

// ContingencyResult summarizes the grid security impact of an N-1 line trip.
type ContingencyResult struct {
	TrippedLineID     string   `json:"tripped_line_id"`
	FromBus           int      `json:"from_bus"`
	ToBus             int      `json:"to_bus"`
	PostTripOverloads []string `json:"post_trip_overloads"`
	MaxLoadingPct     float64  `json:"max_loading_pct"`
	IsSystemSecure    bool     `json:"is_system_secure"`
}

// ContingencyAnalyzer evaluates single and double contingency security (N-1 / N-2).
type ContingencyAnalyzer struct {
	solver *ACPowerFlowSolver
}

// NewContingencyAnalyzer creates an N-1 contingency analyzer.
func NewContingencyAnalyzer() *ContingencyAnalyzer {
	return &ContingencyAnalyzer{
		solver: NewACPowerFlowSolver(),
	}
}

// AnalyzeN1 evaluates the outage of each transmission line in sequence.
func (ca *ContingencyAnalyzer) AnalyzeN1(buses []Bus, lines []TransmissionLine) []ContingencyResult {
	var results []ContingencyResult

	for i := 0; i < len(lines); i++ {
		tripped := lines[i]

		// Form remaining network
		var remaining []TransmissionLine
		for j := 0; j < len(lines); j++ {
			if i != j {
				remaining = append(remaining, lines[j])
			}
		}

		res := ca.solver.Solve(buses, remaining, 20, 1e-3)
		maxLoad := 0.0

		for _, l := range res.Lines {
			if l.MVARating > 0 {
				loadPct := (l.CurrentFlow / l.MVARating) * 100.0
				if loadPct > maxLoad {
					maxLoad = loadPct
				}
			}
		}

		secure := res.Converged && len(res.OverloadedLines) == 0 && maxLoad <= 100.0

		results = append(results, ContingencyResult{
			TrippedLineID:     tripped.ID,
			FromBus:           tripped.FromBus,
			ToBus:             tripped.ToBus,
			PostTripOverloads: res.OverloadedLines,
			MaxLoadingPct:     math.Round(maxLoad*10) / 10,
			IsSystemSecure:    secure,
		})
	}

	return results
}
