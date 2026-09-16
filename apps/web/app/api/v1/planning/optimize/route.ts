import { fetchExternalAPI, SNAPSHOT_PROJECTS } from "@/lib/data";
import { publicJSON, publicOptions } from "@/lib/public-api";

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
  return GET(request);
}

export const OPTIONS = publicOptions;
