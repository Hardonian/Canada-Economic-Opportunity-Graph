package api

import (
	"net/http"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/cegs"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/connector"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// RuntimeInfo is the intentionally small, non-secret runtime projection
// returned by the system status endpoint.
type RuntimeInfo struct {
	Environment       string `json:"environment"`
	StorageMode       string `json:"storage_mode"`
	SchedulerEnabled  bool   `json:"scheduler_enabled"`
	SchedulerInterval string `json:"scheduler_interval,omitempty"`
}

type productCapability struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Endpoint string `json:"endpoint"`
	Status   string `json:"status"`
}

var productCapabilities = []productCapability{
	{ID: "capital-radar", Label: "Capital and project radar", Endpoint: "/api/v1/radar", Status: "ACTIVE"},
	{ID: "project-intelligence", Label: "Project diligence and provenance", Endpoint: "/api/v1/projects", Status: "ACTIVE"},
	{ID: "planning", Label: "National planning and scenario models", Endpoint: "/api/v1/planning/optimize", Status: "ACTIVE"},
	{ID: "finance", Label: "Project finance and transition taxonomy", Endpoint: "/api/v1/finance/transition-taxonomy", Status: "ACTIVE"},
	{ID: "graph", Label: "GraphQL and graph analytics", Endpoint: "/api/v1/graphql", Status: "ACTIVE"},
	{ID: "standards", Label: "CEGS exports and interoperability", Endpoint: "/api/v1/cegs/export", Status: "ACTIVE"},
	{ID: "events", Label: "Server-sent event stream", Endpoint: "/api/v1/stream/events", Status: "ACTIVE"},
}

// handleSystemStatus returns a single operational control-plane view backed by
// the live store and the connector registry used for ingestion.
func (s *Server) handleSystemStatus(w http.ResponseWriter, r *http.Request) {
	stats, err := s.store.GetRadarStats(r.Context())
	if err != nil || stats == nil {
		writeError(w, r, http.StatusServiceUnavailable, "system_status_unavailable", "Operational status is temporarily unavailable.")
		return
	}

	connectors := &connector.HealthSnapshot{Generated: time.Now().UTC()}
	if s.connectorRegistry != nil {
		connectors = s.connectorRegistry.AggregateHealth()
	}

	status := "OPERATIONAL"
	if stats.TotalProjects == 0 || stats.DataStatus == domain.StatusUnavailable || connectors.Broken > 0 {
		status = "DEGRADED"
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": status,
		"service": map[string]string{
			"name":      "CanadaOpportunityGraph API",
			"version":   "1.0.0",
			"cegs_spec": cegs.SpecVersion,
		},
		"runtime":      s.runtimeInfo,
		"data":         stats,
		"ingestion":    connectors,
		"capabilities": productCapabilities,
		"generated_at": time.Now().UTC(),
	})
}
