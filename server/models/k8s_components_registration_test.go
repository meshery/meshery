// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLoadK8sMeshModelMetadata_Success verifies that a valid model template JSON file
// is correctly read, parsed into ModelTemplate, and that the underlying file descriptor
// is closed immediately upon return (Fixes #22752).
func TestLoadK8sMeshModelMetadata_Success(t *testing.T) {
	tempDir := t.TempDir()
	validJSON := `{
		"model": "kubernetes",
		"primaryColor": "#326CE5",
		"secondaryColor": "#7aa1f0",
		"shape": "circle"
	}`

	filePath := filepath.Join(tempDir, "model_template_valid.json")
	err := os.WriteFile(filePath, []byte(validJSON), 0600)
	require.NoError(t, err, "failed to write test file")

	metadata, err := loadK8sMeshModelMetadata(filePath)
	require.NoError(t, err, "loadK8sMeshModelMetadata should succeed for valid JSON")

	assert.Equal(t, "#326CE5", metadata.PrimaryColor)
	require.NotNil(t, metadata.SecondaryColor)
	assert.Equal(t, "#7aa1f0", *metadata.SecondaryColor)
	require.NotNil(t, metadata.Shape)
	assert.Equal(t, "circle", string(*metadata.Shape))

	// Verify the file descriptor was closed: on Windows, an open file descriptor
	// prevents removing or renaming the file. Successful deletion confirms closure.
	removeErr := os.Remove(filePath)
	assert.NoError(t, removeErr, "file descriptor should be closed, allowing removal")
}

// TestLoadK8sMeshModelMetadata_FileNotFound verifies that non-existent paths return
// an error gracefully without panic or resource leaks.
func TestLoadK8sMeshModelMetadata_FileNotFound(t *testing.T) {
	nonExistentPath := filepath.Join(t.TempDir(), "does_not_exist.json")

	metadata, err := loadK8sMeshModelMetadata(nonExistentPath)
	assert.Error(t, err, "expected error for non-existent file")
	assert.Empty(t, metadata.PrimaryColor)
	assert.Nil(t, metadata.SecondaryColor)
}

// TestLoadK8sMeshModelMetadata_CorruptedJSON verifies that malformed JSON files return
// an unmarshal error while still properly closing the open file descriptor via defer (Fixes #22752).
func TestLoadK8sMeshModelMetadata_CorruptedJSON(t *testing.T) {
	tempDir := t.TempDir()
	malformedJSON := `{"model": "kubernetes", "primaryColor": ` // incomplete JSON

	filePath := filepath.Join(tempDir, "model_template_corrupt.json")
	err := os.WriteFile(filePath, []byte(malformedJSON), 0600)
	require.NoError(t, err, "failed to write corrupt test file")

	metadata, err := loadK8sMeshModelMetadata(filePath)
	assert.Error(t, err, "expected error for corrupted JSON content")
	assert.Empty(t, metadata.PrimaryColor)

	// Verify deferred close executed despite JSON unmarshal error:
	removeErr := os.Remove(filePath)
	assert.NoError(t, removeErr, "file descriptor should be closed even on unmarshal error")
}

// TestLoadK8sMeshModelMetadata_LiveTemplate verifies that the actual project-relative
// template file at k8sMeshModelPath can be read and parsed.
func TestLoadK8sMeshModelMetadata_LiveTemplate(t *testing.T) {
	if _, err := os.Stat(k8sMeshModelPath); os.IsNotExist(err) {
		t.Skipf("skipping live template test: %s not found from current directory", k8sMeshModelPath)
	}

	metadata, err := loadK8sMeshModelMetadata(k8sMeshModelPath)
	require.NoError(t, err, "failed to load live model_template.json")
	assert.Equal(t, "#326CE5", metadata.PrimaryColor)
	require.NotNil(t, metadata.SecondaryColor)
	assert.Equal(t, "#7aa1f0", *metadata.SecondaryColor)
	require.NotNil(t, metadata.Shape)
	assert.Equal(t, "circle", string(*metadata.Shape))
}

// TestK8sMeshModelMetadata_PackageVar verifies that the package-level K8sMeshModelMetadata
// is initialized by init() when the template file is available.
func TestK8sMeshModelMetadata_PackageVar(t *testing.T) {
	if _, err := os.Stat(k8sMeshModelPath); os.IsNotExist(err) {
		t.Skipf("skipping package var test: %s not found from current directory", k8sMeshModelPath)
	}

	assert.Equal(t, "#326CE5", K8sMeshModelMetadata.PrimaryColor)
	require.NotNil(t, K8sMeshModelMetadata.SecondaryColor)
	assert.Equal(t, "#7aa1f0", *K8sMeshModelMetadata.SecondaryColor)
	require.NotNil(t, K8sMeshModelMetadata.Shape)
	assert.Equal(t, "circle", string(*K8sMeshModelMetadata.Shape))
}
