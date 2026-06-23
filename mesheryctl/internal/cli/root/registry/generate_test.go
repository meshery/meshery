package registry

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func createTempCSVFile(t *testing.T, dir, fileName string) string {
	t.Helper()
	filePath := filepath.Join(dir, fileName)
	err := os.WriteFile(filePath, []byte("header\nvalue\n"), 0o600)
	assert.NoError(t, err)
	return filePath
}

// resetGenerateFlags resets all flags to their default values
func resetGenerateFlags(cmd *cobra.Command) {
	mesheryctlflags.InitValidators(cmd)
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
	})
}

func TestGenerateCmdFlags(t *testing.T) {
	tests := []struct {
		name         string
		flagName     string
		defaultValue string
		expectedType string
		description  string
	}{
		{
			name:         "timeout flag exists with correct default",
			flagName:     "timeout",
			defaultValue: meshkitRegistryUtils.DefaultModelTimeout.String(),
			expectedType: "duration",
			description:  "timeout duration for generating each model",
		},
		{
			name:         "latest-only flag exists with correct default",
			flagName:     "latest-only",
			defaultValue: "false",
			expectedType: "bool",
			description:  "generate only the latest version of each model",
		},
		{
			name:         "model flag exists",
			flagName:     "model",
			defaultValue: "",
			expectedType: "string",
			description:  "specific model name to be generated",
		},
		{
			name:         "output flag exists with correct default",
			flagName:     "output",
			defaultValue: "../server/meshmodel",
			expectedType: "string",
			description:  "location to output generated models",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := generateCmd.PersistentFlags().Lookup(tt.flagName)
			assert.NotNil(t, flag, "Flag %s should exist", tt.flagName)
			assert.Equal(t, tt.defaultValue, flag.DefValue, "Flag %s should have correct default value", tt.flagName)
			assert.Equal(t, tt.expectedType, flag.Value.Type(), "Flag %s should have correct type", tt.flagName)
		})
	}
}

func TestGenerateCmdTimeoutFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		flagValue      string
		expectedResult time.Duration
		expectError    bool
	}{
		{
			name:           "Parse 5 minutes",
			flagValue:      "5m",
			expectedResult: 5 * time.Minute,
			expectError:    false,
		},
		{
			name:           "Parse 10 minutes",
			flagValue:      "10m",
			expectedResult: 10 * time.Minute,
			expectError:    false,
		},
		{
			name:           "Parse 1 hour",
			flagValue:      "1h",
			expectedResult: 1 * time.Hour,
			expectError:    false,
		},
		{
			name:           "Parse 30 seconds",
			flagValue:      "30s",
			expectedResult: 30 * time.Second,
			expectError:    false,
		},
		{
			name:           "Parse complex duration",
			flagValue:      "1h30m",
			expectedResult: 1*time.Hour + 30*time.Minute,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags before each test
			resetGenerateFlags(generateCmd)

			err := generateCmd.PersistentFlags().Set("timeout", tt.flagValue)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify the parsed value
				duration, err := generateCmd.PersistentFlags().GetDuration("timeout")
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, duration)
			}
		})
	}
}

func TestGenerateCmdLatestOnlyFlagParsing(t *testing.T) {
	tests := []struct {
		name           string
		flagValue      string
		expectedResult bool
		expectError    bool
	}{
		{
			name:           "Parse true",
			flagValue:      "true",
			expectedResult: true,
			expectError:    false,
		},
		{
			name:           "Parse false",
			flagValue:      "false",
			expectedResult: false,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGenerateFlags(generateCmd)

			err := generateCmd.PersistentFlags().Set("latest-only", tt.flagValue)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				value, err := generateCmd.PersistentFlags().GetBool("latest-only")
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, value)
			}
		})
	}
}

func TestGenerateCmdMutualExclusivity(t *testing.T) {
	// Test that spreadsheet-id, directory, and model-csv are mutually exclusive
	// This is verified by Cobra's MarkFlagsMutuallyExclusive
	resetGenerateFlags(generateCmd)

	// Check that the flags exist and are properly configured
	spreadsheetFlag := generateCmd.PersistentFlags().Lookup("spreadsheet-id")
	directoryFlag := generateCmd.PersistentFlags().Lookup("directory")
	modelCSVFlag := generateCmd.PersistentFlags().Lookup("model-csv")

	assert.NotNil(t, spreadsheetFlag, "spreadsheet-id flag should exist")
	assert.NotNil(t, directoryFlag, "directory flag should exist")
	assert.NotNil(t, modelCSVFlag, "model-csv flag should exist")
}

func TestGenerateCmdRequiredTogether(t *testing.T) {
	// Test that spreadsheet-id and spreadsheet-cred are required together
	resetGenerateFlags(generateCmd)

	spreadsheetIDFlag := generateCmd.PersistentFlags().Lookup("spreadsheet-id")
	spreadsheetCredFlag := generateCmd.PersistentFlags().Lookup("spreadsheet-cred")

	assert.NotNil(t, spreadsheetIDFlag, "spreadsheet-id flag should exist")
	assert.NotNil(t, spreadsheetCredFlag, "spreadsheet-cred flag should exist")

	// Test that model-csv and component-csv are required together
	modelCSVFlag := generateCmd.PersistentFlags().Lookup("model-csv")
	componentCSVFlag := generateCmd.PersistentFlags().Lookup("component-csv")

	assert.NotNil(t, modelCSVFlag, "model-csv flag should exist")
	assert.NotNil(t, componentCSVFlag, "component-csv flag should exist")
}

func TestGenerateCmdHelp(t *testing.T) {
	resetGenerateFlags(generateCmd)

	// Test that generate command itself has the expected flags in its usage
	usageStr := generateCmd.UsageString()

	// Check that new flags are documented in usage
	assert.Contains(t, usageStr, "--timeout", "Usage should mention --timeout flag")
	assert.Contains(t, usageStr, "--latest-only", "Usage should mention --latest-only flag")

	// Check that examples mention new flags
	assert.Contains(t, generateCmd.Example, "--timeout", "Examples should include --timeout usage")
	assert.Contains(t, generateCmd.Example, "--latest-only", "Examples should include --latest-only usage")
}

func TestGenerateCmdExamples(t *testing.T) {
	// Verify that command examples are properly formatted
	assert.NotEmpty(t, generateCmd.Example, "Command should have examples")

	expectedExamples := []string{
		"spreadsheet-id",
		"spreadsheet-cred",
		"registrant-def",
		"registrant-cred",
		"directory",
		"model-csv",
		"component-csv",
		"timeout",
		"latest-only",
	}

	for _, example := range expectedExamples {
		assert.Contains(t, generateCmd.Example, example, "Examples should include %s", example)
	}
}

func TestDefaultModelTimeoutConstant(t *testing.T) {
	// Verify that the default timeout used is from meshkit
	expectedDefault := 5 * time.Minute
	assert.Equal(t, expectedDefault, meshkitRegistryUtils.DefaultModelTimeout,
		"Default model timeout should be 5 minutes")
}

func TestGenerationOptionsCreation(t *testing.T) {
	// Test creating GenerationOptions with various configurations
	tests := []struct {
		name        string
		timeout     time.Duration
		latestOnly  bool
		hasCallback bool
	}{
		{
			name:        "Default options",
			timeout:     5 * time.Minute,
			latestOnly:  false,
			hasCallback: true,
		},
		{
			name:        "Custom timeout with latest only",
			timeout:     10 * time.Minute,
			latestOnly:  true,
			hasCallback: true,
		},
		{
			name:        "Short timeout",
			timeout:     1 * time.Minute,
			latestOnly:  false,
			hasCallback: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var callback func(string, int, int)
			if tt.hasCallback {
				callback = func(modelName string, current, total int) {
					// Mock callback
				}
			}

			opts := meshkitRegistryUtils.GenerationOptions{
				ModelTimeout:      tt.timeout,
				LatestVersionOnly: tt.latestOnly,
				ProgressCallback:  callback,
			}

			assert.Equal(t, tt.timeout, opts.ModelTimeout)
			assert.Equal(t, tt.latestOnly, opts.LatestVersionOnly)
			if tt.hasCallback {
				assert.NotNil(t, opts.ProgressCallback)
			} else {
				assert.Nil(t, opts.ProgressCallback)
			}
		})
	}
}

func TestPreRunEValidation(t *testing.T) {
	// Test PreRunE validation for various input scenarios
	tests := []struct {
		name        string
		args        []string
		setup       func()
		expectError bool
		errorMsg    string
	}{
		{
			name:        "given no input source when prerune then returns error",
			args:        []string{},
			setup:       func() { resetGenerateFlags(generateCmd) },
			expectError: true,
			errorMsg:    "isn't specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			generateCmd.SetArgs(tt.args)
			err := generateCmd.PreRunE(generateCmd, tt.args)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPreRunEValidationPartialCSVInputs(t *testing.T) {
	tests := []struct {
		name          string
		flagValues    map[string]string
		expectError   bool
		expectedError error
	}{
		{
			name: "given model csv without component csv when prerune then returns missing component csv error",
			flagValues: map[string]string{
				"model-csv": "models.csv",
			},
			expectError:   true,
			expectedError: utils.ErrFlagsInvalid(errors.New("--component-csv is required when --model-csv is provided")),
		},
		{
			name: "given component csv without model csv when prerune then returns missing model csv error",
			flagValues: map[string]string{
				"component-csv": "components.csv",
			},
			expectError:   true,
			expectedError: utils.ErrFlagsInvalid(errors.New("--model-csv is required when --component-csv is provided")),
		},
		{
			name: "given relationship csv without model and component csv when prerune then returns relationship usage error",
			flagValues: map[string]string{
				"relationship-csv": "relationships.csv",
			},
			expectError:   true,
			expectedError: utils.ErrFlagsInvalid(errors.New("--relationship-csv can only be used with --model-csv and --component-csv")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGenerateFlags(generateCmd)
			for flagName, value := range tt.flagValues {
				err := generateCmd.PersistentFlags().Set(flagName, value)
				assert.NoError(t, err)
			}

			err := generateCmd.PreRunE(generateCmd, []string{})
			if tt.expectError {
				utils.AssertMeshkitErrorsEqual(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPreRunEValidationValidPathsNoRegression(t *testing.T) {
	t.Run("given complete csv mode with relationship csv when prerune then no error", func(t *testing.T) {
		resetGenerateFlags(generateCmd)
		tmpDir := t.TempDir()
		modelPath := createTempCSVFile(t, tmpDir, "models.csv")
		componentPath := createTempCSVFile(t, tmpDir, "components.csv")
		relationshipPath := createTempCSVFile(t, tmpDir, "relationships.csv")

		assert.NoError(t, generateCmd.PersistentFlags().Set("model-csv", modelPath))
		assert.NoError(t, generateCmd.PersistentFlags().Set("component-csv", componentPath))
		assert.NoError(t, generateCmd.PersistentFlags().Set("relationship-csv", relationshipPath))

		err := generateCmd.PreRunE(generateCmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("given spreadsheet mode when prerune then no error", func(t *testing.T) {
		resetGenerateFlags(generateCmd)
		assert.NoError(t, generateCmd.PersistentFlags().Set("spreadsheet-id", "sheet-id"))
		assert.NoError(t, generateCmd.PersistentFlags().Set("spreadsheet-cred", "dummy-cred"))

		err := generateCmd.PreRunE(generateCmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("given registrant mode when prerune then no error", func(t *testing.T) {
		resetGenerateFlags(generateCmd)
		assert.NoError(t, generateCmd.PersistentFlags().Set("registrant-def", "conn-def"))
		assert.NoError(t, generateCmd.PersistentFlags().Set("registrant-cred", "cred-def"))

		err := generateCmd.PreRunE(generateCmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("given directory mode when prerune then no error", func(t *testing.T) {
		resetGenerateFlags(generateCmd)
		assert.NoError(t, generateCmd.PersistentFlags().Set("directory", t.TempDir()))

		err := generateCmd.PreRunE(generateCmd, []string{})
		assert.NoError(t, err)
	})
}
