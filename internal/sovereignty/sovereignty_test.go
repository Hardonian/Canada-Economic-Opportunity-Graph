package sovereignty

import (
	"testing"
)

func TestAISovereigntyEvaluation(t *testing.T) {
	profile := &AIProfile{
		SubjectID:          "sovereign-llm-ca",
		SubjectName:        "MapleLLM Sovereign Model",
		DataResidencyCA:    true,
		ComputeResidencyCA: true,
		CanadianOwnership:  1.0,
		ForeignLegalRisk:   1.0,
		LocalDeployment:    true,
		OfflineCapability:  true,
		OpenWeights:        true,
		BilingualCapacity:  1.0,
		QuebecLaw25Ready:   true,
		CleanEnergySource:  0.95,
	}

	res := EvaluateSovereignty(profile)
	if res.OverallScore < 85.0 {
		t.Errorf("expected high sovereignty score for 100%% Canadian stack, got %f", res.OverallScore)
	}
}

func TestICAScreener(t *testing.T) {
	screener := NewICAScreener()

	// Hostile SOE acquisition of critical mineral project
	req := ICAScreeningRequest{
		ProjectID:                "proj-crawford-nickel",
		Sector:                   "Critical Minerals",
		AcquiringEntity:          "China Minmetals SOE",
		AcquiringJurisdiction:    "CN",
		IsStateOwned:             true,
		AcquisitionSharePct:      20.0,
		IncludesCriticalMinerals: true,
	}

	verdict := screener.ScreenAcquisition(req)
	if verdict.Verdict != "PROHIBITED" {
		t.Errorf("expected PROHIBITED verdict, got %s", verdict.Verdict)
	}
	if verdict.StatutorySection != "Section 25.4 (Order in Council Divestment / Blocking Order)" {
		t.Errorf("expected Section 25.4, got %s", verdict.StatutorySection)
	}
	if verdict.SovereigntyScore > 30.0 {
		t.Errorf("expected low sovereignty score for prohibited SOE takeover, got %f", verdict.SovereigntyScore)
	}
}
