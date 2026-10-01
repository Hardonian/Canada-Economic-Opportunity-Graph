package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/concession"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/transitionfinance"
)

// TransitionTaxonomyRequest allows caller to specify a custom project or list of project IDs to evaluate.
type TransitionTaxonomyRequest struct {
	ProjectIDs []string        `json:"project_ids,omitempty"`
	Project    *domain.Project `json:"project,omitempty"`
}

// handleTransitionTaxonomy evaluates transition taxonomy credentials for the portfolio or requested projects.
func (s *Server) handleTransitionTaxonomy(w http.ResponseWriter, r *http.Request) {
	engine := transitionfinance.NewTransitionEngine()

	// 1. If project query param provided, evaluate single project
	if id := r.URL.Query().Get("id"); id != "" {
		project, err := s.resolveProject(r.Context(), id)
		if err != nil || project == nil {
			writeError(w, r, http.StatusNotFound, "project_not_found", "Project not found.")
			return
		}
		assessment := engine.EvaluateProject(project)
		writeJSON(w, http.StatusOK, assessment)
		return
	}

	// 2. If POST body with custom project or specific project IDs
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(io.LimitReader(r.Body, s.maxRequestBody))
		if err == nil && len(body) > 0 {
			var req TransitionTaxonomyRequest
			if err := json.Unmarshal(body, &req); err == nil {
				if req.Project != nil {
					assessment := engine.EvaluateProject(req.Project)
					writeJSON(w, http.StatusOK, assessment)
					return
				}
				if len(req.ProjectIDs) > 0 {
					var selected []*domain.Project
					for _, pid := range req.ProjectIDs {
						if p, err := s.resolveProject(r.Context(), pid); err == nil && p != nil {
							selected = append(selected, p)
						}
					}
					if len(selected) > 0 {
						summary := engine.EvaluatePortfolio(selected)
						writeJSON(w, http.StatusOK, summary)
						return
					}
				}
			}
		}
	}

	// 3. Default: Evaluate all indexed projects in the database
	projects, _, err := s.store.ListProjects(r.Context(), database.ProjectFilter{Limit: 500})
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "database_error", "Failed to retrieve projects.")
		return
	}

	summary := engine.EvaluatePortfolio(projects)
	writeJSON(w, http.StatusOK, summary)
}

// handleTransitionTaxonomyProject evaluates a specific project by URL path parameter.
func (s *Server) handleTransitionTaxonomyProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := s.resolveProject(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, r, http.StatusNotFound, "project_not_found", "Project not found.")
		return
	}

	engine := transitionfinance.NewTransitionEngine()
	assessment := engine.EvaluateProject(project)
	writeJSON(w, http.StatusOK, assessment)
}

// handleAirportConcessions simulates the National Airports System (NAS) ground lease modernization.
func (s *Server) handleAirportConcessions(w http.ResponseWriter, r *http.Request) {
	params := concession.DefaultSimulationParams()

	if qYears := r.URL.Query().Get("years"); qYears != "" {
		if val, err := strconv.Atoi(qYears); err == nil && val >= 20 && val <= 60 {
			params.ConcessionHorizonYears = val
		}
	}
	if qRoyalty := r.URL.Query().Get("royalty_pct"); qRoyalty != "" {
		if val, err := strconv.ParseFloat(qRoyalty, 64); err == nil && val > 0 && val <= 30 {
			params.FederalRoyaltyRatePercent = val
		}
	}
	if qPension := r.URL.Query().Get("pension_equity_pct"); qPension != "" {
		if val, err := strconv.ParseFloat(qPension, 64); err == nil && val > 0 && val <= 90 {
			params.PensionEquitySharePercent = val
		}
	}
	if qDebt := r.URL.Query().Get("debt_pct"); qDebt != "" {
		if val, err := strconv.ParseFloat(qDebt, 64); err == nil && val > 0 && val <= 80 {
			params.CommercialDebtSharePercent = val
		}
	}
	if qGrowth := r.URL.Query().Get("passenger_cagr_pct"); qGrowth != "" {
		if val, err := strconv.ParseFloat(qGrowth, 64); err == nil && val > 0 && val <= 10 {
			params.PassengerCAGRPercent = val
		}
	}
	if qDiscount := r.URL.Query().Get("discount_rate_pct"); qDiscount != "" {
		if val, err := strconv.ParseFloat(qDiscount, 64); err == nil && val > 0 && val <= 15 {
			params.DiscountRatePercent = val
		}
	}

	sim := concession.SimulateConcessions(params)
	writeJSON(w, http.StatusOK, sim)
}

// handleAirportConcessionSimulate executes a simulation with custom parameters supplied via JSON POST.
func (s *Server) handleAirportConcessionSimulate(w http.ResponseWriter, r *http.Request) {
	params := concession.DefaultSimulationParams()

	body, err := io.ReadAll(io.LimitReader(r.Body, s.maxRequestBody))
	if err == nil && len(body) > 0 {
		var inputParams concession.SimulationParams
		if err := json.Unmarshal(body, &inputParams); err == nil {
			if inputParams.ConcessionHorizonYears > 0 {
				params.ConcessionHorizonYears = inputParams.ConcessionHorizonYears
			}
			if inputParams.FederalRoyaltyRatePercent > 0 {
				params.FederalRoyaltyRatePercent = inputParams.FederalRoyaltyRatePercent
			}
			if inputParams.PensionEquitySharePercent > 0 {
				params.PensionEquitySharePercent = inputParams.PensionEquitySharePercent
			}
			if inputParams.CommercialDebtSharePercent > 0 {
				params.CommercialDebtSharePercent = inputParams.CommercialDebtSharePercent
			}
			if inputParams.FederalSubordinatedSharePct > 0 {
				params.FederalSubordinatedSharePct = inputParams.FederalSubordinatedSharePct
			}
			if inputParams.PassengerCAGRPercent > 0 {
				params.PassengerCAGRPercent = inputParams.PassengerCAGRPercent
			}
			if inputParams.InflationPercent > 0 {
				params.InflationPercent = inputParams.InflationPercent
			}
			if inputParams.DiscountRatePercent > 0 {
				params.DiscountRatePercent = inputParams.DiscountRatePercent
			}
			if len(inputParams.CustomCapexCAD) > 0 {
				params.CustomCapexCAD = inputParams.CustomCapexCAD
			}
		}
	}

	sim := concession.SimulateConcessions(params)
	writeJSON(w, http.StatusOK, sim)
}
