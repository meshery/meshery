package registry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkiterrors "github.com/meshery/meshkit/errors"
)

func TestInvokeCompUpdate_PropagatesComponentCSVReadError(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	originalComponentCSVPath := componentCSVFilePath
	originalLogFile := logFile
	originalLogDirPath := logDirPath
	originalModelLocation := modelLocation
	originalModelName := modelName
	originalSheetGID := sheetGID
	t.Cleanup(func() {
		componentCSVFilePath = originalComponentCSVPath
		logFile = originalLogFile
		logDirPath = originalLogDirPath
		modelLocation = originalModelLocation
		modelName = originalModelName
		sheetGID = originalSheetGID
	})

	testLogFile, err := os.CreateTemp(t.TempDir(), "registry-update-*")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	logFile = testLogFile

	componentCSVFilePath = filepath.Join(t.TempDir(), "missing-components.csv")
	logDirPath = t.TempDir()
	modelLocation = t.TempDir()
	modelName = ""
	sheetGID = 1

	err = InvokeCompUpdate()
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if meshkiterrors.GetCode(err) != ErrUpdateRegistryCode {
		t.Errorf("expected error code %v, got %v", ErrUpdateRegistryCode, meshkiterrors.GetCode(err))
	}
}
