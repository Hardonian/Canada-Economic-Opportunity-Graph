package earthobs

import (
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestEarthObsGroundTruthCorroboration(t *testing.T) {
	evaluator := NewEvaluator()

	t.Run("Construction stage verified by SAR and optical structures", func(t *testing.T) {
		project := &domain.Project{
			ID:           "p-oneida",
			CurrentStage: domain.StageConstruction,
		}
		obs := []GroundTruthObservation{
			{
				Constellation:          ConstellationRCM,
				PassDate:               time.Now().Add(-48 * time.Hour),
				StructuralRadarCrossSec: -8.5, // High backscatter
				EarthworksDetected:     true,
				StructuralPourDetected: true,
			},
		}
		dossier := evaluator.CorroborateProject(project, obs)
		if dossier.CorroborationStatus != domain.ConfidenceVerified {
			t.Fatalf("Expected ConfidenceVerified, got %s", dossier.CorroborationStatus)
		}
		if !dossier.StructuresConfirmed {
			t.Fatal("Expected structures confirmed")
		}
	})

	t.Run("Discrepancy flagged if claimed construction has undisturbed terrain", func(t *testing.T) {
		project := &domain.Project{
			ID:           "p-phantom",
			CurrentStage: domain.StageConstruction,
		}
		obs := []GroundTruthObservation{
			{
				Constellation:          ConstellationSentinel2,
				PassDate:               time.Now().Add(-24 * time.Hour),
				OpticalVegetationIndex: 0.85, // Pristine dense vegetation
				EarthworksDetected:     false,
				StructuralPourDetected: false,
			},
		}
		dossier := evaluator.CorroborateProject(project, obs)
		if dossier.CorroborationStatus != domain.ConfidenceConflict {
			t.Fatalf("Expected ConfidenceConflict for phantom construction, got %s", dossier.CorroborationStatus)
		}
		if dossier.PhysicalProgressScore > 30.0 {
			t.Fatalf("Expected low progress score, got %.1f", dossier.PhysicalProgressScore)
		}
	})

	t.Run("InSAR Coherence and Tailings Dam Stability", func(t *testing.T) {
		sp := NewSARPipeline()
		mDate := time.Now().Add(-12 * 24 * time.Hour)
		sDate := time.Now()

		insar := sp.ComputeInSAR(mDate, sDate, 0.5, 0.45, 0.12)
		if insar.TemporalBaselineDays != 12 {
			t.Errorf("expected 12 baseline days, got %d", insar.TemporalBaselineDays)
		}
		if insar.TailingsStabilityRating != "STABLE" {
			t.Errorf("expected STABLE, got %s", insar.TailingsStabilityRating)
		}
	})

	t.Run("Computer Vision Machinery Detection", func(t *testing.T) {
		cve := NewCVInferenceEngine()
		mach := cve.DetectMachinery("proj-darlington", 3400000000, 0.1, 0.35, 14.5)
		if mach.TotalMachineryUnits == 0 {
			t.Fatalf("expected heavy equipment detected on active megaproject")
		}
		if mach.ConstructionActivityRate == "DORMANT" {
			t.Errorf("expected active rate, got %s", mach.ConstructionActivityRate)
		}
	})

	t.Run("AIS Port Congestion Telemetry", func(t *testing.T) {
		ate := NewAISTrackingEngine()
		vessels := []AISVessel{
			{MMSI: "123456", VesselName: "Vancouver Star", Status: "ANCHORED", SpeedKnots: 0.1, ArrivalDate: time.Now().Add(-72 * time.Hour)},
			{MMSI: "234567", VesselName: "Pacific Leader", Status: "MOORED", SpeedKnots: 0.0, ArrivalDate: time.Now().Add(-24 * time.Hour)},
			{MMSI: "345678", VesselName: "Atlantic Express", Status: "ANCHORED", SpeedKnots: 0.2, ArrivalDate: time.Now().Add(-96 * time.Hour)},
		}

		telemetry := ate.EvaluatePort("CAVAN", "Port of Vancouver", vessels)
		if telemetry.AnchoredVesselCount != 2 {
			t.Errorf("expected 2 anchored vessels, got %d", telemetry.AnchoredVesselCount)
		}
		if telemetry.AverageDwellDays <= 0 {
			t.Errorf("expected positive average dwell days")
		}
	})

	t.Run("Ground Truth Discrepancy Detection", func(t *testing.T) {
		gte := NewGroundTruthDiscrepancyEngine()
		cve := NewCVInferenceEngine()
		mach := cve.DetectMachinery("p-scam", 500000000, 0.8, -0.1, 2.0) // Very low activity

		rep := gte.EvaluateDiscrepancy("p-scam", 85.0, mach, nil) // Claims 85% progress
		if !rep.IsDiscrepancyFlagged {
			t.Errorf("expected discrepancy to be flagged for severe divergence")
		}
		if rep.RiskVerdict != "MATERIAL_DECEPTION_ALERT" {
			t.Errorf("expected MATERIAL_DECEPTION_ALERT, got %s", rep.RiskVerdict)
		}
	})
}
