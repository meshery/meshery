package relationships

import (
	"fmt"
	"sort"
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
	return fmt.Errorf("Relationship definition is invalid\n\n%s\n\n%s",
		formatRelationshipViolations(details.Violations),
		"Review the schema violations above and update the relationship definition file.",
	)
}

func formatRelationshipViolations(violations []meshkitschema.Violation) string {
	if len(violations) == 0 {
		return "No detailed schema violations were reported by the validator."
	}

	// Group violations by instance path
	groups := make(map[string][]string)
	var paths []string
	for _, v := range violations {
		path := v.InstancePath
		if path == "" {
			path = "/"
		}
		msg := strings.TrimSpace(v.Message)
		if msg == "" {
			msg = "validation failed"
		}
		if _, exists := groups[path]; !exists {
			paths = append(paths, path)
		}
		groups[path] = append(groups[path], msg)
	}
	sort.Strings(paths)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Schema violations (%d):\n", len(violations)))
	for _, path := range paths {
		sb.WriteString(fmt.Sprintf("\n  %s\n", path))
		for _, msg := range groups[path] {
			sb.WriteString(fmt.Sprintf("    - %s\n", msg))
		}
	}

	return sb.String()
}
