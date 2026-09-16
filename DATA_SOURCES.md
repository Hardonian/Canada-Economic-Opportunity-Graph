# Authoritative Data Sources & Provenance Policy

CanadaOpportunityGraph treats provenance as a product boundary. Observed facts, source-reported values, deterministic derivations, and user-authored scenarios are kept distinct. A government publisher does not make every value independently audited, and a cryptographic hash proves content integrity—not factual truth or legal authenticity.

## Active global trade and supply-chain adapter

The `world_bank_global_trade` adapter ingests 11 official Canada observations
covering logistics performance, customs, trade infrastructure, international
shipments, logistics services, tracking, timeliness, trade-to-GDP, merchandise
exports/imports, and high-technology export intensity. It defaults to the
reviewed fixture at `data/fixtures/world_bank_trade_canada.json`; setting
`GLOBAL_TRADE_MODE=live` fetches the same fixed allowlist from the public World
Bank Indicators API.

Each observation becomes a `TradeMetric` with a stable ID and a separate
`Evidence` record containing the normalized-record SHA-256, HTTPS URL, source
locator, retrieval/publication/effective dates, parser version, mapping version,
and pipeline version. The adapter rejects detached lineage before ingestion.
See `docs/DATA_CREDENTIALS.md` for the no-key/key-required access register.

## Newly added sources (2026-09-15)

Four additional sources broaden coverage into political disclosure, federal
procurement, monetary policy, and housing. All four default to reviewed
fixtures so the default build stays deterministic; each has an opt-in live mode
documented below.

| Adapter | Publisher | Live mode env | Access | Epistemic treatment |
| :--- | :--- | :--- | :--- | :--- |
| `bank_of_canada_valet` | Bank of Canada / Banque du Canada (Valet API) | `BANKOFCANADA_LIVE=1` | Key-free, public | `REPORTED` |
| `open_canada_federal_contracts` | Open.Canada Federal Contracts (>$10K) | `FEDERAL_CONTRACTS_LIVE=1` | Key-free, public | `REPORTED` |
| `lobbyist_registry` | Office of the Commissioner of Lobbying of Canada | `LOBBYIST_REGISTRY_LIVE=1` | Key-free, public registry | `REPORTED` |
| `cmhc_housing` | Canada Mortgage and Housing Corporation | `CMHC_HOUSING_LIVE=1` | Key-free, public | `REPORTED` |

**Series covered by `bank_of_canada_valet`:** CPI (Canada, Quebec, Ontario),
GDP at market prices (chained 2017 dollars), and the overnight money market
financing rate (policy rate). The Valet endpoint accepts a fixed allowlist; the
adapter refuses any host outside `www.bankofcanada.ca` and bounds the response
body to 2 MiB.

**Records produced by `open_canada_federal_contracts`:** `Procurement` records
with `RequirementClass = CONFIRMED`, bilingual buyer department names, and
award dates parsed to `*time.Time`. This is federal award data, not tender
notices — do not represent these as open opportunities.

**Records produced by `lobbyist_registry`:** `Entity` records for active
registrants only (inactive registrations are skipped, not silently zeroed).
Jurisdiction is `CA:FED`. Subject matters carry both English and French.

**Records produced by `cmhc_housing`:** `TradeMetric` records for housing
starts (Canada, Ontario, Quebec, British Columbia), rental vacancy rate, and
average 2-bedroom rent. Bilingual indicator names are preserved.

Each of the four applies the same provenance contract as the World Bank
adapter: a separately hashed `Evidence` record with HTTPS URL, publisher,
retrieval timestamp, parser version, and SHA-256 of the parsed payload.

## Active snapshot sources

| Source | Coverage in the default build | Mode | Epistemic treatment |
| :--- | :--- | :--- | :--- |
| [NRCan Major Projects Inventory](https://open.canada.ca/data/en/dataset/f5f2db55-31e4-42fb-8c73-23e1c44de9b2) | 2025–2035 point layer; energy, mining, forestry and clean technology across provinces and territories | Pinned official ArcGIS snapshot; refresh script included | `REPORTED` at record and CAPEX level |
| Canadian Nuclear Safety Commission | Darlington construction milestone | Human-reviewed primary-source snapshot | `VERIFIED` source event; project-level support remains explicit |
| Impact Assessment Agency of Canada | Crawford and Contrecœur regulatory milestones | Human-reviewed primary-source snapshot | `VERIFIED` source events |
| Canada Infrastructure Bank and project issuer releases | Oneida capital and operating milestone | Human-reviewed primary/issuer snapshot | `VERIFIED` or `REPORTED` per assertion |

The adapters for CanadaBuys, CER, IAAC, IDEaS and legacy NRCan fixture formats remain implementation scaffolds. Their empty fixtures are not represented as live polling, complete coverage, or confirmed tender feeds.

The four sources in **Newly added sources** above are NOT scaffolds: they ship
with substantive fixtures (not empty arrays) and tested parsers. Their default
mode is `CURATED_SNAPSHOT`, which means the values come from a reviewed
fixture — not live polling — and are labeled `REPORTED`, never `VERIFIED`.

## Refresh and reproduce

```powershell
# Download only from the fixed Government of Canada ArcGIS host, validate the
# response shape/count, and replace the pinned NRCan snapshot.
.\scripts\snapshot_nrcan_mpi.ps1

# Regenerate checksummed CEGS, JSONL, CSV and GeoJSON releases.
go run .\scripts\generate_datasets.go
```

The runtime uses the pinned snapshot by default. `NewLiveNRCanAdapter` is opt-in, HTTPS-only, host/path allowlisted, redirect bounded, response-size bounded, and covered by tests.

## Epistemic states

- **`VERIFIED`**: the cited primary record was reviewed and directly supports the scoped assertion. It does not mean certified, complete, current forever, or endorsed by the publisher.
- **`SUPPORTED`**: supported by a formal proponent or institutional publication.
- **`REPORTED`**: faithfully normalized from a source inventory but not independently audited record-by-record.
- **`INFERRED`**: deterministically derived from disclosed rules or ontology.
- **`CONFLICTED`**: credible sources disagree and no silent winner was selected.
- **`UNKNOWN`**: the source does not establish the value. Zero is never substituted for an unknown monetary fact in analytical totals.
- **`STALE`**: the evidence has crossed the methodology freshness threshold.
- **`RETRACTED`**: a formerly published assertion has been withdrawn and must not drive current decisions.

## Integrity and licensing

Each normalized source record receives a SHA-256 content hash, stable source identifier, retrieval/effective timestamps, publisher, locator and parser version. Release manifests contain checksums for generated artifacts. These controls detect change and support reproducibility; they are not digital signatures from source agencies.

Upstream records retain their source terms, including the [Open Government Licence – Canada](https://open.canada.ca/en/open-government-licence-canada). Consult each release manifest and source record before redistribution. No government agency endorses this independent platform.
