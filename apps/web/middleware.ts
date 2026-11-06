import { NextRequest, NextResponse } from "next/server";
import { auth0, auth0Configured } from "@/lib/auth0";

const privatePath = (path: string) => path === "/deal-room" || path.startsWith("/deal-room/") || path === "/api/deal-room" || path.startsWith("/api/deal-room/");

export async function middleware(request: NextRequest) {
  if (!auth0Configured() || !auth0) return NextResponse.next();
  const authResponse = await auth0.middleware(request);
  if (request.nextUrl.pathname.startsWith("/auth/")) return authResponse;
  if (!privatePath(request.nextUrl.pathname)) return authResponse;
  const session = await auth0.getSession(request);
  if (session) return authResponse;
  const login = new URL("/auth/login", request.url);
  login.searchParams.set("returnTo", `${request.nextUrl.pathname}${request.nextUrl.search}`);
  return NextResponse.redirect(login);
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico|ads.txt).*)"],
};
