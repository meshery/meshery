package handlers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	"github.com/meshery/meshkit/models/registration"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRegistryDB(t *testing.T) (*meshmodel.RegistryManager, *database.Handler) {
	t.Helper()
	db, err := database.New(database.Options{
		Filename: filepath.Join(t.TempDir(), "registry.db"),
		Engine:   "sqlite",
	})
	require.NoError(t, err, "failed to create test database")
	t.Cleanup(func() {
		_ = db.DBClose()
	})
	rm, err := meshmodel.NewRegistryManager(&db)
	require.NoError(t, err, "failed to create registry manager")
	return rm, &db
}

func seedTestExportModel(t *testing.T, rm *meshmodel.RegistryManager, modelName string) {
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
	_, _, err = rm.RegisterEntity(meshmodel.RegistrantHostToV1beta3(conn), m)
	require.NoError(t, err, "seedTestExportModel: RegisterEntity failed")
}

func validateExportedPayload(t *testing.T, filename string, r io.Reader, foundModel, foundComponent *bool) {
	t.Helper()
	switch filepath.Base(filename) {
	case "model.json":
		data, err := io.ReadAll(r)
		require.NoError(t, err)
		var m map[string]interface{}
		require.NoError(t, json.Unmarshal(data, &m), "model.json must be valid JSON")
		assert.Equal(t, "meshery-test-artifacthub-registrant", m["name"])
		*foundModel = true
	case "TestComponent.json":
		data, err := io.ReadAll(r)
		require.NoError(t, err)
		var c map[string]interface{}
		require.NoError(t, json.Unmarshal(data, &c), "TestComponent.json must be valid JSON")
		assert.Equal(t, "Test Component", c["displayName"])
		*foundComponent = true
	}
}

func seedArtifactHubModel(t *testing.T, rm *meshmodel.RegistryManager) {
	t.Helper()
	regErrorLog := models.NewRegistrationFailureLogHandler()
	regHelper := registration.NewRegistrationHelper(t.TempDir(), rm, regErrorLog)
	regHelper.Register(registration.NewDir(artifactHubRegistrantModelDir))

	for registrant, summary := range regErrorLog.GetNonRegisteredEntitiesCount() {
		require.Zero(t, summary.Models, "setup failure: failed to register model for registrant %s", registrant)
		require.Zero(t, summary.Components, "setup failure: failed to register components for registrant %s", registrant)
	}

	entities, count, _, err := rm.GetEntities(&regv1beta1.ModelFilter{Name: "meshery-test-artifacthub-registrant"})
	require.NoError(t, err, "setup failure: failed to query registry manager for test fixture")
	require.EqualValues(t, 1, count, "setup failure: test fixture model was not registered")
	require.NotEmpty(t, entities, "setup failure: registered model entity list is empty")
}

func TestExportModel_NotFound(t *testing.T) {
	rm, _ := newTestRegistryManager(t)
	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:          &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:             newTestLogger(t),
		SystemID:        &systemID,
		registryManager: rm,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/meshmodels/models/export?name=nonexistent", nil)
	rec := httptest.NewRecorder()

	h.ExportModel(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.True(t, strings.Contains(rec.Body.String(), "has not been found"))
}

func TestExportModel_TarGz_Success(t *testing.T) {
	rm, _ := newTestRegistryManager(t)
	seedArtifactHubModel(t, rm)

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:          &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:             newTestLogger(t),
		SystemID:        &systemID,
		registryManager: rm,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/meshmodels/models/export?name=meshery-test-artifacthub-registrant&file_type=tar.gz", nil)
	rec := httptest.NewRecorder()

	h.ExportModel(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/gzip", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "meshery-test-artifacthub-registrant.tar.gz")

	gzr, err := gzip.NewReader(bytes.NewReader(rec.Body.Bytes()))
	require.NoError(t, err)
	defer func() { _ = gzr.Close() }()

	tr := tar.NewReader(gzr)
	var foundModelJSON, foundFixtureComponent bool
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		validateExportedPayload(t, hdr.Name, tr, &foundModelJSON, &foundFixtureComponent)
	}
	assert.True(t, foundModelJSON, "model.json should be present in exported archive")
	assert.True(t, foundFixtureComponent, "fixture component should be present in exported archive")
}

func TestExportModel_OCI_Success(t *testing.T) {
	rm, _ := newTestRegistryManager(t)
	seedArtifactHubModel(t, rm)

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:          &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:             newTestLogger(t),
		SystemID:        &systemID,
		registryManager: rm,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/meshmodels/models/export?name=meshery-test-artifacthub-registrant&file_type=oci", nil)
	rec := httptest.NewRecorder()

	h.ExportModel(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/x-tar", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "meshery-test-artifacthub-registrant.tar")

	tr := tar.NewReader(bytes.NewReader(rec.Body.Bytes()))
	var foundManifestJSON, foundLayerTarGz bool
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		if filepath.Base(hdr.Name) == "manifest.json" {
			foundManifestJSON = true
		}
		if strings.HasSuffix(hdr.Name, ".tar.gz") {
			foundLayerTarGz = true
			gzr, err := gzip.NewReader(tr)
			require.NoError(t, err)
			defer func() { _ = gzr.Close() }()
			layerTr := tar.NewReader(gzr)
			var foundModelInLayer, foundCompInLayer bool
			for {
				layerHdr, err := layerTr.Next()
				if err == io.EOF {
					break
				}
				require.NoError(t, err)
				validateExportedPayload(t, layerHdr.Name, layerTr, &foundModelInLayer, &foundCompInLayer)
			}
			assert.True(t, foundModelInLayer, "model.json should be present inside OCI image layer")
			assert.True(t, foundCompInLayer, "fixture component should be present inside OCI image layer")
		}
	}
	assert.True(t, foundManifestJSON, "manifest.json should be present in OCI archive")
	assert.True(t, foundLayerTarGz, "layer tar.gz should be present in OCI archive")
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
