import "server-only";

import { createCipheriv, createDecipheriv, createHash, randomBytes } from "node:crypto";

export type DecisionStatus = "ADVANCE" | "INVESTIGATE" | "MONITOR";

export interface WorkspaceDecision {
  id: string;
  project_id: string;
  status: DecisionStatus;
  rationale: string;
  created_at: string;
}

export interface WorkspaceScenario {
  catalytic_capital_cad: number;
  permit_months_saved: number;
  domestic_procurement_pct: number;
  indigenous_equity_pct: number;
}

export interface WorkspaceState {
  version: 1;
  selected_project_id: string;
  saved_project_ids: string[];
  scenario: WorkspaceScenario;
  decisions: WorkspaceDecision[];
  updated_at: string;
}

const COOKIE_NAME = "cog_execution_workspace";
const MAX_COOKIE_BYTES = 3_800;
const ID_PATTERN = /^[A-Za-z0-9._~-]{1,200}$/;
const DEFAULT_SCENARIO: WorkspaceScenario = {
  catalytic_capital_cad: 100_000_000,
  permit_months_saved: 6,
  domestic_procurement_pct: 55,
  indigenous_equity_pct: 10,
};

function encryptionSecret(): { value: string; configured: boolean } {
  const configured = process.env.COG_WORKSPACE_SECRET?.trim();
  if (configured && configured.length >= 32) return { value: configured, configured: true };
  // Local development remains usable, but is explicitly reported as DEMO mode by
  // the API. Production deployments must set a unique COG_WORKSPACE_SECRET.
  return { value: "cog-local-development-only-secret-do-not-deploy", configured: false };
}

function key(): Buffer {
  return createHash("sha256").update(encryptionSecret().value).digest();
}

function defaultWorkspace(): WorkspaceState {
  return {
    version: 1,
    selected_project_id: "",
    saved_project_ids: [],
    scenario: { ...DEFAULT_SCENARIO },
    decisions: [],
    updated_at: new Date().toISOString(),
  };
}

function validID(value: unknown): value is string {
  return typeof value === "string" && ID_PATTERN.test(value);
}

function boundedNumber(value: unknown, fallback: number, minimum: number, maximum: number): number {
  return typeof value === "number" && Number.isFinite(value) ? Math.min(maximum, Math.max(minimum, value)) : fallback;
}

function cleanText(value: unknown, maximum: number): string {
  if (typeof value !== "string") return "";
  return value.replace(/[\u0000-\u001F\u007F]/g, " ").replace(/\s+/g, " ").trim().slice(0, maximum);
}

export function normalizeWorkspace(value: unknown): WorkspaceState {
  const base = defaultWorkspace();
  if (!value || typeof value !== "object" || Array.isArray(value)) return base;
  const raw = value as Record<string, unknown>;
  const rawScenario = raw.scenario && typeof raw.scenario === "object" && !Array.isArray(raw.scenario)
    ? raw.scenario as Record<string, unknown>
    : {};
  const saved = Array.isArray(raw.saved_project_ids)
    ? [...new Set(raw.saved_project_ids.filter(validID))].slice(0, 12)
    : [];
  const decisions = Array.isArray(raw.decisions)
    ? raw.decisions.flatMap((item) => {
      if (!item || typeof item !== "object" || Array.isArray(item)) return [];
      const decision = item as Record<string, unknown>;
      const status = decision.status;
      if (!validID(decision.id) || !validID(decision.project_id) ||
          (status !== "ADVANCE" && status !== "INVESTIGATE" && status !== "MONITOR")) return [];
      const rationale = cleanText(decision.rationale, 280);
      if (!rationale) return [];
      const createdAt = typeof decision.created_at === "string" && !Number.isNaN(Date.parse(decision.created_at))
        ? decision.created_at
        : new Date().toISOString();
      return [{ id: decision.id, project_id: decision.project_id, status: status as DecisionStatus, rationale, created_at: createdAt }];
    }).slice(0, 12)
    : [];

  return {
    version: 1,
    selected_project_id: validID(raw.selected_project_id) ? raw.selected_project_id : "",
    saved_project_ids: saved,
    scenario: {
      catalytic_capital_cad: boundedNumber(rawScenario.catalytic_capital_cad, base.scenario.catalytic_capital_cad, 0, 2_000_000_000),
      permit_months_saved: boundedNumber(rawScenario.permit_months_saved, base.scenario.permit_months_saved, 0, 60),
      domestic_procurement_pct: boundedNumber(rawScenario.domestic_procurement_pct, base.scenario.domestic_procurement_pct, 0, 100),
      indigenous_equity_pct: boundedNumber(rawScenario.indigenous_equity_pct, base.scenario.indigenous_equity_pct, 0, 100),
    },
    decisions,
    updated_at: new Date().toISOString(),
  };
}

export function encryptWorkspace(workspace: WorkspaceState): string {
  const serialized = JSON.stringify(workspace);
  const iv = randomBytes(12);
  const cipher = createCipheriv("aes-256-gcm", key(), iv);
  const ciphertext = Buffer.concat([cipher.update(serialized, "utf8"), cipher.final()]);
  return `${iv.toString("base64url")}.${ciphertext.toString("base64url")}.${cipher.getAuthTag().toString("base64url")}`;
}

export function decryptWorkspace(token: string | undefined): WorkspaceState {
  if (!token || token.length > MAX_COOKIE_BYTES) return defaultWorkspace();
  const [ivValue, ciphertextValue, authTagValue, ...extra] = token.split(".");
  if (extra.length || !ivValue || !ciphertextValue || !authTagValue) return defaultWorkspace();
  try {
    const decipher = createDecipheriv("aes-256-gcm", key(), Buffer.from(ivValue, "base64url"));
    decipher.setAuthTag(Buffer.from(authTagValue, "base64url"));
    const plaintext = Buffer.concat([decipher.update(Buffer.from(ciphertextValue, "base64url")), decipher.final()]).toString("utf8");
    return normalizeWorkspace(JSON.parse(plaintext));
  } catch {
    return defaultWorkspace();
  }
}

export function workspaceCookieName(): string {
  return COOKIE_NAME;
}

export function workspaceSecurityMode(): "CONFIGURED" | "DEMO" {
  return encryptionSecret().configured ? "CONFIGURED" : "DEMO";
}

export function workspaceCookieOptions() {
  return {
    httpOnly: true,
    sameSite: "lax" as const,
    secure: process.env.NODE_ENV === "production",
    path: "/",
    maxAge: 60 * 60 * 24 * 30,
  };
}
