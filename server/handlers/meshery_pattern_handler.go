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
	"time"

	"github.com/gofrs/uuid"
	guid "github.com/google/uuid"
	"github.com/gorilla/mux"
	helpers "github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/meshes"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/resource/selector"
	patternutils "github.com/layer5io/meshery/server/models/pattern/utils"
	"github.com/layer5io/meshkit/encoding"
	"github.com/layer5io/meshkit/errors"
	"github.com/layer5io/meshkit/models/converter"
	_errors "github.com/pkg/errors"

	"github.com/layer5io/meshkit/models/catalog/v1alpha1"
	"github.com/layer5io/meshkit/models/events"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/models/oci"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/catalog"

	regv1beta1 "github.com/layer5io/meshkit/models/meshmodel/registry/v1beta1"
	coreV1 "github.com/meshery/schemas/models/v1alpha1/core"
	"github.com/meshery/schemas/models/v1alpha2"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	patternV1beta1 "github.com/meshery/schemas/models/v1beta1/pattern"
	"gopkg.in/yaml.v3"
)

// MesheryPatternRequestBody refers to the type of request body that
// SaveMesheryPattern would receive
// Deprecated
type MesheryPatternPOSTRequestBody struct {
	Name        string             `json:"name,omitempty"`
	URL         string             `json:"url,omitempty"`
	Path        string             `json:"path,omitempty"`
	Save        bool               `json:"save,omitempty"`
	PatternData *DesignPostPayload `json:"pattern_data,omitempty"`
}

type MesheryPatternUPDATERequestBody struct {
	Name          string                 `json:"name,omitempty"`
	URL           string                 `json:"url,omitempty"`
	Path          string                 `json:"path,omitempty"`
	Save          bool                   `json:"save,omitempty"`
	PatternData   *models.MesheryPattern `json:"pattern_data,omitempty"`
	CytoscapeJSON string                 `json:"cytoscape_json,omitempty"`
}

type DesignPostPayload struct {
	ID         *uuid.UUID                 `json:"id,omitempty"`
	Name       string                     `json:"name,omitempty"`
	DesignFile patternV1beta1.PatternFile `json:"design_file"`
	// Meshery doesn't have the user id fields
	// but the remote provider is allowed to provide one
	UserID      *string              `json:"user_id"`
	Visibility  string               `json:"visibility"`
	CatalogData v1alpha1.CatalogData `json:"catalog_data,omitempty"`
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

func (h *Handler) handleProviderPatternSaveError(rw http.ResponseWriter, eventBuilder *events.EventBuilder, userID uuid.UUID, body []byte, err error, provider models.Provider) {

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
		rw.Write(body)
		h.log.Error(&meshkitErr)
		event = eventBuilder.WithSeverity(events.Error).WithDescription(description).WithMetadata(map[string]interface{}{
			"error": meshkitErr,
		}).Build()

	} else {
		h.log.Error(ErrSavePattern(err))
		http.Error(rw, ErrSavePattern(err).Error(), http.StatusBadRequest)
		event = eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrSavePattern(err),
		}).WithDescription(ErrSavePattern(err).Error()).Build()
	}

	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)
	return
}

// swagger:route POST /api/pattern PatternsAPI idPostPatternFile
// Handle POST requests for patterns
//
// Edit/update a meshery pattern
// responses:
//
//	200: mesheryPatternResponseWrapper
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

	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction(models.Create).
		ActedUpon(userID).WithSeverity(events.Informational).WithDescription("Save design ")

	requestPayload := &DesignPostPayload{}
	if err := json.NewDecoder(r.Body).Decode(&requestPayload); err != nil {
		h.logErrorParsingRequestBody(rw, provider, err, userID, eventBuilder)
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.logErrorGettingUserToken(rw, provider, err, userID, eventBuilder)
		return
	}

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
		h.handleProviderPatternSaveError(rw, eventBuilder, userID, savedDesignByt, err, provider)
		return
	}

	if requestPayload.DesignFile.Id != uuid.Nil {
		eventBuilder = eventBuilder.WithAction(models.Update)
	} else {
		eventBuilder = eventBuilder.WithAction(models.Create)
	}
	metadata := map[string]interface{}{
		"design": map[string]interface{}{
			"name": requestPayload.DesignFile.Name,
			"id":   requestPayload.DesignFile.Id.String(),
		},
		"doclink": "https://docs.meshery.io/concepts/logical/designs",
	}
	event := eventBuilder.
		WithMetadata(metadata).
		Build()
	_ = provider.PersistEvent(event)

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
	if mesheryPattern.Type.Valid && mesheryPattern.Type.String != string(coreV1.MesheryDesign) && mesheryPattern.PatternFile == "" {
		token, _ := ctx.Value(models.TokenCtxKey).(string)

		sourceContent := mesheryPattern.SourceContent
		if len(mesheryPattern.SourceContent) == 0 {
			h.log.Info("Pattern file doesn't contain SourceContent, fetching from remote provider")
			sourceContent, err := provider.GetDesignSourceContent(token, mesheryPattern.ID.String())
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
			return err
		}

		bytPattern, err := encoding.Marshal(design)
		if err != nil {
			return err
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

// swagger:route GET /api/pattern PatternsAPI idGetPatternFiles
// Handle GET request for patterns
//
// Returns the list of all the patterns saved by the current user
// This will return all the patterns with their details
//
// ```?order={field}``` orders on the passed field
//
// ```?search=<design name>``` A string matching is done on the specified design name
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 10
//
// ```?visibility={[visibility]}``` Default visibility is public + private; Mulitple visibility filters can be passed as an array
// Eg: ```?visibility=["public", "published"]``` will return public and published designs
//
// ```?metrics``` Returns metrics like deployment/share/clone/view/download count for desings, default is false,
//
// / ```?populate``` Add the design content to the response like pattern_file return design file content
//
//	200: mesheryPatternsResponseWrapper
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
		http.Error(rw, ErrFetchPattern(err).Error(), http.StatusInternalServerError)
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
			http.Error(rw, ErrFetchPattern(err).Error(), http.StatusInternalServerError)
			return
		}
	}

	resp, err := provider.GetMesheryPatterns(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), updateAfter, filter.Visibility, includeMetrics, populate)

	if err != nil {
		h.log.Error(ErrFetchPattern(err))
		http.Error(rw, ErrFetchPattern(err).Error(), http.StatusInternalServerError)
		return
	}

	// mc := NewContentModifier(token, provider, prefObj, user.UserID)
	// //acts like a middleware, modifying the bytes lazily just before sending them back
	// err = mc.AddMetadataForPatterns(r.Context(), &resp)
	// if err != nil {
	// 	fmt.Println("Could not add metadata about pattern's current support ", err.Error())
	// }
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route GET /api/pattern/catalog PatternsAPI idGetCatalogMesheryPatternsHandler
// Handle GET request for catalog patterns
//
// # Patterns can be further filtered through query parameter
//
// ```?order={field}``` orders on the passed field
//
// ```?page={page-number}``` Default page number is 0
//
// ```?pagesize={pagesize}``` Default pagesize is 10.
//
// ```?search={patternname}``` If search is non empty then a greedy search is performed
//
// ```?metrics``` Returns metrics like deployment/share/clone/view/download count for designs, default false,
//
// ```?class={class}``` Filters patterns based on class
//
// ```?technology={technology}``` Filters patterns based on technology
//
// ```?type={type}``` Filters patterns based on type
//
// ```?orgID={orgID}``` Filters patterns based on organization ID
//
// ```?workspaceID={workspaceID}``` Filter patterns based on workspace ID
//
// ```?userid={userid}``` Filters patterns based on user ID
//
// responses:
//
// 200: mesheryPatternsResponseWrapper
func (h *Handler) GetCatalogMesheryPatternsHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()
	tokenString := r.Context().Value(models.TokenCtxKey).(string)

	resp, err := provider.GetCatalogMesheryPatterns(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("metrics"), q["populate"], q["class"], q["technology"], q["type"], q["orgID"], q["workspaceID"], q["userid"])
	if err != nil {
		h.log.Error(ErrFetchPattern(err))
		http.Error(rw, ErrFetchPattern(err).Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route DELETE /api/pattern/{id} PatternsAPI idDeleteMesheryPattern
// Handle Delete for a Meshery Pattern
//
// Deletes a meshery pattern with ID: id
// responses:
//
//	200: noContentWrapper
//
// DeleteMesheryPatternHandler deletes a pattern with the given id
func (h *Handler) DeleteMesheryPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	patternID := mux.Vars(r)["id"]
	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction("delete").ActedUpon(uuid.FromStringOrNil(patternID))

	mesheryPattern := models.MesheryPattern{}

	resp, err := provider.DeleteMesheryPattern(r, patternID)
	if err != nil {
		errPatternDelete := ErrDeletePattern(err)

		h.log.Error(errPatternDelete)
		http.Error(rw, errPatternDelete.Error(), http.StatusInternalServerError)
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": errPatternDelete,
		}).WithDescription("Error deleting pattern.").Build()
		http.Error(rw, errPatternDelete.Error(), http.StatusInternalServerError)
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		return
	}

	_ = json.Unmarshal(resp, &mesheryPattern)
	event := eventBuilder.WithSeverity(events.Informational).WithDescription(fmt.Sprintf("Pattern %s deleted.", mesheryPattern.Name)).Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)
	go h.config.PatternChannel.Publish(uuid.FromStringOrNil(user.ID), struct{}{})

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route GET /api/pattern/download/{id} PatternsAPI idGetMesheryPattern
// Handle GET request for Meshery Pattern with the given id
//
// ?oci={true|false} - If true, returns the pattern in OCI Artifact format
// ?export={Kubernetes Manifest} - exports the pattern file in the specified design format
// ?pkg={true|false} - If true, returns the artifact hub pkg and pattern file in zip file. If "oci" is true, "pkg" is ignored and the export always contains the artifact hub pkg.
//
// Get the pattern with the given id
// responses:
//  200:

// GetMesheryPatternHandler returns the pattern file with the given id

func (h *Handler) DownloadMesheryPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	var formatConverter converter.ConvertFormat
	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction("download").ActedUpon(userID).WithSeverity(events.Informational)

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
			http.Error(rw, err.Error(), http.StatusBadRequest)
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
		http.Error(rw, ErrGetPattern(err).Error(), http.StatusNotFound)
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrGetPattern(err),
		}).WithDescription(fmt.Sprintf("Failed to fetch design file for ID: %s.", patternID)).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		return
	}
	pattern := &models.MesheryPattern{}
	err = json.Unmarshal(resp, &pattern)
	if err != nil {
		obj := "download pattern"
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(rw, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": models.ErrUnmarshal(err, obj),
		}).WithDescription(fmt.Sprintf("Failed to unmarshal design file for ID: %s.", patternID)).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		return
	}

	err = h.VerifyAndConvertToDesign(r.Context(), pattern, provider)
	if err != nil {
		event := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("convert").WithDescription(fmt.Sprintf("The \"%s\" is not in the design format, failed to convert and persist the original source content from \"%s\" to design file format", pattern.Name, pattern.Type.String)).WithMetadata(map[string]interface{}{"error": err}).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// v1beta1
	isOldFormat, err := patternutils.IsDesignInAlpha2Format(pattern.PatternFile)
	if err != nil {
		err = ErrPatternFile(err)
		event := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("download").WithDescription(fmt.Sprintf("Failed to parse design \"%s\".", pattern.Name)).WithMetadata(map[string]interface{}{"error": err, "id": pattern.ID}).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if isOldFormat {

		eventBuilder := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("convert")
		_, patternFileStr, err := h.convertV1alpha2ToV1beta1(pattern, eventBuilder)
		event := eventBuilder.Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		if err != nil {
			h.log.Error(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		pattern.PatternFile = patternFileStr
	}

	if formatConverter != nil {
		patternFile, err := formatConverter.Convert(pattern.PatternFile)
		if err != nil {
			err = ErrExportPatternInFormat(err, exportFormat, pattern.Name)
			h.log.Error(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if exportFormat == string(coreV1.HelmChart) {
			rw.Header().Set("Content-Type", "application/x-gzip")
			rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.tgz", pattern.Name))
		} else {
			rw.Header().Set("Content-Type", "application/yaml")
			rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.yml", pattern.Name))
		}
		_, err = fmt.Fprint(rw, patternFile)
		if err != nil {
			err = ErrWriteResponse(err)
			h.log.Error(err)
			http.Error(rw, _errors.Wrapf(err, "failed to export design \"%s\" in %s format", pattern.Name, exportFormat).Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if ociFormat {
		tmpDir, err := oci.CreateTempOCIContentDir()
		if err != nil {
			h.log.Error(ErrCreateDir(err, "OCI"))
			http.Error(rw, ErrCreateDir(err, "OCI").Error(), http.StatusInternalServerError)

			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrCreateDir(err, "OCI"),
			}).WithDescription("Error creating tmp directory under ~/.meshery/content/").Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}
		defer os.RemoveAll(tmpDir)

		tmpDesignFile := filepath.Join(tmpDir, "design.yml")
		file, err := os.Create(tmpDesignFile)
		if err != nil {
			h.log.Error(ErrCreateFile(err, tmpDesignFile))
			http.Error(rw, ErrCreateFile(err, tmpDesignFile).Error(), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrCreateFile(err, tmpDesignFile),
			}).WithDescription(fmt.Sprintf("Error creating tmp file %s", tmpDesignFile)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}
		defer file.Close()

		var design patternV1beta1.PatternFile

		err = encoding.Unmarshal([]byte(pattern.PatternFile), &design)

		if err != nil {

			err = ErrEncodePattern(err)
			h.log.Error(err)
			http.Error(rw, fmt.Sprintf("Failed to export design \"%s\" as OCI image.", pattern.Name), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": err,
			}).WithDescription(fmt.Sprintf("Failed to export design \"%s\" as OCI image.", pattern.Name)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}

		ymlDesign, err := yaml.Marshal(design)

		if err != nil {
			err = ErrEncodePattern(err)
			h.log.Error(err)
			http.Error(rw, fmt.Sprintf("Failed to export design \"%s\" as OCI image.", pattern.Name), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": err,
			}).WithDescription(fmt.Sprintf("Failed to export design \"%s\" as OCI image.", pattern.Name)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}

		if _, err := file.Write(ymlDesign); err != nil {
			h.log.Error(ErrWritingIntoFile(err, tmpDesignFile))
			http.Error(rw, ErrWritingIntoFile(err, tmpDesignFile).Error(), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrWritingIntoFile(err, tmpDesignFile),
			}).WithDescription(fmt.Sprintf("Error writing into tmp design file %s at %s", pattern.Name, tmpDesignFile)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}

		artifactHubPkgFilePath := filepath.Join(tmpDir, "artifacthub-pkg.yml")
		artifactHubPkgFile, err := os.Create(artifactHubPkgFilePath)
		if err != nil {
			h.log.Error(ErrCreateFile(err, "artifacthub-pkg.yml"))
			eb := *eventBuilder
			event := eb.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to create artifacthub pkg for the design \"%s\"", pattern.Name)).WithMetadata(map[string]interface{}{"error": err}).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)
		}

		data, err := createArtifactHubPkg(pattern, strings.Trim(fmt.Sprintf("%s %s", user.FirstName, user.LastName), " "))
		if err != nil {
			h.log.Error(err)
			eb := *eventBuilder
			event := eb.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to create artifacthub pkg for the design \"%s\"", pattern.Name)).WithMetadata(map[string]interface{}{"error": err}).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)
		}

		_, err = artifactHubPkgFile.Write(data)
		if err != nil {
			err = ErrWritingIntoFile(err, "artifacthub-pkg.yml")
			h.log.Error(err)
			eb := *eventBuilder
			event := eb.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to create artifacthub pkg for the design \"%s\"", pattern.Name)).WithMetadata(map[string]interface{}{"error": err}).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)
		}

		ociImg, err := oci.BuildImage(tmpDir)
		if err != nil {
			h.log.Error(ErrBuildOCIImg(err))
			http.Error(rw, ErrBuildOCIImg(err).Error(), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrBuildOCIImg(err),
			}).WithDescription(fmt.Sprintf("Error building OCI Image from contents at %s", tmpDesignFile)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}

		digest, err := ociImg.Digest()
		if err != nil {
			h.log.Error(ErrBuildOCIImg(err))
			http.Error(rw, ErrBuildOCIImg(err).Error(), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrBuildOCIImg(err),
			}).WithDescription(fmt.Sprintf("Error getting image digest for %s", tmpDesignFile)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}

		size, err := ociImg.Size()
		if err != nil {
			h.log.Error(ErrBuildOCIImg(err))
			http.Error(rw, ErrBuildOCIImg(err).Error(), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrBuildOCIImg(err),
			}).WithDescription(fmt.Sprintf("Error getting calculating image size for %s", tmpDesignFile)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}

		h.log.Info(fmt.Sprintf("OCI Image built. Digest: %v, Size: %v", digest, size))

		eventBuilder.WithSeverity(events.Informational).WithDescription(fmt.Sprintf("OCI Image built. Digest: %v, Size: %v", digest, size))
		event := eventBuilder.Build()
		go h.config.EventBroadcaster.Publish(userID, event)
		_ = provider.PersistEvent(event)

		pretifiedName := strings.ToLower(strings.Replace(pattern.Name, " ", "", -1)) // ensures that tag validation passes
		tmpOCITarFilePath := filepath.Join(tmpDir, pretifiedName+".tar")
		err = oci.SaveOCIArtifact(ociImg, tmpOCITarFilePath, pretifiedName)
		if err != nil {
			h.log.Error(ErrSaveOCIArtifact(err))
			http.Error(rw, ErrSaveOCIArtifact(err).Error(), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrSaveOCIArtifact(err),
			}).WithDescription(fmt.Sprintf("Failed to save OCI Artifact %s temporarily", tmpOCITarFilePath)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}

		file, err = os.OpenFile(tmpOCITarFilePath, os.O_RDONLY, 0444)
		if err != nil {
			h.log.Error(ErrOpenFile(tmpOCITarFilePath))
			http.Error(rw, ErrOpenFile(tmpOCITarFilePath).Error(), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrOpenFile(tmpOCITarFilePath),
			}).WithDescription(fmt.Sprintf("Failed to read contents of OCI Artifact %s", tmpOCITarFilePath)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}
		content, err := io.ReadAll(file)
		if err != nil {
			h.log.Error(ErrIOReader(err))
			http.Error(rw, ErrIOReader(err).Error(), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrIOReader(err),
			}).WithDescription(fmt.Sprintf("Failed to read contents of OCI artifact %s", tmpOCITarFilePath)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}

		h.log.Info("OCI Artifact saved at: ", tmpOCITarFilePath)

		eventBuilder.WithSeverity(events.Informational).WithDescription(fmt.Sprintf("OCI Artifact temporarily saved at: %s", tmpOCITarFilePath))
		event = eventBuilder.Build()
		go h.config.EventBroadcaster.Publish(userID, event)
		_ = provider.PersistEvent(event)

		rw.Header().Set("Content-Type", "application/tar")
		rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.tar", pattern.Name))

		reader := bytes.NewReader(content)
		if _, err := io.Copy(rw, reader); err != nil {
			http.Error(rw, ErrIOReader(err).Error(), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrIOReader(err),
			}).WithDescription("Failed to save contents of OCI Artifact at requested path").Build()
			_ = provider.PersistEvent(event)
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
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)
		}
		ymlDesign, err := yaml.Marshal(pattern.PatternFile)
		if err != nil {
			err = ErrEncodePattern(err)
			h.log.Error(err)
			http.Error(rw, fmt.Sprintf("Failed to export design \"%s\" as OCI image.", pattern.Name), http.StatusInternalServerError)
			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": err,
			}).WithDescription(fmt.Sprintf("Failed to export design \"%s\" as OCI image.", pattern.Name)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

			return
		}
		err = tarWriter.Compress(pattern.Name+".yml", ymlDesign)
		if err != nil {
			h.log.Error(err)
			eb := *eventBuilder
			event := eb.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to compress design \"%s\" and artifacthub pkg.", pattern.Name)).WithMetadata(map[string]interface{}{"error": err}).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			tarWriter.Close()
			return
		}

		tarWriter.Close()

		_, err = io.Copy(rw, tarWriter.Buffer)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	yamlBytes, err := encoding.ToYaml([]byte(pattern.PatternFile))
	if err != nil {
		err = ErrParsePattern(err)
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/yaml")
	rw.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.yml", pattern.Name))

	_, err = rw.Write(yamlBytes)
	if err != nil {
		err = ErrEncodePattern(err)
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route POST /api/pattern/clone/{id} PatternsAPI idCloneMesheryPattern
// Handle Clone for a Meshery Pattern
//
// Creates a local copy of a published pattern with id: id
// responses:
//
//	200 : noContentWrapper
//
// CloneMesheryPatternHandler clones a pattern with the given id
func (h *Handler) CloneMesheryPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	patternID := mux.Vars(r)["id"]
	patternUUID := uuid.FromStringOrNil(patternID)

	userID := uuid.FromStringOrNil(user.ID)
	token, _ := r.Context().Value(models.TokenCtxKey).(string)

	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction("clone").ActedUpon(patternUUID).WithSeverity(events.Informational)

	var parsedBody *models.MesheryClonePatternRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil || patternID == "" {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
		return
	}

	defer func() {
		_ = r.Body.Close()
	}()

	mesheryPatternByt, err := provider.GetMesheryPattern(r, patternID, "false")
	if err != nil {
		h.log.Error(ErrGetPattern(err))
		http.Error(rw, ErrGetPattern(err).Error(), http.StatusNotFound)
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrGetPattern(err),
		}).WithDescription(fmt.Sprintf("Failed to fetch meshery pattern \"%s\" with id: %s.", parsedBody.Name, patternID)).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		return
	}

	pattern := &models.MesheryPattern{}
	err = json.Unmarshal(mesheryPatternByt, &pattern)
	if err != nil {
		obj := "pattern: " + patternID
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(rw, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": models.ErrUnmarshal(err, obj),
		}).WithDescription(fmt.Sprintf("Failed to fetch meshery pattern \"%s\" with ID: %s.", parsedBody.Name, patternID)).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		return
	}

	isOldFormat, err := patternutils.IsDesignInAlpha2Format(pattern.PatternFile)
	if err != nil {
		err = ErrPatternFile(err)
		event := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithDescription(fmt.Sprintf("Failed to parse design \"%s\".", pattern.Name)).WithMetadata(map[string]interface{}{"error": err, "id": patternID}).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if isOldFormat {
		eventBuilder := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("convert")
		_, patternFileStr, err := h.convertV1alpha2ToV1beta1(pattern, eventBuilder)
		event := eventBuilder.Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		if err != nil {
			h.log.Error(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		pattern.PatternFile = patternFileStr
		_, err = provider.SaveMesheryPattern(token, pattern)
		if err != nil {
			h.log.Error(ErrSavePattern(err))
			http.Error(rw, ErrSavePattern(err).Error(), http.StatusInternalServerError)

			event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
				"error": ErrSavePattern(_errors.Wrapf(err, "failed to persist converted v1beta1 design file \"%s\" with id: %s", parsedBody.Name, patternID)),
			}).WithDescription(ErrSavePattern(err).Error()).Build()

			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)
			return
		}
	}

	resp, err := provider.CloneMesheryPattern(r, patternID, parsedBody)
	if err != nil {
		h.log.Error(ErrClonePattern(err))
		http.Error(rw, ErrClonePattern(err).Error(), http.StatusInternalServerError)
		return
	}
	go h.config.PatternChannel.Publish(uuid.FromStringOrNil(user.ID), struct{}{})
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route POST /api/pattern/catalog/publish PatternsAPI idPublishCatalogPatternHandler
// Handle Publish for a Meshery Pattern
//
// Publishes pattern to Meshery Catalog by setting visibility to published and setting catalog data
// responses:
//
//	202: noContentWrapper
//
// PublishCatalogPatternHandler sets visibility of pattern with given id as published
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

	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().
		FromUser(userID).
		FromSystem(*h.SystemID).
		WithCategory("pattern").
		WithAction("publish").
		ActedUpon(userID)

	var parsedBody *models.MesheryCatalogPatternRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRequestBody(err))
		e := eventBuilder.WithSeverity(events.Error).
			WithMetadata(map[string]interface{}{
				"error": ErrRequestBody(err),
			}).
			WithDescription("Error parsing design payload.").Build()
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
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
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrPublishCatalogPattern(err).Error(), http.StatusInternalServerError)
		return
	}

	var respBody *models.CatalogRequest
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		h.log.Error(ErrPublishCatalogPattern(err))
		e := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrPublishCatalogPattern(err),
		}).WithDescription("Error parsing response.").Build()
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrPublishCatalogPattern(err).Error(), http.StatusInternalServerError)
	}

	e := eventBuilder.WithSeverity(events.Informational).ActedUpon(parsedBody.ID).WithDescription(fmt.Sprintf("Request to publish '%s' design submitted with status: %s", respBody.ContentName, respBody.Status)).Build()
	_ = provider.PersistEvent(e)
	go h.config.EventBroadcaster.Publish(userID, e)

	go h.config.PatternChannel.Publish(uuid.FromStringOrNil(user.ID), struct{}{})
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusAccepted)
	fmt.Fprint(rw, string(resp))
}

// swagger:route DELETE /api/pattern/catalog/unpublish PatternsAPI idUnPublishCatalogPatternHandler
// Handle Publish for a Meshery Pattern
//
// Unpublishes pattern from Meshery Catalog by setting visibility to private and removing catalog data from website
// responses:
//
//	200: noContentWrapper
//
// UnPublishCatalogPatternHandler sets visibility of pattern with given id as private
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

	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().
		FromUser(userID).
		FromSystem(*h.SystemID).
		WithCategory("pattern").
		WithAction("unpublish_request").
		ActedUpon(userID)

	var parsedBody *models.MesheryCatalogPatternRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRequestBody(err))
		e := eventBuilder.WithSeverity(events.Error).
			WithMetadata(map[string]interface{}{
				"error": ErrRequestBody(err),
			}).
			WithDescription("Error parsing design payload.").Build()
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
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
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrPublishCatalogPattern(err).Error(), http.StatusInternalServerError)
		return
	}

	var respBody *models.CatalogRequest
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		h.log.Error(ErrPublishCatalogPattern(err))
		e := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrPublishCatalogPattern(err),
		}).WithDescription("Error parsing response.").Build()
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrPublishCatalogPattern(err).Error(), http.StatusInternalServerError)
	}

	e := eventBuilder.WithSeverity(events.Informational).ActedUpon(parsedBody.ID).WithDescription(fmt.Sprintf("'%s' design unpublished", respBody.ContentName)).Build()
	_ = provider.PersistEvent(e)
	go h.config.EventBroadcaster.Publish(userID, e)

	go h.config.PatternChannel.Publish(uuid.FromStringOrNil(user.ID), struct{}{})
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route DELETE /api/patterns PatternsAPI idDeleteMesheryPattern
// Handle Delete for multiple Meshery Patterns
//
// DeleteMultiMesheryPatternsHandler deletes patterns with the given ids
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
		http.Error(rw, fmt.Sprintf("failed to delete the pattern: %s", err), http.StatusInternalServerError)
		return
	}
	go h.config.PatternChannel.Publish(uuid.FromStringOrNil(user.ID), struct{}{})
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route GET /api/pattern/{id} PatternsAPI idGetMesheryPattern
// Handle GET for a Meshery Pattern
//
// ```?metrics``` Returns metrics like deployment/share/clone/view/download count for desings, default false,
//
// Fetches the pattern with the given id
// responses:
// 	200: mesheryPatternResponseWrapper

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
	userID := uuid.FromStringOrNil(user.ID)

	resp, err := provider.GetMesheryPattern(r, patternID, r.URL.Query().Get("metrics"))
	if err != nil {
		h.log.Error(ErrGetPattern(err))
		http.Error(rw, ErrGetPattern(err).Error(), http.StatusNotFound)
		return
	}

	pattern := &models.MesheryPattern{}
	err = json.Unmarshal(resp, &pattern)
	if err != nil {
		h.log.Error(ErrGetPattern(err))
		http.Error(rw, ErrGetPattern(err).Error(), http.StatusInternalServerError)
		return
	}

	err = h.VerifyAndConvertToDesign(r.Context(), pattern, provider)
	if err != nil {
		event := events.NewEvent().ActedUpon(patternUUID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("convert").WithDescription(fmt.Sprintf("The \"%s\" is not in the design format, failed to convert and persist the original source content from \"%s\" to design file format", pattern.Name, pattern.Type.String)).WithMetadata(map[string]interface{}{"error": err}).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	isOldFormat, err := patternutils.IsDesignInAlpha2Format(pattern.PatternFile)
	if err != nil {
		err = ErrPatternFile(err)
		event := events.NewEvent().ActedUpon(patternUUID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("view").WithDescription(fmt.Sprintf("Failed to parse design \"%s\".", pattern.Name)).WithMetadata(map[string]interface{}{"error": err, "id": pattern.ID}).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if isOldFormat {
		eventBuilder := events.NewEvent().ActedUpon(*pattern.ID).FromSystem(*h.SystemID).FromUser(userID).WithCategory("pattern").WithAction("convert")
		_, patternFileStr, err := h.convertV1alpha2ToV1beta1(pattern, eventBuilder)
		event := eventBuilder.Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		if err != nil {
			h.log.Error(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		pattern.PatternFile = patternFileStr
	}

	// deprettify pattern for backward compatibility with older designs which had the configuration in prettified format
	var design patternV1beta1.PatternFile
	err = encoding.Unmarshal([]byte(pattern.PatternFile), &design)

	if err != nil {
		err = ErrParsePattern(err)
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, component := range design.Components {
		component.Configuration = core.Format.DePrettify(component.Configuration, false)
	}

	patternBytes, err := encoding.Marshal(design)
	pattern.PatternFile = string(patternBytes)
	// done deprettifying

	if err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// done deprettifying

	rw.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(pattern); err != nil {
		http.Error(rw, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) formatPatternOutput(rw http.ResponseWriter, content []byte, format, sourcetype string, eventBuilder *events.EventBuilder, URL, action string) {
	contentMesheryPatternSlice := make([]models.MesheryPattern, 0)

	if err := json.Unmarshal(content, &contentMesheryPatternSlice); err != nil {
		http.Error(rw, ErrDecodePattern(err).Error(), http.StatusInternalServerError)
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
		http.Error(rw, models.ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
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
	fmt.Fprint(rw, string(data))
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

// swagger:route PUT /api/pattern/{sourcetype} PatternsAPI idUpdateMesheryPattern
// Handle PUT request for Meshery Pattern with the given payload
//
// Updates the pattern with the given payload
// responses:
//
//	200: mesheryPatternResponseWrapper
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
	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction("update").ActedUpon(userID)

	res := meshes.EventsResponse{
		Component:     "core",
		ComponentName: "Design",
		OperationId:   guid.NewString(),
		EventType:     meshes.EventType_INFO,
	}

	sourcetype := mux.Vars(r)["sourcetype"]
	if sourcetype == "" {
		http.Error(rw, "missing route variable \"source-type\"", http.StatusBadRequest)

		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrSaveApplication(fmt.Errorf("missing route variable \"source-type\" (one of %s, %s, %s)", models.K8sManifest, models.DockerCompose, models.HelmChart)),
		}).WithDescription("Please provide design source-type").Build()

		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		go h.EventsBuffer.Publish(&res)
		return
	}

	var parsedBody *MesheryPatternUPDATERequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		http.Error(rw, ErrRetrieveData(err).Error(), http.StatusBadRequest)
		return
	}

	actedUpon := &userID
	if parsedBody.PatternData != nil && parsedBody.PatternData.ID != nil {
		actedUpon = parsedBody.PatternData.ID
	}

	eventBuilder.ActedUpon(*actedUpon)

	token, err := provider.GetProviderToken(r)
	if err != nil {
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()

		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
		return
	}
	format := r.URL.Query().Get("output")

	mesheryPattern := parsedBody.PatternData
	mesheryPattern.Type = sql.NullString{
		String: sourcetype,
		Valid:  true,
	}
	resp, err := provider.SaveMesheryPattern(token, mesheryPattern)
	if err != nil {
		errAppSave := ErrSaveApplication(err)
		h.log.Error(errAppSave)

		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "%s", err)

		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": errAppSave,
		}).WithDescription(fmt.Sprintf("Error saving design %s", parsedBody.PatternData.Name)).Build()

		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		return
	}
	go h.config.PatternChannel.Publish(userID, struct{}{})

	eventBuilder.WithSeverity(events.Informational)
	h.formatPatternOutput(rw, resp, format, sourcetype, eventBuilder, parsedBody.URL, models.Update)
	event := eventBuilder.Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)

}

// swagger:route POST /api/pattern/{sourcetype} PatternsAPI idPostPatternFileRequest
// Handle POST request for Pattern Files
//
// Creates a new Pattern with source-content
// responses:
//  200: mesheryPatternResponseWrapper

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

// swagger:route GET /api/pattern/types PatternsAPI typeGetMesheryPatternTypesHandler
// Handle GET request for Meshery Pattern types
//
// Get pattern file types
// responses:
//
//	200: mesheryApplicationTypesResponseWrapper
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
		http.Error(rw, models.ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(b))
}

// swagger:route GET /api/pattern/download/{id}/{sourcetype} PatternsAPI typeGetPatternSourceContent
// Handle GET request for Meshery Patterns with of source content
//
// Get the pattern source-content
// responses:
//  200: mesheryPatternSourceContentResponseWrapper

// GetMesheryPatternHandler fetched the design using the given id and sourcetype
func (h *Handler) GetMesheryPatternSourceHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	designID := mux.Vars(r)["id"]
	token, _ := r.Context().Value(models.TokenCtxKey).(string)

	resp, err := provider.GetDesignSourceContent(token, designID)
	if err != nil {
		h.log.Error(ErrGetPattern(err))
		http.Error(rw, ErrGetPattern(err).Error(), http.StatusNotFound)
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
		h.log.Error(ErrApplicationSourceContent(err, "download"))
		http.Error(rw, ErrApplicationSourceContent(err, "download").Error(), http.StatusInternalServerError)
	}
}

func createArtifactHubPkg(pattern *models.MesheryPattern, user string) ([]byte, error) {
	isCatalogItem := pattern.Visibility == models.Published
	var version string
	if isCatalogItem {
		version = pattern.CatalogData.PublishedVersion
	}
	artifactHubPkg := catalog.BuildArtifactHubPkg(pattern.Name, "", user, version, pattern.CreatedAt.Format(time.RFC3339), &pattern.CatalogData)

	data, err := yaml.Marshal(artifactHubPkg)
	if err != nil {
		return nil, models.ErrMarshalYAML(err, "artifacthub-pkg")
	}

	return data, nil
}

func (h *Handler) convertV1alpha2ToV1beta1(mesheryPattern *models.MesheryPattern, eventBuilder *events.EventBuilder) (*pattern.PatternFile, string, error) {

	v1alpha1PatternFile := v1alpha2.PatternFile{}

	v1beta1PatternFile := pattern.PatternFile{}

	err := encoding.Unmarshal([]byte(mesheryPattern.PatternFile), &v1alpha1PatternFile)
	if err != nil {
		return nil, "", ErrParsePattern(err)
	}

	for _, svc := range v1alpha1PatternFile.Services {
		svc.Settings = helpers.RecursiveCastMapStringInterfaceToMapStringInterface(svc.Settings)
		svc.Traits = helpers.RecursiveCastMapStringInterfaceToMapStringInterface(svc.Traits)
	}

	err = v1beta1PatternFile.ConvertFrom(&v1alpha1PatternFile)
	if err != nil {
		return nil, "", err
	}

	v1beta1PatternFile.Id = *mesheryPattern.ID
	v1beta1PatternFile.Version = v1alpha1PatternFile.Version

	h.log.Infof("Converted design file with id \"%s\" to v1beta1 format", *mesheryPattern.ID)

	err = mapModelRelatedData(h.registryManager, &v1beta1PatternFile)
	if err != nil {
		eventBuilder.WithDescription("Design converted to v1beta1 format but failed to assign styles and metadata").
			WithMetadata(map[string]interface{}{"error": ErrGetComponentDefinition(err), "id": *mesheryPattern.ID}).WithSeverity(events.Warning)
		return nil, "", err
	}

	v1beta1PatternByt, err := encoding.Marshal(v1beta1PatternFile)
	if err != nil {
		return nil, "", utils.ErrMarshal(err)
	}
	eventBuilder.WithSeverity(events.Informational).WithDescription(fmt.Sprintf("Converted design file \"%s\" with id \"%s\" to v1beta1 format", mesheryPattern.Name, *mesheryPattern.ID))
	return &v1beta1PatternFile, string(v1beta1PatternByt), nil
}

func mapModelRelatedData(reg *meshmodel.RegistryManager, patternFile *pattern.PatternFile) error {
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
			comp.Model.Registrant.Status = connection.Registered
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

		var styles component.Styles

		if comp.Styles != nil {
			styles = *comp.Styles
		} else {
			comp.Styles = &component.Styles{}
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
