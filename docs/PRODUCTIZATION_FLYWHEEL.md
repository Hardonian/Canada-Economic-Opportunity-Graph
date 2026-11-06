# CanadaOpportunityGraph productization flywheel

## Product thesis

CanadaOpportunityGraph is not a generic dashboard and it does not claim to be
a government system. It is an independent Canadian economic operating system:
a shared public evidence layer, accountable operator workflows, and a guarded
private collaboration layer.

The durable advantage is the boundary between those layers. Public records can
be cited, integrated, and challenged. Tenant records remain tenant-scoped,
identity-gated, encrypted, retention-controlled, and excluded from the public
graph.

## The loop now implemented

1. **Observe** — authoritative and reviewed records enter the public graph with
   source URLs, retrieval times, confidence, and hashes.
2. **Verify** — operators inspect evidence and deterministic score coverage;
   the mission queue exposes its routing heuristic rather than hiding it.
3. **Decide** — an operator exports a traceable JSON brief or saves a target
   into the encrypted execution workspace, where scenarios and rationales are
   persisted server-side.
4. **Mobilize** — sensitive documents and collaborators move into the OIDC/
   SAML-protected deal room, never into public graph endpoints.
5. **Learn** — decision records, audit events, source quality, and workflow
   outcomes identify which connectors and data contracts deserve investment.

## What creates the moat

| Compounding asset | Product mechanism | Why it gets stronger |
| --- | --- | --- |
| Trust | Source-level provenance, transparent scoring, and explicit fallback status | Every successful review reinforces the graph as a shared reference point. |
| Workflow | Mission queue → encrypted execution workspace → private deal room | Users stop exporting context across disconnected tools. |
| Network | CEGS-compatible APIs, exports, and public data contracts | More publishers and consumers make the standard more useful. |
| Security | Server-side tenancy, ABAC, short-lived signed gateway assertions, encrypted storage, audit and retention | Institutions can collaborate without donating sensitive data to the public graph. |
| Canadian fit | Bilingual product surface, Canadian project/economic taxonomy, and locally governed deployment choices | The product is specific enough for Canadian institutions while remaining interoperable. |

## Commercial sequence

1. **Land:** a public mission-control pilot for a regional economic agency,
   infrastructure investor, Crown corporation, Indigenous economic-development
   organization, or strategic supplier network.
2. **Prove:** use a bounded decision workflow with a stated success metric:
   time-to-evidence, number of cross-organization handoffs, quality of sourced
   project records, or time-to-diligence.
3. **Expand:** activate tenant deal rooms and enterprise identity for the teams
   that must share non-public materials.
4. **Embed:** integrate their approved data contracts and export formats, while
   keeping client-private data separate from the shared graph.
5. **Platform:** sell secure workspace seats, controlled data connectors,
   implementation, and governed API access—not speculative intelligence claims.

## Production gates before a paid deployment

- Provision the customer’s Auth0 enterprise connection or equivalent upstream
  OIDC/SAML connection and map tenant and clearance claims.
- Put `COG_WORKSPACE_SECRET`, deal-room gateway secrets, and encryption keys in
  a managed secrets service; run planned key rotation with re-encryption.
- Keep the Go API on a private network; expose only the authenticated Next.js
  gateway to browsers.
- Move deal-room files to a production object store/KMS design with malware
  scanning, backup retention, immutable audit export, incident response, and
  legal-hold policy approved for the customer’s jurisdiction.
- Complete privacy, accessibility, security, and procurement reviews before
  claiming any certification or government affiliation.

## Demo script

1. Open **Mission Control** and filter the transparent public mission queue.
2. Open a project dossier and inspect a source link/content hash.
3. Export a JSON mission brief, then queue the project to the execution
   workspace and record the intervention rationale.
4. Show the private deal-room configuration boundary; sign in with a configured
   enterprise identity and demonstrate tenant isolation, clearance checks,
   audit records, and retention.
5. Close on the operations page to demonstrate freshness, connector health,
   serving mode, and no silent substitution of live data.
