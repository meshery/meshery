package model_provider

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrInvalidModelProviderPayloadCode = "meshery-server-1432"
	ErrMissingCredentialFieldCode      = "meshery-server-1433"
	ErrUnsupportedProviderKindCode     = "meshery-server-1434"
)

func ErrInvalidModelProviderPayload(kind string, reason string) error {
	return errors.New(
		ErrInvalidModelProviderPayloadCode,
		errors.Alert,
		[]string{fmt.Sprintf("invalid payload for model-provider connection of kind %q: %s", kind, reason)},
		[]string{},
		[]string{"Ensure the connection metadata and credential secret match the required schema for the provider."},
		[]string{"Provide all required fields for the chosen provider kind."},
	)
}

func ErrMissingCredentialField(kind, field string) error {
	return errors.New(
		ErrMissingCredentialFieldCode,
		errors.Alert,
		[]string{fmt.Sprintf("missing required credential field %q for provider %q", field, kind)},
		[]string{},
		[]string{"The credential secret is missing a required field."},
		[]string{fmt.Sprintf("Add the %q field to the credential secret for the %q provider.", field, kind)},
	)
}

func ErrUnsupportedProviderKind(kind string) error {
	return errors.New(
		ErrUnsupportedProviderKindCode,
		errors.Alert,
		[]string{fmt.Sprintf("unsupported model-provider kind %q", kind)},
		[]string{},
		[]string{"The requested provider kind is not supported."},
		[]string{"Use one of the supported provider kinds: openai, anthropic, aws-bedrock, ollama."},
	)
}
