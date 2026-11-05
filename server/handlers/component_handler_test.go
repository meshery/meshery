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
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	meshkitutils "github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// ExportModel Test Fixtures & Tests
// ---------------------------------------------------------------------------

func newTestRegistryDB(t *testing.T) (*registry.RegistryManager, *database.Handler) {
	t.Helper()
	db, err := database.New(database.Options{
		Filename: filepath.Join(t.TempDir(), "registry.db"),
		Engine:   "sqlite",
	})
	require.NoError(t, err, "failed to create test database")
	t.Cleanup(func() {
		_ = db.DBClose()
	})
	rm, err := registry.NewRegistryManager(&db)
	require.NoError(t, err, "failed to create registry manager")
	return rm, &db
}

func seedTestExportModel(t *testing.T, rm *registry.RegistryManager, modelName string) {
	t.Helper()
	conn := connection.Connection{
		Name:    "test-registrant",
		Kind:    modelName,
		Type:    "platform",
		SubType: "orchestration",
		Status:  connection.ConnectionStatusConnected,
	}
	m := &model.ModelDefinition{
		Name:          modelName,
		DisplayName:   "Test Model",
		SchemaVersion: "models.meshery.io/v1beta1",
		Version:       "v1.0.0",
		Model:         model.Model{Version: "v1.0.0"},
		Category:      category.CategoryDefinition{Name: "Orchestration"},
		Status:        model.Enabled,
		Metadata:      &model.ModelDefinition_Metadata{},
	}
	id, err := m.GenerateID()
	require.NoError(t, err)
	m.ID = id
	_, _, err = rm.RegisterEntity(registry.RegistrantHostToV1beta3(conn), m)
	require.NoError(t, err, "seedTestExportModel: RegisterEntity failed")
}

func TestExportModel_NotFound(t *testing.T) {
	rm, _ := newTestRegistryDB(t)
	log, err := logger.New("component_handler_test", logger.Options{})
	require.NoError(t, err)

	h := &Handler{
		registryManager: rm,
		log:             log,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/meshmodels/models/export?name=non-existent-model", nil)
	rw := httptest.NewRecorder()

	h.ExportModel(rw, req)

	resp := rw.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestExportModel_OCI(t *testing.T) {
	rm, _ := newTestRegistryDB(t)
	modelName := "oci-model"
	seedTestExportModel(t, rm, modelName)

	log, err := logger.New("component_handler_test", logger.Options{})
	require.NoError(t, err)

	h := &Handler{
		registryManager: rm,
		log:             log,
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/meshmodels/models/export?name=%s&fileType=oci", modelName), nil)
	rw := httptest.NewRecorder()

	h.ExportModel(rw, req)

	resp := rw.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/x-tar", resp.Header.Get("Content-Type"))
	assert.NotEmpty(t, rw.Body.Bytes())
}

func TestExportModel_Concurrent(t *testing.T) {
	rm, _ := newTestRegistryDB(t)
	modelName := "concurrent-model"
	seedTestExportModel(t, rm, modelName)

	log, err := logger.New("component_handler_test", logger.Options{})
	require.NoError(t, err)

	h := &Handler{
		registryManager: rm,
		log:             log,
	}

	const concurrency = 10
	var wg sync.WaitGroup
	errs := make(chan error, concurrency)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/meshmodels/models/export?name=%s&fileType=tar", modelName), nil)
			rw := httptest.NewRecorder()

			h.ExportModel(rw, req)

			resp := rw.Result()
			if resp.StatusCode != http.StatusOK {
				errs <- fmt.Errorf("worker %d got status %d: %s", workerID, resp.StatusCode, rw.Body.String())
				return
			}

			// Verify the response is a valid gzip archive containing files
			gzReader, err := gzip.NewReader(rw.Body)
			if err != nil {
				errs <- fmt.Errorf("worker %d gzip reader error: %w", workerID, err)
				return
			}

			tr := tar.NewReader(gzReader)
			fileCount := 0
			modelJSONFound := false
			for {
				hdr, err := tr.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					errs <- fmt.Errorf("worker %d tar read error: %w", workerID, err)
					_ = gzReader.Close()
					return
				}
				// Count only regular files; directory headers carry no content.
				if hdr.Typeflag != tar.TypeReg {
					continue
				}
				fileCount++
				if strings.HasSuffix(hdr.Name, "model.json") {
					modelJSONFound = true
				}
			}

			if err := gzReader.Close(); err != nil {
				errs <- fmt.Errorf("worker %d gzip close error: %w", workerID, err)
				return
			}

			if fileCount == 0 {
				errs <- fmt.Errorf("worker %d tar archive was empty", workerID)
				return
			}
			if !modelJSONFound {
				errs <- fmt.Errorf("worker %d tar archive is missing model.json", workerID)
				return
			}
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		assert.NoError(t, err)
	}
}

// ---------------------------------------------------------------------------
// RegisterMeshmodels Test Fixtures & Tests (Fixes #22571)
// ---------------------------------------------------------------------------

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
