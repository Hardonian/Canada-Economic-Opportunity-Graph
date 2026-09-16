import { getProcurements, SNAPSHOT_MANIFEST } from "@/lib/data";
import { publicJSON, publicOptions } from "@/lib/public-api";

function boundedInteger(value: string | null, fallback: number, maximum: number): number {
  if (!value || !/^\d+$/.test(value)) return fallback;
  const parsed = Number(value);
  return Number.isSafeInteger(parsed) ? Math.min(maximum, parsed) : fallback;
}

export async function GET(request: Request) {
  const search = new URL(request.url).searchParams;
  const limit = Math.max(1, boundedInteger(search.get("limit"), 50, 200));
  const offset = boundedInteger(search.get("offset"), 0, 10_000);
  const query = search.get("q")?.trim().toLowerCase();
  const stage = search.get("stage")?.trim().toLowerCase();
  const buyerType = search.get("buyer_type")?.trim().toLowerCase();

  const allProcurements = await getProcurements();

  const filtered = allProcurements.filter((p) => {
    if (query && !`${p.title} ${p.tender_id} ${p.buyer}`.toLowerCase().includes(query)) return false;
    if (stage && p.stage.toLowerCase() !== stage) return false;
    if (buyerType && p.buyer_type.toLowerCase() !== buyerType) return false;
    return true;
  });

  return publicJSON({
    procurements: filtered.slice(offset, offset + limit),
    total: filtered.length,
    limit,
    offset,
    source_mode: "LIVE_AND_VERIFIED_SNAPSHOT",
    generated_at: SNAPSHOT_MANIFEST.generated_at,
  });
}

export const OPTIONS = publicOptions;
