package models

// credential_encryption.go holds the cryptographic primitives that protect a
// persisted credential's secret material at rest. The shape-aware wrappers that
// decide *what* gets encrypted live next to the resolution rules they depend on,
// in credential_secret.go.
//
// # What this protects, and what it does not
//
// The key is derived from GlobalTokenForAnonymousResults, the build-time secret
// linked into the server binary by
// `-X main.globalTokenForAnonymousResults=$TOKEN` (install/docker/Dockerfile).
// The key therefore ships inside the binary. This protects a datastore that has
// been separated from the binary that wrote it: a stolen or exfiltrated
// database, a copied ~/.meshery directory, a filesystem backup, a snapshot, a
// support bundle. It is NOT protection against an attacker who holds the
// Meshery Server binary or image, who can recover the key from it. Say that
// plainly wherever this is documented; it is envelope encryption with a shipped
// key, not key-managed encryption.
//
// # Key identifier, and why every ciphertext carries one
//
// $TOKEN is stable in practice - it is the `GLOBAL_TOKEN` repository secret,
// passed as a build arg by the stable and edge release workflows - but it is not
// guaranteed stable. A maintainer can rotate it, and a locally built server
// (`make server`) links no token at all and falls back to the "dev_token"
// default. Both cases produce a different key, and both are reachable by an
// ordinary user: build locally today, run the released image tomorrow, against
// the same ~/.meshery.
//
// So every envelope carries a short identifier of the key that produced it.
// A read that meets an identifier it cannot match reports that specific failure
// instead of returning garbage or a bare "decryption failed", and the operator
// learns which key they need. That is bookkeeping, not key management: Meshery
// stores no keys and can decrypt only what the running binary's own token
// derives.
//
// # Scope
//
// Only DefaultLocalProvider encrypts, because only it owns the datastore the
// credential lands in. A Remote Provider persists credentials in its own store
// and must receive them as plaintext over TLS - handing it an envelope keyed to
// this binary would produce a credential nobody can ever open. Encryption
// therefore belongs in the local provider's persistence methods, not in the
// HTTP handlers above them.

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"sync"
)

const (
	// credentialEnvelopePrefix marks a string as a Meshery credential
	// ciphertext envelope. It is the sole detector for the encrypted form:
	// decryption is never attempted speculatively, and a failure to decrypt is
	// never read as "this must have been plaintext".
	credentialEnvelopePrefix = "meshery.enc.v1:"

	// credentialEnvelopeParts is the number of colon-separated fields in a
	// well-formed envelope: the "meshery.enc.v1" tag, the key id, and the
	// base64 payload.
	credentialEnvelopeParts = 3

	// credentialKDFSalt is a fixed, non-secret domain separator. HKDF's salt
	// need not be secret; its job here is to keep this key derivation distinct
	// from any other use of the same input keying material.
	credentialKDFSalt = "meshery/credential-encryption"

	// credentialKeyInfo and credentialKeyIDInfo are the HKDF `info` labels for
	// the two independent outputs derived from the same secret. Deriving the
	// identifier through HKDF rather than hashing the key keeps the identifier
	// from being a function of the key an attacker would need to confirm a
	// guess against.
	credentialKeyInfo   = "meshery credential secret encryption key v1"
	credentialKeyIDInfo = "meshery credential secret key identifier v1"

	// credentialKeyLen is the AES-256 key length in bytes.
	credentialKeyLen = 32

	// credentialKeyIDLen is the length in bytes of the derived key identifier;
	// it is stored hex-encoded, so 6 bytes render as 12 characters. It exists
	// to tell a handful of keys apart, not to be collision-proof against an
	// adversary who chooses them.
	credentialKeyIDLen = 6
)

// credentialKeyring is the derived key material for one input secret. It is
// computed once per process on first use, not at package init, because
// GlobalTokenForAnonymousResults is assigned by main() after this package's
// variables are initialised.
type credentialKeyring struct {
	keyID string
	aead  cipher.AEAD
	err   error
}

var (
	credentialKeyringOnce  sync.Once
	credentialKeyringValue *credentialKeyring
)

// activeCredentialKeyring returns the keyring derived from the token this
// binary was built with.
func activeCredentialKeyring() *credentialKeyring {
	credentialKeyringOnce.Do(func() {
		credentialKeyringValue = deriveCredentialKeyring(GlobalTokenForAnonymousResults)
	})
	return credentialKeyringValue
}

// deriveCredentialKeyring stretches a build-time secret into an AES-256 key and
// a short identifier for it. The secret is a human-chosen string of unknown
// entropy, so it is never used as key material directly.
func deriveCredentialKeyring(secret string) *credentialKeyring {
	ring := &credentialKeyring{}

	// Condense the secret to a fixed 32 bytes before HKDF. HKDF-Extract hashes
	// its input anyway, so this changes nothing cryptographically, but it keeps
	// the derivation working for a short token (the 9-byte "dev_token" default
	// among them) under `GODEBUG=fips140=only`, which rejects input keying
	// material below 112 bits.
	ikm := sha256.Sum256([]byte(secret))

	key, err := hkdf.Key(sha256.New, ikm[:], []byte(credentialKDFSalt), credentialKeyInfo, credentialKeyLen)
	if err != nil {
		ring.err = ErrCredentialKeyDerivation(err)
		return ring
	}

	id, err := hkdf.Key(sha256.New, ikm[:], []byte(credentialKDFSalt), credentialKeyIDInfo, credentialKeyIDLen)
	if err != nil {
		ring.err = ErrCredentialKeyDerivation(err)
		return ring
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		ring.err = ErrCredentialKeyDerivation(err)
		return ring
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		ring.err = ErrCredentialKeyDerivation(err)
		return ring
	}

	ring.keyID = hex.EncodeToString(id)
	ring.aead = aead
	return ring
}

// IsEncryptedCredentialEnvelope reports whether value is a Meshery credential
// ciphertext envelope. It is a prefix test on a string, deliberately: the
// encrypted and plaintext forms must be distinguishable without attempting
// decryption, so that a genuine decryption failure is always reported as a
// failure rather than mistaken for plaintext.
func IsEncryptedCredentialEnvelope(value interface{}) bool {
	str, ok := value.(string)
	return ok && strings.HasPrefix(str, credentialEnvelopePrefix)
}

// encryptCredentialEnvelope seals plaintext under the active key and renders
// the result as `meshery.enc.v1:<keyID>:<base64(nonce||ciphertext)>`.
//
// The nonce is random per call and stored with the ciphertext, so encrypting
// the same credential twice yields different envelopes and no two credentials
// share a nonce.
func encryptCredentialEnvelope(plaintext []byte) (string, error) {
	ring := activeCredentialKeyring()
	if ring.err != nil {
		return "", ring.err
	}

	nonce := make([]byte, ring.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", ErrCredentialEncrypt(err)
	}

	sealed := ring.aead.Seal(nonce, nonce, plaintext, nil)

	return credentialEnvelopePrefix + ring.keyID + ":" + base64.StdEncoding.EncodeToString(sealed), nil
}

// decryptCredentialEnvelope opens an envelope produced by
// encryptCredentialEnvelope. A key identifier that does not match this binary's
// key is reported as its own error, because it has a cause an operator can act
// on - the datastore was written by a build carrying a different $TOKEN - which
// a generic authentication failure would hide.
func decryptCredentialEnvelope(envelope string) ([]byte, error) {
	ring := activeCredentialKeyring()
	if ring.err != nil {
		return nil, ring.err
	}

	parts := strings.SplitN(envelope, ":", credentialEnvelopeParts)
	if len(parts) != credentialEnvelopeParts || parts[0]+":" != credentialEnvelopePrefix {
		return nil, ErrMalformedCredentialEnvelope("envelope is not in meshery.enc.v1:<keyId>:<payload> form")
	}

	keyID, payload := parts[1], parts[2]
	if keyID != ring.keyID {
		return nil, ErrCredentialKeyMismatch(keyID, ring.keyID)
	}

	sealed, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, ErrMalformedCredentialEnvelope("envelope payload is not valid base64")
	}
	if len(sealed) < ring.aead.NonceSize() {
		return nil, ErrMalformedCredentialEnvelope("envelope payload is shorter than the nonce it must carry")
	}

	nonce, ciphertext := sealed[:ring.aead.NonceSize()], sealed[ring.aead.NonceSize():]
	plaintext, err := ring.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrCredentialDecrypt(err)
	}
	return plaintext, nil
}
