package graphanalytics

import (
	"fmt"
	"sync"
)

// HyperEdge connects three or more nodes simultaneously (e.g. multi-party syndicates).
type HyperEdge struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"` // "MULTI_NATION_EQUITY", "TRI_PARTY_OFFTAKE", "PUBLIC_PRIVATE_GUARANTEE"
	NodeIDs    []string               `json:"node_ids"`
	Attributes map[string]interface{} `json:"attributes"`
}

// HyperGraphManager manages hyper-edges and computes multi-party intersections.
type HyperGraphManager struct {
	mu         sync.RWMutex
	hyperEdges map[string]*HyperEdge
	nodeIndex  map[string][]string // nodeID -> []hyperEdgeID
}

// NewHyperGraphManager initializes the hyper-graph engine.
func NewHyperGraphManager() *HyperGraphManager {
	return &HyperGraphManager{
		hyperEdges: make(map[string]*HyperEdge),
		nodeIndex:  make(map[string][]string),
	}
}

// AddHyperEdge registers a multi-party connection.
func (hgm *HyperGraphManager) AddHyperEdge(edge *HyperEdge) error {
	if edge == nil || edge.ID == "" {
		return fmt.Errorf("hyper-edge cannot be empty")
	}
	if len(edge.NodeIDs) < 2 {
		return fmt.Errorf("hyper-edge must connect at least 2 nodes, got %d", len(edge.NodeIDs))
	}

	hgm.mu.Lock()
	defer hgm.mu.Unlock()

	hgm.hyperEdges[edge.ID] = edge
	for _, nID := range edge.NodeIDs {
		hgm.nodeIndex[nID] = append(hgm.nodeIndex[nID], edge.ID)
	}
	return nil
}

// GetHyperEdgesForNode returns all hyper-edges a node participates in.
func (hgm *HyperGraphManager) GetHyperEdgesForNode(nodeID string) []*HyperEdge {
	hgm.mu.RLock()
	defer hgm.mu.RUnlock()

	edgeIDs := hgm.nodeIndex[nodeID]
	res := make([]*HyperEdge, 0, len(edgeIDs))
	for _, id := range edgeIDs {
		if e, ok := hgm.hyperEdges[id]; ok {
			res = append(res, e)
		}
	}
	return res
}
