package perf

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrReadFilepathCode             = "1056"
	ErrNoProfileNameCode            = "1057"
	ErrNoTestURLCode                = "1058"
	ErrNotValidURLCode              = "1059"
	ErrFailMarshalCode              = "1060"
	ErrFailUnmarshalCode            = "1061"
	ErrNoProfileFoundCode           = "1062"
	ErrInvalidOutputChoiceCode      = "1063"
	ErrUnauthenticatedCode          = "1064"
	ErrFailUnmarshalFileCode        = "1065"
	ErrInvalidTestConfigFileCode    = "1066"
	ErrArgumentOverflowCode         = "1067"
	ErrInvalidJSONFileCode          = "1068"
	ErrHealthCheckerCode            = "1069"
	ErrPerformanceProfileResultCode = "1070"
)

func ErrReadFilepath(err error) error {
	return errors.New(ErrReadFilepathCode, errors.Alert,
		[]string{"Unable to read file"},
		[]string{err.Error()},
		[]string{"The provided file is not present or has an invalid path"},
		[]string{"Provide a valid file path", formatErrorWithReference()})
}

func ErrNoProfileName() error {
	return errors.New(ErrNoProfileNameCode, errors.Alert,
		[]string{"No profile name"},
		[]string{"No profile name provided"},
		[]string{"No profile name provided to perform test"},
		[]string{"Provide a profile-name", formatErrorWithReference()})
}

func ErrNoTestURL() error {
	return errors.New(ErrNoTestURLCode, errors.Alert,
		[]string{"No URL"},
		[]string{"Unable to get URL for performing test"},
		[]string{"No URL provided for performing tests"},
		[]string{"Provide a test URL", formatErrorWithReference()})
}

func ErrNotValidURL() error {
	return errors.New(ErrNotValidURLCode, errors.Alert,
		[]string{"Invalid URL"},
		[]string{"Invalid URL for performing test"},
		[]string{"Invalid test URL provided"},
		[]string{"Input a valid test URL", formatErrorWithReference()})
}

func ErrFailMarshal(err error) error {
	return errors.New(ErrFailMarshalCode, errors.Alert,
		[]string{"Unable to convert JSON object."},
		[]string{err.Error()},
		[]string{"Invalid JSON Format"},
		[]string{"Check the data structure you are providing for marshalling", formatErrorWithReference()})
}

func ErrFailUnmarshal(err error) error {
	return errors.New(ErrFailUnmarshalCode, errors.Alert,
		[]string{"Error unmarshalling response"},
		[]string{"Failed to unmarshal response body", err.Error()},
		[]string{"The JSON format from the response body is not valid."},
		[]string{"Ensure a valid JSON is provided for processing.", formatErrorWithReference()})
}

func ErrNoProfileFound() error {
	return errors.New(ErrNoProfileFoundCode, errors.Alert,
		[]string{"No profile found"},
		[]string{"No profile found with given name"},
		[]string{"Invalid profile name or no profile name provided"},
		[]string{"Ensure a valid profile name is inputted. See https://docs.meshery.io/reference/mesheryctl/perf/apply for more details"})
}

func ErrInvalidOutputChoice() error {
	return errors.New(ErrInvalidOutputChoiceCode, errors.Alert,
		[]string{"Invalid output format choice"},
		[]string{"Output format choice is invalid, use [json|yaml]"},
		[]string{"Invalid JSON or YAML content"},
		[]string{"Check the JSON or YAML structure.", formatErrorWithReference()})
}

func ErrFailUnmarshalFile(err error) error {
	return errors.New(ErrFailUnmarshalFileCode, errors.Alert,
		[]string{"Failed to unmarshal configuration file"},
		[]string{err.Error()},
		[]string{"Unable to covert JSON format from yml file"},
		[]string{"Ensure a valid data is provided in your yml file", formatErrorWithReference()})
}

func ErrInvalidTestConfigFile() error {
	return errors.New(ErrInvalidTestConfigFileCode,
		errors.Alert,
		[]string{"Invalid config file"},
		[]string{"Invalid test configuration file"},
		[]string{"The test configuration is outdated or incorrect"},
		[]string{"See https://docs.meshery.io/guides/performance-management#running-performance-benchmarks-through-mesheryctl for a valid configuration file"})
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
	return errors.New(ErrArgumentOverflowCode,
		errors.Alert,
		[]string{"Invalid arguments"},
		[]string{"Invalid number of arguments"},
		[]string{"Too many arguments passed"},
		[]string{"Use `mesheryctl perf --help` to display usage guide.", formatErrorWithReference()})
}

func ErrInvalidJSONFile() error {
	return errors.New(ErrInvalidJSONFileCode, errors.Alert,
		[]string{"Invalid JSON"},
		[]string{"Unable to create perf profile"},
		[]string{"Invalid JSON passed as options in file"},
		[]string{"Provide a valid JSON object", formatErrorWithReference()})
}

func ErrHealthChecker(err error) error {
	return errors.New(ErrHealthCheckerCode,
		errors.Alert,
		[]string{"Healthchecks Failed"},
		[]string{err.Error()},
		[]string{"Failed to initialize healthchecker"},
		[]string{"Ensure Mesheryctl is running and has the right configurations. See https://docs.meshery.io/reference/mesheryctl/system/check for checking Meshery components"})
}

func ErrPerformanceProfileResult(err error) error {
	return errors.New(ErrPerformanceProfileResultCode, errors.Alert,
		[]string{"Error displaying performance"},
		[]string{"Unable to display performance Profile" + err.Error()},
		[]string{"Failed to fetch results for a specific profile"},
		[]string{"Check your network connection and ensure Meshery is running .", formatErrorWithReference()})
}
