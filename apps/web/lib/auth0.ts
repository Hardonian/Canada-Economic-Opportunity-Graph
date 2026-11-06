import { Auth0Client } from "@auth0/nextjs-auth0/server";

// Auth0 supplies the OIDC session boundary; enterprise SAML connections are
// configured upstream in Auth0 and arrive here as the same verified session.
export function auth0Configured(): boolean {
  return Boolean(
    process.env.AUTH0_DOMAIN?.trim() &&
    process.env.AUTH0_CLIENT_ID?.trim() &&
    process.env.AUTH0_CLIENT_SECRET?.trim() &&
    process.env.AUTH0_SECRET?.trim(),
  );
}

// Avoid constructing the provider (and emitting a misleading SDK warning) in
// public-only deployments. Private routes refuse access until it is configured.
export const auth0 = auth0Configured() ? new Auth0Client() : null;
