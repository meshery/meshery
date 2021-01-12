package handlers

import (
	"net/http"
	"plugin"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/database"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	loadedPlugins = map[string]*Plugin{
		"/api/plugins/meshmap": &Plugin{
			state:      false,
			targetPath: "/Users/abishekk/Documents/layer5/meshery/cmd/plugin.so",
		},
	}

	loadedEndpoints = map[string]*Endpoint{
		"/api/plugins/graphql": &Endpoint{},
	}
)

type Endpoint struct {
	handler *handler.Server
}

type Plugin struct {
	state      bool
	targetPath string
}

func (h *Handler) PluginHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {

	if val, ok := loadedEndpoints[req.URL.Path]; ok {
		val.handler.ServeHTTP(w, req)
		return
	}

	if _, ok := loadedPlugins[req.URL.Path]; ok {
		if loadedPlugins[req.URL.Path].state {
			_, _ = w.Write([]byte("already initialized"))
			return
		}
	}

	err := h.loadPlugins(w, req, provider)
	if err != nil {
		msg := "unable to load plugin"
		err = errors.Wrapf(err, msg)
		logrus.Error(err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	loadedPlugins[req.URL.Path].state = true
	_, _ = w.Write([]byte("initialized"))
}

func (h *Handler) loadPlugins(w http.ResponseWriter, req *http.Request, provider models.Provider) error {
	plug, err := plugin.Open(loadedPlugins[req.URL.Path].targetPath)
	if err != nil {
		return err
	}

	symRun, err := plug.Lookup("Run")
	if err != nil {
		return err
	}

	run := symRun.(func(*database.Handler) (*handler.Server, error))
	graphqlHandler, err := run(provider.GetGenericPersister())
	if err != nil {
		return err
	}

	loadedEndpoints["/api/plugins/graphql"].handler = graphqlHandler

	return nil
}
