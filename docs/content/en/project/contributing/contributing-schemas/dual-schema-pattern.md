---
title: The Dual-Schema Pattern
description: Schema Design Principles and the Dual-Schema Pattern
weight: 1
---

## Schema Design Principles: The Dual-Schema Pattern
Every persisted entity in Meshery follows a strict two-schema contract. Violating this contract causes generated Go structs and API clients to be incorrect.

### Rule 1 — Entity schema = response schema only

The `<construct>.yaml` file represents the **full server-side object** as returned in API responses. It must:

- Include **all** server-generated fields: `id`, `createdAt`, `updatedAt`, `deletedAt`
- List server-generated required fields in `required` (they are always present in responses)
- Have `additionalProperties: false` at the top level

```yaml
# keychain.yaml — response schema ✅
type: object
additionalProperties: false
required:
  - id
  - name
  - owner
  - createdAt
  - updatedAt
properties:
  id:
    $ref: ../../v1alpha1/core/api.yml#/components/schemas/uuid
  name:
    type: string
  owner:
    $ref: ../../v1alpha1/core/api.yml#/components/schemas/uuid
  createdAt:
    $ref: ../../v1alpha1/core/api.yml#/components/schemas/created_at
  updatedAt:
    $ref: ../../v1alpha1/core/api.yml#/components/schemas/updated_at
  deletedAt:
    $ref: ../../v1alpha1/core/api.yml#/components/schemas/nullTime
```

### Rule 2 — Write operations use a separate `*Payload` schema

Every entity that supports `POST` or `PUT` must define a dedicated `{Entity}Payload` schema in `api.yml`. The payload schema:

- Contains **only client-settable fields** (no `createdAt`, `updatedAt`, `deletedAt`)
- Makes `id` optional with `omitempty` for upsert patterns, or omits it entirely for create-only
- Is referenced by `requestBody` in `POST`/`PUT` operations
- Is **never** reused as a response body

```yaml
# In api.yml — write schema ✅
components:
  schemas:
    KeychainPayload:
      type: object
      description: Payload for creating or updating a keychain.
      required:
        - name
      properties:
        id:
          $ref: ../../v1alpha1/core/api.yml#/components/schemas/uuid
          description: Existing keychain ID for updates; omit on create.
          x-oapi-codegen-extra-tags:
            json: "id,omitempty"
        name:
          type: string
          description: Name of the keychain.
        owner:
          $ref: ../../v1alpha1/core/api.yml#/components/schemas/uuid
          description: Owner UUID; set server-side from auth context if omitted.
          x-oapi-codegen-extra-tags:
            json: "owner,omitempty"

paths:
  /api/auth/keychains:
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/KeychainPayload"   # ← Payload, not Keychain
      responses:
        "200":
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Keychain"        # ← Full entity in response
```

### Rule 3 — Never use the entity schema as a POST/PUT request body

Using the full entity schema as a `requestBody` forces clients to supply server-generated fields (`id`, `createdAt`, `updatedAt`) and produces incorrect generated client code.

```yaml
# ❌ Wrong — exposes server-generated required fields to clients
post:
  requestBody:
    content:
      application/json:
        schema:
          $ref: "#/components/schemas/Keychain"

# ✅ Correct — separate payload type for writes
post:
  requestBody:
    content:
      application/json:
        schema:
          $ref: "#/components/schemas/KeychainPayload"
```

### Checklist when adding a new entity

- [ ] `<construct>.yaml` has `additionalProperties: false`
- [ ] `<construct>.yaml` lists all server-generated fields in `properties` and `required`
- [ ] `api.yml` defines a `{Construct}Payload` schema with only client-settable fields
- [ ] All `POST`/`PUT` `requestBody` entries reference `{Construct}Payload`, not `{Construct}`
- [ ] `GET` responses reference the full `{Construct}` entity schema

---