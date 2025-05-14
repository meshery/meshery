package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
	pCore "github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshkit/encoding"
	"github.com/layer5io/meshkit/files"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/utils"
	coreV1 "github.com/meshery/schemas/models/v1alpha1/core"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

type MesheryDesignImportPayload struct {
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	File     []byte `json:"file,omitempty"`
	FileName string `json:"file_name,omitempty"`
}

type FileToImport struct {
	Data     []byte
	FileName string
}

func GetFileToImportFromPayload(payload MesheryDesignImportPayload) (FileToImport, error) {
	if payload.URL != "" {
		resp, err := http.Get(payload.URL)
		if err != nil {
			return FileToImport{}, models.ErrDoRequest(err, "GET", payload.URL)
		}
		defer func() {
			if resp != nil {
				_ = resp.Body.Close()
			}
		}()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return FileToImport{}, err
		}

		// Get filename from Content-Disposition header or URL
		filename := getFileNameFromResponse(resp, payload.URL)

		return FileToImport{
			Data:     body,
			FileName: filename,
		}, nil
	}

	return FileToImport{
		Data:     payload.File,
		FileName: payload.FileName,
	}, nil
}

func getFileNameFromResponse(resp *http.Response, fileURL string) string {
	// Check Content-Disposition header
	contentDisp := resp.Header.Get("Content-Disposition")
	if contentDisp != "" {
		parts := strings.Split(contentDisp, "filename=")
		if len(parts) > 1 {
			filename := strings.Trim(parts[1], `"`)
			return filename
		}
	}

	// Fallback: Extract filename from URL
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "unknown_file"
	}
	return path.Base(parsedURL.Path)
}

func ConvertFileToManifest(identifiedFile files.IdentifiedFile, rawFile FileToImport) (string, error) {

	switch identifiedFile.Type {

	case coreV1.HelmChart:
		return files.ConvertHelmChartToKubernetesManifest(identifiedFile)
	case coreV1.DockerCompose:
		return files.ConvertDockerComposeToKubernetesManifest(identifiedFile)
	case coreV1.K8sManifest:
		return string(rawFile.Data), nil
	case coreV1.K8sKustomize:
		return files.ConvertKustomizeToKubernetesManifest(identifiedFile)
	default:
		return "", files.ErrUnsupportedFileTypeForConversionToDesign(rawFile.FileName, string(identifiedFile.Type))
	}
}

// returns the design file , the type of file that was identified during converion , and any error
func ConvertFileToDesign(fileToImport FileToImport, registry *registry.RegistryManager, logger logger.Handler) (pattern.PatternFile, coreV1.IaCFileTypes, error) {

	defer utils.TrackTime(logger, time.Now(), "ConvertFileToDesign")

	var emptyDesign pattern.PatternFile

	validImportExtensions := map[string]bool{
		".yml":    true,
		".yaml":   true,
		".json":   true,
		".tar.gz": true,
		".tgz":    true,
		".zip":    true,
		".tar":    true,
	}

	tempDir, err := os.MkdirTemp("", "temp-import")

	defer os.RemoveAll(tempDir)

	if err != nil {
		return emptyDesign, "", fmt.Errorf("Failed to create tmp directory %w", err)
	}

	now := time.Now()
	// NOTE: the FileName must also have extension
	sanitizedFile, err := files.SanitizeFile(fileToImport.Data, fileToImport.FileName, tempDir, validImportExtensions)
	utils.TrackTime(logger, now, "SanitizeFile")

	if err != nil {
		return emptyDesign, "", err
	}

	now = time.Now()
	identifiedFile, err := files.IdentifyFile(sanitizedFile)
	utils.TrackTime(logger, now, "IdentifyFile")

	if err != nil {
		return emptyDesign, "", err
	}

	if identifiedFile.Type == coreV1.MesheryDesign {
		design := identifiedFile.ParsedFile.(pattern.PatternFile)
		return design, identifiedFile.Type, nil
	}

	now = time.Now()
	manifest, err := ConvertFileToManifest(identifiedFile, fileToImport)
	utils.TrackTime(logger, now, "ConvertFileToManifest")

	if err != nil {
		return emptyDesign, "", err
	}

	now = time.Now()
	design, err := pCore.NewPatternFileFromK8sManifest(manifest, fileToImport.FileName, true, registry)
	utils.TrackTime(logger, now, "ConvertManifestToDesign")

	if err != nil {
		return emptyDesign, "", err
	}

	return design, identifiedFile.Type, err
}

func (h *Handler) logErrorGettingUserToken(rw http.ResponseWriter, provider models.Provider, err error, userID uuid.UUID, eventBuilder *events.EventBuilder) {

	h.log.Error(ErrRetrieveUserToken(err))
	http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)

	if eventBuilder != nil {
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()

		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
	}

}

func (h *Handler) logErrorParsingRequestBody(rw http.ResponseWriter, provider models.Provider, err error, userID uuid.UUID, eventBuilder *events.EventBuilder) {

	h.log.Error(ErrRequestBody(err))
	http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)

	if eventBuilder != nil {
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrRequestBody(err),
		}).WithDescription("Unable to parse request body").Build()

		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
	}
}

func ImportErrorEvent(eventBuilder events.EventBuilder, importPayload MesheryDesignImportPayload, err error) *events.Event {

	source := importPayload.URL
	if source == "" {
		source = importPayload.FileName
	}
	event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
		"error": err,
	}).WithDescription(fmt.Sprintf("Failed to import design '%s' from %s", importPayload.Name, source)).Build()

	return event

}

// swagger:route POST /api/pattern/import PatternsAPI idPostPatternFile
// Handle design import
//
// responses: [Design]
//
//	200: mesheryPatternResponseWrapper
func (h *Handler) DesignFileImportHandler(
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
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction("create").ActedUpon(userID).WithSeverity(events.Informational)

	var importDesignPayload MesheryDesignImportPayload

	if err := json.NewDecoder(r.Body).Decode(&importDesignPayload); err != nil {
		h.logErrorParsingRequestBody(rw, provider, err, userID, eventBuilder)
		return
	}
	actedUpon := &userID
	eventBuilder.ActedUpon(*actedUpon)

	token, err := provider.GetProviderToken(r)
	if err != nil || token == " " {
		h.logErrorGettingUserToken(rw, provider, err, userID, eventBuilder)
		return
	}

	fileToImport, err := GetFileToImportFromPayload(importDesignPayload)

	if err != nil {
		h.log.Error(fmt.Errorf("Conversion: Failed to get file from payload  %w", err))
		http.Error(rw, err.Error(), http.StatusBadRequest)
		event := ImportErrorEvent(*eventBuilder, importDesignPayload, err)
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		return
	}

	design, sourceFileType, err := ConvertFileToDesign(fileToImport, h.registryManager, h.log)

	if err != nil {
		h.log.Error(fmt.Errorf("Conversion: Failed to convert to design %w", err))
		http.Error(rw, err.Error(), http.StatusBadRequest)
		event := ImportErrorEvent(*eventBuilder, importDesignPayload, err)
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		return
	}

	design.Name = importDesignPayload.Name
	patternFile, err := encoding.Marshal(design)

	if err != nil {
		h.log.Error(ErrSavePattern(err))
		http.Error(rw, ErrSavePattern(err).Error(), http.StatusInternalServerError)

		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrSavePattern(err),
		}).WithDescription(ErrSavePattern(err).Error()).Build()

		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		return
	}

	// pattern to be saved in the database
	designRecord := models.MesheryPattern{
		Name:          design.Name,
		PatternFile:   string(patternFile),
		SourceContent: fileToImport.Data,
		Type: sql.NullString{
			String: string(sourceFileType),
			Valid:  true,
		},
	}

	savedDesignByt, err := provider.SaveMesheryPattern(token, &designRecord)

	if err != nil {

		h.handleProviderPatternSaveError(rw, eventBuilder, userID, savedDesignByt, err, provider)
		return
	}

	_, _ = rw.Write(savedDesignByt)

	event := eventBuilder.WithSeverity(events.Success).WithDescription(fmt.Sprintf("Imported design '%s' of type '%s'", design.Name, sourceFileType)).Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)

}
