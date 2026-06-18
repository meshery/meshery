package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/meshery/meshery/server/meshes"
	"github.com/meshery/meshery/server/models"
)

// SystemStatusHandler returns a structured JSON report of every dependency
// the Meshery Server relies on: database, NATS broker, Kubernetes
// connectivity, MeshSync health, remote provider reachability, and adapter
// status.
//
//	GET /api/system/status
//
// Response shape:
//
//	{
//	  "status": "healthy" | "degraded" | "unhealthy",
//	  "healthy": true,
//	  "dependencies": {
//	    "database":      { "status": "healthy", "error": "" },
//	    "natsBroker":    { "status": "healthy", "error": "" },
//	    "kubernetes":    { "status": "healthy", "clusterCount": 2, "version": "v1.29.0", "error": "" },
//	    "meshsync":      { "status": "healthy", "error": "" },
//	    "remoteProvider":{ "status": "healthy", "error": "" },
//	    "adapters":      [ { "name": "istio", "status": "healthy", "location": "...", "error": "" } ]
//	  }
//	}
func (h *Handler) SystemStatusHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	ctx := r.Context()
	status := systemStatus{
		Dependencies: systemDependencies{
			Database:       checkDatabase(h),
			NATSBroker:     checkNATS(h),
			Kubernetes:     checkKubernetes(h, provider),
			MeshSync:       checkMeshSync(h),
			RemoteProvider: checkRemoteProvider(h),
			Adapters:       checkAdapters(ctx, h),
		},
	}
	status.Status = computeOverallStatus(status.Dependencies)
	status.Healthy = status.Status == "healthy"

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if status.Status == "unhealthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	_ = json.NewEncoder(w).Encode(status)
}

// ---------------------------------------------------------------------------
// Response types
// ---------------------------------------------------------------------------

type systemStatus struct {
	Status       string             `json:"status"`
	Healthy      bool               `json:"healthy"`
	Dependencies systemDependencies `json:"dependencies"`
}

type systemDependencies struct {
	Database       dependencyStatus `json:"database"`
	NATSBroker     dependencyStatus `json:"natsBroker"`
	Kubernetes     k8sStatus        `json:"kubernetes"`
	MeshSync       dependencyStatus `json:"meshsync"`
	RemoteProvider dependencyStatus `json:"remoteProvider"`
	Adapters       []adapterStatus  `json:"adapters"`
}

type dependencyStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type k8sStatus struct {
	Status       string `json:"status"`
	Error        string `json:"error,omitempty"`
	ClusterCount int    `json:"clusterCount,omitempty"`
	Version      string `json:"version,omitempty"`
}

type adapterStatus struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Location string `json:"location,omitempty"`
	Error    string `json:"error,omitempty"`
}

// ---------------------------------------------------------------------------
// Individual check functions
// ---------------------------------------------------------------------------

func checkDatabase(h *Handler) dependencyStatus {
	if h.dbHandler == nil || h.dbHandler.DB == nil {
		return dependencyStatus{Status: "unhealthy", Error: "database handler is nil"}
	}
	sqlDB, err := h.dbHandler.DB.DB()
	if err != nil {
		return dependencyStatus{Status: "unhealthy", Error: err.Error()}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return dependencyStatus{Status: "unhealthy", Error: "database ping timed out after 3s"}
		}
		return dependencyStatus{Status: "unhealthy", Error: err.Error()}
	}
	return dependencyStatus{Status: "healthy"}
}

func checkNATS(h *Handler) dependencyStatus {
	if h.brokerConn == nil || h.brokerConn.IsEmpty() {
		return dependencyStatus{Status: "unhealthy", Error: "NATS broker not connected"}
	}
	return dependencyStatus{Status: "healthy"}
}

func checkKubernetes(h *Handler, provider models.Provider) k8sStatus {
	// Attempt a real K8s API ping via the provider's global KubeClient.
	kubeClient := provider.GetKubeClient()
	if kubeClient != nil && kubeClient.KubeClient != nil {
		version, err := kubeClient.KubeClient.ServerVersion()
		if err != nil {
			return k8sStatus{
				Status: "unhealthy",
				Error:  fmt.Sprintf("Kubernetes API unreachable: %v", err),
			}
		}
		return k8sStatus{Status: "healthy", ClusterCount: 1, Version: version.GitVersion}
	}

	// Fallback: check the kubeconfig folder for configuration files.
	if h.config == nil {
		return k8sStatus{Status: "unknown", Error: "handler config not initialized"}
	}
	if h.config.KubeConfigFolder == "" {
		return k8sStatus{Status: "degraded", Error: "kubeconfig folder not configured"}
	}
	info, err := os.Stat(h.config.KubeConfigFolder)
	if err != nil {
		if os.IsNotExist(err) {
			return k8sStatus{Status: "unhealthy", Error: "kubeconfig folder does not exist"}
		}
		return k8sStatus{Status: "unhealthy", Error: err.Error()}
	}
	if !info.IsDir() {
		return k8sStatus{Status: "unhealthy", Error: "kubeconfig path is not a directory"}
	}
	entries, err := os.ReadDir(h.config.KubeConfigFolder)
	if err != nil {
		return k8sStatus{Status: "degraded", Error: err.Error()}
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}
	return k8sStatus{Status: "degraded", Error: "kubeconfig exists but no active KubeClient", ClusterCount: count}
}

func checkMeshSync(h *Handler) dependencyStatus {
	if h.MeshsyncChannel == nil {
		return dependencyStatus{Status: "unknown", Error: "meshsync channel not configured"}
	}
	// Non-destructive peek: len() returns the number of buffered elements
	// without consuming them. If the channel is unbuffered, len is always 0
	// and we report "healthy" (channel exists, no data to consume).
	if len(h.MeshsyncChannel) > 0 {
		return dependencyStatus{Status: "healthy"}
	}
	return dependencyStatus{Status: "degraded", Error: "meshsync channel empty (no recent sync)"}
}

func checkRemoteProvider(h *Handler) dependencyStatus {
	if len(h.config.Providers) == 0 {
		return dependencyStatus{Status: "healthy", Error: "no providers configured"}
	}
	for name, provider := range h.config.Providers {
		if provider.GetProviderType() == models.RemoteProviderType {
			props := provider.GetProviderProperties()
			if len(props.Capabilities) == 0 {
				return dependencyStatus{Status: "degraded", Error: name + ": provider capabilities not yet loaded"}
			}
			return dependencyStatus{Status: "healthy"}
		}
	}
	return dependencyStatus{Status: "healthy", Error: "no remote provider configured"}
}

func checkAdapters(ctx context.Context, h *Handler) []adapterStatus {
	tracker := h.config.AdapterTracker
	if tracker == nil {
		return []adapterStatus{}
	}
	adapters := tracker.GetAdapters(ctx)
	if len(adapters) == 0 {
		return []adapterStatus{}
	}
	results := make([]adapterStatus, 0, len(adapters))
	for _, a := range adapters {
		s := adapterStatus{
			Name:     a.Name,
			Location: a.Location,
			Status:   "unknown",
		}
		if a.Location != "" {
			probeCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			mClient, err := meshes.CreateClient(probeCtx, a.Location)
			if err != nil {
				s.Status = "unhealthy"
				s.Error = fmt.Sprintf("gRPC dial failed: %v", err)
			} else {
				_, err := mClient.MClient.MeshName(probeCtx, &meshes.MeshNameRequest{})
				if err != nil {
					s.Status = "unhealthy"
					s.Error = fmt.Sprintf("MeshName RPC failed: %v", err)
				} else {
					s.Status = "healthy"
				}
				_ = mClient.Close()
			}
			cancel()
		}
		results = append(results, s)
	}
	return results
}

func computeOverallStatus(deps systemDependencies) string {
	unhealthy := 0
	degraded := 0

	check := func(s string) {
		switch s {
		case "unhealthy":
			unhealthy++
		case "degraded":
			degraded++
		case "unknown":
			degraded++ // treat unknown as degraded
		}
	}

	check(deps.Database.Status)
	check(deps.NATSBroker.Status)
	check(deps.Kubernetes.Status)
	check(deps.MeshSync.Status)
	check(deps.RemoteProvider.Status)
	for _, a := range deps.Adapters {
		check(a.Status)
	}

	if unhealthy > 0 {
		return "unhealthy"
	}
	if degraded > 0 {
		return "degraded"
	}
	return "healthy"
}
