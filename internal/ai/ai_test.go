package ai

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestAIPGuardrails(t *testing.T) {
	g := NewAIPGuardrails()

	// Authorized tool call
	call1 := ToolCall{
		CallID:    "call-1",
		ToolName:  "query_cashflow",
		AgentRole: RoleFinanceAnalyst,
	}
	ok1, _ := g.ValidateToolInvocation(call1)
	if !ok1 {
		t.Errorf("expected FinanceAnalyst to be authorized for query_cashflow")
	}

	// Unauthorized role tool call
	call2 := ToolCall{
		CallID:    "call-2",
		ToolName:  "solve_power_flow",
		AgentRole: RoleFinanceAnalyst,
	}
	ok2, _ := g.ValidateToolInvocation(call2)
	if ok2 {
		t.Errorf("expected FinanceAnalyst to be blocked from solve_power_flow")
	}

	// Mutation tool without sovereign seal
	call3 := ToolCall{
		CallID:    "call-3",
		ToolName:  "mutate_graph_state",
		AgentRole: RoleFinanceAnalyst,
	}
	ok3, _ := g.ValidateToolInvocation(call3)
	if ok3 {
		t.Errorf("expected mutation tool without auth seal to be blocked")
	}

	// Mutation tool WITH sovereign seal
	call4 := ToolCall{
		CallID:         "call-4",
		ToolName:       "mutate_graph_state",
		AgentRole:      RoleFinanceAnalyst,
		PermissionAuth: "AUTH-SOVEREIGN-SEAL-2026",
	}
	ok4, _ := g.ValidateToolInvocation(call4)
	if !ok4 {
		t.Errorf("expected mutation tool with valid auth seal to be authorized")
	}
}

func TestAgentMeshOrchestrator(t *testing.T) {
	g := NewAIPGuardrails()
	orchestrator := NewAgentMeshOrchestrator(g)

	project := &domain.Project{
		ID:           "proj-darlington-smr",
		Name:         "Darlington SMR Unit 1",
		CapexCAD:     3400000000,
		CurrentStage: domain.StageConstruction,
		Province:     "ON",
	}

	report := orchestrator.EvaluateProject(project)
	if report.ConsensusScore <= 0 {
		t.Errorf("expected positive consensus score, got %f", report.ConsensusScore)
	}
	if len(report.EvaluationsByRole) != 4 {
		t.Errorf("expected evaluations from 4 specialized roles, got %d", len(report.EvaluationsByRole))
	}
	if report.ConsensusVerdict == "" {
		t.Errorf("expected non-empty consensus verdict")
	}
}

func TestGroundedRAGEngine(t *testing.T) {
	rag := NewGroundedRAGEngine()
	cit := rag.GenerateCitation("Darlington_IAAC_Decision.pdf", "Section 4.2", 14, "Proponent shall maintain minimum cooling water flow of 12 m3/s")

	if cit.ContentHash == "" {
		t.Fatalf("expected non-empty citation hash")
	}

	if !rag.VerifyCitation(cit) {
		t.Errorf("expected citation verification to pass")
	}

	// Tamper with quote
	cit.ExactQuote = "Tampered quote"
	if rag.VerifyCitation(cit) {
		t.Errorf("expected tampered citation verification to fail")
	}
}

func TestMultiAgentNegotiator(t *testing.T) {
	neg := NewMultiAgentNegotiator()
	prop := neg.OptimizeContractTerms(3400000000, 300.0)

	if prop.PPAStrikePriceCADPerMWh <= 0 {
		t.Errorf("expected positive strike price")
	}
	if prop.DebtTenorYears != 30 {
		t.Errorf("expected 30 year tenor for >$2B project, got %d", prop.DebtTenorYears)
	}
	if prop.ParetoEfficiencyScore < 80.0 {
		t.Errorf("expected high Pareto score, got %f", prop.ParetoEfficiencyScore)
	}
}
