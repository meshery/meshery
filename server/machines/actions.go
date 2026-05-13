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
	"github.com/sirupsen/logrus"
)

// Action to be executed in a given state.
type Action interface {

	// Used as guards/prerequisites checks and actions to be performed when the machine enters a given state.
	ExecuteOnEntry(context context.Context, machinectx interface{}, data interface{}) (EventType, *events.Event, error)

	Execute(context context.Context, machinectx interface{}, data interface{}) (EventType, *events.Event, error)

	// Used for cleanup actions to perform when the machine exits a given state
	ExecuteOnExit(context context.Context, machinectx interface{}, data interface{}) (EventType, *events.Event, error)
}

type DefaultConnectAction struct{}

func (da *DefaultConnectAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}

func (da *DefaultConnectAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {

	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*core.Uuid)
	userUUID := user.ID

	token, _ := ctx.Value(models.TokenCtxKey).(string)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

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

	connection, err := provider.SaveConnection(&connections.ConnectionPayload{
		ID:           payload.ID,
		Kind:         payload.Kind,
		Type:         payload.Type,
		SubType:      payload.SubType,
		Status:       connections.CONNECTED,
		Name:         payload.Name,
		MetaData:     payload.MetaData,
		CredentialID: credentialID,
	}, token, false)
	if err != nil {
		_err := models.ErrPersistConnection(err)
		return NoOp, eventBuilder.WithDescription(fmt.Sprintf("Unable to perisit the \"%s\" connection details", payload.Name)).WithMetadata(map[string]interface{}{"error": _err}).Build(), _err
	}
	logrus.Debug(connection, "grafana connected connection")
	return NoOp, nil, nil
}

func (da *DefaultConnectAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (EventType, *events.Event, error) {
	return NoOp, nil, nil
}
