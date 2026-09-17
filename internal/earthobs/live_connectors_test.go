package earthobs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCopernicusClient_SearchSentinel1Scenes(t *testing.T) {
	// Mock Copernicus OData API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"value": [
				{
					"Id": "s1a-iw-grd-20260901t113000",
					"Name": "S1A_IW_GRDH_1SDV_20260901T113000_Kitimat_BC",
					"ContentType": "application/octet-stream",
					"ContentLength": 1048576000,
					"OriginDate": "2026-09-01T11:30:00Z"
				}
			]
		}`))
	}))
	defer mockServer.Close()

	client := NewCopernicusClient(mockServer.URL, "test-api-token")
	scenes, err := client.SearchSentinel1Scenes(
		context.Background(),
		53.5, 54.5, -129.0, -128.0, // Kitimat LNG Terminal area
		time.Now().Add(-7*24*time.Hour),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("unexpected copernicus query error: %v", err)
	}
	if len(scenes) != 1 {
		t.Fatalf("expected 1 scene, got %d", len(scenes))
	}
	if scenes[0].ID != "s1a-iw-grd-20260901t113000" {
		t.Errorf("expected scene ID s1a-iw-grd-20260901t113000, got %s", scenes[0].ID)
	}
	if scenes[0].ContentLength != 1048576000 {
		t.Errorf("expected length 1048576000, got %d", scenes[0].ContentLength)
	}
}

func TestLiveAISDecoder_IngestStreamAndLookup(t *testing.T) {
	decoder := NewLiveAISDecoder()

	streamData := `# Canadian Coast Guard Maritime Stream Snapshot
316001234,POLAR_PRINCE,RESEARCH_VESSEL,CA,12.4,UNDERWAY,48.4,-123.3
316005678,PACIFIC_HOPE,TANKER,CA,0.2,ANCHORED,49.2,-123.1
316009999,KITIMAT_SPIRIT,LNG_CARRIER,BS,0.0,MOORED,53.9,-128.6
`

	n, err := decoder.IngestStream(strings.NewReader(streamData))
	if err != nil {
		t.Fatalf("unexpected ingest error: %v", err)
	}
	if n != 3 {
		t.Fatalf("expected 3 vessels decoded, got %d", n)
	}

	v, ok := decoder.GetVessel("316009999")
	if !ok {
		t.Fatal("expected vessel 316009999 to be found")
	}
	if v.VesselName != "KITIMAT_SPIRIT" || v.Status != "MOORED" {
		t.Errorf("unexpected vessel details: %+v", v)
	}

	all := decoder.GetAllVessels()
	if len(all) != 3 {
		t.Fatalf("expected 3 tracked vessels, got %d", len(all))
	}

	// Port evaluation integration test
	engine := NewAISTrackingEngine()
	telemetry := engine.EvaluatePort("port-kitimat", "Port of Kitimat", all)
	if telemetry.MooredVesselCount != 1 || telemetry.AnchoredVesselCount != 1 {
		t.Errorf("unexpected port telemetry: %+v", telemetry)
	}
}
