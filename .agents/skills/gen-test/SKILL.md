---
name: gen-test
description: Generate idiomatic tests for Go packages and handlers in the Meshery project.
tools: ['search/changes', 'search/codebase', 'edit/editFiles', 'vscode/extensions', 'web/fetch', 'web/githubRepo', 'vscode/getProjectSetupInfo', 'vscode/installExtension', 'vscode/newWorkspace', 'vscode/runCommand', 'vscode/openSimpleBrowser', 'read/problems', 'execute/getTerminalOutput', 'execute/runInTerminal', 'read/terminalLastCommand', 'read/terminalSelection', 'execute/createAndRunTask', 'execute', 'execute/runTask', 'execute/runTests', 'search', 'search/searchResults', 'execute/testFailure', 'search/usages', 'vscode/vscodeAPI', 'github/*', 'memory']
---

# Skill: gen-test

Generate idiomatic tests for Go packages and handlers in the Meshery project.

## Usage

Invoke this skill with a target file or package path:
- `/gen-test server/handlers/provider_handler.go`
- `/gen-test mesheryctl/internal/cli/root/system/start.go`

## Instructions

1. **Read the target file** to understand its exported functions, methods, and types.
2. **Read existing tests** in the same package (look for `*_test.go` files) to match the project's testing style.
3. **Generate tests** following these conventions:

### Go Test Conventions

- Use table-driven tests with `t.Run` subtests
- Name test functions `Test<FunctionName>` or `Test<Type>_<Method>`
- Use `testify/assert` or `testify/require` if already used in the package; otherwise use standard library
- For HTTP handlers, use `httptest.NewRecorder()` and `httptest.NewRequest()`
- Mock external dependencies using interfaces — check if mock implementations already exist in the package
- Test both success and error paths
- Include edge cases: nil inputs, empty strings, boundary values

### File Placement

- Place test files in the same directory as the source file
- Name: `<source_file>_test.go` (or append to existing test file if one exists)

### What to Test

For **handlers** (`server/handlers/`):
- HTTP status codes for valid and invalid requests
- Response body structure
- Authentication/authorization checks (if applicable)
- Input validation errors

For **models** (`server/models/`):
- Struct method behavior
- Serialization/deserialization
- Validation logic

For **mesheryctl commands** (`mesheryctl/`):
- Command output for various flag combinations
- Error messages for invalid input
- Exit codes

### Example Pattern

```go
func TestHandlerName(t *testing.T) {
    tests := []struct {
        name           string
        method         string
        path           string
        body           string
        expectedStatus int
    }{
        {
            name:           "valid request",
            method:         http.MethodGet,
            path:           "/api/resource",
            expectedStatus: http.StatusOK,
        },
        {
            name:           "missing required field",
            method:         http.MethodPost,
            path:           "/api/resource",
            body:           `{}`,
            expectedStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
            rec := httptest.NewRecorder()

            handler.ServeHTTP(rec, req)

            if rec.Code != tt.expectedStatus {
                t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
            }
        })
    }
}
```
