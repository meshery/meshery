// Package encryption provides optional application-layer AES-256-GCM
// at-rest encryption for Meshery's sensitive datastore columns (credential
// secrets and kubeconfig auth/cluster sections).
//
// Activation
//
// Set one of the following environment variables before starting the server:
//
//	MESHERY_ENCRYPTION_KEY      – 64-char hex string (openssl rand -hex 32)
//	                              or 44-char base64 string (openssl rand -base64 32)
//	MESHERY_ENCRYPTION_KEY_FILE – path to a file whose first line is the
//	                              hex or base64 key (useful with Kubernetes Secrets)
//
// If neither variable is set, NewFromEnv returns (nil, nil) and the feature
// is disabled; all existing plaintext rows are read without modification.
//
// Wire format
//
// Encrypted maps are stored in the existing JSON column as:
//
//	{"__enc__":"enc:v1:<base64(nonce || ciphertext)>"}
//
// The 12-byte GCM nonce is prepended to the ciphertext so each value is
// independently decryptable. A map that does not contain the "__enc__" key is
// treated as plaintext and returned unchanged by DecryptMap.
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	// EncryptionKeyEnv is the environment variable for the hex/base64 key.
	EncryptionKeyEnv = "MESHERY_ENCRYPTION_KEY"

	// EncryptionKeyFileEnv is the environment variable pointing to a file
	// that contains the hex/base64 key.
	EncryptionKeyFileEnv = "MESHERY_ENCRYPTION_KEY_FILE"

	// encSentinelKey is the JSON map key that marks an encrypted value.
	encSentinelKey = "__enc__"

	// encPrefix is a version tag prepended to every encrypted blob so future
	// algorithm migrations can be detected and handled explicitly.
	encPrefix = "enc:v1:"

	// gcmNonceSize is the standard GCM nonce length in bytes.
	gcmNonceSize = 12
)

// Service encrypts and decrypts map[string]interface{} values with AES-256-GCM.
// A nil Service pointer is valid; all methods on a nil receiver are no-ops that
// return the input unchanged, so callers can safely skip nil checks:
//
//	svc, _ := encryption.NewFromEnv()
//	encrypted, err := svc.EncryptMap(secret)  // safe even when svc == nil
type Service struct {
	gcm cipher.AEAD // built once at construction; safe for concurrent use
}

// New creates a Service from the supplied raw 32-byte key.
// Returns an error if the key length is wrong.
func New(key []byte) (*Service, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes (AES-256); got %d bytes", len(key))
	}
	k := make([]byte, 32)
	copy(k, key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Service{gcm: gcm}, nil
}

// NewFromEnv reads the encryption key from environment variables.
// It checks MESHERY_ENCRYPTION_KEY first, then MESHERY_ENCRYPTION_KEY_FILE.
// Returns (nil, nil) when neither variable is set — encryption is disabled.
// Returns a non-nil error when a variable is set but the value is invalid.
func NewFromEnv() (*Service, error) {
	source := EncryptionKeyEnv
	raw := strings.TrimSpace(os.Getenv(EncryptionKeyEnv))
	if raw == "" {
		file := strings.TrimSpace(os.Getenv(EncryptionKeyFileEnv))
		if file == "" {
			// Neither variable set — feature disabled.
			return nil, nil
		}
		source = EncryptionKeyFileEnv
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("encryption: cannot read key file %q: %w", file, err)
		}
		// Take only the first line, strip whitespace.
		raw = strings.TrimSpace(strings.SplitN(string(data), "\n", 2)[0])
	}
	key, err := parseKey(raw)
	if err != nil {
		return nil, fmt.Errorf("encryption: invalid key in %s: %w", source, err)
	}
	return New(key)
}

// parseKey decodes a hex- or base64-encoded 32-byte key.
func parseKey(s string) ([]byte, error) {
	// 64 hex chars → 32 bytes
	if len(s) == 64 {
		b, err := hex.DecodeString(s)
		if err == nil && len(b) == 32 {
			return b, nil
		}
	}
	// Standard base64 with optional padding (44 chars for 32 bytes, or 43 without =)
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		// Try without padding
		b, err = base64.RawStdEncoding.DecodeString(s)
	}
	if err != nil {
		return nil, fmt.Errorf("value is neither 64-char hex nor valid base64: %w", err)
	}
	if len(b) != 32 {
		return nil, fmt.Errorf("decoded key length is %d bytes; must be 32 bytes (AES-256)", len(b))
	}
	return b, nil
}

// EncryptMap JSON-marshals m, encrypts the result with AES-256-GCM, and
// returns a sentinel map {"__enc__":"enc:v1:<base64>"}.
//
// If the receiver is nil (feature disabled), m is returned unchanged.
// If m already contains the "__enc__" key it is returned unchanged to prevent
// double-encryption.
func (s *Service) EncryptMap(m map[string]interface{}) (map[string]interface{}, error) {
	if s == nil {
		return m, nil
	}
	if _, alreadyEncrypted := m[encSentinelKey]; alreadyEncrypted {
		return m, nil
	}

	plaintext, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("encryption: marshal before encrypt: %w", err)
	}

	ciphertext, err := s.seal(plaintext)
	if err != nil {
		return nil, fmt.Errorf("encryption: seal: %w", err)
	}

	return map[string]interface{}{
		encSentinelKey: encPrefix + base64.StdEncoding.EncodeToString(ciphertext),
	}, nil
}

// DecryptMap detects a sentinel map and decrypts it back to the original
// map[string]interface{}.
//
// If the receiver is nil or m does not contain the "__enc__" key the map is
// returned unchanged (plaintext pass-through for backwards compatibility).
func (s *Service) DecryptMap(m map[string]interface{}) (map[string]interface{}, error) {
	if s == nil {
		return m, nil
	}

	enc, ok := m[encSentinelKey]
	if !ok {
		// Plaintext row — return as-is for backwards compatibility.
		return m, nil
	}

	encStr, ok := enc.(string)
	if !ok {
		return nil, fmt.Errorf("encryption: sentinel value is not a string")
	}
	if !strings.HasPrefix(encStr, encPrefix) {
		return nil, fmt.Errorf("encryption: unrecognised encryption prefix in sentinel value")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(encStr, encPrefix))
	if err != nil {
		return nil, fmt.Errorf("encryption: base64 decode: %w", err)
	}

	plaintext, err := s.open(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("encryption: open: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(plaintext, &result); err != nil {
		return nil, fmt.Errorf("encryption: unmarshal after decrypt: %w", err)
	}
	return result, nil
}

// IsEncrypted reports whether m is a sentinel map produced by EncryptMap.
func (s *Service) IsEncrypted(m map[string]interface{}) bool {
	if m == nil {
		return false
	}
	_, ok := m[encSentinelKey]
	return ok
}

// seal encrypts plaintext with AES-256-GCM and returns nonce || ciphertext.
func (s *Service) seal(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, gcmNonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("nonce generation: %w", err)
	}

	// Append ciphertext after nonce so the combined blob is self-contained.
	sealed := s.gcm.Seal(nonce, nonce, plaintext, nil)
	return sealed, nil
}

// open decrypts a nonce || ciphertext blob produced by seal.
func (s *Service) open(data []byte) ([]byte, error) {
	if len(data) < gcmNonceSize {
		return nil, fmt.Errorf("ciphertext too short: %d bytes", len(data))
	}

	nonce := data[:gcmNonceSize]
	ciphertext := data[gcmNonceSize:]

	plaintext, err := s.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("AEAD open failed (wrong key or tampered data): %w", err)
	}
	return plaintext, nil
}
