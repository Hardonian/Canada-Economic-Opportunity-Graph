package indigenouslinker

import (
	"strings"
)

// TreatyClassification categorizes Canadian Crown-Indigenous agreements.
type TreatyClassification string

const (
	TreatyNumbered           TreatyClassification = "NUMBERED_TREATY"
	TreatyHistoricPreConfed  TreatyClassification = "HISTORIC_PRE_CONFEDERATION"
	TreatyModernLandClaim    TreatyClassification = "MODERN_COMPREHENSIVE_CLAIM"
	TreatyUnceded            TreatyClassification = "UNCEDED_TERRITORY"
)

// TreatyTerritory defines sovereign boundary metadata for Indigenous agreements.
type TreatyTerritory struct {
	ID             string                `json:"id"`
	Name           string                `json:"name"`
	YearSigned     int                   `json:"year_signed"`
	Classification TreatyClassification  `json:"classification"`
	Provinces      []string              `json:"provinces"`
	MinLat         float64               `json:"min_lat"`
	MaxLat         float64               `json:"max_lat"`
	MinLng         float64               `json:"min_lng"`
	MaxLng         float64               `json:"max_lng"`
	SignatoryNations []string            `json:"signatory_nations"`
}

// CanadaTreatiesRegistry contains authoritative spatial boundary references across Canada.
var CanadaTreatiesRegistry = []TreatyTerritory{
	{
		ID:             "treaty-1",
		Name:           "Treaty 1 (Stone Fort Treaty)",
		YearSigned:     1871,
		Classification: TreatyNumbered,
		Provinces:      []string{"MB"},
		MinLat:         49.0, MaxLat: 51.0, MinLng: -98.5, MaxLng: -96.0,
		SignatoryNations: []string{"Anishinaabe (Ojibwe)", "Nehiyaw (Cree)"},
	},
	{
		ID:             "treaty-6",
		Name:           "Treaty 6 (Central Alberta & Saskatchewan)",
		YearSigned:     1876,
		Classification: TreatyNumbered,
		Provinces:      []string{"AB", "SK"},
		MinLat:         51.5, MaxLat: 54.5, MinLng: -116.0, MaxLng: -104.0,
		SignatoryNations: []string{"Plains Cree", "Woods Cree", "Nakota", "Saulteaux"},
	},
	{
		ID:             "treaty-7",
		Name:           "Treaty 7 (Southern Alberta / Blackfoot Confederacy)",
		YearSigned:     1877,
		Classification: TreatyNumbered,
		Provinces:      []string{"AB"},
		MinLat:         49.0, MaxLat: 52.0, MinLng: -115.5, MaxLng: -110.0,
		SignatoryNations: []string{"Siksika", "Kainai", "Piikani", "Stoney Nakoda", "Tsuut'ina"},
	},
	{
		ID:             "treaty-8",
		Name:           "Treaty 8 (Athabasca / Peace River / Northern AB & BC)",
		YearSigned:     1899,
		Classification: TreatyNumbered,
		Provinces:      []string{"AB", "BC", "SK", "NT"},
		MinLat:         54.0, MaxLat: 61.0, MinLng: -124.0, MaxLng: -106.0,
		SignatoryNations: []string{"Cree", "Dene", "Dane-zaa (Beaver)", "Chipewyan"},
	},
	{
		ID:             "treaty-9",
		Name:           "Treaty 9 (James Bay Treaty / Northern Ontario)",
		YearSigned:     1905,
		Classification: TreatyNumbered,
		Provinces:      []string{"ON"},
		MinLat:         48.0, MaxLat: 56.5, MinLng: -94.0, MaxLng: -79.5,
		SignatoryNations: []string{"Ojibwe", "Cree", "Oji-Cree (Ring of Fire)"},
	},
	{
		ID:             "treaty-robinson-superior",
		Name:           "Robinson-Superior Treaty",
		YearSigned:     1850,
		Classification: TreatyHistoricPreConfed,
		Provinces:      []string{"ON"},
		MinLat:         47.5, MaxLat: 50.5, MinLng: -90.5, MaxLng: -84.5,
		SignatoryNations: []string{"Anishinaabe of Lake Superior"},
	},
	{
		ID:             "treaty-james-bay",
		Name:           "James Bay and Northern Quebec Agreement (JBNQA)",
		YearSigned:     1975,
		Classification: TreatyModernLandClaim,
		Provinces:      []string{"QC"},
		MinLat:         49.0, MaxLat: 62.5, MinLng: -80.0, MaxLng: -64.0,
		SignatoryNations: []string{"Grand Council of the Crees (Eeyou Istchee)", "Makivik Corporation (Inuit of Nunavik)"},
	},
	{
		ID:             "treaty-nisgaa",
		Name:           "Nisga'a Final Agreement",
		YearSigned:     2000,
		Classification: TreatyModernLandClaim,
		Provinces:      []string{"BC"},
		MinLat:         54.8, MaxLat: 56.0, MinLng: -130.5, MaxLng: -128.5,
		SignatoryNations: []string{"Nisga'a Nation (Nass Valley)"},
	},
	{
		ID:             "treaty-nunavut",
		Name:           "Nunavut Land Claims Agreement",
		YearSigned:     1993,
		Classification: TreatyModernLandClaim,
		Provinces:      []string{"NU"},
		MinLat:         56.0, MaxLat: 83.0, MinLng: -120.0, MaxLng: -60.0,
		SignatoryNations: []string{"Tunngavik Federation of Nunavut (Inuit)"},
	},
}

// FindTreatiesForLocation identifies all Crown-Indigenous treaties and modern agreements covering a coordinate.
func FindTreatiesForLocation(lat, lng float64, province string) []*TreatyTerritory {
	var matches []*TreatyTerritory

	for i := range CanadaTreatiesRegistry {
		t := &CanadaTreatiesRegistry[i]

		// Check province membership
		provMatch := false
		for _, p := range t.Provinces {
			if strings.EqualFold(p, province) {
				provMatch = true
				break
			}
		}
		if !provMatch && province != "" {
			continue
		}

		// Spatial bounding box intersection
		if lat >= t.MinLat && lat <= t.MaxLat && lng >= t.MinLng && lng <= t.MaxLng {
			matches = append(matches, t)
		}
	}

	return matches
}
