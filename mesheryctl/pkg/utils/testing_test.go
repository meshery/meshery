package utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestNewTestHelper(t *testing.T) {
	th := NewTestHelper(t)

	expectedVersion := "v0.5.10"
	if th.Version != expectedVersion {
		t.Errorf("expected version %s, got %s", expectedVersion, th.Version)
	}

	expectedBaseURL := MesheryEndpoint
	if th.BaseURL != expectedBaseURL {
		t.Errorf("expected base URL %s, got %s", expectedBaseURL, th.BaseURL)
	}
}

func TestNewGoldenFile(t *testing.T) {
	dir := t.TempDir()
	name := "testfile.golden"
	gf := NewGoldenFile(t, name, dir)

	if gf.name != name {
		t.Errorf("expected name %s, got %s", name, gf.name)
	}
	if gf.dir != dir {
		t.Errorf("expected dir %s, got %s", dir, gf.dir)
	}
}

func TestGoldenFile_Load(t *testing.T) {
	dir := t.TempDir()
	name := "testfile.golden"
	content := "test content"
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("could not write test file: %v", err)
	}

	gf := NewGoldenFile(t, name, dir)
	loadedContent := gf.Load()

	if loadedContent != content {
		t.Errorf("expected content %s, got %s", content, loadedContent)
	}
}

func TestGoldenFile_Write(t *testing.T) {
	dir := t.TempDir()
	name := "testfile.golden"
	content := "test content"

	gf := NewGoldenFile(t, name, dir)
	gf.Write(content)

	path := filepath.Join(dir, name)
	writtenContent, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read written file: %v", err)
	}

	if string(writtenContent) != content {
		t.Errorf("expected content %s, got %s", content, string(writtenContent))
	}
}

func TestGoldenFile_WriteInByte(t *testing.T) {
	dir := t.TempDir()
	name := "testfile.golden"
	content := []byte("test content")

	gf := NewGoldenFile(t, name, dir)
	gf.WriteInByte(content)

	path := filepath.Join(dir, name)
	writtenContent, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read written file: %v", err)
	}

	if string(writtenContent) != string(content) {
		t.Errorf("expected content %s, got %s", string(content), string(writtenContent))
	}
}

func TestEquals(t *testing.T) {
	t.Run("EqualValues", func(t *testing.T) {
		expected := "test"
		actual := "test"
		Equals(t, expected, actual)
	})
}

func TestGetBasePath(t *testing.T) {
	// Get the expected base path using runtime.Caller
	_, expectedFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	expectedBasePath := filepath.Dir(expectedFile)

	actualBasePath := GetBasePath(t)

	if expectedBasePath != actualBasePath {
		t.Errorf("expected base path %s, got %s", expectedBasePath, actualBasePath)
	}
}

func TestGoldenFile_LoadByte(t *testing.T) {
	dir := t.TempDir()
	name := "testfile.golden"
	content := []byte("test content")

	path := filepath.Join(dir, name)
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		t.Fatalf("could not write test file: %v", err)
	}

	gf := NewGoldenFile(t, name, dir)
	loadedContent := gf.LoadByte()

	if !bytes.Equal(loadedContent, content) {
		t.Errorf("expected content %s, got %s", string(content), string(loadedContent))
	}
}

func TestSetupLogrusGrabTesting(t *testing.T) {
	t.Run("VerboseFalse", func(t *testing.T) {
		buffer := SetupLogrusGrabTesting(t, false)
		if buffer == nil {
			t.Fatal("expected buffer to be non-nil")
		}

		logrus.Info("test message")

		if !strings.Contains(buffer.String(), "test message") {
			t.Errorf("expected buffer to contain 'test message', got %s", buffer.String())
		}
	})

	t.Run("VerboseTrue", func(t *testing.T) {
		buffer := SetupLogrusGrabTesting(t, true)
		if buffer == nil {
			t.Fatal("expected buffer to be non-nil")
		}

		logrus.Info("test message")

		if !strings.Contains(buffer.String(), "test message") {
			t.Errorf("expected buffer to contain 'test message', got %s", buffer.String())
		}
	})
}

func TestSetupMeshkitLoggerTesting(t *testing.T) {
	t.Run("VerboseTrue", func(t *testing.T) {
		buffer := SetupMeshkitLoggerTesting(t, true)
		if buffer == nil {
			t.Fatal("expected buffer to be non-nil")
		}

		Log.SetLevel(logrus.DebugLevel)
		Log.Info("test message")

		if !strings.Contains(buffer.String(), "test message") {
			t.Errorf("expected buffer to contain 'test message', got %s", buffer.String())
		}
	})

	t.Run("VerboseFalse", func(t *testing.T) {
		buffer := SetupMeshkitLoggerTesting(t, false)
		if buffer == nil {
			t.Fatal("expected buffer to be non-nil")
		}

		Log.SetLevel(logrus.InfoLevel)
		Log.Info("test message")

		if !strings.Contains(buffer.String(), "test message") {
			t.Errorf("expected buffer to contain 'test message', got %s", buffer.String())
		}
	})
}

func TestSetupCustomContextEnv(t *testing.T) {
	dir := t.TempDir()
	configFilePath := filepath.Join(dir, "testconfig.yaml")
	configContent := `
contexts:
  local:
    endpoint: http://localhost:9081
    token: Default
current-context: local
`
	err := os.WriteFile(configFilePath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("could not write test config file: %v", err)
	}

	SetupCustomContextEnv(t, configFilePath)

	expectedEndpoint := "http://localhost:9081"
	actualEndpoint := viper.GetString("contexts.local.endpoint")
	if actualEndpoint != expectedEndpoint {
		t.Errorf("expected endpoint %s, got %s", expectedEndpoint, actualEndpoint)
	}

	expectedContext := "local"
	actualContext := viper.GetString("current-context")
	if actualContext != expectedContext {
		t.Errorf("expected current context %s, got %s", expectedContext, actualContext)
	}
}

func TestSetFileLocationTesting(t *testing.T) {
	dir := t.TempDir()

	expectedMesheryFolder := filepath.Join(dir, "fixtures", MesheryFolder)
	expectedDockerComposeFile := filepath.Join(expectedMesheryFolder, DockerComposeFile)
	expectedAuthConfigFile := filepath.Join(expectedMesheryFolder, AuthConfigFile)

	SetFileLocationTesting(dir)

	if MesheryFolder != expectedMesheryFolder {
		t.Errorf("expected MesheryFolder %s, got %s", expectedMesheryFolder, MesheryFolder)
	}
	if DockerComposeFile != expectedDockerComposeFile {
		t.Errorf("expected DockerComposeFile %s, got %s", expectedDockerComposeFile, DockerComposeFile)
	}
	if AuthConfigFile != expectedAuthConfigFile {
		t.Errorf("expected AuthConfigFile %s, got %s", expectedAuthConfigFile, AuthConfigFile)
	}
}

func TestPopulate(t *testing.T) {
	t.Run("TestPopulate_Success", func(t *testing.T) {

		dir := t.TempDir()
		srcFile := filepath.Join(dir, "source.txt")
		dstFile := filepath.Join(dir, "destination.txt")
		content := "test content"

		err := os.WriteFile(srcFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("could not write source file: %v", err)
		}

		err = Populate(srcFile, dstFile)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		dstContent, err := os.ReadFile(dstFile)
		if err != nil {
			t.Fatalf("could not read destination file: %v", err)
		}

		if string(dstContent) != content {
			t.Errorf("expected content %s, got %s", content, string(dstContent))
		}
	})

	t.Run("TestPopulate_SourceFileNotExist", func(t *testing.T) {

		dir := t.TempDir()
		srcFile := filepath.Join(dir, "nonexistent.txt")
		dstFile := filepath.Join(dir, "destination.txt")

		err := Populate(srcFile, dstFile)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !os.IsNotExist(err) {
			t.Errorf("expected file not exist error, got %v", err)
		}
	})

	t.Run("TestPopulate_SourceNotRegularFile", func(t *testing.T) {

		dir := t.TempDir()
		srcDir := filepath.Join(dir, "sourcedir")
		dstFile := filepath.Join(dir, "destination.txt")

		err := os.Mkdir(srcDir, 0755)
		if err != nil {
			t.Fatalf("could not create source directory: %v", err)
		}

		err = Populate(srcDir, dstFile)

		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		expectedErr := fmt.Sprintf("%s is not a regular file", srcDir)
		if err.Error() != expectedErr {
			t.Errorf("expected error %s, got %v", expectedErr, err)
		}
	})
}

