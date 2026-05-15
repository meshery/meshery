// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/utils"
	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
)

// GetConnectionOperatorStatusHandler returns the latest status of the Meshery
// operator for the supplied Kubernetes connection. It replaces the
// `getOperatorStatus` GraphQL query.
func (h *Handler) GetConnectionOperatorStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	connectionID := mux.Vars(req)["connectionID"]
	if connectionID == "" {
		h.log.Error(ErrEmptyConnectionID())
		writeMeshkitError(w, ErrEmptyConnectionID(), http.StatusBadRequest)
		return
	}

	unknown := &model.MesheryControllersStatusListItem{
		ConnectionID: connectionID,
		Status:       model.MesheryControllerStatusUnkown,
		Controller:   model.GetInternalController(models.MesheryOperator),
	}

	connUUID := uuid.FromStringOrNil(connectionID)
	inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(connUUID)
	if !ok || inst == nil {
		writeOperatorStatusListItem(w, unknown)
		return
	}

	machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
	if err != nil {
		h.log.Error(model.ErrMesheryControllersStatusSubscription(err))
		writeOperatorStatusListItem(w, unknown)
		return
	}

	ctrlHandlers := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
	op, ok := ctrlHandlers[models.MesheryOperator]
	if !ok || op == nil {
		writeOperatorStatusListItem(w, unknown)
		return
	}

	resp := &model.MesheryControllersStatusListItem{
		ConnectionID: connectionID,
		Status:       model.GetInternalControllerStatus(op.GetStatus()),
		Controller:   model.GetInternalController(models.MesheryOperator),
	}
	writeOperatorStatusListItem(w, resp)
}

func writeOperatorStatusListItem(w http.ResponseWriter, item *model.MesheryControllersStatusListItem) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(item)
}

// GetConnectionMeshsyncStatusHandler returns the Meshsync controller status for
// the supplied Kubernetes connection. It replaces the `getMeshsyncStatus`
// GraphQL query.
func (h *Handler) GetConnectionMeshsyncStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	connectionID := mux.Vars(req)["connectionID"]
	if connectionID == "" {
		h.log.Error(ErrEmptyConnectionID())
		writeMeshkitError(w, ErrEmptyConnectionID(), http.StatusBadRequest)
		return
	}

	unknown := &model.OperatorControllerStatus{
		ConnectionID: connectionID,
		Status:       model.Status(model.MesheryControllerStatusUnkown.String()),
		Name:         model.GetInternalController(models.Meshsync).String(),
	}

	connUUID := uuid.FromStringOrNil(connectionID)
	inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(connUUID)
	if !ok || inst == nil {
		writeOperatorControllerStatus(w, unknown)
		return
	}

	machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
	if err != nil {
		h.log.Error(model.ErrMesheryControllersStatusSubscription(err))
		writeOperatorControllerStatus(w, unknown)
		return
	}

	handlers := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
	status := model.GetMeshSyncInfo(
		handlers[models.Meshsync],
		handlers[models.MesheryBroker],
		h.log,
	)
	status.ConnectionID = connectionID
	writeOperatorControllerStatus(w, &status)
}

// GetConnectionNatsStatusHandler returns the Meshery-broker (NATS) controller
// status for the supplied Kubernetes connection. It replaces the
// `getNatsStatus` GraphQL query.
func (h *Handler) GetConnectionNatsStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	connectionID := mux.Vars(req)["connectionID"]
	if connectionID == "" {
		h.log.Error(ErrEmptyConnectionID())
		writeMeshkitError(w, ErrEmptyConnectionID(), http.StatusBadRequest)
		return
	}

	unknown := &model.OperatorControllerStatus{
		ConnectionID: connectionID,
		Status:       model.Status(model.MesheryControllerStatusUnkown.String()),
		Name:         model.GetInternalController(models.MesheryBroker).String(),
	}

	connUUID := uuid.FromStringOrNil(connectionID)
	if connUUID == uuid.Nil {
		writeOperatorControllerStatus(w, unknown)
		return
	}

	inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(connUUID)
	if !ok || inst == nil {
		writeOperatorControllerStatus(w, unknown)
		return
	}

	machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
	if err != nil {
		h.log.Error(model.ErrMesheryControllersStatusSubscription(err))
		writeOperatorControllerStatus(w, unknown)
		return
	}

	status := model.GetBrokerInfo(
		machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()[models.MesheryBroker],
		h.log,
	)
	status.ConnectionID = connectionID
	writeOperatorControllerStatus(w, &status)
}

func writeOperatorControllerStatus(w http.ResponseWriter, s *model.OperatorControllerStatus) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s)
}

// operatorChangeRequest is the body for POST /api/system/kubernetes/contexts/{contextID}/operator.
type operatorChangeRequest struct {
	TargetStatus model.Status `json:"targetStatus"`
}

// ChangeOperatorStatusHandler installs or uninstalls the Meshery operator for
// the supplied Kubernetes context. It replaces the `changeOperatorStatus`
// GraphQL mutation.
//
// CRITICAL: this handler intentionally preserves the side effects of the
// resolver — after a successful deploy it calls model.SubscribeToBroker and
// records the broker endpoint into the shared model.ConnectionTrackerSingleton.
// The deploy/undeploy itself runs in a goroutine so the HTTP request can return
// PROCESSING promptly (the resolver does the same).
func (h *Handler) ChangeOperatorStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	ctxID := mux.Vars(req)["contextID"]
	if ctxID == "" {
		writeJSONError(w, "missing contextID", http.StatusBadRequest)
		return
	}

	defer func() {
		_ = req.Body.Close()
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	var input operatorChangeRequest
	if err := json.Unmarshal(body, &input); err != nil {
		h.log.Error(models.ErrUnmarshal(err, "operator status input"))
		writeMeshkitError(w, models.ErrUnmarshal(err, "operator status input"), http.StatusBadRequest)
		return
	}

	deleteOperator := true
	switch input.TargetStatus {
	case model.StatusEnabled:
		h.log.Info("Installing Operator")
		deleteOperator = false
	case model.StatusDisabled:
		h.log.Info("Uninstalling Operator in context ", ctxID)
	default:
		writeJSONError(w, fmt.Sprintf("unsupported targetStatus %q", input.TargetStatus), http.StatusBadRequest)
		return
	}

	// Resolve the kubeconfig for the target context. Mirrors operator.go:56-85,
	// but uses the same AllKubeClusterKey set by KubernetesMiddleware. We try
	// the in-context list first and fall back to provider.GetK8sContext (used
	// by other handlers like KubernetesPingHandler).
	var (
		k8scontext models.K8sContext
		kubeclient *mesherykube.Client
		found      bool
	)

	if allContexts, ok := req.Context().Value(models.AllKubeClusterKey).([]*models.K8sContext); ok {
		for _, c := range allContexts {
			if c != nil && c.ID == ctxID {
				k8scontext = *c
				found = true
				break
			}
		}
	}
	if !found {
		token, ok := req.Context().Value(models.TokenCtxKey).(string)
		if !ok {
			writeJSONError(w, "missing auth token", http.StatusUnauthorized)
			return
		}
		k8scontextPtr, gerr := provider.GetK8sContext(token, ctxID)
		if gerr != nil {
			writeMeshkitError(w, ErrInvalidKubeContext(gerr, ctxID), http.StatusNotFound)
			return
		}
		k8scontext = k8scontextPtr
	}

	kubeclient, err = k8scontext.GenerateKubeHandler()
	if err != nil {
		writeMeshkitError(w, model.ErrMesheryClient(err), http.StatusBadRequest)
		return
	}
	if kubeclient.KubeClient == nil {
		writeJSONError(w, "meshery client not initialised", http.StatusBadRequest)
		return
	}

	// Locate the operator controller for this context via the state machine
	// keyed off the connection ID (NOT contextID — they may diverge). The
	// legacy resolver pulled this map from ctx.Value(MesheryControllerHandlersKey),
	// but that key is never populated anywhere — so the resolver's `op[ctxID]`
	// indirection silently nil-derefs in practice. Using the state-machine
	// tracker is both the documented pattern (see getOperatorStatus) and the
	// only path that actually works at runtime.
	connectionID := k8scontext.ConnectionID
	if connectionID == "" {
		writeJSONError(w, fmt.Sprintf("k8s context %s has no connection ID", ctxID), http.StatusBadRequest)
		return
	}
	inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(uuid.FromStringOrNil(connectionID))
	if !ok || inst == nil {
		writeJSONError(w, fmt.Sprintf("no state machine instance for connection %s", connectionID), http.StatusNotFound)
		return
	}
	machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
	if err != nil {
		h.log.Error(model.ErrMesheryControllersStatusSubscription(err))
		writeMeshkitError(w, model.ErrMesheryControllersStatusSubscription(err), http.StatusInternalServerError)
		return
	}
	ctrlHandlers := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
	opHandler, ok := ctrlHandlers[models.MesheryOperator]
	if !ok || opHandler == nil {
		writeJSONError(w, "operator controller not initialised for this context", http.StatusInternalServerError)
		return
	}

	// Return PROCESSING and continue the actual deploy in a goroutine — matches
	// the resolver's return-then-deploy contract.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]model.Status{"status": model.StatusProcessing})

	go h.runOperatorChange(
		ctxID,
		k8scontext,
		kubeclient,
		opHandler,
		deleteOperator,
		provider,
	)
}

func (h *Handler) runOperatorChange(
	ctxID string,
	k8scontext models.K8sContext,
	kubeclient *mesherykube.Client,
	opHandler controllers.IMesheryController,
	deleteOperator bool,
	provider models.Provider,
) {
	if h.config.OperatorTracker != nil && h.config.OperatorTracker.DisableOperator {
		h.log.Info("skipping operator deployment (in disabled mode)")
		return
	}

	var err error
	if deleteOperator {
		err = opHandler.Undeploy()
	} else {
		err = opHandler.Deploy(true)
	}
	if err != nil {
		h.log.Error(err)
		return
	}

	if h.config.OperatorTracker != nil {
		if deleteOperator {
			h.config.OperatorTracker.Undeployed(ctxID, true)
		} else {
			h.config.OperatorTracker.Undeployed(ctxID, false)
		}
	}

	h.log.Info("Operator operation executed")

	if deleteOperator {
		return
	}

	// brokerChannel mirrors the resolver's r.brokerChannel. The REST path
	// does not consume the messages, but NATS will keep publishing onto the
	// underlying subscription, so we MUST drain the channel — without a
	// reader, once the 32-element buffer fills the publishing goroutine
	// inside meshkit/broker.SubscribeWithChannel will block, stalling the
	// broker and leaking the goroutine. Spawn a discard reader that runs
	// until the channel is closed (i.e. for the lifetime of this process,
	// matching the legacy resolver's lifetime).
	brokerChannel := make(chan *broker.Message, 32)
	go func() {
		for range brokerChannel {
			// Discard: REST/SSE callers don't consume broker events.
		}
	}()
	endpoint, berr := model.SubscribeToBroker(
		provider,
		kubeclient,
		brokerChannel,
		h.brokerConn,
		model.ConnectionTrackerSingleton,
	)
	h.log.Debug("Endpoint: ", endpoint)
	if berr != nil {
		h.log.Error(berr)
		return
	}
	model.ConnectionTrackerSingleton.Set(k8scontext.ID, endpoint)
	h.log.Info("Connected to broker at: ", endpoint)
	model.ConnectionTrackerSingleton.Log(h.log)
}
