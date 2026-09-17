package econometrics

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestSUTCalibration_CoreSovereignSectors(t *testing.T) {
	keySectors := []domain.Sector{
		domain.SectorCriticalMinerals,
		domain.SectorMiningMetals,
		domain.SectorNuclearEnergy,
		domain.SectorCleanEnergy,
		domain.SectorAICompute,
		domain.SectorTransportation,
	}

	cert := GenerateSUTCalibrationCertificate(keySectors)
	if cert.TotalSectors != len(keySectors) {
		t.Fatalf("expected %d sectors evaluated, got %d", len(keySectors), cert.TotalSectors)
	}

	if cert.PassedSectors != cert.TotalSectors {
		for _, r := range cert.SectorReports {
			if !r.EmpiricallyValid {
				t.Errorf("sector %s failed calibration: %v", r.Sector, r.ValidationNotes)
			}
		}
	}

	if cert.IntegrityHash == "" {
		t.Fatal("expected non-empty cryptographic integrity hash on calibration certificate")
	}
}
