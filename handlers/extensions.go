package handlers

import (
	"encoding/json"
	"net/http"
	"path"
	"plugin"

	"github.com/layer5io/meshery/models"
)

var (
	extendedEndpoints = make(map[string]*models.Router)
)

// Defines the version metadata for the extension
type ExtensionVersion struct {
	Version string `json:"version,omitempty"`
}

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
		return ErrPluginOpen(err)
	}

	// Run function
	symRun, err := plug.Lookup("Run")
	if err != nil {
		return ErrPluginLookup(err)
	}
	runFunction := symRun.(func(*models.ExtensionInput) (*models.ExtensionOutput, error))

	output, err := runFunction(&models.ExtensionInput{
		DBHandler:       provider.GetGenericPersister(),
		MeshSyncChannel: h.meshsyncChannel,
		BrokerConn:      h.brokerConn,
		Logger:          h.log,
	})
	if err != nil {
		return ErrPluginRun(err)
	}

	// Add http endpoint to serve
	if output.Router != nil {
		extendedEndpoints[output.Router.Path] = output.Router
	}

	return nil
}

func (h *Handler) ExtensionsVersionHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if provider.GetProviderType() == models.LocalProviderType {
		err := json.NewEncoder(w).Encode("extension not available for current provider")
		if err != nil {
			h.log.Error(ErrEncoding(err, "extension version"))
			http.Error(w, ErrEncoding(err, "extension version").Error(), http.StatusNotFound)
		}
		return
	}

	// gets the extension version from provider properties
	version := provider.GetProviderProperties().PackageVersion

	extensionVersion := &ExtensionVersion{
		Version: version,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(extensionVersion)
	if err != nil {
		h.log.Error(ErrEncoding(err, "extension version"))
		http.Error(w, ErrEncoding(err, "extension version").Error(), http.StatusNotFound)
	}
}
