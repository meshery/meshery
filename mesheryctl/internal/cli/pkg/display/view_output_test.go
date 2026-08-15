package display

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

type viewOutputTestStruct struct {
	Name string `json:"name" yaml:"name"`
}

func TestFormatAndSaveOutput_WhenDisplayOnly_ThenWritesToConfiguredWriter(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	buf := &bytes.Buffer{}
	err := FormatAndSaveOutput(viewOutputTestStruct{Name: "meshery"}, "yaml", buf, false, "")
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "name: meshery")
}

func TestFormatAndSaveOutput_WhenSaveEnabled_ThenWritesToFile(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	tmpFile, err := os.CreateTemp("", "meshery_view_output_*.json")
	assert.NoError(t, err)
	tmpFilePath := tmpFile.Name()
	_ = tmpFile.Close()
	defer func() {
		assert.NoError(t, os.Remove(tmpFilePath))
	}()

	buf := &bytes.Buffer{}
	err = FormatAndSaveOutput(viewOutputTestStruct{Name: "meshery"}, "json", buf, true, tmpFilePath)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "\"name\": \"meshery\"")

	content, err := os.ReadFile(tmpFilePath)
	assert.NoError(t, err)

	var got viewOutputTestStruct
	err = json.Unmarshal(content, &got)
	assert.NoError(t, err)
	assert.Equal(t, "meshery", got.Name)
}

func TestFormatAndSaveOutput_WhenFormatIsInvalid_ThenReturnsUnsupportedFormatError(t *testing.T) {
	err := FormatAndSaveOutput(viewOutputTestStruct{Name: "meshery"}, "invalid", nil, false, "")
	utils.AssertMeshkitErrorsEqual(t, err, ErrUnsupportedFormat("invalid"))
}
