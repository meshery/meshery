---
title: Schema-Driven UI Development in Meshery
description: How to contribute to Meshery Schemas for UI
categories: [contributing]
aliases: [/project/contributing/contributing-ui-schemas]
---

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
   "@meshery/schema": "file:../../schemas"
   ```

3. Now you can import types and schema objects in your UI components:

#### Example: Importing and Using a Type

```ts
import { DesignTypes } from "@meshery/schema";

// Type safety!
const renderDesignCard = (design: Design) => <div>{design.name}</div>;
```

#### Example: Accessing Runtime Schema

```ts
import { DesignSchema } from "@meshery/schema";

const validateDesign = (data) => {
  const isValid = ajv.validate(DesignSchema, data);
  return isValid;
};
```

{{% alert color="info" title="Can't find a Schema Import in the UI?" %}}
<p>If you're trying to import a schema object or type in the UI but it's missing, it's likely because it hasn't been exported yet from the schemas package.</p>
<p>To keep the package lightweight, only actively used types and objects are exported by default. If you need access to a new schema, simply export it in the <a href='https://github.com/meshery/schemas' target='_blank'>Meshery Schemas repository</a> and regenerate the package.</p>
{{% /alert %}}

---

## Integration Points in UI

### A. Generated RTK Query Client

`@meshery/schemas` publishes more than types: `@meshery/schemas/mesheryApi` is a
generated [RTK Query](https://redux-toolkit.js.org/rtk-query/overview) client with a
hook per API operation. `ui/rtk-query/index.ts` re-exports that client as `api`, so
every module under `ui/rtk-query/` injects into the *same* instance and the generated
hooks are already available from those modules.

**Consume the generated hook. Do not re-declare the request.** A hand-written
`builder.query`/`builder.mutation` for an operation schemas already defines forks the
wire contract silently, and Meshery and Layer5 Cloud drift apart with nothing failing.

```ts
// Correct - the generated hook, imported directly.
import { useGetWorkspacesQuery } from "@meshery/schemas/mesheryApi";
```

A local module is the right place only for *ergonomics* over a generated endpoint -
adapting an argument shape, or attaching a Meshery-local cache tag:

```ts
// Argument adaptation: callers pass a bare id, the generated endpoint takes an object.
export const useGetCredentialByIdQuery = (credentialId: string, options?: object) =>
  useSchemasGetCredentialByIdQuery({ credentialId }, options);
```

```ts
// Cache tags: the callback form of enhanceEndpoints APPENDS to the generated tags.
// The object form Object.assigns over them and drops every schemas-side tag.
import { appendInvalidatesTags } from "./utils";

api.enhanceEndpoints({
  addTagTypes: [TAGS.DESIGNS],
  endpoints: {
    importDesign: appendInvalidatesTags("importDesign", { type: TAGS.DESIGNS }),
  },
});
```

#### A local endpoint that shadows a generated one is dead code

`injectEndpoints` **silently discards** an endpoint whose name the generated client
already defines - a dev-only console warning, nothing more - and serves every call
from the generated definition. Neither the typecheck nor a test that only asserts a
hook exists will notice, so a local declaration can look authoritative while a
different request goes over the wire. Meshery has shipped user-visible breakage this
way: callers shaped for the discarded local definition sent `undefined` where the
generated endpoint expected an id.

Before adding an endpoint, check the name against the generated client:

```bash
grep '<operationId>:t\.' ui/node_modules/@meshery/schemas/dist/mesheryApi.js
```

Override deliberately only when the generated operation is genuinely wrong for
Meshery today - for example when schemas has landed a path the server has not
adopted yet. Set `overrideExisting: true` explicitly, explain why in a comment, and
link the `meshery/schemas` issue tracking the divergence - see
`ui/rtk-query/notificationCenter.ts` and
[meshery/schemas#1134](https://github.com/meshery/schemas/issues/1134).

Test the **effective** endpoint, not the declared one: dispatch through a real store
and assert the URL, method and body. A test that reads the module's source shape, or
that calls `fetch` itself and then asserts its own mock, proves nothing. The guard
tests to copy are
`ui/rtk-query/__tests__/{workspace-mutation-wrappers,notificationCenter-effective-endpoints}.test.ts*`;
where a local hook only adapts an argument shape, drive that hook with `renderHook`
so the adaptation itself is exercised rather than bypassed.

---

### B. RJSF JSON Schemas

Meshery uses [react-jsonschema-form](https://github.com/rjsf-team/react-jsonschema-form) to render forms dynamically based on JSON schemas. All of Meshery’s RJSF schemas are defined in the `@sistent/sistent` package, which extends schemas from the `@meshery/schema` package.

This approach enables us to generate forms that automatically adapt to the schema structure without hardcoding field properties like type, enum, description, and others.

```ts
import { DesignSchema } from "@meshery/schema";

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

### C. General Form UI

OpenAPI schemas (especially request bodies for POST/PUT operations) serve as the foundation for building form logic. These definitions include:

- Field validations (e.g. `required`, `format`, `maxLength`)
- Field types and formats (e.g. `string`, `integer`, `date-time`)
- Descriptions and examples
- Enum constraints and conditional logic
- Custom extensions like `x-rjsf-*` for layout

This ensures alignment between frontend form behavior and backend expectations.

```ts
import { DesignTypes } from "@meshery/schema";

const DesignForm = ({ design }: { design: DesignTypes }) => (
  <form>
    <input type="text" value={design.name} />
    {/* More fields derived from the schema */}
  </form>
);
```

---

### D. UI-Specific Descriptions and Enhancements

Any UI-specific metadata—such as `name`, `type`, `hints`, `descriptions`, `defaults`, etc.—is defined directly within the relevant schema object. Elements like tooltips, descriptions, and other metadata are frequently needed across the UI, so having a single source of truth in the schema object ensures consistency and reduces duplication.

For example, if we have a `Design` schema, the UI retrieves details like the design’s name, description, and other properties directly from the schema object.

---

### E. Type Safety for Component Props

Generated TypeScript types from the schema ensure UI components are type-safe and consistent with backend contracts.

```ts
import { DesignTypes } from "@meshery/schema";

const DesignCard = ({ design }: { design: DesignTypes }) => (
  <div>
    <h2>{design.name}</h2>
    <p>{design.description}</p>
  </div>
);
```
