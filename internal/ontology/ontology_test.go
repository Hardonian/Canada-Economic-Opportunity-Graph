package ontology

import (
	"context"
	"testing"
	"time"
)

func TestWriteBackAndRollback(t *testing.T) {
	invalidator := NewCascadeInvalidator()
	engine := NewWriteBackEngine(invalidator)

	ctx := context.Background()
	action := &Action{
		ID:             "ACT-001",
		Type:           ActionAdjustCapitalTranche,
		TargetNodeID:   "proj-darlington-smr",
		TargetNodeType: "PROJECT",
		BranchID:       "main",
		Initiator:      "user-pension-lead",
		Parameters: map[string]interface{}{
			"senior_debt_cad": 1200000000.0,
			"equity_cad":      800000000.0,
		},
	}

	executed, err := engine.ExecuteAction(ctx, action)
	if err != nil {
		t.Fatalf("unexpected error executing action: %v", err)
	}
	if executed.Status != ActionStatusApplied {
		t.Errorf("expected status %s, got %s", ActionStatusApplied, executed.Status)
	}

	state := engine.GetEntityState("proj-darlington-smr", "main")
	if state["senior_debt_cad"] != 1200000000.0 {
		t.Errorf("expected senior_debt_cad 1.2B, got %v", state["senior_debt_cad"])
	}

	// Revert
	reverted, err := engine.RevertAction("ACT-001")
	if err != nil {
		t.Fatalf("failed to revert action: %v", err)
	}
	if reverted.Status != ActionStatusReverted {
		t.Errorf("expected status %s, got %s", ActionStatusReverted, reverted.Status)
	}
}

func TestCascadeInvalidation(t *testing.T) {
	invalidator := NewCascadeInvalidator()
	engine := NewWriteBackEngine(invalidator)

	ctx := context.Background()
	slipAction := &Action{
		ID:             "ACT-SLIP-01",
		Type:           ActionSetMilestoneDelay,
		TargetNodeID:   "proj-crawford-nickel",
		TargetNodeType: "PROJECT",
		BranchID:       "main",
		Initiator:      "epcm-manager",
		Parameters: map[string]interface{}{
			"delay_months": 8.5,
			"milestone":    "COMMERCIAL_OPERATION",
		},
	}

	_, err := engine.ExecuteAction(ctx, slipAction)
	if err != nil {
		t.Fatalf("failed to execute slip action: %v", err)
	}

	alerts := invalidator.GetActiveAlerts()
	if len(alerts) == 0 {
		t.Fatalf("expected at least 1 cascade invalidation alert, got 0")
	}

	found := false
	for _, a := range alerts {
		if a.RuleID == "RULE-MILESTONE-SLIP-OFFTAKE" {
			found = true
			if a.Severity != SeverityCritical {
				t.Errorf("expected CRITICAL severity, got %s", a.Severity)
			}
		}
	}
	if !found {
		t.Errorf("expected RULE-MILESTONE-SLIP-OFFTAKE alert to be triggered")
	}
}

func TestBranchingAndMerging(t *testing.T) {
	invalidator := NewCascadeInvalidator()
	engine := NewWriteBackEngine(invalidator)
	bm := NewBranchManager(engine)

	ctx := context.Background()

	// Base action on main
	baseAct := &Action{
		ID:             "ACT-MAIN-01",
		Type:           ActionUpdateProjectStage,
		TargetNodeID:   "proj-oneida-bess",
		TargetNodeType: "PROJECT",
		BranchID:       "main",
		Initiator:      "system",
		Parameters: map[string]interface{}{
			"stage": "CONSTRUCTION",
		},
	}
	if _, err := engine.ExecuteAction(ctx, baseAct); err != nil {
		t.Fatalf("failed to execute base action: %v", err)
	}

	// Create scenario branch
	branch, err := bm.CreateBranch("scenario-delay", "Delay Scenario", "Testing 12-month grid delay", "main", "analyst")
	if err != nil {
		t.Fatalf("failed to create branch: %v", err)
	}
	if branch.ID != "scenario-delay" {
		t.Errorf("unexpected branch ID: %s", branch.ID)
	}

	// Execute action on scenario branch
	scenAct := &Action{
		ID:             "ACT-SCEN-01",
		Type:           ActionSetMilestoneDelay,
		TargetNodeID:   "proj-oneida-bess",
		TargetNodeType: "PROJECT",
		BranchID:       "scenario-delay",
		Initiator:      "analyst",
		Parameters: map[string]interface{}{
			"stage": "COMMISSIONING",
		},
	}
	if _, err := engine.ExecuteAction(ctx, scenAct); err != nil {
		t.Fatalf("failed to execute scenario action: %v", err)
	}

	// Verify isolation
	mainState := engine.GetEntityState("proj-oneida-bess", "main")
	scenState := engine.GetEntityState("proj-oneida-bess", "scenario-delay")

	if mainState["stage"] != "CONSTRUCTION" {
		t.Errorf("expected main stage to remain CONSTRUCTION, got %v", mainState["stage"])
	}
	if scenState["stage"] != "COMMISSIONING" {
		t.Errorf("expected scenario stage to be COMMISSIONING, got %v", scenState["stage"])
	}

	// Merge branch
	if _, err := bm.MergeBranch("scenario-delay", "main", "lead"); err != nil {
		t.Fatalf("failed to merge branch: %v", err)
	}

	mergedMainState := engine.GetEntityState("proj-oneida-bess", "main")
	if mergedMainState["stage"] != "COMMISSIONING" {
		t.Errorf("expected merged main stage to be COMMISSIONING, got %v", mergedMainState["stage"])
	}
}

func TestBiTemporalTimeTravel(t *testing.T) {
	engine := NewWriteBackEngine(nil)
	evaluator := NewTemporalEvaluator(engine)
	ctx := context.Background()

	t0 := time.Now().Add(-2 * time.Hour).UTC()
	t1 := time.Now().Add(-1 * time.Hour).UTC()

	// Action at t0
	act0 := &Action{
		ID:           "ACT-T0",
		Type:         ActionUpdateProjectStage,
		TargetNodeID: "proj-site-c",
		BranchID:     "main",
		Parameters:   map[string]interface{}{"capacity_mw": 900.0},
		ValidTime:    t0,
	}
	if _, err := engine.ExecuteAction(ctx, act0); err != nil {
		t.Fatalf("err act0: %v", err)
	}

	// Action at t1
	act1 := &Action{
		ID:           "ACT-T1",
		Type:         ActionUpdateProjectStage,
		TargetNodeID: "proj-site-c",
		BranchID:     "main",
		Parameters:   map[string]interface{}{"capacity_mw": 1100.0},
		ValidTime:    t1,
	}
	if _, err := engine.ExecuteAction(ctx, act1); err != nil {
		t.Fatalf("err act1: %v", err)
	}

	// Query as of t0 + 10 mins
	asOfT0 := evaluator.GetNodeAsOf("proj-site-c", t0.Add(10*time.Minute), time.Time{})
	if asOfT0["capacity_mw"] != 900.0 {
		t.Errorf("expected capacity 900.0 at t0, got %v", asOfT0["capacity_mw"])
	}

	// Query as of now
	asOfNow := evaluator.GetNodeAsOf("proj-site-c", time.Now().UTC(), time.Time{})
	if asOfNow["capacity_mw"] != 1100.0 {
		t.Errorf("expected capacity 1100.0 at now, got %v", asOfNow["capacity_mw"])
	}
}
