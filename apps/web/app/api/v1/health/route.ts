import { publicJSON, publicOptions } from "@/lib/public-api";
import { getOperationsStatus } from "@/lib/operations";

export const dynamic = "force-dynamic";

export async function GET() {
  const operations = await getOperationsStatus();
  return publicJSON({
    status: operations.status === "DEGRADED" ? "degraded" : "ok",
    mode: operations.mode,
    checked_at: operations.checked_at,
    generated_at: operations.data.generated_at,
    dataset_version: operations.data.dataset_version,
    projects: operations.data.total_projects,
    canonical_sources: operations.data.canonical_sources,
    evidence_records: operations.data.evidence_records,
    score_records: operations.data.score_records,
    freshness: operations.data.freshness,
    upstream: operations.upstream,
    runtime: operations.runtime,
  }, { headers: { "Cache-Control": "no-store" } });
}

export const OPTIONS = publicOptions;
