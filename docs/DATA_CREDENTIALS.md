# Data access and credential register

This project never commits access tokens, passwords, or subscription keys. Most
official sources in the trade and supply-chain registry are deliberately chosen
because their public interfaces require no credential.

| Publisher service | Access status | Credential / configuration |
| --- | --- | --- |
| World Bank Indicators API | Active adapter; public | No credential. Set `GLOBAL_TRADE_MODE=live` to opt into network ingestion; deterministic snapshot mode is the default. |
| Bank of Canada Valet API | Active adapter; public | No credential. Set `BANKOFCANADA_LIVE=1` to opt in; host allowlisted to `www.bankofcanada.ca`. |
| Open.Canada Federal Contracts (>$10K) | Active adapter; public | No credential. Set `FEDERAL_CONTRACTS_LIVE=1` to opt in; hosts allowlisted to `search.open.canada.ca` / `open.canada.ca`. |
| Office of the Commissioner of Lobbying of Canada | Active adapter; public registry | No credential. Set `LOBBYIST_REGISTRY_LIVE=1` to opt in; host allowlisted to `lobbycanada.gc.ca`. |
| Canada Mortgage and Housing Corporation | Active adapter; public | No credential. Set `CMHC_HOUSING_LIVE=1` to opt in; host allowlisted to `www.cmhc-schl.gc.ca`. |
| Statistics Canada Web Data Service | Registered; public | No credential. |
| ISED Trade Data Online | Registered portal/download | No credential for public reports. |
| OECD Data Explorer SDMX API | Registered; public | No credential for public SDMX data. |
| IMF PortWatch Search API | Registered; public catalogue | No credential for public catalogue access. |
| UNCTADstat Data Centre | Registered portal/download | No credential for public downloads. |
| UN Comtrade public API | Registered; public tier | Public requests can operate without a key subject to publisher limits. `UN_COMTRADE_API_KEY` is reserved for a future promoted adapter; obtain any subscription credential directly from UN Comtrade. |
| WTO Timeseries API | Registered; key required | A free key must be issued to an accountable user at <https://apiportal.wto.org/>. Store it only as `WTO_API_KEY` in the deployment secret manager. |
| Google AdSense Auto ads | Optional monetization | Site approval and a public publisher ID are required. Set `NEXT_PUBLIC_ADSENSE_PUBLISHER_ID=ca-pub-################`; never store the Google account password in this project. |

## Live-mode safety contract

Every adapter that can reach the network in live mode enforces, in code:

- HTTPS only
- An explicit host allowlist (never a suffix match — exact host)
- A redirect bound that refuses any hop off the allowlist
- A response-size bound (2 MiB for BoC Valet; 4 MiB for the others)
- A request timeout (8–10 s)

A live failure falls back to the reviewed fixture, sets
`SourceHealth.Status = DEGRADED`, and records `LastError`. It never returns
fabricated or zero-filled data.

## Secret handling

- Local values belong in an ignored `.env` file copied from `.env.example`.
- Production values belong in the hosting provider's encrypted environment
  variables. Scope keys to the minimum environments that need them.
- Logs, health payloads, source profiles, exported datasets, and browser bundles
  must report only `configured` / `not configured`, never a credential value.
- Rotate a key immediately if it appears in Git history, build output, logs, an
  issue, or a client-side bundle.

## AdSense activation

1. Obtain site approval and the publisher ID from AdSense.
2. Add the publisher ID as `NEXT_PUBLIC_ADSENSE_PUBLISHER_ID` in Vercel.
3. Enable Auto ads and the left/right side-rail format in AdSense.
4. Redeploy and verify `/ads.txt` returns the account declaration.
5. Ads remain consent-gated and use Google's `disablePersonalization` privacy
   treatment. A publisher remains responsible for its privacy notice, regional
   consent requirements, tax reporting, and AdSense policy compliance.

No software change can guarantee advertising income. Earnings depend on site
approval, eligible traffic, geography, advertiser demand, viewability, and
policy-compliant use.
