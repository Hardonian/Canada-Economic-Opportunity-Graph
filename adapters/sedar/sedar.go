// Package sedar provides an ingestion adapter for Canadian public capital market disclosures
// from SEDAR+ (System for Electronic Document Analysis and Retrieval), including NI 43-101
// technical reports, quarterly MD&As, prospectuses, and material change reports.
package sedar

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/adapters"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/identity"
)

const (
	adapterName     = "sedar_plus_disclosures"
	pipelineVersion = "sedar-v1"
	parserVersion   = "sedar-json-v1"
)

var (
	// Capex extraction regex matching CAD amounts in disclosures
	sedarCapexPattern = regexp.MustCompile(`(?i)(?:CAD|C\$|\$)\s*([\d,]+(?:\.\d+)?)\s*(?:million|billion|M|B)`)
)

// FilingRecord represents a raw filing metadata entry from SEDAR+.
type FilingRecord struct {
	FilingID        string   `json:"filing_id"`
	IssuerName      string   `json:"issuer_name"`
	Ticker          string   `json:"ticker"`
	Exchange        string   `json:"exchange"` // TSX, TSXV, CSE
	DocumentType    string   `json:"document_type"` // NI_43_101, MDA, PROSPECTUS, MATERIAL_CHANGE
	ProjectName     string   `json:"project_name"`
	Sector          string   `json:"sector"`
	Province        string   `json:"province"`
	Headline        string   `json:"headline"`
	FilingDate      string   `json:"filing_date"`
	DocumentURL     string   `json:"document_url"`
	DeclaredCapexCAD int64   `json:"declared_capex_cad,omitempty"`
	Summary         string   `json:"summary"`
	Keywords        []string `json:"keywords"`
}

// SEDARAdapter extracts public company market disclosures and capital items.
type SEDARAdapter struct {
	fixtureData []byte
	health      adapters.SourceHealth
}

// NewSEDARAdapter constructs a SEDAR+ adapter instance.
func NewSEDARAdapter(fixtureData []byte) *SEDARAdapter {
	now := time.Now().UTC()
	return &SEDARAdapter{
		fixtureData: fixtureData,
		health: adapters.SourceHealth{
			AdapterName:    adapterName,
			Tier:           domain.SourceTier1,
			Status:         "HEALTHY",
			LastAttempt:    now,
			LastSuccess:    now,
			LastChange:     now,
			Mode:           "CURATED_SNAPSHOT",
			RateLimitState: "NORMAL",
		},
	}
}

// Name returns the canonical adapter name.
func (a *SEDARAdapter) Name() string {
	return adapterName
}

// Tier returns Tier 1 (Official public regulator filing).
func (a *SEDARAdapter) Tier() domain.SourceTier {
	return domain.SourceTier1
}

// Health returns the operational health status.
func (a *SEDARAdapter) Health() *adapters.SourceHealth {
	return &a.health
}

// Fetch retrieves filing data (or returns embedded canonical market disclosures).
func (a *SEDARAdapter) Fetch(ctx context.Context) ([]byte, error) {
	if len(a.fixtureData) > 0 {
		return a.fixtureData, nil
	}
	return CanonicalSEDARFilings()
}

// Parse converts SEDAR+ filings into normalized CEGS domain entities, projects, and capital items.
func (a *SEDARAdapter) Parse(data []byte) (*adapters.IngestionResult, error) {
	var records []FilingRecord
	if err := json.Unmarshal(data, &records); err != nil {
		a.health.ParseFailures++
		a.health.LastError = fmt.Sprintf("failed to parse SEDAR+ json: %v", err)
		return nil, err
	}

	result := &adapters.IngestionResult{}
	now := time.Now().UTC()

	for _, rec := range records {
		a.health.DocumentsSeen++
		docHash := sha256.Sum256([]byte(rec.FilingID + rec.FilingDate + rec.Headline))
		recordHash := hex.EncodeToString(docHash[:])
		evidenceID := identity.StableID("evidence", adapterName, rec.FilingID+":"+recordHash)

		fDate, _ := time.Parse("2006-01-02", rec.FilingDate)
		if fDate.IsZero() {
			fDate = now
		}

		evidence := &domain.Evidence{
			ID:                 evidenceID,
			SourceURL:          rec.DocumentURL,
			Publisher:          "SEDAR+ / Canadian Securities Administrators (CSA)",
			SourceTier:         domain.SourceTier1,
			Visibility:         domain.VisibilityPublicAttribution,
			Publishable:        true,
			RetrievalTimestamp: now,
			PublicationDate:    &fDate,
			Confidence:         domain.ConfidenceVerified,
			ExtractionMethod:   "official_sedar_pipeline",
			ContentHash:        recordHash,
			HashScope:          "normalized_source_record",
			SourceClass:        "SECURITIES_REGULATORY_FILING",
			SourceID:           adapterName,
			SourceRecordID:     rec.FilingID,
			PipelineVersion:    pipelineVersion,
			ParserVersion:      parserVersion,
			RawSnippet:         rec.Headline + " - " + rec.Summary,
			Attribution:        "SEDAR+ Open Capital Markets Data",
		}
		result.Evidence = append(result.Evidence, evidence)

		// Create Proponent Entity
		entityID := identity.NormalizeEntityID(rec.IssuerName)
		entity := &domain.Entity{
			ID:          entityID,
			Slug:        entityID,
			LegalName:   rec.IssuerName,
			CommonName:  rec.IssuerName,
			EntityType:  "Corporation",
			Jurisdiction: "CA",
			Identifiers: map[string]string{
				"Ticker":   rec.Ticker,
				"Exchange": rec.Exchange,
			},
			Description: fmt.Sprintf("Public issuer listed on %s (%s)", rec.Exchange, rec.Ticker),
			SourceIDs:   []string{adapterName},
			EvidenceIDs: []string{evidenceID},
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		result.Entities = append(result.Entities, entity)

		// Create Project if specified
		if rec.ProjectName != "" {
			projID := identity.NormalizeProjectID(rec.ProjectName)
			sec := parseSector(rec.Sector)
			capex := rec.DeclaredCapexCAD
			if capex == 0 {
				capex = extractCapexFromText(rec.Summary)
			}

			project := &domain.Project{
				ID:           projID,
				Slug:         projID,
				Name:         rec.ProjectName,
				Sector:       sec,
				Subsector:    rec.DocumentType,
				Province:     rec.Province,
				CurrentStage: domain.StageFEED,
				CapexCAD:     capex,
				SourceIDs:    []string{adapterName},
				EvidenceIDs:  []string{evidenceID},
				CreatedAt:    now,
				UpdatedAt:    now,
			}
			result.Projects = append(result.Projects, project)

			// Record Capital Item
			if capex > 0 {
				capItem := &domain.CapitalItem{
					ID:               "cap-sedar-" + rec.FilingID,
					ProjectID:        projID,
					Category:         domain.CapitalCategoryPrivate,
					Status:           domain.CapitalStatusCommitted,
					AmountCAD:        capex,
					AmountType:       "estimated",
					ProviderName:     rec.IssuerName,
					ProviderEntityID: entityID,
					Notes:            fmt.Sprintf("%s - %s filing", rec.DocumentType, rec.Headline),
					EvidenceID:       evidenceID,
					Visibility:       domain.VisibilityPublicAttribution,
					Publishable:      true,
					CreatedAt:        now,
				}
				result.CapitalItems = append(result.CapitalItems, capItem)
			}

			// Relationship between Issuer and Project
			rel := &domain.Relationship{
				ID:             fmt.Sprintf("rel-%s-%s", entityID, projID),
				ProjectID:      projID,
				SourceEntityID: entityID,
				TargetEntityID: projID,
				RelationType:   "proponent",
				Confidence:     domain.ConfidenceVerified,
				EvidenceID:     evidenceID,
				CreatedAt:      now,
			}
			result.Relationships = append(result.Relationships, rel)

			// Milestone Event
			event := &domain.Event{
				ID:          "ev-milestone-" + rec.FilingID,
				ProjectID:   projID,
				EventType:   "regulatory_filing",
				EventDate:   fDate,
				Title:       fmt.Sprintf("SEDAR+ %s: %s", rec.DocumentType, rec.Headline),
				Description: rec.Summary,
				EvidenceID:  evidenceID,
				CreatedAt:   now,
			}
			result.Events = append(result.Events, event)
		}
	}

	a.health.LastSuccess = now
	return result, nil
}

func parseSector(s string) domain.Sector {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "CRITICAL_MINERALS", "MINING":
		return domain.SectorCriticalMinerals
	case "CLEAN_ENERGY", "ENERGY":
		return domain.SectorCleanEnergy
	case "NUCLEAR":
		return domain.SectorNuclearEnergy
	case "TRANSPORTATION":
		return domain.SectorTransportation
	case "AI_COMPUTE", "TECHNOLOGY":
		return domain.SectorAICompute
	default:
		return domain.SectorCleanEnergy
	}
}

func extractCapexFromText(text string) int64 {
	matches := sedarCapexPattern.FindStringSubmatch(text)
	if len(matches) < 2 {
		return 0
	}
	clean := strings.ReplaceAll(matches[1], ",", "")
	val, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0
	}
	lower := strings.ToLower(matches[0])
	if strings.Contains(lower, "billion") || strings.Contains(lower, "b") {
		return int64(val * 1_000_000_000)
	}
	if strings.Contains(lower, "million") || strings.Contains(lower, "m") {
		return int64(val * 1_000_000)
	}
	return int64(val)
}

// CanonicalSEDARFilings generates simulated high-fidelity SEDAR+ disclosures.
func CanonicalSEDARFilings() ([]byte, error) {
	filings := []FilingRecord{
		{
			FilingID:         "sedar-cnc-2026-q1",
			IssuerName:       "Canada Nickel Company Inc.",
			Ticker:           "CNC",
			Exchange:         "TSXV",
			DocumentType:     "NI_43_101",
			ProjectName:      "Crawford Nickel Sulphide Project",
			Sector:           "CRITICAL_MINERALS",
			Province:         "ON",
			Headline:         "Updated Bankable Feasibility Study and Carbon Capture Integration",
			FilingDate:       "2026-03-15",
			DocumentURL:      "https://www.sedarplus.ca/crawford-bfs-2026.pdf",
			DeclaredCapexCAD: 2_500_000_000,
			Summary:          "Canada Nickel files updated NI 43-101 technical report confirming initial capital expenditure of $2.5 billion CAD with IPT tailings carbon mineralization.",
			Keywords:         []string{"nickel", "carbon capture", "feasibility", "clean mining"},
		},
		{
			FilingID:         "sedar-fmc-2026-mda",
			IssuerName:       "First Majestic Critical Metals Corp.",
			Ticker:           "FM",
			Exchange:         "TSX",
			DocumentType:     "MDA",
			ProjectName:      "Athabasca Clean Energy Intertie",
			Sector:           "CLEAN_ENERGY",
			Province:         "SK",
			Headline:         "Management Discussion and Analysis for Fiscal Year Ended 2025",
			FilingDate:       "2026-02-28",
			DocumentURL:      "https://www.sedarplus.ca/fm-mda-2025.pdf",
			DeclaredCapexCAD: 450_000_000,
			Summary:          "Company allocates $450 million CAD for transmission interconnection and off-grid renewable power infrastructure in northern Saskatchewan.",
			Keywords:         []string{"intertie", "transmission", "clean energy"},
		},
	}
	return json.Marshal(filings)
}
