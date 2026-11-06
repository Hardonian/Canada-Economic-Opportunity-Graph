import "server-only";

import { createHmac } from "node:crypto";
import { auth0, auth0Configured } from "@/lib/auth0";
import { requestUpstreamAPI } from "@/lib/upstream";

const TENANT_CLAIM = "https://canada-opportunity-graph.ca/tenant_id";
const ORG_CLAIM = "https://canada-opportunity-graph.ca/organization";
const CLEARANCE_CLAIM = "https://canada-opportunity-graph.ca/clearance";
const CITIZENSHIP_CLAIM = "https://canada-opportunity-graph.ca/citizenship";
const MAX_PRIVATE_BODY_BYTES = 25 * 1024 * 1024;

export type DealRoomActor = {
  user_id: string;
  tenant_id: string;
  organization: string;
  clearance: "UNCLASSIFIED" | "PROTECTED_A" | "PROTECTED_B" | "PROTECTED_C" | "CONFIDENTIAL" | "SECRET" | "TOP_SECRET";
  citizenship: string;
};

const CLEARANCES = new Set<DealRoomActor["clearance"]>(["UNCLASSIFIED", "PROTECTED_A", "PROTECTED_B", "PROTECTED_C", "CONFIDENTIAL", "SECRET", "TOP_SECRET"]);
const ID_PATTERN = /^[A-Za-z0-9._|:-]{1,200}$/;

function stringClaim(value: unknown): string {
  return typeof value === "string" ? value.trim() : "";
}

export function dealRoomConfigured(): boolean {
  return auth0Configured() && Boolean(process.env.COG_DEALROOM_GATEWAY_SECRET?.trim() && process.env.COG_DEALROOM_GATEWAY_SECRET.trim().length >= 32);
}

export function actorFromUser(user: unknown): DealRoomActor | null {
  if (!user || typeof user !== "object") return null;
  const claims = user as Record<string, unknown>;
  const userID = stringClaim(claims.sub);
  const tenantID = stringClaim(claims[TENANT_CLAIM]);
  const organization = stringClaim(claims[ORG_CLAIM]) || tenantID;
  const rawClearance = stringClaim(claims[CLEARANCE_CLAIM]).toUpperCase();
  const clearance = (CLEARANCES.has(rawClearance as DealRoomActor["clearance"]) ? rawClearance : "UNCLASSIFIED") as DealRoomActor["clearance"];
  const citizenship = stringClaim(claims[CITIZENSHIP_CLAIM]).toUpperCase();
  if (!ID_PATTERN.test(userID) || !ID_PATTERN.test(tenantID) || !ID_PATTERN.test(organization)) return null;
  return { user_id: userID, tenant_id: tenantID, organization, clearance, citizenship };
}

function privateResponse(status: number, code: string, message: string): Response {
  return Response.json({ error: { code, message } }, { status, headers: { "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff" } });
}

function isSameOrigin(request: Request): boolean {
  const origin = request.headers.get("origin");
  return request.headers.get("sec-fetch-site") !== "cross-site" && (!origin || origin === new URL(request.url).origin);
}

export async function proxyDealRoom(request: Request, path: string[]): Promise<Response> {
  if (!dealRoomConfigured()) return privateResponse(503, "deal_room_not_configured", "Enterprise identity and the private gateway must be configured before using deal rooms.");
  if (request.method !== "GET" && !isSameOrigin(request)) return privateResponse(403, "forbidden_origin", "Cross-site deal-room mutations are not permitted.");
  if (!path.length || path.some((segment) => !/^[A-Za-z0-9._~-]{1,200}$/.test(segment))) return privateResponse(400, "invalid_path", "The requested private API path is invalid.");

  const session = await auth0!.getSession();
  const actor = actorFromUser(session?.user);
  if (!actor) return privateResponse(403, "tenant_claim_required", "Your authenticated account needs the required tenant and access-classification claims.");

  const contentLength = Number(request.headers.get("content-length") || 0);
  if (Number.isFinite(contentLength) && contentLength > MAX_PRIVATE_BODY_BYTES) return privateResponse(413, "body_too_large", "The private request exceeds the document limit.");
  let body: ArrayBuffer | undefined;
  if (request.method !== "GET" && request.method !== "HEAD") {
    body = await request.arrayBuffer();
    if (body.byteLength > MAX_PRIVATE_BODY_BYTES) return privateResponse(413, "body_too_large", "The private request exceeds the document limit.");
  }

  const encodedActor = Buffer.from(JSON.stringify(actor)).toString("base64url");
  const timestamp = new Date().toISOString();
  const signature = createHmac("sha256", process.env.COG_DEALROOM_GATEWAY_SECRET!.trim()).update(`${timestamp}.${encodedActor}`).digest("base64url");
  const headers = new Headers({
    Accept: request.headers.get("accept") || "application/json",
    "X-COG-Actor": encodedActor,
    "X-COG-Actor-Timestamp": timestamp,
    "X-COG-Actor-Signature": signature,
  });
  const contentType = request.headers.get("content-type");
  if (contentType) headers.set("Content-Type", contentType);
  const result = await requestUpstreamAPI(`/deal-rooms/${path.join("/")}${new URL(request.url).search}`, {
    method: request.method,
    headers,
    body,
    cache: "no-store",
  }, 30_000);
  if (!result.response) return privateResponse(503, "deal_room_api_unavailable", "The private deal-room service is unavailable.");

  const responseHeaders = new Headers({ "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff", "X-COG-Gateway": "authenticated" });
  for (const name of ["content-type", "content-disposition", "content-length", "x-request-id"] as const) {
    const value = result.response.headers.get(name);
    if (value) responseHeaders.set(name, value);
  }
  return new Response(request.method === "HEAD" ? null : result.response.body, { status: result.response.status, headers: responseHeaders });
}
