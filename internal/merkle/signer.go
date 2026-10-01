// Package merkle provides hardware-backed cryptographic attestation and KMS signing
// for sovereign transparency log roots.
package merkle

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"time"
)

// KMSSigner abstracts hardware security modules (HSM), cloud KMS (AWS/GCP),
// or local cryptographic keys for signing transparency roots.
type KMSSigner interface {
	Sign(ctx context.Context, digest []byte) ([]byte, error)
	KeyID() string
	Algorithm() string
	PublicKeyPEM() (string, error)
}

// ECDSASigner is an in-memory/software implementation of KMSSigner using NIST P-256.
type ECDSASigner struct {
	keyID      string
	privateKey *ecdsa.PrivateKey
}

// NewECDSASigner generates a new P-256 ECDSA key pair for cryptographic signing.
func NewECDSASigner(keyID string) (*ECDSASigner, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate ecdsa key: %w", err)
	}
	if keyID == "" {
		keyID = "kms-p256-default"
	}
	return &ECDSASigner{
		keyID:      keyID,
		privateKey: key,
	}, nil
}

// KeyID returns the KMS key identifier.
func (s *ECDSASigner) KeyID() string {
	return s.keyID
}

// Algorithm returns the signature algorithm name.
func (s *ECDSASigner) Algorithm() string {
	return "ECDSA_P256_SHA256"
}

// PublicKeyPEM exports the public key in PEM format.
func (s *ECDSASigner) PublicKeyPEM() (string, error) {
	der, err := x509.MarshalPKIXPublicKey(&s.privateKey.PublicKey)
	if err != nil {
		return "", err
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}
	return string(pem.EncodeToMemory(block)), nil
}

// Sign signs the SHA-256 digest using ASN.1 DER-encoded ECDSA signature.
func (s *ECDSASigner) Sign(ctx context.Context, digest []byte) ([]byte, error) {
	return ecdsa.SignASN1(rand.Reader, s.privateKey, digest)
}

// SignedRoot encapsulates a Merkle root with its cryptographic signature and provenance.
type SignedRoot struct {
	Root           *Root     `json:"root"`
	SignatureHex   string    `json:"signature_hex"`
	KeyID          string    `json:"key_id"`
	Algorithm      string    `json:"algorithm"`
	SignerIdentity string    `json:"signer_identity"`
	SignedAt       time.Time `json:"signed_at"`
	PublicKeyPEM   string    `json:"public_key_pem,omitempty"`
}

// CanonicalDigest computes the deterministic byte sequence signed for a Merkle root.
func CanonicalDigest(root *Root) []byte {
	msg := fmt.Sprintf("%s:%s:%d:%s",
		root.MethodologyVersion,
		root.RootHash,
		root.LeafCount,
		root.PublishedAt.UTC().Format(time.RFC3339),
	)
	sum := sha256.Sum256([]byte(msg))
	return sum[:]
}

// SignRoot cryptographically signs a Merkle root using a KMSSigner.
func SignRoot(ctx context.Context, root *Root, signer KMSSigner, signerIdentity string) (*SignedRoot, error) {
	if root == nil {
		return nil, errors.New("cannot sign nil root")
	}
	if signer == nil {
		return nil, errors.New("signer cannot be nil")
	}

	digest := CanonicalDigest(root)
	sigBytes, err := signer.Sign(ctx, digest)
	if err != nil {
		return nil, fmt.Errorf("kms signing failed: %w", err)
	}

	pubPEM, _ := signer.PublicKeyPEM()
	if signerIdentity == "" {
		signerIdentity = "sovereign-attestation.gc.ca"
	}

	return &SignedRoot{
		Root:           root,
		SignatureHex:   hex.EncodeToString(sigBytes),
		KeyID:          signer.KeyID(),
		Algorithm:      signer.Algorithm(),
		SignerIdentity: signerIdentity,
		SignedAt:       time.Now().UTC(),
		PublicKeyPEM:   pubPEM,
	}, nil
}

// VerifySignedRoot verifies the ECDSA ASN.1 signature against the canonical Merkle digest.
func VerifySignedRoot(signed *SignedRoot, pubKeyPEM string) (bool, error) {
	if signed == nil || signed.Root == nil {
		return false, errors.New("signed root or inner root is nil")
	}
	if pubKeyPEM == "" {
		pubKeyPEM = signed.PublicKeyPEM
	}
	if pubKeyPEM == "" {
		return false, errors.New("public key PEM required for verification")
	}

	block, _ := pem.Decode([]byte(pubKeyPEM))
	if block == nil {
		return false, errors.New("failed to decode public key PEM block")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}

	pubKey, ok := pubInterface.(*ecdsa.PublicKey)
	if !ok {
		return false, errors.New("public key is not an ECDSA key")
	}

	sigBytes, err := hex.DecodeString(signed.SignatureHex)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature hex: %w", err)
	}

	digest := CanonicalDigest(signed.Root)
	valid := ecdsa.VerifyASN1(pubKey, digest, sigBytes)
	return valid, nil
}
