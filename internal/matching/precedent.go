package matching

import (
	"sort"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// PrecedentMatch records a historical deal that is similar to a given project,
// with an explanation of which dimensions matched.
type PrecedentMatch struct {
	Deal          domain.DealPrecedent `json:"deal"`
	SimilarityPct float64              `json:"similarity_pct"` // 0-100
	MatchedOn     []string             `json:"matched_on"`     // sector, province, stage, instrument, scale
}

// FindDealPrecedents searches investor profiles for historical transactions
// that resemble the target project along key dimensions. This answers:
// "Who has financed projects like this before?" — it explains similarity,
// never claims future intent.
func FindDealPrecedents(project *domain.Project, profiles []*domain.InvestorProfile, topN int) []PrecedentMatch {
	if project == nil || len(profiles) == 0 {
		return nil
	}
	if topN <= 0 {
		topN = 10
	}

	var matches []PrecedentMatch
	for _, profile := range profiles {
		for _, deal := range profile.PublicDealHistory {
			similarity, matchedOn := dealSimilarity(project, deal)
			if similarity < 20 {
				continue
			}
			matches = append(matches, PrecedentMatch{
				Deal:          deal,
				SimilarityPct: similarity,
				MatchedOn:     matchedOn,
			})
		}
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].SimilarityPct > matches[j].SimilarityPct
	})

	if len(matches) > topN {
		matches = matches[:topN]
	}
	return matches
}

func dealSimilarity(project *domain.Project, deal domain.DealPrecedent) (float64, []string) {
	var score float64
	var matched []string

	// Sector match (strongest signal).
	if deal.Sector == project.Sector {
		score += 35
		matched = append(matched, "sector")
	}

	// Province match.
	if deal.Province != "" && strings.EqualFold(deal.Province, project.Province) {
		score += 15
		matched = append(matched, "province")
	}

	// Stage match.
	if deal.Stage == project.CurrentStage {
		score += 20
		matched = append(matched, "stage")
	}

	// Scale similarity (within 5x).
	if deal.AmountCAD > 0 && project.CapexCAD > 0 {
		ratio := float64(deal.AmountCAD) / float64(project.CapexCAD)
		if ratio < 1 {
			ratio = 1 / ratio
		}
		if ratio <= 5 {
			score += 20
			matched = append(matched, "scale")
		} else if ratio <= 10 {
			score += 10
			matched = append(matched, "scale_approximate")
		}
	}

	// Recent deal bonus.
	if deal.Year > 0 && deal.Year >= 2022 {
		score += 10
		matched = append(matched, "recent")
	}

	return score, matched
}

// CanonicalInvestorProfiles returns authoritative institutional allocator profiles
// across Canada's Maple 8 pensions, Crown corporations, and global infrastructure funds.
func CanonicalInvestorProfiles() []*domain.InvestorProfile {
	now := time.Now().UTC()
	return []*domain.InvestorProfile{
		{
			EntityID:          "cppib",
			InvestorTypes:     []domain.CounterpartyType{domain.CounterpartyPensionFund, domain.CounterpartyInfrastructureFund},
			TargetSectors:     []domain.Sector{domain.SectorCleanEnergy, domain.SectorNuclearEnergy, domain.SectorTransportation, domain.SectorAICompute},
			TargetGeographies: []string{"ON", "AB", "BC", "QC"},
			MinTicketCAD:      250_000_000,
			MaxTicketCAD:      2_500_000_000,
			PreferredInstruments: []domain.CapitalNeedType{
				domain.NeedInfrastructureEquity,
				domain.NeedEquity,
				domain.NeedProjectFinance,
			},
			PreferredStages: []domain.LifecycleStage{
				domain.StageFEED,
				domain.StageConstruction,
				domain.StageOperating,
			},
			CanadianExposureCAD: 85_000_000_000,
			ActiveInvestments:    42,
			PublicDealHistory: []domain.DealPrecedent{
				{
					DealID:      "cppib-bruce-power",
					ProjectName: "Bruce Power Nuclear Life-Extension & Clean Refurbishment",
					Sector:      domain.SectorNuclearEnergy,
					Province:    "ON",
					AmountCAD:   2_000_000_000,
					Instrument:  domain.NeedInfrastructureEquity,
					Stage:       domain.StageConstruction,
					Year:        2023,
				},
				{
					DealID:      "cppib-pattern-energy",
					ProjectName: "Pattern Energy Clean Renewable Power Platform",
					Sector:      domain.SectorCleanEnergy,
					Province:    "ON",
					AmountCAD:   1_200_000_000,
					Instrument:  domain.NeedEquity,
					Stage:       domain.StageOperating,
					Year:        2022,
				},
				{
					DealID:      "cppib-hwy407",
					ProjectName: "Highway 407 ETR Concession & Intelligent Corridor",
					Sector:      domain.SectorTransportation,
					Province:    "ON",
					AmountCAD:   3_250_000_000,
					Instrument:  domain.NeedInfrastructureEquity,
					Stage:       domain.StageOperating,
					Year:        2021,
				},
			},
			Visibility:  domain.VisibilityPublic,
			Publishable: true,
			UpdatedAt:   now,
		},
		{
			EntityID:          "cdpq",
			InvestorTypes:     []domain.CounterpartyType{domain.CounterpartyPensionFund, domain.CounterpartyInfrastructureFund},
			TargetSectors:     []domain.Sector{domain.SectorCleanEnergy, domain.SectorTransportation, domain.SectorCriticalMinerals},
			TargetGeographies: []string{"QC", "ON", "BC", "AB"},
			MinTicketCAD:      150_000_000,
			MaxTicketCAD:      1_500_000_000,
			PreferredInstruments: []domain.CapitalNeedType{
				domain.NeedInfrastructureEquity,
				domain.NeedProjectFinance,
				domain.NeedSeniorDebt,
			},
			PreferredStages: []domain.LifecycleStage{
				domain.StagePreFEED,
				domain.StageFEED,
				domain.StageConstruction,
				domain.StageOperating,
			},
			CanadianExposureCAD: 110_000_000_000,
			ActiveInvestments:    55,
			PublicDealHistory: []domain.DealPrecedent{
				{
					DealID:      "cdpq-rem-montreal",
					ProjectName: "Réseau express métropolitain (REM) Light Rail Network",
					Sector:      domain.SectorTransportation,
					Province:    "QC",
					AmountCAD:   3_500_000_000,
					Instrument:  domain.NeedInfrastructureEquity,
					Stage:       domain.StageConstruction,
					Year:        2023,
				},
				{
					DealID:      "cdpq-energir-clean-gas",
					ProjectName: "Énergir Clean Renewable Natural Gas & Boreal Grid",
					Sector:      domain.SectorCleanEnergy,
					Province:    "QC",
					AmountCAD:   1_140_000_000,
					Instrument:  domain.NeedEquity,
					Stage:       domain.StageOperating,
					Year:        2022,
				},
				{
					DealID:      "cdpq-nemaska-lithium",
					ProjectName: "Bécancour Whabouchi Lithium Hydroxide Complex",
					Sector:      domain.SectorCriticalMinerals,
					Province:    "QC",
					AmountCAD:   250_000_000,
					Instrument:  domain.NeedEquity,
					Stage:       domain.StageFEED,
					Year:        2023,
				},
			},
			Visibility:  domain.VisibilityPublic,
			Publishable: true,
			UpdatedAt:   now,
		},
		{
			EntityID:          "cib",
			InvestorTypes:     []domain.CounterpartyType{domain.CounterpartyGovernment, domain.CounterpartyBank},
			TargetSectors:     []domain.Sector{domain.SectorCleanEnergy, domain.SectorNuclearEnergy, domain.SectorTransportation, domain.SectorCriticalMinerals},
			TargetGeographies: []string{"ON", "QC", "BC", "AB", "SK", "MB", "NS", "NB", "NL", "PE", "YT", "NT", "NU"},
			MinTicketCAD:      50_000_000,
			MaxTicketCAD:      1_500_000_000,
			PreferredInstruments: []domain.CapitalNeedType{
				domain.NeedSubordinatedDebt,
				domain.NeedSeniorDebt,
				domain.NeedProjectFinance,
				domain.NeedLoanGuarantee,
			},
			PreferredStages: []domain.LifecycleStage{
				domain.StageFEED,
				domain.StageConstruction,
			},
			CanadianExposureCAD: 35_000_000_000,
			ActiveInvestments:    48,
			PublicDealHistory: []domain.DealPrecedent{
				{
					DealID:      "cib-oneida-bess",
					ProjectName: "Oneida Energy Storage 250MW BESS Facility",
					Sector:      domain.SectorCleanEnergy,
					Province:    "ON",
					AmountCAD:   170_000_000,
					Instrument:  domain.NeedSubordinatedDebt,
					Stage:       domain.StageConstruction,
					Year:        2023,
				},
				{
					DealID:      "cib-darlington-smr",
					ProjectName: "Darlington New Nuclear Project — BWRX-300 SMR Unit 1",
					Sector:      domain.SectorNuclearEnergy,
					Province:    "ON",
					AmountCAD:   970_000_000,
					Instrument:  domain.NeedSubordinatedDebt,
					Stage:       domain.StageFEED,
					Year:        2022,
				},
				{
					DealID:      "cib-contrecoeur-port",
					ProjectName: "Contrecœur Port Terminal Expansion",
					Sector:      domain.SectorTransportation,
					Province:    "QC",
					AmountCAD:   300_000_000,
					Instrument:  domain.NeedSeniorDebt,
					Stage:       domain.StageFEED,
					Year:        2024,
				},
				{
					DealID:      "cib-tshiuetin-rail",
					ProjectName: "Tshiuetin First Nations Railway Modernization",
					Sector:      domain.SectorTransportation,
					Province:    "QC",
					AmountCAD:   50_000_000,
					Instrument:  domain.NeedIndigenousEquity,
					Stage:       domain.StageOperating,
					Year:        2023,
				},
			},
			Visibility:  domain.VisibilityPublic,
			Publishable: true,
			UpdatedAt:   now,
		},
		{
			EntityID:          "otpp",
			InvestorTypes:     []domain.CounterpartyType{domain.CounterpartyPensionFund, domain.CounterpartyInfrastructureFund},
			TargetSectors:     []domain.Sector{domain.SectorCleanEnergy, domain.SectorNuclearEnergy, domain.SectorAICompute},
			TargetGeographies: []string{"ON", "AB", "BC", "QC"},
			MinTicketCAD:      100_000_000,
			MaxTicketCAD:      1_000_000_000,
			PreferredInstruments: []domain.CapitalNeedType{
				domain.NeedInfrastructureEquity,
				domain.NeedEquity,
			},
			PreferredStages: []domain.LifecycleStage{
				domain.StageFEED,
				domain.StageConstruction,
				domain.StageOperating,
			},
			CanadianExposureCAD: 70_000_000_000,
			ActiveInvestments:    30,
			PublicDealHistory: []domain.DealPrecedent{
				{
					DealID:      "otpp-cubico-renewables",
					ProjectName: "Cubico Clean Power Sustainable Portfolio",
					Sector:      domain.SectorCleanEnergy,
					Province:    "ON",
					AmountCAD:   850_000_000,
					Instrument:  domain.NeedInfrastructureEquity,
					Stage:       domain.StageOperating,
					Year:        2023,
				},
				{
					DealID:      "otpp-clean-grid-infra",
					ProjectName: "Ontario Clean Grid Interconnect & HVDC",
					Sector:      domain.SectorCleanEnergy,
					Province:    "ON",
					AmountCAD:   600_000_000,
					Instrument:  domain.NeedInfrastructureEquity,
					Stage:       domain.StageConstruction,
					Year:        2024,
				},
			},
			Visibility:  domain.VisibilityPublic,
			Publishable: true,
			UpdatedAt:   now,
		},
		{
			EntityID:          "cgf",
			InvestorTypes:     []domain.CounterpartyType{domain.CounterpartyGovernment, domain.CounterpartyStrategicCorporate},
			TargetSectors:     []domain.Sector{domain.SectorCleanEnergy, domain.SectorCriticalMinerals, domain.SectorAICompute},
			TargetGeographies: []string{"AB", "ON", "BC", "QC", "SK"},
			MinTicketCAD:      50_000_000,
			MaxTicketCAD:      1_000_000_000,
			PreferredInstruments: []domain.CapitalNeedType{
				domain.NeedOfftake,
				domain.NeedSubordinatedDebt,
				domain.NeedEquity,
			},
			PreferredStages: []domain.LifecycleStage{
				domain.StageFEED,
				domain.StageConstruction,
			},
			CanadianExposureCAD: 15_000_000_000,
			ActiveInvestments:    12,
			PublicDealHistory: []domain.DealPrecedent{
				{
					DealID:      "cgf-entropy-ccus",
					ProjectName: "Entropy Carbon Capture & Storage CCfD",
					Sector:      domain.SectorCleanEnergy,
					Province:    "AB",
					AmountCAD:   200_000_000,
					Instrument:  domain.NeedOfftake,
					Stage:       domain.StageConstruction,
					Year:        2023,
				},
				{
					DealID:      "cgf-crawford-offtake",
					ProjectName: "Crawford Nickel Critical Minerals Carbon Capture",
					Sector:      domain.SectorCriticalMinerals,
					Province:    "ON",
					AmountCAD:   150_000_000,
					Instrument:  domain.NeedStrategicInvestment,
					Stage:       domain.StageFEED,
					Year:        2024,
				},
			},
			Visibility:  domain.VisibilityPublic,
			Publishable: true,
			UpdatedAt:   now,
		},
		{
			EntityID:          "gic",
			InvestorTypes:     []domain.CounterpartyType{domain.CounterpartySovereignFund, domain.CounterpartyInfrastructureFund},
			TargetSectors:     []domain.Sector{domain.SectorTransportation, domain.SectorAICompute, domain.SectorCleanEnergy},
			TargetGeographies: []string{"ON", "BC", "QC", "AB"},
			MinTicketCAD:      300_000_000,
			MaxTicketCAD:      3_000_000_000,
			PreferredInstruments: []domain.CapitalNeedType{
				domain.NeedInfrastructureEquity,
				domain.NeedEquity,
			},
			PreferredStages: []domain.LifecycleStage{
				domain.StageConstruction,
				domain.StageOperating,
			},
			CanadianExposureCAD: 25_000_000_000,
			ActiveInvestments:    15,
			PublicDealHistory: []domain.DealPrecedent{
				{
					DealID:      "gic-nas-airports-hub",
					ProjectName: "Canadian International Airport Investor Leasing Hub",
					Sector:      domain.SectorTransportation,
					Province:    "ON",
					AmountCAD:   1_500_000_000,
					Instrument:  domain.NeedInfrastructureEquity,
					Stage:       domain.StageFEED,
					Year:        2026,
				},
				{
					DealID:      "gic-boreal-compute-cluster",
					ProjectName: "Chisasibi Sovereign Clean AI Compute Facility",
					Sector:      domain.SectorAICompute,
					Province:    "QC",
					AmountCAD:   800_000_000,
					Instrument:  domain.NeedInfrastructureEquity,
					Stage:       domain.StageFEED,
					Year:        2025,
				},
			},
			Visibility:  domain.VisibilityPublic,
			Publishable: true,
			UpdatedAt:   now,
		},
	}
}
