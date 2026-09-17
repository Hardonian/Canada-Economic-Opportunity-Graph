package ontology

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// ActionType defines the typed domain actions executable on the graph.
type ActionType string

const (
	ActionApproveILGPGuarantee        ActionType = "APPROVE_ILGP_GUARANTEE"
	ActionTriggerSection35Consultation ActionType = "TRIGGER_SECTION_35_CONSULTATION"
	ActionDispatchFirmPower           ActionType = "DISPATCH_FIRM_POWER"
	ActionAdjustCapitalTranche        ActionType = "ADJUST_CAPITAL_TRANCHE"
	ActionRecordEnvironmentalPermit    ActionType = "RECORD_ENVIRONMENTAL_PERMIT"
	ActionUpdateProjectStage          ActionType = "UPDATE_PROJECT_STAGE"
	ActionSetMilestoneDelay           ActionType = "SET_MILESTONE_DELAY"
)

// ActionStatus represents the lifecycle state of a write-back action.
type ActionStatus string

const (
	ActionStatusPending  ActionStatus = "PENDING"
	ActionStatusApplied  ActionStatus = "APPLIED"
	ActionStatusReverted ActionStatus = "REVERTED"
	ActionStatusFailed   ActionStatus = "FAILED"
)

// Action represents an isolated, transactional graph mutation.
type Action struct {
	ID             string                 `json:"id"`
	Type           ActionType             `json:"type"`
	TargetNodeID   string                 `json:"target_node_id"`
	TargetNodeType string                 `json:"target_node_type"`
	BranchID       string                 `json:"branch_id"`
	Initiator      string                 `json:"initiator"` // User / Agent ID
	Parameters     map[string]interface{} `json:"parameters"`
	PreviousState  map[string]interface{} `json:"previous_state,omitempty"`
	NewState       map[string]interface{} `json:"new_state,omitempty"`
	ValidTime      time.Time              `json:"valid_time"`  // Real-world effective time
	SystemTime     time.Time              `json:"system_time"` // Transaction commit time
	Status         ActionStatus           `json:"status"`
	Signature      string                 `json:"signature"`
	ContentHash    string                 `json:"content_hash"`
}

// ComputeHash calculates the canonical SHA-256 fingerprint of the action.
func (a *Action) ComputeHash() string {
	payload := fmt.Sprintf("%s|%s|%s|%s|%s|%d|%d",
		a.ID, a.Type, a.TargetNodeID, a.BranchID, a.Initiator,
		a.ValidTime.UnixNano(), a.SystemTime.UnixNano())
	if a.Parameters != nil {
		if b, err := json.Marshal(a.Parameters); err == nil {
			payload += "|" + string(b)
		}
	}
	hash := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(hash[:])
}

// BiTemporalRecord stores the temporal evolution of an entity property.
type BiTemporalRecord struct {
	PropertyKey string      `json:"property_key"`
	Value       interface{} `json:"value"`
	ValidStart  time.Time   `json:"valid_start"`
	ValidEnd    time.Time   `json:"valid_end"`
	SystemStart time.Time   `json:"system_start"`
	SystemEnd   time.Time   `json:"system_end"`
	ActionID    string      `json:"action_id"`
}

// Branch represents an isolated sandbox workspace for scenario testing.
type Branch struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	BaseBranch  string     `json:"base_branch"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	MergedAt    *time.Time `json:"merged_at,omitempty"`
	Status      string     `json:"status"` // "ACTIVE", "MERGED", "ABANDONED"
	Actions     []string   `json:"actions"` // Ordered list of Action IDs
}

// InvalidationSeverity indicates the operational risk level of a cascading impact.
type InvalidationSeverity string

const (
	SeverityLow      InvalidationSeverity = "LOW"
	SeverityMedium   InvalidationSeverity = "MEDIUM"
	SeverityHigh     InvalidationSeverity = "HIGH"
	SeverityCritical InvalidationSeverity = "CRITICAL"
)

// InvalidationRule defines a condition that triggers downstream graph alerts.
type InvalidationRule struct {
	RuleID         string               `json:"rule_id"`
	Name           string               `json:"name"`
	TriggerAction  ActionType           `json:"trigger_action"`
	ConditionField string               `json:"condition_field"`
	Severity       InvalidationSeverity `json:"severity"`
	ImpactReason   string               `json:"impact_reason"`
}

// InvalidationAlert is a detected cascading risk resulting from an action.
type InvalidationAlert struct {
	AlertID        string               `json:"alert_id"`
	RuleID         string               `json:"rule_id"`
	ActionID       string               `json:"action_id"`
	SourceNodeID   string               `json:"source_node_id"`
	AffectedNodeID string               `json:"affected_node_id"`
	Severity       InvalidationSeverity `json:"severity"`
	Reason         string               `json:"reason"`
	DetectedAt     time.Time            `json:"detected_at"`
}
