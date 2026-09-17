package graphanalytics

import (
	"fmt"
	"sync"
)

// GraphTraversalEngine provides high-speed pathfinding and multi-hop reachability walks.
type GraphTraversalEngine struct {
	mu  sync.RWMutex
	adj map[string][]string
}

// NewGraphTraversalEngine builds the adjacency index.
func NewGraphTraversalEngine(edges map[string][]string) *GraphTraversalEngine {
	copied := make(map[string][]string)
	for k, v := range edges {
		c := make([]string, len(v))
		copy(c, v)
		copied[k] = c
	}
	return &GraphTraversalEngine{
		adj: copied,
	}
}

// FindShortestPath computes the shortest unweighted path using Breadth-First Search.
func (gte *GraphTraversalEngine) FindShortestPath(sourceID, targetID string) ([]string, error) {
	gte.mu.RLock()
	defer gte.mu.RUnlock()

	if sourceID == targetID {
		return []string{sourceID}, nil
	}

	queue := []string{sourceID}
	visited := map[string]bool{sourceID: true}
	parent := make(map[string]string)

	found := false
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == targetID {
			found = true
			break
		}

		for _, neighbor := range gte.adj[curr] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = curr
				queue = append(queue, neighbor)
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("no path found between %s and %s", sourceID, targetID)
	}

	// Reconstruct path
	path := []string{targetID}
	curr := targetID
	for curr != sourceID {
		curr = parent[curr]
		path = append([]string{curr}, path...)
	}

	return path, nil
}

// GetNeighborsK returns all nodes reachable within k degrees of separation with their distance.
func (gte *GraphTraversalEngine) GetNeighborsK(sourceID string, k int) map[string]int {
	gte.mu.RLock()
	defer gte.mu.RUnlock()

	distances := map[string]int{sourceID: 0}
	queue := []string{sourceID}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		d := distances[curr]

		if d >= k {
			continue
		}

		for _, neighbor := range gte.adj[curr] {
			if _, seen := distances[neighbor]; !seen {
				distances[neighbor] = d + 1
				queue = append(queue, neighbor)
			}
		}
	}

	delete(distances, sourceID)
	return distances
}
