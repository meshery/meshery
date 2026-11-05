package machines

import (
	"context"
	"fmt"

	"github.com/meshery/schemas/models/core"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/utils"
)

// Action to be executed in a given state.
type Action interface {

	// Used as guards/prerequisites checks and actions to be performed when the machine enters a given state.
	ExecuteOnEntry(context context.Context, machinectx interface{}, data interface{}) (EventType, *events.Event, error)

	Execute(context context.Context, machinectx interface{}, data interface{}) (EventType, *events.Event, error)

	// Used for cleanup actions to perform when the machine exits a given state
	ExecuteOnExit(context context.Context, machinectx interface{}, data interface{}) (EventType, *events.Event, error)
}

// DefaultAction does nothing on entry, execution or exit.
type DefaultAction struct{}

func (da *DefaultAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

func (da *DefaultAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

func (da *DefaultAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

// WithPersistence wraps the REGISTERED state's action so that registering a
// connection persists it. It runs inner (the kind's register action, e.g. the
// Grafana/Prometheus reachability check; nil for kinds without one) and, only
// when inner succeeds, saves the connection and its credential with status
// REGISTERED.
//
// Persisting on register keeps the invariant SendEvent relies on: once a
// transition completes there is a connection row whose status it can update.
// Without it, the status write after the register transition fails with
// "record not found" (#20983).
func WithPersistence(inner Action) Action {
	return &PersistingAction{inner: inner}
}

// PersistingAction is the Action returned by WithPersistence.
type PersistingAction struct {
	inner Action
}

// Inner returns the wrapped action, or nil when there is none.
func (pa *PersistingAction) Inner() Action {
	return pa.inner
}

func (pa *PersistingAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {
	if pa.inner == nil {
		return NoOp, nil, nil
	}
	return pa.inner.ExecuteOnEntry(ctx, machineCtx, data)
}

func (pa *PersistingAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {
	next := NoOp
	var event *events.Event
	if pa.inner != nil {
		var err error
		next, event, err = pa.inner.Execute(ctx, machineCtx, data)
		if err != nil {
			// Nothing is persisted: the machine halts in its previous state with
			// no row and no credential written.
			return next, event, err
		}
	}

	// An existing row (re-registration, or a connection auto-registered
	// earlier) is kept as-is; SendEvent updates its status.
	if !connectionPersisted(ctx, data) {
		if _, saveEvent, err := saveConnection(ctx, data, connections.REGISTERED); err != nil {
			return NoOp, saveEvent, err
		}
	}
	return next, event, nil
}

func (pa *PersistingAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {
	if pa.inner == nil {
		return NoOp, nil, nil
	}
	return pa.inner.ExecuteOnExit(ctx, machineCtx, data)
}

type DefaultConnectAction struct{}

func (da *DefaultConnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

func (da *DefaultConnectAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {
	// Register already persisted the row and its credential; SendEvent moves it
	// to CONNECTED. Saving again would create a second, unattached credential,
	// because SaveConnection returns an existing row without setting its
	// CredentialID.
	if connectionPersisted(ctx, data) {
		return NoOp, nil, nil
	}
	return saveConnection(ctx, data, connections.CONNECTED)
}

func (da *DefaultConnectAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

// connectionPersisted reports whether the payload's connection already has a
// row. A payload that cannot be read reports false, so saveConnection surfaces
// the error.
func connectionPersisted(ctx context.Context, data interface{}) bool {
	payload, err := utils.Cast[connections.ConnectionPayload](data)
	if err != nil || payload.ID == uuid.Nil {
		return false
	}
	provider, _ := ctx.Value(models.ProviderCtxKey).(models.Provider)
	token, _ := ctx.Value(models.TokenCtxKey).(string)
	existing, _, err := provider.GetConnectionByID(token, payload.ID)
	return err == nil && existing != nil
}

// saveConnection creates the connection's credential (unless the payload
// references an existing one) and the connection row with the given status.
// This is DefaultConnectAction's original persistence logic, shared with
// WithPersistence.
func saveConnection(ctx context.Context, data interface{}, status connections.ConnectionStatus) (EventType, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*core.Uuid)
	userUUID := user.ID

	token, _ := ctx.Value(models.TokenCtxKey).(string)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.")

	provider, _ := ctx.Value(models.ProviderCtxKey).(models.Provider)

	payload, err := utils.Cast[connections.ConnectionPayload](data)
	if err != nil {
		return NoOp, eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}

	existingCredID, ok := payload.CredentialSecret["id"]
	credential := &models.Credential{}
	// If existing credential is used do not persist again
	if !ok {
		credName, _ := payload.CredentialSecret["name"].(string)
		credential, err = provider.SaveUserCredential(token, &models.Credential{
			Name:   credName,
			UserId: userUUID,
			Type:   payload.Kind,
			Secret: payload.CredentialSecret,
		})
	}

	if err != nil {
		_err := models.ErrPersistCredential(err)
		return NoOp, eventBuilder.WithDescription(fmt.Sprintf("Unable to persist credential information for the connection %s", payload.Name)).
			WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": err}).Build(), _err
	}

	var credentialID *core.Uuid
	if ok {
		idStr, isStr := existingCredID.(string)
		if !isStr {
			parseErr := fmt.Errorf("credential id is not a string")
			_err := models.ErrPersistCredential(parseErr)
			return NoOp, eventBuilder.
				WithDescription(fmt.Sprintf("Invalid credential identifier for the connection %s", payload.Name)).
				WithSeverity(events.Error).
				WithMetadata(map[string]interface{}{"error": parseErr}).Build(), _err
		}
		parsed, parseErr := uuid.FromString(idStr)
		if parseErr != nil {
			_err := models.ErrPersistCredential(parseErr)
			return NoOp, eventBuilder.
				WithDescription(fmt.Sprintf("Invalid credential identifier for the connection %s", payload.Name)).
				WithSeverity(events.Error).
				WithMetadata(map[string]interface{}{"error": parseErr}).Build(), _err
		}
		credentialID = &parsed
	} else {
		credentialID = &credential.ID
	}

	_, err = provider.SaveConnection(&connections.ConnectionPayload{
		ID:           payload.ID,
		Kind:         payload.Kind,
		Type:         payload.Type,
		SubType:      payload.SubType,
		Status:       status,
		Name:         payload.Name,
		MetaData:     payload.MetaData,
		CredentialID: credentialID,
	}, token, false)
	if err != nil {
		_err := models.ErrPersistConnection(err)
		return NoOp, eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to persist the \"%s\" connection details", payload.Name)).WithMetadata(map[string]interface{}{"error": _err}).Build(), _err
	}
	return NoOp, nil, nil
}
