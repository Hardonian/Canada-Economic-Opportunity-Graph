package security

import (
	"context"
	"testing"
)

func TestHardwareSecurityModule_Lifecycle(t *testing.T) {
	ctx := context.Background()

	// 1. Initialize HSM token
	hsm, err := NewHardwareSecurityModule(1, "CANADIAN_SOVEREIGN_HSM_01", "pin-secret-1234")
	if err != nil {
		t.Fatalf("failed to init HSM: %v", err)
	}

	keyMaterial := []byte("01234567890123456789012345678901") // 32 bytes
	if err := hsm.ImportHardwareKey("master-lakehouse-key", keyMaterial); err != nil {
		t.Fatalf("failed to import key to HSM: %v", err)
	}

	// 2. Read key and sign digest
	k, err := hsm.GetKey(ctx, "master-lakehouse-key")
	if err != nil || len(k) != 32 {
		t.Fatalf("expected 32-byte key from HSM, err=%v", err)
	}

	sig, err := hsm.Sign(ctx, "master-lakehouse-key", []byte("payload-to-sign"))
	if err != nil || len(sig) == 0 {
		t.Fatalf("expected non-empty HMAC signature, err=%v", err)
	}

	// 3. Zeroize HSM session
	hsm.Zeroize()
	_, err = hsm.GetKey(ctx, "master-lakehouse-key")
	if err == nil {
		t.Fatal("expected error getting key after HSM zeroize")
	}
}

func TestEnclaveAttestation_GenerateAndVerify(t *testing.T) {
	refPCR0 := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	engine := NewAttestationEngine(refPCR0)

	quote := engine.GenerateQuote([]byte("runtime-session-nonce-12345"))
	if quote.PCR0 != refPCR0 {
		t.Fatalf("expected quote PCR0 %s, got %s", refPCR0, quote.PCR0)
	}

	valid, msg := engine.VerifyQuote(quote)
	if !valid {
		t.Fatalf("expected valid quote, got error: %s", msg)
	}

	// Tampered quote PCR0 should fail
	quote.PCR0 = "tampered-hash-value"
	valid, _ = engine.VerifyQuote(quote)
	if valid {
		t.Fatal("expected tampered quote to fail verification")
	}
}
