// Package httputil centralizes the canonical HTTP response helpers used across
// server/handlers and server/models. Factored out of server/handlers/utils.go so
// both packages can emit the same JSON response shape without an import cycle
// (server/models cannot import server/handlers because handlers already depends
// on models).
//
// Reach for:
//
//   - WriteMeshkitError: ANY error path. If err wraps a *meshkiterrors.Error or
//     *ErrorV2, the code, severity, probable cause and remediation survive onto
//     the wire. If it doesn't, the .Error() string is still emitted as JSON.
//   - WriteJSONError: error paths where the message is a bare string with no
//     MeshKit wrapper. Prefer promoting the string to a MeshKit error and using
//     WriteMeshkitError instead.
//   - WriteJSONMessage: success paths that return a small status or result
//     payload (e.g. {"message": "deleted"}).
//   - WriteJSONEmptyObject: success paths that need to return an empty JSON
//     object ({}) with the Content-Type header set.
//
// Never use http.Error from handlers or the provider layer — it emits
// Content-Type: text/plain which crashes RTK Query's default baseQuery on the
// UI (see docs/content/en/project/contributing/error-contract.md).
package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

// WriteJSONError writes a JSON-encoded {"error": message} body with the given
// HTTP status. Using JSON (instead of http.Error's plain text) keeps client
// response parsers — notably RTK Query's default baseQuery, which parses by
// Content-Type and treats application/json as JSON — from choking on error
// bodies that happen to start with a letter (e.g. "WorkspaceID or OrgID ...").
func WriteJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// errorResponse is the wire shape for all non-2xx responses from Meshery Server.
// Fields mirror github.com/meshery/meshkit/errors.Error; omitempty keeps the
// body small for bare-string errors that carry no MeshKit metadata.
type errorResponse struct {
	Error                string   `json:"error"`
	Code                 string   `json:"code,omitempty"`
	Severity             string   `json:"severity,omitempty"`
	ProbableCause        []string `json:"probableCause,omitempty"`
	SuggestedRemediation []string `json:"suggestedRemediation,omitempty"`
	LongDescription      []string `json:"longDescription,omitempty"`
}

// WriteMeshkitError writes a JSON error response that preserves MeshKit error
// metadata (code, severity, probable cause, remediation) when err is (or wraps)
// a *meshkiterrors.Error or *meshkiterrors.ErrorV2. Non-MeshKit errors still
// produce a JSON body — they just carry only the .Error() string, matching
// WriteJSONError's shape so clients never see plain text from this package.
//
// Prefer this over http.Error for every handler error path. RTK Query's
// baseQuery dispatches on Content-Type and crashes on plain-text bodies that
// happen to start with a letter (e.g. "Status Code: 404 ...").
func WriteMeshkitError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)

	resp := errorResponse{}
	if err == nil {
		resp.Error = http.StatusText(status)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Error = err.Error()

	// Populate MeshKit fields only when the error carries them. GetCode etc.
	// return the "None" sentinel for non-MeshKit errors; treat that as absent.
	if code := meshkiterrors.GetCode(err); code != "" && code != "None" {
		resp.Code = code
		resp.Severity = severityString(meshkiterrors.GetSeverity(err))
		if short := meshkiterrors.GetSDescription(err); short != "" && short != "None" {
			// Use ShortDescription as the user-facing `error` when available —
			// err.Error() on a MeshKit error concatenates every field with pipes.
			resp.Error = short
		}

		// Prefer the underlying string slices over MeshKit's Get* accessors.
		// The accessors join every entry with "." into one string, so an error
		// listing three distinct probable causes arrived on the wire as a single
		// run-on line ("...permission..No user exists..The provider..") and the
		// UI rendered it as one bullet. The wire contract has always declared
		// these as arrays; this makes them arrays in fact.
		causes, remedies, longs := meshkitDetailSlices(err)
		resp.ProbableCause = causes
		resp.SuggestedRemediation = remedies
		resp.LongDescription = longs
	}

	_ = json.NewEncoder(w).Encode(resp)
}

// meshkitDetailSlices returns the probable-cause, suggested-remediation and
// long-description slices carried by err, with MeshKit's "None" sentinel and
// blank entries removed. Each result is nil when the error carries nothing for
// that field, so `omitempty` keeps it off the wire entirely.
func meshkitDetailSlices(err error) (causes, remedies, longs []string) {
	var errV2 *meshkiterrors.ErrorV2
	if errors.As(err, &errV2) && errV2 != nil {
		return compactDetails(errV2.ProbableCause), compactDetails(errV2.SuggestedRemediation), compactDetails(errV2.LongDescription)
	}

	var errV1 *meshkiterrors.Error
	if errors.As(err, &errV1) && errV1 != nil {
		return compactDetails(errV1.ProbableCause), compactDetails(errV1.SuggestedRemediation), compactDetails(errV1.LongDescription)
	}

	return nil, nil, nil
}

// compactDetails drops blank entries and MeshKit's "None" placeholder, and
// returns nil rather than an empty slice so `omitempty` keeps the field off the
// wire entirely.
func compactDetails(in []string) []string {
	out := make([]string, 0, len(in))
	for _, entry := range in {
		trimmed := strings.TrimSpace(entry)
		if trimmed == "" || trimmed == "None" {
			continue
		}
		out = append(out, trimmed)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// severityString converts a MeshKit Severity enum to the string label used on
// the wire. Kept here (not in MeshKit) because MeshKit's Severity.String is
// not yet exported in all versions we pin.
func severityString(s meshkiterrors.Severity) string {
	switch s {
	case meshkiterrors.Emergency:
		return "EMERGENCY"
	case meshkiterrors.Alert:
		return "ALERT"
	case meshkiterrors.Critical:
		return "CRITICAL"
	case meshkiterrors.Fatal:
		return "FATAL"
	default:
		return "ERROR"
	}
}

// WriteJSONMessage encodes an arbitrary payload as JSON with the given status
// code. Use for success responses that currently write a bare string (e.g.
// "Database reset successful") — promote them to a structured message.
func WriteJSONMessage(w http.ResponseWriter, payload any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// WriteJSONEmptyObject writes "{}" with Content-Type: application/json.
// Several handlers return an empty object to indicate "no data, but request
// succeeded"; this helper centralizes the pattern and guarantees the header
// so clients (notably RTK Query's Content-Type-dispatching baseQuery) never
// see a bare "{}" with no declared Content-Type.
func WriteJSONEmptyObject(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_, _ = w.Write([]byte("{}"))
}
