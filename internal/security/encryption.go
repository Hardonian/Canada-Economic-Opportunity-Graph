package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// FieldEncryptor provides AES-256-GCM authenticated encryption for sensitive values.
type FieldEncryptor struct {
	key []byte
}

// NewFieldEncryptor creates an encryptor using a 32-byte master key.
func NewFieldEncryptor(key32Bytes []byte) (*FieldEncryptor, error) {
	if len(key32Bytes) != 32 {
		return nil, fmt.Errorf("encryption key must be exactly 32 bytes (256 bits), got %d", len(key32Bytes))
	}
	return &FieldEncryptor{key: key32Bytes}, nil
}

// Encrypt generates AES-256-GCM ciphertext and unique nonce.
func (fe *FieldEncryptor) Encrypt(plaintext []byte) (string, string, error) {
	block, err := aes.NewCipher(fe.key)
	if err != nil {
		return "", "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), hex.EncodeToString(nonce), nil
}

// Decrypt authenticates and decrypts AES-256-GCM ciphertext.
func (fe *FieldEncryptor) Decrypt(ciphertextHex, nonceHex string) ([]byte, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return nil, fmt.Errorf("invalid ciphertext hex: %w", err)
	}
	nonce, err := hex.DecodeString(nonceHex)
	if err != nil {
		return nil, fmt.Errorf("invalid nonce hex: %w", err)
	}

	block, err := aes.NewCipher(fe.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("authentication/decryption failure: %w", err)
	}

	return plaintext, nil
}
