package display

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name    string `json:"name" yaml:"name"`
	Content string `json:"content" yaml:"content"`
}

func TestOutputFormatter_Display(t *testing.T) {
	data := testStruct{
		Name:    "Valid",
		Content: "This is fine.",
	}

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
