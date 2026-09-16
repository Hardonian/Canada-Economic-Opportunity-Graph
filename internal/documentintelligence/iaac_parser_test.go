package documentintelligence

import (
	"testing"
)

func TestExtractIAACDecisionStatement(t *testing.T) {
	t.Run("Extracts conditions and evaluates low injunction risk with signed IBA", func(t *testing.T) {
		fixture := `Project: Marathon Palladium-Copper Mine
Minister: Minister of Environment and Climate Change
Decision: APPROVAL_WITH_CONDITIONS
Conditions: 268
Indigenous Conditions: 42
Aquatic Conditions: 35
Wildlife Conditions: 18
Financial Assurance: $45M
Notes: Formal Impact Benefit Agreement signed with Biigtigong Nishnaabeg First Nation. Full consent agreement established.`

		stmt, err := ExtractIAACDecisionStatement(fixture)
		if err != nil {
			t.Fatalf("ExtractIAACDecisionStatement failed: %v", err)
		}
		if stmt.ProjectName != "Marathon Palladium-Copper Mine" {
			t.Fatalf("unexpected project: %s", stmt.ProjectName)
		}
		if stmt.TotalConditionsCount != 268 {
			t.Fatalf("unexpected conditions: %d", stmt.TotalConditionsCount)
		}
		if stmt.FinancialAssuranceBondCAD != 45_000_000 {
			t.Fatalf("unexpected bond: $%d", stmt.FinancialAssuranceBondCAD)
		}
		if stmt.DutyToConsultDepth != ConsultDeepAccommodation {
			t.Fatalf("expected ConsultDeepAccommodation, got %s", stmt.DutyToConsultDepth)
		}
		if stmt.InjunctionVulnerabilityScore > 40.0 {
			t.Fatalf("expected low/moderate vulnerability with signed IBA, got %.1f", stmt.InjunctionVulnerabilityScore)
		}
	})

	t.Run("Flags critical litigation exposure when dissenting unceded nation objects", func(t *testing.T) {
		fixture := `Project: Remote Northern Pipeline Corridor
Conditions: 120
Indigenous Conditions: 25
Notes: Project traverses unceded aboriginal title lands with unresolved objection from dissenting nation.`

		stmt, err := ExtractIAACDecisionStatement(fixture)
		if err != nil {
			t.Fatalf("ExtractIAACDecisionStatement failed: %v", err)
		}
		if stmt.LitigationExposure != LitigationCritical && stmt.LitigationExposure != LitigationHigh {
			t.Fatalf("expected high/critical exposure for dissenting nation, got %s", stmt.LitigationExposure)
		}
	})
}
