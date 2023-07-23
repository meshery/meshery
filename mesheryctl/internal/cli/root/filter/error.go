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

package filter

import (
	"fmt"
	"strconv"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrInvalidAuthTokenCode   = "1000"
	ErrInvalidAPICallCode     = "1001"
	ErrReadAPIResponseCode    = "1002"
	ErrUnmarshalCode          = "1003"
	ErrMarshalCode            = "1004"
	ErrMarshalIndentCode      = "1005"
	ErrReadFileCode           = "1006"
	ErrReadConfigFileCode     = "1007"
	ErrProcessConfigFileCode  = "1008"
	ErrNewRequestCode         = "1009"
	ErrMakeRequestCode        = "1010"
	ErrReadResponseBodyCode   = "1011"
	ErrResponseStatusBodyCode = "1012"
	ErrResponseStatusCode     = "1013"
	ErrJSONToYAMLCode         = "1014"
	ErrOutFormatFlagCode      = "1015"
	ErrFilterNameOrIDCode     = "1016"
)

func ErrInvalidAuthToken() error {
	return errors.New(ErrInvalidAuthTokenCode, errors.Alert, []string{"Authentication token not found. please supply a valid user token with the --token (or -t) flag"}, []string{}, []string{}, []string{})
}

func ErrInvalidAPICall(statusCode int) error {
	return errors.New(ErrInvalidAPICallCode, errors.Alert, []string{"Response Status Code ", strconv.Itoa(statusCode), " Possible Server Error"}, []string{"Response Status Code " + strconv.Itoa(statusCode) + " possible Server Error"}, []string{}, []string{})
}

func ErrReadAPIResponse(err error) error {
	return errors.New(ErrReadAPIResponseCode, errors.Alert, []string{"Failed to read response body"}, []string{"Failed to read response body", " " + err.Error()}, []string{}, []string{})
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

func ErrProcessConfigFile(err error) error {
	return errors.New(ErrProcessConfigFileCode, errors.Alert, []string{"Unable to process config file"}, []string{"Can't process config file from meshconfig", "" + err.Error()}, []string{}, []string{})
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
	return errors.New(ErrOutFormatFlagCode, errors.Alert, []string{"output-format choice invalid"}, []string{"output-format choice invalid, use [json|yaml]"}, []string{}, []string{})
}

func ErrFilterNameOrID(err error) error {
	return errors.New(ErrFilterNameOrIDCode, errors.Alert, []string{"invalid filter name or ID"}, []string{"invalid filter name or ID", "" + err.Error()}, []string{}, []string{})
}