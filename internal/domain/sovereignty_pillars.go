package domain

import "time"

// --- Pillar A: Critical Minerals Midstream & Allied Supply Chains ---

// CriticalMineralProcessingStage classifies where an asset sits on the value chain.
type CriticalMineralProcessingStage string

const (
	StageExtraction    CriticalMineralProcessingStage = "EXTRACTION"        // Raw spodumene, run-of-mine ore
	StageRefining      CriticalMineralProcessingStage = "REFINING"          // Hydro-metallurgical, chemical refining, nickel briquette
	StageCAMProduction CriticalMineralProcessingStage = "CAM_MANUFACTURING" // Cathode active material, precursor CAM
	StageRecycling     CriticalMineralProcessingStage = "RECYCLING"         // Black mass recycling, circular recovery
)

// CriticalMineralProjectSummary captures midstream value retention metrics.
type CriticalMineralProjectSummary struct {
	ProjectID             string                         `json:"project_id"`
	ProjectSlug           string                         `json:"project_slug"`
	ProjectName           string                         `json:"project_name"`
	Province              string                         `json:"province"`
	PrimaryMineral        string                         `json:"primary_mineral"`
	ProcessingStage       CriticalMineralProcessingStage `json:"processing_stage"`
	DomesticRetentionRate float64                        `json:"domestic_retention_rate"` // 0.0 to 1.0 (Value retained in Canada)
	CapexCAD              int64                          `json:"capex_cad"`
	AlliedOfftakeEligible bool                           `json:"allied_offtake_eligible"` // USMCA / G7 / Five Eyes
	ICAReviewStatus       string                         `json:"ica_review_status"`       // Investment Canada Act review
}

// CriticalMineralsSummary provides aggregate intelligence across Canada's critical mineral assets.
type CriticalMineralsSummary struct {
	TotalProjects        int                             `json:"total_projects"`
	TotalCapexCAD        int64                           `json:"total_capex_cad"`
	ExtractionCount      int                             `json:"extraction_count"`
	RefiningCount        int                             `json:"refining_count"`
	ManufacturingCount   int                             `json:"manufacturing_count"`
	RecyclingCount       int                             `json:"recycling_count"`
	AverageRetentionRate float64                         `json:"average_retention_rate"`
	TopMinerals          map[string]int                  `json:"top_minerals"`
	Projects             []CriticalMineralProjectSummary `json:"projects"`
	GeneratedAt          time.Time                       `json:"generated_at"`
}

// --- Pillar B: Indigenous Economic Sovereignty & $5B Loan Guarantee Mesh ---

// IndigenousLoanGuaranteeReq inputs for the sovereign financial simulation.
type IndigenousLoanGuaranteeReq struct {
	ProjectCapexCAD             int64   `json:"project_capex_cad"`
	IndigenousEquityPct         float64 `json:"indigenous_equity_pct"`          // 5.0 to 50.0 (%)
	LoanTermYears               int     `json:"loan_term_years"`                // e.g. 20-30 years
	BaseSeniorRatePct           float64 `json:"base_senior_rate_pct"`           // e.g. 6.5 (%)
	SovereignSpreadReductionBps int     `json:"sovereign_spread_reduction_bps"` // e.g. 85 bps
}

// IndigenousLoanGuaranteeResult output of the deterministic debt syndication model.
type IndigenousLoanGuaranteeResult struct {
	ProjectCapexCAD                    int64   `json:"project_capex_cad"`
	EquityAmountCAD                    int64   `json:"equity_amount_cad"`
	LoanGuaranteeAmountCAD             int64   `json:"loan_guarantee_amount_cad"`
	SovereignDiscountBps               int     `json:"sovereign_discount_bps"`
	GuaranteedSeniorRatePct            float64 `json:"guaranteed_senior_rate_pct"`
	AnnualDebtServiceSavingsCAD        int64   `json:"annual_debt_service_savings_cad"`
	CumulativeInterestSavingsCAD       int64   `json:"cumulative_interest_savings_cad"`
	ProjectedAnnualCommunityDividendCAD int64   `json:"projected_annual_community_dividend_cad"`
	ThirtyYearCumulativeDividendCAD    int64   `json:"thirty_year_cumulative_dividend_cad"`
	RecommendedFacility                string  `json:"recommended_facility"` // e.g. "FEDERAL_ILGP", "AIOC", "ALGP"
}

// IndigenousProjectOverview represents a partnered capital project.
type IndigenousProjectOverview struct {
	ProjectID          string   `json:"project_id"`
	ProjectName        string   `json:"project_name"`
	Province           string   `json:"province"`
	NationsPartnered   []string `json:"nations_partnered"`
	IBAStatus          string   `json:"iba_status"` // "EXECUTED", "IN_NEGOTIATION", "PLANNED"
	EquityPct          float64  `json:"equity_pct"`
	GuaranteedFacility string   `json:"guaranteed_facility"`
	CapexCAD           int64    `json:"capex_cad"`
}

// IndigenousOverviewSummary aggregates First Nations major projects co-ownership.
type IndigenousOverviewSummary struct {
	TotalPartneredProjects    int                         `json:"total_partnered_projects"`
	TotalEquityGuaranteedCAD  int64                       `json:"total_equity_guaranteed_cad"`
	AverageEquityPct          float64                     `json:"average_equity_pct"`
	FederalFacilityAllocation int64                       `json:"federal_facility_allocation_cad"` // Out of $5B
	Projects                  []IndigenousProjectOverview `json:"projects"`
	GeneratedAt               time.Time                   `json:"generated_at"`
}

// --- Pillar C: Inter-Provincial Trade & Regulatory Friction ---

// InterProvincialFriction models trade barrier costs between provinces.
type InterProvincialFriction struct {
	CorridorID            string   `json:"corridor_id"`
	OriginProvince        string   `json:"origin_province"`
	DestProvince          string   `json:"dest_province"`
	AnnualTradeVolumeCAD  int64    `json:"annual_trade_volume_cad"`
	FrictionTaxCAD        int64    `json:"friction_tax_cad"`       // Estimated cost of internal barriers
	BarrierIndex          float64  `json:"barrier_index"`          // 1.0 (low) to 10.0 (high)
	DivergenceAreas       []string `json:"divergence_areas"`       // "Axle weights", "Environmental duplication", "Labor certification"
	HarmonizationStatus   string   `json:"harmonization_status"`   // "FRAGMENTED", "ACTIVE_MOU", "HARMONIZED"
	PotentialSavingsCAD   int64    `json:"potential_savings_cad"`  // One Project One Assessment dividend
	ScheduleCompressWeeks int      `json:"schedule_compress_weeks"`
}

// TradeFrictionReport summarizes Canada-wide internal market friction.
type TradeFrictionReport struct {
	TotalAnnualFrictionTaxCAD      int64                     `json:"total_annual_friction_tax_cad"`
	NationalHarmonizationDividendCAD int64                   `json:"national_harmonization_dividend_cad"`
	AverageBarrierIndex            float64                   `json:"average_barrier_index"`
	Corridors                      []InterProvincialFriction `json:"corridors"`
	GeneratedAt                    time.Time                 `json:"generated_at"`
}

// --- Pillar D: Clean Baseload Power vs Sovereign AI Hyperscale Nexus ---

// ProvincialGridProfile models clean energy supply vs compute demand by province.
type ProvincialGridProfile struct {
	Province                  string  `json:"province"`
	GridAuthority             string  `json:"grid_authority"` // "IESO", "HYDRO_QUEBEC", "BC_HYDRO", "AESO"
	TotalGenerationCapacityMW float64 `json:"total_generation_capacity_mw"`
	CleanEnergyPct            float64 `json:"clean_energy_pct"`
	FirmBaseloadType          string  `json:"firm_baseload_type"` // "NUCLEAR", "HYDRO", "GAS_CCS", "MIXED"
	AIComputeHeadroomMW       float64 `json:"ai_compute_headroom_mw"`
	PlannedIndustrialDemandMW float64 `json:"planned_industrial_demand_mw"`
	SovereignComputeExaFLOPs  float64 `json:"sovereign_compute_exaflops"` // Domestic compute potential
	CleanFLOPsPerMegawatt     float64 `json:"clean_flops_per_megawatt"`
}

// BaseloadComputeSummary provides national grid vs compute allocation overview.
type BaseloadComputeSummary struct {
	TotalCleanCapacityMW     float64                 `json:"total_clean_capacity_mw"`
	TotalAIComputeHeadroomMW float64                 `json:"total_ai_compute_headroom_mw"`
	TotalSovereignExaFLOPs   float64                 `json:"total_sovereign_exaflops"`
	AverageCleanGridPct      float64                 `json:"average_clean_grid_pct"`
	Grids                    []ProvincialGridProfile `json:"grids"`
	GeneratedAt              time.Time               `json:"generated_at"`
}
