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
}
