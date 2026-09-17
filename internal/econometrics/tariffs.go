package econometrics

import (
	"math"
)

// TariffShockRequest defines parameters for a cross-border trade tariff simulation.
type TariffShockRequest struct {
	TariffRatePct   float64  `json:"tariff_rate_pct"` // e.g. 10.0 or 25.0
	AffectedSectors []string `json:"affected_sectors"` // "STEEL", "ALUMINUM", "ENERGY", "CRITICAL_MINERALS", "AUTOMOTIVE"
}

// TariffShockResult summarizes the national macroeconomic and trade damage from a trade shock.
type TariffShockResult struct {
	TariffRatePct          float64  `json:"tariff_rate_pct"`
	AnnualExportLossCADM   float64  `json:"annual_export_loss_cad_m"`
	AnnualGDPLossCADM      float64  `json:"annual_gdp_loss_cad_m"`
	DirectJobsAtRisk       int      `json:"direct_jobs_at_risk"`
	CADDevaluationPressure float64  `json:"cad_devaluation_pressure_pct"` // e.g. -3.5%
	WorstHitProvinces      []string `json:"worst_hit_provinces"`
	MitigationStrategies   []string `json:"mitigation_strategies"`
}

// TariffShockSimulator models cross-border trade friction under USMCA disputes or tariffs.
type TariffShockSimulator struct{}

// NewTariffShockSimulator creates a tariff shock simulator.
func NewTariffShockSimulator() *TariffShockSimulator {
	return &TariffShockSimulator{}
}

// Simulate evaluates the national economic exposure to U.S. border tariff shocks.
func (tss *TariffShockSimulator) Simulate(req TariffShockRequest) *TariffShockResult {
	baseBilateralExportCADB := 580.0 // ~$580B CAD annual exports to US
	sectorMultiplier := float64(len(req.AffectedSectors)) * 0.18

	// Price elasticity of export demand ~ -0.75 for commodities
	exportDropPct := (req.TariffRatePct * 0.75) * sectorMultiplier
	exportLossB := baseBilateralExportCADB * (exportDropPct / 100.0)
	exportLossM := exportLossB * 1000.0

	// Leontief multiplier for manufacturing/energy is ~1.85x
	gdpLossM := exportLossM * 1.85
	jobsAtRisk := int(gdpLossM * 4.2) // ~4.2 direct/indirect jobs per $1M CAD output lost
	fxPressure := -math.Min(12.0, req.TariffRatePct*0.28)

	var provinces []string
	for _, s := range req.AffectedSectors {
		switch s {
		case "STEEL", "AUTOMOTIVE":
			provinces = append(provinces, "ON")
		case "ALUMINUM":
			provinces = append(provinces, "QC")
		case "ENERGY":
			provinces = append(provinces, "AB", "SK")
		case "CRITICAL_MINERALS":
			provinces = append(provinces, "ON", "QC", "NT", "SK")
		}
	}

	mitigations := []string{
		"Accelerate Accelerated Capital Cost Allowance (ACCA) for domestic downstream processing",
		"Direct bilateral export redirection to European Union (CETA) and Indo-Pacific (CPTPP) markets",
		"Invoke Defense Production Act (DPA) Title III co-investment exemptions with Pentagon for critical minerals",
		"Deploy Canadian Commercial Corporation (CCC) sovereign export guarantees to offset bonding costs",
	}

	return &TariffShockResult{
		TariffRatePct:          req.TariffRatePct,
		AnnualExportLossCADM:   math.Round(exportLossM),
		AnnualGDPLossCADM:      math.Round(gdpLossM),
		DirectJobsAtRisk:       jobsAtRisk,
		CADDevaluationPressure: math.Round(fxPressure*10) / 10,
		WorstHitProvinces:      dedupStrings(provinces),
		MitigationStrategies:   mitigations,
	}
}

func dedupStrings(list []string) []string {
	seen := make(map[string]bool)
	var res []string
	for _, s := range list {
		if !seen[s] && s != "" {
			seen[s] = true
			res = append(res, s)
		}
	}
	return res
}
