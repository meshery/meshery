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

package design

import "encoding/json"

// deploymentMessagePerComp mirrors the server's per-component deployment
// summary (server/models/pattern/patterns.DeploymentMessagePerComp); wire
// field names match Go field names as the server marshals without json tags.
type deploymentMessagePerComp struct {
	Kind     string
	CompName string
	Success  bool
}

type deploymentMessagePerContext struct {
	Summary []deploymentMessagePerComp
}

type deploymentSummary map[string][]deploymentMessagePerContext

// parseDeploymentSummary returns the names of every component that reported a
// failed deployment. ok is false when the body is not the summary shape.
func parseDeploymentSummary(body []byte) (failedComponents []string, ok bool) {
	var summary deploymentSummary
	if err := json.Unmarshal(body, &summary); err != nil {
		return nil, false
	}
	if summary == nil {
		return nil, false
	}

	seen := make(map[string]struct{})
	for _, contexts := range summary {
		for _, ctx := range contexts {
			for _, comp := range ctx.Summary {
				if comp.Success {
					continue
				}
				name := comp.CompName
				if name == "" {
					name = comp.Kind
				}
				if _, exists := seen[name]; exists {
					continue
				}
				seen[name] = struct{}{}
				failedComponents = append(failedComponents, name)
			}
		}
	}
	return failedComponents, true
}
