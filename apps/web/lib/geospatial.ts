/**
 * Comprehensive Geospatial & Global Strategic Topography Data Module
 * Canada Economic Opportunity Graph (CEGS 0.1 / 1.0)
 * 
 * Defines global trade routes, geopolitical friction/conflict markers,
 * sovereign opportunity zones, transmission/logistics corridors,
 * and high-resolution satellite imagery base layers.
 */

export interface LatLngPoint {
  lat: number;
  lng: number;
  name?: string;
}

export interface TradeRoute {
  id: string;
  name: string;
  category: "TRANS_PACIFIC" | "TRANS_ATLANTIC" | "ARCTIC_POLAR" | "CONTINENTAL_USMCA" | "MARITIME_SEAWAY";
  origin: string;
  destination: string;
  primaryCommodities: string[];
  annualVelocity: string;
  strategicSignificance: string;
  color: string;
  dashArray?: string;
  waypoints: LatLngPoint[];
}

export interface ConflictMarker {
  id: string;
  name: string;
  type: "REGULATORY_STANDOFF" | "CHOKE_POINT" | "JURISDICTIONAL_OVERLAP" | "GEOECONOMIC_SANCTION";
  severity: "CRITICAL" | "HIGH" | "ELEVATED";
  coordinates: LatLngPoint;
  jurisdiction: string;
  impactSummary: string;
  affectedCorridors: string[];
  status: string;
}

export interface OpportunityZone {
  id: string;
  name: string;
  sector: string;
  center: LatLngPoint;
  boundsRadiusKm: number;
  estimatedEndowmentCAD: string;
  criticalMineralsOrEnergy: string[];
  keyProjects: string[];
  description: string;
  color: string;
}

export interface MapTileProvider {
  id: string;
  name: string;
  category: "SATELLITE" | "TACTICAL_DARK" | "TOPOGRAPHIC" | "GOOGLE";
  url: string;
  attribution: string;
  maxZoom: number;
  subdomains?: string[];
  badge: string;
}

export const MAP_PROVIDERS: MapTileProvider[] = [
  {
    id: "esri-satellite",
    name: "Esri World Imagery (High-Res Satellite)",
    category: "SATELLITE",
    url: "https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}",
    attribution: "Tiles &copy; Esri &mdash; Source: Esri, i-cubed, USDA, USGS, AEX, GeoEye, Getmapping, Aerogrid, IGN, IGP, UPR-EGP, GIS Community",
    maxZoom: 19,
    badge: "SUB-METER OPTICAL",
  },
  {
    id: "esri-clarity",
    name: "Esri Clarity (High-Definition Cloud-Free)",
    category: "SATELLITE",
    url: "https://clarity.maptiles.arcgis.com/arcgis/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}",
    attribution: "Tiles &copy; Esri Clarity &mdash; Archival Cloud-Free High-Definition",
    maxZoom: 19,
    badge: "CLOUD-FREE ARCHIVE",
  },
  {
    id: "carto-dark",
    name: "CartoDB Dark Matter (Tactical Operations)",
    category: "TACTICAL_DARK",
    url: "https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png",
    attribution: "&copy; <a href=\"https://www.openstreetmap.org/copyright\">OpenStreetMap</a> &copy; <a href=\"https://carto.com/\">CARTO</a>",
    maxZoom: 19,
    subdomains: ["a", "b", "c", "d"],
    badge: "NIGHT OPS",
  },
  {
    id: "osm-topo",
    name: "OpenStreetMap Infrastructure & Topo",
    category: "TOPOGRAPHIC",
    url: "https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png",
    attribution: "&copy; <a href=\"https://www.openstreetmap.org/copyright\">OpenStreetMap</a> contributors",
    maxZoom: 18,
    subdomains: ["a", "b", "c"],
    badge: "ROAD & RAIL GRID",
  },
  {
    id: "google-hybrid",
    name: "Google Maps Hybrid / Satellite (Direct)",
    category: "GOOGLE",
    url: "https://mt1.google.com/vt/lyrs=y&x={x}&y={y}&z={z}",
    attribution: "&copy; Google Maps &mdash; Satellite & Road Labels",
    maxZoom: 20,
    badge: "GOOGLE HYBRID",
  },
];

export const GLOBAL_TRADE_ROUTES: TradeRoute[] = [
  {
    id: "route-trans-pacific-lng-minerals",
    name: "Trans-Pacific Strategic Minerals & LNG Corridor",
    category: "TRANS_PACIFIC",
    origin: "Port of Prince Rupert & Vancouver (BC)",
    destination: "Yokohama / Tokyo / Busan / Kaohsiung",
    primaryCommodities: ["Liquefied Natural Gas (LNG)", "Metallurgical Coal", "Copper Concentrates", "Potash", "Rare Earth Elements"],
    annualVelocity: "$68.4B CAD Trade Velocity",
    strategicSignificance: "Primary sovereign maritime lifeline linking Western Canadian resource basins to Indo-Pacific industrial allies.",
    color: "#00F5A0",
    dashArray: "6, 4",
    waypoints: [
      { lat: 54.315, lng: -130.320, name: "Port of Prince Rupert" },
      { lat: 52.0, lng: -140.0 },
      { lat: 48.0, lng: -160.0 },
      { lat: 45.0, lng: -180.0, name: "International Date Line Crossing" },
      { lat: 42.0, lng: 165.0 },
      { lat: 38.0, lng: 150.0 },
      { lat: 35.443, lng: 139.638, name: "Port of Yokohama (Japan)" },
      { lat: 35.102, lng: 129.040, name: "Port of Busan (South Korea)" },
    ],
  },
  {
    id: "route-trans-atlantic-uranium-metals",
    name: "Trans-Atlantic Nuclear Fuel & Critical Metals Artery",
    category: "TRANS_ATLANTIC",
    origin: "Port of Montreal & Port of Halifax",
    destination: "Rotterdam / Antwerp / Hamburg / Liverpool",
    primaryCommodities: ["Refined Uranium (U3O8)", "Battery-Grade Nickel", "Biochar & Green Ammonia", "Titanium & Direct-Reduction Iron"],
    annualVelocity: "$52.1B CAD Trade Velocity",
    strategicSignificance: "Supplies NATO nuclear energy fleets and EU Net-Zero Industry Act manufacturing clusters.",
    color: "#38BDF8",
    dashArray: "6, 4",
    waypoints: [
      { lat: 45.508, lng: -73.554, name: "Port of Montreal" },
      { lat: 44.648, lng: -63.575, name: "Port of Halifax" },
      { lat: 46.5, lng: -50.0, name: "Grand Banks Atlantic Gateway" },
      { lat: 50.0, lng: -30.0 },
      { lat: 51.5, lng: -10.0, name: "Celtic Sea Approach" },
      { lat: 51.924, lng: 4.477, name: "Port of Rotterdam (Netherlands)" },
      { lat: 51.219, lng: 4.402, name: "Port of Antwerp (Belgium)" },
    ],
  },
  {
    id: "route-arctic-northwest-passage",
    name: "Northwest Passage Polar Sovereignty Vector",
    category: "ARCTIC_POLAR",
    origin: "Churchill / Tuktoyaktuk / Lancaster Sound",
    destination: "Bering Strait & North Atlantic Arctic Hubs",
    primaryCommodities: ["Arctic Defence Supply", "Northern High-Grade Iron", "Sovereign Maritime Surveillance", "Grain via Hudson Bay"],
    annualVelocity: "$8.7B CAD Projected 2030 Capacity",
    strategicSignificance: "Shortens Euro-Asian maritime transit by 7,000 km; vital for Canadian sovereign jurisdiction across high northern waters.",
    color: "#F59E0B",
    dashArray: "8, 5",
    waypoints: [
      { lat: 58.768, lng: -94.165, name: "Port of Churchill (Manitoba)" },
      { lat: 64.0, lng: -83.0, name: "Foxe Basin Approach" },
      { lat: 74.2, lng: -80.5, name: "Lancaster Sound (Eastern Gateway)" },
      { lat: 74.5, lng: -95.0, name: "Resolute / Barrow Strait" },
      { lat: 71.0, lng: -125.0, name: "Amundsen Gulf" },
      { lat: 69.444, lng: -133.032, name: "Tuktoyaktuk Arctic Harbour" },
      { lat: 71.0, lng: -156.0, name: "Point Barrow Arctic Vector" },
      { lat: 65.7, lng: -168.9, name: "Bering Strait Sovereign Boundary" },
    ],
  },
  {
    id: "route-usmca-continental-rail-spine",
    name: "USMCA Continental Heavy Industrial Rail Spine",
    category: "CONTINENTAL_USMCA",
    origin: "Calgary / Edmonton Industrial Heartland",
    destination: "Chicago Hub -> Houston Gulf -> Monterrey / Lázaro Cárdenas (Mexico)",
    primaryCommodities: ["Synthetic Crude", "Automotive Parts & EV Sub-assemblies", "Grain & Agricultural Potash", "Structural Steel"],
    annualVelocity: "$142.6B CAD Bilateral Velocity",
    strategicSignificance: "First single-line rail network connecting Canada, the United States, and Mexico (CPKC & CN networks).",
    color: "#A78BFA",
    dashArray: "5, 3",
    waypoints: [
      { lat: 53.546, lng: -113.493, name: "Edmonton Energy Junction" },
      { lat: 51.044, lng: -114.071, name: "Calgary Rail HQ" },
      { lat: 49.000, lng: -104.000, name: "Portal Border Crossing (US Boundary)" },
      { lat: 44.977, lng: -93.265, name: "Twin Cities Hub" },
      { lat: 41.878, lng: -87.629, name: "Chicago Rail Crossroads" },
      { lat: 39.099, lng: -94.578, name: "Kansas City Interchange" },
      { lat: 29.760, lng: -95.369, name: "Houston Chemical & Port Terminal" },
      { lat: 27.503, lng: -99.507, name: "Laredo / Nuevo Laredo Border" },
      { lat: 25.686, lng: -100.316, name: "Monterrey Industrial Complex (Mexico)" },
    ],
  },
  {
    id: "route-st-lawrence-seaway-spine",
    name: "Great Lakes - St. Lawrence Industrial Marine Corridor",
    category: "MARITIME_SEAWAY",
    origin: "Thunder Bay (Lake Superior)",
    destination: "St. Lawrence Gulf & Atlantic High Seas",
    primaryCommodities: ["Iron Ore Pellets", "Grain & Agri-Food", "Heavy Machinery", "Bauxite & Aluminium Ingots"],
    annualVelocity: "$35.2B CAD Maritime Commerce",
    strategicSignificance: "Deep inland waterway servicing 25% of North America's total industrial and manufacturing capacity.",
    color: "#34D399",
    dashArray: "4, 4",
    waypoints: [
      { lat: 48.380, lng: -89.247, name: "Thunder Bay Terminal" },
      { lat: 46.513, lng: -84.346, name: "Sault Ste. Marie Locks" },
      { lat: 42.314, lng: -83.036, name: "Detroit-Windsor Channel" },
      { lat: 43.255, lng: -79.871, name: "Hamilton Steel Port" },
      { lat: 45.508, lng: -73.554, name: "Montreal Harbour" },
      { lat: 46.813, lng: -71.207, name: "Quebec City Deepwater Port" },
      { lat: 48.500, lng: -64.000, name: "Gulf of St. Lawrence Outflow" },
    ],
  },
];

export const CONFLICT_MARKERS: ConflictMarker[] = [
  {
    id: "conflict-ring-of-fire-regulatory",
    name: "Ring of Fire Jurisdictional & Infrastructure Standoff",
    type: "REGULATORY_STANDOFF",
    severity: "CRITICAL",
    coordinates: { lat: 52.825, lng: -86.155 },
    jurisdiction: "Ontario / Federal (IAAC & Treaty 9 First Nations)",
    impactSummary: "All-weather road access and regional assessment delays hold back $60B+ in critical chromite, nickel, and platinum group metals.",
    affectedCorridors: ["Northern Ontario Mining Trunk", "Critical Mineral - EV Battery Highway"],
    status: "Active Joint Assessment & Consent Protocols",
  },
  {
    id: "conflict-salish-sea-tanker-moratorium",
    name: "Salish Sea Marine Choke Point & Escort Zone",
    type: "CHOKE_POINT",
    severity: "HIGH",
    coordinates: { lat: 48.785, lng: -123.320 },
    jurisdiction: "British Columbia / Federal Transport Canada / US Coast Guard",
    impactSummary: "Tug escort mandates, marine mammal acoustic speed caps, and seasonal whale protection corridors compress Pacific crude export cadence.",
    affectedCorridors: ["Trans-Mountain Pacific Export Corridor", "Trans-Pacific LNG & Minerals Corridor"],
    status: "Active Vessel Traffic Management System (VTMS)",
  },
  {
    id: "conflict-beaufort-sea-moratorium",
    name: "Beaufort Sea Arctic Energy Moratorium Boundary",
    type: "JURISDICTIONAL_OVERLAP",
    severity: "ELEVATED",
    coordinates: { lat: 70.450, lng: -135.000 },
    jurisdiction: "Federal Crown / Inuvialuit Settlement Region / US Boundary Sector",
    impactSummary: "Ongoing federal moratoria on offshore Arctic oil & gas extraction juxtaposed against defense radar expansion and undersea subsea fiber.",
    affectedCorridors: ["Northwest Passage Polar Sovereignty Vector"],
    status: "5-Year Federal Science Review Pending",
  },
  {
    id: "conflict-suez-redsea-diverter",
    name: "Bab-el-Mandeb & Red Sea Maritime Reroute Diverter",
    type: "CHOKE_POINT",
    severity: "CRITICAL",
    coordinates: { lat: 12.585, lng: 43.330 },
    jurisdiction: "International Waters / Red Sea Transit Corridor",
    impactSummary: "Houthi drone/missile attacks divert bulk carriers around Cape of Good Hope (+14 days), sharply escalating Canadian potash and grain demand in Atlantic basins.",
    affectedCorridors: ["Trans-Atlantic Nuclear Fuel Artery"],
    status: "Operation Prosperity Guardian Ongoing",
  },
  {
    id: "conflict-panama-canal-draft-restriction",
    name: "Panama Canal Freshwater Drought Choke Point",
    type: "CHOKE_POINT",
    severity: "HIGH",
    coordinates: { lat: 9.101, lng: -79.695 },
    jurisdiction: "Panama Canal Authority (ACP)",
    impactSummary: "Gatun Lake freshwater drought caps daily vessel transits by 30%, shifting Canadian prairie grains and Alberta sulfur to West Coast rail trunks.",
    affectedCorridors: ["USMCA Continental Heavy Industrial Rail Spine", "Trans-Pacific Strategic Minerals Corridor"],
    status: "Variable Draft Surcharge & Slot Auctions",
  },
  {
    id: "conflict-iaac-provincial-carveout",
    name: "Impact Assessment Act Supreme Court Reference Zone",
    type: "JURISDICTIONAL_OVERLAP",
    severity: "ELEVATED",
    coordinates: { lat: 53.500, lng: -115.000 },
    jurisdiction: "Alberta / Saskatchewan / Federal Attorney General",
    impactSummary: "Bill C-69 constitutional ruling restricts federal intervention to clear interprovincial effects, necessitating newly streamlined federal-provincial permitting accords.",
    affectedCorridors: ["USMCA Continental Rail Spine"],
    status: "Revised Federal IAAC Regulations Enacted",
  },
];

export const OPPORTUNITY_ZONES: OpportunityZone[] = [
  {
    id: "opp-ring-of-fire",
    name: "Ring of Fire Multi-Metal Mineral District",
    sector: "Critical Minerals & Battery Feedstocks",
    center: { lat: 52.825, lng: -86.155 },
    boundsRadiusKm: 120,
    estimatedEndowmentCAD: "$67.0B CAD",
    criticalMineralsOrEnergy: ["Chromite", "Nickel", "Copper", "Platinum", "Palladium"],
    keyProjects: ["proj-crawford-nickel", "Eagle's Nest", "Black Thor Chromite"],
    description: "North America's premier undeveloped chromite and magmatic nickel-copper deposit cluster, crucial for stainless steel and EV battery cathodes.",
    color: "#F59E0B",
  },
  {
    id: "opp-athabasca-basin",
    name: "Athabasca Basin High-Grade Uranium Basin",
    sector: "Nuclear Fuel & Clean Baseload",
    center: { lat: 57.500, lng: -105.500 },
    boundsRadiusKm: 160,
    estimatedEndowmentCAD: "$84.5B CAD",
    criticalMineralsOrEnergy: ["Uranium (U3O8)", "Cobalt", "Vanadium"],
    keyProjects: ["proj-bruce-nuclear", "McArthur River", "Cigar Lake", "Arrow / NexGen"],
    description: "Contains ore grades 10x to 100x above global average, powering over 20% of North America's clean nuclear generation and Western sovereign supply.",
    color: "#00F5A0",
  },
  {
    id: "opp-james-bay-lithium",
    name: "Eeyou Istchee James Bay Lithium Hub",
    sector: "Battery Materials & Chemical Refining",
    center: { lat: 52.000, lng: -76.000 },
    boundsRadiusKm: 140,
    estimatedEndowmentCAD: "$32.4B CAD",
    criticalMineralsOrEnergy: ["Spodumene Lithium", "Beryllium", "Clean Hydro Power"],
    keyProjects: ["proj-chisasibi-ai-compute", "Whabouchi Lithium", "Rose Lithium-Tantalum", "Corvette / Patriot"],
    description: "Low-carbon lithium hard-rock corridor powered by 100% renewable Hydro-Québec transmission, feeding Bécancour battery cathode precursor plants.",
    color: "#38BDF8",
  },
  {
    id: "opp-montney-duvernay",
    name: "Montney-Duvernay Clean Transition & CCS Play",
    sector: "Energy Transition, Geothermal & Carbon Capture",
    center: { lat: 55.000, lng: -119.500 },
    boundsRadiusKm: 180,
    estimatedEndowmentCAD: "$110.0B CAD",
    criticalMineralsOrEnergy: ["Sub-surface Lithium Brines", "Hydrogen Feedstock", "Gigaton Deep CCS Reservoirs"],
    keyProjects: ["Quest CCS", "Pathways Alliance Trunk", "Clearwater Oil Sands"],
    description: "World's most concentrated deep saline aquifer geological formation for permanent carbon dioxide sequestration and lithium extraction from brine.",
    color: "#A78BFA",
  },
  {
    id: "opp-southern-ontario-ev",
    name: "Southern Ontario EV & Battery Gigafactory Belt",
    sector: "Advanced Manufacturing & Mobility",
    center: { lat: 42.750, lng: -81.250 },
    boundsRadiusKm: 110,
    estimatedEndowmentCAD: "$48.0B CAD",
    criticalMineralsOrEnergy: ["Cathode Active Material", "Cell Manufacturing", "Grid Energy Storage"],
    keyProjects: ["proj-oneida-battery", "proj-darlington-smr", "NextStar Windsor", "PowerCo St. Thomas"],
    description: "Anchor of Canada's automotive industrial renewal with over $40B in committed capital for battery cell gigafactories, cathode facilities, and grid storage.",
    color: "#EC4899",
  },
  {
    id: "opp-labrador-trough",
    name: "Labrador Trough Direct-Reduction Iron Ore Belt",
    sector: "Green Steel & Direct Reduction Pellets",
    center: { lat: 54.800, lng: -66.800 },
    boundsRadiusKm: 150,
    estimatedEndowmentCAD: "$55.0B CAD",
    criticalMineralsOrEnergy: ["Super-Concentrated Direct-Reduction Iron Ore", "Clean Churchill Falls Hydro"],
    keyProjects: ["Bloom Lake Iron", "Carol Lake Pellet Plant", "Mont-Wright"],
    description: "Supplies premium low-silica, high-purity direct-reduction (DR) grade iron pellets indispensable for green hydrogen-fuelled electric arc steel furnaces.",
    color: "#F97316",
  },
  {
    id: "opp-alberta-heartland",
    name: "Alberta Industrial Heartland & Hydrogen Valley",
    sector: "Clean Hydrogen, Ammonia & Petrochemicals",
    center: { lat: 53.750, lng: -113.200 },
    boundsRadiusKm: 90,
    estimatedEndowmentCAD: "$41.5B CAD",
    criticalMineralsOrEnergy: ["Net-Zero Blue Hydrogen", "Clean Ammonia", "Circular Petrochemicals"],
    keyProjects: ["Air Products Net-Zero Hydrogen", "Dow Fort Saskatchewan Path2Zero"],
    description: "Canada's largest continuous chemical manufacturing cluster with shared open-access carbon transport pipelines and deep geological storage.",
    color: "#10B981",
  },
];

export const STRATEGIC_CORRIDORS = [
  {
    id: "corridor-baie-james",
    name: "Baie-James Clean Hydro Corridor",
    from: "proj-chisasibi-ai-compute",
    to: "proj-hyperscale-qc",
    fromCoords: { lat: 53.783, lng: -78.916 },
    toCoords: { lat: 45.508, lng: -73.554 },
    capacity: "5,000 MW Transmission",
    color: "#00F5A0",
  },
  {
    id: "corridor-hudson-bay",
    name: "Hudson Bay Strategic Northern Vector",
    from: "proj-churchill-arctic-gateway",
    to: "proj-kivalliq-link",
    fromCoords: { lat: 58.768, lng: -94.165 },
    toCoords: { lat: 62.816, lng: -92.083 },
    capacity: "Subsea Hydro & Fiber Cable",
    color: "#F59E0B",
  },
  {
    id: "corridor-nuclear-spine",
    name: "Ontario Clean Nuclear Baseload Spine",
    from: "proj-bruce-nuclear",
    to: "proj-darlington-smr",
    fromCoords: { lat: 44.325, lng: -81.598 },
    toCoords: { lat: 43.869, lng: -78.718 },
    capacity: "10,000 MW Nuclear Grid",
    color: "#38BDF8",
  },
  {
    id: "corridor-critical-mineral-battery",
    name: "Critical Mineral - EV Battery Highway",
    from: "proj-crawford-nickel",
    to: "proj-oneida-battery",
    fromCoords: { lat: 48.783, lng: -81.333 },
    toCoords: { lat: 42.923, lng: -80.015 },
    capacity: "500-kV Heavy Industrial Feed",
    color: "#EC4899",
  },
];

export const MAP_FOCUS_PRESETS = [
  {
    id: "preset-canada-all",
    name: "🇨🇦 All Canada",
    center: [58.0, -98.0] as [number, number],
    zoom: 4,
    description: "Pan-Canadian sovereign economic geography",
  },
  {
    id: "preset-global-trade",
    name: "🌏 Global Trade Routes",
    center: [38.0, -35.0] as [number, number],
    zoom: 2,
    description: "Pacific, Atlantic, and Arctic international corridors",
  },
  {
    id: "preset-arctic-sovereignty",
    name: "❄️ Arctic Sovereignty",
    center: [71.0, -95.0] as [number, number],
    zoom: 4,
    description: "Northwest Passage, Churchill, and northern security",
  },
  {
    id: "preset-critical-minerals",
    name: "💎 Critical Mineral Belts",
    center: [53.5, -88.0] as [number, number],
    zoom: 5,
    description: "Ring of Fire, Athabasca, and James Bay clusters",
  },
  {
    id: "preset-pacific-gateway",
    name: "🚢 Pacific Gateway",
    center: [49.5, -135.0] as [number, number],
    zoom: 4,
    description: "Prince Rupert & Vancouver maritime access to Asia",
  },
  {
    id: "preset-usmca-rail",
    name: "🚂 USMCA Rail Spine",
    center: [38.0, -98.0] as [number, number],
    zoom: 4,
    description: "Continental cross-border industrial integration",
  },
];

// --- Pillar A: Critical Minerals Refining & Midstream Hubs ---
export interface CriticalMineralHub {
  id: string;
  name: string;
  location: LatLngPoint;
  province: string;
  focusMinerals: string[];
  processingType: "HYDROMETALLURGICAL" | "SMELTING" | "CAM_GIGAFACTORY" | "RECYCLING";
  domesticRetentionRate: number;
  capacityMetric: string;
  annualValueCAD: string;
  color: string;
}

export const CRITICAL_MINERAL_HUBS: CriticalMineralHub[] = [
  {
    id: "hub-becancour",
    name: "Bécancour Battery Valley & CAM Megahub",
    location: { lat: 46.342, lng: -72.435 },
    province: "QC",
    focusMinerals: ["Lithium", "Nickel", "Cobalt", "Graphite"],
    processingType: "CAM_GIGAFACTORY",
    domesticRetentionRate: 0.94,
    capacityMetric: "120,000 tonnes/yr CAM",
    annualValueCAD: "$14.5B CAD",
    color: "#00F5A0",
  },
  {
    id: "hub-temiskaming",
    name: "Temiskaming Shores Cobalt & Nickel Refinery",
    location: { lat: 47.516, lng: -79.678 },
    province: "ON",
    focusMinerals: ["Cobalt", "Nickel", "Black Mass Recycling"],
    processingType: "HYDROMETALLURGICAL",
    domesticRetentionRate: 0.88,
    capacityMetric: "6,500 tonnes/yr Cobalt Sulfate",
    annualValueCAD: "$1.8B CAD",
    color: "#38BDF8",
  },
  {
    id: "hub-sudbury",
    name: "Sudbury Integrated Nickel Smelting Complex",
    location: { lat: 46.491, lng: -80.993 },
    province: "ON",
    focusMinerals: ["Class-1 Nickel", "Copper", "PGMs"],
    processingType: "SMELTING",
    domesticRetentionRate: 0.82,
    capacityMetric: "145,000 tonnes/yr Refined Nickel",
    annualValueCAD: "$7.2B CAD",
    color: "#F59E0B",
  },
  {
    id: "hub-sorel-tracy",
    name: "Sorel-Tracy Critical & Rare Metals Facility",
    location: { lat: 46.033, lng: -73.116 },
    province: "QC",
    focusMinerals: ["Titanium", "Scandium", "Rare Earth Elements"],
    processingType: "HYDROMETALLURGICAL",
    domesticRetentionRate: 0.91,
    capacityMetric: "3 tonnes/yr High-Purity Scandium Oxide",
    annualValueCAD: "$950M CAD",
    color: "#EC4899",
  },
];

// --- Pillar B: Historic & Modern Treaty Territories ---
export interface TreatyTerritory {
  id: string;
  name: string;
  center: LatLngPoint;
  province: string;
  historicalFramework: "HISTORIC_NUMBERED_TREATY" | "MODERN_COMPREHENSIVE_TREATY" | "PEACE_AND_FRIENDSHIP";
  signatories: string[];
  loanGuaranteeEligibility: string;
  color: string;
}

export const TREATY_TERRITORIES: TreatyTerritory[] = [
  {
    id: "treaty-james-bay",
    name: "Grand Council of the Crees (Eeyou Istchee / James Bay)",
    center: { lat: 52.5, lng: -77.5 },
    province: "QC",
    historicalFramework: "MODERN_COMPREHENSIVE_TREATY",
    signatories: ["Cree Nation", "Government of Canada", "Government of Quebec"],
    loanGuaranteeEligibility: "100% Priority Access under $5B Federal ILGP",
    color: "#F59E0B",
  },
  {
    id: "treaty-8-alberta-bc",
    name: "Treaty 8 Traditional Territory (Energy & Mineral Basin)",
    center: { lat: 56.5, lng: -117.5 },
    province: "AB / BC",
    historicalFramework: "HISTORIC_NUMBERED_TREATY",
    signatories: ["Treaty 8 First Nations", "Crown in Right of Canada"],
    loanGuaranteeEligibility: "Full Priority Access — AIOC & Federal ILGP",
    color: "#10B981",
  },
  {
    id: "treaty-3-ontario",
    name: "Grand Council Treaty #3 (Northwestern Ontario Infrastructure)",
    center: { lat: 49.8, lng: -93.5 },
    province: "ON",
    historicalFramework: "HISTORIC_NUMBERED_TREATY",
    signatories: ["Treaty 3 Anishinaabe Nations", "Crown"],
    loanGuaranteeEligibility: "Eligible for ALGP & Federal ILGP",
    color: "#8B5CF6",
  },
  {
    id: "treaty-nisgaa",
    name: "Nisga'a Nation Territory (Nass River / Coastal Gateway)",
    center: { lat: 55.2, lng: -129.2 },
    province: "BC",
    historicalFramework: "MODERN_COMPREHENSIVE_TREATY",
    signatories: ["Nisga'a Lisims Government", "Canada", "British Columbia"],
    loanGuaranteeEligibility: "Full Sovereign Co-Financing Available",
    color: "#06B6D4",
  },
];

// --- Pillar D: Clean Energy Grid Interties & AI Compute Clusters ---
export interface GridIntertieZone {
  id: string;
  name: string;
  center: LatLngPoint;
  gridOperator: string;
  cleanCapacityMW: number;
  aiHeadroomMW: number;
  baseloadType: string;
  cleanFlopsRatio: string;
  color: string;
}

export const GRID_INTERTIE_ZONES: GridIntertieZone[] = [
  {
    id: "grid-darlington-nuclear",
    name: "Darlington Nuclear & SMR AI Clean Compute Hub",
    center: { lat: 43.869, lng: -78.718 },
    gridOperator: "IESO (Ontario)",
    cleanCapacityMW: 3500,
    aiHeadroomMW: 850,
    baseloadType: "CANDU Nuclear & BWRX-300 SMR",
    cleanFlopsRatio: "7.00 ExaFLOPs/GW",
    color: "#38BDF8",
  },
  {
    id: "grid-beauharnois-hydro",
    name: "Beauharnois Hydroelectric Hyperscale Megacampus",
    center: { lat: 45.316, lng: -73.905 },
    gridOperator: "Hydro-Québec",
    cleanCapacityMW: 1900,
    aiHeadroomMW: 600,
    baseloadType: "Run-of-the-River Hydro",
    cleanFlopsRatio: "6.94 ExaFLOPs/GW",
    color: "#00F5A0",
  },
  {
    id: "grid-site-c-peace",
    name: "Site C Clean Hydro Transmission Dispatch",
    center: { lat: 56.196, lng: -120.912 },
    gridOperator: "BC Hydro",
    cleanCapacityMW: 1100,
    aiHeadroomMW: 450,
    baseloadType: "Reservoir Hydro",
    cleanFlopsRatio: "6.95 ExaFLOPs/GW",
    color: "#10B981",
  },
  {
    id: "grid-alberta-heartland",
    name: "Alberta Industrial Heartland Cogeneration & CCS",
    center: { lat: 53.722, lng: -113.217 },
    gridOperator: "AESO (Alberta)",
    cleanCapacityMW: 900,
    aiHeadroomMW: 350,
    baseloadType: "Industrial Cogeneration + CCS",
    cleanFlopsRatio: "6.46 ExaFLOPs/GW",
    color: "#F59E0B",
  },
];

