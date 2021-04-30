package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/models"
	OAM "github.com/layer5io/meshery/models/oam"
)

// MesheryPatternRequestBody refers to the type of request body that
// SaveMesheryPattern would receive
type MesheryPatternRequestBody struct {
	URL           string                 `json:"url,omitempty"`
	Path          string                 `json:"path,omitempty"`
	Save          bool                   `json:"save,omitempty"`
	PatternData   *models.MesheryPattern `json:"pattern_data,omitempty"`
	CytoscapeJSON string                 `json:"cytoscape_json,omitempty"`
}

// PatternFileRequestHandler will handle requests of both type GET and POST
// on the route /api/experimental/patterns
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
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(rw, "failed to read request body: %s", err)
		return
	}

	token, err := provider.GetProviderToken(r)
	if err != nil {
		http.Error(rw, "failed to get user token", http.StatusInternalServerError)
		return
	}

	format := r.URL.Query().Get("output")

	// If Content is not empty then assume it's a local upload
	if parsedBody.PatternData != nil {
		patternName, err := models.GetPatternName(parsedBody.PatternData.PatternFile)
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to save the pattern: %s", err), http.StatusBadRequest)
			return
		}

		// Assign a name if no name is provided
		if parsedBody.PatternData.Name == "" {
			parsedBody.PatternData.Name = patternName
		}
		// Assign a location if no location is specified
		if parsedBody.PatternData.Location == nil {
			parsedBody.PatternData.Location = map[string]interface{}{
				"host": "",
				"path": "",
				"type": "local",
			}
		}

		mesheryPattern := parsedBody.PatternData

		if parsedBody.Save {
			resp, err := provider.SaveMesheryPattern(token, mesheryPattern)
			if err != nil {
				http.Error(rw, fmt.Sprintf("failed to save the pattern: %s", err), http.StatusInternalServerError)
				return
			}

			formatOutput(rw, resp, format)
			return
		}

		byt, err := json.Marshal([]models.MesheryPattern{*mesheryPattern})
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to encode pattern: %s", err), http.StatusInternalServerError)
			return
		}

		formatOutput(rw, byt, format)
		return
	}

	if parsedBody.URL != "" {
		resp, err := provider.RemotePatternFile(r, parsedBody.URL, parsedBody.Path, parsedBody.Save)

		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to import pattern: %s", err), http.StatusInternalServerError)
			return
		}

		formatOutput(rw, resp, format)
		return
	}

	if parsedBody.CytoscapeJSON != "" {
		pf, err := OAM.NewPatternFileFromCytoscapeJSJSON([]byte(parsedBody.CytoscapeJSON))
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
			http.Error(rw, fmt.Sprintf("failed to get the pattern name: %s", err), http.StatusBadRequest)
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
				http.Error(rw, fmt.Sprintf("failed to save the pattern: %s", err), http.StatusInternalServerError)
				return
			}

			formatOutput(rw, resp, format)
			return
		}

		byt, err := json.Marshal([]models.MesheryPattern{*mesheryPattern})
		if err != nil {
			http.Error(rw, fmt.Sprintf("failed to encode pattern: %s", err), http.StatusInternalServerError)
			return
		}

		formatOutput(rw, byt, format)
		return
	}
}

// GetMesheryPatternsHandler returns the list of all the patterns saved by the current user
func (h *Handler) GetMesheryPatternsHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()

	resp, err := provider.GetMesheryPatterns(r, q.Get("page"), q.Get("page_size"), q.Get("search"), q.Get("order"))
	if err != nil {
		http.Error(rw, fmt.Sprintf("failed to fetch the patterns: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

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
		http.Error(rw, fmt.Sprintf("failed to delete the pattern: %s", err), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

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
		http.Error(rw, fmt.Sprintf("failed to get the pattern: %s", err), http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

func formatOutput(rw http.ResponseWriter, content []byte, format string) {
	contentMesheryPatternSlice := make([]models.MesheryPattern, 0)

	if err := json.Unmarshal(content, &contentMesheryPatternSlice); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "failed to decode patterns data into go slice: %s", err)
		return
	}

	result := []models.MesheryPattern{}

	for _, content := range contentMesheryPatternSlice {
		if format == "cytoscape" {
			patternFile, err := OAM.NewPatternFile([]byte(content.PatternFile))
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(rw, "failed to parse to PatternFile: %s", err)
				return
			}

			cyjs, _ := patternFile.ToCytoscapeJS()

			bytes, err := json.Marshal(&cyjs)
			if err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(rw, "failed to convert PatternFile to Cytoscape object: %s", err)
				return
			}

			// Replace the patternfile with cytoscape type data
			content.PatternFile = string(bytes)
		}

		result = append(result, content)
	}

	data, err := json.Marshal(&result)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "failed to marshal pattern file: %s", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(data))
}
