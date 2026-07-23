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

package workspaces

import (
	"github.com/meshery/meshkit/errors"
)

var ErrDeleteWorkspaceCode = "mesheryctl-1245"

func ErrDeleteWorkspace(err error) error {
	return errors.New(ErrDeleteWorkspaceCode, errors.Alert,
		[]string{"Failed to delete workspace"},
		[]string{err.Error()},
		[]string{"The workspace ID may be invalid", "Meshery server may be unreachable"},
		[]string{"Verify the workspace ID is correct and try again", "Ensure Meshery server is running and reachable"},
	)
}
