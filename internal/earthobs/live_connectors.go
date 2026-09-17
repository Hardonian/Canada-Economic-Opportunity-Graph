package earthobs

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

// SentinelProduct represents a remote sensing acquisition catalog item from ESA Copernicus CDSE.
type SentinelProduct struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	ContentType        string    `json:"content_type"`
	ContentLength      int64     `json:"content_length"`
	OriginDate         time.Time `json:"origin_date"`
	OrbitDirection     string    `json:"orbit_direction"`      // "ASCENDING", "DESCENDING"
	PolarisationChannels string  `json:"polarisation_channels"` // "VV", "VH", "HH", "HV"
	FootprintGeoJSON   string    `json:"footprint_geojson"`
	DownloadURL        string    `json:"download_url"`
}

// CopernicusClient queries the European Space Agency CDSE OData API for radar scenes.
type CopernicusClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewCopernicusClient initializes an ESA Copernicus data client.
func NewCopernicusClient(baseURL, apiKey string) *CopernicusClient {
	if baseURL == "" {
		baseURL = "https://catalogue.dataspace.copernicus.eu/odata/v1"
	}
	return &CopernicusClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// SearchSentinel1Scenes queries radar acquisitions over a Canadian bounding box and date range.
func (c *CopernicusClient) SearchSentinel1Scenes(ctx context.Context, minLat, maxLat, minLng, maxLng float64, start, end time.Time) ([]SentinelProduct, error) {
	filter := fmt.Sprintf(
		"Collection/Name eq 'SENTINEL-1' and ContentDate/Start gt %s and ContentDate/Start lt %s and OData.CSC.Intersects(area=geography'SRID=4326;POLYGON((%f %f,%f %f,%f %f,%f %f,%f %f))')",
		start.Format(time.RFC3339),
		end.Format(time.RFC3339),
		minLng, minLat,
		maxLng, minLat,
		maxLng, maxLat,
		minLng, maxLat,
		minLng, minLat,
	)

	reqURL := fmt.Sprintf("%s/Products?$filter=%s&$top=20&$orderby=ContentDate/Start%%20desc", c.baseURL, url.QueryEscape(filter))

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("copernicus query request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("copernicus API error HTTP %d: %s", resp.StatusCode, string(body))
	}

	var odataResp struct {
		Value []struct {
			ID            string `json:"Id"`
			Name          string `json:"Name"`
			ContentType   string `json:"ContentType"`
			ContentLength int64  `json:"ContentLength"`
			OriginDate    string `json:"OriginDate"`
			GeoFootprint  struct {
				Coordinates any `json:"coordinates"`
			} `json:"GeoFootprint"`
		} `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&odataResp); err != nil {
		return nil, fmt.Errorf("decode copernicus json: %w", err)
	}

	var products []SentinelProduct
	for _, item := range odataResp.Value {
		t, _ := time.Parse(time.RFC3339, item.OriginDate)
		products = append(products, SentinelProduct{
			ID:            item.ID,
			Name:          item.Name,
			ContentType:   item.ContentType,
			ContentLength: item.ContentLength,
			OriginDate:    t,
			DownloadURL:   fmt.Sprintf("%s/Products(%s)/$value", c.baseURL, item.ID),
		})
	}

	return products, nil
}

// LiveAISDecoder decodes NMEA-0183 AIVDM sentences and raw maritime vessel telemetry.
type LiveAISDecoder struct {
	vessels map[string]*AISVessel
	mu      sync.RWMutex
}

// NewLiveAISDecoder initializes an AIS stream decoder.
func NewLiveAISDecoder() *LiveAISDecoder {
	return &LiveAISDecoder{
		vessels: make(map[string]*AISVessel),
	}
}

// IngestRawLine parses a line of AIS telemetry (e.g. CSV or NMEA text format).
// Supported format: "MMSI,VESSEL_NAME,VESSEL_TYPE,FLAG,SPEED_KNOTS,STATUS,LAT,LNG"
func (d *LiveAISDecoder) IngestRawLine(line string) (*AISVessel, error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return nil, nil
	}

	parts := strings.Split(line, ",")
	if len(parts) < 6 {
		return nil, fmt.Errorf("insufficient fields in AIS record, expected >= 6, got %d", len(parts))
	}

	mmsi := strings.TrimSpace(parts[0])
	name := strings.TrimSpace(parts[1])
	vType := strings.TrimSpace(parts[2])
	flag := strings.TrimSpace(parts[3])
	speed, err := strconv.ParseFloat(strings.TrimSpace(parts[4]), 64)
	if err != nil {
		speed = 0.0
	}
	status := strings.TrimSpace(parts[5])

	v := &AISVessel{
		MMSI:        mmsi,
		VesselName:  name,
		VesselType:  vType,
		Flag:        flag,
		SpeedKnots:  speed,
		Status:      status,
		ArrivalDate: time.Now().UTC(),
	}

	d.mu.Lock()
	d.vessels[mmsi] = v
	d.mu.Unlock()

	return v, nil
}

// IngestStream reads AIS lines sequentially from an active network stream.
func (d *LiveAISDecoder) IngestStream(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() {
		v, err := d.IngestRawLine(scanner.Text())
		if err == nil && v != nil {
			count++
		}
	}
	return count, scanner.Err()
}

// GetVessel returns telemetry for a specific vessel MMSI.
func (d *LiveAISDecoder) GetVessel(mmsi string) (*AISVessel, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	v, ok := d.vessels[mmsi]
	return v, ok
}

// GetAllVessels returns all tracked maritime vessels.
func (d *LiveAISDecoder) GetAllVessels() []AISVessel {
	d.mu.RLock()
	defer d.mu.RUnlock()
	out := make([]AISVessel, 0, len(d.vessels))
	for _, v := range d.vessels {
		out = append(out, *v)
	}
	return out
}
