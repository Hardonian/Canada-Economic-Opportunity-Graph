package filings

import (
	"testing"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

func TestStreamerParseFiling(t *testing.T) {
	streamer := NewStreamer()
	now := time.Now().UTC()

	rawText := `
		MANAGEMENT'S DISCUSSION AND ANALYSIS - Q3 2026
		The Board of Directors has reached a Positive Final Investment Decision (FID) for the project.
		Early works capital expenditure estimate was revised upward to $1.8B CAD.
		Major EPC contract awarded to Tier-1 contractor for FEED and site preparation infrastructure.
	`

	record := streamer.ParseFiling(
		"Canada Nickel Company Inc.",
		"CNC",
		"TSXV",
		FilingTypeMDA,
		"Q3 2026 MD&A Disclosure",
		rawText,
		"https://sedarplus.ca/filings/10049281.pdf",
		now,
	)

	if record == nil {
		t.Fatal("expected non-nil filing record")
	}
	if !record.StageChangeDetected {
		t.Errorf("expected StageChangeDetected to be true")
	}
	if record.DetectedStage != domain.StageConstruction {
		t.Errorf("expected stage %v, got %v", domain.StageConstruction, record.DetectedStage)
	}
	if len(record.ContractAwards) == 0 {
		t.Errorf("expected at least one contract award detected")
	}
	if record.AuditHash == "" {
		t.Errorf("expected non-empty audit hash")
	}
}

func TestEAParser(t *testing.T) {
	parser := NewEAParser()
	now := time.Now().UTC()

	rawNotice := `
		B.C. Environmental Assessment Office:
		An Environmental Assessment Certificate approval granted for the LNG terminal facility.
	`

	rec := parser.ParseNotice(
		"BC_EAO",
		"BC",
		"Cedar LNG Project",
		"EA-2026-081",
		"Issuance of Environmental Assessment Certificate",
		rawNotice,
		"https://projects.eao.gov.bc.ca/p/cedarlng",
		now,
	)

	if rec == nil {
		t.Fatal("expected non-nil EA record")
	}
	if rec.Milestone != EAMilestoneCertificateIssued {
		t.Errorf("expected milestone %s, got %s", EAMilestoneCertificateIssued, rec.Milestone)
	}
	if !rec.Approved {
		t.Errorf("expected Approved to be true")
	}
	if rec.ConditionsCount == 0 {
		t.Errorf("expected conditions count > 0")
	}
}

func TestAmendmentTracker(t *testing.T) {
	tracker := NewAmendmentTracker()
	now := time.Now().UTC()

	rawText := `
		CanadaBuys Tender Notice DCC-2026-HQ-0199:
		Contract awarded to Aecon-PCL Industrial Joint Venture for $185M CAD.
	`

	amend := tracker.TrackAmendment("WS39482910-Doc29102", 3, rawText, "https://canadabuys.canada.ca/tender/WS39482910", now)
	if amend == nil {
		t.Fatal("expected non-nil amendment")
	}
	if amend.Type != AmendmentAwardNotice {
		t.Errorf("expected AmendmentAwardNotice, got %s", amend.Type)
	}
	if amend.ContractValueCAD != 185_000_000 {
		t.Errorf("expected 185000000 CAD, got %d", amend.ContractValueCAD)
	}
	if amend.AuditHash == "" {
		t.Errorf("expected valid audit hash")
	}
}
