package merkle

import (
	"context"
	"testing"
	"time"
)

func TestSignAndVerifyMerkleRoot(t *testing.T) {
	signer, err := NewECDSASigner("kms-key-ca-prod-01")
	if err != nil {
		t.Fatalf("failed to create signer: %v", err)
	}

	root := &Root{
		MethodologyVersion: MethodologyVersion,
		RootHash:           "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		LeafCount:          10,
		PublishedAt:        time.Now().UTC(),
	}

	ctx := context.Background()
	signed, err := SignRoot(ctx, root, signer, "sovereign-auditor.gc.ca")
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}

	if signed.KeyID != "kms-key-ca-prod-01" {
		t.Fatalf("expected key id kms-key-ca-prod-01, got %s", signed.KeyID)
	}
	if signed.Algorithm != "ECDSA_P256_SHA256" {
		t.Fatalf("expected ECDSA_P256_SHA256, got %s", signed.Algorithm)
	}
	if len(signed.SignatureHex) == 0 {
		t.Fatal("expected non-empty signature hex")
	}

	// Verify valid signature
	valid, err := VerifySignedRoot(signed, "")
	if err != nil {
		t.Fatalf("verify failed with error: %v", err)
	}
	if !valid {
		t.Fatal("expected valid signature verification")
	}

	// Tamper test: modifying root hash must cause verification to fail
	tampered := &SignedRoot{
		Root: &Root{
			MethodologyVersion: root.MethodologyVersion,
			RootHash:           "bad-hash-tampered-content",
			LeafCount:          root.LeafCount,
			PublishedAt:        root.PublishedAt,
		},
		SignatureHex:   signed.SignatureHex,
		PublicKeyPEM:   signed.PublicKeyPEM,
		KeyID:          signed.KeyID,
		Algorithm:      signed.Algorithm,
		SignerIdentity: signed.SignerIdentity,
		SignedAt:       signed.SignedAt,
	}

	tamperValid, err := VerifySignedRoot(tampered, "")
	if err != nil {
		t.Fatalf("unexpected error during tampered verify: %v", err)
	}
	if tamperValid {
		t.Fatal("expected tampered root verification to fail, but it returned true")
	}
}
