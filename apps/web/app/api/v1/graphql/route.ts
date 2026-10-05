import { publicJSON, publicOptions } from "@/lib/public-api";
import { requestUpstreamAPI } from "@/lib/upstream";

export async function POST(request: Request) {
  const body = await request.text();
  if (body.length > 64_000) {
    return publicJSON({ errors: [{ message: "GraphQL request exceeds the 64 KB gateway limit." }] }, { status: 413 });
  }
  const result = await requestUpstreamAPI("/graphql", {
    method: "POST",
    headers: { "Content-Type": "application/json", Accept: "application/json" },
    body,
    cache: "no-store",
  }, 10_000);
  if (!result.response) {
    return publicJSON({ errors: [{ message: "The GraphQL engine is unavailable." }] }, { status: 503 });
  }
  const payload = await result.response.json().catch(() => ({ errors: [{ message: "The GraphQL engine returned an invalid response." }] }));
  return publicJSON(payload, { status: result.response.status });
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
