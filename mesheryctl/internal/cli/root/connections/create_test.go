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

func TestCreateConnectionDeferCleanupPattern(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "kubeconfig.yaml")

	originalConfigPath := utils.ConfigPath
	originalLog := utils.Log
	defer func() {
		utils.ConfigPath = originalConfigPath
		utils.Log = originalLog
	}()

	utils.ConfigPath = configPath

	simulateDefer := func(success bool) error {
		succeed := success
		defer func() {
			if !succeed {
				_ = os.Remove(utils.ConfigPath)
			}
		}()

		testFile, err := os.Create(utils.ConfigPath)
		if err != nil {
			return err
		}
		testFile.Close()

		if succeed {
			succeed = true
		}
		return nil
	}

	if err := simulateDefer(false); err != nil {
		t.Fatalf("simulateDefer(false) should not error: %v", err)
	}

	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatal("config file should be cleaned up when success=false")
	}

	if err := simulateDefer(true); err != nil {
		t.Fatalf("simulateDefer(true) should not error: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file should exist when success=true")
	}

	os.Remove(configPath)
}
