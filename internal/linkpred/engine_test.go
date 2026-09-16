package linkpred

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestLinkPredictor(t *testing.T) {
	lp := NewLinkPredictor()

	target := &domain.Project{
		ID:           "p-darlington",
		Name:         "Darlington SMR Project",
		Sector:       domain.SectorNuclearEnergy,
		Province:     "ON",
		ProponentID:  "e-opg",
	}

	entities := []*domain.Entity{
		{
			ID:           "e-aecon",
			LegalName:    "Aecon Group Inc.",
			CommonName:   "Aecon",
			EntityType:   "Supplier",
			Jurisdiction: "ON",
			Description:  "Nuclear and heavy civil engineering construction",
		},
		{
			ID:           "e-cib",
			LegalName:    "Canada Infrastructure Bank",
			CommonName:   "CIB",
			EntityType:   "Investor",
			Jurisdiction: "Federal",
			Description:  "Concessionary public infrastructure bank",
		},
		{
			ID:           "e-cree",
			LegalName:    "Cree Nation of Chisasibi",
			CommonName:   "Chisasibi First Nation",
			EntityType:   "FirstNation",
			Jurisdiction: "QC",
		},
	}

	projects := []*domain.Project{target}
	relationships := []*domain.Relationship{
		{ProjectID: "p-darlington", SourceEntityID: "e-opg", TargetEntityID: "p-darlington", RelationType: "proponent"},
	}

	predictions := lp.PredictProjectPartners(target, entities, projects, relationships)

	if len(predictions) == 0 {
		t.Fatal("Expected link predictions for candidate entities")
	}

	// Verify Aecon is classified as EPC contractor
	var foundAecon bool
	for _, pred := range predictions {
		if pred.EntityID == "e-aecon" {
			foundAecon = true
			if pred.PredictedRole != RoleEPCContractor {
				t.Fatalf("Expected Aecon to be EPC_CONTRACTOR, got %s", pred.PredictedRole)
			}
			if pred.ConfidenceScore < 0.60 {
				t.Fatalf("Expected high confidence for Aecon in nuclear, got %.2f", pred.ConfidenceScore)
			}
		}
	}
	if !foundAecon {
		t.Fatal("Expected Aecon in link predictions")
	}
}
