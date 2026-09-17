package cegs

import (
	"testing"
)

func TestC59Auditor(t *testing.T) {
	auditor := NewC59Auditor()

	// Compliant project: 12% apprentice hours, certified wages
	res1 := auditor.AuditProject("proj-1", 1000000000.0, 100000.0, 12000.0, true, 30.0)
	if res1.AuditStatus != "COMPLIANT_MAX_CREDIT" {
		t.Errorf("expected COMPLIANT_MAX_CREDIT, got %s", res1.AuditStatus)
	}
	if res1.EligibleITCRatePct != 30.0 {
		t.Errorf("expected 30%% ITC rate, got %f", res1.EligibleITCRatePct)
	}
	if res1.EstimatedITCMonetizedCAD != 300000000.0 {
		t.Errorf("expected $300M monetized ITC, got %f", res1.EstimatedITCMonetizedCAD)
	}

	// Non-compliant project: 5% apprentice hours (violates >=10% mandate)
	res2 := auditor.AuditProject("proj-2", 1000000000.0, 100000.0, 5000.0, true, 30.0)
	if res2.AuditStatus != "PENALIZED_NON_COMPLIANT" {
		t.Errorf("expected PENALIZED_NON_COMPLIANT, got %s", res2.AuditStatus)
	}
	if res2.EligibleITCRatePct != 20.0 {
		t.Errorf("expected 20%% penalized rate, got %f", res2.EligibleITCRatePct)
	}
}

func TestC69ClockAuditor(t *testing.T) {
	auditor := NewC69ClockAuditor()

	// On track planning phase (120 of 180 days)
	s1 := auditor.AuditClock("proj-1", PhaseEarlyPlanning, 120, false)
	if s1.IsStatutoryBreach {
		t.Errorf("expected no breach for 120/180 days")
	}
	if s1.RemainingDays != 60 {
		t.Errorf("expected 60 remaining days, got %d", s1.RemainingDays)
	}

	// Breached impact statement phase (340 of 300 days)
	s2 := auditor.AuditClock("proj-2", PhaseImpactStatement, 340, false)
	if !s2.IsStatutoryBreach {
		t.Errorf("expected statutory breach for 340/300 days")
	}
	if s2.LitigationRisk != "STATUTORY_MANDAMUS_EXPOSURE" {
		t.Errorf("expected mandamus exposure, got %s", s2.LitigationRisk)
	}
}

func TestTRC92Evaluator(t *testing.T) {
	eval := NewTRC92Evaluator()

	// Exemplary partnership: 25% equity, 12% procurement, 100% staff trained
	sc := eval.Evaluate("proj-1", 25.0, 12.0, 100.0)
	if sc.TRCRating != "EXEMPLARY_PARTNERSHIP" {
		t.Errorf("expected EXEMPLARY_PARTNERSHIP, got %s", sc.TRCRating)
	}
	if sc.OverallReconciliationScore < 95.0 {
		t.Errorf("expected score >= 95.0, got %f", sc.OverallReconciliationScore)
	}
}
