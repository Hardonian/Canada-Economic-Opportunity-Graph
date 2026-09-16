package linkpred

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// PredictedRole classifies the forecasted relationship between an entity and project.
type PredictedRole string

const (
	RoleEPCContractor       PredictedRole = "EPC_CONTRACTOR"
	RoleJointVenturePartner PredictedRole = "JOINT_VENTURE_PARTNER"
	RoleCommercialOfftaker  PredictedRole = "OFFTAKER"
	RoleIndigenousCoOwner   PredictedRole = "INDIGENOUS_CO_OWNER"
	RoleConcessionaryLender PredictedRole = "CONCESSIONARY_LENDER"
)

// PredictedLink represents a machine-learned or graph-topology-inferred relationship.
type PredictedLink struct {
	EntityID         string        `json:"entity_id"`
	EntityName       string        `json:"entity_name"`
	EntityType       string        `json:"entity_type"`
	ProjectID        string        `json:"project_id"`
	ProjectName      string        `json:"project_name"`
	PredictedRole    PredictedRole `json:"predicted_role"`
	ConfidenceScore  float64       `json:"confidence_score"` // 0.0 - 1.0
	AdamicAdarScore  float64       `json:"adamic_adar_score"`
	CommonNeighbors  []string      `json:"common_neighbors"`
	SectorAffinity   float64       `json:"sector_affinity"`
	GeographicMatch  bool          `json:"geographic_match"`
	Rationale        string        `json:"rationale"`
	AuditHash        string        `json:"audit_hash"`
	PredictedAt      time.Time     `json:"predicted_at"`
}

// LinkPredictor runs topological graph algorithms to forecast unannounced relationships.
type LinkPredictor struct{}

func NewLinkPredictor() *LinkPredictor {
	return &LinkPredictor{}
}

// PredictProjectPartners forecasts likely supply chain and capital partners for a project.
func (lp *LinkPredictor) PredictProjectPartners(
	target *domain.Project,
	entities []*domain.Entity,
	projects []*domain.Project,
	relationships []*domain.Relationship,
) []PredictedLink {
	// Build adjacency graphs
	projToEnt := make(map[string]map[string]bool)
	entToProj := make(map[string]map[string]bool)
	entDegrees := make(map[string]int)
	entityMap := make(map[string]*domain.Entity)

	for _, e := range entities {
		entityMap[e.ID] = e
	}

	for _, r := range relationships {
		if r.ProjectID != "" && r.SourceEntityID != "" {
			if projToEnt[r.ProjectID] == nil {
				projToEnt[r.ProjectID] = make(map[string]bool)
			}
			projToEnt[r.ProjectID][r.SourceEntityID] = true

			if entToProj[r.SourceEntityID] == nil {
				entToProj[r.SourceEntityID] = make(map[string]bool)
			}
			entToProj[r.SourceEntityID][r.ProjectID] = true
			entDegrees[r.SourceEntityID]++
		}
		if r.ProjectID != "" && r.TargetEntityID != "" {
			if projToEnt[r.ProjectID] == nil {
				projToEnt[r.ProjectID] = make(map[string]bool)
			}
			projToEnt[r.ProjectID][r.TargetEntityID] = true

			if entToProj[r.TargetEntityID] == nil {
				entToProj[r.TargetEntityID] = make(map[string]bool)
			}
			entToProj[r.TargetEntityID][r.ProjectID] = true
			entDegrees[r.TargetEntityID]++
		}
	}

	existingPartners := projToEnt[target.ID]
	if existingPartners == nil {
		existingPartners = make(map[string]bool)
	}

	predictions := make([]PredictedLink, 0)

	for _, ent := range entities {
		// Skip if already linked or is proponent
		if existingPartners[ent.ID] || ent.ID == target.ProponentID {
			continue
		}

		// Calculate Adamic-Adar index over shared project neighbors
		// AA(target, ent) = sum_{p in N(target) \cap N(ent)} 1 / log(|N(p)|)
		commonProjects := make([]string, 0)
		for pID := range entToProj[ent.ID] {
			if pID == target.ID {
				continue
			}
			// Check if pID shares entities with target
			for sharedEnt := range projToEnt[pID] {
				if existingPartners[sharedEnt] {
					commonProjects = append(commonProjects, pID)
					break
				}
			}
		}

		var aaScore float64
		for _, cp := range commonProjects {
			deg := len(projToEnt[cp])
			if deg > 1 {
				aaScore += 1.0 / math.Log(float64(deg))
			} else {
				aaScore += 1.0
			}
		}

		// Sector affinity heuristic
		sectorAffinity := 0.5
		entLower := strings.ToLower(ent.LegalName + " " + ent.CommonName + " " + ent.Description)
		switch target.Sector {
		case domain.SectorNuclearEnergy:
			if strings.Contains(entLower, "nuclear") || strings.Contains(entLower, "opg") || strings.Contains(entLower, "bruce") {
				sectorAffinity = 0.95
			}
		case domain.SectorCriticalMinerals:
			if strings.Contains(entLower, "mining") || strings.Contains(entLower, "metals") || strings.Contains(entLower, "refining") {
				sectorAffinity = 0.90
			}
		case domain.SectorAICompute:
			if strings.Contains(entLower, "compute") || strings.Contains(entLower, "cloud") || strings.Contains(entLower, "hydro") {
				sectorAffinity = 0.90
			}
		}

		// Geo match
		geoMatch := strings.EqualFold(ent.Jurisdiction, target.Province) || strings.EqualFold(ent.Jurisdiction, "CA") || strings.EqualFold(ent.Jurisdiction, "Federal")

		// Classify role
		var role PredictedRole
		switch {
		case strings.Contains(entLower, "first nation") || strings.Contains(entLower, "cree") || strings.Contains(entLower, "inuit") || ent.EntityType == "FirstNation":
			role = RoleIndigenousCoOwner
		case strings.Contains(entLower, "bank") || strings.Contains(entLower, "cib") || strings.Contains(entLower, "infrastructure bank"):
			role = RoleConcessionaryLender
		case strings.Contains(entLower, "engineering") || strings.Contains(entLower, "construction") || strings.Contains(entLower, "snc") || strings.Contains(entLower, "aecon"):
			role = RoleEPCContractor
		case strings.Contains(entLower, "offtaker") || strings.Contains(entLower, "oem") || strings.Contains(entLower, "auto") || strings.Contains(entLower, "utility"):
			role = RoleCommercialOfftaker
		default:
			role = RoleJointVenturePartner
		}

		// Confidence calculation
		confidence := (math.Min(aaScore, 3.0) / 3.0) * 0.30 + (sectorAffinity * 0.45)
		if geoMatch {
			confidence += 0.25
		}
		if confidence > 0.99 {
			confidence = 0.99
		}

		if confidence >= 0.50 {
			name := ent.CommonName
			if name == "" {
				name = ent.LegalName
			}
			rationale := fmt.Sprintf("High topology affinity (AA=%.2f, sector=%.0f%%) across shared capital consortia and %s jurisdiction.",
				aaScore, sectorAffinity*100, ent.Jurisdiction)

			pLink := PredictedLink{
				EntityID:        ent.ID,
				EntityName:      name,
				EntityType:      ent.EntityType,
				ProjectID:       target.ID,
				ProjectName:     target.Name,
				PredictedRole:   role,
				ConfidenceScore: math.Round(confidence*100) / 100,
				AdamicAdarScore: math.Round(aaScore*100) / 100,
				CommonNeighbors: commonProjects,
				SectorAffinity:  sectorAffinity,
				GeographicMatch: geoMatch,
				Rationale:       rationale,
				PredictedAt:     time.Now().UTC(),
			}
			h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%.2f", pLink.EntityID, pLink.ProjectID, pLink.PredictedRole, pLink.ConfidenceScore)))
			pLink.AuditHash = hex.EncodeToString(h[:])

			predictions = append(predictions, pLink)
		}
	}

	sort.Slice(predictions, func(i, j int) bool {
		return predictions[i].ConfidenceScore > predictions[j].ConfidenceScore
	})

	if len(predictions) > 10 {
		predictions = predictions[:10]
	}

	return predictions
}
