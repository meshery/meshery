package codex

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/core"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RegisterAction struct{}

func (ra *RegisterAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	logLevel := viper.GetInt("LOG_LEVEL")
	if viper.GetBool("DEBUG") {
		logLevel = int(logrus.DebugLevel)
	}
	log, err := logger.New("meshery", logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: logLevel,
	})
	if err != nil {
		return machines.NoOp, nil, err
	}

	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*core.Uuid)
	userUUID := user.ID

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.").WithSeverity(events.Error)

	connPayload, err := utils.Cast[connections.ConnectionPayload](data)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	metadata, err := utils.Cast[map[string]interface{}](connPayload.MetaData)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	codexConn, err := utils.MarshalAndUnmarshal[map[string]interface{}, connections.CodexConn](metadata)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	codexCred, err := utils.MarshalAndUnmarshal[map[string]interface{}, connections.CodexCred](connPayload.CredentialSecret)
	if err != nil && !connPayload.SkipCredentialVerification {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	if !connPayload.SkipCredentialVerification {
		log.Debug("executing connectivity check for Codex connection")
		client := &http.Client{Timeout: 10 * time.Second}

		parsedURL, parseErr := url.Parse(codexConn.BaseURL)
		if parseErr != nil || parsedURL.Scheme != "https" || parsedURL.User != nil {
			validationErr := fmt.Errorf("baseUrl must be a valid https URL without embedded credentials")
			return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": models.ErrCodexConnectivity(validationErr)}).Build(), models.ErrCodexConnectivity(validationErr)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/models", codexConn.BaseURL), nil)
		if err != nil {
			return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": models.ErrCodexConnectivity(err)}).Build(), models.ErrCodexConnectivity(err)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", codexCred.APIKey))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": models.ErrCodexConnectivity(err)}).Build(), models.ErrCodexConnectivity(err)
		}
		defer func() {
			if closeErr := resp.Body.Close(); closeErr != nil {
				log.Error(closeErr)
			}
		}()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			err := fmt.Errorf("codex API unreachable, status %d", resp.StatusCode)
			return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": models.ErrCodexConnectivity(err)}).Build(), models.ErrCodexConnectivity(err)
		}
	}

	return machines.Exit, nil, nil
}

func (ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
