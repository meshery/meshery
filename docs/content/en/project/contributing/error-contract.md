---
title: HTTP Error Response Contract
description: How Meshery Server responds to error conditions over HTTP, and how clients should parse those responses.
categories: [contributing]
---

# HTTP Error Response Contract

Every non-2xx HTTP response from Meshery Server carries a JSON body with
`Content-Type: application/json; charset=utf-8`. Clients should parse the body
as JSON before surfacing errors to users.

> **Migration status:** the JSON contract is enforced across the `./server`
> Go module by a `forbidigo` lint rule that blocks any new `http.Error` call
> at CI time, configured in
> [`.github/.golangci.yml`](https://github.com/meshery/meshery/blob/master/.github/.golangci.yml).
> The migration landed in PR
> [#18919](https://github.com/meshery/meshery/pull/18919) and the lint guard
> was flipped from advisory to fully blocking in PR
> [#18927](https://github.com/meshery/meshery/pull/18927). Legitimate
> non-JSON responders are handled in one of three ways: a file-level allowlist
> in `.github/.golangci.yml` for the Kubernetes healthz probe handler
> (`k8s_healthz_handler.go`) and the static-asset UI handler
> (`ui_handler.go`); SSE stream handlers (`events_streamer.go`,
> `load_test_handler.go`) that don't call `http.Error` at all (they emit
> error events through the active SSE writer instead); and binary/tar/YAML
> downloads and HTTP redirects, which use `w.Write` or `http.Redirect`
> rather than `http.Error` and so are not governed by the lint rule.

## Shape

```json
{
  "error": "Human-readable short description",
  "code": "meshery-server-1033",
  "severity": "ERROR",
  "probable_cause": ["Connection to the remote provider timed out."],
  "suggested_remediation": ["Verify that Meshery Cloud is reachable."],
  "long_description": ["Full technical details suitable for logs."]
}
```

### Fields

| Field | Required | Notes |
|-------|----------|-------|
| `error` | yes | User-facing message. For MeshKit errors, this is the ShortDescription. |
| `code` | when available | MeshKit error code (e.g. `meshery-server-1033`). Stable across releases. Use for telemetry, i18n lookup, and programmatic handling. |
| `severity` | when available | One of `EMERGENCY`, `ALERT`, `CRITICAL`, `FATAL`, `ERROR`. |
| `probable_cause` | optional | Array of strings. |
| `suggested_remediation` | optional | Array of strings. Surface to users when present. |
| `long_description` | optional | Array of strings. Suitable for developer logs; may contain stack-style detail. |

Fields marked "when available" are omitted (via `omitempty`) for errors that
originated outside the MeshKit error catalog.

## Client contract

- Do not rely on plain-text error bodies — they are always JSON.
- When `code` is present, prefer it over string matching on `error`.
- When `suggested_remediation` is non-empty, surface it alongside `error`.
- When the body is not valid JSON, treat the response as a bug and report
  the offending endpoint; do not attempt a text fallback.

## Producing errors in handlers

Use `writeMeshkitError(w, err, status)` in `server/handlers/utils.go`:

```go
if err != nil {
    h.log.Error(ErrGetResult(err))
    writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
    return
}
```

For bare-string errors without a MeshKit code, use `writeJSONError(w, msg, status)`.
Every bare-string error is a candidate for promotion to a MeshKit error —
prefer adding a code when fixing an adjacent bug.

Do not use `http.Error` in handlers or provider code. It writes
`Content-Type: text/plain` and strips MeshKit metadata, which crashes
RTK Query's default baseQuery on the UI.

### Streaming JSON responses

When a handler writes a JSON body via `json.NewEncoder(w).Encode(...)`,
encoding into the `ResponseWriter` directly commits headers + status before
the encoder errors. If the encode fails partway, you cannot emit a fresh
error response — the wire already has a partial JSON envelope and a 200 OK
status. **Buffer first, write once:**

```go
var buf bytes.Buffer
if err := json.NewEncoder(&buf).Encode(payload); err != nil {
    // No headers committed yet — safe to emit a fresh error response.
    writeMeshkitError(w, models.ErrEncoding(err, "<object name>"), http.StatusInternalServerError)
    return
}
w.Header().Set("Content-Type", "application/json; charset=utf-8")
if _, err := w.Write(buf.Bytes()); err != nil {
    // Client likely disconnected; headers already committed, can't emit
    // a new error. Log and move on.
    h.log.Warn(fmt.Sprintf("write response: %v", err))
}
```

The provider layer (`server/models/remote_provider.go`,
`server/models/default_local_provider.go`) follows this pattern — see commit
`ed1ce9f25c` for reference call sites.

Legitimate exceptions, organised by how each is enforced. The `forbidigo`
lint rule (in [`.github/.golangci.yml`](https://github.com/meshery/meshery/blob/master/.github/.golangci.yml))
matches `http.Error` codebase-wide; some non-JSON responders sidestep it by
not calling `http.Error` at all, while two files are exempted via a
file-level allowlist. See PR
[#18919](https://github.com/meshery/meshery/pull/18919) (migration) and
[#18927](https://github.com/meshery/meshery/pull/18927) (lint guard
blocking) for context.

- **File-level allowlist in `.github/.golangci.yml`** — entire file is exempt
  because plain text is the contract there, not an internal detail:
  - `k8s_healthz_handler.go` — Kubernetes healthz probe contract
  - `ui_handler.go` — serves static HTML; any `http.Error` call would relate
    to asset resolution, not API surface
- **Doesn't call `http.Error` at all** — these sites emit non-JSON for
  legitimate reasons but use other write mechanisms, so the lint rule does
  not apply:
  - SSE stream handlers (`events_streamer.go`, `load_test_handler.go`) emit
    error events through the active SSE writer (`Content-Type:
    text/event-stream`) rather than `http.Error`, which would corrupt the
    stream by trying to re-set headers after they have been flushed.
  - Binary/tar/YAML downloads use `w.Write(bytes)` with the appropriate
    `Content-Type`.
  - HTTP redirects use `http.Redirect`.
