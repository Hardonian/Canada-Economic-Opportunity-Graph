package ai

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// AgentRole defines specialized autonomous industrial copilot capabilities.
type AgentRole string

const (
	RoleFinanceAnalyst        AgentRole = "FINANCE_ANALYST"
	RoleIndigenousJurist       AgentRole = "INDIGENOUS_JURIST"
	RoleGridPhysicist         AgentRole = "GRID_PHYSICIST"
	RoleSupplyChainSpecialist AgentRole = "SUPPLY_CHAIN_SPECIALIST"
)

// ToolCall represents a requested tool invocation by an agent.
type ToolCall struct {
	CallID         string                 `json:"call_id"`
	ToolName       string                 `json:"tool_name"`
	Arguments      map[string]interface{} `json:"arguments"`
	AgentRole      AgentRole              `json:"agent_role"`
	PermissionAuth string                 `json:"permission_auth"` // Signed token
}

// ToolResult represents the output of an executed tool call.
type ToolResult struct {
	CallID      string                 `json:"call_id"`
	Status      string                 `json:"status"` // "SUCCESS", "DENIED", "ERROR"
	Output      map[string]interface{} `json:"output"`
	ExecutionMs float64                `json:"execution_ms"`
}

// GroundedCitation anchors an analytical assertion to a verifiable source token span.
type GroundedCitation struct {
	CitationID   string `json:"citation_id"`
	DocumentName string `json:"document_name"`
	Section      string `json:"section"`
	PageNumber   int    `json:"page_number"`
	ExactQuote   string `json:"exact_quote"`
	ContentHash  string `json:"content_hash"`
}

// ComputeHash generates SHA-256 fingerprint of the citation quote.
func (c *GroundedCitation) ComputeHash() string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%s", c.DocumentName, c.PageNumber, c.ExactQuote)))
	return hex.EncodeToString(h[:])
}

// AgentEvaluation represents the formal findings of a specialized copilot.
type AgentEvaluation struct {
	AgentID         string             `json:"agent_id"`
	Role            AgentRole          `json:"role"`
	ProjectID       string             `json:"project_id"`
	Recommendation  string             `json:"recommendation"` // "PROCEED", "CONDITIONS_MANDATED", "HOLD", "REJECT"
	ConfidenceScore float64            `json:"confidence_score"` // 0.0 - 100.0
	KeyFindings     []string           `json:"key_findings"`
	CriticalRisks   []string           `json:"critical_risks"`
	Citations       []GroundedCitation `json:"citations"`
	EvaluatedAt     time.Time          `json:"evaluated_at"`
}

// ConsensusReport synthesizes the multi-agent peer review into an executive decision memo.
type ConsensusReport struct {
	ReportID             string            `json:"report_id"`
	ProjectID            string            `json:"project_id"`
	ConsensusVerdict     string            `json:"consensus_verdict"` // "APPROVED", "APPROVED_WITH_CONDITIONS", "REJECTED"
	ConsensusScore       float64           `json:"consensus_score"`   // 0.0 - 100.0
	EvaluationsByRole    map[AgentRole]*AgentEvaluation `json:"evaluations_by_role"`
	ConsensusActionItems []string          `json:"consensus_action_items"`
	GeneratedAt          time.Time         `json:"generated_at"`
}
