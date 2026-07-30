---
name: meshkit-errors
enabled: true
event: file
action: block
conditions:
  - field: file_path
    operator: regex_match
    pattern: /server/.*\.go$
  - field: file_path
    operator: not_contains
    pattern: _test.go
  - field: content
    operator: regex_match
    pattern: fmt\.Errorf\(|errors\.New\(\s*"|errors\.Errorf\(|errors\.Wrapf?\(
---

MeshKit error framework required in the Go backend.

This `server/**/*.go` edit adds an ad-hoc error (`fmt.Errorf`, std-lib `errors.New("…")`, or `pkg/errors`) that bypasses MeshKit. Every backend error MUST be a structured MeshKit error with a unique code (CLAUDE.md Critical Rule 1).

Replace it with a builder over `github.com/meshery/meshkit/errors`:

```go
const ErrFooCode = "meshery-cloud-NNNN"
func ErrFoo(err error) error {
    return errors.New(
        ErrFooCode, errors.Alert,
        []string{"short description"},
        []string{err.Error()},
        []string{"probable cause(s)"},
        []string{"remedy/remedies"},
    )
}
```

Allocate the next free code in the package's `error.go` and verify with `make error`. MeshKit's own `errors.New(ErrCode, …)` is allowed — its first argument is a code constant, not a string literal.

Companion to `.claude/hooks/meshkit-errors.sh` (the wired, net-new-aware enforcer). Fires when this directory is the active working directory.
