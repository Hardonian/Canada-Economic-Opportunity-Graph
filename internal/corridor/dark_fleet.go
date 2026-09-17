package corridor

import (
	"fmt"
	"time"
)

// DarkFleetAnomaly defines suspicious maritime activity types.
type DarkFleetAnomaly string

const (
	AnomalyTransponderBlackout DarkFleetAnomaly = "AIS_TRANSPONDER_BLACKOUT"
	AnomalyFlagOfConvenience   DarkFleetAnomaly = "FLAG_OF_CONVENIENCE_HOPPING"
	AnomalyShipToShipTransfer  DarkFleetAnomaly = "SUSPICIOUS_SHIP_TO_SHIP_TRANSFER"
)

// DarkFleetAlert models a counter-intelligence detection on maritime trade flows.
type DarkFleetAlert struct {
	AlertID     string           `json:"alert_id"`
	MMSI        string           `json:"mmsi"`
	VesselName  string           `json:"vessel_name"`
	Flag        string           `json:"flag"`
	Anomaly     DarkFleetAnomaly `json:"anomaly"`
	RiskScore   float64          `json:"risk_score"` // 0.0 - 100.0
	Description string           `json:"description"`
	DetectedAt  time.Time        `json:"detected_at"`
}

// DarkFleetDetector identifies sanctions evasion and illicit maritime trade in Canadian waters.
type DarkFleetDetector struct{}

// NewDarkFleetDetector creates a dark fleet detector.
func NewDarkFleetDetector() *DarkFleetDetector {
	return &DarkFleetDetector{}
}

// EvaluateVessel assesses whether a vessel trajectory shows signs of dark fleet operations.
func (dfd *DarkFleetDetector) EvaluateVessel(
	mmsi, name, flag string,
	blackoutDurationHours float64,
	proximityDistanceMeters float64,
	flagChangesLast12Months int,
) []DarkFleetAlert {
	var alerts []DarkFleetAlert
	now := time.Now().UTC()

	// 1. Blackout check (> 4 hours intentional transponder silencing)
	if blackoutDurationHours > 4.0 {
		alerts = append(alerts, DarkFleetAlert{
			AlertID:     fmt.Sprintf("ALERT-BLACKOUT-%s", mmsi),
			MMSI:        mmsi,
			VesselName:  name,
			Flag:        flag,
			Anomaly:     AnomalyTransponderBlackout,
			RiskScore:   92.0,
			Description: fmt.Sprintf("Vessel silenced AIS transponder for %.1f hours in contiguous waters", blackoutDurationHours),
			DetectedAt:  now,
		})
	}

	// 2. Ship-to-Ship (STS) transfer in open water
	if proximityDistanceMeters < 50.0 && blackoutDurationHours > 1.0 {
		alerts = append(alerts, DarkFleetAlert{
			AlertID:     fmt.Sprintf("ALERT-STS-%s", mmsi),
			MMSI:        mmsi,
			VesselName:  name,
			Flag:        flag,
			Anomaly:     AnomalyShipToShipTransfer,
			RiskScore:   98.0,
			Description: "Close-proximity rendezvous (<50m) detected consistent with unmonitored oil or mineral transshipment",
			DetectedAt:  now,
		})
	}

	// 3. Flag hopping
	if flagChangesLast12Months >= 2 {
		alerts = append(alerts, DarkFleetAlert{
			AlertID:     fmt.Sprintf("ALERT-FLAG-%s", mmsi),
			MMSI:        mmsi,
			VesselName:  name,
			Flag:        flag,
			Anomaly:     AnomalyFlagOfConvenience,
			RiskScore:   75.0,
			Description: fmt.Sprintf("Vessel changed registry flags %d times within 12 months", flagChangesLast12Months),
			DetectedAt:  now,
		})
	}

	return alerts
}
