package corridor

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// StrategicGateway classifies key Canadian multi-modal hubs.
type StrategicGateway string

const (
	GatewayVancouver   StrategicGateway = "PORT_OF_VANCOUVER"
	GatewayPrinceRupert StrategicGateway = "PORT_OF_PRINCE_RUPERT"
	GatewayMontreal    StrategicGateway = "PORT_OF_MONTREAL"
	GatewayChurchill   StrategicGateway = "PORT_OF_CHURCHILL_ARCTIC"
)

// PortLogisticsProfile models dynamic throughput, rail turnaround, and bottleneck status.
type PortLogisticsProfile struct {
	GatewayID               StrategicGateway `json:"gateway_id"`
	PortName                string           `json:"port_name"`
	Province                string           `json:"province"`
	Class1RailConnections   []string         `json:"class1_rail_connections"` // "CN", "CPKC"
	AnnualThroughputMNTonnes float64         `json:"annual_throughput_mn_tonnes"`
	AverageVesselDwellHours float64          `json:"average_vessel_dwell_hours"`
	AverageRailcarDwellHours float64         `json:"average_railcar_dwell_hours"`
	BerthCapacityUtilizationPercent float64  `json:"berth_capacity_utilization_percent"`
	RailHeadTurnaroundCyclesDays float64     `json:"rail_head_turnaround_cycles_days"`
	IceClassEscortRequired  bool             `json:"ice_class_escort_required"`
	ActiveBottleneckStatus  string           `json:"active_bottleneck_status"` // "ELEVATED", "OPTIMAL", "CONGESTED"
	MitigationAction        string           `json:"mitigation_action"`
	AuditHash               string           `json:"audit_hash"`
}

// CanonicalGateways returns the four major national gateways.
func CanonicalGateways() []PortLogisticsProfile {
	gateways := []PortLogisticsProfile{
		{
			GatewayID:                       GatewayPrinceRupert,
			PortName:                        "Port of Prince Rupert (Fairview & Ridley Terminals)",
			Province:                        "BC",
			Class1RailConnections:           []string{"CN"},
			AnnualThroughputMNTonnes:        32.5,
			AverageVesselDwellHours:         38.0,
			AverageRailcarDwellHours:        22.5,
			BerthCapacityUtilizationPercent: 78.5,
			RailHeadTurnaroundCyclesDays:    5.2,
			IceClassEscortRequired:          false,
			ActiveBottleneckStatus:          "OPTIMAL",
			MitigationAction:                "Fastest trans-Pacific transit route; ongoing Ridley Island Energy Export Terminal expansion.",
		},
		{
			GatewayID:                       GatewayVancouver,
			PortName:                        "Port of Vancouver (Roberts Bank & Burrard Inlet)",
			Province:                        "BC",
			Class1RailConnections:           []string{"CN", "CPKC"},
			AnnualThroughputMNTonnes:        150.0,
			AverageVesselDwellHours:         74.0,
			AverageRailcarDwellHours:        52.0,
			BerthCapacityUtilizationPercent: 91.2,
			RailHeadTurnaroundCyclesDays:    8.4,
			IceClassEscortRequired:          false,
			ActiveBottleneckStatus:          "CONGESTED",
			MitigationAction:                "Cascade priority container blocks to Roberts Bank Terminal 2 (RBT2) rail interchange.",
		},
		{
			GatewayID:                       GatewayMontreal,
			PortName:                        "Port of Montreal (Contrecœur Expansion Hub)",
			Province:                        "QC",
			Class1RailConnections:           []string{"CN", "CPKC"},
			AnnualThroughputMNTonnes:        38.0,
			AverageVesselDwellHours:         44.0,
			AverageRailcarDwellHours:        28.0,
			BerthCapacityUtilizationPercent: 82.0,
			RailHeadTurnaroundCyclesDays:    4.8,
			IceClassEscortRequired:          true, // Winter St. Lawrence seaway ice navigation
			ActiveBottleneckStatus:          "OPTIMAL",
			MitigationAction:                "Contrecœur terminal expansion will add 1.15M TEU capacity to St. Lawrence corridor.",
		},
		{
			GatewayID:                       GatewayChurchill,
			PortName:                        "Port of Churchill (Hudson Bay Arctic Deepwater Gateway)",
			Province:                        "MB",
			Class1RailConnections:           []string{"Arctic Gateway Hudson Bay Railway"},
			AnnualThroughputMNTonnes:        1.2,
			AverageVesselDwellHours:         24.0,
			AverageRailcarDwellHours:        18.0,
			BerthCapacityUtilizationPercent: 35.0,
			RailHeadTurnaroundCyclesDays:    6.5,
			IceClassEscortRequired:          true, // Arctic polar ice navigation
			ActiveBottleneckStatus:          "OPTIMAL",
			MitigationAction:                "Hudson Bay Railway rehabilitation enables direct European critical mineral and wheat exports.",
		},
	}

	for i := range gateways {
		g := &gateways[i]
		auditData := fmt.Sprintf("%s|%s|%.1f|%.1f|%.1f", g.GatewayID, g.PortName, g.AnnualThroughputMNTonnes, g.BerthCapacityUtilizationPercent, g.RailHeadTurnaroundCyclesDays)
		h := sha256.Sum256([]byte(auditData))
		g.AuditHash = hex.EncodeToString(h[:])
	}

	return gateways
}
