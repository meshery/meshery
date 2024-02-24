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
	ErrGettingSessionDataCode                = "1098"
	ErrNoAdaptersCode                        = "1099"
	ErrPromptCode                            = "1100"
	ErrCreatingDeployRequestCode             = "1101"
	ErrCreatingDeployResponseRequestCode     = "1102"
	ErrCreatingDeployResponseStreamCode      = "1103"
	ErrCreatingValidateResponseStreamCode    = "1104"
	ErrTimeoutWaitingForDeployResponseCode   = "1105"
	ErrFailedDeployingMeshCode               = "1106"
	ErrCreatingValidateRequestCode           = "1107"
	ErrCreatingValidateResponseRequestCode   = "1108"
	ErrTimeoutWaitingForValidateResponseCode = "1109"
	ErrSMIConformanceTestsFailedCode         = "1110"
	ErrValidateAdapterCode                   = "1111"
	ErrSendOperationCode                     = "1112"
	ErrValidMeshNameCode                     = "1113"
	ErrWaitValidateResponseCode              = "1114"
)

var (
	// no adapters found
	ErrNoAdapters = errors.New(ErrNoAdaptersCode, errors.Fatal,
		[]string{"Unable to validate adapter"},
		[]string{"Adapter for required mesh not found"},
		[]string{"Unable to fetch Mesh adapter, adapter not valid."},
		[]string{"Deploy a proper Meshery Adapter for your service mesh. See https://docs.meshery.io/concepts/architecture/adapters for list of adapters"})

	ErrFailedDeployingMesh = errors.New(ErrFailedDeployingMeshCode, errors.Fatal,
		[]string{"Failed to deploy the service mesh"},
		[]string{"Failed to deploy the service mesh"},
		[]string{"Network error or wrong environment configurations"},
		[]string{"Ensure your have a strong connection and the right environment configuration"})

	ErrTimeoutWaitingForDeployResponse = errors.New(ErrTimeoutWaitingForDeployResponseCode, errors.Fatal,
		[]string{"Timed out waiting for deploy event"},
		[]string{"Timed out waiting for deployment"},
		[]string{"Network error or wrong environment configurations"},
		[]string{"Ensure your have a strong connection and the right environment configuration"})

	ErrTimeoutWaitingForValidateResponse = errors.New(ErrTimeoutWaitingForValidateResponseCode, errors.Fatal,
		[]string{"Timed out waiting for validate response"},
		[]string{"Timed out waiting for validate response"},
		[]string{"Network error or wrong environment configurations"},
		[]string{"Ensure your have a strong connection and the right environment configuration"})

	ErrSMIConformanceTestsFailed = errors.New(ErrSMIConformanceTestsFailedCode, errors.Fatal, []string{"SMI conformance tests failed"}, []string{"SMI conformance tests failed"}, []string{}, []string{"Join https://layer5io.slack.com/archives/C010H0HE2E6"})
)

// When unable to get release data
func ErrGettingSessionData(err error) error {
	return errors.New(ErrGettingSessionDataCode, errors.Fatal,
		[]string{"Unable to fetch session data"},
		[]string{err.Error()},
		[]string{"Unable to get session from server"},
		[]string{"Check your network connection and verify the status of the Meshery server with `mesheryctl system status`."})
}

func ErrPrompt(err error) error {
	return errors.New(ErrPromptCode, errors.Fatal,
		[]string{"Error while reading selected option"},
		[]string{err.Error()},
		[]string{"Unable to execute Service Mesh from the list"},
		[]string{"Ensure you enter a valid Service Mesh"})
}

func ErrCreatingDeployRequest(err error) error {
	return errors.New(ErrCreatingDeployRequestCode, errors.Fatal,
		[]string{"Error sending deploy request"},
		[]string{err.Error()},
		[]string{"There might be a connection failure with Meshery Server"},
		[]string{"Check your network connection and verify the status of the Meshery server with `mesheryctl system status`."})
}

func ErrCreatingDeployResponseRequest(err error) error {
	return errors.New(ErrCreatingDeployResponseRequestCode,
		errors.Fatal,
		[]string{"Failed to create a response"},
		[]string{err.Error()},
		[]string{"Error occurred while generating a response from deployed request"},
		[]string{"Check your network connection and verify the status of the Meshery server with `mesheryctl system status`."})
}

func ErrCreatingDeployResponseStream(err error) error {
	return errors.New(ErrCreatingDeployResponseStreamCode,
		errors.Fatal,
		[]string{"Error creating deploy event response stream"},
		[]string{err.Error()},
		[]string{"Unable convert a connection to a stream of server sent events"},
		[]string{"Ensure your connection is valid and verify the status of the Meshery server with `mesheryctl system status`."})
}

func ErrCreatingValidateRequest(err error) error {
	return errors.New(ErrCreatingValidateRequestCode, errors.Fatal,
		[]string{"Error sending Validate request"},
		[]string{err.Error()},
		[]string{"An HTTP error occurred due to network connection"},
		[]string{"Check your network connection and verify the status of the Meshery server with `mesheryctl system status`."})
}

func ErrCreatingValidateResponseRequest(err error) error {
	return errors.New(ErrCreatingValidateResponseRequestCode, errors.Fatal,
		[]string{"Error creating request for validate response"},
		[]string{err.Error()},
		[]string{"Unable to create response from request"},
		[]string{"Check your network connection and verify the status of the Meshery server with `mesheryctl system status`."})
}

func ErrCreatingValidateResponseStream(err error) error {
	return errors.New(ErrCreatingValidateResponseStreamCode, errors.Fatal,
		[]string{"Error creating validate event response stream"},
		[]string{err.Error()},
		[]string{"Meshery is not running or there is a network issue."},
		[]string{"Check your network connection and verify the status of the Meshery server with `mesheryctl system status`."})
}

func ErrValidatingAdapters(err error) error {
	return errors.New(ErrValidateAdapterCode, errors.Fatal,
		[]string{"Unable to validate adapter"},
		[]string{err.Error()},
		[]string{"Unable to fetch Mesh adapter, adapter not valid."},
		[]string{"Enter a valid Mesh adapter."})
}

func ErrSendOperation(err error) error {
	return errors.New(ErrSendOperationCode, errors.Alert,
		[]string{"Unable to perform operation"},
		[]string{err.Error()},
		[]string{"Failed to perform an adapter operation due to network connection."},
		[]string{"Check your network connection and the status of Meshery Server via `mesheryctl system status`."})
}

func ErrValidMeshName(meshName string) error {
	return errors.New(ErrValidMeshNameCode, errors.Alert,
		[]string{"Invalid mesh name"},
		[]string{"Unable to validate mesh name"},
		[]string{"%s is not a valid mesh name or is unsupported", meshName},
		[]string{"Enter a valid mesh name"})
}

func ErrWaitValidateResponse(err error) error {
	return errors.New(ErrWaitValidateResponseCode, errors.Fatal,
		[]string{"Unable to validate a response"},
		[]string{err.Error()},
		[]string{"Unable to create responses after verifying operation"},
		[]string{"Ensure your connection is valid and verify the status of the Meshery server with `mesheryctl system status`."})
}
