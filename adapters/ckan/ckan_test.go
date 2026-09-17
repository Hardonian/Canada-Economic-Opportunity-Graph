package ckan

import (
	"testing"
)

func TestClassifyResource(t *testing.T) {
	tests := []struct {
		format       string
		mimeType     string
		resourceType string
		rawURL       string
		expected     Family
	}{
		{
			format:   "geojson",
			mimeType: "application/geo+json",
			rawURL:   "https://open.canada.ca/data/dataset.geojson",
			expected: FamilyGeoJSON,
		},
		{
			format:   "wfs",
			mimeType: "",
			rawURL:   "https://maps.geogratis.gc.ca/wfs?service=wfs&request=GetCapabilities",
			expected: FamilyWFS,
		},
		{
			format:   "ArcGIS",
			mimeType: "",
			rawURL:   "https://services.arcgis.com/rest/services/FeatureServer/0",
			expected: FamilyArcGISFeature,
		},
		{
			format:   "csv",
			mimeType: "text/csv",
			rawURL:   "https://open.canada.ca/data/projects.csv",
			expected: FamilyCSV,
		},
	}

	for _, tt := range tests {
		got := ClassifyResource(tt.format, tt.mimeType, tt.resourceType, tt.rawURL)
		if got != tt.expected {
			t.Errorf("ClassifyResource(%q, %q, %q, %q) = %v; want %v",
				tt.format, tt.mimeType, tt.resourceType, tt.rawURL, got, tt.expected)
		}
	}
}
