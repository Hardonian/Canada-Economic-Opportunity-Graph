package nationalplanning

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// ObjectiveType specifies the optimization target.
type ObjectiveType string

const (
	ObjectiveMaxCrowdingIn    ObjectiveType = "MAX_CROWDING_IN"
	ObjectiveMaxSovereignty   ObjectiveType = "MAX_SOVEREIGNTY"
	ObjectiveMaxDecarb        ObjectiveType = "MAX_DECARBONIZATION"
	ObjectiveBalancedStrategy ObjectiveType = "BALANCED"
)

// BudgetEnvelopes holds public fiscal funding limits in CAD.
type BudgetEnvelopes struct {
	CIBConcessionaryCAD int64 `json:"cib_concessionary_cad"`
	SIFGrantsCAD        int64 `json:"sif_grants_cad"`
	ITCTaxCreditsCAD    int64 `json:"itc_tax_credits_cad"`
	IndigenousLoansCAD  int64 `json:"indigenous_loans_cad"`
}

// OptimizationRequest specifies constraints and objectives for capital deployment.
type OptimizationRequest struct {
	Objective         ObjectiveType   `json:"objective"`
	Envelopes         BudgetEnvelopes `json:"envelopes"`
	MinProvincialCap  float64         `json:"min_provincial_cap"`  // Minimum % required for any funded province
	TargetSectors     []domain.Sector `json:"target_sectors,omitempty"`
}

// ProjectAllocation details public co-investment tranches for a selected project.
type ProjectAllocation struct {
	ProjectID         string        `json:"project_id"`
	ProjectName       string        `json:"project_name"`
	Province          string        `json:"province"`
	Sector            domain.Sector `json:"sector"`
	TotalCapexCAD     int64         `json:"total_capex_cad"`
	CIBAllocatedCAD   int64         `json:"cib_allocated_cad"`
	SIFAllocatedCAD   int64         `json:"sif_allocated_cad"`
	ITCAllocatedCAD   int64         `json:"itc_allocated_cad"`
	IndigAllocatedCAD int64         `json:"indig_allocated_cad"`
	TotalPublicCAD    int64         `json:"total_public_cad"`
	PrivateMobilized  int64         `json:"private_mobilized_cad"`
	UtilityScore      float64       `json:"utility_score"`
}

// OptimizationResult summarizes the optimal sovereign portfolio.
type OptimizationResult struct {
	RequestID               string              `json:"request_id"`
	Objective               ObjectiveType       `json:"objective"`
	AllocatedProjects       []ProjectAllocation `json:"allocated_projects"`
	TotalPublicInvestedCAD  int64               `json:"total_public_invested_cad"`
	TotalPrivateMobilizedCAD int64              `json:"total_private_mobilized_cad"`
	CrowdingInMultiplier    float64             `json:"crowding_in_multiplier"`
	TotalGHGAbatedMtPerYear float64             `json:"total_ghg_abated_mt_yr"`
	ProvincialAllocations   map[string]int64    `json:"provincial_allocations"`
	SectorAllocations       map[string]int64    `json:"sector_allocations"`
	RemainingEnvelopes      BudgetEnvelopes     `json:"remaining_envelopes"`
	AuditHash               string              `json:"audit_hash"`
	SolvedAt                time.Time           `json:"solved_at"`
}

// Optimizer implements mathematical portfolio allocation over the economic graph.
type Optimizer struct{}

func NewOptimizer() *Optimizer {
	return &Optimizer{}
}

type scoredProject struct {
	project *domain.Project
	score   float64
	needed  BudgetEnvelopes
}

// Optimize solves for the highest utility capital allocation respecting budget limits.
func (opt *Optimizer) Optimize(projects []*domain.Project, req OptimizationRequest) *OptimizationResult {
	if req.Objective == "" {
		req.Objective = ObjectiveBalancedStrategy
	}

	// Default envelopes if unset
	if req.Envelopes.CIBConcessionaryCAD == 0 && req.Envelopes.SIFGrantsCAD == 0 {
		req.Envelopes = BudgetEnvelopes{
			CIBConcessionaryCAD: 4_000_000_000,
			SIFGrantsCAD:        2_000_000_000,
			ITCTaxCreditsCAD:    3_000_000_000,
			IndigenousLoansCAD:  1_000_000_000,
		}
	}

	remCIB := req.Envelopes.CIBConcessionaryCAD
	remSIF := req.Envelopes.SIFGrantsCAD
	remITC := req.Envelopes.ITCTaxCreditsCAD
	remIndig := req.Envelopes.IndigenousLoansCAD

	// Score candidates
	scored := make([]scoredProject, 0, len(projects))
	for _, p := range projects {
		if p.CapexCAD <= 0 {
			continue
		}

		// Calculate utility score based on objective
		var u float64
		sovScore := p.Scores["strategicity"]
		if sovScore == 0 {
			sovScore = 70.0
		}
		buildScore := p.Scores["buildability"]
		if buildScore == 0 {
			buildScore = 65.0
		}

		capexBillions := float64(p.CapexCAD) / 1e9

		switch req.Objective {
		case ObjectiveMaxCrowdingIn:
			u = (capexBillions * 2.0) + (buildScore * 0.5)
		case ObjectiveMaxSovereignty:
			u = (sovScore * 2.5) + (capexBillions * 0.5)
		case ObjectiveMaxDecarb:
			decarbBonus := 10.0
			if p.Sector == domain.SectorNuclearEnergy || p.Sector == domain.SectorCleanEnergy {
				decarbBonus = 80.0
			}
			u = decarbBonus + (buildScore * 0.5)
		default: // Balanced
			u = (capexBillions * 1.0) + (sovScore * 1.0) + (buildScore * 1.0)
		}

		// Estimate standard required public support tranches (typically 25%-35% total public support)
		cibNeed := int64(math.Round(float64(p.CapexCAD) * 0.15))
		sifNeed := int64(math.Round(float64(p.CapexCAD) * 0.05))
		itcNeed := int64(math.Round(float64(p.CapexCAD) * 0.10))
		indigNeed := int64(math.Round(float64(p.CapexCAD) * 0.05))

		scored = append(scored, scoredProject{
			project: p,
			score:   u,
			needed: BudgetEnvelopes{
				CIBConcessionaryCAD: cibNeed,
				SIFGrantsCAD:        sifNeed,
				ITCTaxCreditsCAD:    itcNeed,
				IndigenousLoansCAD:  indigNeed,
			},
		})
	}

	// Sort descending by utility score
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	allocations := make([]ProjectAllocation, 0)
	var totalPub, totalPriv int64
	provDist := make(map[string]int64)
	sectorDist := make(map[string]int64)
	var ghgAbated float64

	for _, item := range scored {
		p := item.project
		need := item.needed

		// Check if we can partially or fully allocate
		if remCIB < need.CIBConcessionaryCAD && remSIF < need.SIFGrantsCAD {
			continue // Insufficient public envelope
		}

		allocCIB := int64(math.Min(float64(remCIB), float64(need.CIBConcessionaryCAD)))
		allocSIF := int64(math.Min(float64(remSIF), float64(need.SIFGrantsCAD)))
		allocITC := int64(math.Min(float64(remITC), float64(need.ITCTaxCreditsCAD)))
		allocIndig := int64(math.Min(float64(remIndig), float64(need.IndigenousLoansCAD)))

		pubThis := allocCIB + allocSIF + allocITC + allocIndig
		if pubThis <= 0 {
			continue
		}

		privThis := p.CapexCAD - pubThis
		if privThis < 0 {
			privThis = 0
		}

		remCIB -= allocCIB
		remSIF -= allocSIF
		remITC -= allocITC
		remIndig -= allocIndig

		totalPub += pubThis
		totalPriv += privThis
		provDist[p.Province] += pubThis
		sectorDist[string(p.Sector)] += pubThis

		// GHG abatement estimate: clean energy / nuclear offsets ~1.2 Mt CO2e per $1B CAPEX
		if p.Sector == domain.SectorNuclearEnergy || p.Sector == domain.SectorCleanEnergy {
			ghgAbated += (float64(p.CapexCAD) / 1e9) * 1.25
		}

		allocations = append(allocations, ProjectAllocation{
			ProjectID:         p.ID,
			ProjectName:       p.Name,
			Province:          p.Province,
			Sector:            p.Sector,
			TotalCapexCAD:     p.CapexCAD,
			CIBAllocatedCAD:   allocCIB,
			SIFAllocatedCAD:   allocSIF,
			ITCAllocatedCAD:   allocITC,
			IndigAllocatedCAD: allocIndig,
			TotalPublicCAD:    pubThis,
			PrivateMobilized:  privThis,
			UtilityScore:      math.Round(item.score*10) / 10,
		})
	}

	mult := 0.0
	if totalPub > 0 {
		mult = float64(totalPriv) / float64(totalPub)
	}

	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%d", req.Objective, totalPub, len(allocations))))
	reqID := fmt.Sprintf("opt_%x", sum[:8])

	result := &OptimizationResult{
		RequestID:                reqID,
		Objective:                req.Objective,
		AllocatedProjects:        allocations,
		TotalPublicInvestedCAD:   totalPub,
		TotalPrivateMobilizedCAD: totalPriv,
		CrowdingInMultiplier:     math.Round(mult*10) / 10,
		TotalGHGAbatedMtPerYear:  math.Round(ghgAbated*10) / 10,
		ProvincialAllocations:    provDist,
		SectorAllocations:        sectorDist,
		RemainingEnvelopes: BudgetEnvelopes{
			CIBConcessionaryCAD: remCIB,
			SIFGrantsCAD:        remSIF,
			ITCTaxCreditsCAD:    remITC,
			IndigenousLoansCAD:  remIndig,
		},
		SolvedAt: time.Now().UTC(),
	}

	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%d|%d|%.2f", result.RequestID, result.TotalPublicInvestedCAD, result.TotalPrivateMobilizedCAD, result.CrowdingInMultiplier)))
	result.AuditHash = hex.EncodeToString(h[:])

	return result
}
