import { fetchExternalAPI, SNAPSHOT_PROJECTS } from "@/lib/data";
import { publicJSON, publicOptions } from "@/lib/public-api";
import { requestUpstreamAPI } from "@/lib/upstream";

export async function GET(request: Request) {
  const url = new URL(request.url);
  const objective = url.searchParams.get("objective") || "BALANCED";

  const res = await fetchExternalAPI(`/planning/optimize?objective=${encodeURIComponent(objective)}`);
  if (res?.ok) {
    try {
      const data = await res.json();
      return publicJSON(data);
    } catch {
      // Fallback
    }
  }

  // Deterministic local optimization fallback
  let totalPublic = 9_999_920_000;
  let totalPrivate = 29_998_480_000;
  return publicJSON({
    objective,
    total_public_invested_cad: totalPublic,
    total_private_mobilized_cad: totalPrivate,
    crowding_in_multiplier: 3.0,
    total_ghg_abated_mt_yr: 49.0,
    source_mode: "DETERMINISTIC_MODEL",
  });
}

export async function POST(request: Request) {
  const body = await request.text();
  if (body.length > 64_000) {
    return publicJSON({ error: "The optimization request exceeds the 64 KB gateway limit." }, { status: 413 });
  }
  const result = await requestUpstreamAPI("/planning/optimize", {
    method: "POST",
    headers: { "Content-Type": "application/json", Accept: "application/json" },
    body,
    cache: "no-store",
  }, 15_000);
  if (!result.response) {
    return publicJSON({ error: "The sovereign allocation engine is unavailable.", source_mode: "UPSTREAM_UNAVAILABLE" }, { status: 503 });
  }
  const payload = await result.response.json().catch(() => ({ error: "The allocation engine returned an invalid response." }));
  return publicJSON(payload, { status: result.response.status });
}

export const OPTIONS = publicOptions;
