package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
)

type Router struct {
	s    *http.ServeMux
	port int
}

// New returns a new ServeMux with app routes.
func NewRouter(ctx context.Context, h models.HandlerInterface, meshClient meshes.MeshClient, port int) *Router {
	mux := http.NewServeMux()
	mux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../public/static/img/meshery-logo.png")
	}))
	mux.Handle("/play/static/", http.StripPrefix("/play/static/", http.FileServer(http.Dir("../public/static/"))))
	mux.Handle("/play/dashboard", h.AuthMiddleware(http.HandlerFunc(h.K8SConfigHandler(ctx, meshClient))))
	mux.HandleFunc("/play/", h.IndexHandler)

	mux.Handle("/play/load-test", h.AuthMiddleware(http.HandlerFunc(h.LoadTestHandler)))
	mux.Handle("/play/mesh", h.AuthMiddleware(http.HandlerFunc(h.MeshOpsHandler(ctx, meshClient))))
	mux.HandleFunc("/play/logout", h.LogoutHandler)
	mux.HandleFunc("/play/login", h.LoginHandler)

	return &Router{
		s:    mux,
		port: port,
	}
}

func (r *Router) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", r.port), r.s)
}
