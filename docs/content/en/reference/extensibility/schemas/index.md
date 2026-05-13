---
title: "Extensibility: Schemas"
description: Meshery schemas support x-* vendor extensions as an intentional extensibility mechanism for remote providers to attach provider-specific metadata.
aliases:
  - /extensibility/schemas
---

Meshery follows a [Schema-Driven Development](/project/contributing/contributing-schemas/) approach in which schemas define the structure of data across the platform. To support remote providers, Meshery schemas are intentionally extensible via `x-*` vendor extensions — a standard OpenAPI mechanism that allows custom metadata to be attached to a schema without affecting tools or validators that don't recognize it.

## x-annotations as an Extension Point

Meshery already uses several `x-*` vendor extensions internally to guide code generation, such as `x-go-type`, `x-oapi-codegen-extra-tags`, and `x-internal`. These are consumed by Meshery's own build pipeline.

`x-annotations` is the designated extension point for **remote providers**. It allows a provider to carry provider-specific metadata alongside any Meshery schema at runtime. Meshery core does not read, validate, or act on `x-annotations` — it passes the schema through unmodified. Only provider-side code reads them.

## Usage

Remote providers can add an `x-annotations` block to any Meshery schema. The following rules apply:

- **Additive only.** Do not remove or redefine existing core schema fields.
- **Use a provider-specific prefix** for all keys to avoid collisions with other providers (e.g., `myprovider/key` rather than just `key`).
- **Do not use `x-annotations` to conflict with or override Meshery core schema validation.**

### Example

The following shows how a remote provider annotates a Meshery connection schema with provider-specific metadata:

```yaml
type: object
properties:
  name:
    type: string
    description: Name of the connection
  endpoint:
    type: string
    description: Connection endpoint URL
x-annotations:
  acmemesh/tier: enterprise
  acmemesh/region: us-east-1
  acmemesh/billing-id: acct-00123
```

The core fields (`name`, `endpoint`) are untouched and fully compatible with Meshery core. The provider's metadata sits under `x-annotations` with its own `acmemesh/` prefix. Meshery core ignores the annotations block — only the provider's own code reads and acts on it.

### Reading x-annotations in provider code

```go
if annotations, ok := schema["x-annotations"].(map[string]interface{}); ok {
    tier, _ := annotations["acmemesh/tier"].(string)
    region, _ := annotations["acmemesh/region"].(string)
}
```

## Related

- [Extensibility: Providers](/reference/extensibility/providers/)
- [Contributing to Meshery Schemas](/project/contributing/contributing-schemas/)

{{< discuss >}}
