---
layout: page
title: Contributing to Meshery Schema
permalink: project/contributing/contributing-schema
abstract: How to contribute to Meshery Schema
language: en
type: project
category: contributing
list: include
---

### Overview
Meshery follows schema-driven development. As a project, Meshery has different types of schemas. Some schemas are external facing, and some internal to Meshery itself. This repository serves as a central location for storing schemas from which all Meshery components can take reference.

The schemas follow a versioned approach to maintain backward compatibility while allowing for evolution of the definitions.

{% include alert.html type="info" title="Meshery Documentation Core Concepts" content="To better understand how schemas fit into Meshery's architecture, read about Meshery's core concepts in the <a href='https://docs.meshery.io/concepts/logical'>Meshery documentation</a>.`" %}

### Prerequisites
- **oapi-codegen**: This tool is essential for generating Go code from OpenAPI specifications. Install it using:

```bash
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
```

- **make**: The repository uses Makefiles to automate various tasks. Ensure you have make installed on your system.

### Development Workflow

#### Schema Resolution Process
When you work with the schemas, you'll frequently use this essential command:

```bash
make resolve-ref path="./schemas/constructs/v1beta1"
```

**Key functions:**
1.  Resolves `$ref` references between schema files
2.  Adds code generation metadata tags
3.  Creates complete, self-contained schemas
4.  Validates reference consistency

**Example:**
Consider this schema snippet with an external reference:

```json
"capabilities": {
  "type": "array",
  "description": "Meshery manages components...",
  "items": {
    "$ref": "../v1alpha1/capability.json" // reference here
  }
}
```

After running the command, it becomes a complete, self-contained schema:

```json
"capabilities": {
  "type": "array",
  "description": "Meshery manages components...",
  "items": {
    "$id": "https://schemas.meshery.io/capability.json",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "description": "Meshery manages entities...",
    "additionalProperties": false,
    "type": "object",
    "required": [
      "schemaVersion",
      "version",
      "displayName",
      "kind",
      "type",
      "entityState",
      "status"
    ],
    "x-oapi-codegen-extra-tags": { // additional metadata tag
      "gorm": "type:bytes;serializer:json"
    }
  }
}
```

**When to run this command?**

Whenever you:
-   Modify schema files
-   Add new schema references
-   Before generating Go code
-   When troubleshooting code generation issues

#### Code Generation and Configuration

The code generation process uses two key configuration files:
1. **scripts/config.yml**: Controls oapi-codegen behavior

```yml
package: organization  # Set your desired package name

generate:
  models: true  # Generate model structs

import-mapping:  # Map schema references to Go imports
  "../v1beta1/model.json": "github.com/meshery/schemas/models/v1beta1/model"
  "../v1alpha1/capability.json": "github.com/meshery/schemas/models/v1alpha1/capability"

output: models/v1beta1/organization.go  # Specify output file
output-options:
  skip-prune: true
  include-tags:  # Filter generated code by tags
  - organizations
```

2. **schemas/constructs/openapi/models.yml**: Defines OpenAPI schema references

```yml
openapi: 3.0.0
components:
  schemas:
    component_definition:
      $ref: ../v1beta1/component.json
```

#### Workflow

1.  Update schema references in **models.yml**:
    -   Uncomment/add needed schema references
    -   Each reference generates corresponding Go structs
2.  Modify **config.yml**:
    -   Set appropriate package name
    -   Update output file path
    -   Add required import mappings
    -   Configure include-tags if needed
3. Generate code:

```bash
oapi-codegen -config scripts/config.yml schemas/constructs/openapi/models.yml
```

**Key Points:**
-   Run `make resolve-ref` before code generation
-   Keep import mappings synchronized with schema references
-   Generated code inherits package name from config
-   Use tags to filter generated structs

### Example

Let's walk through a practical example, you made some changes in the **component.json**
1. First, ensure your schema references are resolved:

```yml
make resolve-ref path="./schemas/constructs/v1beta1"
```

2. Update `schemas/constructs/openapi/models.yml` to reference component.json:

```yml
openapi: 3.0.0
components:
  schemas:
    component_definition:
      $ref: ../v1beta1/component.json
```

3. Configure **config.yml**:

<code>
package: component

generate:
  models: true

import-mapping:
  "../v1beta1/model.json": "github.com/meshery/schemas/models/v1beta1/model"
  "../v1alpha1/capability.json": "github.com/meshery/schemas/models/v1alpha1/capability"

output: models/v1beta1/component/component.go
output-options:
  skip-prune: true
  include-tags:
  - components
</code>

4. Generate code

```yml
oapi-codegen -config scripts/config.yml schemas/constructs/openapi/models.yml
```

### Contributing to Documentation

1. **Schema Documentation**
-   Add detailed descriptions in schema fields
-   Include example values where helpful
-   Document validation rules and constraints

```json
{
  "displayName": {
    "type": "string",
    "description": "Human-readable name for the component.", // <-- description here
    "minLength": 1,
    "maxLength": 100,
    "examples": ["nginx-deployment"] // <-- examples
  }
}
```

### Getting Help

- [GitHub Issues](https://github.com/meshery/schemas/issues) - Report bugs or request features
- [Community Slack](https://slack.layer5.io) - Real-time discussions with maintainers
- [Weekly Meetings](https://layer5.io/community/calendar) - Join our community calls

---
> **Community Resources**
> For more contribution guidelines, see the [Meshery Contributing Guide](https://github.com/meshery/meshery/blob/master/CONTRIBUTING.md).
