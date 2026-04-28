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
	mx                sync.Mutex
)

// Defines the version metadata for the extension
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

func (h *Handler) LoadExtensionFromPackage(_ http.ResponseWriter, _ *http.Request, provider models.Provider) error {
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

/*
* ExtensionsHandler is a handler function which works as a proxy to resolve the
* request of any extension point to its remote provider
 */
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
