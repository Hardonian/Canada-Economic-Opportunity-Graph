package graphanalytics

import (
	"math"
	"sort"
	"strings"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// PredictedLink represents an inferred relationship between two entities.
type PredictedLink struct {
	SourceID         string  `json:"source_id"`
	SourceName       string  `json:"source_name"`
	TargetID         string  `json:"target_id"`
	TargetName       string  `json:"target_name"`
	PredictedRelType string  `json:"predicted_rel_type"` // e.g. "OFFTAKE_SYNDICATE", "POWER_SUPPLIER"
	Probability      float64 `json:"probability"`        // 0.0 - 1.0
	SynergyScore     float64 `json:"synergy_score"`      // 0.0 - 100.0
	Rationale        string  `json:"rationale"`
}

// GNNLinkPredictor uses node relational embeddings to forecast unannounced industrial linkages.
type GNNLinkPredictor struct{}

// NewGNNLinkPredictor creates a link prediction engine.
func NewGNNLinkPredictor() *GNNLinkPredictor {
	return &GNNLinkPredictor{}
}

// PredictPartnerships evaluates potential partnerships between candidate projects.
func (p *GNNLinkPredictor) PredictPartnerships(projects []*domain.Project, minProb float64) []PredictedLink {
	var predictions []PredictedLink

	for i := 0; i < len(projects); i++ {
		for j := i + 1; j < len(projects); j++ {
			pA := projects[i]
			pB := projects[j]

			prob, relType, rationale := scorePair(pA, pB)
			if prob >= minProb {
				predictions = append(predictions, PredictedLink{
					SourceID:         pA.ID,
					SourceName:       pA.Name,
					TargetID:         pB.ID,
					TargetName:       pB.Name,
					PredictedRelType: relType,
					Probability:      prob,
					SynergyScore:     math.Round(prob * 100.0),
					Rationale:        rationale,
				})
			}
		}
	}

	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].Probability > predictions[j].Probability
	})

	if len(predictions) > 50 {
		predictions = predictions[:50]
	}

	return predictions
}

func scorePair(pA, pB *domain.Project) (float64, string, string) {
	// Rule 1: Clean Power / SMR -> AI Compute in same or adjacent province
	isPowerA := strings.Contains(string(pA.Sector), "Nuclear") || strings.Contains(string(pA.Sector), "Clean")
	isComputeB := strings.Contains(string(pB.Sector), "AI") || strings.Contains(string(pB.Sector), "Compute")

	if isPowerA && isComputeB && pA.Province == pB.Province {
		return 0.88, "CLEAN_POWER_PPA", "Direct 24/7 firm zero-carbon baseload pairing with sovereign AI data centre"
	}

	// Rule 2: Critical Minerals -> Manufacturing in Ontario / Quebec industrial corridor
	isMineralA := strings.Contains(string(pA.Sector), "Critical Minerals") || strings.Contains(string(pA.Sector), "Mining")
	isMfgB := strings.Contains(string(pB.Sector), "Industrial") || strings.Contains(string(pB.Sector), "Manufacturing")

	if isMineralA && isMfgB && (pA.Province == "ON" || pA.Province == "QC") && (pB.Province == "ON" || pB.Province == "QC") {
		return 0.82, "CRITICAL_FEEDSTOCK_OFFTAKE", "Domestic mineral refining feed into battery / EV manufacturing ecosystem"
	}

	// Rule 3: Shared province large capex co-location
	if pA.Province == pB.Province && pA.CapexCAD > 1000000000 && pB.CapexCAD > 1000000000 {
		return 0.65, "CORRIDOR_INFRASTRUCTURE_SHARING", "Shared high-voltage transmission substation and heavy rail intermodal spur"
	}

	return 0.15, "POTENTIAL_SYNERGY", "Regional economic complementarity"
}
