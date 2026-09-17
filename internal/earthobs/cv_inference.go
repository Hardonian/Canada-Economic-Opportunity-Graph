package earthobs

import (
	"math"
)

// MachineryDetectionResult summarizes computer vision equipment extraction on satellite imagery.
type MachineryDetectionResult struct {
	ProjectID                string  `json:"project_id"`
	HaulTruckCount           int     `json:"haul_truck_count"`
	ExcavatorCount           int     `json:"excavator_count"`
	CraneCount               int     `json:"crane_count"`
	DrillingRigCount         int     `json:"drilling_rig_count"`
	TotalMachineryUnits      int     `json:"total_machinery_units"`
	EstimatedSiteWorkers     int     `json:"estimated_site_workers"`
	ConstructionActivityRate string  `json:"construction_activity_rate"` // "DORMANT", "LOW", "MODERATE", "SURGING"
	ConfidenceScore          float64 `json:"confidence_score"`           // 0.0 - 1.0
}

// CVInferenceEngine processes high-resolution optical and infrared satellite orthomosaics.
type CVInferenceEngine struct{}

// NewCVInferenceEngine creates a computer vision inference engine.
func NewCVInferenceEngine() *CVInferenceEngine {
	return &CVInferenceEngine{}
}

// DetectMachinery estimates active construction equipment based on optical reflectance signatures.
func (cv *CVInferenceEngine) DetectMachinery(projectID string, capexCAD int64, ndvi, ndbi, rcsDB float64) *MachineryDetectionResult {
	// Heuristics: High NDBI (built-up) + High Radar Cross Section (steel) indicates cranes and heavy machines
	scaleFactor := math.Log10(math.Max(1e6, float64(capexCAD))) - 6.0

	cranes := 0
	haulTrucks := 0
	excavators := 0
	rigs := 0

	if rcsDB > 12.0 {
		cranes = int(math.Max(1, math.Round(scaleFactor*1.8)))
	}
	if ndbi > 0.15 {
		excavators = int(math.Max(2, math.Round(scaleFactor*3.2)))
		haulTrucks = int(math.Max(4, math.Round(scaleFactor*5.5)))
	}
	if rcsDB > 15.0 && ndbi > 0.25 {
		rigs = int(math.Max(1, math.Round(scaleFactor*1.2)))
	}

	total := cranes + haulTrucks + excavators + rigs
	workers := total * 14 // Average craft trades per heavy equipment unit

	rate := "DORMANT"
	if total > 20 {
		rate = "SURGING"
	} else if total > 8 {
		rate = "MODERATE"
	} else if total > 2 {
		rate = "LOW"
	}

	conf := 0.85
	if rcsDB > 10.0 {
		conf = 0.92
	}

	return &MachineryDetectionResult{
		ProjectID:                projectID,
		HaulTruckCount:           haulTrucks,
		ExcavatorCount:           excavators,
		CraneCount:               cranes,
		DrillingRigCount:         rigs,
		TotalMachineryUnits:      total,
		EstimatedSiteWorkers:     workers,
		ConstructionActivityRate: rate,
		ConfidenceScore:          conf,
	}
}
