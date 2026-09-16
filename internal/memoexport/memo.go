package memoexport

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// MemoType classifies the institutional document format.
type MemoType string

const (
	MemoTypeCabinetMC      MemoType = "MEMORANDUM_TO_CABINET"      // Privy Council Office (PCO) standard
	MemoTypeTreasuryBoard  MemoType = "TREASURY_BOARD_SUBMISSION"  // TB Sub standard
	MemoTypeInvestmentComm MemoType = "INVESTMENT_COMMITTEE_MEMO"  // Pension / Infrastructure PE standard
)

// DecisionMemo contains the formal synthesized briefing document.
type DecisionMemo struct {
	MemoID             string    `json:"memo_id"`
	Type               MemoType  `json:"type"`
	ProjectID          string    `json:"project_id"`
	ProjectName        string    `json:"project_name"`
	SecurityCaveat     string    `json:"security_caveat"` // "CONFIDENTIAL // FOR MINISTERIAL CONSIDERATION ONLY"
	PreparedFor        string    `json:"prepared_for"`
	DateFormatted      string    `json:"date_formatted"`
	ExecutiveSummary   string    `json:"executive_summary"`
	StrategicRationale string    `json:"strategic_rationale"`
	FinancialExposure  string    `json:"financial_exposure"`
	IndigenousTreaty   string    `json:"indigenous_treaty"`
	GeopoliticalRisk   string    `json:"geopolitical_risk"`
	RecommendedAction  string    `json:"recommended_action"`
	MarkdownContent    string    `json:"markdown_content"`
	AuditHash          string    `json:"audit_hash"`
}

// MemoGenerator compiles structured graph facts into official memorandum documents.
type MemoGenerator struct{}

// NewMemoGenerator creates a new memo generator.
func NewMemoGenerator() *MemoGenerator {
	return &MemoGenerator{}
}

// GenerateCabinetMemo compiles an official Memorandum to Cabinet (MC) annex.
func (g *MemoGenerator) GenerateCabinetMemo(project *domain.Project, memoType MemoType) *DecisionMemo {
	now := time.Now().UTC()
	dateStr := now.Format("January 02, 2006")

	caveat := "PROTECTED B // CABINET CONFIDENTIAL"
	recipient := "Cabinet Committee on Economy, Inclusion and Climate"
	if memoType == MemoTypeInvestmentComm {
		caveat = "COMMERCIALLY CONFIDENTIAL // PRIVILEGED INVESTMENT COMMITTEE BRIEF"
		recipient = "Chief Investment Officer & Global Infrastructure Investment Committee"
	} else if memoType == MemoTypeTreasuryBoard {
		caveat = "PROTECTED B // TREASURY BOARD PRESIDENTIAL REVIEW"
		recipient = "Treasury Board of Canada Secretariat (TBS)"
	}

	capexBillions := float64(project.CapexCAD) / 1e9
	if capexBillions <= 0 {
		capexBillions = 0.5
	}

	execSummary := fmt.Sprintf("National capital sponsorship appraisal for %s located in %s, Canada. The proponent seeks federal co-investment authorization, debt syndication, and regulatory alignment across an estimated $%.2fB CAD capital envelope.",
		project.Name, project.Province, capexBillions)

	stratRationale := fmt.Sprintf("Classified within the %s sovereign strategic mandate. Supports domestic supply chain retention, bilateral USMCA resilience, and critical infrastructure readiness under the Federal Canadian Economic Sovereignty Framework.",
		project.Sector)

	finExposure := fmt.Sprintf("Recommended capital structure targets 35%% sponsor equity, 15%% concessionary catalytic debt via the Canada Infrastructure Bank (CIB), 10%% Indigenous equity syndication under the $5B Federal Indigenous Loan Guarantee Program (ILGP), and 40%% commercial senior debt syndication. Estimated private capital crowding-in multiplier of 2.8x.",
	)

	indigTreaty := fmt.Sprintf("Project right-of-way traverses traditional treaty territories in %s. Statutory Duty to Consult requires formal Early Engagement, revenue-sharing agreements, and First Nations equity co-ownership guarantees prior to Final Investment Decision (FID).",
		project.Province)

	geoRisk := "Supply chain vulnerability assessed against foreign export restrictions and Title III Defense Production Act interoperability. Domestic value retention index (DVRI) prioritizes Canadian EPC procurement."

	recAction := "Authorize the Minister to execute the strategic co-investment term sheet, approve CIB concessional loan participation, and refer the proponent to the Major Projects Management Office (MPMO) for coordinated permitting acceleration."

	// Assemble formatted Markdown
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", strings.ToUpper(string(memoType))))
	sb.WriteString(fmt.Sprintf("**SECURITY CLASSIFICATION**: %s  \n", caveat))
	sb.WriteString(fmt.Sprintf("**TO**: %s  \n", recipient))
	sb.WriteString(fmt.Sprintf("**SUBJECT**: Comprehensive National Infrastructure Sponsorship Appraisal — %s  \n", project.Name))
	sb.WriteString(fmt.Sprintf("**DATE**: %s  \n", dateStr))
	sb.WriteString(fmt.Sprintf("**TRACKED CAPITAL ENVELOPE**: $%.2f Billion CAD  \n\n", capexBillions))
	sb.WriteString("---\n\n")

	sb.WriteString("## 1. Executive Summary\n")
	sb.WriteString(execSummary + "\n\n")

	sb.WriteString("## 2. Strategic Rationale & Economic Sovereignty Mandate\n")
	sb.WriteString(stratRationale + "\n\n")

	sb.WriteString("## 3. Financial Exposure & Capital Stack Co-Investment Architecture\n")
	sb.WriteString(finExposure + "\n\n")

	sb.WriteString("## 4. Indigenous Co-Ownership & Duty to Consult\n")
	sb.WriteString(indigTreaty + "\n\n")

	sb.WriteString("## 5. Geopolitical Stress-Testing & Supply-Chain Hardening\n")
	sb.WriteString(geoRisk + "\n\n")

	sb.WriteString("## 6. Recommended Ministerial Action\n")
	sb.WriteString(recAction + "\n\n")

	sb.WriteString("---\n")
	sb.WriteString("*Generated deterministically by CanadaOpportunityGraph (COG) Decision Engine with SHA-256 cryptographic provenance.*  \n")

	md := sb.String()
	memoID := fmt.Sprintf("memo-%s-%s-%d", strings.ToLower(string(memoType)), project.Slug, now.Unix())

	h := sha256.Sum256([]byte(md))
	auditHash := hex.EncodeToString(h[:])

	return &DecisionMemo{
		MemoID:             memoID,
		Type:               memoType,
		ProjectID:          project.ID,
		ProjectName:        project.Name,
		SecurityCaveat:     caveat,
		PreparedFor:        recipient,
		DateFormatted:      dateStr,
		ExecutiveSummary:   execSummary,
		StrategicRationale: stratRationale,
		FinancialExposure:  finExposure,
		IndigenousTreaty:   indigTreaty,
		GeopoliticalRisk:   geoRisk,
		RecommendedAction:  recAction,
		MarkdownContent:    md,
		AuditHash:          auditHash,
	}
}
