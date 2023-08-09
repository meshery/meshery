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

package mesh

import "github.com/layer5io/meshkit/errors"

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrGettingSessionDataCode                = "1009"
	ErrNoAdaptersCode                        = "1010"
	ErrPromptCode                            = "1011"
	ErrCreatingDeployRequestCode             = "1012"
	ErrCreatingDeployResponseRequestCode     = "1013"
	ErrAddingAuthDetailsCode                 = "1014"
	ErrCreatingDeployResponseStreamCode      = "1015"
	ErrTimeoutWaitingForDeployResponseCode   = "1016"
	ErrFailedDeployingMeshCode               = "1017"
	ErrCreatingValidateRequestCode           = "1018"
	ErrCreatingValidateResponseRequestCode   = "1019"
	ErrTimeoutWaitingForValidateResponseCode = "1020"
	ErrSMIConformanceTestsFailedCode         = "1021"
)

var (
	// no adapters found
	ErrNoAdapters = errors.New(ErrNoAdaptersCode, errors.Fatal, []string{"Adapter for required mesh not found"}, []string{"Adapter for required mesh not found"}, []string{""}, []string{"Deploy the proper Meshery Adapter for your service mesh"})

	ErrFailedDeployingMesh = errors.New(ErrFailedDeployingMeshCode, errors.Fatal, []string{"Failed to deploy the service mesh"}, []string{"Failed to deploy the service mesh"}, []string{}, []string{"Check your environment and try again"})

	ErrTimeoutWaitingForDeployResponse = errors.New(ErrTimeoutWaitingForDeployResponseCode, errors.Fatal, []string{"Timed out waiting for deploy event"}, []string{"Timed out waiting for deployment"}, []string{}, []string{"Check your environment and try again"})

	ErrTimeoutWaitingForValidateResponse = errors.New(ErrTimeoutWaitingForValidateResponseCode, errors.Fatal, []string{"Timed out waiting for validate response"}, []string{"Timed out waiting for validate response"}, []string{""}, []string{"Check your environment and try again"})

	ErrSMIConformanceTestsFailed = errors.New(ErrSMIConformanceTestsFailedCode, errors.Fatal, []string{"SMI conformance tests failed"}, []string{"SMI conformance tests failed"}, []string{}, []string{"Join https://layer5io.slack.com/archives/C010H0HE2E6"})
)

// When unable to get release data
func ErrGettingSessionData(err error) error {
	return errors.New(ErrGettingSessionDataCode, errors.Fatal, []string{"Unable to fetch session data"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPrompt(err error) error {
	return errors.New(ErrPromptCode, errors.Fatal, []string{"Error while reading selected option"}, []string{err.Error()}, []string{}, []string{})
}

func ErrCreatingDeployRequest(err error) error {
	return errors.New(ErrCreatingDeployRequestCode, errors.Fatal, []string{"Error sending deploy request"}, []string{err.Error()}, []string{}, []string{})
}

func ErrCreatingDeployResponseRequest(err error) error {
	return errors.New(ErrCreatingDeployResponseRequestCode, errors.Fatal, []string{"Error creating request for deploy response"}, []string{err.Error()}, []string{}, []string{})
}

func ErrAddingAuthDetails(err error) error {
	return errors.New(ErrAddingAuthDetailsCode, errors.Fatal, []string{"Error adding auth details"}, []string{err.Error()}, []string{}, []string{})
}

func ErrCreatingDeployResponseStream(err error) error {
	return errors.New(ErrCreatingDeployResponseStreamCode, errors.Fatal, []string{"Error creating deploy event response stream"}, []string{err.Error()}, []string{}, []string{})
}

func ErrCreatingValidateRequest(err error) error {
	return errors.New(ErrCreatingValidateRequestCode, errors.Fatal, []string{"Error sending Validate request"}, []string{err.Error()}, []string{}, []string{})
}

func ErrCreatingValidateResponseRequest(err error) error {
	return errors.New(ErrCreatingValidateResponseRequestCode, errors.Fatal, []string{"Error creating request for validate response"}, []string{err.Error()}, []string{}, []string{})
}

func ErrCreatingValidateResponseStream(err error) error {
	return errors.New(ErrCreatingDeployResponseStreamCode, errors.Fatal, []string{"Error creating validate event response stream"}, []string{err.Error()}, []string{}, []string{})
}
