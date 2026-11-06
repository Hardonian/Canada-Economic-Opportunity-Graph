import { ImageResponse } from "next/og";

export const alt = "CanadaOpportunityGraph National Mission Control";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default function OpenGraphImage() {
  return new ImageResponse(
    (
      <div style={{ background: "#050b08", color: "#f7faf8", width: "100%", height: "100%", display: "flex", flexDirection: "column", padding: "64px", position: "relative" }}>
        <div style={{ position: "absolute", width: "560px", height: "560px", borderRadius: "999px", background: "rgba(0, 245, 160, 0.16)", filter: "blur(70px)", right: "-100px", top: "-160px" }} />
        <div style={{ display: "flex", alignItems: "center", color: "#00d488", fontSize: 20, letterSpacing: 4 }}>CANADAOPPORTUNITYGRAPH // CEGS</div>
        <div style={{ marginTop: 54, display: "flex", flexDirection: "column", fontSize: 76, lineHeight: 1.05, fontWeight: 800, letterSpacing: -3 }}>
          <span>National</span>
          <span style={{ color: "#00d488" }}>Mission Control</span>
        </div>
        <div style={{ marginTop: 30, display: "flex", maxWidth: "900px", color: "#c7d7cf", fontSize: 28, lineHeight: 1.35 }}>Public evidence → accountable execution → secure collaboration</div>
        <div style={{ marginTop: "auto", display: "flex", alignItems: "center", gap: 18, fontSize: 20 }}>
          <span style={{ border: "1px solid #00d488", color: "#00d488", padding: "10px 16px", borderRadius: 999 }}>TRACEABLE</span>
          <span style={{ border: "1px solid #f6c453", color: "#f6c453", padding: "10px 16px", borderRadius: 999 }}>CANADIAN</span>
          <span style={{ border: "1px solid #5c7d6d", color: "#f7faf8", padding: "10px 16px", borderRadius: 999 }}>MULTI-TENANT</span>
        </div>
      </div>
    ),
    size,
  );
}
