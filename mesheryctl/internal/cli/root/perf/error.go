package perf

import (
	"strconv"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrMesheryConfigCode       = "replace me"
	ErrReadFilepathCode        = "replace me"
	ErrFailNewRequestCode      = "replace me"
	ErrNoProfileNameCode       = "replace me"
	ErrNoTestURLCode           = "replace me"
	ErrNotValidURLCode         = "replace me"
	ErrFailConvReqCode         = "replace me"
	ErrFailConvQPSCode         = "replace me"
	ErrFailMarshalCode         = "replace me"
	ErrAttachAuthTokenCode     = "replace me"
	ErrFailRequestCode         = "replace me"
	ErrFailReqStatusCode       = "replace me"
	ErrFailReadReqBodyCode     = "replace me"
	ErrFailUnmarshalCode       = "replace me"
	ErrNoProfileFoundCode      = "replace me"
	ErrFailTestRunCode         = "replace me"
	ErrStartTermboxCode        = "replace me"
	ErrConvJSONYAMLCode        = "replace me"
	ErrInvalidOutputChoiceCode = "replace me"
	ErrNoResultFoundCode       = "replace me"
)

func ErrMesheryConfig(err error) error {
	return errors.New(ErrMesheryConfigCode, errors.Alert, []string{"error processing config"}, []string{err.Error()}, []string{}, []string{})
}

func ErrReadFilepath(err error) error {
	return errors.New(ErrReadFilepathCode, errors.Alert, []string{}, []string{err.Error()}, []string{}, []string{})
}

func ErrFailNewRequest(err error) error {
	return errors.New(ErrFailNewRequestCode, errors.Alert, []string{"Failed to create a request"}, []string{err.Error()}, []string{}, []string{})
}

func ErrNoProfileName() error {
	return errors.New(ErrNoProfileNameCode, errors.Alert, []string{"please enter a profile-name"}, []string{}, []string{}, []string{})
}

func ErrNoTestURL() error {
	return errors.New(ErrNoTestURLCode, errors.Alert, []string{"please enter a test URL"}, []string{}, []string{}, []string{})
}

func ErrNotValidURL() error {
	return errors.New(ErrNotValidURLCode, errors.Alert, []string{"please enter a valid test URL"}, []string{}, []string{}, []string{})
}

func ErrFailConvReq(err error) error {
	return errors.New(ErrFailConvReqCode, errors.Alert, []string{"failed to convert concurrent-request"}, []string{err.Error()}, []string{}, []string{})
}

func ErrFailConvQPS(err error) error {
	return errors.New(ErrFailConvQPSCode, errors.Alert, []string{"failed to convert qps"}, []string{err.Error()}, []string{}, []string{})
}

func ErrFailMarshal(err error) error {
	return errors.New(ErrFailMarshalCode, errors.Alert, []string{"failed to marshal values"}, []string{err.Error()}, []string{}, []string{})
}

func ErrAttachAuthToken(err error) error {
	return errors.New(ErrAttachAuthTokenCode, errors.Alert, []string{"authentication token not found. please supply a valid user token with the --token (or -t) flag"}, []string{err.Error()}, []string{}, []string{})
}

func ErrFailRequest(err error) error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{"failed to make a request"}, []string{err.Error()}, []string{}, []string{})
}

func ErrFailReqStatus(statusCode int) error {
	return errors.New(ErrFailReqStatusCode, errors.Alert, []string{"Response Status Code ", strconv.Itoa(statusCode), ", possible Server Error"}, []string{}, []string{}, []string{})
}

func ErrFailReadReqBody(err error) error {
	return errors.New(ErrFailReadReqBodyCode, errors.Alert, []string{"failed to read response body"}, []string{err.Error()}, []string{}, []string{})
}

func ErrFailUnmarshal(err error) error {
	return errors.New(ErrFailUnmarshalCode, errors.Alert, []string{"failed to unmarshal response body"}, []string{err.Error()}, []string{}, []string{})
}

func ErrNoProfileFound() error {
	return errors.New(ErrNoProfileFoundCode, errors.Alert, []string{"no profiles found with given name."}, []string{}, []string{}, []string{})
}

func ErrNoResultFound() error {
	return errors.New(ErrNoResultFoundCode, errors.Alert, []string{"results does not exit. Please run a profile test and try again. Use `mesheryctl perf list` to see a list of performance profiles"}, []string{}, []string{}, []string{})
}

func ErrFailTestRun() error {
	return errors.New(ErrFailTestRunCode, errors.Alert, []string{"failed to run test"}, []string{}, []string{}, []string{})
}

func ErrStartTermbox(err error) error {
	return errors.New(ErrStartTermboxCode, errors.Alert, []string{"failed to start termbox"}, []string{err.Error()}, []string{}, []string{})
}

func ErrConvJSONYAML(err error) error {
	return errors.New(ErrConvJSONYAMLCode, errors.Alert, []string{"failed to convert json to yaml"}, []string{err.Error()}, []string{}, []string{})
}

func ErrInvalidOutputChoice() error {
	return errors.New(ErrConvJSONYAMLCode, errors.Alert, []string{"output-format choice invalid, use [json|yaml]"}, []string{}, []string{}, []string{})
}
