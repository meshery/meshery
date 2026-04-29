---
title: Contributing to Meshery Schemas
description: How to contribute to Meshery's schema-driven development ecosystem — authoring OpenAPI schemas, following casing and design rules, and working with the build pipeline.
categories: [contributing]
---

Meshery follows **Schema-Driven Development (SDD)**: the structure, validation rules, and API contracts for every resource are defined centrally in the [`meshery/schemas`](https://github.com/meshery/schemas) repository and propagate automatically into Go structs, TypeScript types, and RTK Query clients.

## Quick start

```bash
git clone https://github.com/meshery/schemas.git
cd schemas
make setup          # go mod tidy + npm install
make build          # generate Go, TypeScript, RTK, validate
npm run build       # build TypeScript distribution
```

## Where to go next

| I want to… | Start here |
|---|---|
| Understand what schemas are and how they propagate | [Schema Fundamentals](./schema-fundamentals/) |
| Add a new schema construct or modify an existing one | [Authoring a Schema](./authoring-a-schema/) |
| Understand the validation rules or vendor extensions | [Schema Rules and Extensions](./schema-rules-and-extensions/) |
| Work with generated Go/TypeScript code or run the build | [Schema Consumers and Build](./schema-consumers-and-build/) |
| Contribute UI forms driven by schemas (RJSF / @sistent) | [Schema-Driven UI Development](./contributing-ui-schemas/) |

## Prerequisites

- Go — see [`go.mod`](https://github.com/meshery/schemas/blob/master/go.mod) for the required version
- Node.js and npm — see [`package.json`](https://github.com/meshery/schemas/blob/master/package.json) `engines` field
- `oapi-codegen` — install via `go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest`
- `make`

## Getting help

- [GitHub Issues](https://github.com/meshery/schemas/issues)
- [Community Slack](https://slack.meshery.io)
- [Weekly Meetings](https://meshery.io/calendar)

## Further reading

- [meshery/schemas README](https://github.com/meshery/schemas/blob/master/README.md)
- [AGENTS.md — authoritative contributor contract](https://github.com/meshery/schemas/blob/master/AGENTS.md)
- [Core schema definitions](https://github.com/meshery/schemas/blob/master/schemas/constructs/)
