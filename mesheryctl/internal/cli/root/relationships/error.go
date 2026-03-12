package relationships

import (
	"fmt"
	"strings"

	"github.com/meshery/meshkit/errors"
	meshkitschema "github.com/meshery/meshkit/schema"
)

var (
	ErrEmptySheetDataCode               = "mesheryctl-1204"
	ErrRelationshipValidationFailedCode = "mesheryctl-1205"
)

const (
	errInvalidArg = "only one argument must be provided and needs to be enclosed by double quotes if it contains spaces (eg. \"model name\", modelName)"
)

func ErrEmptySheetData(err error) error {
	return errors.New(
		ErrEmptySheetDataCode,
		errors.Alert,
		[]string{"Invalid or empty spreadsheet provided"},
		[]string{err.Error()},
		[]string{"The spreadsheet must contain valid relationship entries with the required headers."},
		[]string{"Ensure the spreadsheet contains at least two headers and corresponding values for it's headers. Refer to the Meshery relationship csv for the expected spreadsheet format \"https://github.com/meshery/meshery/blob/master/mesheryctl/templates/template-csvs/Relationships.csv\""},
	)
}

func ErrRelationshipValidationFailed(details meshkitschema.ValidationDetails) error {
	probableCauses := []string{"The relationship definition violates the Meshery relationship schema"}
	if len(details.Violations) > 0 {
		probableCauses = append(probableCauses, formatRelationshipViolations(details.Violations)...)
	}

	return errors.New(
		ErrRelationshipValidationFailedCode,
		errors.Alert,
		[]string{"Relationship validation failed"},
		[]string{"The supplied relationship definition is invalid"},
		probableCauses,
		[]string{"Review the reported schema violations and update the relationship definition file"},
	)
}

func formatRelationshipViolations(violations []meshkitschema.Violation) []string {
	formatted := make([]string, 0, len(violations))
	for _, violation := range violations {
		location := violation.InstancePath
		if location == "" {
			location = "/"
		}

		message := strings.TrimSpace(violation.Message)
		if message == "" {
			message = fmt.Sprintf("Validation failed for %s", location)
		}

		formatted = append(formatted, fmt.Sprintf(" Violation at %s: %s", location, message))
	}

	return formatted
}
