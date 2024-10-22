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

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrImportDesignCode          = "mesheryctl-1001"
	ErrInValidSourceCode         = "mesheryctl-1002"
	ErrOnboardDesignCode         = "mesheryctl-1003"
	ErrOffboardDesignCode        = "mesheryctl-1005"
	ErrDesignFlagCode            = "mesheryctl-1006"
	ErrDesignManifestCode        = "mesheryctl-1007"
	ErrDesignsNotFoundCode       = "mesheryctl-1037"
	ErrInvalidDesignFileCode     = "mesheryctl-1038"
	ErrDesignInvalidNameOrIDCode = "mesheryctl-1039"
	ErrPatternSourceTypeCode     = "mesheryctl-1121"
	ErrCopyDataCode              = "mesheryctl-1122"
	ErrCreateFileCode            = "mesheryctl-1123"
	ErrRetrieveHomeDirCode       = "mesheryctl-1124"
	ErrReadFromBodyCode          = "mesheryctl-1125"
	ErrMarkFlagRequireCode       = "mesheryctl-1126"
)

const (
	errPatternMsg = `Usage: mesheryctl design import -f [file/url] -s [source-type]
Example: mesheryctl design import -f ./pattern.yml -s "Kubernetes Manifest"`

	errOnboardMsg = `Usage: mesheryctl design onboard -f [filepath] -s [source type]
Example: mesheryctl design onboard -f ./pattern.yml -s "Kubernetes Manifest"
Description: Onboard pattern`
)

func ErrDesignNotFound() error {
	return errors.New(ErrDesignsNotFoundCode, errors.Fatal, []string{"Design Not Found"}, []string{"No Design found with the given name or ID"}, []string{"Design with the given name or ID is not present"}, []string{"Please check if the given design name or ID is present via 'mesheryctl design list'"})
}

func ErrInvalidDesignFile(err error) error {
	return errors.New(ErrInvalidDesignFileCode, errors.Fatal, []string{err.Error()}, []string{"Design appears invalid. Could not parse successfully"}, []string{"Design file provided is not valid"}, []string{"Please check that your design file is a valid yaml file"})
}

func ErrPatternInvalidNameOrID(err error) error {
	return errors.New(
		ErrDesignInvalidNameOrIDCode,
		errors.Alert,
		[]string{"Unable to fetch Design"},
		[]string{err.Error()},
		[]string{"Invalid design name or ID"},
		[]string{"Run `mesheryctl design view -a` to view all available designs."})
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
		[]string{"Invalid design source type due to wrong type/passing"},
		[]string{"design source type (-s) is invalid or not passed."},
		[]string{"Ensure you pass a valid source type. \nAllowed source types: %s", strings.Join(validSourceTypes, ", ")})
}

func ErrDesignManifest() error {
	return errors.New(ErrDesignManifestCode, errors.Alert,
		[]string{"No file path detected"},
		[]string{"No manifest file path detected"},
		[]string{"Manifest path not provided"},
		[]string{"Provide the path to the design manifest. \n\n%v", errPatternMsg})
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
func ErrCopyData(filepath string, err error) error {
	return errors.New(ErrCopyDataCode, errors.Alert,
		[]string{"Error copying data to file"},
		[]string{fmt.Sprintf("Failed to copy data to the file at path: %s", filepath), err.Error()},
		[]string{"Insufficient disk space, or file system errors."},
		[]string{"Check for sufficient disk space, and verify the integrity of the file system."})
}
func ErrCreateFile(filepath string, err error) error {
	return errors.New(ErrCreateFileCode, errors.Alert,
		[]string{"Error creating file"},
		[]string{fmt.Sprintf("Failed to create the file at path: %s", filepath), err.Error()},
		[]string{"Insufficient disk page, filepath could be invalid."},
		[]string{"Verify that the file path is valid, and ensure there is sufficient disk space available."})
}
func ErrRetrieveHomeDir(err error) error {
	return errors.New(ErrRetrieveHomeDirCode, errors.Alert,
		[]string{"Error retrieving user home/root directory"},
		[]string{"Failed to retrieve the home/root directory,", err.Error()},
		[]string{"Operating system environment issue or insufficient permissions."},
		[]string{"Ensure that the operating system environment is set up correctly and run the application with elevated privileges."})
}
func ErrMarkFlagRequire(flagName string, err error) error {
	return errors.New(ErrMarkFlagRequireCode, errors.Alert,
		[]string{fmt.Sprintf("Failed to mark the flag '%s' as required", flagName)},
		[]string{err.Error()},
		[]string{"The flag may not exist or there was some error while specifying the flag."},
		[]string{"Please ensure that the required flag '%s' is correctly specified and set before running the command."})
}
func ErrReadFromBody(err error) error {
	return errors.New(ErrReadFromBodyCode, errors.Alert,
		[]string{"Unable to read data from the response body"},
		[]string{err.Error()},
		[]string{"The data for the pattern (design) file might be corrupted."},
		[]string{"Please ensure that your network connection is stable. If the issue continues, check the server response or data format for potential problems."})
}
