package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	models "github.com/meshery/meshery/server/models"
)

// maxExtensionPayloadSize limits the maximum allowed request body size for
// extension installation and removal requests (4 MB), preventing unbounded
// memory allocation and denial-of-service vulnerabilities via oversized payloads.
const maxExtensionPayloadSize = 4 * 1024 * 1024

// installRequest is the expected request payload
type installRequest struct {
	ExtType           string                 `json:"extType"`
	PackageURL        string                 `json:"packageUrl"`
	ExtensionMetadata map[string]interface{} `json:"extensionMetadata"`
}

type removeRequest struct {
	ExtType string `json:"extType"`
	Title   string `json:"title"`
}

// persistLocalProviderCapabilitiesForUser updates and stores local provider capabilities
// for the given user in the provider capability store. It safely skips unauthenticated
// or nil users, nil providers, or users with empty/nil UUIDs.
func persistLocalProviderCapabilitiesForUser(user *models.User, provider *models.DefaultLocalProvider) error {
	if user == nil || provider == nil {
		return nil
	}
	if user.ID.String() == "" || user.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return nil
	}
	return provider.WriteCapabilitiesForUser(user.ID.String(), &provider.ProviderProperties)
}

// InstallExtensionHandler installs an extension for the local provider.
// It verifies provider initialization, bounds request payload reads to prevent denial-of-service,
// ensures the active provider is a DefaultLocalProvider, and persists updated user capabilities.
func (h *Handler) InstallExtensionHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	if provider == nil {
		err := fmt.Errorf("provider not initialized")
		if h.log != nil {
			h.log.Error(err)
		}
		writeMeshkitError(w, err, http.StatusInternalServerError)
		return
	}

	if provider.GetProviderType() != models.LocalProviderType {
		writeMeshkitError(w, fmt.Errorf("install extension is currently supported only for the local provider"), http.StatusNotImplemented)
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	var req installRequest
	bd, err := io.ReadAll(io.LimitReader(r.Body, maxExtensionPayloadSize+1))
	if err != nil {
		if h.log != nil {
			h.log.Error(err)
		}
		writeMeshkitError(w, fmt.Errorf("failed to read request body: %w", err), http.StatusBadRequest)
		return
	}
	if len(bd) > maxExtensionPayloadSize {
		err := fmt.Errorf("request body exceeds maximum allowed size of 4MB")
		if h.log != nil {
			h.log.Error(err)
		}
		writeMeshkitError(w, err, http.StatusRequestEntityTooLarge)
		return
	}
	if len(bd) == 0 {
		writeMeshkitError(w, fmt.Errorf("empty request body"), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(bd, &req); err != nil {
		if h.log != nil {
			h.log.Error(models.ErrUnmarshal(err, "install extension request"))
		}
		writeMeshkitError(w, models.ErrUnmarshal(err, "install extension request"), http.StatusBadRequest)
		return
	}

	lp, ok := provider.(*models.DefaultLocalProvider)
	if !ok {
		writeMeshkitError(w, fmt.Errorf("local provider does not support InstallExtension"), http.StatusNotImplemented)
		return
	}

	if err := lp.InstallExtension(req.ExtType, req.PackageURL, req.ExtensionMetadata); err != nil {
		if h.log != nil {
			h.log.Error(err)
		}
		writeMeshkitError(w, fmt.Errorf("failed to install extension: %w", err), http.StatusInternalServerError)
		return
	}
	if err := persistLocalProviderCapabilitiesForUser(user, lp); err != nil {
		if h.log != nil {
			h.log.Error(err)
		}
		writeMeshkitError(w, fmt.Errorf("failed to persist updated local provider capabilities: %w", err), http.StatusInternalServerError)
		return
	}

	writeJSONEmptyObject(w, http.StatusOK)
}

// RemoveExtensionHandler removes an installed extension from local provider capabilities.
// It verifies provider initialization, bounds request payload reads to prevent denial-of-service,
// ensures the active provider is a DefaultLocalProvider, and synchronizes updated user capabilities.
func (h *Handler) RemoveExtensionHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	if provider == nil {
		err := fmt.Errorf("provider not initialized")
		if h.log != nil {
			h.log.Error(err)
		}
		writeMeshkitError(w, err, http.StatusInternalServerError)
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	var req removeRequest
	bd, err := io.ReadAll(io.LimitReader(r.Body, maxExtensionPayloadSize+1))
	if err != nil {
		if h.log != nil {
			h.log.Error(err)
		}
		writeMeshkitError(w, fmt.Errorf("failed to read request body: %w", err), http.StatusBadRequest)
		return
	}
	if len(bd) > maxExtensionPayloadSize {
		err := fmt.Errorf("request body exceeds maximum allowed size of 4MB")
		if h.log != nil {
			h.log.Error(err)
		}
		writeMeshkitError(w, err, http.StatusRequestEntityTooLarge)
		return
	}
	if len(bd) == 0 {
		writeMeshkitError(w, fmt.Errorf("empty request body"), http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(bd, &req); err != nil {
		if h.log != nil {
			h.log.Error(models.ErrUnmarshal(err, "remove extension request"))
		}
		writeMeshkitError(w, models.ErrUnmarshal(err, "remove extension request"), http.StatusBadRequest)
		return
	}

	if provider.GetProviderType() != models.LocalProviderType {
		writeMeshkitError(w, fmt.Errorf("remove extension is currently supported only for the local provider"), http.StatusNotImplemented)
		return
	}

	lp, ok := provider.(*models.DefaultLocalProvider)
	if !ok {
		writeMeshkitError(w, fmt.Errorf("local provider does not support RemoveExtension"), http.StatusNotImplemented)
		return
	}

	if err := lp.RemoveExtension(req.ExtType, req.Title); err != nil {
		if h.log != nil {
			h.log.Error(err)
		}
		writeMeshkitError(w, fmt.Errorf("failed to remove extension: %w", err), http.StatusInternalServerError)
		return
	}
	if err := persistLocalProviderCapabilitiesForUser(user, lp); err != nil {
		if h.log != nil {
			h.log.Error(err)
		}
		writeMeshkitError(w, fmt.Errorf("failed to persist updated local provider capabilities: %w", err), http.StatusInternalServerError)
		return
	}

	writeJSONEmptyObject(w, http.StatusOK)
}
