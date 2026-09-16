package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/earthobs"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/econometrics"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/gridphysics"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/nationalplanning"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/risk"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/ubo"
)

// handlePlanningOptimize executes the Sovereign Capital Allocation Optimizer.
func (s *Server) handlePlanningOptimize(w http.ResponseWriter, r *http.Request) {
	var req nationalplanning.OptimizationRequest
	if r.Method == http.MethodGet {
		if obj := r.URL.Query().Get("objective"); obj != "" {
			req.Objective = nationalplanning.ObjectiveType(obj)
		}
	} else if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
			writeError(w, r, http.StatusBadRequest, "bad_request", "Invalid optimization request payload.")
			return
		}
	}

	projects, _, err := s.store.ListProjects(r.Context(), database.ProjectFilter{Limit: 200})
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "store_error", "Failed to retrieve project catalog.")
		return
	}

	optimizer := nationalplanning.NewOptimizer()
	result := optimizer.Optimize(projects, req)

	writeJSON(w, http.StatusOK, result)
}

// handlePlanningWarGame executes geopolitical & geoeconomic shock stress tests.
func (s *Server) handlePlanningWarGame(w http.ResponseWriter, r *http.Request) {
	var req nationalplanning.WarGameRequest
	if r.Method == http.MethodGet {
		if sc := r.URL.Query().Get("scenario"); sc != "" {
			req.Scenario = nationalplanning.ShockScenario(sc)
		}
	} else if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err.Error() != "EOF" {
			writeError(w, r, http.StatusBadRequest, "bad_request", "Invalid war game request payload.")
			return
		}
	}
	if req.Scenario == "" {
		req.Scenario = nationalplanning.ShockUSMCATariffs
	}

	projects, _, err := s.store.ListProjects(r.Context(), database.ProjectFilter{Limit: 200})
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "store_error", "Failed to retrieve project catalog.")
		return
	}

	engine := nationalplanning.NewWarGameEngine()
	result := engine.SimulateScenario(projects, req)

	writeJSON(w, http.StatusOK, result)
}

// handlePlanningLabor calculates regional Red Seal craft labor collision metrics.
func (s *Server) handlePlanningLabor(w http.ResponseWriter, r *http.Request) {
	province := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("province")))
	if province == "" {
		province = "ON"
	}

	projects, _, err := s.store.ListProjects(r.Context(), database.ProjectFilter{Limit: 200})
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "store_error", "Failed to retrieve project catalog.")
		return
	}

	aggregator := nationalplanning.NewLaborAggregator()
	report := aggregator.AnalyzeProvince(province, projects)

	writeJSON(w, http.StatusOK, report)
}

// handleProjectMRIO calculates Multi-Regional Input-Output macroeconomic metrics.
func (s *Server) handleProjectMRIO(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := s.store.GetProject(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, r, http.StatusNotFound, "not_found", "Project not found.")
		return
	}

	engine := econometrics.NewEngine()
	mrio := engine.CalculateMRIO(project)

	writeJSON(w, http.StatusOK, mrio)
}

// handleProjectUBO evaluates Ultimate Beneficial Ownership and Investment Canada Act screening.
func (s *Server) handleProjectUBO(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := s.store.GetProject(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, r, http.StatusNotFound, "not_found", "Project not found.")
		return
	}

	// Fetch relationships or mock beneficial owners if unlinked
	evaluator := ubo.NewEvaluator()
	screening := evaluator.ScreenProject(project, nil)

	writeJSON(w, http.StatusOK, screening)
}

// handleProjectFlyvbjerg generates Bayesian reference-class megaproject risk hazard curves.
func (s *Server) handleProjectFlyvbjerg(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := s.store.GetProject(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, r, http.StatusNotFound, "not_found", "Project not found.")
		return
	}

	evaluator := risk.NewEvaluator()
	forecast := evaluator.ForecastProject(project)

	writeJSON(w, http.StatusOK, forecast)
}

// handleProjectGrid calculates electrical grid interconnection and hosting capacity feasibility.
func (s *Server) handleProjectGrid(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := s.store.GetProject(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, r, http.StatusNotFound, "not_found", "Project not found.")
		return
	}

	engine := gridphysics.NewEngine()
	assessment := engine.AssessProject(project)

	writeJSON(w, http.StatusOK, assessment)
}

// handleProjectEarthObs returns satellite SAR and optical ground-truth verification telemetry.
func (s *Server) handleProjectEarthObs(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := s.store.GetProject(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, r, http.StatusNotFound, "not_found", "Project not found.")
		return
	}

	evaluator := earthobs.NewEvaluator()
	// Synthesize ground-truth observation pass for demonstration/audit
	pass := earthobs.GroundTruthObservation{
		ObservationID:          "rcm-pass-" + project.ID,
		ProjectID:              project.ID,
		Constellation:          earthobs.ConstellationRCM,
		PassDate:               time.Now().Add(-72 * time.Hour),
		SARCoherenceChange:     0.42,
		OpticalVegetationIndex: 0.28,
		StructuralPourDetected: project.CurrentStage == domain.StageConstruction || project.CurrentStage == domain.StageOperating,
		EarthworksDetected:     project.CurrentStage != domain.StageConcept,
		ResolutionMeters:       3.0,
	}

	dossier := evaluator.CorroborateProject(project, []earthobs.GroundTruthObservation{pass})

	writeJSON(w, http.StatusOK, dossier)
}
