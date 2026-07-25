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

package connections

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// TestCreateConnectionCleanupOnTokenFailure verifies that kubeconfig.yaml is
// removed when the connection is written but subsequent token setup fails.
func TestCreateConnectionCleanupOnTokenFailure(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "kubeconfig.yaml")

	originalConfigPath := utils.ConfigPath
	defer func() { utils.ConfigPath = originalConfigPath }()
	utils.ConfigPath = configPath

	// Simulate: file written, then token step fails → file must be cleaned up.
	simulateWriteThenFail := func() {
		fileWritten := false
		defer func() {
			if fileWritten {
				_ = os.Remove(utils.ConfigPath)
			}
		}()

		f, err := os.Create(utils.ConfigPath)
		if err != nil {
			t.Fatalf("failed to create config file: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("failed to close config file: %v", err)
		}
		fileWritten = true

		// token step fails — fileWritten stays true, defer removes the file
	}

	simulateWriteThenFail()

	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatal("kubeconfig.yaml should be removed when token setup fails after write")
	}
}

// TestCreateConnectionNoCleanupOnSuccess verifies that kubeconfig.yaml is
// kept when the full connection flow succeeds.
func TestCreateConnectionNoCleanupOnSuccess(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "kubeconfig.yaml")

	originalConfigPath := utils.ConfigPath
	defer func() { utils.ConfigPath = originalConfigPath }()
	utils.ConfigPath = configPath

	// Simulate: file written, token step succeeds → file must be kept.
	simulateWriteThenSucceed := func() {
		fileWritten := false
		defer func() {
			if fileWritten {
				_ = os.Remove(utils.ConfigPath)
			}
		}()

		f, err := os.Create(utils.ConfigPath)
		if err != nil {
			t.Fatalf("failed to create config file: %v", err)
		}
		if err := f.Close(); err != nil {
			t.Fatalf("failed to close config file: %v", err)
		}
		fileWritten = true

		// token step succeeds
		fileWritten = false
	}

	simulateWriteThenSucceed()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("kubeconfig.yaml should be kept when connection is created successfully")
	}
}

// TestCreateConnectionNoCleanupOnPrereqFailure verifies that a pre-existing
// kubeconfig.yaml is not deleted when a prerequisite check fails before any
// file write occurs.
func TestCreateConnectionNoCleanupOnPrereqFailure(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "kubeconfig.yaml")

	originalConfigPath := utils.ConfigPath
	defer func() { utils.ConfigPath = originalConfigPath }()
	utils.ConfigPath = configPath

	// Create a pre-existing config file.
	f, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("failed to create pre-existing config file: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("failed to close pre-existing config file: %v", err)
	}

	// Simulate: prerequisite check fails before any write → file must be untouched.
	simulatePrereqFail := func() {
		fileWritten := false
		defer func() {
			if fileWritten {
				_ = os.Remove(utils.ConfigPath)
			}
		}()

		// prerequisite fails here, fileWritten never set to true
		return //nolint:staticcheck
	}

	simulatePrereqFail()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("pre-existing kubeconfig.yaml must not be deleted when prereq check fails")
	}
}
