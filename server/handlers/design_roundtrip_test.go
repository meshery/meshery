package handlers

import (
	"archive/tar"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	pattern "github.com/meshery/schemas/models/v1beta3/design"
	"gopkg.in/yaml.v2"
)

type roundtripProvider struct {
	*models.DefaultLocalProvider
	patterns map[string]*models.MesheryPattern
}

func (p *roundtripProvider) GetProviderToken(r *http.Request) (string, error) {
	return "test-token", nil
}

func (p *roundtripProvider) GetMesheryPattern(r *http.Request, id string, _ string) ([]byte, error) {
	pat, ok := p.patterns[id]
	if !ok {
		return nil, fmt.Errorf("pattern not found")
	}
	return json.Marshal(pat)
}

func (p *roundtripProvider) SaveMesheryPattern(token string, pattern *models.MesheryPattern) ([]byte, error) {
	if pattern.ID == nil {
		newID := uuid.Must(uuid.NewV4())
		pattern.ID = &newID
	}
	p.patterns[pattern.ID.String()] = pattern
	return json.Marshal(pattern)
}

func TestDesignExportImportRoundtrip(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("logger init: %v", err)
	}
	systemID := uuid.Must(uuid.NewV4())
	user := &models.User{
		ID:        uuid.Must(uuid.NewV4()),
		FirstName: "Test",
		LastName:  "User",
	}

	db, err := database.New(database.Options{
		Engine:   database.SQLITE,
		Filename: filepath.Join(t.TempDir(), "roundtrip.db"),
	})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(events.Event{}, models.MesheryPattern{}); err != nil {
		t.Fatalf("migrate tables: %v", err)
	}

	rm, err := registry.NewRegistryManager(&db)
	if err != nil {
		t.Fatalf("registry manager init: %v", err)
	}

	h := &Handler{
		config: &models.HandlerConfig{
			EventBroadcaster: models.NewBroadcaster("Events"),
		},
		log:             log,
		SystemID:        &systemID,
		registryManager: rm,
	}

	baseProvider := &models.DefaultLocalProvider{
		EventsPersister:  &models.EventsPersister{DB: &db},
		GenericPersister: &db,
	}

	testCases := []struct {
		name        string
		exportQuery string
		patternFile string
		patternName string
		sourceType  string
		expectError bool
	}{
		{
			name:        "1. Normal input",
			exportQuery: "",
			patternName: "normal-design",
			sourceType:  "design",
			patternFile: `name: normal-design
schemaVersion: designs.meshery.io/v1beta1
version: 0.0.1
components: []
`,
		},
		{
			name:        "2. Empty input",
			exportQuery: "",
			patternName: "empty-design",
			sourceType:  "design",
			patternFile: ``,
			expectError: true,
		},
		{
			name:        "3. Several items",
			exportQuery: "",
			patternName: "multi-service-design",
			sourceType:  "design",
			patternFile: `name: multi-service-design
schemaVersion: designs.meshery.io/v1beta1
version: 0.0.1
components:
  - name: service-a
    component:
      kind: Deployment
      version: apps/v1
  - name: service-b
    component:
      kind: Service
      version: v1
  - name: service-c
    component:
      kind: ConfigMap
      version: v1
`,
		},
		{
			name:        "4. Error case",
			exportQuery: "",
			patternName: "malformed-design",
			sourceType:  "design",
			patternFile: `name: [unclosed list
services: {broken:
`,
			expectError: true,
		},
		{
			name:        "5. Exception override case with pkg export",
			exportQuery: "?pkg=true",
			patternName: "pkg-design",
			sourceType:  "design",
			patternFile: `name: pkg-design
schemaVersion: designs.meshery.io/v1beta1
version: 0.0.1
components: []
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			patID := uuid.Must(uuid.NewV4())
			prov := &roundtripProvider{
				DefaultLocalProvider: baseProvider,
				patterns: map[string]*models.MesheryPattern{
					patID.String(): {
						ID:            &patID,
						Name:          tc.patternName,
						PatternFile:   tc.patternFile,
						SourceContent: []byte(tc.patternFile),
						Type: sql.NullString{
							String: tc.sourceType,
							Valid:  true,
						},
					},
				},
			}

			// --- PATH A: Export ---
			exportURL := "/api/pattern/download/" + patID.String() + tc.exportQuery
			reqA := httptest.NewRequest(http.MethodGet, exportURL, nil)
			reqA = mux.SetURLVars(reqA, map[string]string{"id": patID.String()})
			reqA.Header.Set("token", "test-token")
			recA := httptest.NewRecorder()

			h.DownloadMesheryPatternHandler(recA, reqA, nil, user, prov)

			exportStatus := recA.Code
			exportBody := recA.Body.Bytes()
			contentType := recA.Header().Get("Content-Type")

			t.Logf("[%s] PATH A: Export HTTP Status=%d, Content-Type=%s, BodySize=%d",
				tc.name, exportStatus, contentType, len(exportBody))

			if tc.expectError {
				if exportStatus == http.StatusOK {
					t.Fatalf("[%s] Expected error on malformed input, but export returned HTTP 200.", tc.name)
				}
				t.Logf("[%s] Correctly failed at export stage with status %d.", tc.name, exportStatus)
				return
			}

			if exportStatus != http.StatusOK {
				t.Fatalf("[%s] Export failed: status=%d, body=%s", tc.name, exportStatus, string(exportBody))
			}

			// If export is a tar/archive (as with ?pkg=true), verify its structural integrity
			var extractedDesign []byte
			if tc.exportQuery == "?pkg=true" {
				archiveReader := bytes.NewReader(exportBody)
				tarReader := tar.NewReader(archiveReader)
				foundDesignFile := false
				foundPkgFile := false
				for {
					header, err := tarReader.Next()
					if err == io.EOF {
						break
					}
					if err != nil {
						t.Fatalf("Corrupt tar archive encountered: %v", err)
					}
					if header.Name == tc.patternName+".yml" {
						foundDesignFile = true
						var readErr error
						extractedDesign, readErr = io.ReadAll(tarReader)
						if readErr != nil {
							t.Fatalf("Failed to read design file from tar archive: %v", readErr)
						}
					}
					if header.Name == "artifacthub-pkg.yml" {
						foundPkgFile = true
					}
				}
				if !foundDesignFile {
					t.Fatalf("Tar archive missing design file %s.yml", tc.patternName)
				}
				if !foundPkgFile {
					t.Fatalf("Tar archive missing metadata file artifacthub-pkg.yml")
				}

				// Inspect whether extracted design is double-marshalled YAML string
				var unquoted string
				if err := yaml.Unmarshal(extractedDesign, &unquoted); err == nil && strings.TrimSpace(unquoted) == strings.TrimSpace(tc.patternFile) {
					t.Errorf("Design inside archive was double-serialized as a quoted YAML scalar string instead of valid YAML object. Extracted:\n%s", string(extractedDesign))
				}

				// Check for trailing content after archive EOF, ensuring that any
				// remaining bytes beyond the tar end-of-archive marker consist strictly of zero padding
				trailingBytes, err := io.ReadAll(archiveReader)
				if err != nil {
					t.Fatalf("Error reading trailing content from tar stream: %v", err)
				}
				for i, b := range trailingBytes {
					if b != 0 {
						t.Fatalf("Archive has non-zero trailing content after archive EOF at offset %d: byte=%x", i, b)
					}
				}
				if len(exportBody)%512 != 0 {
					t.Errorf("Archive binary payload has invalid block alignment: length %d not a multiple of 512 (trailing bytes appended)", len(exportBody))
				}
			}

			// --- PATH B: Import ---
			payloadData := exportBody
			importFileName := tc.patternName + ".yml"
			if tc.exportQuery == "?pkg=true" {
				// The artifacthub package export is an archive containing both
				// <pattern-name>.yml and artifacthub-pkg.yml. For direct design import
				// via POST /api/pattern, the import endpoint accepts the design YAML
				// extracted from the package.
				payloadData = extractedDesign
			}
			var importBody pattern.MesheryPatternImportRequestBody
			if err := importBody.FromMesheryPatternImportFilePayload(pattern.MesheryPatternImportFilePayload{
				Name:     &tc.patternName,
				FileName: importFileName,
				File:     payloadData,
			}); err != nil {
				t.Fatalf("Failed to create import body: %v", err)
			}
			importJSON, _ := json.Marshal(importBody)

			reqB := httptest.NewRequest(http.MethodPost, "/api/pattern/import", bytes.NewReader(importJSON))
			reqB.Header.Set("token", "test-token")
			reqB = reqB.WithContext(context.WithValue(reqB.Context(), models.TokenCtxKey, "test-token"))
			recB := httptest.NewRecorder()

			h.DesignFileImportHandler(recB, reqB, nil, user, prov)

			importStatus := recB.Code
			importResp := recB.Body.Bytes()
			t.Logf("[%s] PATH B: Import HTTP Status=%d, Body=%s", tc.name, importStatus, string(importResp))

			if importStatus != http.StatusOK && importStatus != http.StatusCreated {
				t.Fatalf("Exported design could not be re-imported: status=%d body=%s", importStatus, string(importResp))
			}
		})
	}
}
