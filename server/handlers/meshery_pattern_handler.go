package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	guid "github.com/google/uuid"
	"github.com/gorilla/mux"
	helpers "github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/meshes"
	"github.com/meshery/meshery/server/models"
	patterncore "github.com/meshery/meshery/server/models/pattern/core"
	"github.com/meshery/meshery/server/models/pattern/resource/selector"
	patternutils "github.com/meshery/meshery/server/models/pattern/utils"
	"github.com/meshery/meshkit/encoding"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/converter"
	_errors "github.com/pkg/errors"

	"github.com/meshery/meshkit/models/catalog/v1alpha1"
	"github.com/meshery/meshkit/models/events"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/models/oci"
	meshkitPatternHelpers "github.com/meshery/meshkit/models/patterns"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/meshkit/utils/catalog"

	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1alpha2"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta2/component"
	design "github.com/meshery/schemas/models/v1beta3/design"
	"gopkg.in/yaml.v3"
)

// MesheryPatternRequestBody refers to the type of request body that
// SaveMesheryPattern would receive. Canonical wire form for the
// wrapper is `patternData` (camelCase) per the identifier-naming
// migration; the legacy snake_case `pattern_data` spelling is still
// dual-accepted via UnmarshalJSON for the deprecation window so
// unmigrated clients keep working. Canonical wins when both are
// present.
//
// Deprecated: retained for API compatibility; no live call site reads
// this type today (POST /api/pattern decodes directly into
// DesignPostPayload). Kept so external tooling that reflects on the
// server package still finds the shape.
type MesheryPatternPOSTRequestBody struct {
	Name        string             `json:"name,omitempty"`
	URL         string             `json:"url,omitempty"`
	Path        string             `json:"path,omitempty"`
	Save        bool               `json:"save,omitempty"`
	PatternData *DesignPostPayload `json:"patternData,omitempty"`
}

// UnmarshalJSON dual-accepts the canonical camelCase `patternData`
// and the legacy snake_case `pattern_data` wrapper keys for
// PatternData. Canonical wins when both are present. Other fields
// unmarshal via stdlib default rules through the embedded-alias
// pattern; PatternData is explicitly re-zeroed before the precedence
// switch so a reused receiver does not carry a stale pointer when
// the next payload omits both spellings.
func (p *MesheryPatternPOSTRequestBody) UnmarshalJSON(data []byte) error {
	type alias MesheryPatternPOSTRequestBody
	aux := &struct {
		*alias
		PatternDataCanonical *DesignPostPayload `json:"patternData,omitempty"`
		PatternDataLegacy    *DesignPostPayload `json:"pattern_data,omitempty"`
	}{alias: (*alias)(p)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	p.PatternData = nil
	switch {
	case aux.PatternDataCanonical != nil:
		p.PatternData = aux.PatternDataCanonical
	case aux.PatternDataLegacy != nil:
		p.PatternData = aux.PatternDataLegacy
	}
	return nil
}

// MesheryPatternUPDATERequestBody is the request body for the pattern
// UPDATE / save path handled by handlePatternUpdate. Canonical wire
// form for both PatternData (`patternData`) and CytoscapeJSON
// (`cytoscapeJSON`) is camelCase per the identifier-naming migration;
// legacy snake_case spellings (`pattern_data`, `cytoscape_json`) are
// dual-accepted via UnmarshalJSON for the deprecation window.
// Canonical wins when both are present.
type MesheryPatternUPDATERequestBody struct {
	Name          string                 `json:"name,omitempty"`
	URL           string                 `json:"url,omitempty"`
	Path          string                 `json:"path,omitempty"`
	Save          bool                   `json:"save,omitempty"`
	PatternData   *models.MesheryPattern `json:"patternData,omitempty"`
	CytoscapeJSON string                 `json:"cytoscapeJSON,omitempty"`
}

// UnmarshalJSON dual-accepts the canonical camelCase keys
// (`patternData`, `cytoscapeJSON`) and the legacy snake_case keys
// (`pattern_data`, `cytoscape_json`) on the UPDATE request body.
// Canonical wins when both are present for a given field. Other
// fields unmarshal via stdlib default rules through the embedded-alias
// pattern; PatternData and CytoscapeJSON are explicitly reset before
// the precedence switches so a reused receiver does not carry stale
// values when the next payload omits both spellings.
func (p *MesheryPatternUPDATERequestBody) UnmarshalJSON(data []byte) error {
	type alias MesheryPatternUPDATERequestBody
	aux := &struct {
		*alias
		PatternDataCanonical   *models.MesheryPattern `json:"patternData,omitempty"`
		PatternDataLegacy      *models.MesheryPattern `json:"pattern_data,omitempty"`
		CytoscapeJSONCanonical *string                `json:"cytoscapeJSON,omitempty"`
		CytoscapeJSONLegacy    *string                `json:"cytoscape_json,omitempty"`
	}{alias: (*alias)(p)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	p.PatternData = nil
	switch {
	case aux.PatternDataCanonical != nil:
		p.PatternData = aux.PatternDataCanonical
	case aux.PatternDataLegacy != nil:
		p.PatternData = aux.PatternDataLegacy
	}
	p.CytoscapeJSON = ""
	switch {
	case aux.CytoscapeJSONCanonical != nil:
		p.CytoscapeJSON = *aux.CytoscapeJSONCanonical
	case aux.CytoscapeJSONLegacy != nil:
		p.CytoscapeJSON = *aux.CytoscapeJSONLegacy
	}
	return nil
}

// DesignPostPayload is the request body for POST /api/pattern. Canonical
// wire form is camelCase (`designFile`) per the identifier-naming
// migration; legacy snake_case (`design_file`) and the alternate
// "pattern file" vocabulary (`patternFile`, `pattern_file`) are still
// accepted for the deprecation window so unmigrated clients (e.g.
// meshery-extensions' meshmap `saveDesign` and Kanvas' legacy body
// shape) keep working. Custom MarshalJSON emits both canonical and
// legacy spellings so any external consumer still reading either form
// continues to round-trip.
//
// Once every known consumer (meshery-cloud, meshery-extensions, Kanvas)
// has migrated off the legacy spellings, drop MarshalJSON/UnmarshalJSON
// and keep only the `designFile` struct tag.
type DesignPostPayload struct {
	ID         *core.Uuid         `json:"id,omitempty"`
	Name       string             `json:"name,omitempty"`
	DesignFile design.PatternFile `json:"designFile"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID      *string              `json:"user_id"`
	Visibility  string               `json:"visibility"`
	CatalogData v1alpha1.CatalogData `json:"catalog_data,omitempty"`
}

// MarshalJSON emits both the canonical (`designFile`) and legacy
// (`design_file`) spellings for the design payload field so external
// consumers on either vocabulary keep working while they migrate.
// The alternate "pattern file" spellings (`patternFile` / `pattern_file`)
// are accepted on the Unmarshal path but not emitted — they exist only
// so clients that speak the legacy "pattern" vocabulary continue to
// parse, not so Meshery introduces new wire forms.
func (p DesignPostPayload) MarshalJSON() ([]byte, error) {
	type alias DesignPostPayload
	return json.Marshal(struct {
		alias
		DesignFileLegacy design.PatternFile `json:"design_file"`
	}{
		alias:            alias(p),
		DesignFileLegacy: p.DesignFile,
	})
}

// UnmarshalJSON accepts any of `designFile` (canonical), `patternFile`
// (alternate camelCase vocabulary), `design_file` (legacy snake_case),
// or `pattern_file` (legacy alternate). Precedence when multiple are
// present: canonical camelCase wins over legacy; `designFile` wins
// over `patternFile`.
//
// Implementation uses the embedded-alias pattern so every non-custom
// field on DesignPostPayload (including any added later) unmarshals
// via stdlib default rules — only the design-file key-precedence is
// custom-handled here. The inner `*alias` is initialised to point at
// the receiver, so fields absent from the input naturally reset to
// their zero value per stdlib json.Unmarshal semantics; DesignFile is
// explicitly re-zeroed before the precedence switch so a reused
// receiver does not retain stale design data when the next payload
// omits all four spellings.
func (p *DesignPostPayload) UnmarshalJSON(data []byte) error {
	type alias DesignPostPayload
	aux := &struct {
		*alias
		DesignFileCanonical  *design.PatternFile `json:"designFile,omitempty"`
		PatternFileCanonical *design.PatternFile `json:"patternFile,omitempty"`
		DesignFileLegacy     *design.PatternFile `json:"design_file,omitempty"`
		PatternFileLegacy    *design.PatternFile `json:"pattern_file,omitempty"`
	}{alias: (*alias)(p)}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	// Depth-0 aux fields win the `designFile` tag over the embedded
	// alias's DesignFile by Go's json struct-tag precedence rules, so
	// the alias's DesignFile is never populated directly — we apply the
	// precedence-winning spelling below. Reset first so a reused receiver
	// zeros cleanly when all four spellings are absent.
	p.DesignFile = design.PatternFile{}
	switch {
	case aux.DesignFileCanonical != nil:
		p.DesignFile = *aux.DesignFileCanonical
	case aux.PatternFileCanonical != nil:
		p.DesignFile = *aux.PatternFileCanonical
	case aux.DesignFileLegacy != nil:
		p.DesignFile = *aux.DesignFileLegacy
	case aux.PatternFileLegacy != nil:
		p.DesignFile = *aux.PatternFileLegacy
	}
	return nil
}

// PatternFileRequestHandler will handle requests of both type GET and POST
// on the route /api/pattern
func (h *Handler) PatternFileRequestHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	if r.Method == http.MethodGet {
		h.GetMesheryPatternsHandler(rw, r, prefObj, user, provider)
		return
	}

	if r.Method == http.MethodPost {
		h.handlePatternPOST(rw, r, prefObj, user, provider)
		return
	}
}

func (h *Handler) handleProviderPatternSaveError(rw http.ResponseWriter, eventBuilder *events.EventBuilder, userID core.Uuid, body []byte, err error, provider models.Provider, token string) {

	var meshkitErr errors.Error
	var event *events.Event
	errorParsingToMeshkitError := json.Unmarshal(body, &meshkitErr)

	description := ""
	if len(meshkitErr.ShortDescription) > 0 {
		description = fmt.Sprintf("Failed To Save Design, %s", meshkitErr.ShortDescription[0])
	} else {
		description = "Failed to save design"
	}

	if errorParsingToMeshkitError == nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write(body)
		h.log.Error(&meshkitErr)
		event = eventBuilder.WithSeverity(events.Error).WithDescription(description).WithMetadata(map[string]interface{}{
			"error": meshkitErr,
		}).Build()

	} else {
		h.log.Error(ErrSavePattern(err))
		writeMeshkitError(rw, ErrSavePattern(err), http.StatusBadRequest)
		event = eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrSavePattern(err),
		}).WithDescription(ErrSavePattern(err).Error()).Build()
	}

	_ = provider.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)
}

func (h *Handler) handleProviderPatternGetError(rw http.ResponseWriter, eventBuilder *events.EventBuilder, userID core.Uuid, body []byte, err error, provider models.Provider, token string) {

	var meshkitErr errors.Error
	var event *events.Event
	errorParsingToMeshkitError := json.Unmarshal(body, &meshkitErr)

	description := ""
	if len(meshkitErr.ShortDescription) > 0 {
		description = fmt.Sprintf("Failed to fetch Design, %s", meshkitErr.ShortDescription[0])
	} else {
		description = "Failed to fetch design"
	}

	if errorParsingToMeshkitError == nil {
		rw.WriteHeader(http.StatusBadRequest)
		_, _ = rw.Write(body)
		h.log.Error(&meshkitErr)
		event = eventBuilder.WithSeverity(events.Error).WithDescription(description).WithMetadata(map[string]interface{}{
			"error": meshkitErr,
		}).Build()

	} else {
		h.log.Error(ErrGetPattern(err))
		writeMeshkitError(rw, ErrGetPattern(err), http.StatusBadRequest)
		event = eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrGetPattern(err),
		}).WithDescription(ErrGetPattern(err).Error()).Build()
	}

	_ = provider.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)
}

func (h *Handler) handlePatternPOST(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	var err error

	userID := user.ID
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction(models.Create).WithSeverity(events.Informational).WithDescription("Save design ")

	requestPayload := &DesignPostPayload{}
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		h.logErrorParsingRequestBody(rw, provider, err, userID, eventBuilder)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.logErrorGettingUserToken(rw, provider, err, userID, eventBuilder)
		return
	}

	// Set the event builder with the pattern ID if available
	if requestPayload.ID != nil {
		eventBuilder = eventBuilder.ActedUpon(*requestPayload.ID)
	}

	// Dehydrate the pattern before saving to the database to reduce size
	meshkitPatternHelpers.DehydratePattern(&requestPayload.DesignFile)

	designFileBytes, err := encoding.Marshal(requestPayload.DesignFile)

	if err != nil {
		h.logErrorParsingRequestBody(rw, provider, err, userID, eventBuilder)
		return
	}

	designFile := string(designFileBytes)

	mesheryPatternRecord := models.MesheryPattern{
		ID:          requestPayload.ID,
		PatternFile: designFile,
		UserID:      requestPayload.UserID,
		Name:        requestPayload.Name,
		Visibility:  requestPayload.Visibility,
		CatalogData: requestPayload.CatalogData,
	}

	savedDesignByt, err := provider.SaveMesheryPattern(token, &mesheryPatternRecord)

	if err != nil {
		h.handleProviderPatternSaveError(rw, eventBuilder, userID, savedDesignByt, err, provider, token)
		return
	}

	if requestPayload.DesignFile.ID != uuid.Nil {
		eventBuilder = eventBuilder.WithAction(models.Update)
	} else {
		eventBuilder = eventBuilder.WithAction(models.Create)
	}
	description := fmt.Sprintf("Design %s saved at version %s", requestPayload.DesignFile.Name, requestPayload.DesignFile.Version)
	metadata := map[string]interface{}{
		"history_title": fmt.Sprintf("Version %s - %d components and %d relationships",
			requestPayload.DesignFile.Version, len(requestPayload.DesignFile.Components), len(requestPayload.DesignFile.Relationships)),

		"design": map[string]interface{}{
			"name": requestPayload.DesignFile.Name,
			"id":   requestPayload.DesignFile.ID.String(),
		},
		"doclink": "https://docs.meshery.io/concepts/logical/designs",
	}
	event := eventBuilder.
		WithDescription(description).
		WithMetadata(metadata).
		Build()
	_ = provider.PersistEvent(*event, token)

	_, _ = rw.Write(savedDesignByt)

}

// Verifies and converts a pattern to design format if required.
// A pattern is required to be converted to design format iff,
// 1. pattern_file attribute is empty, and
// 2. The "type" (sourcetype/original content) is not Design. [is one of compose/helmchart/manifests]
func (h *Handler) VerifyAndConvertToDesign(
	ctx context.Context,
	mesheryPattern *models.MesheryPattern,
	provider models.Provider,
) error {
	// Only proceed if we need to convert a non-design pattern that doesn't have a pattern file yet
	if mesheryPattern.Type.Valid && mesheryPattern.Type.String != string(core.MesheryDesign) && mesheryPattern.PatternFile == "" {
		token, ok := ctx.Value(models.TokenCtxKey).(string)
		if !ok {
			return ErrRetrieveUserToken(fmt.Errorf("failed to retrieve user token"))
		}

		sourceContent := mesheryPattern.SourceContent
		if len(mesheryPattern.SourceContent) == 0 {
			h.log.Info("Pattern file doesn't contain SourceContent, fetching from remote provider")
			var err error
			sourceContent, err = provider.GetDesignSourceContent(token, mesheryPattern.ID.String())
			if err != nil {
				return ErrDesignSourceContent(err, "get ")
			}
			mesheryPattern.SourceContent = sourceContent
		}

		fileToImport := FileToImport{
			Data:     sourceContent,
			FileName: mesheryPattern.Name, // Use pattern name as filename, make sure extension is there
		}

		// This function requires a valid file extension in the filename to work.
		// Note: Assuming MesheryPattern is already sanitized and identified, we are again going through
		// sanitization and identification (which is a redundant step but it's a one off)
		design, _, err := ConvertFileToDesign(fileToImport, h.registryManager, h.log)
		if err != nil {
			// ConvertFileToDesign returns a mix of bare and MeshKit-wrapped
			// errors depending on which import-pipeline stage failed. Wrap
			// so the caller's JSON error envelope always carries MeshKit
			// metadata, regardless of which stage surfaced the failure.
			return ErrConvertToDesign(err)
		}

		bytPattern, err := encoding.Marshal(design)
		if err != nil {
			return ErrEncodePattern(err)
		}
		mesheryPattern.PatternFile = string(bytPattern)

		resp, err := provider.SaveMesheryPattern(token, mesheryPattern)
		if err != nil {
			return ErrSavePattern(err)
		}

		contentMesheryPatternSlice := make([]models.MesheryPattern, 0)
		if err := json.Unmarshal(resp, &contentMesheryPatternSlice); err != nil {
			return models.ErrUnmarshal(err, "pattern")
		}
	}
	return nil
}

// Commenting out unused function
// func unCompressOCIArtifactIntoDesign(artifact []byte) (*models.MesheryPattern, error) {

// 	// Assume design is in OCI Tarball Format
// 	tmpDir, err := oci.CreateTempOCIContentDir()
// 	if err != nil {
// 		return nil, ErrCreateDir(err, "OCI")
// 	}
// 	defer os.RemoveAll(tmpDir)

// 	tmpInputDesignFile := filepath.Join(tmpDir, "design.tar")
// 	file, err := os.Create(tmpInputDesignFile)
// 	if err != nil {
// 		return nil, ErrCreateFile(err, tmpInputDesignFile)
// 	}
// 	defer file.Close()

// 	reader := bytes.NewReader(artifact)
// 	if _, err := io.Copy(file, reader); err != nil {
// 		return nil, ErrWritingIntoFile(err, tmpInputDesignFile)
// 	}

// 	tmpOutputDesignFile := filepath.Join(tmpDir, "output")
// 	// Extract the tarball
// 	if err := oci.UnCompressOCIArtifact(tmpInputDesignFile, tmpOutputDesignFile); err != nil {
// 		return nil, ErrUnCompressOCIArtifact(err)
// 	}

// 	files, err := walker.WalkLocalDirectory(tmpOutputDesignFile)
// 	if err != nil {
// 		return nil, ErrWaklingLocalDirectory(err)
// 	}

// 	// TODO: Add support to merge multiple designs into one
// 	// Currently, assumes to save only the first design
// 	if len(files) == 0 {
// 		return nil, ErrEmptyOCIImage(fmt.Errorf("no design file detected in the imported OCI image"))
// 	}
// 	design := files[0]

// 	var patternFile pattern.PatternFile

// 	err = encoding.Unmarshal([]byte(design.Content), &patternFile)
// 	if err != nil {
// 		return nil, ErrDecodePattern(err)
// 	}
// 	mesheryPattern := &models.MesheryPattern{
// 		PatternFile: design.Content,
// 		Name:        design.Name,
// 	}

// 	return mesheryPattern, nil
// }

// Commenting out unused function
// func githubRepoDesignScan(
// 	owner,
// 	repo,
// 	path,
// 	branch,
// 	sourceType string,
// 	reg *meshmodel.RegistryManager,
// ) ([]models.MesheryPattern, error) {
// 	var mu sync.Mutex
// 	ghWalker := walker.NewGit()
// 	result := make([]models.MesheryPattern, 0)
// 	err := ghWalker.
// 		Owner(owner).
// 		Repo(repo).
// 		Branch(branch).
// 		Root(path).
// 		RegisterFileInterceptor(func(f walker.File) error {
// 			ext := filepath.Ext(f.Name)
// 			var k8sres string
// 			var err error
// 			k8sres = f.Content
// 			if ext == ".yml" || ext == ".yaml" {
// 				if sourceType == string(models.DockerCompose) {
// 					k8sres, err = kompose.Convert([]byte(f.Content))
// 					if err != nil {
// 						return ErrRemoteApplication(err)
// 					}
// 				}
// 				pattern, err := pCore.NewPatternFileFromK8sManifest(k8sres, "", false, reg)
// 				if err != nil {
// 					return err //always a meshkit error
// 				}

// 				patternByt, _ := encoding.Marshal(pattern)

// 				af := models.MesheryPattern{
// 					Name:        strings.TrimSuffix(f.Name, ext),
// 					PatternFile: string(patternByt),
// 					Location: map[string]interface{}{
// 						"type":   "github",
// 						"host":   fmt.Sprintf("github.com/%s/%s", owner, repo),
// 						"path":   f.Path,
// 						"branch": branch,
// 					},
// 					Type: sql.NullString{
// 						String: string(sourceType),
// 						Valid:  true,
// 					},
// 					SourceContent: []byte(f.Content),
// 				}

// 				mu.Lock()
// 				result = append(result, af)
// 				mu.Unlock()
// 			}

// 			return nil
// 		}).
// 		Walk()

// 	return result, ErrRemoteApplication(err)
// }

// Commenting out unused function
// Always returns a meshery pattern slice of length 1 otherwise an error is returned
// func genericHTTPDesignFile(fileURL, patternName, sourceType string, reg *meshmodel.RegistryManager, log logger.Handler) ([]models.MesheryPattern, error) {
// 	resp, err := http.Get(fileURL)
// 	if err != nil {
// 		return nil, ErrRemoteApplication(err)
// 	}
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, ErrRemoteApplication(fmt.Errorf("file not found"))
// 	}

// 	defer models.SafeClose(resp.Body, log)

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, ErrRemoteApplication(err)
// 	}

// 	res := string(body)

// 	if sourceType == string(models.DockerCompose) {
// 		res, err = kompose.Convert(body)
// 		if err != nil {
// 			return nil, ErrRemoteApplication(err)
// 		}
// 	}

// 	var pattern pattern.PatternFile
// 	if sourceType == string(models.DockerCompose) || sourceType == string(models.K8sManifest) {
// 		var err error
// 		pattern, err = pCore.NewPatternFileFromK8sManifest(res, "", false, reg)
// 		if err != nil {
// 			return nil, err //This error is already a meshkit error
// 		}
// 	} else {
// 		err := encoding.Unmarshal([]byte(res), &pattern)
// 		if err != nil {
// 			return nil, ErrDecodePattern(err)
// 		}
// 	}

// 	if patternName != "" {
// 		pattern.Name = patternName
// 	}

// 	patternByt, _ := encoding.Marshal(pattern)

// 	url := strings.Split(fileURL, "/")
// 	af := models.MesheryPattern{
// 		Name:        url[len(url)-1],
// 		PatternFile: string(patternByt),
// 		Location: map[string]interface{}{
// 			"type":   "http",
// 			"host":   fileURL,
// 			"path":   "",
// 			"branch": "",
// 		},
// 		Type: sql.NullString{
// 			String: string(sourceType),
// 			Valid:  true,
// 		},
// 		SourceContent: body,
// 	}
// 	return []models.MesheryPattern{af}, nil
// }

func (h *Handler) GetMesheryPatternsHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()
	tokenString := r.Context().Value(models.TokenCtxKey).(string)
	updateAfter := q.Get("updated_after")
	includeMetrics := q.Get("metrics")
	err := r.ParseForm() // necessary to get r.Form["visibility"], i.e, ?visibility=public&visbility=private
	if err != nil {
		h.log.Error(ErrFetchPattern(err))
		writeMeshkitError(rw, ErrFetchPattern(err), http.StatusInternalServerError)
		return
	}
	filter := struct {
		Visibility []string `json:"visibility"`
	}{}

	visibility := q.Get("visibility")
	populate := q["populate"]
	if visibility != "" {
		err := json.Unmarshal([]byte(visibility), &filter.Visibility)
		if err != nil {
			h.log.Error(ErrFetchPattern(err))
			// The visibility JSON comes from the query string — a client-side
			// input, so a parse failure is a 400, not a 500.
			writeMeshkitError(rw, ErrFetchPattern(err), http.StatusBadRequest)
			return
		}
	}

	resp, err := provider.GetMesheryPatterns(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), updateAfter, filter.Visibility, includeMetrics, populate)

	if err != nil {
		h.log.Error(ErrFetchPattern(err))
		writeMeshkitError(rw, ErrFetchPattern(err), http.StatusInternalServerError)
		return
	}

	// mc := NewContentModifier(token, provider, prefObj, user.UserId)
	// //acts like a middleware, modifying the bytes lazily just before sending them back
	// err = mc.AddMetadataForPatterns(r.Context(), &resp)
	// if err != nil {
	// 	fmt.Println("Could not add metadata about pattern's current support ", err.Error())
	// }
	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) GetCatalogMesheryPatternsHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()
	tokenString := r.Context().Value(models.TokenCtxKey).(string)

	// Canonical form is `orgId`; `orgID` is dual-accepted during the Phase 2
	// deprecation window. Merge both lists (canonical first, legacy appended
	// unless duplicate) so the provider sees every caller's intended org
	// filter. Retire once Phase 3 consumer migration completes.
	orgIDs := append([]string{}, q["orgId"]...)
	for _, legacy := range q["orgID"] {
		if !slices.Contains(orgIDs, legacy) {
			orgIDs = append(orgIDs, legacy)
		}
	}

	resp, err := provider.GetCatalogMesheryPatterns(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("metrics"), q["populate"], q["class"], q["technology"], q["type"], orgIDs, q["workspaceID"], q["userid"])
	if err != nil {
		h.log.Error(ErrFetchPattern(err))
		writeMeshkitError(rw, ErrFetchPattern(err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) DeleteMesheryPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	patternID := mux.Vars(r)["id"]
	userID := user.ID
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction("delete").ActedUpon(uuid.FromStringOrNil(patternID))

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()
		_ = provider.PersistEvent(*event, "")
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(rw, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}

	mesheryPattern := models.MesheryPattern{}

	resp, err := provider.DeleteMesheryPattern(r, patternID)
	if err != nil {
		errPatternDelete := ErrDeletePattern(err)

		h.log.Error(errPatternDelete)
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": errPatternDelete,
		}).WithDescription("Error: Could not delete design.").Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(rw, errPatternDelete, http.StatusInternalServerError)
		return
	}

	_ = json.Unmarshal(resp, &mesheryPattern)
	event := eventBuilder.WithSeverity(events.Informational).WithDescription(fmt.Sprintf("Pattern %s deleted.", mesheryPattern.Name)).Build()
	_ = provider.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)
	go h.config.PatternChannel.Publish(user.ID, struct{}{})

	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		h.log.Error(err)
	}
}

// GetMesheryPatternHandler returns the pattern file with the given id

func (h *Handler) DownloadMesheryPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	var formatConverter converter.ConvertFormat
	userID := user.ID
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction("download").ActedUpon(userID).WithSeverity(events.Informational)

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()
		_ = provider.PersistEvent(*event, "")
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(rw, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}

	exportFormat := r.URL.Query().Get("export")
	h.log.Info(fmt.Sprintf("Export format received: '%s'", exportFormat))

	if exportFormat != "" {
		var errConvert error
		h.log.Info(fmt.Sprintf("Attempting to create converter for format: '%s'", exportFormat))

		h.log.Info(fmt.Sprintf("Available formats - K8sManifest: '%s', HelmChart: '%s'",
			converter.K8sManifest, converter.HelmChart))

		formatConverter, errConvert = converter.NewFormatConverter(converter.DesignFormat(exportFormat))
		if errConvert != nil {
			h.log.Info(fmt.Sprintf("Failed to create converter: %v", errConvert))
			err := ErrExportPatternInFormat(errConvert, exportFormat, "")
			h.log.Error(err)
			writeMeshkitError(rw, err, http.StatusBadRequest)
			return
		}
		h.log.Info(fmt.Sprintf("Successfully created converter for format: '%s'", exportFormat))
	}

	patternID := mux.Vars(r)["id"]
	ociFormat, _ := strconv.ParseBool(r.URL.Query().Get("oci"))
	ahpkg, _ := strconv.ParseBool(r.URL.Query().Get("pkg"))

	resp, err := provider.GetMesheryPattern(r, patternID, "false")
	if err != nil {
		h.log.Error(ErrGetPattern(err))
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrGetPattern(err),
		}).WithDescription(fmt.Sprintf("Failed to fetch design file for ID: %s.", patternID)).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)

		writeMeshkitError(rw, ErrGetPattern(err), http.StatusNotFound)
		return
	}
	pattern := &models.MesheryPattern{}
	err = json.Unmarshal(resp, &pattern)
	if err != nil {
		obj := "download pattern"
		h.log.Error(models.ErrUnmarshal(err, obj))
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": models.ErrUnmarshal(err, obj),
		}).WithDescription(fmt.Sprintf("Failed to unmarshal design file for ID: %s.", patternID)).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)

		writeMeshkitError(rw, models.ErrUnmarshal(err, obj), http.StatusInternalServerError)
		return
	}

	// publish a download event
	downloadEvent := events.DesignDownloadEvent(*pattern.ID, pattern.Name, userID, *h.SystemID)
	_ = provider.PersistEvent(*downloadEvent, token)
	go h.config.EventBroadcaster.Publish(userID, downloadEvent)

	err = h.VerifyAndConvertToDesign(r.Context(), pattern, provider)
	if err != nil {
		event := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("convert").WithDescription(fmt.Sprintf("The \"%s\" is not in the design format, failed to convert and persist the original source content from \"%s\" to design file format", pattern.Name, pattern.Type.String)).WithMetadata(map[string]interface{}{"error": err}).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		h.log.Error(err)
		writeMeshkitError(rw, err, http.StatusInternalServerError)
		return
	}

	// v1beta1
	isOldFormat, err := patternutils.IsDesignInAlpha2Format(pattern.PatternFile)
	if err != nil {
		err = ErrPatternFile(err)
		event := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("download").WithDescription(fmt.Sprintf("Failed to parse design \"%s\".", pattern.Name)).WithMetadata(map[string]interface{}{"error": err, "id": pattern.ID}).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		h.log.Error(err)
		writeMeshkitError(rw, err, http.StatusInternalServerError)
		return
	}

	if isOldFormat {

		eventBuilder := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("convert")
		_, patternFileStr, err := h.convertV1alpha2ToV1beta3(pattern, eventBuilder)
		event := eventBuilder.Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		if err != nil {
			h.log.Error(err)
			writeMeshkitError(rw, err, http.StatusInternalServerError)
			return
		}

		pattern.PatternFile = patternFileStr
	}

	if formatConverter != nil {
		patternFile, err := formatConverter.Convert(pattern.PatternFile)
		if err != nil {
			err = ErrExportPatternInFormat(err, exportFormat, pattern.Name)
			h.log.Error(err)
			writeMeshkitError(rw, err, http.StatusInternalServerError)
			return
		}
		if exportFormat == string(core.HelmChart) {
			rw.Header().Set("Content-Type", "application/x-gzip")
			rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.tgz", pattern.Name))
		} else {
			rw.Header().Set("Content-Type", "application/yaml")
			rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.yml", pattern.Name))
		}
		_, err = fmt.Fprint(rw, patternFile)
		if err != nil {
			// Headers were already committed above (Content-Type:
			// application/yaml or application/x-gzip), so we cannot send
			// a fresh JSON error response. Log only.
			err = ErrWriteResponse(_errors.Wrapf(err, "failed to export design %q in %s format", pattern.Name, exportFormat))
			h.log.Error(err)
			return
		}
		return
	}

	if ociFormat {
		tmpDir, err := oci.CreateTempOCIContentDir()
		if err != nil {
			h.log.Error(ErrCreateDir(err, "OCI"))

			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrCreateDir(err, "OCI"),
			}).WithDescription("Error creating tmp directory under ~/.meshery/content/").Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, ErrCreateDir(err, "OCI"), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				h.log.Error(err)
			}
		}()

		tmpDesignFile := filepath.Join(tmpDir, "design.yml")
		file, err := os.Create(tmpDesignFile)
		if err != nil {
			h.log.Error(ErrCreateFile(err, tmpDesignFile))
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrCreateFile(err, tmpDesignFile),
			}).WithDescription(fmt.Sprintf("Error creating tmp file %s", tmpDesignFile)).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, ErrCreateFile(err, tmpDesignFile), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				h.log.Error(err)
			}
		}()

		var designPattern design.PatternFile

		err = encoding.Unmarshal([]byte(pattern.PatternFile), &designPattern)

		if err != nil {

			err = ErrEncodePattern(err)
			h.log.Error(err)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": err,
			}).WithDescription(fmt.Sprintf("Failed to export design \"%s\" as OCI image.", pattern.Name)).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, err, http.StatusInternalServerError)
			return
		}

		ymlDesign, err := yaml.Marshal(designPattern)

		if err != nil {
			err = ErrEncodePattern(err)
			h.log.Error(err)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": err,
			}).WithDescription(fmt.Sprintf("Failed to export design \"%s\" as OCI image.", pattern.Name)).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, err, http.StatusInternalServerError)
			return
		}

		if _, err := file.Write(ymlDesign); err != nil {
			h.log.Error(ErrWritingIntoFile(err, tmpDesignFile))
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrWritingIntoFile(err, tmpDesignFile),
			}).WithDescription(fmt.Sprintf("Error writing into tmp design file %s at %s", pattern.Name, tmpDesignFile)).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, ErrWritingIntoFile(err, tmpDesignFile), http.StatusInternalServerError)
			return
		}

		artifactHubPkgFilePath := filepath.Join(tmpDir, "artifacthub-pkg.yml")
		artifactHubPkgFile, err := os.Create(artifactHubPkgFilePath)
		if err != nil {
			h.log.Error(ErrCreateFile(err, "artifacthub-pkg.yml"))
			eb := *eventBuilder
			event := eb.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to create artifacthub pkg for the design \"%s\"", pattern.Name)).WithMetadata(map[string]interface{}{"error": err}).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)
		}

		data, err := createArtifactHubPkg(pattern, strings.Trim(fmt.Sprintf("%s %s", user.FirstName, user.LastName), " "))
		if err != nil {
			h.log.Error(err)
			eb := *eventBuilder
			event := eb.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to create artifacthub pkg for the design \"%s\"", pattern.Name)).WithMetadata(map[string]interface{}{"error": err}).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)
		}

		_, err = artifactHubPkgFile.Write(data)
		if err != nil {
			err = ErrWritingIntoFile(err, "artifacthub-pkg.yml")
			h.log.Error(err)
			eb := *eventBuilder
			event := eb.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to create artifacthub pkg for the design \"%s\"", pattern.Name)).WithMetadata(map[string]interface{}{"error": err}).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)
		}

		ociImg, err := oci.BuildImage(tmpDir)
		if err != nil {
			h.log.Error(ErrBuildOCIImg(err))
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrBuildOCIImg(err),
			}).WithDescription(fmt.Sprintf("Error building OCI Image from contents at %s", tmpDesignFile)).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, ErrBuildOCIImg(err), http.StatusInternalServerError)
			return
		}

		digest, err := ociImg.Digest()
		if err != nil {
			h.log.Error(ErrBuildOCIImg(err))
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrBuildOCIImg(err),
			}).WithDescription(fmt.Sprintf("Error getting image digest for %s", tmpDesignFile)).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, ErrBuildOCIImg(err), http.StatusInternalServerError)
			return
		}

		size, err := ociImg.Size()
		if err != nil {
			h.log.Error(ErrBuildOCIImg(err))
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrBuildOCIImg(err),
			}).WithDescription(fmt.Sprintf("Error getting calculating image size for %s", tmpDesignFile)).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, ErrBuildOCIImg(err), http.StatusInternalServerError)
			return
		}

		h.log.Info(fmt.Sprintf("OCI Image built. Digest: %v, Size: %v", digest, size))

		eventBuilder.WithSeverity(events.Informational).WithDescription(fmt.Sprintf("OCI Image built. Digest: %v, Size: %v", digest, size))
		event := eventBuilder.Build()
		go h.config.EventBroadcaster.Publish(userID, event)
		_ = provider.PersistEvent(*event, token)

		pretifiedName := strings.ToLower(strings.ReplaceAll(pattern.Name, " ", "")) // ensures that tag validation passes
		tmpOCITarFilePath := filepath.Join(tmpDir, pretifiedName+".tar")
		err = oci.SaveOCIArtifact(ociImg, tmpOCITarFilePath, pretifiedName)
		if err != nil {
			h.log.Error(ErrSaveOCIArtifact(err))
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrSaveOCIArtifact(err),
			}).WithDescription(fmt.Sprintf("Failed to save OCI Artifact %s temporarily", tmpOCITarFilePath)).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, ErrSaveOCIArtifact(err), http.StatusInternalServerError)
			return
		}

		file, err = os.OpenFile(tmpOCITarFilePath, os.O_RDONLY, 0444)
		if err != nil {
			h.log.Error(ErrOpenFile(tmpOCITarFilePath))
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrOpenFile(tmpOCITarFilePath),
			}).WithDescription(fmt.Sprintf("Failed to read contents of OCI Artifact %s", tmpOCITarFilePath)).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, ErrOpenFile(tmpOCITarFilePath), http.StatusInternalServerError)
			return
		}
		content, err := io.ReadAll(file)
		if err != nil {
			h.log.Error(ErrIOReader(err))
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrIOReader(err),
			}).WithDescription(fmt.Sprintf("Failed to read contents of OCI artifact %s", tmpOCITarFilePath)).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, ErrIOReader(err), http.StatusInternalServerError)
			return
		}

		h.log.Info("OCI Artifact saved at: ", tmpOCITarFilePath)

		eventBuilder.WithSeverity(events.Informational).WithDescription(fmt.Sprintf("OCI Artifact temporarily saved at: %s", tmpOCITarFilePath))
		event = eventBuilder.Build()
		go h.config.EventBroadcaster.Publish(userID, event)
		_ = provider.PersistEvent(*event, token)

		rw.Header().Set("Content-Type", "application/tar")
		rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.tar", pattern.Name))

		reader := bytes.NewReader(content)
		if _, err := io.Copy(rw, reader); err != nil {
			// Headers have already been committed and the tar stream has
			// started — we cannot send a fresh JSON error response here.
			// Log and emit the event, then return.
			h.log.Error(ErrIOReader(err))
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrIOReader(err),
			}).WithDescription("Failed to save contents of OCI Artifact at requested path").Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}
		return
	}

	if ahpkg {
		rw.Header().Set("Content-Type", "application/zip")
		rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.zip", pattern.Name))

		tarWriter := utils.NewTarWriter()
		data, _ := createArtifactHubPkg(pattern, strings.Trim(fmt.Sprintf("%s %s", user.FirstName, user.LastName), " "))
		err = tarWriter.Compress("artifacthub-pkg.yml", data)
		if err != nil {
			h.log.Error(err)
			eb := *eventBuilder
			event := eb.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to create artifacthub pkg for the design \"%s\"", pattern.Name)).WithMetadata(map[string]interface{}{"error": err}).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)
		}
		ymlDesign, err := yaml.Marshal(pattern.PatternFile)
		if err != nil {
			err = ErrEncodePattern(err)
			h.log.Error(err)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": err,
			}).WithDescription(fmt.Sprintf("Failed to export design \"%s\" as OCI image.", pattern.Name)).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)

			writeMeshkitError(rw, err, http.StatusInternalServerError)
			return
		}
		err = tarWriter.Compress(pattern.Name+".yml", ymlDesign)
		if err != nil {
			wrappedErr := ErrCompressArtifact(err)
			h.log.Error(wrappedErr)
			eb := *eventBuilder
			event := eb.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to compress design \"%s\" and artifacthub pkg.", pattern.Name)).WithMetadata(map[string]interface{}{"error": wrappedErr}).Build()
			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)
			tarWriter.Close()
			writeMeshkitError(rw, wrappedErr, http.StatusInternalServerError)
			return
		}

		tarWriter.Close()

		_, err = io.Copy(rw, tarWriter.Buffer)
		if err != nil {
			// Headers were committed above (Content-Type: application/zip),
			// so a fresh JSON error response cannot be sent. Log only.
			h.log.Error(ErrIOReader(err))
			return
		}
	}

	yamlBytes, err := encoding.ToYaml([]byte(pattern.PatternFile))
	if err != nil {
		err = ErrParsePattern(err)
		h.log.Error(err)
		writeMeshkitError(rw, err, http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/yaml")
	rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.yml", pattern.Name))

	_, err = rw.Write(yamlBytes)
	if err != nil {
		// Headers were committed above (Content-Type: application/yaml),
		// so a fresh JSON error response cannot be sent. Log only.
		h.log.Error(ErrEncodePattern(err))
		return
	}
}

func (h *Handler) CloneMesheryPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	patternID := mux.Vars(r)["id"]
	patternUUID := uuid.FromStringOrNil(patternID)

	userID := user.ID
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction("clone").ActedUpon(patternUUID).WithSeverity(events.Informational)

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()
		_ = provider.PersistEvent(*event, "")
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(rw, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}

	var parsedBody *models.MesheryClonePatternRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil || patternID == "" {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	mesheryPatternByt, err := provider.GetMesheryPattern(r, patternID, "false")
	if err != nil {
		h.log.Error(ErrGetPattern(err))
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrGetPattern(err),
		}).WithDescription(fmt.Sprintf("Failed to fetch meshery pattern \"%s\" with id: %s.", parsedBody.Name, patternID)).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)

		writeMeshkitError(rw, ErrGetPattern(err), http.StatusNotFound)
		return
	}

	pattern := &models.MesheryPattern{}
	err = json.Unmarshal(mesheryPatternByt, &pattern)
	if err != nil {
		obj := "pattern: " + patternID
		h.log.Error(models.ErrUnmarshal(err, obj))
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": models.ErrUnmarshal(err, obj),
		}).WithDescription(fmt.Sprintf("Failed to fetch meshery pattern \"%s\" with ID: %s.", parsedBody.Name, patternID)).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)

		writeMeshkitError(rw, models.ErrUnmarshal(err, obj), http.StatusInternalServerError)
		return
	}

	isOldFormat, err := patternutils.IsDesignInAlpha2Format(pattern.PatternFile)
	if err != nil {
		err = ErrPatternFile(err)
		event := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithDescription(fmt.Sprintf("Failed to parse design \"%s\".", pattern.Name)).WithMetadata(map[string]interface{}{"error": err, "id": patternID}).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		h.log.Error(err)
		writeMeshkitError(rw, err, http.StatusInternalServerError)
		return
	}

	if isOldFormat {
		eventBuilder := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("convert")
		_, patternFileStr, err := h.convertV1alpha2ToV1beta3(pattern, eventBuilder)
		event := eventBuilder.Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		if err != nil {
			h.log.Error(err)
			writeMeshkitError(rw, err, http.StatusInternalServerError)
			return
		}
		pattern.PatternFile = patternFileStr
		_, err = provider.SaveMesheryPattern(token, pattern)
		if err != nil {
			h.log.Error(ErrSavePattern(err))

			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrSavePattern(_errors.Wrapf(err, "failed to persist converted v1beta3 design file \"%s\" with id: %s", parsedBody.Name, patternID)),
			}).WithDescription(ErrSavePattern(err).Error()).Build()

			_ = provider.PersistEvent(*event, token)
			go h.config.EventBroadcaster.Publish(userID, event)
			writeMeshkitError(rw, ErrSavePattern(err), http.StatusInternalServerError)
			return
		}
	}

	resp, err := provider.CloneMesheryPattern(r, patternID, parsedBody)
	if err != nil {
		h.log.Error(ErrClonePattern(err))
		writeMeshkitError(rw, ErrClonePattern(err), http.StatusInternalServerError)
		return
	}

	go h.config.PatternChannel.Publish(user.ID, struct{}{})
	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) PublishCatalogPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	userID := user.ID
	eventBuilder := events.NewEvent().
		FromUser(userID).
		FromSystem(*h.SystemID).
		WithCategory("pattern").
		WithAction("publish").
		ActedUpon(userID)

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()
		_ = provider.PersistEvent(*event, "")
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(rw, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}

	var parsedBody *models.MesheryCatalogPatternRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRequestBody(err))
		e := eventBuilder.WithSeverity(events.Error).
			WithMetadata(map[string]interface{}{
				"error": ErrRequestBody(err),
			}).
			WithDescription("Error parsing design payload.").Build()
		_ = provider.PersistEvent(*e, token)
		go h.config.EventBroadcaster.Publish(userID, e)
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	resp, err := provider.PublishCatalogPattern(r, parsedBody)
	if err != nil {
		h.log.Error(ErrPublishCatalogPattern(err))
		e := eventBuilder.WithSeverity(events.Error).
			WithMetadata(map[string]interface{}{
				"error": ErrPublishCatalogPattern(err),
			}).
			WithDescription("Error publishing design.").Build()
		_ = provider.PersistEvent(*e, token)
		go h.config.EventBroadcaster.Publish(userID, e)
		writeMeshkitError(rw, ErrPublishCatalogPattern(err), http.StatusInternalServerError)
		return
	}

	var respBody *models.CatalogRequest
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		h.log.Error(ErrPublishCatalogPattern(err))
		e := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrPublishCatalogPattern(err),
		}).WithDescription("Error parsing response.").Build()
		_ = provider.PersistEvent(*e, token)
		go h.config.EventBroadcaster.Publish(userID, e)
		writeMeshkitError(rw, ErrPublishCatalogPattern(err), http.StatusInternalServerError)
		return
	}

	e := eventBuilder.WithSeverity(events.Informational).ActedUpon(parsedBody.ID).WithDescription(fmt.Sprintf("Request to publish '%s' design submitted with status: %s", respBody.ContentName, respBody.Status)).Build()
	_ = provider.PersistEvent(*e, token)
	go h.config.EventBroadcaster.Publish(userID, e)

	go h.config.PatternChannel.Publish(user.ID, struct{}{})
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusAccepted)
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) UnPublishCatalogPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	userID := user.ID
	eventBuilder := events.NewEvent().
		FromUser(userID).
		FromSystem(*h.SystemID).
		WithCategory("pattern").
		WithAction("unpublish_request").
		ActedUpon(userID)

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()
		_ = provider.PersistEvent(*event, "")
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(rw, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}

	var parsedBody *models.MesheryCatalogPatternRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRequestBody(err))
		e := eventBuilder.WithSeverity(events.Error).
			WithMetadata(map[string]interface{}{
				"error": ErrRequestBody(err),
			}).
			WithDescription("Error parsing design payload.").Build()
		_ = provider.PersistEvent(*e, token)
		go h.config.EventBroadcaster.Publish(userID, e)
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	resp, err := provider.UnPublishCatalogPattern(r, parsedBody)
	if err != nil {
		h.log.Error(ErrPublishCatalogPattern(err))
		e := eventBuilder.WithSeverity(events.Error).
			WithMetadata(map[string]interface{}{
				"error": ErrPublishCatalogPattern(err),
			}).
			WithDescription("Error publishing design.").Build()
		_ = provider.PersistEvent(*e, token)
		go h.config.EventBroadcaster.Publish(userID, e)
		writeMeshkitError(rw, ErrPublishCatalogPattern(err), http.StatusInternalServerError)
		return
	}

	var respBody *models.CatalogRequest
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		h.log.Error(ErrPublishCatalogPattern(err))
		e := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrPublishCatalogPattern(err),
		}).WithDescription("Error parsing response.").Build()
		_ = provider.PersistEvent(*e, token)
		go h.config.EventBroadcaster.Publish(userID, e)
		writeMeshkitError(rw, ErrPublishCatalogPattern(err), http.StatusInternalServerError)
		return
	}

	e := eventBuilder.WithSeverity(events.Informational).ActedUpon(parsedBody.ID).WithDescription(fmt.Sprintf("'%s' design unpublished", respBody.ContentName)).Build()
	_ = provider.PersistEvent(*e, token)
	go h.config.EventBroadcaster.Publish(userID, e)

	go h.config.PatternChannel.Publish(user.ID, struct{}{})
	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) DeleteMultiMesheryPatternsHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error(models.ErrDataRead(err, "Request Body"))
	}
	var patterns models.MesheryPatternDeleteRequestBody
	err = json.Unmarshal([]byte(body), &patterns)
	if err != nil {
		h.log.Error(models.ErrMarshal(err, "pattern"))
	}

	h.log.Debug("patterns to be deleted: ", patterns)

	resp, err := provider.DeleteMesheryPatterns(r, patterns)

	if err != nil {
		h.log.Error(ErrDeletePattern(err))
		writeMeshkitError(rw, ErrDeletePattern(err), http.StatusInternalServerError)
		return
	}
	go h.config.PatternChannel.Publish(user.ID, struct{}{})
	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(resp)); err != nil {
		h.log.Error(err)
	}
}

// GetMesheryPatternHandler fetched the pattern with the given id
func (h *Handler) GetMesheryPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	patternID := mux.Vars(r)["id"]
	patternUUID := uuid.FromStringOrNil(patternID)
	userID := user.ID
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction("view").ActedUpon(patternUUID)

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()
		_ = provider.PersistEvent(*event, "")
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(rw, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}

	resp, err := provider.GetMesheryPattern(r, patternID, r.URL.Query().Get("metrics"))
	if err != nil {
		h.handleProviderPatternGetError(rw, eventBuilder, userID, resp, err, provider, token)
		return
	}

	pattern := &models.MesheryPattern{}
	err = json.Unmarshal(resp, &pattern)
	if err != nil {
		h.log.Error(ErrGetPattern(err))
		writeMeshkitError(rw, ErrGetPattern(err), http.StatusInternalServerError)
		return
	}

	err = h.VerifyAndConvertToDesign(r.Context(), pattern, provider)
	if err != nil {
		event := events.NewEvent().ActedUpon(patternUUID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("convert").WithDescription(fmt.Sprintf("The \"%s\" is not in the design format, failed to convert and persist the original source content from \"%s\" to design file format", pattern.Name, pattern.Type.String)).WithMetadata(map[string]interface{}{"error": err}).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		h.log.Error(err)
		writeMeshkitError(rw, err, http.StatusInternalServerError)
		return
	}

	isOldFormat, err := patternutils.IsDesignInAlpha2Format(pattern.PatternFile)
	if err != nil {
		err = ErrPatternFile(err)
		event := events.NewEvent().ActedUpon(patternUUID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("view").WithDescription(fmt.Sprintf("Failed to parse design \"%s\".", pattern.Name)).WithMetadata(map[string]interface{}{"error": err, "id": pattern.ID}).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		h.log.Error(err)
		writeMeshkitError(rw, err, http.StatusInternalServerError)
		return
	}

	if isOldFormat {
		eventBuilder := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("convert")
		_, patternFileStr, err := h.convertV1alpha2ToV1beta3(pattern, eventBuilder)
		event := eventBuilder.Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		if err != nil {
			h.log.Error(err)
			writeMeshkitError(rw, err, http.StatusInternalServerError)
			return
		}
		pattern.PatternFile = patternFileStr
	}

	// deprettify pattern for backward compatibility with older designs which had the configuration in prettified format
	var designPattern design.PatternFile
	err = encoding.Unmarshal([]byte(pattern.PatternFile), &designPattern)

	if err != nil {
		err = ErrParsePattern(err)
		h.log.Error(err)
		writeMeshkitError(rw, err, http.StatusInternalServerError)
		return
	}

	for _, component := range designPattern.Components {
		component.Configuration = patterncore.Format.DePrettify(component.Configuration, false)
	}

	patternBytes, err := encoding.Marshal(designPattern)
	pattern.PatternFile = string(patternBytes)
	// done deprettifying

	if err != nil {
		h.log.Error(ErrEncodePattern(err))
		writeMeshkitError(rw, ErrEncodePattern(err), http.StatusInternalServerError)
		return
	}

	// done deprettifying

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(pattern); err != nil {
		// Content-Type header is already set, so the response has started.
		// Log only — a fresh JSON error response cannot be sent.
		h.log.Error(ErrEncodeResponse(err))
		return
	}
}

func (h *Handler) formatPatternOutput(rw http.ResponseWriter, content []byte, format, sourcetype string, eventBuilder *events.EventBuilder, URL, action string) {
	contentMesheryPatternSlice := make([]models.MesheryPattern, 0)

	if err := json.Unmarshal(content, &contentMesheryPatternSlice); err != nil {
		h.log.Error(ErrDecodePattern(err))
		writeMeshkitError(rw, ErrDecodePattern(err), http.StatusInternalServerError)
		return
	}

	result := []models.MesheryPattern{}
	names := []string{}
	for _, content := range contentMesheryPatternSlice {
		if content.ID != nil {
			eventBuilder.ActedUpon(*content.ID)
		}

		result = append(result, content)
		names = append(names, content.Name)
	}

	data, err := json.Marshal(&result)
	if err != nil {
		obj := "pattern file"
		h.log.Error(models.ErrMarshal(err, obj))
		writeMeshkitError(rw, models.ErrMarshal(err, obj), http.StatusInternalServerError)
		return
	}
	var response string
	if URL == "" {
		actionDesc := "updated"
		if action == models.Create {
			actionDesc = "created"
		}
		response = fmt.Sprintf("%s \"%s\" %s.", sourcetype, strings.Join(names, ","), actionDesc)
	} else {
		response = fmt.Sprintf("%s \"%s\" imported from URL %s", sourcetype, strings.Join(names, ","), URL)
	}
	eventBuilder.WithDescription(response)
	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(data)); err != nil {
		h.log.Error(err)
	}
}

// Since the client currently does not support pattern imports and externalized variables, the first(import) stage of pattern engine
// is evaluated here to simplify the pattern file such that it is valid when a deploy takes place

//unsued currently

// func evalImportAndReferenceStage(p *pattern.PatternFile) (newp pattern.PatternFile) {
// 	chain := stages.CreateChain()
// 	chain.
// 		// Add(stages.Import(sip, sap)). enable this
// 		Add(stages.Filler(false)).
// 		Add(func(data *stages.Data, err error, next stages.ChainStageNextFunction) {
// 			data.Lock.Lock()
// 			newp = *data.Pattern
// 			data.Lock.Unlock()
// 		}).
// 		Process(&stages.Data{
// 			Pattern: p,
// 		})
// 	return newp
// }

// Only pass Meshkit err here or there will be a panic
func addMeshkitErr(res *meshes.EventsResponse, err error) {
	if err != nil {
		res.EventType = meshes.EventType_ERROR
		res.ProbableCause = errors.GetCause(err)
		res.SuggestedRemediation = errors.GetRemedy(err)
		res.Details = err.Error()
		res.Summary = errors.GetSDescription(err)
		res.ErrorCode = errors.GetCode(err)
	}
}

func (h *Handler) handlePatternUpdate(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()
	userID := user.ID
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction("update").FromUser(userID)

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()
		_ = provider.PersistEvent(*event, "")
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(rw, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}

	res := meshes.EventsResponse{
		Component:     "core",
		ComponentName: "Design",
		OperationId:   guid.NewString(),
		EventType:     meshes.EventType_INFO,
	}

	sourcetype := mux.Vars(r)["sourcetype"]
	if sourcetype == "" {
		// Domain-correct error: this is the pattern (design) handler, so a
		// missing source-type should surface as a pattern-save failure, not
		// the legacy "application" domain error.
		missingRouteErr := ErrMissingRouteVariable("sourcetype", string(models.K8sManifest), string(models.DockerCompose), string(models.HelmChart))
		h.log.Error(missingRouteErr)

		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": missingRouteErr,
		}).WithDescription("Please provide design source-type").Build()

		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		go h.EventsBuffer.Publish(&res)
		writeMeshkitError(rw, missingRouteErr, http.StatusBadRequest)
		return
	}

	var parsedBody *MesheryPatternUPDATERequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRetrieveData(err))
		writeMeshkitError(rw, ErrRetrieveData(err), http.StatusBadRequest)
		return
	}

	actedUpon := &userID
	if parsedBody.PatternData != nil && parsedBody.PatternData.ID != nil {
		actedUpon = parsedBody.PatternData.ID
	}

	eventBuilder.ActedUpon(*actedUpon)

	format := r.URL.Query().Get("output")

	mesheryPattern := parsedBody.PatternData
	mesheryPattern.Type = sql.NullString{
		String: sourcetype,
		Valid:  true,
	}
	resp, err := provider.SaveMesheryPattern(token, mesheryPattern)
	if err != nil {
		errPatternSave := ErrSavePattern(err)
		h.log.Error(errPatternSave)

		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": errPatternSave,
		}).WithDescription(fmt.Sprintf("Error saving design %s", parsedBody.PatternData.Name)).Build()

		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)

		writeMeshkitError(rw, errPatternSave, http.StatusInternalServerError)
		return
	}
	go h.config.PatternChannel.Publish(userID, struct{}{})

	eventBuilder = eventBuilder.WithSeverity(events.Informational).ActedUpon(*mesheryPattern.ID)

	h.formatPatternOutput(rw, resp, format, sourcetype, eventBuilder, parsedBody.URL, models.Update)
	event := eventBuilder.Build()
	_ = provider.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)

}

// PatternFileRequestHandler will handle requests of both type GET and POST
// on the route /api/pattern
func (h *Handler) DesignFileRequestHandlerWithSourceType(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	if r.Method == http.MethodPost {
		h.handlePatternPOST(rw, r, prefObj, user, provider)
		return
	}

	if r.Method == http.MethodPut {
		h.handlePatternUpdate(rw, r, prefObj, user, provider)
		return
	}
}

func (h *Handler) GetMesheryDesignTypesHandler(
	rw http.ResponseWriter,
	_ *http.Request,
	_ *models.Preference,
	_ *models.User,
	_ models.Provider,
) {
	response := models.GetDesignsTypes()
	b, err := json.Marshal(response)
	if err != nil {
		obj := "available types"
		h.log.Error(models.ErrMarshal(err, obj))
		writeMeshkitError(rw, models.ErrMarshal(err, obj), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(rw, string(b)); err != nil {
		h.log.Error(err)
	}
}

// GetMesheryPatternHandler fetched the design using the given id and sourcetype
func (h *Handler) GetMesheryPatternSourceHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	designID := mux.Vars(r)["id"]
	token, ok := r.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		tokenErr := ErrRetrieveUserToken(fmt.Errorf("failed to retrieve user token"))
		h.log.Error(tokenErr)
		writeMeshkitError(rw, tokenErr, http.StatusInternalServerError)
		return
	}

	resp, err := provider.GetDesignSourceContent(token, designID)
	if err != nil {
		h.log.Error(ErrGetPattern(err))
		writeMeshkitError(rw, ErrGetPattern(err), http.StatusNotFound)
		return
	}

	var mimeType string
	sourcetype := mux.Vars(r)["sourcetype"]

	switch models.DesignType(sourcetype) {
	case models.HelmChart:
		mimeType = "application/x-tar"
	default: // docker-compose, k8smanifest
		mimeType = "application/x-yaml"
	}

	reader := bytes.NewReader(resp)
	rw.Header().Set("Content-Type", mimeType)
	_, err = io.Copy(rw, reader)
	if err != nil {
		// Headers were committed above (Content-Type: application/x-tar or
		// application/x-yaml), so a fresh JSON error response cannot be
		// sent. Log only.
		h.log.Error(ErrApplicationSourceContent(err, "download"))
	}
}

func createArtifactHubPkg(pattern *models.MesheryPattern, user string) ([]byte, error) {
	isCatalogItem := pattern.Visibility == models.Published
	var version string
	if isCatalogItem {
		version = pattern.CatalogData.PublishedVersion
	}
	artifactHubPkg := catalog.BuildArtifactHubPkg(pattern.Name, "", user, version, pattern.CreatedAt, &pattern.CatalogData)

	data, err := yaml.Marshal(artifactHubPkg)
	if err != nil {
		return nil, models.ErrMarshalYAML(err, "artifacthub-pkg")
	}

	return data, nil
}

func (h *Handler) convertV1alpha2ToV1beta3(mesheryPattern *models.MesheryPattern, eventBuilder *events.EventBuilder) (*design.PatternFile, string, error) {

	v1alpha1PatternFile := v1alpha2.PatternFile{}

	v1beta3PatternFile := design.PatternFile{}

	err := encoding.Unmarshal([]byte(mesheryPattern.PatternFile), &v1alpha1PatternFile)
	if err != nil {
		return nil, "", ErrParsePattern(err)
	}

	for _, svc := range v1alpha1PatternFile.Services {
		svc.Settings = helpers.RecursiveCastMapStringInterfaceToMapStringInterface(svc.Settings)
		svc.Traits = helpers.RecursiveCastMapStringInterfaceToMapStringInterface(svc.Traits)
	}

	err = v1beta3PatternFile.ConvertFrom(&v1alpha1PatternFile)
	if err != nil {
		// Wrap so the JSON error envelope returned to the client carries
		// MeshKit metadata (code, severity, suggested_remediation) instead
		// of just the bare conversion message.
		wrappedErr := ErrConvertPattern(err)
		eventBuilder.WithSeverity(events.Error).
			WithDescription(fmt.Sprintf("Failed to migrate design \"%s\" to current schema version", mesheryPattern.Name)).
			WithMetadata(map[string]interface{}{"error": wrappedErr})
		return nil, "", wrappedErr
	}

	v1beta3PatternFile.ID = *mesheryPattern.ID
	v1beta3PatternFile.Version = v1alpha1PatternFile.Version

	h.log.Infof("Converted design file with id \"%s\" to v1beta3 format", *mesheryPattern.ID)

	err = mapModelRelatedData(h.registryManager, &v1beta3PatternFile)
	if err != nil {
		wrappedErr := ErrGetComponentDefinition(err)
		eventBuilder.WithSeverity(events.Error).
			WithDescription(fmt.Sprintf("Failed to enrich converted design \"%s\" with styles and metadata", mesheryPattern.Name)).
			WithMetadata(map[string]interface{}{"error": wrappedErr, "id": *mesheryPattern.ID})
		return nil, "", wrappedErr
	}

	v1beta3PatternByt, err := encoding.Marshal(v1beta3PatternFile)
	if err != nil {
		return nil, "", utils.ErrMarshal(err)
	}
	eventBuilder.WithSeverity(events.Informational).WithDescription(fmt.Sprintf("Converted design file \"%s\" with id \"%s\" to v1beta3 format", mesheryPattern.Name, *mesheryPattern.ID))
	return &v1beta3PatternFile, string(v1beta3PatternByt), nil
}

func mapModelRelatedData(reg *meshmodel.RegistryManager, patternFile *design.PatternFile) error {
	s := selector.New(reg)
	for _, comp := range patternFile.Components {
		if comp == nil {
			continue
		}

		wc, err := s.GetDefinition(comp.Component.Kind, comp.Model.Model.Version, comp.Model.Name, comp.Component.Version, true)
		if err != nil {
			m := []string{"meshery", "meshery-core", "meshery-shapes", "meshery-flowchart"}
			// if model is one of those defined in the slice above as meshery, and no matching defs were found,
			// try to find the component just by name, this ensures the component is upgraded to newer model.
			// Eg: Some old designs contains "Comment" component under "meshery" model instead of "meshery-core"

			// Update the component kind to reflect the current registry.
			// Eg: The Connection component for k8s, had "kind" updated to "KuberntesConnection",hence any designs which has model k8s and kind "Connection" will fail, to ensure it gets converted, update the kind
			if comp.Model.Name == "kubernetes" && comp.Component.Kind == "Connection" {
				comp.Component.Kind = "KubernetesConnection"
			} else if comp.Model.Name == "aws" || comp.Model.Name == "gcp" {
				comp.Component.Kind = fmt.Sprintf("%s %s", strings.ToUpper(comp.Model.Name), comp.Component.Kind)
			} else if !slices.Contains(m, comp.Model.Name) {
				return err
			}

			entities, _, _, _ := reg.GetEntities(&regv1beta1.ComponentFilter{
				Name:       comp.Component.Kind,
				APIVersion: comp.Component.Version,
			})
			comp, found := selector.FindCompDefinitionWithVersion(entities, comp.Model.Model.Version)

			if found {
				wc = *comp
			}

		}

		comp.Model = wc.Model
		if wc.Model.Registrant.Status == "" {
			comp.Model.Registrant.Status = connection.ConnectionStatusRegistered
		}
		comp.Format = wc.Format
		comp.Version = wc.Version
		status := component.ComponentDefinitionStatus(wc.Model.Status)
		comp.Status = &status

		// Replace the SVG value with the  svg path. UI uses the path to fetch the SVG from the server.

		// helpers.WriteSVGsOnFileSystem()
		if comp.Model.Metadata.SvgComplete != nil && *comp.Model.Metadata.SvgComplete == "" {
			comp.Model.Metadata.SvgComplete = nil
		}
		comp.Capabilities = wc.Capabilities
		if comp.Capabilities == nil {
			comp.Capabilities = models.K8sMeshModelMetadata.Capabilities
		}
		comp.Metadata.Genealogy = wc.Metadata.Genealogy
		comp.Metadata.IsAnnotation = wc.Metadata.IsAnnotation
		comp.Metadata.Published = wc.Metadata.Published

		var styles core.ComponentStyles

		if comp.Styles != nil {
			styles = *comp.Styles
		} else {
			comp.Styles = &core.ComponentStyles{}
		}

		// Assign the other styles and reassign the position.
		if wc.Styles != nil {
			comp.Styles = wc.Styles
		}
		if styles.Position != nil {
			comp.Styles.Position = styles.Position
		}

	}

	return nil
}
