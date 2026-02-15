package utils

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type fakeInspector struct {
	resp map[string]container.InspectResponse
	errs map[string]error
}

func (f fakeInspector) ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error) {
	if err, ok := f.errs[containerID]; ok {
		return container.InspectResponse{}, err
	}
	if r, ok := f.resp[containerID]; ok {
		return r, nil
	}
	return container.InspectResponse{
		ContainerJSONBase: &container.ContainerJSONBase{
			State: &container.State{Status: container.StateRunning},
		},
	}, nil
}

func Test_ToComposeSummaries(t *testing.T) {
	tests := []struct {
		name                    string
		inputContainerSummary   []container.Summary
		inputContainerInspector containerInspector
		wantContainerSummary    []api.ContainerSummary
		wantErrContains         string
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
			inputContainerInspector: fakeInspector{
				resp: map[string]container.InspectResponse{
					"0123456789abcdef0123456789abcdef": {
						ContainerJSONBase: &container.ContainerJSONBase{
							State: &container.State{
								Status: container.StateRunning,
								Health: &container.Health{Status: container.Healthy},
							},
						},
					},
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					ID:       "0123456789abcdef0123456789abcdef",
					Name:     "meshery",
					Names:    []string{"/meshery", "/meshery_alias"},
					Image:    "meshery/meshery:stable",
					Command:  "meshery-server",
					Project:  "meshery",
					Service:  "meshery",
					Created:  1700000000,
					State:    "running",
					Status:   "Up 10 seconds",
					Health:   "healthy",
					ExitCode: 0,
					Publishers: api.PortPublishers{
						{URL: "0.0.0.0", TargetPort: 9081, PublishedPort: 9081, Protocol: "tcp"},
						{URL: "0.0.0.0", TargetPort: 9443, PublishedPort: 9443, Protocol: "tcp"},
					},
					Labels:       map[string]string{api.ProjectLabel: "meshery", api.ServiceLabel: "meshery"},
					SizeRw:       100,
					SizeRootFs:   200,
					Mounts:       []string{"meshery-data", "/tmp/meshery"},
					Networks:     []string{"bridge"},
					LocalVolumes: 1,
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
			inputContainerInspector: fakeInspector{
				resp: map[string]container.InspectResponse{
					"abcdef1234567890abcdef1234567890": {
						ContainerJSONBase: &container.ContainerJSONBase{
							State: &container.State{Status: container.StateRunning},
						},
					},
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					ID:           "abcdef1234567890abcdef1234567890",
					Name:         "standalone-nginx",
					Names:        []string{"/standalone-nginx"},
					Image:        "nginx:latest",
					Command:      "nginx -g 'daemon off;'",
					Project:      "",
					Service:      "",
					Created:      1700001000,
					State:        "running",
					Status:       "Up 2 minutes",
					Health:       "",
					ExitCode:     0,
					Labels:       map[string]string{},
					LocalVolumes: 0,
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
			inputContainerInspector: fakeInspector{
				resp: map[string]container.InspectResponse{
					"oa7jFYpO239rZS6La05mrQpPHLlxAXUf": {
						ContainerJSONBase: &container.ContainerJSONBase{
							State: &container.State{Status: container.StateExited, ExitCode: 137},
						},
					},
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					ID:       "oa7jFYpO239rZS6La05mrQpPHLlxAXUf",
					Name:     "job-runner",
					Names:    []string{"/job-runner"},
					Image:    "busybox",
					Command:  "sh -c 'exit 137'",
					Created:  1700002000,
					State:    "exited",
					Status:   "Exited (137) 1 minute ago",
					ExitCode: 137,
					Labels:   map[string]string{},
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
			inputContainerInspector: fakeInspector{
				resp: map[string]container.InspectResponse{
					"11112222333344445555666677778888": {
						ContainerJSONBase: &container.ContainerJSONBase{
							State: &container.State{Status: container.StateRunning},
						},
					},
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					ID:      "11112222333344445555666677778888",
					Name:    "111122223333",
					Image:   "redis:7",
					Command: "redis-server",
					State:   "running",
					Labels:  map[string]string{},
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
			inputContainerInspector: fakeInspector{
				resp: map[string]container.InspectResponse{
					"shortid7": {
						ContainerJSONBase: &container.ContainerJSONBase{
							State: &container.State{Status: container.StateRunning},
						},
					},
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					ID:      "shortid7",
					Name:    "shortid7",
					Image:   "alpine",
					Command: "sleep 10",
					State:   "running",
					Labels:  map[string]string{},
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
			inputContainerInspector: fakeInspector{
				resp: map[string]container.InspectResponse{
					"deadbeefdeadbeefdeadbeefdeadbeef": {
						ContainerJSONBase: &container.ContainerJSONBase{
							State: nil,
						},
					},
				},
			},
			wantContainerSummary: []api.ContainerSummary{
				{
					ID:       "86myXcLitwGNarOO15jqbyAejTM04UP7",
					Name:     "state-nil",
					Names:    []string{"/state-nil"},
					Image:    "alpine",
					Command:  "sleep 5",
					State:    "running",
					Status:   "Up 5 seconds",
					Health:   "",
					ExitCode: 0,
					Labels:   map[string]string{},
				},
			},
		},
		{
			name: "inspect error returns failure",
			inputContainerSummary: []container.Summary{
				{
					ID:      "oa7jFYpO239rZS6La05mrQpPHLlxAXUf",
					Names:   []string{"/inspect-fail"},
					Image:   "busybox",
					Command: "sleep 1",
					Labels:  map[string]string{},
					State:   container.StateRunning,
					Status:  "Up 1 second",
				},
			},
			inputContainerInspector: fakeInspector{
				errs: map[string]error{
					"oa7jFYpO239rZS6La05mrQpPHLlxAXUf": errors.New("inspect failed for test"),
				},
			},
			wantErrContains: "inspect failed for test",
		},
		{
			name:                    "empty input returns empty output",
			inputContainerSummary:   []container.Summary{},
			inputContainerInspector: fakeInspector{},
			wantContainerSummary:    []api.ContainerSummary{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotContainerSummary, gotErr := ToComposeSummaries(context.Background(), test.inputContainerInspector, test.inputContainerSummary)
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
