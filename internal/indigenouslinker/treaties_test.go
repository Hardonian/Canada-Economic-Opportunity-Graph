package indigenouslinker

import (
	"testing"
)

func TestFindTreatiesForLocation(t *testing.T) {
	tests := []struct {
		name         string
		lat          float64
		lng          float64
		province     string
		expectedID   string
		expectFound  bool
	}{
		{
			name:        "Ring of Fire (Treaty 9 Ontario)",
			lat:         52.7,
			lng:         -86.2,
			province:    "ON",
			expectedID:  "treaty-9",
			expectFound: true,
		},
		{
			name:        "Edmonton / Central Alberta (Treaty 6)",
			lat:         53.5461,
			lng:         -113.4938,
			province:    "AB",
			expectedID:  "treaty-6",
			expectFound: true,
		},
		{
			name:        "Calgary / Southern Alberta (Treaty 7)",
			lat:         51.0447,
			lng:         -114.0719,
			province:    "AB",
			expectedID:  "treaty-7",
			expectFound: true,
		},
		{
			name:        "Nass Valley / Nisga'a (Modern Land Claim BC)",
			lat:         55.2,
			lng:         -129.2,
			province:    "BC",
			expectedID:  "treaty-nisgaa",
			expectFound: true,
		},
		{
			name:        "James Bay / Northern Quebec (JBNQA)",
			lat:         53.0,
			lng:         -75.0,
			province:    "QC",
			expectedID:  "treaty-james-bay",
			expectFound: true,
		},
		{
			name:        "Iqaluit (Nunavut Agreement)",
			lat:         63.7467,
			lng:         -68.5170,
			province:    "NU",
			expectedID:  "treaty-nunavut",
			expectFound: true,
		},
		{
			name:        "Out of bounds coordinate (Atlantic Ocean)",
			lat:         45.0,
			lng:         -40.0,
			province:    "NS",
			expectedID:  "",
			expectFound: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			matches := FindTreatiesForLocation(tc.lat, tc.lng, tc.province)
			if !tc.expectFound {
				if len(matches) > 0 {
					t.Fatalf("expected no matches for out of bounds, got %d", len(matches))
				}
				return
			}

			if len(matches) == 0 {
				t.Fatalf("expected at least one match for %s, got 0", tc.name)
			}

			found := false
			for _, m := range matches {
				if m.ID == tc.expectedID {
					found = true
					if m.YearSigned <= 0 {
						t.Errorf("invalid YearSigned for %s: %d", m.ID, m.YearSigned)
					}
					if len(m.SignatoryNations) == 0 {
						t.Errorf("expected signatory nations for %s, got none", m.ID)
					}
					break
				}
			}

			if !found {
				t.Errorf("expected treaty %s in matches for %s, but was not found", tc.expectedID, tc.name)
			}
		})
	}
}

func TestTreatyRegistryIntegrity(t *testing.T) {
	if len(CanadaTreatiesRegistry) < 5 {
		t.Fatalf("expected at least 5 registered treaties, got %d", len(CanadaTreatiesRegistry))
	}

	seenIDs := make(map[string]bool)
	for _, tr := range CanadaTreatiesRegistry {
		if tr.ID == "" {
			t.Errorf("treaty has empty ID: %+v", tr)
		}
		if seenIDs[tr.ID] {
			t.Errorf("duplicate treaty ID: %s", tr.ID)
		}
		seenIDs[tr.ID] = true

		if tr.MinLat >= tr.MaxLat {
			t.Errorf("invalid lat bounds for %s: min %f >= max %f", tr.ID, tr.MinLat, tr.MaxLat)
		}
		if tr.MinLng >= tr.MaxLng {
			t.Errorf("invalid lng bounds for %s: min %f >= max %f", tr.ID, tr.MinLng, tr.MaxLng)
		}
	}
}
