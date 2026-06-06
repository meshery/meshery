package model_provider

import (
	"context"
	"strings"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/core"
)

type RegisterAction struct{}

func (ra *RegisterAction) ExecuteOnEntry(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

func (ra *RegisterAction) Execute(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	user, ok := ctx.Value(models.UserCtxKey).(*models.User)
	if !ok || user == nil {
		return machines.NoOp, nil, ErrInvalidModelProviderPayload("unknown", "missing user context")
	}
	sysID, ok := ctx.Value(models.SystemIDKey).(*core.Uuid)
	if !ok || sysID == nil {
		return machines.NoOp, nil, ErrInvalidModelProviderPayload("unknown", "missing system ID context")
	}

	eventBuilder := events.NewEvent().
		ActedUpon(user.ID).
		WithCategory("connection").
		WithAction("update").
		FromSystem(*sysID).
		FromUser(user.ID).
		WithDescription("Failed to validate model-provider connection.").
		WithSeverity(events.Error)

	connPayload, err := utils.Cast[connections.ConnectionPayload](data)
	if err != nil {
		return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}

	// Structural validation (required fields, key format) always runs regardless
	// of SkipCredentialVerification. The flag only bypasses network/liveness checks.
	kind := strings.ToLower(connPayload.Kind)
	if err := validateStructure(kind, connPayload); err != nil {
		return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}

	if connPayload.SkipCredentialVerification {
		return machines.Exit, nil, nil
	}

	if err := validateModelProviderPayload(kind, connPayload); err != nil {
		return machines.NoOp, eventBuilder.WithMetadata(map[string]interface{}{"error": err}).Build(), err
	}

	return machines.Exit, nil, nil
}

func (ra *RegisterAction) ExecuteOnExit(ctx context.Context, machineCtx interface{}, data interface{}) (machines.EventType, *events.Event, error) {
	return machines.NoOp, nil, nil
}

// validateStructure checks required fields and key formats without any network I/O.
// It runs unconditionally, even when SkipCredentialVerification is set.
func validateStructure(kind string, payload connections.ConnectionPayload) error {
	switch kind {
	case connections.KindOpenAI:
		return validateOpenAI(payload)
	case connections.KindAnthropic:
		return validateAnthropic(payload)
	case connections.KindAWSBedrock:
		return validateAWSBedrock(payload)
	case connections.KindOllama:
		return validateOllama(payload)
	default:
		return ErrUnsupportedProviderKind(kind)
	}
}

// validateModelProviderPayload performs network/liveness verification for the
// provider (e.g. a test API call). It is skipped when SkipCredentialVerification
// is set. Structural validation via validateStructure always runs first.
func validateModelProviderPayload(kind string, payload connections.ConnectionPayload) error {
	// Network-level verification will be added per provider in a follow-up.
	return nil
}

func validateOpenAI(payload connections.ConnectionPayload) error {
	cred, err := utils.MarshalAndUnmarshal[map[string]interface{}, connections.OpenAICred](payload.CredentialSecret)
	if err != nil {
		return ErrInvalidModelProviderPayload(connections.KindOpenAI, "credential secret could not be parsed")
	}
	if cred.APIKey == "" {
		return ErrMissingCredentialField(connections.KindOpenAI, "apiKey")
	}
	if !strings.HasPrefix(cred.APIKey, "sk-") {
		return ErrInvalidModelProviderPayload(connections.KindOpenAI, "apiKey must start with \"sk-\"")
	}
	return nil
}

func validateAnthropic(payload connections.ConnectionPayload) error {
	cred, err := utils.MarshalAndUnmarshal[map[string]interface{}, connections.AnthropicCred](payload.CredentialSecret)
	if err != nil {
		return ErrInvalidModelProviderPayload(connections.KindAnthropic, "credential secret could not be parsed")
	}
	if cred.APIKey == "" {
		return ErrMissingCredentialField(connections.KindAnthropic, "apiKey")
	}
	if !strings.HasPrefix(cred.APIKey, "sk-ant-") {
		return ErrInvalidModelProviderPayload(connections.KindAnthropic, "apiKey must start with \"sk-ant-\"")
	}
	return nil
}

func validateAWSBedrock(payload connections.ConnectionPayload) error {
	cred, err := utils.MarshalAndUnmarshal[map[string]interface{}, connections.AWSBedrockCred](payload.CredentialSecret)
	if err != nil {
		return ErrInvalidModelProviderPayload(connections.KindAWSBedrock, "credential secret could not be parsed")
	}
	if cred.AccessKeyID == "" {
		return ErrMissingCredentialField(connections.KindAWSBedrock, "accessKeyId")
	}
	if cred.SecretAccessKey == "" {
		return ErrMissingCredentialField(connections.KindAWSBedrock, "secretAccessKey")
	}
	if cred.Region == "" {
		return ErrMissingCredentialField(connections.KindAWSBedrock, "region")
	}
	return nil
}

func validateOllama(payload connections.ConnectionPayload) error {
	conn, err := utils.MarshalAndUnmarshal[map[string]interface{}, connections.ModelProviderConn](payload.MetaData)
	if err != nil {
		return ErrInvalidModelProviderPayload(connections.KindOllama, "connection metadata could not be parsed")
	}
	if conn.BaseURL == "" {
		return ErrInvalidModelProviderPayload(connections.KindOllama, "missing required metadata field \"baseURL\"")
	}
	return nil
}
