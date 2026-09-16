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
}
