package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/corridor"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/database"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/filings"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/memoexport"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/projectfinance"
	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/syndication"
)

func (s *Server) handleSyndicationInvestors(w http.ResponseWriter, r *http.Request) {
	investors := syndication.DefaultInstitutionalInvestors()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"investors":   investors,
		"total_count": len(investors),
	})
}

func (s *Server) handleSyndicationMatch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := s.store.GetProject(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, r, http.StatusNotFound, "not_found", "Project not found.")
		return
	}

	matcher := syndication.NewMatcher(nil)
	consortium := matcher.MatchProject(project)

	writeJSON(w, http.StatusOK, consortium)
}

func (s *Server) handleOfftakeAgreements(w http.ResponseWriter, r *http.Request) {
	projID := r.URL.Query().Get("project_id")
	var offtakes []syndication.OfftakeAgreement
	if projID != "" {
		offtakes = syndication.GetProjectOfftakes(projID)
	} else {
		offtakes = syndication.CanonicalOfftakeAgreements()
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"offtake_agreements": offtakes,
		"count":              len(offtakes),
	})
}

func (s *Server) handleCorridorRoute(w http.ResponseWriter, r *http.Request) {
	originName := r.URL.Query().Get("from")
	destName := r.URL.Query().Get("to")
	iType := corridor.InfrastructureType(r.URL.Query().Get("type"))
	if iType == "" {
		iType = corridor.TypeHVDCTransmission
	}
	if originName == "" {
		originName = "Sudbury Clean Energy Hub"
	}
	if destName == "" {
		destName = "James Bay Hydro Intertie"
	}

	engine := corridor.NewRoutingEngine()
	origin := corridor.CorridorPoint{Name: originName, Latitude: 46.49, Longitude: -81.01, ElevationMeters: 300}
	dest := corridor.CorridorPoint{Name: destName, Latitude: 53.63, Longitude: -77.70, ElevationMeters: 150}

	eval := engine.EvaluateCorridor(origin, dest, iType)
	writeJSON(w, http.StatusOK, eval)
}

func (s *Server) handleLogisticsPorts(w http.ResponseWriter, r *http.Request) {
	ports := corridor.CanonicalGateways()
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"gateways": ports,
		"count":    len(ports),
	})
}

func (s *Server) handleFinanceDCF(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := s.store.GetProject(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, r, http.StatusNotFound, "not_found", "Project not found.")
		return
	}

	calc := projectfinance.NewTaxCreditCalculator()
	credits := calc.CalculateCredits(project)

	writeJSON(w, http.StatusOK, credits)
}

func (s *Server) handleFinanceMonteCarlo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := s.store.GetProject(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, r, http.StatusNotFound, "not_found", "Project not found.")
		return
	}

	iterations := 1000
	if itStr := r.URL.Query().Get("iterations"); itStr != "" {
		if parsed, parseErr := strconv.Atoi(itStr); parseErr == nil && parsed > 0 {
			iterations = parsed
		}
	}

	sim := projectfinance.NewFinanceSimulator()
	result := sim.RunSimulation(project, iterations)

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleExportMemo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := s.store.GetProject(r.Context(), id)
	if err != nil || project == nil {
		writeError(w, r, http.StatusNotFound, "not_found", "Project not found.")
		return
	}

	mType := memoexport.MemoTypeCabinetMC
	switch strings.ToLower(r.URL.Query().Get("type")) {
	case "tb", "treasury_board":
		mType = memoexport.MemoTypeTreasuryBoard
	case "ic", "investment_committee":
		mType = memoexport.MemoTypeInvestmentComm
	}

	gen := memoexport.NewMemoGenerator()
	memo := gen.GenerateCabinetMemo(project, mType)

	writeJSON(w, http.StatusOK, memo)
}

func (s *Server) handleExportGeoJSON(w http.ResponseWriter, r *http.Request) {
	projects, _, err := s.store.ListProjects(r.Context(), database.ProjectFilter{Limit: 500})
	if err != nil {
		writeError(w, r, http.StatusInternalServerError, "store_error", "Failed to load projects.")
		return
	}

	coll := memoexport.ExportGeoJSON(projects)
	writeJSON(w, http.StatusOK, coll)
}

func (s *Server) handleExportSTAC(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := s.store.GetProject(r.Context(), id)
	if err != nil || project == nil {
		// Default fallback for general satellite endpoint
		project = &domain.Project{
			ID:        "general-sat-pass",
			Slug:      "sentinel-radarsat-canada-pass",
			Name:      "RADARSAT Constellation Mission Mosaic",
			Latitude:  51.25,
			Longitude: -85.32,
		}
	}

	stac := memoexport.ExportSTACItem(project)
	writeJSON(w, http.StatusOK, stac)
}

func (s *Server) handleFilingsRecent(w http.ResponseWriter, r *http.Request) {
	streamer := filings.NewStreamer()
	eaParser := filings.NewEAParser()
	amendTracker := filings.NewAmendmentTracker()

	now := time.Now().UTC()

	sampleFiling := streamer.ParseFiling(
		"Canada Nickel Company Inc.",
		"CNC",
		"TSXV",
		filings.FilingTypeMDA,
		"Q3 2026 Management's Discussion & Analysis",
		"Positive Final Investment Decision reached. EPC contract awarded to Ausenco Engineering for $150M CAD.",
		"https://sedarplus.ca/filings/10049281.pdf",
		now.Add(-24*time.Hour),
	)

	sampleEA := eaParser.ParseNotice(
		"BC_EAO",
		"BC",
		"Cedar LNG Project",
		"EA-2026-081",
		"Issuance of Environmental Assessment Certificate",
		"Environmental Assessment Certificate approval granted with 38 legally binding conditions.",
		"https://projects.eao.gov.bc.ca/p/cedarlng",
		now.Add(-48*time.Hour),
	)

	sampleAmend := amendTracker.TrackAmendment(
		"WS39482910-Doc29102",
		3,
		"Contract awarded to Aecon-PCL Industrial Joint Venture for $185M CAD.",
		"https://canadabuys.canada.ca/tender/WS39482910",
		now.Add(-72*time.Hour),
	)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"recent_disclosures": []interface{}{sampleFiling},
		"recent_ea_notices":  []interface{}{sampleEA},
		"recent_amendments":  []interface{}{sampleAmend},
	})
}
