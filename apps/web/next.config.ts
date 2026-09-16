import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: true,
  // Windows developer machines commonly reject pnpm's traced symlink layout.
  // Linux/Docker keeps the self-contained standalone package used in sovereign deployments.
  output: process.platform === "win32" ? undefined : "standalone",
  poweredByHeader: false,
  compress: true,
  async headers() {
    return [
      {
        source: "/:path*",
        headers: [
          { key: "X-Content-Type-Options", value: "nosniff" },
          { key: "X-Frame-Options", value: "DENY" },
          { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
          { key: "Cross-Origin-Opener-Policy", value: "same-origin" },
          { key: "Cross-Origin-Resource-Policy", value: "same-origin" },
          {
            key: "Permissions-Policy",
            value: "camera=(), microphone=(), geolocation=(), payment=(), usb=(), browsing-topics=()",
          },
          {
            key: "Content-Security-Policy",
            value: [
              "default-src 'self'",
              "base-uri 'self'",
              "form-action 'self'",
              "frame-ancestors 'none'",
              "object-src 'none'",
              "script-src 'self' 'unsafe-inline' https://maps.googleapis.com",
              "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com",
              "img-src 'self' data: blob: https://server.arcgisonline.com https://clarity.maptiles.arcgis.com https://*.basemaps.cartocdn.com https://*.tile.openstreetmap.org https://mt0.google.com https://mt1.google.com https://mt2.google.com https://mt3.google.com https://*.googleapis.com https://*.gstatic.com",
              "font-src 'self' data: https://fonts.gstatic.com",
              "connect-src 'self' https://server.arcgisonline.com https://clarity.maptiles.arcgis.com https://*.basemaps.cartocdn.com https://*.tile.openstreetmap.org https://*.googleapis.com",
              "upgrade-insecure-requests",
            ].join("; "),
          },
        ],
      },
    ];
  },
};

export default nextConfig;
