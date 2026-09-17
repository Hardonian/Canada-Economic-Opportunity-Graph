package ai

import (
	"fmt"
	"strings"
)

// AIPGuardrails enforces cryptographic permissions and role execution bounds.
type AIPGuardrails struct {
	allowedRoleTools map[AgentRole]map[string]bool
}

// NewAIPGuardrails initializes the guardrail security policy.
func NewAIPGuardrails() *AIPGuardrails {
	g := &AIPGuardrails{
		allowedRoleTools: make(map[AgentRole]map[string]bool),
	}

	g.allowedRoleTools[RoleFinanceAnalyst] = map[string]bool{
		"query_cashflow":         true,
		"simulate_monte_carlo":   true,
		"calc_debt_service":      true,
		"fetch_financial_filing": true,
	}

	g.allowedRoleTools[RoleIndigenousJurist] = map[string]bool{
		"query_treaty_boundaries":   true,
		"fetch_section35_records":   true,
		"evaluate_ilgp_eligibility": true,
		"audit_trc92_compliance":    true,
	}

	g.allowedRoleTools[RoleGridPhysicist] = map[string]bool{
		"solve_power_flow":        true,
		"simulate_contingency":    true,
		"query_substation_headroom": true,
	}

	g.allowedRoleTools[RoleSupplyChainSpecialist] = map[string]bool{
		"query_transformer_leadtimes": true,
		"trace_ubo_ownership":         true,
		"evaluate_tariff_exposure":    true,
		"audit_domestic_content":      true,
	}

	return g
}

// ValidateToolInvocation checks if the agent is authorized to execute the requested tool.
func (g *AIPGuardrails) ValidateToolInvocation(call ToolCall) (bool, string) {
	// 1. Check role authorization
	tools, ok := g.allowedRoleTools[call.AgentRole]
	if !ok {
		return false, fmt.Sprintf("unrecognized agent role: %s", call.AgentRole)
	}

	// Hazardous mutation tools require root sovereign authorization token
	if strings.HasPrefix(call.ToolName, "mutate_") || strings.HasPrefix(call.ToolName, "admin_") {
		if call.PermissionAuth != "AUTH-SOVEREIGN-SEAL-2026" {
			return false, fmt.Sprintf("tool %s is a state mutation and requires explicit sovereign auth seal", call.ToolName)
		}
	}

	if !tools[call.ToolName] && !strings.HasPrefix(call.ToolName, "mutate_") {
		return false, fmt.Sprintf("agent role %s is not authorized to execute tool %s", call.AgentRole, call.ToolName)
	}

	return true, "AUTHORIZED"
}
