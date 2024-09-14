package machines

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/utils"
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
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)

	token, _ := ctx.Value(models.TokenCtxKey).(string)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	provider, _ := ctx.Value(models.ProviderCtxKey).(models.Provider)

	payload, err := utils.Cast[connections.ConnectionPayload](data)
	if err != nil {
		return NoOp, eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}
	_, ok := payload.CredentialSecret["id"]
	credential := &models.Credential{}
	// If existing credential is used do not persist again
	if !ok {
		credName, _ := payload.CredentialSecret["name"].(string)
		credential, err = provider.SaveUserCredential(token, &models.Credential{
			Name:   credName,
			UserID: &userUUID,
			Type:   payload.Kind,
			Secret: payload.CredentialSecret,
		})
	}

	if err != nil {
		_err := models.ErrPersistCredential(err)
		return NoOp, eventBuilder.WithDescription(fmt.Sprintf("Unable to persist credential information for the connection %s", payload.Name)).
			WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": err}).Build(), _err
	}

	connection, err := provider.SaveConnection(&connections.ConnectionPayload{
		ID:           payload.ID,
		Kind:         payload.Kind,
		Type:         payload.Type,
		SubType:      payload.SubType,
		Status:       connections.CONNECTED,
		Name:         payload.Name,
		MetaData:     payload.MetaData,
		CredentialID: &credential.ID,
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
