package nationalplanning

import (
	"fmt"
	"math"
)

// ChokepointID defines critical infrastructure single points of failure.
type ChokepointID string

const (
	ChokepointSeawayLocks       ChokepointID = "ST_LAWRENCE_SEAWAY_LOCKS"
	ChokepointVancouverGateway   ChokepointID = "PORT_OF_VANCOUVER_BURRARD"
	ChokepointFraserCanyonRail  ChokepointID = "FRASER_CANYON_RAIL_CORRIDOR"
	ChokepointAmbassadorBridge  ChokepointID = "DETROIT_WINDSOR_AMBASSADOR"
	ChokepointCansoStrait       ChokepointID = "CANSO_STRAIT_CAUSEWAY"
)

// ChokepointProfile models the capacity and throughput of a national logistics chokepoint.
type ChokepointProfile struct {
	ID                  ChokepointID `json:"id"`
	Name                string       `json:"name"`
	DailyThroughputCADM float64      `json:"daily_throughput_cad_m"`
	PrimaryCommodities  []string     `json:"primary_commodities"`
	AlternativeModes    []string     `json:"alternative_modes"`
	DiversionFrictionX  float64      `json:"diversion_friction_x"` // Cost multiplier to divert
}

// ChokepointSimulationResult summarizes the national macro fallout of a chokepoint outage.
type ChokepointSimulationResult struct {
	ChokepointID          ChokepointID `json:"chokepoint_id"`
	Name                  string       `json:"name"`
	DurationDays          int          `json:"duration_days"`
	DirectStrandedLossM   float64      `json:"direct_stranded_loss_cad_m"`
	TotalEconomicShockM   float64      `json:"total_economic_shock_cad_m"`
	NationalSupplyFreeze  string       `json:"national_supply_freeze"` // "LOCALIZED", "SEVERE", "SYSTEMIC_PARALYSIS"
	AffectedSupplyChains  []string     `json:"affected_supply_chains"`
	EmergencyActionVector string       `json:"emergency_action_vector"`
}

// ChokepointWargamer simulates national economic fallout from chokepoint disruptions.
type ChokepointWargamer struct {
	profiles map[ChokepointID]ChokepointProfile
}

// NewChokepointWargamer initializes canonical Canadian chokepoints.
func NewChokepointWargamer() *ChokepointWargamer {
	cw := &ChokepointWargamer{
		profiles: make(map[ChokepointID]ChokepointProfile),
	}

	cw.profiles[ChokepointSeawayLocks] = ChokepointProfile{
		ID:                  ChokepointSeawayLocks,
		Name:                "St. Lawrence Seaway Welland & Montreal-Lake Ontario Locks",
		DailyThroughputCADM: 140.0,
		PrimaryCommodities:  []string{"Grain", "Iron Ore", "Steel", "Petroleum Products"},
		AlternativeModes:    []string{"CN/CPKC Rail", "Highway 401 Trucking"},
		DiversionFrictionX:  2.8,
	}

	cw.profiles[ChokepointVancouverGateway] = ChokepointProfile{
		ID:                  ChokepointVancouverGateway,
		Name:                "Port of Vancouver Roberts Bank & Burrard Inlet Gateway",
		DailyThroughputCADM: 850.0,
		PrimaryCommodities:  []string{"Potash", "Metallurgical Coal", "Grain", "Containerized Consumer Goods"},
		AlternativeModes:    []string{"Port of Prince Rupert", "U.S. Puget Sound Ports"},
		DiversionFrictionX:  2.2,
	}

	cw.profiles[ChokepointFraserCanyonRail] = ChokepointProfile{
		ID:                  ChokepointFraserCanyonRail,
		Name:                "Fraser Canyon Dual-Track Rail Chokepoint (CN/CPKC directional)",
		DailyThroughputCADM: 420.0,
		PrimaryCommodities:  []string{"Export Agri-Food", "Forestry", "Bulk Minerals"},
		AlternativeModes:    []string{"Highway 1 Trucking (severely capacity-constrained)"},
		DiversionFrictionX:  3.5,
	}

	cw.profiles[ChokepointAmbassadorBridge] = ChokepointProfile{
		ID:                  ChokepointAmbassadorBridge,
		Name:                "Detroit-Windsor Ambassador Bridge & Gordie Howe Corridor",
		DailyThroughputCADM: 450.0,
		PrimaryCommodities:  []string{"Automotive Just-In-Time Parts", "Machinery", "Agriculture"},
		AlternativeModes:    []string{"Blue Water Bridge (Sarnia)", "Peace Bridge (Fort Erie)"},
		DiversionFrictionX:  1.9,
	}

	return cw
}

// SimulateDisruption runs an emergency wargame scenario of chokepoint closure.
func (cw *ChokepointWargamer) SimulateDisruption(chokepointID ChokepointID, days int) (*ChokepointSimulationResult, error) {
	prof, ok := cw.profiles[chokepointID]
	if !ok {
		return nil, fmt.Errorf("chokepoint profile %s not found", chokepointID)
	}

	directLoss := prof.DailyThroughputCADM * float64(days)
	// Indirect multiplier grows non-linearly with outage duration
	multiplier := 1.0 + (float64(days)/7.0)*prof.DiversionFrictionX*0.35
	totalShock := directLoss * multiplier

	freeze := "LOCALIZED"
	if days >= 14 || totalShock > 10000.0 {
		freeze = "SYSTEMIC_PARALYSIS"
	} else if days >= 5 || totalShock > 2500.0 {
		freeze = "SEVERE"
	}

	action := fmt.Sprintf("Activate National Supply Chain Taskforce Emergency Protocol. Redirect freight across %s.",
		prof.AlternativeModes[0])

	return &ChokepointSimulationResult{
		ChokepointID:          chokepointID,
		Name:                  prof.Name,
		DurationDays:          days,
		DirectStrandedLossM:   math.Round(directLoss),
		TotalEconomicShockM:   math.Round(totalShock),
		NationalSupplyFreeze:  freeze,
		AffectedSupplyChains:  prof.PrimaryCommodities,
		EmergencyActionVector: action,
	}, nil
}
