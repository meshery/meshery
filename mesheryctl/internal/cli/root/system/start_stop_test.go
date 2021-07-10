package system

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestSystemStartStopIntegration(t *testing.T) {
	// skipping this integration test with --short flag
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// setup current context
	utils.SetupContextEnv(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// update all location
	utils.SetFileLocationTesting(t, currDir)

	tests := []struct {
		Name        string
		Action      string
		Args        []string
		ExpectError bool
	}{
		{
			Name:        "Start Meshery with Docker platform",
			Action:      "start",
			Args:        []string{"start", "-p", "docker", "-y"},
			ExpectError: false,
		},
		{
			Name:        "Stop Meshery with Docker platform",
			Action:      "stop",
			Args:        []string{"stop", "-y"},
			ExpectError: false,
		},
		{
			Name:        "Start Meshery with Kubernetes platform",
			Action:      "start",
			Args:        []string{"start", "-p", "kubernetes", "-y"},
			ExpectError: false,
		},
		{
			Name:        "Stop Meshery with Kubernetes platform",
			Action:      "stop",
			Args:        []string{"stop", "-y"},
			ExpectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SystemCmd.SetArgs(tt.Args)

			t.Logf("%sing meshery", tt.Action)
			err := SystemCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("Meshery %sed", tt.Action)
			// Three minute timeout
			t.Log("Sleeping for 2 minutes...")
			time.Sleep(2 * time.Minute)
			t.Log("Sleeping finished")
		})
	}
}
