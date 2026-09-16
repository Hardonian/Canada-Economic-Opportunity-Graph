# CanadaOpportunityGraph & CEGS Strategic Roadmap

This document outlines the phased milestone roadmap for CanadaOpportunityGraph and the Canada Economic Graph Schema (CEGS).

---

## Phase 1: Foundation & Baseline (v0.1) — Current Baseline (Completed)

* [x] **Core Ontology**: Project, Organization, Event, Relationship, Evidence, CapitalItem, Procurement, Opportunity, Signal.
* [x] **CEGS Specification 0.1**: JSON Schemas (Draft 2020-12), normative vocabularies, reference examples, RFC process.
* [x] **Reference Go Engine**: Ingestion pipeline, deterministic scoring engines (Buildability, Investability, Supplierability, AI Sovereignty).
* [x] **Authoritative Adapters**: IAAC, CanadaBuys, NRCan Major Projects, IDEaS Defence, CER Facilities.
* [x] **Deterministic Scoring & Propagation**: Downstream requirements graph generator and momentum signal tracking.
* [x] **Zero-Dependency CLI & REST API**: `cog` binary, OpenAPI 3.0 documentation, JSONL public datasets.
* [x] **Interactive Web Platform**: Next.js 15 institutional interface with Capital Radar, Projects Directory, Map, and Sovereignty Radar.

---

## Phase 2: Live Expansion & Provable Provenance (v0.2) — Next 90 Days

* [x] **Automated Gazette Polling**: Recurring scheduled workers scraping provincial gazettes (Ontario Gazette, Gazette officielle du Québec, BC Gazette). Implemented in `internal/gazettepoll/` with change detection, wired into `cmd/worker/main.go`.
* [x] **Cryptographic Merkle Proofs**: Publish daily Merkle tree roots of all ingested evidence hashes to an append-only transparency log. Implemented in `internal/merkle/` with deterministic empty-tree root, wired into `cmd/worker/main.go` and `cmd/api/main.go`.
* [x] **Expanded Capital Stack Intelligence**: Detailed provincial financing programs integrated (Emissions Reduction Alberta, Investissement Québec, BC InBC, and 30+ more). Implemented in `internal/capitalstack/engine.go` with stacking evaluation, mutual exclusivity checks, and CAPEX conflict detection.
* [x] **Indigenous Business Directory Cross-Referencing**: Automated linkage to Indigenous Services Canada (ISC) Indigenous Business Directory for procurement set-asides. Implemented in `adapters/indigenous/` and `internal/indigenouslinker/`, wired into `cmd/worker/main.go`.
* [x] **Multi-Jurisdiction Reconciliation**: Automated resolution of overlapping federal-provincial regulatory filings. Implemented in `internal/reconciliation/` with merge/link/conflict actions, wired into `cmd/worker/main.go` and `cmd/api/main.go`.

---

## Phase 3: CEGS 1.0 Stability & Ecosystem Adoption (v1.0) — Q3/Q4

* [x] **CEGS 1.0 Migration Toolkit**: Migration toolkit for converting legacy JSON Schemas to canonical CEGS 1.0 vocabulary. Implemented in `internal/cegs/migration.go`.
* [x] **CEGS 1.0 Final Standardization**: Lock core vocabulary and release formal migration toolkits. Add conformance test suite for all CEGS resource types.
* [x] **Python & Rust SDKs**: Standalone CEGS client libraries for data science, notebooks, and backend pipelines.
* [x] **Enterprise GraphQL API**: High-throughput graph query API with fine-grained subscription webhooks for moving projects. Implemented in `internal/graphql/` with GET+POST handler, introspection, and all root Query fields wired into `cmd/api/main.go` at `/api/v1/graphql`.
* [x] **Community-Contributed Adapters**: Sandbox registry for community-maintained regional and municipal adapters. Implemented in `internal/adaptersandbox/registry.go` with validation.
* [x] **Decentralized Verifier Nodes**: Multi-party notarization of major project milestone occurrences. Implemented in `internal/verifier/verifier.go` with quorum-based attestation.

---

## Phase 4: Geospatial & Global Trade Engine + Full QA Closure (Completed)

* [x] **High-Resolution Satellite & GIS Mapping**: Dynamic multi-layer map integrating sub-meter optical satellite imagery (Esri World Imagery + Google Maps Hybrid), CartoDB Dark Matter, and OpenStreetMap with live coordinate telemetry HUD.
* [x] **Global Trade Corridors**: Spherical Mercator geodesic routes covering Trans-Pacific minerals/LNG, Trans-Atlantic nuclear fuel, Arctic Northwest Passage sovereignty vector, St. Lawrence Seaway, and USMCA rail arteries.
* [x] **Geopolitical Friction & Choke Point Surveillance**: Critical warning radar for maritime choke points (Bab-el-Mandeb, Panama Canal drought draft limits, Salish Sea) and inter-jurisdictional environmental review standoffs.
* [x] **Strategic Opportunity Zones**: Polygon overlay geometry for primary mineral & clean energy basins (Ring of Fire, Athabasca Basin, James Bay Lithium, Montney-Duvernay CCS, Labrador Trough Green Iron).
* [x] **Canadian Government FIP UX/UI Compliance**: Official bilingual Federal Identity Program banner with Canadian Flag motif, Canada Wordmark, WCAG 2.2 AA bilingual navigation, and OGL-Canada licensing.
* [x] **Zero Open TODOs & Full QA Verification**: 100% test pass with race detector across 40+ packages, benchmark suites, deterministic CEGS specification conformance, and release manifest validation.

---

## Phase 5: Sovereign Intelligence, Predictive Econometrics & National Planning Engine (v2.0) (Completed)

* [x] **Ultimate Beneficial Ownership (UBO) & Sovereign Screening**: Automated tracing of multi-tiered shell structures, foreign state-owned enterprise (SOE) screening, and Investment Canada Act (ICA) national security review tiers (`internal/ubo/`).
* [x] **StatCan Multi-Regional Input-Output (MRIO) Multipliers**: Deterministic Direct, Indirect, and Induced GDP impact, person-years of employment generated, and 3-tier fiscal tax yields under Statistics Canada SUT multipliers (`internal/econometrics/`).
* [x] **Graph Topology & Critical Infrastructure Centrality**: Betweenness centrality, PageRank, single-point-of-failure (SPOF) detection, and cascading failure propagation simulation (`internal/graphanalytics/`).
* [x] **Bayesian Reference Class Megaproject Forecasting**: Bent Flyvbjerg empirical heavy-tailed lognormal hazard models predicting P10–P90 probability density distributions for cost overruns (%) and schedule delay (months) (`internal/risk/`).
* [x] **Electrical Grid Interconnect & Hosting Capacity Engine**: Provincial balancing authority queue delay estimations (IESO, AESO, Hydro-Québec, BC Hydro) and transmission headroom constraints (`internal/gridphysics/`).
* [x] **Earth Observation SAR & Optical Ground-Truth Telemetry**: Sentinel-1 SAR coherence displacement and Sentinel-2 NDVI surface reflection analysis to independently verify claimed construction stage against physical ground reality (`internal/earthobs/`).
* [x] **Multi-Modal NI 43-101 & Capital Waterfall Parsers**: Automated extraction of mining reserves (tonnage, grade, recovery, mine life, strip ratio, NPV8, IRR) and capital stack waterfalls (`internal/documentintelligence/`).
* [x] **Sovereign Capital Allocation Optimizer (MILP)**: Knapsack/MILP solver optimizing federal fiscal deployment across CIB loans, SIF grants, ITCs, and Indigenous loan guarantees to maximize private capital crowding-in and decarbonization (`internal/nationalplanning/optimizer.go`).
* [x] **Geopolitical & Geoeconomic War Game Simulator**: Macroeconomic shock stress-testing (USMCA 25% Tariffs, Critical Mineral Embargo, Arctic Choke Point Closure, Grid Transformer Seizure) and sovereign countermeasure recommendations (`internal/nationalplanning/wargame.go`).
* [x] **National Craft Labor & Trades Collision Aggregator**: Regional quarterly demand vs union hall capacity for Red Seal certified trades (electricians, linemen, nuclear welders, boilermakers, millwrights) (`internal/nationalplanning/labor.go`).
* [x] **First Nations & Inuit Equity Co-Investment Modeler**: Simulating 25%–50% equity co-development, CIB Indigenous Equity Loan debt amortization, and multi-generational sovereign wealth flows (`internal/nationalplanning/indigenous.go`).
* [x] **Institutional Command Center UI**: Fully productized Next.js Command Console on `/planning` with 6 interactive real-time panels and full backward compatibility with the v1.0 workbench (`apps/web/components/NationalPlanningWorkbench.tsx`).

---

## Phase 6: Relational Graph ML, Statutory Intelligence & Enterprise Graph API (v2.1) (Completed)

* [x] **Relational Link Prediction & Graph ML**: Topological graph algorithms (Adamic-Adar index, Resource Allocation, and Jaccard sector affinity) forecasting unannounced EPC partnerships, offtake agreements, and joint ventures (`internal/linkpred/`).
* [x] **Statutory Impact Assessment & Injunction Risk Engine**: Deterministic extraction of Impact Assessment Act (IAA) Section 54 conditions, environmental covenants, and Haida Nation Duty-to-Consult legal risk scores (`internal/documentintelligence/iaac_parser.go`).
* [x] **Universal CLI Subcommand Wiring**: Full native terminal access for `cog planning` (`optimize`, `wargame`, `labor`), `cog project` (`mrio`, `flyvbjerg`, `ubo`, `grid`, `earthobs`), and `cog extract` (`cards`, `ni43101`, `waterfall`, `iaac`) (`cmd/cog/main.go`).
* [x] **Enterprise GraphQL Intelligence API**: Extended GraphQL SDL schema and high-throughput zero-dependency resolvers for `mrio`, `flyvbjerg`, `ubo`, `grid`, `earthobs`, `planningOptimize`, `planningWarGame`, `planningLabor`, and `predictedLinks` (`internal/graphql/`).
* [x] **Executive Project Intelligence Dossier**: Full-width interactive dossier on all project profile pages (`/projects/[id]`) presenting multi-regional input-output GDP yields, Bayesian overrun hazard curves, UBO national security posture, electrical grid headroom, and Sentinel-1 SAR ground-truth telemetry (`apps/web/components/ProjectIntelligenceDossier.tsx`).
* [x] **Canadian International Airport Global Investor Leasing Hubs**: Ingestion and indexing of the National Airports System (NAS) commercial ground lease modernization initiative announced by Mark Carney at the Canada Economic Growth Summit ($18B CAD target capex across YYZ, YVR, YUL, YYC, YEG) with downstream cargo and clean fueling dependency propagation rules.
* [x] **Release 2026-09-16-1 & Full QA Closure**: Cryptographic release verification via `cmd/releasecheck`, complete race-detected test pass (`go test -race ./...`), live REST and GraphQL server query validation, and Next.js 15 production build pass (27/27 pages).



