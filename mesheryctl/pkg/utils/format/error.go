// Package format provides CLI commands and utilities for mesheryctl.
package format

import "github.com/meshery/meshkit/errors"

var (
	ErrOutputToJSONCode = "mesheryctl-1147"
	ErrOutputToYamlCode = "mesheryctl-1158"
)

func ErrOutputToJSON() error {
	return errors.New(ErrOutputToJSONCode, errors.Alert,
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
