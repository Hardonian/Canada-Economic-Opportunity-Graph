package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/concession"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/transitionfinance"
)

func TestCarneyEndpoints(t *testing.T) {
	store := database.NewMemoryStore()
	p1 := &domain.Project{
		ID:        "darlington-smr",
		Slug:      "darlington-smr",
		Name:      "Darlington New Nuclear Project — Unit 1",
		Sector:    domain.SectorNuclearEnergy,
		Subsector: "Small Modular Reactors",
		Province:  "ON",
		CapexCAD:  7_700_000_000,
	}
	p2 := &domain.Project{
		ID:        "nas-airport-hubs",
		Slug:      "canadian-international-airports-global-investor-leasing-infrastructure-hubs",
		Name:      "Canadian International Airports — Global Investor Leasing & Infrastructure Hubs",
		Sector:    domain.SectorTransportation,
		Subsector: "Airport Commercial Ground Leases & Intermodal Logistics",
		Province:  "ON",
		CapexCAD:  18_000_000_000,
	}
	if err := store.SaveProject(t.Context(), p1); err != nil {
		t.Fatalf("Failed to save p1: %v", err)
	}
	if err := store.SaveProject(t.Context(), p2); err != nil {
		t.Fatalf("Failed to save p2: %v", err)
	}

	server := NewServer(store)

	t.Run("GET /api/v1/finance/transition-taxonomy returns portfolio summary", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/finance/transition-taxonomy", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var summary transitionfinance.PortfolioTransitionSummary
		if err := json.NewDecoder(w.Body).Decode(&summary); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if summary.TotalProjects != 2 {
			t.Errorf("Expected 2 projects, got %d", summary.TotalProjects)
		}
		if summary.TotalGreenCapexCAD <= 0 {
			t.Errorf("Expected positive green capex, got %d", summary.TotalGreenCapexCAD)
		}
		if summary.WeightedAverageCredibilityIndex <= 0 {
			t.Errorf("Expected positive credibility, got %.1f", summary.WeightedAverageCredibilityIndex)
		}
	})

	t.Run("GET /api/v1/finance/transition-taxonomy/{id} returns single project assessment", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/finance/transition-taxonomy/darlington-smr", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var assessment transitionfinance.TransitionAssessment
		if err := json.NewDecoder(w.Body).Decode(&assessment); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if assessment.TaxonomyCategory != transitionfinance.CategoryGreen {
			t.Errorf("Expected CategoryGreen, got %s", assessment.TaxonomyCategory)
		}
		if assessment.CredibilityIndex.OverallScore < 80 {
			t.Errorf("Expected credibility >= 80, got %.1f", assessment.CredibilityIndex.OverallScore)
		}
	})

	t.Run("POST /api/v1/finance/transition-taxonomy evaluates custom project", func(t *testing.T) {
		custom := &domain.Project{
			ID:        "crawford-custom",
			Name:      "Crawford Nickel Carbon Capture Project",
			Summary:   "carbon capture and ultramafic direct mineral sequestration",
			Sector:    domain.SectorCriticalMinerals,
			Subsector: "Nickel Mining & CCUS",
			CapexCAD:  3_500_000_000,
		}
		reqBody, _ := json.Marshal(TransitionTaxonomyRequest{Project: custom})
		req := httptest.NewRequest("POST", "/api/v1/finance/transition-taxonomy", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var assessment transitionfinance.TransitionAssessment
		if err := json.NewDecoder(w.Body).Decode(&assessment); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if assessment.TaxonomyCategory != transitionfinance.CategoryTransition {
			t.Errorf("Expected CategoryTransition, got %s", assessment.TaxonomyCategory)
		}
	})

	t.Run("GET /api/v1/finance/concession/airports returns NAS simulation", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/finance/concession/airports", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var sim concession.NationalConcessionSimulation
		if err := json.NewDecoder(w.Body).Decode(&sim); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if sim.TotalNASModernizationCapexCAD != 18_000_000_000 {
			t.Errorf("Expected 18B capex, got %d", sim.TotalNASModernizationCapexCAD)
		}
		if len(sim.AirportResults) != 5 {
			t.Errorf("Expected 5 airports, got %d", len(sim.AirportResults))
		}
		if sim.CrowdingInMultiplier <= 3.0 {
			t.Errorf("Expected crowding in > 3.0, got %.1f", sim.CrowdingInMultiplier)
		}
	})

	t.Run("POST /api/v1/finance/concession/airports/simulate applies custom inputs", func(t *testing.T) {
		customParams := concession.SimulationParams{
			ConcessionHorizonYears:    50,
			FederalRoyaltyRatePercent: 12.0,
			PensionEquitySharePercent: 55.0,
		}
		reqBody, _ := json.Marshal(customParams)
		req := httptest.NewRequest("POST", "/api/v1/finance/concession/airports/simulate", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var sim concession.NationalConcessionSimulation
		if err := json.NewDecoder(w.Body).Decode(&sim); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if sim.Params.ConcessionHorizonYears != 50 {
			t.Errorf("Expected 50 year horizon, got %d", sim.Params.ConcessionHorizonYears)
		}
		if sim.Params.FederalRoyaltyRatePercent != 12.0 {
			t.Errorf("Expected 12%% royalty, got %.1f", sim.Params.FederalRoyaltyRatePercent)
		}
	})
}
