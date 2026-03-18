// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrFetchFilterCode          = "mesheryctl-1019"
	ErrDeleteFilterCode         = "mesheryctl-1018"
	ErrFilterNotFoundCode       = "mesheryctl-1167"
	ErrFilterURIRequiredCode    = "mesheryctl-1168"
	ErrViewAllWithNameCode      = "mesheryctl-1169"
	ErrMultiWordFilterNameCode  = "mesheryctl-1170"
	ErrInvalidFilterCommandCode = "mesheryctl-1171"
)

const (
	errFilterNameOrIDNotProvided = "filter name or ID not provided"
)

func formatErrorWithReference(subCmdUsed string) string {
	baseURL := "reference/mesheryctl/filter"
	switch subCmdUsed {
	case "import":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/import")
	case "view":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/view")
	case "delete":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/delete")
	case "list":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/list")

	}
	return fmt.Sprintf("\nSee %s for usage details\n", baseURL)
}

func ErrFetchFilter(err error) error {
	return errors.New(
		ErrFetchFilterCode,
		errors.Alert,
		[]string{"Unable to Fetch Filter"},
		[]string{err.Error()},
		[]string{"Filter name or id doesn't exist"},
		[]string{"Run `mesheryctl filter view -a` to view all available filters."})
}

func ErrDeleteFilter(err error, filterName string, subCmdUsed string) error {
	return errors.New(ErrDeleteFilterCode, errors.Alert,
		[]string{"Failed to delete filter"},
		[]string{fmt.Sprintf("failed to delete filter %s", filterName), err.Error()},
		[]string{"Unable to delete the specified filter"},
		[]string{"Verify the filter exists using `mesheryctl filter list`", formatErrorWithReference(subCmdUsed)})
}

func ErrFilterNotFound(filterNameOrID string, subCmdUsed string) error {
	return errors.New(ErrFilterNotFoundCode, errors.Alert,
		[]string{"Filter not found"},
		[]string{fmt.Sprintf("filter with name or ID having prefix %s does not exist", filterNameOrID)},
		[]string{"The specified filter does not exist"},
		[]string{"Use `mesheryctl filter list` to see available filters", formatErrorWithReference(subCmdUsed)})
}

func ErrFilterURIRequired(subCmdUsed string) error {
	return errors.New(ErrFilterURIRequiredCode, errors.Alert,
		[]string{"URI is required"},
		[]string{"No URI provided for filter import"},
		[]string{"Filter import command requires a URI"},
		[]string{"Provide a valid file path or URL", formatErrorWithReference(subCmdUsed)})
}

func ErrViewAllWithName(subCmdUsed string) error {
	return errors.New(ErrViewAllWithNameCode, errors.Alert,
		[]string{"Invalid flag combination"},
		[]string{"--all cannot be used when filter name or ID is specified"},
		[]string{"Conflicting options provided"},
		[]string{"Use either --all flag OR specify a filter name/ID", formatErrorWithReference(subCmdUsed)})
}

func ErrMultiWordFilterName(subCmdUsed string) error {
	return errors.New(ErrMultiWordFilterNameCode, errors.Alert,
		[]string{"Invalid filter name format"},
		[]string{"multi-word filter names must be enclosed in double quotes"},
		[]string{"Filter name with spaces requires quoting"},
		[]string{"Enclose multi-word filter names in double quotes", formatErrorWithReference(subCmdUsed)})
}

func ErrInvalidFilterCommand(invalidCmd string, subCmdUsed string) error {
	return errors.New(ErrInvalidFilterCommandCode, errors.Alert,
		[]string{"Invalid command"},
		[]string{fmt.Sprintf("'%s' is an invalid command", invalidCmd)},
		[]string{"The provided subcommand is not recognized"},
		[]string{"Run `mesheryctl filter --help` to see available commands", formatErrorWithReference(subCmdUsed)})
}
