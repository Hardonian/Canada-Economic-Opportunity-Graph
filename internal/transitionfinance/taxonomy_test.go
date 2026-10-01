package transitionfinance

import (
	"testing"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestTransitionEngine_EvaluateProject(t *testing.T) {
	engine := NewTransitionEngine()

	tests := []struct {
		name             string
		project          *domain.Project
		expectedCategory TaxonomyCategory
		minCredibility   float64
		expectedRating   CredibilityRating
	}{
		{
			name: "Darlington New Nuclear SMR (Green)",
			project: &domain.Project{
				ID:        "darlington-smr",
				Name:      "Darlington New Nuclear Project — Unit 1",
				Sector:    domain.SectorNuclearEnergy,
				Subsector: "Small Modular Reactors",
				CapexCAD:  7_700_000_000,
			},
			expectedCategory: CategoryGreen,
			minCredibility:   85.0,
			expectedRating:   RatingHigh,
		},
		{
			name: "Crawford Nickel Carbon Capture (Transition)",
			project: &domain.Project{
				ID:        "crawford-nickel",
				Name:      "Crawford Nickel-Cobalt Sulphide Project",
				Summary:   "2nd largest nickel reserve globally with net-zero carbon capture via ultramafic tailings mineralization.",
				Sector:    domain.SectorCriticalMinerals,
				Subsector: "Nickel Sulphide",
				CapexCAD:  3_500_000_000,
			},
			expectedCategory: CategoryTransition,
			minCredibility:   80.0,
			expectedRating:   RatingHigh,
		},
		{
			name: "Oneida Energy Storage (Green)",
			project: &domain.Project{
				ID:        "oneida-storage",
				Name:      "Oneida Energy Storage Project",
				Sector:    domain.SectorCleanEnergy,
				Subsector: "Battery Energy Storage",
				CapexCAD:  500_000_000,
			},
			expectedCategory: CategoryGreen,
			minCredibility:   85.0,
			expectedRating:   RatingHigh,
		},
		{
			name: "Canadian International Airports Leasing Hubs (Enabling)",
			project: &domain.Project{
				ID:        "airport-leasing-hubs",
				Name:      "Canadian International Airports — Global Investor Leasing & Infrastructure Hubs",
				Sector:    domain.SectorTransportation,
				Subsector: "Airport Commercial Ground Leases & Intermodal Logistics",
				CapexCAD:  18_000_000_000,
			},
			expectedCategory: CategoryEnabling,
			minCredibility:   75.0,
			expectedRating:   RatingMedium,
		},
		{
			name: "Chisasibi Sovereign AI Compute (Enabling)",
			project: &domain.Project{
				ID:        "chisasibi-ai",
				Name:      "Chisasibi Sovereign AI Hyperscale Cluster",
				Sector:    domain.SectorAICompute,
				Subsector: "Clean Hydro Compute",
				CapexCAD:  2_800_000_000,
			},
			expectedCategory: CategoryEnabling,
			minCredibility:   80.0,
			expectedRating:   RatingHigh,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assessment := engine.EvaluateProject(tc.project)
			if assessment == nil {
				t.Fatalf("expected assessment, got nil")
			}
			if assessment.TaxonomyCategory != tc.expectedCategory {
				t.Errorf("expected taxonomy category %s, got %s", tc.expectedCategory, assessment.TaxonomyCategory)
			}
			if assessment.CredibilityIndex.OverallScore < tc.minCredibility {
				t.Errorf("expected credibility >= %.1f, got %.1f", tc.minCredibility, assessment.CredibilityIndex.OverallScore)
			}
			if assessment.CredibilityIndex.Rating != tc.expectedRating {
				t.Errorf("expected rating %s, got %s", tc.expectedRating, assessment.CredibilityIndex.Rating)
			}
			if assessment.Abatement.LifetimeAbatementTonnesCO2e <= 0 {
				t.Errorf("expected lifetime abatement > 0, got %.0f", assessment.Abatement.LifetimeAbatementTonnesCO2e)
			}
			if assessment.Abatement.MarginalAbatementCostCAD <= 0 {
				t.Errorf("expected MAC > 0, got %.2f", assessment.Abatement.MarginalAbatementCostCAD)
			}
			if assessment.AuditHash == "" {
				t.Errorf("expected non-empty audit hash")
			}
			if !assessment.NoLockInConditionMet {
				t.Errorf("expected no-lock-in condition met")
			}
		})
	}
}

func TestTransitionEngine_EvaluatePortfolio(t *testing.T) {
	engine := NewTransitionEngine()

	projects := []*domain.Project{
		{
			ID:        "darlington-smr",
			Name:      "Darlington New Nuclear Project — Unit 1",
			Sector:    domain.SectorNuclearEnergy,
			CapexCAD:  7_700_000_000,
		},
		{
			ID:        "crawford-nickel",
			Name:      "Crawford Nickel Project",
			Summary:   "carbon capture direct mineral sequestration",
			Sector:    domain.SectorCriticalMinerals,
			CapexCAD:  3_500_000_000,
		},
		{
			ID:        "airport-leasing-hubs",
			Name:      "Canadian International Airports — Global Investor Leasing Hubs",
			Sector:    domain.SectorTransportation,
			CapexCAD:  18_000_000_000,
		},
	}

	summary := engine.EvaluatePortfolio(projects)
	if summary.TotalProjects != 3 {
		t.Fatalf("expected 3 projects, got %d", summary.TotalProjects)
	}
	if summary.TotalGreenCapexCAD <= 0 {
		t.Errorf("expected total green capex > 0, got %d", summary.TotalGreenCapexCAD)
	}
	if summary.TotalTransitionCapexCAD <= 0 {
		t.Errorf("expected total transition capex > 0, got %d", summary.TotalTransitionCapexCAD)
	}
	if summary.AggregateLifetimeAbatementTonnes <= 0 {
		t.Errorf("expected aggregate lifetime abatement > 0, got %.0f", summary.AggregateLifetimeAbatementTonnes)
	}
	if summary.WeightedAverageCredibilityIndex <= 70.0 {
		t.Errorf("expected weighted credibility > 70.0, got %.1f", summary.WeightedAverageCredibilityIndex)
	}
	if summary.AuditHash == "" {
		t.Errorf("expected non-empty portfolio audit hash")
	}
}
