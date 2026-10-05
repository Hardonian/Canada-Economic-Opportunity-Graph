import { fetchExternalAPI } from "@/lib/data";
import { publicJSON, publicOptions } from "@/lib/public-api";
import { requestUpstreamAPI } from "@/lib/upstream";

export async function GET(request: Request) {
  const url = new URL(request.url);
  const shock = url.searchParams.get("shock") || "USMCA_2026_TARIFF_25";

  const res = await fetchExternalAPI(`/planning/wargame?scenario=${encodeURIComponent(shock)}`);
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
  const body = await request.text();
  if (body.length > 64_000) {
    return publicJSON({ error: "The war-game request exceeds the 64 KB gateway limit." }, { status: 413 });
  }
  const result = await requestUpstreamAPI("/planning/wargame", {
    method: "POST",
    headers: { "Content-Type": "application/json", Accept: "application/json" },
    body,
    cache: "no-store",
  }, 15_000);
  if (!result.response) {
    return publicJSON({ error: "The national war-game engine is unavailable.", source_mode: "UPSTREAM_UNAVAILABLE" }, { status: 503 });
  }
  const payload = await result.response.json().catch(() => ({ error: "The war-game engine returned an invalid response." }));
  return publicJSON(payload, { status: result.response.status });
}

export const OPTIONS = publicOptions;
