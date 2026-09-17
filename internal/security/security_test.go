package security

import (
	"bytes"
	"testing"
)

func TestABACAuthorization(t *testing.T) {
	evaluator := NewABACEvaluator()

	// Canadian Secret Analyst
	subCA := SecuritySubject{
		UserID:       "user-dnd-01",
		Clearance:    ClassSecret,
		Citizenship:  "CA",
		Organization: "Department of National Defence",
		Compartments: []string{"NUCLEAR", "ARCTIC"},
	}

	// Foreign Partner (US) with Secret
	subUS := SecuritySubject{
		UserID:       "user-us-01",
		Clearance:    ClassSecret,
		Citizenship:  "US",
		Organization: "DoD",
		Compartments: []string{"NUCLEAR"},
	}

	// Secret + Canadian Eyes Only Resource
	objCanOnly := SecurityObject{
		ResourceID:          "doc-arctic-radar",
		ResourceType:        "DOCUMENT",
		Classification:      ClassSecret,
		RequiredCompartment: "ARCTIC",
		Caveats:             []Caveat{CaveatCanadianEyesOnly},
	}

	// Canadian should be permitted
	decCA := evaluator.Authorize(subCA, objCanOnly, "READ")
	if !decCA.Permitted {
		t.Errorf("expected Canadian subject to be permitted, got denied: %s", decCA.DeniedReason)
	}

	// US partner should be denied due to Canadian Eyes Only
	decUS := evaluator.Authorize(subUS, objCanOnly, "READ")
	if decUS.Permitted {
		t.Errorf("expected foreign subject to be denied on CANADIAN_EYES_ONLY resource")
	}

	// Top Secret object should be denied to Secret subject
	objTS := SecurityObject{
		ResourceID:     "doc-alien-tech",
		Classification: ClassTopSecret,
	}
	decTS := evaluator.Authorize(subCA, objTS, "READ")
	if decTS.Permitted {
		t.Errorf("expected Secret subject to be denied on TOP_SECRET resource")
	}
}

func TestSecurityRedactionAndBanner(t *testing.T) {
	banner := GenerateSecurityBanner(ClassSecret, []Caveat{CaveatCanadianEyesOnly})
	expectedBanner := "// CLASSIFICATION: SECRET // CANADIAN_EYES_ONLY //"
	if banner != expectedBanner {
		t.Errorf("expected %s, got %s", expectedBanner, banner)
	}

	redactor := NewRedactor()
	data := map[string]interface{}{
		"project_name": "SMR Fuel Fabrication",
		"capex_cad":    500000000.0,
		"enrichment_purity_pct": 19.75, // Secret
	}
	fieldLevels := map[string]ClassificationLevel{
		"enrichment_purity_pct": ClassSecret,
		"capex_cad":             ClassProtectedB,
	}

	subProtB := SecuritySubject{Clearance: ClassProtectedB}
	redacted := redactor.RedactRecord(data, subProtB, fieldLevels)

	if redacted["enrichment_purity_pct"] != "[REDACTED - SECRET]" {
		t.Errorf("expected enrichment_purity_pct to be redacted, got %v", redacted["enrichment_purity_pct"])
	}
	if redacted["capex_cad"] != 500000000.0 {
		t.Errorf("expected capex_cad to remain unredacted for Protected B subject")
	}
}

func TestFieldEncryption(t *testing.T) {
	key := []byte("01234567890123456789012345678901") // 32 bytes
	fe, err := NewFieldEncryptor(key)
	if err != nil {
		t.Fatalf("failed to create encryptor: %v", err)
	}

	secretPayload := []byte("Confidential PPA Offtake: $64.50/MWh fixed for 25 years")
	cHex, nHex, err := fe.Encrypt(secretPayload)
	if err != nil {
		t.Fatalf("encryption err: %v", err)
	}

	decrypted, err := fe.Decrypt(cHex, nHex)
	if err != nil {
		t.Fatalf("decryption err: %v", err)
	}

	if !bytes.Equal(decrypted, secretPayload) {
		t.Errorf("decrypted payload does not match original plaintext")
	}
}

func TestDLPScanner(t *testing.T) {
	dlp := NewDLPScanner()
	raw := "Employee lead John Doe SIN: 123-456-789 deployed to secret location NORAD_SITE_ALPHA4 for work."

	scrubbed, findings := dlp.ScrubText(raw)
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
	if findings[0] != "CANADIAN_SIN" || findings[1] != "DEFENSE_SECRET" {
		t.Errorf("unexpected findings: %v", findings)
	}
	if !bytes.Contains([]byte(scrubbed), []byte("[REDACTED-SIN]")) {
		t.Errorf("expected [REDACTED-SIN] in scrubbed text: %s", scrubbed)
	}
}
