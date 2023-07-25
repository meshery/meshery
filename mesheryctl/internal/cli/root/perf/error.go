package perf

import (
	"fmt"
	"strconv"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrMesheryConfigCode         = "1027"
	ErrReadFilepathCode          = "1028"
	ErrNoProfileNameCode         = "1029"
	ErrNoTestURLCode             = "1030"
	ErrNotValidURLCode           = "1031"
	ErrFailMarshalCode           = "1032"
	ErrAttachAuthTokenCode       = "1033"
	ErrFailRequestCode           = "1034"
	ErrFailReqStatusCode         = "1035"
	ErrFailUnmarshalCode         = "1036"
	ErrNoProfileFoundCode        = "1037"
	ErrFailTestRunCode           = "1038"
	ErrInvalidOutputChoiceCode   = "1039"
	ErrUnauthenticatedCode       = "1040"
	ErrFailUnmarshalFileCode     = "1041"
	ErrInvalidTestConfigFileCode = "1042"
	ErrArgumentOverflowCode      = "1043"
)

func ErrMesheryConfig(err error) error {
	return errors.New(ErrMesheryConfigCode, errors.Alert, []string{},
		[]string{"error processing config", err.Error(), formatErrorWithReference()}, []string{}, []string{})
}

func ErrReadFilepath(err error) error {
	return errors.New(ErrReadFilepathCode, errors.Alert, []string{},
		[]string{err.Error(), formatErrorWithReference()}, []string{}, []string{})
}

func ErrNoProfileName() error {
	return errors.New(ErrNoProfileNameCode, errors.Alert, []string{},
		[]string{"please enter a profile-name", formatErrorWithReference()}, []string{}, []string{})
}

func ErrNoTestURL() error {
	return errors.New(ErrNoTestURLCode, errors.Alert, []string{},
		[]string{"please enter a test URL", formatErrorWithReference()}, []string{}, []string{})
}

func ErrNotValidURL() error {
	return errors.New(ErrNotValidURLCode, errors.Alert, []string{},
		[]string{"please enter a valid test URL", formatErrorWithReference()}, []string{}, []string{})
}

func ErrFailMarshal(err error) error {
	return errors.New(ErrFailMarshalCode, errors.Alert, []string{},
		[]string{"failed to marshal values", err.Error(), formatErrorWithReference()}, []string{}, []string{})
}

func ErrAttachAuthToken(err error) error {
	return errors.New(ErrAttachAuthTokenCode, errors.Alert, []string{err.Error()},
		[]string{"authentication token not found. Run `mesheryctl system login` to generate a new token or supply a valid token with the `--token` flag.", formatErrorWithReference()}, []string{}, []string{})
}

func ErrFailRequest(err error) error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{},
		[]string{"Failed to make a request", err.Error(), formatErrorWithReference()}, []string{}, []string{})
}

func ErrUnauthenticated() error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{},
		[]string{"Authentication token is invalid. Please supply a valid user token. Login with `mesheryctl system login`", formatErrorWithReference()}, []string{}, []string{})
}

func ErrExpired() error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{},
		[]string{"Authentication token is expired. Please login with `mesheryctl system login` to generate a new token", formatErrorWithReference()}, []string{}, []string{})
}

func ErrFailReqStatus(statusCode int) error {
	return errors.New(ErrFailReqStatusCode, errors.Alert, []string{},
		[]string{"Response Status Code " + strconv.Itoa(statusCode) + ", possible Server Error", formatErrorWithReference()}, []string{}, []string{})
}

func ErrFailUnmarshal(err error) error {
	return errors.New(ErrFailUnmarshalCode, errors.Alert, []string{},
		[]string{"failed to unmarshal response body", err.Error(), formatErrorWithReference()}, []string{}, []string{})
}

func ErrNoProfileFound() error {
	return errors.New(ErrNoProfileFoundCode, errors.Alert, []string{},
		[]string{"no profiles found with given name", formatErrorWithReference()}, []string{}, []string{})
}

func ErrFailTestRun() error {
	return errors.New(ErrFailTestRunCode, errors.Alert, []string{},
		[]string{"failed to run test", formatErrorWithReference()}, []string{}, []string{})
}

func ErrInvalidOutputChoice() error {
	return errors.New(ErrInvalidOutputChoiceCode, errors.Alert, []string{}, []string{"output-format choice invalid, use [json|yaml]", formatErrorWithReference()}, []string{}, []string{})
}

func ErrFailUnmarshalFile(err error) error {
	return errors.New(ErrFailUnmarshalFileCode, errors.Alert, []string{},
		[]string{"failed to unmarshal configuration file", err.Error(), formatErrorWithReference()}, []string{}, []string{})
}

func ErrInvalidTestConfigFile() error {
	return errors.New(ErrInvalidTestConfigFileCode, errors.Alert, []string{},
		[]string{"invalid test configuration file", formatErrorWithReference()}, []string{"the test configuration is outdated or incorrect"}, []string{"see https://docs.meshery.io/guides/performance-management#running-performance-benchmarks-through-mesheryctl for a valid configuration file"})
}

func formatErrorWithReference() string {
	baseURL := "https://docs.meshery.io/reference/mesheryctl/perf"
	switch cmdUsed {
	case "apply":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/apply")
	case "profile":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/profile")
	case "result":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/result")
	}
	return fmt.Sprintf("\nSee %s for usage details\n", baseURL)
}

func ErrorArgumentOverflow() error {
	return errors.New(ErrArgumentOverflowCode, errors.Alert, []string{},
		[]string{"Invalid number of arguments", formatErrorWithReference()}, []string{}, []string{"Use 'mesheryctl --help' to display usage guide."})
}
