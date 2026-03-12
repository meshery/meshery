package relationships

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/meshery/meshkit/errors"
	meshkitschema "github.com/meshery/meshkit/schema"
)

var (
	ErrEmptySheetDataCode = "mesheryctl-1204"
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
<<<<<<< Updated upstream
	return fmt.Errorf("Relationship definition is invalid\n\n%s\n\n%s",
		formatRelationshipViolations(details.Violations),
		"Review the schema violations above and update the relationship definition file.",
			errors.Alert,
			[]string{"Relationship validation failed"},
			[]string{"The supplied relationship definition is invalid"},
			[]string{"The relationship definition violates the Meshery relationship schema"},
			[]string{"Review the reported schema violations and update the relationship definition file"},
		),
	}
}

func ErrRelationshipValidationFailedForFile(file string, details meshkitschema.ValidationDetails) error {
	err := ErrRelationshipValidationFailed(details)
	validationErr, ok := err.(relationshipValidationError)
	if !ok {
		return err
	}

	file = strings.TrimSpace(file)
	if file != "" {
		validationErr.file = filepath.Base(file)
	}
	return validationErr
}

type relationshipValidationError struct {
	file    string
	details meshkitschema.ValidationDetails
	cause   error
}

func (e relationshipValidationError) Error() string {
	var builder strings.Builder

	builder.WriteString("Relationship validation failed")
	if summary := formatRelationshipValidationSummary(e.file, len(e.details.Violations)); summary != "" {
		builder.WriteString(": ")
		builder.WriteString(summary)
	}

	if len(e.details.Violations) == 0 {
		return builder.String()
	}

	builder.WriteString("\n\nViolations:")
	for index, violation := range e.details.Violations {
		fmt.Fprintf(&builder, "\n  %d. %s", index+1, formatRelationshipViolation(violation))
	}

	return builder.String()
}

func (e relationshipValidationError) Unwrap() error {
	return e.cause
}

func formatRelationshipValidationSummary(file string, count int) string {
	var builder strings.Builder

	switch count {
	case 0:
		builder.WriteString("no schema violations reported")
	case 1:
		builder.WriteString("1 schema violation found")
	default:
		fmt.Fprintf(&builder, "%d schema violations found", count)
	}

	if file != "" {
		builder.WriteString(" in ")
		builder.WriteString(file)
	}

	return builder.String()
}

func formatRelationshipViolation(violation meshkitschema.Violation) string {
	location := strings.TrimSpace(violation.InstancePath)
	if location == "" || location == "/" {
		location = "document"
	}

	message := strings.TrimSpace(violation.Message)
	if message == "" {
		message = fmt.Sprintf("validation failed for %s", location)
	}

	return fmt.Sprintf("%s: %s", location, message)
>>>>>>> Stashed changes
}
