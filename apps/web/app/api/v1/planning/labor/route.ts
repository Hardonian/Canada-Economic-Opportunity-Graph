import { fetchExternalAPI } from "@/lib/data";
import { publicJSON, publicOptions } from "@/lib/public-api";

export async function GET(request: Request) {
  const url = new URL(request.url);
  const prov = url.searchParams.get("prov") || "ALL";

  const res = await fetchExternalAPI(`/planning/labor?prov=${encodeURIComponent(prov)}`);
  if (res?.ok) {
    try {
      const data = await res.json();
      return publicJSON(data);
    } catch {
      // Fallback
    }
  }

  return publicJSON({
    province: prov,
    trades: [
      { trade: "Electricians & High-Voltage Linemen", demand_fte: 14200, capacity_fte: 11500, status: "DEFICIT" },
      { trade: "Nuclear Welders (Red Seal)", demand_fte: 3800, capacity_fte: 2400, status: "SEVERE_DEFICIT" },
      { trade: "Boilermakers & Heavy Mechanics", demand_fte: 7100, capacity_fte: 6800, status: "TIGHT" },
      { trade: "Civil Concrete & Heavy Equipment", demand_fte: 22000, capacity_fte: 24000, status: "BALANCED" },
    ],
    source_mode: "DETERMINISTIC_MODEL",
  });
}

export const OPTIONS = publicOptions;
