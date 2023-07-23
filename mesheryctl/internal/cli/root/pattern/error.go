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

import "github.com/layer5io/meshkit/errors"

const (
	ErrPatternsNotFoundCode   = "replace_me"
	ErrInvalidPatternFileCode = "replace_me"
)

func ErrPatternsNotFound() error {
	return errors.New(ErrPatternsNotFoundCode, errors.Fatal, []string{"Patterns Not Found"}, []string{"No Patterns found with the given name"}, []string{"Patterns with the given name is not present to apply"}, []string{"Please check if the given pattern name is present via 'mesheryctl pattern list'"})
}

func ErrInvalidPatternFile(err error) error {
	return errors.New(ErrInvalidPatternFileCode, errors.Fatal, []string{err.Error()}, []string{"Pattern appears invalid. Could not parse successfully"}, []string{"Pattern file provided is not valid"}, []string{"Please check that your pattern file is a valid yaml file"})
}
