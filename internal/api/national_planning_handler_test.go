package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/nationalplanning"
)

func TestNationalPlanningEndpoints(t *testing.T) {
	store := database.NewMemoryStore()
	project := &domain.Project{
		ID:           "test-p1",
		Slug:         "test-p1-slug",
		Name:         "Darlington SMR Facility",
		Sector:       domain.SectorNuclearEnergy,
		Province:     "ON",
		CapexCAD:     7_700_000_000,
		CurrentStage: domain.StageConstruction,
		Scores:       map[string]float64{"strategicity": 95.0, "buildability": 85.0},
	}
	if err := store.SaveProject(t.Context(), project); err != nil {
		t.Fatalf("Failed to save project: %v", err)
	}

	server := NewServer(store)

	t.Run("POST /api/v1/planning/optimize returns optimized portfolio", func(t *testing.T) {
		reqBody, _ := json.Marshal(nationalplanning.OptimizationRequest{
			Objective: nationalplanning.ObjectiveMaxCrowdingIn,
		})
		req := httptest.NewRequest("POST", "/api/v1/planning/optimize", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}

		var res nationalplanning.OptimizationResult
		if err := json.NewDecoder(w.Body).Decode(&res); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if len(res.AllocatedProjects) == 0 {
			t.Fatal("Expected allocated projects")
		}
	})

	t.Run("POST /api/v1/planning/wargame runs scenario shock", func(t *testing.T) {
		reqBody, _ := json.Marshal(nationalplanning.WarGameRequest{
			Scenario: nationalplanning.ShockUSMCATariffs,
		})
		req := httptest.NewRequest("POST", "/api/v1/planning/wargame", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("GET /api/v1/planning/labor returns regional report", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/planning/labor?province=ON", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("GET /api/v1/projects/{id}/mrio returns input-output multipliers", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/projects/test-p1/mrio", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("GET /api/v1/projects/{id}/flyvbjerg returns overrun hazard curve", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/projects/test-p1/flyvbjerg", nil)
		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d", w.Code)
		}
	})
}
