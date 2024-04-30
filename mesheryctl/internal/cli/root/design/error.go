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
	ErrDesignFoundCode           = "mesheryctl-1004"
	ErrOffboardDesignCode        = "mesheryctl-1005"
	ErrDesignFlagCode            = "mesheryctl-1006"
	ErrDesignManifestCode        = "mesheryctl-1007"
	ErrDesignsNotFoundCode       = "mesheryctl-1037"
	ErrInvalidDesignFileCode     = "mesheryctl-1038"
	ErrDesignInvalidNameOrIDCode = "mesheryctl-1039"
)

const (
	errDesignMsg = `Usage: mesheryctl design import -f [file/url] -s [source-type]
Example: mesheryctl design import -f ./design.yml -s "Kubernetes Manifest"`

	errOnboardMsg = `Usage: mesheryctl design onboard -f [filepath] -s [source type]
Example: mesheryctl design onboard -f ./design.yml -s "Kubernetes Manifest"
Description: Onboard design`
)

func ErrDesignNotFound() error {
	return errors.New(ErrDesignsNotFoundCode, errors.Fatal, []string{"Design Not Found"}, []string{"No Design found with the given name or ID"}, []string{"Design with the given name or ID is not present"}, []string{"Please check if the given design name or ID is present via 'mesheryctl design list'"})
}

func ErrInvalidDesignFile(err error) error {
	return errors.New(ErrInvalidDesignFileCode, errors.Fatal, []string{err.Error()}, []string{"Design appears invalid. Could not parse successfully"}, []string{"Design file provided is not valid"}, []string{"Please check that your design file is a valid yaml file"})
}

func ErrDesignInvalidNameOrID(err error) error {
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
		[]string{"Provide the path to the design manifest. \n\n%v", errDesignMsg})
}

func ErrOnboardDesign() error {
	return errors.New(ErrOnboardDesignCode, errors.Alert,
		[]string{"Error Onboarding design"},
		[]string{"Unable to onboard design due to empty path"},
		[]string{"File path or design name not provided."},
		[]string{"Provide a file path/design name. \n\n%v", errOnboardMsg})
}

func ErrDesignFound() error {
	return errors.New(ErrDesignFoundCode, errors.Alert,
		[]string{"design not found"},
		[]string{"No design found with a given name"},
		[]string{"design name not provided"},
		[]string{"Provide an design name. Use `mesheryctl design list`to display list of designlications"},
	)
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
