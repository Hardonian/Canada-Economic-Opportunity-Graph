"use client";

import React, { useEffect, useRef, useState, useCallback } from "react";
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
  Info,
  Ruler,
  Search,
  Download,
  Trash2,
  FileSpreadsheet,
  FileCode,
  Crosshair,
  Radio,
  Landmark,
  Cpu,
  Factory,
} from "lucide-react";
import type { Project } from "@/lib/types";
import {
  MAP_PROVIDERS,
  GLOBAL_TRADE_ROUTES,
  CONFLICT_MARKERS,
  OPPORTUNITY_ZONES,
  STRATEGIC_CORRIDORS,
  MAP_FOCUS_PRESETS,
  CRITICAL_MINERAL_HUBS,
  TREATY_TERRITORIES,
  GRID_INTERTIE_ZONES,
  type MapTileProvider,
  type TradeRoute,
  type ConflictMarker,
  type OpportunityZone,
  type CriticalMineralHub,
  type TreatyTerritory,
  type GridIntertieZone,
} from "@/lib/geospatial";

interface GeospatialMapProps {
  projects: Project[];
  activeProject: Project;
  onSelectProject: (p: Project) => void;
  selectedSector: string;
  selectedStage: string;
}

// Great circle Haversine formula for calculating geodesic distance in kilometers
function haversineDistanceKm(lat1: number, lon1: number, lat2: number, lon2: number): number {
  const R = 6371; // Earth's mean radius in km
  const dLat = ((lat2 - lat1) * Math.PI) / 180;
  const dLon = ((lon2 - lon1) * Math.PI) / 180;
  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos((lat1 * Math.PI) / 180) *
      Math.cos((lat2 * Math.PI) / 180) *
      Math.sin(dLon / 2) *
      Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
  return R * c;
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
    measurement: null,
    mineralHubs: null,
    treatyTerritories: null,
    gridInterties: null,
    buffer: null,
  });

  // State
  const [activeProviderId, setActiveProviderId] = useState<string>("esri-satellite");
  const [showTradeRoutes, setShowTradeRoutes] = useState<boolean>(true);
  const [showConflictMarkers, setShowConflictMarkers] = useState<boolean>(true);
  const [showOpportunityZones, setShowOpportunityZones] = useState<boolean>(true);
  const [showCorridors, setShowCorridors] = useState<boolean>(true);
  const [showMineralHubs, setShowMineralHubs] = useState<boolean>(true);
  const [showTreatyTerritories, setShowTreatyTerritories] = useState<boolean>(true);
  const [showGridInterties, setShowGridInterties] = useState<boolean>(true);
  const [cursorCoords, setCursorCoords] = useState<{ lat: number; lng: number } | null>(null);
  const [currentZoom, setCurrentZoom] = useState<number>(4);
  const [isFullscreen, setIsFullscreen] = useState<boolean>(false);
  const [selectedRoute, setSelectedRoute] = useState<TradeRoute | null>(null);
  const [selectedConflict, setSelectedConflict] = useState<ConflictMarker | null>(null);
  const [selectedZone, setSelectedZone] = useState<OpportunityZone | null>(null);
  const [selectedMineralHub, setSelectedMineralHub] = useState<CriticalMineralHub | null>(null);
  const [selectedTreaty, setSelectedTreaty] = useState<TreatyTerritory | null>(null);
  const [selectedIntertie, setSelectedIntertie] = useState<GridIntertieZone | null>(null);
  const [googleApiKey, setGoogleApiKey] = useState<string>("");
  const [showKeyModal, setShowKeyModal] = useState<boolean>(false);
  const [isLeafletReady, setIsLeafletReady] = useState<boolean>(false);

  // Advanced Geodesic Measure Tool State
  const [isMeasuring, setIsMeasuring] = useState<boolean>(false);
  const [measurePoints, setMeasurePoints] = useState<{ lat: number; lng: number }[]>([]);
  const isMeasuringRef = useRef<boolean>(false);
  const measurePointsRef = useRef<{ lat: number; lng: number }[]>([]);

  // Corridor & Infrastructure Proximity Buffer Tool State (25km / 50km / 100km)
  const [isBufferActive, setIsBufferActive] = useState<boolean>(false);
  const [bufferRadiusKm, setBufferRadiusKm] = useState<number>(50);
  const [bufferCenter, setBufferCenter] = useState<{ lat: number; lng: number } | null>(null);
  const isBufferActiveRef = useRef<boolean>(false);
  const bufferRadiusKmRef = useRef<number>(50);
  const bufferCenterRef = useRef<{ lat: number; lng: number } | null>(null);

  // Instant Search Autocomplete State
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [isSearchOpen, setIsSearchOpen] = useState<boolean>(false);
  const searchContainerRef = useRef<HTMLDivElement>(null);

  // Spatial Export Menu State
  const [showExportMenu, setShowExportMenu] = useState<boolean>(false);
  const exportContainerRef = useRef<HTMLDivElement>(null);

  // Sync refs with state for Leaflet event handlers
  useEffect(() => {
    isMeasuringRef.current = isMeasuring;
  }, [isMeasuring]);

  useEffect(() => {
    measurePointsRef.current = measurePoints;
  }, [measurePoints]);

  useEffect(() => {
    isBufferActiveRef.current = isBufferActive;
  }, [isBufferActive]);

  useEffect(() => {
    bufferRadiusKmRef.current = bufferRadiusKm;
  }, [bufferRadiusKm]);

  useEffect(() => {
    bufferCenterRef.current = bufferCenter;
  }, [bufferCenter]);

  // Close search/export dropdowns on outside click
  useEffect(() => {
    const handleOutsideClick = (e: MouseEvent) => {
      if (searchContainerRef.current && !searchContainerRef.current.contains(e.target as Node)) {
        setIsSearchOpen(false);
      }
      if (exportContainerRef.current && !exportContainerRef.current.contains(e.target as Node)) {
        setShowExportMenu(false);
      }
    };
    document.addEventListener("mousedown", handleOutsideClick);
    return () => document.removeEventListener("mousedown", handleOutsideClick);
  }, []);

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

      // Handle map click for buffer tool or measurement tool
      map.on("click", (e: any) => {
        if (isBufferActiveRef.current) {
          const newCenter = {
            lat: parseFloat(e.latlng.lat.toFixed(4)),
            lng: parseFloat(e.latlng.lng.toFixed(4)),
          };
          bufferCenterRef.current = newCenter;
          setBufferCenter(newCenter);
          return;
        }

        if (isMeasuringRef.current) {
          const newPoint = {
            lat: parseFloat(e.latlng.lat.toFixed(4)),
            lng: parseFloat(e.latlng.lng.toFixed(4)),
          };
          const updated = [...measurePointsRef.current, newPoint];
          measurePointsRef.current = updated;
          setMeasurePoints(updated);
        }
      });

      // Layer groups for dynamic filtering
      layersGroupRef.current.projects = L.layerGroup().addTo(map);
      layersGroupRef.current.tradeRoutes = L.layerGroup().addTo(map);
      layersGroupRef.current.conflictMarkers = L.layerGroup().addTo(map);
      layersGroupRef.current.opportunityZones = L.layerGroup().addTo(map);
      layersGroupRef.current.corridors = L.layerGroup().addTo(map);
      layersGroupRef.current.measurement = L.layerGroup().addTo(map);
      layersGroupRef.current.mineralHubs = L.layerGroup().addTo(map);
      layersGroupRef.current.treatyTerritories = L.layerGroup().addTo(map);
      layersGroupRef.current.gridInterties = L.layerGroup().addTo(map);
      layersGroupRef.current.buffer = L.layerGroup().addTo(map);

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
              <div style="width: ${size}px; height: ${size}px; border-radius: 9999px; background: ${color}; border: 2px solid #040806; box-shadow: 0 0 10px ${color}; display: flex; align-items: center; justify-content: center; cursor: pointer;">
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

  // Render Critical Mineral Hubs
  useEffect(() => {
    if (!isLeafletReady || !layersGroupRef.current.mineralHubs) return;

    import("leaflet").then((L) => {
      const group = layersGroupRef.current.mineralHubs;
      group.clearLayers();

      if (!showMineralHubs) return;

      CRITICAL_MINERAL_HUBS.forEach((hub) => {
        const hubIcon = L.divIcon({
          className: "mineral-hub-marker-icon",
          iconSize: [22, 22],
          iconAnchor: [11, 11],
          html: `
            <div style="position: relative; width: 22px; height: 22px; cursor: pointer;">
              <div style="position: absolute; inset: -4px; border-radius: 6px; background: ${hub.color}; opacity: 0.35; transform: rotate(45deg);" class="animate-pulse"></div>
              <div style="width: 22px; height: 22px; background: #1a0b2e; border: 2px solid ${hub.color}; border-radius: 4px; display: flex; align-items: center; justify-content: center; transform: rotate(45deg); box-shadow: 0 0 10px ${hub.color};">
                <div style="transform: rotate(-45deg); font-size: 10px; font-weight: bold; color: ${hub.color};">◆</div>
              </div>
            </div>
          `,
        });

        const marker = L.marker([hub.location.lat, hub.location.lng], { icon: hubIcon });

        marker.bindTooltip(
          `<b>${hub.name}</b><br/><span style="color:${hub.color}">${hub.capacityMetric}</span><br/><span style="color:#00F5A0">DVRI: ${(hub.domesticRetentionRate * 100).toFixed(0)}%</span>`,
          {
            direction: "top",
            className: "trade-route-tooltip",
          }
        );

        marker.on("click", () => {
          setSelectedMineralHub(hub);
          setSelectedRoute(null);
          setSelectedConflict(null);
          setSelectedZone(null);
          setSelectedTreaty(null);
          setSelectedIntertie(null);
        });

        group.addLayer(marker);
      });
    });
  }, [isLeafletReady, showMineralHubs]);

  // Render Treaty Territories
  useEffect(() => {
    if (!isLeafletReady || !layersGroupRef.current.treatyTerritories) return;

    import("leaflet").then((L) => {
      const group = layersGroupRef.current.treatyTerritories;
      group.clearLayers();

      if (!showTreatyTerritories) return;

      TREATY_TERRITORIES.forEach((t) => {
        // Broad regional jurisdictional aura
        const circle = L.circle([t.center.lat, t.center.lng], {
          radius: 110000,
          color: t.color,
          fillColor: t.color,
          fillOpacity: 0.1,
          weight: 1.5,
          dashArray: "6, 4",
        });

        circle.on("click", () => {
          setSelectedTreaty(t);
          setSelectedRoute(null);
          setSelectedConflict(null);
          setSelectedZone(null);
          setSelectedMineralHub(null);
          setSelectedIntertie(null);
        });

        group.addLayer(circle);

        // Center Emblem Pin
        const treatyIcon = L.divIcon({
          className: "treaty-marker-icon",
          iconSize: [20, 20],
          iconAnchor: [10, 10],
          html: `
            <div style="position: relative; width: 20px; height: 20px; cursor: pointer;">
              <div style="width: 20px; height: 20px; border-radius: 9999px; background: #061c12; border: 2px solid ${t.color}; display: flex; align-items: center; justify-content: center; box-shadow: 0 0 8px ${t.color}; font-size: 9px; font-weight: bold; color: ${t.color};">
                ⚖
              </div>
            </div>
          `,
        });

        const marker = L.marker([t.center.lat, t.center.lng], { icon: treatyIcon });
        marker.bindTooltip(
          `<b>${t.name}</b><br/><span style="color:${t.color}">${t.historicalFramework}</span><br/><span style="color:#00F5A0">${t.loanGuaranteeEligibility}</span>`,
          {
            direction: "top",
            className: "opportunity-tooltip",
          }
        );

        marker.on("click", () => {
          setSelectedTreaty(t);
          setSelectedRoute(null);
          setSelectedConflict(null);
          setSelectedZone(null);
          setSelectedMineralHub(null);
          setSelectedIntertie(null);
        });

        group.addLayer(marker);
      });
    });
  }, [isLeafletReady, showTreatyTerritories]);

  // Render Grid Intertie Zones & AI Compute Hubs
  useEffect(() => {
    if (!isLeafletReady || !layersGroupRef.current.gridInterties) return;

    import("leaflet").then((L) => {
      const group = layersGroupRef.current.gridInterties;
      group.clearLayers();

      if (!showGridInterties) return;

      GRID_INTERTIE_ZONES.forEach((z) => {
        // Clean power basin aura
        const circle = L.circle([z.center.lat, z.center.lng], {
          radius: 50000,
          color: z.color,
          fillColor: z.color,
          fillOpacity: 0.12,
          weight: 2,
          dashArray: "4, 4",
        });

        circle.on("click", () => {
          setSelectedIntertie(z);
          setSelectedRoute(null);
          setSelectedConflict(null);
          setSelectedZone(null);
          setSelectedMineralHub(null);
          setSelectedTreaty(null);
        });

        group.addLayer(circle);

        // Power Center Marker
        const powerIcon = L.divIcon({
          className: "power-intertie-icon",
          iconSize: [22, 22],
          iconAnchor: [11, 11],
          html: `
            <div style="position: relative; width: 22px; height: 22px; cursor: pointer;">
              <div style="position: absolute; inset: -4px; border-radius: 9999px; background: ${z.color}; opacity: 0.3;" class="animate-ping"></div>
              <div style="width: 22px; height: 22px; border-radius: 9999px; background: #071926; border: 2px solid ${z.color}; display: flex; align-items: center; justify-content: center; box-shadow: 0 0 10px ${z.color}; font-size: 11px; font-weight: bold; color: ${z.color};">
                ⚡
              </div>
            </div>
          `,
        });

        const marker = L.marker([z.center.lat, z.center.lng], { icon: powerIcon });
        marker.bindTooltip(
          `<b>${z.name}</b><br/><span style="color:${z.color}">${z.gridOperator} (${z.cleanCapacityMW} MW)</span><br/><span style="color:#38BDF8">AI Headroom: ${z.aiHeadroomMW} MW</span>`,
          {
            direction: "top",
            className: "trade-route-tooltip",
          }
        );

        marker.on("click", () => {
          setSelectedIntertie(z);
          setSelectedRoute(null);
          setSelectedConflict(null);
          setSelectedZone(null);
          setSelectedMineralHub(null);
          setSelectedTreaty(null);
        });

        group.addLayer(marker);
      });
    });
  }, [isLeafletReady, showGridInterties]);

  // Render Infrastructure Proximity Buffer Layer
  useEffect(() => {
    if (!isLeafletReady || !layersGroupRef.current.buffer) return;

    import("leaflet").then((L) => {
      const group = layersGroupRef.current.buffer;
      group.clearLayers();

      if (!isBufferActive || !bufferCenter) return;

      // 1. Buffer Radius Circle
      const circle = L.circle([bufferCenter.lat, bufferCenter.lng], {
        radius: bufferRadiusKm * 1000,
        color: "#00F5A0",
        fillColor: "#00F5A0",
        fillOpacity: 0.16,
        weight: 2.5,
        dashArray: "6, 6",
      });
      group.addLayer(circle);

      // 2. Buffer Center Marker
      const centerIcon = L.divIcon({
        className: "buffer-center-icon",
        iconSize: [24, 24],
        iconAnchor: [12, 12],
        html: `
          <div style="position: relative; width: 24px; height: 24px;">
            <div style="position: absolute; inset: -4px; border-radius: 9999px; background: #00F5A0; opacity: 0.5;" class="animate-ping"></div>
            <div style="width: 24px; height: 24px; border-radius: 9999px; background: #040806; border: 2px solid #00F5A0; display: flex; align-items: center; justify-content: center; box-shadow: 0 0 12px #00F5A0;">
              <div style="width: 6px; height: 6px; border-radius: 9999px; background: #00F5A0;"></div>
            </div>
          </div>
        `,
      });

      const marker = L.marker([bufferCenter.lat, bufferCenter.lng], { icon: centerIcon });
      marker.bindTooltip(
        `<b>Proximity Buffer Epicenter</b><br/>Radius: ${bufferRadiusKm} km`,
        {
          permanent: true,
          direction: "top",
          className: "trade-route-tooltip",
        }
      );
      group.addLayer(marker);
    });
  }, [isLeafletReady, isBufferActive, bufferCenter, bufferRadiusKm]);

  // Render Geodesic Measurement Overlay
  useEffect(() => {
    if (!isLeafletReady || !layersGroupRef.current.measurement) return;

    import("leaflet").then((L) => {
      const group = layersGroupRef.current.measurement;
      group.clearLayers();

      if (measurePoints.length === 0) return;

      // Draw polyline connecting measurement points
      if (measurePoints.length > 1) {
        const latLngs = measurePoints.map((p) => [p.lat, p.lng]);
        const polyline = L.polyline(latLngs as any, {
          color: "#F59E0B",
          weight: 3,
          dashArray: "6, 6",
          opacity: 0.9,
        });
        group.addLayer(polyline);
      }

      // Draw waypoint pins
      measurePoints.forEach((p, idx) => {
        const isEndpoint = idx === measurePoints.length - 1;
        const icon = L.divIcon({
          className: "measurement-waypoint-icon",
          iconSize: [22, 22],
          iconAnchor: [11, 11],
          html: `
            <div style="width: 22px; height: 22px; border-radius: 9999px; background: ${
              isEndpoint ? "#00F5A0" : "#F59E0B"
            }; border: 2px solid #040806; display: flex; align-items: center; justify-content: center; font-weight: bold; font-size: 10px; color: #040806; font-family: monospace; box-shadow: 0 0 8px rgba(0,0,0,0.8);">
              ${idx + 1}
            </div>
          `,
        });

        const marker = L.marker([p.lat, p.lng], { icon });

        // Calculate distance from previous point if applicable
        if (idx > 0) {
          const prev = measurePoints[idx - 1];
          const legDist = haversineDistanceKm(prev.lat, prev.lng, p.lat, p.lng);
          marker.bindTooltip(
            `<b>Leg ${idx}</b>: ${legDist.toFixed(1)} km (${(legDist / 1.852).toFixed(1)} NM)`,
            { permanent: true, direction: "top", className: "trade-route-tooltip" }
          );
        }

        group.addLayer(marker);
      });
    });
  }, [isLeafletReady, measurePoints]);

  // Compute Cumulative Geodesic Distance
  const totalMeasureDistanceKm = React.useMemo(() => {
    if (measurePoints.length < 2) return 0;
    let sum = 0;
    for (let i = 1; i < measurePoints.length; i++) {
      sum += haversineDistanceKm(
        measurePoints[i - 1].lat,
        measurePoints[i - 1].lng,
        measurePoints[i].lat,
        measurePoints[i].lng
      );
    }
    return sum;
  }, [measurePoints]);

  const clearMeasurement = useCallback(() => {
    setMeasurePoints([]);
    measurePointsRef.current = [];
  }, []);

  // Infrastructure Proximity Buffer Analytics
  const bufferIntersectingProjects = React.useMemo(() => {
    if (!bufferCenter) return [];
    return projects.filter((p) => {
      if (typeof p.latitude !== "number" || typeof p.longitude !== "number") return false;
      const dist = haversineDistanceKm(bufferCenter.lat, bufferCenter.lng, p.latitude, p.longitude);
      return dist <= bufferRadiusKm;
    });
  }, [bufferCenter, bufferRadiusKm, projects]);

  const bufferTotalCapex = React.useMemo(() => {
    return bufferIntersectingProjects.reduce((sum, p) => sum + (p.capex_cad || 0), 0);
  }, [bufferIntersectingProjects]);

  const clearBuffer = useCallback(() => {
    setBufferCenter(null);
    bufferCenterRef.current = null;
  }, []);

  // Search Results
  const searchResults = React.useMemo(() => {
    if (!searchQuery.trim()) return [];
    const q = searchQuery.toLowerCase().trim();
    return projects
      .filter(
        (p) =>
          p.name.toLowerCase().includes(q) ||
          p.province.toLowerCase().includes(q) ||
          p.sector.toLowerCase().includes(q)
      )
      .slice(0, 6);
  }, [projects, searchQuery]);

  // Spatial Export Handlers
  const handleExportGeoJSON = useCallback(() => {
    const validProjects = projects.filter(
      (p) => typeof p.latitude === "number" && typeof p.longitude === "number"
    );

    const geojson = {
      type: "FeatureCollection",
      crs: {
        type: "name",
        properties: { name: "urn:ogc:def:crs:OGC:1.3:CRS84" },
      },
      metadata: {
        title: "Canada Economic Opportunity Graph — Active Geospatial Layer",
        export_timestamp: new Date().toISOString(),
        total_features: validProjects.length,
        authority: "Government of Canada / Gouvernement du Canada (CEGS Standard)",
      },
      features: validProjects.map((p) => ({
        type: "Feature",
        geometry: {
          type: "Point",
          coordinates: [p.longitude, p.latitude],
        },
        properties: {
          id: p.id,
          name: p.name,
          slug: p.slug,
          sector: p.sector,
          province: p.province,
          current_stage: p.current_stage,
          capex_cad: p.capex_cad,
          proponent_name: p.proponent_name,
        },
      })),
    };

    const blob = new Blob([JSON.stringify(geojson, null, 2)], { type: "application/geo+json" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `canada_economic_graph_spatial_${new Date().toISOString().slice(0, 10)}.geojson`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }, [projects]);

  const handleExportCSV = useCallback(() => {
    const validProjects = projects.filter(
      (p) => typeof p.latitude === "number" && typeof p.longitude === "number"
    );

    const headers = [
      "id",
      "name",
      "slug",
      "sector",
      "province",
      "current_stage",
      "capex_cad",
      "latitude",
      "longitude",
      "proponent_name",
    ];

    const rows = validProjects.map((p) => [
      `"${p.id}"`,
      `"${p.name.replace(/"/g, '""')}"`,
      `"${p.slug}"`,
      `"${p.sector}"`,
      `"${p.province}"`,
      `"${p.current_stage}"`,
      p.capex_cad,
      p.latitude,
      p.longitude,
      `"${(p.proponent_name || "").replace(/"/g, '""')}"`,
    ]);

    const csvContent = [headers.join(","), ...rows.map((r) => r.join(","))].join("\n");
    const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `canada_economic_graph_projects_${new Date().toISOString().slice(0, 10)}.csv`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }, [projects]);

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
      } ${isMeasuring ? "cursor-crosshair" : ""}`}
    >
      {/* Top Interactive HUD Bar */}
      <div className="absolute top-3 left-3 right-3 z-[1000] flex flex-wrap items-center justify-between gap-2 pointer-events-none">
        {/* Left: Base Map, Presets & Instant Search */}
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

          {/* Preset Buttons */}
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

          {/* Instant Search Autocomplete Bar */}
          <div ref={searchContainerRef} className="relative pl-1 border-l border-border/60">
            <div className="flex items-center bg-[#040806] border border-border/80 rounded-lg px-2 py-1 text-xs">
              <Search className="h-3.5 w-3.5 text-text-subtle mr-1.5 shrink-0" />
              <input
                type="text"
                placeholder="Search map..."
                value={searchQuery}
                onChange={(e) => {
                  setSearchQuery(e.target.value);
                  setIsSearchOpen(true);
                }}
                onFocus={() => setIsSearchOpen(true)}
                className="bg-transparent border-none outline-none text-white text-[11px] w-24 sm:w-32 placeholder-text-subtle/50 font-mono"
              />
              {searchQuery && (
                <button
                  onClick={() => {
                    setSearchQuery("");
                    setIsSearchOpen(false);
                  }}
                  className="text-text-subtle hover:text-white text-[10px] ml-1"
                >
                  ✕
                </button>
              )}
            </div>

            {/* Search Dropdown Results */}
            {isSearchOpen && searchResults.length > 0 && (
              <div className="absolute top-full mt-1.5 left-0 w-72 bg-[#0C1812]/98 backdrop-blur-md border border-aurora/50 rounded-xl p-1.5 shadow-2xl z-[1200] font-mono text-xs space-y-1 max-h-72 overflow-y-auto">
                <div className="px-2 py-1 text-[9px] uppercase tracking-wider text-text-subtle font-bold border-b border-border/50">
                  Matching Capital Projects ({searchResults.length})
                </div>
                {searchResults.map((p) => (
                  <button
                    key={p.id}
                    onClick={() => {
                      if (typeof p.latitude === "number" && typeof p.longitude === "number") {
                        mapInstanceRef.current?.flyTo([p.latitude, p.longitude], 9, {
                          duration: 1.2,
                        });
                      }
                      onSelectProject(p);
                      setIsSearchOpen(false);
                    }}
                    className="w-full text-left p-2 rounded-lg hover:bg-surface transition-colors space-y-0.5 group"
                  >
                    <div className="font-bold text-[11px] text-text-main group-hover:text-aurora truncate">
                      {p.name}
                    </div>
                    <div className="flex items-center justify-between text-[10px] text-text-subtle">
                      <span>
                        {p.sector} • {p.province}
                      </span>
                      <span className="text-aurora font-semibold">
                        {p.capex_cad >= 1e9
                          ? `$${(p.capex_cad / 1e9).toFixed(1)}B`
                          : p.capex_cad > 0
                          ? `$${(p.capex_cad / 1e6).toFixed(0)}M`
                          : "N/R"}
                      </span>
                    </div>
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>

        {/* Right: Layer Toggles, Measure, Export & Utility Controls */}
        <div className="flex items-center gap-1.5 pointer-events-auto bg-card/90 backdrop-blur-md border border-border/80 p-1.5 rounded-xl shadow-xl">
          <button
            onClick={() => setShowTradeRoutes(!showTradeRoutes)}
            className={`px-2 py-1 rounded text-[10px] font-mono flex items-center gap-1 transition-colors ${
              showTradeRoutes
                ? "bg-aurora/20 text-aurora border border-aurora/40"
                : "text-text-subtle hover:text-white"
            }`}
            title="Toggle Global Maritime & Continental Trade Routes"
          >
            <Ship className="h-3 w-3" />
            <span className="hidden sm:inline">Trade Routes</span>
          </button>

          <button
            onClick={() => setShowConflictMarkers(!showConflictMarkers)}
            className={`px-2 py-1 rounded text-[10px] font-mono flex items-center gap-1 transition-colors ${
              showConflictMarkers
                ? "bg-red-500/20 text-red-400 border border-red-500/40"
                : "text-text-subtle hover:text-white"
            }`}
            title="Toggle Conflict & Regulatory Friction Markers"
          >
            <AlertTriangle className="h-3 w-3" />
            <span className="hidden sm:inline">Choke Points</span>
          </button>

          <button
            onClick={() => setShowOpportunityZones(!showOpportunityZones)}
            className={`px-2 py-1 rounded text-[10px] font-mono flex items-center gap-1 transition-colors ${
              showOpportunityZones
                ? "bg-amber-500/20 text-amber-300 border border-amber-500/40"
                : "text-text-subtle hover:text-white"
            }`}
            title="Toggle Strategic Opportunity Zones"
          >
            <Sparkles className="h-3 w-3" />
            <span className="hidden sm:inline">Opportunity Zones</span>
          </button>

          <button
            onClick={() => setShowCorridors(!showCorridors)}
            className={`px-2 py-1 rounded text-[10px] font-mono flex items-center gap-1 transition-colors ${
              showCorridors
                ? "bg-sky-500/20 text-sky-300 border border-sky-500/40"
                : "text-text-subtle hover:text-white"
            }`}
            title="Toggle Clean Transmission Corridors"
          >
            <Zap className="h-3 w-3" />
            <span className="hidden sm:inline">Corridors</span>
          </button>

          <button
            onClick={() => setShowMineralHubs(!showMineralHubs)}
            className={`px-2 py-1 rounded text-[10px] font-mono flex items-center gap-1 transition-colors ${
              showMineralHubs
                ? "bg-purple-500/20 text-purple-300 border border-purple-500/40"
                : "text-text-subtle hover:text-white"
            }`}
            title="Toggle Critical Mineral Refining & Processing Hubs"
          >
            <Factory className="h-3 w-3" />
            <span className="hidden sm:inline">Mineral Hubs</span>
          </button>

          <button
            onClick={() => setShowTreatyTerritories(!showTreatyTerritories)}
            className={`px-2 py-1 rounded text-[10px] font-mono flex items-center gap-1 transition-colors ${
              showTreatyTerritories
                ? "bg-emerald-500/20 text-emerald-300 border border-emerald-500/40"
                : "text-text-subtle hover:text-white"
            }`}
            title="Toggle Treaty Territories & $5B Loan Guarantee Jurisdiction"
          >
            <Landmark className="h-3 w-3" />
            <span className="hidden sm:inline">Treaty Lands</span>
          </button>

          <button
            onClick={() => setShowGridInterties(!showGridInterties)}
            className={`px-2 py-1 rounded text-[10px] font-mono flex items-center gap-1 transition-colors ${
              showGridInterties
                ? "bg-blue-500/20 text-blue-300 border border-blue-500/40"
                : "text-text-subtle hover:text-white"
            }`}
            title="Toggle Clean Energy Baseload & AI Interties"
          >
            <Cpu className="h-3 w-3" />
            <span className="hidden sm:inline">Interties</span>
          </button>

          {/* Infrastructure Proximity Buffer Tool (25km / 50km / 100km) */}
          <button
            onClick={() => {
              const next = !isBufferActive;
              setIsBufferActive(next);
              if (!next) {
                clearBuffer();
              }
            }}
            className={`p-1.5 rounded-lg border transition-colors flex items-center gap-1 ${
              isBufferActive
                ? "bg-emerald-500/25 border-aurora text-aurora shadow-[0_0_10px_rgba(0,245,160,0.5)]"
                : "bg-[#040806] border-border text-text-muted hover:text-aurora"
            }`}
            title={isBufferActive ? "Disable Proximity Buffer Tool" : "Enable Infrastructure Proximity Buffer (25-100km)"}
          >
            <Radio className="h-3.5 w-3.5" />
            <span className="hidden lg:inline text-[10px] font-mono font-bold">Buffer Tool</span>
          </button>

          {/* Interactive Geodesic Measure Tool Button */}
          <button
            onClick={() => {
              const next = !isMeasuring;
              setIsMeasuring(next);
              if (!next) {
                clearMeasurement();
              }
            }}
            className={`p-1.5 rounded-lg border transition-colors flex items-center gap-1 ${
              isMeasuring
                ? "bg-amber-500/25 border-amber-500 text-amber-300 shadow-[0_0_10px_rgba(245,158,11,0.5)]"
                : "bg-[#040806] border-border text-text-muted hover:text-aurora"
            }`}
            title={isMeasuring ? "Disable Distance Measure Tool" : "Enable Geodesic Distance Measure Tool"}
          >
            <Ruler className="h-3.5 w-3.5" />
            <span className="hidden lg:inline text-[10px] font-mono">Measure</span>
          </button>

          {/* Spatial Data Export Dropdown */}
          <div ref={exportContainerRef} className="relative">
            <button
              onClick={() => setShowExportMenu(!showExportMenu)}
              className="p-1.5 rounded-lg bg-[#040806] border border-border text-text-muted hover:text-aurora transition-colors flex items-center gap-1"
              title="Export Spatial Geospatial Layer (GeoJSON / CSV)"
            >
              <Download className="h-3.5 w-3.5" />
              <span className="hidden lg:inline text-[10px] font-mono">Export</span>
            </button>

            {showExportMenu && (
              <div className="absolute top-full mt-1.5 right-0 w-52 bg-[#0C1812]/98 backdrop-blur-md border border-border rounded-xl p-1.5 shadow-2xl z-[1200] font-mono text-xs space-y-1">
                <div className="px-2.5 py-1 text-[9px] uppercase tracking-wider text-text-subtle font-bold border-b border-border/50">
                  Download Spatial Layer
                </div>
                <button
                  onClick={() => {
                    handleExportGeoJSON();
                    setShowExportMenu(false);
                  }}
                  className="w-full text-left px-2.5 py-1.5 rounded-lg hover:bg-surface text-text-main hover:text-aurora text-[11px] flex items-center gap-2"
                >
                  <FileCode className="h-3.5 w-3.5 text-aurora" />
                  Export GeoJSON (.geojson)
                </button>
                <button
                  onClick={() => {
                    handleExportCSV();
                    setShowExportMenu(false);
                  }}
                  className="w-full text-left px-2.5 py-1.5 rounded-lg hover:bg-surface text-text-main hover:text-aurora text-[11px] flex items-center gap-2"
                >
                  <FileSpreadsheet className="h-3.5 w-3.5 text-sky-400" />
                  Export Table (.csv)
                </button>
              </div>
            )}
          </div>

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

      {/* Floating Measurement HUD Banner */}
      {isMeasuring && (
        <div className="absolute top-16 left-3 z-[1000] bg-[#040806]/95 backdrop-blur-md border border-amber-500/60 rounded-xl px-4 py-2.5 shadow-2xl flex flex-wrap items-center gap-3 font-mono text-xs text-white animate-in fade-in slide-in-from-top-2">
          <div className="flex items-center gap-1.5 text-amber-400 font-bold">
            <Ruler className="h-4 w-4 animate-pulse" />
            <span>GEODESIC MEASURE</span>
          </div>
          <div className="h-4 w-px bg-border/80 hidden sm:block" />
          <div>
            DISTANCE:{" "}
            <span className="text-aurora font-bold">{totalMeasureDistanceKm.toFixed(1)} km</span>{" "}
            <span className="text-text-muted text-[11px]">
              ({(totalMeasureDistanceKm / 1.852).toFixed(1)} NM)
            </span>
          </div>
          <div className="text-[11px] text-text-subtle hidden sm:block">
            {measurePoints.length === 0
              ? "Click map to set initial waypoint"
              : `${measurePoints.length} waypoint${measurePoints.length > 1 ? "s" : ""}`}
          </div>
          <div className="flex items-center gap-1.5 ml-auto">
            {measurePoints.length > 0 && (
              <button
                onClick={clearMeasurement}
                className="px-2 py-1 rounded bg-surface hover:bg-surface/80 text-text-muted hover:text-white text-[10px] flex items-center gap-1"
                title="Reset points"
              >
                <Trash2 className="h-3 w-3" />
                Clear
              </button>
            )}
            <button
              onClick={() => {
                setIsMeasuring(false);
                clearMeasurement();
              }}
              className="px-2 py-1 rounded bg-amber-500/20 text-amber-300 hover:bg-amber-500/30 text-[10px] font-bold"
            >
              Exit Tool
            </button>
          </div>
        </div>
      )}

      {/* Floating Infrastructure Proximity Buffer HUD Banner */}
      {isBufferActive && (
        <div className="absolute top-16 left-3 z-[1000] bg-[#040806]/95 backdrop-blur-md border border-aurora/60 rounded-xl px-4 py-2.5 shadow-2xl flex flex-wrap items-center gap-3 font-mono text-xs text-white animate-in fade-in slide-in-from-top-2">
          <div className="flex items-center gap-1.5 text-aurora font-bold">
            <Radio className="h-4 w-4 animate-pulse" />
            <span>CORRIDOR PROXIMITY BUFFER</span>
          </div>

          <div className="h-4 w-px bg-border/80 hidden sm:block" />

          {/* Radius selector */}
          <div className="flex items-center gap-1 bg-surface/60 p-0.5 rounded-lg border border-border/60">
            {[25, 50, 100].map((r) => (
              <button
                key={r}
                onClick={() => setBufferRadiusKm(r)}
                className={`px-2 py-0.5 rounded text-[10px] font-mono transition-colors ${
                  bufferRadiusKm === r
                    ? "bg-aurora text-black font-bold shadow-[0_0_8px_#00F5A0]"
                    : "text-text-subtle hover:text-white"
                }`}
              >
                {r}km
              </button>
            ))}
          </div>

          <div className="h-4 w-px bg-border/80 hidden sm:block" />

          {/* Coordinates or Click Prompt */}
          <div className="text-[11px]">
            {bufferCenter ? (
              <span>
                EPICENTER: <span className="text-white font-bold">{bufferCenter.lat}°, {bufferCenter.lng}°</span>
              </span>
            ) : (
              <span className="text-aurora animate-pulse">Click map to drop radius epicenter</span>
            )}
          </div>

          {/* Analytics Results */}
          {bufferCenter && (
            <div className="flex items-center gap-2 bg-[#0C1812] px-2 py-1 rounded-lg border border-aurora/40 text-[11px]">
              <div>
                INTERSECTIONS:{" "}
                <span className="text-aurora font-bold">{bufferIntersectingProjects.length}</span>
              </div>
              <span className="text-border">|</span>
              <div>
                SUM CAPEX:{" "}
                <span className="text-aurora font-bold">
                  {bufferTotalCapex >= 1e9
                    ? `$${(bufferTotalCapex / 1e9).toFixed(2)}B CAD`
                    : bufferTotalCapex > 0
                    ? `$${(bufferTotalCapex / 1e6).toFixed(0)}M CAD`
                    : "$0 CAD"}
                </span>
              </div>
            </div>
          )}

          {/* Action Buttons */}
          <div className="flex items-center gap-1.5 ml-auto">
            {bufferCenter && (
              <button
                onClick={clearBuffer}
                className="px-2 py-1 rounded bg-surface hover:bg-surface/80 text-text-muted hover:text-white text-[10px] flex items-center gap-1"
                title="Reset buffer epicenter"
              >
                <Trash2 className="h-3 w-3" />
                Clear
              </button>
            )}
            <button
              onClick={() => {
                setIsBufferActive(false);
                clearBuffer();
              }}
              className="px-2 py-1 rounded bg-aurora/20 text-aurora hover:bg-aurora/30 text-[10px] font-bold"
            >
              Exit Tool
            </button>
          </div>
        </div>
      )}

      {/* Main Map Container */}
      <div ref={mapContainerRef} className="w-full h-full z-0 select-none" />

      {/* Bottom Telemetry HUD Bar */}
      <div className="absolute bottom-3 left-3 right-3 z-[1000] pointer-events-none flex items-center justify-between">
        <div className="pointer-events-auto bg-[#040806]/90 backdrop-blur-md border border-border/80 px-3 py-1.5 rounded-xl text-[10px] font-mono text-text-subtle flex flex-wrap items-center gap-3 shadow-xl">
          {/* Canadian Government FIP Emblem */}
          <div className="flex items-center gap-1.5 text-white">
            <svg className="h-3 w-5 shrink-0 rounded-[1px] overflow-hidden border border-white/20" viewBox="0 0 100 50">
              <rect width="25" height="50" fill="#D8292F" />
              <rect x="25" width="50" height="50" fill="#FFFFFF" />
              <rect x="75" width="25" height="50" fill="#D8292F" />
              <path
                d="M 50 10 L 52 18 L 59 15 L 56 22 L 64 22 L 59 27 L 66 33 L 57 33 L 54 36 L 53 43 L 51 43 L 50 41 L 49 43 L 47 43 L 46 36 L 43 33 L 34 33 L 41 27 L 36 22 L 44 22 L 41 15 L 48 18 Z"
                fill="#D8292F"
              />
            </svg>
            <span className="font-bold text-white tracking-tight">NRCAN / ISED</span>
          </div>

          <div className="flex items-center gap-1 text-aurora">
            <Navigation className="h-3 w-3 animate-pulse" />
            <span>GEODETIC TELEMETRY</span>
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
            DATUM: <span className="text-white">NAD83 / WGS84 (EPSG:3857)</span>
          </div>
        </div>

        {/* Layer Counter Indicator */}
        <div className="hidden md:flex pointer-events-auto bg-[#040806]/90 backdrop-blur-md border border-border/80 px-3 py-1.5 rounded-xl text-[10px] font-mono text-text-subtle items-center gap-2 shadow-xl">
          <span className="text-aurora font-bold">{projects.length}</span> projects
          <span>•</span>
          <span className="text-purple-400 font-bold">{CRITICAL_MINERAL_HUBS.length}</span> hubs
          <span>•</span>
          <span className="text-emerald-400 font-bold">{TREATY_TERRITORIES.length}</span> treaties
          <span>•</span>
          <span className="text-sky-400 font-bold">{GRID_INTERTIE_ZONES.length}</span> interties
          <span>•</span>
          <span className="text-amber-400 font-bold">{OPPORTUNITY_ZONES.length}</span> zones
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

      {/* Dynamic Detail Flyout: Selected Critical Mineral Hub */}
      {selectedMineralHub && (
        <div className="absolute top-16 right-3 z-[1000] w-84 bg-[#140a24]/95 backdrop-blur-md border border-purple-500/50 rounded-xl p-4 text-xs font-mono shadow-2xl space-y-2.5 animate-in fade-in slide-in-from-right-4">
          <div className="flex items-center justify-between border-b border-purple-900/60 pb-2">
            <div className="text-[10px] font-bold text-purple-400 uppercase tracking-wider flex items-center gap-1.5">
              <Factory className="h-3.5 w-3.5" />
              Critical Mineral Midstream Hub
            </div>
            <button
              onClick={() => setSelectedMineralHub(null)}
              className="text-text-subtle hover:text-white text-sm"
            >
              ✕
            </button>
          </div>
          <div className="font-bold text-sm text-purple-100">{selectedMineralHub.name}</div>
          <div className="flex flex-wrap gap-1.5">
            <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-purple-900/40 text-purple-300 border border-purple-700/50">
              {selectedMineralHub.province} • {selectedMineralHub.processingType}
            </span>
            <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-950/60 text-aurora border border-aurora/30">
              DVRI Retention: {(selectedMineralHub.domesticRetentionRate * 100).toFixed(0)}%
            </span>
          </div>
          <div className="space-y-1 bg-[#07030e]/80 p-2 rounded-lg border border-purple-900/40 text-[11px]">
            <div>
              <span className="text-text-subtle">Capacity:</span>{" "}
              <span className="text-white font-semibold">{selectedMineralHub.capacityMetric}</span>
            </div>
            <div>
              <span className="text-text-subtle">Annual Value:</span>{" "}
              <span className="text-aurora font-semibold">{selectedMineralHub.annualValueCAD}</span>
            </div>
          </div>
          <div className="text-[10px] text-text-subtle">
            <span className="font-bold uppercase text-white">Focus Minerals:</span>{" "}
            {selectedMineralHub.focusMinerals.join(", ")}
          </div>
        </div>
      )}

      {/* Dynamic Detail Flyout: Selected Treaty Territory */}
      {selectedTreaty && (
        <div className="absolute top-16 right-3 z-[1000] w-84 bg-[#081a12]/95 backdrop-blur-md border border-emerald-500/50 rounded-xl p-4 text-xs font-mono shadow-2xl space-y-2.5 animate-in fade-in slide-in-from-right-4">
          <div className="flex items-center justify-between border-b border-emerald-900/60 pb-2">
            <div className="text-[10px] font-bold text-emerald-400 uppercase tracking-wider flex items-center gap-1.5">
              <Landmark className="h-3.5 w-3.5" />
              Treaty Land & Economic Jurisdiction
            </div>
            <button
              onClick={() => setSelectedTreaty(null)}
              className="text-text-subtle hover:text-white text-sm"
            >
              ✕
            </button>
          </div>
          <div className="font-bold text-sm text-emerald-100">{selectedTreaty.name}</div>
          <div className="flex flex-wrap gap-1.5">
            <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-emerald-900/40 text-emerald-300 border border-emerald-700/50">
              {selectedTreaty.historicalFramework.replace(/_/g, " ")}
            </span>
            <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-[#040806] text-aurora border border-aurora/40">
              {selectedTreaty.province}
            </span>
          </div>
          <div className="space-y-1 bg-[#040806]/80 p-2 rounded-lg border border-emerald-900/40 text-[11px]">
            <div className="text-aurora font-semibold">
              {selectedTreaty.loanGuaranteeEligibility}
            </div>
            <div className="text-text-muted text-[10px]">
              Eligible for up to 100% debt guarantee under Canada&apos;s $5B ILGP facility.
            </div>
          </div>
          <div className="text-[10px] text-text-subtle">
            <span className="font-bold uppercase text-white">Signatories:</span>{" "}
            {selectedTreaty.signatories.join(", ")}
          </div>
        </div>
      )}

      {/* Dynamic Detail Flyout: Selected Clean Grid Intertie */}
      {selectedIntertie && (
        <div className="absolute top-16 right-3 z-[1000] w-84 bg-[#071520]/95 backdrop-blur-md border border-sky-500/50 rounded-xl p-4 text-xs font-mono shadow-2xl space-y-2.5 animate-in fade-in slide-in-from-right-4">
          <div className="flex items-center justify-between border-b border-sky-900/60 pb-2">
            <div className="text-[10px] font-bold text-sky-400 uppercase tracking-wider flex items-center gap-1.5">
              <Cpu className="h-3.5 w-3.5" />
              Clean Baseload & AI Intertie
            </div>
            <button
              onClick={() => setSelectedIntertie(null)}
              className="text-text-subtle hover:text-white text-sm"
            >
              ✕
            </button>
          </div>
          <div className="font-bold text-sm text-sky-100">{selectedIntertie.name}</div>
          <div className="flex flex-wrap gap-1.5">
            <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-sky-900/40 text-sky-300 border border-sky-700/50">
              {selectedIntertie.gridOperator}
            </span>
            <span className="px-2 py-0.5 rounded text-[10px] font-bold bg-[#040806] text-white border border-border">
              {selectedIntertie.baseloadType}
            </span>
          </div>
          <div className="space-y-1 bg-[#040806]/80 p-2 rounded-lg border border-sky-900/40 text-[11px]">
            <div className="flex justify-between">
              <span className="text-text-subtle">Clean Capacity:</span>{" "}
              <span className="text-white font-bold">{selectedIntertie.cleanCapacityMW} MW</span>
            </div>
            <div className="flex justify-between">
              <span className="text-text-subtle">AI Headroom:</span>{" "}
              <span className="text-sky-300 font-bold">{selectedIntertie.aiHeadroomMW} MW</span>
            </div>
            <div className="flex justify-between">
              <span className="text-text-subtle">Clean Compute Efficiency:</span>{" "}
              <span className="text-aurora font-bold">{selectedIntertie.cleanFlopsRatio}</span>
            </div>
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
