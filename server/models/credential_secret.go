package models

// credential_secret.go is the single place Meshery decides what a persisted
// credential's `secret` map actually contains. Four shapes exist in production
// and every one of them has to keep reading:
//
//	canonical      {"prometheusURL": "..."}                     the secret map IS the payload
//	kubernetes     {"auth": {...}, "cluster": {...}}            the secret map IS the payload
//	legacy nested  {"credentialName": "x", "secret": {...}}     the payload is one level down
//	legacy string  {"secret": "<token>"}                        the payload is a bare string
//
// The canonical shape is the one meshery/schemas declares
// (schemas/constructs/v1beta1/credential/forms/*.json: top-level `name` plus a
// `secret` object holding the kind-specific fields), and the one Layer5 Cloud is
// moving to. The legacy nested shape is what Meshery UI's credential form still
// writes. Readers must not care which they are handed, so they go through
// CredentialPayload / CredentialAuthSecret rather than reaching into the map.
//
// ui/utils/credentialSecret.ts is the TypeScript mirror of this file. What must
// stay in step is the *resolution rules* - which shapes are recognized, which
// keys make up the wrapper, and that ambiguity resolves toward canonical - not
// the return types. The two deliberately differ on the legacy string shape:
// CredentialPayload returns nil because it is typed to a map, and Go callers
// reach for CredentialAuthSecret instead, while the TypeScript
// resolveCredentialPayload returns the bare string. Port behaviour across, not
// signatures.

import "encoding/json"

// legacyWrapperKeys are the only keys the legacy double-nested wrapper carries.
// An outer map made up of nothing but these, with an object or string under
// `secret`, is a wrapper rather than a payload.
//
// "name" is here because the registration path persists the wrapper as
// {name, secret} (server/machines/actions.go), not just as
// {credentialName, secret}. The cost is that a canonical payload whose only
// fields were "name" and "secret" would be unwrapped as a wrapper. No canonical
// form has that shape - prometheus, grafana and kubernetes carry neither field -
// so if you add a credential kind, keep its payload from being exactly
// {name, secret}.
var legacyWrapperKeys = map[string]struct{}{
	"credentialName": {},
	"name":           {},
	"secret":         {},
}

// canonicalAuthSecretKeys names the canonical credential fields that hold string
// auth material, per meshery/schemas .../credential/forms/*.json. Currently
// grafana only. Prometheus and Kubernetes have none: a canonical Prometheus
// credential is anonymous, and a Kubernetes credential's auth is a structured
// object read via CredentialPayload. Adding a kind whose canonical form carries
// string auth under a new property means adding it here *and* to
// CANONICAL_AUTH_SECRET_KEYS in the TypeScript mirror.
var canonicalAuthSecretKeys = []string{"grafanaAPIKey"}

// isLegacyWrapper reports whether secret is the legacy double-nested wrapper
// rather than a payload that merely happens to carry a `secret` field. Ambiguity
// resolves toward canonical: only a map consisting solely of wrapper keys, and
// carrying a `secret` entry, is unwrapped.
func isLegacyWrapper(secret map[string]interface{}) bool {
	if _, ok := secret["secret"]; !ok {
		return false
	}
	for key := range secret {
		if _, ok := legacyWrapperKeys[key]; !ok {
			return false
		}
	}
	return true
}

// CredentialPayload resolves a persisted credential secret map to the object
// carrying the credential's fields, unwrapping the legacy double-nested shape
// and returning the canonical and Kubernetes shapes untouched.
//
// It returns nil when the credential stores a bare string rather than an object
// (the legacy string shape) - use CredentialAuthSecret for that.
func CredentialPayload(secret map[string]interface{}) map[string]interface{} {
	if secret == nil {
		return nil
	}
	if isLegacyWrapper(secret) {
		nested, _ := secret["secret"].(map[string]interface{})
		return nested
	}
	return secret
}

// CredentialAuthSecret resolves the string auth material a credential carries -
// an API key, a service-account token, or "username:password" for basic auth -
// tolerating every persisted shape. It returns "" for credentials that carry no
// string auth material, which is the correct result for an anonymous Prometheus
// or a Kubernetes credential and is what the telemetry clients read as "no auth".
func CredentialAuthSecret(secret map[string]interface{}) string {
	if secret == nil {
		return ""
	}

	if isLegacyWrapper(secret) {
		// Legacy string shape: {"secret": "<token>"}.
		if nested, ok := secret["secret"].(string); ok {
			return nested
		}
		nested, _ := secret["secret"].(map[string]interface{})
		return authSecretFromPayload(nested)
	}

	return authSecretFromPayload(secret)
}

// authSecretFromPayload reads the string auth material out of an already
// unwrapped payload, preferring the canonical field names over the legacy
// string-valued `secret` field.
func authSecretFromPayload(payload map[string]interface{}) string {
	if payload == nil {
		return ""
	}
	for _, key := range canonicalAuthSecretKeys {
		if value, ok := payload[key].(string); ok && value != "" {
			return value
		}
	}
	// A payload that is not a pure wrapper can still carry a string `secret`
	// field (e.g. {"apiKey": "...", "secret": "<token>"}); keep reading it so
	// tolerance is never narrower than the pre-existing behaviour.
	value, _ := payload["secret"].(string)
	return value
}

// credentialCiphertextKey is the reserved property that holds a credential's
// ciphertext envelope. A persisted secret map carrying this key is encrypted;
// one that does not is plaintext, whatever else it contains. No credential form
// in meshery/schemas declares a property in this namespace, and none may: it is
// reserved for Meshery's own at-rest envelope.
//
// The marker is what makes encrypted and plaintext unambiguous to tell apart.
// Nothing in the read path infers the answer by trying to decrypt, so a
// credential that genuinely fails to decrypt is reported as a failure rather
// than mistaken for a plaintext row.
const credentialCiphertextKey = "__mesheryEncryptedSecret"

// credentialSecretEnvelope is the plaintext Meshery seals. It carries the
// resolved payload plus the one bit that resolution discards: whether the
// stored secret was the legacy double-nested wrapper.
//
// Recording that bit rather than inferring it on the way back is what makes the
// round trip exactly lossless for all four shapes. Inference cannot separate the
// wrapper {"secret": {...}} from a canonical payload that happens to be a map,
// and would rewrite the pathological-but-legal {"secret": {"secret": "tok"}}
// into something CredentialPayload resolves differently.
type credentialSecretEnvelope struct {
	// Wrapped records that the payload was found one level down, under the
	// legacy wrapper's `secret` key.
	Wrapped bool `json:"wrapped"`

	// Payload is the resolved credential payload: an object for the canonical,
	// Kubernetes and legacy nested shapes, a bare string for the legacy string
	// shape.
	Payload interface{} `json:"payload"`
}

// EncryptCredentialSecret converts a plaintext persisted credential secret into
// its encrypted form, ready to be written to the datastore.
//
// The four shapes converge here: each one is resolved to its payload first, so
// the payload is the only thing encryption ever sees and no shape is handled
// separately. What survives in plaintext alongside the ciphertext is the legacy
// wrapper's non-secret bookkeeping (`credentialName`/`name`), which is the
// credential's display name and is already stored in the clear in the
// credentials table's `name` column.
//
// The input map is never mutated: callers hand in a payload they go on to use
// (server/machines/actions.go reuses the connection payload's credential secret
// after saving), and turning that into ciphertext under them would be a
// surprise. A fresh map is returned instead.
//
// A secret arriving here already carrying the reserved ciphertext property is
// rejected, not passed through. Every read path hands its callers plaintext -
// SaveUserCredential, UpdateUserCredential, GetCredentialByID,
// GetUserCredentials and DeleteUserCredential all decrypt before returning - so
// no internal caller can produce that shape, and the only thing that can is a
// request body, which handlers unmarshal straight into Credential.Secret.
// Treating it as proof the map is already sealed would let a client store the
// rest of its secret in the clear and leave the row unreadable forever. The
// marker exists to make encrypted and plaintext unambiguous; accepting a
// caller-supplied one, silently or by stripping it, is exactly what it exists
// to prevent.
func EncryptCredentialSecret(secret map[string]interface{}) (map[string]interface{}, error) {
	if len(secret) == 0 {
		// nil and empty carry nothing to protect. Encrypting them would only
		// put an envelope in the datastore that decrypts back to nothing.
		return secret, nil
	}
	if _, reserved := secret[credentialCiphertextKey]; reserved {
		return nil, ErrReservedCredentialProperty(credentialCiphertextKey)
	}

	envelope := credentialSecretEnvelope{Payload: secret}
	encrypted := make(map[string]interface{}, 2)

	if isLegacyWrapper(secret) {
		envelope.Wrapped = true
		envelope.Payload = secret["secret"]
		for key, value := range secret {
			if key != "secret" {
				encrypted[key] = value
			}
		}
	}

	plaintext, err := json.Marshal(envelope)
	if err != nil {
		return nil, ErrCredentialEncrypt(err)
	}

	ciphertext, err := encryptCredentialEnvelope(plaintext)
	if err != nil {
		return nil, err
	}
	encrypted[credentialCiphertextKey] = ciphertext

	return encrypted, nil
}

// DecryptCredentialSecret converts a persisted credential secret back to the
// plaintext shape callers expect, and is the inverse of
// EncryptCredentialSecret.
//
// A secret that carries no ciphertext envelope is returned untouched. That is
// the upgrade path and it is deliberate: every row written before credential
// encryption shipped is plaintext, keeps reading with no migration, and
// converts to ciphertext the next time it is written.
func DecryptCredentialSecret(secret map[string]interface{}) (map[string]interface{}, error) {
	if secret == nil {
		return nil, nil
	}

	raw, ok := secret[credentialCiphertextKey]
	if !ok {
		return secret, nil
	}
	ciphertext, ok := raw.(string)
	if !ok || !IsEncryptedCredentialEnvelope(ciphertext) {
		// The reserved key is present but does not hold an envelope. Something
		// wrote a shape Meshery does not produce; refuse it rather than hand
		// the caller a "payload" that is really Meshery's own marker.
		return nil, ErrMalformedCredentialEnvelope("reserved property " + credentialCiphertextKey + " does not hold a credential ciphertext envelope")
	}

	plaintext, err := decryptCredentialEnvelope(ciphertext)
	if err != nil {
		return nil, err
	}

	var envelope credentialSecretEnvelope
	if err := json.Unmarshal(plaintext, &envelope); err != nil {
		return nil, ErrCredentialDecrypt(err)
	}

	if !envelope.Wrapped {
		payload, ok := envelope.Payload.(map[string]interface{})
		if !ok {
			return nil, ErrMalformedCredentialEnvelope("unwrapped credential payload decrypted to a non-object")
		}
		return payload, nil
	}

	decrypted := make(map[string]interface{}, len(secret))
	for key, value := range secret {
		if key != credentialCiphertextKey {
			decrypted[key] = value
		}
	}
	decrypted["secret"] = envelope.Payload

	return decrypted, nil
}

// IsEncryptedCredentialSecret reports whether a persisted credential secret is
// stored encrypted.
func IsEncryptedCredentialSecret(secret map[string]interface{}) bool {
	if secret == nil {
		return false
	}
	return IsEncryptedCredentialEnvelope(secret[credentialCiphertextKey])
}
