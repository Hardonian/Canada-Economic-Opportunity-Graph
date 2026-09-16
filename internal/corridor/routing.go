package corridor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

// InfrastructureType classifies linear right-of-way corridors.
type InfrastructureType string

const (
	TypeHVDCTransmission InfrastructureType = "HVDC_TRANSMISSION_LINE"
	TypeHydrogenPipeline InfrastructureType = "CLEAN_HYDROGEN_PIPELINE"
	TypeCO2Pipeline      InfrastructureType = "CARBON_CAPTURE_SEQUESTRATION_PIPELINE"
	TypeArcticHeavyRail  InfrastructureType = "ARCTIC_HEAVY_HAUL_RAIL"
	TypeAllWeatherRoad   InfrastructureType = "NORTHERN_ALL_WEATHER_ROAD"
)

// CorridorPoint represents a waypoint along a planned linear route.
type CorridorPoint struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	ElevationMeters float64 `json:"elevation_meters"`
}

// EnvironmentalSensitivity quantifies biological and geotechnical constraints.
type EnvironmentalSensitivity struct {
	PermafrostThawHazardScore float64 `json:"permafrost_thaw_hazard_score"` // 0.0 - 1.0
	CaribouRangeIntersectKM   float64 `json:"caribou_range_intersect_km"`
	WetlandCrossingCount      int     `json:"wetland_crossing_count"`
	WatercourseCrossingsMajor int     `json:"watercourse_crossings_major"`
}

// LandTenureDistribution breaks down ownership across the corridor.
type LandTenureDistribution struct {
	CrownProvincialPercent float64 `json:"crown_provincial_percent"`
	CrownFederalPercent    float64 `json:"crown_federal_percent"`
	IndigenousSettlementPercent float64 `json:"indigenous_settlement_percent"`
	PrivateFreeholdPercent float64 `json:"private_freehold_percent"`
}

// RouteEvaluation evaluates cost, geotechnical risk, and timing for a linear corridor.
type RouteEvaluation struct {
	CorridorID          string                   `json:"corridor_id"`
	Type                InfrastructureType       `json:"type"`
	OriginName          string                   `json:"origin_name"`
	DestinationName     string                   `json:"destination_name"`
	TotalLengthKM       float64                  `json:"total_length_km"`
	BaseConstructionCapexCAD int64               `json:"base_construction_capex_cad"`
	RightOfWayAcquisitionCAD int64               `json:"right_of_way_acquisition_cad"`
	EnvironmentalMitigationCAD int64             `json:"environmental_mitigation_cad"`
	TotalEstimatedCapexCAD int64                 `json:"total_estimated_capex_cad"`
	EstimatedScheduleMonths int                  `json:"estimated_schedule_months"`
	ImpedanceIndex      float64                  `json:"impedance_index"` // Lower is easier
	Sensitivity         EnvironmentalSensitivity `json:"sensitivity"`
	LandTenure          LandTenureDistribution   `json:"land_tenure"`
	KeyWaypoints        []CorridorPoint          `json:"key_waypoints"`
	AuditHash           string                   `json:"audit_hash"`
}

// RoutingEngine evaluates linear rights-of-way.
type RoutingEngine struct{}

// NewRoutingEngine instantiates the corridor pathfinder.
func NewRoutingEngine() *RoutingEngine {
	return &RoutingEngine{}
}

// EvaluateCorridor calculates geotechnical, land-tenure, and CAPEX metrics for a linear alignment.
func (r *RoutingEngine) EvaluateCorridor(origin, dest CorridorPoint, infraType InfrastructureType) *RouteEvaluation {
	// Great-circle distance approximation (Haversine)
	const earthRadiusKM = 6371.0
	dLat := (dest.Latitude - origin.Latitude) * (math.Pi / 180.0)
	dLon := (dest.Longitude - origin.Longitude) * (math.Pi / 180.0)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(origin.Latitude*(math.Pi/180.0))*math.Cos(dest.Latitude*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	directDistanceKM := earthRadiusKM * c

	// Linear winding factor (terrain avoidance ~1.22x)
	actualLengthKM := directDistanceKM * 1.22

	var perKMBaseCAD int64
	var scheduleMonths int

	switch infraType {
	case TypeHVDCTransmission:
		perKMBaseCAD = 3_200_000 // $3.2M / km for 500kV bipole HVDC
		scheduleMonths = int(actualLengthKM*0.12) + 24
	case TypeHydrogenPipeline:
		perKMBaseCAD = 4_500_000 // $4.5M / km for high-pressure embrittlement-resistant steel
		scheduleMonths = int(actualLengthKM*0.15) + 30
	case TypeCO2Pipeline:
		perKMBaseCAD = 3_800_000
		scheduleMonths = int(actualLengthKM*0.14) + 24
	case TypeArcticHeavyRail:
		perKMBaseCAD = 6_500_000 // $6.5M / km for sub-Arctic ballasted rail
		scheduleMonths = int(actualLengthKM*0.20) + 36
	default:
		perKMBaseCAD = 2_800_000 // $2.8M / km for all-weather 2-lane industrial road
		scheduleMonths = int(actualLengthKM*0.10) + 18
	}

	baseCapex := int64(actualLengthKM * float64(perKMBaseCAD))
	rowAcquisition := int64(float64(baseCapex) * 0.08)
	envMitigation := int64(float64(baseCapex) * 0.12)
	totalCapex := baseCapex + rowAcquisition + envMitigation

	// North of 53° Latitude triggers permafrost thaw risks
	permafrostRisk := 0.15
	if (origin.Latitude+dest.Latitude)/2 > 53.0 {
		permafrostRisk = 0.72
	}

	eval := &RouteEvaluation{
		CorridorID:               fmt.Sprintf("corridor-%s-%s", strings.ToLower(origin.Name), strings.ToLower(dest.Name)),
		Type:                     infraType,
		OriginName:               origin.Name,
		DestinationName:          dest.Name,
		TotalLengthKM:            math.Round(actualLengthKM*10) / 10,
		BaseConstructionCapexCAD: baseCapex,
		RightOfWayAcquisitionCAD: rowAcquisition,
		EnvironmentalMitigationCAD: envMitigation,
		TotalEstimatedCapexCAD:   totalCapex,
		EstimatedScheduleMonths:  scheduleMonths,
		ImpedanceIndex:           math.Round((actualLengthKM*1.15+permafrostRisk*100)*10) / 10,
		Sensitivity: EnvironmentalSensitivity{
			PermafrostThawHazardScore: permafrostRisk,
			CaribouRangeIntersectKM:   math.Round(actualLengthKM * 0.35),
			WetlandCrossingCount:      int(actualLengthKM / 25.0),
			WatercourseCrossingsMajor: int(actualLengthKM / 45.0),
		},
		LandTenure: LandTenureDistribution{
			CrownProvincialPercent:      58.0,
			CrownFederalPercent:         12.0,
			IndigenousSettlementPercent: 24.0,
			PrivateFreeholdPercent:      6.0,
		},
		KeyWaypoints: []CorridorPoint{
			origin,
			{
				Name:            fmt.Sprintf("%s-%s Midpoint Hub", origin.Name, dest.Name),
				Latitude:        (origin.Latitude + dest.Latitude) / 2,
				Longitude:       (origin.Longitude + dest.Longitude) / 2,
				ElevationMeters: (origin.ElevationMeters + dest.ElevationMeters) / 2,
			},
			dest,
		},
	}

	auditData := fmt.Sprintf("%s|%s|%.1f|%d|%d|%.2f",
		eval.CorridorID, eval.Type, eval.TotalLengthKM, eval.TotalEstimatedCapexCAD, eval.EstimatedScheduleMonths, eval.ImpedanceIndex)
	h := sha256.Sum256([]byte(auditData))
	eval.AuditHash = hex.EncodeToString(h[:])

	return eval
}
