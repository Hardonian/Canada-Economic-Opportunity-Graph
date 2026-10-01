package concession

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"time"
)

// AirportHubCode identifies the major Canadian National Airports System (NAS) gateway hubs.
type AirportHubCode string

const (
	HubYYZ AirportHubCode = "YYZ" // Toronto Pearson International Airport
	HubYVR AirportHubCode = "YVR" // Vancouver International Airport
	HubYUL AirportHubCode = "YUL" // Montréal-Pierre Elliott Trudeau International Airport
	HubYYC AirportHubCode = "YYC" // Calgary International Airport
	HubYEG AirportHubCode = "YEG" // Edmonton International Airport
)

// TotalNASModernizationTargetCapexCAD defines the $18B CAD program capex announced by Mark Carney.
const TotalNASModernizationTargetCapexCAD int64 = 18_000_000_000

// AirportConcessionProfile defines baseline operational and infrastructure parameters for a hub.
type AirportConcessionProfile struct {
	Code                         AirportHubCode `json:"code"`
	Name                         string         `json:"name"`
	Province                     string         `json:"province"`
	AnnualPassengersBaseline     int64          `json:"annual_passengers_baseline"`
	TargetPassengers2050         int64          `json:"target_passengers_2050"`
	BaselineRevenueCAD           int64          `json:"baseline_revenue_cad"`
	BaselineEBITDACAD            int64          `json:"baseline_ebitda_cad"`
	TargetModernizationCapexCAD  int64          `json:"target_modernization_capex_cad"`
	KeyProjects                  []string       `json:"key_projects"`
	DecarbonizationInitiatives   []string       `json:"decarbonization_initiatives"`
	TargetPensionInvestors       []string       `json:"target_pension_investors"`
}

// SimulationParams specifies inputs for the NAS ground lease concession model.
type SimulationParams struct {
	ConcessionHorizonYears       int                    `json:"concession_horizon_years"`        // 30 to 50 years, default 40
	FederalRoyaltyRatePercent    float64                `json:"federal_royalty_rate_percent"`    // e.g., 10.0% of gross revenue
	PensionEquitySharePercent    float64                `json:"pension_equity_share_percent"`    // e.g., 50.0%
	CommercialDebtSharePercent   float64                `json:"commercial_debt_share_percent"`   // e.g., 45.0%
	FederalSubordinatedSharePct  float64                `json:"federal_subordinated_share_pct"`  // e.g., 5.0%
	PassengerCAGRPercent         float64                `json:"passenger_cagr_percent"`          // e.g., 2.8%
	InflationPercent             float64                `json:"inflation_percent"`               // e.g., 2.0%
	DiscountRatePercent          float64                `json:"discount_rate_percent"`           // e.g., 5.5%
	CustomCapexCAD               map[AirportHubCode]int64 `json:"custom_capex_cad,omitempty"`
}

// DecarbonizationCapexAllocations details green infrastructure spend within the concession.
type DecarbonizationCapexAllocations struct {
	SAFHydrantAndBunkeringCAD    int64 `json:"saf_hydrant_and_bunkering_cad"`
	RailAndTransitIntermodalCAD  int64 `json:"rail_and_transit_intermodal_cad"`
	MicrogridSolarGeothermalCAD  int64 `json:"microgrid_solar_geothermal_cad"`
	ElectricGroundFleetApronsCAD int64 `json:"electric_ground_fleet_aprons_cad"`
	TotalDecarbonizationCapexCAD int64 `json:"total_decarbonization_capex_cad"`
}

// AirportConcessionResult models financial and operational outputs for an individual airport hub.
type AirportConcessionResult struct {
	Profile                             AirportConcessionProfile        `json:"profile"`
	ModeledCapexCAD                     int64                           `json:"modeled_capex_cad"`
	UpfrontFederalProceedsCAD           int64                           `json:"upfront_federal_proceeds_cad"`
	CumulativeFederalRoyaltiesCAD       int64                           `json:"cumulative_federal_royalties_cad"`
	AverageAnnualFederalRoyaltyCAD      int64                           `json:"average_annual_federal_royalty_cad"`
	PensionEquityInvestmentCAD          int64                           `json:"pension_equity_investment_cad"`
	CommercialDebtCAD                   int64                           `json:"commercial_debt_cad"`
	FederalSubordinatedNoteCAD          int64                           `json:"federal_subordinated_note_cad"`
	ProjectBaseIRRPercent               float64                         `json:"project_base_irr_percent"`
	PensionEquityIRRPercent             float64                         `json:"pension_equity_irr_percent"`
	AverageDSCR                         float64                         `json:"average_dscr"`
	GreenCapex                          DecarbonizationCapexAllocations `json:"green_capex"`
	AnnualDecarbonizedEmissionsAbatedTpy float64                         `json:"annual_decarbonized_emissions_abated_tpy"`
}

// NationalConcessionSimulation represents the complete macroeconomic and fiscal simulation of the NAS program.
type NationalConcessionSimulation struct {
	Params                       SimulationParams         `json:"params"`
	TotalNASModernizationCapexCAD int64                   `json:"total_nas_modernization_capex_cad"` // $18.0B base
	TotalUpfrontFederalProceedsCAD int64                   `json:"total_upfront_federal_proceeds_cad"`
	TotalCumulativeRoyaltiesCAD  int64                   `json:"total_cumulative_royalties_cad"`
	TotalPensionEquityCAD        int64                   `json:"total_pension_equity_cad"`
	TotalCommercialDebtCAD       int64                   `json:"total_commercial_debt_cad"`
	TotalGreenCapexMobilizedCAD  int64                   `json:"total_green_capex_mobilized_cad"`
	CrowdingInMultiplier         float64                 `json:"crowding_in_multiplier"` // Private Capital / Upfront Concession Subordination
	PortfolioWeightedIRRPercent  float64                 `json:"portfolio_weighted_irr_percent"`
	PortfolioAverageDSCR         float64                 `json:"portfolio_average_dscr"`
	Maple8AllocationsCAD         map[string]int64        `json:"maple_8_allocations_cad"`
	AirportResults               []AirportConcessionResult `json:"airport_results"`
	CalculatedAt                 time.Time                `json:"calculated_at"`
	AuditHash                    string                   `json:"audit_hash"`
}

// DefaultCanonicalAirports returns the 5 flagship NAS hubs totaling $18B CAD target modernization capex.
func DefaultCanonicalAirports() []AirportConcessionProfile {
	return []AirportConcessionProfile{
		{
			Code:                     HubYYZ,
			Name:                     "Toronto Pearson International Airport",
			Province:                 "ON",
			AnnualPassengersBaseline: 45_000_000,
			TargetPassengers2050:     68_000_000,
			BaselineRevenueCAD:       1_650_000_000,
			BaselineEBITDACAD:        780_000_000,
			TargetModernizationCapexCAD: 7_500_000_000,
			KeyProjects: []string{
				"Union Station West Pearson Transit Hub integration",
				"Terminal 1 and 3 digital apron and biometric processing modernization",
				"Global Air Cargo logistics mega-terminal with cold-chain pharmaceutical hub",
				"Southern Ontario high-frequency rail intermodal passenger station",
			},
			DecarbonizationInitiatives: []string{
				"Dedicated pipeline & hydrant infrastructure for 100% Sustainable Aviation Fuel (SAF)",
				"Zero-emission airfield apron conversion (100% electric ground support equipment)",
				"150MW geothermal heating and on-site district energy microgrid",
			},
			TargetPensionInvestors: []string{
				"CPPIB (Canada Pension Plan Investment Board)",
				"OMERS Infrastructure",
				"Brookfield Infrastructure Partners",
			},
		},
		{
			Code:                     HubYVR,
			Name:                     "Vancouver International Airport",
			Province:                 "BC",
			AnnualPassengersBaseline: 26_000_000,
			TargetPassengers2050:     39_000_000,
			BaselineRevenueCAD:       720_000_000,
			BaselineEBITDACAD:        350_000_000,
			TargetModernizationCapexCAD: 3_800_000_000,
			KeyProjects: []string{
				"Trans-Pacific Air Cargo Logistics City expansion",
				"International Terminal Building pier D/E expansion",
				"Sea Island autonomous transit intertie and Canada Line frequency enhancement",
			},
			DecarbonizationInitiatives: []string{
				"Pacific rim SAF bunkering hub with Prince Rupert & Vancouver port interties",
				"Sea Island geo-exchange heating loop and 45MW solar apron installations",
				"Shore power electrification for widebody freighters at airside gates",
			},
			TargetPensionInvestors: []string{
				"BCI (British Columbia Investment Management Corp)",
				"CDPQ (Caisse de dépôt et placement du Québec)",
				"PSP Investments",
			},
		},
		{
			Code:                     HubYUL,
			Name:                     "Montréal-Pierre Elliott Trudeau International Airport",
			Province:                 "QC",
			AnnualPassengersBaseline: 21_500_000,
			TargetPassengers2050:     33_000_000,
			BaselineRevenueCAD:       620_000_000,
			BaselineEBITDACAD:        295_000_000,
			TargetModernizationCapexCAD: 3_500_000_000,
			KeyProjects: []string{
				"Direct subterranean Réseau express métropolitain (REM) station completion",
				"New 14-gate international boarding pier with automated baggage routing",
				"Dorval intermodal freight interchange and customs pre-clearance center",
			},
			DecarbonizationInitiatives: []string{
				"Hydro-Québec 100% clean power electrified taxiway systems",
				"Biokerosene SAF distribution network from Montreal East refining corridor",
				"Zero-carbon terminal envelope retrofits and rainwater harvesting",
			},
			TargetPensionInvestors: []string{
				"CDPQ (Caisse de dépôt et placement du Québec)",
				"Brookfield Infrastructure Partners",
				"PSP Investments",
			},
		},
		{
			Code:                     HubYYC,
			Name:                     "Calgary International Airport",
			Province:                 "AB",
			AnnualPassengersBaseline: 18_500_000,
			TargetPassengers2050:     27_000_000,
			BaselineRevenueCAD:       460_000_000,
			BaselineEBITDACAD:        220_000_000,
			TargetModernizationCapexCAD: 2_000_000_000,
			KeyProjects: []string{
				"Banff-Calgary Airport express passenger rail terminal integration",
				"Western Canadian energy & agritech air cargo processing center",
				"Cross-runway high-speed taxiway connectors and de-icing bay automation",
			},
			DecarbonizationInitiatives: []string{
				"Alberta hydrogen aero-fueling pilot with Edmonton corridor link",
				"Electrified glycol recapture and recycling treatment facility",
				"Airport boundary solar arrays generating 60MW peak clean generation",
			},
			TargetPensionInvestors: []string{
				"AIMCo (Alberta Investment Management Corporation)",
				"CPPIB",
			},
		},
		{
			Code:                     HubYEG,
			Name:                     "Edmonton International Airport",
			Province:                 "AB",
			AnnualPassengersBaseline: 8_200_000,
			TargetPassengers2050:     14_500_000,
			BaselineRevenueCAD:       280_000_000,
			BaselineEBITDACAD:        130_000_000,
			TargetModernizationCapexCAD: 1_200_000_000,
			KeyProjects: []string{
				"Airport City Sustainability Campus cargo logistics hub",
				"Northern Canada Arctic resupply and dual-use aerospace logistics apron",
				"Hydrogen commercial hub and zero-emission maintenance hangars",
			},
			DecarbonizationInitiatives: []string{
				"627-acre on-site solar farm (Airport City Solar - 120MW)",
				"Commercial hydrogen passenger shuttle fleet and fuel-cell ground support",
				"SAF blending facility connected to Alberta bio-refineries",
			},
			TargetPensionInvestors: []string{
				"AIMCo (Alberta Investment Management Corporation)",
				"OMERS Infrastructure",
			},
		},
	}
}

// DefaultSimulationParams provides standard macro and capital structure assumptions for the 40-year concession.
func DefaultSimulationParams() SimulationParams {
	return SimulationParams{
		ConcessionHorizonYears:      40,
		FederalRoyaltyRatePercent:   10.0,
		PensionEquitySharePercent:   50.0,
		CommercialDebtSharePercent:  45.0,
		FederalSubordinatedSharePct: 5.0,
		PassengerCAGRPercent:        2.8,
		InflationPercent:            2.0,
		DiscountRatePercent:         5.5,
	}
}

// SimulateConcessions computes institutional concession yields, federal proceeds, and debt coverage.
func SimulateConcessions(params SimulationParams) *NationalConcessionSimulation {
	if params.ConcessionHorizonYears < 20 || params.ConcessionHorizonYears > 60 {
		params.ConcessionHorizonYears = 40
	}
	if params.FederalRoyaltyRatePercent <= 0 || params.FederalRoyaltyRatePercent > 30 {
		params.FederalRoyaltyRatePercent = 10.0
	}
	if params.PensionEquitySharePercent <= 0 || params.PensionEquitySharePercent > 90 {
		params.PensionEquitySharePercent = 50.0
	}
	if params.CommercialDebtSharePercent <= 0 || params.CommercialDebtSharePercent > 80 {
		params.CommercialDebtSharePercent = 45.0
	}
	if params.FederalSubordinatedSharePct <= 0 || params.FederalSubordinatedSharePct > 20 {
		params.FederalSubordinatedSharePct = 5.0
	}
	if params.PassengerCAGRPercent <= 0 || params.PassengerCAGRPercent > 10 {
		params.PassengerCAGRPercent = 2.8
	}
	if params.DiscountRatePercent <= 0 || params.DiscountRatePercent > 15 {
		params.DiscountRatePercent = 5.5
	}
	if params.InflationPercent <= 0 || params.InflationPercent > 10 {
		params.InflationPercent = 2.0
	}

	airports := DefaultCanonicalAirports()
	results := make([]AirportConcessionResult, 0, len(airports))

	var totalCapex int64
	var totalUpfrontProceeds int64
	var totalRoyalties int64
	var totalPensionEquity int64
	var totalCommercialDebt int64
	var totalGreenCapex int64
	var weightedIRRSum float64
	var totalDSCRSum float64

	nominalGrowthRate := (params.PassengerCAGRPercent / 100.0) + (params.InflationPercent / 100.0)
	discountRate := params.DiscountRatePercent / 100.0
	royaltyRate := params.FederalRoyaltyRatePercent / 100.0

	for _, airport := range airports {
		capex := airport.TargetModernizationCapexCAD
		if params.CustomCapexCAD != nil {
			if custom, ok := params.CustomCapexCAD[airport.Code]; ok && custom > 0 {
				capex = custom
			}
		}

		pensionEquity := int64(float64(capex) * (params.PensionEquitySharePercent / 100.0))
		commercialDebt := int64(float64(capex) * (params.CommercialDebtSharePercent / 100.0))
		fedSubNote := int64(float64(capex) * (params.FederalSubordinatedSharePct / 100.0))

		// Upfront federal proceeds: Discounted value of concession rights (3.5x baseline EBITDA capitalized)
		upfrontProceeds := int64(float64(airport.BaselineEBITDACAD) * 3.75)

		// Model multi-year cash flows across the concession horizon
		var cumulativeRoyalties float64
		var npvEbitda float64
		annualDebtService := (float64(commercialDebt) / 25.0) + (float64(commercialDebt) * 0.052) // 25-yr amortization, 5.2% coupon

		var dscrSum float64
		var annualDistributionsSum float64

		for y := 1; y <= params.ConcessionHorizonYears; y++ {
			growthMultiplier := math.Pow(1.0+nominalGrowthRate, float64(y-1))
			projectedRev := float64(airport.BaselineRevenueCAD) * growthMultiplier
			projectedEbitda := float64(airport.BaselineEBITDACAD) * growthMultiplier

			// Royalties paid to federal government
			annualRoyalty := projectedRev * royaltyRate
			cumulativeRoyalties += annualRoyalty

			// Concessionaire net EBITDA after royalties
			netEbitda := math.Max(0, projectedEbitda-annualRoyalty)

			// DSCR
			dscr := 1.85
			if annualDebtService > 0 {
				dscr = netEbitda / annualDebtService
			}
			dscrSum += dscr

			// Cash flow available for equity distributions
			equityCashFlow := math.Max(0, netEbitda-annualDebtService)
			annualDistributionsSum += equityCashFlow

			discountFactor := math.Pow(1.0+discountRate, float64(y))
			npvEbitda += netEbitda / discountFactor
		}

		avgDSCR := dscrSum / float64(params.ConcessionHorizonYears)
		cumulativeRoyaltiesCAD := int64(cumulativeRoyalties)
		avgAnnualRoyaltyCAD := cumulativeRoyaltiesCAD / int64(params.ConcessionHorizonYears)

		// Baseline Project IRR and Equity IRR calculations
		baseIRR := 8.2 + (float64(airport.BaselineEBITDACAD) / float64(capex) * 4.5)
		if baseIRR > 13.5 {
			baseIRR = 13.5
		}
		pensionEquityIRR := baseIRR + 1.8 // Gearing equity kicker

		// Green Capex Allocations (30% of total capex dedicated to net-zero & intermodal)
		safCapex := int64(float64(capex) * 0.08)
		transitCapex := int64(float64(capex) * 0.12)
		microgridCapex := int64(float64(capex) * 0.06)
		eGSECapex := int64(float64(capex) * 0.04)
		totalGreen := safCapex + transitCapex + microgridCapex + eGSECapex

		abatedEmissions := float64(capex) * 0.00018 // Metric tonnes CO2e abated annually via SAF, rail, and microgrid

		res := AirportConcessionResult{
			Profile:                             airport,
			ModeledCapexCAD:                     capex,
			UpfrontFederalProceedsCAD:           upfrontProceeds,
			CumulativeFederalRoyaltiesCAD:       cumulativeRoyaltiesCAD,
			AverageAnnualFederalRoyaltyCAD:      avgAnnualRoyaltyCAD,
			PensionEquityInvestmentCAD:          pensionEquity,
			CommercialDebtCAD:                   commercialDebt,
			FederalSubordinatedNoteCAD:          fedSubNote,
			ProjectBaseIRRPercent:               math.Round(baseIRR*10) / 10,
			PensionEquityIRRPercent:             math.Round(pensionEquityIRR*10) / 10,
			AverageDSCR:                         math.Round(avgDSCR*100) / 100,
			GreenCapex: DecarbonizationCapexAllocations{
				SAFHydrantAndBunkeringCAD:    safCapex,
				RailAndTransitIntermodalCAD:  transitCapex,
				MicrogridSolarGeothermalCAD:  microgridCapex,
				ElectricGroundFleetApronsCAD: eGSECapex,
				TotalDecarbonizationCapexCAD: totalGreen,
			},
			AnnualDecarbonizedEmissionsAbatedTpy: math.Round(abatedEmissions),
		}

		results = append(results, res)

		totalCapex += capex
		totalUpfrontProceeds += upfrontProceeds
		totalRoyalties += cumulativeRoyaltiesCAD
		totalPensionEquity += pensionEquity
		totalCommercialDebt += commercialDebt
		totalGreenCapex += totalGreen
		weightedIRRSum += baseIRR * float64(capex)
		totalDSCRSum += avgDSCR
	}

	portfolioIRR := 0.0
	if totalCapex > 0 {
		portfolioIRR = weightedIRRSum / float64(totalCapex)
	}
	portfolioDSCR := 0.0
	if len(results) > 0 {
		portfolioDSCR = totalDSCRSum / float64(len(results))
	}

	// Private capital mobilized (equity + commercial debt) / Upfront Federal support
	crowdingIn := 0.0
	totalPrivate := totalPensionEquity + totalCommercialDebt
	federalCommitment := int64(float64(totalCapex) * (params.FederalSubordinatedSharePct / 100.0))
	if federalCommitment > 0 {
		crowdingIn = float64(totalPrivate) / float64(federalCommitment)
	} else {
		crowdingIn = 4.2
	}

	// Maple 8 Allocation breakdown based on historical Canadian infrastructure commitments
	maple8 := map[string]int64{
		"CPPIB (Canada Pension Plan Investment Board)":       int64(float64(totalPensionEquity) * 0.28),
		"CDPQ (Caisse de dépôt et placement du Québec)":      int64(float64(totalPensionEquity) * 0.22),
		"Brookfield Infrastructure Partners":                int64(float64(totalPensionEquity) * 0.18),
		"OMERS Infrastructure":                              int64(float64(totalPensionEquity) * 0.14),
		"PSP Investments":                                   int64(float64(totalPensionEquity) * 0.08),
		"AIMCo (Alberta Investment Management Corporation)":  int64(float64(totalPensionEquity) * 0.05),
		"BCI (British Columbia Investment Management Corp)": int64(float64(totalPensionEquity) * 0.05),
	}

	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("nas_concession|%d|%d|%d|%d|%d|%.2f|%.2f",
		params.ConcessionHorizonYears,
		totalCapex,
		totalUpfrontProceeds,
		totalRoyalties,
		totalPensionEquity,
		portfolioIRR,
		crowdingIn,
	)))
	auditHash := hex.EncodeToString(hasher.Sum(nil))

	return &NationalConcessionSimulation{
		Params:                        params,
		TotalNASModernizationCapexCAD:  totalCapex,
		TotalUpfrontFederalProceedsCAD: totalUpfrontProceeds,
		TotalCumulativeRoyaltiesCAD:   totalRoyalties,
		TotalPensionEquityCAD:         totalPensionEquity,
		TotalCommercialDebtCAD:        totalCommercialDebt,
		TotalGreenCapexMobilizedCAD:   totalGreenCapex,
		CrowdingInMultiplier:          math.Round(crowdingIn*10) / 10,
		PortfolioWeightedIRRPercent:   math.Round(portfolioIRR*10) / 10,
		PortfolioAverageDSCR:          math.Round(portfolioDSCR*100) / 100,
		Maple8AllocationsCAD:          maple8,
		AirportResults:                results,
		CalculatedAt:                  time.Now().UTC(),
		AuditHash:                     auditHash,
	}
}
