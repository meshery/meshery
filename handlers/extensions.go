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
	plug, err := plugin.Open(path.Join(provider.PackageLocation(), provider.GetProviderProperties().Extensions.GraphQL[0].Path))
	if err != nil {
		return err
	}

	symRun, err := plug.Lookup("Run")
	if err != nil {
		return err
	}
	runFunction := symRun.(func(*models.ExtensionInput) (*models.ExtensionOutput, error))

	output, err := runFunction(&models.ExtensionInput{
		DBHandler: provider.GetGenericPersister(),
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
