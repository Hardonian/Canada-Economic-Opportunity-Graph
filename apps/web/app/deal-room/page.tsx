import { redirect } from "next/navigation";
import DealRoomClient from "@/components/DealRoomClient";
import { auth0, auth0Configured } from "@/lib/auth0";
import { actorFromUser, dealRoomConfigured } from "@/lib/deal-room-gateway";
import { getProjects } from "@/lib/data";

export const dynamic = "force-dynamic";

export default async function DealRoomPage() {
  if (!dealRoomConfigured()) return <DealRoomConfiguration />;
  const session = await auth0!.getSession();
  if (!session) redirect("/auth/login?returnTo=/deal-room");
  const actor = actorFromUser(session.user);
  if (!actor) return <ClaimConfiguration />;
  const projects = await getProjects();
  return <DealRoomClient actor={actor} projects={projects} />;
}

function DealRoomConfiguration() {
  return <main className="mx-auto max-w-3xl px-4 py-16 sm:px-6"><section className="glass-panel rounded-3xl border border-gold/35 p-8"><p className="font-mono text-xs font-bold tracking-[0.18em] text-gold">PRIVATE DEAL ROOM</p><h1 className="mt-3 text-3xl font-black text-text-main">Enterprise identity is not configured.</h1><p className="mt-4 text-sm leading-relaxed text-text-muted">This deployment intentionally fails closed. Configure Auth0 OIDC plus the deal-room encryption and trusted-gateway secrets to activate tenant-isolated collaboration.</p><pre className="mt-6 overflow-x-auto rounded-xl border border-borderSubtle bg-background p-4 text-xs text-aurora">AUTH0_DOMAIN=…{`\n`}AUTH0_CLIENT_ID=…{`\n`}AUTH0_CLIENT_SECRET=…{`\n`}AUTH0_SECRET=…{`\n`}COG_DEALROOM_GATEWAY_SECRET=…{`\n`}COG_DEALROOM_ENCRYPTION_KEY=…</pre></section></main>;
}

function ClaimConfiguration() {
  return <main className="mx-auto max-w-3xl px-4 py-16 sm:px-6"><section className="glass-panel rounded-3xl border border-gold/35 p-8"><p className="font-mono text-xs font-bold tracking-[0.18em] text-gold">ACCESS CLAIM REQUIRED</p><h1 className="mt-3 text-3xl font-black text-text-main">Your account has no deal-room tenant claim.</h1><p className="mt-4 text-sm leading-relaxed text-text-muted">Ask your identity administrator to include <code className="text-aurora">https://canada-opportunity-graph.ca/tenant_id</code> and the approved clearance claim in the Auth0 ID token.</p></section></main>;
}
