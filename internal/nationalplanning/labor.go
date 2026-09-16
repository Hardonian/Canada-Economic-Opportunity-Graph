package nationalplanning

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// TradeCategory represents a Red Seal certified craft labor trade.
type TradeCategory string

const (
	TradeElectricians    TradeCategory = "Industrial Electricians"
	TradePowerLinemen    TradeCategory = "High-Voltage Powerline Technicians"
	TradeNuclearWelders  TradeCategory = "Nuclear-Certified Welders & Fitters"
	TradeBoilermakers    TradeCategory = "Boilermakers & Pressure Vessel Techs"
	TradeMillwrights     TradeCategory = "Industrial Millwrights"
	TradeHeavyEquipment  TradeCategory = "Heavy Equipment Operators"
	TradeCivilCarpenters TradeCategory = "Formwork & Heavy Civil Carpenters"
)

// TradeDemandDetail holds quarterly demand and pinch-point metrics for a specific craft.
type TradeDemandDetail struct {
	Trade             TradeCategory `json:"trade"`
	PeakDemandFTE     int           `json:"peak_demand_fte"`
	RegionalSupplyFTE int           `json:"regional_supply_fte"`
	UtilizationPct    float64       `json:"utilization_pct"` // > 90% triggers collision warning
	CollisionStatus   string        `json:"collision_status"` // NORMAL, ELEVATED, CRITICAL_SHORTAGE
	WageInflationRisk string        `json:"wage_inflation_risk"` // LOW, MODERATE, SEVERE
}

// RegionalLaborReport provides labor constraint analysis for an economic region/province.
type RegionalLaborReport struct {
	Province            string              `json:"province"`
	TotalActiveCapexCAD int64               `json:"total_active_capex_cad"`
	ConcurrentProjects  int                 `json:"concurrent_projects"`
	TotalLaborDemandFTE int                 `json:"total_labor_demand_fte"`
	Trades              []TradeDemandDetail `json:"trades"`
	CollisionDetected   bool                `json:"collision_detected"`
	StrategicAdvice     string              `json:"strategic_advice"`
	AuditHash           string              `json:"audit_hash"`
	AnalyzedAt          time.Time           `json:"analyzed_at"`
}

// LaborAggregator calculates regional craft labor demand across concurrent capital builds.
type LaborAggregator struct{}

func NewLaborAggregator() *LaborAggregator {
	return &LaborAggregator{}
}

// RegionalTradeSupply estimates available union hall and trade capacity by province.
var regionalTradeSupply = map[string]map[TradeCategory]int{
	"ON": {
		TradeElectricians:    12000,
		TradePowerLinemen:    3500,
		TradeNuclearWelders:  2800,
		TradeBoilermakers:    4200,
		TradeMillwrights:     6500,
		TradeHeavyEquipment:  9000,
		TradeCivilCarpenters: 11000,
	},
	"QC": {
		TradeElectricians:    10000,
		TradePowerLinemen:    4000,
		TradeNuclearWelders:  1200,
		TradeBoilermakers:    3000,
		TradeMillwrights:     5500,
		TradeHeavyEquipment:  8500,
		TradeCivilCarpenters: 9500,
	},
	"AB": {
		TradeElectricians:    9000,
		TradePowerLinemen:    2500,
		TradeNuclearWelders:  1800,
		TradeBoilermakers:    4500,
		TradeMillwrights:     6000,
		TradeHeavyEquipment:  12000,
		TradeCivilCarpenters: 8000,
	},
	"BC": {
		TradeElectricians:    7500,
		TradePowerLinemen:    2200,
		TradeNuclearWelders:  1100,
		TradeBoilermakers:    2500,
		TradeMillwrights:     4200,
		TradeHeavyEquipment:  7000,
		TradeCivilCarpenters: 7500,
	},
}

// AnalyzeProvince evaluates all projects in a province for trades bottlenecks.
func (la *LaborAggregator) AnalyzeProvince(province string, projects []*domain.Project) *RegionalLaborReport {
	provSupply, hasSupply := regionalTradeSupply[province]
	if !hasSupply {
		provSupply = map[TradeCategory]int{
			TradeElectricians:    3000,
			TradePowerLinemen:    800,
			TradeNuclearWelders:  400,
			TradeBoilermakers:    1000,
			TradeMillwrights:     1500,
			TradeHeavyEquipment:  2500,
			TradeCivilCarpenters: 3000,
		}
	}

	demand := map[TradeCategory]int{
		TradeElectricians:    0,
		TradePowerLinemen:    0,
		TradeNuclearWelders:  0,
		TradeBoilermakers:    0,
		TradeMillwrights:     0,
		TradeHeavyEquipment:  0,
		TradeCivilCarpenters: 0,
	}

	var provCapex int64
	projectCount := 0

	for _, p := range projects {
		if p.Province != province || p.CapexCAD <= 0 {
			continue
		}
		provCapex += p.CapexCAD
		projectCount++

		// Base FTE per $1B capex
		billions := float64(p.CapexCAD) / 1e9

		switch p.Sector {
		case domain.SectorNuclearEnergy:
			demand[TradeNuclearWelders] += int(math.Round(billions * 450))
			demand[TradeElectricians] += int(math.Round(billions * 600))
			demand[TradeBoilermakers] += int(math.Round(billions * 400))
			demand[TradeCivilCarpenters] += int(math.Round(billions * 500))
		case domain.SectorCriticalMinerals:
			demand[TradeHeavyEquipment] += int(math.Round(billions * 800))
			demand[TradeMillwrights] += int(math.Round(billions * 500))
			demand[TradePowerLinemen] += int(math.Round(billions * 300))
			demand[TradeElectricians] += int(math.Round(billions * 450))
		case domain.SectorCleanEnergy:
			demand[TradePowerLinemen] += int(math.Round(billions * 650))
			demand[TradeElectricians] += int(math.Round(billions * 500))
			demand[TradeHeavyEquipment] += int(math.Round(billions * 400))
		case domain.SectorAICompute:
			demand[TradeElectricians] += int(math.Round(billions * 900))
			demand[TradePowerLinemen] += int(math.Round(billions * 400))
			demand[TradeCivilCarpenters] += int(math.Round(billions * 350))
		default:
			demand[TradeHeavyEquipment] += int(math.Round(billions * 400))
			demand[TradeCivilCarpenters] += int(math.Round(billions * 400))
			demand[TradeElectricians] += int(math.Round(billions * 300))
		}
	}

	tradeDetails := make([]TradeDemandDetail, 0, len(demand))
	var totalDemand int
	var collisionFound bool

	for trade, dem := range demand {
		totalDemand += dem
		sup := provSupply[trade]
		if sup <= 0 {
			sup = 1000
		}
		util := (float64(dem) / float64(sup)) * 100.0

		status := "NORMAL"
		inflation := "LOW"

		switch {
		case util > 85.0:
			status = "CRITICAL_SHORTAGE"
			inflation = "SEVERE"
			collisionFound = true
		case util > 65.0:
			status = "ELEVATED"
			inflation = "MODERATE"
		}

		tradeDetails = append(tradeDetails, TradeDemandDetail{
			Trade:             trade,
			PeakDemandFTE:     dem,
			RegionalSupplyFTE: sup,
			UtilizationPct:    math.Round(util*10) / 10,
			CollisionStatus:   status,
			WageInflationRisk: inflation,
		})
	}

	sort.Slice(tradeDetails, func(i, j int) bool {
		return tradeDetails[i].UtilizationPct > tradeDetails[j].UtilizationPct
	})

	advice := "Regional labor supply sufficient across all standard construction trades."
	if collisionFound {
		advice = "CRITICAL COLLISION: Concurrent megaprojects exceed 85% of regional Red Seal capacity. Recommend inter-provincial trade mobility pacts, dedicated apprenticeship wage subsidies, and staggered civil construction windows."
	}

	rep := &RegionalLaborReport{
		Province:            province,
		TotalActiveCapexCAD: provCapex,
		ConcurrentProjects:  projectCount,
		TotalLaborDemandFTE: totalDemand,
		Trades:              tradeDetails,
		CollisionDetected:   collisionFound,
		StrategicAdvice:     advice,
		AnalyzedAt:          time.Now().UTC(),
	}

	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%d|%d", rep.Province, rep.TotalActiveCapexCAD, rep.TotalLaborDemandFTE, len(rep.Trades))))
	rep.AuditHash = hex.EncodeToString(h[:])

	return rep
}
