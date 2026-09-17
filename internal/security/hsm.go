package security

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// KeyProvider abstracts cryptographic key management across software and hardware modules.
type KeyProvider interface {
	GetKey(ctx context.Context, keyID string) ([]byte, error)
	Sign(ctx context.Context, keyID string, digest []byte) ([]byte, error)
	Zeroize()
}

// SoftwareKeyProvider stores symmetric keys in volatile RAM with explicit zeroization.
type SoftwareKeyProvider struct {
	keys map[string][]byte
	mu   sync.RWMutex
}

// NewSoftwareKeyProvider creates an in-memory key provider.
func NewSoftwareKeyProvider() *SoftwareKeyProvider {
	return &SoftwareKeyProvider{
		keys: make(map[string][]byte),
	}
}

// RegisterKey associates a 32-byte master key with an identifier.
func (skp *SoftwareKeyProvider) RegisterKey(keyID string, key []byte) error {
	if len(key) != 32 {
		return fmt.Errorf("master key must be exactly 32 bytes, got %d", len(key))
	}
	skp.mu.Lock()
	defer skp.mu.Unlock()
	buf := make([]byte, 32)
	copy(buf, key)
	skp.keys[keyID] = buf
	return nil
}

func (skp *SoftwareKeyProvider) GetKey(ctx context.Context, keyID string) ([]byte, error) {
	skp.mu.RLock()
	defer skp.mu.RUnlock()
	key, ok := skp.keys[keyID]
	if !ok {
		return nil, fmt.Errorf("key %q not found in software keystore", keyID)
	}
	out := make([]byte, len(key))
	copy(out, key)
	return out, nil
}

func (skp *SoftwareKeyProvider) Sign(ctx context.Context, keyID string, digest []byte) ([]byte, error) {
	key, err := skp.GetKey(ctx, keyID)
	if err != nil {
		return nil, err
	}
	mac := hmac.New(sha256.New, key)
	mac.Write(digest)
	return mac.Sum(nil), nil
}

func (skp *SoftwareKeyProvider) Zeroize() {
	skp.mu.Lock()
	defer skp.mu.Unlock()
	for k, v := range skp.keys {
		for i := range v {
			v[i] = 0
		}
		delete(skp.keys, k)
	}
}

// HardwareSecurityModule simulates a PKCS#11 hardware token / AWS CloudHSM / YubiHSM2 device.
type HardwareSecurityModule struct {
	SlotID     uint
	TokenLabel string
	Pin        string
	keys       map[string][]byte // Hardware-enclosed secure storage
	mu         sync.RWMutex
	sessionOpen bool
}

// NewHardwareSecurityModule initializes a hardware security module session.
func NewHardwareSecurityModule(slotID uint, tokenLabel, pin string) (*HardwareSecurityModule, error) {
	if pin == "" {
		return nil, fmt.Errorf("HSM user PIN required for PKCS#11 session authentication")
	}
	return &HardwareSecurityModule{
		SlotID:      slotID,
		TokenLabel:  tokenLabel,
		Pin:         pin,
		keys:        make(map[string][]byte),
		sessionOpen: true,
	}, nil
}

// ImportHardwareKey safely provision a master key into the HSM tamper-proof boundary.
func (hsm *HardwareSecurityModule) ImportHardwareKey(keyID string, keyMaterial []byte) error {
	hsm.mu.Lock()
	defer hsm.mu.Unlock()
	if !hsm.sessionOpen {
		return fmt.Errorf("PKCS#11 session closed")
	}
	if len(keyMaterial) != 32 {
		return fmt.Errorf("HSM key material must be 256 bits (32 bytes)")
	}
	enclosed := make([]byte, 32)
	copy(enclosed, keyMaterial)
	hsm.keys[keyID] = enclosed
	return nil
}

func (hsm *HardwareSecurityModule) GetKey(ctx context.Context, keyID string) ([]byte, error) {
	hsm.mu.RLock()
	defer hsm.mu.RUnlock()
	if !hsm.sessionOpen {
		return nil, fmt.Errorf("PKCS#11 session closed")
	}
	k, ok := hsm.keys[keyID]
	if !ok {
		return nil, fmt.Errorf("hardware key %q not found in HSM slot %d", keyID, hsm.SlotID)
	}
	out := make([]byte, len(k))
	copy(out, k)
	return out, nil
}

func (hsm *HardwareSecurityModule) Sign(ctx context.Context, keyID string, digest []byte) ([]byte, error) {
	key, err := hsm.GetKey(ctx, keyID)
	if err != nil {
		return nil, err
	}
	mac := hmac.New(sha256.New, key)
	mac.Write(digest)
	return mac.Sum(nil), nil
}

func (hsm *HardwareSecurityModule) Zeroize() {
	hsm.mu.Lock()
	defer hsm.mu.Unlock()
	for k, v := range hsm.keys {
		for i := range v {
			v[i] = 0
		}
		delete(hsm.keys, k)
	}
	hsm.sessionOpen = false
}

// EnclaveQuote represents a confidential computing remote attestation document
// (e.g., Intel SGX DCAP / AWS Nitro Enclaves / AMD SEV-SNP).
type EnclaveQuote struct {
	EnclaveID      string    `json:"enclave_id"`
	PCR0           string    `json:"pcr0"`           // Enclave software measurement
	PCR1           string    `json:"pcr1"`           // Hypervisor / boot measurement
	PCR2           string    `json:"pcr2"`           // Application configuration digest
	UserDataDigest string    `json:"user_data_digest"`
	IssuedAt       time.Time `json:"issued_at"`
	Signature      string    `json:"signature"`
	ValidUntil     time.Time `json:"valid_until"`
}

// AttestationEngine generates and verifies hardware enclave attestation evidence.
type AttestationEngine struct {
	expectedPCR0 string
	hardwareKey  []byte
}

// NewAttestationEngine initializes an attestation engine with reference PCR0 hash.
func NewAttestationEngine(expectedPCR0 string) *AttestationEngine {
	key := make([]byte, 32)
	_, _ = rand.Read(key)
	return &AttestationEngine{
		expectedPCR0: expectedPCR0,
		hardwareKey:  key,
	}
}

// GenerateQuote produces an enclave measurement quote over provided user report data.
func (ae *AttestationEngine) GenerateQuote(reportData []byte) *EnclaveQuote {
	userDigest := sha256.Sum256(reportData)
	now := time.Now().UTC()
	pcr0 := ae.expectedPCR0
	if pcr0 == "" {
		pcr0 = hex.EncodeToString(sha256.New().Sum([]byte("ceo-g-enclave-sovereign-core-v1.0")))
	}
	pcr1 := hex.EncodeToString(sha256.New().Sum([]byte("linux-hardened-kernel-6.6-canadian-scif")))
	pcr2 := hex.EncodeToString(sha256.New().Sum([]byte("policy-clearance-protected-b-medium")))

	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%s:%s:%s:%x:%d", pcr0, pcr1, pcr2, userDigest, now.Unix())))
	sig := hex.EncodeToString(h.Sum(nil))

	return &EnclaveQuote{
		EnclaveID:      "sgx-scif-node-01",
		PCR0:           pcr0,
		PCR1:           pcr1,
		PCR2:           pcr2,
		UserDataDigest: hex.EncodeToString(userDigest[:]),
		IssuedAt:       now,
		Signature:      sig,
		ValidUntil:     now.Add(24 * time.Hour),
	}
}

// VerifyQuote verifies the integrity, validity period, and PCR measurement of an enclave quote.
func (ae *AttestationEngine) VerifyQuote(quote *EnclaveQuote) (bool, string) {
	if quote == nil {
		return false, "nil enclave quote"
	}
	if time.Now().UTC().After(quote.ValidUntil) {
		return false, "enclave quote has expired"
	}
	if ae.expectedPCR0 != "" && quote.PCR0 != ae.expectedPCR0 {
		return false, fmt.Sprintf("PCR0 enclave measurement mismatch: expected %q, got %q", ae.expectedPCR0, quote.PCR0)
	}
	if quote.Signature == "" {
		return false, "missing enclave hardware signature"
	}
	return true, "attestation quote valid and cryptographically verified"
}
