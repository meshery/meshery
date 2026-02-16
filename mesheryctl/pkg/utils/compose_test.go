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

package utils

import (
	"strings"
	"testing"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_ConvertToComposeSummaries(t *testing.T) {
	tests := []struct {
		name                  string
		inputContainerSummary []container.Summary
		wantContainerSummary  []api.ContainerSummary
		wantErrContains       string
	}{
		{
			name: "compose container with labels health exitcode ports mounts networks",
			inputContainerSummary: []container.Summary{
				{
					ID:      "0123456789abcdef0123456789abcdef",
					Names:   []string{"/meshery", "/meshery_alias"},
					Image:   "meshery/meshery:stable",
					Command: "meshery-server",
					Created: 1700000000,
					Ports: []container.Port{
						{IP: "0.0.0.0", PrivatePort: 9443, PublicPort: 9443, Type: "tcp"},
						{IP: "0.0.0.0", PrivatePort: 9081, PublicPort: 9081, Type: "tcp"},
					},
					SizeRw:     100,
					SizeRootFs: 200,
					Labels: map[string]string{
						api.ProjectLabel: "meshery",
						api.ServiceLabel: "meshery",
					},
					State:  container.StateRunning,
					Status: "Up 10 seconds",
					NetworkSettings: &container.NetworkSettingsSummary{
						Networks: map[string]*network.EndpointSettings{
							"bridge": {},
						},
					},
					Mounts: []container.MountPoint{
						{Name: "meshery-data", Driver: "local"},
						{Source: "/tmp/meshery", Driver: "bind"},
					},
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					Name:    "meshery",
					Image:   "meshery/meshery:stable",
					Command: "meshery-server",
					Service: "meshery",
					Created: 1700000000,
					State:   "running",
					Publishers: api.PortPublishers{
						{URL: "0.0.0.0", TargetPort: 9081, PublishedPort: 9081, Protocol: "tcp"},
						{URL: "0.0.0.0", TargetPort: 9443, PublishedPort: 9443, Protocol: "tcp"},
					},
				},
			},
		},
		{
			name: "non compose container has empty project and service",
			inputContainerSummary: []container.Summary{
				{
					ID:      "abcdef1234567890abcdef1234567890",
					Names:   []string{"/standalone-nginx"},
					Image:   "nginx:latest",
					Command: "nginx -g 'daemon off;'",
					Created: 1700001000,
					Labels:  map[string]string{},
					State:   container.StateRunning,
					Status:  "Up 2 minutes",
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					Name:    "standalone-nginx",
					Image:   "nginx:latest",
					Command: "nginx -g 'daemon off;'",
					Created: 1700001000,
					State:   "running",
				},
			},
		},
		{
			name: "inspect provides exited exitcode",
			inputContainerSummary: []container.Summary{
				{
					ID:      "oa7jFYpO239rZS6La05mrQpPHLlxAXUf",
					Names:   []string{"/job-runner"},
					Image:   "busybox",
					Command: "sh -c 'exit 137'",
					Created: 1700002000,
					Labels:  map[string]string{},
					State:   container.StateExited,
					Status:  "Exited (137) 1 minute ago",
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					Name:    "job-runner",
					Image:   "busybox",
					Command: "sh -c 'exit 137'",
					Created: 1700002000,
					State:   "exited",
				},
			},
		},
		{
			name: "names empty long id falls back to first 12 chars",
			inputContainerSummary: []container.Summary{
				{
					ID:      "11112222333344445555666677778888",
					Image:   "redis:7",
					Command: "redis-server",
					Labels:  map[string]string{},
					State:   container.StateRunning,
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					Name:    "111122223333",
					Image:   "redis:7",
					Command: "redis-server",
					State:   "running",
				},
			},
		},
		{
			name: "names empty short id falls back to full id",
			inputContainerSummary: []container.Summary{
				{
					ID:      "shortid7",
					Image:   "alpine",
					Command: "sleep 10",
					Labels:  map[string]string{},
					State:   container.StateRunning,
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					Name:    "shortid7",
					Image:   "alpine",
					Command: "sleep 10",
					State:   "running",
				},
			},
		},
		{
			name: "inspect state nil keeps health empty and exitcode zero",
			inputContainerSummary: []container.Summary{
				{
					ID:      "86myXcLitwGNarOO15jqbyAejTM04UP7",
					Names:   []string{"/state-nil"},
					Image:   "alpine",
					Command: "sleep 5",
					Labels:  map[string]string{},
					State:   container.StateRunning,
					Status:  "Up 5 seconds",
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					Name:    "state-nil",
					Image:   "alpine",
					Command: "sleep 5",
					State:   "running",
				},
			},
		},
		{
			name:                  "empty input returns empty output",
			inputContainerSummary: []container.Summary{},
			wantContainerSummary:  []api.ContainerSummary{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotContainerSummary, gotErr := ConvertToComposeSummaries(test.inputContainerSummary)
			if gotErr != nil {
				if !(strings.Contains(gotErr.Error(), test.wantErrContains)) {
					t.Errorf("Got Err: (%v) does not contain (%v)", gotErr, test.wantErrContains)
				}
			}
			if diff := cmp.Diff(gotContainerSummary, test.wantContainerSummary, cmpopts.EquateEmpty()); diff != "" {
				t.Errorf("ToComposeSummaries() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_canonicalContainerName(t *testing.T) {
	tests := []struct {
		name  string
		input container.Summary
		want  string
	}{
		{
			name: "Canonical name found among multiple names",
			input: container.Summary{
				Names: []string{"/project/linked_container", "/meshery", "/alias"},
				ID:    "0123456789abcdef",
			},
			want: "meshery",
		},
		{
			name: "No canonical name found, fallback to first name",
			input: container.Summary{
				Names: []string{"/complex/path/name", "/another/path"},
			},
			want: "complex/path/name",
		},
		{
			name: "Empty names slice, long ID (12-char fallback)",
			input: container.Summary{
				Names: []string{},
				ID:    "0123456789abcdefgh",
			},
			want: "0123456789ab",
		},
		{
			name: "Empty names slice, short ID (full ID fallback)",
			input: container.Summary{
				Names: nil,
				ID:    "abc-123",
			},
			want: "abc-123",
		},
		{
			name: "Name without leading slash (fallback logic)",
			input: container.Summary{
				Names: []string{"standalone-name"},
			},
			want: "standalone-name",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := canonicalContainerName(tc.input)
			if got != tc.want {
				t.Errorf("canonicalContainerName() = %v, want %v", got, tc.want)
			}
		})
	}
}
