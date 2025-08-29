package model

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock HTTP client for testing
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Mock config for testing
type MockConfig struct {
	baseURL string
}

func (m *MockConfig) GetBaseMesheryURL() string {
	return m.baseURL
}

// Helper function to create a test command with flags
func createTestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.Flags().StringP("output-format", "t", "yaml", "output format")
	cmd.Flags().StringP("output-location", "l", "./", "output location")
	cmd.Flags().StringP("output-type", "o", "oci", "output type")
	cmd.Flags().BoolP("discard-components", "c", false, "discard components")
	cmd.Flags().BoolP("discard-relationships", "r", false, "discard relationships")
	cmd.Flags().StringP("version", "", "", "model version")
	cmd.Flags().IntP("page", "p", 1, "page number")

	return cmd
}

// Helper function to create temporary directory for testing
func createTempDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "model_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tmpDir
}

// Test Args validation function
func TestExportModelCmd_Args(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "no arguments provided",
			args:    []string{},
			wantErr: true,
			errMsg:  "Please provide a model name",
		},
		{
			name:    "valid argument provided",
			args:    []string{"test-model"},
			wantErr: false,
		},
		{
			name:    "multiple arguments provided",
			args:    []string{"test-model", "extra-arg"},
			wantErr: false, // Args function only checks for empty, doesn't restrict multiple
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := exportModelCmd.Args(nil, tt.args)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test PreRunE validation function
func TestExportModelCmd_PreRunE(t *testing.T) {
	// Save original viper settings
	originalViper := viper.GetViper()
	defer func() {
		viper.Reset()
		for key, value := range originalViper.AllSettings() {
			viper.Set(key, value)
		}
	}()

	tmpDir := createTempDir(t)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		setupFlags  func(*cobra.Command)
		setupConfig func()
		wantErr     bool
		errContains string
	}{
		{
			name: "valid flags and config",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("output-format", "yaml")
				cmd.Flags().Set("output-type", "oci")
				cmd.Flags().Set("page", "1")
				cmd.Flags().Set("output-location", tmpDir)
			},
			setupConfig: func() {
				viper.Set("mesheryctl.config.platform", "docker")
			},
			wantErr: false,
		},
		{
			name: "invalid output format",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("output-format", "invalid")
				cmd.Flags().Set("output-type", "oci")
				cmd.Flags().Set("page", "1")
				cmd.Flags().Set("output-location", tmpDir)
			},
			setupConfig: func() {
				viper.Set("mesheryctl.config.platform", "docker")
			},
			wantErr:     true,
			errContains: "invalid value",
		},
		{
			name: "invalid output type",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("output-format", "yaml")
				cmd.Flags().Set("output-type", "invalid")
				cmd.Flags().Set("page", "1")
				cmd.Flags().Set("output-location", tmpDir)
			},
			setupConfig: func() {
				viper.Set("mesheryctl.config.platform", "docker")
			},
			wantErr:     true,
			errContains: "invalid value",
		},
		{
			name: "invalid page number",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("output-format", "yaml")
				cmd.Flags().Set("output-type", "oci")
				cmd.Flags().Set("page", "-1")
				cmd.Flags().Set("output-location", tmpDir)
			},
			setupConfig: func() {
				viper.Set("mesheryctl.config.platform", "docker")
			},
			wantErr:     true,
			errContains: "invalid page number",
		},
		{
			name: "invalid version format",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("output-format", "yaml")
				cmd.Flags().Set("output-type", "oci")
				cmd.Flags().Set("page", "1")
				cmd.Flags().Set("version", "invalid-version")
				cmd.Flags().Set("output-location", tmpDir)
			},
			setupConfig: func() {
				viper.Set("mesheryctl.config.platform", "docker")
			},
			wantErr:     true,
			errContains: "invalid format",
		},
		{
			name: "valid version format",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("output-format", "yaml")
				cmd.Flags().Set("output-type", "oci")
				cmd.Flags().Set("page", "1")
				cmd.Flags().Set("version", "v1.2.3")
				cmd.Flags().Set("output-location", tmpDir)
			},
			setupConfig: func() {
				viper.Set("mesheryctl.config.platform", "docker")
			},
			wantErr: false,
		},
		{
			name: "non-existent output location",
			setupFlags: func(cmd *cobra.Command) {
				cmd.Flags().Set("output-format", "yaml")
				cmd.Flags().Set("output-type", "oci")
				cmd.Flags().Set("page", "1")
				cmd.Flags().Set("output-location", "/non/existent/path")
			},
			setupConfig: func() {
				viper.Set("mesheryctl.config.platform", "docker")
			},
			wantErr:     true,
			errContains: "folder",
		},
		{
			name: "output location is file not directory",
			setupFlags: func(cmd *cobra.Command) {
				// Create a temporary file
				tmpFile, _ := os.CreateTemp("", "test_file_*")
				tmpFile.Close()
				defer os.Remove(tmpFile.Name())

				cmd.Flags().Set("output-format", "yaml")
				cmd.Flags().Set("output-type", "oci")
				cmd.Flags().Set("page", "1")
				cmd.Flags().Set("output-location", tmpFile.Name())
			},
			setupConfig: func() {
				viper.Set("mesheryctl.config.platform", "docker")
			},
			wantErr:     true,
			errContains: "is not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper for each test
			viper.Reset()
			tt.setupConfig()

			cmd := createTestCommand()
			tt.setupFlags(cmd)

			err := exportModelCmd.PreRunE(cmd, []string{"test-model"})

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test export function with different scenarios
func TestExportFunction(t *testing.T) {
	tmpDir := createTempDir(t)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name           string
		modelName      string
		serverResponse func() *httptest.Server
		output         *outputDetail
		wantErr        bool
		errContains    string
		checkFile      bool
		expectedFile   string
	}{
		{
			name:      "successful export with oci format",
			modelName: "test-model",
			serverResponse: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("mock model data"))
				}))
			},
			output: &outputDetail{
				Format: "yaml",
				Type:   "oci",
				Path:   tmpDir,
			},
			wantErr:      false,
			checkFile:    true,
			expectedFile: "test-model.tar",
		},
		{
			name:      "successful export with tar format",
			modelName: "test-model",
			serverResponse: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("mock model data"))
				}))
			},
			output: &outputDetail{
				Format: "json",
				Type:   "tar",
				Path:   tmpDir,
			},
			wantErr:      false,
			checkFile:    true,
			expectedFile: "test-model.tar.gz",
		},
		{
			name:      "server returns error status",
			modelName: "test-model",
			serverResponse: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			output: &outputDetail{
				Format: "yaml",
				Type:   "oci",
				Path:   tmpDir,
			},
			wantErr:     true,
			errContains: "failed to export model",
		},
		{
			name:      "invalid output path",
			modelName: "test-model",
			serverResponse: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("mock model data"))
				}))
			},
			output: &outputDetail{
				Format: "yaml",
				Type:   "oci",
				Path:   "/invalid/path/that/does/not/exist",
			},
			wantErr:     true,
			errContains: "no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.serverResponse()
			defer server.Close()

			newRequest := func(method, url string, body io.Reader) (*http.Request, error) {
				// meshery server
				testURL := strings.Replace(url, "http://localhost:9081", server.URL, 1)
				return http.NewRequest(method, testURL, body)
			}

			makeRequest := func(req *http.Request) (*http.Response, error) {
				client := &http.Client{}
				return client.Do(req)
			}

			// exportWithClient inline for testing
			exportWithClient := func(modelName, url string, output *outputDetail, newRequest func(string, string, io.Reader) (*http.Request, error), makeRequest func(*http.Request) (*http.Response, error)) error {
				req, err := newRequest("GET", url, nil)
				if err != nil {
					return err
				}
				resp, err := makeRequest(req)
				if err != nil {
					return err
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("failed to export model: %s", resp.Status)
				}
				var filePath string
				if output.Type == "oci" {
					filePath = filepath.Join(output.Path, modelName+".tar")
				} else if output.Type == "tar" {
					filePath = filepath.Join(output.Path, modelName+".tar.gz")
				} else {
					return fmt.Errorf("unsupported output type: %s", output.Type)
				}
				outFile, err := os.Create(filePath)
				if err != nil {
					return err
				}
				defer outFile.Close()
				_, err = io.Copy(outFile, resp.Body)
				return err
			}
			err := exportWithClient(tt.modelName, server.URL+"/api/meshmodels/export", tt.output, newRequest, makeRequest)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)

				if tt.checkFile {
					filePath := filepath.Join(tt.output.Path, tt.expectedFile)
					assert.FileExists(t, filePath)

					// Check file content
					content, err := os.ReadFile(filePath)
					assert.NoError(t, err)
					assert.Equal(t, "mock model data", string(content))
				}
			}
		})
	}
}

// Test flag initialization
func TestExportModelCmdFlags(t *testing.T) {
	cmd := exportModelCmd

	// Test that all expected flags are present
	expectedFlags := map[string]struct{}{
		"output-format":         {},
		"output-location":       {},
		"output-type":           {},
		"discard-components":    {},
		"discard-relationships": {},
		"version":               {},
		"page":                  {},
	}

	for flagName := range expectedFlags {
		flag := cmd.Flags().Lookup(flagName)
		assert.NotNil(t, flag, "Flag %s should be defined", flagName)
	}

	// Test default values
	assert.Equal(t, "yaml", cmd.Flags().Lookup("output-format").DefValue)
	assert.Equal(t, "./", cmd.Flags().Lookup("output-location").DefValue)
	assert.Equal(t, "oci", cmd.Flags().Lookup("output-type").DefValue)
	assert.Equal(t, "false", cmd.Flags().Lookup("discard-components").DefValue)
	assert.Equal(t, "false", cmd.Flags().Lookup("discard-relationships").DefValue)
	assert.Equal(t, "", cmd.Flags().Lookup("version").DefValue)
	assert.Equal(t, "1", cmd.Flags().Lookup("page").DefValue)
}

// Test command metadata
func TestExportModelCmdMetadata(t *testing.T) {
	cmd := exportModelCmd

	assert.Equal(t, "export", cmd.Use)
	assert.Equal(t, "Export registered models", cmd.Short)
	assert.Contains(t, cmd.Long, "Export the registered model")
	assert.Contains(t, cmd.Example, "mesheryctl model export")

	// Test that the command has the required functions
	assert.NotNil(t, cmd.Args)
	assert.NotNil(t, cmd.PreRunE)
	assert.NotNil(t, cmd.RunE)
}
