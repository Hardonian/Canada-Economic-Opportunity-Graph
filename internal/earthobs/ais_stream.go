package earthobs

import (
	"math"
	"time"
)

// AISVessel represents a single vessel transponder record in Canadian territorial waters.
type AISVessel struct {
	MMSI        string    `json:"mmsi"`
	VesselName  string    `json:"vessel_name"`
	VesselType  string    `json:"vessel_type"` // "CONTAINER", "BULK_CARRIER", "TANKER", "TUG"
	Flag        string    `json:"flag"`
	SpeedKnots  float64   `json:"speed_knots"`
	Status      string    `json:"status"` // "UNDERWAY", "ANCHORED", "MOORED"
	ArrivalDate time.Time `json:"arrival_date"`
}

// PortCongestionTelemetry summarizes multi-modal logistics velocity at critical maritime gateways.
type PortCongestionTelemetry struct {
	PortID                   string    `json:"port_id"`
	PortName                 string    `json:"port_name"`
	AnchoredVesselCount      int       `json:"anchored_vessel_count"`
	MooredVesselCount        int       `json:"moored_vessel_count"`
	AverageDwellDays         float64   `json:"average_dwell_days"`
	BerthUtilizationPct      float64   `json:"berth_utilization_pct"`
	SupplyChainFrictionScore float64   `json:"supply_chain_friction_score"` // 0.0 - 100.0
	FrictionRating           string    `json:"friction_rating"`             // "FLUID", "ELEVATED", "CHOKED"
	EvaluatedAt              time.Time `json:"evaluated_at"`
}

// AISTrackingEngine evaluates maritime congestion and trade corridor velocity.
type AISTrackingEngine struct{}

// NewAISTrackingEngine creates an AIS telemetry engine.
func NewAISTrackingEngine() *AISTrackingEngine {
	return &AISTrackingEngine{}
}

// EvaluatePort computes logistics throughput metrics for a target Canadian deepwater port.
func (ate *AISTrackingEngine) EvaluatePort(portID, portName string, vessels []AISVessel) *PortCongestionTelemetry {
	anchored := 0
	moored := 0
	totalDwellHours := 0.0
	now := time.Now().UTC()

	for _, v := range vessels {
		dwellHours := now.Sub(v.ArrivalDate).Hours()
		if dwellHours > 0 {
			totalDwellHours += dwellHours
		}

		if v.Status == "MOORED" {
			moored++
		} else if v.Status == "ANCHORED" || v.SpeedKnots < 0.5 {
			anchored++
		}
	}

	avgDwellDays := 0.0
	if len(vessels) > 0 {
		avgDwellDays = math.Round((totalDwellHours/float64(len(vessels))/24.0)*10) / 10
	}

	berthUtil := math.Min(100.0, float64(moored)*12.5) // Assuming typical 8-12 berth deepwater terminal
	friction := math.Min(100.0, float64(anchored)*5.0+avgDwellDays*8.0)

	rating := "FLUID"
	if friction > 65.0 {
		rating = "CHOKED"
	} else if friction > 35.0 {
		rating = "ELEVATED"
	}

	return &PortCongestionTelemetry{
		PortID:                   portID,
		PortName:                 portName,
		AnchoredVesselCount:      anchored,
		MooredVesselCount:        moored,
		AverageDwellDays:         avgDwellDays,
		BerthUtilizationPct:      berthUtil,
		SupplyChainFrictionScore: math.Round(friction*10) / 10,
		FrictionRating:           rating,
		EvaluatedAt:              now,
	}
}
