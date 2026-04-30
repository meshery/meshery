---
title: Authoring a Schema
description: Step-by-step guide to adding and modifying OpenAPI schemas in the meshery/schemas repository, including a worked example, pre-PR checklist, and common mistakes.
categories: [contributing]
weight: 30
---

This page walks through the full lifecycle of adding or modifying a schema construct: directory layout, entity definition, Payload schema, API endpoints, template files, Go helpers, TypeScript exports, and the pre-PR checklist.

{{% alert color="info" title="Read schema fundamentals first" %}}
Before authoring, read [Schema Fundamentals](../schema-fundamentals/) for the dual-schema pattern, casing rules, and HTTP design principles.
{{% /alert %}}

## Before you start {#before-you-start}

Confirm the tools are available (see the [Prerequisites](../#prerequisites) section in the index for the authoritative versions):

```bash
go version
node --version
npm --version
make --version
oapi-codegen --version  # go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

Clone and set up the repository:

```bash
git clone https://github.com/meshery/schemas.git
cd schemas
make setup   # go mod tidy + npm install
```

## Adding a new construct {#add-a-construct}

### Step 1 — Create the directory structure {#step-1}

```bash
mkdir -p schemas/constructs/v1beta1/<construct>/templates
```

### Step 2 — Create the entity schema (`<construct>.yaml`) {#step-2}

This is the **response schema** — the full server-side object. Follow the [dual-schema pattern Rule 1](../schema-fundamentals/#rule-1).

```yaml
# schemas/constructs/v1beta1/keychain/keychain.yaml
$id: https://schemas.meshery.io/keychain.yaml
$schema: http://json-schema.org/draft-07/schema#
title: Keychain
description: A keychain groups related keys for a team or organization.
additionalProperties: false
type: object
required:
  - id
  - name
  - owner
  - created_at
  - updated_at
properties:
  id:
    $ref: ../core/api.yml#/components/schemas/uuid
    description: Keychain ID.
    x-order: 1
    x-oapi-codegen-extra-tags:
      db: id
  name:
    type: string
    maxLength: 100
    description: Keychain name.
    x-order: 2
    x-oapi-codegen-extra-tags:
      db: name
  owner:
    $ref: ../core/api.yml#/components/schemas/uuid
    description: User ID of the keychain owner.
    x-order: 3
    x-oapi-codegen-extra-tags:
      db: owner
  created_at:
    $ref: ../core/api.yml#/components/schemas/created_at
    x-order: 10
  updated_at:
    $ref: ../core/api.yml#/components/schemas/updated_at
    x-order: 11
  deleted_at:
    $ref: ../core/api.yml#/components/schemas/nullTime
    description: Soft-deletion timestamp; null while active.
    x-order: 12
    x-oapi-codegen-extra-tags:
      db: deleted_at
```

Key requirements:

- `additionalProperties: false` at the top level
- All server-generated fields in `properties` (`id`, `created_at`, `updated_at`, `deleted_at`)
- Always-present fields in `required` (`id`, `created_at`, `updated_at`)
- DB-backed fields use snake_case property names matching the column name

For a complete real-world example with additional fields, see [`schemas/constructs/v1beta1/environment/environment.yaml`](https://github.com/meshery/schemas/blob/master/schemas/constructs/v1beta1/environment/environment.yaml).

### Step 3 — Create `api.yml` with endpoints and Payload schema {#step-3}

The `api.yml` file has three roles: (1) reference all subschemas, (2) define API endpoints, (3) act as entry point for code generators.

```yaml
# schemas/constructs/v1beta1/keychain/api.yml
openapi: "3.0.0"
info:
  title: Keychain API
  version: v1beta1

paths:
  /api/auth/keychains:
    get:
      operationId: getKeychains
      summary: Get all keychains
      tags:
        - Keychain
      parameters:
        - $ref: "../core/api.yml#/components/parameters/page"
        - $ref: "../core/api.yml#/components/parameters/pagesize"
      responses:
        "200":
          description: Keychains
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/KeychainPage"
    post:
      operationId: createKeychain
      summary: Create a keychain
      tags:
        - Keychain
      requestBody:
        $ref: "#/components/requestBodies/keychainPayload"
      responses:
        "201":
          description: Keychain created
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Keychain"

  /api/auth/keychains/{keychainId}:
    put:
      operationId: updateKeychain
      summary: Update a keychain
      tags:
        - Keychain
      parameters:
        - name: keychainId
          in: path
          required: true
          schema:
            type: string
            format: uuid
      requestBody:
        $ref: "#/components/requestBodies/keychainPayload"
      responses:
        "200":
          description: Keychain updated
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/Keychain"
    delete:
      operationId: deleteKeychain
      summary: Delete a keychain
      tags:
        - Keychain
      parameters:
        - name: keychainId
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        "204":
          description: Keychain deleted

components:
  schemas:
    Keychain:
      $ref: "./keychain.yaml"
    KeychainPayload:
      type: object
      description: Payload for creating or updating a keychain.
      required:
        - name
      properties:
        id:
          $ref: ../core/api.yml#/components/schemas/uuid
          description: Existing keychain ID for updates; omit on create.
          x-oapi-codegen-extra-tags:
            json: "id,omitempty"
        name:
          type: string
          maxLength: 100
          description: Keychain name.
        owner:
          $ref: ../core/api.yml#/components/schemas/uuid
          description: User ID of the keychain owner.
    KeychainPage:
      type: object
      properties:
        page:
          $ref: ../core/api.yml#/components/schemas/number
        page_size:
          $ref: ../core/api.yml#/components/schemas/number
        total_count:
          $ref: ../core/api.yml#/components/schemas/number
        keychains:
          type: array
          items:
            $ref: "#/components/schemas/Keychain"

  requestBodies:
    keychainPayload:
      description: Keychain create/update body
      required: true
      content:
        application/json:
          schema:
            $ref: "#/components/schemas/KeychainPayload"
```

Notable conventions:

- `operationId` is lower camelCase verbNoun: `getKeychains`, `createKeychain`, `updateKeychain`
- Path parameter is `keychainId` (camelCase + `Id`, not `keychainID`)
- `POST` for creation-only returns `201`; `POST` for upsert returns `200`
- `DELETE` for single resource returns `204` with no body
- Every operation has at least one `tags` entry
- Pagination list uses `page_size` and `total_count` (snake_case, published contract)

### Step 4 — Create template files {#step-4}

Template files provide example instances with valid default values. Value types must match the schema type exactly.

```json
// templates/keychain_template.json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "name": "Example Keychain",
  "owner": "00000000-0000-0000-0000-000000000000",
  "createdAt": "0001-01-01T00:00:00Z",
  "updatedAt": "0001-01-01T00:00:00Z",
  "deleted_at": null
}
```

If a field is `type: array`, its template value must be `[]`, not `{}`.

### Step 5 — Build and verify {#step-5}

```bash
make build     # validate + bundle + generate Go + generate TypeScript + test
npm run build  # build TypeScript distribution (dist/)

# Verify output
ls models/v1beta1/keychain/
ls typescript/generated/v1beta1/keychain/
```

If `make build` fails, it reports the violated rule and the file path. Fix the violation before proceeding.

### Step 6 — Update `typescript/index.ts` {#step-6}

Add your construct to the manually maintained public API surface:

```typescript
// In typescript/index.ts — add with the other v1beta1 imports

// 1. Components import (no .d.ts extension)
import { components as KeychainComponents } from "./generated/v1beta1/keychain/Keychain";

// 2. Schema import
import KeychainDefinitionV1Beta1OpenApiSchema from "./generated/v1beta1/keychain/KeychainSchema";

// 3. Type export in namespace
export namespace v1beta1 {
  // ... existing exports ...
  export type Keychain = KeychainComponents["schemas"]["Keychain"];
  export type KeychainPayload = KeychainComponents["schemas"]["KeychainPayload"];
}

// 4. Schema export
export { KeychainDefinitionV1Beta1OpenApiSchema };
```

### Step 7 — Create a Go helper file (if needed) {#step-7}

If the construct needs database CRUD operations, implement the `entity.Entity` interface in a manually maintained helper file. Always add `// This is not autogenerated.` at the top. Helper files ARE committed; generated `<construct>.go` files are NOT.

```go
// models/v1beta1/keychain/keychain_helper.go
// This is not autogenerated.
package keychain

import (
    "github.com/gofrs/uuid"
    "github.com/meshery/meshkit/database"
    "github.com/meshery/meshkit/models/meshmodel/entity"
    "gorm.io/gorm/clause"
)

func (k Keychain) TableName() string {
    return "keychains"
}

func (k Keychain) Type() entity.EntityType {
    return entity.Keychain
}

func (k *Keychain) GenerateID() (uuid.UUID, error) {
    return uuid.NewV4()
}

func (k Keychain) GetID() uuid.UUID {
    return k.Id
}

func (k *Keychain) Create(db *database.Handler, hostID uuid.UUID) (uuid.UUID, error) {
    k.Id, _ = k.GenerateID()
    return k.Id, db.Omit(clause.Associations).Create(k).Error
}
```

## Modifying an existing construct {#modify-existing}

### Adding a new field

1. Add the property to `<construct>.yaml` with correct casing, `description`, and `db:` tag if DB-backed:

```yaml
properties:
  description:
    type: string
    maxLength: 1000
    description: Human-readable description of the keychain.
    x-order: 4
    x-oapi-codegen-extra-tags:
      db: description
```

2. Update `templates/keychain_template.json` and `templates/keychain_template.yaml` with a default value.

3. Run `make build` and `npm run build`.

### Adding a new API endpoint

Edit `api.yml` to add new paths. Follow the same conventions: lower camelCase `operationId`, correct HTTP method and response code, `tags` on every operation:

```yaml
paths:
  /api/auth/keychains/{keychainId}/keys:
    get:
      operationId: getKeychainKeys
      summary: List keys in a keychain
      tags:
        - Keychain
      parameters:
        - name: keychainId
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        "200":
          description: Keys in keychain
```

## Pre-PR checklist {#pre-pr-checklist}

Run the validator before every PR:

```bash
cd ../schemas && make validate-schemas && make consumer-audit
```

Then verify each item:

- [ ] Modified only schema YAML/JSON files (not generated code in `models/` or `typescript/generated/`)
- [ ] Updated corresponding template files in `templates/` with correct default values
- [ ] Template value types match the schema type (`[]` for array, `""` for string, `{}` for object)
- [ ] Used non-deprecated `core/api.yml` references (not `core.json`)
- [ ] If adding new schemas, referenced them from `api.yml`
- [ ] Removed redundant `x-oapi-codegen-extra-tags` when the property already inherits them via `$ref`
- [ ] If a schema type is stored as a JSON blob in a single DB column AND has fixed properties, used `x-generate-db-helpers: true` at the schema component level (not per-property)
- [ ] Ran `make build` successfully
- [ ] Ran `go test ./...` successfully
- [ ] Ran `npm run build` successfully
- [ ] Verified only source schema files are in the commit (no generated files)
- [ ] If updating `typescript/index.ts`, verified import paths do not include `.d.ts` extension
- [ ] (New entity) `<construct>.yaml` has `additionalProperties: false`
- [ ] (New entity) `<construct>.yaml` includes all server-generated fields in `properties` and always-present ones in `required`
- [ ] (New entity with writes) `api.yml` defines `{Construct}Payload` with only client-settable fields
- [ ] (New entity with writes) All `POST`/`PUT` `requestBody` entries reference `{Construct}Payload`, not `{Construct}`
- [ ] (New SQL driver) `Value()` always marshals — never returns `(nil, nil)`
- [ ] (New SQL driver) `Scan()` sets the receiver to zero/nil when `src` is nil
- [ ] (New endpoint) `operationId` is lower camelCase verbNoun (`createKeychain`, not `CreateKeychain`)
- [ ] (New endpoint) Path parameters are camelCase with `Id` suffix (`{keychainId}`, not `{keychainID}`)
- [ ] (New endpoint) No `DELETE` operation has a `requestBody`; bulk deletes use `POST .../delete`
- [ ] (New `POST` for creation only) Response code is 201, not 200
- [ ] (New property) String properties have `description` and `maxLength`
- [ ] (New property) Numeric properties have `minimum`, `maximum`, or `const`
- [ ] (New property) ID properties have `format: uuid` or `$ref` to UUID type, OR `x-id-format: external` for non-UUID external identifiers
- [ ] (New property) Page-size properties have `minimum: 1`
- [ ] (New endpoint) Operation has at least one `tags` entry

## Common mistakes to avoid {#common-mistakes}

1. ❌ Hand-editing generated Go code in `models/<pkg>/<pkg>.go` — it is overwritten on next `make build`
2. ❌ Hand-editing generated TypeScript in `typescript/generated/` — same reason
3. ❌ Hand-editing built files in `dist/`
4. ❌ Using deprecated `core.json` references — use the version's `core/api.yml` instead
5. ❌ Adding redundant `x-oapi-codegen-extra-tags` on a property that already inherits them via `$ref`
6. ❌ Forgetting to update template files in `templates/` when adding new required fields
7. ❌ Not running `make build` after schema changes
8. ❌ Placing template files outside the `templates/` subdirectory
9. ❌ Using `ts` extension in TypeScript import paths in `index.ts`
10. ❌ Assuming schema property names are PascalCase — check the actual generated `.ts` files
11. ❌ Adding `x-generate-db-helpers` on individual properties — it must be at the schema component level
12. ❌ Using `x-generate-db-helpers` on amorphous types with no fixed property list — use `x-go-type: "core.Map"` instead
13. ❌ Using the full entity schema as a `POST`/`PUT` `requestBody` — always use a separate `*Payload` schema
14. ❌ Omitting `additionalProperties: false` from entity `<construct>.yaml` files
15. ❌ Adding new `Value()` implementations that return `(nil, nil)` — always marshal; return SQL NULL only when explicitly required
16. ❌ In new `Scan()` implementations, returning early without zeroing the receiver when `src` is nil — leaves stale data
17. ❌ Using PascalCase for new `operationId` values — always lower camelCase (`getPatterns`, not `GetPatterns`)
18. ❌ Using ALL-CAPS `ID` suffix in path parameters (`{orgID}`) — use camelCase + `Id` (`{orgId}`)
19. ❌ Using `DELETE` with a request body for bulk operations — use `POST /api/{resources}/delete` instead
20. ❌ Returning 200 from a `POST` that exclusively creates a new resource — use 201
21. ❌ Using all-lowercase suffix in parameter names (`workspaceid`, `pageurl`) — capitalize the suffix (`workspaceId`, `pageUrl`)
22. ❌ Template files with wrong value types — `type: array` → `[]`; `type: string` → `""`; not `{}`
23. ❌ Adding `format: uuid` to ID properties holding external system identifiers — use `x-id-format: external` instead
24. ❌ Setting `minimum: 0` on page-size properties — use `minimum: 1`
25. ❌ Omitting `tags` from operations — every operation must have at least one tag

## Intentional design decisions {#intentional-decisions}

The following patterns are deliberate. Do not flag them as issues during code review:

**`SqlNullTime` vs `NullTime`** — Some entities use `SqlNullTime` for backward compatibility with v1beta1 and downstream GORM/Pop consumers. Do not suggest switching unless the entire entity is being migrated.

**Single core Go package** — All core types live in a single unversioned package: `github.com/meshery/schemas/models/core`. Schema `x-go-type-import` for any core type must use `models/core` with alias `core`. Do not reference `models/v1alpha1/core`.

**`x-enum-casing-exempt: true`** — Enums with this annotation contain published values that will never be lowercased (e.g., `PlanName`: `"Free"`, `"Team Designer"`). Do not suggest lowercasing.

**`page_size` and `total_count`** — These use snake_case as a published API contract, not because they are DB-backed. Do not suggest renaming to `pageSize`/`totalCount`.

**Deprecated v1beta1 constructs** — Files with `x-deprecated: true` in their `info` section are kept for backward compatibility. Do not flag casing issues in deprecated constructs.

**Same field name, different casing across constructs** — A property may be `sub_type` (snake_case, DB-backed) in one construct and `subType` (camelCase, not DB-backed) in another. Both are correct. Casing is determined per-property by whether it maps to a database column in that specific construct.

**`x-id-format: external`** — ID properties annotated with this hold non-UUID external identifiers (Stripe IDs, coupon codes). Do not add `format: uuid` or remove the annotation.

**`deleted_at` in `required` lists** — Server-generated fields that are always present in API responses belong in `required`, even when the value is null for non-deleted resources.

## Contributor flows {#contributor-flows}

### mesheryctl contributor {#flow-mesheryctl}

When adding a command that uses a schema construct:

1. Verify the construct exists in `meshery/schemas`; if not, add it there first
2. Run `make build` in `meshery/schemas` to regenerate Go structs
3. Import the generated type in your mesheryctl command
4. Test: `go test --short ./...` from `mesheryctl/`

### Meshery Server contributor {#flow-server}

When adding a new field to a construct:

1. Add the property in `<construct>.yaml` with correct casing, `description`, and `db:` tag
2. Run `make build` in `meshery/schemas`
3. The generated Go struct is used directly in the server handler
4. Do not manually edit the generated struct — it will be overwritten on next regeneration

### Meshery UI contributor {#flow-ui}

When rendering a new field in the UI:

1. Verify the field exists in the schema and has been exported from `typescript/index.ts`
2. Import via the namespace: `import { v1beta1 } from "@meshery/schemas"`
3. Use RTK Query hooks from `typescript/rtk/` for data fetching
4. See [Schema-Driven UI Development](../contributing-ui-schemas/) for RJSF integration

## Next steps {#next-steps}

- [Schema Rules and Extensions](../schema-rules-and-extensions/) — all validation rules and `x-*` extensions
- [Schema Consumers and Build](../schema-consumers-and-build/) — Go helpers, TypeScript, build pipeline details
