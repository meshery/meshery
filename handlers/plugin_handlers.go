package handlers

import (
	"encoding/json"
	"net/http"
	"plugin"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/database"
	"github.com/sirupsen/logrus"
)

var plugins map[string]*handler.Server

// PluginHandler will handle things related to the plugin
func (h *Handler) PluginHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if plugins == nil {
		plugins = make(map[string]*handler.Server)
	}
	path := req.URL.Path
	logrus.Println(path)
	switch path {
	case "/plugins/load":
		h.loadPlugins(w, req, provider)
	default:
		plugins[path].ServeHTTP(w, req)
	}
}

func (h *Handler) loadPlugins(w http.ResponseWriter, r *http.Request, provider models.Provider) {
	plug, err := plugin.Open("path-to-so-file")
	if err != nil {
		logrus.Printf("Cannot open plugin: %v", err)
	}
	symRun, err := plug.Lookup("Run")
	if err != nil {
		logrus.Println("Cannot open Run()")
	}
	run := symRun.(func(*database.Handler) *handler.Server)
	plugins["/plugins/graphql"] = run(provider.GetGenericPersister())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("/plugins/graphql")
}
