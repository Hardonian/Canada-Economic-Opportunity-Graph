package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
)

func TestPalantirGradeEndpoints(t *testing.T) {
	store := database.NewMemoryStore()
	server := NewServer(store)

	endpoints := []struct {
		name           string
		method         string
		path           string
		body           map[string]interface{}
		expectedStatus int
	}{
		// Pillar 1: Ontology
		{
			name:           "POST /api/v1/ontology/actions/execute",
			method:         http.MethodPost,
			path:           "/api/v1/ontology/actions/execute",
			body:           map[string]interface{}{"action_type": "UPDATE_PROJECT_STAGE", "target_node_id": "crawford-nickel"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/v1/ontology/branches",
			method:         http.MethodGet,
			path:           "/api/v1/ontology/branches",
			expectedStatus: http.StatusOK,
		},
		// Pillar 2: Lakehouse
		{
			name:           "GET /api/v1/lakehouse/query",
			method:         http.MethodGet,
			path:           "/api/v1/lakehouse/query",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/v1/lakehouse/quality",
			method:         http.MethodGet,
			path:           "/api/v1/lakehouse/quality",
			expectedStatus: http.StatusOK,
		},
		// Pillar 3: Graph & GQL
		{
			name:           "POST /api/v1/graph/gql",
			method:         http.MethodPost,
			path:           "/api/v1/graph/gql",
			body:           map[string]interface{}{"query": "MATCH (p:Project) RETURN p"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/v1/graph/gnn-predictions",
			method:         http.MethodGet,
			path:           "/api/v1/graph/gnn-predictions",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/v1/graph/communities",
			method:         http.MethodGet,
			path:           "/api/v1/graph/communities",
			expectedStatus: http.StatusOK,
		},
		// Pillar 4: Earth Observation
		{
			name:           "GET /api/v1/earthobs/insar/crawford-nickel",
			method:         http.MethodGet,
			path:           "/api/v1/earthobs/insar/crawford-nickel",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/v1/earthobs/ais/congestion",
			method:         http.MethodGet,
			path:           "/api/v1/earthobs/ais/congestion",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/v1/earthobs/discrepancies",
			method:         http.MethodGet,
			path:           "/api/v1/earthobs/discrepancies",
			expectedStatus: http.StatusOK,
		},
		// Pillar 5: AIP Multi-Agent
		{
			name:           "POST /api/v1/ai/copilot/run",
			method:         http.MethodPost,
			path:           "/api/v1/ai/copilot/run",
			body:           map[string]interface{}{"role": "FINANCE", "query": "Assess IRR"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /api/v1/ai/negotiate/ppa",
			method:         http.MethodPost,
			path:           "/api/v1/ai/negotiate/ppa",
			expectedStatus: http.StatusOK,
		},
		// Pillar 7: Sovereign Security & ABAC
		{
			name:           "POST /api/v1/security/authorize",
			method:         http.MethodPost,
			path:           "/api/v1/security/authorize",
			body:           map[string]interface{}{"subject_clearance": "SECRET", "object_classification": "PROTECTED_B"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /api/v1/security/dlp/scan",
			method:         http.MethodPost,
			path:           "/api/v1/security/dlp/scan",
			body:           map[string]interface{}{"content": "Secret SIN 123-456-789"},
			expectedStatus: http.StatusOK,
		},
		// Pillar 8: Counter-Intelligence & UBO
		{
			name:           "GET /api/v1/counter-intel/ubo/trace",
			method:         http.MethodGet,
			path:           "/api/v1/counter-intel/ubo/trace",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/v1/counter-intel/dark-fleet",
			method:         http.MethodGet,
			path:           "/api/v1/counter-intel/dark-fleet",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /api/v1/counter-intel/ica/screen",
			method:         http.MethodPost,
			path:           "/api/v1/counter-intel/ica/screen",
			expectedStatus: http.StatusOK,
		},
		// Pillar 10: Compliance & Merkle
		{
			name:           "GET /api/v1/compliance/c59/audit/crawford-nickel",
			method:         http.MethodGet,
			path:           "/api/v1/compliance/c59/audit/crawford-nickel",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/v1/compliance/c69/clock/crawford-nickel",
			method:         http.MethodGet,
			path:           "/api/v1/compliance/c69/clock/crawford-nickel",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/v1/compliance/trc92/scorecard/crawford-nickel",
			method:         http.MethodGet,
			path:           "/api/v1/compliance/trc92/scorecard/crawford-nickel",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/v1/compliance/merkle/root",
			method:         http.MethodGet,
			path:           "/api/v1/compliance/merkle/root",
			expectedStatus: http.StatusOK,
		},
	}

	for _, ep := range endpoints {
		t.Run(ep.name, func(t *testing.T) {
			var req *http.Request
			if ep.body != nil {
				b, _ := json.Marshal(ep.body)
				req = httptest.NewRequest(ep.method, ep.path, bytes.NewReader(b))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(ep.method, ep.path, nil)
			}

			rr := httptest.NewRecorder()
			server.ServeHTTP(rr, req)

			if rr.Code != ep.expectedStatus {
				t.Fatalf("endpoint %s failed: expected status %d, got %d. Body: %s",
					ep.path, ep.expectedStatus, rr.Code, rr.Body.String())
			}
		})
	}
}
