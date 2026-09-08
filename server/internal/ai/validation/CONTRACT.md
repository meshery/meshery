# AI-generated Design Validation Contract

Provider output is untrusted until it passes `validation.Validate`.

1. The document must be a JSON or YAML object with a non-empty `name`, a
   `components` array, and `schemaVersion` equal to
   `designs.meshery.io/v1beta1` or `designs.meshery.io/v1beta3`.
   A single complete Markdown fence labelled `json` or `yaml` is also
   accepted; surrounding prose, incomplete fences, and multiple documents are
   rejected.
2. Each component entry must be an object containing `component.kind`,
   `component.version`, `model.name`, and `configuration`. Registry-backed
   compatibility checks remain the responsibility of the existing design
   validation and hydration pipeline.
3. Credentials and operational secrets are rejected before registry hydration
   or persistence. This includes bearer tokens, passwords, API tokens,
   kubeconfigs, `client-key-data`, and PEM private keys. The validator never
   returns the secret value in an error message.
4. Instance relationships that provide `source`/`target` (or the equivalent
   `from`/`to`) must reference component IDs in the same document. Relationship
   definitions that only contain selectors are left for registry validation;
   an explicitly empty or incorrectly typed endpoint is rejected.
5. Validation errors are deterministic and contain a stable `code`, `path`,
   and redacted message. Malformed output is rejected before any downstream
   work.

The package is a preflight trust-boundary check, not a replacement for
Meshery's existing `NewPatternFile`, registry hydration, relationship policy,
dry-run, or deployment validation. Callers must route accepted output through
those existing mechanisms in their documented order. Provider calls,
embeddings, persistence, and deployment are outside this offline boundary.
