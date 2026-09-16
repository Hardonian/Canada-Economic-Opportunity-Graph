package projectfinance

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// ITCType classifies federal clean economy tax credits.
type ITCType string

const (
	ITCCleanTechnology         ITCType = "CLEAN_TECHNOLOGY_ITC_30"
	ITCCleanHydrogen           ITCType = "CLEAN_HYDROGEN_ITC_40"
	ITCCleanElectricity        ITCType = "CLEAN_ELECTRICITY_ITC_15"
	ITCCleanManufacturing      ITCType = "CLEAN_TECH_MANUFACTURING_ITC_30"
	ITCCarbonCaptureStorage    ITCType = "CCUS_ITC_50"
)

// CleanTaxCreditProfile models refundable federal tax benefits and CCfD underwriting.
type CleanTaxCreditProfile struct {
	ProjectID                 string  `json:"project_id"`
	ApplicableITC             ITCType `json:"applicable_itc"`
	EligibleCapexCAD          int64   `json:"eligible_capex_cad"`
	BaseCreditRatePercent     float64 `json:"base_credit_rate_percent"`
	LaborConditionBonusPercent float64 `json:"labor_condition_bonus_percent"` // +10% for prevailing wage and apprentice ratios
	EffectiveCreditRatePercent float64 `json:"effective_credit_rate_percent"`
	TotalTaxCreditYieldCAD    int64   `json:"total_tax_credit_yield_cad"`
	CCfDEligible              bool    `json:"ccfd_eligible"`
	CCfDStrikePriceCADTonne   float64 `json:"ccfd_strike_price_cad_tonne"`
	EstimatedAnnualCCfDSubsidyCAD int64 `json:"estimated_annual_ccfd_subsidy_cad"`
	AuditHash                 string  `json:"audit_hash"`
}

// TaxCreditCalculator evaluates federal tax incentives under Bills C-59 and C-69.
type TaxCreditCalculator struct{}

// NewTaxCreditCalculator creates a tax credit analyzer.
func NewTaxCreditCalculator() *TaxCreditCalculator {
	return &TaxCreditCalculator{}
}

// CalculateCredits computes refundable tax credits and carbon contract-for-difference benefits.
func (c *TaxCreditCalculator) CalculateCredits(project *domain.Project) *CleanTaxCreditProfile {
	capex := project.CapexCAD
	if capex <= 0 {
		capex = 500_000_000
	}

	itcType := ITCCleanTechnology
	baseRate := 20.0
	laborBonus := 10.0 // Meets prevailing wage and apprenticeships
	eligibleRatio := 0.85
	isCCfD := false
	strikePrice := 0.0
	annualCCfD := int64(0)

	switch project.Sector {
	case domain.SectorNuclearEnergy:
		itcType = ITCCleanElectricity
		baseRate = 5.0
		laborBonus = 10.0 // 15% total Clean Electricity ITC
		eligibleRatio = 0.90
	case domain.SectorCleanEnergy:
		if project.Subsector == "Hydrogen" || project.Subsector == "Clean Fuels" {
			itcType = ITCCleanHydrogen
			baseRate = 30.0
			laborBonus = 10.0 // 40% total
			eligibleRatio = 0.85
			isCCfD = true
			strikePrice = 95.0 // $95/tonne CO2e guaranteed floor
			annualCCfD = 18_000_000
		} else {
			itcType = ITCCleanTechnology
			baseRate = 20.0
			laborBonus = 10.0 // 30% total
			eligibleRatio = 0.88
		}
	case domain.SectorCriticalMinerals:
		itcType = ITCCleanManufacturing
		baseRate = 20.0
		laborBonus = 10.0 // 30% Critical Minerals Extraction & Processing
		eligibleRatio = 0.80
	default:
		itcType = ITCCleanTechnology
		baseRate = 20.0
		laborBonus = 10.0
		eligibleRatio = 0.70
	}

	effectiveRate := baseRate + laborBonus
	eligibleCapex := int64(float64(capex) * eligibleRatio)
	totalCredit := int64(float64(eligibleCapex) * (effectiveRate / 100.0))

	profile := &CleanTaxCreditProfile{
		ProjectID:                  project.ID,
		ApplicableITC:              itcType,
		EligibleCapexCAD:           eligibleCapex,
		BaseCreditRatePercent:      baseRate,
		LaborConditionBonusPercent: laborBonus,
		EffectiveCreditRatePercent: effectiveRate,
		TotalTaxCreditYieldCAD:     totalCredit,
		CCfDEligible:               isCCfD,
		CCfDStrikePriceCADTonne:    strikePrice,
		EstimatedAnnualCCfDSubsidyCAD: annualCCfD,
	}

	auditData := fmt.Sprintf("%s|%s|%d|%.1f|%d|%t",
		profile.ProjectID, profile.ApplicableITC, profile.EligibleCapexCAD, profile.EffectiveCreditRatePercent, profile.TotalTaxCreditYieldCAD, profile.CCfDEligible)
	h := sha256.Sum256([]byte(auditData))
	profile.AuditHash = hex.EncodeToString(h[:])

	return profile
}
