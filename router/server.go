package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/models"
)

type Router struct {
	s    *http.ServeMux
	port int
}

// New returns a new ServeMux with app routes.
func NewRouter(ctx context.Context, h models.HandlerInterface, port int) *Router {
	mux := http.NewServeMux()

	mux.Handle("/api/user", h.AuthMiddleware(http.HandlerFunc(h.UserHandler(ctx))))
	mux.Handle("/api/k8sconfig", h.AuthMiddleware(http.HandlerFunc(h.K8SConfigHandler(ctx))))
	mux.Handle("/api/load-test", h.AuthMiddleware(http.HandlerFunc(h.LoadTestHandler)))
	mux.Handle("/api/mesh", h.AuthMiddleware(http.HandlerFunc(h.MeshOpsHandler(ctx))))

	mux.HandleFunc("/logout", h.LogoutHandler)
	mux.HandleFunc("/login", h.LoginHandler)

	// TODO: have to change this too
	mux.Handle("/favicon.ico", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600") // 1 hr
		http.ServeFile(w, r, "../public/static/img/meshery-logo.png")
	}))
	mux.Handle("/", h.AuthMiddleware(http.FileServer(http.Dir("../ui/out/"))))
	// mux.Handle("/static/", h.AuthMiddleware(http.StripPrefix("/static/", http.FileServer(http.Dir("../ui/out/_next/static/")))))
	// mux.Handle("/", h.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	logrus.Debugf("requesting index.html file")
	// 	_, err := os.Stat("../ui/out/index.html")
	// 	if err != nil {
	// 		if os.IsNotExist(err) {
	// 			logrus.Debugf("index file not exists")
	// 		}
	// 	}
	// 	http.ServeFile(w, r, "../ui/out/index.html")
	// })))

	return &Router{
		s:    mux,
		port: port,
	}
}

func (r *Router) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", r.port), r.s)
}
