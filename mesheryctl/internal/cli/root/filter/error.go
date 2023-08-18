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
	"strconv"

	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrInvalidAuthTokenCode = "1000"
	ErrInvalidAPICallCode   = "1001"
	ErrReadAPIResponseCode  = "1002"
	ErrFilterNameOrIDCode   = "replace_me"
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

func ErrFilterNameOrID(err error) error {
	return errors.New(ErrFilterNameOrIDCode, errors.Alert, []string{"Invalid filter name or ID"}, []string{"Invalid filter name or ID", "" + err.Error()}, []string{"Probable invalid filter name|ID"}, []string{"Please run `mesheryctl filter list` to view all filters avilable"})
}