package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

// K8sHealthzHandler implements Kubernetes-style health check endpoints (/healthz/live and /healthz/ready).
// Following Kubernetes best practices, it returns:
// - HTTP 200 with plain text "ok" when healthy
// - HTTP 503 with plain text error details when unhealthy
// - Supports ?verbose=1 query parameter for detailed check results
func (h *Handler) K8sHealthzHandler(w http.ResponseWriter, r *http.Request) {
	// Parse verbose flag from query parameters
	verbose := r.URL.Query().Get("verbose") == "1"

	// Collect check results
	checks := []healthCheck{}

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

	if capabilitiesLoaded {
		checks = append(checks, healthCheck{name: "capabilities", status: checkOK, reason: "provider capabilities loaded"})
	} else {
		checks = append(checks, healthCheck{name: "capabilities", status: checkFailed, reason: "provider capabilities not loaded"})
	}

	// Check 2: Check for extensions (informational only, does not affect health status)
	extensionExists := false

	// Check all providers for navigator extensions
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

	// Add extension check as informational (uses checkInfo status)
	if extensionExists {
		checks = append(checks, healthCheck{name: "extension", status: checkInfo, reason: "extension package found"})
	} else {
		checks = append(checks, healthCheck{name: "extension", status: checkInfo, reason: "extension package not found"})
	}

	// Determine overall health status (only checkFailed affects health, checkInfo is informational)
	allHealthy := true
	for _, check := range checks {
		if check.status == checkFailed {
			allHealthy = false
			break
		}
	}

	// Format response following Kubernetes conventions
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if allHealthy {
		w.WriteHeader(http.StatusOK)
		if verbose {
			// Return detailed check results
			for _, check := range checks {
				switch check.status {
				case checkOK:
					_, _ = fmt.Fprintf(w, "[+]%s ok\n", check.name)

				case checkInfo:
					// Informational checks show status without +/- prefix
					_, _ = fmt.Fprintf(w, "[i]%s %s\n", check.name, check.reason)
				}
			}
			fmt.Fprint(w, "healthz check passed\n")
		} else {
			// Simple "ok" response
			fmt.Fprint(w, "ok")
		}
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		if verbose {
			// Return detailed check results with failures
			for _, check := range checks {
				switch check.status {
				case checkOK:
					_, _ = fmt.Fprintf(w, "[+]%s ok\n", check.name)

				case checkFailed:
					_, _ = fmt.Fprintf(w, "[-]%s failed: %s\n", check.name, check.reason)

				case checkInfo:
					// Show informational checks even in failure state
					_, _ = fmt.Fprintf(w, "[i]%s %s\n", check.name, check.reason)
				}
			}
		} else {
			// Return simple error message (only failed checks, not informational)
			var failedChecks []string
			for _, check := range checks {
				if check.status == checkFailed {
					failedChecks = append(failedChecks, fmt.Sprintf("%s: %s", check.name, check.reason))
				}
			}
			_, _ = fmt.Fprintf(w, "healthz check failed: %s", strings.Join(failedChecks, "; "))
		}
	}
}

// healthCheck represents the result of a health check
type healthCheck struct {
	name   string
	status checkStatus
	reason string
}

// checkStatus represents the status of a health check
type checkStatus int

const (
	checkOK checkStatus = iota
	checkFailed
	checkInfo // Informational check that doesn't affect health status
)
