package ontology

import (
	"fmt"
	"sync"
	"time"
)

// BranchManager provides Git-like branch isolation and merging for graph state.
type BranchManager struct {
	mu          sync.RWMutex
	branches    map[string]*Branch
	writeEngine *WriteBackEngine
}

// NewBranchManager initializes the branch manager with the canonical 'main' branch.
func NewBranchManager(writeEngine *WriteBackEngine) *BranchManager {
	bm := &BranchManager{
		branches:    make(map[string]*Branch),
		writeEngine: writeEngine,
	}

	mainBranch := &Branch{
		ID:          "main",
		Name:        "Main Sovereign Trunk",
		Description: "Canonical baseline production state of Canadian Economic Opportunity Graph",
		BaseBranch:  "",
		CreatedBy:   "SYSTEM",
		CreatedAt:   time.Now().UTC(),
		Status:      "ACTIVE",
		Actions:     make([]string, 0),
	}
	bm.branches["main"] = mainBranch
	return bm
}

// CreateBranch forks a new isolated workspace from an existing base branch.
func (bm *BranchManager) CreateBranch(id, name, desc, baseBranch, createdBy string) (*Branch, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if id == "" {
		return nil, fmt.Errorf("branch ID cannot be empty")
	}
	if _, exists := bm.branches[id]; exists {
		return nil, fmt.Errorf("branch %s already exists", id)
	}
	if baseBranch == "" {
		baseBranch = "main"
	}
	if _, exists := bm.branches[baseBranch]; !exists {
		return nil, fmt.Errorf("base branch %s does not exist", baseBranch)
	}

	newBranch := &Branch{
		ID:          id,
		Name:        name,
		Description: desc,
		BaseBranch:  baseBranch,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now().UTC(),
		Status:      "ACTIVE",
		Actions:     make([]string, 0),
	}

	// Copy base branch state into new branch in the write engine
	if bm.writeEngine != nil {
		bm.writeEngine.mu.Lock()
		if bm.writeEngine.branchStates[baseBranch] != nil {
			bm.writeEngine.branchStates[id] = make(map[string]map[string]interface{})
			for nodeID, props := range bm.writeEngine.branchStates[baseBranch] {
				bm.writeEngine.branchStates[id][nodeID] = make(map[string]interface{})
				for k, v := range props {
					bm.writeEngine.branchStates[id][nodeID][k] = v
				}
			}
		}
		bm.writeEngine.mu.Unlock()
	}

	bm.branches[id] = newBranch
	return newBranch, nil
}

// GetBranch retrieves a branch by ID.
func (bm *BranchManager) GetBranch(id string) (*Branch, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	b, exists := bm.branches[id]
	if !exists {
		return nil, fmt.Errorf("branch %s not found", id)
	}
	return b, nil
}

// ListBranches lists all active or merged branches.
func (bm *BranchManager) ListBranches() []*Branch {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	list := make([]*Branch, 0, len(bm.branches))
	for _, b := range bm.branches {
		list = append(list, b)
	}
	return list
}

// MergeBranch applies all mutations from sourceBranch into targetBranch.
func (bm *BranchManager) MergeBranch(sourceBranchID, targetBranchID, mergedBy string) (*Branch, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	src, ok := bm.branches[sourceBranchID]
	if !ok {
		return nil, fmt.Errorf("source branch %s not found", sourceBranchID)
	}
	target, ok := bm.branches[targetBranchID]
	if !ok {
		return nil, fmt.Errorf("target branch %s not found", targetBranchID)
	}
	if src.Status != "ACTIVE" {
		return nil, fmt.Errorf("source branch is %s, must be ACTIVE", src.Status)
	}

	// Apply node state mutations
	if bm.writeEngine != nil {
		bm.writeEngine.mu.Lock()
		if bm.writeEngine.branchStates[sourceBranchID] != nil {
			if bm.writeEngine.branchStates[targetBranchID] == nil {
				bm.writeEngine.branchStates[targetBranchID] = make(map[string]map[string]interface{})
			}
			for nodeID, props := range bm.writeEngine.branchStates[sourceBranchID] {
				if bm.writeEngine.branchStates[targetBranchID][nodeID] == nil {
					bm.writeEngine.branchStates[targetBranchID][nodeID] = make(map[string]interface{})
				}
				for k, v := range props {
					bm.writeEngine.branchStates[targetBranchID][nodeID][k] = v
				}
			}
		}
		bm.writeEngine.mu.Unlock()
	}

	now := time.Now().UTC()
	src.Status = "MERGED"
	src.MergedAt = &now
	target.Actions = append(target.Actions, src.Actions...)

	return target, nil
}
