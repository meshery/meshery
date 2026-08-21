package httputil

import (
	stderrors "errors"
	"fmt"
	"net/http"
)

// ProviderStatus is the HTTP status a failure should surface as, carried as a
// first-class value inside the error chain.
//
// It started as "the status the remote provider responded with" and that is
// still its commonest source, but it is equally the mechanism by which a
// locally raised error names the status it deserves - a rejected request body
// is the caller's fault wherever it was detected. So the value asserts a
// status, not an origin, and its message must not claim one.
//
// Why this exists: before it, the upstream status survived only as prose inside
// the MeshKit cause ("Status Code: 403 ..."). Nothing could read it back, so
// every handler that received a provider failure hardcoded a status of its own
// - overwhelmingly http.StatusNotFound. A remote-provider 403 therefore reached
// the browser as a 404, which is why a forbidden environment/workspace create
// never looked like a failure to the UI.
//
// This is internal error plumbing, not a wire construct: it is never
// serialized. The serialized error shape remains errorResponse.
type ProviderStatus int

func (s ProviderStatus) Error() string {
	return fmt.Sprintf("failure carries HTTP status %d %s", int(s), http.StatusText(int(s)))
}

// WithProviderStatus tags err with the HTTP status it should surface as - the
// one a remote provider returned, or the one a locally raised failure has
// decided on - so a handler downstream propagates it instead of fabricating a
// status of its own.
//
// errors.Join is used rather than a wrapper type so that errors.Is/errors.As
// and every meshkiterrors.Get* accessor still reach the original error: the
// MeshKit code, severity, probable cause and remediation all survive intact.
//
// A status that is not a client or server error is not worth propagating (a
// success status here would mean the provider layer mislabelled a failure), so
// err is returned unchanged and ProviderStatusCode reports "unknown", letting
// the caller pick its own defensible default.
func WithProviderStatus(err error, statusCode int) error {
	if err == nil || statusCode < http.StatusBadRequest || statusCode > 599 {
		return err
	}
	return stderrors.Join(err, ProviderStatus(statusCode))
}

// ProviderStatusCode reports the HTTP status recorded for err, if any link in
// the error chain recorded one. The boolean is false when nothing decided a
// status - a marshalling error, an unreachable provider, an unsupported
// capability - in which case the caller must choose its own.
func ProviderStatusCode(err error) (int, bool) {
	var status ProviderStatus
	if stderrors.As(err, &status) && status != 0 {
		return int(status), true
	}
	return 0, false
}

// StatusForProviderError maps a provider-layer error onto the HTTP status the
// handler should surface. When the provider recorded its upstream status
// (403, 404, 409, 422, 500 ...) that status is propagated verbatim so the
// client sees the real failure instead of a fabricated one. Otherwise the
// caller's fallback is used - http.StatusBadGateway is the defensible default
// for "we could not get a usable answer out of the remote provider".
func StatusForProviderError(err error, fallback int) int {
	if status, ok := ProviderStatusCode(err); ok {
		return status
	}
	return fallback
}
