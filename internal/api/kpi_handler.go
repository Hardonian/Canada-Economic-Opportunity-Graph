package api

import (
	"net/http"
	"strconv"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/indicators"
)

// handleKPIList returns all registered KPI definitions, optional category filtering,
// and current national baseline observations.
func (s *Server) handleKPIList(w http.ResponseWriter, r *http.Request) {
	definitions := indicators.CanonicalRegistry()
	categoryFilter := r.URL.Query().Get("category")

	filteredDefs := make([]*indicators.KPIDefinition, 0, len(definitions))
	for _, def := range definitions {
		if categoryFilter == "" || string(def.Category) == categoryFilter {
			filteredDefs = append(filteredDefs, def)
		}
	}

	obs := s.kpiFeedEngine.GetObservations("NATIONAL")

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"total_count":         len(filteredDefs),
		"definitions":         filteredDefs,
		"national_benchmarks": obs,
	})
}

// handleKPIFeeds returns the latest streaming commodity and macroeconomic indicator ticks.
func (s *Server) handleKPIFeeds(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	ticks := s.kpiFeedEngine.GetRecentTicks(limit)
	latest := s.kpiFeedEngine.GetLatestTicks()

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"latest_ticks": latest,
		"recent_ticks": ticks,
		"total_count":  len(ticks),
	})
}

// handleKPIProject computes the complete 8-pillar indicator scorecard for a specific project.
func (s *Server) handleKPIProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, r, http.StatusBadRequest, "invalid_parameter", "Project id or slug is required.")
		return
	}

	project, err := s.resolveProject(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, r, http.StatusNotFound, "not_found", "Project not found.")
		return
	}

	scorecard := s.kpiEvaluator.Evaluate(project)
	writeJSON(w, http.StatusOK, scorecard)
}

// handleKPISummary returns aggregate national metrics and cross-sector rollups.
func (s *Server) handleKPISummary(w http.ResponseWriter, r *http.Request) {
	summary := s.kpiFeedEngine.BuildMacroSummary()
	writeJSON(w, http.StatusOK, summary)
}

// handleKPISnapshot exports an immutable cryptographic snapshot containing all definitions,
// observations, and macro metrics.
func (s *Server) handleKPISnapshot(w http.ResponseWriter, r *http.Request) {
	snapshot, err := s.kpiFeedEngine.GenerateSnapshot()
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "snapshot_error", "Failed to generate KPI snapshot.")
		return
	}
	writeJSON(w, http.StatusOK, snapshot)
}
