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

package root

import "github.com/layer5io/meshkit/errors"

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrProcessingConfigCode        = "1072"
	ErrCreatingConfigFileCode      = "1073"
	ErrAddingTokenToConfigCode     = "1074"
	ErrAddingContextToConfigCode   = "1075"
	ErrUnmarshallingConfigFileCode = "1076"
	ErrGettingRequestContextCode   = "1077"
	ErrUnmarshallingAPIDataCode    = "1078"
	ErrConnectingToServerCode      = "1079"
)

var (
	ErrCreatingConfigFile = errors.New(ErrCreatingConfigFileCode, errors.Alert, []string{"Unable to create config file"}, []string{"Unable to create config file"}, []string{"Default config Path doesn't exist"}, []string{"Please ensure your meshery folder exist "})

	ErrAddingTokenToConfig = errors.New(ErrAddingTokenToConfigCode, errors.Alert, []string{"Unable to add token to config"}, []string{"Unable to add token to config"}, []string{"Invalid Config Path"}, []string{"Run `Mesheryctl system --help` to see config options"})

	ErrAddingContextToConfig = errors.New(ErrAddingContextToConfigCode, errors.Alert, []string{"Unable to add context to config"}, []string{"Unable to add context to config"}, []string{"Invalid Config Path"}, []string{"Run `Mesheryctl system --help` to see config options"})

	ErrUnmarshallingConfigFile = errors.New(ErrUnmarshallingConfigFileCode, errors.Alert, []string{"Error processing json in config file"}, []string{"Error processing json in config file"}, []string{"Couldn't load conig file to json object"}, []string{"Please ensure meshery configurations are valid in the config file "})
)

func ErrProcessingConfig(err error) error {
	return errors.New(ErrProcessingConfigCode, errors.Alert, []string{"Error processing config"},
		[]string{"Error processing config", err.Error()}, []string{"Couldn't load mesh config"}, []string{"Ensure mesheryctl has the right configurations"})
}

func ErrConnectingToServer(err error) error {
	return errors.New(ErrConnectingToServerCode, errors.Fatal, []string{"Unable to communicate with Meshery server"},
		[]string{"Unable to communicate with Meshery server", err.Error()}, []string{"There might be connection failure to Meshery Server"}, []string{"See https://docs.meshery.io for help getting started with Meshery"})
}

func ErrGettingRequestContext(err error) error {
	return errors.New(ErrGettingRequestContextCode, errors.Fatal, []string{"Unable to add token to config"},
		[]string{"Unable to add token to config", err.Error()}, []string{"Meshery is not running or there is a network issue"}, []string{"Check your network connection and check the status of meshery server via 'mesheryctl system status'"})
}

<<<<<<< HEAD
func ErrInvalidAPIResponse(err error) error {
	return errors.New(ErrInvalidAPIResponseCode, errors.Fatal, []string{"Invalid API response encountered"},
		[]string{"Invalid API response encountered", err.Error()}, []string{"Error occured while generating a response body"}, []string{"Check your network connection and check the status of meshery server via 'mesheryctl system status'"})
}

=======
>>>>>>> master
func ErrUnmarshallingAPIData(err error) error {
	return errors.New(ErrUnmarshallingAPIDataCode, errors.Fatal, []string{"Error processing json API data"},
		[]string{"Error processing json API data", err.Error()}, []string{"The json format from Api Data is not valid"}, []string{"Check if valid json is given to process"})
}
