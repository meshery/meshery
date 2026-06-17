package utils

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestVersionForMesheryctl(t *testing.T) {
	tests := []struct {
		name          string
		mockStatus    int
		mockBody      string
		expectVersion string
		expectError   bool
	}{
		{
			name:          "Successful version fetch",
			mockStatus:    http.StatusOK,
			mockBody:      "v0.8.0",
			expectVersion: "v0.8.0",
			expectError:   false,
		},
		{
			name:          "Empty response body",
			mockStatus:    http.StatusOK,
			mockBody:      "",
			expectVersion: "",
			expectError:   false,
		},
		{
			name:          "Whitespace in response",
			mockStatus:    http.StatusOK,
			mockBody:      "  v0.8.0  ",
			expectVersion: "v0.8.0",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder(
				http.MethodGet,
				"https://docs.meshery.io/project/releases/latest",
				httpmock.NewStringResponder(tt.mockStatus, tt.mockBody),
			)

			version, err := GetLatestVersionForMesheryctl()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectVersion, version)
			}
		})
	}
}

func TestGetLatestVersionForMesheryctl_NetworkError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		"https://docs.meshery.io/project/releases/latest",
		httpmock.NewErrorResponder(errors.New("network error")),
	)

	_, err := GetLatestVersionForMesheryctl()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "network error")
}

func TestCheckMesheryctlClientVersion(t *testing.T) {
	tests := []struct {
		name           string
		build          string
		mockStatus     int
		mockBody       string
		expectContains string
	}{
		{
			name:           "Latest version matches build",
			build:          "v0.8.0",
			mockStatus:     http.StatusOK,
			mockBody:       "v0.8.0",
			expectContains: "is the latest release",
		},
		{
			name:           "New version available",
			build:          "v0.7.0",
			mockStatus:     http.StatusOK,
			mockBody:       "v0.8.0",
			expectContains: "A new release of mesheryctl is available",
		},
		{
			name:           "HTTP error returns warning",
			build:          "v0.8.0",
			mockStatus:     http.StatusInternalServerError,
			mockBody:       "",
			expectContains: "unable to check for latest version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder(
				http.MethodGet,
				"https://docs.meshery.io/project/releases/latest",
				httpmock.NewStringResponder(tt.mockStatus, tt.mockBody),
			)

			buf := SetupMeshkitLoggerTesting(t, false)
			CheckMesheryctlClientVersion(tt.build)
			output := buf.String()

			assert.Contains(t, output, tt.expectContains)
		})
	}
}

func TestCheckMesheryctlClientVersion_TimeoutWarning(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		"https://docs.meshery.io/project/releases/latest",
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(50 * time.Millisecond)
			return nil, errors.New("timeout: request took too long")
		},
	)

	buf := SetupMeshkitLoggerTesting(t, false)
	CheckMesheryctlClientVersion("v0.8.0")
	output := buf.String()

	assert.Contains(t, output, "unable to check for latest version")
	assert.Contains(t, output, "timeout")
}