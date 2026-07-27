package format

import "github.com/meshery/meshkit/errors"

var (
	ErrOutputToJsonCode = "mesheryctl-1147"
	ErrOutputToYamlCode = "mesheryctl-1158"
	ErrOutputToToonCode = "mesheryctl-1253"
)

func ErrOutputToJson() error {
	return errors.New(ErrOutputToJsonCode, errors.Alert,
		[]string{"Failed to format output in JSON"},
		[]string{"Error occurred while formatting to JSON"},
		[]string{"The content provided for formatting failed."},
		[]string{"Check the JSON structure you are providing for formatting."})
}

func ErrOutputToYaml() error {
	return errors.New(ErrOutputToYamlCode, errors.Alert,
		[]string{"Failed to format output in YAML"},
		[]string{"Error occurred while formatting to YAML"},
		[]string{"The content provided for formatting failed."},
		[]string{"Check the Yaml structure you are providing for formatting."})
}

// ErrOutputToToon is the error for failing to format output in TOON
func ErrOutputToToon() error {
	return errors.New(ErrOutputToToonCode, errors.Alert,
		[]string{"Failed to format output in TOON"},
		[]string{"Error occurred while formatting to TOON"},
		[]string{"The content provided for formatting failed."},
		[]string{"Check the TOON structure you are providing for formatting."})
}
