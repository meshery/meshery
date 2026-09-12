/**
 * credentialSecret.ts is the single place Meshery UI decides what a persisted
 * credential's `secret` map actually contains. Four shapes exist in production
 * and every one of them has to keep reading:
 *
 * | shape         | stored `secret`                             | payload lives      |
 * | ------------- | ------------------------------------------- | ------------------ |
 * | canonical     | `{ prometheusURL: '...' }`                  | the map itself     |
 * | kubernetes    | `{ auth: {...}, cluster: {...} }`           | the map itself     |
 * | legacy nested | `{ credentialName: 'x', secret: {...} }`    | one level down     |
 * | legacy string | `{ secret: '<token>' }`                     | a bare string      |
 *
 * The canonical shape is the one meshery/schemas declares
 * (schemas/constructs/v1beta1/credential/forms/*.json: top-level `name` plus a
 * `secret` object holding the kind-specific fields), and the one Layer5 Cloud is
 * moving to. The legacy nested shape is what Meshery UI's credential form still
 * writes. Readers must not care which they are handed, so they go through these
 * helpers rather than reaching into the map.
 *
 * server/models/credential_secret.go is the Go mirror of this module. What must
 * stay in step is the *resolution rules* - which shapes are recognized, which
 * keys make up the wrapper, and that ambiguity resolves toward canonical - not
 * the return types. The two deliberately differ on the legacy string shape:
 * `resolveCredentialPayload` returns the bare string, while Go's
 * `CredentialPayload` returns nil because it is typed to a map and Go callers
 * reach for `CredentialAuthSecret` instead. Port behaviour across, not
 * signatures.
 */

/** Loosely-typed persisted credential secret map. */
export type CredentialSecret = Record<string, unknown>;

/**
 * The only keys the legacy double-nested wrapper carries. An outer object made
 * up of nothing but these, with a `secret` entry, is a wrapper rather than a
 * payload.
 *
 * `name` is here because the registration path persists the wrapper as
 * `{name, secret}` (server/machines/actions.go), not just as `{credentialName,
 * secret}`. The cost is that a canonical payload whose only fields were `name`
 * and `secret` would be unwrapped as a wrapper. No canonical form has that
 * shape - prometheus, grafana and kubernetes carry neither field - so if you
 * add a credential kind, keep its payload from being exactly `{name, secret}`.
 */
const LEGACY_WRAPPER_KEYS = new Set(['credentialName', 'name', 'secret']);

/**
 * Canonical credential fields that hold string auth material, per
 * meshery/schemas .../credential/forms/*.json. Currently grafana only.
 * Prometheus and Kubernetes have none: a canonical Prometheus credential is
 * anonymous, and a Kubernetes credential's auth is a structured object read via
 * resolveCredentialPayload. Adding a kind whose canonical form carries string
 * auth under a new property means adding it here *and* to
 * `canonicalAuthSecretKeys` in the Go mirror.
 */
const CANONICAL_AUTH_SECRET_KEYS = ['grafanaAPIKey'];

const isPlainObject = (value: unknown): value is CredentialSecret =>
  typeof value === 'object' && value !== null && !Array.isArray(value);

/**
 * Whether `secret` is the legacy double-nested wrapper rather than a payload
 * that merely happens to carry a `secret` field. Ambiguity resolves toward
 * canonical: only an object consisting solely of wrapper keys, and carrying a
 * `secret` entry, is unwrapped.
 *
 * The `secret` probe is an own-property check to match `Object.keys` below (and
 * the Go mirror, where a map lookup has no prototype chain to fall through):
 * `'secret' in secret` would also answer true for an inherited `secret`, which
 * would unwrap an object whose own keys never included it.
 */
const isLegacyWrapper = (secret: CredentialSecret): boolean =>
  Object.prototype.hasOwnProperty.call(secret, 'secret') &&
  Object.keys(secret).every((key) => LEGACY_WRAPPER_KEYS.has(key));

/**
 * Resolves a persisted credential secret to the value carrying the credential's
 * fields: the object itself for the canonical and Kubernetes shapes, the nested
 * object for the legacy double-nested shape, and the bare string for the legacy
 * string shape.
 */
export const resolveCredentialPayload = (secret?: unknown): unknown => {
  if (!isPlainObject(secret)) {
    return secret;
  }
  return isLegacyWrapper(secret) ? secret.secret : secret;
};

/**
 * Resolves the string auth material a credential carries - an API key, a
 * service-account token, or `username:password` for basic auth - tolerating
 * every persisted shape.
 *
 * This is what the connection registration payload's `credentialSecret.secret`
 * must be: the server rehydrates it into `PromCred`/`GrafanaCred`, whose
 * `secret` field is a plain string, so handing it an object fails the register
 * step outright. `undefined` is the correct answer for credentials with no auth
 * material (an anonymous Prometheus, a Kubernetes credential) and is read
 * server-side as "no auth".
 */
export const resolveCredentialAuthSecret = (secret?: unknown): string | undefined => {
  const payload = resolveCredentialPayload(secret);

  if (typeof payload === 'string') {
    return payload;
  }
  if (!isPlainObject(payload)) {
    return undefined;
  }

  for (const key of CANONICAL_AUTH_SECRET_KEYS) {
    const value = payload[key];
    if (typeof value === 'string' && value) {
      return value;
    }
  }

  // A payload that is not a pure wrapper can still carry a string `secret`
  // field (e.g. `{ apiKey: '...', secret: '<token>' }`); keep reading it so
  // tolerance is never narrower than the pre-existing behaviour.
  return typeof payload.secret === 'string' ? payload.secret : undefined;
};
