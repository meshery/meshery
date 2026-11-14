package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	serverutils "github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models"
	pCore "github.com/meshery/meshery/server/models/pattern/core"
	"github.com/meshery/meshkit/encoding"
	"github.com/meshery/meshkit/files"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/core"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

// FileToImport is the internal tuple of bytes + filename that the
// import pipeline works with after the request-body variant has been
// resolved. It is not part of the wire contract — see the schemas
// repo's MesheryPatternImportRequestBody for that.
type FileToImport struct {
	Data     []byte
	FileName string
}

// importVariant captures the content extracted from one arm of the
// MesheryPatternImportRequestBody oneOf. Keeping the display name
// alongside FileToImport lets callers derive event descriptions without
// re-inspecting the raw payload.
type importVariant struct {
	Name string
	File FileToImport
}

// designImportHTTPClient is the shared client used for URL-variant
// imports. It enforces a bounded 60s timeout and uses a safe transport
// that rejects private, loopback, and cloud metadata addresses (SSRF protection).
var designImportHTTPClient = serverutils.NewSafeHTTPClient(60 * time.Second)

// resolveImportVariant decodes the request body against the two oneOf
// variants published in the schema, enforces the "exactly one of"
// invariant, and — for the URL variant — fetches the remote file
// before returning. On error, the returned importVariant carries any
// context we could extract (the caller's supplied Name plus, for the
// URL variant, the URL echoed in FileName) so that the downstream
// ImportErrorEvent does not emit `Failed to import design '' from ''`
// for unhappy-path requests.
func resolveImportVariant(body pattern.MesheryPatternImportRequestBody) (importVariant, error) {
	filePayload, fileErr := body.AsMesheryPatternImportFilePayload()
	urlPayload, urlErr := body.AsMesheryPatternImportURLPayload()

	hasFile := fileErr == nil && len(filePayload.File) > 0 && filePayload.FileName != ""
	hasURL := urlErr == nil && urlPayload.Url != ""

	switch {
	case hasFile && hasURL:
		// Echo the filename so the error event can name what was
		// rejected; both variants look plausible to the decoder.
		return importVariant{
				Name: stringFromPtr(filePayload.Name),
				File: FileToImport{FileName: filePayload.FileName},
			},
			ErrInvalidImportRequest(errors.New("request body must contain exactly one of File Import (file + fileName) or URL Import (url), not both"))
	case !hasFile && !hasURL:
		return importVariant{}, ErrInvalidImportRequest(errors.New("request body must contain either a File Import (file + fileName) or a URL Import (url)"))
	case hasFile:
		return importVariant{
			Name: stringFromPtr(filePayload.Name),
			File: FileToImport{Data: filePayload.File, FileName: filePayload.FileName},
		}, nil
	default:
		fetched, err := fetchFileFromURL(urlPayload.Url)
		if err != nil {
			// Preserve the URL as the "source" for the failure event.
			return importVariant{
				Name: stringFromPtr(urlPayload.Name),
				File: FileToImport{FileName: urlPayload.Url},
			}, err
		}
		return importVariant{Name: stringFromPtr(urlPayload.Name), File: fetched}, nil
	}
}

// fetchFileFromURL performs the remote GET behind a URL-variant import
// and derives a filename from either Content-Disposition or the path
// segment of the URL. Uses designImportHTTPClient so every fetch is
// bounded by a timeout and guarded against SSRF, and rejects non-2xx
// responses explicitly so a 404 HTML body isn't handed to the design
// parser as if it were a valid import file.
func fetchFileFromURL(fileURL string) (FileToImport, error) {
	if _, err := serverutils.ValidateSafeURL(fileURL); err != nil {
		return FileToImport{}, models.ErrDoRequest(err, "GET", fileURL)
	}

	resp, err := designImportHTTPClient.Get(fileURL)
	if err != nil {
		return FileToImport{}, models.ErrDoRequest(err, "GET", fileURL)
	}
	defer func() {
		if resp != nil {
			_ = resp.Body.Close()
		}
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Wrap so every URL-fetch failure path surfaces with MeshKit
		// metadata, matching the http.Client.Get error path above.
		return FileToImport{}, models.ErrDoRequest(fmt.Errorf("returned HTTP %d %s", resp.StatusCode, resp.Status), "GET", fileURL)
	}
	// Cap read at 10 MB to prevent unbounded memory consumption from large remote responses.
	// Reading up to maxImportFileSize+1 allows detecting whether the remote stream exceeds the cap,
	// rejecting truncated downloads instead of accepting partial files.
	const maxImportFileSize = 10 * 1024 * 1024
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(maxImportFileSize)+1))
	if err != nil {
		return FileToImport{}, models.ErrDoRequest(err, "GET", fileURL)
	}
	if len(body) > maxImportFileSize {
		return FileToImport{}, models.ErrDoRequest(fmt.Errorf("file size exceeds maximum allowed limit of %d bytes", maxImportFileSize), "GET", fileURL)
	}
	return FileToImport{Data: body, FileName: getFileNameFromResponse(resp, fileURL)}, nil
}

func stringFromPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func getFileNameFromResponse(resp *http.Response, fileURL string) string {
	// Content-Disposition can carry more than the name: other parameters, or
	// a filename* form next to filename (RFC 6266, section 4.3).
	// ParseMediaType returns the filename parameter on its own, decoding
	// filename* into it when present.
	if contentDisp := resp.Header.Get("Content-Disposition"); contentDisp != "" {
		if _, params, err := mime.ParseMediaType(contentDisp); err == nil && params["filename"] != "" {
			return params["filename"]
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

	case core.HelmChart:
		return files.ConvertHelmChartToKubernetesManifest(identifiedFile)
	case core.DockerCompose:
		return files.ConvertDockerComposeToKubernetesManifest(identifiedFile)
	case core.K8sManifest:
		return string(rawFile.Data), nil
	case core.K8sKustomize:
		return files.ConvertKustomizeToKubernetesManifest(identifiedFile)
	default:
		return "", files.ErrUnsupportedFileTypeForConversionToDesign(rawFile.FileName, string(identifiedFile.Type))
	}
}

// ConvertFileToDesign returns the design file, the type of file that was identified during conversion, and any error
func ConvertFileToDesign(fileToImport FileToImport, registry *registry.RegistryManager, logger logger.Handler) (pattern.PatternFile, core.IaCFileTypes, error) {

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
	if err != nil {
		return emptyDesign, "", fmt.Errorf("failed to create tmp directory %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			logger.Error(err)
		}
	}()

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

	if identifiedFile.Type == core.MesheryDesign {
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
	design, err := pCore.NewPatternFileFromK8sManifest(manifest, fileToImport.FileName, false, registry)
	utils.TrackTime(logger, now, "ConvertManifestToDesign")

	if err != nil {
		return emptyDesign, "", err
	}

	return design, identifiedFile.Type, err
}

func (h *Handler) logErrorGettingUserToken(rw http.ResponseWriter, provider models.Provider, err error, userID core.Uuid, eventBuilder *events.EventBuilder) {

	h.log.Error(ErrRetrieveUserToken(err))
	writeMeshkitError(rw, ErrRetrieveUserToken(err), http.StatusInternalServerError)

	if eventBuilder != nil {
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()

		h.persistAndPublish(provider, "", userID, event)
	}

}

func (h *Handler) logErrorParsingRequestBody(rw http.ResponseWriter, provider models.Provider, err error, userID core.Uuid, eventBuilder *events.EventBuilder) {

	h.log.Error(ErrRequestBody(err))
	writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)

	if eventBuilder != nil {
		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrRequestBody(err),
		}).WithDescription("Unable to parse request body").Build()

		h.persistAndPublish(provider, "", userID, event)
	}
}

func ImportErrorEvent(eventBuilder events.EventBuilder, variant importVariant, err error) *events.Event {

	source := variant.File.FileName
	event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
		"error": err,
	}).WithDescription(fmt.Sprintf("Failed to import design '%s' from %s", variant.Name, source)).Build()

	return event

}

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
	userID := user.ID
	eventBuilder := events.NewEvent().FromOwner(userID).FromSystem(*h.SystemID).WithCategory("pattern").WithAction("create").ActedUpon(userID).WithSeverity(events.Informational)

	var importBody pattern.MesheryPatternImportRequestBody

	if err := json.NewDecoder(r.Body).Decode(&importBody); err != nil {
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

	variant, err := resolveImportVariant(importBody)
	if err != nil {
		// resolveImportVariant failures are 400-class — either the body
		// violated the oneOf contract (already wrapped as
		// ErrInvalidImportRequest) or the URL variant couldn't be
		// fetched (already wrapped as models.ErrDoRequest). Either way
		// the caller needs to correct the request, not the server to
		// recover.
		h.log.Error(fmt.Errorf("resolve import variant: %w", err))
		writeMeshkitError(rw, err, http.StatusBadRequest)
		event := ImportErrorEvent(*eventBuilder, variant, err)
		h.persistAndPublish(provider, token, userID, event)
		return
	}

	fileToImport := variant.File

	design, sourceFileType, err := ConvertFileToDesign(fileToImport, h.registryManager, h.log)

	if err != nil {
		// ConvertFileToDesign returns a mix of bare and MeshKit-wrapped
		// errors depending on which import-pipeline stage failed
		// (sanitize / identify / convert manifest / build pattern).
		// Wrap so the JSON envelope returned to the client always
		// carries MeshKit metadata.
		wrappedErr := ErrConvertToDesign(err)
		h.log.Error(fmt.Errorf("conversion: failed to convert to design %w", err))
		writeMeshkitError(rw, wrappedErr, http.StatusBadRequest)
		event := ImportErrorEvent(*eventBuilder, variant, wrappedErr)
		h.persistAndPublish(provider, token, userID, event)
		return
	}

	// Only overwrite the design name when the caller supplied one in
	// the request — leaving it unset should preserve whatever name the
	// import pipeline parsed out of the source file (e.g. `metadata.name`
	// from a design YAML, or the derived name from a Kubernetes manifest).
	// Fall back to the filename (without extension) when neither the caller
	// nor the conversion produced a name.
	if variant.Name != "" {
		design.Name = variant.Name
	} else if design.Name == "" {
		if name := pCore.DesignNameFromFileName(fileToImport.FileName); name != "" {
			design.Name = name
		}
	}
	patternFile, err := encoding.Marshal(design)

	if err != nil {
		h.log.Error(ErrSavePattern(err))
		writeMeshkitError(rw, ErrSavePattern(err), http.StatusInternalServerError)

		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrSavePattern(err),
		}).WithDescription(ErrSavePattern(err).Error()).Build()

		h.persistAndPublish(provider, token, userID, event)
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
		h.handleProviderPatternSaveError(rw, eventBuilder, userID, savedDesignByt, err, provider, token)
		return
	}

	_, _ = rw.Write(savedDesignByt)

	event := eventBuilder.WithSeverity(events.Success).WithDescription(fmt.Sprintf("Imported design '%s' of type '%s'", design.Name, sourceFileType)).Build()
	h.persistAndPublish(provider, token, userID, event)

}
