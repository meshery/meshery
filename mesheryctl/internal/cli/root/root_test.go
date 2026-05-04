package root

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

type RootCmdTestInput struct {
	Name             string
	Args             []string
	ExpectedResponse string
}

func resetFlags() {
	cfgFile = utils.DefaultConfigPath
	verbose = false
}

func TestRootCmdIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	fixturesDirPath := filepath.Join(currDir, "fixtures")
	utils.MesheryFolder = filepath.Join(fixturesDirPath, ".meshery")
	utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
	testConfigPath := filepath.Join(fixturesDirPath, "testConfig.yaml")

	tests := []RootCmdTestInput{
		{
			Name: "Given version command, When root command executes, Then version is returned",
			Args: []string{"version"},
		},
		{
			Name: "Given default config, When channel view is executed, Then stable channel is returned",
			Args:             []string{"system", "channel", "view"},
			ExpectedResponse: "Context: local\nChannel: stable\nVersion: latest\n\n",
		},
		{
			Name: "Given custom config, When channel view is executed, Then edge channel is returned",
			Args:             []string{"system", "channel", "view", "--config", testConfigPath},
			ExpectedResponse: "Context: local2\nChannel: edge\nVersion: latest\n\n",
		},
		{
			Name: "Given verbose flag, When channel view is executed, Then config path is printed with output",
			Args: 			  []string{"system", "channel", "view", "--verbose"},
			ExpectedResponse: fmt.Sprintf(
				"Using config file:%s\nContext: local\nChannel: stable\nVersion: latest\n\n",
				utils.DefaultConfigPath,
			),
		},
		{
			Name: "Given verbose flag and custom config, When channel view is executed, Then correct config path and output are shown",
			Args:				[]string{"system", "channel", "view", "-v", "--config", testConfigPath},
			ExpectedResponse:	fmt.Sprintf(
				"Using config file:%s\nContext: local2\nChannel: edge\nVersion: latest\n\n",
				testConfigPath,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			resetFlags()
			viper.Reset()

			b := utils.SetupMeshkitLoggerTesting(t, false)
			RootCmd.SetOut(b)
			RootCmd.SetArgs(tt.Args)

			err := RootCmd.Execute()
			if err != nil {
				t.Error(err)
			}

			if tt.ExpectedResponse == "" {
				return
			}

			actualResponse := b.String()
			expectedResponse := tt.ExpectedResponse

			utils.Equals(t, expectedResponse, actualResponse)
		})
	}

	if err := os.RemoveAll(utils.MesheryFolder); err != nil {
		t.Fatal("Couldn't remove files generated while testing")
	}
}
