package risk

import (
	"math"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// TailRiskProfile quantifies extreme fat-tailed megaproject distributions beyond Gaussian assumptions.
type TailRiskProfile struct {
	ProjectID             string    `json:"project_id"`
	TailIndexAlpha        float64   `json:"tail_index_alpha"`        // Pareto power-law parameter (1.2 to 1.8 for megaprojects)
	ProbExceeding100Pct   float64   `json:"prob_exceeding_100pct"`   // Probability of doubling base Capex (0.0 to 1.0)
	ProbExceedingDoubling float64   `json:"prob_exceeding_doubling"` // Probability of 2x schedule extension
	VaR99CAD              int64     `json:"var_99_cad"`              // 99th percentile Value at Risk (Capex)
	ExpectedShortfall99CAD int64    `json:"expected_shortfall_99_cad"` // CVaR-99: Average loss given overrun is in worst 1%
	FragilityRating       string    `json:"fragility_rating"`        // "ANTIFRAGILE", "ROBUST", "FRAGILE", "EXTREME_HAZARD"
	DeRiskingMeasures     []string  `json:"de_risking_measures"`
	EvaluatedAt           time.Time `json:"evaluated_at"`
}

// EvaluateTailRisk computes power-law extreme hazard estimates calibrated from historical megaprojects.
func (e *Evaluator) EvaluateTailRisk(forecast *FlyvbjergRiskForecast) *TailRiskProfile {
	baseCapex := float64(forecast.BaseCapexCAD)
	if baseCapex <= 0 {
		baseCapex = 1_000_000 // Fallback minimum $1M CAD
	}

	// Tail index alpha: smaller alpha = fatter tails
	// Base infrastructure alpha ~ 1.6; nuclear and megamining drop to ~1.3
	alpha := 1.65
	switch forecast.Sector {
	case domain.SectorNuclearEnergy:
		alpha = 1.25
	case domain.SectorCriticalMinerals, domain.SectorMiningMetals:
		alpha = 1.35
	case domain.SectorTransportation:
		alpha = 1.45
	case domain.SectorCleanEnergy:
		alpha = 1.70
	}

	if forecast.RemoteGeographyPenalty > 0 {
		alpha -= 0.10
	}
	if forecast.TechNoveltyPenalty > 0 {
		alpha -= 0.10
	}
	if alpha < 1.05 {
		alpha = 1.05 // Bound to avoid division by zero in mean calculation
	}

	// Empirical probability of doubling cost based on mu and sigma
	// Derived from lognormal tail prob: P(X >= 2.0) = 1 - Phi((ln(2) - mu)/sigma)
	// Plus power-law adjustment
	ref := DefaultReferenceClasses(forecast.Sector)
	mu := ref.CostMeanMu + forecast.RemoteGeographyPenalty
	sigma := ref.CostStdSigma + forecast.TechNoveltyPenalty

	zDoubling := (math.Log(2.0) - mu) / sigma
	probLognormal := 0.5 * math.Erfc(zDoubling/math.Sqrt2)

	// Megaproject fat-tail uplift
	probDouble := math.Min(probLognormal*1.35, 0.65)

	// VaR-99: 99th percentile factor under generalized Pareto tail
	// x_99 = x_min * (1 - 0.99)^(-1/alpha)
	scaleFactor := math.Pow(0.01, -1.0/alpha)
	// Cap scale factor to realistic maximum (e.g. 5x base capex)
	if scaleFactor > 4.5 {
		scaleFactor = 4.5
	}
	var99 := int64(math.Round(baseCapex * scaleFactor))

	// Expected Shortfall (CVaR-99): E[X | X >= x_99] = (alpha / (alpha - 1)) * x_99
	cvarFactor := alpha / (alpha - 1.0)
	if cvarFactor > 3.0 {
		cvarFactor = 3.0
	}
	cvar99 := int64(math.Round(float64(var99) * cvarFactor * 0.5)) // Blended conditional expectation
	if cvar99 < var99 {
		cvar99 = int64(float64(var99) * 1.15)
	}

	// Fragility categorization
	fragility := "ROBUST"
	var measures []string

	if probDouble >= 0.30 || alpha < 1.30 {
		fragility = "EXTREME_HAZARD"
		measures = append(measures, "Bifurcate monolith into modularized serial units to prevent catastrophic systemic lock-in.")
		measures = append(measures, "Establish guaranteed maximum price (GMP) contracts with Tier-1 EPC contractors.")
		measures = append(measures, "Maintain contingency capital reserve of at least 45% of base Capex.")
	} else if probDouble >= 0.15 {
		fragility = "FRAGILE"
		measures = append(measures, "Front-load engineering and environmental review prior to Final Investment Decision (FID).")
		measures = append(measures, "Deploy digital twin and continuous satellite earth-observation milestone tracking.")
	} else {
		fragility = "ROBUST"
		measures = append(measures, "Maintain standard project risk register and quarterly reference class audits.")
	}

	return &TailRiskProfile{
		ProjectID:              forecast.ProjectID,
		TailIndexAlpha:         math.Round(alpha*100) / 100,
		ProbExceeding100Pct:    math.Round(probDouble*1000) / 1000,
		ProbExceedingDoubling:  math.Round(probDouble*0.8*1000) / 1000,
		VaR99CAD:               var99,
		ExpectedShortfall99CAD: cvar99,
		FragilityRating:        fragility,
		DeRiskingMeasures:      measures,
		EvaluatedAt:            time.Now().UTC(),
	}
}
