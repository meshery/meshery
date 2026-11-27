package handlers

import (
	"archive/tar"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDownloadMesheryPatternHandler_OCIExportDesignNames exports designs as OCI
// images. The design name becomes the image's repository name, so a name with
// characters a repository name does not allow used to fail the export.
func TestDownloadMesheryPatternHandler_OCIExportDesignNames(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	require.NoError(t, err)

	db, err := database.New(database.Options{
		Engine:   database.SQLITE,
		Filename: filepath.Join(t.TempDir(), "oci-export.db"),
	})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(events.Event{}, models.MesheryPattern{}))

	rm, err := registry.NewRegistryManager(&db)
	require.NoError(t, err)

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config: &models.HandlerConfig{
			EventBroadcaster: models.NewBroadcaster("Events"),
		},
		log:             log,
		SystemID:        &systemID,
		registryManager: rm,
	}
	user := &models.User{
		ID:        uuid.Must(uuid.NewV4()),
		FirstName: "Test",
		LastName:  "User",
	}
	baseProvider := &models.DefaultLocalProvider{
		EventsPersister:  &models.EventsPersister{DB: &db},
		GenericPersister: &db,
	}

	tests := []struct {
		designName     string
		wantRepository string
	}{
		// Names that already exported keep the repository name they had.
		{"nginx-demo", "nginx-demo"},
		{"Nginx Demo v2.0", "nginxdemov2.0"},
		// Names that used to fail the export.
		{"My Design (1)", "mydesign-1"},
		{"nginx: prod", "nginx-prod"},
		{"nginx/prod", "nginx-prod"},
		{"a", "a-design"},
	}

	for _, tt := range tests {
		t.Run(tt.designName, func(t *testing.T) {
			patID := uuid.Must(uuid.NewV4())
			designFile := fmt.Sprintf("name: %q\nschemaVersion: designs.meshery.io/v1beta1\nversion: 0.0.1\ncomponents: []\n", tt.designName)
			prov := &roundtripProvider{
				DefaultLocalProvider: baseProvider,
				patterns: map[string]*models.MesheryPattern{
					patID.String(): {
						ID:            &patID,
						Name:          tt.designName,
						PatternFile:   designFile,
						SourceContent: []byte(designFile),
						Type:          sql.NullString{String: "design", Valid: true},
					},
				},
			}

			req := httptest.NewRequest(http.MethodGet, "/api/pattern/download/"+patID.String()+"?oci=true", nil)
			req = mux.SetURLVars(req, map[string]string{"id": patID.String()})
			req.Header.Set("token", "test-token")
			rec := httptest.NewRecorder()

			h.DownloadMesheryPatternHandler(rec, req, nil, user, prov)

			require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
			assert.Equal(t, "application/tar", rec.Header().Get("Content-Type"))
			assert.Equal(t, []string{tt.wantRepository + ":latest"}, ociTarballRepoTags(t, rec.Body.Bytes()))
		})
	}
}

func TestDesignOCIRepositoryName(t *testing.T) {
	tests := []struct {
		designName string
		want       string
	}{
		// Names that already gave a valid repository name keep it.
		{"nginx-demo", "nginx-demo"},
		{"Untitled Design", "untitleddesign"},
		{"Nginx Demo v2.0", "nginxdemov2.0"},
		{"my_design", "my_design"},
		{"Copy of nginx-demo", "copyofnginx-demo"},
		// Names with characters a repository name does not allow.
		{"My Design (1)", "mydesign-1"},
		{"Design #2", "design-2"},
		{"nginx: prod", "nginx-prod"},
		{"nginx/prod", "nginx-prod"},
		{"nginx, redis", "nginx-redis"},
		{"nginx & redis", "nginx-redis"},
		{"nginx+redis", "nginx-redis"},
		{"café", "caf"},
		{"-nginx", "nginx"},
		// Names too short or with nothing usable left.
		{"a", "a-design"},
		{"", "design"},
		{"(!)", "design"},
		{strings.Repeat("a", 300), strings.Repeat("a", 255)},
	}

	for _, tt := range tests {
		t.Run(tt.designName, func(t *testing.T) {
			got := designOCIRepositoryName(tt.designName)
			assert.Equal(t, tt.want, got)
			assert.Regexp(t, ociRepositoryNamePattern, got)
			assert.GreaterOrEqual(t, len(got), ociRepositoryNameMinLength)
			assert.LessOrEqual(t, len(got), ociRepositoryNameMaxLength)
		})
	}
}

// ociTarballRepoTags returns the RepoTags recorded in the manifest.json of an
// image tarball.
func ociTarballRepoTags(t *testing.T, tarball []byte) []string {
	t.Helper()
	var manifest []struct {
		RepoTags []string `json:"RepoTags"`
	}
	tr := tar.NewReader(bytes.NewReader(tarball))
	for {
		header, err := tr.Next()
		require.NoError(t, err, "manifest.json not found in the image tarball")
		if header.Name == "manifest.json" {
			require.NoError(t, json.NewDecoder(tr).Decode(&manifest))
			break
		}
	}
	require.Len(t, manifest, 1)
	return manifest[0].RepoTags
}
