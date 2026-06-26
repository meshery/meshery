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
10. [The Dual-Schema Pattern]({{< ref "dual-schema-pattern.md" >}})
11. [Build and Generation]({{< ref "build-and-generation.md" >}})
12. [Schema-Driven UI Development]({{< ref "ui-schemas.md" >}})
13. [Getting Help](#getting-help)

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

### DB-Mirrored Fields

DB-mirrored fields such as `created_at`, `updated_at`, and `user_id` intentionally remain **snake_case** to mirror existing database columns. Do not rename these to camelCase.

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
  "createdAt": "0001-01-01T00:00:00Z",
  "updatedAt": "0001-01-01T00:00:00Z"
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

## Getting Help

- [GitHub Issues](https://github.com/meshery/schemas/issues) - Report bugs or request features
- [Community Slack](https://slack.meshery.io) - Real-time discussions with maintainers
- [Weekly Meetings](https://meshery.io/calendar) - Join our community calls

## Further Reading

- [meshery/schemas README](https://github.com/meshery/schemas/blob/master/README.md) - Full reference for schema authoring
- [AGENTS.md](https://github.com/meshery/schemas/blob/master/AGENTS.md) - Contributor checklist
- [Core schema definitions](https://github.com/meshery/schemas/blob/master/schemas/constructs/v1alpha1/core/api.yml) - Reusable building blocks
- [Academy construct](https://github.com/meshery/schemas/blob/master/schemas/constructs/v1beta1/academy/api.yml) - Exemplar for advanced patterns including `x-generate-db-helpers`

---
> **Community Resources**
> For more contribution guidelines, see the [Meshery Contributing Guide](https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md).
