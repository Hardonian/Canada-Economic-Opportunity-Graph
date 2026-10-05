import { getOperationsStatus } from "@/lib/operations";
import { publicJSON, publicOptions } from "@/lib/public-api";

export const dynamic = "force-dynamic";

export async function GET() {
  return publicJSON(await getOperationsStatus(), { headers: { "Cache-Control": "private, no-store" } });
}

export const OPTIONS = publicOptions;
