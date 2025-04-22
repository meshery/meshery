package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasCSVs(t *testing.T) {
	tests := []struct {
		name           string
		dirPath        string
		expectedResult bool
	}{
		{
			name:           "directory with CSVs",
			dirPath:        "./fixtures/with_csvs",
			expectedResult: true,
		},
		{
			name:           "directory without CSVs",
			dirPath:        "./fixtures/without_csvs",
			expectedResult: false,
		},
		{
			name:           "non-existent directory",
			dirPath:        "./fixtures/invalid_path",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := hasCSVs(tc.dirPath)

			assert.Equal(t, tc.expectedResult, res)
		})
	}
}
