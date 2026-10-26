package handlers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	"github.com/meshery/meshkit/models/registration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	var foundModelJSON, foundFixtureComponent bool
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		if filepath.Base(hdr.Name) == "model.json" {
			foundModelJSON = true
		}
		if filepath.Base(hdr.Name) == "TestComponent.json" {
			foundFixtureComponent = true
		}
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
			layerTr := tar.NewReader(gzr)
			var foundModelInLayer, foundCompInLayer bool
			for {
				layerHdr, err := layerTr.Next()
				if err == io.EOF {
					break
				}
				require.NoError(t, err)
				if filepath.Base(layerHdr.Name) == "model.json" {
					foundModelInLayer = true
				}
				if filepath.Base(layerHdr.Name) == "TestComponent.json" {
					foundCompInLayer = true
				}
			}
			assert.True(t, foundModelInLayer, "model.json should be present inside OCI image layer")
			assert.True(t, foundCompInLayer, "fixture component should be present inside OCI image layer")
		}
	}
	assert.True(t, foundManifestJSON, "manifest.json should be present in OCI archive")
	assert.True(t, foundLayerTarGz, "layer tar.gz should be present in OCI archive")
}
