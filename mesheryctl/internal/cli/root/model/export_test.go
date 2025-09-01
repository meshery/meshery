package model

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock HTTP client for testing - implements HTTPClient interface
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	args := m.Called(method, url, body)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *MockHTTPClient) MakeRequest(req *http.Request) (*http.Response, error) {
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

// Helper function to create a mock response
func createMockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
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
			errContains: "output location does not exist",
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

func TestExportWithClient(t *testing.T) {
	tmpDir := createTempDir(t)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name            string
		modelName       string
		url             string
		output          *outputDetail
		setupMock       func(*MockHTTPClient)
		wantErr         bool
		errContains     string
		checkFile       bool
		expectedFile    string
		expectedContent string
	}{
		{
			name:      "successful export with oci format",
			modelName: "test-model",
			url:       "http://localhost:9081/api/meshmodels/export",
			output: &outputDetail{
				Format: "yaml",
				Type:   "oci",
				Path:   tmpDir,
			},
			setupMock: func(m *MockHTTPClient) {
				req, _ := http.NewRequest("GET", "http://localhost:9081/api/meshmodels/export", nil)
				m.On("NewRequest", "GET", "http://localhost:9081/api/meshmodels/export", mock.Anything).Return(req, nil)
				m.On("MakeRequest", req).Return(createMockResponse(200, "mock oci model data"), nil)
			},
			wantErr:         false,
			checkFile:       true,
			expectedFile:    "test-model.tar",
			expectedContent: "mock oci model data",
		},
		{
			name:      "successful export with tar format",
			modelName: "test-model",
			url:       "http://localhost:9081/api/meshmodels/export",
			output: &outputDetail{
				Format: "json",
				Type:   "tar",
				Path:   tmpDir,
			},
			setupMock: func(m *MockHTTPClient) {
				req, _ := http.NewRequest("GET", "http://localhost:9081/api/meshmodels/export", nil)
				m.On("NewRequest", "GET", "http://localhost:9081/api/meshmodels/export", mock.Anything).Return(req, nil)
				m.On("MakeRequest", req).Return(createMockResponse(200, "mock tar model data"), nil)
			},
			wantErr:         false,
			checkFile:       true,
			expectedFile:    "test-model.tar.gz",
			expectedContent: "mock tar model data",
		},
		{
			name:      "server returns error status",
			modelName: "test-model",
			url:       "http://localhost:9081/api/meshmodels/export",
			output: &outputDetail{
				Format: "yaml",
				Type:   "oci",
				Path:   tmpDir,
			},
			setupMock: func(m *MockHTTPClient) {
				req, _ := http.NewRequest("GET", "http://localhost:9081/api/meshmodels/export", nil)
				m.On("NewRequest", "GET", "http://localhost:9081/api/meshmodels/export", mock.Anything).Return(req, nil)
				m.On("MakeRequest", req).Return(createMockResponse(404, "Not Found"), nil)
			},
			wantErr:     true,
			errContains: "failed to export model: status 404",
		},
		{
			name:      "request creation fails",
			modelName: "test-model",
			url:       "http://localhost:9081/api/meshmodels/export",
			output: &outputDetail{
				Format: "yaml",
				Type:   "oci",
				Path:   tmpDir,
			},
			setupMock: func(m *MockHTTPClient) {
				m.On("NewRequest", "GET", "http://localhost:9081/api/meshmodels/export", mock.Anything).Return((*http.Request)(nil), fmt.Errorf("request creation failed"))
			},
			wantErr:     true,
			errContains: "request creation failed",
		},
		{
			name:      "request execution fails",
			modelName: "test-model",
			url:       "http://localhost:9081/api/meshmodels/export",
			output: &outputDetail{
				Format: "yaml",
				Type:   "oci",
				Path:   tmpDir,
			},
			setupMock: func(m *MockHTTPClient) {
				req, _ := http.NewRequest("GET", "http://localhost:9081/api/meshmodels/export", nil)
				m.On("NewRequest", "GET", "http://localhost:9081/api/meshmodels/export", mock.Anything).Return(req, nil)
				m.On("MakeRequest", req).Return((*http.Response)(nil), fmt.Errorf("network error"))
			},
			wantErr:     true,
			errContains: "network error",
		},
		{
			name:      "unsupported output type",
			modelName: "test-model",
			url:       "http://localhost:9081/api/meshmodels/export",
			output: &outputDetail{
				Format: "yaml",
				Type:   "invalid",
				Path:   tmpDir,
			},
			setupMock: func(m *MockHTTPClient) {
				req, _ := http.NewRequest("GET", "http://localhost:9081/api/meshmodels/export", nil)
				m.On("NewRequest", "GET", "http://localhost:9081/api/meshmodels/export", mock.Anything).Return(req, nil)
				m.On("MakeRequest", req).Return(createMockResponse(200, "mock model data"), nil)
			},
			wantErr:     true,
			errContains: "unsupported output type: invalid",
		},
		{
			name:      "invalid output path",
			modelName: "test-model",
			url:       "http://localhost:9081/api/meshmodels/export",
			output: &outputDetail{
				Format: "yaml",
				Type:   "oci",
				Path:   "/invalid/path/that/does/not/exist",
			},
			setupMock: func(m *MockHTTPClient) {
				req, _ := http.NewRequest("GET", "http://localhost:9081/api/meshmodels/export", nil)
				m.On("NewRequest", "GET", "http://localhost:9081/api/meshmodels/export", mock.Anything).Return(req, nil)
				m.On("MakeRequest", req).Return(createMockResponse(200, "mock model data"), nil)
			},
			wantErr:     true,
			errContains: "no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockHTTPClient)
			tt.setupMock(mockClient)

			err := exportWithClient(tt.modelName, tt.url, tt.output, mockClient)

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
					assert.Equal(t, tt.expectedContent, string(content))
				}
			}

			// Verify that all expected calls were made
			mockClient.AssertExpectations(t)
		})
	}
}

// Test the backward compatibility export function
func TestExport(t *testing.T) {
	tmpDir := createTempDir(t)
	defer os.RemoveAll(tmpDir)

	// Create a mock client
	mockClient := new(MockHTTPClient)

	// Temporarily replace the default client for this test
	originalClient := defaultHTTPClient
	defaultHTTPClient = mockClient
	defer func() {
		defaultHTTPClient = originalClient
	}()

	// Setup mock expectations
	req, _ := http.NewRequest("GET", "http://test-url", nil)
	mockClient.On("NewRequest", "GET", "http://test-url", mock.Anything).Return(req, nil)
	mockClient.On("MakeRequest", req).Return(createMockResponse(200, "test data"), nil)

	output := &outputDetail{
		Format: "yaml",
		Type:   "oci",
		Path:   tmpDir,
	}

	err := export("test-model", "http://test-url", output)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)

	// Verify file was created
	filePath := filepath.Join(tmpDir, "test-model.tar")
	assert.FileExists(t, filePath)
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "test data", string(content))
}
