package encryption_test

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/pkg/encryption"
)

// newTestKey returns a deterministic 32-byte key suitable for tests.
func newTestKey() []byte {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i + 1)
	}
	return key
}

func newTestService(t *testing.T) *encryption.Service {
	t.Helper()
	svc, err := encryption.New(newTestKey())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return svc
}

// ── New ──────────────────────────────────────────────────────────────────────

func TestNew_ValidKey(t *testing.T) {
	_, err := encryption.New(newTestKey())
	if err != nil {
		t.Fatalf("expected no error for 32-byte key, got: %v", err)
	}
}

func TestNew_ShortKey(t *testing.T) {
	_, err := encryption.New(make([]byte, 16))
	if err == nil {
		t.Fatal("expected error for 16-byte key, got nil")
	}
}

func TestNew_LongKey(t *testing.T) {
	_, err := encryption.New(make([]byte, 64))
	if err == nil {
		t.Fatal("expected error for 64-byte key, got nil")
	}
}

func TestNew_EmptyKey(t *testing.T) {
	_, err := encryption.New([]byte{})
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

// ── NewFromEnv ───────────────────────────────────────────────────────────────

func TestNewFromEnv_NoVarsReturnsNil(t *testing.T) {
	t.Setenv(encryption.EncryptionKeyEnv, "")
	t.Setenv(encryption.EncryptionKeyFileEnv, "")

	svc, err := encryption.NewFromEnv()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if svc != nil {
		t.Fatal("expected nil service when no env var is set")
	}
}

func TestNewFromEnv_HexKey(t *testing.T) {
	hexKey := hex.EncodeToString(newTestKey())
	t.Setenv(encryption.EncryptionKeyEnv, hexKey)
	t.Setenv(encryption.EncryptionKeyFileEnv, "")

	svc, err := encryption.NewFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc == nil {
		t.Fatal("expected non-nil service for valid hex key")
	}
}

func TestNewFromEnv_Base64Key(t *testing.T) {
	b64Key := base64.StdEncoding.EncodeToString(newTestKey())
	t.Setenv(encryption.EncryptionKeyEnv, b64Key)
	t.Setenv(encryption.EncryptionKeyFileEnv, "")

	svc, err := encryption.NewFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc == nil {
		t.Fatal("expected non-nil service for valid base64 key")
	}
}

func TestNewFromEnv_KeyFile(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "meshery-key-*.txt")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	hexKey := hex.EncodeToString(newTestKey())
	if _, err := fmt.Fprintln(f, hexKey); err != nil {
		t.Fatalf("write key file: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("close key file: %v", err)
	}

	t.Setenv(encryption.EncryptionKeyEnv, "")
	t.Setenv(encryption.EncryptionKeyFileEnv, f.Name())

	svc, err := encryption.NewFromEnv()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if svc == nil {
		t.Fatal("expected non-nil service when key file is valid")
	}
}

func TestNewFromEnv_InvalidKeyReturnsError(t *testing.T) {
	t.Setenv(encryption.EncryptionKeyEnv, "not-a-key")
	t.Setenv(encryption.EncryptionKeyFileEnv, "")

	_, err := encryption.NewFromEnv()
	if err == nil {
		t.Fatal("expected error for invalid key value, got nil")
	}
}

// ── EncryptMap / DecryptMap round-trip ───────────────────────────────────────

func TestEncryptDecryptRoundTrip(t *testing.T) {
	svc := newTestService(t)

	original := map[string]interface{}{
		"auth":    map[string]interface{}{"token": "super-secret"},
		"cluster": map[string]interface{}{"server": "https://k8s.example.com"},
	}

	encrypted, err := svc.EncryptMap(original)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}

	decrypted, err := svc.DecryptMap(encrypted)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}

	// Check that the round-trip preserved the auth token.
	authMap, ok := decrypted["auth"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected auth to be map, got %T", decrypted["auth"])
	}
	if got := authMap["token"]; got != "super-secret" {
		t.Errorf("auth.token = %v, want %q", got, "super-secret")
	}
}

func TestEncryptedMapContainsSentinel(t *testing.T) {
	svc := newTestService(t)
	m := map[string]interface{}{"key": "value"}

	enc, err := svc.EncryptMap(m)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	if len(enc) != 1 {
		t.Errorf("encrypted map should have exactly 1 key, got %d", len(enc))
	}
	blob, ok := enc["__enc__"].(string)
	if !ok {
		t.Fatal("sentinel key __enc__ missing or not a string")
	}
	if !strings.HasPrefix(blob, "enc:v1:") {
		t.Errorf("sentinel value should start with enc:v1:, got %q", blob)
	}
}

func TestPlaintextPassThrough(t *testing.T) {
	svc := newTestService(t)

	// A map without __enc__ key should be returned unchanged by DecryptMap.
	plain := map[string]interface{}{"foo": "bar"}
	result, err := svc.DecryptMap(plain)
	if err != nil {
		t.Fatalf("DecryptMap on plaintext: %v", err)
	}
	if result["foo"] != "bar" {
		t.Errorf("expected plaintext pass-through, got %v", result)
	}
}

func TestNonceUniqueness(t *testing.T) {
	svc := newTestService(t)
	m := map[string]interface{}{"secret": "same"}

	enc1, _ := svc.EncryptMap(m)
	enc2, _ := svc.EncryptMap(m)

	blob1 := enc1["__enc__"].(string)
	blob2 := enc2["__enc__"].(string)

	if blob1 == blob2 {
		t.Error("two encryptions of the same plaintext should produce different ciphertexts (nonce reuse)")
	}
}

func TestDoubleEncryptIsNoop(t *testing.T) {
	svc := newTestService(t)
	m := map[string]interface{}{"k": "v"}

	enc1, err := svc.EncryptMap(m)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	// Encrypting an already-encrypted map must be a no-op.
	enc2, err := svc.EncryptMap(enc1)
	if err != nil {
		t.Fatalf("EncryptMap (second): %v", err)
	}
	if enc1["__enc__"] != enc2["__enc__"] {
		t.Error("double-encrypt changed the ciphertext blob")
	}
}

func TestEmptyMapRoundTrip(t *testing.T) {
	svc := newTestService(t)
	m := map[string]interface{}{}

	enc, err := svc.EncryptMap(m)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	dec, err := svc.DecryptMap(enc)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}
	if len(dec) != 0 {
		t.Errorf("expected empty map after round-trip, got %v", dec)
	}
}

func TestTamperedCiphertextReturnsError(t *testing.T) {
	svc := newTestService(t)
	m := map[string]interface{}{"s": "secret"}

	enc, _ := svc.EncryptMap(m)

	// Corrupt the ciphertext by flipping bits in the base64 blob.
	blob := enc["__enc__"].(string)
	corrupted := blob[:len(blob)-4] + "XXXX"
	enc["__enc__"] = corrupted

	_, err := svc.DecryptMap(enc)
	if err == nil {
		t.Fatal("expected error for tampered ciphertext, got nil")
	}
}

func TestDecryptMap_NonStringSentinel(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.DecryptMap(map[string]interface{}{"__enc__": 42})
	if err == nil {
		t.Fatal("expected error for non-string sentinel value, got nil")
	}
}

func TestDecryptMap_UnknownPrefix(t *testing.T) {
	svc := newTestService(t)
	_, err := svc.DecryptMap(map[string]interface{}{"__enc__": "enc:v2:AAAA"})
	if err == nil {
		t.Fatal("expected error for unrecognised encryption prefix, got nil")
	}
}

// ── Nil receiver (disabled feature) ─────────────────────────────────────────

func TestNilServiceEncryptMapIsNoop(t *testing.T) {
	var svc *encryption.Service
	m := map[string]interface{}{"k": "v"}

	result, err := svc.EncryptMap(m)
	if err != nil {
		t.Fatalf("expected no error from nil service, got: %v", err)
	}
	if result["k"] != "v" {
		t.Error("nil service EncryptMap should return input unchanged")
	}
}

func TestNilServiceDecryptMapIsNoop(t *testing.T) {
	var svc *encryption.Service
	m := map[string]interface{}{"k": "v"}

	result, err := svc.DecryptMap(m)
	if err != nil {
		t.Fatalf("expected no error from nil service, got: %v", err)
	}
	if result["k"] != "v" {
		t.Error("nil service DecryptMap should return input unchanged")
	}
}

// ── IsEncrypted ──────────────────────────────────────────────────────────────

func TestIsEncrypted(t *testing.T) {
	svc := newTestService(t)

	plain := map[string]interface{}{"k": "v"}
	enc, _ := svc.EncryptMap(plain)

	if svc.IsEncrypted(plain) {
		t.Error("plaintext map should not be detected as encrypted")
	}
	if !svc.IsEncrypted(enc) {
		t.Error("encrypted map should be detected as encrypted")
	}
}
