package utils

import (
	"errors"
	"net/http"
	"testing"

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
			name:          "successful version fetch",
			mockStatus:    http.StatusOK,
			mockBody:      "v0.8.0",
			expectVersion: "v0.8.0",
			expectError:   false,
		},
		{
			name:          "empty response body",
			mockStatus:    http.StatusOK,
			mockBody:      "",
			expectVersion: "",
			expectError:   false,
		},
		{
			name:          "whitespace in response",
			mockStatus:    http.StatusOK,
			mockBody:      "  v0.8.0  ",
			expectVersion: "v0.8.0",
			expectError:   false,
		},
		{
			name:          "server error returns error",
			mockStatus:    http.StatusInternalServerError,
			mockBody:      "Internal Server Error",
			expectVersion: "",
			expectError:   true,
		},
		{
			name:          "not found returns error",
			mockStatus:    http.StatusNotFound,
			mockBody:      "Not Found",
			expectVersion: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			t.Cleanup(httpmock.DeactivateAndReset)

			httpmock.RegisterResponder(
				http.MethodGet,
				latestVersionURL,
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
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder(
		http.MethodGet,
		latestVersionURL,
		httpmock.NewErrorResponder(errors.New("network error")),
	)

	_, err := GetLatestVersionForMesheryctl()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "network error")
}

func TestGetLatestVersionForMesheryctl_Non2xxStatus(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder(
		http.MethodGet,
		latestVersionURL,
		httpmock.NewStringResponder(http.StatusServiceUnavailable, "Service Unavailable"),
	)

	_, err := GetLatestVersionForMesheryctl()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "503")
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
			name:           "latest version matches build",
			build:          "v0.8.0",
			mockStatus:     http.StatusOK,
			mockBody:       "v0.8.0",
			expectContains: "is the latest release",
		},
		{
			name:           "new version available",
			build:          "v0.7.0",
			mockStatus:     http.StatusOK,
			mockBody:       "v0.8.0",
			expectContains: "A new release of mesheryctl is available",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			t.Cleanup(httpmock.DeactivateAndReset)

			httpmock.RegisterResponder(
				http.MethodGet,
				latestVersionURL,
				httpmock.NewStringResponder(tt.mockStatus, tt.mockBody),
			)

			buf := SetupMeshkitLoggerTesting(t, false)
			CheckMesheryctlClientVersion(tt.build)
			output := buf.String()

			assert.Contains(t, output, tt.expectContains)
		})
	}
}

func TestCheckMesheryctlClientVersion_RequestError(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder(
		http.MethodGet,
		latestVersionURL,
		httpmock.NewErrorResponder(errors.New("connection refused")),
	)

	buf := SetupMeshkitLoggerTesting(t, false)
	CheckMesheryctlClientVersion("v0.8.0")
	output := buf.String()

	assert.Contains(t, output, "unable to check for latest version")
	assert.Contains(t, output, "connection refused")
}

func TestCheckMesheryctlClientVersion_Non2xxResponse(t *testing.T) {
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder(
		http.MethodGet,
		latestVersionURL,
		httpmock.NewStringResponder(http.StatusInternalServerError, ""),
	)

	buf := SetupMeshkitLoggerTesting(t, false)
	CheckMesheryctlClientVersion("v0.8.0")
	output := buf.String()

	assert.Contains(t, output, "unable to check for latest version")
}