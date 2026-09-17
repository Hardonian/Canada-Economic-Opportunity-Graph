package ontology

import (
	"fmt"
	"sync"
	"time"
)

// CascadeInvalidator monitors actions and detects cascading graph invalidations.
type CascadeInvalidator struct {
	mu     sync.RWMutex
	rules  []InvalidationRule
	alerts []*InvalidationAlert
}

// NewCascadeInvalidator initializes the invalidator with canonical sovereign risk rules.
func NewCascadeInvalidator() *CascadeInvalidator {
	ci := &CascadeInvalidator{
		rules:  make([]InvalidationRule, 0),
		alerts: make([]*InvalidationAlert, 0),
	}
	ci.registerDefaultRules()
	return ci
}

func (ci *CascadeInvalidator) registerDefaultRules() {
	ci.rules = append(ci.rules,
		InvalidationRule{
			RuleID:         "RULE-MILESTONE-SLIP-OFFTAKE",
			Name:           "Milestone Slip Invalidates Commercial PPA",
			TriggerAction:  ActionSetMilestoneDelay,
			ConditionField: "delay_months",
			Severity:       SeverityCritical,
			ImpactReason:   "Critical-path schedule slip exceeding 6 months triggers commercial PPA liquidated damages clause.",
		},
		InvalidationRule{
			RuleID:         "RULE-SEC35-FINANCING-RISK",
			Name:           "Contested Section 35 Consultations Hold FID",
			TriggerAction:  ActionTriggerSection35Consultation,
			ConditionField: "contested",
			Severity:       SeverityHigh,
			ImpactReason:   "Unresolved Section 35 consultation challenges freeze senior syndicate debt drawdown.",
		},
		InvalidationRule{
			RuleID:         "RULE-CAPITAL-EQUITY-ILGP",
			Name:           "Indigenous Equity Tranche Below Mandate",
			TriggerAction:  ActionAdjustCapitalTranche,
			ConditionField: "indigenous_equity_pct",
			Severity:       SeverityMedium,
			ImpactReason:   "First Nations co-investment equity tranche dropped below minimum 10% ILGP qualification threshold.",
		},
	)
}

// RegisterRule appends a custom invalidation rule.
func (ci *CascadeInvalidator) RegisterRule(rule InvalidationRule) {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.rules = append(ci.rules, rule)
}

// EvaluateCascades inspects an action against registered invalidation rules.
func (ci *CascadeInvalidator) EvaluateCascades(action *Action) []*InvalidationAlert {
	ci.mu.Lock()
	defer ci.mu.Unlock()

	var newAlerts []*InvalidationAlert
	for _, rule := range ci.rules {
		if rule.TriggerAction != action.Type {
			continue
		}

		triggered := false
		reason := rule.ImpactReason

		switch rule.ConditionField {
		case "delay_months":
			if val, ok := action.Parameters["delay_months"]; ok {
				if fval, ok := val.(float64); ok && fval >= 6.0 {
					triggered = true
					reason = fmt.Sprintf("%s (Observed slip: %.1f months)", rule.ImpactReason, fval)
				}
			}
		case "contested":
			if val, ok := action.Parameters["contested"]; ok {
				if bval, ok := val.(bool); ok && bval {
					triggered = true
				}
			}
		case "indigenous_equity_pct":
			if val, ok := action.Parameters["indigenous_equity_pct"]; ok {
				if fval, ok := val.(float64); ok && fval < 10.0 {
					triggered = true
					reason = fmt.Sprintf("%s (Observed: %.1f%%, Required: >=10%%)", rule.ImpactReason, fval)
				}
			}
		default:
			if _, exists := action.Parameters[rule.ConditionField]; exists {
				triggered = true
			}
		}

		if triggered {
			alert := &InvalidationAlert{
				AlertID:        fmt.Sprintf("ALERT-%d-%s", time.Now().UnixNano(), rule.RuleID),
				RuleID:         rule.RuleID,
				ActionID:       action.ID,
				SourceNodeID:   action.TargetNodeID,
				AffectedNodeID: action.TargetNodeID, // Or connected dependent node
				Severity:       rule.Severity,
				Reason:         reason,
				DetectedAt:     time.Now().UTC(),
			}
			ci.alerts = append(ci.alerts, alert)
			newAlerts = append(newAlerts, alert)
		}
	}

	return newAlerts
}

// GetActiveAlerts returns all recorded cascading invalidation alerts.
func (ci *CascadeInvalidator) GetActiveAlerts() []*InvalidationAlert {
	ci.mu.RLock()
	defer ci.mu.RUnlock()

	result := make([]*InvalidationAlert, len(ci.alerts))
	copy(result, ci.alerts)
	return result
}
