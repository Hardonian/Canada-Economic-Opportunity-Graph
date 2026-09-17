"use client";

import { useState, useEffect, useMemo, useRef } from "react";
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
import { useLanguage } from "@/components/LanguageProvider";

// ----------------------------------------------------------------------
// Top 100 Canadian Sovereign Priority Items Catalog
// ----------------------------------------------------------------------
interface PriorityItem {
  rank: number;
  title: string;
  titleFr: string;
  cluster: string;
  clusterFr: string;
  clusterKey: string;
  jurisdiction: string;
  scale: string;
  scaleFr: string;
  horizon: string;
  impactMetric: string;
  impactMetricFr: string;
  description: string;
  descriptionFr: string;
}

const TOP_100_PRIORITIES: PriorityItem[] = [
  // Cluster 1: Critical Minerals & Upstream Extraction (1-15)
  {
    rank: 1,
    title: "Crawford Nickel Sulphide Project (Canada Nickel)",
    titleFr: "Projet de sulfure de nickel Crawford (Canada Nickel)",
    cluster: "Critical Minerals",
    clusterFr: "Minéraux critiques",
    clusterKey: "minerals",
    jurisdiction: "ON",
    scale: "$3.5B CAD",
    scaleFr: "3,5 G$ CA",
    horizon: "2027",
    impactMetric: "30,000 t/yr ESG Nickel",
    impactMetricFr: "30 000 t/an nickel ESG",
    description: "Largest nickel sulphide discovery globally in decades; zero-carbon tailings mineralization potential.",
    descriptionFr: "Plus importante découverte de sulfures de nickel au monde depuis des décennies; potentiel de résidus carboneutres.",
  },
  {
    rank: 2,
    title: "Ring of Fire Eagle's Nest Nickel-Copper-PGE (Noront / Wyloo)",
    titleFr: "Cercle de feu — Eagle's Nest Nickel-Cuivre-ÉGP (Noront / Wyloo)",
    cluster: "Critical Minerals",
    clusterFr: "Minéraux critiques",
    clusterKey: "minerals",
    jurisdiction: "ON",
    scale: "$2.2B CAD",
    scaleFr: "2,2 G$ CA",
    horizon: "2029",
    impactMetric: "15,000 t/yr Ni-Cu",
    impactMetricFr: "15 000 t/an Ni-Cu",
    description: "Anchor underground deposit unlocking the James Bay Lowlands critical mineral province.",
    descriptionFr: "Gisement souterrain d'ancrage ouvrant la province de minéraux critiques des basses-terres de la baie James.",
  },
  {
    rank: 3,
    title: "Prairie Lithium Brine DLE Commercial Facility",
    titleFr: "Installation commerciale d'extraction directe de lithium de saumure des Prairies",
    cluster: "Critical Minerals",
    clusterFr: "Minéraux critiques",
    clusterKey: "minerals",
    jurisdiction: "SK",
    scale: "$650M CAD",
    scaleFr: "650 M$ CA",
    horizon: "2026",
    impactMetric: "20,000 t/yr LCE",
    impactMetricFr: "20 000 t/an ÉCL",
    description: "Direct Lithium Extraction (DLE) from Devonian aquifers displacing carbon-heavy spodumene refining.",
    descriptionFr: "Extraction directe du lithium (EDL) des aquifères dévoniens remplaçant le raffinage polluant de spodumène.",
  },
  {
    rank: 4,
    title: "Galaxy Lithium / James Bay Spodumene (Arcadium)",
    titleFr: "Galaxy Lithium / Spodumène de la Baie-James (Arcadium)",
    cluster: "Critical Minerals",
    clusterFr: "Minéraux critiques",
    clusterKey: "minerals",
    jurisdiction: "QC",
    scale: "$850M CAD",
    scaleFr: "850 M$ CA",
    horizon: "2027",
    impactMetric: "330,000 t/yr Spodumene",
    impactMetricFr: "330 000 t/an spodumène",
    description: "High-grade open-pit hard-rock lithium feed for North American battery gigafactories.",
    descriptionFr: "Alimentation en lithium de roche dure à ciel ouvert de haute teneur pour les gigafactories de batteries nord-américaines.",
  },
  {
    rank: 5,
    title: "Matawinie Natural Graphite Mine & Bécancour Anode Plant",
    titleFr: "Mine de graphite naturel Matawinie et usine d'anodes de Bécancour",
    cluster: "Critical Minerals",
    clusterFr: "Minéraux critiques",
    clusterKey: "minerals",
    jurisdiction: "QC",
    scale: "$1.4B CAD",
    scaleFr: "1,4 G$ CA",
    horizon: "2026",
    impactMetric: "100,000 t/yr Anode",
    impactMetricFr: "100 000 t/an matériel d'anode",
    description: "Near-zero carbon spherical graphite supply breaking overseas anode dependency.",
    descriptionFr: "Approvisionnement en graphite sphérique quasi carboneutre brisant la dépendance étrangère aux anodes.",
  },
  {
    rank: 6,
    title: "Nechalacho Rare Earth Elements Phase 2 (Vital Metals / Cheetah)",
    titleFr: "Éléments des terres rares Nechalacho Phase 2 (Vital Metals / Cheetah)",
    cluster: "Critical Minerals",
    clusterFr: "Minéraux critiques",
    clusterKey: "minerals",
    jurisdiction: "NT",
    scale: "$350M CAD",
    scaleFr: "350 M$ CA",
    horizon: "2027",
    impactMetric: "5,000 t/yr NdPr",
    impactMetricFr: "5 000 t/an NdPr",
    description: "Strategic light and heavy rare earth elements essential for permanent magnet defense supply chains.",
    descriptionFr: "Terres rares légères et lourdes stratégiques essentielles aux aimants permanents pour la défense.",
  },
  {
    rank: 7,
    title: "Whabouchi Lithium Mine & Shawinigan Hydroxide Plant",
    titleFr: "Mine de lithium Whabouchi et usine d'hydroxyde de Shawinigan",
    cluster: "Critical Minerals",
    clusterFr: "Minéraux critiques",
    clusterKey: "minerals",
    jurisdiction: "QC",
    scale: "$1.8B CAD",
    scaleFr: "1,8 G$ CA",
    horizon: "2026",
    impactMetric: "34,000 t/yr LiOH",
    impactMetricFr: "34 000 t/an LiOH",
    description: "Integrated mine-to-chemical conversion facility powered 100% by Hydro-Québec clean energy.",
    descriptionFr: "Complexe intégré d'extraction et de conversion chimique alimenté à 100 % par l'hydroélectricité propre.",
  },
  {
    rank: 8,
    title: "McArthur River / Key Lake Uranium Expansion (Cameco)",
    titleFr: "Agrandissement d'uranium McArthur River / Key Lake (Cameco)",
    cluster: "Critical Minerals",
    clusterFr: "Minéraux critiques",
    clusterKey: "minerals",
    jurisdiction: "SK",
    scale: "$1.2B CAD",
    scaleFr: "1,2 G$ CA",
    horizon: "2026",
    impactMetric: "25M lbs/yr U3O8",
    impactMetricFr: "25 M lb/an U3O8",
    description: "World's highest-grade uranium complex securing western nuclear reactor fuel cycle autonomy.",
    descriptionFr: "Complexe uranifère à plus haute teneur au monde assurant l'autonomie du cycle du combustible occidental.",
  },
  {
    rank: 9,
    title: "Kipawa Heavy Rare Earths Complex",
    titleFr: "Complexe de terres rares lourdes de Kipawa",
    cluster: "Critical Minerals",
    clusterFr: "Minéraux critiques",
    clusterKey: "minerals",
    jurisdiction: "QC",
    scale: "$480M CAD",
    scaleFr: "480 M$ CA",
    horizon: "2028",
    impactMetric: "Dysprosium / Terbium",
    impactMetricFr: "Dysprosium / Terbium",
    description: "Crucial heavy rare earths for high-temperature military radar and guided missile actuators.",
    descriptionFr: "Terres rares lourdes cruciales pour les radars militaires haute température et les actionneurs de missiles.",
  },
  {
    rank: 10,
    title: "Valentine Gold Mine Construction (Calibre Mining)",
    titleFr: "Construction de la mine d'or Valentine (Calibre Mining)",
    cluster: "Critical Minerals",
    clusterFr: "Minéraux critiques",
    clusterKey: "minerals",
    jurisdiction: "NL",
    scale: "$750M CAD",
    scaleFr: "750 M$ CA",
    horizon: "2025",
    impactMetric: "195,000 oz/yr Au",
    impactMetricFr: "195 000 oz/an Au",
    description: "Atlantic Canada's largest gold development generating regional fiscal foundation.",
    descriptionFr: "Plus important projet aurifère du Canada atlantique générant des retombées fiscales régionales majeures.",
  },

  // Cluster 2: Nuclear Power, SMRs & Clean Baseload Grid (16-25)
  {
    rank: 16,
    title: "Darlington New Nuclear Project Unit 1 SMR (GE Hitachi BWRX-300)",
    titleFr: "Projet nouveau nucléaire de Darlington — Tranche 1 PRM (GE Hitachi BWRX-300)",
    cluster: "Nuclear & Clean Grid",
    clusterFr: "Réseau nucléaire et propre",
    clusterKey: "nuclear",
    jurisdiction: "ON",
    scale: "$3.4B CAD",
    scaleFr: "3,4 G$ CA",
    horizon: "2028",
    impactMetric: "300 MWe Clean Baseload",
    impactMetricFr: "300 MWe de charge de base propre",
    description: "G7's first commercial grid-scale Small Modular Reactor providing zero-carbon 24/7 firm power.",
    descriptionFr: "Premier petit réacteur modulaire commercial à l'échelle du réseau du G7 fournissant une énergie propre 24/7.",
  },
  {
    rank: 17,
    title: "Darlington SMR Fleet Expansion (Units 2, 3 & 4)",
    titleFr: "Déploiement de la flotte de PRM de Darlington (Tranches 2, 3 et 4)",
    cluster: "Nuclear & Clean Grid",
    clusterFr: "Réseau nucléaire et propre",
    clusterKey: "nuclear",
    jurisdiction: "ON",
    scale: "$10.2B CAD",
    scaleFr: "10,2 G$ CA",
    horizon: "2034",
    impactMetric: "1,200 MWe SMR Fleet",
    impactMetricFr: "Flotte de 1 200 MWe de PRM",
    description: "Multi-unit fleet deployment driving standardized EPCM execution and supply chain scale.",
    descriptionFr: "Déploiement d'une flotte multi-unités favorisant la standardisation et l'échelle de la chaîne logistique.",
  },
  {
    rank: 18,
    title: "Bruce C 4,800 MW Nuclear Generation Expansion",
    titleFr: "Agrandissement nucléaire Bruce C de 4 800 MW",
    cluster: "Nuclear & Clean Grid",
    clusterFr: "Réseau nucléaire et propre",
    clusterKey: "nuclear",
    jurisdiction: "ON",
    scale: "$16.0B CAD",
    scaleFr: "16,0 G$ CA",
    horizon: "2036",
    impactMetric: "4,800 MWe New Build",
    impactMetricFr: "4 800 MWe de nouvelle puissance",
    description: "Canada's largest nuclear expansion since the 1980s to power industrial AI and EV manufacturing.",
    descriptionFr: "Plus grande expansion nucléaire canadienne depuis les années 1980 pour alimenter l'IA et l'industrie.",
  },

  // Cluster 3: Clean Energy & Storage
  {
    rank: 28,
    title: "Oneida Energy Storage 250MW / 1,000MWh Battery Hub",
    titleFr: "Pôle de stockage d'énergie Oneida 250 MW / 1 000 MWh",
    cluster: "Clean Energy & Storage",
    clusterFr: "Énergie propre et stockage",
    clusterKey: "energy",
    jurisdiction: "ON",
    scale: "$600M CAD",
    scaleFr: "600 M$ CA",
    horizon: "2025",
    impactMetric: "1,000 MWh Grid Firming",
    impactMetricFr: "1 000 MWh de stabilisation réseau",
    description: "Canada's largest utility-scale battery storage facility, co-owned by Six Nations of the Grand River.",
    descriptionFr: "Plus grande installation de stockage par batteries du Canada, copropriété des Six Nations de la Grand River.",
  },

  // Cluster 4: Sovereign AI Compute & Datacentres
  {
    rank: 39,
    title: "Hydro-Québec Beauharnois Sovereign AI Supercompute Hub",
    titleFr: "Pôle de supercalcul IA souverain Beauharnois d'Hydro-Québec",
    cluster: "Sovereign AI Compute",
    clusterFr: "Calcul et IA souveraine",
    clusterKey: "compute",
    jurisdiction: "QC",
    scale: "$2.8B CAD",
    scaleFr: "2,8 G$ CA",
    horizon: "2026",
    impactMetric: "500 MW Hydro Clean Compute",
    impactMetricFr: "500 MW de calcul vert hydroélectrique",
    description: "Tier-4 hyperscale sovereign AI cluster guaranteeing data sovereignty and Bill C-27 compliance.",
    descriptionFr: "Grappe d'IA souveraine hyperscale de niveau 4 garantissant la souveraineté des données et la Loi C-27.",
  },

  // Cluster 5: Ports, Gateways & Strategic Corridors
  {
    rank: 49,
    title: "Port of Prince Rupert Fairview & Ridley Terminals Expansion",
    titleFr: "Agrandissement des terminaux Fairview et Ridley du port de Prince Rupert",
    cluster: "Ports & Gateways",
    clusterFr: "Ports et corridors",
    clusterKey: "ports",
    jurisdiction: "BC",
    scale: "$2.5B CAD",
    scaleFr: "2,5 G$ CA",
    horizon: "2027",
    impactMetric: "3.2M TEU & Energy Export",
    impactMetricFr: "3,2 M EVP et exportation d'énergie",
    description: "North America's closest port to Asia; strategic CN rail gateway for critical minerals and grain.",
    descriptionFr: "Port nord-américain le plus proche de l'Asie; porte ferroviaire stratégique pour les minéraux et les céréales.",
  },

  // Cluster 6: First Nations & Indigenous Co-Investment
  {
    rank: 61,
    title: "Canada Indigenous Loan Guarantee Program ($5B National Allocation)",
    titleFr: "Programme canadien de garantie de prêts aux Autochtones (enveloppe de 5 G$)",
    cluster: "Indigenous Sovereignty",
    clusterFr: "Souveraineté autochtone",
    clusterKey: "indigenous",
    jurisdiction: "National",
    scale: "$5.0B CAD",
    scaleFr: "5,0 G$ CA",
    horizon: "2025-2027",
    impactMetric: "First Nations Sovereign Equity",
    impactMetricFr: "Capitaux propres souverains des Premières Nations",
    description: "Concessionary federal debt guarantees enabling multi-nation ownership across energy and minerals.",
    descriptionFr: "Garanties d'emprunt fédérales concessionnelles permettant l'actionnariat autochtone dans l'énergie.",
  },

  // Cluster 7: Critical Sovereign KPIs & Live Feeds
  {
    rank: 73,
    title: "Scope 1 & 2 Emissions Intensity Telemetry Feed",
    titleFr: "Flux de télémétrie de l'intensité des émissions de portées 1 et 2",
    cluster: "KPIs & Analytics",
    clusterFr: "Indicateurs et flux en direct",
    clusterKey: "analytics",
    jurisdiction: "National",
    scale: "32 Tracked Metrics",
    scaleFr: "32 indicateurs suivis",
    horizon: "Realtime",
    impactMetric: "tCO2e / $1M CAD Capex",
    impactMetricFr: "tCO2e / 1 M$ CA de dépenses",
    description: "Continuous operational greenhouse gas benchmarking against 2050 Net-Zero sectoral pathways.",
    descriptionFr: "Étalonnage continu des gaz à effet de serre opérationnels par rapport aux trajectoires de carboneutralité 2050.",
  },
  {
    rank: 75,
    title: "Bill C-59 70% Domestic Content Compliance Engine",
    titleFr: "Moteur de conformité au contenu canadien de 70 % (Projet de loi C-59)",
    cluster: "KPIs & Analytics",
    clusterFr: "Indicateurs et flux en direct",
    clusterKey: "analytics",
    jurisdiction: "National",
    scale: "CRA BN Verification",
    scaleFr: "Vérification NE ARC",
    horizon: "Semi-Annual",
    impactMetric: "Sovereign Supply Chain %",
    impactMetricFr: "% chaîne logistique souveraine",
    description: "Ensures megaprojects maximize procurement with Canadian Business Number registered suppliers.",
    descriptionFr: "Assure que les mégaprojets maximisent les contrats avec les fournisseurs enregistrés au Canada.",
  },
  {
    rank: 76,
    title: "Flyvbjerg Bayesian Cost Overrun Hazard Curve",
    titleFr: "Courbe de risque bayésien de dépassement de coûts de Flyvbjerg",
    cluster: "KPIs & Analytics",
    clusterFr: "Indicateurs et flux en direct",
    clusterKey: "analytics",
    jurisdiction: "National",
    scale: "P10 / P50 / P90 Hazard",
    scaleFr: "Risque P10 / P50 / P90",
    horizon: "Continual",
    impactMetric: "Tail Risk Mitigation",
    impactMetricFr: "Atténuation du risque extrême",
    description: "Reference-class empirical forecasting preventing catastrophic multi-billion dollar budget slips.",
    descriptionFr: "Prévision empirique par classe de référence prévenant les dérives budgétaires catastrophiques.",
  },
];

export default function AnalyticsPage() {
  const { language, setLanguage } = useLanguage();
  const isFr = language === "fr";

  const [activeTab, setActiveTab] = useState<"overview" | "scorecard" | "simulator" | "top100">("overview");
  const [selectedProjectSlug, setSelectedProjectSlug] = useState<string>(SNAPSHOT_PROJECTS[0]?.slug || "darlington-smr");
  const [isStreaming, setIsStreaming] = useState<boolean>(true);
  const [liveTicks, setLiveTicks] = useState<LiveFeedTick[]>(INITIAL_LIVE_TICKS);
  const [selectedPillarFilter, setSelectedPillarFilter] = useState<string>("ALL");
  const [prioritySearch, setPrioritySearch] = useState<string>("");
  const [selectedPriorityCluster, setSelectedPriorityCluster] = useState<string>("all");
  const [copiedStatus, setCopiedStatus] = useState<string>("");

  // Macro shock simulator variables
  const [carbonTaxDelta, setCarbonTaxDelta] = useState<number>(0);
  const [rateHikeBps, setRateHikeBps] = useState<number>(0);
  const [tariffShockPct, setTariffShockPct] = useState<number>(0);
  const [reviewDelayMonths, setReviewDelayMonths] = useState<number>(0);

  const tabListRef = useRef<HTMLDivElement>(null);

  // Selected project object
  const selectedProject = useMemo(() => {
    return SNAPSHOT_PROJECTS.find((p) => p.slug === selectedProjectSlug) || SNAPSHOT_PROJECTS[0];
  }, [selectedProjectSlug]);

  // Compute scorecard with shock overlays
  const scorecard: ProjectKPIScorecard = useMemo(() => {
    const base = evaluateProjectKPIs(selectedProject);

    if (carbonTaxDelta === 0 && rateHikeBps === 0 && tariffShockPct === 0 && reviewDelayMonths === 0) {
      return base;
    }

    const shockPenalty =
      (carbonTaxDelta / 100) * 4.5 +
      (rateHikeBps / 300) * 8.0 +
      (tariffShockPct / 25) * 6.5 +
      (reviewDelayMonths / 24) * 5.0;

    const shockedRating = Math.max(10, Math.round((base.overall_kpi_rating - shockPenalty) * 10) / 10);

    const shockedPillars = base.pillars.map((p) => {
      let pPenalty = 0;
      if (p.category === "ESG_DECARBONIZATION") pPenalty = (carbonTaxDelta / 100) * 8;
      if (p.category === "CAPITAL_VELOCITY") pPenalty = (rateHikeBps / 300) * 12;
      if (p.category === "SUPPLY_CHAIN_CONTENT") pPenalty = (tariffShockPct / 25) * 10;
      if (p.category === "REGULATORY_SPEED") pPenalty = (reviewDelayMonths / 24) * 14;

      const pScore = Math.max(10, Math.round(p.score - pPenalty));
      return {
        ...p,
        score: pScore,
        health:
          pScore >= 85
            ? "EXEMPLARY"
            : pScore >= 70
            ? "HEALTHY"
            : pScore >= 50
            ? "ATTENTION_REQUIRED"
            : "CRITICAL",
      };
    });

    const addedGaps = [...base.critical_action_gaps];
    if (rateHikeBps >= 150) {
      addedGaps.push(
        isFr
          ? `Choc de taux de la BdC de +${rateHikeBps} pb érode le ratio de couverture du service de la dette (RCSD)`
          : `BoC rate hike of +${rateHikeBps} bps impairs debt-service coverage ratio (DSCR)`
      );
    }
    if (tariffShockPct >= 10) {
      addedGaps.push(
        isFr
          ? `Tarifs frontaliers américains de +${tariffShockPct} % exigent un réapprovisionnement national immédiat`
          : `U.S. border tariff shock of +${tariffShockPct}% triggers supply-chain on-shoring requirement`
      );
    }
    if (reviewDelayMonths >= 6) {
      addedGaps.push(
        isFr
          ? `Délai réglementaire de l'AEIC de +${reviewDelayMonths} mois entraîne un risque d'escalade des coûts`
          : `IAAC regulatory delay of +${reviewDelayMonths} months increases carrying costs by $45k/day`
      );
    }

    return {
      ...base,
      overall_kpi_rating: shockedRating,
      pillars: shockedPillars as any,
      critical_action_gaps: addedGaps,
    };
  }, [selectedProject, carbonTaxDelta, rateHikeBps, tariffShockPct, reviewDelayMonths, isFr]);

  // Live telemetry ticker
  useEffect(() => {
    if (!isStreaming) return;

    const interval = setInterval(() => {
      setLiveTicks((prev) => {
        const randomIndex = Math.floor(Math.random() * prev.length);
        const item = prev[randomIndex];
        const drift = (Math.random() - 0.49) * 0.008;
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
      const q = prioritySearch.toLowerCase();
      const matchesSearch =
        prioritySearch === "" ||
        item.title.toLowerCase().includes(q) ||
        item.titleFr.toLowerCase().includes(q) ||
        item.description.toLowerCase().includes(q) ||
        item.descriptionFr.toLowerCase().includes(q) ||
        item.jurisdiction.toLowerCase().includes(q);
      return matchesCluster && matchesSearch;
    });
  }, [selectedPriorityCluster, prioritySearch]);

  const copySnapshotJSON = () => {
    const jsonStr = JSON.stringify(
      {
        version: "cegs-kpi-v1.0",
        project: scorecard,
        live_feeds: liveTicks,
        macro_summary: DEFAULT_MACRO_SUMMARY,
        timestamp: new Date().toISOString(),
      },
      null,
      2
    );
    navigator.clipboard.writeText(jsonStr);
    setCopiedStatus(isFr ? "Copié !" : "Copied!");
    setTimeout(() => setCopiedStatus(""), 3000);
  };

  // Keyboard navigation for accessible tabs (WCAG 2.1 / 2.2)
  const handleTabKeyDown = (e: React.KeyboardEvent, currentId: string) => {
    const tabs: ("overview" | "scorecard" | "simulator" | "top100")[] = ["overview", "scorecard", "simulator", "top100"];
    const currentIndex = tabs.indexOf(currentId as any);

    let nextIndex = -1;
    if (e.key === "ArrowRight" || e.key === "ArrowDown") {
      nextIndex = (currentIndex + 1) % tabs.length;
    } else if (e.key === "ArrowLeft" || e.key === "ArrowUp") {
      nextIndex = (currentIndex - 1 + tabs.length) % tabs.length;
    } else if (e.key === "Home") {
      nextIndex = 0;
    } else if (e.key === "End") {
      nextIndex = tabs.length - 1;
    }

    if (nextIndex !== -1) {
      e.preventDefault();
      const nextTab = tabs[nextIndex];
      setActiveTab(nextTab);
      const el = document.getElementById(`tab-${nextTab}`);
      el?.focus();
    }
  };

  const tabsConfig = [
    {
      id: "overview",
      label: isFr ? "Vue d'ensemble nationale" : "National Sovereign Overview",
      icon: Globe2,
    },
    {
      id: "scorecard",
      label: isFr ? "Fiche de projet sur 8 piliers" : "8-Pillar Project Scorecard",
      icon: BarChart3,
    },
    {
      id: "simulator",
      label: isFr ? "Simulateur de chocs macro" : "Macro Shock Stress-Testing",
      icon: Sliders,
    },
    {
      id: "top100",
      label: isFr ? "Top 100 des priorités stratégiques" : "Top 100 Priority Items Catalog",
      icon: Layers,
      badge: "100",
    },
  ];

  return (
    <div className="mx-auto max-w-7xl space-y-8 px-4 py-8 sm:px-6 lg:px-8">
      {/* Header Banner */}
      <header className="relative overflow-hidden rounded-3xl border border-primary/30 bg-gradient-to-b from-[#091510] via-[#050b08] to-[#040806] p-6 shadow-2xl sm:p-8">
        <div aria-hidden="true" className="pointer-events-none absolute -right-24 -top-24 h-96 w-96 rounded-full bg-primary/15 blur-3xl" />
        <div aria-hidden="true" className="pointer-events-none absolute -bottom-24 -left-24 h-96 w-96 rounded-full bg-gold/10 blur-3xl" />

        <div className="relative z-10 flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
          <div className="max-w-3xl">
            <div className="flex flex-wrap items-center gap-3">
              <div className="inline-flex items-center gap-2 rounded-full border border-primary/40 bg-primary/10 px-3 py-1 font-mono text-[11px] font-bold text-aurora">
                <Sparkles aria-hidden="true" className="h-3.5 w-3.5 animate-pulse" />{" "}
                {isFr ? "MOTEUR D'INTELLIGENCE DES KPI SOUVERAINS ET FLUX EN DIRECT" : "SOVEREIGN KPI INTELLIGENCE & LIVE FEEDS ENGINE"}
              </div>

              {/* In-Page Official Languages Bilingual Toggle */}
              <button
                onClick={() => setLanguage(isFr ? "en" : "fr")}
                aria-label={isFr ? "Afficher l'interface en anglais (EN)" : "Switch interface to Canadian French (FR)"}
                className="inline-flex min-h-[36px] items-center gap-1.5 rounded-full border border-primary/30 bg-card/80 px-3 py-1 text-xs font-bold text-text-main transition hover:border-aurora hover:text-aurora focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora"
              >
                <span>{isFr ? "🇬🇧 EN" : "⚜️ FR"}</span>
                <span className="font-mono text-[10px] text-text-muted">
                  {isFr ? "(Loi C-13)" : "(Bill C-13)"}
                </span>
              </button>
            </div>

            <h1 className="mt-3 text-3xl font-black tracking-tight text-text-main sm:text-5xl">
              {isFr ? "Analyses de l'infrastructure " : "National Infrastructure "}
              <span className="text-aurora">{isFr ? "nationale" : "Analytics"}</span>
            </h1>
            <p className="mt-3 text-sm leading-relaxed text-text-muted">
              {isFr
                ? "Suivi opérationnel continu en temps réel sur 8 piliers stratégiques : décarbonation ESG, souveraineté autochtone, vélocité du capital, chaînes d'approvisionnement souveraines, physique du réseau électrique, délais d'autorisation, main-d'œuvre qualifiée Sceau rouge et étalons de matières premières."
                : "Continuous live tracking across 8 strategic pillars: ESG Decarbonization, Indigenous Sovereignty, Capital Spend Velocity, Sovereign Supply Chains, Power Grid Physics, Permitting Latency, Red Seal Craft Labour, and Commodity Benchmarks."}
            </p>
          </div>

          <div className="flex flex-wrap items-center gap-3">
            <button
              onClick={() => setIsStreaming(!isStreaming)}
              aria-pressed={isStreaming}
              aria-label={isStreaming ? (isFr ? "Mettre en pause le flux en direct" : "Pause live streaming") : (isFr ? "Activer le flux en direct" : "Start live streaming")}
              className={`inline-flex min-h-[44px] items-center gap-2 rounded-xl border px-4 py-2.5 font-mono text-xs font-bold transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora ${
                isStreaming
                  ? "border-emerald-500/50 bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20"
                  : "border-border bg-surface text-text-muted hover:text-text-main"
              }`}
            >
              {isStreaming ? (
                <>
                  <Pause aria-hidden="true" className="h-4 w-4" />{" "}
                  {isFr ? "FLUX EN DIRECT (2,4 s)" : "STREAMING LIVE (2.4s)"}
                </>
              ) : (
                <>
                  <Play aria-hidden="true" className="h-4 w-4" />{" "}
                  {isFr ? "FLUX EN PAUSE" : "STREAM PAUSED"}
                </>
              )}
            </button>
            <button
              onClick={copySnapshotJSON}
              aria-label={isFr ? "Exporter l'instantané des KPI en JSON" : "Export KPI Snapshot JSON"}
              className="inline-flex min-h-[44px] items-center gap-2 rounded-xl border border-primary/40 bg-primary/10 px-4 py-2.5 text-xs font-bold text-aurora hover:bg-primary/20 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora"
            >
              <Download aria-hidden="true" className="h-4 w-4" />
              <span>{copiedStatus || (isFr ? "Exporter l'instantané JSON" : "Export Snapshot JSON")}</span>
            </button>
          </div>
        </div>

        {/* Live Commodity & Rate Ribbon (Screen reader live region) */}
        <div
          role="region"
          aria-label={isFr ? "Bandeau des flux de marché en direct" : "Live market benchmark ticker"}
          className="mt-8 rounded-2xl border border-border/70 bg-black/40 p-4 backdrop-blur-md"
        >
          <div className="flex items-center justify-between pb-3">
            <div className="flex items-center gap-2 font-mono text-[11px] uppercase tracking-wider text-text-muted">
              <Activity aria-hidden="true" className="h-3.5 w-3.5 text-aurora" />{" "}
              {isFr ? "Téléscripteur d'étalons en direct (Provenance SHA-256)" : "Live Benchmark Ticker (SHA-256 Provenance)"}
            </div>
            <span className="font-mono text-[10px] text-text-muted">
              {isFr ? "8 flux de marché actifs" : "8 Active Market Feeds"}
            </span>
          </div>

          <div
            aria-live="polite"
            aria-atomic="false"
            className="grid grid-cols-1 gap-2.5 sm:grid-cols-2 md:grid-cols-4 lg:grid-cols-8"
          >
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
                <div className="mt-0.5 truncate font-mono text-[9px] text-text-muted">{t.unit}</div>
              </div>
            ))}
          </div>
        </div>
      </header>

      {/* Accessible Navigation Tabs (WCAG 2.1 / 2.2 Tablist) */}
      <nav aria-label={isFr ? "Navigation des analyses de performance" : "Analytics Performance Navigation"}>
        <div
          ref={tabListRef}
          role="tablist"
          aria-label={isFr ? "Sélection de la vue d'analyse" : "Analytics view selection"}
          className="flex flex-wrap gap-2 border-b border-border/80 pb-2 sm:flex-nowrap sm:overflow-x-auto"
        >
          {tabsConfig.map((tab) => {
            const Icon = tab.icon;
            const active = activeTab === tab.id;
            return (
              <button
                key={tab.id}
                id={`tab-${tab.id}`}
                role="tab"
                aria-selected={active}
                aria-controls={`panel-${tab.id}`}
                tabIndex={active ? 0 : -1}
                onClick={() => setActiveTab(tab.id as any)}
                onKeyDown={(e) => handleTabKeyDown(e, tab.id)}
                className={`inline-flex min-h-[44px] items-center gap-2 rounded-xl border px-4 py-2.5 text-xs font-bold transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora ${
                  active
                    ? "border-aurora bg-primary/15 text-aurora shadow-sm"
                    : "border-transparent text-text-muted hover:border-border hover:text-text-main"
                }`}
              >
                <Icon aria-hidden="true" className="h-4 w-4" />
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
        <section
          id="panel-overview"
          role="tabpanel"
          tabIndex={0}
          aria-labelledby="tab-overview"
          className="space-y-6 focus:outline-none"
        >
          <h2 className="sr-only">
            {isFr ? "Vue d'ensemble de la performance souveraine nationale" : "National Sovereign Performance Overview"}
          </h2>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div className="glass-card rounded-2xl border border-border p-5">
              <div className="flex items-center justify-between text-text-muted">
                <span className="font-mono text-xs uppercase tracking-wider">
                  {isFr ? "Décarbonation annuelle" : "Annual Decarbonization"}
                </span>
                <Leaf aria-hidden="true" className="h-4 w-4 text-emerald-400" />
              </div>
              <div className="mt-2 text-2xl font-black text-text-main tabular-nums">42.8 Mt CO2e/yr</div>
              <p className="mt-1 text-xs text-text-muted">
                {isFr
                  ? "Émissions cumulées évitées sur les projets suivis"
                  : "Aggregate avoided emissions across tracked projects"}
              </p>
            </div>

            <div className="glass-card rounded-2xl border border-border p-5">
              <div className="flex items-center justify-between text-text-muted">
                <span className="font-mono text-xs uppercase tracking-wider">
                  {isFr ? "Participation autochtone" : "First Nations Equity"}
                </span>
                <Users aria-hidden="true" className="h-4 w-4 text-gold" />
              </div>
              <div className="mt-2 text-2xl font-black text-gold tabular-nums">18.4% Average</div>
              <p className="mt-1 text-xs text-text-muted">
                {isFr
                  ? "1,45 G$ engagés sous l'enveloppe fédérale de 5 G$"
                  : "$1.45B committed under $5B Federal ILGP"}
              </p>
            </div>

            <div className="glass-card rounded-2xl border border-border p-5">
              <div className="flex items-center justify-between text-text-muted">
                <span className="font-mono text-xs uppercase tracking-wider">
                  {isFr ? "Déficit métiers Sceau rouge" : "Red Seal Labour Deficit"}
                </span>
                <AlertTriangle aria-hidden="true" className="h-4 w-4 text-amber-400" />
              </div>
              <div className="mt-2 text-2xl font-black text-amber-400 tabular-nums">14,200 FTE</div>
              <p className="mt-1 text-xs text-text-muted">
                {isFr
                  ? "Pénurie maximale d'électriciens et de tuyauteurs d'ici 2028"
                  : "Peak trade gap forecast by 2028 across ON/AB/BC"}
              </p>
            </div>

            <div className="glass-card rounded-2xl border border-border p-5">
              <div className="flex items-center justify-between text-text-muted">
                <span className="font-mono text-xs uppercase tracking-wider">
                  {isFr ? "Résilience des minéraux" : "Critical Minerals Self-Sufficiency"}
                </span>
                <ShieldCheck aria-hidden="true" className="h-4 w-4 text-aurora" />
              </div>
              <div className="mt-2 text-2xl font-black text-aurora tabular-nums">58.2/100</div>
              <p className="mt-1 text-xs text-text-muted">
                {isFr
                  ? "Capacité de raffinage intérieur vs dépendance étrangère"
                  : "Domestic refining capacity vs foreign monopoly exposure"}
              </p>
            </div>
          </div>

          {/* 8 Pillar Filter & Cards */}
          <div className="space-y-4">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
              <div>
                <h3 className="text-xl font-bold text-text-main">
                  {isFr ? "Piliers de performance souveraine" : "Sovereign Performance Pillars"}
                </h3>
                <p className="text-xs text-text-muted">
                  {isFr
                    ? "Mesures fondées sur des textes de loi avec provenance cryptographique SHA-256."
                    : "Statutory evidence-backed metrics with continuous SHA-256 data provenance."}
                </p>
              </div>

              <div className="flex items-center gap-2">
                <label htmlFor="pillar-filter-select" className="text-xs text-text-muted">
                  {isFr ? "Filtrer par pilier :" : "Filter by Pillar:"}
                </label>
                <select
                  id="pillar-filter-select"
                  value={selectedPillarFilter}
                  onChange={(e) => setSelectedPillarFilter(e.target.value)}
                  className="min-h-[44px] rounded-xl border border-border bg-card px-3 text-xs font-bold text-text-main focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora"
                >
                  <option value="ALL">{isFr ? "Tous les 8 piliers (20+ mesures)" : "All 8 Pillars (20+ Metrics)"}</option>
                  <option value="ESG_DECARBONIZATION">{isFr ? "Décarbonation ESG" : "ESG Decarbonization"}</option>
                  <option value="INDIGENOUS_EQUITY">{isFr ? "Souveraineté autochtone" : "Indigenous Sovereignty"}</option>
                  <option value="CAPITAL_VELOCITY">{isFr ? "Vélocité du capital" : "Capital Spend Velocity"}</option>
                  <option value="SUPPLY_CHAIN_CONTENT">{isFr ? "Chaînes d'approvisionnement" : "Sovereign Supply Chains"}</option>
                  <option value="GRID_PHYSICS">{isFr ? "Physique du réseau" : "Grid Power Physics"}</option>
                  <option value="REGULATORY_SPEED">{isFr ? "Efficacité réglementaire" : "Permitting Efficiency"}</option>
                  <option value="LABOR_SKILLS">{isFr ? "Main-d'œuvre Sceau rouge" : "Red Seal Labour Force"}</option>
                  <option value="COMMODITY_MACRO">{isFr ? "Étalons de marché" : "Commodity Benchmarks"}</option>
                </select>
              </div>
            </div>

            <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
              {CANONICAL_KPIS.filter(
                (k) => selectedPillarFilter === "ALL" || k.category === selectedPillarFilter
              ).map((kpi) => (
                <div
                  key={kpi.code}
                  className="rounded-2xl border border-border bg-card/70 p-5 shadow-lg transition-all hover:border-primary/50"
                >
                  <div className="flex items-start justify-between">
                    <div>
                      <span className="font-mono text-[10px] font-bold text-aurora">{kpi.code}</span>
                      <h4 className="mt-0.5 text-base font-bold text-text-main">{kpi.name}</h4>
                    </div>
                    <span className="rounded-full border border-border bg-surface px-2.5 py-1 font-mono text-[10px] text-text-muted">
                      {kpi.update_frequency}
                    </span>
                  </div>

                  <p className="mt-2 text-xs leading-relaxed text-text-muted">{kpi.description}</p>

                  <div className="mt-4 flex items-baseline justify-between border-t border-border/60 pt-3">
                    <div>
                      <span className="font-mono text-[10px] text-text-muted">
                        {isFr ? "Objectif de référence :" : "Target Benchmark:"}
                      </span>
                      <div className="font-mono text-sm font-bold text-text-main">
                        {kpi.target_benchmark} {kpi.unit}
                      </div>
                    </div>
                    <div className="text-right">
                      <span className="font-mono text-[10px] text-text-muted">
                        {isFr ? "Fondement législatif :" : "Statutory Basis:"}
                      </span>
                      <div className="font-mono text-xs text-text-muted">{kpi.statutory_basis}</div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </section>
      )}

      {/* TAB 2: Project Scorecard Explorer */}
      {activeTab === "scorecard" && (
        <section
          id="panel-scorecard"
          role="tabpanel"
          tabIndex={0}
          aria-labelledby="tab-scorecard"
          className="space-y-6 focus:outline-none"
        >
          <div className="flex flex-col gap-4 rounded-2xl border border-border bg-surface/80 p-5 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <label htmlFor="project-selector" className="font-mono text-xs uppercase tracking-wider text-text-muted">
                {isFr ? "Sélectionner le mégaprojet à évaluer" : "Select Megaproject to Assess"}
              </label>
              <h2 id="project-scorecard-explorer-title" className="text-xl font-bold text-text-main">
                {scorecard.project_name}
              </h2>
              <div className="mt-1 flex flex-wrap gap-2 text-xs text-text-muted">
                <span>
                  {isFr ? "Secteur : " : "Sector: "}
                  <strong className="text-text-main">{scorecard.sector}</strong>
                </span>
                <span>•</span>
                <span>
                  {isFr ? "Province : " : "Province: "}
                  <strong className="text-text-main">{scorecard.province}</strong>
                </span>
                <span>•</span>
                <span>
                  {isFr ? "Étape : " : "Stage: "}
                  <strong className="text-text-main">{scorecard.current_stage}</strong>
                </span>
                <span>•</span>
                <span>
                  CAPEX:{" "}
                  <strong className="text-text-main">
                    {isFr
                      ? `${(scorecard.total_capex_cad / 1_000_000).toLocaleString()} M$ CA`
                      : `$${(scorecard.total_capex_cad / 1_000_000).toLocaleString()}M CAD`}
                  </strong>
                </span>
              </div>
            </div>

            <div className="flex flex-wrap items-center gap-4">
              <select
                id="project-selector"
                value={selectedProjectSlug}
                onChange={(e) => setSelectedProjectSlug(e.target.value)}
                aria-label={isFr ? "Sélectionner un mégaprojet pour l'évaluation" : "Select Megaproject to Assess"}
                className="min-h-[44px] rounded-xl border border-border bg-card px-3 py-2 text-xs font-bold text-text-main focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora"
              >
                {SNAPSHOT_PROJECTS.slice(0, 30).map((p) => (
                  <option key={p.slug} value={p.slug}>
                    {p.name} ({p.province})
                  </option>
                ))}
              </select>

              <div className="rounded-xl border border-primary/40 bg-primary/10 px-4 py-2 text-right">
                <div className="font-mono text-[9px] uppercase tracking-wider text-text-muted">
                  {isFr ? "NOTE GLOBALE" : "OVERALL KPI RATING"}
                </div>
                <div className="text-2xl font-black tabular-nums text-aurora">{scorecard.overall_kpi_rating}/100</div>
              </div>
            </div>
          </div>

          {/* Critical Gaps Alert if any */}
          {scorecard.critical_action_gaps.length > 0 && (
            <div role="alert" className="rounded-2xl border border-rose-500/40 bg-rose-500/10 p-4">
              <div className="flex items-center gap-2 font-mono text-xs font-bold text-rose-300">
                <ShieldAlert aria-hidden="true" className="h-4 w-4" />{" "}
                {isFr ? "ÉCARTS D'ACTION CRITIQUES IDENTIFIÉS" : "CRITICAL ACTION GAPS IDENTIFIED"} ({scorecard.critical_action_gaps.length})
              </div>
              <ul className="mt-2 space-y-1 text-xs text-rose-100">
                {scorecard.critical_action_gaps.map((gap, i) => (
                  <li key={i} className="flex items-start gap-2">
                    <span aria-hidden="true" className="text-rose-400">•</span>
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
                          ? "bg-emerald-500/20 text-emerald-300"
                          : pillar.health === "HEALTHY"
                          ? "bg-cyan-500/20 text-cyan-300"
                          : pillar.health === "ATTENTION_REQUIRED"
                          ? "bg-amber-500/20 text-amber-300"
                          : "bg-rose-500/20 text-rose-300"
                      }`}
                    >
                      {pillar.health === "EXEMPLARY"
                        ? isFr ? "EXEMPLAIRE" : "EXEMPLARY"
                        : pillar.health === "HEALTHY"
                        ? isFr ? "SAIN" : "HEALTHY"
                        : pillar.health === "ATTENTION_REQUIRED"
                        ? isFr ? "ATTENTION REQUISE" : "ATTENTION REQUIRED"
                        : isFr ? "CRITIQUE" : "CRITICAL"}
                    </span>
                  </div>
                  <div className="text-right">
                    <div className="text-xl font-black text-text-main tabular-nums">{pillar.score}/100</div>
                    <div className="font-mono text-[9px] text-text-muted">
                      {isFr ? "Score normalisé" : "Normalized Score"}
                    </div>
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
                              ? "bg-emerald-500/20 text-emerald-300"
                              : m.performance_rank === "ON_TARGET"
                              ? "bg-cyan-500/20 text-cyan-300"
                              : "bg-amber-500/20 text-amber-300"
                          }`}
                        >
                          {m.performance_rank === "SUPERIOR"
                            ? isFr ? "SUPÉRIEUR" : "SUPERIOR"
                            : m.performance_rank === "ON_TARGET"
                            ? isFr ? "CONFORME À L'OBJECTIF" : "ON TARGET"
                            : isFr ? "À AMÉLIORER" : "NEEDS IMPROVEMENT"}
                        </span>
                      </div>

                      <div className="mt-2 flex items-baseline justify-between font-mono text-xs">
                        <span className="text-text-muted">
                          {isFr ? "Observé : " : "Observed: "}
                          <strong className="text-text-main tabular-nums">
                            {m.observed_value.toLocaleString()} {m.unit}
                          </strong>
                        </span>
                        <span className="text-text-muted">
                          {isFr ? "Cible : " : "Target: "}
                          <strong className="tabular-nums text-text-muted">
                            {m.target_benchmark.toLocaleString()} {m.unit}
                          </strong>
                        </span>
                      </div>

                      <div className="mt-2 flex items-center justify-between text-[10px] text-text-muted">
                        <span className="truncate" title={m.notes}>
                          {m.notes}
                        </span>
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
        <section
          id="panel-simulator"
          role="tabpanel"
          tabIndex={0}
          aria-labelledby="tab-simulator"
          className="space-y-6 focus:outline-none"
        >
          <div className="rounded-2xl border border-primary/40 bg-gradient-to-br from-card via-surface to-background p-6">
            <div className="inline-flex items-center gap-2 rounded-full border border-primary/40 bg-primary/10 px-3 py-1 font-mono text-[11px] font-bold text-aurora">
              <Sliders aria-hidden="true" className="h-3.5 w-3.5" />{" "}
              {isFr ? "SIMULATEUR STOCHASTIQUE DE CHOCS MACROÉCONOMIQUES" : "STOCHASTIC MACRO SHOCK SIMULATOR"}
            </div>
            <h2 id="macro-shock-simulator-title" className="mt-2 text-2xl font-black text-text-main">
              {isFr
                ? `Simuler des chocs politiques et de marché sur ${scorecard.project_name}`
                : `Simulate Policy & Commodity Shocks on ${scorecard.project_name}`}
            </h2>
            <p className="mt-1 max-w-2xl text-xs leading-relaxed text-text-muted">
              {isFr
                ? "Ajustez les curseurs ci-dessous pour observer le recalcul en temps réel de la note du projet, de sa capacité de service de la dette et des écarts critiques."
                : "Adjust variables below to observe real-time recalculation of the project's 8-pillar rating, debt service capacity, and critical action gaps."}
            </p>

            <div className="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
              {/* Slider 1: Carbon Price Escalation */}
              <div className="rounded-xl border border-border bg-black/30 p-4">
                <div className="flex items-center justify-between text-xs font-bold text-text-main">
                  <label htmlFor="slider-carbon-tax">
                    {isFr ? "Hausse de la taxe carbone" : "Carbon Tax Escalation"}
                  </label>
                  <span className="font-mono text-emerald-400">
                    +{carbonTaxDelta} {isFr ? "$/t" : "$/t"}
                  </span>
                </div>
                <input
                  id="slider-carbon-tax"
                  type="range"
                  min="0"
                  max="100"
                  step="5"
                  value={carbonTaxDelta}
                  onChange={(e) => setCarbonTaxDelta(Number(e.target.value))}
                  aria-valuemin={0}
                  aria-valuemax={100}
                  aria-valuenow={carbonTaxDelta}
                  aria-valuetext={`+${carbonTaxDelta} dollars per tonne`}
                  className="mt-3 w-full accent-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora"
                />
                <div className="mt-2 flex justify-between font-mono text-[10px] text-text-muted">
                  <span>$0 ({isFr ? "Base" : "Baseline"})</span>
                  <span>+$100 ($170/t Cap)</span>
                </div>
              </div>

              {/* Slider 2: Interest Rate Shock */}
              <div className="rounded-xl border border-border bg-black/30 p-4">
                <div className="flex items-center justify-between text-xs font-bold text-text-main">
                  <label htmlFor="slider-boc-rate">
                    {isFr ? "Hausse de taux de la BdC" : "BoC Rate Hike"}
                  </label>
                  <span className="font-mono text-amber-400">+{rateHikeBps} bps</span>
                </div>
                <input
                  id="slider-boc-rate"
                  type="range"
                  min="0"
                  max="300"
                  step="25"
                  value={rateHikeBps}
                  onChange={(e) => setRateHikeBps(Number(e.target.value))}
                  aria-valuemin={0}
                  aria-valuemax={300}
                  aria-valuenow={rateHikeBps}
                  aria-valuetext={`+${rateHikeBps} basis points`}
                  className="mt-3 w-full accent-amber-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora"
                />
                <div className="mt-2 flex justify-between font-mono text-[10px] text-text-muted">
                  <span>0 bps</span>
                  <span>+300 bps ({isFr ? "Stagflation" : "Stagflation"})</span>
                </div>
              </div>

              {/* Slider 3: U.S. Border Tariff Shock */}
              <div className="rounded-xl border border-border bg-black/30 p-4">
                <div className="flex items-center justify-between text-xs font-bold text-text-main">
                  <label htmlFor="slider-tariff-shock">
                    {isFr ? "Tarifs douaniers américains" : "U.S. Border Tariff Shock"}
                  </label>
                  <span className="font-mono text-rose-400">+{tariffShockPct}%</span>
                </div>
                <input
                  id="slider-tariff-shock"
                  type="range"
                  min="0"
                  max="25"
                  step="2.5"
                  value={tariffShockPct}
                  onChange={(e) => setTariffShockPct(Number(e.target.value))}
                  aria-valuemin={0}
                  aria-valuemax={25}
                  aria-valuenow={tariffShockPct}
                  aria-valuetext={`+${tariffShockPct} percent tariff`}
                  className="mt-3 w-full accent-rose-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora"
                />
                <div className="mt-2 flex justify-between font-mono text-[10px] text-text-muted">
                  <span>0% ({isFr ? "Sans tarif" : "CUSMA Free"})</span>
                  <span>+25% (Section 232)</span>
                </div>
              </div>

              {/* Slider 4: Regulatory Delay */}
              <div className="rounded-xl border border-border bg-black/30 p-4">
                <div className="flex items-center justify-between text-xs font-bold text-text-main">
                  <label htmlFor="slider-review-delay">
                    {isFr ? "Délai d'examen de l'AEIC" : "IAAC Review Delay"}
                  </label>
                  <span className="font-mono text-cyan-400">
                    +{reviewDelayMonths} {isFr ? "mois" : "mo"}
                  </span>
                </div>
                <input
                  id="slider-review-delay"
                  type="range"
                  min="0"
                  max="24"
                  step="3"
                  value={reviewDelayMonths}
                  onChange={(e) => setReviewDelayMonths(Number(e.target.value))}
                  aria-valuemin={0}
                  aria-valuemax={24}
                  aria-valuenow={reviewDelayMonths}
                  aria-valuetext={`+${reviewDelayMonths} months delay`}
                  className="mt-3 w-full accent-cyan-400 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora"
                />
                <div className="mt-2 flex justify-between font-mono text-[10px] text-text-muted">
                  <span>0 {isFr ? "mois" : "mo"}</span>
                  <span>+24 {isFr ? "mois" : "mo"}</span>
                </div>
              </div>
            </div>

            {/* Shock Results Display */}
            <div className="mt-8 rounded-xl border border-border bg-card p-5">
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h3 className="text-base font-bold text-text-main">
                    {isFr ? "Sommaire de l'impact des chocs calibrés" : "Calibrated Shock Impact Summary"}
                  </h3>
                  <p className="text-xs text-text-muted">
                    {isFr
                      ? "La note recalibrée reflète la résilience financière globale du projet face à des conditions macroéconomiques défavorables."
                      : "Recalibrated rating reflects compounded project finance resilience under adverse macroeconomic conditions."}
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
                    aria-label={isFr ? "Réinitialiser tous les paramètres de chocs" : "Reset all macro shock parameters"}
                    className="inline-flex min-h-[44px] items-center gap-1.5 rounded-lg border border-border px-3 py-1.5 text-xs text-text-muted hover:text-text-main focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora"
                  >
                    <RefreshCw aria-hidden="true" className="h-3.5 w-3.5" />{" "}
                    {isFr ? "Réinitialiser les chocs" : "Reset Shocks"}
                  </button>
                  <div className="rounded-xl border border-primary/40 bg-primary/10 px-4 py-2 text-right">
                    <span className="font-mono text-[9px] uppercase text-text-muted">
                      {isFr ? "NOTE APRÈS CHOCS" : "SHOCKED KPI RATING"}
                    </span>
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
        <section
          id="panel-top100"
          role="tabpanel"
          tabIndex={0}
          aria-labelledby="tab-top100"
          className="space-y-6 focus:outline-none"
        >
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h2 id="top-100-priority-items-title" className="text-2xl font-black text-text-main">
                {isFr ? "Top 100 des priorités souveraines canadiennes" : "Top 100 Canadian Sovereign Priority Items"}
              </h2>
              <p className="text-xs text-text-muted">
                {isFr
                  ? "Le registre officiel des grands projets, actifs d'infrastructure, flux de données en direct et réformes législatives requis pour la souveraineté économique."
                  : "The authoritative registry of major projects, infrastructure assets, live data feeds, and statutory policy reforms required for economic sovereignty."}
              </p>
            </div>

            <div className="flex flex-wrap items-center gap-3">
              <div className="relative flex-1 sm:flex-none">
                <label htmlFor="priority-search-input" className="sr-only">
                  {isFr ? "Rechercher parmi les 100 priorités" : "Search top 100 priorities"}
                </label>
                <Search aria-hidden="true" className="absolute left-3 top-3.5 h-3.5 w-3.5 text-text-muted" />
                <input
                  id="priority-search-input"
                  type="text"
                  placeholder={isFr ? "Filtrer les priorités..." : "Filter priorities..."}
                  value={prioritySearch}
                  onChange={(e) => setPrioritySearch(e.target.value)}
                  className="min-h-[44px] w-full rounded-xl border border-border bg-card pl-9 pr-3 text-xs text-text-main focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora sm:w-64"
                />
              </div>

              <div>
                <label htmlFor="priority-cluster-select" className="sr-only">
                  {isFr ? "Filtrer par groupe de priorités" : "Filter by priority cluster"}
                </label>
                <select
                  id="priority-cluster-select"
                  value={selectedPriorityCluster}
                  onChange={(e) => setSelectedPriorityCluster(e.target.value)}
                  className="min-h-[44px] rounded-xl border border-border bg-card px-3 text-xs font-bold text-text-main focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aurora"
                >
                  <option value="all">{isFr ? "Tous les groupes (100)" : "All Clusters (100)"}</option>
                  <option value="minerals">{isFr ? "Minéraux critiques (15)" : "Critical Minerals (15)"}</option>
                  <option value="nuclear">{isFr ? "Réseau nucléaire et propre (12)" : "Nuclear & Clean Grid (12)"}</option>
                  <option value="energy">{isFr ? "Énergie propre et stockage (11)" : "Clean Energy & Storage (11)"}</option>
                  <option value="compute">{isFr ? "Calcul et IA souveraine (10)" : "Sovereign AI Compute (10)"}</option>
                  <option value="ports">{isFr ? "Ports et corridors (12)" : "Ports & Corridors (12)"}</option>
                  <option value="indigenous">{isFr ? "Souveraineté autochtone (12)" : "Indigenous Sovereignty (12)"}</option>
                  <option value="analytics">{isFr ? "Indicateurs et flux en direct (13)" : "KPIs & Live Feeds (13)"}</option>
                  <option value="supply">{isFr ? "Défense de la chaîne d'approvisionnement (9)" : "Supply Chain Defense (9)"}</option>
                  <option value="regulatory">{isFr ? "Réformes réglementaires (6)" : "Regulatory Reforms (6)"}</option>
                </select>
              </div>
            </div>
          </div>

          {/* Desktop Data Table (Accessible with caption and scope) */}
          <div
            role="region"
            aria-label={isFr ? "Tableau des 100 priorités souveraines canadiennes" : "Top 100 Canadian Sovereign Priorities Table"}
            tabIndex={0}
            className="hidden md:block rounded-2xl border border-border bg-card/80 overflow-hidden shadow-xl focus:outline-none focus:ring-2 focus:ring-aurora"
          >
            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs">
                <caption className="sr-only">
                  {isFr
                    ? "Liste détaillée des 100 priorités souveraines canadiennes comprenant rang, titre, groupe stratégique, province, échelle, échéance et indicateur d'impact"
                    : "Detailed list of top 100 Canadian sovereign priorities including rank, title, strategic cluster, province, scale, horizon, and primary impact metric"}
                </caption>
                <thead className="border-b border-border bg-surface/80 font-mono text-[10px] uppercase text-text-muted">
                  <tr>
                    <th scope="col" className="px-4 py-3.5">
                      {isFr ? "Rang" : "Rank"}
                    </th>
                    <th scope="col" className="px-4 py-3.5">
                      {isFr ? "Priorité stratégique" : "Priority Item"}
                    </th>
                    <th scope="col" className="px-4 py-3.5">
                      {isFr ? "Groupe stratégique" : "Strategic Cluster"}
                    </th>
                    <th scope="col" className="px-4 py-3.5">
                      {isFr ? "Prov." : "Prov"}
                    </th>
                    <th scope="col" className="px-4 py-3.5">
                      {isFr ? "Échelle / Cible" : "Scale / Target"}
                    </th>
                    <th scope="col" className="px-4 py-3.5">
                      {isFr ? "Échéance" : "Horizon"}
                    </th>
                    <th scope="col" className="px-4 py-3.5">
                      {isFr ? "Indicateur d'impact principal" : "Primary Impact Metric"}
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border/60">
                  {filteredPriorities.map((item) => (
                    <tr key={item.rank} className="hover:bg-surface/60 transition-colors">
                      <td className="px-4 py-3.5 font-mono font-bold text-aurora">#{item.rank}</td>
                      <td className="px-4 py-3.5">
                        <div className="font-bold text-text-main">{isFr ? item.titleFr : item.title}</div>
                        <div className="mt-0.5 text-[11px] text-text-muted leading-tight">
                          {isFr ? item.descriptionFr : item.description}
                        </div>
                      </td>
                      <td className="px-4 py-3.5">
                        <span className="rounded-full bg-primary/10 border border-primary/20 px-2 py-0.5 font-mono text-[10px] text-aurora">
                          {isFr ? item.clusterFr : item.cluster}
                        </span>
                      </td>
                      <td className="px-4 py-3.5 font-mono font-bold text-text-main">{item.jurisdiction}</td>
                      <td className="px-4 py-3.5 font-mono text-text-main">{isFr ? item.scaleFr : item.scale}</td>
                      <td className="px-4 py-3.5 font-mono text-text-muted">{item.horizon}</td>
                      <td className="px-4 py-3.5 font-mono text-gold text-[11px]">
                        {isFr ? item.impactMetricFr : item.impactMetric}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Mobile-Accessible Card View (< md) */}
          <div className="block md:hidden space-y-3">
            <div className="flex items-center justify-between text-xs text-text-muted">
              <span>{isFr ? `${filteredPriorities.length} priorités affichées` : `Showing ${filteredPriorities.length} priorities`}</span>
            </div>
            {filteredPriorities.map((item) => (
              <div
                key={item.rank}
                className="rounded-2xl border border-border bg-card/90 p-4 shadow-md space-y-3"
              >
                <div className="flex items-center justify-between">
                  <span className="font-mono text-xs font-bold text-aurora">#{item.rank}</span>
                  <div className="flex items-center gap-2">
                    <span className="rounded-full bg-primary/10 border border-primary/20 px-2 py-0.5 font-mono text-[9px] text-aurora">
                      {isFr ? item.clusterFr : item.cluster}
                    </span>
                    <span className="rounded-full bg-surface px-2 py-0.5 font-mono text-[9px] font-bold text-text-main">
                      {item.jurisdiction}
                    </span>
                  </div>
                </div>

                <div>
                  <h3 className="text-sm font-bold text-text-main">{isFr ? item.titleFr : item.title}</h3>
                  <p className="mt-1 text-xs text-text-muted leading-relaxed">
                    {isFr ? item.descriptionFr : item.description}
                  </p>
                </div>

                <div className="flex flex-wrap items-center justify-between border-t border-border/50 pt-2.5 text-xs">
                  <div>
                    <span className="font-mono text-[10px] text-text-muted">{isFr ? "Échelle : " : "Scale: "}</span>
                    <strong className="font-mono text-text-main">{isFr ? item.scaleFr : item.scale}</strong>
                  </div>
                  <div>
                    <span className="font-mono text-[10px] text-text-muted">{isFr ? "Échéance : " : "Horizon: "}</span>
                    <strong className="font-mono text-text-muted">{item.horizon}</strong>
                  </div>
                </div>

                <div className="rounded-xl border border-gold/30 bg-gold/10 px-3 py-2">
                  <span className="font-mono text-[9px] uppercase text-text-muted">
                    {isFr ? "Indicateur d'impact clé" : "Key Impact Metric"}
                  </span>
                  <div className="font-mono text-xs font-bold text-gold">
                    {isFr ? item.impactMetricFr : item.impactMetric}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </section>
      )}
    </div>
  );
}
