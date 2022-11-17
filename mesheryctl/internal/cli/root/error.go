package root

import "github.com/layer5io/meshkit/errors"

const (
	ErrProcessingConfigCode        = "1050"
	ErrCreatingConfigFileCode      = "1051"
	ErrAddingTokenToConfigCode     = "1052"
	ErrAddingContextToConfigCode   = "1053"
	ErrUnmarshallingConfigFileCode = "1054"
	ErrGettingRequestContextCode   = "1055"
	ErrInvalidAPIResponseCode      = "1056"
	ErrUnmarshallingAPIDataCode    = "1057"
	ErrConnectingToServerCode      = "1058"
)

var (
	ErrCreatingConfigFile = errors.New(ErrCreatingConfigFileCode, errors.Alert, []string{"Unable to create config file"}, []string{"Unable to create config file"}, []string{}, []string{})

	ErrAddingTokenToConfig = errors.New(ErrAddingTokenToConfigCode, errors.Alert, []string{"Unable to add token to config"}, []string{"Unable to add token to config"}, []string{}, []string{})

	ErrAddingContextToConfig = errors.New(ErrAddingContextToConfigCode, errors.Alert, []string{"Unable to add context to config"}, []string{"Unable to add context to config"}, []string{}, []string{})

	ErrUnmarshallingConfigFile = errors.New(ErrUnmarshallingConfigFileCode, errors.Alert, []string{"Error unmarshalling config file"}, []string{"Error unmarshalling config file"}, []string{}, []string{})
)

func ErrProcessingConfig(err error) error {
	return errors.New(ErrProcessingConfigCode, errors.Alert, []string{"Error processing config"}, []string{"Error processing config", err.Error()}, []string{}, []string{})
}

func ErrConnectingToServer(err error) error {
	return errors.New(ErrConnectingToServerCode, errors.Fatal, []string{"Unable to communicate with Meshery server"}, []string{"Unable to communicate with Meshery server", err.Error(), "See https://docs.meshery.io for help getting started with Meshery"}, []string{}, []string{"See https://docs.meshery.io for help getting started with Meshery"})
}

func ErrGettingRequestContext(err error) error {
	return errors.New(ErrGettingRequestContextCode, errors.Fatal, []string{"Unable to add token to config"}, []string{"Unable to add token to config", err.Error()}, []string{}, []string{})
}

func ErrInvalidAPIResponse(err error) error {
	return errors.New(ErrInvalidAPIResponseCode, errors.Fatal, []string{"Invalid API response"}, []string{"Invalid API response", err.Error()}, []string{}, []string{})
}

func ErrUnmarshallingAPIData(err error) error {
	return errors.New(ErrUnmarshallingAPIDataCode, errors.Fatal, []string{"Error unmarshalling API data"}, []string{"Error unmarshalling API data", err.Error()}, []string{}, []string{})
}
