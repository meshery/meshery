// Copyright 2023 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License
package app

import (
	"strings"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrImportAppCode          = "1080"
	ErrValidSourceCode        = "1081"
	ErrAppManifestCode        = "1189"
	ErrOnboardAppCode         = "1082"
	ErrAppFoundCode           = "1083"
	ErrViewAppFlagCode        = "1190"
	ErrInvalidAppNameOrIDCode = "1084"
)
const (
	errAppMsg = `Usage: mesheryctl app import -f [file/url] -s [source-type]
Example: mesheryctl app import -f ./application.yml -s "Kubernetes Manifest"`

	errOnboardMsg = `Usage: mesheryctl app onboard -f [filepath] -s [source type]
Example: mesheryctl app onboard -f ./application.yml -s "Kubernetes Manifest"
Description: Onboard application`
)

func ErrImportApp(err error) error {
	return errors.New(ErrImportAppCode, errors.Fatal,
		[]string{"Unable to import app"},
		[]string{err.Error()},
		[]string{"Invalid application file"},
		[]string{"Check your application URL/file path"})
}

func ErrValidSource(validSourceTypes []string) error {
	return errors.New(ErrValidSourceCode, errors.Fatal,
		[]string{"Invalid application source type"},
		[]string{"Invalid application source type due to wrong type/passing"},
		[]string{"Application source type (-s) is invalid or not passed."},
		[]string{"Ensure you pass a valid source type. \nAllowed source types: %s", strings.Join(validSourceTypes, ", ")})
}

func ErrAppManifest() error {
	return errors.New(ErrAppManifestCode, errors.Alert,
		[]string{"No file path detected"},
		[]string{"No manifest file path detected"},
		[]string{"Manifest path not provided"},
		[]string{"Provide the path to the app manifest. \n\n%v", errAppMsg})
}

func ErrOnboardApp() error {
	return errors.New(ErrOnboardAppCode, errors.Alert,
		[]string{"Error Onboarding app"},
		[]string{"Unable to onboard app due to empty path"},
		[]string{"File path or application name not provided."},
		[]string{"Provide a file path/app name. \n\n%v", errOnboardMsg})
}

func ErrAppFound() error {
	return errors.New(ErrAppFoundCode, errors.Alert,
		[]string{"App not found"},
		[]string{"No app found with a given name"},
		[]string{"Application name not provided"},
		[]string{"Provide an app name. Use `mesheryctl app list`to display list of applications"},
	)
}

func ErrInvalidAppNameOrID(err error) error {
	return errors.New(ErrInvalidAppNameOrIDCode, errors.Alert,
		[]string{"Invalid application"},
		[]string{"Failed to get application based on input"},
		[]string{"Application name|id is invalid"},
		[]string{"Run `mesheryctl app view --all` to view all applications"})
}

func ErrViewAppFlag() error {
	return errors.New(ErrViewAppFlagCode, errors.Alert,
		[]string{"Invalid command"},
		[]string{"Wrong use of command flags"},
		[]string{"-a/all flag is being used while an app is specified"},
		[]string{"-a/-all cannot be used when [application name|id] is specified"})
}
