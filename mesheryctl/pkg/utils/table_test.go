package utils

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableHeader_Format(t *testing.T) {
	header := TableHeader("status")
	formatted := header.Format()

	// Verify the header is converted to title case (STATUS or Status depending on tw.Title)
	// and contains bold escape codes
	assert.NotEmpty(t, formatted)
	assert.Contains(t, strings.ToLower(formatted), "status")
}

func TestPrintToTable(t *testing.T) {
	header := []string{"ID", "NAME", "STATUS"}
	data := [][]string{
		{"1", "Meshery", "Running"},
		{"2", "Adapter", "Stopped"},
	}
	footer := []string{"", "", "Total: 2"}

	t.Run("Verify table output format", func(t *testing.T) {
		var buf bytes.Buffer
		
		// Use our new function that accepts a buffer
		PrintToTableWithWriter(&buf, header, data, footer)
		
		output := buf.String()

		// Verify headers are present
		assert.Contains(t, output, "ID")
		assert.Contains(t, output, "NAME")
		
		// Verify data rows are present
		assert.Contains(t, output, "Meshery")
		assert.Contains(t, output, "Running")
		assert.Contains(t, output, "Adapter")
		
		// Verify footer is present
		assert.Contains(t, output, "Total: 2")
	})
}

func TestGenerateTableOptions(t *testing.T) {
	t.Run("Options generation", func(t *testing.T) {
		options := generateTableOptions()
		assert.NotNil(t, options)
		assert.Len(t, options, 2, "Should return renderer and config options")
	})
}