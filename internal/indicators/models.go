package indicators

import (
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// KPICategory represents the strategic intelligence domain of an indicator.
type KPICategory string

const (
	CategoryESGDecarbonization KPICategory = "ESG_DECARBONIZATION"
	CategoryIndigenousEquity   KPICategory = "INDIGENOUS_EQUITY"
	CategoryCapitalVelocity    KPICategory = "CAPITAL_VELOCITY"
	CategorySupplyChainContent KPICategory = "SUPPLY_CHAIN_CONTENT"
	CategoryGridPhysics        KPICategory = "GRID_PHYSICS"
	CategoryRegulatorySpeed    KPICategory = "REGULATORY_SPEED"
	CategoryLaborSkills        KPICategory = "LABOR_SKILLS"
	CategoryCommodityMacro     KPICategory = "COMMODITY_MACRO"
)

// PreferredDirection specifies whether higher or lower values indicate stronger performance.
type PreferredDirection string

const (
	DirectionHigherBetter PreferredDirection = "HIGHER_BETTER"
	DirectionLowerBetter  PreferredDirection = "LOWER_BETTER"
	DirectionNeutral      PreferredDirection = "NEUTRAL"
)

// KPICode is a unique alphanumeric identifier for an official or derived metric.
type KPICode string

const (
	// Pillar 1: ESG & Decarbonization
	KPIEmissionsIntensityScope12 KPICode = "ESG.GHG.INTENSITY.SCOPE1_2"
	KPIEmissionsAbatementAnnual  KPICode = "ESG.GHG.ABATEMENT.ANNUAL_MT"
	KPICarbonTaxSensitivity      KPICode = "ESG.CARBON.TAX_SENSITIVITY.CAD"
	KPIWaterRecyclingRatio       KPICode = "ESG.WATER.RECYCLING.RATIO"
	KPICleanPowerPurityPct       KPICode = "ESG.POWER.CLEAN_PURITY.PCT"

	// Pillar 2: Indigenous Economic Sovereignty
	KPIIndigenousEquityPct         KPICode = "INDIG.EQUITY.OWNERSHIP.PCT"
	KPIIndigenousProcurementShare  KPICode = "INDIG.PROCUREMENT.SHARE.PCT"
	KPIIndigenousWorkforcePct      KPICode = "INDIG.WORKFORCE.PARTICIPATION.PCT"
	KPIIndigenousILGPGuaranteeCAD  KPICode = "INDIG.ILGP.GUARANTEE.CAD"
	KPIIndigenousCommunityDividend KPICode = "INDIG.DIVIDEND.COMMUNITY.CAD_YR"

	// Pillar 3: Capital Spend Velocity & Hazards
	KPICapexSpendVelocityCADMo  KPICode = "CAPEX.SPEND.VELOCITY.CAD_MO"
	KPIScheduleSlipMonths       KPICode = "SCHEDULE.SLIP.MONTHS"
	KPICostOverrunHazardPct     KPICode = "RISK.COST_OVERRUN.HAZARD.PCT"
	KPIPublicCrowdingMultiplier KPICode = "FINANCE.PUBLIC_CROWDING.MULTIPLIER"

	// Pillar 4: Sovereign Supply Chain & Domestic Content
	KPIDomesticCanadianContentPct KPICode = "SUPPLY.DOMESTIC_CONTENT.PCT"
	KPICriticalBottleneckScore    KPICode = "SUPPLY.BOTTLENECK.LEAD_TIME.SCORE"
	KPIUSTariffExposurePct        KPICode = "TRADE.US_TARIFF.EXPOSURE.PCT"
	KPIDomesticImportReplacement  KPICode = "TRADE.IMPORT_REPLACEMENT.CAD_YR"

	// Pillar 5: Power, Grid Interconnection & Hosting
	KPIPeakPowerDemandMW        KPICode = "GRID.PEAK_POWER.DEMAND.MW"
	KPIGridInterconnectQueueMo  KPICode = "GRID.QUEUE.WAIT_TIME.MONTHS"
	KPIGridHostingHeadroomMW    KPICode = "GRID.HOSTING.HEADROOM.MW"
	KPISystemReinforcementCost  KPICode = "GRID.REINFORCEMENT.COST.CAD"

	// Pillar 6: Regulatory, Permitting & Legal Latency
	KPIIAACStatutoryDurationMo KPICode = "REG.IAAC.REVIEW_DURATION.MONTHS"
	KPIPermitCompletionPct     KPICode = "REG.PERMIT.COMPLETION.PCT"
	KPIJudicialReviewRiskIndex KPICode = "REG.LITIGATION.RISK_INDEX"
	KPIConditionsCount         KPICode = "REG.CONDITIONS.LEGALLY_BINDING.COUNT"

	// Pillar 7: Labour Dynamics & Critical Skills Pipeline
	KPIRedSealTradesGap        KPICode = "LABOR.RED_SEAL.TRADES_GAP.FTE"
	KPIApprenticeJourneyRatio  KPICode = "LABOR.APPRENTICE.RATIO.PCT"
	KPIHousingAbsorptionDeficit KPICode = "LABOR.HOUSING.ABSORPTION_DEFICIT.UNITS"
	KPITotalPersonYearsCreated KPICode = "LABOR.TOTAL.PERSON_YEARS.COUNT"

	// Pillar 8: Live Commodity & Macro Benchmark Feeds
	KPIWCSWTIDifferential   KPICode = "COMMODITY.WCS_WTI.DIFF.USD_BBL"
	KPIAECOGasSpotCAD       KPICode = "COMMODITY.AECO.GAS.CAD_GJ"
	KPILMENickelCashUSD     KPICode = "COMMODITY.LME.NICKEL.USD_T"
	KPIUraniumUXU3O8USD     KPICode = "COMMODITY.UX.U3O8.USD_LB"
	KPILithiumCarbonateUSD  KPICode = "COMMODITY.LITHIUM.CARBONATE.USD_T"
	KPICADUSDExchangeRate   KPICode = "MACRO.CAD_USD.SPOT_FX"
	KPIBoCOvernightRate     KPICode = "MACRO.BOC.POLICY_RATE.PCT"
	KPIGoC10YearBondYield   KPICode = "MACRO.GOC.10Y_YIELD.PCT"
)

// KPIDefinition defines the metadata, targets, and statutory reference of an indicator.
type KPIDefinition struct {
	Code               KPICode            `json:"code"`
	Name               string             `json:"name"`
	Category           KPICategory        `json:"category"`
	Description        string             `json:"description"`
	Unit               string             `json:"unit"`
	PreferredDirection PreferredDirection `json:"preferred_direction"`
	TargetBenchmark    float64            `json:"target_benchmark"`
	BenchmarkUnit      string             `json:"benchmark_unit"`
	BenchmarkLabel     string             `json:"benchmark_label"`
	StatutoryBasis     string             `json:"statutory_basis"`
	Publisher          string             `json:"publisher"`
	UpdateFrequency    string             `json:"update_frequency"`
	IsLiveFeed         bool               `json:"is_live_feed"`
}

// KPIObservation is an individual observed data point with complete cryptographic provenance.
type KPIObservation struct {
	ID              string                 `json:"id"`
	MetricCode      KPICode                `json:"metric_code"`
	MetricName      string                 `json:"metric_name"`
	Category        KPICategory            `json:"category"`
	Scope           string                 `json:"scope"` // "NATIONAL" or Project ID/Slug
	ProjectID       string                 `json:"project_id,omitempty"`
	Value           float64                `json:"value"`
	Unit            string                 `json:"unit"`
	ReferencePeriod string                 `json:"reference_period"`
	ObservedAt      time.Time              `json:"observed_at"`
	SourceURL       string                 `json:"source_url"`
	Publisher       string                 `json:"publisher"`
	Confidence      domain.ConfidenceLevel `json:"confidence"`
	ContentHash     string                 `json:"content_hash"`
	EvidenceID      string                 `json:"evidence_id,omitempty"`
}

// ProjectMetricAssessment represents a single assessed metric within a project scorecard.
type ProjectMetricAssessment struct {
	Code            KPICode            `json:"code"`
	Name            string             `json:"name"`
	Category        KPICategory        `json:"category"`
	ObservedValue   float64            `json:"observed_value"`
	Unit            string             `json:"unit"`
	TargetBenchmark float64            `json:"target_benchmark"`
	Variance        float64            `json:"variance"`
	PerformanceRank string             `json:"performance_rank"` // "SUPERIOR", "ON_TARGET", "NEEDS_IMPROVEMENT", "CRITICAL_GAP"
	ScoreNormalized float64            `json:"score_normalized"` // 0 - 100
	Confidence      domain.ConfidenceLevel `json:"confidence"`
	Notes           string             `json:"notes"`
	Source          string             `json:"source"`
}

// PillarScore summarizes one of the 8 strategic pillars.
type PillarScore struct {
	Category    KPICategory `json:"category"`
	Title       string      `json:"title"`
	Score       float64     `json:"score"` // 0 - 100
	Health      string      `json:"health"` // "EXEMPLARY", "HEALTHY", "ATTENTION_REQUIRED", "HIGH_RISK"
	Metrics     []*ProjectMetricAssessment `json:"metrics"`
	KeyFindings []string    `json:"key_findings"`
}

// ProjectKPIScorecard provides the complete 8-pillar indicator intelligence dossier for a project.
type ProjectKPIScorecard struct {
	ProjectID          string         `json:"project_id"`
	ProjectSlug        string         `json:"project_slug"`
	ProjectName        string         `json:"project_name"`
	Sector             domain.Sector  `json:"sector"`
	Province           string         `json:"province"`
	CurrentStage       domain.LifecycleStage `json:"current_stage"`
	TotalCapexCAD      int64          `json:"total_capex_cad"`
	OverallKPIRating   float64        `json:"overall_kpi_rating"` // 0 - 100
	Pillars            []*PillarScore `json:"pillars"`
	CriticalActionGaps []string       `json:"critical_action_gaps"`
	AuditHash          string         `json:"audit_hash"`
	EvaluatedAt        time.Time      `json:"evaluated_at"`
}

// LiveFeedTick represents a real-time streaming market or indicator tick.
type LiveFeedTick struct {
	SequenceID     int64     `json:"sequence_id"`
	MetricCode     KPICode   `json:"metric_code"`
	Name           string    `json:"name"`
	Category       KPICategory `json:"category"`
	Value          float64   `json:"value"`
	Unit           string    `json:"unit"`
	ChangeAbsolute float64   `json:"change_absolute"`
	ChangePercent  float64   `json:"change_percent"`
	Direction      string    `json:"direction"` // "UP", "DOWN", "FLAT"
	Source         string    `json:"source"`
	Timestamp      time.Time `json:"timestamp"`
	Hash           string    `json:"hash"`
}

// MacroKPISummary presents national aggregate metrics and cross-sector rollups.
type MacroKPISummary struct {
	NationalEmissionsAbatementMt float64   `json:"national_emissions_abatement_mt"`
	AverageIndigenousEquityPct   float64   `json:"average_indigenous_equity_pct"`
	TotalILGPAllocatedCAD        int64     `json:"total_ilgp_allocated_cad"`
	NationalSpendRunRateCADMo    float64   `json:"national_spend_run_rate_cad_mo"`
	AverageDomesticContentPct    float64   `json:"average_domestic_content_pct"`
	AverageGridQueueWaitMonths   float64   `json:"average_grid_queue_wait_months"`
	AverageIAACReviewDurationMo  float64   `json:"average_iaac_review_duration_mo"`
	NationalRedSealDeficitFTE    int       `json:"national_red_seal_deficit_fte"`
	ActiveFeedTicksCount         int       `json:"active_feed_ticks_count"`
	LatestCommodityTicks         []*LiveFeedTick `json:"latest_commodity_ticks"`
	GeneratedAt                  time.Time `json:"generated_at"`
}

// KPISnapshot holds the serialized state for external consumption and web bundling.
type KPISnapshot struct {
	Version      string            `json:"version"`
	GeneratedAt  time.Time         `json:"generated_at"`
	Definitions  []*KPIDefinition  `json:"definitions"`
	Observations []*KPIObservation `json:"observations"`
	MacroSummary *MacroKPISummary  `json:"macro_summary"`
	AuditHash    string            `json:"audit_hash"`
}
