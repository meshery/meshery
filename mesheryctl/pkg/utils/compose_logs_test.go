package utils

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/docker/compose/v2/pkg/api"
)

// logsRecorder stubs the compose service and only records the options Logs receives.
type logsRecorder struct {
	api.Compose
	project string
	options api.LogOptions
}

func (r *logsRecorder) Logs(_ context.Context, projectName string, _ api.LogConsumer, options api.LogOptions) error {
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
