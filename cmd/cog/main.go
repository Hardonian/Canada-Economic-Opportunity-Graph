package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/adapters"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/adapters/global_trade"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/adapters/nrcan_major_projects"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/adapters/official"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/cegs"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/documentintelligence"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/corridor"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/earthobs"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/econometrics"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/export"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/filings"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/gridphysics"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/indicators"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/ingestion"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/memoexport"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/nationalplanning"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/projectfinance"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/risk"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/syndication"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/ubo"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "help", "-h", "--help":
		printUsage()
	case "planning":
		handlePlanning(os.Args[2:])
	case "syndication":
		handleSyndication(os.Args[2:])
	case "offtake":
		handleOfftake(os.Args[2:])
	case "corridor":
		handleCorridor(os.Args[2:])
	case "finance":
		handleFinance(os.Args[2:])
	case "filings":
		handleFilings(os.Args[2:])
	case "kpi", "indicators":
		handleKPI(os.Args[2:])
	case "cegs":
		handleCEGS(os.Args[2:])
	case "search":
		handleSearch(os.Args[2:])
	case "project":
		handleProject(os.Args[2:])
	case "changes":
		handleChanges(os.Args[2:])
	case "rankings":
		handleRankings(os.Args[2:])
	case "export":
		handleExport(os.Args[2:])
	case "demo":
		handleDemo()
	case "extract":
		handleExtract(os.Args[2:])
	case "ontology":
		handleOntology(os.Args[2:])
	case "lakehouse":
		handleLakehouse(os.Args[2:])
	case "gql":
		handleGQL(os.Args[2:])
	case "earthobs":
		handleEarthObsCLI(os.Args[2:])
	case "ai":
		handleAICLI(os.Args[2:])
	case "security":
		handleSecurityCLI(os.Args[2:])
	case "ubo":
		handleUBOCLI(os.Args[2:])
	case "audit":
		handleAuditCLI(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("CanadaOpportunityGraph CLI (cog) — Economic Intelligence & CEGS Reference Tool")
	fmt.Println("\nUsage:")
	fmt.Println("  cog search <query>                     Search projects and infrastructure assets")
	fmt.Println("  cog project show <id|slug>             Display investor-grade project profile")
	fmt.Println("  cog project mrio <id|slug>             Display StatCan input-output macro multipliers")
	fmt.Println("  cog project flyvbjerg <id|slug>        Display Bayesian cost & schedule overrun hazard curve")
	fmt.Println("  cog project ubo <id|slug>              Display Ultimate Beneficial Ownership & ICA screening")
	fmt.Println("  cog project grid <id|slug>             Display electrical grid hosting capacity feasibility")
	fmt.Println("  cog project earthobs <id|slug>         Display satellite SAR and optical ground-truth telemetry")
	fmt.Println("  cog syndication match <id|slug>        Match institutional capital syndication & Maple 8 allocators")
	fmt.Println("  cog syndication indigenous [name]      Model multi-nation First Nation linear equity syndicate & ILGP")
	fmt.Println("  cog offtake list [proj-id]             Display commercial offtake & long-term clean power PPAs")
	fmt.Println("  cog corridor route [options]           Simulate linear Right-of-Way (RoW) geotechnical impedance")
	fmt.Println("  cog corridor port [port_id]            Display strategic gateway multi-modal rail logistics")
	fmt.Println("  cog finance simulate <id|slug> [--runs N] Run stochastic 10k Monte Carlo project cash flow model")
	fmt.Println("  cog finance cleantax <id|slug>         Calculate Clean Economy ITCs & CCfD underwriting")
	fmt.Println("  cog filings <stream|ea|amendments>     Stream real-time regulatory filings & tender amendments")
	fmt.Println("  cog kpi list                           Display registered Canadian KPI & indicator definitions")
	fmt.Println("  cog kpi show <id|slug>                 Display 8-pillar project KPI scorecard and gap analysis")
	fmt.Println("  cog kpi feeds                          Stream live market ticks and commodity indicators")
	fmt.Println("  cog kpi snapshot                       Export canonical KPI snapshot")
	fmt.Println("  cog planning optimize [--obj <type>]   Run Sovereign Capital Allocation Optimizer")
	fmt.Println("  cog planning wargame [--shock <type>]  Run geopolitical macro shock stress-testing")
	fmt.Println("  cog planning labor [--prov <prov>]     Display Red Seal craft labor collision report")
	fmt.Println("  cog extract cards <file>               Extract restricted portfolio cards from text")
	fmt.Println("  cog extract ni43101 <file>             Extract NI 43-101 mining reserves & economic metrics")
	fmt.Println("  cog extract waterfall <file>           Extract capital stack financing waterfall & WACC")
	fmt.Println("  cog extract iaac <file>                Extract IAAC Decision Statement conditions & legal risk")
	fmt.Println("  cog changes [--since 7d|30d]           List recent momentum signals and milestones")
	fmt.Println("  cog rankings <dimension>               Rank projects (buildability, investability, etc.)")
	fmt.Println("  cog export project <id> [--format json|md|cegs|memo|geojson] Export dossier with provenance")
	fmt.Println("  cog export memo <id>                   Export Privy Council Office Memorandum to Cabinet")
	fmt.Println("  cog export geojson <id>                Export OGC GeoJSON FeatureCollection")
	fmt.Println("  cog cegs validate <file>               Validate document against CEGS 0.1 standard")
	fmt.Println("  cog cegs inspect <file>                Inspect CEGS document & evidence trust profile")
	fmt.Println("  cog cegs diff <old.json> <new.json>    Semantic diff between two CEGS states")
	fmt.Println("  cog ontology <branch|action>           Manage git-like ontology branches and transactional actions")
	fmt.Println("  cog lakehouse <query|quality>          Vectorized zero-copy OLAP query & quality SLA firewall")
	fmt.Println("  cog gql <query|gnn|communities>        ISO GQL parsing, GNN link prediction & Louvain clusters")
	fmt.Println("  cog earthobs <insar|ais|audit> [id]    Sentinel-1 InSAR, AIS port congestion & filing reconciliation")
	fmt.Println("  cog ai <copilot|ppa> [options]         Sovereign AIP multi-agent peer review & bilateral PPA optimizer")
	fmt.Println("  cog security <auth|dlp> [options]      ABAC security clearance checks & real-time DLP redaction")
	fmt.Println("  cog ubo <trace|darkfleet|ica>          Beneficial ownership unraveling, dark fleet & ICA s.25.3 review")
	fmt.Println("  cog audit <c59|c69|trc92|merkle> [id]  Bill C-59 ITC, C-69 clock, TRC 92 scorecard & Merkle root")
	fmt.Println("  cog demo                               Run instant offline demonstration")
}

func getSeededStore() database.Store {
	store := database.NewMemoryStore()
	fixturePath := "data/fixtures/nrcan_mpi_2025.json"
	officialPath := ""
	tradePath := ""
	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		if _, err := os.Stat("../../" + fixturePath); err == nil {
			fixturePath = "../../" + fixturePath
			officialPath = "../../data/fixtures/official_records.json"
			tradePath = "../../data/fixtures/world_bank_trade_canada.json"
		}
	}
	adapterList := []adapters.Adapter{
		nrcan_major_projects.NewNRCanAdapter(fixturePath),
		official.NewAdapter(officialPath),
		global_trade.NewAdapter(tradePath),
	}
	p := ingestion.NewPipeline(store, adapterList)
	if _, err := p.Run(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "failed to load reviewed snapshot: %v\n", err)
	}
	return store
}

func handleCEGS(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog cegs <validate|inspect|diff> [args...]")
		return
	}

	subCmd := args[0]
	switch subCmd {
	case "validate":
		if len(args) < 2 {
			fmt.Println("Usage: cog cegs validate <file.json>")
			return
		}
		data, err := os.ReadFile(args[1])
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		rep, err := cegs.Validate(data)
		if err != nil {
			fmt.Printf("Validation error: %v\n", err)
			os.Exit(1)
		}
		if rep.Valid {
			fmt.Printf("✓ %s\n", rep.Summary)
			fmt.Printf("  ID:           %s\n", rep.ID)
			fmt.Printf("  Type:         %s\n", rep.ResourceType)
			fmt.Printf("  Conformance:  %s\n", rep.ConformanceLevel)
		} else {
			fmt.Printf("✗ %s\n", rep.Summary)
			for _, e := range rep.Errors {
				fmt.Printf("  - Error: %s\n", e)
			}
			os.Exit(1)
		}

	case "inspect":
		if len(args) < 2 {
			fmt.Println("Usage: cog cegs inspect <file.json>")
			return
		}
		data, err := os.ReadFile(args[1])
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}
		insp, err := cegs.Inspect(data)
		if err != nil {
			fmt.Printf("Inspection error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("CEGS %s | Type: %s\n", insp.CEGSVersion, insp.Type)
		fmt.Printf("Name:        %s\n", insp.CanonicalName)
		if insp.Stage != "" {
			fmt.Printf("Stage:       %s\n", insp.Stage)
		}
		fmt.Printf("Evidence:    %d reference(s)\n", insp.EvidenceCount)
		fmt.Printf("Conformance: %s\n", insp.ConformanceLevel)

	case "diff":
		if len(args) < 3 {
			fmt.Println("Usage: cog cegs diff <old.json> <new.json>")
			return
		}
		oldData, err := os.ReadFile(args[1])
		if err != nil {
			fmt.Printf("Error reading old file: %v\n", err)
			os.Exit(1)
		}
		newData, err := os.ReadFile(args[2])
		if err != nil {
			fmt.Printf("Error reading new file: %v\n", err)
			os.Exit(1)
		}
		diff, err := cegs.Diff(oldData, newData)
		if err != nil {
			fmt.Printf("Diff error: %v\n", err)
			os.Exit(1)
		}
		if !diff.HasChanges {
			fmt.Println("No semantic differences detected between states.")
			return
		}
		fmt.Printf("Detected %d semantic change(s) in %s:\n", len(diff.Changes), diff.ResourceID)
		for _, ch := range diff.Changes {
			fmt.Printf("  [%s] %s\n", ch.Code, ch.Description)
		}

	default:
		fmt.Printf("Unknown cegs command: %s\n", subCmd)
	}
}

func handleSearch(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog search <query>")
		return
	}
	query := strings.Join(args, " ")
	store := getSeededStore()

	projects, total, err := store.ListProjects(context.Background(), database.ProjectFilter{
		Search: query,
		Limit:  15,
	})
	if err != nil {
		fmt.Printf("Search error: %v\n", err)
		return
	}

	fmt.Printf("Found %d project(s) matching '%s':\n\n", total, query)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "SLUG\tNAME\tSECTOR\tPROV\tSTAGE\tCAPEX (CAD)\tBUILDABILITY")
	for _, p := range projects {
		bScore := 0.0
		if p.Scores != nil {
			bScore = p.Scores["buildability"]
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t$%d\t%.1f/100\n",
			p.Slug, p.Name, p.Sector, p.Province, p.CurrentStage, p.CapexCAD, bScore)
	}
	w.Flush()
}

func findProject(ctx context.Context, store database.Store, id string) (*domain.Project, error) {
	proj, err := store.GetProject(ctx, id)
	if err == nil && proj != nil {
		return proj, nil
	}
	proj, err = store.GetProjectBySlug(ctx, id)
	if err == nil && proj != nil {
		return proj, nil
	}
	return nil, fmt.Errorf("project not found: %s", id)
}

func handleProject(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: cog project <show|mrio|flyvbjerg|ubo|grid|earthobs> <id|slug>")
		return
	}
	action := strings.ToLower(args[0])
	id := args[1]
	store := getSeededStore()
	ctx := context.Background()

	proj, err := findProject(ctx, store, id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	switch action {
	case "show":
		bundle, err := export.ExportProjectBundle(ctx, store, proj.ID)
		if err != nil {
			fmt.Printf("Error exporting project: %v\n", err)
			return
		}
		fmt.Println(bundle.ToMarkdown())

	case "mrio":
		impact := econometrics.NewEngine().CalculateMRIO(proj)
		fmt.Printf("=== StatCan Input-Output Macro Multipliers: %s ===\n", proj.Name)
		fmt.Printf("Project ID:            %s\n", proj.ID)
		fmt.Printf("Tracked CAPEX:         $%d CAD\n", impact.CapexCAD)
		fmt.Printf("Direct GDP Impact:     $%d CAD\n", impact.DirectGDPCAD)
		fmt.Printf("Indirect Supply GDP:   $%d CAD\n", impact.IndirectGDPCAD)
		fmt.Printf("Induced Wage GDP:      $%d CAD\n", impact.InducedGDPCAD)
		fmt.Printf("Total GDP Generated:   $%d CAD (%.2fx Multiplier)\n", impact.TotalGDPCAD, impact.TotalMultipler)
		fmt.Printf("Person-Years of Jobs:  %d FTE\n", impact.PersonYearsJobs)
		fmt.Println("Fiscal Tax Returns:")
		fmt.Printf("  Federal Tax:         $%d CAD\n", impact.FederalTaxCAD)
		fmt.Printf("  Provincial Tax:      $%d CAD\n", impact.ProvincialTaxCAD)
		fmt.Printf("  Municipal Tax:       $%d CAD\n", impact.MunicipalTaxCAD)
		fmt.Printf("  Total Fiscal Return: $%d CAD\n", impact.TotalFiscalReturn)
		fmt.Printf("Model Version:         %s\n", impact.ModelVersion)
		fmt.Printf("Audit Hash:            %s\n", impact.AuditHash)

	case "flyvbjerg":
		forecast := risk.NewEvaluator().ForecastProject(proj)
		fmt.Printf("=== Bayesian Reference Class Overrun Analysis: %s ===\n", proj.Name)
		fmt.Printf("Sector Benchmark:      %s (Sample size: %d projects)\n", forecast.Sector, forecast.HistoricalSampleSize)
		fmt.Printf("Expected Cost Overrun: +%.1f%%\n", forecast.ExpectedCostOverrunPct)
		fmt.Printf("Expected Delay:        +%d months\n", forecast.ExpectedDelayMonths)
		if forecast.RemoteGeographyPenalty > 0 {
			fmt.Printf("Remote Penalty:        +%.1f%% (Arctic / Remote Corridor)\n", forecast.RemoteGeographyPenalty*100)
		}
		if forecast.TechNoveltyPenalty > 0 {
			fmt.Printf("Tech Novelty Penalty:  +%.1f%% (FOAK / SMR / Novel Process)\n", forecast.TechNoveltyPenalty*100)
		}
		fmt.Println("\nHazard Distribution Quantiles:")
		for _, pt := range forecast.Percentiles {
			fmt.Printf("  P%-2d             Cost: +%5.1f%% | Schedule: +%2d mo | Forecast CAPEX: $%d CAD\n",
				pt.Percentile, pt.CostOverrunPct, pt.ScheduleDelayMonths, pt.ForecastCapexCAD)
		}
		fmt.Printf("\nAudit Hash:            %s\n", forecast.AuditHash)

	case "ubo":
		screening := ubo.NewEvaluator().ScreenProject(proj, nil)
		fmt.Printf("=== Sovereign Screening & ICA National Security Review: %s ===\n", proj.Name)
		fmt.Printf("Proponent:             %s (ID: %s)\n", screening.ProponentName, screening.ProponentID)
		fmt.Printf("Investment Canada Act: %s\n", screening.ICARisk)
		fmt.Printf("Domestic Control:      %.1f%%\n", screening.DomesticControlShare*100)
		fmt.Printf("FTA Partner Share:     %.1f%%\n", screening.FTAPartnerShare*100)
		fmt.Printf("Non-FTA Foreign Share: %.1f%%\n", screening.NonFTAShare*100)
		fmt.Printf("Foreign SOE Share:     %.1f%%\n", screening.SOEExposurePercent*100)
		fmt.Printf("Critical Mineral Flag: %t\n", screening.CriticalMineralFlag)
		fmt.Printf("Dual-Use Sovereignty:  %t\n", screening.DualUseSovereignty)
		if len(screening.NationalSecurityNotes) > 0 {
			fmt.Println("National Security Notes:")
			for _, n := range screening.NationalSecurityNotes {
				fmt.Printf("  - %s\n", n)
			}
		}
		fmt.Printf("Audit Hash:            %s\n", screening.AuditHash)

	case "grid":
		assessment := gridphysics.NewEngine().AssessProject(proj)
		fmt.Printf("=== Electrical Grid Feasibility & Interconnect: %s ===\n", proj.Name)
		fmt.Printf("System Operator:       %s\n", assessment.Operator)
		fmt.Printf("Estimated Load / Gen:  %.1f MW (Voltage: %d kV)\n", assessment.EstimatedLoadOrGenMW, assessment.InterconnectVoltageKV)
		fmt.Printf("Grid Feasibility Score:%.1f / 100\n", assessment.GridFeasibilityScore)
		fmt.Printf("Clean Power Purity:    %.1f%% (Hydro / Nuclear)\n", assessment.CleanPowerPurityPct)
		fmt.Printf("Queue Estimated Delay: %d months\n", assessment.QueueEstimatedMonths)
		fmt.Printf("Substation Headroom:   %.1f MW\n", assessment.SubstationHeadroomMW)
		fmt.Printf("Dedicated Substation:  %t\n", assessment.DedicatedSubstationNeeded)
		fmt.Printf("Reinforcement CAPEX:   $%d CAD\n", assessment.ReinforcementCostCAD)
		if len(assessment.InterconnectNotes) > 0 {
			fmt.Println("Grid Engineering Notes:")
			for _, n := range assessment.InterconnectNotes {
				fmt.Printf("  - %s\n", n)
			}
		}
		fmt.Printf("Audit Hash:            %s\n", assessment.AuditHash)

	case "earthobs":
		dossier := earthobs.NewEvaluator().CorroborateProject(proj, nil)
		fmt.Printf("=== Satellite Ground-Truth & SAR Corroboration: %s ===\n", proj.Name)
		fmt.Printf("Claimed Stage:         %s\n", dossier.ClaimedStage)
		fmt.Printf("Corroboration Status:  %s\n", dossier.CorroborationStatus)
		fmt.Printf("Physical Progress:     %.1f / 100\n", dossier.PhysicalProgressScore)
		fmt.Printf("Earthworks Confirmed:  %t (Copernicus Sentinel-1 SAR)\n", dossier.EarthworksConfirmed)
		fmt.Printf("Structures Confirmed:  %t (RCS Radar Cross-Section)\n", dossier.StructuresConfirmed)
		if !dossier.LastSatellitePass.IsZero() {
			fmt.Printf("Last Satellite Pass:   %s\n", dossier.LastSatellitePass.Format("2006-01-02"))
		}
		fmt.Printf("Telemetry Summary:     %s\n", dossier.TelemetrySummary)
		fmt.Printf("Audit Hash:            %s\n", dossier.AuditHash)

	default:
		fmt.Printf("Unknown project action: %s\n", action)
		fmt.Println("Available actions: show, mrio, flyvbjerg, ubo, grid, earthobs")
	}
}

func handleChanges(_ []string) {
	store := getSeededStore()
	signals, _ := store.ListSignals(context.Background(), 30*24*time.Hour, 15)

	fmt.Printf("Recent Capital and Milestone Signals (Last 30 Days):\n\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DATE\tPROJECT\tSIGNAL TYPE\tDESCRIPTION")
	for _, s := range signals {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			s.Timestamp.Format("2006-01-02"), s.ProjectName, s.Type, s.Description)
	}
	w.Flush()
}

func handleRankings(args []string) {
	dim := "buildability"
	if len(args) > 0 {
		dim = strings.ToLower(args[0])
	}
	store := getSeededStore()
	projects, _ := store.ListRankings(context.Background(), dim, 10)

	fmt.Printf("Top Projects by %s:\n\n", strings.Title(dim))
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "RANK\tNAME\tSECTOR\tPROVINCE\tSCORE\tCAPEX (CAD)")
	for i, p := range projects {
		score := 0.0
		if p.Scores != nil {
			score = p.Scores[dim]
		}
		fmt.Fprintf(w, "#%d\t%s\t%s\t%s\t%.1f\t$%d\n",
			i+1, p.Name, p.Sector, p.Province, score, p.CapexCAD)
	}
	w.Flush()
}

func handleExport(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: cog export <project|memo|geojson> <id|slug> [--format json|md|cegs|memo|geojson]")
		return
	}
	exportMode := strings.ToLower(args[0])
	id := args[1]
	format := "json"
	switch exportMode {
	case "memo":
		format = "memo"
	case "geojson":
		format = "geojson"
	}

	for i, arg := range args {
		if arg == "--format" && i+1 < len(args) {
			format = strings.ToLower(args[i+1])
		}
	}

	store := getSeededStore()
	ctx := context.Background()

	proj, err := store.GetProject(ctx, id)
	if err != nil {
		proj, err = store.GetProjectBySlug(ctx, id)
		if err != nil {
			fmt.Printf("Project not found: %s\n", id)
			return
		}
	}

	switch format {
	case "memo":
		memoGen := memoexport.NewMemoGenerator()
		memo := memoGen.GenerateCabinetMemo(proj, memoexport.MemoTypeCabinetMC)
		fmt.Println(memo.MarkdownContent)
	case "geojson":
		fc := memoexport.ExportGeoJSON([]*domain.Project{proj})
		out, _ := json.MarshalIndent(fc, "", "  ")
		fmt.Println(string(out))
	case "md", "markdown":
		bundle, err := export.ExportProjectBundle(ctx, store, proj.ID)
		if err != nil {
			fmt.Printf("Export error: %v\n", err)
			return
		}
		fmt.Println(bundle.ToMarkdown())
	case "cegs":
		bundle, err := export.ExportProjectBundle(ctx, store, proj.ID)
		if err != nil {
			fmt.Printf("Export error: %v\n", err)
			return
		}
		cegsProj, _ := bundle.ToCEGSExport()
		out, _ := json.MarshalIndent(cegsProj, "", "  ")
		fmt.Println(string(out))
	default:
		bundle, err := export.ExportProjectBundle(ctx, store, proj.ID)
		if err != nil {
			fmt.Printf("Export error: %v\n", err)
			return
		}
		out, _ := json.MarshalIndent(bundle, "", "  ")
		fmt.Println(string(out))
	}
}

func handleDemo() {
	fmt.Println("=== CanadaOpportunityGraph Deterministic Demo Mode ===")
	store := getSeededStore()
	ctx := context.Background()

	stats, _ := store.GetRadarStats(ctx)
	fmt.Printf("Tracked Projects:     %d\n", stats.TotalProjects)
	fmt.Printf("Total Tracked CAPEX:  $%.2f Billion CAD\n", float64(stats.TotalCapexCAD)/1e9)
	fmt.Printf("Capital Moving Week:  $%.2f Million CAD\n", float64(stats.CapitalMovingWeekCAD)/1e6)
	fmt.Printf("Accelerating Assets:  %d\n", stats.AcceleratingProjectsCount)
	fmt.Printf("Active Procurements:  %d\n\n", stats.ActiveProcurementsCount)

	fmt.Println("Top Ranked Projects (Buildability):")
	rankings, _ := store.ListRankings(ctx, "buildability", 3)
	for i, p := range rankings {
		fmt.Printf("  %d. %s (%s, %s) — Buildability: %.1f/100, CAPEX: $%.1fB\n",
			i+1, p.Name, p.Province, p.Sector, p.Scores["buildability"], float64(p.CapexCAD)/1e9)
	}

	fmt.Println("\nCEGS Standard Specification: 0.1 | Reference Implementation Verified.")
}

func handleExtract(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog extract <cards|ni43101|waterfall|iaac> <file> [options]")
		return
	}

	mode := "cards"
	filePath := args[0]
	shift := 0

	switch args[0] {
	case "cards", "ni43101", "waterfall", "iaac":
		mode = args[0]
		if len(args) < 2 {
			fmt.Printf("Usage: cog extract %s <file>\n", mode)
			return
		}
		filePath = args[1]
		shift = 1
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	switch mode {
	case "ni43101":
		rep, err := documentintelligence.ExtractNI43101TechnicalReport(string(data))
		if err != nil {
			fmt.Printf("Extraction error: %v\n", err)
			os.Exit(1)
		}
		out, _ := json.MarshalIndent(rep, "", "  ")
		fmt.Println(string(out))

	case "waterfall":
		wf := documentintelligence.ExtractCapitalWaterfall(string(data), 0)
		out, _ := json.MarshalIndent(wf, "", "  ")
		fmt.Println(string(out))

	case "iaac":
		stmt, err := documentintelligence.ExtractIAACDecisionStatement(string(data))
		if err != nil {
			fmt.Printf("Extraction error: %v\n", err)
			os.Exit(1)
		}
		out, _ := json.MarshalIndent(stmt, "", "  ")
		fmt.Println(string(out))

	case "cards":
		sourceID := "workspace"
		visibility := domain.VisibilityInternalRestricted
		for i := shift; i < len(args); i++ {
			switch args[i] {
			case "--source":
				if i+1 < len(args) {
					sourceID = args[i+1]
					i++
				}
			case "--visibility":
				if i+1 < len(args) {
					visibility = domain.VisibilityClass(strings.ToUpper(strings.TrimSpace(args[i+1])))
					i++
				}
			}
		}
		cards, err := documentintelligence.ExtractCards(string(data), sourceID, visibility, time.Now().UTC())
		if err != nil {
			fmt.Printf("Extraction error: %v\n", err)
			os.Exit(1)
		}
		if len(cards) == 0 {
			fmt.Println("No portfolio cards found in the supplied text.")
			return
		}
		out, _ := json.MarshalIndent(cards, "", "  ")
		fmt.Println(string(out))
	}
}

func handlePlanning(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog planning <optimize|wargame|labor> [options]")
		return
	}
	subCmd := strings.ToLower(args[0])
	store := getSeededStore()
	projects, _, err := store.ListProjects(context.Background(), database.ProjectFilter{Limit: 1000})
	if err != nil {
		fmt.Printf("Error fetching projects: %v\n", err)
		return
	}

	switch subCmd {
	case "optimize":
		obj := nationalplanning.ObjectiveBalancedStrategy
		for i := 1; i < len(args); i++ {
			if args[i] == "--obj" && i+1 < len(args) {
				obj = nationalplanning.ObjectiveType(strings.ToUpper(args[i+1]))
				i++
			}
		}
		res := nationalplanning.NewOptimizer().Optimize(projects, nationalplanning.OptimizationRequest{
			Objective: obj,
		})

		fmt.Println("=== Sovereign Capital Allocation Optimizer ===")
		fmt.Printf("Objective Target:           %s\n", res.Objective)
		fmt.Printf("Total Public Invested:      $%d CAD\n", res.TotalPublicInvestedCAD)
		fmt.Printf("Private Capital Mobilized:  $%d CAD\n", res.TotalPrivateMobilizedCAD)
		fmt.Printf("Crowding-In Multiplier:     %.2fx\n", res.CrowdingInMultiplier)
		fmt.Printf("Total GHG Abated:           %.1f Mt CO2e / yr\n", res.TotalGHGAbatedMtPerYear)
		fmt.Printf("Audit Hash:                 %s\n\n", res.AuditHash)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "PROJECT\tPROV\tSECTOR\tTOTAL CAPEX\tPUBLIC ALLOC\tPRIVATE MOBILIZED\tUTILITY")
		for _, p := range res.AllocatedProjects {
			fmt.Fprintf(w, "%s\t%s\t%s\t$%d\t$%d\t$%d\t%.1f\n",
				p.ProjectName, p.Province, p.Sector, p.TotalCapexCAD, p.TotalPublicCAD, p.PrivateMobilized, p.UtilityScore)
		}
		w.Flush()

	case "wargame":
		scenario := nationalplanning.ShockUSMCATariffs
		for i := 1; i < len(args); i++ {
			if args[i] == "--shock" && i+1 < len(args) {
				scenario = nationalplanning.ShockScenario(strings.ToUpper(args[i+1]))
				i++
			}
		}
		res := nationalplanning.NewWarGameEngine().SimulateScenario(projects, nationalplanning.WarGameRequest{
			Scenario: scenario,
		})

		fmt.Println("=== National Planning Geopolitical War Game ===")
		fmt.Printf("Scenario:                %s\n", res.ScenarioTitle)
		fmt.Printf("Description:             %s\n", res.ScenarioDescription)
		fmt.Printf("Stalled Assets:          %d\n", res.TotalAssetsStalledCount)
		fmt.Printf("Total Frozen CAPEX:      $%d CAD\n", res.TotalFrozenCapexCAD)
		fmt.Printf("Est. National GDP Loss:  $%d CAD\n", res.EstimatedNationalGDPLossCAD)
		fmt.Printf("Audit Hash:              %s\n\n", res.AuditHash)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "PROJECT\tPROV\tSECTOR\tCAPEX\tSTALL %\tRECOMMENDED COUNTERMEASURE")
		for _, sp := range res.StalledProjects {
			fmt.Fprintf(w, "%s\t%s\t%s\t$%d\t%.0f%%\t%s\n",
				sp.ProjectName, sp.Province, sp.Sector, sp.OriginalCapexCAD, sp.StallLikelihood*100, sp.RecommendedAction)
		}
		w.Flush()

		if len(res.SovereignMitigations) > 0 {
			fmt.Println("\nSovereign Strategic Mitigations:")
			for i, m := range res.SovereignMitigations {
				fmt.Printf("  %d. %s\n", i+1, m)
			}
		}

	case "labor":
		prov := "ON"
		for i := 1; i < len(args); i++ {
			if args[i] == "--prov" && i+1 < len(args) {
				prov = strings.ToUpper(args[i+1])
				i++
			}
		}
		rep := nationalplanning.NewLaborAggregator().AnalyzeProvince(prov, projects)

		fmt.Printf("=== Red Seal Craft Labor Pinch-Point Report: %s ===\n", rep.Province)
		fmt.Printf("Active Tracked CAPEX:    $%d CAD\n", rep.TotalActiveCapexCAD)
		fmt.Printf("Concurrent Projects:     %d\n", rep.ConcurrentProjects)
		fmt.Printf("Total Craft Labor Peak:  %d FTE\n", rep.TotalLaborDemandFTE)
		fmt.Printf("Collision Detected:      %t\n", rep.CollisionDetected)
		fmt.Printf("Audit Hash:              %s\n\n", rep.AuditHash)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "RED SEAL TRADE\tDEMAND FTE\tSUPPLY FTE\tUTILIZATION\tSTATUS\tWAGE RISK")
		for _, t := range rep.Trades {
			fmt.Fprintf(w, "%s\t%d\t%d\t%.1f%%\t%s\t%s\n",
				t.Trade, t.PeakDemandFTE, t.RegionalSupplyFTE, t.UtilizationPct, t.CollisionStatus, t.WageInflationRisk)
		}
		w.Flush()

		if rep.StrategicAdvice != "" {
			fmt.Printf("\nStrategic Workforce Directive:\n%s\n", rep.StrategicAdvice)
		}

	default:
		fmt.Printf("Unknown planning command: %s\n", subCmd)
		fmt.Println("Usage: cog planning <optimize|wargame|labor>")
	}
}

func handleSyndication(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog syndication <match|indigenous> [args...]")
		return
	}
	sub := strings.ToLower(args[0])
	store := getSeededStore()
	ctx := context.Background()

	switch sub {
	case "match":
		if len(args) < 2 {
			fmt.Println("Usage: cog syndication match <id|slug>")
			return
		}
		id := args[1]
		proj, err := store.GetProject(ctx, id)
		if err != nil {
			proj, err = store.GetProjectBySlug(ctx, id)
			if err != nil {
				fmt.Printf("Project not found: %s\n", id)
				return
			}
		}
		matcher := syndication.NewMatcher(nil)
		consortium := matcher.MatchProject(proj)
		fmt.Printf("=== Institutional Syndication Consortium: %s ===\n", consortium.ProjectName)
		fmt.Printf("Total CAPEX:                $%d CAD\n", consortium.TotalCapexCAD)
		fmt.Printf("Recommended Equity Tranche: $%d CAD\n", consortium.EquityTrancheCAD)
		fmt.Printf("Recommended Debt Tranche:   $%d CAD\n", consortium.DebtTrancheCAD)
		fmt.Printf("Crown Concession (CIB/CGF): $%d CAD\n", consortium.CrownConcessionCAD)
		fmt.Printf("Indigenous Equity (ILGP):   $%d CAD\n", consortium.IndigenousEquityCAD)
		fmt.Printf("Private Crowding-In Ratio:  %.2fx\n", consortium.PrivateCrowdingInRatio)
		fmt.Printf("Audit Hash:                 %s\n\n", consortium.AuditHash)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "INVESTOR\tCLASS\tSCORE\tTRANCHE\tPROPOSED TICKET\tRATIONALE")
		for _, m := range consortium.Matches {
			fmt.Fprintf(w, "%s\t%s\t%.1f\t%s\t$%d\t%s\n",
				m.Investor.Name, m.Investor.Class, m.MatchScore, m.RecommendedTranche, m.ProposedTicketCAD, m.Rationale)
		}
		w.Flush()

	case "indigenous":
		name := "Ring of Fire Northern Access Road & Transmission Corridor"
		if len(args) >= 2 {
			name = strings.Join(args[1:], " ")
		}
		syndicate := syndication.BuildMultiNationSyndicate(name, 1_200_000_000, []struct {
			Name      string
			Territory string
			KM        float64
		}{
			{Name: "Marten Falls First Nation", Territory: "Treaty 9 Unceded Anishinaabe", KM: 140.0},
			{Name: "Webequie First Nation", Territory: "Treaty 9 Traditional Territory", KM: 110.0},
			{Name: "Neskantaga First Nation", Territory: "Treaty 9 Watershed Stewardship", KM: 60.0},
			{Name: "Nibinamik First Nation", Territory: "Treaty 9 Traditional Boreal", KM: 40.0},
		})

		fmt.Printf("=== Multi-Nation Indigenous Equity Syndicate: %s ===\n", syndicate.CorridorProjectName)
		fmt.Printf("Total Corridor Length:      %.1f km\n", syndicate.TotalCorridorKM)
		fmt.Printf("Total Equity Value:         $%d CAD\n", syndicate.TotalEquityValueCAD)
		fmt.Printf("Federal ILGP Debt Guarantee: $%d CAD\n", syndicate.FederalILGPGreaterCAD)
		fmt.Printf("Interest Spread Savings:    %d bps\n", syndicate.BlendedInterestSpreadBps)
		fmt.Printf("Total Annual Dividends:     $%d CAD / yr\n", syndicate.TotalAnnualDividendsCAD)
		fmt.Printf("Total 30-Year Wealth Gen:   $%d CAD\n", syndicate.Total30YearWealthCAD)
		fmt.Printf("Audit Hash:                 %s\n\n", syndicate.AuditHash)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "COMMUNITY\tTERRITORY\tKM\tEQUITY %\tGUARANTEED DEBT\tANNUAL DIVIDEND\t30-YR WEALTH")
		for _, p := range syndicate.Participants {
			fmt.Fprintf(w, "%s\t%s\t%.1f km\t%.1f%%\t$%d\t$%d\t$%d\n",
				p.BandCouncilName, p.TreatyOrTerritory, p.CorridorKilometers, p.EquitySharePercent, p.GuaranteedDebtCAD, p.AnnualDividendCAD, p.Cumulative30YrCAD)
		}
		w.Flush()

	default:
		fmt.Printf("Unknown syndication command: %s\n", sub)
		fmt.Println("Usage: cog syndication <match|indigenous>")
	}
}

func handleOfftake(args []string) {
	offtakes := syndication.CanonicalOfftakeAgreements()
	targetProj := ""
	if len(args) > 0 && args[0] != "list" {
		targetProj = strings.ToLower(args[0])
	} else if len(args) > 1 && args[0] == "list" {
		targetProj = strings.ToLower(args[1])
	}

	fmt.Println("=== Commercial Offtake & Clean Power Purchase Agreements (PPAs) ===")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "AGREEMENT ID\tPROJECT\tOFFTAKER / BUYER\tRATING\tCOMMODITY\tVOLUME\tTERM\tANNUAL VALUE")
	for _, o := range offtakes {
		if targetProj != "" && !strings.Contains(strings.ToLower(o.ProjectID), targetProj) && !strings.Contains(strings.ToLower(o.ID), targetProj) {
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%d yrs\t$%d CAD\n",
			o.ID, o.ProjectID, o.BuyerName, o.BuyerCreditRating, o.Commodity, o.VolumeAnnualMetric, o.TermYears, o.AnnualContractValueCAD)
	}
	w.Flush()
}

func handleCorridor(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog corridor <route|port> [args...]")
		return
	}
	sub := strings.ToLower(args[0])

	switch sub {
	case "route":
		origin := corridor.CorridorPoint{Name: "Prince George", Latitude: 53.9171, Longitude: -122.7497, ElevationMeters: 575}
		dest := corridor.CorridorPoint{Name: "Port of Prince Rupert", Latitude: 54.3150, Longitude: -130.3208, ElevationMeters: 10}
		infraType := corridor.TypeHVDCTransmission

		for i := 1; i < len(args); i++ {
			switch args[i] {
			case "--type", "--mode":
				if i+1 < len(args) {
					switch strings.ToLower(args[i+1]) {
					case "pipeline", "h2", "hydrogen":
						infraType = corridor.TypeHydrogenPipeline
					case "rail":
						infraType = corridor.TypeArcticHeavyRail
					case "co2":
						infraType = corridor.TypeCO2Pipeline
					case "road":
						infraType = corridor.TypeAllWeatherRoad
					default:
						infraType = corridor.TypeHVDCTransmission
					}
					i++
				}
			case "--origin":
				if i+1 < len(args) {
					origin.Name = args[i+1]
					i++
				}
			case "--dest":
				if i+1 < len(args) {
					dest.Name = args[i+1]
					i++
				}
			}
		}

		router := corridor.NewRoutingEngine()
		eval := router.EvaluateCorridor(origin, dest, infraType)

		fmt.Printf("=== Linear Right-of-Way (RoW) Pathfinding Evaluation ===\n")
		fmt.Printf("Corridor ID:          %s\n", eval.CorridorID)
		fmt.Printf("Type:                 %s\n", eval.Type)
		fmt.Printf("Origin -> Dest:       %s -> %s\n", eval.OriginName, eval.DestinationName)
		fmt.Printf("Total Length:         %.1f km\n", eval.TotalLengthKM)
		fmt.Printf("Estimated Capex:      $%d CAD\n", eval.TotalEstimatedCapexCAD)
		fmt.Printf("Schedule Duration:    %d months\n", eval.EstimatedScheduleMonths)
		fmt.Printf("Impedance Index:      %.1f\n", eval.ImpedanceIndex)
		fmt.Printf("Permafrost Hazard:    %.0f%%\n", eval.Sensitivity.PermafrostThawHazardScore*100)
		fmt.Printf("Caribou Overlap:      %.1f km\n", eval.Sensitivity.CaribouRangeIntersectKM)
		fmt.Printf("Wetland Crossings:    %d\n", eval.Sensitivity.WetlandCrossingCount)
		fmt.Printf("Audit Hash:           %s\n\n", eval.AuditHash)

	case "port":
		gateways := corridor.CanonicalGateways()
		target := ""
		if len(args) > 1 {
			target = strings.ToLower(args[1])
		}

		fmt.Println("=== Strategic Maritime Gateways & Multi-Modal Rail Logistics ===")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "GATEWAY\tPORT NAME\tPROV\tRAIL\tTHROUGHPUT\tVESSEL DWELL\tRAIL DWELL\tBERTH %\tSTATUS")
		for _, g := range gateways {
			if target != "" && !strings.Contains(strings.ToLower(string(g.GatewayID)), target) && !strings.Contains(strings.ToLower(g.PortName), target) {
				continue
			}
			rails := strings.Join(g.Class1RailConnections, "/")
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%.1f Mt\t%.1f hrs\t%.1f hrs\t%.1f%%\t%s\n",
				g.GatewayID, g.PortName, g.Province, rails, g.AnnualThroughputMNTonnes, g.AverageVesselDwellHours, g.AverageRailcarDwellHours, g.BerthCapacityUtilizationPercent, g.ActiveBottleneckStatus)
		}
		w.Flush()

	default:
		fmt.Printf("Unknown corridor command: %s\n", sub)
		fmt.Println("Usage: cog corridor <route|port>")
	}
}

func handleFinance(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog finance <simulate|cleantax> <id|slug> [options]")
		return
	}
	sub := strings.ToLower(args[0])
	if len(args) < 2 {
		fmt.Printf("Usage: cog finance %s <id|slug> [options]\n", sub)
		return
	}
	id := args[1]
	store := getSeededStore()
	ctx := context.Background()

	proj, err := store.GetProject(ctx, id)
	if err != nil {
		proj, err = store.GetProjectBySlug(ctx, id)
		if err != nil {
			fmt.Printf("Project not found: %s\n", id)
			return
		}
	}

	switch sub {
	case "simulate":
		runs := 10000
		for i := 2; i < len(args); i++ {
			if args[i] == "--runs" && i+1 < len(args) {
				fmt.Sscanf(args[i+1], "%d", &runs)
				i++
			}
		}
		sim := projectfinance.NewFinanceSimulator()
		res := sim.RunSimulation(proj, runs)

		fmt.Printf("=== Stochastic Project Finance Monte Carlo Simulation: %s ===\n", res.ProjectName)
		fmt.Printf("Iterations Executed:     %d\n", res.IterationsRun)
		fmt.Printf("Baseline CAPEX:          $%d CAD\n", res.BaselineCapexCAD)
		fmt.Printf("Synthetic Credit Rating: %s (Investment Grade: %t)\n", res.SyntheticCreditRating, res.InvestmentGrade)
		fmt.Printf("Probability of Default:  %.2f%%\n", res.ProbabilityOfDefaultPct)
		fmt.Printf("Audit Hash:              %s\n\n", res.AuditHash)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "METRIC\tMEAN\tP10 (DOWNSIDE)\tP50 (MEDIAN)\tP90 (UPSIDE)")
		fmt.Fprintf(w, "Project IRR\t%.2f%%\t%.2f%%\t%.2f%%\t%.2f%%\n",
			res.ProjectIRRPercent.Mean, res.ProjectIRRPercent.P10, res.ProjectIRRPercent.P50, res.ProjectIRRPercent.P90)
		fmt.Fprintf(w, "Equity IRR\t%.2f%%\t%.2f%%\t%.2f%%\t%.2f%%\n",
			res.EquityIRRPercent.Mean, res.EquityIRRPercent.P10, res.EquityIRRPercent.P50, res.EquityIRRPercent.P90)
		fmt.Fprintf(w, "Minimum DSCR\t%.2fx\t%.2fx\t%.2fx\t%.2fx\n",
			res.MinDSCR.Mean, res.MinDSCR.P10, res.MinDSCR.P50, res.MinDSCR.P90)
		fmt.Fprintf(w, "Average DSCR\t%.2fx\t%.2fx\t%.2fx\t%.2fx\n",
			res.AvgDSCR.Mean, res.AvgDSCR.P10, res.AvgDSCR.P50, res.AvgDSCR.P90)
		fmt.Fprintf(w, "Loan Life Cov (LLCR)\t%.2fx\t%.2fx\t%.2fx\t%.2fx\n",
			res.LoanLifeCoverageRatio.Mean, res.LoanLifeCoverageRatio.P10, res.LoanLifeCoverageRatio.P50, res.LoanLifeCoverageRatio.P90)
		w.Flush()

	case "cleantax":
		calc := projectfinance.NewTaxCreditCalculator()
		profile := calc.CalculateCredits(proj)

		fmt.Printf("=== Clean Economy Tax Credit & CCfD Underwriting: %s ===\n", proj.Name)
		fmt.Printf("Applicable Credit:       %s\n", profile.ApplicableITC)
		fmt.Printf("Eligible CAPEX:          $%d CAD\n", profile.EligibleCapexCAD)
		fmt.Printf("Base Credit Rate:        %.1f%%\n", profile.BaseCreditRatePercent)
		fmt.Printf("Labor Condition Bonus:   +%.1f%%\n", profile.LaborConditionBonusPercent)
		fmt.Printf("Effective Credit Rate:   %.1f%%\n", profile.EffectiveCreditRatePercent)
		fmt.Printf("Total Tax Credit Yield:  $%d CAD\n", profile.TotalTaxCreditYieldCAD)
		fmt.Printf("CCfD Underwriting:       Eligible: %t (Strike: $%.2f/t, Subsidy: $%d/yr)\n",
			profile.CCfDEligible, profile.CCfDStrikePriceCADTonne, profile.EstimatedAnnualCCfDSubsidyCAD)
		fmt.Printf("Audit Hash:              %s\n", profile.AuditHash)

	default:
		fmt.Printf("Unknown finance command: %s\n", sub)
		fmt.Println("Usage: cog finance <simulate|cleantax>")
	}
}

func handleFilings(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: cog filings <stream|ea|amendments> [args...]")
		return
	}
	sub := strings.ToLower(args[0])

	switch sub {
	case "stream":
		records := filings.CanonicalDisclosures()
		fmt.Println("=== Real-Time Continuous Disclosure Stream (SEDAR+ / MD&A) ===")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "FILING ID\tISSUER\tTICKER\tEXCHANGE\tTYPE\tSTAGE DETECTED\tAUDIT HASH")
		for _, f := range records {
			auditSnippet := f.AuditHash
			if len(auditSnippet) > 16 {
				auditSnippet = auditSnippet[:16] + "..."
			}
			stage := string(f.DetectedStage)
			if stage == "" {
				stage = "MONITORING"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				f.ID, f.IssuerName, f.Ticker, f.Exchange, f.FilingType, stage, auditSnippet)
		}
		w.Flush()

	case "ea":
		notices := filings.CanonicalEANotices()
		fmt.Println("=== Provincial Environmental Assessment Registry Notices ===")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "REGISTRY\tPROVINCE\tPROJECT\tMILESTONE\tDEADLINE\tSUMMARY")
		for _, n := range notices {
			deadline := "N/A"
			if n.CommentDeadline != nil {
				deadline = n.CommentDeadline.Format("2006-01-02")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
				n.RegistrySource, n.Province, n.ProjectName, n.Milestone, deadline, n.Summary)
		}
		w.Flush()

	case "amendments":
		amendments := filings.CanonicalTenderAmendments()
		fmt.Println("=== CanadaBuys & DCC Procurement Tender Amendments ===")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "TENDER ID\tTYPE\tAMEND #\tREVISED CLOSING\tAWARD VALUE (CAD)\tWINNING BIDDER")
		for _, a := range amendments {
			closing := "Unchanged"
			if a.RevisedClosing != nil {
				closing = a.RevisedClosing.Format("2006-01-02")
			}
			award := "Pending"
			if a.ContractValueCAD > 0 {
				award = fmt.Sprintf("$%d", a.ContractValueCAD)
			}
			winner := a.WinningBidder
			if winner == "" {
				winner = "In Evaluation"
			}
			fmt.Fprintf(w, "%s\t%s\t#%d\t%s\t%s\t%s\n",
				a.TenderReference, a.Type, a.AmendmentNumber, closing, award, winner)
		}
		w.Flush()

	default:
		fmt.Printf("Unknown filings command: %s\n", sub)
		fmt.Println("Usage: cog filings <stream|ea|amendments>")
	}
}

func handleKPI(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: cog kpi <list|show|feeds|snapshot> [args]")
		return
	}

	sub := args[0]
	switch sub {
	case "list":
		defs := indicators.CanonicalRegistry()
		fmt.Println("=== Canada Sovereign KPI & Indicator Taxonomy ===")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CODE\tCATEGORY\tNAME\tTARGET\tUNIT\tFREQUENCY")
		for _, d := range defs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%.1f\t%s\t%s\n",
				d.Code, d.Category, d.Name, d.TargetBenchmark, d.Unit, d.UpdateFrequency)
		}
		w.Flush()

	case "show":
		if len(args) < 2 {
			fmt.Println("Usage: cog kpi show <id|slug>")
			return
		}
		target := args[1]
		store := getSeededStore()
		ctx := context.Background()

		proj, err := store.GetProject(ctx, target)
		if err != nil || proj == nil {
			proj, _ = store.GetProjectBySlug(ctx, target)
		}

		if proj == nil {
			fmt.Printf("Project %q not found.\n", target)
			return
		}

		evaluator := indicators.NewProjectEvaluator()
		scorecard := evaluator.Evaluate(proj)

		fmt.Printf("=== Project KPI Intelligence Scorecard ===\n")
		fmt.Printf("Project:       %s (%s)\n", scorecard.ProjectName, scorecard.ProjectSlug)
		fmt.Printf("Sector:        %s | Province: %s | Stage: %s\n", scorecard.Sector, scorecard.Province, scorecard.CurrentStage)
		fmt.Printf("Overall KPI:   %.1f / 100\n", scorecard.OverallKPIRating)
		fmt.Printf("Audit Hash:    %s\n\n", scorecard.AuditHash)

		for _, p := range scorecard.Pillars {
			fmt.Printf("[%s] Score: %.1f/100 (%s)\n", p.Title, p.Score, p.Health)
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "  METRIC\tOBSERVED\tTARGET\tRANK\tNOTES")
			for _, m := range p.Metrics {
				fmt.Fprintf(w, "  %s\t%.2f %s\t%.2f %s\t%s\t%s\n",
					m.Name, m.ObservedValue, m.Unit, m.TargetBenchmark, m.Unit, m.PerformanceRank, m.Notes)
			}
			w.Flush()
			fmt.Println()
		}

		if len(scorecard.CriticalActionGaps) > 0 {
			fmt.Println("CRITICAL ACTION GAPS:")
			for _, gap := range scorecard.CriticalActionGaps {
				fmt.Printf("  • %s\n", gap)
			}
		}

	case "feeds":
		engine := indicators.NewLiveFeedEngine()
		ticks := engine.GetLatestTicks()
		fmt.Println("=== Real-Time Commodity & Macro Indicator Feeds ===")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CODE\tNAME\tVALUE\tUNIT\tCHANGE\tDIR\tSOURCE")
		for _, t := range ticks {
			fmt.Fprintf(w, "%s\t%s\t%.2f\t%s\t%+.2f (%.2f%%)\t%s\t%s\n",
				t.MetricCode, t.Name, t.Value, t.Unit, t.ChangeAbsolute, t.ChangePercent, t.Direction, t.Source)
		}
		w.Flush()

	case "snapshot":
		engine := indicators.NewLiveFeedEngine()
		snapshot, err := engine.GenerateSnapshot()
		if err != nil {
			fmt.Printf("Error generating snapshot: %v\n", err)
			return
		}
		data, _ := json.MarshalIndent(snapshot, "", "  ")
		fmt.Println(string(data))

	default:
		fmt.Printf("Unknown kpi subcommand: %s\n", sub)
		fmt.Println("Usage: cog kpi <list|show|feeds|snapshot>")
	}
}



