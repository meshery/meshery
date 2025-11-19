package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path"

	"github.com/spf13/viper"
)

// HealthStatus represents the health status response for the healthz endpoints.
// It includes information about the system's health including whether capabilities
// are loaded and if running in Kanvas mode, whether extensions are present.
type HealthStatus struct {
	Status             string `json:"status"`                         // "healthy" or "unhealthy"
	CapabilitiesLoaded bool   `json:"capabilities_loaded,omitempty"`  // Whether provider capabilities are loaded
	ExtensionExists    bool   `json:"extension_exists,omitempty"`     // Whether extension package exists (Kanvas mode only)
	ReleaseChannel     string `json:"release_channel,omitempty"`      // Current release channel (e.g., "stable", "kanvas")
	Message            string `json:"message,omitempty"`              // Error message if unhealthy
}

func (h *Handler) K8sHealthzHandler(w http.ResponseWriter, r *http.Request) {
	healthStatus := HealthStatus{
		Status: "healthy",
	}

	// Get the release channel
	releaseChannel := viper.GetString("RELEASE_CHANNEL")
	healthStatus.ReleaseChannel = releaseChannel

	// Check 1: Verify capabilities are loaded by checking provider properties
	capabilitiesLoaded := false
	for _, provider := range h.config.Providers {
		providerProps := provider.GetProviderProperties()
		// Check if capabilities have been loaded (non-empty capabilities object)
		if len(providerProps.Capabilities) > 0 {
			capabilitiesLoaded = true
			break
		}
	}
	healthStatus.CapabilitiesLoaded = capabilitiesLoaded

	// Check 2: If running in Kanvas mode, verify extension exists in filesystem
	if releaseChannel == "kanvas" {
		extensionExists := false

		// Check all providers for navigator extensions in Kanvas mode
		for _, provider := range h.config.Providers {
			providerProps := provider.GetProviderProperties()

			// Check if navigator extensions are configured
			if len(providerProps.Extensions.Navigator) > 0 {
				// Check if the extension package exists on filesystem
				packageLocation := provider.PackageLocation()
				if packageLocation != "" {
					// Check if GraphQL extension path exists
					if len(providerProps.Extensions.GraphQL) > 0 {
						packagePath := providerProps.Extensions.GraphQL[0].Path
						fullPath := path.Join(packageLocation, packagePath)
						if _, err := os.Stat(fullPath); err == nil {
							extensionExists = true
							break
						}
					}
				}
			}
		}

		healthStatus.ExtensionExists = extensionExists

		// If in Kanvas mode, we must have extensions
		if !extensionExists {
			healthStatus.Status = "unhealthy"
			healthStatus.Message = "Extension package not found in filesystem"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(healthStatus)
			return
		}
	}

	// If capabilities are not loaded, return unhealthy status
	if !capabilitiesLoaded {
		healthStatus.Status = "unhealthy"
		healthStatus.Message = "Provider capabilities not loaded"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(healthStatus)
		return
	}

	// All checks passed
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthStatus)
}
