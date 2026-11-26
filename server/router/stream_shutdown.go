package router

import (
	"context"
	"net/http"
)

// endOnShutdown wraps a Server-Sent Events handler that serves until its
// client leaves, so that its request context is also cancelled once
// shuttingDown is done. NewRouter cancels shuttingDown as soon as the server
// starts shutting down.
//
// http.Server.Shutdown closes the listener and then waits for every active
// connection to go idle, but it never cancels a request's context. Such a
// handler returns only when its request context is done, so without this
// wrapper a single open stream holds Shutdown for its whole timeout. Every
// open Meshery UI tab holds the events stream, and with an operator-mode
// Kubernetes connection the controller status stream too. Once that timeout
// expires, Shutdown returns with the handler still running. Other requests,
// including load-test runs, are not wrapped, so Shutdown still lets them
// finish.
func endOnShutdown(shuttingDown context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		stop := context.AfterFunc(shuttingDown, cancel)
		defer stop()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
