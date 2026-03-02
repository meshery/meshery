package root

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
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

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// update all locations
	fixturesDirPath := filepath.Join(currDir, "fixtures")
	utils.MesheryFolder = filepath.Join(fixturesDirPath, ".meshery")
	utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
	testConfigPath := filepath.Join(fixturesDirPath, "testConfig.yaml")
	// Explicitly create config for integration tests
	token := config.Token{
		Name:     utils.TemplateToken.Name,
		Location: utils.TemplateToken.Location,
	}

	ctx := config.Context{
		Endpoint:   utils.TemplateContext.Endpoint,
		Platform:   utils.TemplateContext.Platform,
		Channel:    utils.TemplateContext.Channel,
		Version:    utils.TemplateContext.Version,
		Components: utils.TemplateContext.Components,
		Provider:   utils.TemplateContext.Provider,
	}

	_ = config.MutateConfigIfNeeded(utils.DefaultConfigPath, utils.MesheryFolder, token, ctx)

	tests := []RootCmdTestInput{
		{
			Name: "Gets version",
			Args: []string{"version", "--config", utils.DefaultConfigPath},
		},
		{
			Name:             "view the channel to test if meshconfig is correctly created & loaded",
			Args:             []string{"system", "channel", "view", "--config", utils.DefaultConfigPath},
			ExpectedResponse: "Context: local\nChannel: stable\nVersion: latest\n\n",
		},
		{
			Name:             "view the channel with other config",
			Args:             []string{"system", "channel", "view", "--config", testConfigPath},
			ExpectedResponse: "Context: local2\nChannel: edge\nVersion: latest\n\n",
		},
		{
			Name:             "view the channel with verbose flag",
			Args:             []string{"system", "channel", "view", "--verbose"},
			ExpectedResponse: fmt.Sprintf("Using config file:%s\nContext: local\nChannel: stable\nVersion: latest\n\n", utils.DefaultConfigPath),
		},
		{
			Name:             "view the channel with short verbose flag and different config",
			Args:             []string{"system", "channel", "view", "-v", "--config", testConfigPath},
			ExpectedResponse: fmt.Sprintf("Using config file:%s\nContext: local2\nChannel: edge\nVersion: latest\n\n", testConfigPath),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// So that flag values from any previous test are not used.
			resetFlags()
			viper.Reset()

			b := utils.SetupMeshkitLoggerTesting(t, false)
			RootCmd.SetOut(b)
			RootCmd.SetArgs(tt.Args)

			err := RootCmd.Execute()
			RootCmd.InitDefaultHelpCmd()
			RootCmd.InitDefaultHelpFlag()
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
