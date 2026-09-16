import { getRadarData } from "@/lib/data";
import { publicJSON, publicOptions } from "@/lib/public-api";

export async function GET() {
  const radarData = await getRadarData();
  return publicJSON(radarData);
}


export const OPTIONS = publicOptions;
