\---

title: "AI Adapter: Status, Error, and Event Semantics"

description: "Design specification defining error categories, HTTP/GraphQL response codes, event severity, log levels, redaction rules, retry eligibility, and operation ID correlation for AI Adapter operations in Meshery Server."

\---



\## Overview



This document defines the contract between Meshery Server and the AI Adapter

for all error, status, and event scenarios. It is the authoritative reference

for implementers of the AI Adapter handler, the server-side AI operation

dispatcher, and any mesheryctl or UI code consuming AI operation responses.



The AI Adapter introduces a new class of errors distinct from existing

infrastructure adapter errors: they are non-deterministic (LLM outputs vary),

involve external provider credentials, and carry a risk of leaking sensitive

data (API keys, system prompts, user infrastructure descriptions) in logs or

event streams.



\### Scope



\- Server-side handler for AI Adapter operations (`server/handlers/`)

\- Event stream (`EventsResponse`) fields for AI operations

\- HTTP response codes for AI Adapter REST endpoints

\- GraphQL error payloads (where applicable)

\- Log levels and redaction rules

\- mesheryctl exit-code and display behavior

\- Retry and fallback eligibility per error category



\### References



\- \[meshery/meshery#19092](https://github.com/meshery/meshery/issues/19092) — parent feature issue

\- \[meshery/meshery#19220](https://github.com/meshery/meshery/issues/19220) — Connection lifecycle for AI providers

\- \[meshery/meshery#19221](https://github.com/meshery/meshery/issues/19221) — Credential lifecycle for AI providers

\- MeshKit errors contract: `github.com/layer5io/meshkit/errors`

\- `server/handlers/error.go` — existing error code registry

\- `meshops.proto` — `EventsResponse` field definitions



\---



\## Error Code Allocation



AI Adapter error codes are allocated in the range \*\*meshery-server-1200 –

meshery-server-1299\*\* to avoid collision with existing handler error codes

(currently reaching \~1199 in `server/handlers/error.go`).



| Code | Symbolic Constant | Short Description |

|------|------------------|-------------------|

| meshery-server-1200 | `ErrAIProviderUnreachableCode` | Provider endpoint unreachable |

| meshery-server-1201 | `ErrAIAuthFailureCode` | Authentication / API key rejected |

| meshery-server-1202 | `ErrAIRateLimitCode` | Provider rate limit exceeded |

| meshery-server-1203 | `ErrAIContextWindowCode` | Prompt exceeds provider context window |

| meshery-server-1204 | `ErrAIInvalidResponseCode` | Provider returned unparseable response |

| meshery-server-1205 | `ErrAIDesignValidationCode` | Generated design fails schema validation |

| meshery-server-1206 | `ErrAIPromptEmptyCode` | User prompt is empty or whitespace-only |

| meshery-server-1207 | `ErrAIProviderNotConfiguredCode` | No AI Connection configured for this user |

| meshery-server-1208 | `ErrAIOperationTimeoutCode` | Provider response exceeded deadline |

| meshery-server-1209 | `ErrAIStreamAbortedCode` | Client disconnected during streaming |

| meshery-server-1210 | `ErrAISchemaInjectionCode` | Failed to inject Meshery Model schema into prompt |

| meshery-server-1211 | `ErrAIDesignImportCode` | Generated design could not be imported into Meshery |

| meshery-server-1212 | `ErrAIOperationIDMissingCode` | Operation ID absent from request |



\---



\## Error Category Table



Each row defines the full behavioral contract for one error category.



\### 1. Provider Unreachable (`ErrAIProviderUnreachable`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `503 Service Unavailable` |

| \*\*GraphQL error code\*\* | `SERVICE\_UNAVAILABLE` |

| \*\*EventsResponse.Summary\*\* | `"AI provider unreachable"` |

| \*\*EventsResponse.Details\*\* | `"Could not connect to <provider>. Verify the provider URL and network connectivity."` (provider name only — no URL that may contain credentials) |

| \*\*Log level\*\* | `WARN` |

| \*\*Log message\*\* | `"AI provider unreachable: <redacted-url>"` (URL logged only at DEBUG with credentials stripped) |

| \*\*mesheryctl exit code\*\* | `1` |

| \*\*mesheryctl display\*\* | `Error: AI provider unreachable. Check your Connection configuration with 'mesheryctl exp ai config'.` |

| \*\*Retry eligible\*\* | Yes — exponential backoff, max 3 attempts, 2 s initial delay |

| \*\*Redaction required\*\* | Provider URL must have credentials stripped before logging at WARN or above |



\---



\### 2. Authentication Failure (`ErrAIAuthFailure`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `401 Unauthorized` |

| \*\*GraphQL error code\*\* | `UNAUTHENTICATED` |

| \*\*EventsResponse.Summary\*\* | `"AI provider authentication failed"` |

| \*\*EventsResponse.Details\*\* | `"The API key or token for <provider> was rejected. Verify your Credential in Meshery."` |

| \*\*Log level\*\* | `ERROR` |

| \*\*Log message\*\* | `"AI auth failure for provider=<name> connection\_id=<id>"` — API key value MUST NOT appear at any log level |

| \*\*mesheryctl exit code\*\* | `4` |

| \*\*mesheryctl display\*\* | `Error: Authentication failed. Update your AI provider Credential with 'mesheryctl exp ai config'.` |

| \*\*Retry eligible\*\* | No — retrying with the same key will not succeed |

| \*\*Redaction required\*\* | API key / token MUST be redacted in all log levels, event streams, and GraphQL error extensions |



\---



\### 3. Rate Limit Exceeded (`ErrAIRateLimit`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `429 Too Many Requests` |

| \*\*GraphQL error code\*\* | `RESOURCE\_EXHAUSTED` |

| \*\*EventsResponse.Summary\*\* | `"AI provider rate limit exceeded"` |

| \*\*EventsResponse.Details\*\* | `"Request quota exceeded for <provider>. Retry after <retry-after> seconds."` (retry-after sourced from provider response header if available) |

| \*\*Log level\*\* | `WARN` |

| \*\*Log message\*\* | `"AI rate limit: provider=<name> retry\_after=<seconds>"` |

| \*\*mesheryctl exit code\*\* | `1` |

| \*\*mesheryctl display\*\* | `Error: Rate limit exceeded. Retry after <N> seconds.` |

| \*\*Retry eligible\*\* | Yes — honor `Retry-After` header; default 60 s if absent |

| \*\*Redaction required\*\* | None beyond standard credential redaction |



\---



\### 4. Context Window Exceeded (`ErrAIContextWindow`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `422 Unprocessable Entity` |

| \*\*GraphQL error code\*\* | `BAD\_USER\_INPUT` |

| \*\*EventsResponse.Summary\*\* | `"Prompt exceeds model context window"` |

| \*\*EventsResponse.Details\*\* | `"The combined prompt and schema context exceeds the model's token limit. Simplify your request or select a model with a larger context window."` |

| \*\*Log level\*\* | `INFO` |

| \*\*Log message\*\* | `"AI context window exceeded: provider=<name> estimated\_tokens=<n>"` — prompt text MUST NOT be logged at INFO or above |

| \*\*mesheryctl exit code\*\* | `2` |

| \*\*mesheryctl display\*\* | `Error: Prompt too long for selected model. Simplify your request.` |

| \*\*Retry eligible\*\* | No — same prompt will fail again |

| \*\*Redaction required\*\* | Prompt content MUST NOT appear in logs at INFO or above |



\---



\### 5. Invalid / Unparseable Provider Response (`ErrAIInvalidResponse`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `502 Bad Gateway` |

| \*\*GraphQL error code\*\* | `INTERNAL` |

| \*\*EventsResponse.Summary\*\* | `"AI provider returned an invalid response"` |

| \*\*EventsResponse.Details\*\* | `"The response from <provider> could not be parsed. The provider may be experiencing issues."` |

| \*\*Log level\*\* | `ERROR` |

| \*\*Log message\*\* | `"AI invalid response: provider=<name> status=<http-status>"` — response body logged only at DEBUG, truncated to 512 chars |

| \*\*mesheryctl exit code\*\* | `1` |

| \*\*mesheryctl display\*\* | `Error: Unexpected response from AI provider. Try again or switch providers.` |

| \*\*Retry eligible\*\* | Yes — once, immediately |

| \*\*Redaction required\*\* | Response body may contain injected schema with sensitive model names; log only at DEBUG, truncated |



\---



\### 6. Generated Design Fails Schema Validation (`ErrAIDesignValidation`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `422 Unprocessable Entity` |

| \*\*GraphQL error code\*\* | `BAD\_USER\_INPUT` |

| \*\*EventsResponse.Summary\*\* | `"Generated design failed schema validation"` |

| \*\*EventsResponse.Details\*\* | `"The design produced by the AI does not conform to Meshery's schema. Validation errors: <list>."` |

| \*\*Log level\*\* | `WARN` |

| \*\*Log message\*\* | `"AI design validation failed: operation\_id=<id> errors=<count>"` |

| \*\*mesheryctl exit code\*\* | `2` |

| \*\*mesheryctl display\*\* | `Warning: Generated design has validation errors. Review before applying.` |

| \*\*Retry eligible\*\* | Yes — with prompt refinement (automated retry not recommended) |

| \*\*Redaction required\*\* | None |



\---



\### 7. Empty Prompt (`ErrAIPromptEmpty`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `400 Bad Request` |

| \*\*GraphQL error code\*\* | `BAD\_USER\_INPUT` |

| \*\*EventsResponse.Summary\*\* | `"Prompt is empty"` |

| \*\*EventsResponse.Details\*\* | `"A non-empty natural language prompt is required to generate a design."` |

| \*\*Log level\*\* | `DEBUG` |

| \*\*Log message\*\* | `"AI prompt empty: operation\_id=<id>"` |

| \*\*mesheryctl exit code\*\* | `2` |

| \*\*mesheryctl display\*\* | `Error: Prompt cannot be empty. Provide a description of the infrastructure to generate.` |

| \*\*Retry eligible\*\* | No |

| \*\*Redaction required\*\* | None |



\---



\### 8. Provider Not Configured (`ErrAIProviderNotConfigured`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `424 Failed Dependency` |

| \*\*GraphQL error code\*\* | `FAILED\_PRECONDITION` |

| \*\*EventsResponse.Summary\*\* | `"No AI provider configured"` |

| \*\*EventsResponse.Details\*\* | `"No AI Connection is active for this user. Configure a provider with 'mesheryctl exp ai config'."` |

| \*\*Log level\*\* | `INFO` |

| \*\*Log message\*\* | `"AI provider not configured: user\_id=<id>"` |

| \*\*mesheryctl exit code\*\* | `3` |

| \*\*mesheryctl display\*\* | `Error: No AI provider configured. Run 'mesheryctl exp ai config' to add one.` |

| \*\*Retry eligible\*\* | No |

| \*\*Redaction required\*\* | None |



\---



\### 9. Operation Timeout (`ErrAIOperationTimeout`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `504 Gateway Timeout` |

| \*\*GraphQL error code\*\* | `DEADLINE\_EXCEEDED` |

| \*\*EventsResponse.Summary\*\* | `"AI operation timed out"` |

| \*\*EventsResponse.Details\*\* | `"The AI provider did not respond within the configured deadline. Try again or increase the timeout."` |

| \*\*Log level\*\* | `WARN` |

| \*\*Log message\*\* | `"AI timeout: provider=<name> operation\_id=<id> deadline=<duration>"` |

| \*\*mesheryctl exit code\*\* | `1` |

| \*\*mesheryctl display\*\* | `Error: AI operation timed out. Try again or switch to a faster model.` |

| \*\*Retry eligible\*\* | Yes — once, after 5 s |

| \*\*Redaction required\*\* | None |



\---



\### 10. Client Disconnected During Streaming (`ErrAIStreamAborted`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | N/A — connection already closed |

| \*\*GraphQL error code\*\* | N/A |

| \*\*EventsResponse.Summary\*\* | Not emitted — client has disconnected |

| \*\*EventsResponse.Details\*\* | Not emitted |

| \*\*Log level\*\* | `DEBUG` |

| \*\*Log message\*\* | `"AI stream aborted: client disconnected operation\_id=<id>"` |

| \*\*mesheryctl exit code\*\* | N/A — client-initiated |

| \*\*mesheryctl display\*\* | N/A |

| \*\*Retry eligible\*\* | No — client must re-initiate |

| \*\*Redaction required\*\* | None |



\---



\### 11. Schema Injection Failure (`ErrAISchemaInjection`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `500 Internal Server Error` |

| \*\*GraphQL error code\*\* | `INTERNAL` |

| \*\*EventsResponse.Summary\*\* | `"Failed to build AI prompt context"` |

| \*\*EventsResponse.Details\*\* | `"Meshery could not retrieve the Model schema definitions needed to construct the AI prompt. The registry may be unavailable."` |

| \*\*Log level\*\* | `ERROR` |

| \*\*Log message\*\* | `"AI schema injection failed: operation\_id=<id> err=<message>"` |

| \*\*mesheryctl exit code\*\* | `1` |

| \*\*mesheryctl display\*\* | `Error: Internal error building AI prompt. Check server logs.` |

| \*\*Retry eligible\*\* | Yes — once, after registry health check |

| \*\*Redaction required\*\* | Schema content MUST NOT be logged at WARN or above |



\---



\### 12. Design Import Failure (`ErrAIDesignImport`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `500 Internal Server Error` |

| \*\*GraphQL error code\*\* | `INTERNAL` |

| \*\*EventsResponse.Summary\*\* | `"Failed to import generated design"` |

| \*\*EventsResponse.Details\*\* | `"The AI-generated design could not be saved to Meshery. The design content was valid but the import step failed."` |

| \*\*Log level\*\* | `ERROR` |

| \*\*Log message\*\* | `"AI design import failed: operation\_id=<id> err=<message>"` |

| \*\*mesheryctl exit code\*\* | `1` |

| \*\*mesheryctl display\*\* | `Error: Failed to save generated design. Try again.` |

| \*\*Retry eligible\*\* | Yes — idempotent with same operation ID |

| \*\*Redaction required\*\* | None |



\---



\### 13. Operation ID Missing (`ErrAIOperationIDMissing`)



| Field | Value |

|-------|-------|

| \*\*HTTP status\*\* | `400 Bad Request` |

| \*\*GraphQL error code\*\* | `BAD\_USER\_INPUT` |

| \*\*EventsResponse.Summary\*\* | `"Operation ID missing"` |

| \*\*EventsResponse.Details\*\* | `"All AI Adapter requests must include a unique operationId for event correlation."` |

| \*\*Log level\*\* | `DEBUG` |

| \*\*Log message\*\* | `"AI request missing operation\_id"` |

| \*\*mesheryctl exit code\*\* | `2` |

| \*\*mesheryctl display\*\* | `Error: Internal error — missing operation ID. Please report this bug.` |

| \*\*Retry eligible\*\* | No — client must generate a valid UUID |

| \*\*Redaction required\*\* | None |



\---



\## Redaction Rules (Summary)



The following data classes MUST be redacted before writing to any log at

`INFO` or above, and before populating `EventsResponse.Details` or

`EventsResponse.ErrorCode`:



| Data Class | Redaction Method | Max Log Level |

|------------|-----------------|---------------|

| AI provider API key / token | Replace entirely with `\[REDACTED]` | Never logged at any level |

| Provider base URL with embedded credentials | Strip userinfo component per RFC 3986 §3.2.1 | DEBUG only (stripped) |

| Raw LLM prompt text | Omit entirely | DEBUG only |

| Raw LLM response body | Truncate to 512 chars | DEBUG only |

| Meshery Model schema injected into prompt | Omit entirely | DEBUG only |



\---



\## Operation ID Correlation



Every AI Adapter operation MUST carry a client-generated `operationId` (UUID

v4). The server:



1\. Validates presence of `operationId` on receipt — returns `ErrAIOperationIDMissing` if absent.

2\. Propagates `operationId` as `EventsResponse.OperationId` in every event emitted for that operation (start, progress, success, failure).

3\. Logs `operation\_id=<id>` in every server-side log line related to that operation.

4\. Returns `operationId` in the HTTP response body so clients can correlate async events.



\---



\## Retry Eligibility Summary



| Category | Retry | Strategy |

|----------|-------|----------|

| Provider Unreachable | ✅ | Exponential backoff, max 3, 2 s initial |

| Auth Failure | ❌ | — |

| Rate Limit | ✅ | Honor `Retry-After`, default 60 s |

| Context Window | ❌ | — |

| Invalid Response | ✅ | Once, immediate |

| Design Validation | ⚠️ | Manual only |

| Empty Prompt | ❌ | — |

| Not Configured | ❌ | — |

| Timeout | ✅ | Once, 5 s delay |

| Stream Aborted | ❌ | Client re-initiates |

| Schema Injection | ✅ | Once, after registry check |

| Design Import | ✅ | Idempotent retry |

| Operation ID Missing | ❌ | — |



\---



\## Acceptance Criteria



\- \[ ] All 13 error categories have a defined HTTP status, event summary/details template, log level, and redaction rule

\- \[ ] Error codes allocated in `meshery-server-1200–1212` range, registered in `server/handlers/error.go`

\- \[ ] Redaction rules prevent API keys from appearing in any log level or event stream

\- \[ ] Operation ID correlation is enforced server-side (missing ID → 400)

\- \[ ] mesheryctl exit codes distinguish user errors (2), auth errors (4), server errors (1), and precondition failures (3)

\- \[ ] This spec is approved before any implementation PRs are opened



\---



\*Spec authored in response to \[meshery/meshery#19665](https://github.com/meshery/meshery/issues/19665).\*

