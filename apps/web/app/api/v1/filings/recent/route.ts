import { getRecentFilings } from "@/lib/data";
import { publicJSON, publicOptions } from "@/lib/public-api";

export async function GET() {
  const filings = await getRecentFilings();
  return publicJSON(filings);
}

export const OPTIONS = publicOptions;
