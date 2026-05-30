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

package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectEntityType(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected string
	}{
		{
			name:     "model definition JSON",
			data:     `{"schemaVersion":"models.meshery.io/v1beta1","name":"test-model","version":"v0.1.0"}`,
			expected: "model",
		},
		{
			name:     "component definition JSON",
			data:     `{"schemaVersion":"components.meshery.io/v1beta1","displayName":"TestComp"}`,
			expected: "component",
		},
		{
			name:     "relationship definition JSON",
			data:     `{"schemaVersion":"relationships.meshery.io/v1alpha3","kind":"edge"}`,
			expected: "relationship",
		},
		{
			name:     "unknown schemaVersion",
			data:     `{"schemaVersion":"unknown.meshery.io/v1"}`,
			expected: "",
		},
		{
			name:     "missing schemaVersion",
			data:     `{"name":"test"}`,
			expected: "",
		},
		{
			name:     "invalid JSON",
			data:     `not json`,
			expected: "",
		},
		{
			name:     "empty data",
			data:     ``,
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := DetectEntityType([]byte(tc.data))
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestAnyEntityFilter(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected bool
	}{
		{
			name:     "model entity passes filter",
			data:     `{"schemaVersion":"models.meshery.io/v1beta1"}`,
			expected: true,
		},
		{
			name:     "component entity passes filter",
			data:     `{"schemaVersion":"components.meshery.io/v1beta1"}`,
			expected: true,
		},
		{
			name:     "relationship entity passes filter",
			data:     `{"schemaVersion":"relationships.meshery.io/v1alpha3"}`,
			expected: true,
		},
		{
			name:     "unknown entity fails filter",
			data:     `{"schemaVersion":"unknown.meshery.io/v1"}`,
			expected: false,
		},
		{
			name:     "non-entity data fails filter",
			data:     `{"name":"just some data"}`,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := AnyEntityFilter([]byte(tc.data))
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRelationshipOnlyFilter(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected bool
	}{
		{
			name:     "relationship passes filter",
			data:     `{"schemaVersion":"relationships.meshery.io/v1alpha3"}`,
			expected: true,
		},
		{
			name:     "model fails filter",
			data:     `{"schemaVersion":"models.meshery.io/v1beta1"}`,
			expected: false,
		},
		{
			name:     "component fails filter",
			data:     `{"schemaVersion":"components.meshery.io/v1beta1"}`,
			expected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := RelationshipOnlyFilter([]byte(tc.data))
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateHTTPScheme(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		expectErr bool
	}{
		{
			name:      "https URL is valid",
			url:       "https://example.com/model.json",
			expectErr: false,
		},
		{
			name:      "http URL is valid",
			url:       "http://example.com/model.json",
			expectErr: false,
		},
		{
			name:      "ftp URL is invalid",
			url:       "ftp://example.com/model.json",
			expectErr: true,
		},
		{
			name:      "file URL is invalid",
			url:       "file:///etc/passwd",
			expectErr: true,
		},
		{
			name:      "HTTP uppercase is valid",
			url:       "HTTP://example.com/model.json",
			expectErr: false,
		},
		{
			name:      "HTTPS uppercase is valid",
			url:       "HTTPS://example.com/model.json",
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateHTTPScheme(tc.url)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSortedPaths(t *testing.T) {
	files := map[string][]byte{
		"c.json": []byte("c"),
		"a.json": []byte("a"),
		"b.json": []byte("b"),
	}

	result := SortedPaths(files)
	assert.Equal(t, []string{"a.json", "b.json", "c.json"}, result)
}

func TestSortedPathsEmpty(t *testing.T) {
	files := map[string][]byte{}
	result := SortedPaths(files)
	assert.Empty(t, result)
}

func TestCollectFiles_SingleFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "model.json")
	content := `{"schemaVersion":"models.meshery.io/v1beta1","name":"test"}`
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

	result, err := CollectFiles([]string{filePath}, AnyEntityFilter)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, []byte(content), result[filePath])
}

func TestCollectFiles_ExplicitFileAlwaysIncluded(t *testing.T) {
	// An explicitly specified file should always be included, even if the
	// filter would reject it (filter only applies to directory traversal)
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "non-entity.json")
	content := `{"name":"not-an-entity"}`
	require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))

	result, err := CollectFiles([]string{filePath}, AnyEntityFilter)
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, []byte(content), result[filePath])
}

func TestCollectFiles_Directory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a model file
	modelContent := `{"schemaVersion":"models.meshery.io/v1beta1","name":"test"}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "model.json"), []byte(modelContent), 0644))

	// Create a component file
	compContent := `{"schemaVersion":"components.meshery.io/v1beta1","displayName":"TestComp"}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "component.yaml"), []byte(compContent), 0644))

	// Create a non-entity file that should be filtered out by AnyEntityFilter
	nonEntityContent := `{"name":"just-data"}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "other.json"), []byte(nonEntityContent), 0644))

	// Create a non-JSON/YAML file that should be skipped
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("hello"), 0644))

	result, err := CollectFiles([]string{tmpDir}, AnyEntityFilter)
	require.NoError(t, err)

	// Should include model.json and component.yaml, skip other.json (no entity schemaVersion) and readme.txt (wrong ext)
	assert.Len(t, result, 2)
	assert.Contains(t, result, filepath.Join(tmpDir, "model.json"))
	assert.Contains(t, result, filepath.Join(tmpDir, "component.yaml"))
}

func TestCollectFiles_DirectoryWithRelationshipFilter(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a model file
	modelContent := `{"schemaVersion":"models.meshery.io/v1beta1","name":"test"}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "model.json"), []byte(modelContent), 0644))

	// Create a relationship file
	relContent := `{"schemaVersion":"relationships.meshery.io/v1alpha3","kind":"edge"}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "relationship.json"), []byte(relContent), 0644))

	result, err := CollectFiles([]string{tmpDir}, RelationshipOnlyFilter)
	require.NoError(t, err)

	// Only the relationship file should be included
	assert.Len(t, result, 1)
	assert.Contains(t, result, filepath.Join(tmpDir, "relationship.json"))
}

func TestCollectFiles_RecursiveDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	require.NoError(t, os.MkdirAll(subDir, 0755))

	// Create files in root and subdirectory
	relContent := `{"schemaVersion":"relationships.meshery.io/v1alpha3","kind":"edge"}`
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "rel1.json"), []byte(relContent), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(subDir, "rel2.json"), []byte(relContent), 0644))

	result, err := CollectFiles([]string{tmpDir}, RelationshipOnlyFilter)
	require.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestCollectFiles_NonexistentPath(t *testing.T) {
	_, err := CollectFiles([]string{"/nonexistent/path/file.json"}, AnyEntityFilter)
	assert.Error(t, err)
}

func TestCollectFiles_EmptyPaths(t *testing.T) {
	result, err := CollectFiles([]string{}, AnyEntityFilter)
	require.NoError(t, err)
	assert.Empty(t, result)
}
