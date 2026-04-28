---
title: Schema Rules and Extensions
description: All schema validation rules, vendor extensions (x-* annotations), versioning, and the deprecation lifecycle in meshery/schemas.
categories: [contributing]
weight: 40
---

The `meshery/schemas` repository enforces a set of validation rules via the `validation/` Go package. Rules in the 1–36 range are **blocking** and must pass before a PR can merge. Rules in the 37–42 range are **advisory** and reported separately.

## Validator CLI targets {#validator-targets}

```bash
make validate-schemas          # blocking rules only (CI gate)
make validate-schemas-strict   # all rules including advisory
make audit-schemas             # advisory report — new issues only vs baseline
make audit-schemas-full        # full advisory backlog
make consumer-audit            # check downstream consumers (meshery/meshery, layer5io/meshery-cloud)
```

Run `make validate-schemas && make consumer-audit` before every PR.

## Blocking rules (1–36) {#blocking-rules}

These rules must pass for `make validate-schemas` to succeed. A selection of key rules:

| Rule | What it checks |
|---|---|
| 1–5 | OpenAPI structural validity — info, paths, operations present |
| 6–10 | `operationId` present on every operation; unique across the spec |
| 11–15 | `$ref` paths resolve; no circular references |
| 16–20 | Required fields present; `components/schemas` names are PascalCase |
| 21–24 | Path segments are kebab-case; path params end in `Id` |
| 25 | List endpoints reference shared pagination parameters (`page`, `pagesize`) |
| 26–28 | `POST`/`PUT` `requestBody` references `*Payload`, not the entity schema directly |
| 29–31 | Response codes match HTTP method semantics (201 for create-only POST, 204 for DELETE) |
| 32 | DB-backed property names match their `db:` tag value exactly |
| 33 | Pagination envelopes use `page_size` and `total_count`, not `pageSize`/`totalCount` |
| 34–36 | Template files exist; `additionalProperties: false` on entity YAML |

### Rule 32 — DB-backed properties must match their `db:` tag {#rule-32}

If a property has `x-oapi-codegen-extra-tags: { db: "column_name" }` and `column_name` is snake_case, then the schema property name MUST be `column_name`:

```yaml
# WRONG — property name createdAt does not match db tag created_at
createdAt:
  $ref: ../core/api.yml#/components/schemas/created_at
  x-oapi-codegen-extra-tags:
    db: created_at

# CORRECT
created_at:
  $ref: ../core/api.yml#/components/schemas/created_at
```

### Rule 33 — Pagination envelopes use snake_case {#rule-33}

```yaml
# WRONG
pageSize:
  type: integer
totalCount:
  type: integer

# CORRECT
page_size:
  type: integer
  minimum: 1
total_count:
  type: integer
  minimum: 0
```

## Advisory rules (37–42) {#advisory-rules}

These rules are reported by `make audit-schemas` but do not block CI by default. They become blocking with `make validate-schemas-strict`.

| Rule | What it checks |
|---|---|
| 37 | Every property has a `description` |
| 38 | String properties have at least one constraint: `minLength`, `maxLength`, `pattern`, `format`, or `const` |
| 39 | Numeric properties have `minimum`, `maximum`, or `const` |
| 40 | ID-like properties (`id`, `*_id`, `*Id`) have `format: uuid`, a `$ref` to a UUID type, or `x-id-format: external` |
| 41 | Page-size properties (`page_size`, `pagesize`, `pageSize`) have `minimum: 1` |
| 42 | `format` values are from the known OpenAPI 3.0 set (`date-time`, `email`, `uri`, `uuid`, etc.) — not Swagger 2.0 relics like `file` |

### Rule 37 — Every property needs a description {#rule-37}

```yaml
# WRONG — no description
name:
  type: string
  maxLength: 100

# CORRECT
name:
  type: string
  maxLength: 100
  description: Display name of the keychain.
```

### Rule 38 — String properties need a constraint {#rule-38}

At least one of: `minLength`, `maxLength`, `pattern`, `format`, or `const`:

```yaml
# WRONG — unconstrained string
label:
  type: string
  description: Resource label.

# CORRECT
label:
  type: string
  description: Resource label.
  maxLength: 255
```

### Rule 39 — Numeric properties need bounds {#rule-39}

```yaml
# WRONG — unbounded integer
retryCount:
  type: integer
  description: Number of retries attempted.

# CORRECT
retryCount:
  type: integer
  description: Number of retries attempted.
  minimum: 0
  maximum: 100
```

### Rule 40 — ID properties need `format: uuid` or `x-id-format: external` {#rule-40}

Properties matching `id`, `*_id`, or `*Id`:

```yaml
# CORRECT — using $ref to UUID type
owner_id:
  $ref: ../core/api.yml#/components/schemas/uuid
  description: Owner user ID.

# CORRECT — external (non-UUID) identifier
billing_id:
  type: string
  description: Stripe subscription ID.
  x-id-format: external
  maxLength: 500
  pattern: '^[A-Za-z0-9_\-]+$'
```

### Rule 41 — Page-size properties need `minimum: 1` {#rule-41}

A page size of zero is never valid:

```yaml
# WRONG
page_size:
  type: integer
  minimum: 0

# CORRECT
page_size:
  type: integer
  minimum: 1
```

## Vendor extensions (`x-*` annotations) {#vendor-extensions}

### `x-oapi-codegen-extra-tags` {#x-oapi-codegen-extra-tags}

Injects custom Go struct tags into generated code. Use for `db`, `json`, `yaml`, `gorm`, and related ORM tags.

```yaml
credential_id:
  $ref: ../core/api.yml#/components/schemas/uuid
  description: Associated credential ID.
  x-oapi-codegen-extra-tags:
    db: credential_id
    yaml: credential_id
```

Do **not** add this to properties that already inherit it through a `$ref` to a core schema — the tags are already defined in the referenced definition.

### `x-go-type` and `x-go-type-import` {#x-go-type}

Directs `oapi-codegen` to use a specific Go type instead of generating a new struct:

```yaml
# For the core.Map type (amorphous JSON blob)
metadata:
  type: object
  description: Additional metadata.
  x-go-type: core.Map
  x-go-type-skip-optional-pointer: true
  x-oapi-codegen-extra-tags:
    db: metadata

# For a cross-package reference
model:
  $ref: ../model/api.yml#/components/schemas/ModelReference
  x-go-type: model.ModelReference
  x-go-type-import:
    path: github.com/meshery/schemas/models/v1beta1/model
    name: model
  description: Reference to the parent model.
```

### `x-go-type-skip-optional-pointer` {#x-go-type-skip-optional-pointer}

Prevents wrapping the field in a pointer even when it is optional:

```yaml
metadata:
  x-go-type: core.Map
  x-go-type-skip-optional-pointer: true
```

### `x-go-name` {#x-go-name}

Overrides the generated Go field name:

```yaml
id:
  $ref: ../core/api.yml#/components/schemas/uuid
  x-go-name: ID
```

### `x-internal` {#x-internal}

Scopes an API operation to a specific deployment target. The build pipeline uses this to split the merged spec into `cloud_openapi.yml` and `meshery_openapi.yml`:

```yaml
paths:
  /api/entitlement/plans:
    get:
      x-internal: ["cloud"]       # cloud only
      operationId: getPlans

  /api/content/designs:
    get:
      x-internal: ["meshery"]     # Meshery OSS only
      operationId: getDesigns

  /api/integrations/connections:
    get:
      x-internal: ["cloud", "meshery"]   # both (or omit x-internal entirely)
      operationId: getConnections
```

### `x-generate-db-helpers` {#x-generate-db-helpers}

Placed at the **schema component level** (not per-property). Directs the Go generator to auto-generate `Scan()` and `Value()` SQL driver methods in `zz_generated.helpers.go`:

```yaml
components:
  schemas:
    Quiz:
      x-generate-db-helpers: true   # schema-level, NOT on a property
      type: object
      required:
        - id
        - title
      properties:
        id:
          $ref: ../../v1alpha1/core/api.yml#/components/schemas/uuid
        title:
          type: string
          maxLength: 200
          description: Quiz title.
```

Use **only** when both conditions are true:

1. The type has a dedicated schema component with **fixed, named properties**
2. The type is persisted as a **JSON blob in a single DB column** (not a relational table)

Do **not** use for amorphous types like `metadata` — use `x-go-type: "core.Map"` instead.

### `x-id-format: external` {#x-id-format-external}

Exempts an ID property from Rule 40 when it holds a non-UUID external identifier:

```yaml
billing_id:
  type: string
  description: Stripe billing ID.
  x-id-format: external
  maxLength: 500
  pattern: '^[A-Za-z0-9_\-]+$'
```

### `x-enum-casing-exempt: true` {#x-enum-casing-exempt}

Marks an enum schema as containing published values that must not be lowercased:

```yaml
PlanName:
  type: string
  x-enum-casing-exempt: true
  enum:
    - Free
    - Team Designer
    - Enterprise
```

### `x-deprecated: true` {#x-deprecated}

Placed in the `info` block to mark an entire API version as deprecated:

```yaml
info:
  title: Connection API
  version: v1beta1
  x-deprecated: true
  x-superseded-by: v1beta2/connection/api.yml
```

Deprecated files are frozen — do not modify them. Known casing violations in deprecated files are intentional and expected.

### `x-order` {#x-order}

Controls the order of fields in generated code:

```yaml
properties:
  id:
    x-order: 1
  name:
    x-order: 2
  created_at:
    x-order: 10
```

## Versioning and the deprecation lifecycle {#versioning}

### Version map

The authoritative source for what each version contains is the directory layout under [`schemas/constructs/`](https://github.com/meshery/schemas/tree/master/schemas/constructs) in the `meshery/schemas` repository.

| Version | Status | Contents |
|---|---|---|
| `v1alpha1` | Active | Shared core types (`core/`), `capability/`, `catalog_data/` and `catalog_data.yaml`, `component.yaml`, `model.yaml`, `relationship.yaml`, `selector.yaml` |
| `v1alpha2` | Active | `catalog/`, `design.yaml`, `relationship.yaml`, `selector.yaml` |
| `v1alpha3` | Active | `relationship/`, `selector/` |
| `v1beta1` | Partially deprecated | The bulk of constructs (academy, badge, capability, catalog, category, component, connection, core, credential, design, environment, evaluation, event, feature, invitation, key, keychain, model, organization, plan, relationship, role, schedule, selector, subcategory, subscription, team, token, user, view). Some are deprecated in favour of v1beta2. |
| `v1beta2` | Active | Migrated constructs with canonical camelCase wire names |

New versions are introduced selectively per construct. A construct can exist in multiple versions simultaneously while consumers migrate.

### When to introduce a new API version

The canonical trigger is a published wire-format change: if a wire-format field name (or any other published wire contract) must change, introduce a new API version (e.g., v1beta2) and migrate the resource consistently there.

Do **not** partially migrate a resource within the same version. A partial migration (renaming some fields but not others in the same construct) is forbidden. It creates an inconsistent wire format.

### Deprecation lifecycle

1. Copy the construct directory to the new version (`v1beta2/<construct>/`)
2. Apply the breaking changes in the new version
3. Add `x-deprecated: true` and `x-superseded-by` to the old version's `info` block
4. Downstream consumers must migrate to the new version before the old version can be removed
5. Do not touch deprecated files after marking them — they are frozen

### Legacy enum values

Do not recase published enum values within the same API version. For example, if `v1beta1` has `PlanName: "Free"`, that value stays as `"Free"` in `v1beta1` even though new enum values must be lowercase. The `x-enum-casing-exempt: true` annotation documents this intent for the validator.

## Next steps {#next-steps}

- [Schema Consumers and Build](../schema-consumers-and-build/) — build pipeline, Go helpers, TypeScript consumers
- [Authoring a Schema](../authoring-a-schema/) — step-by-step guide with pre-PR checklist
