package cegs

import (
	"fmt"
)

// C69Phase defines formal statutory phases under the Impact Assessment Act.
type C69Phase string

const (
	PhaseEarlyPlanning       C69Phase = "EARLY_PLANNING"       // 180-day statutory limit
	PhaseImpactStatement     C69Phase = "IMPACT_STATEMENT"     // 300-day statutory limit
	PhaseMinisterialDecision C69Phase = "MINISTERIAL_DECISION" // 30-day statutory limit
)

// C69ClockStatus tracks regulatory review duration against statutory caps.
type C69ClockStatus struct {
	ProjectID         string   `json:"project_id"`
	Phase             C69Phase `json:"phase"`
	ElapsedDays       int      `json:"elapsed_days"`
	StatutoryCapDays  int      `json:"statutory_cap_days"`
	RemainingDays     int      `json:"remaining_days"`
	IsClockSuspended  bool     `json:"is_clock_suspended"`
	IsStatutoryBreach bool     `json:"is_statutory_breach"`
	LitigationRisk    string   `json:"litigation_risk"` // "LOW", "ELEVATED", "STATUTORY_MANDAMUS_EXPOSURE"
	StatusReport      string   `json:"status_report"`
}

// C69ClockAuditor monitors statutory deadlines under Bill C-69.
type C69ClockAuditor struct{}

// NewC69ClockAuditor creates an Impact Assessment Act clock auditor.
func NewC69ClockAuditor() *C69ClockAuditor {
	return &C69ClockAuditor{}
}

// AuditClock evaluates whether a federal impact assessment is on track or exceeding statutory limits.
func (ca *C69ClockAuditor) AuditClock(projectID string, phase C69Phase, elapsedDays int, suspended bool) *C69ClockStatus {
	capDays := 180
	switch phase {
	case PhaseImpactStatement:
		capDays = 300
	case PhaseMinisterialDecision:
		capDays = 30
	}

	remaining := capDays - elapsedDays
	breach := remaining < 0
	risk := "LOW"
	report := fmt.Sprintf("Phase %s on track: %d of %d statutory days elapsed (%d remaining).", phase, elapsedDays, capDays, remaining)

	if breach {
		risk = "STATUTORY_MANDAMUS_EXPOSURE"
		report = fmt.Sprintf("Statutory limit breached: %d days elapsed vs %d statutory cap. Proponent entitled to judicial review / order of mandamus.", elapsedDays, capDays)
	} else if remaining <= 30 {
		risk = "ELEVATED"
		report = fmt.Sprintf("Statutory deadline approaching: only %d days remaining to issue Impact Assessment decision statement.", remaining)
	}

	if suspended {
		report += " Note: IAAC statutory clock is currently paused under Section 18 information request suspension."
	}

	return &C69ClockStatus{
		ProjectID:         projectID,
		Phase:             phase,
		ElapsedDays:       elapsedDays,
		StatutoryCapDays:  capDays,
		RemainingDays:     remaining,
		IsClockSuspended:  suspended,
		IsStatutoryBreach: breach,
		LitigationRisk:    risk,
		StatusReport:      report,
	}
}
