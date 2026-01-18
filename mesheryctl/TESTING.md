# mesheryctl testing

This document describes how to run `mesheryctl` tests locally.

## Prerequisites

- Go (use the repository's Go version; see the root `go.mod`)
- From the repo root: `cd /path/to/meshery`

## Go unit tests

Run all Go tests under `mesheryctl`:

- From repo root: `go test ./mesheryctl/...`
- From `mesheryctl/`: `cd mesheryctl && go test ./...`

Run a single package:

- `go test ./mesheryctl/internal/cli/root/components -count=1`

Run a single test (or set of tests) by name:

- `go test ./mesheryctl/internal/cli/root/components -run TestListComponent -count=1`

Notes:

- Many CLI command tests use HTTP mocking and the shared harness in `mesheryctl/pkg/utils/testing.go`.
- Some output is normalized (ANSI codes stripped, pagination cleaned) before assertions.

## Coverage

From `mesheryctl/`:

- `make coverage`

This writes `coverage.out` and generates `coverage.html`.

## Integration tests

Meshery has multiple kinds of integration tests:

- Go integration tests (tagged or named):
  - From `mesheryctl/`: `go test -run Integration ./...`
- BATS-based CLI tests:
  - Located under `mesheryctl/tests/`
  - See the `mesheryctl/Makefile` variables under the BATS section (e.g. `BATS_FOLDER_PATTERN`, `BATS_FILE_PATTERN`).

## Common troubleshooting

- **Auth/token errors in CLI tests**: the test harness writes a temporary JSON token file (cookies) and sets `utils.TokenFlag`.
- **Hanging tests**: check for commands writing to stdout/stderr without restoring pipes; prefer the shared harness helpers.
