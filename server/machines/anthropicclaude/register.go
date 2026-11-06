package anthropicclaude

import (
	"bytes"
	"context"
	"encoding/json"
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

const AnthropicAPIVersion = "2023-06-01"

type AnthropicConn struct {
	BaseURL      string `json:"baseUrl"`
	DefaultModel string `json:"defaultModel"`
}

type AnthropicCred struct {
	APIKey string `json:"apiKey"`
}

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
		logrus.Error(err)
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

	anthropicConn, err := utils.MarshalAndUnmarshal[map[string]interface{}, AnthropicConn](metadata)
	if err != nil {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	anthropicCred, err := utils.MarshalAndUnmarshal[map[string]interface{}, AnthropicCred](connPayload.CredentialSecret)
	if err != nil && !connPayload.SkipCredentialVerification {
		eventBuilder.WithMetadata(map[string]interface{}{"error": err})
		return machines.NoOp, eventBuilder.Build(), err
	}

	if !connPayload.SkipCredentialVerification {
		parsedURL, err := url.Parse(anthropicConn.BaseURL)
		if err != nil || parsedURL.Scheme != "https" || parsedURL.Host != "api.anthropic.com" {
			err = fmt.Errorf("invalid Anthropic Base URL: must be an https URL for api.anthropic.com")
			err = models.ErrAnthropicConnectivity(err)
			log.Error(err)
			return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": err}).Build(), err
		}

		log.Debug("executing ping test for Anthropic Claude connection")
		client := &http.Client{Timeout: 10 * time.Second}

		body := map[string]interface{}{
			"model":      anthropicConn.DefaultModel,
			"max_tokens": 1,
			"messages": []map[string]interface{}{
				{"role": "user", "content": "ping"},
			},
		}

		bodyBytes, err := json.Marshal(body)
		if err != nil {
			err = models.ErrAnthropicConnectivity(err)
			log.Error(err)
			return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": err}).Build(), err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/messages", anthropicConn.BaseURL), bytes.NewReader(bodyBytes))
		if err != nil {
			err = models.ErrAnthropicConnectivity(err)
			log.Error(err)
			return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": err}).Build(), err
		}

		req.Header.Set("x-api-key", anthropicCred.APIKey)
		req.Header.Set("anthropic-version", AnthropicAPIVersion)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			err = models.ErrAnthropicConnectivity(err)
			log.Error(err)
			return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": err}).Build(), err
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Error(fmt.Errorf("error closing response body: %w", err))
			}
		}()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			err := fmt.Errorf("anthropic API unreachable, status %d", resp.StatusCode)
			err = models.ErrAnthropicConnectivity(err)
			log.Error(err)
			return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": err}).Build(), err
		}
	}

	return machines.Exit, nil, nil
}

func (ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}
