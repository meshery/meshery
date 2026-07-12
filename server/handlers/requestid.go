package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/httputil"
)

// RequestIDMiddleware generates or propagates a unique request ID for every
// incoming HTTP request. The ID is:
//  1. Read from the X-Request-ID header if present (for distributed tracing
//     where an upstream proxy or ingress already assigned one).
//  2. Generated as a v4 UUID otherwise.
//
// The ID is stored in the request context and set on the response via the
// X-Request-ID header so clients and downstream middleware can correlate
// error bodies, log entries, and traces.
//
// Downstream propagation: handlers that make outbound calls (gRPC adapters,
// remote providers, database queries) should extract the request ID via
// httputil.RequestIDFromRequest(r) and forward it in the outgoing context or
// metadata headers. Error responses can surface the ID by using
// httputil.WriteMeshkitErrorWithRequest(w, r, err, status) which populates
// the "requestId" field in the error envelope.
func (h *Handler) RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			uid, err := uuid.NewV4()
			if err != nil {
				h.log.Warn(fmt.Errorf("failed to generate request ID for %s %s: %w", r.Method, r.URL.Path, err))
				reqID = ""
			} else {
				reqID = uid.String()
			}
		}

		// Set the response header early so all downstream writes carry it.
		w.Header().Set("X-Request-ID", reqID)

		// Store in context using httputil's key so that
		// httputil.RequestIDFromContext and httputil.RequestIDFromRequest
		// can retrieve it without a circular import into server/models.
		ctx := httputil.NewContextWithRequestID(r.Context(), reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
