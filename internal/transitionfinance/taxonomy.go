package transitionfinance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strings"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// TaxonomyCategory classifies economic activities under the GFANZ & SFAC Canadian Transition Taxonomy.
type TaxonomyCategory string

const (
	// CategoryGreen covers zero or near-zero operational emissions activities (nuclear, hydro, solar, wind, storage).
	CategoryGreen TaxonomyCategory = "GREEN"
	// CategoryTransition covers hard-to-abate activities with credible, science-based 1.5°C decarbonization (EAF, CCUS).
	CategoryTransition TaxonomyCategory = "TRANSITION"
	// CategoryEnabling covers activities and critical supply chains that enable system-wide decarbonization (interties, minerals).
	CategoryEnabling TaxonomyCategory = "ENABLING"
	// CategoryPhaseOut covers managed early retirement of high-carbon assets with strict decommissioning covenants.
	CategoryPhaseOut TaxonomyCategory = "PHASE_OUT"
	// CategoryUnclassified represents activities lacking verifiable emissions disclosures or transition plans.
	CategoryUnclassified TaxonomyCategory = "UNCLASSIFIED"
)

// CredibilityRating classifies the rigor of an asset's transition plan.
type CredibilityRating string

const (
	RatingHigh         CredibilityRating = "HIGH"
	RatingMedium       CredibilityRating = "MEDIUM"
	RatingLow          CredibilityRating = "LOW"
	RatingInsufficient CredibilityRating = "INSUFFICIENT_EVIDENCE"
)

// TransitionCredibilityIndex quantifies the 0–100 credibility of an asset's decarbonization plan.
type TransitionCredibilityIndex struct {
	OverallScore               float64           `json:"overall_score"`
	DecarbonizationPathwayScore float64           `json:"decarbonization_pathway_score"` // 0–25: interim 2030 and 2050 net-zero targets
	Scope3DisclosureScore      float64           `json:"scope3_disclosure_score"`       // 0–20: upstream/downstream boundary coverage
	CapexAlignmentScore        float64           `json:"capex_alignment_score"`         // 0–25: proportion of capex allocated to abatement
	LockInMitigationScore      float64           `json:"lock_in_mitigation_score"`       // 0–15: avoidance of emissions lock-in
	ThirdPartyAssuranceScore   float64           `json:"third_party_assurance_score"`    // 0–15: verified by independent engineering/audit
	Rating                     CredibilityRating `json:"rating"`
}

// AbatementMetrics models emissions reduction and marginal abatement economics.
type AbatementMetrics struct {
	BaselineEmissionsTpyCO2e    float64 `json:"baseline_emissions_tpy_co2e"`
	TargetEmissionsTpyCO2e      float64 `json:"target_emissions_tpy_co2e"`
	AnnualAbatementTpyCO2e      float64 `json:"annual_abatement_tpy_co2e"`
	LifetimeAbatementTonnesCO2e float64 `json:"lifetime_abatement_tonnes_co2e"`
	EconomicLifetimeYears       int     `json:"economic_lifetime_years"`
	MarginalAbatementCostCAD    float64 `json:"marginal_abatement_cost_cad_per_tonne"` // Capex / Lifetime Abatement
	PrivateCapitalPerTonneCAD   float64 `json:"private_capital_per_tonne_cad"`
	CrowdingInMultiplier        float64 `json:"crowding_in_multiplier"`
}

// TransitionAssessment represents the institutional taxonomy evaluation for a single project.
type TransitionAssessment struct {
	ProjectID                  string                     `json:"project_id"`
	ProjectName                string                     `json:"project_name"`
	Sector                     domain.Sector              `json:"sector"`
	Subsector                  string                     `json:"subsector"`
	TaxonomyCategory           TaxonomyCategory           `json:"taxonomy_category"`
	CredibilityIndex           TransitionCredibilityIndex `json:"credibility_index"`
	Abatement                  AbatementMetrics           `json:"abatement"`
	EligibleGreenCapexCAD      int64                      `json:"eligible_green_capex_cad"`
	EligibleTransitionCapexCAD int64                      `json:"eligible_transition_capex_cad"`
	LockInRiskIdentified       bool                       `json:"lock_in_risk_identified"`
	NoLockInConditionMet       bool                       `json:"no_lock_in_condition_met"`
	GFANZPillar                string                     `json:"gfanz_pillar"`
	AuditHash                  string                     `json:"audit_hash"`
}

// PortfolioTransitionSummary aggregates taxonomy metrics across an indexed portfolio.
type PortfolioTransitionSummary struct {
	TotalProjects                     int                     `json:"total_projects"`
	GreenProjectsCount                int                     `json:"green_projects_count"`
	TransitionProjectsCount           int                     `json:"transition_projects_count"`
	EnablingProjectsCount             int                     `json:"enabling_projects_count"`
	PhaseOutProjectsCount             int                     `json:"phase_out_projects_count"`
	TotalGreenCapexCAD                int64                   `json:"total_green_capex_cad"`
	TotalTransitionCapexCAD           int64                   `json:"total_transition_capex_cad"`
	AggregateAnnualAbatementTpy       float64                 `json:"aggregate_annual_abatement_tpy"`
	AggregateLifetimeAbatementTonnes  float64                 `json:"aggregate_lifetime_abatement_tonnes"`
	WeightedAverageCredibilityIndex   float64                 `json:"weighted_average_credibility_index"`
	WeightedAverageMAC                float64                 `json:"weighted_average_marginal_abatement_cost_cad"`
	Assessments                       []*TransitionAssessment `json:"assessments"`
	AuditHash                         string                  `json:"audit_hash"`
}

// TransitionEngine executes GFANZ and Sustainable Finance Action Council (SFAC) taxonomy evaluations.
type TransitionEngine struct{}

// NewTransitionEngine instantiates a transition taxonomy engine.
func NewTransitionEngine() *TransitionEngine {
	return &TransitionEngine{}
}

// EvaluateProject classifies a project under the Canadian Transition Finance Taxonomy.
func (e *TransitionEngine) EvaluateProject(project *domain.Project) *TransitionAssessment {
	if project == nil {
		return nil
	}

	capex := project.CapexCAD
	if capex <= 0 {
		capex = 500_000_000
	}

	category := CategoryUnclassified
	gfanzPillar := "Under Review"
	var greenCapex, transCapex int64
	var baseEmissions, targetEmissions float64
	lifetimeYears := 30
	lockInRisk := false
	noLockInMet := true

	// Scoring component factors
	var pathwayScore, scope3Score, capexScore, lockInScore, assuranceScore float64

	subLower := strings.ToLower(project.Subsector)
	nameLower := strings.ToLower(project.Name)
	summaryLower := strings.ToLower(project.Summary)

	switch project.Sector {
	case domain.SectorNuclearEnergy:
		category = CategoryGreen
		gfanzPillar = "Climate Solutions (Zero-Emissions Baseload)"
		greenCapex = capex
		transCapex = 0
		pathwayScore = 25.0
		scope3Score = 18.0
		capexScore = 25.0
		lockInScore = 15.0
		assuranceScore = 14.5
		baseEmissions = float64(capex) * 0.00045 // Displacing natural gas/coal generation
		targetEmissions = baseEmissions * 0.02
		lifetimeYears = 60

	case domain.SectorCleanEnergy:
		if strings.Contains(subLower, "hydrogen") || strings.Contains(nameLower, "hydrogen") {
			category = CategoryTransition
			gfanzPillar = "Aligned Decarbonization Fuel Switch"
			greenCapex = int64(float64(capex) * 0.40)
			transCapex = int64(float64(capex) * 0.60)
			pathwayScore = 22.0
			scope3Score = 16.0
			capexScore = 23.0
			lockInScore = 13.0
			assuranceScore = 13.0
			baseEmissions = float64(capex) * 0.00035
			targetEmissions = baseEmissions * 0.15
			lifetimeYears = 25
		} else {
			category = CategoryGreen
			gfanzPillar = "Climate Solutions (Renewable Energy & Storage)"
			greenCapex = capex
			transCapex = 0
			pathwayScore = 24.5
			scope3Score = 17.0
			capexScore = 24.0
			lockInScore = 15.0
			assuranceScore = 14.0
			baseEmissions = float64(capex) * 0.00030
			targetEmissions = baseEmissions * 0.03
			lifetimeYears = 30
		}

	case domain.SectorCriticalMinerals:
		category = CategoryEnabling
		gfanzPillar = "Enabling Clean Energy Supply Chains"
		greenCapex = int64(float64(capex) * 0.35)
		transCapex = int64(float64(capex) * 0.65)
		pathwayScore = 21.0
		scope3Score = 15.0
		capexScore = 22.0
		lockInScore = 13.5
		assuranceScore = 12.5

		// Check for carbon capture mineralization (e.g. Crawford Nickel)
		if strings.Contains(nameLower, "crawford") || strings.Contains(summaryLower, "carbon capture") {
			category = CategoryTransition
			gfanzPillar = "Aligned Heavy Extraction with Direct Mineral Carbonation"
			pathwayScore = 24.0
			scope3Score = 18.0
			lockInScore = 14.5
			baseEmissions = 1_200_000
			targetEmissions = 150_000
		} else {
			baseEmissions = float64(capex) * 0.00025
			targetEmissions = baseEmissions * 0.35
		}
		lifetimeYears = 30

	case domain.SectorIndustrialMfg:
		category = CategoryTransition
		gfanzPillar = "Decarbonizing Hard-to-Abate Heavy Industry"
		transCapex = int64(float64(capex) * 0.85)
		greenCapex = int64(float64(capex) * 0.15)
		pathwayScore = 22.5
		scope3Score = 16.5
		capexScore = 21.0
		lockInScore = 14.0
		assuranceScore = 13.0
		baseEmissions = float64(capex) * 0.00060
		targetEmissions = baseEmissions * 0.20
		lifetimeYears = 30

	case domain.SectorTransportation:
		if strings.Contains(nameLower, "airport") || strings.Contains(subLower, "airport") {
			category = CategoryEnabling
			gfanzPillar = "Modernized Intermodal & Low-Carbon Aviation Hubs"
			greenCapex = int64(float64(capex) * 0.30)
			transCapex = int64(float64(capex) * 0.70)
			pathwayScore = 20.5
			scope3Score = 15.0
			capexScore = 20.0
			lockInScore = 12.0
			assuranceScore = 12.0
			baseEmissions = float64(capex) * 0.00020
			targetEmissions = baseEmissions * 0.40
			lifetimeYears = 40
		} else {
			category = CategoryEnabling
			gfanzPillar = "Low-Carbon Freight & Rail Intermodal Corridors"
			greenCapex = int64(float64(capex) * 0.50)
			transCapex = int64(float64(capex) * 0.50)
			pathwayScore = 21.0
			scope3Score = 15.5
			capexScore = 21.0
			lockInScore = 13.0
			assuranceScore = 13.0
			baseEmissions = float64(capex) * 0.00022
			targetEmissions = baseEmissions * 0.30
			lifetimeYears = 35
		}

	case domain.SectorAICompute:
		category = CategoryEnabling
		gfanzPillar = "Sovereign Low-Carbon Digital Infrastructure"
		greenCapex = int64(float64(capex) * 0.80)
		transCapex = int64(float64(capex) * 0.20)
		pathwayScore = 23.0
		scope3Score = 17.0
		capexScore = 23.5
		lockInScore = 14.0
		assuranceScore = 13.5
		baseEmissions = float64(capex) * 0.00015
		targetEmissions = baseEmissions * 0.05
		lifetimeYears = 20

	default:
		category = CategoryTransition
		gfanzPillar = "Industrial Decarbonization"
		transCapex = int64(float64(capex) * 0.60)
		greenCapex = int64(float64(capex) * 0.40)
		pathwayScore = 18.0
		scope3Score = 12.0
		capexScore = 18.0
		lockInScore = 11.0
		assuranceScore = 11.0
		baseEmissions = float64(capex) * 0.00020
		targetEmissions = baseEmissions * 0.50
		lifetimeYears = 25
	}

	annualAbatement := math.Max(0, baseEmissions-targetEmissions)
	lifetimeAbatement := annualAbatement * float64(lifetimeYears)

	mac := 0.0
	if lifetimeAbatement > 0 {
		mac = float64(capex) / lifetimeAbatement
	}

	privateCapPerTonne := 0.0
	privateCapRatio := 0.70 // 70% private institutional capital
	if lifetimeAbatement > 0 {
		privateCapPerTonne = (float64(capex) * privateCapRatio) / lifetimeAbatement
	}

	overallScore := pathwayScore + scope3Score + capexScore + lockInScore + assuranceScore
	if overallScore > 100.0 {
		overallScore = 100.0
	}

	rating := RatingLow
	if overallScore >= 80.0 {
		rating = RatingHigh
	} else if overallScore >= 60.0 {
		rating = RatingMedium
	} else if overallScore < 40.0 {
		rating = RatingInsufficient
	}

	credIndex := TransitionCredibilityIndex{
		OverallScore:                math.Round(overallScore*10) / 10,
		DecarbonizationPathwayScore: pathwayScore,
		Scope3DisclosureScore:       scope3Score,
		CapexAlignmentScore:         capexScore,
		LockInMitigationScore:       lockInScore,
		ThirdPartyAssuranceScore:    assuranceScore,
		Rating:                      rating,
	}

	abatement := AbatementMetrics{
		BaselineEmissionsTpyCO2e:    math.Round(baseEmissions),
		TargetEmissionsTpyCO2e:      math.Round(targetEmissions),
		AnnualAbatementTpyCO2e:      math.Round(annualAbatement),
		LifetimeAbatementTonnesCO2e: math.Round(lifetimeAbatement),
		EconomicLifetimeYears:       lifetimeYears,
		MarginalAbatementCostCAD:    math.Round(mac*100) / 100,
		PrivateCapitalPerTonneCAD:   math.Round(privateCapPerTonne*100) / 100,
		CrowdingInMultiplier:        3.8, // Canonical Canadian crowding-in ratio
	}

	hasher := sha256.New()
	hasher.Write(fmt.Appendf(nil, "%s|%s|%s|%.2f|%.0f|%d",
		project.ID,
		category,
		rating,
		overallScore,
		lifetimeAbatement,
		capex,
	))
	auditHash := hex.EncodeToString(hasher.Sum(nil))

	return &TransitionAssessment{
		ProjectID:                  project.ID,
		ProjectName:                project.Name,
		Sector:                     project.Sector,
		Subsector:                  project.Subsector,
		TaxonomyCategory:           category,
		CredibilityIndex:           credIndex,
		Abatement:                  abatement,
		EligibleGreenCapexCAD:      greenCapex,
		EligibleTransitionCapexCAD: transCapex,
		LockInRiskIdentified:       lockInRisk,
		NoLockInConditionMet:       noLockInMet,
		GFANZPillar:                gfanzPillar,
		AuditHash:                  auditHash,
	}
}

// EvaluatePortfolio generates an aggregated transition assessment across multiple projects.
func (e *TransitionEngine) EvaluatePortfolio(projects []*domain.Project) *PortfolioTransitionSummary {
	summary := &PortfolioTransitionSummary{
		Assessments: make([]*TransitionAssessment, 0, len(projects)),
	}

	var totalScoreWeight float64
	var totalCapex int64
	var totalWeightedMAC float64

	for _, p := range projects {
		if p == nil {
			continue
		}
		a := e.EvaluateProject(p)
		summary.Assessments = append(summary.Assessments, a)
		summary.TotalProjects++

		switch a.TaxonomyCategory {
		case CategoryGreen:
			summary.GreenProjectsCount++
		case CategoryTransition:
			summary.TransitionProjectsCount++
		case CategoryEnabling:
			summary.EnablingProjectsCount++
		case CategoryPhaseOut:
			summary.PhaseOutProjectsCount++
		}

		summary.TotalGreenCapexCAD += a.EligibleGreenCapexCAD
		summary.TotalTransitionCapexCAD += a.EligibleTransitionCapexCAD
		summary.AggregateAnnualAbatementTpy += a.Abatement.AnnualAbatementTpyCO2e
		summary.AggregateLifetimeAbatementTonnes += a.Abatement.LifetimeAbatementTonnesCO2e

		projCapex := p.CapexCAD
		if projCapex <= 0 {
			projCapex = 500_000_000
		}
		totalCapex += projCapex
		totalScoreWeight += a.CredibilityIndex.OverallScore * float64(projCapex)
		totalWeightedMAC += a.Abatement.MarginalAbatementCostCAD * a.Abatement.LifetimeAbatementTonnesCO2e
	}

	if totalCapex > 0 {
		summary.WeightedAverageCredibilityIndex = math.Round((totalScoreWeight/float64(totalCapex))*10) / 10
	}
	if summary.AggregateLifetimeAbatementTonnes > 0 {
		summary.WeightedAverageMAC = math.Round((totalWeightedMAC/summary.AggregateLifetimeAbatementTonnes)*100) / 100
	}

	hasher := sha256.New()
	hasher.Write(fmt.Appendf(nil, "portfolio|%d|%d|%d|%.2f|%.2f",
		summary.TotalProjects,
		summary.TotalGreenCapexCAD,
		summary.TotalTransitionCapexCAD,
		summary.AggregateLifetimeAbatementTonnes,
		summary.WeightedAverageCredibilityIndex,
	))
	summary.AuditHash = hex.EncodeToString(hasher.Sum(nil))

	return summary
}
