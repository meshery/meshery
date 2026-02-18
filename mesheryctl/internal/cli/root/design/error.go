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

package design

import (
	"fmt"
	"strings"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
)

const (
	ErrImportDesignCode               = "mesheryctl-1001"
	ErrInValidSourceCode              = "mesheryctl-1002"
	ErrOnboardDesignCode              = "mesheryctl-1003"
	ErrOffboardDesignCode             = "mesheryctl-1005"
	ErrDesignFlagCode                 = "mesheryctl-1006"
	ErrDesignManifestCode             = "mesheryctl-1007"
	ErrDesignFileNotProvidedCode      = "mesheryctl-1140"
	ErrDesignsNotFoundCode            = "mesheryctl-1037"
	ErrInvalidDesignFileCode          = "mesheryctl-1038"
	ErrPatternSourceTypeCode          = "mesheryctl-1121"
	ErrParseDesignFileCode            = "mesheryctl-1163"
	ErrDeleteDesignCode               = "mesheryctl-1164"
	ErrInvalidCommandCode             = "mesheryctl-1191"
	ErrDesignNameOrIDNotSpecifiedCode = "mesheryctl-1192"
	ErrDesignInvalidApiResponseCode   = "mesheryctl-1199"
)

const (
	errPatternMsg = `Usage: mesheryctl design import -f [file/url] -s [source-type]
Example: mesheryctl design import -f ./pattern.yml -s "Kubernetes Manifest"`

	errOnboardMsg = `Usage: mesheryctl design onboard -f [filepath] -s [source type]
Example: mesheryctl design onboard -f ./pattern.yml -s "Kubernetes Manifest"
Description: Onboard pattern`
	errInvalidPathMsg = "file path %s is invalid. Enter a valid path"
)

func ErrDesignNotFound() error {
	return errors.New(ErrDesignsNotFoundCode, errors.Fatal, []string{"Design Not Found"}, []string{"No Design found with the given name or ID"}, []string{"Design with the given name or ID is not present"}, []string{"Please check if the given design name or ID is present via 'mesheryctl design list'"})
}

func ErrInvalidDesignFile(err error) error {
	return errors.New(ErrInvalidDesignFileCode, errors.Fatal, []string{err.Error()}, []string{"Design appears invalid. Could not parse provided design"}, []string{"Design file provided is not valid"}, []string{"Please check that your design file is a valid yaml file"})
}

func ErrImportDesign(err error) error {
	return errors.New(ErrImportDesignCode, errors.Fatal,
		[]string{"Unable to import design"},
		[]string{err.Error()},
		[]string{"Invalid design file"},
		[]string{"Check your design URL/file path"})
}

func ErrInValidSource(invalidSourceType string, validSourceTypes []string) error {
	return errors.New(ErrInValidSourceCode, errors.Fatal,
		[]string{fmt.Sprintf("Invalid design source type: `%s`", invalidSourceType)},
		[]string{"Invalid design source type was provided"},
		[]string{"Provided design source type (-s) is invalid"},
		[]string{"Ensure you pass a valid source type", fmt.Sprintf("\nAllowed source types: %s", strings.Join(validSourceTypes, ", "))})
}

func ErrDesignManifest() error {
	return errors.New(ErrDesignManifestCode, errors.Alert,
		[]string{"No file path detected"},
		[]string{"No manifest file path detected"},
		[]string{"Manifest path not provided"},
		[]string{"Provide the path to the design manifest. \n\n%v", errPatternMsg})
}

func ErrDesignFileNotProvided() error {
	return errors.New(ErrDesignFileNotProvidedCode, errors.Alert,
		[]string{"Design file path not provided"},
		[]string{"Design [File | File Path | URL] isn't specified"},
		[]string{"[File | File Path | URL] not detected. The '-f' flag is missing or empty in the 'mesheryctl design import' command"},
		[]string{"Provide the path to the design file using the '-f' flag. Ensure that the file path or URL is correctly specified and accessible. \n\n%v", errPatternMsg})
}

func ErrOnboardDesign() error {
	return errors.New(ErrOnboardDesignCode, errors.Alert,
		[]string{"Error Onboarding design"},
		[]string{"Unable to onboard design due to empty path"},
		[]string{"File path or design name not provided."},
		[]string{"Provide a file path/design name. \n\n%v", errOnboardMsg})
}

func ErrDesignSourceType() error {
	return errors.New(ErrPatternSourceTypeCode, errors.Alert,
		[]string{"Source type for the design to import not specified"},
		[]string{"Empty  source type detected"},
		[]string{"Design source type not provided"},
		[]string{"Provide one of the supported source type for the design to import. \n\n%v", errPatternMsg})
}

func ErrViewDesignFlag() error {
	return errors.New(ErrDesignFlagCode, errors.Alert,
		[]string{"Invalid command"},
		[]string{"Wrong use of command flags"},
		[]string{"-a/all flag is being used while an design is specified"},
		[]string{"-a/-all cannot be used when [design name|id] is specified"})
}

func ErrOffboardDesign(err error) error {
	return errors.New(ErrOffboardDesignCode, errors.Alert,
		[]string{"Error Offboarding design"},
		[]string{"Unable to offboard design due to empty path"},
		[]string{"File path or design name not provided."},
		[]string{"Provide a file path/design name. \n\n%v", errOnboardMsg})
}

func ErrParseDesignFile(err error) error {
	return errors.New(ErrParseDesignFileCode, errors.Alert,
		[]string{"Failed to parse design file"},
		[]string{err.Error()},
		[]string{"The design file format is invalid or corrupted", "The YAML/JSON syntax may be incorrect"},
		[]string{"Ensure the design file is a valid Meshery design format", "Check for YAML/JSON syntax errors", "Validate the design file structure"})
}

func ErrDeleteDesign(err error, designName string) error {
	return errors.New(ErrDeleteDesignCode, errors.Alert,
		[]string{"Unable to delete design"},
		[]string{fmt.Sprintf("%s: %s", utils.DesignError(fmt.Sprintf("failed to delete design %s", designName)), err.Error())},
		[]string{"Design may not exist", "Network connection issue", "Meshery server unreachable"},
		[]string{"Verify the design exists using 'mesheryctl design list'", "Check your network connection and Meshery server status"})
}

func ErrInvalidCommand(cmd string, suggestions []string) error {
	var longDesc string
	if len(suggestions) > 0 {
		longDesc = fmt.Sprintf("'%s' is an invalid command. Did you mean one of these?\n\t%s", cmd, strings.Join(suggestions, "\n\t"))
	} else {
		longDesc = fmt.Sprintf("'%s' is an invalid command", cmd)
	}

	return errors.New(ErrInvalidCommandCode, errors.Alert,
		[]string{"Invalid command"},
		[]string{longDesc},
		[]string{"The provided command is not recognized"},
		[]string{"Run 'mesheryctl design --help' to see available commands"},
	)
}

func ErrDesignNameOrIDNotSpecified() error {
	return errors.New(ErrDesignNameOrIDNotSpecifiedCode, errors.Alert,
		[]string{"Design name or ID not specified"},
		[]string{"No design name or ID was provided"},
		[]string{"Command requires a design name or ID as argument"},
		[]string{"Provide a design name or ID, or use '-a' flag to view all designs.\nRun 'mesheryctl design view --help' for usage details"})
}

func ErrDesignInvalidApiResponse(message string) error {
	return errors.New(ErrDesignInvalidApiResponseCode, errors.Alert,
		[]string{"Invalid API response"},
		[]string{message},
		[]string{"The API response is missing expected fields or has an unexpected format"},
		[]string{"Ensure the Meshery server is running a compatible version", "Check for any issues with the Meshery server that may cause it to return malformed responses"})
}
