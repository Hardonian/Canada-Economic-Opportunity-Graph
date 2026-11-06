# Private deal room deployment

The deal room is intentionally disabled until both enterprise identity and the
server-side gateway are configured. It must not be deployed with placeholder
values.

## 1. Configure enterprise identity

Create an Auth0 **Regular Web Application** for the web application, set its
base URL to the public web origin, and configure the callback URL
`https://<web-origin>/auth/callback`. Configure the organization's upstream
OIDC or SAML enterprise connection in Auth0 and require it for the deal-room
users.

The ID token must contain these namespaced custom claims. Claims are read on
the server and are never accepted from the browser.

| Claim | Example | Purpose |
| --- | --- | --- |
| `https://canada-opportunity-graph.ca/tenant_id` | `northstar-capital` | Immutable tenant boundary |
| `https://canada-opportunity-graph.ca/organization` | `Northstar Capital` | Audit context |
| `https://canada-opportunity-graph.ca/clearance` | `PROTECTED_A` | Classification ceiling |
| `https://canada-opportunity-graph.ca/citizenship` | `CA` | Optional audit context |

Clearance values are `UNCLASSIFIED`, `PROTECTED_A`, `PROTECTED_B`,
`PROTECTED_C`, `SECRET`, and `TOP_SECRET`. A missing or invalid clearance is
treated as `UNCLASSIFIED`, so it cannot unlock protected content.

## 2. Set production secrets

Set these values in the web and API deployment environment, using a secrets
manager. Generate distinct random values of at least 32 characters for every
environment. Rotate the gateway and encryption keys through a coordinated
maintenance process; rotating the encryption key without re-encrypting stored
records makes existing records unreadable.

```env
AUTH0_DOMAIN=
AUTH0_CLIENT_ID=
AUTH0_CLIENT_SECRET=
AUTH0_SECRET=
APP_BASE_URL=https://app.example.ca

# Set in both web and API environments.
COG_DEALROOM_GATEWAY_SECRET=

# Set in the API environment only; never expose this to the browser.
COG_DEALROOM_ENCRYPTION_KEY=
COG_DEALROOM_MAX_DOCUMENT_BYTES=26214400
```

`COG_DEALROOM_GATEWAY_SECRET` is the short-lived signed assertion shared by
the web gateway and API. The API rejects direct private-route callers, stale
assertions, unsigned actors, cross-tenant access, and any request without a
valid membership. `COG_DEALROOM_ENCRYPTION_KEY` encrypts document bytes and
deal-room metadata with AES-GCM before they are written to persistent storage.

## 3. Network and storage boundaries

Keep the Go API on a private network or loopback address. The browser should
only reach Next.js, which authenticates the session and signs a narrow actor
assertion for the API. Do not place a CDN, object store, or public static
directory in front of `COG_STORAGE_DIR/dealrooms`.

The persistent API volume must be backed up with encryption and access logs.
Loss of that volume loses private records; loss of the encryption key makes
existing records unrecoverable by design.

## 4. Retention and operating model

Room creators choose a one-, three-, or seven-year retention period. At
startup, every hour, and before every private deal-room request, the API
purges rooms whose retention deadline has passed, including encrypted
documents, membership, and collaboration metadata. A protected purge event is
then recorded. Ensure backup-retention policies are aligned with the same
legal schedule.

Use `OWNER` or `ADMIN` only for tenant personnel who may grant membership.
`MEMBER` can upload and collaborate; `VIEWER` can only read. Every action is
tenant-scoped, membership-scoped, classification-checked, and audit logged.

## Pre-launch checks

- Verify an enterprise SAML/OIDC login receives the four claims above.
- Verify a user from another tenant receives `403` for a room URL and a direct
  API call.
- Upload a permitted document, download it as a viewer with adequate
  clearance, and confirm the audit event.
- Verify the API is not publicly reachable and that the private API returns
  `401` without its signed assertion.
- Test retention in a non-production environment with an expired room.
