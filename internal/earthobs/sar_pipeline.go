package earthobs

import (
	"fmt"
	"math"
	"time"
)

// InSARCoherenceAnalysis represents interferometric radar analysis across repeat satellite passes.
type InSARCoherenceAnalysis struct {
	PairID                string    `json:"pair_id"`
	MasterPassDate        time.Time `json:"master_pass_date"`
	SlavePassDate         time.Time `json:"slave_pass_date"`
	BaselineDistanceM     float64   `json:"baseline_distance_m"`
	TemporalBaselineDays  int       `json:"temporal_baseline_days"`
	InterferometricCoherence float64 `json:"interferometric_coherence"` // 0.0 to 1.0
	LineOfSightDisplacementMM float64 `json:"los_displacement_mm"`      // Surface displacement in mm
	TailingsStabilityRating   string  `json:"tailings_stability_rating"` // "STABLE", "MONITOR", "UNSTABLE"
	SoilMoistureEstimatePct   float64 `json:"soil_moisture_estimate_pct"`
}

// SARPipeline analyzes multi-temporal Sentinel-1 and RADARSAT Constellation radar data.
type SARPipeline struct{}

// NewSARPipeline creates a SAR processing pipeline.
func NewSARPipeline() *SARPipeline {
	return &SARPipeline{}
}

// ComputeInSAR evaluates phase displacement between two SAR radar acquisitions.
func (sp *SARPipeline) ComputeInSAR(masterDate, slaveDate time.Time, rawPhaseDiff float64, backscatterVV, backscatterVH float64) *InSARCoherenceAnalysis {
	days := int(slaveDate.Sub(masterDate).Hours() / 24)
	if days < 0 {
		days = -days
	}

	// Coherence degrades with temporal baseline
	coherence := math.Max(0.1, 0.95-float64(days)*0.015)
	// Radar wavelength for C-Band is ~55.5 mm. Phase diff 2pi = 28 mm displacement
	losDispMM := (rawPhaseDiff / (2 * math.Pi)) * 28.0

	stability := "STABLE"
	if math.Abs(losDispMM) > 25.0 {
		stability = "UNSTABLE"
	} else if math.Abs(losDispMM) > 10.0 {
		stability = "MONITOR"
	}

	// Depolarization ratio (VH / VV) correlates with surface roughness and volumetric soil moisture
	depolRatio := backscatterVH / math.Max(0.01, backscatterVV)
	moisture := math.Min(100.0, math.Max(5.0, depolRatio*120.0))

	return &InSARCoherenceAnalysis{
		PairID:                   fmt.Sprintf("INSAR-%d", slaveDate.Unix()),
		MasterPassDate:           masterDate,
		SlavePassDate:            slaveDate,
		BaselineDistanceM:        85.4,
		TemporalBaselineDays:     days,
		InterferometricCoherence: math.Round(coherence*100) / 100,
		LineOfSightDisplacementMM: math.Round(losDispMM*10) / 10,
		TailingsStabilityRating:  stability,
		SoilMoistureEstimatePct:  math.Round(moisture*10) / 10,
	}
}
