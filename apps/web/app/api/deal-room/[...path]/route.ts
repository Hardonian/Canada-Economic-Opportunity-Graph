import { proxyDealRoom } from "@/lib/deal-room-gateway";

export const runtime = "nodejs";
export const dynamic = "force-dynamic";

type Context = { params: Promise<{ path: string[] }> };

async function proxy(request: Request, context: Context) {
  const { path } = await context.params;
  return proxyDealRoom(request, path);
}

export const GET = proxy;
export const POST = proxy;
export const HEAD = proxy;
