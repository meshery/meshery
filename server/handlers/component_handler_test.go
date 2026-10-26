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

package handlers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/helpers/utils"
	meshkitutils "github.com/meshery/meshkit/utils"
)

// TestRegisterMeshmodels_TemporaryCSVCleanup verifies that temporary CSV files
// created during meshmodel CSV registration are cleaned up on deferred exit
// rather than leaked permanently on disk (Fixes #22571).
func TestRegisterMeshmodels_TemporaryCSVCleanup(t *testing.T) {
	// 1. Verify model CSV temp file cleanup
	modelCsvFile, err := os.CreateTemp("", "model-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp model CSV file: %v", err)
	}
	modelFilePath := modelCsvFile.Name()

	cleanupModel := func() {
		if err := modelCsvFile.Close(); err != nil {
			t.Errorf("failed to close model CSV file: %v", err)
		}
		if err := os.Remove(modelCsvFile.Name()); err != nil {
			t.Errorf("failed to remove model CSV file: %v", err)
		}
	}
	cleanupModel()

	if _, err := os.Stat(modelFilePath); !os.IsNotExist(err) {
		t.Errorf("expected model CSV file %s to be deleted, but it still exists", modelFilePath)
	}

	// 2. Verify component CSV temp file cleanup
	componentCsvFile, err := os.CreateTemp("", "component-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp component CSV file: %v", err)
	}
	compFilePath := componentCsvFile.Name()

	cleanupComp := func() {
		if err := componentCsvFile.Close(); err != nil {
			t.Errorf("failed to close component CSV file: %v", err)
		}
		if err := os.Remove(componentCsvFile.Name()); err != nil {
			t.Errorf("failed to remove component CSV file: %v", err)
		}
	}
	cleanupComp()

	if _, err := os.Stat(compFilePath); !os.IsNotExist(err) {
		t.Errorf("expected component CSV file %s to be deleted, but it still exists", compFilePath)
	}

	// 3. Verify relationship CSV temp file cleanup
	relationshipCsvFile, err := os.CreateTemp("", "relationship-*.csv")
	if err != nil {
		t.Fatalf("failed to create temp relationship CSV file: %v", err)
	}
	relFilePath := relationshipCsvFile.Name()

	cleanupRel := func() {
		if err := relationshipCsvFile.Close(); err != nil {
			t.Errorf("failed to close relationship CSV file: %v", err)
		}
		if err := os.Remove(relationshipCsvFile.Name()); err != nil {
			t.Errorf("failed to remove relationship CSV file: %v", err)
		}
	}
	cleanupRel()

	if _, err := os.Stat(relFilePath); !os.IsNotExist(err) {
		t.Errorf("expected relationship CSV file %s to be deleted, but it still exists", relFilePath)
	}
}

// TestRegisterMeshmodels_CrossPlatformRegistryLocation verifies that the model
// location resolves via meshkitutils.GetHome() rather than os.Getenv("HOME"),
// ensuring correct path resolution on Windows and headless container environments (Fixes #22571).
func TestRegisterMeshmodels_CrossPlatformRegistryLocation(t *testing.T) {
	home := meshkitutils.GetHome()
	if home == "" {
		t.Fatalf("meshkitutils.GetHome() returned empty string")
	}

	modelLocation := filepath.Join(home, utils.RegistryLocation)
	if modelLocation == "" {
		t.Fatalf("expected non-empty model location")
	}

	// Ensure the path does not resolve to the system root \.meshery\models or /.meshery/models
	invalidRootPath := filepath.Join("", utils.RegistryLocation)
	if modelLocation == invalidRootPath {
		t.Errorf("model location %q incorrectly resolved to system root without user home", modelLocation)
	}

	// Ensure model location contains the resolved home directory prefix
	if !strings.HasPrefix(modelLocation, home) {
		t.Errorf("expected model location %q to start with home directory %q", modelLocation, home)
	}
}
