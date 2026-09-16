import { getSignals, SNAPSHOT_MANIFEST } from "@/lib/data";
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
  const type = search.get("type")?.trim().toLowerCase();

  const allSignals = await getSignals();

  const filtered = allSignals.filter((s) => {
    if (type && !s.type.toLowerCase().includes(type)) return false;
    return true;
  });

  return publicJSON({
    signals: filtered.slice(offset, offset + limit),
    total: filtered.length,
    limit,
    offset,
    source_mode: "LIVE_AND_VERIFIED_SNAPSHOT",
    generated_at: SNAPSHOT_MANIFEST.generated_at,
  });
}

export const OPTIONS = publicOptions;
