# Operations Runbook

This runbook covers the checked-in, self-hosted reference deployment. It does
not authorize publishing a new data release or enabling a live credentialed
adapter; those changes remain subject to the source and release policies.

## Prerequisites

- Go `1.26.3` (as declared in `go.mod`)
- Node.js `22` and pnpm `11.8.0` for the web application
- Docker Compose v2 for the container deployment

Install dependencies with `make bootstrap` (or `pnpm --dir apps/web install
--frozen-lockfile` after `go mod download`). Copy `.env.example` to `.env` only
when overriding local Docker Compose defaults. Never place credentials in a
committed environment file.

## Local verification

The standard checks are intentionally offline except for a package manager or
an explicitly live adapter:

```bash
make lint
make test
make release-check
make cegs-validate
make web-test
make web-build
```

`make verify` runs the complete release-quality sequence, including benchmark
tests and immutable-release validation. It does not regenerate a published
snapshot. On Windows, use the equivalent
targets in `scripts/task.ps1`; `lint` and `web-test` are available alongside
the existing build, release, and API targets.

## Data release guardrails

`go run ./scripts/generate_datasets.go` regenerates the mutable current
snapshot from the configured adapters. It validates all artifacts against an
existing immutable release before changing any current artifact. If the
version already exists and any content differs, generation stops without
refreshing `data/public` or `data/cegs`.

Before committing a proposed release, make a controlled update to the declared
dataset version and collection timestamps, then run the generation command and regenerate the web projection with
`pnpm --dir apps/web snapshot:generate`, inspect every diff, and run `make
release-check`. Regenerating the currently published version is intentionally
rejected. The release checker verifies checksums,
counts, ordering, CEGS conformance, cross-record references, immutable-release
parity, and the bundled web snapshots.

## Container deployment

Build and start the least-privileged local stack:

```bash
docker compose up --build
```

The API listens on `127.0.0.1:${API_PORT:-8080}` and exposes `/health`,
`/ready`, and `/metrics`. The web application listens on
`127.0.0.1:${WEB_PORT:-3000}`. The web control plane is available at
`/operations`; its machine-readable projection is `/api/v1/operations/status`.
Compose uses read-only containers, dropped Linux
capabilities, `no-new-privileges`, and bounded temporary filesystems. Configure
the exact browser origin with `CORS_ORIGIN` (or `CORS_ORIGINS` in the API
runtime); do not use a wildcard origin for an authenticated deployment.

Compose defaults to the single-writer persistent profile and runs scheduled
ingestion inside the API process every five minutes. This is deliberate: a
standalone worker with an in-memory store cannot update a different process's
serving store. Set `INGEST_INTERVAL=0` for an immutable snapshot-only runtime.
Live credential-free publisher fetches remain opt-in through the individual
adapter flags in `.env.example`; reviewed fixtures remain the deterministic
fallback.

## Web/API integration contract

Set `COG_API_BASE` on the web runtime to the internal API URL ending in
`/api/v1`. The Next.js same-origin gateway forwards API routes that do not have
a local reviewed-snapshot implementation, including KPI, opportunity,
forecast, finance, graph, compliance, and event-stream endpoints. The gateway
bounds request bodies, forwards only an allowlist of public headers, never
forwards browser-supplied clearance or administration credentials, and returns
an explicit `503` when the Go service is unavailable.

The following endpoints are the fastest deployment checks:

```text
GET /api/v1/health               Web, snapshot, and upstream readiness
GET /api/v1/operations/status    Unified product control plane
GET /api/v1/integrations/status  Connector and optional integration modes
GET /api/v1/system/status        Go runtime, store, scheduler, and connector mesh
```

Helm uses the executable's actual variable names: `ENV`, `STORAGE_MODE`,
`STORAGE_DIR`, `INGEST_INTERVAL`, and `COG_API_BASE`. Persistent mode uses a
`Recreate` API deployment strategy because the built-in WAL is single-writer.

For a production rollout, pin the published image by immutable digest, verify
the Cosign signature, run the health and readiness endpoints behind the
intended proxy, and retain the released `data/releases/<version>` directory
with the image provenance.

## Incident triage and recovery

1. Check `/health` for process availability and `/ready` for loaded data.
   Check the web `/operations` page to distinguish an upstream outage from an
   intentionally resilient snapshot deployment.
2. Review application logs using the request ID returned by the API error
   response. Do not log raw restricted evidence or access tokens.
3. Run `go run ./cmd/releasecheck` before blaming an API projection; a failure
   identifies the affected artifact or reference.
4. Roll back by redeploying the previously signed image and its matching,
   immutable `data/releases/<version>` snapshot. Do not edit a historical
   release in place.
5. Record the cause, affected dataset version, and remediation in the normal
   governance process before publishing a replacement version.

For data handling, security controls, and source-specific terms, see
[SECURITY.md](SECURITY.md), [SOURCE_POLICY.md](SOURCE_POLICY.md), and
[DATA_SOURCES.md](../DATA_SOURCES.md).
