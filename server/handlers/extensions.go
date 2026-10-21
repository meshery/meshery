package handlers

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"net/http"
	"path"
	"plugin"
	"sync"

	"github.com/meshery/meshery/server/extensions"
	"github.com/meshery/meshery/server/models"
)

var (
	//USE WITH CAUTION: Wherever read/write is performed, use this in a thread safe way, using the global mutex
	extendedEndpoints = make(map[string]*extensions.Router)
	// loadedExtensionPackages records the outcome of loading each extension
	// plugin, keyed by its path, so a plugin is opened at most once per process.
	// Go cannot unload a plugin, and re-opening one whose first load failed
	// leaves the runtime without fresh module data, which is a fatal,
	// unrecoverable error ("runtime: no plugin module data"). Guarded by mx.
	loadedExtensionPackages = make(map[string]error)
	mx                      sync.Mutex
)

// ExtensionVersion defines the version metadata for the extension
type ExtensionVersion struct {
	Version string `json:"version,omitempty"`
}

func (h *Handler) ExtensionsEndpointHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	mx.Lock()
	val, ok := extendedEndpoints[req.URL.Path]
	mx.Unlock()
	if ok {
		val.HTTPHandler.ServeHTTP(w, req)
		return
	}

	writeMeshkitError(w, ErrExtensionEndpointNotRegistered(req.URL.Path), http.StatusNotFound)
}

// LoadExtensionFromPackage loads the provider's extension plugin once and
// returns the outcome of that first attempt on every later call, so a request
// never triggers a second plugin.Open of the same package.
func (h *Handler) LoadExtensionFromPackage(_ http.ResponseWriter, _ *http.Request, provider models.Provider) error {
	mx.Lock()
	defer mx.Unlock()
	packagePath := ""
	if len(provider.GetProviderProperties().Extensions.GraphQL) > 0 {
		packagePath = provider.GetProviderProperties().Extensions.GraphQL[0].Path
	}
	pluginPath := path.Join(provider.PackageLocation(), packagePath)

	if err, loaded := loadedExtensionPackages[pluginPath]; loaded {
		return err
	}
	err := h.loadExtensionPlugin(pluginPath, provider)
	loadedExtensionPackages[pluginPath] = err
	return err
}

// loadExtensionPlugin opens the plugin at pluginPath, runs it and registers
// its HTTP endpoint. Callers must hold mx.
func (h *Handler) loadExtensionPlugin(pluginPath string, provider models.Provider) error {
	plug, err := plugin.Open(pluginPath)
	if err != nil {
		return ErrPluginOpen(err)
	}

	// Run function
	symRun, err := plug.Lookup("Run")
	if err != nil {
		return ErrPluginLookup(err)
	}
	runFunction := symRun.(func(*extensions.ExtensionInput) (*extensions.ExtensionOutput, error))

	output, err := runFunction(&extensions.ExtensionInput{
		DBHandler:            provider.GetGenericPersister(),
		MeshSyncChannel:      h.MeshsyncChannel,
		BrokerConn:           h.brokerConn,
		Logger:               h.log,
		K8sConnectionTracker: h.ConnectionToStateMachineInstanceTracker,
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

func (h *Handler) ExtensionsVersionHandler(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	if provider.GetProviderType() == models.LocalProviderType {
		err := json.NewEncoder(w).Encode("extension not available for current provider")
		if err != nil {
			h.log.Error(models.ErrEncoding(err, "extension version"))
			writeMeshkitError(w, models.ErrEncoding(err, "extension version"), http.StatusNotFound)
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
		h.log.Error(models.ErrEncoding(err, "extension version"))
		writeMeshkitError(w, models.ErrEncoding(err, "extension version"), http.StatusNotFound)
	}
}

// ExtensionsHandler is a handler function which works as a proxy to resolve the
// request of any extension point to its remote provider
func (h *Handler) ExtensionsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	resp, err := provider.ExtensionProxy(req)
	if err != nil {
		// ExtensionProxy can fail for two distinct classes of reasons:
		//   1. The active provider doesn't support the call at all — the
		//      local provider always returns ErrLocalProviderSupport,
		//      which is a configuration mismatch (extensions only work
		//      against a remote provider) rather than an upstream
		//      failure. Surface as 501 Not Implemented.
		//   2. Anything else surfaced by the remote provider —
		//      unreachable upstream, request-construction failure,
		//      token retrieval failure, body-read failure. These all
		//      represent failure to obtain a response from the
		//      configured upstream, so 502 Bad Gateway is appropriate.
		// Wrap with ErrExtensionProxy in both cases so the JSON
		// envelope carries MeshKit metadata.
		status := http.StatusBadGateway
		if stderrors.Is(err, models.ErrLocalProviderSupport) {
			status = http.StatusNotImplemented
		}
		wrappedErr := ErrExtensionProxy(err)
		h.log.Error(wrappedErr)
		writeMeshkitError(w, wrappedErr, status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	if _, err := fmt.Fprint(w, string(resp.Body)); err != nil {
		h.log.Error(err)
	}
}
