# AI-generated Design Validation Contract

Provider output is untrusted until it passes `validation.Validate`.

1. The document must be a JSON or YAML object with a non-empty `name`, a
   `components` array, and `schemaVersion` equal to
   `designs.meshery.io/v1beta1` or `designs.meshery.io/v1beta3`.
2. Credentials and operational secrets are rejected before registry hydration
   or persistence. This includes bearer tokens, passwords, API tokens,
   kubeconfigs, `client-key-data`, and PEM private keys. The validator never
   returns the secret value in an error message.
3. Instance relationships that provide `source`/`target` (or the equivalent
   `from`/`to`) must reference component IDs in the same document. Relationship
   definitions that only contain selectors are left for registry validation.
4. Validation errors are deterministic and contain a stable `code`, `path`,
   and redacted message. Malformed output is rejected before any downstream
   work.

The contract is intentionally offline. Provider calls, embeddings, registry
lookups, persistence, and deployment are outside this boundary.
