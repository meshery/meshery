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

package pattern

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrImportPatternCode          = "mesheryctl-1001"
	ErrInValidSourceCode          = "mesheryctl-1002"
	ErrOnboardPatternCode         = "mesheryctl-1003"
	ErrPatternFoundCode           = "mesheryctl-1004"
	ErrOffboardPatternCode        = "mesheryctl-1005"
	ErrPatternFlagCode            = "mesheryctl-1006"
	ErrPatternManifestCode        = "mesheryctl-1007"
	ErrPatternSourceTypeCode      = "mesheryctl-1121"
	ErrPatternsNotFoundCode       = "mesheryctl-1037"
	ErrInvalidPatternFileCode     = "mesheryctl-1038"
	ErrPatternInvalidNameOrIDCode = "mesheryctl-1039"
	ErrCopyDataCode               = "mesheryctl-1122"
	ErrCreateFileCode             = "mesheryctl-1123"
	ErrRetrieveHomeDirCode        = "mesheryctl-1124"
	ErrReadFromBodyCode           = "mesheryctl-1125"
	ErrMarkFlagRequireCode        = "mesheryctl-1126"
)

const (
	errPatternMsg = `Usage: mesheryctl pattern import -f [file/url] -s [source-type]
Example: mesheryctl pattern import -f ./pattern.yml -s "Kubernetes Manifest"`

	errOnboardMsg = `Usage: mesheryctl pattern onboard -f [filepath] -s [source type]
Example: mesheryctl pattern onboard -f ./pattern.yml -s "Kubernetes Manifest"
Description: Onboard pattern`
)

func ErrPatternNotFound() error {
	return errors.New(ErrPatternsNotFoundCode, errors.Fatal, []string{"Pattern Not Found"}, []string{"No Pattern found with the given name or ID"}, []string{"Pattern with the given name or ID is not present"}, []string{"Please check if the given pattern name or ID is present via 'mesheryctl pattern list'"})
}

func ErrInvalidPatternFile(err error) error {
	return errors.New(ErrInvalidPatternFileCode, errors.Fatal, []string{err.Error()}, []string{"Pattern appears invalid. Could not parse successfully"}, []string{"Pattern file provided is not valid"}, []string{"Please check that your pattern file is a valid yaml file"})
}

func ErrPatternInvalidNameOrID(err error) error {
	return errors.New(
		ErrPatternInvalidNameOrIDCode,
		errors.Alert,
		[]string{"Unable to fetch Pattern"},
		[]string{err.Error()},
		[]string{"Invalid pattern name or ID"},
		[]string{"Run `mesheryctl pattern view -a` to view all available patterns."})
}

func ErrImportPattern(err error) error {
	return errors.New(ErrImportPatternCode, errors.Fatal,
		[]string{"Unable to import pattern"},
		[]string{err.Error()},
		[]string{"Invalid pattern file"},
		[]string{"Check your pattern URL/file path"})
}

func ErrInValidSource(invalidSourceType string, validSourceTypes []string) error {
	return errors.New(ErrInValidSourceCode, errors.Fatal,
		[]string{fmt.Sprintf("Invalid pattern source type: `%s`", invalidSourceType)},
		[]string{"Invalid pattern source type due to wrong type/passing"},
		[]string{"pattern source type (-s) is invalid or not passed."},
		[]string{"Ensure you pass a valid source type. \nAllowed source types: %s", strings.Join(validSourceTypes, ", ")})
}

func ErrPatternManifest() error {
	return errors.New(ErrPatternManifestCode, errors.Alert,
		[]string{"No file path detected"},
		[]string{"No manifest file path detected"},
		[]string{"Manifest path not provided"},
		[]string{"Provide the path to the pattern manifest. \n\n%v", errPatternMsg})
}

func ErrPatternSourceType() error {
	return errors.New(ErrPatternSourceTypeCode, errors.Alert,
		[]string{"Source type for the design to import not specified"},
		[]string{"Empty  source type detected"},
		[]string{"Design source type not provided"},
		[]string{"Provide one of the supported source type for the design to import. \n\n%v", errPatternMsg})
}

func ErrOnboardPattern() error {
	return errors.New(ErrOnboardPatternCode, errors.Alert,
		[]string{"Error Onboarding pattern"},
		[]string{"Unable to onboard pattern due to empty path"},
		[]string{"File path or pattern name not provided."},
		[]string{"Provide a file path/pattern name. \n\n%v", errOnboardMsg})
}

func ErrPatternFound() error {
	return errors.New(ErrPatternFoundCode, errors.Alert,
		[]string{"pattern not found"},
		[]string{"No pattern found with a given name"},
		[]string{"pattern name not provided"},
		[]string{"Provide an pattern name. Use `mesheryctl pattern list`to display list of patternlications"},
	)
}

func ErrViewPatternFlag() error {
	return errors.New(ErrPatternFlagCode, errors.Alert,
		[]string{"Invalid command"},
		[]string{"Wrong use of command flags"},
		[]string{"-a/all flag is being used while an pattern is specified"},
		[]string{"-a/-all cannot be used when [pattern name|id] is specified"})
}

func ErrOffboardPattern(err error) error {
	return errors.New(ErrOffboardPatternCode, errors.Alert,
		[]string{"Error Offboarding pattern"},
		[]string{"Unable to offboard pattern due to empty path"},
		[]string{"File path or pattern name not provided."},
		[]string{"Provide a file path/pattern name. \n\n%v", errOnboardMsg})
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
