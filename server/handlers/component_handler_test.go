package handlers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	compv1beta3 "github.com/meshery/schemas/models/v1beta3/component"
	connv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
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

// seedTestExportModelVersion seeds a model with explicit version strings,
// so a traversal version can be registered for export testing.
func seedTestExportModelVersion(t *testing.T, rm *registry.RegistryManager, modelName, modelVersion, innerVersion string) {
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
		Version:       modelVersion,
		Model:         model.Model{Version: innerVersion},
		Category:      category.CategoryDefinition{Name: "Orchestration"},
		Status:        model.Enabled,
		Metadata:      &model.ModelDefinition_Metadata{},
	}
	id, err := m.GenerateID()
	require.NoError(t, err)
	m.ID = id
	_, _, err = rm.RegisterEntity(registry.RegistrantHostToV1beta3(conn), m)
	require.NoError(t, err, "seedTestExportModelVersion: RegisterEntity failed")
}

// TestExportModel_VersionTraversal verifies that a model version containing
// path separators cannot make the export write outside the per-request temp
// directory (CWE-22). The export joins model.Model.Version and model.Version
// into the working directory, so either field carrying ".." escapes it.
func TestExportModel_VersionTraversal(t *testing.T) {
	suffix := make([]byte, 8)
	_, err := rand.Read(suffix)
	require.NoError(t, err)
	probeName := "meshery-traversal-probe-" + hex.EncodeToString(suffix)
	probePath := filepath.Join(os.TempDir(), probeName)
	t.Cleanup(func() {
		_ = os.RemoveAll(probePath)
	})

	tests := []struct {
		name         string
		modelVersion string
		innerVersion string
	}{
		{"model version traverses", "../../../" + probeName, "v1"},
		{"inner model version traverses", "v1", "../../../" + probeName},
		{"both versions traverse", "../../../" + probeName, "../../../" + probeName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, _ := newTestRegistryDB(t)
			modelName := "traversal-model"
			seedTestExportModelVersion(t, rm, modelName, tt.modelVersion, tt.innerVersion)

			log, err := logger.New("component_handler_test", logger.Options{})
			require.NoError(t, err)

			h := &Handler{
				registryManager: rm,
				log:             log,
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/meshmodels/models/export?name=%s&fileType=tar", modelName), nil)
			rw := httptest.NewRecorder()

			h.ExportModel(rw, req)

			resp := rw.Result()
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			_, statErr := os.Stat(probePath)
			assert.True(t, os.IsNotExist(statErr), "export wrote outside the temp dir: %s exists", probePath)
		})
	}
}

// TestRegisterMeshmodelComponents_VersionTraversal verifies that registration
// rejects a component whose model version carries a traversal string, in
// either the outer version or the embedded model version, so a hostile
// version cannot be planted through the registration endpoint.
func TestRegisterMeshmodelComponents_VersionTraversal(t *testing.T) {
	tests := []struct {
		name         string
		version      string
		innerVersion string
	}{
		{"outer version traverses", "../../../escaped", "v1"},
		{"embedded version traverses", "v1", "../../../escaped"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, _ := newTestRegistryDB(t)

			log, err := logger.New("component_handler_test", logger.Options{})
			require.NoError(t, err)

			h := &Handler{
				registryManager: rm,
				log:             log,
			}

			comp := compv1beta3.ComponentDefinition{
				Model: &model.ModelDefinition{
					Name:    "traversal-model",
					Version: tt.version,
					Model:   model.Model{Version: tt.innerVersion},
				},
			}
			entityBytes, err := json.Marshal(comp)
			require.NoError(t, err)

			payload := registry.MeshModelRegistrantData{
				Connection: connv1beta3.Connection{},
				EntityType: entity.ComponentDefinition,
				Entity:     entityBytes,
			}
			body, err := json.Marshal(payload)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/registry/components", bytes.NewReader(body))
			rw := httptest.NewRecorder()

			h.RegisterMeshmodelComponents(rw, req)

			resp := rw.Result()
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	}
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

func TestOCIRepositoryName(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"istio-base", "istio-base"},
		{"aws-api-gateway-operator", "aws-api-gateway-operator"},
		{"cert_manager.v1", "cert_manager.v1"},
		{"confidential containers", "confidential-containers"},
		{"container network interface (cni)", "container-network-interface-cni"},
		{"open policy agent (opa)", "open-policy-agent-opa"},
		{"Piraeus Datastore", "piraeus-datastore"},
		{"open cluster management", "open-cluster-management"},
		{"(!)", "model"},
		{" (istio) ", "istio"},
		{"Istio.Base", "istio-base"},
		{"A", "a-model"},
		{strings.Repeat("a", 300), strings.Repeat("a", 255)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ociRepositoryName(tt.name)
			assert.Equal(t, tt.expected, got)
			assert.Regexp(t, ociRepositoryNamePattern, got)
			assert.GreaterOrEqual(t, len(got), ociRepositoryNameMinLength)
			assert.LessOrEqual(t, len(got), ociRepositoryNameMaxLength)
		})
	}
}

func TestExportModel_OCINameNotValidRepository(t *testing.T) {
	for _, modelName := range []string{"container network interface (cni)", "A"} {
		t.Run(modelName, func(t *testing.T) {
			rm, _ := newTestRegistryDB(t)
			seedTestExportModel(t, rm, modelName)

			log, err := logger.New("component_handler_test", logger.Options{})
			require.NoError(t, err)

			h := &Handler{
				registryManager: rm,
				log:             log,
			}

			req := httptest.NewRequest(http.MethodGet, "/api/meshmodels/models/export?fileType=oci&name="+url.QueryEscape(modelName), nil)
			rw := httptest.NewRecorder()

			h.ExportModel(rw, req)

			resp := rw.Result()
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, "application/x-tar", resp.Header.Get("Content-Type"))
			assert.NotEmpty(t, rw.Body.Bytes())

			// The tar must carry the sanitized name, not just any valid reference.
			var manifest []struct {
				RepoTags []string `json:"RepoTags"`
			}
			tr := tar.NewReader(bytes.NewReader(rw.Body.Bytes()))
			for {
				header, err := tr.Next()
				require.NoError(t, err)
				if header.Name == "manifest.json" {
					require.NoError(t, json.NewDecoder(tr).Decode(&manifest))
					break
				}
			}
			require.Len(t, manifest, 1)
			assert.Equal(t, []string{ociRepositoryName(modelName) + ":latest"}, manifest[0].RepoTags)
		})
	}
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
