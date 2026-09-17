package indigenouslinker

import (
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestOCAPComplianceEngine_ConsentLifecycle(t *testing.T) {
	engine := NewOCAPComplianceEngine()

	business := &domain.Entity{
		ID:           "ent-indig-01",
		EntityType:   "IndigenousBusiness",
		LegalName:    "Cree Nation Engineering & Environmental Ltd.",
		Jurisdiction: "Grand Council of the Crees (Eeyou Istchee)",
	}

	project := &domain.Project{
		ID:       "proj-james-bay-transmission",
		Name:     "James Bay Hydro Transmission Line Upgrade",
		Province: "QC",
	}

	// 1. Unregistered default should permit public procurement matching
	res1 := engine.ValidateLink(business, project, "PROCUREMENT_MATCHING")
	if !res1.Compliant {
		t.Fatalf("expected default compliance for public business, got: %s", res1.ViolationReason)
	}

	// 2. Register explicit consent with restricted use cases
	consent := &CommunityConsent{
		NationOrBandName:  "Grand Council of the Crees (Eeyou Istchee)",
		Status:            ConsentGranted,
		PermittedUseCases: []string{"PROCUREMENT_MATCHING"},
		AuthorizedSigner:  "Grand Chief Office",
		EffectiveDate:     time.Now().Add(-30 * 24 * time.Hour),
		ExpiryDate:        time.Now().Add(365 * 24 * time.Hour),
	}
	engine.RegisterConsent(consent)

	// Permitted use case should pass
	res2 := engine.ValidateLink(business, project, "PROCUREMENT_MATCHING")
	if !res2.Compliant {
		t.Fatalf("expected compliance for permitted use case, got: %s", res2.ViolationReason)
	}

	// Unauthorized use case should fail under OCAP
	res3 := engine.ValidateLink(business, project, "UNAUTHORIZED_DATA_MONETIZATION")
	if res3.Compliant {
		t.Fatal("expected non-compliance for unauthorized data monetization use case")
	}

	// 3. Withdrawn consent should immediately block linking
	consent.Status = ConsentWithdrawn
	res4 := engine.ValidateLink(business, project, "PROCUREMENT_MATCHING")
	if res4.Compliant {
		t.Fatal("expected non-compliance when consent is withdrawn")
	}
}
