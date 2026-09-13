package core

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

type mockCloser struct {
	closeFunc func() error
}

func (m *mockCloser) Close() error {
	return m.closeFunc()
}

func TestSafeClose(t *testing.T) {
	testCases := []struct {
		name      string
		closeFunc func() error
		wantPanic bool
	}{
		{
			name: "given closer without error when SafeClose then do not panic",
			closeFunc: func() error {
				return nil
			},
			wantPanic: false,
		},
		{
			name: "given closer with error when SafeClose then do not panic",
			closeFunc: func() error {
				return errors.New("mock close error")
			},
			wantPanic: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mc := &mockCloser{
				closeFunc: tc.closeFunc,
			}

			defer func() {
				r := recover()
				if tc.wantPanic {
					if r == nil {
						t.Errorf("expected a panic")
					}
				} else {
					if r != nil {
						t.Errorf("did not expect a panic: %v", r)
					}
				}
			}()

			SafeClose(mc)
		})
	}
}

func TestCreateManifestsFolder(t *testing.T) {
	originalMesheryFolder := MesheryFolder
	originalManifestsFolder := ManifestsFolder

	defer func() {
		MesheryFolder = originalMesheryFolder
		ManifestsFolder = originalManifestsFolder
	}()

	tempDir := t.TempDir()
	MesheryFolder = tempDir
	ManifestsFolder = "test_manifests"

	manifestPath := filepath.Join(MesheryFolder, ManifestsFolder)

	t.Run("given no existing folder when CreateManifestsFolder then create successfully", func(t *testing.T) {
		_ = os.RemoveAll(manifestPath)

		err := CreateManifestsFolder()
		if err != nil {
			t.Fatalf("CreateManifestsFolder() error = %v", err)
		}

		info, err := os.Stat(manifestPath)
		if err != nil {
			t.Fatalf("os.Stat(%s) error = %v", manifestPath, err)
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", manifestPath)
		}
	})

	t.Run("given existing folder when CreateManifestsFolder then recreate successfully", func(t *testing.T) {
		err := CreateManifestsFolder()
		if err != nil {
			t.Fatalf("CreateManifestsFolder() error = %v", err)
		}

		dummyFile := filepath.Join(manifestPath, "dummy.txt")
		err = os.WriteFile(dummyFile, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("os.WriteFile() error = %v", err)
		}

		err = CreateManifestsFolder()
		if err != nil {
			t.Fatalf("CreateManifestsFolder() error = %v", err)
		}

		_, err = os.Stat(dummyFile)
		if !os.IsNotExist(err) {
			t.Errorf("expected dummy file to be deleted, but os.Stat error = %v", err)
		}
	})
}
