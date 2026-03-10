package display

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

type testStruct struct {
	Name    string `json:"name" yaml:"name"`
	Content string `json:"content" yaml:"content"`
}

var data = testStruct{
	Name:    "Valid",
	Content: "This is fine.",
}

func Test_Given_OutputFormatter_When_Display_Then_Content_Is_Displayed_Without_Error(t *testing.T) {

	testCases := []struct {
		name          string
		input         testStruct
		formatter     OutputFormatter[testStruct]
		buf           *bytes.Buffer
		expectError   bool
		expectedError error
	}{
		{
			name:        "Given json data when encoded then content is displayed without error",
			input:       data,
			formatter:   NewJSONOutputFormatter(data),
			buf:         &bytes.Buffer{},
			expectError: false,
		},
		{
			name:        "Given yaml data when encoded then content is displayed without error",
			input:       data,
			formatter:   NewYAMLOutputFormatter(data),
			buf:         &bytes.Buffer{},
			expectError: false,
		},
		// Add more cases for specific unsupported characters or encoding issues
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			formatter := tc.formatter.WithOutput(tc.buf)
			err := formatter.Display()
			assert.Nil(t, err)
			output := tc.buf.String()
			assert.Contains(t, output, data.Name, "Output should contain name")
			assert.Contains(t, output, data.Content, "Output should contain content")
			tc.buf.Reset()
		})
	}

}

func setUpTestLogger(t *testing.T) func(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)
	return func(t *testing.T) {
		utils.Log = nil
	}
}

func Test_Given_JSONOutputFormatterSaver_With_Filepath_When_Save_Then_File_Is_Created(t *testing.T) {
	tearDown := setUpTestLogger(t)
	defer tearDown(t)

	tmpFile, err := os.CreateTemp("", "meshery_output_*.json")
	assert.NoError(t, err)
	tmpFilePath := tmpFile.Name()
	_ = tmpFile.Close()
	defer os.Remove(tmpFilePath)

	jsonFormatter := JSONOutputFormatter[testStruct]{
		Data: data,
		EncoderSettings: JsonEncoderSettings{
			SetEscapeHTML: false,
			IndentPrefix:  "",
			IndentValue:   "  ",
		},
	}

	saver := NewJSONOutputFormatterSaver(jsonFormatter).WithFilePath(tmpFilePath)
	err = saver.Save()
	assert.NoError(t, err)

	content, err := os.ReadFile(tmpFilePath)
	assert.NoError(t, err)

	var got testStruct
	err = json.Unmarshal(content, &got)
	assert.NoError(t, err)
	assert.Equal(t, data.Name, got.Name)
	assert.Equal(t, data.Content, got.Content)
}

func Test_Given_JSONOutputFormatterSaver_With_NoFilepath_When_Save_Then_Error_Is_Returned(t *testing.T) {
	jsonFormatter := JSONOutputFormatter[testStruct]{
		Data: data,
		EncoderSettings: JsonEncoderSettings{
			SetEscapeHTML: false,
			IndentPrefix:  "",
			IndentValue:   "  ",
		},
	}

	saver := NewJSONOutputFormatterSaver(jsonFormatter) // no file path set
	err := saver.Save()
	assert.Error(t, err)
}

func Test_Given_YAMLOutputFormatterSaver_With_Filepath_When_Save_Then_File_Is_Created(t *testing.T) {
	tearDown := setUpTestLogger(t)
	defer tearDown(t)

	tmpFile, err := os.CreateTemp("", "meshery_output_*.yaml")
	assert.NoError(t, err)
	tmpFilePath := tmpFile.Name()
	_ = tmpFile.Close()
	defer os.Remove(tmpFilePath)

	yamlFormatter := YAMLOutputFormatter[testStruct]{
		Data: data,
	}

	saver := NewYAMLOutputFormatterSaver(yamlFormatter).WithFilePath(tmpFilePath)
	err = saver.Save()
	assert.NoError(t, err)

	content, err := os.ReadFile(tmpFilePath)
	assert.NoError(t, err)

	var got testStruct
	err = yaml.Unmarshal(content, &got)
	assert.NoError(t, err)
	assert.Equal(t, data.Name, got.Name)
	assert.Equal(t, data.Content, got.Content)
}

func Test_Given_YAMLOutputFormatterSaver_With_NoFilepath_When_Save_Then_Error_Is_Returned(t *testing.T) {
	yamlFormatter := YAMLOutputFormatter[testStruct]{
		Data: data,
	}

	saver := NewYAMLOutputFormatterSaver(yamlFormatter) // no file path set
	err := saver.Save()
	assert.Error(t, err)
}

func TestYAMLOutputFormatterSaver_Save_WriteError(t *testing.T) {
	// Use a directory path to force a write error (os.WriteFile will error on writing to directory)
	tmpDir, err := os.MkdirTemp("", "meshery_output_dir_*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	yamlFormatter := YAMLOutputFormatter[testStruct]{
		Data: data,
	}

	saver := NewYAMLOutputFormatterSaver(yamlFormatter).WithFilePath(tmpDir)
	err = saver.Save()
	assert.Error(t, err)
}
