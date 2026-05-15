package handlers

// kubernetes_context_handler.go — per-context REST endpoints that replace the
// GraphQL operations removed in the GraphQL→SSE migration.
//
// Handler methods:
//   GetKubernetesNamespacesHandler   GET  /api/system/kubernetes/contexts/{contextID}/namespaces
//   GetMesheryOperatorStatusHandler  GET  /api/system/kubernetes/contexts/{contextID}/operator/status
//   GetMeshsyncStatusHandler         GET  /api/system/kubernetes/contexts/{contextID}/meshsync/status
//   GetNatsStatusHandler             GET  /api/system/kubernetes/contexts/{contextID}/nats/status
//   ResyncClusterHandler             POST /api/system/kubernetes/contexts/{contextID}/resync
//
// TODO(schemas-canonical): once meshery/schemas publishes OpenAPI specs for
//   these five paths (see docs/openapi/kubernetes-context-endpoints.yaml) and
//   generates Go types, replace the locally-defined response structs below
//   with imports from the generated package.
//   Tracked in: meshery/meshery#19424

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/utils"
)

// ──────────────────────────────────────────────────────────────────────────────
// Response types
//
// TODO(schemas-canonical): defined locally to avoid an import cycle with
//   server/internal/graphql/model (which itself imports server/handlers).
//   Replace with generated types from github.com/meshery/schemas once available.
// ──────────────────────────────────────────────────────────────────────────────

// NamespacesResponse is returned by GetKubernetesNamespacesHandler.
type NamespacesResponse struct {
	Namespaces []string `json:"namespaces"`
}

// ControllerStatusResponse is returned by GetMesheryOperatorStatusHandler and
// is the per-item shape pushed over the controller-status SSE stream.
type ControllerStatusResponse struct {
	ConnectionID string `json:"connectionID"`
	Controller   string `json:"controller"`
	Status       string `json:"status"`
	Version      string `json:"version"`
}

// OperatorControllerStatusResponse is returned by GetMeshsyncStatusHandler
// and GetNatsStatusHandler.
type OperatorControllerStatusResponse struct {
	ConnectionID string `json:"connectionID"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Status       string `json:"status"`
}

// ResyncClusterRequest is the JSON body accepted by ResyncClusterHandler.
type ResyncClusterRequest struct {
	ClearDb   bool `json:"clearDb"`
	HardReset bool `json:"hardReset"`
	ReSync    bool `json:"reSync"`
}

// ──────────────────────────────────────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────────────────────────────────────

// fetchNamespacesForCluster queries the MeshSync DB for all Namespace
// resources belonging to the given cluster IDs.
// Mirrors graphql/model.SelectivelyFetchNamespaces — kept here to avoid an
// import cycle with that package.
func fetchNamespacesForCluster(provider models.Provider, clusterIDs []string) ([]string, error) {
	const query = `SELECT DISTINCT rom.name as name
		FROM kubernetes_resources kr
		LEFT JOIN kubernetes_resource_object_meta rom ON kr.id = rom.id
		WHERE kr.kind = 'Namespace' AND kr.cluster_id IN ?`

	var rows *sql.Rows
	var err error
	rows, err = provider.GetGenericPersister().Raw(query, clusterIDs).Rows()
	if err != nil {
		return nil, fmt.Errorf("fetchNamespacesForCluster: %w", err)
	}
	defer func() {
		if closeErr := rows.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	namespaces := make([]string, 0)
	for rows.Next() {
		var name string
		if scanErr := rows.Scan(&name); scanErr != nil {
			return nil, fmt.Errorf("fetchNamespacesForCluster scan: %w", scanErr)
		}
		namespaces = append(namespaces, name)
	}
	return namespaces, err
}

// controllerStatusString maps a meshkit MesheryControllerStatus to the string
// form expected by the UI (mirrors graphql/model.MesheryControllerStatus enums).
func controllerStatusString(s controllers.MesheryControllerStatus) string {
	return s.String()
}

// ──────────────────────────────────────────────────────────────────────────────
// Handlers
// ──────────────────────────────────────────────────────────────────────────────

// GetKubernetesNamespacesHandler returns all Kubernetes namespaces that
// MeshSync has observed for the context identified by {contextID}.
//
// swagger:route GET /api/system/kubernetes/contexts/{contextID}/namespaces KubernetesAPI idGetKubernetesNamespaces
// responses:
//
//	200: namespacesResponseWrapper
func (h *Handler) GetKubernetesNamespacesHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	contextID := mux.Vars(req)["contextID"]
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	k8sCtx, err := provider.GetK8sContext(token, contextID)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "context not found", http.StatusNotFound)
		return
	}

	namespaces, err := fetchNamespacesForCluster(provider, []string{k8sCtx.KubernetesServerID.String()})
	if err != nil {
		h.log.Error(err)
		http.Error(w, "failed to fetch namespaces", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(NamespacesResponse{Namespaces: namespaces}); err != nil {
		h.log.Error(models.ErrMarshal(err, "namespaces"))
	}
}

// GetMesheryOperatorStatusHandler returns the operator deployment status for
// the connection backing the context identified by {contextID}.
//
// swagger:route GET /api/system/kubernetes/contexts/{contextID}/operator/status KubernetesAPI idGetMesheryOperatorStatus
// responses:
//
//	200: controllerStatusResponseWrapper
func (h *Handler) GetMesheryOperatorStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	contextID := mux.Vars(req)["contextID"]
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	k8sCtx, err := provider.GetK8sContext(token, contextID)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "context not found", http.StatusNotFound)
		return
	}

	connectionID := k8sCtx.ConnectionID
	unknown := ControllerStatusResponse{
		ConnectionID: connectionID,
		Controller:   "OPERATOR",
		Status:       controllers.Unknown.String(),
	}

	connectionUUID := uuid.FromStringOrNil(connectionID)
	if connectionUUID == uuid.Nil {
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "operator status"))
		}
		return
	}

	inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(connectionUUID)
	if !ok || inst == nil {
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "operator status"))
		}
		return
	}

	machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
	if err != nil {
		h.log.Error(err)
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "operator status"))
		}
		return
	}

	ctrlHandlers := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
	operatorHandler, ok := ctrlHandlers[models.MesheryOperator]
	if !ok || operatorHandler == nil {
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "operator status"))
		}
		return
	}
	version, _ := operatorHandler.GetVersion()
	result := ControllerStatusResponse{
		ConnectionID: connectionID,
		Controller:   "OPERATOR",
		Status:       controllerStatusString(operatorHandler.GetStatus()),
		Version:      version,
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.log.Error(models.ErrMarshal(err, "operator status"))
	}
}

// GetMeshsyncStatusHandler returns the MeshSync controller status for the
// connection backing the context identified by {contextID}.
//
// swagger:route GET /api/system/kubernetes/contexts/{contextID}/meshsync/status KubernetesAPI idGetMeshsyncStatus
// responses:
//
//	200: operatorControllerStatusResponseWrapper
func (h *Handler) GetMeshsyncStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	contextID := mux.Vars(req)["contextID"]
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	k8sCtx, err := provider.GetK8sContext(token, contextID)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "context not found", http.StatusNotFound)
		return
	}

	connectionID := k8sCtx.ConnectionID
	unknown := OperatorControllerStatusResponse{
		ConnectionID: connectionID,
		Status:       controllers.Unknown.String(),
		Name:         "meshsync",
	}

	connectionUUID := uuid.FromStringOrNil(connectionID)
	if connectionUUID == uuid.Nil {
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "meshsync status"))
		}
		return
	}

	inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(connectionUUID)
	if !ok || inst == nil {
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "meshsync status"))
		}
		return
	}

	machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
	if err != nil {
		h.log.Error(err)
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "meshsync status"))
		}
		return
	}

	ctrlHandlers := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
	meshsyncHandler, ok := ctrlHandlers[models.Meshsync]
	if !ok || meshsyncHandler == nil {
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "meshsync status"))
		}
		return
	}
	version, _ := meshsyncHandler.GetVersion()
	result := OperatorControllerStatusResponse{
		ConnectionID: connectionID,
		Name:         meshsyncHandler.GetName(),
		Status:       controllerStatusString(meshsyncHandler.GetStatus()),
		Version:      version,
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.log.Error(models.ErrMarshal(err, "meshsync status"))
	}
}

// GetNatsStatusHandler returns the NATS broker controller status for the
// connection backing the context identified by {contextID}.
//
// swagger:route GET /api/system/kubernetes/contexts/{contextID}/nats/status KubernetesAPI idGetNatsStatus
// responses:
//
//	200: operatorControllerStatusResponseWrapper
func (h *Handler) GetNatsStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	contextID := mux.Vars(req)["contextID"]
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	k8sCtx, err := provider.GetK8sContext(token, contextID)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "context not found", http.StatusNotFound)
		return
	}

	connectionID := k8sCtx.ConnectionID
	unknown := OperatorControllerStatusResponse{
		ConnectionID: connectionID,
		Status:       controllers.Unknown.String(),
		Name:         "meshery-broker",
	}

	connectionUUID := uuid.FromStringOrNil(connectionID)
	if connectionUUID == uuid.Nil {
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "nats status"))
		}
		return
	}

	inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(connectionUUID)
	if !ok || inst == nil {
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "nats status"))
		}
		return
	}

	machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
	if err != nil {
		h.log.Error(err)
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "nats status"))
		}
		return
	}

	ctrlHandlers := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
	brokerHandler, ok := ctrlHandlers[models.MesheryBroker]
	if !ok || brokerHandler == nil {
		if err := json.NewEncoder(w).Encode(unknown); err != nil {
			h.log.Error(models.ErrMarshal(err, "nats status"))
		}
		return
	}
	version, _ := brokerHandler.GetVersion()
	result := OperatorControllerStatusResponse{
		ConnectionID: connectionID,
		Name:         brokerHandler.GetName(),
		Status:       controllerStatusString(brokerHandler.GetStatus()),
		Version:      version,
	}
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.log.Error(models.ErrMarshal(err, "nats status"))
	}
}

// ResyncClusterHandler triggers a cluster resync for the context identified
// by {contextID}.
//
// swagger:route POST /api/system/kubernetes/contexts/{contextID}/resync KubernetesAPI idResyncCluster
// responses:
//
//	200: resyncResponseWrapper
func (h *Handler) ResyncClusterHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	contextID := mux.Vars(req)["contextID"]
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	var body ResyncClusterRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if body.ClearDb {
		k8sCtx, err := provider.GetK8sContext(token, contextID)
		if err != nil {
			h.log.Error(err)
			http.Error(w, "context not found", http.StatusNotFound)
			return
		}
		models.FlushMeshSyncData(req.Context(), k8sCtx, provider, h.config.EventBroadcaster, token, h.SystemID, h.log)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "processing"}); err != nil {
		h.log.Error(models.ErrMarshal(err, "resync response"))
	}
}
