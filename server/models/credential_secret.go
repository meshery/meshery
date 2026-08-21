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
