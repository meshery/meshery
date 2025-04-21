package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasCSVs(t *testing.T) {
	tests := []struct {
		name           string
		dirPath        string
		expectedError  bool
		expectedResult bool
	}{
		{
			name:           "directory with CSVs",
			dirPath:        "./fixtures/with_csvs",
			expectedError:  false,
			expectedResult: true,
		},
		{
			name:           "directory without CSVs",
			dirPath:        "./fixtures/without_csvs",
			expectedError:  false,
			expectedResult: false,
		},
		{
			name:           "non-existent directory",
			dirPath:        "./fixtures/invalid_path",
			expectedError:  true,
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res, err := hasCSVs(tc.dirPath)

			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectedResult, res)
		})
	}
}
