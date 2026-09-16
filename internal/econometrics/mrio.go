package econometrics

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// MacroMultipliers defines Input-Output coefficients for a specific infrastructure sector.
type MacroMultipliers struct {
	DirectGDPPerCAD   float64 // Direct GDP impact per $1 CAD of CAPEX
	IndirectGDPPerCAD float64 // Supply-chain suppliers GDP per $1 CAD
	InducedGDPPerCAD  float64 // Household spending from wages per $1 CAD
	JobsPerMillionCAD float64 // Full-Time Equivalent (FTE) Person-Years per $1M CAD
	FederalTaxRate    float64 // Federal tax yield ratio (corporate + personal + GST)
	ProvincialTaxRate float64 // Provincial tax yield ratio (PST/HST + corporate + royalties)
	MunicipalTaxRate  float64 // Municipal tax yield ratio (property taxes, permitting fees)
}

// DefaultSectorMultipliers returns Statistics Canada SUT-calibrated coefficients.
func DefaultSectorMultipliers(sector domain.Sector) MacroMultipliers {
	switch sector {
	case domain.SectorCriticalMinerals, domain.SectorMiningMetals:
		return MacroMultipliers{
			DirectGDPPerCAD:   0.65,
			IndirectGDPPerCAD: 0.48,
			InducedGDPPerCAD:  0.36,
			JobsPerMillionCAD: 6.8,
			FederalTaxRate:    0.14,
			ProvincialTaxRate: 0.12,
			MunicipalTaxRate:  0.02,
		}
	case domain.SectorNuclearEnergy:
		return MacroMultipliers{
			DirectGDPPerCAD:   0.72,
			IndirectGDPPerCAD: 0.58,
			InducedGDPPerCAD:  0.42,
			JobsPerMillionCAD: 7.5,
			FederalTaxRate:    0.16,
			ProvincialTaxRate: 0.13,
			MunicipalTaxRate:  0.025,
		}
	case domain.SectorCleanEnergy:
		return MacroMultipliers{
			DirectGDPPerCAD:   0.58,
			IndirectGDPPerCAD: 0.44,
			InducedGDPPerCAD:  0.31,
			JobsPerMillionCAD: 5.9,
			FederalTaxRate:    0.13,
			ProvincialTaxRate: 0.10,
			MunicipalTaxRate:  0.02,
		}
	case domain.SectorAICompute:
		return MacroMultipliers{
			DirectGDPPerCAD:   0.50,
			IndirectGDPPerCAD: 0.42,
			InducedGDPPerCAD:  0.28,
			JobsPerMillionCAD: 4.8,
			FederalTaxRate:    0.15,
			ProvincialTaxRate: 0.11,
			MunicipalTaxRate:  0.03,
		}
	case domain.SectorTransportation:
		return MacroMultipliers{
			DirectGDPPerCAD:   0.62,
			IndirectGDPPerCAD: 0.46,
			InducedGDPPerCAD:  0.34,
			JobsPerMillionCAD: 6.4,
			FederalTaxRate:    0.14,
			ProvincialTaxRate: 0.11,
			MunicipalTaxRate:  0.02,
		}
	default:
		return MacroMultipliers{
			DirectGDPPerCAD:   0.60,
			IndirectGDPPerCAD: 0.45,
			InducedGDPPerCAD:  0.32,
			JobsPerMillionCAD: 6.0,
			FederalTaxRate:    0.14,
			ProvincialTaxRate: 0.11,
			MunicipalTaxRate:  0.02,
		}
	}
}

// MRIOImpact summarizes the macro-econometric value created by a project.
type MRIOImpact struct {
	ProjectID         string    `json:"project_id"`
	ProjectName       string    `json:"project_name"`
	CapexCAD          int64     `json:"capex_cad"`
	DirectGDPCAD      int64     `json:"direct_gdp_cad"`
	IndirectGDPCAD    int64     `json:"indirect_gdp_cad"`
	InducedGDPCAD     int64     `json:"induced_gdp_cad"`
	TotalGDPCAD       int64     `json:"total_gdp_cad"`
	TotalMultipler    float64   `json:"total_multiplier"`
	PersonYearsJobs   int64     `json:"person_years_jobs"`
	FederalTaxCAD     int64     `json:"federal_tax_cad"`
	ProvincialTaxCAD  int64     `json:"provincial_tax_cad"`
	MunicipalTaxCAD   int64     `json:"municipal_tax_cad"`
	TotalFiscalReturn int64     `json:"total_fiscal_return_cad"`
	ModelVersion      string    `json:"model_version"`
	AuditHash         string    `json:"audit_hash"`
	CalculatedAt      time.Time `json:"calculated_at"`
}

// SensitivityShockRequest specifies macroeconomic and commodity shocks.
type SensitivityShockRequest struct {
	InterestRateDeltaBps float64 `json:"interest_rate_delta_bps"` // e.g. +150 bps
	FXDepreciationCAD    float64 `json:"fx_depreciation_cad"`    // e.g. -0.05 (CAD weakens vs USD)
	CommodityPriceDelta  float64 `json:"commodity_price_delta"`  // e.g. -0.20 (20% drop in commodity price)
	ImportedCapexShare   float64 `json:"imported_capex_share"`   // Share of CAPEX reliant on foreign equipment
}

// SensitivityShockResult details the stress-tested impact on CAPEX and debt service.
type SensitivityShockResult struct {
	OriginalCapexCAD      int64   `json:"original_capex_cad"`
	AdjustedCapexCAD      int64   `json:"adjusted_capex_cad"`
	CapexDeltaCAD         int64   `json:"capex_delta_cad"`
	AnnualDebtServiceCAD  int64   `json:"annual_debt_service_cad"`
	DebtServiceDeltaCAD   int64   `json:"debt_service_delta_cad"`
	ProjectViabilityShift float64 `json:"project_viability_shift"` // Negative indicates heightened default/stall risk
	StressSummary         string  `json:"stress_summary"`
}

// Engine performs deterministic MRIO and financial sensitivity calculations.
type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

// CalculateMRIO generates the complete multi-regional input-output econometric assessment.
func (e *Engine) CalculateMRIO(project *domain.Project) *MRIOImpact {
	m := DefaultSectorMultipliers(project.Sector)
	capex := float64(project.CapexCAD)

	direct := capex * m.DirectGDPPerCAD
	indirect := capex * m.IndirectGDPPerCAD
	induced := capex * m.InducedGDPPerCAD
	total := direct + indirect + induced

	totalMult := 0.0
	if capex > 0 {
		totalMult = total / capex
	}

	jobs := int64(math.Round((capex / 1_000_000.0) * m.JobsPerMillionCAD))

	fedTax := total * m.FederalTaxRate
	provTax := total * m.ProvincialTaxRate
	munTax := total * m.MunicipalTaxRate
	totalTax := fedTax + provTax + munTax

	impact := &MRIOImpact{
		ProjectID:         project.ID,
		ProjectName:       project.Name,
		CapexCAD:          project.CapexCAD,
		DirectGDPCAD:      int64(math.Round(direct)),
		IndirectGDPCAD:    int64(math.Round(indirect)),
		InducedGDPCAD:     int64(math.Round(induced)),
		TotalGDPCAD:       int64(math.Round(total)),
		TotalMultipler:    math.Round(totalMult*100) / 100,
		PersonYearsJobs:   jobs,
		FederalTaxCAD:     int64(math.Round(fedTax)),
		ProvincialTaxCAD:  int64(math.Round(provTax)),
		MunicipalTaxCAD:   int64(math.Round(munTax)),
		TotalFiscalReturn: int64(math.Round(totalTax)),
		ModelVersion:      "statcan-sut-mrio-v1.0",
		CalculatedAt:      time.Now().UTC(),
	}

	raw := fmt.Sprintf("%s|%d|%d|%d|%s", impact.ProjectID, impact.CapexCAD, impact.TotalGDPCAD, impact.PersonYearsJobs, impact.ModelVersion)
	h := sha256.Sum256([]byte(raw))
	impact.AuditHash = hex.EncodeToString(h[:])

	return impact
}

// EvaluateShock models macro rate, FX, and commodity price stress on project economics.
func (e *Engine) EvaluateShock(project *domain.Project, shock SensitivityShockRequest) *SensitivityShockResult {
	capex := float64(project.CapexCAD)
	if capex <= 0 {
		return &SensitivityShockResult{OriginalCapexCAD: 0, AdjustedCapexCAD: 0}
	}

	// FX Inflation on imported equipment: if imported capex share is 40% and CAD depreciates by 5%, imported cost rises by ~5.26%
	importShare := shock.ImportedCapexShare
	if importShare <= 0 {
		importShare = 0.35 // Default 35% imported content for Canadian industrial assets
	}

	fxCostMultiplier := 1.0
	if shock.FXDepreciationCAD < 0 {
		fxCostMultiplier = 1.0 + (importShare * (math.Abs(shock.FXDepreciationCAD) / (1.0 + shock.FXDepreciationCAD)))
	}

	adjustedCapex := capex * fxCostMultiplier
	capexDelta := adjustedCapex - capex

	// Debt Service: Assuming 60% leverage, 15-year amortizing facility
	debtPortion := adjustedCapex * 0.60
	baseRate := 0.055 // 5.5% base all-in rate
	stressedRate := baseRate + (shock.InterestRateDeltaBps / 10000.0)

	baseAnnualDebt := debtPortion * (baseRate / (1.0 - math.Pow(1.0+baseRate, -15)))
	stressedAnnualDebt := debtPortion * (stressedRate / (1.0 - math.Pow(1.0+stressedRate, -15)))
	debtServiceDelta := stressedAnnualDebt - baseAnnualDebt

	// Viability shift
	viabilityScore := 0.0
	if shock.InterestRateDeltaBps > 0 {
		viabilityScore -= (shock.InterestRateDeltaBps / 100.0) * 0.05 // -5% per 100 bps
	}
	if shock.CommodityPriceDelta < 0 {
		viabilityScore += shock.CommodityPriceDelta * 0.8 // -16% for 20% drop
	}
	if capexDelta > 0 {
		viabilityScore -= (capexDelta / capex) * 0.5
	}

	summary := fmt.Sprintf("Under shock (+%.0fbps interest, %.1f%% FX, %.1f%% commodity), CAPEX shifts by +$%.1fM CAD and annual debt service shifts by +$%.1fM CAD/yr.",
		shock.InterestRateDeltaBps, shock.FXDepreciationCAD*100, shock.CommodityPriceDelta*100, capexDelta/1e6, debtServiceDelta/1e6)

	return &SensitivityShockResult{
		OriginalCapexCAD:      project.CapexCAD,
		AdjustedCapexCAD:      int64(math.Round(adjustedCapex)),
		CapexDeltaCAD:         int64(math.Round(capexDelta)),
		AnnualDebtServiceCAD:  int64(math.Round(stressedAnnualDebt)),
		DebtServiceDeltaCAD:   int64(math.Round(debtServiceDelta)),
		ProjectViabilityShift: math.Round(viabilityScore*100) / 100,
		StressSummary:         summary,
	}
}
