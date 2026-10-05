import { SNAPSHOT_MANIFEST, SNAPSHOT_PROJECTS, SNAPSHOT_RADAR_STATS } from "./data";
import { VETTED_SOURCES } from "./source-data";
import { getUpstreamLocation, requestUpstreamAPI, requestUpstreamRoot } from "./upstream";

export interface ConnectorStatus {
  name: string;
  tier: string;
  status: string;
  mode: string;
  last_attempt: string;
  last_success: string;
  last_hash: string;
  duplicate: boolean;
  documents_seen: number;
  documents_changed: number;
  parse_failures: number;
}

export interface ProductCapability {
  id: string;
  label: string;
  endpoint: string;
  href: string;
  status: string;
}

export interface OperationsStatus {
  status: "OPERATIONAL" | "RESILIENT_SNAPSHOT" | "DEGRADED";
  mode: "LIVE_UPSTREAM_API" | "BUNDLED_REVIEWED_SNAPSHOT";
  checked_at: string;
  upstream: {
    configured: boolean;
    reachable: boolean;
    ready: boolean;
    latency_ms: number;
    version?: string;
  };
  runtime: {
    environment: string;
    storage_mode: string;
    scheduler_enabled: boolean;
    scheduler_interval?: string;
  };
  data: {
    total_projects: number;
    total_capex_cad: number;
    active_procurements_count: number;
    evidence_records: number;
    score_records: number;
    canonical_sources: number;
    dataset_version: string;
    generated_at: string;
    age_hours: number;
    freshness: "CURRENT" | "STALE" | "INVALID";
  };
  ingestion: {
    healthy: number;
    degraded: number;
    broken: number;
    connectors: ConnectorStatus[];
  };
  capabilities: ProductCapability[];
}

const CAPABILITY_PAGES: Record<string, string> = {
  "capital-radar": "/",
  "project-intelligence": "/projects",
  planning: "/planning",
  finance: "/finance",
  graph: "/investigate",
  standards: "/cegs",
  events: "/procurement",
};

const FALLBACK_CAPABILITIES: ProductCapability[] = [
  { id: "capital-radar", label: "Capital and project radar", endpoint: "/api/v1/radar", href: "/", status: "ACTIVE" },
  { id: "project-intelligence", label: "Project diligence and provenance", endpoint: "/api/v1/projects", href: "/projects", status: "ACTIVE" },
  { id: "planning", label: "National planning and scenario models", endpoint: "/api/v1/planning/optimize", href: "/planning", status: "ACTIVE" },
  { id: "finance", label: "Project finance and transition taxonomy", endpoint: "/api/v1/finance/transition-taxonomy", href: "/finance", status: "ACTIVE" },
  { id: "graph", label: "GraphQL and graph analytics", endpoint: "/api/v1/graphql", href: "/investigate", status: "ACTIVE" },
  { id: "standards", label: "CEGS exports and interoperability", endpoint: "/api/v1/cegs/export", href: "/cegs", status: "ACTIVE" },
  { id: "events", label: "Server-sent event stream", endpoint: "/api/v1/stream/events", href: "/procurement", status: "ACTIVE" },
];

function record(value: unknown): Record<string, unknown> | null {
  return value && typeof value === "object" && !Array.isArray(value) ? value as Record<string, unknown> : null;
}

async function jsonRecord(response: Response | null): Promise<Record<string, unknown> | null> {
  if (!response?.ok) return null;
  try {
    return record(await response.json());
  } catch {
    return null;
  }
}

function numeric(value: unknown, fallback: number): number {
  return typeof value === "number" && Number.isFinite(value) && value >= 0 ? value : fallback;
}

function normalizeConnectors(value: unknown): ConnectorStatus[] {
  if (!Array.isArray(value)) return [];
  return value.flatMap((candidate) => {
    const item = record(candidate);
    if (!item || typeof item.name !== "string" || typeof item.status !== "string") return [];
    return [{
      name: item.name,
      tier: typeof item.tier === "string" ? item.tier : "",
      status: item.status,
      mode: typeof item.mode === "string" && item.mode ? item.mode : "UNKNOWN",
      last_attempt: typeof item.last_attempt === "string" ? item.last_attempt : "",
      last_success: typeof item.last_success === "string" ? item.last_success : "",
      last_hash: typeof item.last_hash === "string" ? item.last_hash : "",
      duplicate: item.duplicate === true,
      documents_seen: numeric(item.documents_seen, 0),
      documents_changed: numeric(item.documents_changed, 0),
      parse_failures: numeric(item.parse_failures, 0),
    }];
  });
}

function normalizeCapabilities(value: unknown): ProductCapability[] {
  if (!Array.isArray(value)) return FALLBACK_CAPABILITIES;
  const capabilities = value.flatMap((candidate) => {
    const item = record(candidate);
    if (!item || typeof item.id !== "string" || typeof item.label !== "string" || typeof item.endpoint !== "string") return [];
    return [{
      id: item.id,
      label: item.label,
      endpoint: item.endpoint,
      href: CAPABILITY_PAGES[item.id] || "/operations",
      status: typeof item.status === "string" ? item.status : "ACTIVE",
    }];
  });
  return capabilities.length ? capabilities : FALLBACK_CAPABILITIES;
}

export async function getOperationsStatus(): Promise<OperationsStatus> {
  const location = getUpstreamLocation();
  const [healthResult, readyResult, systemResult] = await Promise.all([
    requestUpstreamRoot("/health", { cache: "no-store", headers: { Accept: "application/json" } }, 3_000),
    requestUpstreamRoot("/ready", { cache: "no-store", headers: { Accept: "application/json" } }, 3_000),
    requestUpstreamAPI("/system/status", { cache: "no-store", headers: { Accept: "application/json" } }, 4_000),
  ]);

  const [health, system] = await Promise.all([
    jsonRecord(healthResult.response),
    jsonRecord(systemResult.response),
  ]);
  const runtime = record(system?.runtime);
  const liveData = record(system?.data);
  const ingestion = record(system?.ingestion);
  const connectors = normalizeConnectors(ingestion?.connectors);
  const generatedAt = String(SNAPSHOT_MANIFEST.generated_at);
  const generatedMillis = Date.parse(generatedAt);
  const ageHours = Number.isFinite(generatedMillis) ? Math.max(0, (Date.now() - generatedMillis) / 3_600_000) : Number.POSITIVE_INFINITY;
  const upstreamReachable = Boolean(healthResult.response?.ok && health);
  const upstreamReady = readyResult.response?.ok === true;
  const systemOperational = system?.status === "OPERATIONAL";

  return {
    status: upstreamReady && systemOperational ? "OPERATIONAL" : ageHours <= 720 ? "RESILIENT_SNAPSHOT" : "DEGRADED",
    mode: upstreamReady && system ? "LIVE_UPSTREAM_API" : "BUNDLED_REVIEWED_SNAPSHOT",
    checked_at: new Date().toISOString(),
    upstream: {
      configured: location.configured,
      reachable: upstreamReachable,
      ready: upstreamReady,
      latency_ms: Math.max(healthResult.latencyMs, readyResult.latencyMs, systemResult.latencyMs),
      version: typeof health?.version === "string" ? health.version : undefined,
    },
    runtime: {
      environment: typeof runtime?.environment === "string" ? runtime.environment : process.env.NODE_ENV || "development",
      storage_mode: typeof runtime?.storage_mode === "string" ? runtime.storage_mode : "snapshot",
      scheduler_enabled: runtime?.scheduler_enabled === true,
      scheduler_interval: typeof runtime?.scheduler_interval === "string" ? runtime.scheduler_interval : undefined,
    },
    data: {
      total_projects: numeric(liveData?.total_projects, SNAPSHOT_PROJECTS.length),
      total_capex_cad: numeric(liveData?.total_capex_cad, SNAPSHOT_RADAR_STATS.total_capex_cad),
      active_procurements_count: numeric(liveData?.active_procurements_count, SNAPSHOT_RADAR_STATS.active_procurements_count),
      evidence_records: SNAPSHOT_MANIFEST.record_counts.evidence,
      score_records: SNAPSHOT_MANIFEST.record_counts.scores,
      canonical_sources: VETTED_SOURCES.length,
      dataset_version: String(SNAPSHOT_MANIFEST.dataset_version),
      generated_at: generatedAt,
      age_hours: Number.isFinite(ageHours) ? Math.round(ageHours * 10) / 10 : -1,
      freshness: !Number.isFinite(ageHours) ? "INVALID" : ageHours <= 720 ? "CURRENT" : "STALE",
    },
    ingestion: {
      healthy: numeric(ingestion?.healthy, connectors.filter((item) => item.status === "HEALTHY").length),
      degraded: numeric(ingestion?.degraded, connectors.filter((item) => item.status === "DEGRADED").length),
      broken: numeric(ingestion?.broken, connectors.filter((item) => !["HEALTHY", "DEGRADED"].includes(item.status)).length),
      connectors,
    },
    capabilities: normalizeCapabilities(system?.capabilities),
  };
}
