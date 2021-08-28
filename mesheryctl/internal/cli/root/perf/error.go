package perf

import (
	"fmt"
	"strconv"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrMesheryConfigCode       = "replace me"
	ErrReadFilepathCode        = "replace me"
	ErrNoProfileNameCode       = "replace me"
	ErrNoTestURLCode           = "replace me"
	ErrNotValidURLCode         = "replace me"
	ErrFailMarshalCode         = "replace me"
	ErrAttachAuthTokenCode     = "replace me"
	ErrFailRequestCode         = "replace me"
	ErrFailReqStatusCode       = "replace me"
	ErrFailUnmarshalCode       = "replace me"
	ErrNoProfileFoundCode      = "replace me"
	ErrFailTestRunCode         = "replace me"
	ErrInvalidOutputChoiceCode = "replace me"
	ErrUnauthenticatedCode     = "replace me"
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
		[]string{"authentication token not found. please supply a valid user token with the --token (or -t) flag. or login with `mesheryctl system login`", formatErrorWithReference()}, []string{}, []string{})
}

func ErrFailRequest(err error) error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{},
		[]string{"failed to make a request", err.Error(), formatErrorWithReference()}, []string{}, []string{})
}

func ErrUnauthenticated() error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{},
		[]string{"Invalid/expired authetication token", formatErrorWithReference()}, []string{}, []string{})

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
