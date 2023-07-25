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

package ctlerrors

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrProcessingConfigCode        = "1050"
	ErrCreatingConfigFileCode      = "1051"
	ErrAddingTokenToConfigCode     = "1052"
	ErrAddingContextToConfigCode   = "1053"
	ErrUnmarshallingConfigFileCode = "1054"
	ErrGettingRequestContextCode   = "1055"
	ErrInvalidAPIResponseCode      = "1056"
	ErrUnmarshallingAPIDataCode    = "1057"
	ErrConnectingToServerCode      = "1058"
	ErrMarshalCode                 = "replace_me"
	ErrMarshalIndentCode           = "replace_me"
	ErrUnmarshalCode               = "replace_me"
	ErrReadFileCode                = "replace_me"
	ErrReadConfigFileCode          = "replace_me"
	ErrNewRequestCode              = "replace_me"
	ErrMakeRequestCode             = "replace_me"
	ErrReadResponseBodyCode        = "replace_me"
	ErrResponseStatusBodyCode      = "replace_me"
	ErrResponseStatusCode          = "replace_me"
	ErrJSONToYAMLCode              = "replace_me"
	ErrOutFormatFlagCode           = "replace_me"
)

var (
	ErrCreatingConfigFile = errors.New(ErrCreatingConfigFileCode, errors.Alert, []string{"Unable to create config file"}, []string{"Unable to create config file"}, []string{}, []string{})

	ErrAddingTokenToConfig = errors.New(ErrAddingTokenToConfigCode, errors.Alert, []string{"Unable to add token to config"}, []string{"Unable to add token to config"}, []string{}, []string{})

	ErrAddingContextToConfig = errors.New(ErrAddingContextToConfigCode, errors.Alert, []string{"Unable to add context to config"}, []string{"Unable to add context to config"}, []string{}, []string{})

	ErrUnmarshallingConfigFile = errors.New(ErrUnmarshallingConfigFileCode, errors.Alert, []string{"Error processing json in config file"}, []string{"Error processing json in config file"}, []string{}, []string{})
)

func ErrProcessingConfig(err error) error {
	return errors.New(ErrProcessingConfigCode, errors.Alert, []string{"Error processing config"}, []string{"Error processing config file from meshconfig", err.Error()}, []string{}, []string{})
}

func ErrConnectingToServer(err error) error {
	return errors.New(ErrConnectingToServerCode, errors.Fatal, []string{"Unable to communicate with Meshery server"}, []string{"Unable to communicate with Meshery server", err.Error(), "See https://docs.meshery.io for help getting started with Meshery"}, []string{}, []string{"See https://docs.meshery.io for help getting started with Meshery"})
}

func ErrGettingRequestContext(err error) error {
	return errors.New(ErrGettingRequestContextCode, errors.Fatal, []string{"Unable to add token to config"}, []string{"Unable to add token to config", err.Error()}, []string{}, []string{})
}

func ErrInvalidAPIResponse(err error) error {
	return errors.New(ErrInvalidAPIResponseCode, errors.Fatal, []string{"Invalid API response encountered"}, []string{"Invalid API response encountered", err.Error()}, []string{}, []string{})
}

func ErrUnmarshallingAPIData(err error) error {
	return errors.New(ErrUnmarshallingAPIDataCode, errors.Fatal, []string{"Error processing json API data"}, []string{"Error processing json API data", err.Error()}, []string{}, []string{})
}

func ErrUnmarshal(err error) error {
	return errors.New(ErrUnmarshalCode, errors.Alert, []string{"Error unmarshalling response"}, []string{"Error processing JSON response from server.\n" + err.Error()}, []string{}, []string{})
}

func ErrMarshal(err error) error {
	return errors.New(ErrMarshalCode, errors.Alert, []string{"Error marshalling request"}, []string{"Error coverting requestbody to JSON.\n" + err.Error()}, []string{}, []string{})
}

func ErrMarshalIndent(err error) error {
	return errors.New(ErrMarshalIndentCode, errors.Alert, []string{"Error indenting json body"}, []string{"Can't marshal indent body filters", "" + err.Error()}, []string{}, []string{})
}

func ErrReadFile(err error) error {
	return errors.New(ErrReadFileCode, errors.Alert, []string{"Unable to read file"}, []string{"unable to read file", "" + err.Error()}, []string{}, []string{})
}

func ErrReadConfigFile(err error) error {
	return errors.New(ErrReadConfigFileCode, errors.Alert, []string{"Unable to read config file"}, []string{"Can't read config file from Path", "" + err.Error()}, []string{}, []string{})
}

func ErrNewRequest(err error) error {
	return errors.New(ErrNewRequestCode, errors.Alert, []string{"Error creating new request"}, []string{"Error in creating new Request", "" + err.Error()}, []string{}, []string{})
}

func ErrMakeRequest(err error) error {
	return errors.New(ErrMakeRequestCode, errors.Alert, []string{"Error in making request"}, []string{"Can't return response from the new request made", "" + err.Error()}, []string{}, []string{})
}

func ErrReadResponseBody(err error) error {
	return errors.New(ErrReadResponseBodyCode, errors.Alert, []string{"Error in reading response"}, []string{"Can't read response body", "" + err.Error()}, []string{}, []string{})
}

func ErrResponseStatusBody(statusCode int, body string) error {
	return errors.New(ErrResponseStatusBodyCode, errors.Alert, []string{"Wrong status code"}, []string{"Server returned with status code: " + fmt.Sprint(statusCode) + "\n" + "Response: " + body}, []string{}, []string{})
}

func ErrResponseStatus(statusCode int) error {
	return errors.New(ErrResponseStatusCode, errors.Alert, []string{"Wrong status code"}, []string{"Server returned with status code: " + fmt.Sprint(statusCode)}, []string{}, []string{})
}

func ErrJSONToYAML(err error) error {
	return errors.New(ErrJSONToYAMLCode, errors.Alert, []string{"Can't convert JSON to yaml"}, []string{"Failed to convert json to yaml", "" + err.Error()}, []string{}, []string{})
}

func ErrOutFormatFlag() error {
	return errors.New(ErrOutFormatFlagCode, errors.Alert, []string{"Output-format choice invalid"}, []string{"Output-format choice invalid, use [json|yaml]"}, []string{}, []string{})
}
