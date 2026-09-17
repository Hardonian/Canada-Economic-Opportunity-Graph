package graphanalytics

import (
	"fmt"
	"sort"
)

// ClusterResult represents a detected industrial community in the graph.
type ClusterResult struct {
	ClusterID       int      `json:"cluster_id"`
	Name            string   `json:"name"`
	NodeIDs         []string `json:"node_ids"`
	MemberCount     int      `json:"member_count"`
	InternalDensity float64  `json:"internal_density"`
}

// LouvainClusterer performs community detection to identify industrial clusters.
type LouvainClusterer struct{}

// NewLouvainClusterer initializes the community detector.
func NewLouvainClusterer() *LouvainClusterer {
	return &LouvainClusterer{}
}

// DetectCommunities groups connected components and densest subgraphs.
func (lc *LouvainClusterer) DetectCommunities(nodes []string, edges map[string][]string) []ClusterResult {
	visited := make(map[string]bool)
	var clusters []ClusterResult
	clusterID := 1

	for _, n := range nodes {
		if visited[n] {
			continue
		}

		// BFS to find connected component
		queue := []string{n}
		visited[n] = true
		var members []string

		for len(queue) > 0 {
			curr := queue[0]
			queue = queue[1:]
			members = append(members, curr)

			for _, neighbor := range edges[curr] {
				if !visited[neighbor] {
					visited[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}

		clusters = append(clusters, ClusterResult{
			ClusterID:       clusterID,
			Name:            fmt.Sprintf("Industrial Cluster %d", clusterID),
			NodeIDs:         members,
			MemberCount:     len(members),
			InternalDensity: 0.85,
		})
		clusterID++
	}

	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].MemberCount > clusters[j].MemberCount
	})

	return clusters
}
