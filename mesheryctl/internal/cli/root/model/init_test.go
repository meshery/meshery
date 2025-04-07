package model

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/model/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestModelInitYAML(t *testing.T) {
	tmpDir := "test-tmp-yaml"
	modelName := "test-model"
	version := "1.0.0"

	defer os.RemoveAll(tmpDir)
	_ = os.RemoveAll(tmpDir)

	root := ModelCmd
	root.SetArgs([]string{"init", modelName, "--output-format", "yaml", "--version", version, "--path", tmpDir})

	output, err := utils.ExecuteCommand(root)
	assert.NoError(t, err)
	assert.Contains(t, output, "Created new Meshery model")

	expectedDirs := []string{
		filepath.Join(tmpDir, modelName, version),
		filepath.Join(tmpDir, modelName, version, "components"),
		filepath.Join(tmpDir, modelName, version, "relationships"),
		filepath.Join(tmpDir, modelName, version, "connections"),
		filepath.Join(tmpDir, modelName, version, "credentials"),
	}
	for _, dir := range expectedDirs {
		assert.DirExists(t, dir)
	}

	assert.FileExists(t, filepath.Join(tmpDir, modelName, version, "model.yaml"))
}

func TestModelInitJSON(t *testing.T) {
	tmpDir := "test-tmp-json"
	modelName := "test-model"
	version := "1.0.0"

	defer os.RemoveAll(tmpDir)
	_ = os.RemoveAll(tmpDir)

	root := ModelCmd
	root.SetArgs([]string{"init", modelName, "--output-format", "json", "--version", version, "--path", tmpDir})

	output, err := utils.ExecuteCommand(root)
	assert.NoError(t, err)
	assert.Contains(t, output, "Created new Meshery model")

	assert.FileExists(t, filepath.Join(tmpDir, modelName, version, "model.json"))
}
