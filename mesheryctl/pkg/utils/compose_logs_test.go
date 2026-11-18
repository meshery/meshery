package utils

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/docker/compose/v2/pkg/api"
)

// logsRecorder stubs the compose service and only records the options Logs receives.
type logsRecorder struct {
	api.Compose
	called  bool
	project string
	options api.LogOptions
}

func (r *logsRecorder) Logs(_ context.Context, projectName string, _ api.LogConsumer, options api.LogOptions) error {
	r.called = true
	r.project = projectName
	r.options = options
	return nil
}

func TestComposeClientLogsPassesFollow(t *testing.T) {
	composeFile := filepath.Join(t.TempDir(), "docker-compose.yaml")
	content := "services:\n  meshery:\n    image: meshery/meshery:stable\n"
	if err := os.WriteFile(composeFile, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, follow := range []bool{false, true} {
		recorder := &logsRecorder{}
		client := &ComposeClient{service: recorder}

		if err := client.Logs(context.Background(), composeFile, follow, io.Discard); err != nil {
			t.Fatalf("Logs(follow=%v) returned error: %v", follow, err)
		}
		if recorder.options.Follow != follow {
			t.Errorf("Logs(follow=%v) sent Follow=%v to compose", follow, recorder.options.Follow)
		}
		if recorder.project != "meshery" {
			t.Errorf("Logs used project %q, want %q", recorder.project, "meshery")
		}
	}
}

// writeTestComposeFile writes a compose file with the meshery and meshery-istio
// services.
func writeTestComposeFile(t *testing.T) string {
	t.Helper()
	composeFile := filepath.Join(t.TempDir(), "docker-compose.yaml")
	content := "services:\n  meshery:\n    image: meshery/meshery:stable\n  meshery-istio:\n    image: meshery/meshery-istio:stable\n"
	if err := os.WriteFile(composeFile, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return composeFile
}

// TestComposeClientLogsPassesServices verifies that only the requested
// services' logs are asked for, and that no services means all of them.
func TestComposeClientLogsPassesServices(t *testing.T) {
	composeFile := writeTestComposeFile(t)

	tests := []struct {
		name     string
		services []string
	}{
		{name: "no services means every service", services: nil},
		{name: "one service", services: []string{"meshery-istio"}},
		{name: "several services", services: []string{"meshery", "meshery-istio"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := &logsRecorder{}
			client := &ComposeClient{service: recorder}

			if err := client.Logs(context.Background(), composeFile, false, io.Discard, tt.services...); err != nil {
				t.Fatalf("Logs returned error: %v", err)
			}
			if !reflect.DeepEqual(recorder.options.Services, tt.services) {
				t.Errorf("Logs sent services %v to compose, want %v", recorder.options.Services, tt.services)
			}
		})
	}
}

// TestComposeClientLogsRejectsUnknownService verifies that a service that isn't
// in the deployment fails instead of printing every service's logs.
func TestComposeClientLogsRejectsUnknownService(t *testing.T) {
	composeFile := writeTestComposeFile(t)
	recorder := &logsRecorder{}
	client := &ComposeClient{service: recorder}

	err := client.Logs(context.Background(), composeFile, false, io.Discard, "meshery-linkerd")
	if err == nil {
		t.Fatal("expected an unknown service to be rejected")
	}
	if !strings.Contains(err.Error(), "meshery-linkerd") {
		t.Errorf("expected the error to name the service, got %q", err.Error())
	}
	if recorder.called {
		t.Error("expected no logs to be requested for an unknown service")
	}
}
