package cegs

import (
	"fmt"
	"math"
)

// C59AuditResult holds the statutory evaluation for Clean Economy ITCs.
type C59AuditResult struct {
	ProjectID                string  `json:"project_id"`
	TotalCraftHoursWorked    float64 `json:"total_craft_hours_worked"`
	ApprenticeHoursWorked    float64 `json:"apprentice_hours_worked"`
	ApprenticeRatioPct       float64 `json:"apprentice_ratio_pct"`
	Meets10PctMandate        bool    `json:"meets_10_pct_mandate"`
	PrevailingWagesCertified bool    `json:"prevailing_wages_certified"`
	EligibleITCRatePct       float64 `json:"eligible_itc_rate_pct"` // e.g. 30.0% or 10.0% (penalized)
	EstimatedITCMonetizedCAD float64 `json:"estimated_itc_monetized_cad"`
	AuditStatus              string  `json:"audit_status"` // "COMPLIANT_MAX_CREDIT", "PENALIZED_NON_COMPLIANT"
	AuditSummary             string  `json:"audit_summary"`
}

// C59Auditor audits construction labor compliance against Bill C-59 prevailing wage and apprenticeship rules.
type C59Auditor struct{}

// NewC59Auditor creates a Bill C-59 statutory auditor.
func NewC59Auditor() *C59Auditor {
	return &C59Auditor{}
}

// AuditProject evaluates payroll records to determine if project qualifies for maximum Clean Economy ITC.
func (a *C59Auditor) AuditProject(
	projectID string,
	capexEligibleCAD float64,
	totalCraftHours float64,
	apprenticeHours float64,
	wagesCertified bool,
	baseITCRatePct float64, // Typically 30% for Clean Tech / 40% for Clean Hydrogen
) *C59AuditResult {
	if totalCraftHours <= 0 {
		totalCraftHours = 1.0
	}

	ratio := (apprenticeHours / totalCraftHours) * 100.0
	meetsApprentice := ratio >= 10.0

	// Under Bill C-59, if either prevailing wage or apprentice requirements fail, ITC drops by 10 percentage points or to reduced rate
	compliant := meetsApprentice && wagesCertified
	effectiveRate := baseITCRatePct
	status := "COMPLIANT_MAX_CREDIT"
	summary := fmt.Sprintf("Project satisfies Bill C-59 requirements: Apprentice ratio %.1f%% (>=10%% mandate) and certified prevailing union wages.", ratio)

	if !compliant {
		effectiveRate = math.Max(10.0, baseITCRatePct-10.0) // Penalized rate
		status = "PENALIZED_NON_COMPLIANT"
		summary = fmt.Sprintf("Bill C-59 non-compliance detected: Apprentice ratio %.1f%% or wage certification missing. ITC rate penalized from %.1f%% to %.1f%%.",
			ratio, baseITCRatePct, effectiveRate)
	}

	monetized := capexEligibleCAD * (effectiveRate / 100.0)

	return &C59AuditResult{
		ProjectID:                projectID,
		TotalCraftHoursWorked:    totalCraftHours,
		ApprenticeHoursWorked:    apprenticeHours,
		ApprenticeRatioPct:       math.Round(ratio*10) / 10,
		Meets10PctMandate:        meetsApprentice,
		PrevailingWagesCertified: wagesCertified,
		EligibleITCRatePct:       effectiveRate,
		EstimatedITCMonetizedCAD: math.Round(monetized),
		AuditStatus:              status,
		AuditSummary:             summary,
	}
}
