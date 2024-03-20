// Copyright 2023 Layer5, Inc.
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
	ErrImportpatternCode          = "mesheryctl-1001"
	ErrInValidSourceCode          = "mesheryctl-1002"
	ErrOnboardpatternCode         = "mesheryctl-1003"
	ErrpatternFoundCode           = "mesheryctl-1004"
	ErrInvalidpatternNameOrIDCode = "mesheryctl-1005"
	ErrpatternFlagCode            = "mesheryctl-1006"
	ErrpatternManifestCode        = "mesheryctl-1007"
	ErrPatternsNotFoundCode       = "mesheryctl-1037"
	ErrInvalidPatternFileCode     = "mesheryctl-1038"
	ErrPatternInvalidNameOrIDCode = "mesheryctl-1039"
)

const (
	errpatternMsg = `Usage: mesheryctl pattern import -f [file/url] -s [source-type]
Example: mesheryctl pattern import -f ./patternlication.yml -s "Kubernetes Manifest"`

	errOnboardMsg = `Usage: mesheryctl pattern onboard -f [filepath] -s [source type]
Example: mesheryctl pattern onboard -f ./patternlication.yml -s "Kubernetes Manifest"
Description: Onboard patternlication`
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

func ErrImportpattern(err error) error {
	return errors.New(ErrImportpatternCode, errors.Fatal,
		[]string{"Unable to import pattern"},
		[]string{err.Error()},
		[]string{"Invalid patternlication file"},
		[]string{"Check your patternlication URL/file path"})
}

func ErrInValidSource(invalidSourceType string, validSourceTypes []string) error {
	return errors.New(ErrInValidSourceCode, errors.Fatal,
		[]string{fmt.Sprintf("Invalid patternlication source type: `%s`", invalidSourceType)},
		[]string{"Invalid patternlication source type due to wrong type/passing"},
		[]string{"patternlication source type (-s) is invalid or not passed."},
		[]string{"Ensure you pass a valid source type. \nAllowed source types: %s", strings.Join(validSourceTypes, ", ")})
}

func ErrpatternManifest() error {
	return errors.New(ErrpatternManifestCode, errors.Alert,
		[]string{"No file path detected"},
		[]string{"No manifest file path detected"},
		[]string{"Manifest path not provided"},
		[]string{"Provide the path to the pattern manifest. \n\n%v", errpatternMsg})
}

func ErrOnboardpattern() error {
	return errors.New(ErrOnboardpatternCode, errors.Alert,
		[]string{"Error Onboarding pattern"},
		[]string{"Unable to onboard pattern due to empty path"},
		[]string{"File path or patternlication name not provided."},
		[]string{"Provide a file path/pattern name. \n\n%v", errOnboardMsg})
}

func ErrpatternFound() error {
	return errors.New(ErrpatternFoundCode, errors.Alert,
		[]string{"pattern not found"},
		[]string{"No pattern found with a given name"},
		[]string{"patternlication name not provided"},
		[]string{"Provide an pattern name. Use `mesheryctl pattern list`to display list of patternlications"},
	)
}

func ErrInvalidpatternNameOrID(err error) error {
	return errors.New(ErrInvalidpatternNameOrIDCode, errors.Alert,
		[]string{"Invalid patternlication"},
		[]string{"Failed to get patternlication based on input"},
		[]string{"patternlication name|id is invalid"},
		[]string{"Run `mesheryctl pattern view --all` to view all patternlications"})
}

func ErrViewpatternFlag() error {
	return errors.New(ErrpatternFlagCode, errors.Alert,
		[]string{"Invalid command"},
		[]string{"Wrong use of command flags"},
		[]string{"-a/all flag is being used while an pattern is specified"},
		[]string{"-a/-all cannot be used when [patternlication name|id] is specified"})
}
