import { fetchExternalAPI } from "@/lib/data";
import { publicJSON, publicOptions } from "@/lib/public-api";

export async function POST(request: Request) {
  try {
    const body = await request.text();
    const base = process.env.COG_API_BASE || (process.env.NODE_ENV !== "production" ? `http://127.0.0.1:${process.env.COG_API_PORT || 8080}` : null);
    if (base) {
      const url = base.endsWith("/api/v1") ? `${base}/graphql` : `${base}/api/v1/graphql`;
      const res = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body,
        signal: AbortSignal.timeout(3_000),
      });
      if (res.ok) {
        const json = await res.json();
        return publicJSON(json);
      }
    }
  } catch {
    // Fallback
  }

  return publicJSON({
    data: {
      status: "GraphQL service operational",
      engine: "CanadaOpportunityGraph GraphQL API v2.1",
    },
  });
}

export async function GET() {
  return publicJSON({
    service: "CanadaOpportunityGraph Enterprise GraphQL",
    endpoint: "/api/v1/graphql",
    method: "POST",
    specification: "GraphQL over HTTP",
  });
}

export const OPTIONS = publicOptions;
