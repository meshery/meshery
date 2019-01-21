package handlers

import (
	"context"
	"net/http"
)

// New returns a new ServeMux with app routes.
func New(ctx context.Context, config *ServerConfig) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/play/static/", http.StripPrefix("/play/static/", http.FileServer(http.Dir("../public/static/"))))
	mux.Handle("/play/dashboard", authMiddleware(http.HandlerFunc(dashboardHandler(ctx, config.MeshClient))))
	mux.HandleFunc("/play/", indexHandler)

	mux.Handle("/play/load-test", authMiddleware(http.HandlerFunc(loadTestHandler)))
	mux.Handle("/play/mesh", authMiddleware(http.HandlerFunc(meshOpsHandler(ctx, config.MeshClient))))
	mux.HandleFunc("/play/logout", logoutHandler)
	mux.Handle("/play/tweet", authMiddleware(http.HandlerFunc(tweetHandler)))
	mux.HandleFunc("/play/login", loginHandler)

	return mux
}
