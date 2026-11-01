package handlers

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			defer gzReader.Close()

			tr := tar.NewReader(gzReader)
			fileCount := 0
			for {
				hdr, err := tr.Next()
				if err == io.EOF {
					break
				}
				if err != nil {
					errs <- fmt.Errorf("worker %d tar read error: %w", workerID, err)
					return
				}
				if hdr.Name != "" {
					fileCount++
				}
			}

			if fileCount == 0 {
				errs <- fmt.Errorf("worker %d tar archive was empty", workerID)
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
