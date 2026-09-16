package memoexport

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/Hardonian/CEO-G-Canada-Economic-Opportunity-Graph/internal/domain"
)

// GeoJSONFeatureCollection models standard OGC-compliant GeoJSON export.
type GeoJSONFeatureCollection struct {
	Type      string           `json:"type"` // "FeatureCollection"
	Features  []GeoJSONFeature `json:"features"`
	AuditHash string           `json:"audit_hash"`
}

// GeoJSONFeature represents a spatial project node.
type GeoJSONFeature struct {
	Type       string                 `json:"type"` // "Feature"
	Geometry   GeoJSONGeometry        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

// GeoJSONGeometry represents point coordinates [lon, lat].
type GeoJSONGeometry struct {
	Type        string    `json:"type"` // "Point"
	Coordinates []float64 `json:"coordinates"`
}

// STACItem models a SpatioTemporal Asset Catalog metadata record for satellite passes.
type STACItem struct {
	Type       string                 `json:"type"` // "Feature"
	STACVersion string                `json:"stac_version"`
	ID         string                 `json:"id"`
	Bbox       []float64              `json:"bbox"`
	Geometry   GeoJSONGeometry        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
	Assets     map[string]interface{} `json:"assets"`
}

// ExportGeoJSON converts a list of projects into an OGC GeoJSON collection.
func ExportGeoJSON(projects []*domain.Project) *GeoJSONFeatureCollection {
	features := make([]GeoJSONFeature, 0, len(projects))

	for _, p := range projects {
		if p.Latitude == 0 && p.Longitude == 0 {
			continue
		}

		props := map[string]interface{}{
			"id":            p.ID,
			"slug":          p.Slug,
			"name":          p.Name,
			"sector":        string(p.Sector),
			"province":      p.Province,
			"stage":         string(p.CurrentStage),
			"capex_cad":     p.CapexCAD,
			"location_name": p.LocationName,
		}

		features = append(features, GeoJSONFeature{
			Type: "Feature",
			Geometry: GeoJSONGeometry{
				Type:        "Point",
				Coordinates: []float64{p.Longitude, p.Latitude},
			},
			Properties: props,
		})
	}

	coll := &GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}

	auditData := fmt.Sprintf("geojson|%d", len(features))
	h := sha256.Sum256([]byte(auditData))
	coll.AuditHash = hex.EncodeToString(h[:])

	return coll
}

// ExportSTACItem creates a STAC satellite metadata asset for a project.
func ExportSTACItem(project *domain.Project) *STACItem {
	lat := project.Latitude
	lon := project.Longitude
	if lat == 0 && lon == 0 {
		lat = 51.25
		lon = -85.32
	}

	bbox := []float64{lon - 0.05, lat - 0.05, lon + 0.05, lat + 0.05}

	return &STACItem{
		Type:        "Feature",
		STACVersion: "1.0.0",
		ID:          fmt.Sprintf("stac-%s-sar", project.Slug),
		Bbox:        bbox,
		Geometry: GeoJSONGeometry{
			Type:        "Point",
			Coordinates: []float64{lon, lat},
		},
		Properties: map[string]interface{}{
			"datetime":            "2026-09-15T00:00:00Z",
			"platform":            "sentinel-1b-radarsat-constellation",
			"constellation":       "RADARSAT",
			"instruments":         []string{"c-sar"},
			"sar:instrument_mode": "IW",
			"sar:polarizations":   []string{"VV", "VH"},
			"project_name":        project.Name,
			"project_slug":        project.Slug,
		},
		Assets: map[string]interface{}{
			"backscatter_vv": map[string]string{
				"href": fmt.Sprintf("https://earthobs.cog.cegs.dev/passes/%s_vv.tif", project.Slug),
				"type": "image/tiff; application=geotiff; profile=cloud-optimized",
			},
			"coherence": map[string]string{
				"href": fmt.Sprintf("https://earthobs.cog.cegs.dev/passes/%s_coh.tif", project.Slug),
				"type": "image/tiff; application=geotiff; profile=cloud-optimized",
			},
		},
	}
}
