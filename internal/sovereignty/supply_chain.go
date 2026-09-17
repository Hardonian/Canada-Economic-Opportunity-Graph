package sovereignty

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// CriticalMineral represents an official Canadian Critical Mineral as designated by NRCan.
type CriticalMineral struct {
	Symbol                string   `json:"symbol"`
	Name                  string   `json:"name"`
	CanadianReservesRank  int      `json:"canadian_reserves_rank"`  // Global rank in reserves (1-based, 0 if unranked)
	DomesticExtractionPct float64  `json:"domestic_extraction_pct"` // % of national demand extracted domestically
	DomesticRefiningPct   float64  `json:"domestic_refining_pct"`   // % of refined material processed domestically
	ForeignMonopolyShare  float64  `json:"foreign_monopoly_share"`  // % global refining held by top single foreign nation
	DominantForeignNation string   `json:"dominant_foreign_nation"` // ISO code of foreign choke-point nation
	StrategicApplications []string `json:"strategic_applications"`
}

// CanonicalNRCanCriticalMinerals provides reference baselines for key Canadian strategic resources.
var CanonicalNRCanCriticalMinerals = map[string]CriticalMineral{
	"LI": {
		Symbol:                "LI",
		Name:                  "Lithium",
		CanadianReservesRank:  5,
		DomesticExtractionPct: 25.0,
		DomesticRefiningPct:   10.0,
		ForeignMonopolyShare:  65.0,
		DominantForeignNation: "CN",
		StrategicApplications: []string{"EV Batteries", "Grid Energy Storage", "Aerospace Alloys"},
	},
	"NI": {
		Symbol:                "NI",
		Name:                  "Nickel",
		CanadianReservesRank:  6,
		DomesticExtractionPct: 85.0,
		DomesticRefiningPct:   55.0,
		ForeignMonopolyShare:  45.0,
		DominantForeignNation: "ID",
		StrategicApplications: []string{"Class-1 EV Cathodes", "Stainless Steel", "Defense Armor"},
	},
	"CO": {
		Symbol:                "CO",
		Name:                  "Cobalt",
		CanadianReservesRank:  7,
		DomesticExtractionPct: 60.0,
		DomesticRefiningPct:   30.0,
		ForeignMonopolyShare:  75.0,
		DominantForeignNation: "CN",
		StrategicApplications: []string{"Superalloys", "Lithium-Ion Batteries", "Magnets"},
	},
	"GRAPHITE": {
		Symbol:                "GRAPHITE",
		Name:                  "Natural Graphite",
		CanadianReservesRank:  8,
		DomesticExtractionPct: 35.0,
		DomesticRefiningPct:   5.0,
		ForeignMonopolyShare:  90.0,
		DominantForeignNation: "CN",
		StrategicApplications: []string{"Battery Anodes", "Nuclear Moderators", "Refractories"},
	},
	"REE": {
		Symbol:                "REE",
		Name:                  "Rare Earth Elements",
		CanadianReservesRank:  10,
		DomesticExtractionPct: 15.0,
		DomesticRefiningPct:   2.0,
		ForeignMonopolyShare:  85.0,
		DominantForeignNation: "CN",
		StrategicApplications: []string{"Permanent Magnets", "Wind Turbines", "Precision Guided Munitions", "MRI"},
	},
	"U": {
		Symbol:                "U",
		Name:                  "Uranium",
		CanadianReservesRank:  2,
		DomesticExtractionPct: 95.0,
		DomesticRefiningPct:   80.0,
		ForeignMonopolyShare:  20.0,
		DominantForeignNation: "KZ",
		StrategicApplications: []string{"CANDU Reactors", "SMR Small Modular Reactors", "Medical Isotopes"},
	},
	"CU": {
		Symbol:                "CU",
		Name:                  "Copper",
		CanadianReservesRank:  10,
		DomesticExtractionPct: 70.0,
		DomesticRefiningPct:   40.0,
		ForeignMonopolyShare:  45.0,
		DominantForeignNation: "CN",
		StrategicApplications: []string{"Grid Electrification", "EV Motors", "Renewable Transmission"},
	},
	"POTASH": {
		Symbol:                "POTASH",
		Name:                  "Potash (Potassium Chloride)",
		CanadianReservesRank:  1,
		DomesticExtractionPct: 100.0,
		DomesticRefiningPct:   100.0,
		ForeignMonopolyShare:  18.0,
		DominantForeignNation: "RU",
		StrategicApplications: []string{"Global Fertilizer Security", "Agricultural Yields", "Food Sovereignty"},
	},
}

// MineralProjectProfile profiles an economic project's contribution to Canadian supply sovereignty.
type MineralProjectProfile struct {
	ProjectID             string   `json:"project_id"`
	ProjectName           string   `json:"project_name"`
	TargetMineral         string   `json:"target_mineral"`
	HasExtraction         bool     `json:"has_extraction"`
	HasDomesticRefining   bool     `json:"has_domestic_refining"`
	HasEndProductMfg      bool     `json:"has_end_product_mfg"`
	AnnualCapacityTonnes  float64  `json:"annual_capacity_tonnes"`
	AlliedOfftakePct      float64  `json:"allied_offtake_pct"` // % off-take committed to Canada / CUSMA / NATO
	IndigenousEquityShare float64  `json:"indigenous_equity_share"` // 0.0 to 1.0
}

// SupplyChainVulnerability flags critical single points of failure in value chains.
type SupplyChainVulnerability struct {
	Severity    string `json:"severity"` // "CRITICAL", "HIGH", "MODERATE", "LOW"
	Category    string `json:"category"` // "REFINING_BOTTLENECK", "RAW_EXPORT_RISK", "FOREIGN_CHOKEPOINT"
	Description string `json:"description"`
}

// SelfSufficiencyEvaluation stores output metrics for Canadian supply chain resilience.
type SelfSufficiencyEvaluation struct {
	ProjectID              string                     `json:"project_id"`
	TargetMineral          string                     `json:"target_mineral"`
	SelfSufficiencyIndex   float64                    `json:"self_sufficiency_index"` // 0.0 to 100.0
	ExtractionScore        float64                    `json:"extraction_score"`
	RefiningScore          float64                    `json:"refining_score"`
	SovereignRetentionScore float64                   `json:"sovereign_retention_score"`
	ChokePointMitigation   float64                    `json:"choke_point_mitigation"`
	Vulnerabilities        []SupplyChainVulnerability `json:"vulnerabilities"`
	StrategicDirectives    []string                   `json:"strategic_directives"`
	EvaluatedAt            time.Time                  `json:"evaluated_at"`
}

// SupplyChainEvaluator calculates sovereign critical mineral and value-chain self-sufficiency.
type SupplyChainEvaluator struct{}

// NewSupplyChainEvaluator returns an initialized evaluator.
func NewSupplyChainEvaluator() *SupplyChainEvaluator {
	return &SupplyChainEvaluator{}
}

// Evaluate evaluates a project's self-sufficiency index and resilience impact.
func (e *SupplyChainEvaluator) Evaluate(profile MineralProjectProfile) *SelfSufficiencyEvaluation {
	sym := strings.ToUpper(strings.TrimSpace(profile.TargetMineral))
	refMineral, exists := CanonicalNRCanCriticalMinerals[sym]
	if !exists {
		// Fallback for custom minerals
		refMineral = CriticalMineral{
			Symbol:                sym,
			Name:                  sym,
			DomesticExtractionPct: 30.0,
			DomesticRefiningPct:   15.0,
			ForeignMonopolyShare:  60.0,
			DominantForeignNation: "CN",
		}
	}

	var vulns []SupplyChainVulnerability
	var directives []string

	// 1. Extraction Score (30%)
	extractionScore := 20.0
	if profile.HasExtraction {
		extractionScore = 80.0
		if profile.AnnualCapacityTonnes > 10000 {
			extractionScore = 100.0
		}
	} else {
		vulns = append(vulns, SupplyChainVulnerability{
			Severity:    "HIGH",
			Category:    "UPSTREAM_GAP",
			Description: fmt.Sprintf("Project lacks domestic extraction for %s, relying on upstream global imports.", refMineral.Name),
		})
	}

	// 2. Refining / Processing Score (35%)
	refiningScore := 10.0
	if profile.HasDomesticRefining {
		refiningScore = 85.0
		if profile.HasEndProductMfg {
			refiningScore = 100.0
		}
	} else if profile.HasExtraction {
		vulns = append(vulns, SupplyChainVulnerability{
			Severity:    "CRITICAL",
			Category:    "RAW_EXPORT_RISK",
			Description: fmt.Sprintf("Extracted %s raw ore will be exported unrefined abroad, leaving downstream multiplier value overseas.", refMineral.Name),
		})
		directives = append(directives, fmt.Sprintf("Mandate domestic processing or CUSMA refining joint venture before granting export permits for %s.", refMineral.Name))
	}

	// 3. Sovereign Retention & Allied Offtake (20%)
	offtake := profile.AlliedOfftakePct
	if offtake > 100.0 {
		offtake = 100.0
	}
	retentionScore := offtake * 0.8
	if profile.IndigenousEquityShare > 0.05 {
		retentionScore += math.Min(profile.IndigenousEquityShare*100.0*0.2, 20.0)
	}

	// 4. Foreign Choke-point Mitigation (15%)
	chokeMitigation := 50.0
	if refMineral.ForeignMonopolyShare > 60.0 {
		if profile.HasDomesticRefining {
			chokeMitigation = 100.0
			directives = append(directives, fmt.Sprintf("Project breaks foreign refining choke-point held by %s (%.0f%% global control).", refMineral.DominantForeignNation, refMineral.ForeignMonopolyShare))
		} else {
			chokeMitigation = 20.0
			vulns = append(vulns, SupplyChainVulnerability{
				Severity:    "CRITICAL",
				Category:    "FOREIGN_CHOKEPOINT",
				Description: fmt.Sprintf("%s is subject to extreme foreign concentration (%s holds %.0f%%). Domestic refining urgently needed.", refMineral.Name, refMineral.DominantForeignNation, refMineral.ForeignMonopolyShare),
			})
		}
	}

	// Composite Self-Sufficiency Index (0-100)
	totalIndex := (extractionScore * 0.30) +
		(refiningScore * 0.35) +
		(retentionScore * 0.20) +
		(chokeMitigation * 0.15)

	if totalIndex > 100.0 {
		totalIndex = 100.0
	}

	return &SelfSufficiencyEvaluation{
		ProjectID:              profile.ProjectID,
		TargetMineral:          refMineral.Name,
		SelfSufficiencyIndex:   math.Round(totalIndex*10) / 10,
		ExtractionScore:        math.Round(extractionScore*10) / 10,
		RefiningScore:          math.Round(refiningScore*10) / 10,
		SovereignRetentionScore: math.Round(retentionScore*10) / 10,
		ChokePointMitigation:   math.Round(chokeMitigation*10) / 10,
		Vulnerabilities:        vulns,
		StrategicDirectives:    directives,
		EvaluatedAt:            time.Now(),
	}
}
