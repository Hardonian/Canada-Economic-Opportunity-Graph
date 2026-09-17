package cegs

import (
	"math"
)

// TRC92Scorecard evaluates corporate reconciliation compliance under TRC Call to Action #92.
type TRC92Scorecard struct {
	ProjectID                  string  `json:"project_id"`
	Pillar1ConsentAndEquity    float64 `json:"pillar1_consent_and_equity"`    // 0.0 - 100.0 (FPIC, Equity Co-ownership)
	Pillar2JobsAndProcurement  float64 `json:"pillar2_jobs_and_procurement"`  // 0.0 - 100.0 (Local hiring, contracting)
	Pillar3EducationAndUNDRIP  float64 `json:"pillar3_education_and_undrip"`  // 0.0 - 100.0 (Staff training on intercultural competency)
	OverallReconciliationScore float64 `json:"overall_reconciliation_score"` // 0.0 - 100.0
	TRCRating                  string  `json:"trc_rating"`                   // "EXEMPLARY_PARTNERSHIP", "COMMITTED", "NEEDS_IMPROVEMENT"
	Summary                    string  `json:"summary"`
}

// TRC92Evaluator audits corporate practices against Truth and Reconciliation Commission Call to Action 92.
type TRC92Evaluator struct{}

// NewTRC92Evaluator creates a TRC 92 evaluator.
func NewTRC92Evaluator() *TRC92Evaluator {
	return &TRC92Evaluator{}
}

// Evaluate computes the 3-pillar TRC Call to Action 92 business reconciliation scorecard.
func (e *TRC92Evaluator) Evaluate(
	projectID string,
	equityOwnershipPct float64, // Target: >= 25%
	indigenousProcurementPct float64, // Target: >= 5% federal minimum (10% preferred)
	staffTrainedOnUNDRIPPct float64, // Target: 100% executive and site staff
) *TRC92Scorecard {
	// Pillar 1: FPIC & Equity (Weight: 40%)
	p1 := math.Min(100.0, (equityOwnershipPct/25.0)*100.0)

	// Pillar 2: Jobs & Procurement (Weight: 35%)
	p2 := math.Min(100.0, (indigenousProcurementPct/10.0)*100.0)

	// Pillar 3: Education & Intercultural Competency (Weight: 25%)
	p3 := math.Min(100.0, staffTrainedOnUNDRIPPct)

	overall := math.Round((p1*0.40 + p2*0.35 + p3*0.25)*10) / 10

	rating := "NEEDS_IMPROVEMENT"
	summary := "First Nations co-investment and workforce participation below recommended reconciliation thresholds."

	if overall >= 85.0 {
		rating = "EXEMPLARY_PARTNERSHIP"
		summary = "Model project implementing genuine economic self-determination and intergenerational wealth creation."
	} else if overall >= 60.0 {
		rating = "COMMITTED"
		summary = "Demonstrated commitment to Indigenous economic participation with active agreements in place."
	}

	return &TRC92Scorecard{
		ProjectID:                  projectID,
		Pillar1ConsentAndEquity:    math.Round(p1*10) / 10,
		Pillar2JobsAndProcurement:  math.Round(p2*10) / 10,
		Pillar3EducationAndUNDRIP:  math.Round(p3*10) / 10,
		OverallReconciliationScore: overall,
		TRCRating:                  rating,
		Summary:                    summary,
	}
}
