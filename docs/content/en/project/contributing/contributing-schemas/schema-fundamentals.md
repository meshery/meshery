---
title: Schema Fundamentals
description: Schema-Driven Development in Meshery — how schemas propagate, the dual-schema pattern, casing rules, and HTTP API design principles.
categories: [contributing]
weight: 20
---

Meshery uses **Schema-Driven Development (SDD)**: every API resource is defined once as an OpenAPI schema in [`meshery/schemas`](https://github.com/meshery/schemas). That single definition propagates automatically into Go structs, TypeScript types, RTK Query clients, and runtime validation. Changes to a schema file update all downstream consumers on the next `make build`.

## Schema propagation {#schema-propagation}

```
schemas/constructs/v1beta1/<construct>/
├── <construct>.yaml        ← entity definition
└── api.yml                 ← endpoints + Payload schemas
         │
         ▼  make build
    bundle-openapi
         │
         ├─▶ generate-golang  →  models/v1beta1/<construct>/<construct>.go
         ├─▶ generate-ts      →  typescript/generated/v1beta1/<Construct>/<Construct>.ts
         └─▶ generate-rtk     →  typescript/rtk/  (RTK Query hooks)
```

The pipeline runs in four stages:

1. **bundle-openapi** — resolves all `$ref` references and merges specs into `_openapi_build/`
2. **generate-golang** — generates Go structs in `models/` via oapi-codegen
3. **generate-ts** — generates TypeScript types in `typescript/generated/` via openapi-typescript
4. **generate-rtk** — generates RTK Query client hooks in `typescript/rtk/` from the bundled spec

Stages 2–4 each consume the bundled spec from stage 1 and run independently. Run `make build` to execute all stages plus tests, or run individual targets when iterating.

## File structure for a construct {#file-structure}

Every construct follows the same directory layout:

```
schemas/constructs/v1beta1/<construct>/
├── api.yml                          # OpenAPI spec: endpoints + all schema definitions
├── <construct>.yaml                 # Entity (response) schema
└── templates/
    ├── <construct>_template.json    # Example instance (JSON)
    └── <construct>_template.yaml    # Example instance (YAML)
```
&nbsp;

| File | Role | Edit? |
|---|---|---|
| `<construct>.yaml` | Response schema — the full server-side object | Yes |
| `api.yml` | Endpoints + `{Construct}Payload` + list response wrappers | Yes |
| `templates/<construct>_template.json` | Example instance with default values | Yes |
| `templates/<construct>_template.yaml` | Same example in YAML | Yes |
| `models/v1beta1/<construct>/<construct>.go` | Generated Go structs | **Never** |
| `models/v1beta1/<construct>/<construct>_helper.go` | Manual SQL helpers + interface impl | Yes |
| `typescript/generated/v1beta1/<Construct>/<Construct>.ts` | Generated TypeScript types | **Never** |

## The dual-schema pattern {#dual-schema-pattern}

This is the most critical design rule. Every writable entity needs **two schemas** with distinct responsibilities.

### Rule 1 — `<construct>.yaml` is the response schema only {#rule-1}

The entity YAML represents the **full server-side object** as returned by API responses. It must:

- Have `additionalProperties: false` at the top level
- Define all server-generated fields in `properties`: `id`, `created_at`, `updated_at`, `deleted_at`
- List always-present server-generated fields in `required`

See the real file at [schemas/constructs/v1beta1/environment/environment.yaml](https://github.com/meshery/schemas/blob/master/schemas/constructs/v1beta1/environment/environment.yaml) for the full property set (which includes additional fields)

```yaml
# environment.yaml — simplified illustration
$id: https://schemas.meshery.io/environment.yaml
$schema: http://json-schema.org/draft-07/schema#
title: Environment
additionalProperties: false
type: object
required:
  - id
  - schemaVersion
  - name
  - description
  - organization_id
properties:
  id:
    $ref: ../core/api.yml#/components/schemas/uuid
    x-oapi-codegen-extra-tags:
      db: id
  name:
    type: string
    maxLength: 100
    description: Environment name
    x-oapi-codegen-extra-tags:
      db: name
  organization_id:
    $ref: ../core/api.yml#/components/schemas/uuid
    x-oapi-codegen-extra-tags:
      db: organization_id
  created_at:
    $ref: ../core/api.yml#/components/schemas/created_at
  updated_at:
    $ref: ../core/api.yml#/components/schemas/updated_at
  deleted_at:
    $ref: ../core/api.yml#/components/schemas/nullTime
```

### Rule 2 — `{Construct}Payload` in `api.yml` is for writes only {#rule-2}

The Payload schema is defined in `api.yml` under `components/schemas`. It:

- Contains **only client-settable fields** — never `created_at`, `updated_at`, `deleted_at`
- Makes `id` optional (`omitempty`) for upsert patterns
- Is the schema referenced by all `POST`/`PUT` `requestBody` entries

See the real `EnvironmentPayload` in [`schemas/constructs/v1beta1/environment/api.yml`](https://github.com/meshery/schemas/blob/master/schemas/constructs/v1beta1/environment/api.yml).

```yaml
# in api.yml — components/schemas
EnvironmentPayload:
  type: object
  required:
    - name
    - organization_id
  properties:
    name:
      $ref: ../core/api.yml#/components/schemas/Text
      description: Environment name.
    description:
      $ref: ../core/api.yml#/components/schemas/Text
      description: Environment description.
    organization_id:
      type: string
      format: uuid
      maxLength: 500
      description: Organization this environment belongs to.
      x-oapi-codegen-extra-tags:
        json: organization_id
```

Note that the `requestBody` object name is `environmentPayload` (camelCase) while the schema component name is `EnvironmentPayload` (PascalCase) — `requestBodies` keys are camelCase per [OpenAPI naming convention](#casing-rules), schema component names are PascalCase.

### Rule 3 — POST/PUT `requestBody` must reference `*Payload` only {#rule-3}

```yaml
# WRONG — forces clients to supply server-generated fields
post:
  requestBody:
    content:
      application/json:
        schema:
          $ref: "#/components/schemas/Environment"

# CORRECT
post:
  requestBody:
    $ref: "#/components/requestBodies/environmentPayload"
```

### Reference implementations {#reference-implementations}

Model new schemas on these constructs:

- `schemas/constructs/v1beta1/connection/` — connection + ConnectionPayload
- `schemas/constructs/v1beta1/key/` — key + KeyPayload
- `schemas/constructs/v1beta1/team/` — team + teamPayload + teamUpdatePayload
- `schemas/constructs/v1beta1/environment/` — environment + EnvironmentPayload

## Casing rules {#casing-rules}

### The core principle {#casing-core-principle}

**The OpenAPI property name IS the JSON wire-format name.** Both code generators (`oapi-codegen` for Go, `openapi-typescript` for TypeScript) use the schema property name verbatim as the JSON key. No casing transformation occurs. The casing decision is made once — in the schema — and propagates identically to Go, TypeScript, and all API consumers.

In one sentence: **wire is camelCase everywhere; DB is snake_case; the ORM layer is the sole translation boundary.**

Every property is governed by exactly one authority:

- **DB-backed field**: the database column name is the authority — use the exact snake_case column name
- **Non-DB field**: API design convention is the authority — use camelCase

### Casing reference table {#casing-table}

| Element | Casing | Example | Counter-example |
|---|---|---|---|
| DB-backed schema property | `snake_case` | `created_at`, `org_id`, `credential_id` | ❌ `createdAt`, `orgId` |
| Non-DB schema property | `camelCase` | `schemaVersion`, `displayName` | ❌ `schema_version`, `display_name` |
| Pagination envelope fields | `snake_case` (published contract) | `page_size`, `total_count` | ❌ `pageSize`, `totalCount` |
| `components/schemas` names | `PascalCase` | `Connection`, `KeychainPayload` | ❌ `connection`, `keychainPayload` |
| `operationId` | lower `camelCase` verbNoun | `createConnection`, `getPatterns` | ❌ `CreateConnection`, `get_patterns` |
| Path segments | `kebab-case` plural nouns | `/api/role-holders` | ❌ `/api/roleHolders` |
| Path parameters | `camelCase` + `Id` suffix | `{orgId}`, `{connectionId}` | ❌ `{orgID}`, `{org_id}` |
| Query parameters | `camelCase` | `pageSize`, `orgId` | ❌ `page_size`, `orgID` |
| Enum values (new) | `lowercase` | `enabled`, `ignored` | ❌ `Enabled`, `ENABLED` |
| Go struct field names | `PascalCase` (generated) | `CreatedAt`, `CredentialID` | — |
| TypeScript property names | Verbatim from schema (generated) | `credential_id`, `schemaVersion` | — |
| `json:"..."` tag | Verbatim from schema property name | `json:"credential_id"` | — |
| `db:"..."` tag | Database column name | `db:"credential_id"` | — |
| `gorm:"..."` tag | Hand-coded in helper files only | `gorm:"column:type"` | — |

### Why DB-backed fields stay snake_case {#db-casing-rationale}

The database column is the compatibility boundary. Both GORM (Meshery Server) and Buffalo Pop (Meshery Cloud) read `db:"column_name"` tags to map struct fields to columns. If the schema property were `credentialId` but the column is `credential_id`, the `json` tag and `db` tag would disagree — the API response delivers `credentialId` while the ORM writes `credential_id`. A single struct cannot serve both roles without manual translation.

**The rule**: if a field has a `db:` tag whose value is snake_case, the schema property name MUST be that exact snake_case value.

### MUST NOT {#casing-forbidden}

- MUST NOT use `ID` (ALL CAPS) in URL parameters, JSON tags, or TypeScript properties. `Id` is canonical: `{orgId}`, not `{orgID}`.
- MUST NOT partially migrate casing within a single resource. If the wire format must change, introduce a new API version (v1beta2, v1beta3) and migrate the resource completely there.
- MUST NOT introduce a new `json:` tag that matches the `db:` tag on a new DB-backed field — wire is camel, DB is snake; they differ by design.

## HTTP API design {#http-api-design}

### HTTP method semantics {#http-methods}

| Use case | Method | Response code |
|---|---|---|
| Create a new resource (creation only) | `POST` | **201** Created |
| Upsert (resource may already exist) | `POST` | **200** OK |
| Update an existing resource | `PUT` or `PATCH` | **200** OK |
| Non-CRUD action on a resource | `POST` to a sub-resource path | **200** OK |
| Bulk delete | `POST /api/{resources}/delete` | **200** OK |
| Single delete | `DELETE /api/.../{resourceId}` | **204** No Content |

### Bulk delete: never `DELETE` with a body {#bulk-delete}

REST semantics do not define a request body for `DELETE`. Many HTTP clients and proxies strip it silently. Use a `POST` sub-resource instead:

```yaml
# WRONG — DELETE with a request body
delete:
  operationId: deletePatterns
  requestBody:
    content:
      application/json:
        schema:
          $ref: '#/components/schemas/PatternIds'

# CORRECT — POST sub-resource for bulk delete
/api/content/designs/delete:
  post:
    operationId: deletePatterns
    summary: Bulk delete designs by ID
    requestBody:
      content:
        application/json:
          schema:
            $ref: '#/components/schemas/PatternIds'
    responses:
      "200":
        description: Designs deleted
```

### Response description wording {#response-wording}

Response descriptions must **not** contain the word "successfully". Use neutral wording: `Designs deleted`, `Connection updated`, `Webhook processed`.

### Resource grouping {#resource-grouping}

Endpoints are grouped under `/api` by domain:

| Prefix | Domain |
|---|---|
| `/api/identity/` | Users, orgs, roles, teams, invitations |
| `/api/integrations/` | Connections, environments, credentials |
| `/api/content/` | Designs, views, components, models |
| `/api/entitlement/` | Plans, subscriptions, features |
| `/api/auth/` | Tokens, keychains, keys |

New endpoints must be placed in the appropriate category. Path segments must be kebab-case plural nouns.

## Common schema patterns {#common-schema-patterns}

### Core schema references {#core-refs}

Always reference shared types from the version's `core/api.yml`. From a v1beta1 construct, the path is `../core/api.yml` (sibling directory). Adjust the relative path for other versions (`../../v1alpha1/core/api.yml` from a v1alpha1 sub-construct, etc.).

| Type | Reference (from v1beta1 construct) |
|---|---|
| UUID | `../core/api.yml#/components/schemas/uuid` |
| Timestamp (created) | `../core/api.yml#/components/schemas/created_at` |
| Timestamp (updated) | `../core/api.yml#/components/schemas/updated_at` |
| Nullable time | `../core/api.yml#/components/schemas/nullTime` |
| Version string | `../core/api.yml#/components/schemas/versionString` |
| Semver string | `../core/api.yml#/components/schemas/semverString` |
| Input string | `../core/api.yml#/components/schemas/inputString` |

### Field ordering with `x-order` {#x-order}

Use `x-order` to ensure fields appear in a consistent order in generated code:

```yaml
properties:
  id:
    $ref: ../core/api.yml#/components/schemas/uuid
    x-order: 1
  name:
    type: string
    x-order: 2
  created_at:
    $ref: ../core/api.yml#/components/schemas/created_at
    x-order: 10
```

### Cross-construct references {#cross-construct-refs}

When referencing another construct, use `x-go-type` and `x-go-type-import` to avoid redundant struct generation:

```yaml
model:
  $ref: ../model/api.yml#/components/schemas/ModelReference
  x-go-type: model.ModelReference
  x-go-type-import:
    path: github.com/meshery/schemas/models/v1beta1/model
  description: Reference to the parent model.
```

## Next steps {#next-steps}

- [Authoring a Schema](../authoring-a-schema/) — add a new construct step by step with a worked example
- [Schema Rules and Extensions](../schema-rules-and-extensions/) — all validation rules and vendor extensions
- [Schema Consumers and Build](../schema-consumers-and-build/) — build pipeline, Go helpers, TypeScript
