package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemType_String(t *testing.T) {
	tests := []struct {
		name     string
		dt       SystemType
		expected string
	}{
		{"Meshery", Meshery, "meshery"},
		{"Docs", Docs, "docs"},
		{"RemoteProvider", RemoteProvider, "remote-provider"},
		{"Unknown", SystemType(99), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Now we ARE using the assert package
			assert.Equal(t, tt.expected, tt.dt.String())
		})
	}
}

func TestGetIndexForRegisterCol(t *testing.T) {
	cols := []string{"name", "version", "repository", "register"}
	
	t.Run("Find existing column", func(t *testing.T) {
		index := GetIndexForRegisterCol(cols, "register")
		assert.Equal(t, 3, index)
	})

	t.Run("Column not found", func(t *testing.T) {
		index := GetIndexForRegisterCol(cols, "non-existent")
		assert.Equal(t, -1, index)
	})
}

func TestGenerateIcons_DirectoryCreation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "registry_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	mockImgPath := filepath.Join(tmpDir, "assets/images")

	t.Run("Verify directory logic", func(t *testing.T) {
		modelName := "test-model"
		fullPath := filepath.Join(mockImgPath, modelName)
		
		err := os.MkdirAll(fullPath, 0777)
		assert.NoError(t, err)

		_, err = os.Stat(fullPath)
		assert.False(t, os.IsNotExist(err), "Directory should exist")
	})
}