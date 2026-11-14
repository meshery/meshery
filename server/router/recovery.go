package router

import (
	"bufio"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/httputil"
	"github.com/meshery/meshkit/logger"
)

// Option configures NewRouter.
type Option func(*routerOptions)

type routerOptions struct {
	log logger.Handler
}

// WithLogger sets the logger the recovery middleware writes contained panics
// to. Without it a panic is still contained, counted and passed to the
// registered hooks; it is just not logged.
func WithLogger(log logger.Handler) Option {
	return func(o *routerOptions) { o.log = log }
}

// recoveryMiddleware answers a panicking HTTP handler with a structured 500
// instead of an aborted connection.
//
// net/http already stops a handler panic from crashing the process: it logs the
// panic and closes the connection. What the client gets then is a reset with no
// body, and the panic is invisible to Meshery's own counters and hooks. This
// middleware turns it into the same MeshKit error body every other failure
// uses, logs the stack server-side only, and reports it through
// utils.ReportRecoveredPanic like any other contained panic.
//
// It does not help a goroutine a handler starts: that needs its own guard
// (utils.SafeGo or utils.Guard).
func recoveryMiddleware(log logger.Handler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &recoveryWriter{ResponseWriter: w}
			defer func() {
				recovered := recover()
				if recovered == nil {
					return
				}
				// net/http's own signal to abort a response quietly; it is not
				// a defect and must reach net/http untouched.
				if recovered == http.ErrAbortHandler {
					panic(recovered)
				}

				utils.ReportRecoveredPanic("http "+r.Method+" "+routeLabel(r), recovered, log)

				// Once the response has started, or the connection has been
				// taken over (WebSocket), the status can no longer change and a
				// body would be spliced into the middle of what was sent.
				if rw.wroteHeader || rw.hijacked {
					return
				}
				httputil.WriteMeshkitError(w, models.ErrHandlerPanic(fmt.Errorf("panic: %v", recovered)), http.StatusInternalServerError)
			}()
			next.ServeHTTP(rw, r)
		})
	}
}

// routeLabel names a request by its route template, so the label stays
// low-cardinality (no ids) when it is used as a counter key.
func routeLabel(r *http.Request) string {
	if route := mux.CurrentRoute(r); route != nil {
		if tmpl, err := route.GetPathTemplate(); err == nil && tmpl != "" {
			return tmpl
		}
	}
	return "unmatched"
}

// recoveryWriter records whether the response has started, so the middleware
// knows if it can still send a 500. It passes Flush and Hijack through, because
// Meshery serves streaming responses and GraphQL subscriptions over WebSocket;
// wrapping the writer must not take either away.
type recoveryWriter struct {
	http.ResponseWriter
	wroteHeader bool
	hijacked    bool
}

func (w *recoveryWriter) WriteHeader(code int) {
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *recoveryWriter) Write(b []byte) (int, error) {
	w.wroteHeader = true
	return w.ResponseWriter.Write(b)
}

// Flush implements http.Flusher. Flushing sends the headers, so the response
// has started.
func (w *recoveryWriter) Flush() {
	w.wroteHeader = true
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements http.Hijacker for WebSocket upgrades.
func (w *recoveryWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the underlying ResponseWriter does not support hijacking")
	}
	conn, rw, err := h.Hijack()
	if err == nil {
		w.hijacked = true
	}
	return conn, rw, err
}

// Unwrap lets http.NewResponseController reach the underlying writer.
func (w *recoveryWriter) Unwrap() http.ResponseWriter { return w.ResponseWriter }
