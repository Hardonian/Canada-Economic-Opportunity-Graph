package earthobs

import (
	"fmt"
	"math"
	"time"
)

// DiscrepancyReport exposes discrepancies between corporate filings and satellite reality.
type DiscrepancyReport struct {
	ProjectID           string    `json:"project_id"`
	ClaimedProgressPct  float64   `json:"claimed_progress_pct"`
	ObservedProgressPct float64   `json:"observed_progress_pct"`
	VariancePct         float64   `json:"variance_pct"`
	IsDiscrepancyFlagged bool     `json:"is_discrepancy_flagged"`
	RiskVerdict         string    `json:"risk_verdict"` // "VERIFIED_CONCURRENT", "SCHEDULE_RISK", "MATERIAL_DECEPTION_ALERT"
	Explanation         string    `json:"explanation"`
	EvaluatedAt         time.Time `json:"evaluated_at"`
}

// GroundTruthDiscrepancyEngine cross-audits regulatory filings with orbital telemetry.
type GroundTruthDiscrepancyEngine struct{}

// NewGroundTruthDiscrepancyEngine creates an orbital discrepancy engine.
func NewGroundTruthDiscrepancyEngine() *GroundTruthDiscrepancyEngine {
	return &GroundTruthDiscrepancyEngine{}
}

// EvaluateDiscrepancy cross-checks claimed construction progress against satellite detections.
func (ge *GroundTruthDiscrepancyEngine) EvaluateDiscrepancy(
	projectID string,
	claimedPct float64,
	mach *MachineryDetectionResult,
	sar *InSARCoherenceAnalysis,
) *DiscrepancyReport {
	// Synthesize physical observed progress
	machScore := 0.0
	if mach != nil {
		machScore = math.Min(100.0, float64(mach.TotalMachineryUnits)*4.5)
	}

	sarScore := 50.0
	if sar != nil {
		if sar.TailingsStabilityRating == "STABLE" && sar.InterferometricCoherence > 0.5 {
			sarScore = 80.0
		} else if sar.TailingsStabilityRating == "UNSTABLE" {
			sarScore = 20.0
		}
	}

	observed := math.Round((machScore*0.6+sarScore*0.4)*10) / 10
	variance := math.Round((claimedPct-observed)*10) / 10

	flagged := false
	verdict := "VERIFIED_CONCURRENT"
	explanation := "Satellite ground-truth observations corroborate claimed corporate progress reports."

	if variance > 35.0 {
		flagged = true
		verdict = "MATERIAL_DECEPTION_ALERT"
		explanation = fmt.Sprintf("Severe divergence: Proponent claims %.1f%% progress, but satellite CV and InSAR observe only %.1f%% active physical infrastructure.", claimedPct, observed)
	} else if variance > 15.0 {
		flagged = true
		verdict = "SCHEDULE_RISK"
		explanation = fmt.Sprintf("Moderate divergence: Proponent claims %.1f%% progress; satellite telemetry observes %.1f%% physical site activity.", claimedPct, observed)
	}

	return &DiscrepancyReport{
		ProjectID:            projectID,
		ClaimedProgressPct:   claimedPct,
		ObservedProgressPct:  observed,
		VariancePct:          variance,
		IsDiscrepancyFlagged: flagged,
		RiskVerdict:          verdict,
		Explanation:          explanation,
		EvaluatedAt:          time.Now().UTC(),
	}
}
