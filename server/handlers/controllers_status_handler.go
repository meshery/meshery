package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/utils"
	system "github.com/meshery/schemas/models/v1beta1/system"
)

// Controller-status SSE + REST handlers.
//
// These replace the former GraphQL surface for controller status:
//   - subscribeMesheryControllersStatus  -> SubscribeMesheryControllersStatusHandler (SSE)
//   - getOperatorStatus                  -> OperatorStatusHandler   (REST)
//   - getMeshsyncStatus                  -> MeshsyncStatusHandler   (REST)
//   - getNatsStatus                      -> BrokerStatusHandler     (REST)
//
// The data source is unchanged: controller status is read on-demand from the
// meshkit controller handlers reached through the per-connection FSM
// (ConnectionToStateMachineInstanceTracker -> kubernetes.MachineCtx ->
// MesheryCtrlsHelper). Missing / not-yet-ready instances degrade to an
// "unknown" status, matching the old resolver behavior.
//
// Wire payloads use canonical camelCase JSON tags (connectionId, not
// connectionID) per the identifier naming guide. The status/controller string
// values are reproduced here to match exactly what the old GraphQL enums put on
// the wire, so the frontend contract is unchanged.

const (
	// controllersStatusKeepAliveInterval bounds how long the SSE connection can
	// sit idle before we write a comment line, so the browser and any proxy in
	// front of Meshery notice a dead peer and don't time the stream out.
	controllersStatusKeepAliveInterval = 15 * time.Second
	// controllersStatusPollInterval is how often the server re-reads controller
	// status behind the single SSE connection. Polling is intentional for this
	// phase; a future phase replaces it with an event-driven source (see the
	// connections-sse migration plan).
	controllersStatusPollInterval = 5 * time.Second
	// controllersStatusWriteTimeout caps how long a single frame write may block
	// on the socket, so a client that stops reading can't wedge this goroutine.
	controllersStatusWriteTimeout = 30 * time.Second
)

// ControllerStatusItem is one controller's status for one connection. It is the
// element type of both the SSE snapshot array and the operator REST response.
type ControllerStatusItem struct {
	ConnectionID string `json:"connectionId"`
	Controller   string `json:"controller"` // OPERATOR | MESHSYNC | BROKER
	Status       string `json:"status"`
	Version      string `json:"version"`
}

// ControllerInfo is the MeshSync / Broker one-shot status payload. Its status
// string may be composed (e.g. "Connected <endpoint>"), mirroring the old
// getMeshsyncStatus / getNatsStatus queries the UI branches on.
type ControllerInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Status       string `json:"status"`
	ConnectionID string `json:"connectionId"`
}

// controllersStatusUnknown is the status string emitted when a connection has no
// ready FSM instance or its context can't be read.
const controllersStatusUnknown = "UNKOWN"

// internalControllerName reproduces model.GetInternalController: the wire name
// for a controller enum.
func internalControllerName(c models.MesheryController) string {
	switch c {
	case models.MesheryBroker:
		return "BROKER"
	case models.MesheryOperator:
		return "OPERATOR"
	case models.Meshsync:
		return "MESHSYNC"
	}
	return ""
}

// internalControllerStatus reproduces model.GetInternalControllerStatus: the
// wire status string for a meshkit controller status.
func internalControllerStatus(status controllers.MesheryControllerStatus) string {
	switch status {
	case controllers.Deployed:
		return "DEPLOYED"
	case controllers.NotDeployed:
		return "NOTDEPLOYED"
	case controllers.Deploying:
		return "DEPLOYING"
	case controllers.Unknown:
		return controllersStatusUnknown
	case controllers.Undeployed:
		return "UNDEPLOYED"
	case controllers.Enabled:
		return "ENABLED"
	case controllers.Running:
		return "RUNNING"
	case controllers.Connected:
		return "CONNECTED"
	}
	return ""
}

// machineCtxForConnection resolves the kubernetes FSM machine context for a
// connection. Returns (nil, false) when the connection has no ready instance.
func (h *Handler) machineCtxForConnection(connectionID string) (*kubernetes.MachineCtx, bool) {
	connUUID := uuid.FromStringOrNil(connectionID)
	if connUUID == uuid.Nil || h.ConnectionToStateMachineInstanceTracker == nil {
		return nil, false
	}
	inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(connUUID)
	if !ok || inst == nil {
		return nil, false
	}
	machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
	if err != nil || machinectx == nil || machinectx.MesheryCtrlsHelper == nil {
		if err != nil {
			h.log.Error(err)
		}
		return nil, false
	}
	return machinectx, true
}

// controllerHandlersForConnection resolves the per-controller meshkit handlers
// for a connection through its FSM instance. Returns (nil, false) when the
// connection has no ready instance — callers degrade to an unknown status.
func (h *Handler) controllerHandlersForConnection(connectionID string) (map[models.MesheryController]controllers.IMesheryController, bool) {
	machinectx, ok := h.machineCtxForConnection(connectionID)
	if !ok {
		return nil, false
	}
	return machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext(), true
}

// mesheryHoldsLiveBrokerConnection reports whether Meshery currently holds a
// live broker connection (via the MeshSync data handler) for this connection —
// i.e. it is actually receiving MeshSync data. This is the authoritative
// "connected" signal for the MeshSync/broker controllers, as opposed to
// meshkit's status which re-probes the broker's monitoring endpoint from Meshery
// and false-negatives when that endpoint isn't reachable even though the data
// path is up (Docker Desktop, cluster-internal broker, etc.).
func (h *Handler) mesheryHoldsLiveBrokerConnection(machinectx *kubernetes.MachineCtx) bool {
	if machinectx == nil || machinectx.MesheryCtrlsHelper == nil {
		return false
	}
	dataHandler := machinectx.MesheryCtrlsHelper.GetMeshSyncDataHandlersForEachContext()
	return dataHandler != nil && dataHandler.IsConnected()
}

// deriveControllerStatus upgrades a MeshSync/broker controller's status to
// CONNECTED when Meshery holds a live broker connection. We only upgrade a
// controller that is already present (running/deployed/enabled) — a
// not-deployed controller is never reported as connected — and the operator's
// status is left untouched (it is unrelated to broker connectivity).
func deriveControllerStatus(controller models.MesheryController, status string, brokerConnected bool) string {
	if !brokerConnected {
		return status
	}
	if controller != models.Meshsync && controller != models.MesheryBroker {
		return status
	}
	switch status {
	case "RUNNING", "DEPLOYED", "ENABLED":
		return "CONNECTED"
	}
	return status
}

// collectControllersStatus builds the full status list for the requested
// connections. The result is sorted (connectionId, controller) so callers can
// compare successive snapshots byte-for-byte to detect changes.
func (h *Handler) collectControllersStatus(connectionIDs []string) []ControllerStatusItem {
	items := make([]ControllerStatusItem, 0)
	for _, connectionID := range connectionIDs {
		machinectx, ok := h.machineCtxForConnection(connectionID)
		if !ok {
			continue
		}
		ctrlHandlers := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
		brokerConnected := h.mesheryHoldsLiveBrokerConnection(machinectx)
		for controller, ctrlHandler := range ctrlHandlers {
			version, err := ctrlHandler.GetVersion()
			if err != nil {
				h.log.Debugf("controllers status: version for %s on %s: %v", internalControllerName(controller), connectionID, err)
			}
			status := deriveControllerStatus(controller, internalControllerStatus(ctrlHandler.GetStatus()), brokerConnected)
			items = append(items, ControllerStatusItem{
				ConnectionID: connectionID,
				Controller:   internalControllerName(controller),
				Status:       status,
				Version:      version,
			})
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].ConnectionID != items[j].ConnectionID {
			return items[i].ConnectionID < items[j].ConnectionID
		}
		return items[i].Controller < items[j].Controller
	})
	return items
}

// SubscribeMesheryControllersStatusHandler streams controller status (operator,
// MeshSync, broker) for the requested connections over Server-Sent Events. It
// replaces the subscribeMesheryControllersStatus GraphQL subscription.
//
// Connections are passed as repeatable camelCase query params:
// ?connectionIds=<id>&connectionIds=<id>. The handler emits the full status
// snapshot once immediately, then re-polls every controllersStatusPollInterval
// and re-emits the full snapshot only when it changed. Sending the full list
// (not per-controller deltas) keeps the client idempotent — it just replaces
// its controller state — removing the fragile client-side merge the old Relay
// path needed.
//
// Each frame is the JSON array framed as an unnamed SSE event
// (data: <json>\n\n) so the browser's EventSource.onmessage receives it. The
// subscription lives for the duration of the request.
func (h *Handler) SubscribeMesheryControllersStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, _ models.Provider) {
	if user == nil {
		writeJSONError(w, "user unauthorized", http.StatusUnauthorized)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSONError(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	connectionIDs := req.URL.Query()["connectionIds"]

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// X-Accel-Buffering disables buffering at any nginx hop in front of Meshery
	// so events reach the browser immediately.
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	responseController := http.NewResponseController(w)

	writeFrame := func(payload []byte) bool {
		_ = responseController.SetWriteDeadline(time.Now().Add(controllersStatusWriteTimeout))
		if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
			return false
		}
		flusher.Flush()
		return true
	}

	// Initial snapshot so the UI renders current state without waiting a tick.
	last, err := json.Marshal(h.collectControllersStatus(connectionIDs))
	if err != nil {
		h.log.Error(models.ErrMarshal(err, "controllers status"))
		return
	}
	if !writeFrame(last) {
		return
	}

	poll := time.NewTicker(controllersStatusPollInterval)
	defer poll.Stop()
	keepAlive := time.NewTicker(controllersStatusKeepAliveInterval)
	defer keepAlive.Stop()

	ctx := req.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case <-poll.C:
			payload, err := json.Marshal(h.collectControllersStatus(connectionIDs))
			if err != nil {
				h.log.Error(models.ErrMarshal(err, "controllers status"))
				continue
			}
			// Snapshots are sorted, so equal state marshals to equal bytes —
			// suppress no-op frames.
			if string(payload) == string(last) {
				continue
			}
			if !writeFrame(payload) {
				return
			}
			last = payload
		case <-keepAlive.C:
			_ = responseController.SetWriteDeadline(time.Now().Add(controllersStatusWriteTimeout))
			if _, err := fmt.Fprint(w, ": keepalive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

// OperatorStatusHandler returns the operator's current status for a connection.
// Replaces the getOperatorStatus GraphQL query.
// GET /api/system/controllers/operator/status?connectionId=<id>
func (h *Handler) OperatorStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	connectionID := req.URL.Query().Get("connectionId")
	item := ControllerStatusItem{
		ConnectionID: connectionID,
		Controller:   internalControllerName(models.MesheryOperator),
		Status:       controllersStatusUnknown,
	}
	if ctrlHandlers, ok := h.controllerHandlersForConnection(connectionID); ok {
		if operator, ok := ctrlHandlers[models.MesheryOperator]; ok {
			item.Status = internalControllerStatus(operator.GetStatus())
			item.Version, _ = operator.GetVersion()
		}
	}
	writeJSONMessage(w, item, http.StatusOK)
}

// MeshsyncStatusHandler returns MeshSync's current status for a connection.
// Replaces the getMeshsyncStatus GraphQL query.
// GET /api/system/controllers/meshsync/status?connectionId=<id>
func (h *Handler) MeshsyncStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	connectionID := req.URL.Query().Get("connectionId")
	info := ControllerInfo{
		ConnectionID: connectionID,
		Name:         "MeshSync",
		Status:       controllersStatusUnknown,
	}
	if machinectx, ok := h.machineCtxForConnection(connectionID); ok {
		ctrlHandlers := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
		info = h.meshsyncInfo(ctrlHandlers[models.Meshsync], ctrlHandlers[models.MesheryBroker], h.mesheryHoldsLiveBrokerConnection(machinectx))
		info.ConnectionID = connectionID
	}
	writeJSONMessage(w, info, http.StatusOK)
}

// BrokerStatusHandler returns the Meshery Broker (NATS) status for a connection.
// Replaces the getNatsStatus GraphQL query.
// GET /api/system/controllers/broker/status?connectionId=<id>
func (h *Handler) BrokerStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	connectionID := req.URL.Query().Get("connectionId")
	info := ControllerInfo{
		ConnectionID: connectionID,
		Name:         "MesheryBroker",
		Status:       controllersStatusUnknown,
	}
	if machinectx, ok := h.machineCtxForConnection(connectionID); ok {
		ctrlHandlers := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
		info = h.brokerInfo(ctrlHandlers[models.MesheryBroker], h.mesheryHoldsLiveBrokerConnection(machinectx))
		info.ConnectionID = connectionID
	}
	writeJSONMessage(w, info, http.StatusOK)
}

// isBrokerReachableStatus reports whether a meshkit status string means the
// controller is present (deployed/running) — the states we upgrade to Connected
// when Meshery already holds a live broker connection.
func isBrokerReachableStatus(status string) bool {
	return status == controllers.Running.String() ||
		status == controllers.Deployed.String() ||
		status == controllers.Enabled.String()
}

// composeConnectedBrokerStatus builds the "Connected <endpoint>" status string,
// appending the broker's public endpoint when it can be resolved.
func composeConnectedBrokerStatus(broker controllers.IMesheryController) string {
	status := controllers.Connected.String()
	if broker != nil {
		if endpoint, err := broker.GetPublicEndpoint(); err == nil && endpoint != "" {
			status = fmt.Sprintf("%s %s", status, endpoint)
		}
	}
	return status
}

// brokerInfo reproduces model.GetBrokerInfo without the gqlgen model dependency
// (handlers cannot import internal/graphql/model — it imports handlers).
// brokerConnected upgrades a present-but-unverified broker to Connected when
// Meshery already holds a live broker connection.
func (h *Handler) brokerInfo(broker controllers.IMesheryController, brokerConnected bool) ControllerInfo {
	if broker == nil {
		return ControllerInfo{Status: controllersStatusUnknown}
	}
	status := broker.GetStatus().String()
	if status == controllers.Connected.String() {
		endpoint, _ := broker.GetPublicEndpoint()
		status = fmt.Sprintf("%s %s", status, endpoint)
	} else if brokerConnected && isBrokerReachableStatus(status) {
		status = composeConnectedBrokerStatus(broker)
	}
	version, _ := broker.GetVersion()
	return ControllerInfo{
		Name:    broker.GetName(),
		Status:  status,
		Version: version,
	}
}

// meshsyncInfo reproduces model.GetMeshSyncInfo without the gqlgen model
// dependency. brokerConnected upgrades a present-but-unverified MeshSync to
// Connected when Meshery already holds a live broker connection.
func (h *Handler) meshsyncInfo(meshsync, broker controllers.IMesheryController, brokerConnected bool) ControllerInfo {
	if meshsync == nil {
		return ControllerInfo{Status: controllersStatusUnknown}
	}
	status := meshsync.GetStatus().String()
	if broker == nil {
		status = controllers.Unknown.String()
	} else if status == controllers.Connected.String() {
		endpoint, err := broker.GetPublicEndpoint()
		if err != nil {
			h.log.Warn(err)
		} else if endpoint == "" {
			h.log.Warnf("broker public endpoint is empty while composing meshsync status")
		} else {
			status = fmt.Sprintf("%s %s", status, endpoint)
		}
	} else if brokerConnected && isBrokerReachableStatus(status) {
		status = composeConnectedBrokerStatus(broker)
	}
	version, _ := meshsync.GetVersion()
	return ControllerInfo{
		Name:    meshsync.GetName(),
		Status:  status,
		Version: version,
	}
}

// ---- Controller diagnostics ----
//
// The diagnostics API turns the raw controller statuses (and Meshery's live
// broker connection) into human-actionable problems with remediation steps, so
// the connection detail view can render a "Diagnostics" section. The wording is
// shared with the connect-time events (models.BrokerUnreachable*), so the same
// guidance appears whether the user reads the notification or the detail view.

// diagnosticControllerPtr maps a Meshery controller to the schemas enum pointer,
// or nil when it doesn't map.
func diagnosticControllerPtr(c models.MesheryController) *system.ControllerDiagnosticController {
	var v system.ControllerDiagnosticController
	switch c {
	case models.MesheryBroker:
		v = system.ControllerDiagnosticControllerBROKER
	case models.Meshsync:
		v = system.ControllerDiagnosticControllerMESHSYNC
	case models.MesheryOperator:
		v = system.ControllerDiagnosticControllerOPERATOR
	default:
		return nil
	}
	return &v
}

func strPtr(s string) *string { return &s }

// computeConnectionDiagnostics derives the diagnostics list for a connection
// from its controller statuses and Meshery's live broker connection. Healthy is
// true when no warning/error diagnostic was raised.
func (h *Handler) computeConnectionDiagnostics(connectionID string) system.ConnectionDiagnostics {
	result := system.ConnectionDiagnostics{
		ConnectionId: uuid.FromStringOrNil(connectionID),
		Healthy:      true,
		Diagnostics:  []system.ControllerDiagnostic{},
	}

	add := func(d system.ControllerDiagnostic) {
		if d.Severity == system.Warning || d.Severity == system.Error {
			result.Healthy = false
		}
		result.Diagnostics = append(result.Diagnostics, d)
	}

	machinectx, ok := h.machineCtxForConnection(connectionID)
	if !ok {
		add(system.ControllerDiagnostic{
			Severity:    system.Info,
			Code:        "connection_inactive",
			Summary:     "Connection is not active",
			Description: strPtr("Meshery has no active session for this connection yet, so its controller status can't be read. Connect the cluster to begin monitoring the operator, MeshSync, and broker."),
		})
		return result
	}

	// Embedded MeshSync runs in-process and does not use the in-cluster operator
	// or broker, so those being undeployed is expected — no diagnostics apply.
	if machinectx.MesheryCtrlsHelper.GetMeshsyncDeploymentMode() == connections.MeshsyncDeploymentModeEmbedded {
		return result
	}

	ctrlHandlers := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
	brokerConnected := h.mesheryHoldsLiveBrokerConnection(machinectx)

	// Surface *why* operator setup failed (kubeconfig unreadable, Kubernetes
	// client creation, or a failed Deploy) — not merely that the operator is
	// absent. When client creation fails the controller handlers map is empty, so
	// the status-based checks below are skipped entirely; this error-level
	// diagnostic is what tells the user what actually went wrong. Mirrors the
	// broker diagnostics: an actionable Error carrying the underlying cause and
	// remediation.
	if opErr := machinectx.MesheryCtrlsHelper.GetOperatorError(); opErr != nil {
		add(system.ControllerDiagnostic{
			Severity:    system.Error,
			Controller:  diagnosticControllerPtr(models.MesheryOperator),
			Code:        "operator_deploy_failed",
			Summary:     "Meshery Operator deployment failed",
			Description: strPtr(fmt.Sprintf("Meshery could not deploy the Meshery Operator for this connection, so MeshSync and the Meshery Broker are unavailable. Underlying cause: %s", opErr.Error())),
			Remediation: &[]string{
				"Verify Meshery can read the kubeconfig for this cluster — it needs read permission on the kubeconfig file.",
				"Ensure Meshery has permission to create resources in the 'meshery' namespace.",
				"Re-connect the cluster to retry operator deployment.",
			},
		})
	}

	// Operator must be deployed for operator-mode MeshSync/broker to exist.
	if operator := ctrlHandlers[models.MesheryOperator]; operator != nil {
		switch internalControllerStatus(operator.GetStatus()) {
		case "NOTDEPLOYED", "UNDEPLOYED":
			add(system.ControllerDiagnostic{
				Severity:    system.Warning,
				Controller:  diagnosticControllerPtr(models.MesheryOperator),
				Code:        "operator_not_deployed",
				Summary:     "Meshery Operator is not deployed",
				Description: strPtr("The Meshery Operator manages MeshSync and the Meshery Broker inside the cluster. Without it, Meshery cannot collect cluster state for this connection."),
				Remediation: &[]string{
					"Deploy the operator by re-connecting the cluster, or switch MeshSync to operator mode from the connection's actions.",
					"Ensure Meshery has permission to create resources in the 'meshery' namespace.",
				},
			})
		}
	}

	// Broker present but Meshery holds no live connection => unreachable. This is
	// also why MeshSync would show "running but not connected".
	if broker := ctrlHandlers[models.MesheryBroker]; broker != nil {
		if !brokerConnected && isBrokerReachableStatus(broker.GetStatus().String()) {
			rem := append([]string(nil), models.BrokerUnreachableRemediation...)
			d := system.ControllerDiagnostic{
				Severity:    system.Warning,
				Controller:  diagnosticControllerPtr(models.MesheryBroker),
				Code:        "broker_unreachable",
				Summary:     "Meshery Broker unreachable",
				Description: strPtr(models.BrokerUnreachableLongDescription),
				Remediation: &rem,
			}
			if endpoint, err := broker.GetPublicEndpoint(); err == nil && endpoint != "" {
				d.Endpoint = strPtr(endpoint)
			}
			add(d)
		}

		// Informational: surface HOW Meshery reaches the broker (managed
		// port-forward vs in-cluster ClusterIP vs a direct endpoint) so the
		// networking is visible in the UI.
		if net, ok := h.brokerNetworkingDiagnostic(machinectx, broker); ok {
			add(net)
		}
	}

	return result
}

// brokerNetworkingDiagnostic describes the transport Meshery uses to reach the
// broker. It returns ok=false when the broker isn't deployed yet (nothing to
// describe).
func (h *Handler) brokerNetworkingDiagnostic(
	machinectx *kubernetes.MachineCtx,
	broker controllers.IMesheryController,
) (system.ControllerDiagnostic, bool) {
	// Only describe networking for a broker that is actually present.
	if !isBrokerReachableStatus(broker.GetStatus().String()) &&
		broker.GetStatus() != controllers.Connected {
		return system.ControllerDiagnostic{}, false
	}

	managedAddr := machinectx.MesheryCtrlsHelper.GetBrokerPortForwardAddr()
	endpoint, _ := broker.GetPublicEndpoint()
	inCluster := os.Getenv("KUBERNETES_SERVICE_HOST") != ""

	var transport, description, ep string
	switch {
	case managedAddr != "":
		transport = "Managed port-forward"
		ep = managedAddr
		description = fmt.Sprintf(
			"Meshery reaches the Meshery Broker (NATS) through a managed port-forward at %s, tunneled to the in-cluster NATS pod through the Kubernetes API server. This is automatic for out-of-cluster Meshery; set MESHERY_MANAGED_BROKER_PORTFORWARD=false on the server to disable it. The transport follows the connection's MeshSync mode — operator mode uses the in-cluster broker; embedded mode uses none.",
			managedAddr,
		)
	case inCluster:
		transport = "In-cluster (ClusterIP)"
		ep = endpoint
		description = "Meshery runs inside the cluster and reaches the Meshery Broker (NATS) directly at its cluster-internal (ClusterIP) address."
	default:
		transport = "Direct endpoint"
		ep = endpoint
		description = "Meshery reaches the Meshery Broker (NATS) at the resolved endpoint (e.g. a NodePort/LoadBalancer or a manual port-forward)."
	}

	d := system.ControllerDiagnostic{
		Severity:    system.Info,
		Controller:  diagnosticControllerPtr(models.MesheryBroker),
		Code:        "broker_networking",
		Summary:     fmt.Sprintf("Broker networking: %s", transport),
		Description: strPtr(description),
	}
	if ep != "" {
		d.Endpoint = strPtr(ep)
	}
	return d, true
}

// ControllerDiagnosticsHandler returns human-actionable diagnostics and
// remediation for a kubernetes connection's controllers. It powers the
// "Diagnostics" section of the connection detail view.
// GET /api/system/controllers/diagnostics?connectionId=<id>
func (h *Handler) ControllerDiagnosticsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	connectionID := req.URL.Query().Get("connectionId")
	writeJSONMessage(w, h.computeConnectionDiagnostics(connectionID), http.StatusOK)
}
