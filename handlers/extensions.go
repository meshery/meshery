package handlers

import (
	"net/http"
	"path"
	"plugin"

	"github.com/layer5io/meshery/models"
)

var (
	extendedEndpoints = make(map[string]*models.Router, 0)
)

func (h *Handler) ExtensionsEndpointHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if val, ok := extendedEndpoints[req.URL.Path]; ok {
		val.HTTPHandler.ServeHTTP(w, req)
		return
	}

	http.Error(w, "Invalid endpoint", http.StatusInternalServerError)
}

func (h *Handler) LoadExtensionFromPackage(w http.ResponseWriter, req *http.Request, provider models.Provider) error {
	packagePath := ""
	if len(provider.GetProviderProperties().Extensions.GraphQL) > 0 {
		packagePath = provider.GetProviderProperties().Extensions.GraphQL[0].Path
	}

	plug, err := plugin.Open(path.Join(provider.PackageLocation(), packagePath))
	if err != nil {
		return err
	}

	// Run function
	symRun, err := plug.Lookup("Run")
	if err != nil {
		return err
	}
	runFunction := symRun.(func(*models.ExtensionInput) (*models.ExtensionOutput, error))

	output, err := runFunction(&models.ExtensionInput{
		DBHandler:       provider.GetGenericPersister(),
		MeshSyncChannel: h.meshsyncChannel,
		Logger:          h.log,
	})
	if err != nil {
		return err
	}

	// Add http endpoint to serve
	if output.Router != nil {
		extendedEndpoints[output.Router.Path] = output.Router
	}

	return nil
}
