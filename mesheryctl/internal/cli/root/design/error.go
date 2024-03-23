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

package design

import "github.com/layer5io/meshkit/errors"

const (
	ErrDesignsNotFoundCode       = "mesheryctl-1037"
	ErrInvalidDesignFileCode     = "mesheryctl-1038"
	ErrDesignInvalidNameOrIDCode = "mesheryctl-1039"
)

func ErrDesignNotFound() error {
	return errors.New(ErrDesignsNotFoundCode, errors.Fatal, []string{"Design Not Found"}, []string{"No Design found with the given name or ID"}, []string{"Design with the given name or ID is not present"}, []string{"Please check if the given Design name or ID is present via 'mesheryctl design list'"})
}

func ErrInvalidDesignFile(err error) error {
	return errors.New(ErrInvalidDesignFileCode, errors.Fatal, []string{err.Error()}, []string{"Design appears invalid. Could not parse successfully"}, []string{"Design file provided is not valid"}, []string{"Please check that your Design file is a valid yaml file"})
}

func ErrDesignInvalidNameOrID(err error) error {
	return errors.New(
		ErrDesignInvalidNameOrIDCode,
		errors.Alert,
		[]string{"Unable to fetch design"},
		[]string{err.Error()},
		[]string{"Invalid Design name or ID"},
		[]string{"Run `mesheryctl design view -a` to view all available designs."})
}
