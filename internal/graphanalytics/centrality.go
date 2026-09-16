package graphanalytics

import (
	"sort"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// NodeCentrality holds graph-theoretic metrics for an entity or project.
type NodeCentrality struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Type                string   `json:"type"` // "PROJECT", "INFRASTRUCTURE", "SUPPLIER", "OFFTAKER"
	InDegree            int      `json:"in_degree"`
	OutDegree           int      `json:"out_degree"`
	TotalDegree         int      `json:"total_degree"`
	BetweennessScore    float64  `json:"betweenness_score"` // 0.0 - 1.0
	IsSinglePointOfFail bool     `json:"is_single_point_of_failure"`
	CriticalityRating   string   `json:"criticality_rating"` // LOW, MEDIUM, HIGH, SOVEREIGN_CRITICAL
	DependentIDs        []string `json:"dependent_ids"`
}

// CascadingFailureResult models the system-wide fallout of a compromised or delayed node.
type CascadingFailureResult struct {
	TriggerNodeID       string           `json:"trigger_node_id"`
	TriggerNodeName     string           `json:"trigger_node_name"`
	DirectlyAffectedIDs []string         `json:"directly_affected_ids"`
	TotalCompromisedIDs []string         `json:"total_compromised_ids"`
	FrozenCapexCAD      int64            `json:"frozen_capex_cad"`
	CascadingDepth      int              `json:"cascading_depth"`
	VulnerableSectors   []domain.Sector  `json:"vulnerable_sectors"`
	SystemicRiskScore   float64          `json:"systemic_risk_score"` // 0.0 - 100.0
	EvaluatedAt         time.Time        `json:"evaluated_at"`
}

// GraphAnalyzer computes topological properties and simulates cascades.
type GraphAnalyzer struct{}

func NewGraphAnalyzer() *GraphAnalyzer {
	return &GraphAnalyzer{}
}

// AnalyzeNetwork calculates degree and betweenness centrality over projects and relationships.
func (ga *GraphAnalyzer) AnalyzeNetwork(projects []*domain.Project, relationships []*domain.Relationship) []NodeCentrality {
	inMap := make(map[string]int)
	outMap := make(map[string]int)
	adj := make(map[string][]string)
	names := make(map[string]string)
	capexMap := make(map[string]int64)

	for _, p := range projects {
		names[p.ID] = p.Name
		capexMap[p.ID] = p.CapexCAD
	}

	for _, rel := range relationships {
		outMap[rel.SourceEntityID]++
		inMap[rel.TargetEntityID]++
		adj[rel.SourceEntityID] = append(adj[rel.SourceEntityID], rel.TargetEntityID)
		adj[rel.TargetEntityID] = append(adj[rel.TargetEntityID], rel.SourceEntityID)
	}

	results := make([]NodeCentrality, 0, len(projects))
	maxDeg := 1

	for _, p := range projects {
		inD := inMap[p.ID]
		outD := outMap[p.ID]
		tot := inD + outD
		if tot > maxDeg {
			maxDeg = tot
		}
	}

	for _, p := range projects {
		inD := inMap[p.ID]
		outD := outMap[p.ID]
		tot := inD + outD

		// Normalized centrality proxy
		betweenness := float64(tot) / float64(maxDeg)
		isSPOF := tot >= 3 || p.CapexCAD >= 3_000_000_000

		crit := "LOW"
		switch {
		case isSPOF && betweenness > 0.6:
			crit = "SOVEREIGN_CRITICAL"
		case betweenness > 0.4:
			crit = "HIGH"
		case betweenness > 0.2:
			crit = "MEDIUM"
		}

		results = append(results, NodeCentrality{
			ID:                  p.ID,
			Name:                p.Name,
			Type:                "PROJECT",
			InDegree:            inD,
			OutDegree:           outD,
			TotalDegree:         tot,
			BetweennessScore:    betweenness,
			IsSinglePointOfFail: isSPOF,
			CriticalityRating:   crit,
			DependentIDs:        adj[p.ID],
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].BetweennessScore > results[j].BetweennessScore
	})

	return results
}

// SimulateCascade computes BFS propagation of a node failure across the supply and dependency graph.
func (ga *GraphAnalyzer) SimulateCascade(triggerID string, projects []*domain.Project, relationships []*domain.Relationship) *CascadingFailureResult {
	projectMap := make(map[string]*domain.Project)
	for _, p := range projects {
		projectMap[p.ID] = p
	}

	adjDownstream := make(map[string][]string)
	for _, rel := range relationships {
		adjDownstream[rel.SourceEntityID] = append(adjDownstream[rel.SourceEntityID], rel.TargetEntityID)
		if rel.ProjectID != "" && rel.ProjectID != rel.SourceEntityID {
			adjDownstream[rel.ProjectID] = append(adjDownstream[rel.ProjectID], rel.TargetEntityID)
		}
	}

	visited := make(map[string]bool)
	queue := []string{triggerID}
	visited[triggerID] = true

	directlyAffected := make([]string, 0)
	for _, next := range adjDownstream[triggerID] {
		if !visited[next] {
			directlyAffected = append(directlyAffected, next)
		}
	}

	depth := 0
	for len(queue) > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			curr := queue[0]
			queue = queue[1:]

			for _, neighbor := range adjDownstream[curr] {
				if !visited[neighbor] {
					visited[neighbor] = true
					queue = append(queue, neighbor)
				}
			}
		}
		if len(queue) > 0 {
			depth++
		}
	}

	allCompromised := make([]string, 0, len(visited))
	var frozenCapex int64
	sectorSet := make(map[domain.Sector]bool)

	for id := range visited {
		if id != triggerID {
			allCompromised = append(allCompromised, id)
			if p, ok := projectMap[id]; ok {
				frozenCapex += p.CapexCAD
				sectorSet[p.Sector] = true
			}
		}
	}

	sectors := make([]domain.Sector, 0, len(sectorSet))
	for s := range sectorSet {
		sectors = append(sectors, s)
	}

	triggerName := triggerID
	if p, ok := projectMap[triggerID]; ok {
		triggerName = p.Name
		frozenCapex += p.CapexCAD
		sectors = append(sectors, p.Sector)
	}

	systemicRisk := float64(len(allCompromised))*10.0 + (float64(frozenCapex)/1e9)*5.0
	if systemicRisk > 100.0 {
		systemicRisk = 100.0
	}

	return &CascadingFailureResult{
		TriggerNodeID:       triggerID,
		TriggerNodeName:     triggerName,
		DirectlyAffectedIDs: directlyAffected,
		TotalCompromisedIDs: allCompromised,
		FrozenCapexCAD:      frozenCapex,
		CascadingDepth:      depth,
		VulnerableSectors:   sectors,
		SystemicRiskScore:   systemicRisk,
		EvaluatedAt:         time.Now().UTC(),
	}
}
