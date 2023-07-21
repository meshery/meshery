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
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrFileReadCode  = "1000"
	ErrUnmarshalCode = "1003"
)

func ErrFileRead(err error) error {
	return errors.New(ErrFileReadCode, errors.Alert, []string{"Unable to read file"}, []string{err.Error()}, []string{"Provided file is not present or invalid path"}, []string{"Please provide a valid file path with a valid file"})
}

func ErrUnmarshal(err error) error {
	return errors.New(ErrUnmarshalCode, errors.Alert, []string{"Error unmarshalling response"}, []string{"Error processing JSON response from server.\n" + err.Error()}, []string{}, []string{})
}
