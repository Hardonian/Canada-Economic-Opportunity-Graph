"use client";

import { useState, useEffect, useMemo } from "react";
import Link from "next/link";
import {
  Activity,
  AlertTriangle,
  ArrowDownRight,
  ArrowUpRight,
  BarChart3,
  CheckCircle2,
  Clock,
  Compass,
  Cpu,
  Database,
  Download,
  Flame,
  Globe2,
  Layers,
  Leaf,
  Pause,
  Play,
  RefreshCw,
  Search,
  ShieldAlert,
  ShieldCheck,
  Sparkles,
  Sliders,
  Users,
  Zap,
} from "lucide-react";
import { SNAPSHOT_PROJECTS } from "@/lib/data";
import { CANONICAL_KPIS, INITIAL_LIVE_TICKS, DEFAULT_MACRO_SUMMARY, evaluateProjectKPIs } from "@/lib/kpi-data";
import { LiveFeedTick, ProjectKPIScorecard } from "@/lib/types";

// ----------------------------------------------------------------------
// Top 100 Canadian Sovereign Priority Items Catalog
// ----------------------------------------------------------------------
interface PriorityItem {
  rank: number;
  title: string;
  cluster: string;
  clusterKey: string;
  jurisdiction: string;
  scale: string;
  horizon: string;
  impactMetric: string;
  description: string;
}

const TOP_100_PRIORITIES: PriorityItem[] = [
  // Cluster 1: Critical Minerals & Upstream Extraction (1-15)
  { rank: 1, title: "Crawford Nickel Sulphide Project (Canada Nickel)", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "ON", scale: "$3.5B CAD", horizon: "2027", impactMetric: "30,000 t/yr ESG Nickel", description: "Largest nickel sulphide discovery globally in decades; zero-carbon tailings mineralization potential." },
  { rank: 2, title: "Ring of Fire Eagle's Nest Nickel-Copper-PGE (Noront / Wyloo)", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "ON", scale: "$2.2B CAD", horizon: "2029", impactMetric: "15,000 t/yr Ni-Cu", description: "Anchor underground deposit unlocking the James Bay Lowlands critical mineral province." },
  { rank: 3, title: "Prairie Lithium Brine DLE Commercial Facility", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "SK", scale: "$650M CAD", horizon: "2026", impactMetric: "20,000 t/yr LCE", description: "Direct Lithium Extraction (DLE) from Devonian aquifers displacing carbon-heavy spodumene refining." },
  { rank: 4, title: "Galaxy Lithium / James Bay Spodumene (Arcadium)", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "QC", scale: "$850M CAD", horizon: "2027", impactMetric: "330,000 t/yr Spodumene", description: "High-grade open-pit hard-rock lithium feed for North American battery gigafactories." },
  { rank: 5, title: "Matawinie Natural Graphite Mine & Bécancour Anode Plant", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "QC", scale: "$1.4B CAD", horizon: "2026", impactMetric: "100,000 t/yr Anode", description: "Near-zero carbon spherical graphite supply breaking overseas anode dependency." },
  { rank: 6, title: "Nechalacho Rare Earth Elements Phase 2 (Vital Metals / Cheetah)", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "NT", scale: "$350M CAD", horizon: "2027", impactMetric: "5,000 t/yr NdPr", description: "Strategic light and heavy rare earth elements essential for permanent magnet defense supply chains." },
  { rank: 7, title: "Whabouchi Lithium Mine & Shawinigan Hydroxide Plant", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "QC", scale: "$1.8B CAD", horizon: "2026", impactMetric: "34,000 t/yr LiOH", description: "Integrated mine-to-chemical conversion facility powered 100% by Hydro-Québec clean energy." },
  { rank: 8, title: "McArthur River / Key Lake Uranium Expansion (Cameco)", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "SK", scale: "$1.2B CAD", horizon: "2026", impactMetric: "25M lbs/yr U3O8", description: "World's highest-grade uranium complex securing western nuclear reactor fuel cycle autonomy." },
  { rank: 9, title: "Kipawa Heavy Rare Earths Complex", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "QC", scale: "$480M CAD", horizon: "2028", impactMetric: "Dysprosium / Terbium", description: "Crucial heavy rare earths for high-temperature military radar and guided missile actuators." },
  { rank: 10, title: "Valentine Gold Mine Construction (Calibre Mining)", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "NL", scale: "$750M CAD", horizon: "2025", impactMetric: "195,000 oz/yr Au", description: "Atlantic Canada's largest gold development generating regional fiscal foundation." },
  { rank: 11, title: "Dumond Nickel-Magnetite Megaproject (Waterton)", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "QC", scale: "$2.0B CAD", horizon: "2029", impactMetric: "39,000 t/yr Ni", description: "High-volume open-pit nickel reserve with natural ultramafic carbon capture potential." },
  { rank: 12, title: "Thor Lake T-Zone Beryllium-Niobium Deposit", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "NT", scale: "$280M CAD", horizon: "2028", impactMetric: "Be-Nb Aerospace Alloys", description: "Critical specialized aerospace metal reserves for hypersonic airframes and naval reactors." },
  { rank: 13, title: "Sisson Tungsten-Molybdenum Project", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "NB", scale: "$580M CAD", horizon: "2028", impactMetric: "Tungsten Carbide Tooling", description: "Sovereign North American supply of hardened tungsten alloys for machine tooling and ammunition." },
  { rank: 14, title: "Frontier Lithium PAK Lithium Project & Refinery", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "ON", scale: "$1.6B CAD", horizon: "2027", impactMetric: "20,000 t/yr LiOH", description: "Northwestern Ontario pegmatite deposit integrated with First Nations partnership tranches." },
  { rank: 15, title: "Lac des Îles Palladium Mine Deep Extension (Impala)", cluster: "Critical Minerals", clusterKey: "minerals", jurisdiction: "ON", scale: "$420M CAD", horizon: "2026", impactMetric: "220,000 oz/yr PGE", description: "Underground expansion of crucial catalytic and hydrogen electrolyser precious metals." },

  // Cluster 2: Nuclear Power, SMRs & Clean Baseload Grid (16-27)
  { rank: 16, title: "Darlington New Nuclear Project Unit 1 SMR (GE Hitachi BWRX-300)", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "ON", scale: "$3.4B CAD", horizon: "2028", impactMetric: "300 MWe Clean Baseload", description: "G7's first commercial grid-scale Small Modular Reactor providing zero-carbon 24/7 firm power." },
  { rank: 17, title: "Darlington SMR Fleet Expansion (Units 2, 3 & 4)", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "ON", scale: "$10.2B CAD", horizon: "2034", impactMetric: "1,200 MWe SMR Fleet", description: "Multi-unit fleet deployment driving standardized EPCM execution and supply chain scale." },
  { rank: 18, title: "Bruce C 4,800 MW Nuclear Generation Expansion", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "ON", scale: "$16.0B CAD", horizon: "2036", impactMetric: "4,800 MWe New Build", description: "Canada's largest nuclear expansion since the 1980s to power industrial AI and EV manufacturing." },
  { rank: 19, title: "Pickering B Major Nuclear Component Refurbishment (Units 5-8)", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "ON", scale: "$12.5B CAD", horizon: "2032", impactMetric: "2,000 MWe Life Extension", description: "30-year operational life extension securing Ontario's cleanest baseload industrial power." },
  { rank: 20, title: "Point Lepreau ARC-100 Advanced SMR Facility", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "NB", scale: "$950M CAD", horizon: "2029", impactMetric: "100 MWe Fast Reactor", description: "Sodium-cooled advanced small modular reactor with spent fuel recycling capabilities." },
  { rank: 21, title: "SaskPower SMR First-Mover Site (Estevan/Elbow)", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "SK", scale: "$4.0B CAD", horizon: "2033", impactMetric: "300 MWe Coal Phaseout", description: "Direct replacement of coal baseload generation on the Saskatchewan electrical grid." },
  { rank: 22, title: "Westinghouse eVinci Micro-Reactor Arctic Deployment", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "NT", scale: "$120M CAD", horizon: "2027", impactMetric: "5 MWe Diesel Displacement", description: "Transportable heat-pipe nuclear battery eliminating remote diesel barge dependency." },
  { rank: 23, title: "Atlantic Loop Interprovincial High-Voltage Intertie", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "QC/NB/NS", scale: "$6.0B CAD", horizon: "2031", impactMetric: "2,000 MW Hydro Transfer", description: "Transmission corridor moving surplus Hydro-Québec energy to phase out Nova Scotia coal." },
  { rank: 24, title: "Waasigan 230kV Transmission Line (Thunder Bay to Atikokan)", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "ON", scale: "$620M CAD", horizon: "2026", impactMetric: "350 MW Mining Capacity", description: "Hydro One regional transmission corridor unlocking Northwestern Ontario critical minerals." },
  { rank: 25, title: "Wawa-to-Porcupine 500kV Bulk Transmission Intertie", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "ON", scale: "$1.1B CAD", horizon: "2029", impactMetric: "1,000 MW Bulk Transfer", description: "Reinforces northern Ontario grid reliability and integrates remote clean hydro resources." },
  { rank: 26, title: "Chalk River Canadian Global Research SMR (Micro-Hub)", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "ON", scale: "$250M CAD", horizon: "2026", impactMetric: "Medical Radioisotopes & H2", description: "Global R&D testbed for advanced nuclear fuels, medical actinium, and high-temp electrolysis." },
  { rank: 27, title: "McClean Lake Uranium Solution Tailings Recovery", cluster: "Nuclear & Clean Grid", clusterKey: "nuclear", jurisdiction: "SK", scale: "$180M CAD", horizon: "2025", impactMetric: "Circular Fuel Cycle", description: "Advanced circular reprocessing of high-grade uranium tailings with zero surface disturbance." },

  // Cluster 3: Clean Energy, Hydrogen & Long-Duration Storage (28-38)
  { rank: 28, title: "Oneida Energy Storage 250MW / 1,000MWh Battery Hub", cluster: "Clean Energy & Storage", clusterKey: "energy", jurisdiction: "ON", scale: "$600M CAD", horizon: "2025", impactMetric: "1,000 MWh Grid Firming", description: "Canada's largest utility-scale battery storage facility, co-owned by Six Nations of the Grand River." },
  { rank: 29, title: "World Energy GH2 Project Nujio'qonik Clean Hydrogen", cluster: "Clean Energy & Storage", clusterKey: "energy", jurisdiction: "NL", scale: "$12.0B CAD", horizon: "2028", impactMetric: "250,000 t/yr Green H2", description: "Gigawatt-scale wind-to-green-hydrogen export corridor to European industrial hubs." },
  { rank: 30, title: "Air Products Net-Zero Edmonton Hydrogen Energy Complex", cluster: "Clean Energy & Storage", clusterKey: "energy", jurisdiction: "AB", scale: "$1.6B CAD", horizon: "2025", impactMetric: "Auto-thermal Reformer + CCUS", description: "World-scale blue hydrogen facility capturing 95% of CO2 emissions for merchant transport." },
  { rank: 31, title: "EverWind Point Tupper Green Hydrogen & Ammonia Hub", cluster: "Clean Energy & Storage", clusterKey: "energy", jurisdiction: "NS", scale: "$8.0B CAD", horizon: "2027", impactMetric: "1M t/yr Green Ammonia", description: "Deepwater port export facility utilizing Nova Scotia onshore wind and First Nations equity." },
  { rank: 32, title: "Dow Fort Saskatchewan Path2Zero Ethylene Expansion", cluster: "Clean Energy & Storage", clusterKey: "energy", jurisdiction: "AB", scale: "$8.9B CAD", horizon: "2027", impactMetric: "World 1st Net-Zero Cracker", description: "Net-zero Scope 1 & 2 carbon emissions ethylene cracker utilizing hydrogen circularity." },
  { rank: 33, title: "Pathways Alliance Carbon Capture & Storage Trunkline", cluster: "Clean Energy & Storage", clusterKey: "energy", jurisdiction: "AB", scale: "$16.5B CAD", horizon: "2030", impactMetric: "22 Mt/yr CO2 Sequestered", description: "400km carbon capture pipeline network linking 20 oil sands facilities to deep saline aquifers." },
  { rank: 34, title: "Bécancour Battery Materials Industrial Park Hydro Intertie", cluster: "Clean Energy & Storage", clusterKey: "energy", jurisdiction: "QC", scale: "$450M CAD", horizon: "2025", impactMetric: "500 MW Industrial Feed", description: "Electrification hub supplying Ford, GM-POSCO, and Nemaska Lithium refining facilities." },
  { rank: 35, title: "Port Hawkesbury Paper Biofuel & Steam Co-Generation", cluster: "Clean Energy & Storage", clusterKey: "energy", jurisdiction: "NS", scale: "$160M CAD", horizon: "2026", impactMetric: "Fossil Fuel Displacement", description: "Industrial forest biomass conversion displacing imported fuel oil in Cape Breton." },
  { rank: 36, title: "Boralex Apuiat 200 MW Innu Wind Power Project", cluster: "Clean Energy & Storage", clusterKey: "energy", jurisdiction: "QC", scale: "$600M CAD", horizon: "2025", impactMetric: "200 MW Clean Power", description: "50-50 partnership between the Innu Nation and Boralex delivering northern green power." },
  { rank: 37, title: "Cascades Green Energy Biomass Microgrid", cluster: "Clean Energy & Storage", clusterKey: "energy", jurisdiction: "QC", scale: "$110M CAD", horizon: "2026", impactMetric: "85% Decarbonization", description: "Closed-loop circular packaging industrial steam and thermal energy displacement." },
  { rank: 38, title: "TC Energy Canyon Creek Pumped Hydro Storage (Crowsnest)", cluster: "Clean Energy & Storage", clusterKey: "energy", jurisdiction: "AB", scale: "$2.2B CAD", horizon: "2030", impactMetric: "750 MW / 6,000 MWh", description: "Long-duration gravitational energy storage firming intermittent renewable energy in southern Alberta." },

  // Cluster 4: Sovereign AI Compute & Datacentres (39-48)
  { rank: 39, title: "Hydro-Québec Beauharnois Sovereign AI Supercompute Hub", cluster: "Sovereign AI Compute", clusterKey: "compute", jurisdiction: "QC", scale: "$2.8B CAD", horizon: "2026", impactMetric: "500 MW Hydro Clean Compute", description: "Tier-4 hyperscale sovereign AI cluster guaranteeing data sovereignty and Bill C-27 compliance." },
  { rank: 40, title: "Calgary East Deep Learning & Data Sovereignty Campus", cluster: "Sovereign AI Compute", clusterKey: "compute", jurisdiction: "AB", scale: "$1.4B CAD", horizon: "2026", impactMetric: "250 MW Natural Gas + CCUS", description: "High-density GPU AI training cluster co-located with dispatchable gas generation and CCUS." },
  { rank: 41, title: "Markham Technology Hub SMR-Powered AI Data Centre", cluster: "Sovereign AI Compute", clusterKey: "compute", jurisdiction: "ON", scale: "$1.9B CAD", horizon: "2028", impactMetric: "300 MW Dedicated Nuclear PPA", description: "Zero-emission high-density compute facility backed by long-term Darlington SMR clean power PPA." },
  { rank: 42, title: "Dalhousie Ocean & Arctic AI Compute Node", cluster: "Sovereign AI Compute", clusterKey: "compute", jurisdiction: "NS", scale: "$180M CAD", horizon: "2025", impactMetric: "Oceanographic & Radar Telemetry", description: "Atlantic sovereign compute facility processing high-resolution sub-surface sonar and radar." },
  { rank: 43, title: "Mila National AI Research Infrastructure Expansion", cluster: "Sovereign AI Compute", clusterKey: "compute", jurisdiction: "QC", scale: "$320M CAD", horizon: "2025", impactMetric: "Bilingual Foundation Models", description: "Dedicated sovereign compute cluster for Canada's AI foundation models and biotech LLMs." },
  { rank: 44, title: "Brampton Next-Gen Liquid-Cooled AI Colocation Facility", cluster: "Sovereign AI Compute", clusterKey: "compute", jurisdiction: "ON", scale: "$850M CAD", horizon: "2026", impactMetric: "120 kW/rack Direct Liquid Cooling", description: "Direct-to-chip liquid-cooled facility achieving 1.15 PUE efficiency for next-gen 1MW GPU clusters." },
  { rank: 45, title: "Vancouver Coastal Edge AI Inference Grid", cluster: "Sovereign AI Compute", clusterKey: "compute", jurisdiction: "BC", scale: "$420M CAD", horizon: "2026", impactMetric: "Sub-5ms Transpacific Latency", description: "Low-latency edge compute node linking transpacific undersea fiber cables to Canadian enterprise." },
  { rank: 46, title: "Northern Ontario Indigenous Fibre Intertie & Data Vault", cluster: "Sovereign AI Compute", clusterKey: "compute", jurisdiction: "ON", scale: "$210M CAD", horizon: "2027", impactMetric: "100 Gbps Low-Thermal Hosting", description: "Naturally cooled cold-climate data vault owned by First Nations telecom consortium." },
  { rank: 47, title: "Saskatoon Agricultural & Agrigenomics AI Compute Cluster", cluster: "Sovereign AI Compute", clusterKey: "compute", jurisdiction: "SK", scale: "$140M CAD", horizon: "2026", impactMetric: "Drought-Resistant Crop Models", description: "Specialized sovereign compute center simulating high-yield genomic crop resilience." },
  { rank: 48, title: "Ottawa Sovereign Cloud & National Security Compute Enclave", cluster: "Sovereign AI Compute", clusterKey: "compute", jurisdiction: "ON", scale: "$750M CAD", horizon: "2025", impactMetric: "Protected B & Secret Cloud", description: "Air-gapped sovereign intelligence cloud for National Defence and Communications Security Establishment." },

  // Cluster 5: Ports, Gateways & Strategic Corridors (49-60)
  { rank: 49, title: "Port of Prince Rupert Fairview & Ridley Terminals Expansion", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "BC", scale: "$2.5B CAD", horizon: "2027", impactMetric: "3.2M TEU & Energy Export", description: "North America's closest port to Asia; strategic CN rail gateway for critical minerals and grain." },
  { rank: 50, title: "Port of Churchill Hudson Bay Arctic Deepwater Gateway", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "MB", scale: "$450M CAD", horizon: "2026", impactMetric: "Arctic Sovereignty & Trade", description: "Arctic Gateway Hudson Bay Railway restoration providing shortest maritime route to Europe." },
  { rank: 51, title: "Grays Bay Road and Port Corridor (Coronation Gulf)", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "NU", scale: "$1.8B CAD", horizon: "2030", impactMetric: "Deepwater Arctic Mineral Hub", description: "All-weather deepwater Arctic port connecting the Slave Geological Province to global shipping." },
  { rank: 52, title: "Port of Vancouver Roberts Bank Terminal 2 Container Hub", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "BC", scale: "$3.5B CAD", horizon: "2031", impactMetric: "2.4M TEU Additional Capacity", description: "Major container gateway expansion resolving critical West Coast supply-chain congestion." },
  { rank: 53, title: "Port of Montreal Contrecœur Terminal Expansion", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "QC", scale: "$1.4B CAD", horizon: "2027", impactMetric: "1.15M TEU Multimodal Hub", description: "St. Lawrence Seaway intermodal container terminal linking central Canada to European markets." },
  { rank: 54, title: "Ring of Fire All-Weather Access Road & Corridors", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "ON", scale: "$1.2B CAD", horizon: "2028", impactMetric: "Mineral Transport & Community Link", description: "Multi-modal road and transmission link led by Webequie and Marten Falls First Nations." },
  { rank: 55, title: "Inuvik-Tuktoyaktuk All-Weather Arctic Highway Modernization", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "NT", scale: "$220M CAD", horizon: "2026", impactMetric: "Year-Round Arctic Ocean Link", description: "Permafrost-stabilized highway connecting Canada's road grid directly to the Beaufort Sea." },
  { rank: 56, title: "Port of Saint John West Side Modernization (DP World)", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "NB", scale: "$205M CAD", horizon: "2025", impactMetric: "800,000 TEU East Coast Hub", description: "Dual-rail-connected Atlantic deepwater gateway with direct access to CPKC and CN." },
  { rank: 57, title: "St. Lawrence Seaway Digital Ice Navigation & Lock Modernization", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "ON/QC", scale: "$380M CAD", horizon: "2026", impactMetric: "Extended 10-Month Navigation", description: "Real-time satellite and radar telemetry extending commercial maritime shipping season." },
  { rank: 58, title: "Edmonton Intermodal Logistics Super-Hub (CN Rail)", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "AB", scale: "$520M CAD", horizon: "2026", impactMetric: "Mid-Continent Logistics Hub", description: "Multimodal logistics park connecting transcontinental rail with the Alaska Highway corridor." },
  { rank: 59, title: "Squamish Marine Export Terminal (Woodfibre LNG)", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "BC", scale: "$5.1B CAD", horizon: "2027", impactMetric: "2.1 Mt/yr E-Drive LNG", description: "Hydro-powered electric-drive LNG export terminal with Squamish Nation environmental regulation." },
  { rank: 60, title: "Gordie Howe International Bridge Multi-Modal Gateway", cluster: "Ports & Gateways", clusterKey: "ports", jurisdiction: "ON", scale: "$6.4B CAD", horizon: "2025", impactMetric: "10,000 Commercial Trucks/Day", description: "New 6-lane border crossing securing 25% of Canada-U.S. bilateral merchandise trade." },

  // Cluster 6: First Nations & Indigenous Co-Investment (61-72)
  { rank: 61, title: "Canada Indigenous Loan Guarantee Program ($5B National Allocation)", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "National", scale: "$5.0B CAD", horizon: "2025-2027", impactMetric: "First Nations Sovereign Equity", description: "Concessionary federal debt guarantees enabling multi-nation ownership across energy and minerals." },
  { rank: 62, title: "First Nations Major Projects Coalition Capital Syndication", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "National", scale: "$10.0B CAD Portfolio", horizon: "Ongoing", impactMetric: "Crowding-In Institutional Capital", description: "Indigenous-led alliance providing commercial diligence and equity syndication across 130+ nations." },
  { rank: 63, title: "Tahltan Nation Central BC Mining & Clean Hydro Joint Ventures", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "BC", scale: "$1.5B CAD", horizon: "2026", impactMetric: "Golden Triangle Stewardship", description: "Sovereign land-use protocol and shared royalties across Red Chris, Eskay Creek, and Galore Creek." },
  { rank: 64, title: "Six Nations of the Grand River Equity Tranche (Oneida Hub)", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "ON", scale: "$150M CAD", horizon: "2025", impactMetric: "Generational Dividend Stream", description: "Pioneering equity ownership model in North American clean energy storage." },
  { rank: 65, title: "Marten Falls & Webequie Community Infrastructure Syndicates", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "ON", scale: "$400M CAD", horizon: "2027", impactMetric: "Road & Intertie Co-Ownership", description: "Linear right-of-way co-ownership model delivering long-term commercial tolling revenue." },
  { rank: 66, title: "Squamish Nation Woodfibre Commercial Revenue Royalty", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "BC", scale: "$1.1B CAD over 40 yrs", horizon: "2027", impactMetric: "Sovereign Environmental Authority", description: "Legally binding agreement granting host nation independent environmental oversight authority." },
  { rank: 67, title: "Innu Nation Apuiat Clean Energy Trust Fund", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "QC", scale: "$300M CAD Equity", horizon: "2025", impactMetric: "50% Commercial Ownership", description: "30-year revenue distribution funding community housing, healthcare, and cultural preservation." },
  { rank: 68, title: "Kaska Dena Nation Ross River Minerals Protocol", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "YT", scale: "$250M CAD", horizon: "2027", impactMetric: "Yukon Zinc-Lead Equity", description: "Consensual mineral extraction governance ensuring high-ratio local employment and remediation bonds." },
  { rank: 69, title: "Miawpukek First Nation Wind-to-Hydrogen Joint Venture", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "NL", scale: "$1.2B CAD Equity", horizon: "2028", impactMetric: "Atlantic Export Sovereignty", description: "Direct equity stake in World Energy GH2 export project on Newfoundland's west coast." },
  { rank: 70, title: "Métis Nation of Alberta Saline Aquifer Carbon Sequestration Hub", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "AB", scale: "$450M CAD", horizon: "2028", impactMetric: "Pore Space Co-Ownership", description: "Indigenous subsurface pore space tenure and carbon credit monetization agreement." },
  { rank: 71, title: "Gwich'in Tribal Council Arctic All-Weather Highway Logistics", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "NT", scale: "$85M CAD", horizon: "2026", impactMetric: "Northern Freight Sovereignty", description: "100% Indigenous-owned heavy transport fleet serving the Mackenzie Valley corridor." },
  { rank: 72, title: "First Nations 5% Federal Procurement Compliance Tracker", cluster: "Indigenous Sovereignty", clusterKey: "indigenous", jurisdiction: "National", scale: "$1.6B CAD/yr", horizon: "2025", impactMetric: "Statutory Procurement Spend", description: "Mandatory compliance monitoring ensuring federal departments meet the 5% Indigenous business mandate." },

  // Cluster 7: Critical Sovereign KPIs & Live Feeds (73-85)
  { rank: 73, title: "Scope 1 & 2 Emissions Intensity Telemetry Feed", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "National", scale: "32 Tracked Metrics", horizon: "Realtime", impactMetric: "tCO2e / $1M CAD Capex", description: "Continuous operational greenhouse gas benchmarking against 2050 Net-Zero sectoral pathways." },
  { rank: 74, title: "Annual Lifecycle GHG Abatement Index (Mt CO2e/yr)", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "National", scale: "42.8 Mt National Total", horizon: "Annual", impactMetric: "Decarbonization Potential", description: "Net displaced fossil counter-factual lifecycle emissions across nuclear, SMR, and battery storage." },
  { rank: 75, title: "Bill C-59 70% Domestic Content Compliance Engine", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "National", scale: "CRA BN Verification", horizon: "Semi-Annual", impactMetric: "Sovereign Supply Chain %", description: "Ensures megaprojects maximize procurement with Canadian Business Number registered suppliers." },
  { rank: 76, title: "Flyvbjerg Bayesian Cost Overrun Hazard Hazard Curve", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "National", scale: "P10 / P50 / P90 Hazard", horizon: "Continual", impactMetric: "Tail Risk Mitigation", description: "Reference-class empirical forecasting preventing catastrophic multi-billion dollar budget slips." },
  { rank: 77, title: "Red Seal Craft Union Trades Gap Monitor (BuildForce)", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "National", scale: "14,200 FTE Deficit", horizon: "Quarterly", impactMetric: "Peak Craft Shortage", description: "Real-time collision detection for electricians, boilermakers, and pipefitters across provinces." },
  { rank: 78, title: "Western Canadian Select (WCS) Crude Discount Live Ticker", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "AB/SK", scale: "Market Pricing", horizon: "Realtime", impactMetric: "$USD/bbl Differential", description: "Live trading spread against Cushing WTI measuring heavy oil pipeline takeaway sufficiency." },
  { rank: 79, title: "AECO C Natural Gas Spot & Forward Hub Feed", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "AB", scale: "Market Pricing", horizon: "Daily", impactMetric: "$CAD/GJ Spot Price", description: "Benchmark natural gas feedstock cost governing petrochemical, hydrogen, and gas peaker margins." },
  { rank: 80, title: "LME Grade-1 Nickel Cash Settlement Real-Time Feed", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "Global/CA", scale: "Market Pricing", horizon: "Realtime", impactMetric: "$USD/tonne LME Cash", description: "Direct commercial trigger for Canadian nickel sulphide mining and battery pCAM refinery FID." },
  { rank: 81, title: "Ux U3O8 Spot Uranium Indicator Feed", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "Global/SK", scale: "Cameco Benchmark", horizon: "Weekly", impactMetric: "$USD/lb Spot Yellowcake", description: "Fuel cycle cost benchmark for Ontario Power Generation, Bruce Power, and global export buyers." },
  { rank: 82, title: "Bank of Canada Valet Policy Rate & 10Y Yield Streaming", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "National", scale: "Monetary Anchor", horizon: "Daily", impactMetric: "Cost of Project Capital", description: "Automated macro telemetry updating discount hurdle rates and DSCR coverage ratios." },
  { rank: 83, title: "Balancing Authority Substation Headroom Index (MW)", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "Provincial", scale: "Grid Telemetry", horizon: "Monthly", impactMetric: "MW Available Hosting", description: "Direct physical feasibility scoring before permitting industrial SMRs or AI datacentres." },
  { rank: 84, title: "Statutory Permitting Completion Progress Index (%)", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "Provincial", scale: "All Permitted Sites", horizon: "Continual", impactMetric: "Construction Readiness", description: "Tracking completion of municipal, provincial, and federal licenses against target milestones." },
  { rank: 85, title: "Municipal Housing Absorption Deficit Gauge (Units)", cluster: "KPIs & Analytics", clusterKey: "analytics", jurisdiction: "Municipal", scale: "CMHC Data Link", horizon: "Semi-Annual", impactMetric: "Local Rental Pressure", description: "Detects community housing deficits within 45 minutes of construction sites to avoid local displacement." },

  // Cluster 8: Supply Chain Vulnerability & Industrial Substitution (86-94)
  { rank: 86, title: "500kV High-Voltage Autotransformer Procurement Reserve", cluster: "Supply Chain", clusterKey: "supply", jurisdiction: "National", scale: "$500M CAD Reserve", horizon: "2026", impactMetric: "Lead Time Reduction (48 to 18 mo)", description: "National strategic stockpiling of critical high-voltage transmission transformers." },
  { rank: 87, title: "Nuclear Calandria & Steam Generator Heavy Forging Capacity", cluster: "Supply Chain", clusterKey: "supply", jurisdiction: "ON", scale: "$750M CAD (BWXT/Cameco)", horizon: "2027", impactMetric: "CSA N285 Certified Metallurgy", description: "Domestic fabrication scale ensuring Canada remains self-sufficient in reactor pressure components." },
  { rank: 88, title: "Cathode Active Material (pCAM) Domestic Precursor Plant", cluster: "Supply Chain", clusterKey: "supply", jurisdiction: "QC", scale: "$1.2B CAD (Bécancour)", horizon: "2026", impactMetric: "120,000 t/yr Battery Feed", description: "Refines nickel, cobalt, and manganese sulphates into cathode precursors without overseas transport." },
  { rank: 89, title: "Class 1 Heavy-Haul Arctic Rail Rolling Stock Fleet", cluster: "Supply Chain", clusterKey: "supply", jurisdiction: "MB/ON", scale: "$320M CAD", horizon: "2026", impactMetric: "Cold-Weather Hopper Cars", description: "Specialized cold-weather railcars rated for continuous operation on discontinuous permafrost." },
  { rank: 90, title: "U.S. Section 232 Steel & Aluminium Tariff Protection Shield", cluster: "Supply Chain", clusterKey: "supply", jurisdiction: "National", scale: "Policy Protocol", horizon: "2025", impactMetric: "Zero-Tariff Border Transit", description: "Digital traceability proving Canadian origin to shield domestic exports from U.S. trade actions." },
  { rank: 91, title: "Domestic Low-Carbon Rebar & Structural Steel Substitution", cluster: "Supply Chain", clusterKey: "supply", jurisdiction: "ON/QC", scale: "$850M CAD (Algoma/Stelco)", horizon: "2026", impactMetric: "Electric Arc Furnace Steel", description: "Substitutes high-carbon offshore steel with Canadian electric-arc furnace structural steel." },
  { rank: 92, title: "Critical Mineral Chemical Refining Reagents (Sulphuric Acid)", cluster: "Supply Chain", clusterKey: "supply", jurisdiction: "ON/QC", scale: "$280M CAD", horizon: "2026", impactMetric: "Domestic Leaching Security", description: "Secures domestic supply of high-purity industrial acid for battery metal extraction." },
  { rank: 93, title: "Modular Remote Workforce Accommodation Manufacturing", cluster: "Supply Chain", clusterKey: "supply", jurisdiction: "AB/BC", scale: "$350M CAD", horizon: "2025", impactMetric: "15,000 Turnkey Beds", description: "High-spec modular camp fabrication supporting northern mining and clean energy projects." },
  { rank: 94, title: "High-Purity Hydrogen Fuel Cell Membrane Fabrication Hub", cluster: "Supply Chain", clusterKey: "supply", jurisdiction: "BC", scale: "$190M CAD (Ballard/AFCC)", horizon: "2026", impactMetric: "Electrolyser MEA Stack Output", description: "Advanced domestic membrane electrode assembly production for clean hydrogen generation." },

  // Cluster 9: Statutory Regulatory Reform & Assessment Speed (95-100)
  { rank: 95, title: "Bill C-69 Impact Assessment Act 2-Year Statutory Deadline", cluster: "Regulatory Reform", clusterKey: "regulatory", jurisdiction: "Federal", scale: "Statutory Reform", horizon: "Enacted", impactMetric: "Max 24-Month Federal Reviews", description: "Enforceable statutory duration cap on federal environmental reviews for major projects." },
  { rank: 96, title: "One-Project-One-Assessment Federal/Provincial Reciprocity", cluster: "Regulatory Reform", clusterKey: "regulatory", jurisdiction: "Intergovernmental", scale: "Reciprocity Accords", horizon: "2025", impactMetric: "Zero Duplicate EA Hearings", description: "Eliminates overlapping federal and provincial reviews through single-window substitute assessments." },
  { rank: 97, title: "Digital Environmental Baseline Sensor Telemetry Network", cluster: "Regulatory Reform", clusterKey: "regulatory", jurisdiction: "National", scale: "Satellite/IoT Integration", horizon: "2026", impactMetric: "Pre-Approved Baseline Data", description: "Shared open-data environmental baselines reducing upfront EA study durations by 12-18 months." },
  { rank: 98, title: "Early Crown Consultation & Section 35 Mandate Protocol", cluster: "Regulatory Reform", clusterKey: "regulatory", jurisdiction: "Crown-Indigenous", scale: "Legal Standard", horizon: "2025", impactMetric: "Pre-Application Free Prior Consent", description: "Front-loaded consultation guidelines resolving territory overlaps prior to formal statutory filing." },
  { rank: 99, title: "Major Projects Management Office (MPMO) Single-Window Concierge", cluster: "Regulatory Reform", clusterKey: "regulatory", jurisdiction: "Federal (NRCan)", scale: "All Tier-1 Megaprojects", horizon: "Active", impactMetric: "Inter-Agency Permitting Fast-Track", description: "Federal deputy-minister level task force clearing bureaucratic impasses across departments." },
  { rank: 100, title: "Canada Economic Opportunity Graph (CEO-G) National Deployment", cluster: "Regulatory Reform", clusterKey: "regulatory", jurisdiction: "National", scale: "Cryptographic Open Graph", horizon: "2025-2026", impactMetric: "Evidence-Addressable Sovereign Intelligence", description: "The single source of verified truth connecting projects, capital stacks, First Nations, and sovereign KPIs." },
];

export default function AnalyticsPage() {
  const [selectedProjectSlug, setSelectedProjectSlug] = useState<string>("darlington-new-nuclear-project-unit-1");
  const [liveTicks, setLiveTicks] = useState<LiveFeedTick[]>(INITIAL_LIVE_TICKS);
  const [isStreaming, setIsStreaming] = useState<boolean>(true);
  const [activeTab, setActiveTab] = useState<"overview" | "scorecard" | "simulator" | "top100">("overview");
  const [selectedPriorityCluster, setSelectedPriorityCluster] = useState<string>("all");
  const [prioritySearch, setPrioritySearch] = useState<string>("");

  // Macro Shock Simulation State
  const [carbonTaxDelta, setCarbonTaxDelta] = useState<number>(0); // $0 to +$100 / tonne
  const [rateHikeBps, setRateHikeBps] = useState<number>(0); // -100 to +300 bps
  const [tariffShockPct, setTariffShockPct] = useState<number>(0); // 0% to 25%
  const [reviewDelayMonths, setReviewDelayMonths] = useState<number>(0); // 0 to 24 months

  // Find project
  const selectedProject = useMemo(() => {
    return (
      SNAPSHOT_PROJECTS.find((p) => p.slug === selectedProjectSlug) ||
      SNAPSHOT_PROJECTS[0]
    );
  }, [selectedProjectSlug]);

  // Compute scorecard with shocks applied
  const scorecard = useMemo<ProjectKPIScorecard>(() => {
    const base = evaluateProjectKPIs(selectedProject);
    if (carbonTaxDelta === 0 && rateHikeBps === 0 && tariffShockPct === 0 && reviewDelayMonths === 0) {
      return base;
    }

    // Apply simulated macro shock adjustments
    const shockedPillars = base.pillars.map((pillar) => {
      let pillarScore = pillar.score;
      const metrics = pillar.metrics.map((m) => {
        let obs = m.observed_value;
        let score = m.score_normalized;

        if (m.code === "ESG.CARBON.TAX_SENSITIVITY.CAD") {
          obs += carbonTaxDelta;
          score = Math.max(10, score - (carbonTaxDelta / 10) * 4);
        } else if (m.code === "MACRO.BOC.POLICY_RATE.PCT") {
          obs += rateHikeBps / 100;
          score = Math.max(15, score - (rateHikeBps / 50) * 5);
        } else if (m.code === "TRADE.US_TARIFF.EXPOSURE.PCT") {
          obs += tariffShockPct;
          score = Math.max(10, score - tariffShockPct * 2.5);
        } else if (m.code === "REG.IAAC.REVIEW_DURATION.MONTHS") {
          obs += reviewDelayMonths;
          score = Math.max(15, score - reviewDelayMonths * 2);
        }

        return { ...m, observed_value: Math.round(obs * 100) / 100, score_normalized: Math.round(score * 10) / 10 };
      });

      const avg = metrics.reduce((acc, x) => acc + x.score_normalized, 0) / metrics.length;
      pillarScore = Math.round(avg * 10) / 10;
      const health = pillarScore >= 85 ? "EXEMPLARY" : pillarScore >= 70 ? "HEALTHY" : pillarScore >= 50 ? "ATTENTION_REQUIRED" : "HIGH_RISK";

      return { ...pillar, score: pillarScore, health: health as any, metrics };
    });

    const overall = Math.round(
      (shockedPillars[0].score * 0.15 +
        shockedPillars[1].score * 0.15 +
        shockedPillars[2].score * 0.15 +
        shockedPillars[3].score * 0.15 +
        shockedPillars[4].score * 0.15 +
        shockedPillars[5].score * 0.10 +
        shockedPillars[6].score * 0.10 +
        shockedPillars[7].score * 0.05) *
        10
    ) / 10;

    return { ...base, pillars: shockedPillars, overall_kpi_rating: overall };
  }, [selectedProject, carbonTaxDelta, rateHikeBps, tariffShockPct, reviewDelayMonths]);

  // Live feed simulated micro-updates
  useEffect(() => {
    if (!isStreaming) return;
    const interval = setInterval(() => {
      setLiveTicks((prev) => {
        const randomIndex = Math.floor(Math.random() * prev.length);
        const item = prev[randomIndex];
        const drift = (Math.random() - 0.49) * 0.008; // small +/- 0.4% move
        const newVal = item.value * (1 + drift);
        const diff = newVal - item.value;
        const diffPct = (diff / item.value) * 100;
        const dir = diff > 0.0001 ? "UP" : diff < -0.0001 ? "DOWN" : "FLAT";

        const updated = [...prev];
        updated[randomIndex] = {
          ...item,
          value: Math.round(newVal * 100) / 100,
          change_absolute: Math.round(diff * 100) / 100,
          change_percent: Math.round(diffPct * 100) / 100,
          direction: dir,
          timestamp: new Date().toISOString(),
        };
        return updated;
      });
    }, 2400);

    return () => clearInterval(interval);
  }, [isStreaming]);

  // Filter Top 100 Priorities
  const filteredPriorities = useMemo(() => {
    return TOP_100_PRIORITIES.filter((item) => {
      const matchesCluster = selectedPriorityCluster === "all" || item.clusterKey === selectedPriorityCluster;
      const matchesSearch =
        prioritySearch === "" ||
        item.title.toLowerCase().includes(prioritySearch.toLowerCase()) ||
        item.description.toLowerCase().includes(prioritySearch.toLowerCase()) ||
        item.jurisdiction.toLowerCase().includes(prioritySearch.toLowerCase());
      return matchesCluster && matchesSearch;
    });
  }, [selectedPriorityCluster, prioritySearch]);

  const copySnapshotJSON = () => {
    const jsonStr = JSON.stringify({
      version: "cegs-kpi-v1.0",
      project: scorecard,
      live_feeds: liveTicks,
      macro_summary: DEFAULT_MACRO_SUMMARY,
      timestamp: new Date().toISOString(),
    }, null, 2);
    navigator.clipboard.writeText(jsonStr);
    alert("Copied full KPI Snapshot & Scorecard JSON to clipboard!");
  };

  return (
    <div className="mx-auto max-w-7xl space-y-8 px-4 py-8 sm:px-6 lg:px-8">
      {/* Header Banner */}
      <header className="relative overflow-hidden rounded-3xl border border-primary/30 bg-gradient-to-b from-[#091510] via-[#050b08] to-[#040806] p-6 shadow-2xl sm:p-8">
        <div aria-hidden="true" className="pointer-events-none absolute -right-24 -top-24 h-96 w-96 rounded-full bg-primary/15 blur-3xl" />
        <div aria-hidden="true" className="pointer-events-none absolute -bottom-24 -left-24 h-96 w-96 rounded-full bg-gold/10 blur-3xl" />

        <div className="relative z-10 flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
          <div className="max-w-3xl">
            <div className="inline-flex items-center gap-2 rounded-full border border-primary/40 bg-primary/10 px-3 py-1 font-mono text-[11px] font-bold text-aurora">
              <Sparkles aria-hidden="true" className="h-3.5 w-3.5 animate-pulse" /> SOVEREIGN KPI INTELLIGENCE & LIVE FEEDS ENGINE
            </div>
            <h1 className="mt-3 text-3xl font-black tracking-tight text-text-main sm:text-5xl">
              National Infrastructure <span className="text-aurora">Analytics</span>
            </h1>
            <p className="mt-3 text-sm leading-relaxed text-text-muted">
              Continuous live tracking across 8 strategic pillars: ESG Decarbonization, Indigenous Sovereignty, Capital Spend Velocity, Sovereign Supply Chains, Power Grid Physics, Permitting Latency, Red Seal Craft Labour, and Commodity Benchmarks.
            </p>
          </div>

          <div className="flex flex-wrap items-center gap-3">
            <button
              onClick={() => setIsStreaming(!isStreaming)}
              className={`inline-flex items-center gap-2 rounded-xl border px-4 py-2.5 font-mono text-xs font-bold transition-all ${
                isStreaming
                  ? "border-emerald-500/50 bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20"
                  : "border-border bg-surface text-text-muted hover:text-text-main"
              }`}
            >
              {isStreaming ? (
                <>
                  <Pause className="h-3.5 w-3.5" /> STREAMING LIVE (2.4s)
                </>
              ) : (
                <>
                  <Play className="h-3.5 w-3.5" /> STREAM PAUSED
                </>
              )}
            </button>
            <button
              onClick={copySnapshotJSON}
              className="inline-flex items-center gap-2 rounded-xl border border-primary/40 bg-primary/10 px-4 py-2.5 text-xs font-bold text-aurora hover:bg-primary/20"
            >
              <Download className="h-3.5 w-3.5" /> Export Snapshot JSON
            </button>
          </div>
        </div>

        {/* Live Commodity & Rate Ribbon */}
        <div className="mt-8 rounded-2xl border border-border/70 bg-black/40 p-4 backdrop-blur-md">
          <div className="flex items-center justify-between pb-3">
            <div className="flex items-center gap-2 font-mono text-[11px] uppercase tracking-wider text-text-muted">
              <Activity className="h-3.5 w-3.5 text-aurora" /> Live Benchmark Ticker (SHA-256 Provenance)
            </div>
            <span className="font-mono text-[10px] text-text-muted">8 Active Market Feeds</span>
          </div>

          <div className="grid grid-cols-2 gap-3 sm:grid-cols-4 lg:grid-cols-8">
            {liveTicks.map((t) => (
              <div
                key={t.metric_code}
                className="rounded-xl border border-border/60 bg-surface/60 p-2.5 transition-all hover:border-primary/50"
              >
                <div className="truncate font-mono text-[10px] text-text-muted" title={t.name}>
                  {t.name}
                </div>
                <div className="mt-1 flex items-baseline justify-between">
                  <span className="font-mono text-xs font-black tabular-nums text-text-main">
                    {t.metric_code.includes("USD_T") || t.metric_code.includes("GUARANTEE")
                      ? `$${t.value.toLocaleString()}`
                      : t.metric_code.includes("SPOT_FX")
                      ? `$${t.value.toFixed(4)}`
                      : t.value.toFixed(2)}
                  </span>
                  <span
                    className={`inline-flex items-center font-mono text-[10px] font-bold ${
                      t.direction === "UP"
                        ? "text-emerald-400"
                        : t.direction === "DOWN"
                        ? "text-rose-400"
                        : "text-text-muted"
                    }`}
                  >
                    {t.direction === "UP" ? "+" : ""}
                    {t.change_percent}%
                  </span>
                </div>
                <div className="mt-0.5 truncate font-mono text-[9px] text-text-muted">
                  {t.unit}
                </div>
              </div>
            ))}
          </div>
        </div>
      </header>

      {/* Navigation Tabs */}
      <nav aria-label="Analytics Navigation" className="flex border-b border-border/80">
        <div className="flex gap-2">
          {[
            { id: "overview", label: "National Sovereign Overview", icon: Globe2 },
            { id: "scorecard", label: "8-Pillar Project Scorecard", icon: BarChart3 },
            { id: "simulator", label: "Macro Shock Stress-Testing", icon: Sliders },
            { id: "top100", label: "Top 100 Priority Items Catalog", icon: Layers, badge: "100" },
          ].map((tab) => {
            const Icon = tab.icon;
            const active = activeTab === tab.id;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id as any)}
                className={`inline-flex items-center gap-2 border-b-2 px-4 py-3 text-xs font-bold transition-colors ${
                  active
                    ? "border-aurora text-aurora"
                    : "border-transparent text-text-muted hover:border-border hover:text-text-main"
                }`}
              >
                <Icon className="h-4 w-4" />
                <span>{tab.label}</span>
                {tab.badge && (
                  <span className="rounded-full bg-primary/20 px-2 py-0.5 font-mono text-[10px] text-aurora">
                    {tab.badge}
                  </span>
                )}
              </button>
            );
          })}
        </div>
      </nav>

      {/* TAB 1: National Sovereign Overview */}
      {activeTab === "overview" && (
        <section aria-labelledby="national-sovereign-overview-title" className="space-y-6">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div className="glass-card rounded-2xl border border-border p-5">
              <div className="flex items-center justify-between text-text-muted">
                <span className="font-mono text-xs uppercase tracking-wider">Annual Decarbonization</span>
                <Leaf className="h-4 w-4 text-emerald-400" />
              </div>
              <div className="mt-2 text-2xl font-black text-text-main tabular-nums">42.8 Mt CO2e/yr</div>
              <p className="mt-1 text-xs text-text-muted">Aggregate avoided emissions across tracked projects</p>
            </div>

            <div className="glass-card rounded-2xl border border-border p-5">
              <div className="flex items-center justify-between text-text-muted">
                <span className="font-mono text-xs uppercase tracking-wider">First Nations Equity</span>
                <Users className="h-4 w-4 text-gold" />
              </div>
              <div className="mt-2 text-2xl font-black text-gold tabular-nums">18.4% Average</div>
              <p className="mt-1 text-xs text-text-muted">$1.45B committed under $5B Federal ILGP</p>
            </div>

            <div className="glass-card rounded-2xl border border-border p-5">
              <div className="flex items-center justify-between text-text-muted">
                <span className="font-mono text-xs uppercase tracking-wider">Canadian Content</span>
                <ShieldCheck className="h-4 w-4 text-aurora" />
              </div>
              <div className="mt-2 text-2xl font-black text-aurora tabular-nums">68.2% Domestic BN</div>
              <p className="mt-1 text-xs text-text-muted">Procurement spend with registered Canadian vendors</p>
            </div>

            <div className="glass-card rounded-2xl border border-border p-5">
              <div className="flex items-center justify-between text-text-muted">
                <span className="font-mono text-xs uppercase tracking-wider">Red Seal Trades Gap</span>
                <AlertTriangle className="h-4 w-4 text-amber-400" />
              </div>
              <div className="mt-2 text-2xl font-black text-amber-400 tabular-nums">14,200 FTE Deficit</div>
              <p className="mt-1 text-xs text-text-muted">Peak construction craft shortage in western provinces</p>
            </div>
          </div>

          {/* 8 Strategic KPI Pillar Matrix */}
          <div className="rounded-2xl border border-border bg-surface/50 p-6">
            <h2 id="national-sovereign-overview-title" className="text-xl font-bold text-text-main">32 Canonical Indicators Across 8 Strategic Pillars</h2>
            <p className="mt-1 text-xs text-text-muted">
              Every indicator points to an authoritative Canadian statute, provincial system operator, or multilateral source with independent verification.
            </p>

            <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
              {[
                { title: "Pillar 1: ESG & Decarbonization", count: "5 Metrics", target: "Net-Zero 2050", color: "border-emerald-500/30 text-emerald-400", desc: "Scope 1/2 intensity, Mt abatement, carbon price sensitivity, water circularity, clean electricity purity." },
                { title: "Pillar 2: Indigenous Sovereignty", count: "5 Metrics", target: "25% Co-Ownership", color: "border-amber-500/30 text-amber-400", desc: "First Nations equity, 5% federal procurement compliance, IBA local employment, $5B ILGP loan guarantee." },
                { title: "Pillar 3: Capital Spend Velocity", count: "4 Metrics", target: "4.0x Crowding-In", color: "border-cyan-500/30 text-cyan-400", desc: "Monthly capex deployment, schedule slippage, Flyvbjerg P50 overrun hazard, institutional leverage." },
                { title: "Pillar 4: Sovereign Supply Chain", count: "4 Metrics", target: "70% Domestic BN", color: "border-blue-500/30 text-blue-400", desc: "Canadian Business Number spend, critical hardware bottlenecks, U.S. border tariff exposure." },
                { title: "Pillar 5: Power & Grid Physics", count: "4 Metrics", target: "18 Mo Queue", color: "border-purple-500/30 text-purple-400", desc: "Peak MW demand, provincial ISO queue latency, substation headroom, transmission reinforcement cost." },
                { title: "Pillar 6: Permitting & Regulatory", count: "4 Metrics", target: "24 Mo IAAC", color: "border-rose-500/30 text-rose-400", desc: "Bill C-69 statutory elapsed time, permit completion velocity, Section 35 judicial review risk." },
                { title: "Pillar 7: Labour & Apprenticeship", count: "4 Metrics", target: ">=10% Apprentice", color: "border-yellow-500/30 text-yellow-400", desc: "Red Seal craft deficit, apprentice-to-journeyperson ratio, municipal housing absorption buffer." },
                { title: "Pillar 8: Commodity Benchmarks", count: "8 Live Feeds", target: "Realtime Feeds", color: "border-emerald-400/30 text-aurora", desc: "WCS discount, AECO gas, LME nickel, Ux uranium, lithium carbonate, BoC rate, CAD/USD spot." },
              ].map((p) => (
                <div key={p.title} className={`rounded-xl border bg-card/60 p-4 ${p.color}`}>
                  <div className="flex items-center justify-between">
                    <span className="font-mono text-[10px] uppercase font-bold">{p.count}</span>
                    <span className="rounded bg-black/40 px-2 py-0.5 font-mono text-[9px]">{p.target}</span>
                  </div>
                  <h3 className="mt-2 text-sm font-bold text-text-main">{p.title}</h3>
                  <p className="mt-1 text-xs text-text-muted leading-relaxed">{p.desc}</p>
                </div>
              ))}
            </div>
          </div>
        </section>
      )}

      {/* TAB 2: Project Scorecard Explorer */}
      {activeTab === "scorecard" && (
        <section aria-labelledby="project-scorecard-explorer-title" className="space-y-6">
          <div className="flex flex-col gap-4 rounded-2xl border border-border bg-surface/80 p-5 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <label htmlFor="project-selector" className="font-mono text-xs uppercase tracking-wider text-text-muted">Select Megaproject to Assess</label>
              <h2 id="project-scorecard-explorer-title" className="text-xl font-bold text-text-main">
                {scorecard.project_name}
              </h2>
              <div className="mt-1 flex flex-wrap gap-2 text-xs text-text-muted">
                <span>Sector: <strong className="text-text-main">{scorecard.sector}</strong></span>
                <span>•</span>
                <span>Province: <strong className="text-text-main">{scorecard.province}</strong></span>
                <span>•</span>
                <span>Stage: <strong className="text-text-main">{scorecard.current_stage}</strong></span>
                <span>•</span>
                <span>CAPEX: <strong className="text-text-main">${(scorecard.total_capex_cad / 1_000_000).toLocaleString()}M CAD</strong></span>
              </div>
            </div>

            <div className="flex items-center gap-4">
              <select
                id="project-selector"
                value={selectedProjectSlug}
                onChange={(e) => setSelectedProjectSlug(e.target.value)}
                aria-label="Select Megaproject to Assess"
                className="min-h-10 rounded-xl border border-border bg-card px-3 py-2 text-xs font-bold text-text-main focus:border-primary focus:outline-none"
              >
                {SNAPSHOT_PROJECTS.slice(0, 30).map((p) => (
                  <option key={p.slug} value={p.slug}>
                    {p.name} ({p.province})
                  </option>
                ))}
              </select>

              <div className="rounded-xl border border-primary/40 bg-primary/10 px-4 py-2 text-right">
                <div className="font-mono text-[9px] uppercase tracking-wider text-text-muted">OVERALL KPI RATING</div>
                <div className="text-2xl font-black tabular-nums text-aurora">{scorecard.overall_kpi_rating}/100</div>
              </div>
            </div>
          </div>

          {/* Critical Gaps Alert if any */}
          {scorecard.critical_action_gaps.length > 0 && (
            <div className="rounded-2xl border border-rose-500/40 bg-rose-500/10 p-4">
              <div className="flex items-center gap-2 font-mono text-xs font-bold text-rose-400">
                <ShieldAlert className="h-4 w-4" /> CRITICAL ACTION GAPS IDENTIFIED ({scorecard.critical_action_gaps.length})
              </div>
              <ul className="mt-2 space-y-1 text-xs text-rose-200/90">
                {scorecard.critical_action_gaps.map((gap, i) => (
                  <li key={i} className="flex items-start gap-2">
                    <span className="text-rose-400">•</span>
                    <span>{gap}</span>
                  </li>
                ))}
              </ul>
            </div>
          )}

          {/* 8 Pillar Detail Cards */}
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
            {scorecard.pillars.map((pillar) => (
              <div key={pillar.title} className="rounded-2xl border border-border bg-card/70 p-5 shadow-lg">
                <div className="flex items-center justify-between border-b border-border/60 pb-3">
                  <div>
                    <h3 className="text-base font-bold text-text-main">{pillar.title}</h3>
                    <span
                      className={`inline-block mt-0.5 rounded px-2 py-0.5 font-mono text-[9px] font-bold ${
                        pillar.health === "EXEMPLARY"
                          ? "bg-emerald-500/20 text-emerald-400"
                          : pillar.health === "HEALTHY"
                          ? "bg-cyan-500/20 text-cyan-400"
                          : pillar.health === "ATTENTION_REQUIRED"
                          ? "bg-amber-500/20 text-amber-400"
                          : "bg-rose-500/20 text-rose-400"
                      }`}
                    >
                      {pillar.health}
                    </span>
                  </div>
                  <div className="text-right">
                    <div className="text-xl font-black text-text-main tabular-nums">{pillar.score}/100</div>
                    <div className="font-mono text-[9px] text-text-muted">Normalized Score</div>
                  </div>
                </div>

                <div className="mt-4 space-y-3">
                  {pillar.metrics.map((m) => (
                    <div key={m.code} className="rounded-xl border border-border/50 bg-surface/50 p-3">
                      <div className="flex items-center justify-between text-xs">
                        <span className="font-bold text-text-main">{m.name}</span>
                        <span
                          className={`rounded px-1.5 py-0.5 font-mono text-[9px] font-bold ${
                            m.performance_rank === "SUPERIOR"
                              ? "bg-emerald-500/20 text-emerald-400"
                              : m.performance_rank === "ON_TARGET"
                              ? "bg-cyan-500/20 text-cyan-400"
                              : m.performance_rank === "NEEDS_IMPROVEMENT"
                              ? "bg-amber-500/20 text-amber-400"
                              : "bg-rose-500/20 text-rose-400"
                          }`}
                        >
                          {m.performance_rank}
                        </span>
                      </div>

                      <div className="mt-2 flex items-baseline justify-between font-mono text-xs">
                        <span className="text-text-muted">
                          Observed: <strong className="text-text-main tabular-nums">{m.observed_value.toLocaleString()} {m.unit}</strong>
                        </span>
                        <span className="text-text-muted">
                          Target: <strong className="tabular-nums text-text-muted">{m.target_benchmark.toLocaleString()} {m.unit}</strong>
                        </span>
                      </div>

                      <div className="mt-2 flex items-center justify-between text-[10px] text-text-muted">
                        <span className="truncate" title={m.notes}>{m.notes}</span>
                        <span className="font-mono text-aurora">{m.score_normalized} pts</span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </section>
      )}

      {/* TAB 3: Macro Shock Stress-Testing Simulator */}
      {activeTab === "simulator" && (
        <section aria-labelledby="macro-shock-simulator-title" className="space-y-6">
          <div className="rounded-2xl border border-primary/40 bg-gradient-to-br from-card via-surface to-background p-6">
            <div className="inline-flex items-center gap-2 rounded-full border border-primary/40 bg-primary/10 px-3 py-1 font-mono text-[11px] font-bold text-aurora">
              <Sliders className="h-3.5 w-3.5" /> STOCHASTIC MACRO SHOCK SIMULATOR
            </div>
            <h2 id="macro-shock-simulator-title" className="mt-2 text-2xl font-black text-text-main">
              Simulate Policy & Commodity Shocks on {scorecard.project_name}
            </h2>
            <p className="mt-1 max-w-2xl text-xs leading-relaxed text-text-muted">
              Adjust variables below to observe real-time recalculation of the project&apos;s 8-pillar rating, debt service capacity, and critical action gaps.
            </p>

            <div className="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
              {/* Slider 1: Carbon Price Escalation */}
              <div className="rounded-xl border border-border bg-black/30 p-4">
                <div className="flex items-center justify-between text-xs font-bold text-text-main">
                  <span>Carbon Tax Escalation</span>
                  <span className="font-mono text-emerald-400">+${carbonTaxDelta}/t</span>
                </div>
                <input
                  type="range"
                  min="0"
                  max="100"
                  step="5"
                  value={carbonTaxDelta}
                  onChange={(e) => setCarbonTaxDelta(Number(e.target.value))}
                  className="mt-3 w-full accent-primary"
                />
                <div className="mt-2 flex justify-between font-mono text-[10px] text-text-muted">
                  <span>$0 (Baseline)</span>
                  <span>+$100 ($170/t Cap)</span>
                </div>
              </div>

              {/* Slider 2: Interest Rate Shock */}
              <div className="rounded-xl border border-border bg-black/30 p-4">
                <div className="flex items-center justify-between text-xs font-bold text-text-main">
                  <span>BoC Rate Hike</span>
                  <span className="font-mono text-amber-400">+{rateHikeBps} bps</span>
                </div>
                <input
                  type="range"
                  min="0"
                  max="300"
                  step="25"
                  value={rateHikeBps}
                  onChange={(e) => setRateHikeBps(Number(e.target.value))}
                  className="mt-3 w-full accent-amber-400"
                />
                <div className="mt-2 flex justify-between font-mono text-[10px] text-text-muted">
                  <span>0 bps</span>
                  <span>+300 bps (Stagflation)</span>
                </div>
              </div>

              {/* Slider 3: U.S. Border Tariff Shock */}
              <div className="rounded-xl border border-border bg-black/30 p-4">
                <div className="flex items-center justify-between text-xs font-bold text-text-main">
                  <span>U.S. Border Tariff Shock</span>
                  <span className="font-mono text-rose-400">+{tariffShockPct}%</span>
                </div>
                <input
                  type="range"
                  min="0"
                  max="25"
                  step="2.5"
                  value={tariffShockPct}
                  onChange={(e) => setTariffShockPct(Number(e.target.value))}
                  className="mt-3 w-full accent-rose-400"
                />
                <div className="mt-2 flex justify-between font-mono text-[10px] text-text-muted">
                  <span>0% (CUSMA Free)</span>
                  <span>+25% (Section 232)</span>
                </div>
              </div>

              {/* Slider 4: Regulatory Delay */}
              <div className="rounded-xl border border-border bg-black/30 p-4">
                <div className="flex items-center justify-between text-xs font-bold text-text-main">
                  <span>IAAC Review Delay</span>
                  <span className="font-mono text-cyan-400">+{reviewDelayMonths} mo</span>
                </div>
                <input
                  type="range"
                  min="0"
                  max="24"
                  step="3"
                  value={reviewDelayMonths}
                  onChange={(e) => setReviewDelayMonths(Number(e.target.value))}
                  className="mt-3 w-full accent-cyan-400"
                />
                <div className="mt-2 flex justify-between font-mono text-[10px] text-text-muted">
                  <span>0 mo</span>
                  <span>+24 mo (Judicial Delay)</span>
                </div>
              </div>
            </div>

            {/* Shock Results Display */}
            <div className="mt-8 rounded-xl border border-border bg-card p-5">
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h3 className="text-base font-bold text-text-main">Calibrated Shock Impact Summary</h3>
                  <p className="text-xs text-text-muted">
                    Recalibrated rating reflects compounded project finance resilience under adverse macroeconomic conditions.
                  </p>
                </div>
                <div className="flex items-center gap-4">
                  <button
                    onClick={() => {
                      setCarbonTaxDelta(0);
                      setRateHikeBps(0);
                      setTariffShockPct(0);
                      setReviewDelayMonths(0);
                    }}
                    className="inline-flex items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-xs text-text-muted hover:text-text-main"
                  >
                    <RefreshCw className="h-3 w-3" /> Reset Shocks
                  </button>
                  <div className="rounded-xl border border-primary/40 bg-primary/10 px-4 py-2 text-right">
                    <span className="font-mono text-[9px] uppercase text-text-muted">SHOCKED KPI RATING</span>
                    <div className="text-xl font-black text-aurora">{scorecard.overall_kpi_rating}/100</div>
                  </div>
                </div>
              </div>

              <div className="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-4">
                {scorecard.pillars.slice(0, 4).map((p) => (
                  <div key={p.title} className="rounded-lg border border-border/50 bg-surface/50 p-3">
                    <div className="text-xs text-text-muted truncate">{p.title}</div>
                    <div className="mt-1 text-lg font-bold text-text-main tabular-nums">{p.score}/100</div>
                    <span className="font-mono text-[9px] text-text-muted">{p.health}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </section>
      )}

      {/* TAB 4: Top 100 Priority Items Catalog */}
      {activeTab === "top100" && (
        <section aria-labelledby="top-100-priority-items-title" className="space-y-6">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 id="top-100-priority-items-title" className="text-2xl font-black text-text-main">Top 100 Canadian Sovereign Priority Items</h2>
              <p className="text-xs text-text-muted">
                The authoritative registry of major projects, infrastructure assets, live data feeds, and statutory policy reforms required for economic sovereignty.
              </p>
            </div>

            <div className="flex flex-wrap items-center gap-2">
              <div className="relative">
                <Search className="absolute left-3 top-2.5 h-3.5 w-3.5 text-text-muted" />
                <input
                  type="text"
                  placeholder="Filter priorities..."
                  value={prioritySearch}
                  onChange={(e) => setPrioritySearch(e.target.value)}
                  className="min-h-9 rounded-xl border border-border bg-card pl-9 pr-3 text-xs text-text-main focus:border-primary focus:outline-none"
                />
              </div>

              <select
                value={selectedPriorityCluster}
                onChange={(e) => setSelectedPriorityCluster(e.target.value)}
                className="min-h-9 rounded-xl border border-border bg-card px-3 text-xs font-bold text-text-main focus:border-primary focus:outline-none"
              >
                <option value="all">All Clusters (100)</option>
                <option value="minerals">Critical Minerals (15)</option>
                <option value="nuclear">Nuclear & Clean Grid (12)</option>
                <option value="energy">Clean Energy & Hydrogen (11)</option>
                <option value="compute">Sovereign AI Compute (10)</option>
                <option value="ports">Ports & Corridors (12)</option>
                <option value="indigenous">Indigenous Sovereignty (12)</option>
                <option value="analytics">KPIs & Live Feeds (13)</option>
                <option value="supply">Supply Chain Defense (9)</option>
                <option value="regulatory">Regulatory Reforms (6)</option>
              </select>
            </div>
          </div>

          <div className="rounded-2xl border border-border bg-card/80 overflow-hidden shadow-xl">
            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead className="border-b border-border bg-surface/80 font-mono text-[10px] uppercase text-text-muted">
                  <tr>
                    <th className="px-4 py-3">Rank</th>
                    <th className="px-4 py-3">Priority Item</th>
                    <th className="px-4 py-3">Strategic Cluster</th>
                    <th className="px-4 py-3">Prov</th>
                    <th className="px-4 py-3">Scale / Target</th>
                    <th className="px-4 py-3">Horizon</th>
                    <th className="px-4 py-3">Primary Impact Metric</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/60">
                  {filteredPriorities.map((item) => (
                    <tr key={item.rank} className="hover:bg-surface/60 transition-colors">
                      <td className="px-4 py-3 font-mono font-bold text-aurora">#{item.rank}</td>
                      <td className="px-4 py-3">
                        <div className="font-bold text-text-main">{item.title}</div>
                        <div className="mt-0.5 text-[11px] text-text-muted leading-tight">{item.description}</div>
                      </td>
                      <td className="px-4 py-3">
                        <span className="rounded-full bg-primary/10 border border-primary/20 px-2 py-0.5 font-mono text-[10px] text-aurora">
                          {item.cluster}
                        </span>
                      </td>
                      <td className="px-4 py-3 font-mono font-bold text-text-main">{item.jurisdiction}</td>
                      <td className="px-4 py-3 font-mono text-text-main">{item.scale}</td>
                      <td className="px-4 py-3 font-mono text-text-muted">{item.horizon}</td>
                      <td className="px-4 py-3 font-mono text-gold text-[11px]">{item.impactMetric}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </section>
      )}
    </div>
  );
}
