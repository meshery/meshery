package llm

import (
	"context"
	stderrors "errors"
	"fmt"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrInvalidConfigCode         = "meshery-server-1484"
	ErrInvalidCredentialsCode    = "meshery-server-1485"
	ErrAuthFailureCode           = "meshery-server-1486"
	ErrProviderUnavailableCode   = "meshery-server-1487"
	ErrMalformedResponseCode     = "meshery-server-1488"
	ErrUnsupportedModelCode      = "meshery-server-1489"
	ErrContextCancelledCode      = "meshery-server-1490"
	ErrSSRFValidationCode        = "meshery-server-1491"
	ErrLLMInferenceCode          = "meshery-server-1492"
)

func ErrInvalidConfig(err error) error {
	return errors.New(ErrInvalidConfigCode, errors.Alert,
		[]string{"Invalid AI provider configuration"},
		[]string{err.Error()},
		[]string{"The provider configuration is missing required fields or has malformed data."},
		[]string{"Verify the connection configuration for the AI provider."})
}

func ErrInvalidCredentials(err error) error {
	return errors.New(ErrInvalidCredentialsCode, errors.Alert,
		[]string{"Invalid or missing AI provider credentials"},
		[]string{err.Error()},
		[]string{"The provider credentials are missing, improperly formatted, or wrong type."},
		[]string{"Verify the credential attached to the AI connection."})
}

func ErrAuthFailure(err error) error {
	return errors.New(ErrAuthFailureCode, errors.Alert,
		[]string{"Authentication failed with the AI provider"},
		[]string{err.Error()},
		[]string{"The provider rejected the credentials (e.g., 401 Unauthorized)."},
		[]string{"Ensure the API key or token is active and valid for the selected provider."})
}

func ErrProviderUnavailable(err error) error {
	return errors.New(ErrProviderUnavailableCode, errors.Alert,
		[]string{"AI provider is currently unavailable or timed out"},
		[]string{err.Error()},
		[]string{"The provider could not be reached, the network timed out, or the provider returned 5xx."},
		[]string{"Check the provider's status page, verify network connectivity, and check timeout settings."})
}

func ErrMalformedResponse(err error) error {
	return errors.New(ErrMalformedResponseCode, errors.Alert,
		[]string{"Received malformed response from AI provider"},
		[]string{err.Error()},
		[]string{"The AI provider returned a response that could not be parsed or did not match the expected structure."},
		[]string{"Verify the model supports structured output if required, and check the raw response."})
}

func ErrUnsupportedModel(model string) error {
	return errors.New(ErrUnsupportedModelCode, errors.Alert,
		[]string{"Unsupported model requested from AI provider"},
		[]string{fmt.Sprintf("Model requested: %s", model)},
		[]string{"The specified model is not supported by the configured AI provider."},
		[]string{"Select a supported model for this provider."})
}

func ErrContextCancelled(err error) error {
	return errors.New(ErrContextCancelledCode, errors.Alert,
		[]string{"AI request context was cancelled"},
		[]string{err.Error()},
		[]string{"The request was cancelled by the user or the system before it could complete."},
		[]string{"Retry the request if it was unintentionally cancelled."})
}

func ErrSSRFValidation(err error) error {
	return errors.New(ErrSSRFValidationCode, errors.Alert,
		[]string{"AI provider URL failed security validation"},
		[]string{err.Error()},
		[]string{"The URL points to a restricted address (e.g. metadata service) or uses an unsupported scheme."},
		[]string{"Use a safe, routable endpoint. Cloud providers require HTTPS."})
}

func ErrLLMInference(err error) error {
	return errors.New(ErrLLMInferenceCode, errors.Alert,
		[]string{"LLM inference request failed"},
		[]string{err.Error()},
		[]string{"The provider returned an error during text or structure generation."},
		[]string{"Check the provider logs and ensure the prompt/context sizes are within model limits."})
}

// NormalizeError converts external provider errors into secret-safe internal errors.
// It ensures no API keys, Authorization headers, or full URLs leak into the returned error.
func NormalizeError(err error) error {
	if err == nil {
		return nil
	}

	if stderrors.Is(err, context.DeadlineExceeded) {
		return ErrProviderUnavailable(stderrors.New("context deadline exceeded"))
	}
	if stderrors.Is(err, context.Canceled) {
		return ErrContextCancelled(stderrors.New("context cancelled"))
	}

	// We return a generic error if it's an unrecognized error so we don't accidentally leak URLs with inline creds.
	return ErrLLMInference(stderrors.New("provider inference error (details redacted for safety)"))
}
