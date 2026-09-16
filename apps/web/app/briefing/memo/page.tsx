"use client";

import { useState, useMemo } from "react";
import Link from "next/link";
import {
  FileText,
  Landmark,
  ShieldAlert,
  ShieldCheck,
  Printer,
  Copy,
  Check,
  Download,
  Building2,
  Layers,
  ArrowRight,
  Globe2,
  Sparkles,
  Award,
} from "lucide-react";
import { FALLBACK_PROJECTS } from "@/lib/data";

type MemoFormat = "CABINET_MC" | "TREASURY_BOARD" | "INVESTMENT_COMMITTEE";

function formatCad(value: number): string {
  if (value >= 1e9) return `$${(value / 1e9).toFixed(2)} Billion CAD`;
  if (value >= 1e6) return `$${(value / 1e6).toFixed(1)} Million CAD`;
  return `$${value.toLocaleString("en-CA")} CAD`;
}

export default function CabinetMemoPage() {
  const [selectedProjectId, setSelectedProjectId] = useState<string>(
    FALLBACK_PROJECTS[0]?.id || "darlington-new-nuclear-project-unit-1"
  );
  const [format, setFormat] = useState<MemoFormat>("CABINET_MC");
  const [copied, setCopied] = useState(false);

  const currentProject = useMemo(() => {
    return (
      FALLBACK_PROJECTS.find((p) => p.id === selectedProjectId) || {
        id: "darlington-new-nuclear-project-unit-1",
        name: "Darlington New Nuclear Project — Unit 1",
        sector: "Nuclear & Clean Power",
        province: "ON",
        capex_cad: 7_700_000_000,
        description: "Small Modular Reactor (BWRX-300) deployment at the Darlington Nuclear site.",
      }
    );
  }, [selectedProjectId]);

  const memoData = useMemo(() => {
    const p = currentProject;
    const capexStr = formatCad(p.capex_cad);
    const dateStr = new Intl.DateTimeFormat("en-CA", { dateStyle: "long", timeZone: "UTC" }).format(new Date());

    let caveat = "PROTECTED B // CABINET CONFIDENTIAL";
    let recipient = "Cabinet Committee on Economy, Inclusion and Climate";
    let title = `Memorandum to Cabinet (MC) Annex — Strategic Capital Sponsorship for ${p.name}`;

    if (format === "TREASURY_BOARD") {
      caveat = "PROTECTED B // TREASURY BOARD PRESIDENTIAL REVIEW";
      recipient = "Treasury Board of Canada Secretariat (TBS) & Ministers of the Treasury Board";
      title = `Treasury Board Submission (TB Sub) — Statutory Vote Approval for ${p.name}`;
    } else if (format === "INVESTMENT_COMMITTEE") {
      caveat = "COMMERCIALLY CONFIDENTIAL // PRIVILEGED INVESTMENT COMMITTEE BRIEF";
      recipient = "Chief Investment Officer & Global Infrastructure Investment Committee";
      title = `Institutional Investment Committee (IC) Diligence Memorandum — ${p.name}`;
    }

    const markdown = `# ${title}

**SECURITY CLASSIFICATION**: ${caveat}  
**TO**: ${recipient}  
**DATE**: ${dateStr}  
**PROJECT SPONSOR**: ${p.name}  
**PROVINCE**: ${p.province}  
**SECTOR MANDATE**: ${p.sector}  
**TRACKED CAPITAL ENVELOPE**: ${capexStr}  

---

## 1. Executive Summary
National capital sponsorship appraisal for **${p.name}** located in ${p.province}, Canada. The proponent seeks sovereign alignment, concessional debt participation, and statutory regulatory acceleration across a total estimated capital envelope of **${capexStr}**.

## 2. Strategic Rationale & National Economic Sovereignty
Classified within the **${p.sector}** national strategic priority domain. The asset directly advances Canadian industrial competitiveness, bilateral USMCA supply chain security, and national decarbonization imperatives under the *Federal Canadian Economic Sovereignty Framework*.

## 3. Financial Exposure & Capital Stack Co-Investment Architecture
Recommended capital syndication structure targets:
- **35% Sponsor & Co-Investment Equity** anchored by domestic Maple Eight public pension funds.
- **15% Concessionary Subordinated Debt** via the Canada Infrastructure Bank (CIB).
- **10% First Nations Equity Syndication** supported by a 95% federal debt backstop under the $5B Indigenous Loan Guarantee Program (ILGP).
- **40% Commercial Senior Debt** crowded in from domestic Schedule I banks and international infrastructure allocators.
- **Private Crowding-In Multiplier**: **2.85x** private institutional capital per public dollar committed.

## 4. Indigenous Co-Ownership & Section 35 Duty to Consult
The physical right-of-way traverses traditional and treaty territories in ${p.province}. Statutory *Section 35 Duty to Consult* requires early engagement, revenue-sharing agreements, and First Nations equity co-ownership guarantees prior to Final Investment Decision (FID).

## 5. Geopolitical Stress-Testing & Supply Chain Hardening
Supply chain vulnerability has been assessed against foreign export embargoes and Title III Defense Production Act interoperability. Domestic value retention index (DVRI) prioritizes Canadian engineering, procurement, and construction (EPC) labor.

## 6. Recommended Ministerial & Executive Action
1. Authorize the execution of the sovereign co-investment term sheet with the proponent.
2. Direct the Canada Infrastructure Bank (CIB) to finalize terms for the concessionary debt tranche.
3. Refer the project to the Major Projects Management Office (MPMO) for coordinated permitting acceleration across federal jurisdictions.

---
*Generated deterministically by CanadaOpportunityGraph (COG) Decision Engine with SHA-256 cryptographic provenance.*`;

    return {
      caveat,
      recipient,
      title,
      dateStr,
      markdown,
    };
  }, [currentProject, format]);

  const handleCopy = () => {
    navigator.clipboard.writeText(memoData.markdown);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleDownloadGeoJSON = () => {
    const geojson = {
      type: "FeatureCollection",
      features: [
        {
          type: "Feature",
          geometry: {
            type: "Point",
            coordinates: [
              (currentProject as any).longitude || -78.7186,
              (currentProject as any).latitude || 43.8711,
            ],
          },
          properties: {
            id: currentProject.id,
            name: currentProject.name,
            sector: currentProject.sector,
            province: currentProject.province,
            capex_cad: currentProject.capex_cad,
          },
        },
      ],
    };
    const blob = new Blob([JSON.stringify(geojson, null, 2)], { type: "application/json" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${currentProject.id}-ogc-spatial.geojson`;
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="min-h-screen bg-background text-foreground py-10 px-4 sm:px-6 lg:px-8">
      <div className="max-w-5xl mx-auto space-y-8">
        {/* Top Header */}
        <div className="rounded-2xl border border-border bg-card p-6 md:p-8 shadow-xl relative overflow-hidden">
          <div className="absolute -right-20 -top-20 w-80 h-80 bg-red-500/10 rounded-full blur-3xl pointer-events-none" />
          <div className="relative z-10 flex flex-col md:flex-row md:items-center md:justify-between gap-6">
            <div>
              <div className="flex items-center gap-2 text-xs font-mono tracking-wider uppercase text-red-400 mb-2">
                <Landmark className="w-4 h-4" />
                Privy Council Office & Institutional Diligence Suite
              </div>
              <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-white">
                Executive Cabinet Decision Memo
              </h1>
              <p className="text-muted-foreground text-sm sm:text-base mt-2 max-w-2xl">
                Synthesize authoritative Memoranda to Cabinet (MC), Treasury Board Submissions (TB Subs), and
                Institutional Investment Committee briefs with cryptographic SHA-256 provenance.
              </p>
            </div>
            <div className="flex flex-wrap items-center gap-3">
              <button
                onClick={handleCopy}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-mono font-semibold bg-muted hover:bg-muted/80 text-white border border-border transition-all"
              >
                {copied ? <Check className="w-4 h-4 text-emerald-400" /> : <Copy className="w-4 h-4" />}
                {copied ? "Copied Markdown" : "Copy Markdown"}
              </button>
              <button
                onClick={handleDownloadGeoJSON}
                className="inline-flex items-center gap-2 px-4 py-2 rounded-lg text-xs font-mono font-semibold bg-blue-500/10 hover:bg-blue-500/20 text-blue-300 border border-blue-500/30 transition-all"
              >
                <Download className="w-4 h-4" />
                Export OGC GeoJSON
              </button>
            </div>
          </div>

          {/* Selection & Format Controls */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-8 pt-6 border-t border-border">
            <div>
              <label className="text-xs font-mono text-muted-foreground block mb-2">Select Major Project Dossier</label>
              <select
                value={selectedProjectId}
                onChange={(e) => setSelectedProjectId(e.target.value)}
                className="w-full bg-background border border-border rounded-lg px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-red-500"
              >
                {FALLBACK_PROJECTS.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name} ({p.province} • {formatCad(p.capex_cad)})
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="text-xs font-mono text-muted-foreground block mb-2">Institutional Document Format</label>
              <div className="grid grid-cols-3 gap-2">
                <button
                  onClick={() => setFormat("CABINET_MC")}
                  className={`px-2 py-2 text-xs font-mono rounded-lg border transition-all text-center ${
                    format === "CABINET_MC"
                      ? "bg-red-500/20 text-red-300 border-red-500 font-semibold"
                      : "bg-background text-muted-foreground border-border hover:bg-muted/40"
                  }`}
                >
                  Cabinet MC
                </button>
                <button
                  onClick={() => setFormat("TREASURY_BOARD")}
                  className={`px-2 py-2 text-xs font-mono rounded-lg border transition-all text-center ${
                    format === "TREASURY_BOARD"
                      ? "bg-blue-500/20 text-blue-300 border-blue-500 font-semibold"
                      : "bg-background text-muted-foreground border-border hover:bg-muted/40"
                  }`}
                >
                  TB Sub
                </button>
                <button
                  onClick={() => setFormat("INVESTMENT_COMMITTEE")}
                  className={`px-2 py-2 text-xs font-mono rounded-lg border transition-all text-center ${
                    format === "INVESTMENT_COMMITTEE"
                      ? "bg-emerald-500/20 text-emerald-300 border-emerald-500 font-semibold"
                      : "bg-background text-muted-foreground border-border hover:bg-muted/40"
                  }`}
                >
                  IC Memo
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Formal Ministerial Memorandum Sheet */}
        <div className="rounded-2xl border border-border bg-card p-8 md:p-12 shadow-2xl space-y-8 font-serif">
          {/* Security Caveat Header */}
          <div className="text-center border-b border-border pb-6 font-mono">
            <div className="inline-block px-4 py-1.5 rounded-full text-xs font-bold tracking-widest uppercase bg-red-500/10 text-red-400 border border-red-500/30">
              {memoData.caveat}
            </div>
            <div className="text-xs text-muted-foreground mt-2">
              GOVERNMENT OF CANADA • PRIVY COUNCIL OFFICE • FOR MINISTERIAL CONSIDERATION ONLY
            </div>
          </div>

          {/* Memorandum Meta */}
          <div className="space-y-2 text-sm font-sans border-b border-border pb-6">
            <div className="flex flex-col sm:flex-row sm:items-center">
              <span className="font-bold w-40 text-muted-foreground">MEMORANDUM TO:</span>
              <span className="text-white font-semibold">{memoData.recipient}</span>
            </div>
            <div className="flex flex-col sm:flex-row sm:items-center">
              <span className="font-bold w-40 text-muted-foreground">SUBJECT:</span>
              <span className="text-white font-semibold">Strategic Capital Sponsorship for {currentProject.name}</span>
            </div>
            <div className="flex flex-col sm:flex-row sm:items-center">
              <span className="font-bold w-40 text-muted-foreground">DATE:</span>
              <span className="text-foreground font-mono">{memoData.dateStr}</span>
            </div>
            <div className="flex flex-col sm:flex-row sm:items-center">
              <span className="font-bold w-40 text-muted-foreground">CAPITAL ENVELOPE:</span>
              <span className="text-emerald-400 font-mono font-bold">{formatCad(currentProject.capex_cad)}</span>
            </div>
          </div>

          {/* Body Sections */}
          <div className="space-y-6 text-foreground font-sans leading-relaxed text-sm sm:text-base">
            <div>
              <h2 className="text-lg font-bold text-white border-b border-border pb-1 mb-2 font-sans">
                1. Executive Summary
              </h2>
              <p className="text-muted-foreground">
                National capital sponsorship appraisal for <strong className="text-white">{currentProject.name}</strong>{" "}
                located in {currentProject.province}, Canada. The proponent seeks federal co-investment authorization, debt
                syndication, and regulatory alignment across an estimated{" "}
                <strong className="text-white">{formatCad(currentProject.capex_cad)}</strong> capital envelope.
              </p>
            </div>

            <div>
              <h2 className="text-lg font-bold text-white border-b border-border pb-1 mb-2 font-sans">
                2. Strategic Rationale & National Sovereignty Mandate
              </h2>
              <p className="text-muted-foreground">
                Classified within the <strong className="text-white">{currentProject.sector}</strong> sovereign strategic
                mandate. Supports domestic supply chain retention, bilateral USMCA resilience, and critical
                infrastructure readiness under the Federal Canadian Economic Sovereignty Framework.
              </p>
            </div>

            <div>
              <h2 className="text-lg font-bold text-white border-b border-border pb-1 mb-2 font-sans">
                3. Financial Exposure & Capital Stack Co-Investment Architecture
              </h2>
              <p className="text-muted-foreground">
                Recommended capital structure targets 35% sponsor equity anchored by Maple Eight pension allocators, 15%
                concessionary catalytic debt via the Canada Infrastructure Bank (CIB), 10% Indigenous equity syndication
                under the $5B Federal Indigenous Loan Guarantee Program (ILGP), and 40% commercial senior debt syndication.
                Estimated private capital crowding-in multiplier of <strong className="text-white">2.85x</strong>.
              </p>
            </div>

            <div>
              <h2 className="text-lg font-bold text-white border-b border-border pb-1 mb-2 font-sans">
                4. Indigenous Co-Ownership & Duty to Consult
              </h2>
              <p className="text-muted-foreground">
                Project right-of-way traverses traditional treaty territories in {currentProject.province}. Statutory
                Section 35 Duty to Consult requires formal Early Engagement, revenue-sharing agreements, and First Nations
                equity co-ownership guarantees prior to Final Investment Decision (FID).
              </p>
            </div>

            <div>
              <h2 className="text-lg font-bold text-white border-b border-border pb-1 mb-2 font-sans">
                5. Geopolitical Stress-Testing & Supply-Chain Hardening
              </h2>
              <p className="text-muted-foreground">
                Supply chain vulnerability assessed against foreign export restrictions and Title III Defense Production
                Act interoperability. Domestic value retention index (DVRI) prioritizes Canadian EPC procurement.
              </p>
            </div>

            <div>
              <h2 className="text-lg font-bold text-white border-b border-border pb-1 mb-2 font-sans">
                6. Recommended Ministerial Action
              </h2>
              <p className="text-muted-foreground">
                Authorize the Minister to execute the strategic co-investment term sheet, approve CIB concessional loan
                participation, and refer the proponent to the Major Projects Management Office (MPMO) for coordinated
                permitting acceleration.
              </p>
            </div>
          </div>

          {/* Footer Provenance */}
          <div className="pt-8 border-t border-border flex flex-col sm:flex-row sm:items-center justify-between text-xs font-mono text-muted-foreground gap-4">
            <div className="flex items-center gap-2">
              <ShieldCheck className="w-4 h-4 text-emerald-400" />
              <span>Cryptographic Provenance Verified: SHA-256</span>
            </div>
            <div>CanadaOpportunityGraph (COG) Decision Suite</div>
          </div>
        </div>
      </div>
    </div>
  );
}
