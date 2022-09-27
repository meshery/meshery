package system

import (
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestSystemCmdIntegration(t *testing.T) {
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
	utils.SetFileLocationTesting(currDir)

	tests := []struct {
		Name            string
		Action          string
		Args            []string
		ExpectError     bool
		TimeoutRequired int
	}{
		// Docker platform testing
		//start
		{
			Name:            "Start Meshery with Docker platform",
			Action:          "start",
			Args:            []string{"start", "-p", "docker", "-y"},
			ExpectError:     false,
			TimeoutRequired: 1,
		},
		//update
		{
			Name:            "Update Meshery",
			Action:          "update",
			Args:            []string{"update"},
			ExpectError:     false,
			TimeoutRequired: 1,
		},
		//restart
		{
			Name:            "Restart Meshery",
			Action:          "restart",
			Args:            []string{"restart"},
			ExpectError:     false,
			TimeoutRequired: 1,
		},
		//status
		{
			Name:            "Printing Meshery status with Docker platform",
			Action:          "status",
			Args:            []string{"status"},
			ExpectError:     false,
			TimeoutRequired: 0,
		},
		//stop
		{
			Name:            "Stop Meshery with Docker platform",
			Action:          "stop",
			Args:            []string{"stop", "-y"},
			ExpectError:     false,
			TimeoutRequired: 0,
		},

		// Kubernetes platform testing
		//start
		{
			Name:            "Start Meshery with Kubernetes platform",
			Action:          "start",
			Args:            []string{"start", "-p", "kubernetes", "-y"},
			ExpectError:     false,
			TimeoutRequired: 1,
		},
		//update
		{
			Name:            "Update Meshery",
			Action:          "update",
			Args:            []string{"update"},
			ExpectError:     false,
			TimeoutRequired: 1,
		},
		//restart
		{
			Name:            "Restart Meshery",
			Action:          "restart",
			Args:            []string{"restart"},
			ExpectError:     false,
			TimeoutRequired: 1,
		},
		//status
		{
			Name:            "Printing Meshery status with Kubernetes platform",
			Action:          "status",
			Args:            []string{"status"},
			ExpectError:     false,
			TimeoutRequired: 0,
		},
		//logs
		{
			Name:            "Printing Meshery logs with Kubernetes platform",
			Action:          "logs",
			Args:            []string{"logs"},
			ExpectError:     false,
			TimeoutRequired: 0,
		},
		//stop
		{
			Name:            "Stop Meshery with Kubernetes platform",
			Action:          "stop",
			Args:            []string{"stop", "-y"},
			ExpectError:     false,
			TimeoutRequired: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// setting up utils.Log
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			SystemCmd.SetArgs(tt.Args)

			t.Logf("Performing %s", tt.Action)
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
