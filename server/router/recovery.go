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

				// A hijacked connection (WebSocket) is no longer HTTP, so there is
				// nothing to answer on it.
				if rw.hijacked {
					return
				}
				// The client already holds part of the response. Returning now
				// would let net/http finish it as if it were complete, so a
				// truncated body would look whole. Abort it instead, quietly,
				// as net/http itself does for http.ErrAbortHandler.
				if rw.wroteHeader {
					panic(http.ErrAbortHandler)
				}
				// Headers the handler set describe the body it meant to send,
				// not this error body. A stale length or encoding would stop the
				// client from reading it.
				h := w.Header()
				h.Del("Content-Length")
				h.Del("Content-Encoding")
				httputil.WriteMeshkitError(w, models.ErrHandlerPanic(), http.StatusInternalServerError)
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
	// Interim 1xx statuses (100, 103) are followed by the real one, so the
	// response has not started and a 500 can still be sent. 101 Switching
	// Protocols is final. The flag is set only after the delegated call
	// returns: net/http panics on a status it rejects (outside 100-999) before
	// committing anything, and then a 500 can still be sent.
	w.ResponseWriter.WriteHeader(code)
	if code >= http.StatusOK || code == http.StatusSwitchingProtocols {
		w.wroteHeader = true
	}
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
