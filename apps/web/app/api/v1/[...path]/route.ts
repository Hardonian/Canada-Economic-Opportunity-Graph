import { publicOptions } from "@/lib/public-api";
import { requestUpstreamAPI } from "@/lib/upstream";

export const dynamic = "force-dynamic";

const MAX_PROXY_BODY_BYTES = 1 << 20;
const REQUEST_HEADERS = ["accept", "content-type", "if-none-match", "last-event-id", "x-request-id"] as const;
const RESPONSE_HEADERS = [
  "cache-control",
  "content-disposition",
  "content-type",
  "etag",
  "last-modified",
  "ratelimit-limit",
  "ratelimit-policy",
  "ratelimit-remaining",
  "retry-after",
  "x-dlp-redactions-applied",
  "x-request-id",
  "x-security-classification",
] as const;

type RouteContext = { params: Promise<{ path: string[] }> };

function gatewayError(status: number, code: string, message: string, latencyMs = 0): Response {
  return Response.json(
    { error: { code, message }, gateway: { status: "DEGRADED", upstream_latency_ms: latencyMs } },
    { status, headers: { "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff" } },
  );
}

async function proxy(request: Request, context: RouteContext): Promise<Response> {
  const { path } = await context.params;
  if (!path.length || path.some((segment) => !/^[A-Za-z0-9._~-]+$/.test(segment))) {
    return gatewayError(400, "invalid_path", "The requested API path is invalid.");
  }

  const contentLength = Number(request.headers.get("content-length") || 0);
  if (Number.isFinite(contentLength) && contentLength > MAX_PROXY_BODY_BYTES) {
    return gatewayError(413, "body_too_large", "The request body exceeds the gateway limit.");
  }

  if (request.method === "POST") {
    const requestURL = new URL(request.url);
    const origin = request.headers.get("origin");
    const fetchSite = request.headers.get("sec-fetch-site");
    if ((origin && origin !== requestURL.origin) || fetchSite === "cross-site") {
      return gatewayError(403, "forbidden_origin", "Cross-site API mutations are not permitted.");
    }
    const contentType = request.headers.get("content-type") || "";
    if (contentLength > 0 && !contentType.toLowerCase().startsWith("application/json")) {
      return gatewayError(415, "unsupported_media_type", "API mutations require application/json.");
    }
  }

  const headers = new Headers();
  for (const name of REQUEST_HEADERS) {
    const value = request.headers.get(name);
    if (value) headers.set(name, value);
  }

  let body: ArrayBuffer | undefined;
  if (request.method !== "GET" && request.method !== "HEAD") {
    body = await request.arrayBuffer();
    if (body.byteLength > MAX_PROXY_BODY_BYTES) {
      return gatewayError(413, "body_too_large", "The request body exceeds the gateway limit.");
    }
  }

  const incomingURL = new URL(request.url);
  const upstreamPath = `/${path.join("/")}${incomingURL.search}`;
  const isEventStream = upstreamPath.startsWith("/stream/events");
  const result = await requestUpstreamAPI(
    upstreamPath,
    {
      method: request.method,
      headers,
      body,
      cache: "no-store",
      signal: isEventStream ? request.signal : undefined,
    },
    isEventStream ? 300_000 : 15_000,
  );

  if (!result.response) {
    const code = result.error === "NOT_CONFIGURED" ? "upstream_not_configured" : "upstream_unavailable";
    const message = result.error === "NOT_CONFIGURED"
      ? "The application API is not configured for this deployment."
      : "The application API could not be reached.";
    return gatewayError(503, code, message, result.latencyMs);
  }

  const responseHeaders = new Headers();
  for (const name of RESPONSE_HEADERS) {
    const value = result.response.headers.get(name);
    if (value) responseHeaders.set(name, value);
  }
  if (!responseHeaders.has("cache-control")) {
    responseHeaders.set("Cache-Control", request.method === "GET" && !isEventStream
      ? "public, max-age=0, s-maxage=30, stale-while-revalidate=300"
      : "no-store");
  }
  responseHeaders.set("X-COG-Gateway", "upstream");
  responseHeaders.set("X-Upstream-Latency-Ms", String(result.latencyMs));
  responseHeaders.set("X-Content-Type-Options", "nosniff");

  return new Response(request.method === "HEAD" ? null : result.response.body, {
    status: result.response.status,
    statusText: result.response.statusText,
    headers: responseHeaders,
  });
}

export const GET = proxy;
export const HEAD = proxy;
export const POST = proxy;
export const OPTIONS = publicOptions;
