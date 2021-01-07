package handlers

import (
	"encoding/json"
	"net/http"
	"plugin"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var plugins map[string]*handler.Server

// PluginHandler will handle things related to the plugin
func (h *Handler) PluginHandler(w http.ResponseWriter, r *http.Request) {
	if plugins == nil {
		plugins = make(map[string]*handler.Server)
	}
	path := r.URL.Path
	logrus.Println(path)
	switch path {
	case "/plugins/load":
		h.loadPlugins(w, r)
	default:
		plugins[path].ServeHTTP(w, r)
	}
}

func (h *Handler) loadPlugins(w http.ResponseWriter, r *http.Request) {
	plug, err := plugin.Open("/home/devkalra/Desktop/layer5labs/meshery-extensions/graphql/graphql.so")
	if err != nil {
		logrus.Fatalf("Cannot open plugin: %v", err)
	}
	symRun, err := plug.Lookup("Run")
	if err != nil {
		logrus.Fatalln("Cannot open Run()")
	}
	run := symRun.(func(*gorm.DB) *handler.Server)
	plugins["/plugins/graphql"] = run(nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("/plugins/graphql")
}
