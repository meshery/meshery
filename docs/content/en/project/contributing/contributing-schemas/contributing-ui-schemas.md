---
title: Schema-Driven UI Development
description: How to integrate Meshery schema definitions into the Meshery UI using react-jsonschema-form (RJSF), @sistent/sistent, and the @meshery/schemas TypeScript package.
categories: [contributing]
weight: 60
aliases:
  - /project/contributing/contributing-ui-schemas/
---

Meshery aims to decouple UI logic from hardcoded structures so the UI dynamically adapts when schemas change, without requiring constant manual updates. This page explains how to consume schema definitions from the [`meshery/schemas`](https://github.com/meshery/schemas) repository in the Meshery UI.

## Overview {#overview}

Meshery's schemas define the **structure of data (constructs)** using JSON Schema and **their API operations** using OpenAPI specs. From these definitions, Meshery auto-generates:

- Go structs for backend validation and API handling
- TypeScript types and objects for frontend development
- RTK Query hooks for type-safe data fetching
- Templates in JSON and YAML formats

The schema is the **single source of truth** — the UI should read types, validation rules, form metadata, and component props from the generated TypeScript package, not from hardcoded constants.

## Installing the schema package locally {#local-install}

When developing a UI feature that requires changes to a schema, install the local `schemas` checkout before pushing to npm:

```bash
# From the meshery/meshery/ui directory, with schemas checked out at ../schemas
npm install ../../schemas
```

This sets your `package.json` to:

```json
"@meshery/schemas": "file:../../schemas"
```

Build the schemas package first so the local install resolves correctly:

```bash
cd ../../schemas
npm run build
```

## Importing types and schema objects {#imports}

### Importing types {#importing-types}

Types are exported from the `v1beta1` or `v1beta2` namespaces:

```typescript
import { v1beta1 } from "@meshery/schemas";

// Use a named type
type Design = v1beta1.Design;
type Environment = v1beta1.Environment;

const renderDesignCard = (design: v1beta1.Design) => (
  <div>{design.display_name}</div>
);
```

{{% alert color="info" title="Property casing follows the schema" %}}
TypeScript property names mirror the schema exactly. DB-backed fields like `organization_id` and `display_name` are snake_case in both the schema and the generated TypeScript. Non-DB fields like `schemaVersion` are camelCase.
{{% /alert %}}

### Importing schema objects (runtime) {#importing-schemas}

Schema objects are named `{Construct}Definition{Version}OpenApiSchema`:

```typescript
import { EnvironmentDefinitionV1Beta1OpenApiSchema } from "@meshery/schemas";

const isValid = ajv.validate(
  EnvironmentDefinitionV1Beta1OpenApiSchema.components.schemas.Environment,
  data
);
```

{{% alert color="info" title="Can't find a schema import in the UI?" %}}
If a type or schema object is missing, it likely hasn't been exported from `typescript/index.ts` yet. Export it in the schemas repository and regenerate the package. Only actively used types are exported by default to keep the package lightweight.
{{% /alert %}}

## Integration points in the UI {#integration-points}

### A. RJSF JSON Schemas {#rjsf}

Meshery uses [react-jsonschema-form](https://github.com/rjsf-team/react-jsonschema-form) to render forms dynamically based on JSON schemas. All of Meshery's RJSF schemas are defined in the `@sistent/sistent` package, which extends schemas from `@meshery/schemas`.

This enables automatically-generated forms that adapt to schema changes without hardcoding field properties like type, enum, description, and constraints.

```typescript
import { EnvironmentDefinitionV1Beta1OpenApiSchema } from "@meshery/schemas";

const EnvironmentPayloadSchema =
  EnvironmentDefinitionV1Beta1OpenApiSchema.components.schemas.EnvironmentPayload;

const environmentFormSchema = {
  ...EnvironmentPayloadSchema,
  properties: {
    ...EnvironmentPayloadSchema.properties,
    name: {
      ...EnvironmentPayloadSchema.properties.name,
      ui: {
        label: EnvironmentPayloadSchema.properties.name.title,
        placeholder: EnvironmentPayloadSchema.properties.name.description,
      },
    },
  },
};
```

### B. Custom `x-rjsf-*` extensions {#x-rjsf}

Schemas may carry `x-rjsf-*` custom extensions to influence RJSF rendering:

| Extension | Effect |
|---|---|
| `x-rjsf-ui:widget` | Override the default input widget (`textarea`, `password`, `range`, etc.) |
| `x-rjsf-ui:order` | Control field display order |
| `x-rjsf-ui:options` | Pass arbitrary options to the widget |

These extensions are defined in the schema construct YAML and preserved through code generation into the schema object.

### C. General form UI from OpenAPI request bodies {#form-ui}

OpenAPI `{Construct}Payload` schemas serve as the foundation for form logic. These schemas include:

- Field validations (`required`, `format`, `maxLength`, `minimum`)
- Field types and formats (`string`, `integer`, `date-time`)
- Descriptions and examples for tooltips
- Enum constraints and conditional logic (`if`/`then`/`else`)
- Custom `x-rjsf-*` extensions for layout

Using the Payload schema (not the full entity schema) ensures form fields match exactly what the API accepts — the Payload schema contains only client-settable fields, without server-generated fields like `created_at`.

```typescript
import { v1beta1 } from "@meshery/schemas";

type EnvironmentFormValues = v1beta1.EnvironmentPayload;

const EnvironmentForm = ({ onSubmit }: { onSubmit: (v: EnvironmentFormValues) => void }) => (
  <form onSubmit={...}>
    <input type="text" name="name" />
    <input type="text" name="organization_id" />
  </form>
);
```

### D. UI-specific metadata {#ui-metadata}

Any UI metadata — names, hints, descriptions, defaults — should be derived from the schema, not defined separately in UI code. If a construct's schema includes `description`, `title`, `default`, or `examples` on a property, the UI should use them directly.

```typescript
import { EnvironmentDefinitionV1Beta1OpenApiSchema } from "@meshery/schemas";

const nameField = EnvironmentDefinitionV1Beta1OpenApiSchema
  .components.schemas.EnvironmentPayload.properties.name;

const tooltip = nameField.description;   // From schema — no duplication
const maxLen = nameField.maxLength;       // From schema — validated consistently
```

### E. Type safety for component props {#type-safety}

Generated TypeScript types keep UI components consistent with backend contracts:

```typescript
import { v1beta1 } from "@meshery/schemas";

const EnvironmentCard = ({ environment }: { environment: v1beta1.Environment }) => (
  <div>
    <h2>{environment.name}</h2>
    <p>{environment.description}</p>
    <small>Org: {environment.organization_id}</small>
  </div>
);
```

When the backend adds a new field to `Environment`, TypeScript immediately flags any UI component that uses the type but doesn't handle the new field — preventing silent runtime errors.

### F. RTK Query hooks {#rtk-hooks}

The schemas package ships RTK Query hooks that provide type-safe access to all API endpoints:

```typescript
import { mesheryApi } from "@meshery/schemas/mesheryApi";

// Type-safe data fetching
const { data, isLoading } = mesheryApi.useGetEnvironmentsQuery({
  page: 0,
  pagesize: 10,
  orgId: currentOrg.id,
});

// Type-safe mutation
const [createEnvironment] = mesheryApi.useCreateEnvironmentMutation();
await createEnvironment({ name: "Production", organization_id: orgId });
```

Do not hand-roll RTK query endpoints when `mesheryApi` provides a canonical equivalent. Hand-rolled endpoints diverge from the canonical schema types and break the consumer-audit CI gate.

## Next steps {#next-steps}

- [Schema Fundamentals](../schema-fundamentals/) — dual-schema pattern, casing rules, HTTP design
- [Authoring a Schema](../authoring-a-schema/) — step-by-step schema contribution guide
- [Schema Consumers and Build](../schema-consumers-and-build/) — Go helpers, TypeScript types, build pipeline
