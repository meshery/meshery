package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	guid "github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/meshes"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/utils"
	"github.com/layer5io/meshkit/models/events"
	regv1beta1 "github.com/layer5io/meshkit/models/meshmodel/registry/v1beta1"
	"github.com/meshery/schemas/models/v1beta1"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/model"
)

// swagger:route GET /api/filter/file/{id} FiltersAPI idGetFilterFile
// Handle GET request for filter file with given id
//
// Returns the Meshery Filter file saved by the current user with the given id
// responses:
//
//	200: mesheryFilterResponseWrapper
func (h *Handler) GetMesheryFilterFileHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	filterID := mux.Vars(r)["id"]

	resp, err := provider.GetMesheryFilterFile(r, filterID)
	if err != nil {
		h.log.Error(ErrGetFilter(err))
		http.Error(rw, ErrGetFilter(err).Error(), http.StatusNotFound)
		return
	}

	reader := bytes.NewReader(resp)
	rw.Header().Set("Content-Type", "application/wasm")
	_, err = io.Copy(rw, reader)
	if err != nil {
		h.log.Error(ErrDownloadWASMFile(err, "download"))
		http.Error(rw, ErrDownloadWASMFile(err, "download").Error(), http.StatusInternalServerError)
	}
}

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
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {

	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("filter").WithAction("update")

	defer func() {
		_ = r.Body.Close()
	}()
	res := meshes.EventsResponse{
		Component:     "core",
		ComponentName: "Filters",
		OperationId:   guid.NewString(),
		EventType:     meshes.EventType_INFO,
	}
	var parsedBody *models.MesheryFilterRequestBody

	actedUpon := &userID
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		invalidReqBody := ErrRequestBody(err)
		h.log.Error(invalidReqBody)

		event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": invalidReqBody,
		}).WithDescription(fmt.Sprintf("Filter %s is corrupted.", parsedBody.FilterData.Name)).Build()

		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		http.Error(rw, ErrSaveFilter(err).Error(), http.StatusBadRequest)
		addMeshkitErr(&res, ErrGetFilter(err))
		go h.EventsBuffer.Publish(&res)
		return
	}

	if parsedBody.FilterData != nil && parsedBody.FilterData.ID != nil {
		actedUpon = parsedBody.FilterData.ID
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
		addMeshkitErr(&res, ErrRetrieveUserToken(err))
		go h.EventsBuffer.Publish(&res)
		return
	}

	format := r.URL.Query().Get("output")

	filterResource, err := h.generateFilterComponent(parsedBody.Config)
	if err != nil {
		h.log.Error(ErrEncodeFilter(err))
		http.Error(rw, ErrEncodeFilter(err).Error(), http.StatusInternalServerError)
		return
	}

	// If Content is not empty then assume it's a local upload
	if parsedBody.FilterData != nil {
		// Assign a name if no name is provided
		if parsedBody.FilterData.Name == "" {
			parsedBody.FilterData.Name = "meshery-filter-" + utils.GetRandomAlphabetsOfDigit(5)
		}
		// Assign a location if no location is specified
		if len(parsedBody.FilterData.Location) == 0 {
			parsedBody.FilterData.Location = map[string]interface{}{
				"host":   "",
				"path":   "",
				"type":   "local",
				"branch": "",
			}
		}

		mesheryFilter := models.MesheryFilter{
			FilterFile:     parsedBody.FilterData.FilterFile,
			Name:           parsedBody.FilterData.Name,
			ID:             parsedBody.FilterData.ID,
			UserID:         parsedBody.FilterData.UserID,
			UpdatedAt:      parsedBody.FilterData.UpdatedAt,
			Location:       parsedBody.FilterData.Location,
			FilterResource: filterResource,
			CatalogData:    parsedBody.FilterData.CatalogData,
		}

		if parsedBody.Save {
			resp, err := provider.SaveMesheryFilter(token, &mesheryFilter)
			if err != nil {
				errFilterSave := ErrSaveFilter(err)
				h.log.Error(errFilterSave)
				http.Error(rw, errFilterSave.Error(), http.StatusInternalServerError)

				event := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
					"error": errFilterSave,
				}).WithDescription(fmt.Sprintf("Failed persisting filter %s", parsedBody.FilterData.Name)).Build()

				_ = provider.PersistEvent(event)
				go h.config.EventBroadcaster.Publish(userID, event)
				addMeshkitErr(&res, ErrSaveFilter(err))
				go h.EventsBuffer.Publish(&res)
				return
			}

			go h.config.FilterChannel.Publish(userID, struct{}{})
			h.formatFilterOutput(rw, resp, format, &res, eventBuilder)

			eventBuilder.WithSeverity(events.Informational).Build()
			return
		}

		byt, err := json.Marshal([]models.MesheryFilter{mesheryFilter})
		if err != nil {
			h.log.Error(ErrEncodeFilter(err))
			http.Error(rw, ErrEncodeFilter(err).Error(), http.StatusInternalServerError)
			return
		}

		h.formatFilterOutput(rw, byt, format, &res, eventBuilder)
		_ = provider.PersistEvent(eventBuilder.Build())
		return
	}

	if parsedBody.URL != "" {
		resp, err := provider.RemoteFilterFile(r, parsedBody.URL, parsedBody.Path, parsedBody.Save, filterResource)

		if err != nil {
			h.log.Error(ErrImportFilter(err))
			http.Error(rw, ErrImportFilter(err).Error(), http.StatusInternalServerError)
			return
		}

		h.formatFilterOutput(rw, resp, format, &res, eventBuilder)
		_ = provider.PersistEvent(eventBuilder.Build())
		return
	}
}

// swagger:route GET /api/filter FiltersAPI idGetFilterFiles
// Handle GET request for filters
//
// # Returns the list of all the filters saved by the current user
//
// ```?order={field}``` orders on the passed field
//
// ```?search=<filter name>``` A string matching is done on the specified filter name
//
// ```?page={page-number}``` Default page number is 0
//
// ```?pagesize={pagesize}``` Default pagesize is 10
//
// ```?visibility={[visibility]}``` Default visibility is public + private; Mulitple visibility filters can be passed as an array
// Eg: ```?visibility=["public", "published"]``` will return public and published filters
//
// responses:
//
//	200: mesheryFiltersResponseWrapper
func (h *Handler) GetMesheryFiltersHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()
	tokenString := r.Context().Value(models.TokenCtxKey).(string)

	filter := struct {
		Visibility []string `json:"visibility"`
	}{}

	visibility := q.Get("visibility")
	if visibility != "" {
		err := json.Unmarshal([]byte(visibility), &filter.Visibility)
		if err != nil {
			h.log.Error(ErrFetchFilter(err))
			http.Error(rw, ErrFetchFilter(err).Error(), http.StatusInternalServerError)
			return
		}
	}

	resp, err := provider.GetMesheryFilters(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), filter.Visibility)
	if err != nil {
		h.log.Error(ErrFetchFilter(err))
		http.Error(rw, ErrFetchFilter(err).Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route GET /api/filter/catalog FiltersAPI idGetCatalogMesheryFiltersHandler
// Handle GET request for catalog filters
//
// # Filters can be further filtered through query parameter
//
// ```?order={field}``` orders on the passed field
//
// ```?page={page-number}``` Default page number is 0
//
// ```?pagesize={pagesize}``` Default pagesize is 10.
//
// ```?search={filtername}``` If search is non empty then a greedy search is performed
// responses:
//
//	200: mesheryFiltersResponseWrapper
func (h *Handler) GetCatalogMesheryFiltersHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
) {
	q := r.URL.Query()
	tokenString := r.Context().Value(models.TokenCtxKey).(string)

	resp, err := provider.GetCatalogMesheryFilters(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"))
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
	_ *models.Preference,
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

	go h.config.FilterChannel.Publish(uuid.FromStringOrNil(user.ID), struct{}{})
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route POST /api/filter/clone/{id} FiltersAPI idCloneMesheryFilter
// Handle Clone for a Meshery Filter
//
// Creates a local copy of a published filter with id: id
// responses:
//
//	200: noContentWrapper
//
// CloneMesheryFilterHandler clones a filter with the given id
func (h *Handler) CloneMesheryFilterHandler(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	filterID := mux.Vars(r)["id"]
	var parsedBody *models.MesheryCloneFilterRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil || filterID == "" {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
		return
	}

	resp, err := provider.CloneMesheryFilter(r, filterID, parsedBody)
	if err != nil {
		h.log.Error(ErrCloneFilter(err))
		http.Error(rw, ErrCloneFilter(err).Error(), http.StatusInternalServerError)
		return
	}

	go h.config.FilterChannel.Publish(uuid.FromStringOrNil(user.ID), struct{}{})
	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(resp))
}

// swagger:route POST /api/filter/catalog/publish FiltersAPI idPublishCatalogFilterHandler
// Handle Publish for a Meshery Filter
//
// Publishes filter to Meshery Catalog by setting visibility to published and setting catalog data
// responses:
//
//	202: noContentWrapper
//
// PublishCatalogFilterHandler set visibility of filter with given id as published
func (h *Handler) PublishCatalogFilterHandler(
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
		WithCategory("filter").
		WithAction("publish").
		ActedUpon(userID)

	var parsedBody *models.MesheryCatalogFilterRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRequestBody(err))
		e := eventBuilder.WithSeverity(events.Error).
			WithMetadata(map[string]interface{}{
				"error": ErrRequestBody(err),
			}).
			WithDescription("Error parsing filter payload.").Build()
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
		return
	}

	resp, err := provider.PublishCatalogFilter(r, parsedBody)
	if err != nil {
		h.log.Error(ErrPublishCatalogFilter(err))
		e := eventBuilder.WithSeverity(events.Error).
			WithMetadata(map[string]interface{}{
				"error": ErrPublishCatalogFilter(err),
			}).
			WithDescription("Error publishing filter.").Build()
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrPublishCatalogFilter(err).Error(), http.StatusInternalServerError)
		return
	}

	var respBody *models.CatalogRequest
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		h.log.Error(ErrPublishCatalogFilter(err))
		e := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrPublishCatalogFilter(err),
		}).WithDescription("Error parsing response.").Build()
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrPublishCatalogFilter(err).Error(), http.StatusInternalServerError)
	}

	e := eventBuilder.WithSeverity(events.Informational).ActedUpon(parsedBody.ID).WithDescription(fmt.Sprintf("Request to publish '%s' filter submitted with status: %s", respBody.ContentName, respBody.Status)).Build()
	_ = provider.PersistEvent(e)
	go h.config.EventBroadcaster.Publish(userID, e)

	go h.config.FilterChannel.Publish(uuid.FromStringOrNil(user.ID), struct{}{})
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusAccepted)
	fmt.Fprint(rw, string(resp))
}

// swagger:route DELETE /api/filter/catalog/unpublish FiltersAPI idUnPublishCatalogFilterHandler
// Handle UnPublish for a Meshery Filter
//
// Unpublishes filter from Meshery Catalog by setting visibility to private and removing catalog data from website
// responses:
//
//	200: noContentWrapper
//
// UnPublishCatalogFilterHandler sets visibility of filter with given id as private
func (h *Handler) UnPublishCatalogFilterHandler(
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
		WithCategory("filter").
		WithAction("unpublish_request").
		ActedUpon(userID)

	var parsedBody *models.MesheryCatalogFilterRequestBody
	if err := json.NewDecoder(r.Body).Decode(&parsedBody); err != nil {
		h.log.Error(ErrRequestBody(err))
		e := eventBuilder.WithSeverity(events.Error).
			WithMetadata(map[string]interface{}{
				"error": ErrRequestBody(err),
			}).
			WithDescription("Error parsing filter payload.").Build()
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)
		return
	}
	resp, err := provider.UnPublishCatalogFilter(r, parsedBody)
	if err != nil {
		h.log.Error(ErrPublishCatalogFilter(err))
		e := eventBuilder.WithSeverity(events.Error).
			WithMetadata(map[string]interface{}{
				"error": ErrPublishCatalogFilter(err),
			}).
			WithDescription("Error publishing filter.").Build()
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrPublishCatalogFilter(err).Error(), http.StatusInternalServerError)
		return
	}

	var respBody *models.CatalogRequest
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		h.log.Error(ErrPublishCatalogFilter(err))
		e := eventBuilder.WithSeverity(events.Error).WithMetadata(map[string]interface{}{
			"error": ErrPublishCatalogFilter(err),
		}).WithDescription("Error parsing response.").Build()
		_ = provider.PersistEvent(e)
		go h.config.EventBroadcaster.Publish(userID, e)
		http.Error(rw, ErrPublishCatalogFilter(err).Error(), http.StatusInternalServerError)
	}

	e := eventBuilder.WithSeverity(events.Informational).ActedUpon(parsedBody.ID).WithDescription(fmt.Sprintf("'%s' filter unpublished", respBody.ContentName)).Build()
	_ = provider.PersistEvent(e)
	go h.config.EventBroadcaster.Publish(userID, e)

	go h.config.FilterChannel.Publish(uuid.FromStringOrNil(user.ID), struct{}{})
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
	_ *models.Preference,
	_ *models.User,
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

func (h *Handler) formatFilterOutput(rw http.ResponseWriter, content []byte, _ string, res *meshes.EventsResponse, eventBuilder *events.EventBuilder) {
	contentMesheryFilterSlice := make([]models.MesheryFilter, 0)
	names := []string{}
	if err := json.Unmarshal(content, &contentMesheryFilterSlice); err != nil {
		http.Error(rw, ErrDecodeFilter(err).Error(), http.StatusInternalServerError)

		return
	}

	result := []models.MesheryFilter{}

	data, err := json.Marshal(&result)
	if err != nil {
		obj := "filter file"
		http.Error(rw, models.ErrMarshal(err, obj).Error(), http.StatusInternalServerError)

		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(data))
	for _, filter := range contentMesheryFilterSlice {
		names = append(names, filter.Name)
		if filter.ID != nil {
			eventBuilder.ActedUpon(*filter.ID)
		}
	}
	res.Details = "filters saved"
	res.Summary = "following filters were saved: " + strings.Join(names, ",")
	go h.EventsBuffer.Publish(res)
	eventBuilder.WithDescription(fmt.Sprintf("Filter %s saved", strings.Join(names, ",")))
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

func (h *Handler) generateFilterComponent(config string) (string, error) {
	res, _, _, _ := h.registryManager.GetEntities(&regv1beta1.ComponentFilter{
		Name:       "WASMFilter",
		Trim:       false,
		APIVersion: v1beta1.ComponentSchemaVersion,
		Version:    "v1.0.0",
		Limit:      1,
	})

	if len(res) > 0 {
		filterEntity := res[0]
		filterCompDef, ok := filterEntity.(*component.ComponentDefinition)
		if ok {
			filterID, _ := uuid.NewV4()
			filterSvc := component.ComponentDefinition{
				Id:          filterID,
				DisplayName: strings.ToLower(filterCompDef.Component.Kind) + utils.GetRandomAlphabetsOfDigit(5),
				Component: component.Component{
					Kind:    filterCompDef.Component.Kind,
					Version: filterCompDef.Component.Version,
				},
				Model: model.ModelDefinition{
					Name: filterCompDef.Model.Name,
					Model: model.Model{
						Version: filterCompDef.Model.Model.Version,
					},
				},
				Metadata: component.ComponentDefinition_Metadata{
					IsAnnotation: true,
				},
				Configuration: map[string]interface{}{
					"config": config,
				},
			}
			marshalledFilter, err := json.Marshal(filterSvc)
			if err != nil {
				return string(marshalledFilter), err
			}
			return string(marshalledFilter), nil
		}
	}
	return "", nil
}
