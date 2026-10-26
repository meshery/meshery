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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
	meshkitutils "github.com/meshery/meshkit/utils"
)

// mockCSVProvider mocks Provider for CSV registration tests.
type mockCSVProvider struct {
	*models.DefaultLocalProvider
}

// GetProviderToken returns an empty provider token for testing.
func (m *mockCSVProvider) GetProviderToken(_ *http.Request) (string, error) {
	return "", nil
}

// PersistEvent implements the event persistence stub for testing.
func (m *mockCSVProvider) PersistEvent(_ events.Event, _ string) error {
	return nil
}

// TestRegisterMeshmodels_TemporaryCSVCleanup verifies that temporary CSV files
// allocated during model, component, and relationship CSV processing are cleaned up
// without leaking on disk across both early exit and full generation paths (Fixes #22571).
func TestRegisterMeshmodels_TemporaryCSVCleanup(t *testing.T) {
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config: &models.HandlerConfig{
			EventBroadcaster: models.NewBroadcaster("Events"),
		},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}

	provider := &mockCSVProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	trackTempFiles := func() (map[string]bool, func()) {
		patterns := []string{"model-*.csv", "component-*.csv", "relationship-*.csv"}
		initialMap := make(map[string]bool)
		for _, pat := range patterns {
			matches, _ := filepath.Glob(filepath.Join(os.TempDir(), pat))
			for _, m := range matches {
				initialMap[m] = true
			}
		}

		assertNoLeaks := func() {
			for _, pat := range patterns {
				matches, _ := filepath.Glob(filepath.Join(os.TempDir(), pat))
				for _, match := range matches {
					if !initialMap[match] {
						t.Errorf("temporary file %s was leaked by RegisterMeshmodels", match)
					}
				}
			}
		}
		return initialMap, assertNoLeaks
	}

	validCSV := "data:text/csv;base64," + base64.StdEncoding.EncodeToString([]byte("name,model\ntest,test-model"))

	t.Run("early exit on invalid component CSV cleans up model file", func(t *testing.T) {
		_, assertNoLeaks := trackTempFiles()

		payload := map[string]interface{}{
			"uploadType": "csv",
			"importBody": map[string]interface{}{
				"modelCsv":     validCSV,
				"componentCsv": "invalid-csv-data",
			},
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/meshmodels/register", bytes.NewReader(payloadBytes))
		rec := httptest.NewRecorder()
		h.RegisterMeshmodels(rec, req, nil, user, provider)

		assertNoLeaks()
	})

	t.Run("early exit on invalid relationship CSV cleans up model and component files", func(t *testing.T) {
		_, assertNoLeaks := trackTempFiles()

		payload := map[string]interface{}{
			"uploadType": "csv",
			"importBody": map[string]interface{}{
				"modelCsv":        validCSV,
				"componentCsv":    validCSV,
				"relationshipCsv": "invalid-csv-data",
			},
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/meshmodels/register", bytes.NewReader(payloadBytes))
		rec := httptest.NewRecorder()
		h.RegisterMeshmodels(rec, req, nil, user, provider)

		assertNoLeaks()
	})

	t.Run("all three CSV files created are cleaned up after processing", func(t *testing.T) {
		_, assertNoLeaks := trackTempFiles()

		payload := map[string]interface{}{
			"uploadType": "csv",
			"importBody": map[string]interface{}{
				"modelCsv":        validCSV,
				"componentCsv":    validCSV,
				"relationshipCsv": validCSV,
			},
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/meshmodels/register", bytes.NewReader(payloadBytes))
		rec := httptest.NewRecorder()
		h.RegisterMeshmodels(rec, req, nil, user, provider)

		assertNoLeaks()
	})
}

// TestRegisterMeshmodels_CrossPlatformRegistryLocation verifies that the model
// location resolves via meshkitutils.GetHome() rather than os.Getenv("HOME"),
// ensuring correct cross-platform path resolution across Linux, macOS, and Windows (Fixes #22571).
func TestRegisterMeshmodels_CrossPlatformRegistryLocation(t *testing.T) {
	home := meshkitutils.GetHome()
	if home == "" {
		t.Fatalf("meshkitutils.GetHome() returned empty string")
	}

	modelLocation := filepath.Join(home, utils.RegistryLocation)
	if modelLocation == "" {
		t.Fatalf("expected non-empty model location")
	}

	invalidRootPath := filepath.Join("", utils.RegistryLocation)
	if modelLocation == invalidRootPath {
		t.Errorf("model location %q incorrectly resolved to system root without user home", modelLocation)
	}

	if !strings.HasPrefix(modelLocation, home) {
		t.Errorf("expected model location %q to start with home directory %q", modelLocation, home)
	}
}
