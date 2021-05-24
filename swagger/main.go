package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"os/signal"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Router represents router
type Router struct {
	S    *mux.Router
	port int
}

// newRouterDocs returns a new ServeMux with app routes.
func newRouterDocs(ctx context.Context, port int) *Router {
	gMux := mux.NewRouter()

	swaggerOpts := middleware.SwaggerUIOpts{SpecURL: "./swagger.yaml"}
	// opts := middleware.RedocOpts{SpecURL: "./swagger.yaml"}

	swaggerSh := middleware.SwaggerUI(swaggerOpts, nil)
	// sh := middleware.Redoc(opts, nil)

	gMux.Handle("/docs", swaggerSh)
	// gMux.Handle("/swagger", swaggerSh)

	gMux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/docs", http.StatusFound)
	}))

	gMux.Handle("/swagger.yaml", http.FileServer(http.Dir("./swagger")))

	return &Router{
		S:    gMux,
		port: port,
	}
}

// Run starts the http server
func main() {
	ctx := context.Background()

	port := 9090
	r := newRouterDocs(ctx, port)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		logrus.Infof("Starting Server listening on :%d", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", r.port), r.S); err != nil {
			logrus.Fatalf("ListenAndServe Error: %v", err)
		}
	}()
	<-c
	logrus.Info("Shutting down Meshery-OpenAPI")
}
