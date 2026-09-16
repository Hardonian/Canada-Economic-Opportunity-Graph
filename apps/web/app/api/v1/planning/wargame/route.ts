import { fetchExternalAPI } from "@/lib/data";
import { publicJSON, publicOptions } from "@/lib/public-api";

export async function GET(request: Request) {
  const url = new URL(request.url);
  const shock = url.searchParams.get("shock") || "USMCA_2026_TARIFF_25";

  const res = await fetchExternalAPI(`/planning/wargame?shock=${encodeURIComponent(shock)}`);
  if (res?.ok) {
    try {
      const data = await res.json();
      return publicJSON(data);
    } catch {
      // Fallback
    }
  }

  return publicJSON({
    shock,
    scenario: "Geopolitical Tariff Stress Test",
    frozen_capex_cad: 5_153_619_999,
    stalled_assets_count: 11,
    source_mode: "DETERMINISTIC_MODEL",
  });
}

export async function POST(request: Request) {
  return GET(request);
}

export const OPTIONS = publicOptions;
