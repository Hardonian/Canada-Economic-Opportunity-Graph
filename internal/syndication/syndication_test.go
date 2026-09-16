package syndication

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestMatcherMatchProject(t *testing.T) {
	matcher := NewMatcher(nil) // uses default investors
	proj := &domain.Project{
		ID:           "test-nuclear-smr",
		Name:         "Darlington SMR Unit 1",
		Sector:       domain.SectorNuclearEnergy,
		CurrentStage: domain.StageConstruction,
		CapexCAD:     2_500_000_000,
	}

	consortium := matcher.MatchProject(proj)
	if consortium == nil {
		t.Fatal("expected non-nil consortium")
	}
	if len(consortium.Matches) == 0 {
		t.Fatal("expected investor matches")
	}
	if consortium.PrivateCrowdingInRatio <= 0 {
		t.Errorf("expected positive crowding-in ratio, got %f", consortium.PrivateCrowdingInRatio)
	}
	if consortium.CrownConcessionCAD <= 0 {
		t.Errorf("expected non-zero Crown concessionary allocation")
	}
	if consortium.AuditHash == "" {
		t.Errorf("expected non-empty audit hash")
	}
}

func TestCanonicalOfftakes(t *testing.T) {
	offtakes := GetProjectOfftakes("darlington-small-modular-reactor-deployment-project")
	if len(offtakes) == 0 {
		t.Fatal("expected offtakes for darlington SMR")
	}
	o := offtakes[0]
	if o.BuyerName == "" {
		t.Errorf("expected non-empty buyer name")
	}
	if o.AnnualContractValueCAD <= 0 {
		t.Errorf("expected positive annual contract value")
	}
	if !o.TakeOrPayObligation {
		t.Errorf("expected take or pay obligation to be true")
	}
}

func TestMultiNationSyndicate(t *testing.T) {
	communities := []struct {
		Name      string
		Territory string
		KM        float64
	}{
		{Name: "Marten Falls First Nation", Territory: "Treaty 9", KM: 120.0},
		{Name: "Webequie First Nation", Territory: "Treaty 9", KM: 180.0},
		{Name: "Aroland First Nation", Territory: "Treaty 9", KM: 100.0},
	}

	syndicate := BuildMultiNationSyndicate("Ring of Fire Northern Road Corridor", 300_000_000, communities)
	if syndicate == nil {
		t.Fatal("expected non-nil syndicate")
	}
	if len(syndicate.Participants) != 3 {
		t.Fatalf("expected 3 participants, got %d", len(syndicate.Participants))
	}
	if syndicate.TotalCorridorKM != 400.0 {
		t.Errorf("expected 400km total, got %f", syndicate.TotalCorridorKM)
	}
	if syndicate.FederalILGPGreaterCAD != 285_000_000 {
		t.Errorf("expected 285M CAD guarantee, got %d", syndicate.FederalILGPGreaterCAD)
	}
	if syndicate.AuditHash == "" {
		t.Errorf("expected non-empty audit hash")
	}
}
