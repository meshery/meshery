---
name: api-doc
description: Document REST API endpoints and GraphQL operations in the Meshery server.
tools: ['search/changes', 'search/codebase', 'edit/editFiles', 'vscode/extensions', 'web/fetch', 'web/githubRepo', 'vscode/getProjectSetupInfo', 'vscode/installExtension', 'vscode/newWorkspace', 'vscode/runCommand', 'vscode/openSimpleBrowser', 'read/problems', 'execute/getTerminalOutput', 'execute/runInTerminal', 'read/terminalLastCommand', 'read/terminalSelection', 'execute/createAndRunTask', 'execute', 'execute/runTask', 'execute/runTests', 'search', 'search/searchResults', 'execute/testFailure', 'search/usages', 'vscode/vscodeAPI', 'github/*', 'memory']
---

# Skill: api-doc

Document REST API endpoints and GraphQL operations in the Meshery server.

## Usage

Invoke this skill with a target handler file or directory:
- `/api-doc server/handlers/design_handler.go`
- `/api-doc server/handlers/` (document all handlers)

## Instructions

1. **Read the target handler file(s)** to identify all HTTP endpoint functions.
2. **Read the router configuration** in `server/router/` to find URL paths and HTTP methods mapped to each handler.
3. **Read the request/response model types** in `server/models/` referenced by the handlers.
4. **Generate documentation** in the format below.

### Output Format

For each endpoint, document:

````markdown
### `METHOD /api/path`

**Description**: Brief description of what the endpoint does.

**Authentication**: Required / Not required

**Parameters**:

| Name | In | Type | Required | Description |
|------|----|------|----------|-------------|
| id | path | string | yes | Resource identifier |
| page | query | int | no | Page number (default: 1) |

**Request Body** (if applicable):

```json
{
  "field": "type — description"
}
```

**Response** `200 OK`:

```json
{
  "field": "type — description"
}
```

**Error Responses**:

| Status | Description |
|--------|-------------|
| 400 | Invalid request body |
| 401 | Authentication required |
| 404 | Resource not found |
````

### Guidelines

- Derive descriptions from handler logic and comments — do not invent behavior
- Include all query parameters, path parameters, and headers the handler reads
- Document the actual JSON structure by reading the Go struct tags (`json:"field_name"`)
- Note any middleware applied (auth, rate limiting) visible in the router setup
- For GraphQL, document queries and mutations with their input/output types
- Group endpoints by resource type (designs, patterns, filters, connections, etc.)
