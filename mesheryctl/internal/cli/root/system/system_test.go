package system

import (
	"os"
	"testing"
	"time"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestSystemCmdIntegration(t *testing.T) {
	// skipping this integration test with --short flag
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	err := utils.SetFileLocation()
	if err != nil {
		t.Fatal(err)
	}

	// create .meshery directory if not present
	if _, err := os.Stat(utils.MesheryFolder); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(utils.MesheryFolder, 0775)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Create config file if not present in meshery folder
	err = utils.CreateConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	// Add Token to context file
	err = config.AddTokenToConfig(utils.TemplateToken, utils.DefaultConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	// Add Context to context file
	err = config.AddContextToConfig("local", utils.TemplateContext, utils.DefaultConfigPath, true)
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigFile(utils.DefaultConfigPath)

	viper.AutomaticEnv() // read in environment variables that match

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
