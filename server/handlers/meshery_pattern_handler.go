package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/meshes"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/errors"

	pCore "github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/stages"
	"github.com/sirupsen/logrus"
)

// MesheryPatternRequestBody refers to the type of request body that
// SaveMesheryPattern would receive
type MesheryPatternRequestBody struct {
	Name          string                 `json:"name,omitempty"`
	URL           string                 `json:"url,omitempty"`
	Path          string                 `json:"path,omitempty"`
	Save          bool                   `json:"save,omitempty"`
	PatternData   *models.MesheryPattern `json:"pattern_data,omitempty"`
	CytoscapeJSON string                 `json:"cytoscape_json,omitempty"`
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

// swagger:route POST /api/pattern PatternsAPI idPostPatternFile
// Handle POST requests for patterns
//
// Edit/update a meshery pattern
// responses:
// 	200: mesheryPatternResponseWrapper

func (h *Handler) handlePatternPOST(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	var err error
	res := meshes.EventsResponse{
		Component:     "core",
		ComponentName: "Design",
		OperationId:   uuid.NewString(),
		EventType:     meshes.EventType_INFO,
	}
	var parsedBody *MesheryPatternRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
		addMeshkitErr(&res, ErrRequestBody(err))
		go h.EventsBuffer.Publish(&res)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
		addMeshkitErr(&res, ErrRequestBody(err))
		go h.EventsBuffer.Publish(&res)
		return
	}

	format := r.URL.Query().Get("output")

	if parsedBody.CytoscapeJSON != "" {
		pf, err := pCore.NewPatternFileFromCytoscapeJSJSON(parsedBody.Name, []byte(parsedBody.CytoscapeJSON))
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "%s", err)
			addMeshkitErr(&res, ErrSavePattern(err))
			go h.EventsBuffer.Publish(&res)
			return
		}

		pfByt, err := pf.ToYAML()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "%s", err)
			addMeshkitErr(&res, ErrSavePattern(err))
			go h.EventsBuffer.Publish(&res)
			return
		}

		patternName, err := models.GetPatternName(string(pfByt))
		if err != nil {
			h.log.Error(ErrGetPattern(err))
			http.Error(rw, ErrGetPattern(err).Error(), http.StatusBadRequest)
			addMeshkitErr(&res, ErrGetPattern(err))
			go h.EventsBuffer.Publish(&res)
			return
		}

		mesheryPattern := &models.MesheryPattern{
			Name:        patternName,
			PatternFile: string(pfByt),
			Location: map[string]interface{}{
				"host": "",
				"path": "",
				"type": "local",
			},
		}
		if parsedBody.PatternData != nil {
			mesheryPattern.ID = parsedBody.PatternData.ID
		}
		if parsedBody.Save {
			resp, err := provider.SaveMesheryPattern(token, mesheryPattern)
			if err != nil {
				h.log.Error(ErrSavePattern(err))
				http.Error(rw, ErrSavePattern(err).Error(), http.StatusInternalServerError)
				addMeshkitErr(&res, ErrSavePattern(err))
				go h.EventsBuffer.Publish(&res)
				return
			}

			go h.config.ConfigurationChannel.PublishPatterns()
			h.formatPatternOutput(rw, resp, format, &res)
			return
		}

		byt, err := json.Marshal([]models.MesheryPattern{*mesheryPattern})
		if err != nil {
			h.log.Error(ErrEncodePattern(err))
			http.Error(rw, ErrEncodePattern(err).Error(), http.StatusInternalServerError)
			addMeshkitErr(&res, ErrEncodePattern(err))
			go h.EventsBuffer.Publish(&res)
			return
		}

		h.formatPatternOutput(rw, byt, format, &res)
		return
	}
	// If Content is not empty then assume it's a local upload
	if parsedBody.PatternData != nil {
		// Assign a name if no name is provided
		if parsedBody.PatternData.Name == "" {
			patternName, err := models.GetPatternName(parsedBody.PatternData.PatternFile)
			if err != nil {
				h.log.Error(ErrSavePattern(err))
				http.Error(rw, ErrSavePattern(err).Error(), http.StatusBadRequest)
				addMeshkitErr(&res, ErrSavePattern(err))
				go h.EventsBuffer.Publish(&res)
				return
			}
			parsedBody.PatternData.Name = patternName
		}

		// Assign a location if no location is specified
		if parsedBody.PatternData.Location == nil {
			parsedBody.PatternData.Location = map[string]interface{}{
				"host":   "",
				"path":   "",
				"type":   "local",
				"branch": "",
			}
		}

		mesheryPattern := parsedBody.PatternData

		if parsedBody.Save {
			resp, err := provider.SaveMesheryPattern(token, mesheryPattern)
			if err != nil {
				h.log.Error(ErrSavePattern(err))
				http.Error(rw, ErrSavePattern(err).Error(), http.StatusInternalServerError)
				addMeshkitErr(&res, ErrSavePattern(err))
				go h.EventsBuffer.Publish(&res)
				return
			}
			go h.config.ConfigurationChannel.PublishPatterns()
			h.formatPatternOutput(rw, resp, format, &res)
			return
		}

		byt, err := json.Marshal([]models.MesheryPattern{*mesheryPattern})
		if err != nil {
			h.log.Error(ErrEncodePattern(err))
			http.Error(rw, ErrEncodePattern(err).Error(), http.StatusInternalServerError)
			addMeshkitErr(&res, ErrSavePattern(err))
			go h.EventsBuffer.Publish(&res)
			return
		}

		h.formatPatternOutput(rw, byt, format, &res)
		return
	}

	if parsedBody.URL != "" {
		resp, err := provider.RemotePatternFile(r, parsedBody.URL, parsedBody.Path, parsedBody.Save)

		if err != nil {
			h.log.Error(ErrImportPattern(err))
			http.Error(rw, ErrImportPattern(err).Error(), http.StatusInternalServerError)
			addMeshkitErr(&res, ErrImportPattern(err))
			go h.EventsBuffer.Publish(&res)
			return
		}

		h.formatPatternOutput(rw, resp, format, &res)
		return
	}
	//Depracated: The below logic was used when applications were stored as k8s_manifests.
	// if parsedBody.K8sManifest != "" {
	// 	pattern, err := pCore.NewPatternFileFromK8sManifest(parsedBody.K8sManifest, false)
	// 	if err != nil {
	// 		http.Error(rw, fmt.Sprintf("failed to convert to pattern: %s", err), http.StatusBadRequest)
	// 		fmt.Println("err: ", err)
	// 		addMeshkitErr(&res, err) //this error is already a meshkit error so no further wrapping required
	// 		go h.EventsBuffer.Publish(&res)
	// 		return
	// 	}

	// 	patternYAML, err := pattern.ToYAML()
	// 	if err != nil {
	// 		http.Error(rw, fmt.Sprintf("failed to generate pattern: %s", err), http.StatusInternalServerError)
	// 		addMeshkitErr(&res, ErrSavePattern(err))
	// 		go h.EventsBuffer.Publish(&res)
	// 		return
	// 	}

	// 	name, err := models.GetPatternName(string(patternYAML))
	// 	if err != nil {
	// 		http.Error(rw, fmt.Sprintf("failed to get pattern name: %s", err), http.StatusInternalServerError)
	// 		addMeshkitErr(&res, ErrSavePattern(err))
	// 		go h.EventsBuffer.Publish(&res)
	// 		return
	// 	}

	// 	patternModel := &models.MesheryPattern{
	// 		Name:        name,
	// 		PatternFile: string(patternYAML),
	// 		Location: sql.Map{
	// 			"host": "",
	// 			"path": "",
	// 			"type": "local",
	// 		},
	// 	}

	// 	if parsedBody.Save {
	// 		resp, err := provider.SaveMesheryPattern(token, patternModel)
	// 		if err != nil {
	// 			http.Error(rw, fmt.Sprintf("failed to save the pattern: %s", err), http.StatusInternalServerError)
	// 			addMeshkitErr(&res, ErrSavePattern(err))
	// 			go h.EventsBuffer.Publish(&res)
	// 			return
	// 		}

	// 		go h.config.ConfigurationChannel.PublishPatterns()
	// 		h.formatPatternOutput(rw, resp, format, &res)
	// 		return
	// 	}

	// 	byt, err := json.Marshal([]models.MesheryPattern{*patternModel})
	// 	if err != nil {
	// 		http.Error(rw, fmt.Sprintf("failed to encode pattern: %s", err), http.StatusInternalServerError)
	// 		addMeshkitErr(&res, ErrSavePattern(err))
	// 		go h.EventsBuffer.Publish(&res)
	// 		return
	// 	}

	// 	h.formatPatternOutput(rw, byt, format, &res)
	// 	return
	// }
}

// swagger:route GET /api/pattern PatternsAPI idGetPatternFiles
// Handle GET request for patterns
//
// Returns the list of all the patterns saved by the current user
// This will return all the patterns with their details
//
// ```?updated_after=<timestamp>``` timestamp should be of the format "YYYY-MM-DD HH:MM:SS"
//
// ```?order={field}``` orders on the passed field
//
// ```?search=<design name>``` A string matching is done on the specified design name
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 10
// responses:
// 	200: mesheryPatternsResponseWrapper

// GetMesheryPatternsHandler returns the list of all the patterns saved by the current user
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
	resp, err := provider.GetMesheryPatterns(tokenString, q.Get("page"), q.Get("page_size"), q.Get("search"), q.Get("order"), updateAfter)
	if err != nil {
		h.log.Error(ErrFetchPattern(err))
		http.Error(rw, ErrFetchPattern(err).Error(), http.StatusInternalServerError)
		return
	}

	// token, err := provider.GetProviderToken(r)
	if err != nil {
		http.Error(rw, "failed to get user token", http.StatusInternalServerError)
		return
	}
	// mc := NewContentModifier(token, provider, prefObj, user.UserID)
	// //acts like a middleware, modifying the bytes lazily just before sending them back
	// err = mc.AddMetadataForPatterns(r.Context(), &resp)
	if err != nil {
		fmt.Println("Could not add metadata about pattern's current support ", err.Error())
	}
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route POST /api/pattern/catalog PatternsAPI idGetCatalogMesheryPatternsHandler
// Handle GET request for catalog patterns
//
// Used to get catalog patterns
// responses:
//
//	200: mesheryPatternResponseWrapper
func (h *Handler) GetCatalogMesheryPatternsHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()
	tokenString := r.Context().Value(models.TokenCtxKey).(string)

	resp, err := provider.GetCatalogMesheryPatterns(tokenString, q.Get("search"), q.Get("order"))
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
	_ *models.User,
	provider models.Provider,
) {
	patternID := mux.Vars(r)["id"]
	resp, err := provider.DeleteMesheryPattern(r, patternID)
	if err != nil {
		h.log.Error(ErrDeletePattern(err))
		http.Error(rw, ErrDeletePattern(err).Error(), http.StatusInternalServerError)
		return
	}

	go h.config.ConfigurationChannel.PublishPatterns()
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route GET /api/pattern/{id} PatternsAPI idGetMesheryPattern
// Handle GET request for Meshery Pattern with the given id
//
// Get the pattern with the given id
// responses:
//  200:

// GetMesheryPatternHandler returns the pattern file with the given id

func (h *Handler) DownloadMesheryPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	patternID := mux.Vars(r)["id"]
	resp, err := provider.GetMesheryPattern(r, patternID)
	if err != nil {
		h.log.Error(ErrGetPattern(err))
		http.Error(rw, ErrGetPattern(err).Error(), http.StatusNotFound)
		return
	}

	pattern := &models.MesheryPattern{}

	err = json.Unmarshal(resp, &pattern)
	if err != nil {
		obj := "download pattern"
		h.log.Error(ErrUnmarshal(err, obj))
		http.Error(rw, ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/x-yaml")
	if _, err := io.Copy(rw, strings.NewReader(pattern.PatternFile)); err != nil {
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
//	200: noContentWrapper
//
// CloneMesheryPatternHandler clones a pattern with the given id
func (h *Handler) CloneMesheryPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	patternID := mux.Vars(r)["id"]
	var parsedBody *models.MesheryClonePatternRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil || patternID == "" {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
		return
	}

	resp, err := provider.CloneMesheryPattern(r, patternID, parsedBody)
	if err != nil {
		h.log.Error(ErrClonePattern(err))
		http.Error(rw, ErrClonePattern(err).Error(), http.StatusInternalServerError)
		return
	}

	go h.config.ConfigurationChannel.PublishPatterns()
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route POST /api/pattern/catalog/publish PatternsAPI idPublishCatalogPatternHandler
// Handle Publish for a Meshery Pattern
//
// Publishes pattern to Meshery Catalog by setting visibility to published and setting catalog data
// responses:
//
//	200: noContentWrapper
//
// PublishCatalogPatternHandler sets visibility of pattern with given id as published
func (h *Handler) PublishCatalogPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	var parsedBody *models.MesheryCatalogPatternRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
		return
	}
	resp, err := provider.PublishCatalogPattern(r, parsedBody)
	if err != nil {
		h.log.Error(ErrPublishCatalogPattern(err))
		http.Error(rw, ErrPublishCatalogPattern(err).Error(), http.StatusInternalServerError)
		return
	}

	go h.config.ConfigurationChannel.PublishPatterns()
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
	_ *models.User,
	provider models.Provider,
) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error(rw, "err deleting pattern, converting bytes: ", err)
	}

	var patterns models.MesheryPatternDeleteRequestBody
	err = json.Unmarshal([]byte(body), &patterns)
	if err != nil {
		logrus.Error("error marshaling patterns json: ", err)
	}

	logrus.Debugf("patterns to be deleted: %+v", patterns)

	resp, err := provider.DeleteMesheryPatterns(r, patterns)

	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to delete the pattern: %s", err), http.StatusInternalServerError)
		return
	}

	go h.config.ConfigurationChannel.PublishPatterns()
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route GET /api/pattern/{id} PatternsAPI idGetMesheryPattern
// Handle GET for a Meshery Pattern
//
// Fetches the pattern with the given id
// responses:
// 	200: mesheryPatternResponseWrapper

// GetMesheryPatternHandler fetched the pattern with the given id
func (h *Handler) GetMesheryPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	patternID := mux.Vars(r)["id"]

	resp, err := provider.GetMesheryPattern(r, patternID)
	if err != nil {
		h.log.Error(ErrGetPattern(err))
		http.Error(rw, ErrGetPattern(err).Error(), http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

func (h *Handler) formatPatternOutput(rw http.ResponseWriter, content []byte, format string, res *meshes.EventsResponse) {
	contentMesheryPatternSlice := make([]models.MesheryPattern, 0)

	if err := json.Unmarshal(content, &contentMesheryPatternSlice); err != nil {
		http.Error(rw, ErrDecodePattern(err).Error(), http.StatusInternalServerError)
		addMeshkitErr(res, ErrDecodePattern(err))
		go h.EventsBuffer.Publish(res)
		return
	}

	result := []models.MesheryPattern{}
	names := []string{}
	for _, content := range contentMesheryPatternSlice {
		if format == "cytoscape" {
			patternFile, err := pCore.NewPatternFile([]byte(content.PatternFile))
			if err != nil {
				http.Error(rw, ErrParsePattern(err).Error(), http.StatusBadRequest)
				addMeshkitErr(res, ErrParsePattern(err))
				go h.EventsBuffer.Publish(res)
				return
			}

			//TODO: The below line has to go away once the client fully supports referencing variables  and pattern imports inside design
			newpatternfile := evalImportAndReferenceStage(&patternFile)

			cyjs, _ := newpatternfile.ToCytoscapeJS()

			bytes, err := json.Marshal(&cyjs)
			if err != nil {
				http.Error(rw, ErrConvertPattern(err).Error(), http.StatusInternalServerError)
				addMeshkitErr(res, ErrConvertPattern(err))
				go h.EventsBuffer.Publish(res)
				return
			}

			// Replace the patternfile with cytoscape type data
			content.PatternFile = string(bytes)
		}

		result = append(result, content)
		names = append(names, content.Name)
	}

	data, err := json.Marshal(&result)
	if err != nil {
		obj := "pattern file"
		http.Error(rw, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		addMeshkitErr(res, ErrMarshal(err, obj))
		go h.EventsBuffer.Publish(res)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(data))
	res.Details = "\"" + strings.Join(names, ",") + "\" design saved"
	res.Summary = "Changes to the \"" + strings.Join(names, ",") + "\" design have been saved."
	// go h.EventsBuffer.Publish(res)
}

// Since the client currently does not support pattern imports and externalized variables, the first(import) stage of pattern engine
// is evaluated here to simplify the pattern file such that it is valid when a deploy takes place
func evalImportAndReferenceStage(p *pCore.Pattern) (newp pCore.Pattern) {
	sap := &serviceActionProvider{}
	sip := &serviceInfoProvider{}
	chain := stages.CreateChain()
	chain.
		Add(stages.Import(sip, sap)).
		Add(stages.Filler(false)).
		Add(func(data *stages.Data, err error, next stages.ChainStageNextFunction) {
			data.Lock.Lock()
			newp = *data.Pattern
			data.Lock.Unlock()
		}).
		Process(&stages.Data{
			Pattern: p,
		})
	return newp
}

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
