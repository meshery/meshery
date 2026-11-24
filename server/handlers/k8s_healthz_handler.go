package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/meshery/meshkit/logger"
)

// K8sHealthzHandler is a backward-compatible entry point that dispatches explicitly
// to K8sLivenessHandler for /healthz/live and /healthz/live/, or to
// K8sReadinessHandler for /healthz/ready and /healthz/ready/.
// Noncanonical, traversal, or unrecognized paths are rejected with HTTP 404 Not Found.
func (h *Handler) K8sHealthzHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/healthz/live", "/healthz/live/":
		h.K8sLivenessHandler(w, r)
	case "/healthz/ready", "/healthz/ready/":
		h.K8sReadinessHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

// K8sLivenessHandler implements the Kubernetes liveness probe (/healthz/live).
// Following Kubernetes conventions, it returns:
// - HTTP 200 with plain text "ok" when healthy
// - Supports ?verbose=1 query parameter for detailed check results
// Liveness verifies process viability and responsiveness without depending on
// asynchronously initialized or remote dependencies like provider capabilities or databases.
func (h *Handler) K8sLivenessHandler(w http.ResponseWriter, r *http.Request) {
	verbose := r.URL.Query().Get("verbose") == "1"

	checks := []healthCheck{
		{name: "live", status: checkOK, reason: "server is responsive"},
	}

	formatHealthResponse(w, checks, verbose, h.log)
}

// K8sReadinessHandler implements the Kubernetes readiness probe (/healthz/ready).
// It returns:
// - HTTP 200 with plain text "ok" when ready to serve user traffic
// - HTTP 503 with plain text error details when dependencies are not ready
// - Supports ?verbose=1 query parameter for detailed check results
// Readiness verifies:
// 1. Database connectivity via PingContext (bounded timeout <= 2s)
// 2. Provider capabilities loaded
// 3. Extension package existence on filesystem (informational only)
func (h *Handler) K8sReadinessHandler(w http.ResponseWriter, r *http.Request) {
	verbose := r.URL.Query().Get("verbose") == "1"

	checks := []healthCheck{}

	// Check 1: Database connectivity
	if h.dbHandler == nil || h.dbHandler.DB == nil {
		checks = append(checks, healthCheck{name: "database", status: checkFailed, reason: "database handler not initialized"})
	} else {
		sqlDB, err := h.dbHandler.DB.DB()
		if err != nil {
			checks = append(checks, healthCheck{name: "database", status: checkFailed, reason: fmt.Sprintf("failed to get sql.DB: %v", err)})
		} else {
			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			err := sqlDB.PingContext(ctx)
			cancel()
			if err != nil {
				checks = append(checks, healthCheck{name: "database", status: checkFailed, reason: fmt.Sprintf("database ping failed: %v", err)})
			} else {
				checks = append(checks, healthCheck{name: "database", status: checkOK, reason: "database responsive"})
			}
		}
	}

	// Check 2: Verify capabilities are loaded by checking provider properties
	capabilitiesLoaded := false
	if h.config != nil {
		for _, provider := range h.config.Providers {
			providerProps := provider.GetProviderProperties()
			// Check if capabilities have been loaded (non-empty capabilities object)
			if len(providerProps.Capabilities) > 0 {
				capabilitiesLoaded = true
				break
			}
		}
	}

	if capabilitiesLoaded {
		checks = append(checks, healthCheck{name: "capabilities", status: checkOK, reason: "provider capabilities loaded"})
	} else {
		checks = append(checks, healthCheck{name: "capabilities", status: checkFailed, reason: "provider capabilities not loaded"})
	}

	// Check 3: Check for extensions (informational only, does not affect health status)
	extensionExists := false
	if h.config != nil {
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
	}

	// Add extension check as informational (uses checkInfo status)
	if extensionExists {
		checks = append(checks, healthCheck{name: "extension", status: checkInfo, reason: "extension package found"})
	} else {
		checks = append(checks, healthCheck{name: "extension", status: checkInfo, reason: "extension package not found"})
	}

	formatHealthResponse(w, checks, verbose, h.log)
}

// formatHealthResponse formats the HTTP response following Kubernetes health check conventions.
func formatHealthResponse(w http.ResponseWriter, checks []healthCheck, verbose bool, log logger.Handler) {
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
			if _, err := fmt.Fprint(w, "healthz check passed\n"); err != nil && log != nil {
				log.Error(err)
			}
		} else {
			// Simple "ok" response
			if _, err := fmt.Fprint(w, "ok"); err != nil && log != nil {
				log.Error(err)
			}
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
