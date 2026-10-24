package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/models/registration"
	"github.com/meshery/meshkit/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	require.NoError(t, err)

	regm, err := meshmodel.NewRegistryManager(&db)
	require.NoError(t, err)

	regHelper := registration.NewRegistrationHelper(t.TempDir(), regm, models.NewRegistrationFailureLogHandler())
	regHelper.Register(registration.NewDir(artifactHubRegistrantModelDir))

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:          &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:             newTestLogger(t),
		SystemID:        &systemID,
		registryManager: regm,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/meshmodels/models/export?name=meshery-test-artifacthub-registrant&file_type=tar.gz", nil)
	rec := httptest.NewRecorder()

	h.ExportModel(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/gzip", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "meshery-test-artifacthub-registrant.tar.gz")
	assert.NotEmpty(t, rec.Body.Bytes())
}

func TestExportModel_OCI_Success(t *testing.T) {
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	require.NoError(t, err)

	regm, err := meshmodel.NewRegistryManager(&db)
	require.NoError(t, err)

	regHelper := registration.NewRegistrationHelper(t.TempDir(), regm, models.NewRegistrationFailureLogHandler())
	regHelper.Register(registration.NewDir(artifactHubRegistrantModelDir))

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:          &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:             newTestLogger(t),
		SystemID:        &systemID,
		registryManager: regm,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/meshmodels/models/export?name=meshery-test-artifacthub-registrant&file_type=oci", nil)
	rec := httptest.NewRecorder()

	h.ExportModel(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/x-tar", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Header().Get("Content-Disposition"), "meshery-test-artifacthub-registrant.tar")
	assert.NotEmpty(t, rec.Body.Bytes())
}

