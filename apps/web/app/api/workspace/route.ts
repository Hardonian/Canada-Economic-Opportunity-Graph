import { cookies } from "next/headers";
import {
  decryptWorkspace,
  encryptWorkspace,
  normalizeWorkspace,
  workspaceCookieName,
  workspaceCookieOptions,
  workspaceSecurityMode,
} from "@/lib/workspace-session";

export const runtime = "nodejs";
export const dynamic = "force-dynamic";

const MAX_BODY_BYTES = 16_384;

function response(body: unknown, status = 200): Response {
  return Response.json(body, {
    status,
    headers: {
      "Cache-Control": "no-store",
      "X-Content-Type-Options": "nosniff",
    },
  });
}

function isSameOrigin(request: Request): boolean {
  const origin = request.headers.get("origin");
  const fetchSite = request.headers.get("sec-fetch-site");
  if (fetchSite === "cross-site") return false;
  return !origin || origin === new URL(request.url).origin;
}

export async function GET() {
  const jar = await cookies();
  const workspace = decryptWorkspace(jar.get(workspaceCookieName())?.value);
  return response({ workspace, security_mode: workspaceSecurityMode() });
}

export async function POST(request: Request) {
  if (!isSameOrigin(request)) return response({ error: "Cross-site workspace updates are not permitted." }, 403);
  const contentType = request.headers.get("content-type")?.toLowerCase() ?? "";
  if (!contentType.startsWith("application/json")) return response({ error: "Workspace updates require application/json." }, 415);
  const contentLength = Number(request.headers.get("content-length") ?? 0);
  if (Number.isFinite(contentLength) && contentLength > MAX_BODY_BYTES) return response({ error: "Workspace update is too large." }, 413);

  const raw = await request.text();
  if (raw.length > MAX_BODY_BYTES) return response({ error: "Workspace update is too large." }, 413);
  let body: unknown;
  try {
    body = JSON.parse(raw);
  } catch {
    return response({ error: "Workspace update is not valid JSON." }, 400);
  }
  const workspace = normalizeWorkspace(body);
  const token = encryptWorkspace(workspace);
  if (token.length > 3_800) return response({ error: "Workspace data exceeds the secure session limit." }, 413);

  const jar = await cookies();
  jar.set(workspaceCookieName(), token, workspaceCookieOptions());
  return response({ workspace, security_mode: workspaceSecurityMode() });
}
