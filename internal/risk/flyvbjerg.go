package risk

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// ReferenceClassParam holds empirically calibrated lognormal distribution parameters.
type ReferenceClassParam struct {
	Sector             domain.Sector
	CostMeanMu         float64 // Mean of ln(1 + cost_overrun_pct)
	CostStdSigma       float64 // Std dev of ln(1 + cost_overrun_pct)
	ScheduleMeanMonths float64 // Mean schedule slip in months
	ScheduleStdMonths  float64 // Std dev schedule slip in months
	SampleCount        int     // Historical megaproject sample size in benchmark set
}

// DefaultReferenceClasses returns Flyvbjerg-calibrated historical parameters for Canadian infrastructure.
func DefaultReferenceClasses(sector domain.Sector) ReferenceClassParam {
	switch sector {
	case domain.SectorNuclearEnergy:
		return ReferenceClassParam{
			Sector:             sector,
			CostMeanMu:         0.58, // ~80% median cost overrun historically
			CostStdSigma:       0.45,
			ScheduleMeanMonths: 38.0, // ~3.2 years average delay
			ScheduleStdMonths:  18.0,
			SampleCount:        145,
		}
	case domain.SectorCriticalMinerals, domain.SectorMiningMetals:
		return ReferenceClassParam{
			Sector:             sector,
			CostMeanMu:         0.42, // ~52% median cost overrun
			CostStdSigma:       0.38,
			ScheduleMeanMonths: 24.0, // ~2.0 years delay
			ScheduleStdMonths:  14.0,
			SampleCount:        210,
		}
	case domain.SectorTransportation:
		return ReferenceClassParam{
			Sector:             sector,
			CostMeanMu:         0.35, // ~42% median cost overrun
			CostStdSigma:       0.30,
			ScheduleMeanMonths: 18.0,
			ScheduleStdMonths:  10.0,
			SampleCount:        340,
		}
	case domain.SectorCleanEnergy:
		return ReferenceClassParam{
			Sector:             sector,
			CostMeanMu:         0.22, // ~25% median cost overrun
			CostStdSigma:       0.25,
			ScheduleMeanMonths: 12.0,
			ScheduleStdMonths:  8.0,
			SampleCount:        180,
		}
	default:
		return ReferenceClassParam{
			Sector:             sector,
			CostMeanMu:         0.30,
			CostStdSigma:       0.30,
			ScheduleMeanMonths: 16.0,
			ScheduleStdMonths:  10.0,
			SampleCount:        150,
		}
	}
}

// PercentilePoint holds a probability threshold prediction.
type PercentilePoint struct {
	Percentile         int   `json:"percentile"` // e.g. 10, 50, 80, 90
	CostOverrunPct     float64 `json:"cost_overrun_pct"`
	ForecastCapexCAD   int64   `json:"forecast_capex_cad"`
	ScheduleDelayMonths int    `json:"schedule_delay_months"`
}

// FlyvbjergRiskForecast provides full reference-class risk distributions.
type FlyvbjergRiskForecast struct {
	ProjectID              string            `json:"project_id"`
	ProjectName            string            `json:"project_name"`
	Sector                 domain.Sector     `json:"sector"`
	BaseCapexCAD           int64             `json:"base_capex_cad"`
	ReferenceClass         string            `json:"reference_class"`
	HistoricalSampleSize   int               `json:"historical_sample_size"`
	RemoteGeographyPenalty float64           `json:"remote_geography_penalty"` // Added mu
	TechNoveltyPenalty     float64           `json:"tech_novelty_penalty"`     // Added sigma
	ExpectedCostOverrunPct float64           `json:"expected_cost_overrun_pct"`
	ExpectedDelayMonths    int               `json:"expected_delay_months"`
	Percentiles            []PercentilePoint `json:"percentiles"` // P10, P50, P80, P90
	AuditHash              string            `json:"audit_hash"`
	CalculatedAt           time.Time         `json:"calculated_at"`
}

// Evaluator performs Bayesian reference-class calculations.
type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// ForecastProject generates Flyvbjerg risk hazard distributions.
func (e *Evaluator) ForecastProject(project *domain.Project) *FlyvbjergRiskForecast {
	ref := DefaultReferenceClasses(project.Sector)

	mu := ref.CostMeanMu
	sigma := ref.CostStdSigma
	delayMean := ref.ScheduleMeanMonths

	// Remote geography penalty: Arctic, Northern Ontario, Northern Quebec
	remotePenalty := 0.0
	loc := strings.ToLower(project.LocationName + " " + project.Province)
	if strings.Contains(loc, "arctic") || strings.Contains(loc, "nu") || strings.Contains(loc, "nt") ||
		strings.Contains(loc, "yt") || strings.Contains(loc, "james bay") || strings.Contains(loc, "ring of fire") {
		remotePenalty = 0.12
		mu += remotePenalty
		delayMean += 8.0
	}

	// Tech novelty penalty
	techPenalty := 0.0
	sub := strings.ToLower(project.Subsector + " " + project.Summary)
	if strings.Contains(sub, "smr") || strings.Contains(sub, "hydrogen") || strings.Contains(sub, "first-of-a-kind") ||
		strings.Contains(sub, "fook") || strings.Contains(sub, "novel") {
		techPenalty = 0.10
		sigma += techPenalty
		delayMean += 6.0
	}

	// Standard normal quantiles: Z for P10 (-1.282), P50 (0.0), P80 (0.842), P90 (1.282)
	zVals := []struct {
		p int
		z float64
	}{
		{p: 10, z: -1.282},
		{p: 50, z: 0.0},
		{p: 80, z: 0.842},
		{p: 90, z: 1.282},
	}

	percentiles := make([]PercentilePoint, 0, len(zVals))
	capex := float64(project.CapexCAD)

	for _, pt := range zVals {
		costFactor := math.Exp(mu + pt.z*sigma)
		overrunPct := (costFactor - 1.0) * 100.0
		if overrunPct < 0 {
			overrunPct = 0
		}
		fcCapex := capex * costFactor
		delayM := int(math.Round(delayMean + pt.z*ref.ScheduleStdMonths))
		if delayM < 0 {
			delayM = 0
		}

		percentiles = append(percentiles, PercentilePoint{
			Percentile:          pt.p,
			CostOverrunPct:      math.Round(overrunPct*10) / 10,
			ForecastCapexCAD:    int64(math.Round(fcCapex)),
			ScheduleDelayMonths: delayM,
		})
	}

	expectedOverrun := (math.Exp(mu+0.5*sigma*sigma) - 1.0) * 100.0

	forecast := &FlyvbjergRiskForecast{
		ProjectID:              project.ID,
		ProjectName:            project.Name,
		Sector:                 project.Sector,
		BaseCapexCAD:           project.CapexCAD,
		ReferenceClass:         string(project.Sector),
		HistoricalSampleSize:   ref.SampleCount,
		RemoteGeographyPenalty: remotePenalty,
		TechNoveltyPenalty:     techPenalty,
		ExpectedCostOverrunPct: math.Round(expectedOverrun*10) / 10,
		ExpectedDelayMonths:    int(math.Round(delayMean)),
		Percentiles:            percentiles,
		CalculatedAt:           time.Now().UTC(),
	}

	h := sha256.Sum256(fmt.Appendf(nil, "%s|%d|%.2f|%d", forecast.ProjectID, forecast.BaseCapexCAD, forecast.ExpectedCostOverrunPct, forecast.ExpectedDelayMonths))
	forecast.AuditHash = hex.EncodeToString(h[:])

	return forecast
}
