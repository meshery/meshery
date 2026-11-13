package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/meshery/schemas/models/core"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/helpers"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/events"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	"github.com/meshery/meshkit/utils"
)

// ProcessConnectionRegistration drives the connection registration state
// machine (POST /api/integrations/connections/register). The wire contract is
// the schemas-defined ConnectionRegistrationEvent: an `initialize` event
// returns the registration bootstrap, every other event advances the tracked
// registration and responds with an empty body.
func (h *Handler) ProcessConnectionRegistration(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method == http.MethodDelete {
		// Deprecated: cancelling via DELETE on this route with the tracker id in
		// the body predates the schemas contract (v1.3.32). New clients cancel via
		// DELETE /api/integrations/connections/register/{registrationId}
		// (CancelConnectionRegister). Retire this branch once no supported client
		// sends the legacy shape.
		h.handleProcessTermination(w, req)
		return
	}

	registrationEvent := connectionv1beta3.ConnectionRegistrationEvent{}
	userUUID := user.ID
	token, err := provider.GetProviderToken(req)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		writeMeshkitError(w, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}
	err = json.NewDecoder(req.Body).Decode(&registrationEvent)
	if err != nil {
		writeMeshkitError(w, models.ErrUnmarshal(err, "connection registration payload"), http.StatusBadRequest)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*h.SystemID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.")

	if registrationEvent.Status == connectionv1beta3.ConnectionRegistrationEventStatusInitialize {
		h.handleRegistrationInitEvent(w, &registrationEvent)
	} else {
		smInstanceTracker := h.ConnectionToStateMachineInstanceTracker
		connectionRegisterPayload := registrationEventToConnectionPayload(&registrationEvent)

		var inst *machines.StateMachine
		if smInstanceTracker != nil {
			if tracked, ok := smInstanceTracker.Get(connectionRegisterPayload.ID); ok {
				inst = tracked
			}
		}

		if inst == nil {
			var machineCtx any
			var initFunc connections.InitFunc

			if strings.EqualFold(connectionRegisterPayload.Kind, "kubernetes") {
				contextID := connectionRegisterPayload.ID.String()
				if provider.GetProviderType() != models.RemoteProviderType && connectionRegisterPayload.MetaData != nil {
					for _, key := range []string{"contextId", "context_id", "k8sContextId", "k8s_context_id"} {
						if val, ok := connectionRegisterPayload.MetaData[key].(string); ok && val != "" {
							contextID = val
							break
						}
					}
				}

				k8scontext, err := provider.GetK8sContext(token, contextID)
				if err != nil {
					_err := ErrGetK8sContexts(err)
					h.log.Error(_err)
					writeMeshkitError(w, _err, http.StatusInternalServerError)
					return
				}

				if k8scontext.ConnectionID == "" {
					k8scontext.ConnectionID = connectionRegisterPayload.ID.String()
					if persister := provider.GetGenericPersister(); persister != nil {
						_ = persister.Model(&models.K8sContext{}).Where("id = ? AND (connection_id IS NULL OR connection_id = '')", k8scontext.ID).Update("connection_id", k8scontext.ConnectionID).Error
					}
				}

				var operatorTracker *models.OperatorTracker
				var k8scontextChannel *models.K8scontextChan
				var eventBroadcaster *models.Broadcast
				if h.config != nil {
					operatorTracker = h.config.OperatorTracker
					k8scontextChannel = h.config.K8scontextChannel
					eventBroadcaster = h.config.EventBroadcaster
				}

				machineCtx = &kubernetes.MachineCtx{
					K8sContext:         k8scontext,
					MesheryCtrlsHelper: h.MesheryCtrlsHelper,
					K8sCompRegHelper:   h.K8sCompRegHelper,
					OperatorTracker:    operatorTracker,
					K8scontextChannel:  k8scontextChannel,
					EventBroadcaster:   eventBroadcaster,
					RegistryManager:    h.registryManager,
				}
				initFunc = kubernetes.AssignInitialCtx
			}

			var err error
			inst, err = helpers.InitializeMachineWithContext(
				machineCtx,
				req.Context(),
				connectionRegisterPayload.ID,
				userUUID,
				smInstanceTracker,
				h.log,
				provider,
				machines.DISCOVERED,
				strings.ToLower(connectionRegisterPayload.Kind),
				initFunc,
			)
			if err != nil {
				wrappedErr := ErrInitializeMachine(err)
				event := eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to persist the \"%s\" connection details", connectionRegisterPayload.Kind)).WithMetadata(map[string]interface{}{
					"error": wrappedErr,
				}).Build()
				if event != nil {
					h.persistAndPublish(provider, token, userUUID, event)
				}
				h.log.Error(wrappedErr)
				writeMeshkitError(w, wrappedErr, http.StatusInternalServerError)
				return
			}
		}

		event, err := inst.SendEvent(req.Context(), machines.EventType(registrationEvent.Status), connectionRegisterPayload)
		if err != nil {
			wrappedErr := ErrSendMachineEvent(err)
			h.log.Error(wrappedErr)
			writeMeshkitError(w, wrappedErr, http.StatusInternalServerError)
			if event != nil {
				go h.persistAndPublish(provider, token, userUUID, event)
			}
			return
		}
	}
}

// registrationEventToConnectionPayload converts the schemas wire event into the
// internal payload the connection state machines consume
// (machines/*/register.go casts the event data to connections.ConnectionPayload).
// The two shapes share the same wire fields; only the Go types differ.
func registrationEventToConnectionPayload(event *connectionv1beta3.ConnectionRegistrationEvent) connections.ConnectionPayload {
	payload := connections.ConnectionPayload{
		Kind:                       event.Kind,
		SubType:                    event.SubType,
		Type:                       event.Type,
		MetaData:                   event.Metadata,
		Status:                     connections.ConnectionStatus(event.Status),
		CredentialSecret:           event.CredentialSecret,
		Name:                       event.Name,
		CredentialID:               event.CredentialID,
		Model:                      event.Model,
		SkipCredentialVerification: event.SkipCredentialVerification,
	}
	if event.ID != nil {
		payload.ID = *event.ID
	}
	return payload
}

// handleProcessTermination is the deprecated body-based cancel path; see the
// deprecation note in ProcessConnectionRegistration.
func (h *Handler) handleProcessTermination(w http.ResponseWriter, req *http.Request) {
	body := make(map[string]string, 0)
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		_err := models.ErrUnmarshal(err, "request body")
		h.log.Error(_err)
		writeMeshkitError(w, _err, http.StatusBadRequest)
		return
	}
	smInstancetracker := h.ConnectionToStateMachineInstanceTracker

	id, ok := body["id"]
	if ok {
		smInstancetracker.Remove(uuid.FromStringOrNil(id))
	}
}

// CancelConnectionRegister discards the in-progress registration state machine
// tracked by the {registrationId} path parameter
// (DELETE /api/integrations/connections/register/{registrationId}). Idempotent:
// unknown ids are ignored. Nothing is persisted for the abandoned process.
func (h *Handler) CancelConnectionRegister(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	registrationID, err := uuid.FromString(mux.Vars(req)["registrationId"])
	if err != nil {
		invalidIDErr := models.ErrInvalidUUID(err)
		h.log.Error(invalidIDErr)
		writeMeshkitError(w, invalidIDErr, http.StatusBadRequest)
		return
	}

	h.ConnectionToStateMachineInstanceTracker.Remove(registrationID)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleRegistrationInitEvent(w http.ResponseWriter, event *connectionv1beta3.ConnectionRegistrationEvent) {
	compFilter := &regv1beta1.ComponentFilter{
		Name:  fmt.Sprintf("%sConnection", event.Kind),
		Limit: 1,
	}
	connectionComponent, _, _, err := h.registryManager.GetEntities(compFilter)
	if err != nil {
		h.log.Error(ErrWriteResponse(err))
		writeMeshkitError(w, ErrWriteResponse(err), http.StatusInternalServerError)
		return
	}

	if len(connectionComponent) == 0 {
		writeMeshkitError(w, ErrUnknownConnectionKind(event.Kind), http.StatusBadRequest)
		return
	}

	connectionDefinition, err := utils.MarshalAndUnmarshal[interface{}, core.Map](connectionComponent[0])
	if err != nil {
		h.log.Error(ErrWriteResponse(err))
		writeMeshkitError(w, ErrWriteResponse(err), http.StatusInternalServerError)
		return
	}

	// The bootstrap id acts as the connection registration process tracker.
	// Clients echo it on every subsequent event until the process completes or
	// is cancelled.
	bootstrap := connectionv1beta3.ConnectionRegistrationBootstrap{
		Connection: connectionDefinition,
		ID:         uuid.Must(uuid.NewV4()),
	}

	credential, _, _, err := h.registryManager.GetEntities(&regv1beta1.ComponentFilter{
		Name:  fmt.Sprintf("%sCredential", event.Kind),
		Limit: 1,
	})
	if err != nil {
		h.log.Error(ErrWriteResponse(err))
		writeMeshkitError(w, ErrWriteResponse(err), http.StatusInternalServerError)
		return
	}
	if len(credential) > 0 {
		credentialDefinition, err := utils.MarshalAndUnmarshal[interface{}, core.Map](credential[0])
		if err != nil {
			h.log.Error(ErrWriteResponse(err))
			writeMeshkitError(w, ErrWriteResponse(err), http.StatusInternalServerError)
			return
		}
		bootstrap.Credential = credentialDefinition
	}

	err = json.NewEncoder(w).Encode(&bootstrap)
	if err != nil {
		h.log.Error(ErrWriteResponse(err))
	}
}

func (h *Handler) SaveConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	userID := user.ID
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	connection := connections.ConnectionPayload{}
	err = json.Unmarshal(bd, &connection)
	obj := "connection"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusBadRequest)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(userID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("create")

	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	if token == "" {
		if ck, err := req.Cookie(models.TokenCookieName); err == nil {
			token = ck.Value
		}
	}

	_, err = provider.SaveConnection(&connection, token, false)
	if err != nil {
		_err := ErrFailToSave(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Error creating connection %s", connection.Name)).WithMetadata(metadata).Build()
		h.persistAndPublish(provider, token, userID, event)

		h.log.Error(_err)
		writeMeshkitError(w, _err, http.StatusInternalServerError)
		return
	}

	description := fmt.Sprintf("Connection %s created.", connection.Name)

	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).Build()
	h.persistAndPublish(provider, token, userID, event)

	h.log.Info(description)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetConnections(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	q := req.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	order := q.Get("order")
	search := q.Get("search")
	// Canonical camelCase `pageSize` (matches the schemas wire contract); fall
	// back to legacy lowercase `pagesize` for older clients.
	pageSizeStr := q.Get("pageSize")
	if pageSizeStr == "" {
		pageSizeStr = q.Get("pagesize")
	}
	filter := q.Get("filter")
	name := q.Get("name")

	var pageSize int
	if pageSizeStr == "all" {
		pageSize = 100
	} else {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}

	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if page < 0 {
		page = 0
	}
	if order == "" {
		order = "updated_at desc"
	}

	err := req.ParseForm()
	if err != nil {
		h.log.Error(ErrGetConnections(err))
		writeMeshkitError(w, ErrGetConnections(err), http.StatusInternalServerError)
		return
	}

	// Filters are passed as repeated query params — the standard convention
	// across the API (e.g. ?kind=kubernetes&kind=meshery&status=connected). No
	// per-param JSON decoding.
	queryParam := struct {
		Status []string
		Kind   []string
		Type   []string
	}{
		Status: q["status"],
		Kind:   q["kind"],
		Type:   q["type"],
	}

	connectionsPage, err := provider.GetConnections(req, user.ID.String(), page, pageSize, search, order, filter, queryParam.Status, queryParam.Kind, queryParam.Type, name)
	obj := "connections"

	if err != nil {
		h.log.Error(ErrQueryGet(obj))
		writeMeshkitError(w, ErrQueryGet(obj), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(connectionsPage); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
		writeMeshkitError(w, models.ErrEncoding(err, obj), http.StatusInternalServerError)
		return
	}
}

// GetConnectionsByKind is an internal handler that fetches connections filtered by kind.
// Note: This handler is used internally by other handlers (e.g., GrafanaConfigHandler, PrometheusConfigHandler)
// and is not exposed as an HTTP route since the route was deprecated in schemas v0.8.115.
func (h *Handler) GetConnectionsByKind(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	q := req.URL.Query()
	connectionKind := mux.Vars(req)["connectionKind"]
	page, _ := strconv.Atoi(q.Get("page"))
	order := q.Get("order")
	search := q.Get("search")
	pageSize, _ := strconv.Atoi(q.Get("pagesize"))

	if pageSize > 25 {
		pageSize = 25
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if page < 0 {
		page = 0
	}
	if order == "" {
		order = "updated_at desc"
	}

	h.log.Debug(fmt.Sprintf("page: %d, page size: %d, search: %s, order: %s, kind: %s", page+1, pageSize, search, order, connectionKind))

	// Use GetConnections with kind filter
	connectionsPage, err := provider.GetConnections(req, user.ID.String(), page, pageSize, search, order, "", nil, []string{connectionKind}, nil, "")
	obj := "connections"

	if err != nil {
		// Provider implementations return a mix of bare errors and
		// MeshKit-wrapped ones depending on whether the failure was
		// inside DoRequest, in unmarshal, or in the local DAO. Wrap
		// uniformly so the JSON envelope always carries MeshKit
		// metadata.
		wrappedErr := ErrGetConnections(err)
		h.log.Error(wrappedErr)
		writeMeshkitError(w, wrappedErr, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(connectionsPage); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
		writeMeshkitError(w, models.ErrEncoding(err, obj), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetConnectionByID(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionId"])
	if connectionID == uuid.Nil {
		invalidIDErr := models.ErrInvalidUUID(fmt.Errorf("invalid connection ID"))
		h.log.Error(invalidIDErr)
		writeMeshkitError(w, invalidIDErr, http.StatusBadRequest)
		return
	}

	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connection, statusCode, err := provider.GetConnectionByID(token, connectionID)
	obj := "connection"

	if err != nil {
		h.log.Error(ErrQueryGet(obj))
		writeMeshkitError(w, ErrQueryGet(obj), statusCode)
		return
	}

	if err := json.NewEncoder(w).Encode(connection); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
		writeMeshkitError(w, models.ErrEncoding(err, obj), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateConnectionById(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionID, err := uuid.FromString(mux.Vars(req)["connectionId"])
	if err != nil {
		invalidIDErr := models.ErrInvalidUUID(fmt.Errorf("invalid connection ID: %w", err))
		h.log.Error(invalidIDErr)
		writeMeshkitError(w, invalidIDErr, http.StatusBadRequest)
		return
	}
	// No connection has the nil ID, so there is nothing to update. Like a
	// malformed ID, it is a client error: issue #22999 requires a nil ID
	// to answer 400, not 404.
	if connectionID == uuid.Nil {
		invalidIDErr := models.ErrInvalidUUID(fmt.Errorf("invalid connection ID: %s", connectionID))
		h.log.Error(invalidIDErr)
		writeMeshkitError(w, invalidIDErr, http.StatusBadRequest)
		return
	}
	userID := user.ID

	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("update")

	connection := &connections.ConnectionPayload{}
	err = json.Unmarshal(bd, connection)
	obj := "connection"
	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusBadRequest)
		return
	}

	// MeshSync deployment-mode changes are handled by the dedicated
	// POST /api/integrations/connections/{connectionId}/actions endpoint
	// (PerformConnectionAction), which owns the metadata merge and cluster-side
	// redeploy. PUT here only updates connection fields; it no longer sniffs the
	// metadata for a mode change.

	token, err := provider.GetProviderToken(req)
	if err != nil {
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()

		h.persistAndPublish(provider, token, userID, event)
		writeMeshkitError(w, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}

	// An update never creates a connection: answer 404 for an ID the provider
	// does not know, before anything is saved.
	existing, statusCode, gerr := provider.GetConnectionByID(token, connectionID)
	if gerr != nil && statusCode == http.StatusNotFound {
		h.log.Warnf("No connection with ID %q found to update", connectionID)
		writeMeshkitError(w, models.ErrResultNotFound(gerr), http.StatusNotFound)
		return
	}

	// A lookup that fails for any other reason leaves the connection's
	// existence unknown: stop here rather than saving over a row the
	// provider could not read back. The local provider backfills omitted
	// fields from the lookup, so proceeding would let a partial record
	// clobber the stored one.
	if gerr != nil {
		// Log the real lookup error for diagnosis; the client gets a
		// lookup-specific error, not ErrQueryGet's "pass the query
		// parameter" guidance, which is wrong here — this PUT carries
		// its id in the path, and a lookup failure is a provider
		// failure, not a missing query parameter.
		h.log.Error(gerr)
		lookupErr := models.ErrFetch(fmt.Errorf("connection lookup failed"), obj, statusCode)
		writeMeshkitError(w, lookupErr, http.StatusInternalServerError)
		return
	}

	// A PUT here is a partial update (the UI's connect action sends only
	// {status}); the provider's UpdateConnectionById backfills any omitted field
	// from the persisted row so a partial payload never clobbers columns like
	// kind/name/type/metadata.
	// The status the connection had before this update: the machine reads
	// the saved status, already the new one, so it is told the old one to
	// resolve the unreachable condition of a connection leaving NOTFOUND.
	statusCtx := req.Context()
	if connection.Status != "" && existing != nil {
		statusCtx = machines.WithPriorStatus(statusCtx, existing.Status)
	}

	// The URL identifies the connection being updated: a body carrying a
	// different non-nil id must not redirect the update onto another row (or
	// create one under it). Both providers only fall back to the URL id when
	// the payload id is nil, so pin it here for both of them.
	connection.ID = connectionID

	updatedConnection, err := provider.UpdateConnectionById(token, connection, mux.Vars(req)["connectionId"])
	if err != nil {
		_err := ErrFailToSave(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Error updating connection").WithMetadata(metadata).Build()
		h.persistAndPublish(provider, token, userID, event)

		h.log.Error(_err)
		writeMeshkitError(w, _err, http.StatusInternalServerError)
		return
	}

	// TODO enhance event with information about meshsync deployment mode change
	description := fmt.Sprintf("Connection %s updated.", updatedConnection.Name)
	eventBuilder = eventBuilder.WithDescription(description)

	if connection.Status != "" {
		kind := updatedConnection.Kind
		if kind == "" && existing != nil {
			kind = existing.Kind
		}
		event, notifyErr := h.NotifySmOfConnectionStatusChange(statusCtx, userID, provider, token, connection, kind)
		if notifyErr != nil {
			// The new status is already saved, so the response stays 200, but the
			// failure must not vanish: log it as well as persisting the event.
			h.log.Error(notifyErr)
		}
		_, _ = provider.PersistEvent(event, token)
	}

	event := eventBuilder.WithSeverity(events.Informational).Build()
	h.persistAndPublish(provider, token, userID, event)
	h.log.Info(description)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) NotifySmOfConnectionStatusChange(ctx context.Context, userID core.Uuid, provider models.Provider, token string, connection *connections.ConnectionPayload, kind string) (events.Event, error) {
	connectionID := connection.ID

	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("update")

	if connection.Status != "" {
		// Everything below is Kubernetes specific: the K8s context lookup, the
		// kubernetes.MachineCtx and the "kubernetes" machine. The caller has
		// already persisted the new status, so any other kind must not enter
		// this path. On the local provider it fails the context lookup (and the
		// caller used to discard that error); on a remote provider it would run
		// Kubernetes actions against a non-Kubernetes connection.
		// An empty kind keeps the previous behaviour.
		if kind != "" && !strings.EqualFold(kind, "kubernetes") {
			h.log.Debug(fmt.Sprintf("status of %s connection %s changed to %q, no Kubernetes machine to notify", kind, connectionID, connection.Status))
			return *eventBuilder.
				WithSeverity(events.Informational).
				WithDescription(fmt.Sprintf("Status of %s connection %s updated to %q", kind, connectionID, connection.Status)).
				Build(), nil
		}
		smInstanceTracker := h.ConnectionToStateMachineInstanceTracker
		// token, _ := req.Context().Value(models.TokenCtxKey).(string)
		k8scontext, err := provider.GetK8sContext(token, connectionID.String())

		if err != nil {
			eventBuilder = eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Failed to update connection status for %s", connectionID)).WithMetadata(map[string]interface{}{
				"error": err,
			})

			return *eventBuilder.Build(), err
		}

		eventBuilder = eventBuilder.WithSeverity(events.Informational).
			WithDescription(fmt.Sprintf("Processing status update to \"%s\" for connection %s", connection.Status, k8scontext.Name)).
			WithMetadata(map[string]interface{}{
				"connectionName": k8scontext.Name,
			})

		machineCtx := &kubernetes.MachineCtx{
			K8sContext:         k8scontext,
			MesheryCtrlsHelper: h.MesheryCtrlsHelper,
			K8sCompRegHelper:   h.K8sCompRegHelper,
			OperatorTracker:    h.config.OperatorTracker,
			K8scontextChannel:  h.config.K8scontextChannel,
			EventBroadcaster:   h.config.EventBroadcaster,
			RegistryManager:    h.registryManager,
		}

		inst, err := helpers.InitializeMachineWithContext(
			machineCtx,
			ctx,
			connectionID,
			userID,
			smInstanceTracker,
			h.log,
			provider,
			machines.InitialState,
			"kubernetes",
			kubernetes.AssignInitialCtx,
		)

		// A connection being deleted must not leave a tracker entry behind,
		// whichever way its machine failed. If an uninitialized or broken instance
		// was previously present in the tracker, both the error return below and
		// the no-context return after it would otherwise strand the entry.
		// Nothing can drive that machine afterwards, and the delete paths that
		// normally Remove it (the goroutine below, DeleteContext) both do so only
		// after a SendEvent that cannot succeed without a Context.
		usable := err == nil && helpers.HasMachineContext(inst)
		if !usable && connection.Status == connections.DELETED {
			smInstanceTracker.Remove(connectionID)
		}

		if err != nil {
			eventBuilder = eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Failed to update connection status for %s", connectionID)).WithMetadata(map[string]interface{}{
				"error": err,
			})
			return *eventBuilder.Build(), err
		}
		// A machine whose Context was never assigned cannot service the event:
		// SendEvent would only fail on ErrAssertMachineCtx and publish an error
		// the user can do nothing about. Same shape as the DeleteContext guard in
		// contexts_handler.go; see helpers.HasMachineContext.
		if !usable {
			h.log.Debug(fmt.Sprintf("machine instance for connection %s has no context assigned, skipping the %q event", connectionID, connection.Status))
			return *eventBuilder.Build(), nil
		}

		// detach from the http request lifecycle so that the goroutine isn't cancelled when
		// the handler returns, while preserving context values and enforcing a bounded execution deadline.
		detachedCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), defaultFSMTransitionTimeout)
		go func(inst *machines.StateMachine, status connections.ConnectionStatus) {
			defer cancel()
			event, err := inst.SendEvent(detachedCtx, machines.EventType(helpers.StatusToEvent(status)), nil)
			if err != nil {
				h.log.Error(err)
				h.reportMachineOutcome(provider, token, userID, event, err)
				return
			}

			if status == connections.DELETED {
				smInstanceTracker.Remove(inst.ID)
			}

			h.reportMachineOutcome(provider, token, userID, event, nil)
		}(inst, connection.Status)
	}

	return *eventBuilder.Build(), nil
}

func (h *Handler) DeleteConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionId"])
	userID := user.ID

	token, err := provider.GetProviderToken(req)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		writeMeshkitError(w, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}
	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("delete")

	// A Kubernetes connection deleted while NOTFOUND retires its unreachable
	// condition, so the deletion event resolves it.
	resolvesUnreachable := false
	if existing, _, gerr := provider.GetConnectionByID(token, connectionID); gerr == nil && existing != nil {
		resolvesUnreachable = existing.Kind == "kubernetes" && existing.Status == connections.ConnectionStatus(machines.NOTFOUND)
	}

	deletedConnection, err := provider.DeleteConnection(req, connectionID)
	if err != nil {
		// Deleting an id the provider never issued is "not found", not a
		// failure to write. The local provider says so with ErrResultNotFound;
		// the remote provider says so with its own 404, which ErrDelete now
		// carries (models.ErrDelete -> httputil.WithProviderStatus). Answer 404
		// for either, and answer it *before* the failure event: a delete
		// addressed to a stale id is a client mistake, not an incident worth a
		// persisted error notification.
		//
		// Previously every failure - including a plain 404 - was wrapped in
		// ErrFailToSave and answered 500, so `mesheryctl connection delete
		// <unknown-uuid>` reported "Failed to Save: .connection"
		// (meshery-server-1051), "Provider Database could be down or not
		// reachable". That is the wrong verb, the wrong code and the wrong
		// diagnosis for a connection that simply is not there.
		if errors.GetCode(err) == models.ErrResultNotFoundCode || providerStatus(err) == http.StatusNotFound {
			h.log.Warnf("No connection with ID %q found to delete", connectionID)
			writeMeshkitError(w, err, http.StatusNotFound)
			return
		}

		_err := ErrFailToDelete(err, "connection")
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Error deleting connection").WithMetadata(metadata).Build()
		h.persistAndPublish(provider, token, userID, event)

		h.log.Error(_err)
		// Surface the status the provider actually returned (403, 409, 5xx ...)
		// rather than reporting every upstream refusal as a fault inside
		// Meshery Server. The fallback stays 500 because this path is also the
		// local provider's, where an untagged failure really is ours.
		writeMeshkitError(w, _err, providerStatusOrInternal(err))
		return
	}
	// The initial lookup may fail even when deletion succeeds.
	// Use the deleted connection's details when the provider returns them.
	if deletedConnection != nil &&
		deletedConnection.Kind == "kubernetes" &&
		deletedConnection.Status == connections.ConnectionStatus(machines.NOTFOUND) {
		resolvesUnreachable = true
	}

	// Nothing can drive a deleted connection's machine any more; left in the
	// tracker, the Kubernetes middlewares would keep probing it.
	if h.ConnectionToStateMachineInstanceTracker != nil {
		h.ConnectionToStateMachineInstanceTracker.Remove(connectionID)
	}

	// A provider is free to answer a successful delete without echoing the
	// deleted row back; dereferencing it unconditionally panicked the request.
	deletedName := connectionID.String()
	if deletedConnection != nil && deletedConnection.Name != "" {
		deletedName = deletedConnection.Name
	}
	description := fmt.Sprintf("Connection %s deleted.", deletedName)
	if resolvesUnreachable {
		eventBuilder = eventBuilder.WithCategory("connection").WithSeverity(events.Success).
			Resolves(kubernetes.UnreachableConditionKey(connectionID), time.Now().UTC())
	} else {
		eventBuilder = eventBuilder.WithSeverity(events.Informational)
	}
	event := eventBuilder.WithDescription(description).Build()

	h.reportMachineOutcome(provider, token, userID, event, nil)

	h.log.Info("connection deleted.")
	w.WriteHeader(http.StatusOK)
}

// handleMeshSyncDeploymentModeChange retrieves existing connection, compares meshsync deployment modes
// between existing and new connections, and performs necessary actions when they differ
// Returns: oldMode, newMode, changed, error
func (h *Handler) handleMeshSyncDeploymentModeChange(
	ctx context.Context,
	connectionID core.Uuid,
	newConnection *connections.ConnectionPayload,
	token string,
	userID core.Uuid,
	provider models.Provider,
) (connections.MeshsyncDeploymentMode, connections.MeshsyncDeploymentMode, bool, error) {
	if newConnection == nil {
		return connections.MeshsyncDeploymentModeUndefined, connections.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("new connection is nil, cannot compare meshsync deployment modes")
	}

	if h.SystemID == nil {
		return connections.MeshsyncDeploymentModeUndefined, connections.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("system ID is not configured in handler")
	}
	// TODO is h.SystemID a correct instance id here?
	mesheryInstanceID := *h.SystemID

	// Retrieve existing connection for mode comparison
	existingConnection, statusCode, err := provider.GetConnectionByID(token, connectionID)
	if err != nil {
		return connections.MeshsyncDeploymentModeUndefined, connections.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("failed to retrieve existing connection (status %d): %w", statusCode, err)
	}

	if existingConnection == nil {
		return connections.MeshsyncDeploymentModeUndefined, connections.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("existing connection is nil, cannot compare meshsync deployment modes")
	}

	if existingConnection.Kind != "kubernetes" {
		return connections.MeshsyncDeploymentModeUndefined, connections.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("connection is not of kind kubernetes")
	}

	existingMeshSyncMode := connections.MeshsyncDeploymentModeFromMetadata(existingConnection.Metadata)
	newMeshSyncMode := connections.MeshsyncDeploymentModeFromMetadata(newConnection.MetaData)

	// draw back to default mode
	if newMeshSyncMode == connections.MeshsyncDeploymentModeUndefined {
		newMeshSyncMode = h.MeshsyncDefaultDeploymentMode
	}

	meshSyncModeChanged := existingMeshSyncMode != newMeshSyncMode
	if meshSyncModeChanged {
		instanceTracker := h.ConnectionToStateMachineInstanceTracker
		if instanceTracker == nil {
			return existingMeshSyncMode, newMeshSyncMode, false, fmt.Errorf("instance tracker is nil in handler instance")
		}

		machine, ok := instanceTracker.Get(connectionID)
		if !ok || machine == nil {
			return existingMeshSyncMode, newMeshSyncMode, false, fmt.Errorf("instance tracker does not contain machine for connection %s", connectionID)
		}

		machineCtx, err := kubernetes.GetMachineCtx(machine.Context, nil)
		if err != nil {
			return existingMeshSyncMode, newMeshSyncMode, false, fmt.Errorf("failed to get machine context for connection %s: %w", connectionID, err)
		}

		if machineCtx == nil {
			return existingMeshSyncMode, newMeshSyncMode, false, fmt.Errorf("machine context is nil for connection %s", connectionID)
		}

		ctrlHelper := machineCtx.MesheryCtrlsHelper
		if ctrlHelper == nil {
			return existingMeshSyncMode, newMeshSyncMode, false, fmt.Errorf("machine context does not contain reference to MesheryCtrlsHelper for connection %s", connectionID)
		}

		// Resolve BEFORE tearing anything down. A resolution failure means the
		// intended configuration is unknown, and undeploying the operators only
		// to rebuild them from the last-known-good document is a partial
		// reconciliation performed on unreadable input - the same fault as
		// reconciling a connection whose override will not parse. Fail with the
		// connection left exactly as it is instead.
		merged, _, resolveErr := ctrlHelper.ResolveControllersConfigForConnection(newConnection.MetaData)
		if resolveErr != nil {
			h.log.Error(resolveErr)
			return existingMeshSyncMode, newMeshSyncMode, false, resolveErr
		}

		// disconnect
		{
			contextID := machineCtx.K8sContext.ID
			ctrlHelper.
				UpdateOperatorsStatusMap(machineCtx.OperatorTracker).
				UndeployDeployedOperators(machineCtx.OperatorTracker, contextID).
				RemoveCtxControllerHandler(ctx, contextID)
			ctrlHelper.RemoveMeshSyncDataHandler(ctx, contextID)
		}
		// connect. The controllers configuration is refreshed and stashed
		// before AddCtxControllerHandlers, which constructs the operator
		// controller handler with the Helm chart version that document
		// resolves to (operator.version) and captures it there. Resolving from
		// the incoming payload's metadata rather than the helper's stale copy
		// is what lets one save change the mode and the chart version at once.
		{
			ctrlHelper.SetControllersConfig(merged)
			ctrlHelper.
				SetMeshsyncDeploymentMode(newMeshSyncMode).
				AddCtxControllerHandlers(machineCtx.K8sContext).
				UpdateOperatorsStatusMap(machineCtx.OperatorTracker).
				DeployUndeployedOperators(machineCtx.OperatorTracker, machineCtx.K8sContext.ID).
				AddMeshsyncDataHandlers(ctx, machineCtx.K8sContext, userID, mesheryInstanceID, provider)
		}

	}

	return existingMeshSyncMode, newMeshSyncMode, meshSyncModeChanged, nil
}
