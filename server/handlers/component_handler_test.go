package handlers

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	_models "github.com/meshery/meshkit/models/meshmodel/core/v1beta1"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	_model "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta3/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	m := &_model.ModelDefinition{
		Name:          modelName,
		DisplayName:   "Test Model",
		SchemaVersion: "models.meshery.io/v1beta1",
		Version:       "v1.0.0",
		Model:         _model.Model{Version: "v1.0.0"},
		Category:      category.CategoryDefinition{Name: "Orchestration"},
		Status:        _model.Enabled,
		Metadata:      &_model.ModelDefinition_Metadata{},
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
	defer func() {
		_ = resp.Body.Close()
	}()
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
	defer func() {
		_ = resp.Body.Close()
	}()
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
			defer func() {
				_ = resp.Body.Close()
			}()
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

func newComponentModelTestDB(t *testing.T) *database.Handler {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		if strings.Contains(err.Error(), "requires cgo") {
			t.Skip("skipping database test: sqlite requires cgo")
		}
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	dbHandler := &database.Handler{DB: db, Mutex: &sync.Mutex{}}
	if err := dbHandler.AutoMigrate(
		&registry.Registry{},
		&component.ComponentDefinition{},
		&relationship.RelationshipDefinition{},
		&_models.PolicyDefinition{},
		&_model.ModelDefinition{},
	); err != nil {
		t.Fatalf("failed to auto migrate model definition and dependencies: %v", err)
	}
	return dbHandler
}

func TestDeleteModel_InvalidUUID(t *testing.T) {
	h := &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/meshmodels/models/invalid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "invalid-uuid",
	})
	rw := httptest.NewRecorder()

	h.DeleteModel(rw, req, nil, nil, nil)
	resp := rw.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteModel_NotFound(t *testing.T) {
	dbHandler := newComponentModelTestDB(t)
	h := &Handler{
		config:    &models.HandlerConfig{},
		log:       newTestLogger(t),
		dbHandler: dbHandler,
	}

	missingID := uuid.Must(uuid.NewV4()).String()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/meshmodels/models/%s", missingID), nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": missingID,
	})
	rw := httptest.NewRecorder()

	h.DeleteModel(rw, req, nil, nil, nil)
	resp := rw.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "deleting non-existent model must return HTTP 404")
	assert.Contains(t, rw.Body.String(), fmt.Sprintf("model with id %s not found", missingID))
}

func TestDeleteModel_Success(t *testing.T) {
	dbHandler := newComponentModelTestDB(t)
	h := &Handler{
		config:    &models.HandlerConfig{},
		log:       newTestLogger(t),
		dbHandler: dbHandler,
	}

	modelID := uuid.Must(uuid.NewV4())
	modelRecord := _model.ModelDefinition{
		ID:   modelID,
		Name: "test-model",
	}
	require.NoError(t, dbHandler.Create(&modelRecord).Error)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/meshmodels/models/%s", modelID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": modelID.String(),
	})
	rw := httptest.NewRecorder()

	h.DeleteModel(rw, req, nil, nil, nil)
	resp := rw.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "deleting existing model must return HTTP 204")

	var count int64
	dbHandler.Model(&_model.ModelDefinition{}).Where("id = ?", modelID).Count(&count)
	assert.Equal(t, int64(0), count, "model record must be deleted from database")
}
