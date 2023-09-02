package perf

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrMesheryConfigCode         = "1033"
	ErrReadFilepathCode          = "1034"
	ErrNoProfileNameCode         = "1035"
	ErrNoTestURLCode             = "1036"
	ErrNotValidURLCode           = "1037"
	ErrFailMarshalCode           = "1038"
	ErrFailUnmarshalCode         = "1039"
	ErrNoProfileFoundCode        = "1040"
	ErrFailTestRunCode           = "1041"
	ErrInvalidOutputChoiceCode   = "1042"
	ErrUnauthenticatedCode       = "1043"
	ErrFailUnmarshalFileCode     = "1044"
	ErrInvalidTestConfigFileCode = "1045"
	ErrArgumentOverflowCode      = "1046"
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
