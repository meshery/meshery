package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
)

// GetControllersDefaultConfig handles GET /api/system/controllers/config.
// Returns the server-wide default Meshery Operator / MeshSync / Broker
// configuration. Fields that have never been set are absent, meaning the
// built-in defaults apply.
func (h *Handler) GetControllersDefaultConfig(w http.ResponseWriter, _ *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	cfg, err := models.GetControllersConfigDefaults(h.dbHandler)
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, http.StatusInternalServerError)
		return
	}
	if cfg == nil {
		cfg = &controllersconfig.MesheryControllersConfig{}
	}
	cfg.SchemaVersion = connections.ControllersConfigSchemaVersion

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cfg); err != nil {
		h.log.Error(models.ErrMarshal(err, "controllers config defaults"))
	}
}

// UpdateControllersDefaultConfig handles PUT /api/system/controllers/config.
// Persists the server-wide defaults, then re-applies the effective
// configuration to every connected Kubernetes connection that inherits the
// changed fields. An empty document clears the stored defaults.
func (h *Handler) UpdateControllersDefaultConfig(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	userID := user.ID
	token, _ := req.Context().Value(models.TokenCtxKey).(string)

	cfg, ok := h.readControllersConfigPayload(w, req)
	if !ok {
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(userID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("update")

	if err := models.SaveControllersConfigDefaults(h.dbHandler, cfg); err != nil {
		h.log.Error(err)
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Failed to persist server-wide controllers configuration defaults.").WithMetadata(map[string]interface{}{"error": err}).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(w, err, http.StatusInternalServerError)
		return
	}

	event := eventBuilder.WithSeverity(events.Informational).WithDescription("Server-wide Meshery Operator, MeshSync, and Broker configuration defaults updated. Re-applying to connected clusters.").Build()
	_ = provider.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)

	// Re-apply the resolved configuration to every tracked Kubernetes
	// connection: per-connection overrides still win, so inheriting
	// connections pick up the new defaults and overriding connections are
	// unaffected for the overridden fields.
	go h.reapplyControllersConfigToTrackedConnections(token, userID, provider)

	stored, err := models.GetControllersConfigDefaults(h.dbHandler)
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, http.StatusInternalServerError)
		return
	}
	if stored == nil {
		stored = &controllersconfig.MesheryControllersConfig{}
	}
	stored.SchemaVersion = connections.ControllersConfigSchemaVersion
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stored); err != nil {
		h.log.Error(models.ErrMarshal(err, "controllers config defaults"))
	}
}

// GetConnectionControllersConfig handles
// GET /api/integrations/connections/{connectionId}/controllers/config.
// Returns the layered controllers configuration for the connection: the
// per-connection override, the server-wide default, and the resolved
// effective configuration.
func (h *Handler) GetConnectionControllersConfig(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)

	connection, ok := h.fetchKubernetesConnection(w, req, token, provider)
	if !ok {
		return
	}

	response, err := h.buildConnectionControllersConfig(connection)
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error(models.ErrMarshal(err, "connection controllers config"))
	}
}

// UpdateConnectionControllersConfig handles
// PUT /api/integrations/connections/{connectionId}/controllers/config.
// Persists the override into the connection's metadata and applies the
// resolved effective configuration to that connection's cluster. Absent
// fields inherit the server-wide default; an empty document removes the
// override entirely.
func (h *Handler) UpdateConnectionControllersConfig(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	userID := user.ID
	token, _ := req.Context().Value(models.TokenCtxKey).(string)

	override, ok := h.readControllersConfigPayload(w, req)
	if !ok {
		return
	}

	connection, ok := h.fetchKubernetesConnection(w, req, token, provider)
	if !ok {
		return
	}
	connectionID := connection.ID

	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("update")

	metadata := connection.Metadata
	if metadata == nil {
		metadata = core.Map{}
	}
	if err := connections.SetControllersConfigToMetadata(metadata, override); err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, http.StatusBadRequest)
		return
	}
	// Refresh the materialized meshsync_deployment_mode cache so every consumer
	// that predates the layered document (state machine, header chips,
	// kubeconfig flows) keeps seeing the truth. The cache is recomputed from
	// the layered document on every update, not only when the override sets a
	// mode: returning the field to Inherit therefore replaces a stale
	// materialization with the inherited mode, and the mode-change machinery
	// below sees the correct desired mode. Resolution deliberately reads the
	// *pre-update* cache only as a compatibility floor, so it can never mask
	// the value the user just saved.
	serverDefaults, err := models.GetControllersConfigDefaults(h.dbHandler)
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, http.StatusInternalServerError)
		return
	}
	mergedForMode, _ := connections.ResolveControllersConfig(override, serverDefaults)
	desiredMode := connections.ResolveDeploymentMode(mergedForMode, metadata, h.MeshsyncDefaultDeploymentMode)
	connections.MaterializeMeshsyncDeploymentMode(metadata, desiredMode.Mode)

	payload := &connections.ConnectionPayload{
		ID:           connection.ID,
		Kind:         connection.Kind,
		SubType:      connection.SubType,
		Type:         connection.ConnectionType,
		Name:         connection.Name,
		MetaData:     metadata,
		Status:       connection.Status,
		CredentialID: connection.CredentialID,
	}

	// A deployment-mode change tears down and reattaches the controller
	// machinery for the connection (deploy/undeploy operator, restart the
	// MeshSync data pipeline). Delegated to the existing mode-change path so
	// both entry points behave identically.
	if _, _, modeChanged, err := h.handleMeshSyncDeploymentModeChange(req.Context(), connectionID, payload, token, userID, provider); err != nil {
		h.log.Error(err)
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Failed to apply MeshSync deployment mode change.").WithMetadata(map[string]interface{}{"error": err}).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(w, err, http.StatusInternalServerError)
		return
	} else if modeChanged {
		event := eventBuilder.WithSeverity(events.Informational).WithDescription("MeshSync deployment mode changed as part of the controllers configuration update.").Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
	}

	if _, err := provider.UpdateConnectionById(token, payload, connectionID.String()); err != nil {
		h.log.Error(err)
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Failed to persist the connection's controllers configuration override.").WithMetadata(map[string]interface{}{"error": err}).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(w, err, http.StatusInternalServerError)
		return
	}

	// Apply the resolved configuration to this connection's cluster (or
	// restart the embedded run so it picks the new knobs up).
	h.applyControllersConfigToConnection(req.Context(), connectionID, metadata, token, userID, provider, eventBuilder)

	connection.Metadata = metadata
	response, err := h.buildConnectionControllersConfig(connection)
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error(models.ErrMarshal(err, "connection controllers config"))
	}
}

// readControllersConfigPayload decodes and validates a controllers
// configuration write payload. On failure it writes the HTTP error response
// and returns ok=false.
func (h *Handler) readControllersConfigPayload(w http.ResponseWriter, req *http.Request) (*controllersconfig.MesheryControllersConfig, bool) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return nil, false
	}
	payload := &controllersconfig.MesheryControllersConfigPayload{}
	if err := json.Unmarshal(body, payload); err != nil {
		h.log.Error(models.ErrUnmarshal(err, "controllers config"))
		writeMeshkitError(w, models.ErrUnmarshal(err, "controllers config"), http.StatusBadRequest)
		return nil, false
	}
	cfg := &controllersconfig.MesheryControllersConfig{
		SchemaVersion: connections.ControllersConfigSchemaVersion,
		Operator:      payload.Operator,
		Meshsync:      payload.Meshsync,
		Broker:        payload.Broker,
	}
	if err := connections.ValidateControllersConfig(cfg); err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, http.StatusBadRequest)
		return nil, false
	}
	return cfg, true
}

// fetchKubernetesConnection resolves the {connectionId} path parameter into a
// Kubernetes connection. On failure it writes the HTTP error response and
// returns ok=false.
func (h *Handler) fetchKubernetesConnection(w http.ResponseWriter, req *http.Request, token string, provider models.Provider) (*connections.Connection, bool) {
	connectionID, err := uuid.FromString(mux.Vars(req)["connectionId"])
	if err != nil || connectionID == uuid.Nil {
		idErr := ErrEmptyConnectionID()
		h.log.Error(idErr)
		writeMeshkitError(w, idErr, http.StatusBadRequest)
		return nil, false
	}
	connection, statusCode, err := provider.GetConnectionByID(token, connectionID)
	if err != nil {
		h.log.Error(err)
		if statusCode < http.StatusContinue {
			statusCode = http.StatusInternalServerError
		}
		writeMeshkitError(w, err, statusCode)
		return nil, false
	}
	if connection.Kind != "kubernetes" {
		err := connections.ErrControllersConfigInvalid("controllers configuration applies to Kubernetes connections only")
		h.log.Error(err)
		writeMeshkitError(w, err, http.StatusBadRequest)
		return nil, false
	}
	return connection, true
}

// buildConnectionControllersConfig assembles the layered view (override,
// server default, effective) for a connection.
//
// The effective document reports the deployment mode the connection actually
// runs, not just the one the two editable layers set: clients decide from it
// which settings can reach anything on this cluster (Meshery Operator manages
// MeshSync and Meshery Broker, so in embedded mode most of the document is
// inert), and that decision has to be made against the mode the apply path
// resolves.
func (h *Handler) buildConnectionControllersConfig(connection *connections.Connection) (*controllersconfig.ConnectionControllersConfig, error) {
	override, err := connections.ControllersConfigFromMetadata(connection.Metadata)
	if err != nil {
		return nil, err
	}
	serverDefaults, err := models.GetControllersConfigDefaults(h.dbHandler)
	if err != nil {
		return nil, err
	}
	_, effective, _ := connections.ResolveConnectionControllersConfig(
		override, serverDefaults, connection.Metadata, h.MeshsyncDefaultDeploymentMode,
	)
	return &controllersconfig.ConnectionControllersConfig{
		Override:  override,
		Default:   serverDefaults,
		Effective: *effective,
	}, nil
}

// applyControllersConfigToConnection propagates the resolved configuration
// for one connection: operator mode patches the cluster's custom resources
// and the MeshSync deployment overlay; embedded mode restarts the in-process
// MeshSync run so it picks up the new knobs. Failures are surfaced as events
// rather than failing the request: the override is already persisted and
// re-applies on the next connect.
func (h *Handler) applyControllersConfigToConnection(
	ctx context.Context,
	connectionID core.Uuid,
	metadata core.Map,
	token string,
	userID core.Uuid,
	provider models.Provider,
	eventBuilder *events.EventBuilder,
) {
	if h.ConnectionToStateMachineInstanceTracker == nil {
		return
	}
	machine, ok := h.ConnectionToStateMachineInstanceTracker.Get(connectionID)
	if !ok || machine == nil {
		// Not currently tracked (never connected in this server session):
		// nothing to apply now; ConnectAction applies on connect.
		return
	}
	machineCtx, err := kubernetes.GetMachineCtx(machine.Context, nil)
	if err != nil || machineCtx == nil {
		h.log.Warnf("controllers config: no machine context for connection %s; configuration applies on next connect", connectionID)
		return
	}
	ctrlHelper := machineCtx.MesheryCtrlsHelper
	if ctrlHelper == nil {
		return
	}

	merged, _, err := ctrlHelper.ResolveControllersConfigForConnection(metadata)
	if err != nil {
		h.emitControllersConfigApplyEvent(eventBuilder, provider, token, userID, events.Error, "Failed to resolve the connection's controllers configuration.", map[string]interface{}{"error": err, "connectionId": connectionID})
		return
	}
	ctrlHelper.SetControllersConfig(merged)

	// The layered document outranks the materialized meshsync_deployment_mode
	// cache, so a server-wide default change reaches connections that inherit
	// it (see connections.ResolveDeploymentMode).
	resolvedMode := connections.ResolveDeploymentMode(merged, metadata, h.MeshsyncDefaultDeploymentMode)

	switch resolvedMode.Mode {
	case connections.MeshsyncDeploymentModeEmbedded:
		// Restart the in-process run so libmeshsync options (output
		// filters) and, when the target cluster carries a MeshSync CR, the
		// watch-list are re-read.
		contextID := machineCtx.K8sContext.ID
		ctrlHelper.RemoveMeshSyncDataHandler(ctx, contextID)
		ctrlHelper.AddMeshsyncDataHandlers(ctx, machineCtx.K8sContext, userID, *h.SystemID, provider)
		h.emitControllersConfigApplyEvent(eventBuilder, provider, token, userID, events.Informational, "Controllers configuration applied: embedded MeshSync restarted with the updated configuration.", map[string]interface{}{"connectionId": connectionID})
	case connections.MeshsyncDeploymentModeOperator:
		// operator.version fixes the Helm chart version the operator is
		// deployed at, and the controller handler captures it at construction,
		// so a changed version has to be redeployed here rather than waiting
		// for the connection to reconnect.
		if chartVersion, redeployed, verr := ctrlHelper.ReconcileOperatorChartVersion(machineCtx.K8sContext, machineCtx.OperatorTracker); verr != nil {
			h.log.Error(verr)
			h.emitControllersConfigApplyEvent(eventBuilder, provider, token, userID, events.Error, "Failed to deploy Meshery Operator at the configured chart version.", map[string]interface{}{"error": verr, "connectionId": connectionID, "operatorChartVersion": chartVersion})
		} else if redeployed {
			h.emitControllersConfigApplyEvent(eventBuilder, provider, token, userID, events.Informational, fmt.Sprintf("Meshery Operator redeployed at chart version %q.", chartVersion), map[string]interface{}{"connectionId": connectionID, "operatorChartVersion": chartVersion})
		}

		kubeClient, err := machineCtx.K8sContext.GenerateKubeHandler()
		if err != nil {
			h.emitControllersConfigApplyEvent(eventBuilder, provider, token, userID, events.Error, "Failed to reach the connection's cluster to apply the controllers configuration.", map[string]interface{}{"error": err, "connectionId": connectionID})
			return
		}
		applyCtx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		result, err := models.ApplyControllersConfigToCluster(applyCtx, h.log, kubeClient, merged)
		if err != nil {
			h.emitControllersConfigApplyEvent(eventBuilder, provider, token, userID, events.Error, "Failed to apply the controllers configuration to the cluster.", map[string]interface{}{"error": err, "connectionId": connectionID, "result": result})
			return
		}
		h.emitControllersConfigApplyEvent(eventBuilder, provider, token, userID, events.Informational, "Controllers configuration applied to the cluster.", map[string]interface{}{"connectionId": connectionID, "result": result})
	}
}

// reapplyControllersConfigToTrackedConnections fans a server-wide defaults
// change out to every tracked Kubernetes connection.
func (h *Handler) reapplyControllersConfigToTrackedConnections(token string, userID core.Uuid, provider models.Provider) {
	if h.ConnectionToStateMachineInstanceTracker == nil {
		return
	}
	serverDefaults, err := models.GetControllersConfigDefaults(h.dbHandler)
	if err != nil {
		h.log.Error(err)
		return
	}
	h.ConnectionToStateMachineInstanceTracker.Range(func(connectionID core.Uuid, _ *machines.StateMachine) bool {
		connection, _, err := provider.GetConnectionByID(token, connectionID)
		if err != nil {
			h.log.Warnf("controllers config: skipping re-apply for connection %s: %v", connectionID, err)
			return true
		}
		if connection.Kind != "kubernetes" {
			return true
		}
		eventBuilder := events.NewEvent().ActedUpon(connectionID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("update")
		// A deployment-mode change is not just a document to re-apply: it
		// deploys or undeploys the operator and re-attaches the MeshSync data
		// pipeline. Reconcile it first so the subsequent apply targets the
		// cluster shape the new mode implies.
		h.reconcileInheritedDeploymentMode(connection, serverDefaults, token, userID, provider, eventBuilder)
		h.applyControllersConfigToConnection(context.Background(), connectionID, connection.Metadata, token, userID, provider, eventBuilder)
		return true
	})
}

// reconcileInheritedDeploymentMode brings one connection into line with the
// server-wide defaults after they change. A connection that overrides the
// deployment mode is untouched (its override still wins); a connection that
// inherits follows the new default, which means persisting the refreshed
// materialization and actually tearing down and re-attaching the controller
// machinery. Without this, a server-wide mode change reached the document but
// never the cluster.
func (h *Handler) reconcileInheritedDeploymentMode(
	connection *connections.Connection,
	serverDefaults *controllersconfig.MesheryControllersConfig,
	token string,
	userID core.Uuid,
	provider models.Provider,
	eventBuilder *events.EventBuilder,
) {
	metadata := connection.Metadata
	if metadata == nil {
		metadata = core.Map{}
		connection.Metadata = metadata
	}

	override, err := connections.ControllersConfigFromMetadata(metadata)
	if err != nil {
		// An unreadable override is unknown intent, not absent intent. Treating
		// it as absent used to resolve the connection to the server-wide
		// default, find that it differed from the materialized mode, and
		// reconcile - tearing down and redeploying MeshSync into a mode the
		// user never chose, on the strength of corrupt data. Skip this
		// connection instead and leave it exactly as it is; the per-connection
		// GET surfaces the parse failure so it can be corrected deliberately.
		h.log.Error(err)
		h.emitControllersConfigApplyEvent(
			eventBuilder, provider, token, userID, events.Error,
			"Skipped reconciling this connection's MeshSync deployment mode: its stored controllers configuration could not be parsed.",
			map[string]interface{}{"error": err, "connectionId": connection.ID},
		)
		return
	}
	merged, _ := connections.ResolveControllersConfig(override, serverDefaults)
	desired := connections.ResolveDeploymentMode(merged, metadata, h.MeshsyncDefaultDeploymentMode)

	previous := connections.MeshsyncDeploymentModeFromMetadata(metadata)
	if previous == desired.Mode {
		return
	}
	connections.MaterializeMeshsyncDeploymentMode(metadata, desired.Mode)

	payload := &connections.ConnectionPayload{
		ID:           connection.ID,
		Kind:         connection.Kind,
		SubType:      connection.SubType,
		Type:         connection.ConnectionType,
		Name:         connection.Name,
		MetaData:     metadata,
		Status:       connection.Status,
		CredentialID: connection.CredentialID,
	}
	if _, err := provider.UpdateConnectionById(token, payload, connection.ID.String()); err != nil {
		h.log.Error(err)
		h.emitControllersConfigApplyEvent(eventBuilder, provider, token, userID, events.Error, "Failed to persist the MeshSync deployment mode inherited from the server-wide default.", map[string]interface{}{"error": err, "connectionId": connection.ID})
		return
	}

	if err := h.reconcileMeshsyncDeploymentMode(context.Background(), connection.ID, desired.Mode, merged, userID, provider); err != nil {
		h.log.Error(err)
		h.emitControllersConfigApplyEvent(eventBuilder, provider, token, userID, events.Error, "Failed to redeploy MeshSync for the deployment mode inherited from the server-wide default.", map[string]interface{}{"error": err, "connectionId": connection.ID})
		return
	}
	h.emitControllersConfigApplyEvent(eventBuilder, provider, token, userID, events.Informational, fmt.Sprintf("MeshSync deployment mode changed from '%s' to '%s' following the server-wide default.", previous, desired.Mode), map[string]interface{}{
		"meshsyncDeploymentModeOld": previous,
		"meshsyncDeploymentModeNew": desired.Mode,
		"connectionId":              connection.ID,
	})
}

// emitControllersConfigApplyEvent persists and broadcasts a controllers
// configuration apply event.
func (h *Handler) emitControllersConfigApplyEvent(
	eventBuilder *events.EventBuilder,
	provider models.Provider,
	token string,
	userID core.Uuid,
	severity events.EventSeverity,
	description string,
	metadata map[string]interface{},
) {
	event := eventBuilder.WithSeverity(severity).WithDescription(description).WithMetadata(metadata).Build()
	_ = provider.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)
}
