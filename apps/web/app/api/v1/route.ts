import { publicJSON, publicOptions } from "@/lib/public-api";

export const dynamic = "force-static";

export function GET() {
  return publicJSON({
    name: "CanadaOpportunityGraph public snapshot API",
    version: "v1",
    cegs_version: "0.1",
    documentation: "/apis",
    endpoints: [
      "/api/v1/health",
      "/api/v1/integrations/status",
      "/api/v1/radar",
      "/api/v1/cegs/export",
      "/api/v1/projects",
      "/api/v1/projects/{slug}",
      "/api/v1/projects/{slug}/earthobs",
      "/api/v1/projects/{slug}/flyvbjerg",
      "/api/v1/projects/{slug}/grid",
      "/api/v1/projects/{slug}/mrio",
      "/api/v1/projects/{slug}/ubo",
      "/api/v1/procurements",
      "/api/v1/signals",
      "/api/v1/filings/recent",
      "/api/v1/graphql",
      "/api/v1/planning/optimize",
      "/api/v1/planning/wargame",
      "/api/v1/planning/labor",
      "/api/v1/sources",
      "/api/v1/sources/coverage",
      "/api/v1/sources/{id}",
      "/api/v1/trade/metrics",
      "/api/v1/export/project/{slug}?format=cegs|markdown",
    ],
  });
}

export const OPTIONS = publicOptions;
