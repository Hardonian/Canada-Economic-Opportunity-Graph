package corridor

import (
	"testing"
)

func TestRoutingEngineEvaluateCorridor(t *testing.T) {
	engine := NewRoutingEngine()

	origin := CorridorPoint{
		Name:            "Sudbury Clean Energy Hub",
		Latitude:        46.49,
		Longitude:       -81.01,
		ElevationMeters: 300,
	}
	dest := CorridorPoint{
		Name:            "James Bay Hydro Intertie",
		Latitude:        53.63,
		Longitude:       -77.70,
		ElevationMeters: 150,
	}

	eval := engine.EvaluateCorridor(origin, dest, TypeHVDCTransmission)
	if eval == nil {
		t.Fatal("expected non-nil route evaluation")
	}
	if eval.TotalLengthKM <= 0 {
		t.Errorf("expected positive length, got %f", eval.TotalLengthKM)
	}
	if eval.TotalEstimatedCapexCAD <= 0 {
		t.Errorf("expected positive CAPEX, got %d", eval.TotalEstimatedCapexCAD)
	}
	if eval.Sensitivity.PermafrostThawHazardScore <= 0 {
		t.Errorf("expected positive permafrost risk for Northern route")
	}
	if len(eval.KeyWaypoints) != 3 {
		t.Errorf("expected 3 waypoints, got %d", len(eval.KeyWaypoints))
	}
	if eval.AuditHash == "" {
		t.Errorf("expected non-empty audit hash")
	}
}

func TestCanonicalGateways(t *testing.T) {
	gateways := CanonicalGateways()
	if len(gateways) != 4 {
		t.Fatalf("expected 4 gateways, got %d", len(gateways))
	}

	for _, g := range gateways {
		if g.PortName == "" {
			t.Errorf("expected non-empty port name")
		}
		if g.AnnualThroughputMNTonnes <= 0 {
			t.Errorf("expected positive throughput for %s", g.PortName)
		}
		if g.AuditHash == "" {
			t.Errorf("expected valid audit hash for %s", g.PortName)
		}
	}
}
