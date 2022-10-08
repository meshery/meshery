package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
)

// MesheryFilterRequestBody refers to the type of request body that
// SaveMesheryFilter would receive
type MesheryFilterRequestBody struct {
	URL        string                `json:"url,omitempty"`
	Path       string                `json:"path,omitempty"`
	Save       bool                  `json:"save,omitempty"`
	FilterData *models.MesheryFilter `json:"filter_data,omitempty"`
}

// swagger:route GET /api/filter/file/{id} FiltersAPI idGetFilterFiles
// Handle GET request for filter file with given id
//
// Returns the Meshery Filter file saved by the current user with the given id
// responses:
//
//	200: mesheryFilterResponseWrapper
func (h *Handler) GetMesheryFilterFileHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	filterID := mux.Vars(r)["id"]

	resp, err := provider.GetMesheryFilterFile(r, filterID)
	if err != nil {
		h.log.Error(ErrGetFilter(err))
		http.Error(rw, ErrGetFilter(err).Error(), http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route GET /api/filter FiltersAPI idGetFilterFile
// Handle GET request for all filters
//
// Returns all the Meshery Filters saved by the current user
// responses:
//
//	200: mesheryFiltersResponseWrapper
//
// FilterFileRequestHandler will handle requests of both type GET and POST
// on the route /api/filter
func (h *Handler) FilterFileRequestHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	if r.Method == http.MethodGet {
		h.GetMesheryFiltersHandler(rw, r, prefObj, user, provider)
		return
	}

	if r.Method == http.MethodPost {
		h.handleFilterPOST(rw, r, prefObj, user, provider)
		return
	}
}

// swagger:route POST /api/filter FiltersAPI idPostFilterFile
// Handle POST requests for Meshery Filters
//
// Used to save/update a Meshery Filter
// responses:
//
//	200: mesheryFilterResponseWrapper
func (h *Handler) handleFilterPOST(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	var parsedBody *MesheryFilterRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrGetFilter(err).Error(), http.StatusBadRequest)
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
	if parsedBody.FilterData != nil {
		// Assign a name if no name is provided
		if parsedBody.FilterData.Name == "" {
			// TODO: Dynamically generate names or get the name of the file from the UI (@navendu-pottekkat)
			parsedBody.FilterData.Name = "Test Filter"
		}
		// Assign a location if no location is specified
		if parsedBody.FilterData.Location == nil {
			parsedBody.FilterData.Location = map[string]interface{}{
				"host":   "",
				"path":   "",
				"type":   "local",
				"branch": "",
			}
		}

		mesheryFilter := parsedBody.FilterData

		if parsedBody.Save {
			resp, err := provider.SaveMesheryFilter(token, mesheryFilter)
			if err != nil {
				h.log.Error(ErrSaveFilter(err))
				http.Error(rw, ErrSaveFilter(err).Error(), http.StatusInternalServerError)
				return
			}

			go h.config.ConfigurationChannel.PublishFilters()
			formatFilterOutput(rw, resp, format)
			return
		}

		byt, err := json.Marshal([]models.MesheryFilter{*mesheryFilter})
		if err != nil {
			h.log.Error(ErrEncodeFilter(err))
			http.Error(rw, ErrEncodeFilter(err).Error(), http.StatusInternalServerError)
			return
		}

		formatFilterOutput(rw, byt, format)
		return
	}

	if parsedBody.URL != "" {
		resp, err := provider.RemoteFilterFile(r, parsedBody.URL, parsedBody.Path, parsedBody.Save)

		if err != nil {
			h.log.Error(ErrImportFilter(err))
			http.Error(rw, ErrImportFilter(err).Error(), http.StatusInternalServerError)
			return
		}

		formatFilterOutput(rw, resp, format)
		return
	}
}

// GetMesheryFiltersHandler returns the list of all the filters saved by the current user
func (h *Handler) GetMesheryFiltersHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()
	tokenString := r.Context().Value(models.TokenCtxKey).(string)

	resp, err := provider.GetMesheryFilters(tokenString, q.Get("page"), q.Get("page_size"), q.Get("search"), q.Get("order"))
	if err != nil {
		h.log.Error(ErrFetchFilter(err))
		http.Error(rw, ErrFetchFilter(err).Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route POST /api/filter/catalog FiltersAPI idGetCatalogMesheryFiltersHandler
// Handle GET request for catalog filters
//
// Used to get catalog filters
// responses:
//
//	200: mesheryFilterResponseWrapper
func (h *Handler) GetCatalogMesheryFiltersHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()
	tokenString := r.Context().Value(models.TokenCtxKey).(string)

	resp, err := provider.GetCatalogMesheryFilters(tokenString, q.Get("search"), q.Get("order"))
	if err != nil {
		h.log.Error(ErrFetchFilter(err))
		http.Error(rw, ErrFetchFilter(err).Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route DELETE /api/filter/{id} FiltersAPI idDeleteMesheryFilter
// Handle Delete for a Meshery Filter
//
// Deletes a meshery filter with ID: id
// responses:
//
//	200: noContentWrapper
//
// DeleteMesheryFilterHandler deletes a filter with the given id
func (h *Handler) DeleteMesheryFilterHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	filterID := mux.Vars(r)["id"]

	resp, err := provider.DeleteMesheryFilter(r, filterID)
	if err != nil {
		h.log.Error(ErrDeleteFilter(err))
		http.Error(rw, ErrDeleteFilter(err).Error(), http.StatusInternalServerError)
		return
	}

	go h.config.ConfigurationChannel.PublishFilters()
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route POST /api/filter/clone/{id} FiltersAPI idCloneMesheryFilter
// Handle Clone for a Meshery Filter
//
// Creates a local copy of a public filter with id: id
// responses:
//
//	200: noContentWrapper
//
// CloneMesheryFilterHandler clones a filter with the given id
func (h *Handler) CloneMesheryFilterHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	filterID := mux.Vars(r)["id"]

	resp, err := provider.CloneMesheryFilter(r, filterID)
	if err != nil {
		h.log.Error(ErrCloneFilter(err))
		http.Error(rw, ErrCloneFilter(err).Error(), http.StatusInternalServerError)
		return
	}

	go h.config.ConfigurationChannel.PublishFilters()
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route GET /api/filter/{id} FiltersAPI idGetMesheryFilter
// Handle GET request for a Meshery Filter
//
// Fetches the Meshery Filter with the given id
// responses:
// 	200: mesheryFilterResponseWrapper

// GetMesheryFilterHandler fetched the filter with the given id
func (h *Handler) GetMesheryFilterHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	filterID := mux.Vars(r)["id"]

	resp, err := provider.GetMesheryFilter(r, filterID)
	if err != nil {
		h.log.Error(ErrGetFilter(err))
		http.Error(rw, ErrGetFilter(err).Error(), http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

func formatFilterOutput(rw http.ResponseWriter, content []byte, format string) {
	contentMesheryFilterSlice := make([]models.MesheryFilter, 0)

	if err := json.Unmarshal(content, &contentMesheryFilterSlice); err != nil {
		http.Error(rw, ErrDecodeFilter(err).Error(), http.StatusInternalServerError)
		// rw.WriteHeader(http.StatusInternalServerError)
		// fmt.Fprintf(rw, "failed to decode filters data into go slice: %s", err)
		return
	}

	result := []models.MesheryFilter{}

	data, err := json.Marshal(&result)
	if err != nil {
		obj := "filter file"
		http.Error(rw, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		// rw.WriteHeader(http.StatusInternalServerError)
		// fmt.Fprintf(rw, "failed to marshal filter file: %s", err)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(data))
}

// swagger:route POST /api/filter/deploy FilterAPI idPostDeployFilterFile
// Handle POST request for Filter File Deploy
//
// Deploy an attached filter file with the request
// responses:
//  200: FilterFilesResponseWrapper

// swagger:route DELETE /api/filter/deploy FilterAPI idDeleteFilterFile
// Handle DELETE request for Filter File Deploy
//
// Delete a deployed filter file with the request
// responses:
//  200:

// FilterFileHandler handles the requested related to filter files
func (h *Handler) FilterFileHandler(
	rw http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	// Filter files are just pattern files
	h.PatternFileHandler(rw, r, prefObj, user, provider)
}
