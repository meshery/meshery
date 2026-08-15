package display

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

type viewOutputTestStruct struct {
	Name string `json:"name" yaml:"name"`
}

func TestFormatAndSaveOutput_WhenDisplayOnly_ThenWritesToConfiguredWriter(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	buf := &bytes.Buffer{}
	err := FormatAndSaveOutput(viewOutputTestStruct{Name: "meshery"}, "yaml", buf, false, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("name: meshery")) {
		t.Fatalf("expected output to contain %q, got %q", "name: meshery", buf.String())
	}
}

func TestFormatAndSaveOutput_WhenSaveEnabled_ThenWritesToFile(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	tmpFile, err := os.CreateTemp("", "meshery_view_output_*.json")
	if err != nil {
		t.Fatalf("expected temp file to be created, got %v", err)
	}
	tmpFilePath := tmpFile.Name()
	_ = tmpFile.Close()
	defer func() {
		if removeErr := os.Remove(tmpFilePath); removeErr != nil {
			t.Fatalf("expected temp file to be removed, got %v", removeErr)
		}
	}()

	buf := &bytes.Buffer{}
	err = FormatAndSaveOutput(viewOutputTestStruct{Name: "meshery"}, "json", buf, true, tmpFilePath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("\"name\": \"meshery\"")) {
		t.Fatalf("expected output to contain %q, got %q", "\"name\": \"meshery\"", buf.String())
	}

	content, err := os.ReadFile(tmpFilePath)
	if err != nil {
		t.Fatalf("expected saved file to be readable, got %v", err)
	}

	var got viewOutputTestStruct
	err = json.Unmarshal(content, &got)
	if err != nil {
		t.Fatalf("expected saved file to unmarshal, got %v", err)
	}
	if got.Name != "meshery" {
		t.Fatalf("expected saved file name %q, got %q", "meshery", got.Name)
	}
}

func TestFormatAndSaveOutput_WhenFormatIsInvalid_ThenReturnsUnsupportedFormatError(t *testing.T) {
	err := FormatAndSaveOutput(viewOutputTestStruct{Name: "meshery"}, "invalid", nil, false, "")
	utils.AssertMeshkitErrorsEqual(t, err, ErrUnsupportedFormat("invalid"))
}

func TestFormatAndSaveOutput_AcceptsMixedCaseFormat(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "mixed-case.yaml")

	err := FormatAndSaveOutput(viewOutputTestStruct{Name: "meshery"}, "YAML", nil, true, tmpFile)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
