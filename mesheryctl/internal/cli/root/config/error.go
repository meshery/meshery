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

package config

import "github.com/meshery/meshkit/errors"

const (
	ErrInvalidMeshConfigCode          = "mesheryctl-1008"
	ErrUpdateConfigCode               = "mesheryctl-1009"
	ErrReadMesheryConfigCode          = "mesheryctl-1220"
	ErrDuplicateContextCode           = "mesheryctl-1221"
	ErrWriteMeshConfigCode            = "mesheryctl-1222"
	ErrInvalidVersionCode             = "mesheryctl-1256"
	ErrValidateVersionRequestCode     = "mesheryctl-1257"
	ErrValidateVersionFailRequestCode = "mesheryctl-1258"
)

func ErrInvalidMeshConfig(err error) error {
	return errors.New(ErrInvalidMeshConfigCode, errors.Alert, []string{"Invalid Meshconfig"},
		[]string{err.Error()}, []string{"Loading the Invalid MeshConfig data caused error"}, []string{"Make sure that your 'config.yaml' file in your '.meshery' is valid"})
}

func ErrUpdateConfig(err error) error {
	return errors.New(ErrUpdateConfigCode, errors.Fatal,
		[]string{"Error in writing config"},
		[]string{err.Error()},
		[]string{"Unable to Update config file"},
		[]string{"Ensure that you have the correct context in your meshconfig at `$HOME/.meshery/config.yaml`."})
}

func ErrReadMesheryConfig(err error) error {
	return errors.New(ErrReadMesheryConfigCode, errors.Fatal,
		[]string{"Error in reading config"},
		[]string{err.Error()},
		[]string{"Unable to read config file"},
		[]string{"Ensure that you have the correct context in your meshconfig at `$HOME/.meshery/config.yaml`."})
}

func ErrDuplicateContext(err error) error {
	return errors.New(ErrDuplicateContextCode, errors.Fatal,
		[]string{"Duplicate Context"},
		[]string{err.Error()},
		[]string{"A context with the same name already exists in the config file"},
		[]string{"Use a different name for the context or delete the existing context with the same name"})
}

func ErrWriteMeshConfig(err error) error {
	return errors.New(ErrWriteMeshConfigCode, errors.Fatal,
		[]string{"Error in writing config"},
		[]string{err.Error()},
		[]string{"Unable to write to config file"},
		[]string{"Ensure that you have the correct permissions to write to the config file at `$HOME/.meshery/config.yaml`."})
}

// ErrValidateVersionRequest returns an error when an HTTP request to validate release version cannot be created.
func ErrValidateVersionRequest(err error) error {
	return errors.New(
		ErrValidateVersionRequestCode,
		errors.Fatal,
		[]string{"Error occurred while making an HTTP request."},
		[]string{err.Error()},
		[]string{"Could not construct HTTP request to validate release version."},
		[]string{"Check your request values or network configuration."},
	)
}

// ErrValidateVersionFailRequest returns an error when the HTTP request to validate release version fails.
func ErrValidateVersionFailRequest(err error) error {
	return errors.New(
		ErrValidateVersionFailRequestCode,
		errors.Alert,
		[]string{"Failed to make a request"},
		[]string{err.Error()},
		[]string{"Could not reach GitHub to validate the Meshery release version."},
		[]string{"Check your internet connection and verify that github.com is reachable."},
	)
}

// ErrInvalidVersion returns an error when version validation fails.
func ErrInvalidVersion(err error) error {
	return errors.New(
		ErrInvalidVersionCode,
		errors.Alert,
		[]string{"Error validating Meshery version"},
		[]string{err.Error()},
		[]string{"The version in the context is invalid, unsupported, or unreachable on GitHub"},
		[]string{"Verify the version in your context is valid, or check your network connection"},
	)
}
