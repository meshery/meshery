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

package filter

import (
	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrFilterNameOrIDCode = "mesheryctl-1018"
	ErrFetchFilterCode    = "mesheryctl-1019"
)

func ErrFilterNameOrID(err error) error {
	return errors.New(
		ErrFilterNameOrIDCode,
		errors.Alert,
		[]string{"Unable to fetch filter"},
		[]string{err.Error()},
		[]string{"Probable invalid filter name or id"},
		[]string{"Run `mesheryctl filter list` to view all available filters."})
}

func ErrFetchFilter(err error) error {
	return errors.New(
		ErrFetchFilterCode,
		errors.Alert,
		[]string{"Unable to Fetch Filter"},
		[]string{err.Error()},
		[]string{"Filter name or id doesn't exist"},
		[]string{"Run `mesheryctl filter view -a` to view all available filters."})
}
