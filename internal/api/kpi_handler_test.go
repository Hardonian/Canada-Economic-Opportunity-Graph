package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestKPIEndpoints(t *testing.T) {
	store := database.NewMemoryStore()

	// Seed a test project
	proj := &domain.Project{
		ID:           "test-darlington-smr",
		Slug:         "darlington-new-nuclear-project-unit-1",
		Name:         "Darlington New Nuclear Project (Unit 1)",
		Sector:       domain.SectorNuclearEnergy,
		Province:     "ON",
		CurrentStage: domain.StageConstruction,
		CapexCAD:     3_400_000_000,
		Confidence:   domain.ConfidenceVerified,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := store.SaveProject(context.Background(), proj); err != nil {
		t.Fatalf("failed to seed test project: %v", err)
	}

	server := mustServer(t, store, testOptions())

	// 1. GET /api/v1/kpis
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/kpis", nil)
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/v1/kpis status = %d: %s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, "definitions") || !strings.Contains(body, "ESG.GHG.INTENSITY.SCOPE1_2") {
			t.Fatalf("expected KPI definitions in response: %s", body)
		}
	}

	// 2. GET /api/v1/kpis?category=INDIGENOUS_EQUITY
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/kpis?category=INDIGENOUS_EQUITY", nil)
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/v1/kpis?category=INDIGENOUS_EQUITY status = %d: %s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, "INDIG.EQUITY.OWNERSHIP.PCT") {
			t.Fatalf("expected indigenous KPI in filtered response: %s", body)
		}
	}

	// 3. GET /api/v1/kpis/feeds
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/kpis/feeds?limit=10", nil)
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/v1/kpis/feeds status = %d: %s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, "COMMODITY.WCS_WTI.DIFF.USD_BBL") || !strings.Contains(body, "latest_ticks") {
			t.Fatalf("expected commodity ticks in feeds response: %s", body)
		}
	}

	// 4. GET /api/v1/kpis/project/{id}
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/kpis/project/darlington-new-nuclear-project-unit-1", nil)
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/v1/kpis/project/{id} status = %d: %s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, "overall_kpi_rating") || !strings.Contains(body, "ESG_DECARBONIZATION") {
			t.Fatalf("expected 8-pillar scorecard in project response: %s", body)
		}
	}

	// 5. GET /api/v1/kpis/summary
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/kpis/summary", nil)
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/v1/kpis/summary status = %d: %s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, "national_emissions_abatement_mt") {
			t.Fatalf("expected national summary in response: %s", body)
		}
	}

	// 6. GET /api/v1/kpis/snapshot
	{
		req := httptest.NewRequest(http.MethodGet, "/api/v1/kpis/snapshot", nil)
		rec := httptest.NewRecorder()
		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("GET /api/v1/kpis/snapshot status = %d: %s", rec.Code, rec.Body.String())
		}
		body := rec.Body.String()
		if !strings.Contains(body, "audit_hash") || !strings.Contains(body, "cegs-kpi-v1.0") {
			t.Fatalf("expected snapshot with audit hash in response: %s", body)
		}
	}
}
