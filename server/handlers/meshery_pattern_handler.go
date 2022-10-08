package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/internal/sql"
	"github.com/layer5io/meshery/server/models"

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
	K8sManifest   string                 `json:"K8sManifest,omitempty"`
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
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	var parsedBody *MesheryPatternRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
		// rw.WriteHeader(http.StatusBadRequest)
		// fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
		return
	}

	format := r.URL.Query().Get("output")

	// If Content is not empty then assume it's a local upload
	if parsedBody.PatternData != nil {
		// Assign a name if no name is provided
		if parsedBody.PatternData.Name == "" {
			patternName, err := models.GetPatternName(parsedBody.PatternData.PatternFile)
			if err != nil {
				h.log.Error(ErrSavePattern(err))
				http.Error(rw, ErrSavePattern(err).Error(), http.StatusBadRequest)
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
				return
			}
			go h.config.ConfigurationChannel.PublishPatterns()
			formatPatternOutput(rw, resp, format)
			return
		}

		byt, err := json.Marshal([]models.MesheryPattern{*mesheryPattern})
		if err != nil {
			h.log.Error(ErrEncodePattern(err))
			http.Error(rw, ErrEncodePattern(err).Error(), http.StatusInternalServerError)
			return
		}

		formatPatternOutput(rw, byt, format)
		return
	}

	if parsedBody.URL != "" {
		resp, err := provider.RemotePatternFile(r, parsedBody.URL, parsedBody.Path, parsedBody.Save)

		if err != nil {
			h.log.Error(ErrImportPattern(err))
			http.Error(rw, ErrImportPattern(err).Error(), http.StatusInternalServerError)
			return
		}

		formatPatternOutput(rw, resp, format)
		return
	}

	if parsedBody.CytoscapeJSON != "" {
		pf, err := pCore.NewPatternFileFromCytoscapeJSJSON(parsedBody.Name, []byte(parsedBody.CytoscapeJSON))
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(rw, "%s", err)
			return
		}

		pfByt, err := pf.ToYAML()
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(rw, "%s", err)
			return
		}

		patternName, err := models.GetPatternName(string(pfByt))
		if err != nil {
			h.log.Error(ErrGetPattern(err))
			http.Error(rw, ErrGetPattern(err).Error(), http.StatusBadRequest)
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

		if parsedBody.Save {
			resp, err := provider.SaveMesheryPattern(token, mesheryPattern)
			if err != nil {
				h.log.Error(ErrSavePattern(err))
				http.Error(rw, ErrSavePattern(err).Error(), http.StatusInternalServerError)
				return
			}

			go h.config.ConfigurationChannel.PublishPatterns()
			formatPatternOutput(rw, resp, format)
			return
		}

		byt, err := json.Marshal([]models.MesheryPattern{*mesheryPattern})
		if err != nil {
			h.log.Error(ErrEncodePattern(err))
			http.Error(rw, ErrEncodePattern(err).Error(), http.StatusInternalServerError)
			return
		}

		formatPatternOutput(rw, byt, format)
		return
	}

	if parsedBody.K8sManifest != "" {
		pattern, err := pCore.NewPatternFileFromK8sManifest(parsedBody.K8sManifest, false)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to convert to pattern: %s", err), http.StatusBadRequest)
			return
		}

		patternYAML, err := pattern.ToYAML()
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to generate pattern: %s", err), http.StatusInternalServerError)
			return
		}

		name, err := models.GetPatternName(string(patternYAML))
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to get pattern name: %s", err), http.StatusInternalServerError)
			return
		}

		patternModel := &models.MesheryPattern{
			Name:        name,
			PatternFile: string(patternYAML),
			Location: sql.Map{
				"host": "",
				"path": "",
				"type": "local",
			},
		}

		if parsedBody.Save {
			resp, err := provider.SaveMesheryPattern(token, patternModel)
			if err != nil {
				http.Error(rw, fmt.Sprintf("failed to save the pattern: %s", err), http.StatusInternalServerError)
				return
			}

			go h.config.ConfigurationChannel.PublishPatterns()
			formatPatternOutput(rw, resp, format)
			return
		}

		byt, err := json.Marshal([]models.MesheryPattern{*patternModel})
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to encode pattern: %s", err), http.StatusInternalServerError)
			return
		}

		formatPatternOutput(rw, byt, format)
		return
	}
}

// swagger:route GET /api/pattern PatternsAPI idGetPatternFiles
// Handle GET request for patterns
//
// Returns the list of all the patterns saved by the current user
// This will return all the patterns with their details
// responses:
// 	200: mesheryPatternsResponseWrapper

// GetMesheryPatternsHandler returns the list of all the patterns saved by the current user
func (h *Handler) GetMesheryPatternsHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()
	tokenString := r.Context().Value(models.TokenCtxKey).(string)

	resp, err := provider.GetMesheryPatterns(tokenString, q.Get("page"), q.Get("page_size"), q.Get("search"), q.Get("order"))
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
	prefObj *models.Preference,
	user *models.User,
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
	prefObj *models.Preference,
	user *models.User,
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

// swagger:route POST /api/pattern/clone/{id} PatternsAPI idCloneMesheryPattern
// Handle Clone for a Meshery Pattern
//
// Creates a local copy of a public pattern with id: id
// responses:
//
//	200: noContentWrapper
//
// CloneMesheryPatternHandler clones a pattern with the given id
func (h *Handler) CloneMesheryPatternHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	patternID := mux.Vars(r)["id"]

	resp, err := provider.CloneMesheryPattern(r, patternID)
	if err != nil {
		h.log.Error(ErrClonePattern(err))
		http.Error(rw, ErrClonePattern(err).Error(), http.StatusInternalServerError)
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
	prefObj *models.Preference,
	user *models.User,
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
	prefObj *models.Preference,
	user *models.User,
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

func formatPatternOutput(rw http.ResponseWriter, content []byte, format string) {
	contentMesheryPatternSlice := make([]models.MesheryPattern, 0)

	if err := json.Unmarshal(content, &contentMesheryPatternSlice); err != nil {
		http.Error(rw, ErrDecodePattern(err).Error(), http.StatusInternalServerError)
		// rw.WriteHeader(http.StatusInternalServerError)
		// fmt.Fprintf(rw, "failed to decode patterns data into go slice: %s", err)
		return
	}

	result := []models.MesheryPattern{}

	for _, content := range contentMesheryPatternSlice {
		if format == "cytoscape" {
			patternFile, err := pCore.NewPatternFile([]byte(content.PatternFile))
			if err != nil {
				http.Error(rw, ErrParsePattern(err).Error(), http.StatusBadRequest)
				// rw.WriteHeader(http.StatusBadRequest)
				// fmt.Fprintf(rw, "failed to parse to PatternFile: %s", err)
				return
			}

			//TODO: The below line has to go away once the client fully supports referencing variables  and pattern imports inside design
			newpatternfile := evalImportAndReferenceStage(&patternFile)

			cyjs, _ := newpatternfile.ToCytoscapeJS()

			bytes, err := json.Marshal(&cyjs)
			if err != nil {
				http.Error(rw, ErrConvertPattern(err).Error(), http.StatusInternalServerError)
				// rw.WriteHeader(http.StatusInternalServerError)
				// fmt.Fprintf(rw, "failed to convert PatternFile to Cytoscape object: %s", err)
				return
			}

			// Replace the patternfile with cytoscape type data
			content.PatternFile = string(bytes)
		}

		result = append(result, content)
	}

	data, err := json.Marshal(&result)
	if err != nil {
		obj := "pattern file"
		http.Error(rw, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		// rw.WriteHeader(http.StatusInternalServerError)
		// fmt.Fprintf(rw, "failed to marshal pattern file: %s", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(data))
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
