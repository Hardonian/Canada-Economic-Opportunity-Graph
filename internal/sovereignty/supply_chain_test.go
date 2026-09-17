package sovereignty

import (
	"testing"
)

func TestSupplyChainEvaluator_LithiumIntegrated(t *testing.T) {
	evaluator := NewSupplyChainEvaluator()

	// Integrated Canadian lithium extraction + cathode refining project
	project := MineralProjectProfile{
		ProjectID:             "li-nemaska-whabouchi",
		ProjectName:           "Nemaska Whabouchi & Becancour Hydroxide",
		TargetMineral:         "LI",
		HasExtraction:         true,
		HasDomesticRefining:   true,
		HasEndProductMfg:      true,
		AnnualCapacityTonnes:  34000,
		AlliedOfftakePct:      100.0,
		IndigenousEquityShare: 0.10,
	}

	res := evaluator.Evaluate(project)

	if res.SelfSufficiencyIndex < 90.0 {
		t.Errorf("expected high self sufficiency index for fully integrated project, got %.1f", res.SelfSufficiencyIndex)
	}

	if res.ExtractionScore != 100.0 {
		t.Errorf("expected extraction score 100.0, got %.1f", res.ExtractionScore)
	}

	if res.RefiningScore != 100.0 {
		t.Errorf("expected refining score 100.0, got %.1f", res.RefiningScore)
	}

	if res.ChokePointMitigation != 100.0 {
		t.Errorf("expected choke point mitigation 100.0, got %.1f", res.ChokePointMitigation)
	}

	// Should not have RAW_EXPORT_RISK vulnerability
	for _, v := range res.Vulnerabilities {
		if v.Category == "RAW_EXPORT_RISK" {
			t.Errorf("did not expect RAW_EXPORT_RISK for integrated project: %+v", v)
		}
	}
}

func TestSupplyChainEvaluator_RawExportVulnerability(t *testing.T) {
	evaluator := NewSupplyChainEvaluator()

	// Project that digs graphite in Quebec/Ontario but ships raw concentrate overseas without refining
	project := MineralProjectProfile{
		ProjectID:             "graphite-unrefined",
		ProjectName:           "Northern Graphite Concentrate Only",
		TargetMineral:         "GRAPHITE",
		HasExtraction:         true,
		HasDomesticRefining:   false,
		HasEndProductMfg:      false,
		AnnualCapacityTonnes:  15000,
		AlliedOfftakePct:      50.0,
		IndigenousEquityShare: 0.0,
	}

	res := evaluator.Evaluate(project)

	if res.SelfSufficiencyIndex > 65.0 {
		t.Errorf("expected lower index due to lack of domestic refining, got %.1f", res.SelfSufficiencyIndex)
	}

	foundExportRisk := false
	foundChokePoint := false
	for _, v := range res.Vulnerabilities {
		if v.Category == "RAW_EXPORT_RISK" {
			foundExportRisk = true
		}
		if v.Category == "FOREIGN_CHOKEPOINT" {
			foundChokePoint = true
		}
	}

	if !foundExportRisk {
		t.Error("expected RAW_EXPORT_RISK vulnerability to be flagged")
	}
	if !foundChokePoint {
		t.Error("expected FOREIGN_CHOKEPOINT vulnerability to be flagged for unrefined graphite")
	}
}

func TestSupplyChainEvaluator_UnknownMineral(t *testing.T) {
	evaluator := NewSupplyChainEvaluator()

	project := MineralProjectProfile{
		ProjectID:           "unknown-mineral-1",
		ProjectName:         "Custom Scandium Project",
		TargetMineral:       "SCANDIUM",
		HasExtraction:       true,
		HasDomesticRefining: true,
	}

	res := evaluator.Evaluate(project)
	if res.TargetMineral != "SCANDIUM" {
		t.Errorf("expected target mineral SCANDIUM, got %s", res.TargetMineral)
	}
	if res.SelfSufficiencyIndex <= 0.0 {
		t.Errorf("expected positive self sufficiency index, got %.1f", res.SelfSufficiencyIndex)
	}
}
