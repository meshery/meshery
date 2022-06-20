package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"plugin"
	"sync"

	"github.com/layer5io/meshery/models"
)

var (
	//USE WITH CAUTION: Wherever read/write is performed, use this in a thread safe way, using the global mutex
	extendedEndpoints = make(map[string]*models.Router)
	mx                sync.Mutex
)

// Defines the version metadata for the extension
type ExtensionVersion struct {
	Version string `json:"version,omitempty"`
}

func (h *Handler) ExtensionsEndpointHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	mx.Lock()
	val, ok := extendedEndpoints[req.URL.Path]
	mx.Unlock()
	if ok {
		val.HTTPHandler.ServeHTTP(w, req)
		return
	}

	http.Error(w, "Invalid endpoint", http.StatusInternalServerError)
}

func (h *Handler) LoadExtensionFromPackage(w http.ResponseWriter, req *http.Request, provider models.Provider) error {
	mx.Lock()
	defer mx.Unlock()
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

/*
* ExtensionsHandler is a handler function which works as a proxy to resolve the
* request of any extension point to its remote provider
 */
func (h *Handler) ExtensionsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	resp, err := provider.ExtensionProxy(req)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}
