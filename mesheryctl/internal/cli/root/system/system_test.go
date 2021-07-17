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
		Name            string
		Action          string
		Args            []string
		ExpectError     bool
		TimeoutRequired int
	}{
		// Docker platform testing
		{
			Name:            "Start Meshery with Docker platform",
			Action:          "start",
			Args:            []string{"start", "-p", "docker", "-y"},
			ExpectError:     false,
			TimeoutRequired: 2,
		},
		{
			Name:            "Printing Meshery status with Docker platform",
			Action:          "status",
			Args:            []string{"status"},
			ExpectError:     false,
			TimeoutRequired: 0,
		},
		{
			Name:            "Stop Meshery with Docker platform",
			Action:          "stop",
			Args:            []string{"stop", "-y"},
			ExpectError:     false,
			TimeoutRequired: 2,
		},

		// Kubernetes platform testing
		{
			Name:            "Start Meshery with Kubernetes platform",
			Action:          "start",
			Args:            []string{"start", "-p", "kubernetes", "-y"},
			ExpectError:     false,
			TimeoutRequired: 2,
		},
		{
			Name:            "Printing Meshery status with Kubernetes platform",
			Action:          "status",
			Args:            []string{"status"},
			ExpectError:     false,
			TimeoutRequired: 0,
		},
		{
			Name:            "Stop Meshery with Kubernetes platform",
			Action:          "stop",
			Args:            []string{"stop", "-y"},
			ExpectError:     false,
			TimeoutRequired: 2,
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
			// Sleeping for required timeout
			t.Logf("Sleeping for %v minutes...", tt.TimeoutRequired)
			time.Sleep(time.Duration(tt.TimeoutRequired) * time.Minute)
			t.Log("Sleeping finished")
		})
	}
}
