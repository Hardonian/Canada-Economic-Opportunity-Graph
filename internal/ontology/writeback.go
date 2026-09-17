package ontology

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// WriteBackEngine coordinates transactional graph mutations and isolated rollbacks.
type WriteBackEngine struct {
	mu           sync.RWMutex
	actions      map[string]*Action
	nodeActions  map[string][]string               // nodeID -> actionIDs
	branchStates map[string]map[string]map[string]interface{} // branchID -> nodeID -> properties
	invalidator  *CascadeInvalidator
}

// NewWriteBackEngine initializes a thread-safe write-back engine.
func NewWriteBackEngine(invalidator *CascadeInvalidator) *WriteBackEngine {
	return &WriteBackEngine{
		actions:      make(map[string]*Action),
		nodeActions:  make(map[string][]string),
		branchStates: make(map[string]map[string]map[string]interface{}),
		invalidator:  invalidator,
	}
}

// ExecuteAction applies an action transactionally to the graph branch.
func (e *WriteBackEngine) ExecuteAction(ctx context.Context, action *Action) (*Action, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}
	if action.ID == "" {
		return nil, fmt.Errorf("action ID is required")
	}
	if action.TargetNodeID == "" {
		return nil, fmt.Errorf("target node ID is required")
	}
	if action.BranchID == "" {
		action.BranchID = "main"
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.actions[action.ID]; exists {
		return nil, fmt.Errorf("action %s already executed", action.ID)
	}

	// Prepare current state
	branch := action.BranchID
	if e.branchStates[branch] == nil {
		e.branchStates[branch] = make(map[string]map[string]interface{})
	}
	if e.branchStates[branch][action.TargetNodeID] == nil {
		e.branchStates[branch][action.TargetNodeID] = make(map[string]interface{})
	}

	// Capture previous state
	prevState := make(map[string]interface{})
	for k, v := range e.branchStates[branch][action.TargetNodeID] {
		prevState[k] = v
	}
	action.PreviousState = prevState

	// Apply mutations
	newState := make(map[string]interface{})
	for k, v := range prevState {
		newState[k] = v
	}
	for k, v := range action.Parameters {
		newState[k] = v
		e.branchStates[branch][action.TargetNodeID][k] = v
	}
	action.NewState = newState

	if action.ValidTime.IsZero() {
		action.ValidTime = time.Now().UTC()
	}
	action.SystemTime = time.Now().UTC()
	action.Status = ActionStatusApplied
	action.ContentHash = action.ComputeHash()

	e.actions[action.ID] = action
	e.nodeActions[action.TargetNodeID] = append(e.nodeActions[action.TargetNodeID], action.ID)

	// Trigger cascade invalidation evaluation if invalidator configured
	if e.invalidator != nil {
		e.invalidator.EvaluateCascades(action)
	}

	return action, nil
}

// RevertAction rolls back an applied action to its previous state.
func (e *WriteBackEngine) RevertAction(actionID string) (*Action, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	action, exists := e.actions[actionID]
	if !exists {
		return nil, fmt.Errorf("action %s not found", actionID)
	}
	if action.Status != ActionStatusApplied {
		return nil, fmt.Errorf("action %s is not in APPLIED state (current: %s)", actionID, action.Status)
	}

	branch := action.BranchID
	if e.branchStates[branch] != nil && e.branchStates[branch][action.TargetNodeID] != nil {
		// Restore previous state
		for k := range action.NewState {
			delete(e.branchStates[branch][action.TargetNodeID], k)
		}
		for k, v := range action.PreviousState {
			e.branchStates[branch][action.TargetNodeID][k] = v
		}
	}

	action.Status = ActionStatusReverted
	action.SystemTime = time.Now().UTC()
	action.ContentHash = action.ComputeHash()

	return action, nil
}

// GetEntityState returns the current properties of a node on a given branch.
func (e *WriteBackEngine) GetEntityState(nodeID string, branchID string) map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if branchID == "" {
		branchID = "main"
	}

	result := make(map[string]interface{})
	if b, ok := e.branchStates[branchID]; ok {
		if nodeProps, exists := b[nodeID]; exists {
			for k, v := range nodeProps {
				result[k] = v
			}
		}
	}
	return result
}

// GetAction retrieves an individual action by ID.
func (e *WriteBackEngine) GetAction(actionID string) (*Action, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	a, ok := e.actions[actionID]
	return a, ok
}

// GetActionHistory returns all actions executed against a node in chronological order.
func (e *WriteBackEngine) GetActionHistory(nodeID string) []*Action {
	e.mu.RLock()
	defer e.mu.RUnlock()

	actionIDs := e.nodeActions[nodeID]
	history := make([]*Action, 0, len(actionIDs))
	for _, id := range actionIDs {
		if a, ok := e.actions[id]; ok {
			history = append(history, a)
		}
	}
	return history
}
