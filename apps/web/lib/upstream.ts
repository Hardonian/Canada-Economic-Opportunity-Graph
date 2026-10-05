export interface UpstreamLocation {
  configured: boolean;
  apiBase: string | null;
  rootBase: string | null;
}

export interface UpstreamResult {
  response: Response | null;
  latencyMs: number;
  error: "NOT_CONFIGURED" | "INVALID_PATH" | "TIMEOUT" | "UNREACHABLE" | null;
}

type NextFetchInit = RequestInit & {
  next?: { revalidate?: number; tags?: string[] };
};

function configuredBase(): string | undefined {
  return process.env.COG_API_BASE || process.env.NEXT_PUBLIC_API_BASE || process.env.NEXT_PUBLIC_API_URL;
}

export function getUpstreamLocation(): UpstreamLocation {
  const raw = configuredBase() || (process.env.NODE_ENV !== "production" ? `http://localhost:${process.env.COG_API_PORT || 8080}` : undefined);
  if (!raw) return { configured: false, apiBase: null, rootBase: null };

  try {
    const parsed = new URL(raw);
    if ((parsed.protocol !== "http:" && parsed.protocol !== "https:") || parsed.username || parsed.password) {
      return { configured: false, apiBase: null, rootBase: null };
    }
    parsed.hash = "";
    parsed.search = "";
    const normalizedPath = parsed.pathname.replace(/\/+$/, "");
    const hasAPIPrefix = normalizedPath.endsWith("/api/v1");
    const rootPath = hasAPIPrefix ? normalizedPath.slice(0, -"/api/v1".length) : normalizedPath;
    const origin = parsed.origin;
    const rootBase = `${origin}${rootPath}`.replace(/\/$/, "");
    return {
      configured: true,
      rootBase,
      apiBase: hasAPIPrefix ? `${origin}${normalizedPath}` : `${rootBase}/api/v1`,
    };
  } catch {
    return { configured: false, apiBase: null, rootBase: null };
  }
}

function safeRelativePath(path: string): string | null {
  const normalized = path.startsWith("/") ? path : `/${path}`;
  const pathname = normalized.split("?", 1)[0];
  if (normalized.startsWith("//") || normalized.includes("\\") || pathname.split("/").includes("..")) return null;
  return normalized;
}

async function upstreamRequest(base: string | null, path: string, init: NextFetchInit = {}, timeoutMs = 4_000): Promise<UpstreamResult> {
  if (!base) return { response: null, latencyMs: 0, error: "NOT_CONFIGURED" };
  const relative = safeRelativePath(path);
  if (!relative) return { response: null, latencyMs: 0, error: "INVALID_PATH" };

  const started = performance.now();
  try {
    const signal = init.signal || AbortSignal.timeout(timeoutMs);
    const response = await fetch(`${base}${relative}`, { ...init, signal });
    return { response, latencyMs: Math.round(performance.now() - started), error: null };
  } catch (error) {
    const timedOut = error instanceof Error && (error.name === "TimeoutError" || error.name === "AbortError");
    return {
      response: null,
      latencyMs: Math.round(performance.now() - started),
      error: timedOut ? "TIMEOUT" : "UNREACHABLE",
    };
  }
}

export function requestUpstreamAPI(path: string, init: NextFetchInit = {}, timeoutMs?: number): Promise<UpstreamResult> {
  return upstreamRequest(getUpstreamLocation().apiBase, path, init, timeoutMs);
}

export function requestUpstreamRoot(path: string, init: NextFetchInit = {}, timeoutMs?: number): Promise<UpstreamResult> {
  return upstreamRequest(getUpstreamLocation().rootBase, path, init, timeoutMs);
}
