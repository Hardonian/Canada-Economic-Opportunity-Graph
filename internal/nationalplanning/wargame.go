package nationalplanning

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// ShockScenario identifies pre-calibrated geopolitical/economic disruption tests.
type ShockScenario string

const (
	ShockUSMCATariffs       ShockScenario = "USMCA_2026_TARIFF_25"
	ShockCriticalMineralBan ShockScenario = "FOREIGN_CRITICAL_MINERAL_EXPORT_BAN"
	ShockArcticChokeDisrupt ShockScenario = "ARCTIC_TRANSIT_CHOKEPOINT_SEIZURE"
	ShockTransformerCrisis  ShockScenario = "GLOBAL_500KV_TRANSFORMER_SHORTAGE"
	ShockInterestSpike250   ShockScenario = "STAGFLATION_RATE_SPIKE_250BPS"
)

// WarGameRequest specifies shock scenario parameters.
type WarGameRequest struct {
	Scenario        ShockScenario `json:"scenario"`
	TariffRatePct   float64       `json:"tariff_rate_pct,omitempty"`   // e.g. 25.0
	DurationMonths  int           `json:"duration_months,omitempty"`  // e.g. 24
	CustomNotes     string        `json:"custom_notes,omitempty"`
}

// StalledProjectDetail holds project-specific shock vulnerability.
type StalledProjectDetail struct {
	ProjectID         string        `json:"project_id"`
	ProjectName       string        `json:"project_name"`
	Sector            domain.Sector `json:"sector"`
	Province          string        `json:"province"`
	OriginalCapexCAD  int64         `json:"original_capex_cad"`
	StallLikelihood   float64       `json:"stall_likelihood"`   // 0.0 - 1.0
	VulnerabilityNote string        `json:"vulnerability_note"`
	RecommendedAction string        `json:"recommended_action"`
}

// WarGameSimulationResult holds the system-wide fallout and tactical countermeasures.
type WarGameSimulationResult struct {
	SimulationID            string                 `json:"simulation_id"`
	Scenario                ShockScenario          `json:"scenario"`
	ScenarioTitle           string                 `json:"scenario_title"`
	ScenarioDescription     string                 `json:"scenario_description"`
	TotalAssetsStalledCount int                    `json:"total_assets_stalled_count"`
	TotalFrozenCapexCAD     int64                  `json:"total_frozen_capex_cad"`
	EstimatedNationalGDPLossCAD int64              `json:"estimated_national_gdp_loss_cad"`
	AffectedSectors         []domain.Sector        `json:"affected_sectors"`
	StalledProjects         []StalledProjectDetail `json:"stalled_projects"`
	SovereignMitigations    []string               `json:"sovereign_mitigations"`
	AuditHash               string                 `json:"audit_hash"`
	SimulatedAt             time.Time              `json:"simulated_at"`
}

// WarGameEngine runs systemic geopolitical shock stress-testing over the economic graph.
type WarGameEngine struct{}

func NewWarGameEngine() *WarGameEngine {
	return &WarGameEngine{}
}

// SimulateScenario executes the chosen macro shock against the portfolio of projects.
func (wge *WarGameEngine) SimulateScenario(projects []*domain.Project, req WarGameRequest) *WarGameSimulationResult {
	var title, desc string
	stalled := make([]StalledProjectDetail, 0)
	var frozenCapex int64
	sectorMap := make(map[domain.Sector]bool)
	mitigations := make([]string, 0)

	switch req.Scenario {
	case ShockUSMCATariffs:
		title = "USMCA 2026 Comprehensive Tariff Shock (25%)"
		desc = "Unilateral 25% tariffs imposed on Canadian non-defense metals, automotive precursors, battery chemicals, and clean power exports."
		mitigations = []string{
			"Activate bilateral CETA/CPTPP emergency trade corridors to pivot export quotas toward Europe and Japan.",
			"Expand the Canada Growth Fund (CGF) contract-for-difference (CFD) backstop to absorb export margin compressions.",
			"Mandate domestic Canadian procurement preference for federally subsidized public infrastructure builds.",
		}
		for _, p := range projects {
			if p.Sector == domain.SectorCriticalMinerals || p.Sector == domain.SectorIndustrialMfg {
				stalled = append(stalled, StalledProjectDetail{
					ProjectID:         p.ID,
					ProjectName:       p.Name,
					Sector:            p.Sector,
					Province:          p.Province,
					OriginalCapexCAD:  p.CapexCAD,
					StallLikelihood:   0.85,
					VulnerabilityNote: "High export dependence on US automotive and OEM battery manufacturers; border tariff erodes project IRR.",
					RecommendedAction: "Provide emergency CIB export guarantee and negotiate temporary defense-critical exemption under Title III DPA.",
				})
				frozenCapex += p.CapexCAD
				sectorMap[p.Sector] = true
			}
		}

	case ShockCriticalMineralBan:
		title = "Foreign Critical Mineral & Chemical Precursor Embargo"
		desc = "Foreign state embargo on synthetic graphite, refined rare earth oxides, and battery cathode catalysts."
		mitigations = []string{
			"Fast-track domestic refining approvals for Ontario Crawford and Quebec Baie-James pCAM facilities under emergency Cabinet directive.",
			"Establish a National Strategic Critical Minerals Reserve financed via CIB concessionary facilities.",
			"Authorize DND IDEaS dual-use capital injection for domestic hydrometallurgical recycling.",
		}
		for _, p := range projects {
			if p.Sector == domain.SectorCriticalMinerals || strings.Contains(strings.ToLower(p.Summary), "battery") {
				stalled = append(stalled, StalledProjectDetail{
					ProjectID:         p.ID,
					ProjectName:       p.Name,
					Sector:            p.Sector,
					Province:          p.Province,
					OriginalCapexCAD:  p.CapexCAD,
					StallLikelihood:   0.75,
					VulnerabilityNote: "Processing mill reliant on offshore proprietary catalysts and specialized chemical reagents.",
					RecommendedAction: "Accelerate domestic pilot hydrometallurgy facilities to replace offshore inputs.",
				})
				frozenCapex += p.CapexCAD
				sectorMap[p.Sector] = true
			}
		}

	case ShockTransformerCrisis:
		title = "Global 500kV High-Voltage Transformer Shortage (Lead Time 260 Wks)"
		desc = "Catastrophic global supply freeze on high-voltage autotransformers and grain-oriented electrical steel (GOES)."
		mitigations = []string{
			"Nationalize or provide 100% capex matching for Canadian domestic power transformer manufacturing (e.g. PTI Transformers, Northern Transformer).",
			"Implement provincial grid-sharing pact: pool strategic spare transformer inventories between OPG, Hydro-Québec, and BC Hydro.",
		}
		for _, p := range projects {
			if p.Sector == domain.SectorAICompute || p.Sector == domain.SectorNuclearEnergy || p.Sector == domain.SectorCleanEnergy {
				stalled = append(stalled, StalledProjectDetail{
					ProjectID:         p.ID,
					ProjectName:       p.Name,
					Sector:            p.Sector,
					Province:          p.Province,
					OriginalCapexCAD:  p.CapexCAD,
					StallLikelihood:   0.90,
					VulnerabilityNote: "Substation energization delayed by 4+ years due to overseas transformer procurement queue.",
					RecommendedAction: "Reallocate federal CIB funding to domestic transformer manufacturing priority queue.",
				})
				frozenCapex += p.CapexCAD
				sectorMap[p.Sector] = true
			}
		}

	default: // Arctic / Northern Disruption
		title = "Arctic Transit Corridor & Deepwater Port Contestation"
		desc = "Severe maritime choke point disruption and sovereign defense alerts in Churchill and the Northwest Passage."
		mitigations = []string{
			"Deploy Canadian Armed Forces Operation NANOOK logistics assets to secure civil northern transport.",
			"Accelerate all-weather road links linking remote northern mining basins directly to transcontinental rail heads.",
		}
		for _, p := range projects {
			if p.Sector == domain.SectorDefenceArctic || strings.Contains(strings.ToLower(p.LocationName), "arctic") || p.Province == "NU" || p.Province == "NT" {
				stalled = append(stalled, StalledProjectDetail{
					ProjectID:         p.ID,
					ProjectName:       p.Name,
					Sector:            p.Sector,
					Province:          p.Province,
					OriginalCapexCAD:  p.CapexCAD,
					StallLikelihood:   0.70,
					VulnerabilityNote: "Remote logistics lifeline cut; air freight costs exceed economic viability.",
					RecommendedAction: "Deploy federal Arctic sovereign infrastructure defense package.",
				})
				frozenCapex += p.CapexCAD
				sectorMap[p.Sector] = true
			}
		}
	}

	sectors := make([]domain.Sector, 0, len(sectorMap))
	for s := range sectorMap {
		sectors = append(sectors, s)
	}

	gdpLoss := int64(math.Round(float64(frozenCapex) * 1.45)) // 1.45x macro multiplier loss

	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%d", req.Scenario, len(stalled), frozenCapex)))
	simID := fmt.Sprintf("war_%x", sum[:8])

	res := &WarGameSimulationResult{
		SimulationID:            simID,
		Scenario:                req.Scenario,
		ScenarioTitle:           title,
		ScenarioDescription:     desc,
		TotalAssetsStalledCount: len(stalled),
		TotalFrozenCapexCAD:     frozenCapex,
		EstimatedNationalGDPLossCAD: gdpLoss,
		AffectedSectors:         sectors,
		StalledProjects:         stalled,
		SovereignMitigations:    mitigations,
		SimulatedAt:             time.Now().UTC(),
	}

	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%d", res.SimulationID, res.TotalAssetsStalledCount, res.TotalFrozenCapexCAD)))
	res.AuditHash = hex.EncodeToString(h[:])

	return res
}
