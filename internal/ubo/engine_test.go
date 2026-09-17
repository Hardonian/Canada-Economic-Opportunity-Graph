package ubo

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestUBOSovereignScreening(t *testing.T) {
	evaluator := NewEvaluator()

	t.Run("Pure domestic project is CLEAR", func(t *testing.T) {
		project := &domain.Project{
			ID:          "p-domestic",
			ProponentID: "e-hydro",
			Sector:      domain.SectorCleanEnergy,
		}
		res := evaluator.ScreenProject(project, nil)
		if res.ICARisk != RiskClear {
			t.Fatalf("expected RiskClear, got %s", res.ICARisk)
		}
		if res.DomesticControlShare != 1.0 {
			t.Fatalf("expected 1.0 domestic, got %.2f", res.DomesticControlShare)
		}
	})

	t.Run("Non-FTA SOE in critical minerals is PROHIBITED", func(t *testing.T) {
		project := &domain.Project{
			ID:          "p-lithium",
			ProponentID: "e-junior-mine",
			Sector:      domain.SectorCriticalMinerals,
			Subsector:   "Lithium Spodumene",
		}
		owners := []BeneficialOwner{
			{Name: "Canadian Founders", Jurisdiction: "CA", OwnershipPercent: 0.70, IsSOE: false},
			{Name: "Sinomine Mining Corp", Jurisdiction: "CN", OwnershipPercent: 0.30, IsSOE: true},
		}
		res := evaluator.ScreenProject(project, owners)
		if res.ICARisk != RiskProhibited {
			t.Fatalf("expected RiskProhibited for non-FTA SOE in critical minerals, got %s", res.ICARisk)
		}
		if len(res.NationalSecurityNotes) == 0 {
			t.Fatal("expected national security notes")
		}
	})

	t.Run("Allied FTA investment is CLEAR", func(t *testing.T) {
		project := &domain.Project{
			ID:          "p-nickel",
			ProponentID: "e-nickel-co",
			Sector:      domain.SectorCriticalMinerals,
		}
		owners := []BeneficialOwner{
			{Name: "Toronto Holdings", Jurisdiction: "CA", OwnershipPercent: 0.60},
			{Name: "US Critical Minerals Fund", Jurisdiction: "US", OwnershipPercent: 0.40},
		}
		res := evaluator.ScreenProject(project, owners)
		if res.ICARisk != RiskClear {
			t.Fatalf("expected RiskClear for US allied investment, got %s", res.ICARisk)
		}
	})

	t.Run("UBOUnraveler traces through offshore secrecy havens", func(t *testing.T) {
		unraveler := NewUBOUnraveler()
		holdings := []HoldingNode{
			{EntityID: "hold-1", Name: "Canadian Asset Corp", Jurisdiction: "CA", DirectShare: 1.0, ParentID: "proj-1"},
			{EntityID: "hold-2", Name: "Cayman HoldCo", Jurisdiction: "KY", DirectShare: 0.8, ParentID: "hold-1"},
			{EntityID: "hold-3", Name: "BVI Nominee", Jurisdiction: "VG", DirectShare: 0.5, ParentID: "hold-2"},
			{EntityID: "owner-true", Name: "Zhaojin Gold Group", Jurisdiction: "CN", DirectShare: 1.0, ParentID: "hold-3"},
		}

		unraveled := unraveler.UnravelChain("proj-1", holdings)
		if len(unraveled) == 0 {
			t.Fatalf("expected unraveled leaf owner, got 0")
		}
		top := unraveled[0]
		if top.UltimateEntityID != "owner-true" {
			t.Errorf("expected true owner Zhaojin, got %s", top.Name)
		}
		if !top.IsOffshoreObfuscated {
			t.Errorf("expected offshore obfuscation to be flagged")
		}
		if top.OffshoreSecrecyHops < 2 {
			t.Errorf("expected at least 2 secrecy hops (KY, VG), got %d", top.OffshoreSecrecyHops)
		}
		expectedShare := 0.8 * 0.5 * 1.0
		if top.EffectiveShare != expectedShare {
			t.Errorf("expected effective share %f, got %f", expectedShare, top.EffectiveShare)
		}
	})

	t.Run("SOEDetector flags foreign state entities", func(t *testing.T) {
		detector := NewSOEDetector()
		isSOE, state, risk := detector.DetectSOE("SASAC Overseas Investment Co", "CN")
		if !isSOE {
			t.Errorf("expected SASAC to be detected as SOE")
		}
		if state != "CHINA" {
			t.Errorf("expected CHINA controlling state, got %s", state)
		}
		if risk < 90.0 {
			t.Errorf("expected high risk, got %f", risk)
		}
	})
}
