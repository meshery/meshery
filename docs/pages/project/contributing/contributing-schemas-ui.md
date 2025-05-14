---
layout: page
title: Contributing to Meshery Schemas
permalink: project/contributing/contributing-schemas-ui
abstract:
language: en
type: project
category: contributing
list: include
---

# Schema-Driven UI Development in Meshery

Meshery aims to decouple the UI logic from hardcoded structures and become fully schema-driven. This allows the UI to dynamically adapt based on changes in the underlying schema without requiring constant manual updates. This document explains how to integrate schema definitions from the [`meshery/schemas`](https://github.com/meshery/schemas) repository into the Meshery UI using a code-generation approach.

---

## Overview

Meshery’s schemas define the **structure of data (constructs)** using JSON Schema and **their behavior (API operations)** using OpenAPI specifications. From these definitions, Meshery auto-generates:

- **Go structs** for backend validation and API handling
- **TypeScript types and objects** for frontend development
- **Templates** in JSON and YAML formats

This approach ensures the **schemas remain the single source of truth**.

---

## Repository Structure

All schema definitions live in the [`meshery/schemas`](https://github.com/meshery/schemas) repository.

```
schemas/
  constructs/
    v1beta1/
      design/
        design.json              # JSON schema for the noun
        design_template.json     # JSON template (generated)
        design_template.yaml     # YAML template (generated)
        openapi.yml               # OpenAPI spec defining operations
        subschemas/               # Optional reusable schema parts
```

- **`<construct>.json`**: Defines structure (noun) — e.g., what a `Design` looks like.
- **`openapi.yml`**: Describes operations (verbs) — e.g., how to `GET`, `POST`, or `DELETE` a design.
- **`Templates`**: Valid, default-filled, resolved objects from the schema.

## Schema-Driven UI Development Workflow

Meshery follows a **schema-first** approach where the UI is driven by JSON schemas and OpenAPI specifications. Here's how you can contribute to and use these schemas in the Meshery UI.

---

### Step 1: Define or Update the Schema

If you're introducing or modifying a UI feature that requires a new schema:

1. Navigate to the appropriate schema directory:

   ```
   schemas/constructs/<version>/<construct>/
   ```

2. Define the schema structure in **`<construct>.json`** using **JSON Schema v7**.

3. Define API operations related to this schema in **`openapi.yml`** using **OpenAPI v3**. This helps generate consistent API types and models that align with your UI needs.

---

### Step 2: Generate TypeScript Types and Schema Objects

Use the following command to generate TypeScript types and JavaScript schema objects:

```bash
make generate-types
```

This will:

- Generate `.ts` and `.d.ts` files under:

  ```
  typescript/constructs/<version>/<construct>/
  ```

- Create:

  - **Typed interfaces** for schema validation and component props
  - **Schema objects** usable for dynamic UI (e.g., forms, validations)

---

### Step 3: Build and Export Types/Schema Objects

After generation:

1. Open [`typescript/index.ts`](https://github.com/meshery/schemas/blob/master/typescript/index.ts)

   - Export the newly generated types and schema objects from here.
   - Example:

     ```ts
     export * from "./constructs/v1beta1/design";
     ```

2. Build the TypeScript package to make the changes usable:

   ```bash
   npm run build
   ```

---

### Step 4: Use Schema Package in UI

To consume the schema in the UI:

1. Install the schema package locally:

   ```bash
   npm install <path-to-schemas-repo>
   ```

   Examples:

   - Relative path:

     ```bash
     npm install ../../schemas
     ```

   - Absolute path:

     ```bash
     npm install /home/user/code/schemas
     ```

2. This will update your `package.json` to something like:

   ```json
   "@layer5/schema": "file:../../schemas"
   ```

3. Now you can import types and schema objects in your UI components:

#### Example: Importing and Using a Type

```ts
import { DesignTypes } from "@layer5/schema";

// Type safety!
const renderDesignCard = (design: Design) => <div>{design.name}</div>;
```

#### Example: Accessing Runtime Schema

```ts
import { DesignSchema } from "@layer5/schema";

const validateDesign = (data) => {
  const isValid = ajv.validate(DesignSchema, data);
  return isValid;
};
```

---

## Integration Points in UI

### A. RJSF JSON Schemas

Meshery uses [react-jsonschema-form](https://github.com/rjsf-team/react-jsonschema-form) to render forms dynamically based on JSON schemas. All of Meshery’s RJSF schemas are defined in the `@layer5/sistent` package, which extends schemas from the `@layer5/schema` package.

This approach enables us to generate forms that automatically adapt to the schema structure without hardcoding field properties like type, enum, description, and others.

```ts
import { DesignSchema } from "@layer5/schema";

const designSchema = {
  ...DesignSchema,
  properties: {
    ...DesignSchema.properties,
    name: {
      ...DesignSchema.properties.name,
      description: DesignSchema.properties.description,
      ui: {
        label: DesignSchema.properties.name.title,
        placeholder: DesignSchema.properties.name.description,
      },
    },
    // Other properties with UI-specific enhancements
  },
};
```

---

### B. General Form UI

OpenAPI schemas (especially request bodies for POST/PUT operations) serve as the foundation for building form logic. These definitions include:

- Field validations (e.g. `required`, `format`, `maxLength`)
- Field types and formats (e.g. `string`, `integer`, `date-time`)
- Descriptions and examples
- Enum constraints and conditional logic
- Custom extensions like `x-rjsf-*` for layout

This ensures alignment between frontend form behavior and backend expectations.

```ts
import { DesignTypes } from "@layer5/schema";

const DesignForm = ({ design }: { design: DesignTypes }) => (
  <form>
    <input type="text" value={design.name} />
    {/* More fields derived from the schema */}
  </form>
);
```

---

### C. UI-Specific Descriptions and Enhancements

Any UI-specific metadata—such as `name`, `type`, `hints`, `descriptions`, `defaults`, etc.—is defined directly within the relevant schema object. Elements like tooltips, descriptions, and other metadata are frequently needed across the UI, so having a single source of truth in the schema object ensures consistency and reduces duplication.

For example, if we have a `Design` schema, the UI retrieves details like the design’s name, description, and other properties directly from the schema object.

---

### D. Type Safety for Component Props

Generated TypeScript types from the schema ensure UI components are type-safe and consistent with backend contracts.

```ts
import { DesignTypes } from "@layer5/schema";

const DesignCard = ({ design }: { design: DesignTypes }) => (
  <div>
    <h2>{design.name}</h2>
    <p>{design.description}</p>
  </div>
);
```
