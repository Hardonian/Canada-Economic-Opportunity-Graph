"use client";

import React, { useEffect, useRef, useState } from "react";
import Link from "next/link";
import {
  MapPin,
  Layers,
  Globe2,
  AlertTriangle,
  Compass,
  Sparkles,
  Maximize2,
  Minimize2,
  Navigation,
  KeyRound,
  Check,
  ChevronRight,
  Shield,
  Activity,
  Zap,
  ExternalLink,
  Flame,
  Ship,
  Info
} from "lucide-react";
import type { Project } from "@/lib/types";
import {
  MAP_PROVIDERS,
  GLOBAL_TRADE_ROUTES,
  CONFLICT_MARKERS,
  OPPORTUNITY_ZONES,
  STRATEGIC_CORRIDORS,
  MAP_FOCUS_PRESETS,
  type MapTileProvider,
  type TradeRoute,
  type ConflictMarker,
  type OpportunityZone
} from "@/lib/geospatial";

interface GeospatialMapProps {
  projects: Project[];
  activeProject: Project;
  onSelectProject: (p: Project) => void;
  selectedSector: string;
  selectedStage: string;
}

export default function GeospatialMap({
  projects,
  activeProject,
  onSelectProject,
  selectedSector,
  selectedStage,
}: GeospatialMapProps) {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const mapInstanceRef = useRef<any>(null);
  const tileLayerRef = useRef<any>(null);
  const layersGroupRef = useRef<any>({
    projects: null,
    tradeRoutes: null,
    conflictMarkers: null,
    opportunityZones: null,
    corridors: null,
  });

  // State
  const [activeProviderId, setActiveProviderId] = useState<string>("esri-satellite");
  const [showTradeRoutes, setShowTradeRoutes] = useState<boolean>(true);
  const [showConflictMarkers, setShowConflictMarkers] = useState<boolean>(true);
  const [showOpportunityZones, setShowOpportunityZones] = useState<boolean>(true);
  const [showCorridors, setShowCorridors] = useState<boolean>(true);
  const [cursorCoords, setCursorCoords] = useState<{ lat: number; lng: number } | null>(null);
  const [currentZoom, setCurrentZoom] = useState<number>(4);
  const [isFullscreen, setIsFullscreen] = useState<boolean>(false);
  const [selectedRoute, setSelectedRoute] = useState<TradeRoute | null>(null);
  const [selectedConflict, setSelectedConflict] = useState<ConflictMarker | null>(null);
  const [selectedZone, setSelectedZone] = useState<OpportunityZone | null>(null);
  const [googleApiKey, setGoogleApiKey] = useState<string>("");
  const [showKeyModal, setShowKeyModal] = useState<boolean>(false);
  const [isLeafletReady, setIsLeafletReady] = useState<boolean>(false);

  // Initialize Map
  useEffect(() => {
    let isMounted = true;

    async function initLeaflet() {
      if (!mapContainerRef.current || mapInstanceRef.current) return;

      const L = await import("leaflet");

      if (!isMounted || !mapContainerRef.current) return;

      // Create Leaflet map instance centered over Canada
      const map = L.map(mapContainerRef.current, {
        center: [58.0, -98.0],
        zoom: 4,
        minZoom: 2,
        maxZoom: 19,
        zoomControl: false,
        attributionControl: false,
      });

      mapInstanceRef.current = map;

      // Add zoom control to bottom right
      L.control.zoom({ position: "bottomright" }).addTo(map);

      // Track cursor coordinates
      map.on("mousemove", (e: any) => {
        setCursorCoords({
          lat: parseFloat(e.latlng.lat.toFixed(4)),
          lng: parseFloat(e.latlng.lng.toFixed(4)),
        });
      });

      map.on("zoomend", () => {
        setCurrentZoom(map.getZoom());
      });

      // Layer groups for dynamic filtering
      layersGroupRef.current.projects = L.layerGroup().addTo(map);
      layersGroupRef.current.tradeRoutes = L.layerGroup().addTo(map);
      layersGroupRef.current.conflictMarkers = L.layerGroup().addTo(map);
      layersGroupRef.current.opportunityZones = L.layerGroup().addTo(map);
      layersGroupRef.current.corridors = L.layerGroup().addTo(map);

      setIsLeafletReady(true);
    }

    initLeaflet();

    return () => {
      isMounted = false;
      if (mapInstanceRef.current) {
        mapInstanceRef.current.remove();
        mapInstanceRef.current = null;
      }
    };
  }, []);

  // Update Base Tile Layer
  useEffect(() => {
    if (!isLeafletReady || !mapInstanceRef.current) return;

    import("leaflet").then((L) => {
      const map = mapInstanceRef.current;
      if (tileLayerRef.current) {
        map.removeLayer(tileLayerRef.current);
      }

      const provider = MAP_PROVIDERS.find((p) => p.id === activeProviderId) || MAP_PROVIDERS[0];
      
      let tileUrl = provider.url;
      if (provider.id === "google-hybrid" && googleApiKey) {
        tileUrl = `https://mt1.google.com/vt/lyrs=y&x={x}&y={y}&z={z}&key=${googleApiKey}`;
      }

      tileLayerRef.current = L.tileLayer(tileUrl, {
        attribution: provider.attribution,
        maxZoom: provider.maxZoom,
        subdomains: provider.subdomains || "abc",
      }).addTo(map);
    });
  }, [isLeafletReady, activeProviderId, googleApiKey]);

  // Render Projects Markers
  useEffect(() => {
    if (!isLeafletReady || !layersGroupRef.current.projects) return;

    import("leaflet").then((L) => {
      const group = layersGroupRef.current.projects;
      group.clearLayers();

      function hasValidCoords(proj: Project): proj is Project & { latitude: number; longitude: number } {
        return (
          typeof proj.latitude === "number" &&
          typeof proj.longitude === "number" &&
          Number.isFinite(proj.latitude) &&
          Number.isFinite(proj.longitude)
        );
      }

      const mappable = projects.filter(hasValidCoords);

      mappable.forEach((p) => {
        const isSelected = activeProject.id === p.id;

        let color = "#00F5A0"; // Green: Construction / Operating
        if (p.current_stage === "PERMITTING" || p.current_stage === "ENVIRONMENTAL_REVIEW") {
          color = "#F59E0B"; // Gold: Permitting
        } else if (p.current_stage === "FEASIBILITY" || p.current_stage === "ANNOUNCED") {
          color = "#38BDF8"; // Cyan: Feasibility
        }

        const size = isSelected ? 22 : 14;
        const pingClass = isSelected ? "animate-ping opacity-75" : "";

        const customIcon = L.divIcon({
          className: "custom-project-marker",
          iconSize: [size, size],
          iconAnchor: [size / 2, size / 2],
          html: `
            <div style="position: relative; width: ${size}px; height: ${size}px;">
              ${
                isSelected
                  ? `<div style="position: absolute; inset: -6px; border-radius: 9999px; background: ${color}; opacity: 0.4;" class="${pingClass}"></div>`
                  : ""
              }
              <div style="width: ${size}px; height: ${size}px; border-radius: 9999px; background: ${color}; border: 2px solid #040806; box-shadow: 0 0 10px ${color}; display: flex; items-center; justify-content: center; cursor: pointer;">
                <div style="width: 4px; height: 4px; border-radius: 9999px; background: #ffffff;"></div>
              </div>
            </div>
          `,
        });

        const marker = L.marker([p.latitude, p.longitude], { icon: customIcon });

        const capexFormatted =
          p.capex_cad >= 1e9
            ? `$${(p.capex_cad / 1e9).toFixed(2)}B CAD`
            : `$${(p.capex_cad / 1e6).toFixed(0)}M CAD`;

        const popupContent = `
          <div style="min-width: 220px; font-family: monospace; padding: 2px;">
            <div style="font-size: 9px; color: ${color}; font-weight: bold; text-transform: uppercase; letter-spacing: 0.05em; margin-bottom: 4px;">
              ${p.current_stage || "PROJECT"} • ${p.province}
            </div>
            <div style="font-size: 13px; font-weight: bold; color: #f7faf8; margin-bottom: 6px; line-height: 1.3;">
              ${p.name}
            </div>
            <div style="font-size: 11px; color: #a8c4b8; margin-bottom: 6px;">
              Sector: <span style="color: #f7faf8;">${p.sector}</span><br/>
              CAPEX: <span style="color: #00F5A0; font-weight: bold;">${capexFormatted}</span>
            </div>
            <div style="border-top: 1px solid #1e3a2b; padding-top: 6px; font-size: 10px; color: #a8c4b8;">
              Proponent: ${p.proponent_name || "Lead Proponent"}
            </div>
          </div>
        `;

        marker.bindPopup(popupContent, {
          className: "institutional-popup",
        });

        marker.on("click", () => {
          onSelectProject(p);
        });

        group.addLayer(marker);
      });
    });
  }, [isLeafletReady, projects, activeProject, onSelectProject]);

  // Render Trade Routes
  useEffect(() => {
    if (!isLeafletReady || !layersGroupRef.current.tradeRoutes) return;

    import("leaflet").then((L) => {
      const group = layersGroupRef.current.tradeRoutes;
      group.clearLayers();

      if (!showTradeRoutes) return;

      GLOBAL_TRADE_ROUTES.forEach((route) => {
        const latLngs = route.waypoints.map((w) => [w.lat, w.lng]);

        const polyline = L.polyline(latLngs as any, {
          color: route.color,
          weight: 3,
          opacity: 0.85,
          dashArray: route.dashArray || "5, 5",
        });

        polyline.on("click", () => {
          setSelectedRoute(route);
          setSelectedConflict(null);
          setSelectedZone(null);
        });

        group.addLayer(polyline);

        // Waypoint markers
        route.waypoints.forEach((w) => {
          if (w.name) {
            const circle = L.circleMarker([w.lat, w.lng], {
              radius: 4,
              color: route.color,
              fillColor: "#040806",
              fillOpacity: 1,
              weight: 2,
            });
            circle.bindTooltip(`<b>${w.name}</b><br/>${route.name}`, {
              direction: "top",
              className: "trade-route-tooltip",
            });
            group.addLayer(circle);
          }
        });
      });
    });
  }, [isLeafletReady, showTradeRoutes]);

  // Render Conflict Markers
  useEffect(() => {
    if (!isLeafletReady || !layersGroupRef.current.conflictMarkers) return;

    import("leaflet").then((L) => {
      const group = layersGroupRef.current.conflictMarkers;
      group.clearLayers();

      if (!showConflictMarkers) return;

      CONFLICT_MARKERS.forEach((c) => {
        const iconColor = c.severity === "CRITICAL" ? "#EF4444" : "#F59E0B";

        const conflictIcon = L.divIcon({
          className: "conflict-marker-icon",
          iconSize: [20, 20],
          iconAnchor: [10, 10],
          html: `
            <div style="position: relative; width: 20px; height: 20px; cursor: pointer;">
              <div style="position: absolute; inset: -4px; border-radius: 9999px; background: ${iconColor}; opacity: 0.3;" class="animate-ping"></div>
              <div style="width: 20px; height: 20px; background: #180909; border: 2px solid ${iconColor}; border-radius: 6px; display: flex; align-items: center; justify-content: center; transform: rotate(45deg); box-shadow: 0 0 8px ${iconColor};">
                <div style="transform: rotate(-45deg); font-size: 10px; font-weight: bold; color: ${iconColor};">!</div>
              </div>
            </div>
          `,
        });

        const marker = L.marker([c.coordinates.lat, c.coordinates.lng], { icon: conflictIcon });

        marker.on("click", () => {
          setSelectedConflict(c);
          setSelectedRoute(null);
          setSelectedZone(null);
        });

        group.addLayer(marker);
      });
    });
  }, [isLeafletReady, showConflictMarkers]);

  // Render Opportunity Zones
  useEffect(() => {
    if (!isLeafletReady || !layersGroupRef.current.opportunityZones) return;

    import("leaflet").then((L) => {
      const group = layersGroupRef.current.opportunityZones;
      group.clearLayers();

      if (!showOpportunityZones) return;

      OPPORTUNITY_ZONES.forEach((zone) => {
        const circle = L.circle([zone.center.lat, zone.center.lng], {
          radius: zone.boundsRadiusKm * 1000,
          color: zone.color,
          fillColor: zone.color,
          fillOpacity: 0.12,
          weight: 1.5,
          dashArray: "4, 4",
        });

        circle.on("click", () => {
          setSelectedZone(zone);
          setSelectedRoute(null);
          setSelectedConflict(null);
        });

        circle.bindTooltip(
          `<b>${zone.name}</b><br/><span style="color:${zone.color}">${zone.estimatedEndowmentCAD}</span>`,
          {
            direction: "center",
            permanent: false,
            className: "opportunity-tooltip",
          }
        );

        group.addLayer(circle);
      });
    });
  }, [isLeafletReady, showOpportunityZones]);

  // Render Strategic Corridors
  useEffect(() => {
    if (!isLeafletReady || !layersGroupRef.current.corridors) return;

    import("leaflet").then((L) => {
      const group = layersGroupRef.current.corridors;
      group.clearLayers();

      if (!showCorridors) return;

      STRATEGIC_CORRIDORS.forEach((c) => {
        const polyline = L.polyline(
          [
            [c.fromCoords.lat, c.fromCoords.lng],
            [c.toCoords.lat, c.toCoords.lng],
          ],
          {
            color: c.color,
            weight: 2.5,
            dashArray: "4, 6",
            opacity: 0.75,
          }
        );

        polyline.bindTooltip(`<b>${c.name}</b><br/>${c.capacity}`, {
          direction: "center",
          className: "corridor-tooltip",
        });

        group.addLayer(polyline);
      });
    });
  }, [isLeafletReady, showCorridors]);

  // Pan to preset function
  const handlePresetSelect = (preset: (typeof MAP_FOCUS_PRESETS)[0]) => {
    if (mapInstanceRef.current) {
      mapInstanceRef.current.flyTo(preset.center, preset.zoom, {
        duration: 1.2,
      });
    }
  };

  // Fly to active project when changed from outside
  useEffect(() => {
    if (
      mapInstanceRef.current &&
      typeof activeProject.latitude === "number" &&
      typeof activeProject.longitude === "number"
    ) {
      mapInstanceRef.current.flyTo([activeProject.latitude, activeProject.longitude], 7, {
        duration: 1.0,
      });
    }
  }, [activeProject]);

  return (
    <div
      className={`relative w-full rounded-2xl border border-border/80 overflow-hidden bg-[#040806] shadow-2xl transition-all duration-300 ${
        isFullscreen ? "fixed inset-0 z-50 rounded-none border-none h-screen" : "h-[680px]"
      }`}
    >
      {/* Top Interactive HUD Bar */}
      <div className="absolute top-3 left-3 right-3 z-[1000] flex flex-wrap items-center justify-between gap-2 pointer-events-none">
        {/* Left: Base Map & Presets */}
        <div className="flex flex-wrap items-center gap-1.5 pointer-events-auto bg-card/90 backdrop-blur-md border border-border/80 p-1.5 rounded-xl shadow-xl">
          {/* Base Layer Switcher */}
          <div className="flex items-center gap-1 bg-[#040806]/80 p-1 rounded-lg border border-border/50">
            {MAP_PROVIDERS.map((p) => (
              <button
                key={p.id}
                onClick={() => setActiveProviderId(p.id)}
                className={`px-2 py-1 rounded text-[10px] font-mono transition-colors flex items-center gap-1 ${
                  activeProviderId === p.id
                    ? "bg-aurora text-black font-bold shadow-[0_0_8px_#00F5A0]"
                    : "text-text-muted hover:text-white"
                }`}
                title={p.attribution}
              >
                {p.id.includes("satellite") && "🛰️"}
                {p.id.includes("clarity") && "🔭"}
                {p.id.includes("dark") && "🌃"}
                {p.id.includes("topo") && "🗺️"}
                {p.id.includes("google") && "🇬"}
                <span className="hidden sm:inline">{p.badge}</span>
              </button>
            ))}
          </div>

          {/* Preset Buttons Dropdown */}
          <div className="hidden md:flex items-center gap-1 pl-1 border-l border-border/60">
            {MAP_FOCUS_PRESETS.slice(0, 4).map((preset) => (
              <button
                key={preset.id}
                onClick={() => handlePresetSelect(preset)}
                className="px-2 py-1 rounded text-[10px] font-mono text-text-subtle hover:text-aurora hover:bg-surface transition-colors"
                title={preset.description}
              >
                {preset.name}
              </button>
            ))}
          </div>
        </div>

        {/* Right: Layer Toggles & Utility Controls */}
        <div className="flex items-center gap-1.5 pointer-events-auto bg-card/90 backdrop-blur-md border border-border/80 p-1.5 rounded-xl shadow-xl">
          <button
            onClick={() => setShowTradeRoutes(!showTradeRoutes)}
            className={`px-2 py-1 rounded text-[10px] font-mono flex items-center gap-1 transition-colors ${
              showTradeRoutes ? "bg-aurora/20 text-aurora border border-aurora/40" : "text-text-subtle hover:text-white"
            }`}
            title="Toggle Global Maritime & Continental Trade Routes"
          >
            <Ship className="h-3 w-3" />
            <span className="hidden sm:inline">Trade Routes</span>
          </button>

          <button
            onClick={() => setShowConflictMarkers(!showConflictMarkers)}
            className={`px-2 py-1 rounded text-[10px] font-mono flex items-center gap-1 transition-colors ${
              showConflictMarkers ? "bg-red-500/20 text-red-400 border border-red-500/40" : "text-text-subtle hover:text-white"
            }`}
            title="Toggle Conflict & Regulatory Friction Markers"
          >
            <AlertTriangle className="h-3 w-3" />
            <span className="hidden sm:inline">Choke Points</span>
          </button>

          <button
            onClick={() => setShowOpportunityZones(!showOpportunityZones)}
            className={`px-2 py-1 rounded text-[10px] font-mono flex items-center gap-1 transition-colors ${
              showOpportunityZones ? "bg-amber-500/20 text-amber-300 border border-amber-500/40" : "text-text-subtle hover:text-white"
            }`}
            title="Toggle Strategic Opportunity Zones"
          >
            <Sparkles className="h-3 w-3" />
            <span className="hidden sm:inline">Opportunity Zones</span>
          </button>

          <button
            onClick={() => setShowCorridors(!showCorridors)}
            className={`px-2 py-1 rounded text-[10px] font-mono flex items-center gap-1 transition-colors ${
              showCorridors ? "bg-sky-500/20 text-sky-300 border border-sky-500/40" : "text-text-subtle hover:text-white"
            }`}
            title="Toggle Clean Transmission Corridors"
          >
            <Zap className="h-3 w-3" />
            <span className="hidden sm:inline">Corridors</span>
          </button>

          {/* Google Maps API Key Modal Trigger */}
          <button
            onClick={() => setShowKeyModal(!showKeyModal)}
            className="p-1.5 rounded-lg bg-[#040806] border border-border text-text-muted hover:text-aurora transition-colors"
            title="Configure Google Maps API Key"
          >
            <KeyRound className="h-3.5 w-3.5" />
          </button>

          {/* Fullscreen Button */}
          <button
            onClick={() => setIsFullscreen(!isFullscreen)}
            className="p-1.5 rounded-lg bg-[#040806] border border-border text-text-muted hover:text-aurora transition-colors"
            title={isFullscreen ? "Exit Fullscreen" : "Fullscreen Map"}
          >
            {isFullscreen ? <Minimize2 className="h-3.5 w-3.5" /> : <Maximize2 className="h-3.5 w-3.5" />}
          </button>
        </div>
      </div>

      {/* Main Map Container */}
      <div ref={mapContainerRef} className="w-full h-full z-0 select-none" />

      {/* Bottom Telemetry HUD Bar */}
      <div className="absolute bottom-3 left-3 z-[1000] pointer-events-none">
        <div className="pointer-events-auto bg-[#040806]/85 backdrop-blur-md border border-border/80 px-3 py-1.5 rounded-xl text-[10px] font-mono text-text-subtle flex items-center gap-3 shadow-xl">
          <div className="flex items-center gap-1 text-aurora">
            <Navigation className="h-3 w-3 animate-pulse" />
            <span>GEO-RADAR ACTIVE</span>
          </div>
          {cursorCoords ? (
            <div>
              LAT: <span className="text-white">{cursorCoords.lat}°</span> | LNG:{" "}
              <span className="text-white">{cursorCoords.lng}°</span>
            </div>
          ) : (
            <div>SCANNING GEOGRAPHY</div>
          )}
          <div className="hidden sm:inline text-text-muted">
            ZOOM: <span className="text-aurora">{currentZoom}x</span>
          </div>
          <div className="hidden md:inline text-text-subtle/80">
            PROJECTION: <span className="text-white">EPSG:3857 (SPHERICAL MERCATOR)</span>
          </div>
        </div>
      </div>

      {/* Dynamic Detail Flyout: Selected Trade Route */}
      {selectedRoute && (
        <div className="absolute top-16 right-3 z-[1000] w-80 bg-[#0C1812]/95 backdrop-blur-md border border-aurora/40 rounded-xl p-4 text-xs font-mono shadow-2xl space-y-2.5 animate-in fade-in slide-in-from-right-4">
          <div className="flex items-center justify-between border-b border-border/60 pb-2">
            <div className="text-[10px] font-bold text-aurora uppercase tracking-wider flex items-center gap-1.5">
              <Ship className="h-3.5 w-3.5" />
              Strategic Maritime Trade Route
            </div>
            <button
              onClick={() => setSelectedRoute(null)}
              className="text-text-subtle hover:text-white text-sm"
            >
              ✕
            </button>
          </div>
          <div className="font-bold text-sm text-text-main">{selectedRoute.name}</div>
          <div className="text-text-muted text-[11px] leading-relaxed">
            {selectedRoute.strategicSignificance}
          </div>
          <div className="space-y-1 bg-[#040806]/80 p-2 rounded-lg border border-border/50 text-[11px]">
            <div>
              <span className="text-text-subtle">Origin:</span>{" "}
              <span className="text-white">{selectedRoute.origin}</span>
            </div>
            <div>
              <span className="text-text-subtle">Destination:</span>{" "}
              <span className="text-white">{selectedRoute.destination}</span>
            </div>
            <div>
              <span className="text-text-subtle">Velocity:</span>{" "}
              <span className="text-aurora font-semibold">{selectedRoute.annualVelocity}</span>
            </div>
          </div>
          <div className="text-[10px] text-text-subtle">
            <span className="font-bold uppercase text-white">Commodities:</span>{" "}
            {selectedRoute.primaryCommodities.join(", ")}
          </div>
        </div>
      )}

      {/* Dynamic Detail Flyout: Selected Conflict / Choke Point */}
      {selectedConflict && (
        <div className="absolute top-16 right-3 z-[1000] w-80 bg-[#180909]/95 backdrop-blur-md border border-red-500/50 rounded-xl p-4 text-xs font-mono shadow-2xl space-y-2.5 animate-in fade-in slide-in-from-right-4">
          <div className="flex items-center justify-between border-b border-red-900/60 pb-2">
            <div className="text-[10px] font-bold text-red-400 uppercase tracking-wider flex items-center gap-1.5">
              <AlertTriangle className="h-3.5 w-3.5" />
              Choke Point / Regulatory Friction
            </div>
            <button
              onClick={() => setSelectedConflict(null)}
              className="text-text-subtle hover:text-white text-sm"
            >
              ✕
            </button>
          </div>
          <div className="font-bold text-sm text-red-100">{selectedConflict.name}</div>
          <div className="inline-block px-2 py-0.5 rounded text-[10px] font-bold bg-red-900/40 text-red-300 border border-red-700/50">
            SEVERITY: {selectedConflict.severity}
          </div>
          <div className="text-text-muted text-[11px] leading-relaxed">
            {selectedConflict.impactSummary}
          </div>
          <div className="space-y-1 bg-[#040806]/80 p-2 rounded-lg border border-red-900/40 text-[11px]">
            <div>
              <span className="text-text-subtle">Jurisdiction:</span>{" "}
              <span className="text-white">{selectedConflict.jurisdiction}</span>
            </div>
            <div>
              <span className="text-text-subtle">Regulatory Status:</span>{" "}
              <span className="text-amber-300">{selectedConflict.status}</span>
            </div>
          </div>
        </div>
      )}

      {/* Dynamic Detail Flyout: Selected Opportunity Zone */}
      {selectedZone && (
        <div className="absolute top-16 right-3 z-[1000] w-80 bg-[#141208]/95 backdrop-blur-md border border-amber-500/50 rounded-xl p-4 text-xs font-mono shadow-2xl space-y-2.5 animate-in fade-in slide-in-from-right-4">
          <div className="flex items-center justify-between border-b border-amber-900/60 pb-2">
            <div className="text-[10px] font-bold text-amber-400 uppercase tracking-wider flex items-center gap-1.5">
              <Sparkles className="h-3.5 w-3.5" />
              Strategic Opportunity Zone
            </div>
            <button
              onClick={() => setSelectedZone(null)}
              className="text-text-subtle hover:text-white text-sm"
            >
              ✕
            </button>
          </div>
          <div className="font-bold text-sm text-amber-100">{selectedZone.name}</div>
          <div className="inline-block px-2 py-0.5 rounded text-[10px] font-bold bg-amber-900/40 text-amber-300 border border-amber-700/50">
            EST. ENDOWMENT: {selectedZone.estimatedEndowmentCAD}
          </div>
          <div className="text-text-muted text-[11px] leading-relaxed">
            {selectedZone.description}
          </div>
          <div className="text-[10px] text-text-subtle">
            <span className="font-bold uppercase text-white">Target Resources:</span>{" "}
            {selectedZone.criticalMineralsOrEnergy.join(", ")}
          </div>
        </div>
      )}

      {/* Google Maps API Key Modal */}
      {showKeyModal && (
        <div className="absolute inset-0 z-[2000] bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-[#0C1812] border border-border/90 rounded-2xl max-w-md w-full p-6 space-y-4 shadow-2xl font-mono">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2 text-aurora font-bold text-sm">
                <KeyRound className="h-4 w-4" />
                Google Maps Free API Key Integration
              </div>
              <button
                onClick={() => setShowKeyModal(false)}
                className="text-text-subtle hover:text-white text-sm"
              >
                ✕
              </button>
            </div>

            <p className="text-xs text-text-muted leading-relaxed">
              Google Maps provides a monthly $200 recurring credit (covering ~28,000 map loads/month for free).
              Enter your Google Maps API key below to enable direct Google Satellite & Photogrammetry tiles, or continue using our zero-credential high-resolution Esri World Imagery.
            </p>

            <div className="space-y-1.5">
              <label className="text-[11px] text-text-subtle">API Key (Optional)</label>
              <input
                type="text"
                value={googleApiKey}
                onChange={(e) => setGoogleApiKey(e.target.value)}
                placeholder="AIzaSy..."
                className="w-full bg-[#040806] border border-border rounded-xl px-3 py-2 text-xs text-white placeholder-text-subtle/50 focus:outline-none focus:border-primary"
              />
            </div>

            <div className="flex items-center justify-between pt-2">
              <span className="text-[10px] text-text-subtle">
                Key saved in local session memory
              </span>
              <button
                onClick={() => {
                  if (googleApiKey) {
                    setActiveProviderId("google-hybrid");
                  }
                  setShowKeyModal(false);
                }}
                className="px-4 py-2 rounded-xl bg-aurora text-black font-bold text-xs hover:bg-aurora/90 transition-colors"
              >
                Apply Key
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
