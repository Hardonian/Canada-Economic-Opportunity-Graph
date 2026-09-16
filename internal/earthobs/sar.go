package earthobs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// Constellation identifies the satellite source.
type Constellation string

const (
	ConstellationRCM       Constellation = "RADARSAT Constellation Mission (RCM)"
	ConstellationSentinel1 Constellation = "Copernicus Sentinel-1 (C-Band SAR)"
	ConstellationSentinel2 Constellation = "Copernicus Sentinel-2 (Multispectral Optical)"
	ConstellationPlanet    Constellation = "PlanetScope High-Res"
)

// GroundTruthObservation represents a single satellite pass observation.
type GroundTruthObservation struct {
	ObservationID          string        `json:"observation_id"`
	ProjectID              string        `json:"project_id"`
	Constellation          Constellation `json:"constellation"`
	PassDate               time.Time     `json:"pass_date"`
	SARCoherenceChange     float64       `json:"sar_coherence_change"`     // -1.0 to 1.0 (High positive indicates massive surface displacement/earthworks)
	OpticalVegetationIndex float64       `json:"optical_vegetation_index"` // NDVI (Drop indicates site clearing)
	StructuralRadarCrossSec float64      `json:"structural_rcs_db"`        // Backscatter in dB (Increase indicates steel/concrete erect structures)
	EarthworksDetected     bool          `json:"earthworks_detected"`
	StructuralPourDetected bool          `json:"structural_pour_detected"`
	ResolutionMeters       float64       `json:"resolution_meters"`
	EvidenceID             string        `json:"evidence_id"`
}

// GroundTruthDossier provides satellite verification of physical progress vs claimed milestone.
type GroundTruthDossier struct {
	ProjectID                string                   `json:"project_id"`
	ClaimedStage             domain.LifecycleStage    `json:"claimed_stage"`
	CorroborationStatus      domain.ConfidenceLevel   `json:"corroboration_status"` // VERIFIED, REPORTED, CONFLICTED
	PhysicalProgressScore    float64                  `json:"physical_progress_score"` // 0.0 - 100.0
	EarthworksConfirmed      bool                     `json:"earthworks_confirmed"`
	StructuresConfirmed      bool                     `json:"structures_confirmed"`
	LastSatellitePass        time.Time                `json:"last_satellite_pass"`
	Observations             []GroundTruthObservation `json:"observations"`
	TelemetrySummary         string                   `json:"telemetry_summary"`
	AuditHash                string                   `json:"audit_hash"`
	EvaluatedAt              time.Time                `json:"evaluated_at"`
}

// Evaluator correlates satellite passes with economic graph milestones.
type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// CorroborateProject validates physical site telemetry against claimed lifecycle stage.
func (e *Evaluator) CorroborateProject(project *domain.Project, obs []GroundTruthObservation) *GroundTruthDossier {
	dossier := &GroundTruthDossier{
		ProjectID:    project.ID,
		ClaimedStage: project.CurrentStage,
		Observations: obs,
		EvaluatedAt:  time.Now().UTC(),
	}

	if len(obs) == 0 {
		dossier.CorroborationStatus = domain.ConfidenceReported
		dossier.PhysicalProgressScore = 50.0
		dossier.TelemetrySummary = "No recent SAR or optical satellite passes indexed for this project boundary."
		dossier.AuditHash = calculateHash(dossier)
		return dossier
	}

	var hasEarthworks, hasStructures bool
	var latestPass time.Time

	for _, o := range obs {
		if o.PassDate.After(latestPass) {
			latestPass = o.PassDate
		}
		if o.EarthworksDetected || o.SARCoherenceChange > 0.30 || o.OpticalVegetationIndex < 0.20 {
			hasEarthworks = true
		}
		if o.StructuralPourDetected || (o.StructuralRadarCrossSec != 0 && o.StructuralRadarCrossSec > -10.0) {
			hasStructures = true
		}
	}

	dossier.LastSatellitePass = latestPass
	dossier.EarthworksConfirmed = hasEarthworks
	dossier.StructuresConfirmed = hasStructures

	// Compare with claimed stage
	stage := project.CurrentStage
	isAdvanced := stage == domain.StageConstruction || stage == domain.StageCommissioning || stage == domain.StageOperating

	if isAdvanced {
		if hasStructures {
			dossier.CorroborationStatus = domain.ConfidenceVerified
			dossier.PhysicalProgressScore = 95.0
			dossier.TelemetrySummary = "Satellite telemetry independently corroborates active civil construction and structural erection."
		} else if hasEarthworks {
			dossier.CorroborationStatus = domain.ConfidenceSupported
			dossier.PhysicalProgressScore = 70.0
			dossier.TelemetrySummary = "Satellite telemetry confirms site clearing and earthworks; structural foundation phase pending."
		} else {
			dossier.CorroborationStatus = domain.ConfidenceConflict
			dossier.PhysicalProgressScore = 20.0
			dossier.TelemetrySummary = "DISCREPANCY: Proponent claims active construction, but multi-temporal SAR and optical sensors detect undisturbed terrain."
		}
	} else {
		// Planning / Permitting / Feasibility
		if hasEarthworks || hasStructures {
			dossier.CorroborationStatus = domain.ConfidenceSupported
			dossier.PhysicalProgressScore = 80.0
			dossier.TelemetrySummary = "Early pre-construction preparatory works or road access construction detected ahead of announced FID."
		} else {
			dossier.CorroborationStatus = domain.ConfidenceVerified
			dossier.PhysicalProgressScore = 100.0
			dossier.TelemetrySummary = "Site remains in baseline condition, fully consistent with pre-construction regulatory/planning stage."
		}
	}

	dossier.AuditHash = calculateHash(dossier)
	return dossier
}

func calculateHash(d *GroundTruthDossier) string {
	raw := fmt.Sprintf("%s|%s|%s|%.1f|%t|%t",
		d.ProjectID, d.ClaimedStage, d.CorroborationStatus, d.PhysicalProgressScore, d.EarthworksConfirmed, d.StructuresConfirmed)
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}
