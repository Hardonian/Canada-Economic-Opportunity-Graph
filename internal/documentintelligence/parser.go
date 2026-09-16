// Package documentintelligence extracts restricted candidate records without
// promoting them into the public project graph. It is intentionally a bounded,
// deterministic parser suitable for evaluation fixtures and human review.
package documentintelligence

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

const ExtractorVersion = "portfolio-cards-v1.0"

var (
	cardBoundary = regexp.MustCompile(`(?im)^\s*---\s*PROJECT\s*---\s*$`)
	fieldLine    = regexp.MustCompile(`(?m)^([A-Za-z][A-Za-z ]{1,40}):\s*(.+)$`)
	amountToken  = regexp.MustCompile(`(?i)(US\$|USD\s*|C\$|CAD\s*|\$)\s*([0-9]+(?:\.[0-9]+)?)\s*([BM])`)
	phaseToken   = regexp.MustCompile(`(?i)\b(phase\s+[0-9A-Za-z]+)\b`)
)

type ExtractedCard struct {
	Candidate          *domain.CandidateProject     `json:"candidate"`
	CapitalRequirement *domain.CapitalRequirement  `json:"capital_requirement,omitempty"`
	CapitalNeed        *domain.CapitalNeed          `json:"capital_need,omitempty"`
}

func ExtractCards(text, sourceID string, visibility domain.VisibilityClass, observedAt time.Time) ([]ExtractedCard, error) {
	if !visibility.Valid() || visibility.Public() {
		return nil, fmt.Errorf("portfolio extraction requires an explicit private or restricted visibility")
	}
	parts := cardBoundary.Split(strings.TrimSpace(text), -1)
	results := make([]ExtractedCard, 0, len(parts))
	for _, part := range parts {
		fields := parseFields(part)
		name := strings.TrimSpace(first(fields, "project", "project name", "name"))
		if name == "" {
			continue
		}
		id := stableID(sourceID + "|" + name)
		candidate := &domain.CandidateProject{
			ID: "candidate-" + id, SourceID: sourceID, Name: name,
			Proponent: first(fields, "proponent", "developer"), Location: first(fields, "location", "geography"),
			OriginalStage: first(fields, "stage", "development stage"), Visibility: visibility,
			CreatedAt: observedAt.UTC(),
		}
		card := ExtractedCard{Candidate: candidate}
		if raw := first(fields, "capex", "capital cost", "project cost"); raw != "" {
			amount := ParseCapex(raw)
			phase := ""
			if match := phaseToken.FindStringSubmatch(raw); len(match) > 1 {
				phase = strings.TrimSpace(match[1])
			}
			card.CapitalRequirement = &domain.CapitalRequirement{
				ID: "capital-requirement-" + id, ProjectID: candidate.ID, Purpose: phase,
				Amount: amount, Status: domain.ConfidenceReported, Visibility: visibility,
				Publishable: false, CreatedAt: observedAt.UTC(),
			}
		}
		if raw := first(fields, "financing objective", "capital need", "partners sought"); raw != "" {
			card.CapitalNeed = &domain.CapitalNeed{
				ID: "capital-need-" + id, ProjectID: candidate.ID, Types: ClassifyCapitalNeeds(raw),
				Counterparties: ClassifyCounterparties(raw), Status: domain.CapitalSeeking,
				OriginalLanguage: raw, Visibility: visibility, Publishable: false,
				PublicationState: domain.PublicationPrivateOnly, CreatedAt: observedAt.UTC(), UpdatedAt: observedAt.UTC(),
			}
		}
		results = append(results, card)
	}
	return results, nil
}

func parseFields(text string) map[string]string {
	fields := map[string]string{}
	for _, match := range fieldLine.FindAllStringSubmatch(text, -1) {
		fields[strings.ToLower(strings.TrimSpace(match[1]))] = strings.TrimSpace(match[2])
	}
	return fields
}

func first(fields map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := fields[key]; value != "" {
			return value
		}
	}
	return ""
}

func ParseCapex(raw string) domain.MonetaryAmount {
	result := domain.MonetaryAmount{AmountType: domain.AmountNotAvailable, OriginalText: strings.TrimSpace(raw)}
	lower := strings.ToLower(strings.TrimSpace(raw))
	if lower == "" || strings.Contains(lower, "not available") || lower == "n/a" || lower == "unknown" {
		return result
	}
	matches := amountToken.FindAllStringSubmatch(raw, -1)
	if len(matches) == 0 {
		return result
	}
	currency := parseCurrency(matches[0][1])
	values := make([]int64, 0, len(matches))
	for _, match := range matches {
		if parseCurrency(match[1]) != currency {
			return result // mixed-currency text requires human review
		}
		value, _ := strconv.ParseFloat(match[2], 64)
		multiplier := float64(1_000_000)
		if strings.EqualFold(match[3], "B") {
			multiplier = 1_000_000_000
		}
		values = append(values, int64(value*multiplier+0.5))
	}
	result.Currency = currency
	if len(values) > 1 {
		min, max := values[0], values[1]
		if min > max {
			min, max = max, min
		}
		result.Minimum, result.Maximum, result.AmountType = &min, &max, domain.AmountRange
		return result
	}
	value := values[0]
	result.Amount = &value
	switch {
	case strings.Contains(lower, "+") || strings.Contains(lower, "at least"):
		result.AmountType = domain.AmountMinimum
		result.Minimum = &value
	case strings.Contains(lower, "approximately") || strings.Contains(lower, "approx.") || strings.Contains(lower, "about") || strings.Contains(lower, "~"):
		result.AmountType = domain.AmountApproximate
	case strings.Contains(lower, "up to"):
		result.AmountType = domain.AmountMaximum
		result.Maximum = &value
	default:
		result.AmountType = domain.AmountExact
	}
	return result
}

func parseCurrency(prefix string) string {
	upper := strings.ToUpper(strings.TrimSpace(prefix))
	switch upper {
	case "US$", "USD":
		return "USD"
	case "C$", "CAD":
		return "CAD"
	default:
		return "UNSPECIFIED"
	}
}

func ClassifyCapitalNeeds(raw string) []domain.CapitalNeedType {
	lower := strings.ToLower(raw)
	mapping := []struct {
		phrases []string
		kind    domain.CapitalNeedType
	}{
		{[]string{"infrastructure equity"}, domain.NeedInfrastructureEquity},
		{[]string{"project finance", "project financing"}, domain.NeedProjectFinance},
		{[]string{"private credit"}, domain.NeedPrivateCredit},
		{[]string{"joint venture", " jv "}, domain.NeedJointVenture},
		{[]string{"strategic invest"}, domain.NeedStrategicInvestment},
		{[]string{"government co-invest", "government support"}, domain.NeedGovernmentSupport},
		{[]string{"loan guarantee"}, domain.NeedLoanGuarantee},
		{[]string{"export credit", "eca financing"}, domain.NeedExportCredit},
		{[]string{"indigenous equity"}, domain.NeedIndigenousEquity},
		{[]string{"pension capital"}, domain.NeedPensionCapital},
		{[]string{"sovereign capital"}, domain.NeedSovereignCapital},
		{[]string{"offtake"}, domain.NeedOfftake},
		{[]string{"anchor tenant"}, domain.NeedAnchorTenant},
		{[]string{"prepayment"}, domain.NeedPrepayment},
		{[]string{"streaming"}, domain.NeedStreaming},
		{[]string{"royalty"}, domain.NeedRoyalty},
		{[]string{"senior debt"}, domain.NeedSeniorDebt},
		{[]string{"subordinated debt", "mezzanine"}, domain.NeedSubordinatedDebt},
		{[]string{"development capital"}, domain.NeedDevelopmentCapital},
		{[]string{"equity"}, domain.NeedEquity},
		{[]string{"debt"}, domain.NeedDebt},
	}
	seen := map[domain.CapitalNeedType]bool{}
	result := []domain.CapitalNeedType{}
	for _, item := range mapping {
		for _, phrase := range item.phrases {
			if strings.Contains(" "+lower+" ", phrase) && !seen[item.kind] {
				seen[item.kind] = true
				result = append(result, item.kind)
				break
			}
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func ClassifyCounterparties(raw string) []domain.CounterpartyType {
	lower := strings.ToLower(raw)
	mapping := map[string]domain.CounterpartyType{
		"infrastructure fund": domain.CounterpartyInfrastructureFund,
		"pension": domain.CounterpartyPensionFund,
		"private equity": domain.CounterpartyPrivateEquity,
		"bank": domain.CounterpartyBank,
		"private credit": domain.CounterpartyPrivateCredit,
		"export credit": domain.CounterpartyECA,
		"sovereign": domain.CounterpartySovereignFund,
		"strategic partner": domain.CounterpartyStrategicCorporate,
		"epc": domain.CounterpartyEPC,
		"offtaker": domain.CounterpartyOfftaker,
		"anchor tenant": domain.CounterpartyAnchorTenant,
		"indigenous partner": domain.CounterpartyIndigenousPartner,
		"joint venture": domain.CounterpartyJVPartner,
	}
	seen := map[domain.CounterpartyType]bool{}
	result := []domain.CounterpartyType{}
	for phrase, kind := range mapping {
		if strings.Contains(lower, phrase) && !seen[kind] {
			seen[kind] = true
			result = append(result, kind)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result
}

func NormalizeStages(raw string) []domain.LifecycleStage {
	lower := strings.ToLower(raw)
	mapping := []struct {
		phrase string
		stage  domain.LifecycleStage
	}{
		{"under construction", domain.StageConstruction}, {"construction-ready", domain.StageConstructionReady},
		{"shovel-ready", domain.StageConstructionReady}, {"pre-fid", domain.StageFIDLikely},
		{"detailed engineering", domain.StageDetailedEngineering}, {"pre-feed", domain.StagePreFEED},
		{"feed", domain.StageFEED}, {"pre-feasibility", domain.StagePreDevelopment},
		{"feasibility", domain.StageFeasibility}, {"permitting", domain.StagePermitting},
		{"financing", domain.StageFinancing}, {"commissioning", domain.StageCommissioning},
		{"operational", domain.StageOperating}, {"operating", domain.StageOperating},
		{"expansion", domain.StageExpansion}, {"concept", domain.StageConcept},
	}
	seen := map[domain.LifecycleStage]bool{}
	result := []domain.LifecycleStage{}
	for _, item := range mapping {
		if strings.Contains(lower, item.phrase) && !seen[item.stage] {
			seen[item.stage] = true
			result = append(result, item.stage)
		}
	}
	return result
}

func stableID(value string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(value))))
	return hex.EncodeToString(sum[:8])
}

// TechnicalReport43101 holds parsed parameters from a National Instrument 43-101 technical filing.
type TechnicalReport43101 struct {
	ProjectName      string    `json:"project_name"`
	Commodity        string    `json:"commodity"`
	ReserveTonnageMt float64   `json:"reserve_tonnage_mt"`
	Grade            string    `json:"grade"`
	RecoveryRatePct  float64   `json:"recovery_rate_pct"`
	StripRatio       float64   `json:"strip_ratio"`
	MineLifeYears    int       `json:"mine_life_years"`
	InitialCapexCAD  int64     `json:"initial_capex_cad"`
	AfterTaxNPV8CAD  int64     `json:"after_tax_npv8_cad"`
	AfterTaxIRRPct   float64   `json:"after_tax_irr_pct"`
	AuthorFirm       string    `json:"author_firm"`
	EffectiveDate    time.Time `json:"effective_date"`
}

var (
	tonnageRegex  = regexp.MustCompile(`(?i)(?:reserve|resource|tonnage|deposit)[:\s]+([0-9]+(?:\.[0-9]+)?)\s*(?:mt|million tonnes|million tons)`)
	gradeRegex    = regexp.MustCompile(`(?i)(?:grade|average grade)[:\s]+([0-9]+(?:\.[0-9]+)?%?\s*[A-Za-z0-9_]+)`)
	recoveryRegex = regexp.MustCompile(`(?i)(?:recovery|metallurgical recovery)[:\s]+([0-9]+(?:\.[0-9]+)?)\s*%`)
	mineLifeRegex = regexp.MustCompile(`(?i)(?:mine life|life of mine|lom)[:\s]+([0-9]+)\s*years`)
	npvRegex      = regexp.MustCompile(`(?i)(?:npv8%?|after-tax npv)[:\s]+(?:US\$|USD|C\$|CAD|\$)?\s*([0-9]+(?:\.[0-9]+)?)\s*([BM])`)
	irrRegex      = regexp.MustCompile(`(?i)(?:after-tax irr|irr)[:\s]+([0-9]+(?:\.[0-9]+)?)\s*%`)
)

// ExtractNI43101TechnicalReport parses mining reserves and economic metrics from filing text.
func ExtractNI43101TechnicalReport(text string) (*TechnicalReport43101, error) {
	lines := strings.Split(text, "\n")
	fields := make(map[string]string)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if parts := strings.SplitN(l, ":", 2); len(parts) == 2 {
			fields[strings.ToLower(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
		}
	}

	report := &TechnicalReport43101{
		ProjectName:   first(fields, "project", "property", "project name"),
		Commodity:     first(fields, "commodity", "primary metal", "mineral"),
		Grade:         first(fields, "grade", "average grade"),
		AuthorFirm:    first(fields, "author", "qp", "engineering firm", "consultant"),
		EffectiveDate: time.Now().UTC(),
	}

	if match := tonnageRegex.FindStringSubmatch(text); len(match) > 1 {
		if val, err := strconv.ParseFloat(match[1], 64); err == nil {
			report.ReserveTonnageMt = val
		}
	}

	if report.Grade == "" {
		if match := gradeRegex.FindStringSubmatch(text); len(match) > 1 {
			report.Grade = strings.TrimSpace(match[1])
		}
	}

	if match := recoveryRegex.FindStringSubmatch(text); len(match) > 1 {
		if val, err := strconv.ParseFloat(match[1], 64); err == nil {
			report.RecoveryRatePct = val
		}
	}

	if match := mineLifeRegex.FindStringSubmatch(text); len(match) > 1 {
		if val, err := strconv.Atoi(match[1]); err == nil {
			report.MineLifeYears = val
		}
	}

	if match := npvRegex.FindStringSubmatch(text); len(match) > 2 {
		if val, err := strconv.ParseFloat(match[1], 64); err == nil {
			mult := 1_000_000.0
			if strings.EqualFold(match[2], "b") {
				mult = 1_000_000_000.0
			}
			report.AfterTaxNPV8CAD = int64(math.Round(val * mult))
		}
	}

	if match := irrRegex.FindStringSubmatch(text); len(match) > 1 {
		if val, err := strconv.ParseFloat(match[1], 64); err == nil {
			report.AfterTaxIRRPct = val
		}
	}

	if rawCapex := first(fields, "capex", "initial capex"); rawCapex != "" {
		amount := ParseCapex(rawCapex)
		report.InitialCapexCAD = amount.AmountInCents / 100 // whole CAD
	}

	return report, nil
}

// CapitalWaterfallDossier models parsed financing tranches and weighted cost of capital.
type CapitalWaterfallDossier struct {
	SeniorDebtCAD        int64   `json:"senior_debt_cad"`
	SponsorEquityCAD     int64   `json:"sponsor_equity_cad"`
	ConcessionaryDebtCAD int64   `json:"concessionary_debt_cad"` // CIB
	TaxCreditEquityCAD   int64   `json:"tax_credit_equity_cad"`   // ITCs
	IndigenousEquityCAD  int64   `json:"indigenous_equity_cad"`
	TotalCapexCAD        int64   `json:"total_capex_cad"`
	BlendedWACCPct       float64 `json:"blended_wacc_pct"`
	IsBalanced           bool    `json:"is_balanced"`
}

var trancheRegex = regexp.MustCompile(`(?i)([A-Za-z ]+):\s*(?:C\$|CAD|\$)?\s*([0-9]+(?:\.[0-9]+)?)\s*([BM])`)

// ExtractCapitalWaterfall parses capital stack breakdowns and computes WACC.
func ExtractCapitalWaterfall(text string, reportedCapexCAD int64) *CapitalWaterfallDossier {
	waterfall := &CapitalWaterfallDossier{
		TotalCapexCAD: reportedCapexCAD,
	}

	matches := trancheRegex.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		name := strings.ToLower(strings.TrimSpace(m[1]))
		val, err := strconv.ParseFloat(m[2], 64)
		if err != nil {
			continue
		}
		mult := 1_000_000.0
		if strings.EqualFold(m[3], "b") {
			mult = 1_000_000_000.0
		}
		cad := int64(math.Round(val * mult))

		switch {
		case strings.Contains(name, "senior debt") || strings.Contains(name, "bank debt"):
			waterfall.SeniorDebtCAD += cad
		case strings.Contains(name, "sponsor equity") || strings.Contains(name, "equity"):
			waterfall.SponsorEquityCAD += cad
		case strings.Contains(name, "cib") || strings.Contains(name, "concessionary"):
			waterfall.ConcessionaryDebtCAD += cad
		case strings.Contains(name, "itc") || strings.Contains(name, "tax credit"):
			waterfall.TaxCreditEquityCAD += cad
		case strings.Contains(name, "indigenous"):
			waterfall.IndigenousEquityCAD += cad
		}
	}

	totalIdentified := waterfall.SeniorDebtCAD + waterfall.SponsorEquityCAD + waterfall.ConcessionaryDebtCAD +
		waterfall.TaxCreditEquityCAD + waterfall.IndigenousEquityCAD

	if reportedCapexCAD > 0 && totalIdentified == reportedCapexCAD {
		waterfall.IsBalanced = true
	}

	// Calculate WACC
	// Typical costs: Senior Debt 6.5%, Sponsor Equity 14.0%, Concessionary 4.5%, ITCs 0%, Indigenous 4.0%
	if totalIdentified > 0 {
		totF := float64(totalIdentified)
		wacc := (float64(waterfall.SeniorDebtCAD)/totF)*6.5 +
			(float64(waterfall.SponsorEquityCAD)/totF)*14.0 +
			(float64(waterfall.ConcessionaryDebtCAD)/totF)*4.5 +
			(float64(waterfall.TaxCreditEquityCAD)/totF)*0.0 +
			(float64(waterfall.IndigenousEquityCAD)/totF)*4.0
		waterfall.BlendedWACCPct = math.Round(wacc*100) / 100
	}

	return waterfall
}

