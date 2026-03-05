---
title: Contributing to Meshery Schemas
description: How to contribute to Meshery Schemas
categories: [contributing]
---

This comprehensive guide covers everything you need to know to contribute to the Meshery Schemas repository. Meshery follows **Schema-Driven Development (SDD)**, where the structure of data is centrally defined using schemas that power consistency, validation, and code generation across the platform.

---

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Quick Start](#quick-start)
4. [Schema Directory Structure](#schema-directory-structure)
5. [Understanding `api.yml` - The Construct Index File](#understanding-apiyml---the-construct-index-file)
6. [Schema File Roles](#schema-file-roles)
7. [Naming Conventions](#naming-conventions)
8. [Adding a New Schema](#adding-a-new-schema)
9. [Modifying Existing Schemas](#modifying-existing-schemas)
10. [Code Generation](#code-generation)
11. [Go Helper Files](#go-helper-files)
12. [TypeScript Integration](#typescript-integration)
13. [Common Schema Patterns](#common-schema-patterns)
14. [Template Files](#template-files)
15. [What NOT to Commit](#what-not-to-commit)
16. [Testing Your Changes](#testing-your-changes)
17. [Common Mistakes to Avoid](#common-mistakes-to-avoid)
18. [Checklist for Schema Changes](#checklist-for-schema-changes)
19. [Getting Help](#getting-help)

---

## Overview

Meshery schemas offer a powerful system designed for:

- **Model-Driven Management**: Meshery uses explicit models for describing infrastructure and applications.
- **Dynamic Discovery**: Process different kinds of relationships and styles for adaptive configurations.
- **Lifecycle Management**: Track status and lifecycle of resources via schema properties.
- **Extensibility**: Open-ended metadata and modular schema components enable expansion.
- **Visual Representation**: Properties for styling edges and nodes create user-friendly visuals.
- **Automated Operations**: Support validation, automated configuration, and patching.

Meshery uses the **OpenAPI v3** specification with a **modular, versioned, and extensible** schema strategy:

- ✅ **Versioned schemas** for backward compatibility
- 🧩 **Modular constructs** for maintainability and reuse
- 🧪 **Schemas are used** for validation, API documentation, and automatic code generation

---

## Prerequisites

Before contributing, ensure you have the following installed:

### 1. Go (v1.24.0+)

```bash
# Verify installation
go version
```

### 2. oapi-codegen

Essential for generating Go code from OpenAPI specifications:

```bash
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
export PATH="${GOPATH:-$HOME/go}/bin:$PATH"
```

### 3. Node.js & npm

Required for TypeScript generation and build process:

```bash
# Verify installation
node --version
npm --version
```

### 4. make

The repository uses Makefiles for automation:

```bash
# Verify installation
make --version
```

---

## Quick Start

```bash
# Clone the repository
git clone https://github.com/meshery/schemas.git
cd schemas

# Install dependencies
make setup
npm install

# Generate all code (Go, TypeScript, RTK Query)
make build

# Build TypeScript distribution
npm run build
```

---

## Schema Directory Structure

All schemas are located in the `schemas/constructs/` directory:

```
schemas/
├── constructs/
│   ├── <schema-version>/               # e.g., v1alpha1, v1beta1
│   │   └── <construct>/                # e.g., model, component, design
│   │       ├── api.yml                 # Index file: refs subschemas + defines API endpoints
│   │       ├── <construct>.yaml        # Subschema: data model definition
│   │       ├── <construct>_core.yml    # Subschema: core/shared types (optional)
│   │       └── templates/              # Manually defined template files
│   │           ├── <construct>_template.json
│   │           └── <construct>_template.yaml
│   │
│   ├── v1alpha1/
│   │   ├── core/
│   │   │   └── api.yml                 # Core schema definitions (timestamps, UUIDs, etc.)
│   │   └── capability/
│   │       └── api.yml
│   │
│   ├── v1alpha3/
│   │   └── relationship/
│   │       ├── api.yml
│   │       ├── relationship_core.yml
│   │       └── templates/
│   │
│   └── v1beta1/
│       ├── model/
│       │   ├── api.yml
│       │   ├── model.yaml
│       │   ├── model_core.yml
│       │   └── templates/
│       ├── component/
│       │   ├── api.yml
│       │   ├── component.yaml
│       │   └── templates/
│       ├── design/
│       ├── environment/
│       ├── workspace/
│       └── ...
│
├── models/                             # Auto-generated Go code (do NOT commit)
│   └── <version>/<package>/<package>.go
│
├── typescript/
│   ├── index.ts                        # Manually maintained - public API surface
│   ├── generated/                      # Auto-generated TypeScript (do NOT commit)
│   │   └── <version>/<package>/
│   │       ├── <Package>.d.ts          # Type definitions
│   │       └── <Package>Schema.ts      # Schema as JS object
│   └── rtk/                            # RTK Query client configurations
│
├── dist/                               # Built distribution (do NOT commit)
│   ├── index.js, index.d.ts
│   ├── cloudApi.js, mesheryApi.js
│   └── constructs/<version>/<package>/<Package>Schema.js
│
└── _openapi_build/                     # Bundled OpenAPI specs (do NOT commit)
    ├── merged_openapi.yml
    ├── cloud_openapi.yml
    └── meshery_openapi.yml
```

---

## Understanding `api.yml` - The Construct Index File

Each construct directory contains an `api.yml` file that serves as the **index file** for that construct. This is the entry point for code generation tools.

### The Three Roles of `api.yml`

1. **References all subschemas** - Aggregates and references all related schema definitions via `$ref`
2. **Defines API endpoints** - Contains all REST operations (GET, POST, PUT, DELETE) for the construct
3. **Acts as the entry point** - Used by code generators (oapi-codegen, openapi-typescript)

### Example `api.yml` Structure

```yaml
openapi: 3.0.0
info:
  title: Model API
  version: v1beta1

paths:
  /api/models:
    get:
      operationId: getModels
      summary: Get all models
      responses:
        "200":
          description: Success
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/ModelDefinition"
    post:
      operationId: createModel
      summary: Create a new model
      # ...

components:
  schemas:
    ModelDefinition:
      $ref: "./model.yaml#/ModelDefinition"    # Reference to subschema
    
    ModelReference:
      $ref: "./model_core.yml#/ModelReference"  # Reference to another subschema
```

---

## Schema File Roles

| File | Purpose |
|------|---------|
| `api.yml` | **Index file** - aggregates all subschemas via `$ref` and defines API endpoints for the construct |
| `<construct>.yaml` | **Subschema** - defines the main data model (noun) for the construct |
| `<other-subschemas>.yml` | **Subschema** - defines core/shared types used by the main schema |
| `templates/*.json` | **Templates** - example instances with default values |

---

## Naming Conventions

### Property Names

- Use **camelCase** for property fields: `schemaVersion`, `displayName`, `componentsCount`
- Identifier fields use **lowerCamelCase with "Id" suffix**: `modelId`, `registrantId`, `categoryId`
- Enums use **lowercase words**: `enabled`, `ignored`, `duplicate`

### OpenAPI Schema Names

- **PascalCase nouns** under `components/schemas`: `Model`, `Component`, `Design`
- **Files/folders are lowercase**: `api.yml`, `model.yaml`, `component.yaml`
- Template files: `templates/<construct>_template.json`

### Endpoints and Operations

- Paths under `/api` with **kebab-case, plural nouns**: `/api/workspaces`, `/api/environments`
- Path params are **camelCase**: `{subscriptionId}`, `{connectionId}`
- Non-CRUD actions append a verb segment: `.../register`, `.../export`, `.../cancel`
- `operationId` is **camelCase VerbNoun**: `getModels`, `createDesign`, `registerMeshmodels`

### Versioning

- `schemaVersion` uses **group/version**: `models.meshery.io/v1beta1`, `components.meshery.io/v1beta1`
- Version strings follow **k8s-style**: `v1`, `v1alpha1`, `v1beta1`
- Semver fields use **standard SemVer**: `1.0.0`, `2.3.1`

---

## Adding a New Schema

### Step 1: Create the Directory Structure

```bash
mkdir -p schemas/constructs/v1beta1/mypackage/templates
```

### Step 2: Create the Index File (`api.yml`)

Create `schemas/constructs/v1beta1/mypackage/api.yml`:

```yaml
openapi: 3.0.0
info:
  title: MyPackage API
  version: v1beta1
  description: API for managing MyPackage resources

paths:
  /api/mypackages:
    get:
      operationId: getMyPackages
      summary: Get all mypackages
      tags:
        - MyPackage
      responses:
        "200":
          description: Success
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: "#/components/schemas/MyPackage"
    post:
      operationId: createMyPackage
      summary: Create a new mypackage
      tags:
        - MyPackage
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/MyPackage"
      responses:
        "201":
          description: Created

components:
  schemas:
    MyPackage:
      $ref: "./mypackage.yaml#/MyPackage"
```

### Step 3: Create Subschema Files (Optional)

Create `schemas/constructs/v1beta1/mypackage/mypackage.yaml`:

```yaml
MyPackage:
  type: object
  required:
    - id
    - name
  properties:
    id:
      $ref: ../../v1alpha1/core/api.yml#/components/schemas/uuid
      x-order: 1
    
    name:
      type: string
      description: Name of the package
      minLength: 1
      maxLength: 100
      x-order: 2
    
    description:
      type: string
      description: Description of the package
      x-order: 3
    
    created_at:
      $ref: ../../v1alpha1/core/api.yml#/components/schemas/created_at
      x-order: 10
    
    updated_at:
      $ref: ../../v1alpha1/core/api.yml#/components/schemas/updated_at
      x-order: 11
```

### Step 4: Create Template Files

Create `schemas/constructs/v1beta1/mypackage/templates/mypackage_template.json`:

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "name": "Example Package",
  "description": "An example package instance",
  "created_at": "0001-01-01T00:00:00Z",
  "updated_at": "0001-01-01T00:00:00Z"
}
```

### Step 5: Build and Test

```bash
# Run the build - your schema will be automatically discovered
make build
npm run build

# Verify Go code was generated
ls models/v1beta1/mypackage/

# Verify TypeScript was generated
ls typescript/generated/v1beta1/mypackage/
```

---

## Modifying Existing Schemas

### Adding a New Field

1. **Update the schema / yml file** (e.g., `model.yaml or api.yml`):

```yaml
properties:
  # ... existing properties
  
  newField:
    type: string
    description: Description of the new field
    x-order: 20
    x-oapi-codegen-extra-tags:
      yaml: "newField"
      json: "newField"
```

2. **Update template files** with default values:

```json
{
  "newField": "default value"
}
```

3. **Run the build**:

```bash
make build
npm run build
```

### Adding a New API Endpoint

Edit the `api.yml` file to add new paths:

```yaml
paths:
  # ... existing paths
  
  /api/models/{id}/export:
    post:
      operationId: exportModel
      summary: Export a model
      parameters:
        - $ref: "../../v1alpha1/core/api.yml#/components/parameters/id"
      responses:
        "200":
          description: Model exported successfully
```

---

## Code Generation

The build system generates code from your schemas automatically.

### What Gets Generated

| Output | Location | Description |
|--------|----------|-------------|
| **Go structs** | `models/<version>/<package>/` | Strongly-typed models for backend |
| **TypeScript types** | `typescript/generated/<version>/<package>/<Package>.d.ts` | Interface definitions |
| **TypeScript schemas** | `typescript/generated/<version>/<package>/<Package>Schema.ts` | OpenAPI schema as const JS object |
| **RTK Query clients** | `typescript/rtk/` | Auto-generated API hooks for Redux |
| **Bundled OpenAPI** | `_openapi_build/` | Merged API specifications |

### Build Commands

| Command | Description |
|---------|-------------|
| `make build` | Full build: bundles OpenAPI + generates Go/TypeScript |
| `make bundle-openapi` | Bundle and merge OpenAPI specs only |
| `make generate-golang` | Generate Go code (requires bundled specs) |
| `make generate-ts` | Generate TypeScript types and schemas |
| `npm run build` | Build TypeScript distribution with tsup |

### Dynamic Schema Discovery

Schemas are discovered automatically by scanning `schemas/constructs/` for directories containing an `api.yml` file. No manual configuration needed!

---

## Go Helper Files

While Go structs are auto-generated from schemas, you often need to add **custom methods** to make these structs compatible with databases, implement interfaces, or add utility functions. This is done through manually created **helper files**.

### When to Create Helper Files

Create a helper file (`*_helper.go` or `helpers.go`) in the generated package when you need:

| Use Case | Description |
|----------|-------------|
| **SQL Driver Compatibility** | Implement `database/sql/driver.Scanner` and `driver.Valuer` interfaces |
| **Entity Interface** | Implement the `entity.Entity` interface for database CRUD operations |
| **GORM Table Names** | Define custom table names via `TableName()` method |
| **Utility Methods** | Add helper functions for serialization, validation, or business logic |
| **Type Conversions** | Add methods to convert between related types |

### Helper File Location

```
models/
├── core/
│   ├── core.go                    # Auto-generated (do NOT edit)
│   ├── helpers.go                 # Manual: utility functions
│   ├── datatype_map.go            # Manual: Map type with SQL driver methods
│   └── datatype_null_time.go      # Manual: NullTime with SQL driver methods
├── v1beta1/
│   ├── model/
│   │   ├── model.go               # Auto-generated (do NOT edit)
│   │   └── model_helper.go        # Manual: Entity interface, TableName, etc.
│   ├── component/
│   │   ├── component.go           # Auto-generated (do NOT edit)
│   │   └── component_helper.go    # Manual: Entity interface, TableName, etc.
│   └── category/
│       ├── category.go            # Auto-generated (do NOT edit)
│       └── category_helper.go     # Manual: Entity interface, TableName, etc.
```

### SQL Driver Interface Implementation

To store complex types (structs, maps, slices) in SQL databases, implement `Scan` and `Value` methods:

```go
// helpers.go - This is NOT autogenerated
package mypackage

import (
    "database/sql/driver"
    "encoding/json"
    
    "github.com/meshery/schemas/models/core"
)

// Scan implements sql.Scanner interface for reading from database
func (m *MyComplexType) Scan(value interface{}) error {
    mapVal := core.Map{}
    err := mapVal.Scan(value)
    if err != nil {
        return err
    }
    return core.MapToStruct(mapVal, m)
}

// Value implements driver.Valuer interface for writing to database
func (m MyComplexType) Value() (driver.Value, error) {
    mapVal, err := core.StructToMap(m)
    if err != nil {
        return nil, err
    }
    return core.Map(mapVal).Value()
}
```

### Entity Interface Implementation

For structs that need database CRUD operations, implement the `entity.Entity` interface:

```go
// component_helper.go - This is NOT autogenerated
package component

import (
    "fmt"
    "sync"
    
    "github.com/gofrs/uuid"
    "github.com/meshery/meshkit/database"
    "github.com/meshery/meshkit/models/meshmodel/entity"
    "gorm.io/gorm/clause"
)

// TableName returns the database table name for GORM
func (c ComponentDefinition) TableName() string {
    return "component_definition_dbs"
}

// Type returns the entity type identifier
func (c ComponentDefinition) Type() entity.EntityType {
    return entity.ComponentDefinition
}

// GenerateID generates a new UUID for the entity
func (c *ComponentDefinition) GenerateID() (uuid.UUID, error) {
    return uuid.NewV4()
}

// GetID returns the entity's ID
func (c ComponentDefinition) GetID() uuid.UUID {
    return c.Id
}

// GetEntityDetail returns a human-readable description
func (c *ComponentDefinition) GetEntityDetail() string {
    return fmt.Sprintf("type: %s, name: %s, model: %s", 
        c.Type(), c.DisplayName, c.Model.Name)
}

// Create inserts the entity into the database
func (c *ComponentDefinition) Create(db *database.Handler, hostID uuid.UUID) (uuid.UUID, error) {
    c.Id, _ = c.GenerateID()
    err := db.Omit(clause.Associations).Create(&c).Error
    return c.Id, err
}

// UpdateStatus updates the entity's status in the database
func (c *ComponentDefinition) UpdateStatus(db *database.Handler, status entity.EntityStatus) error {
    return nil
}
```

### Deterministic ID Generation

For entities that need consistent IDs based on their content (to prevent duplicates):

```go
// model_helper.go
package model

import (
    "crypto/md5"
    "encoding/hex"
    "encoding/json"
    
    "github.com/gofrs/uuid"
)

func (m *ModelDefinition) GenerateID() (uuid.UUID, error) {
    // Create identifier from unique fields
    modelIdentifier := ModelDefinition{
        Registrant:    m.Registrant,
        Version:       m.Version,
        SchemaVersion: m.SchemaVersion,
        Name:          m.Name,
        Model: Model{
            Version: m.Model.Version,
        },
    }
    byt, err := json.Marshal(modelIdentifier)
    if err != nil {
        return uuid.UUID{}, err
    }
    
    hash := md5.Sum(byt)
    return uuid.FromString(hex.EncodeToString(hash[:]))
}
```

### Type Alias and Conversion Helpers

Add type aliases and conversion methods for convenience:

```go
// model_helper.go
package model

// Type alias for cleaner code
type Styles = ComponentDefinition_Styles

// ToReference converts full definition to a lightweight reference
func (m ModelDefinition) ToReference() ModelReference {
    return ModelReference{
        Name:        m.Name,
        Version:     m.Version,
        DisplayName: m.DisplayName,
        Model:       m.Model,
        Registrant: RegistrantReference{
            Kind: m.Registrant.Kind,
        },
    }
}
```

### Core Utility Types

The `models/core/` package provides reusable types with built-in SQL compatibility:

| Type | Purpose | Use Case |
|------|---------|----------|
| `core.Map` | `map[string]any` with SQL support | Storing JSON objects in database |
| `core.NullTime` | Nullable time with JSON/YAML support | Optional timestamp fields (e.g., `deleted_at`) |
| `core.Time` | Time wrapper with custom formatting | Required timestamp fields |

**Using Core Types:**

```go
package mypackage

import "github.com/meshery/schemas/models/core"

// For nullable timestamps (e.g., deleted_at)
type MyStruct struct {
    DeletedAt core.NullTime `json:"deleted_at" gorm:"column:deleted_at"`
}

// For JSON metadata stored as blob
type MyStruct struct {
    Metadata core.Map `json:"metadata" gorm:"type:bytes;serializer:json"`
}
```



### Important Notes for Helper Files

1. **File Header Comment**: Always add `// This is not autogenerated.` at the top
2. **Same Package**: Helper files must be in the same package as the generated code
3. **DO Commit Helper Files**: Unlike generated `.go` files, helper files ARE committed to the repository
4. **Naming Convention**: Use `<package>_helper.go` or `helpers.go`
5. **Pointer vs Value Receivers**: Use pointer receivers for methods that modify the struct


## TypeScript Integration

### Using Generated Types

```typescript
import { v1beta1, v1alpha1 } from "@meshery/schemas";

const component: v1beta1.Component = { /* ... */ };
const model: v1beta1.Model = { /* ... */ };
const design: v1beta1.Design = { /* ... */ };
```

### Using Schema Exports

```typescript
// From main index
import {
  ModelDefinitionV1Beta1OpenApiSchema,
  ComponentDefinitionV1Beta1OpenApiSchema,
} from "@meshery/schemas";

// Direct import
import ModelSchema from "@meshery/schemas/dist/constructs/v1beta1/model/ModelSchema";
```

### Maintaining `typescript/index.ts`

The `typescript/index.ts` file is **manually maintained** and defines the public API surface. When adding new constructs:

1. Import components from the generated `.d.ts` file
2. Import the schema from the generated `*Schema.ts` file
3. Add type exports to the appropriate namespace

```typescript
// Type imports (no .d.ts extension)
import { components as MyPackageComponents } from "./generated/v1beta1/mypackage/MyPackage";

// Schema imports
import MyPackageV1Beta1Schema from "./generated/v1beta1/mypackage/MyPackageSchema";

// Export in namespace
export namespace v1beta1 {
  export type MyPackage = MyPackageComponents["schemas"]["MyPackage"];
}

// Export schema
export { MyPackageV1Beta1Schema };
```

---

## Common Schema Patterns

### Core Schema References

Always use the non-deprecated references from `v1alpha1/core/api.yml`:

| Type | Reference |
|------|-----------|
| **UUID** | `../../v1alpha1/core/api.yml#/components/schemas/uuid` |
| **Timestamp (created)** | `../../v1alpha1/core/api.yml#/components/schemas/created_at` |
| **Timestamp (updated)** | `../../v1alpha1/core/api.yml#/components/schemas/updated_at` |
| **Version String** | `../../v1alpha1/core/api.yml#/components/schemas/versionString` |
| **Semver String** | `../../v1alpha1/core/api.yml#/components/schemas/semverString` |
| **Input String** | `../../v1alpha1/core/api.yml#/components/schemas/inputString` |

### Timestamp Fields Pattern

```yaml
properties:
  createdAt:
    $ref: ../../v1alpha1/core/api.yml#/components/schemas/created_at
    x-order: 14
  
  updatedAt:
    $ref: ../../v1alpha1/core/api.yml#/components/schemas/updated_at
    x-order: 15
```



### Preserving Field Order with `x-order`

Use the `x-order` tag to ensure fields appear in a specific order in generated code:

```yaml
properties:
  id:
    type: string
    x-order: 1
  name:
    type: string
    x-order: 2
  description:
    type: string
    x-order: 3
```

### Referencing Other Constructs

When referencing models or other constructs, add `x-go-type` and `x-go-import-path` to avoid generating redundant Go structs:

```yaml
model:
  $ref: ../model/api.yml#/components/schemas/ModelReference
  x-go-type: model.ModelReference
  x-go-type-import:
    path: github.com/meshery/schemas/models/v1beta1/model
  description: Reference to the model
```

### Annotating API Paths for Filtering

Use `x-internal` to control which bundled output includes the path:

```yaml
paths:
  /api/entitlement/plans:
    get:
      x-internal: ["cloud"]  # Only included in cloud_openapi.yml
      operationId: getPlans
      # ...
```

- **With `x-internal`**: Included only in the specified clients
- **Without `x-internal`**: Included in **all** clients

---

## Template Files

Templates are manually defined files in the `templates/` subdirectory. They provide example instances with default values.

### Template Structure

```
constructs/v1beta1/model/templates/
├── model_template.json         # Default JSON template
├── model_template.yaml         # Default YAML template
├── model_minimal_template.json # Minimal variant (optional)
└── model_full_template.yaml    # Full variant (optional)
```

### Template Content Example

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "schemaVersion": "models.meshery.io/v1beta1",
  "version": "1.0.0",
  "name": "example-model",
  "displayName": "Example Model",
  "description": "An example model template",
  "created_at": "0001-01-01T00:00:00Z",
  "updated_at": "0001-01-01T00:00:00Z"
}
```

### When to Update Templates

- When adding new required fields to a schema
- When changing default values
- When adding new optional fields that should have example values

---

## What NOT to Commit

**CRITICAL**: Do not commit generated files. Only commit source schema files.

### ❌ Do NOT Commit

| Directory/File | Description |
|----------------|-------------|
| `models/<version>/<package>/<package>.go` | Auto-generated Go structs |
| `typescript/generated/` | Generated TypeScript types and schemas |
| `dist/` | Built distribution files |
| `_openapi_build/` | Bundled OpenAPI specs |
| `merged_openapi.yml` | Generated merged spec |
| `cloud_openapi.yml` | Generated cloud spec |
| `meshery_openapi.yml` | Generated meshery spec |

### ✅ DO Commit

| Directory/File | Description |
|----------------|-------------|
| `constructs/<version>/<package>/api.yml` | Index file for each construct |
| `constructs/<version>/<package>/*.yaml` | Subschema files |
| `constructs/<version>/<package>/*.json` | Schema files in JSON format |
| `constructs/<version>/<package>/templates/` | Template files |
| `typescript/index.ts` | Manually maintained public API |
| `models/<version>/<package>/*_helper.go` | **Manual** Go helper files (SQL drivers, Entity interface) |
| `models/<version>/<package>/helpers.go` | **Manual** Go utility functions |
| `models/core/datatype_*.go` | **Manual** Core data type definitions |

> **Note**: The `models/` directory contains both auto-generated files (e.g., `model.go`) and manually created helper files (e.g., `model_helper.go`). Only the auto-generated struct files should NOT be committed. Helper files that implement interfaces, SQL drivers, and utility methods ARE committed and maintained manually.

---

## Testing Your Changes

### Run Full Build

```bash
make build
npm run build
```

### Run Go Tests

```bash
go test ./...
```



### Verify Generated Code

```bash
# Check Go code compiles
go build ./...

# Check TypeScript compiles
npm run build
```

---

## Common Mistakes to Avoid

1. ❌ **Committing auto-generated Go structs** (`models/<pkg>/<pkg>.go`) - but DO commit helper files
2. ❌ **Using deprecated `core.json`** references instead of `v1alpha1/core/api.yml`
3. ❌ **Adding redundant `x-oapi-codegen-extra-tags`** when using core schema references
4. ❌ **Forgetting to update template files** when adding new fields
5. ❌ **Not testing the build** after schema changes
6. ❌ **Placing template files outside** the `templates/` subdirectory
7. ❌ **Using `.d.ts` extension** in TypeScript import paths
8. ❌ **Assuming schema property names are PascalCase** (check actual generated files)
9. ❌ **Not adding new schemas to `typescript/index.ts`** for public API exposure
10. ❌ **Forgetting `x-go-type`** when referencing other constructs
11. ❌ **Editing auto-generated `.go` files** instead of creating helper files
12. ❌ **Forgetting `// This is not autogenerated.`** comment in helper files
13. ❌ **Missing `TableName()` method** in helper files for GORM entities
14. ❌ **Not implementing `Scan()`/`Value()`** for complex types stored in SQL

---

## Checklist for Schema Changes

Before submitting a PR, verify:

- [ ] Modified only schema YAML/JSON files (not auto-generated code)
- [ ] Created/updated `api.yml` as the index file if adding new construct
- [ ] Referenced all subschemas from `api.yml`
- [ ] Used non-deprecated `v1alpha1/core/api.yml` references
- [ ] Updated corresponding template files with default values
- [ ] Removed redundant `x-oapi-codegen-extra-tags` when using core refs
- [ ] Created helper files (`*_helper.go`) if Go structs need SQL compatibility
- [ ] Added `// This is not autogenerated.` comment to helper files
- [ ] Implemented `TableName()`, `Scan()`, `Value()` as needed in helper files
- [ ] Used `sync.Mutex` for thread-safe `Create()` methods
- [ ] Added `x-order` tags for consistent field ordering
- [ ] Ran `make build` successfully
- [ ] Ran `go test ./...` successfully
- [ ] Ran `npm run build` successfully
- [ ] Updated `typescript/index.ts` if adding new public types
- [ ] Verified only source schema files are in the commit



### 1. Mesheryctl Contributor Flow

  a. Add a new schema on a new command

  **Example:** You want to add a `mesheryctl model build` command. 
  **Steps:**
  - Add the new verb in `openapi.yaml` under the appropriate construct (e.g., `model/`)
  - Update `<construct>.json` if new properties are needed
  - Run:
  ```bash 
  make generate-types
  make golang-generate 
  ```
  - Implement the CLI logic
  - Add tests (Check existing unit tests for format)
  
  b. Add an existing schema on an existing command

  **Example**: 
  You detect a part of existing code that is not following the schema driven development principle (model is a struct created in mesheryctl command), you have two options:
  - If you know how to implement, update the existing code to use a proper struct generated from the `meshery/schemas` repository
    - Update the CLI logic
    - Add/Adjust tests if needed
  - If you don't know how to implement it, open an issue on Github using either a mesheryctl issue template ([feature](https://github.com/meshery/meshery/issues/new?template=mesheryctl_feature.md), [bug](https://github.com/meshery/meshery/issues/new?template=mesheryctl_bug.md))

  c. Add a new schema on an existing command

  **Steps:**
  - Add the new verb in `openapi.yaml` under the appropriate construct (e.g., `model/`)
  - Update `<construct>.json` if new properties are needed
  - Run:
  ```bash 
  make generate-types
  make golang-generate 
  ```
  - Update the CLI logic
  - Add/Adjust tests if needed

> *Why it matters:* This reduces drift between backend logic and API contract, enforces consistency between Meshery's components (Server, UI, CLI) and , resulting in higher quality code.

### 2. Meshery Server Contributor Flow
**Example:** Add a new `status` field to `component`.
**Steps:**
- Add the new property in `component.json`
- Run:
```bash
make validate-schemas
make golang-generate
```
- The generated Go structs (from `oapi-codegen`) are used in the backend.
- If the backend uses GORM with auto-migration enabled, these structs may be used to update the DB schema.
- Avoid manually editing the generated models, as they will be overwritten when schemas are regenerated.

> *Why it matters:* This reduces drift between backend logic and API contract, enforces consistency between Meshery's components (Server, UI, CLI) and , resulting in higher quality code.

### 3. Meshery UI Contributor Flow
**Example:** Show the new `version` field on the Model dashboard. 
**Steps:**
- Check `openapi.yaml` to verify the new field exists
- Wait for the backend to regenerate and expose the property
- Use RTK + TypeScript types to access and render data

> **Note**: `make generate-types` now generates only TypeScript types and schema-related objects. `_template.json` / `_template.yaml` files are no longer auto-generated.

*Why it matters:* UI stays in sync with the backend - fewer bugs, fewer mismatches, easier onboarding.

### 4. Meshery Docs Contributor Flow
**Example:** You are writing a guide!
**Steps:**
- Read the schema structure and workflows
- Walk through the scenarios above
- Write a guide that's accurate, actionable, and friendly

*Why it matters:* Docs are often the first impression contributors get. Schema-driven clarity starts here.

## Getting Help
- [GitHub Issues](https://github.com/meshery/schemas/issues) - Report bugs or request features
- [Community Slack](https://slack.meshery.io) - Real-time discussions with maintainers
- [Weekly Meetings](https://meshery.io/calendar) - Join our community calls

---
> **Community Resources**
> For more contribution guidelines, see the [Meshery Contributing Guide](https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md).
