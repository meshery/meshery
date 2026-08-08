package registry

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/sheets/v4"
)

func TestCheckSpreadsheetFetch(t *testing.T) {
	tests := []struct {
		name        string
		resp        *sheets.Spreadsheet
		err         error
		expectError bool
		errorMsg    string
	}{
		{
			name:        "network error is propagated",
			resp:        nil,
			err:         errors.New("network failure"),
			expectError: true,
			errorMsg:    "network failure",
		},
		{
			name:        "nil response with nil error is treated as failure",
			resp:        nil,
			err:         nil,
			expectError: true,
			errorMsg:    "nil response",
		},
		{
			name: "non-200 status with nil err is reported as an error",
			resp: &sheets.Spreadsheet{
				ServerResponse: googleapi.ServerResponse{HTTPStatusCode: 403},
			},
			err:         nil,
			expectError: true,
			errorMsg:    "unexpected HTTP status 403",
		},
		{
			name: "429 rate limit is reported as an error",
			resp: &sheets.Spreadsheet{
				ServerResponse: googleapi.ServerResponse{HTTPStatusCode: 429},
			},
			err:         nil,
			expectError: true,
			errorMsg:    "unexpected HTTP status 429",
		},
		{
			name: "200 status returns no error",
			resp: &sheets.Spreadsheet{
				ServerResponse: googleapi.ServerResponse{HTTPStatusCode: 200},
			},
			err:         nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkSpreadsheetFetch(tt.resp, tt.err, "./models")

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
