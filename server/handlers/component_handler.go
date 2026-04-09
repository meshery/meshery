package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/meshery/meshery/server/helpers"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models"

	// "github.com/meshery/meshkit/errors"
	// "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/events"

	meshkitOci "github.com/meshery/meshkit/models/oci"
	"github.com/meshery/meshkit/models/registration"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	meshkitutils "github.com/meshery/meshkit/utils"

	_models "github.com/meshery/meshkit/models/meshmodel/core/v1beta1"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	schemav1beta1 "github.com/meshery/schemas/models/v1beta1"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/connection"
	_model "github.com/meshery/schemas/models/v1beta1/model"

	"github.com/meshery/meshkit/models/meshmodel/entity"
	"github.com/meshery/meshkit/models/meshmodel/registry"

	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	"gorm.io/gorm"
)

/**Meshmodel endpoints **/
const DefaultPageSizeForMeshModelComponents = 25

func (h *Handler) GetMeshmodelModelsByCategories(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	cat := mux.Vars(r)["category"]
	queryParams := r.URL.Query()
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	returnAnnotationComp := queryParams.Get("annotations")

	filter := &regv1beta1.ModelFilter{
		Category:    cat,
		Version:     queryParams.Get("version"),
		Limit:       limit,
		Offset:      offset,
		OrderOn:     order,
		Sort:        sort,
		Annotations: returnAnnotationComp,
	}
	if search != "" {
		filter.Greedy = true
		filter.DisplayName = search
	}
	entities, count, _, _ := h.registryManager.GetEntities(filter)
	var modelDefs []_model.ModelDefinition
	for _, model := range entities {
		model, ok := model.(*_model.ModelDefinition)
		if ok {
			modelDefs = append(modelDefs, *model)
		}
	}

	var pgSize int64
	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	res := models.MeshmodelsDuplicateAPIResponse{
		Page:     page,
		PageSize: int(pgSize),
		Count:    count,
		Models:   models.FindDuplicateModels(modelDefs),
	}

	if err := enc.Encode(res); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetMeshmodelModelsByCategoriesByModel(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	cat := mux.Vars(r)["category"]
	model := mux.Vars(r)["model"]
	queryParams := r.URL.Query()
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	var greedy bool
	if search == queryParamTrue {
		greedy = true
	}
	returnAnnotationComp := queryParams.Get("annotations")

	entities, count, _, _ := h.registryManager.GetEntities(&regv1beta1.ModelFilter{
		Category:    cat,
		Name:        model,
		Version:     queryParams.Get("version"),
		Limit:       limit,
		Offset:      offset,
		Greedy:      greedy,
		OrderOn:     order,
		Sort:        sort,
		Annotations: returnAnnotationComp,
	})

	var modelDefs []_model.ModelDefinition
	for _, model := range entities {
		model, ok := model.(*_model.ModelDefinition)
		if ok {
			modelDefs = append(modelDefs, *model)
		}
	}

	var pgSize int64
	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	res := models.MeshmodelsDuplicateAPIResponse{
		Page:     page,
		PageSize: int(pgSize),
		Count:    count,
		Models:   models.FindDuplicateModels(modelDefs),
	}

	if err := enc.Encode(res); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetMeshmodelModels(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	queryParams := r.URL.Query()
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	v := queryParams.Get("version")
	returnAnnotationComp := queryParams.Get("annotations")

	filter := &regv1beta1.ModelFilter{
		Id:          queryParams.Get("id"),
		Registrant:  queryParams.Get("registrant"),
		Version:     v,
		Limit:       limit,
		Offset:      offset,
		OrderOn:     order,
		Sort:        sort,
		Annotations: returnAnnotationComp,

		Components:    queryParams.Get("components") == queryParamTrue,
		Relationships: queryParams.Get("relationships") == queryParamTrue,
		Status:        queryParams.Get("status"),
		Trim:          queryParams.Get("trim") == queryParamTrue,
	}
	if search != "" {
		filter.DisplayName = search
		filter.Name = search
		filter.Greedy = true
	}

	entities, count, _, _ := h.registryManager.GetEntities(filter)
	var modelDefs []_model.ModelDefinition
	for _, model := range entities {
		model, ok := model.(*_model.ModelDefinition)
		if ok {
			modelDefs = append(modelDefs, *model)
		}
	}
	var pgSize int64
	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	res := models.MeshmodelsDuplicateAPIResponse{
		Page:     page,
		PageSize: int(pgSize),
		Count:    count,
		Models:   models.FindDuplicateModels(modelDefs),
	}

	if err := enc.Encode(res); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetMeshmodelModelsByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	name := mux.Vars(r)["model"]
	queryParams := r.URL.Query()
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	var greedy bool
	if search == queryParamTrue {
		greedy = true
	}
	v := queryParams.Get("version")
	returnAnnotationComp := queryParams.Get("annotations")
	entities, count, _, _ := h.registryManager.GetEntities(&regv1beta1.ModelFilter{
		Name:        name,
		Version:     v,
		Limit:       limit,
		Offset:      offset,
		Greedy:      greedy,
		OrderOn:     order,
		Sort:        sort,
		Annotations: returnAnnotationComp,

		Components:    queryParams.Get("components") == queryParamTrue,
		Relationships: queryParams.Get("relationships") == queryParamTrue,
	})

	var modelDefs []_model.ModelDefinition
	for _, model := range entities {
		model, ok := model.(*_model.ModelDefinition)
		if ok {
			modelDefs = append(modelDefs, *model)
		}
	}

	var pgSize int64
	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	res := models.MeshmodelsDuplicateAPIResponse{
		Page:     page,
		PageSize: int(pgSize),
		Count:    count,
		Models:   models.FindDuplicateModels(modelDefs),
	}

	if err := enc.Encode(res); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetMeshmodelCategories(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	filter := &regv1beta1.CategoryFilter{
		Limit:   limit,
		Offset:  offset,
		OrderOn: order,
		Sort:    sort,
	}
	if search != "" {
		filter.Greedy = true
		filter.Name = search
	}

	categories, count, _, _ := h.registryManager.GetEntities(filter)

	var pgSize int64

	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	res := models.MeshmodelCategoriesAPIResponse{
		Page:       page,
		PageSize:   int(pgSize),
		Count:      count,
		Categories: categories,
	}

	if err := enc.Encode(res); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetMeshmodelCategoriesByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	name := mux.Vars(r)["category"]
	var greedy bool
	if search == queryParamTrue {
		greedy = true
	}
	categories, count, _, _ := h.registryManager.GetEntities(&regv1beta1.CategoryFilter{
		Name:    name,
		Limit:   limit,
		Greedy:  greedy,
		Offset:  offset,
		OrderOn: order,
		Sort:    sort,
	})

	var pgSize int64

	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	res := models.MeshmodelCategoriesAPIResponse{
		Page:       page,
		PageSize:   int(pgSize),
		Count:      count,
		Categories: categories,
	}

	if err := enc.Encode(res); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetMeshmodelComponentsByNameByModelByCategory(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	name := mux.Vars(r)["name"]

	queryParams := r.URL.Query()
	var greedy bool
	if search == queryParamTrue {
		greedy = true
	}
	typ := mux.Vars(r)["model"]
	cat := mux.Vars(r)["category"]
	v := queryParams.Get("version")
	returnAnnotationComp := queryParams.Get("annotations")
	entities, count, _, _ := h.registryManager.GetEntities(&regv1beta1.ComponentFilter{
		Name:         name,
		CategoryName: cat,
		ModelName:    typ,
		APIVersion:   queryParams.Get("apiVersion"),
		Version:      v,
		Offset:       offset,
		Greedy:       greedy,
		Limit:        limit,
		OrderOn:      order,
		Sort:         sort,
		Annotations:  returnAnnotationComp,
	})

	comps := processComponentDefinitions(entities)

	var pgSize int64
	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelComponentsDuplicateAPIResponse{
		Page:       page,
		PageSize:   int(pgSize),
		Count:      count,
		Components: models.FindDuplicateComponents(comps),
	}

	if err := enc.Encode(response); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetMeshmodelComponentsByNameByCategory(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	name := mux.Vars(r)["name"]
	var greedy bool
	queryParams := r.URL.Query()
	if search == queryParamTrue {
		greedy = true
	}
	cat := mux.Vars(r)["category"]
	v := queryParams.Get("version")
	returnAnnotationComp := queryParams.Get("annotations")

	entities, count, _, _ := h.registryManager.GetEntities(&regv1beta1.ComponentFilter{
		Name:         name,
		ModelName:    queryParams.Get("model"),
		CategoryName: cat,
		APIVersion:   queryParams.Get("apiVersion"),
		Version:      v,
		Offset:       offset,
		Limit:        limit,
		Greedy:       greedy,
		OrderOn:      order,
		Sort:         sort,
		Annotations:  returnAnnotationComp,
	})
	comps := processComponentDefinitions(entities)

	var pgSize int64
	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelComponentsDuplicateAPIResponse{
		Page:       page,
		PageSize:   int(pgSize),
		Count:      count,
		Components: models.FindDuplicateComponents(comps),
	}

	if err := enc.Encode(response); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetMeshmodelComponentsByNameByModel(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	name := mux.Vars(r)["name"]
	var greedy bool
	queryParams := r.URL.Query()

	if search == queryParamTrue {
		greedy = true
	}
	typ := mux.Vars(r)["model"]
	v := queryParams.Get("version")

	returnAnnotationComp := queryParams.Get("annotations")

	entities, count, _, _ := h.registryManager.GetEntities(&regv1beta1.ComponentFilter{
		Name:        name,
		ModelName:   typ,
		APIVersion:  queryParams.Get("apiVersion"),
		Version:     v,
		Offset:      offset,
		Greedy:      greedy,
		Limit:       limit,
		OrderOn:     order,
		Sort:        sort,
		Annotations: returnAnnotationComp,
	})
	comps := processComponentDefinitions(entities)

	var pgSize int64
	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelComponentsDuplicateAPIResponse{
		Page:       page,
		PageSize:   int(pgSize),
		Count:      count,
		Components: models.FindDuplicateComponents(comps),
	}

	if err := enc.Encode(response); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetAllMeshmodelComponentsByName(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	name := mux.Vars(r)["name"]
	var greedy bool
	queryParams := r.URL.Query()
	if search == queryParamTrue {
		greedy = true
	}
	v := queryParams.Get("version")
	returnAnnotationComp := queryParams.Get("annotations")
	entities, count, _, _ := h.registryManager.GetEntities(&regv1beta1.ComponentFilter{
		Name:        name,
		Trim:        queryParams.Get("trim") == queryParamTrue,
		APIVersion:  queryParams.Get("apiVersion"),
		Version:     v,
		ModelName:   queryParams.Get("model"),
		Offset:      offset,
		Limit:       limit,
		Greedy:      greedy,
		OrderOn:     order,
		Sort:        sort,
		Annotations: returnAnnotationComp,
	})

	comps := processComponentDefinitions(entities)

	var pgSize int64
	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelComponentsDuplicateAPIResponse{
		Page:       page,
		PageSize:   int(pgSize),
		Count:      count,
		Components: models.FindDuplicateComponents(comps),
	}

	if err := enc.Encode(response); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetMeshmodelComponentByModel(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	typ := mux.Vars(r)["model"]
	queryParams := r.URL.Query()
	v := queryParams.Get("version")

	returnAnnotationComp := queryParams.Get("annotations")
	filter := &regv1beta1.ComponentFilter{
		Id:          queryParams.Get("id"),
		ModelName:   typ,
		Version:     v,
		Trim:        queryParams.Get("trim") == queryParamTrue,
		APIVersion:  queryParams.Get("apiVersion"),
		Limit:       limit,
		Offset:      offset,
		OrderOn:     order,
		Sort:        sort,
		Annotations: returnAnnotationComp,
	}
	if search != "" {
		filter.Greedy = true
		filter.DisplayName = search
	}
	entities, count, _, _ := h.registryManager.GetEntities(filter)
	comps := processComponentDefinitions(entities)

	var pgSize int64
	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelComponentsDuplicateAPIResponse{
		Page:       page,
		PageSize:   int(pgSize),
		Count:      count,
		Components: models.FindDuplicateComponents(comps),
	}

	if err := enc.Encode(response); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetMeshmodelComponentByModelByCategory(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	typ := mux.Vars(r)["model"]
	cat := mux.Vars(r)["category"]
	queryParams := r.URL.Query()
	v := queryParams.Get("version")
	returnAnnotationComp := queryParams.Get("annotations")
	filter := &regv1beta1.ComponentFilter{
		CategoryName: cat,
		ModelName:    typ,
		Version:      v,
		Trim:         queryParams.Get("trim") == queryParamTrue,
		APIVersion:   queryParams.Get("apiVersion"),
		Limit:        limit,
		Offset:       offset,
		OrderOn:      order,
		Sort:         sort,
		Annotations:  returnAnnotationComp,
	}
	if search != "" {
		filter.Greedy = true
		filter.DisplayName = search
	}
	entities, count, _, _ := h.registryManager.GetEntities(filter)
	comps := processComponentDefinitions(entities)

	var pgSize int64
	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelComponentsDuplicateAPIResponse{
		Page:       page,
		PageSize:   int(pgSize),
		Count:      count,
		Components: models.FindDuplicateComponents(comps),
	}

	if err := enc.Encode(response); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetMeshmodelComponentByCategory(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	cat := mux.Vars(r)["category"]
	queryParams := r.URL.Query()
	v := queryParams.Get("version")
	returnAnnotationComp := queryParams.Get("annotations")
	filter := &regv1beta1.ComponentFilter{
		CategoryName: cat,
		Version:      v,
		Trim:         queryParams.Get("trim") == queryParamTrue,
		APIVersion:   queryParams.Get("apiVersion"),
		Limit:        limit,
		Offset:       offset,
		OrderOn:      order,
		Sort:         sort,
		Annotations:  returnAnnotationComp,
	}
	if search != "" {
		filter.Greedy = true
		filter.DisplayName = search
	}
	entities, count, _, _ := h.registryManager.GetEntities(filter)
	comps := processComponentDefinitions(entities)

	var pgSize int64
	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	response := models.MeshmodelComponentsDuplicateAPIResponse{
		Page:       page,
		PageSize:   int(pgSize),
		Count:      count,
		Components: models.FindDuplicateComponents(comps),
	}

	if err := enc.Encode(response); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func (h *Handler) GetAllMeshmodelComponents(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	queryParams := r.URL.Query()
	v := queryParams.Get("version")
	returnAnnotationComp := queryParams.Get("annotations")
	filter := &regv1beta1.ComponentFilter{
		Id:          queryParams.Get("id"),
		Version:     v,
		Trim:        queryParams.Get("trim") == queryParamTrue,
		APIVersion:  queryParams.Get("apiVersion"),
		Limit:       limit,
		Offset:      offset,
		OrderOn:     order,
		Sort:        sort,
		Annotations: returnAnnotationComp,
	}
	if search != "" {
		filter.Greedy = true
		filter.DisplayName = search
	}
	entities, count, _, _ := h.registryManager.GetEntities(filter)
	comps := processComponentDefinitions(entities)

	var pgSize int64

	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}

	res := models.MeshmodelComponentsDuplicateAPIResponse{
		Page:       page,
		PageSize:   int(pgSize),
		Count:      count,
		Components: models.FindDuplicateComponents(comps),
	}

	if err := enc.Encode(res); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

// request body should be json
// request body should be of ComponentCapability format
func (h *Handler) RegisterMeshmodelComponents(rw http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var cc registry.MeshModelRegistrantData
	err := dec.Decode(&cc)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var c component.ComponentDefinition
	switch cc.EntityType {
	case entity.ComponentDefinition:
		var isModelError bool
		var isRegistranError bool
		err = json.Unmarshal(cc.Entity, &c)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		utils.WriteSVGsOnFileSystem(&c)
		isRegistranError, isModelError, err = h.registryManager.RegisterEntity(cc.Connection, &c)
		helpers.HandleError(cc.Connection, &c, err, isModelError, isRegistranError)
	}
	err = helpers.WriteLogsToFiles()
	if err != nil {
		h.log.Error(err)
	}
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	go h.config.MeshModelSummaryChannel.Publish()
}

func (h *Handler) GetMeshmodelRegistrants(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)

	filter := &_models.HostFilter{
		Limit:   limit,
		Offset:  offset,
		Sort:    sort,
		OrderOn: order,
	}
	if search != "" {
		filter.Greedy = true
		filter.DisplayName = search
	}
	hosts, count, err := h.registryManager.GetRegistrants(filter)
	if err != nil {
		h.log.Error(ErrGetMeshModels(err))
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
		return
	}

	var pgSize int64

	if limit == 0 {
		pgSize = count
	} else {
		pgSize = int64(limit)
	}
	res := models.MeshmodelRegistrantsAPIResponse{
		Page:        page,
		PageSize:    int(pgSize),
		Count:       count,
		Registrants: hosts,
	}

	if err := enc.Encode(res); err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

// request body should be json
// request body should be of struct containing ID and Status fields
func (h *Handler) UpdateEntityStatus(rw http.ResponseWriter, r *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	dec := json.NewDecoder(r.Body)
	userID := user.ID
	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
		return
	}
	entityType := mux.Vars(r)["entityType"]
	var updateData struct {
		ID          string `json:"id"`
		Status      string `json:"status"`
		DisplayName string `json:"displayname"`
	}
	err = dec.Decode(&updateData)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(userID).FromUser(userID).FromSystem(*h.SystemID).WithCategory(entityType).WithAction("update")
	err = h.registryManager.UpdateEntityStatus(updateData.ID, updateData.Status, entityType)
	if err != nil {
		eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Failed to update '%s' status to %s", updateData.DisplayName, updateData.Status)).WithMetadata(map[string]interface{}{
			"error": err,
		})
		_event := eventBuilder.Build()
		_ = provider.PersistEvent(*_event, token)
		go h.config.EventBroadcaster.Publish(userID, _event)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	description := fmt.Sprintf("Status of '%s' updated to %s.", updateData.DisplayName, updateData.Status)

	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).Build()
	_ = provider.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)

	// Respond with success status
	rw.WriteHeader(http.StatusNoContent)
}

// processComponentDefinitions processes a list of entities and extracts component definitions,
// it also sets the ModelReference field for each component definition.
func processComponentDefinitions(entities []entity.Entity) []component.ComponentDefinition {
	var comps []component.ComponentDefinition
	for _, r := range entities {
		comp, ok := r.(*component.ComponentDefinition)
		if ok {
			if comp.Model != nil {
				comp.ModelReference = comp.Model.ToReference()
			}

			comps = append(comps, *comp)
		}

	}
	return comps
}

// request content byte in form value and header of the type in form

func (h *Handler) RegisterMeshmodels(rw http.ResponseWriter, r *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	var response models.RegistryAPIResponse
	regErrorStore := models.NewRegistrationFailureLogHandler()
	var tempFile *os.File
	var mu sync.Mutex
	defer func() {
		_ = r.Body.Close()
	}()
	userID := user.ID
	var message string

	token, err := provider.GetProviderToken(r)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		http.Error(rw, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
		return
	}

	//Here the codes handles to decode and store the data from the payload
	var importRequest schemav1beta1.ImportRequest

	err = json.NewDecoder(r.Body).Decode(&importRequest)
	if err != nil {
		h.log.Info("Error in unmarshalling request body")
		h.sendErrorEvent(userID, provider, "Error in unmarshalling request body", err, token)
		http.Error(rw, "Invalid request format", http.StatusBadRequest)
		return
	}

	registrationHelper := registration.NewRegistrationHelper(
		utils.UI,
		h.registryManager,
		regErrorStore,
	)
	var dir registration.Dir
	switch importRequest.UploadType {
	case "csv":
		err := meshkitRegistryUtils.SetLogger(false)
		if err != nil {
			h.handleError(rw, err, "Error setting logger")
			h.sendErrorEvent(userID, provider, "Error setting logger", err, token)
		}
		fetchBase64DataFromDataURL := func(dataURL string) ([]byte, error) {
			if strings.HasPrefix(dataURL, "data:text/csv;base64,") {
				base64Data := strings.TrimPrefix(dataURL, "data:text/csv;base64,")
				return base64.StdEncoding.DecodeString(base64Data)
			}
			return nil, ErrFileType("file is not of type csv")
		}
		modelCSVData, err := fetchBase64DataFromDataURL(importRequest.ImportBody.ModelCsv)
		if err != nil {
			h.handleError(rw, err, "Error fetching or decoding Model CSV")
			h.sendErrorEvent(userID, provider, "Error fetching or decoding Model CSV", err, token)
			return
		}
		modelCsvFile, err := os.CreateTemp("", "model-*.csv")
		if err != nil {
			err = ErrCreateFile(err, "Error creating temp file for Model CSV")
			h.handleError(rw, err, "Error creating temp file for Model CSV")
			h.sendErrorEvent(userID, provider, "Error creating temp file for Model CSV", err, token)
			return
		}
		defer func() {
			if err := modelCsvFile.Close(); err != nil {
				h.log.Error(err)
			}
		}()

		_, err = modelCsvFile.Write(modelCSVData)
		if err != nil {
			err = ErrWritingIntoFile(err, "Error writing Model CSV to temp file")
			h.handleError(rw, err, "Error writing Model CSV to temp file")
			h.sendErrorEvent(userID, provider, "Error writing Model CSV to temp file", err, token)
			return
		}

		componentCSVData, err := fetchBase64DataFromDataURL(importRequest.ImportBody.ComponentCsv)
		if err != nil {
			h.handleError(rw, err, "Error fetching or decoding Component CSV")
			h.sendErrorEvent(userID, provider, "Error fetching or decoding Component CSV", err, token)
			return
		}

		componentCsvFile, err := os.CreateTemp("", "component-*.csv")
		if err != nil {
			err = ErrCreateFile(err, "Error creating temp file for Component CSV")
			h.handleError(rw, err, "Error creating temp file for Component CSV")
			h.sendErrorEvent(userID, provider, "Error creating temp file for Component CSV", err, token)
			return
		}
		defer func() {
			if err := componentCsvFile.Close(); err != nil {
				h.log.Error(err)
			}
		}()

		_, err = componentCsvFile.Write(componentCSVData)
		if err != nil {
			err = ErrWritingIntoFile(err, "Error writing Component CSV to temp file")
			h.handleError(rw, err, "Error writing Component CSV to temp file")
			h.sendErrorEvent(userID, provider, "Error writing Component CSV to temp file", err, token)
			return
		}

		relationshipCSVData, err := fetchBase64DataFromDataURL(importRequest.ImportBody.RelationshipCSV)
		if err != nil {
			h.handleError(rw, err, "Error fetching or decoding Model CSV")
			h.sendErrorEvent(userID, provider, "Error fetching or decoding Model CSV", err, token)
			return
		}
		relationshipCsvFile, err := os.CreateTemp("", "relationship-*.csv")
		if err != nil {
			err = ErrCreateFile(err, "Error creating temp file for Model CSV")
			h.handleError(rw, err, "Error creating temp file for Model CSV")
			h.sendErrorEvent(userID, provider, "Error creating temp file for Model CSV", err, token)
			return
		}
		defer func() {
			if err := relationshipCsvFile.Close(); err != nil {
				h.log.Error(err)
			}
		}()

		_, err = relationshipCsvFile.Write(relationshipCSVData)
		if err != nil {
			err = ErrWritingIntoFile(err, "Error writing Model CSV to temp file")
			h.handleError(rw, err, "Error writing Model CSV to temp file")
			h.sendErrorEvent(userID, provider, "Error writing Model CSV to temp file", err, token)
			return
		}

		var wg sync.WaitGroup
		modelLocation := filepath.Join(os.Getenv("HOME"), utils.RegistryLocation)
		if _, err := os.Stat(modelLocation); os.IsNotExist(err) {
			_ = os.MkdirAll(modelLocation, 0755)
		}

		tempDir, err := os.MkdirTemp("", "tempData")
		if err != nil {
			h.handleError(rw, err, "Error creating temporary directory")
			h.sendErrorEvent(userID, provider, "Error creating temporary directory", err, token)
			return
		}
		defer func() {
			if err := os.RemoveAll(tempDir); err != nil {
				h.log.Error(err)
			}
		}()

		err = meshkitRegistryUtils.InvokeGenerationFromSheet(&wg, tempDir, 0, 0, "", "", modelCsvFile.Name(), componentCsvFile.Name(), "", relationshipCsvFile.Name(), 0, nil)
		if err != nil {
			h.handleError(rw, err, "Error invoking generation from sheet")
			h.sendErrorEvent(userID, provider, "Error invoking generation from sheet", err, token)
			return
		}

		h.sendEventForImport(userID, provider, 0, "", true, token)
		modelDirPaths, err := models.GetModelDirectoryPaths(tempDir)
		if err != nil {
			h.log.Error(models.ErrSeedingComponents(err))
		}
		if importRequest.Register {
			for _, dirPath := range modelDirPaths {
				dir := registration.NewDir(dirPath)
				registrationHelper.Register(dir)
			}
		} else {
			return
		}

		if _, err := os.Stat(modelLocation); os.IsNotExist(err) {
			_ = os.MkdirAll(modelLocation, 0755)
		}

		err = utils.CopyDirectory(tempDir, modelLocation)
		if err != nil {
			h.handleError(rw, err, "Error copying data to model location")
			h.sendErrorEvent(userID, provider, "Error copying data to model location", err, token)
			return
		}

	//Case when it is URL and them the model is generated from the URL
	case "url":

		model := &meshkitRegistryUtils.ModelCSV{
			Model:             importRequest.ImportBody.Model.Model,
			ModelDisplayName:  importRequest.ImportBody.Model.ModelDisplayName,
			PrimaryColor:      importRequest.ImportBody.Model.PrimaryColor,
			SecondaryColor:    importRequest.ImportBody.Model.SecondaryColor,
			Category:          importRequest.ImportBody.Model.Category,
			Registrant:        importRequest.ImportBody.Model.Registrant,
			Shape:             importRequest.ImportBody.Model.Shape,
			SubCategory:       importRequest.ImportBody.Model.SubCategory,
			SVGColor:          importRequest.ImportBody.Model.SvgColor,
			SVGWhite:          importRequest.ImportBody.Model.SvgWhite,
			SVGComplete:       importRequest.ImportBody.Model.SvgComplete,
			IsAnnotation:      strconv.FormatBool(importRequest.ImportBody.Model.IsAnnotation),
			PublishToRegistry: strconv.FormatBool(importRequest.ImportBody.Model.PublishToRegistry),
		}
		setDefaultValues(model)
		//Model generation strats from here
		model.Model = strings.ToLower(model.Model)

		pkg, version, err := meshkitRegistryUtils.GenerateModels(model.Registrant, importRequest.ImportBody.Url, model.Model)
		if err != nil {
			h.handleError(rw, err, "Error generating model")
			h.sendErrorEvent(userID, provider, "Error generating model", err, token)
			return
		}
		modelDirPath, compDirPath, err := utils.CreateVersionedDirectoryForModelAndComp(version, model.Model)
		if err != nil {
			h.handleError(rw, err, "Error decoding JSON into ModelCSV")
			h.sendErrorEvent(userID, provider, "Error decoding JSON into ModelCSV", err, token)
			return
		}
		filePath := filepath.Join(modelDirPath, model.Model+".json")
		modelDef := model.CreateModelDefinition(version, utils.DefVersion)
		err = modelDef.WriteModelDefinition(filePath, "json")
		if err != nil {
			h.handleError(rw, err, "Error decoding JSON into ModelCSV")
			h.sendErrorEvent(userID, provider, "Error decoding JSON into ModelCSV", err, token)
			return
		}

		//Component generation starts here
		lengthofComps, _, err := meshkitRegistryUtils.GenerateComponentsFromPkg(pkg, compDirPath, utils.DefVersion, modelDef, model.Group)
		if err != nil {
			h.handleError(rw, err, "Error generating components")
			h.sendErrorEvent(userID, provider, "Error generating components", err, token)
			return
		}

		//Event when the URL is used to show that we g
		h.sendEventForImport(userID, provider, lengthofComps, model.Model, false, token)
		if importRequest.Register {
			dir = registration.NewDir(modelDirPath)
			registrationHelper.Register(dir)
		} else {
			return
		}

	case "file":
		base64Data, err := json.Marshal(importRequest.ImportBody.ModelFile)
		if err != nil {
			http.Error(rw, "Internal server error", http.StatusInternalServerError)
			return
		}
		base64String := string(base64Data)
		// Remove double quotes
		base64String = base64String[1 : len(base64String)-1]

		decodedBytes, err := base64.StdEncoding.DecodeString(base64String)
		if err != nil {
			http.Error(rw, "Invalid base64 data", http.StatusBadRequest)
			return
		}
		tempFile, err = CreateTemp(importRequest.ImportBody.FileName, decodedBytes)
		if err != nil {
			err = meshkitutils.ErrCreateFile(err, "Error creating temp file")
			h.handleError(rw, err, err.Error())
			h.sendErrorEvent(userID, provider, "Error creating temp file", err, token)
			return
		}
		defer func() {
			if err := os.Remove(tempFile.Name()); err != nil {
				h.log.Error(err)
			}
		}()

		dir = registration.NewDir(tempFile.Name())
		if importRequest.Register {
			registrationHelper.Register(dir)
			if err := tempFile.Close(); err != nil {
				h.log.Error(err)
			}
		}
	case "urlImport":
		downloadFile := func(url string) ([]byte, error) {
			resp, err := http.Get(url)
			if err != nil {
				return nil, fmt.Errorf("error downloading file from URL: %v", err)
			}
			fileData, err := io.ReadAll(resp.Body)
			if err != nil {
				if closeErr := resp.Body.Close(); closeErr != nil {
					return nil, fmt.Errorf("error reading downloaded file: %v (close error: %v)", err, closeErr)
				}
				return nil, fmt.Errorf("error reading downloaded file: %v", err)
			}

			if err := resp.Body.Close(); err != nil {
				return nil, fmt.Errorf("error closing response body: %v", err)
			}

			if resp.StatusCode != http.StatusOK {
				return nil, fmt.Errorf("failed to download file. Status code: %d", resp.StatusCode)
			}

			return fileData, nil
		}

		// Download the file from the provided URL
		fileData, err := downloadFile(importRequest.ImportBody.Url)
		if err != nil {
			h.handleError(rw, err, "Error downloading file from URL")
			h.sendErrorEvent(userID, provider, "Error downloading file from URL", err, token)
			return
		}
		isOCI := meshkitOci.IsOCIArtifact(fileData)
		fileType := ".tar"
		if !isOCI {
			fileType = detectFileType(fileData)
		}
		name := "model" + fileType
		//write the file to a temp file
		tempFile, err = CreateTemp(name, fileData)
		if err != nil {
			err = meshkitutils.ErrCreateFile(err, "Error creating temp file")
			h.handleError(rw, err, "Error creating temp file")
			h.sendErrorEvent(userID, provider, "Error creating temp file", err, token)
			return
		}
		defer func() {
			if err := os.Remove(tempFile.Name()); err != nil {
				h.log.Error(err)
			}
		}()

		dir = registration.NewDir(tempFile.Name())
		if importRequest.Register {
			registrationHelper.Register(dir)
			if err := tempFile.Close(); err != nil {
				h.log.Error(err)
			}
		}
	}

	h.handleRegistrationAndError(registrationHelper, &mu, &response, regErrorStore)
	var errMsg string
	message = writeMessageString(&response)
	if response.EntityCount.TotalErrCount > 0 {
		errMsg = ErrMsgContruct(&response)
	}

	h.sendSuccessResponse(rw, userID, provider, message, errMsg, &response, token)

}

func (h *Handler) ExportModel(rw http.ResponseWriter, r *http.Request) {
	modelId := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	version := r.URL.Query().Get("version")
	outputFormat := r.URL.Query().Get("output_format")
	fileTypes := r.URL.Query().Get("file_type")
	if fileTypes == "" {
		fileTypes = "oci"
	}
	if outputFormat == "" {
		outputFormat = "json"
	}
	hasComponents, err := strconv.ParseBool(r.URL.Query().Get("components"))
	if err != nil {
		hasComponents = true
	}

	hasRelationships, err := strconv.ParseBool(r.URL.Query().Get("relationships"))
	if err != nil {
		hasRelationships = true
	}

	// 1. Get the model data
	modelFilter := &regv1beta1.ModelFilter{
		Id:            modelId,
		Name:          name,
		Components:    hasComponents,
		Relationships: hasRelationships,
		Greedy:        true,
		Version:       version,
	}
	e, _, _, err := h.registryManager.GetEntities(modelFilter)
	if err != nil {
		h.log.Error(ErrGetMeshModels(err))
		http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusInternalServerError)
		return
	}

	// No entities found for the given filter
	// This can happen if the model is not found or if the model is not registered
	// or if the model is not in the registry
	// In this case, we return a 404 error
	if len(e) == 0 {
		message := "model with "
		if modelId != "" {
			message += fmt.Sprintf("id %s ", modelId)
		}
		if name != "" {
			message += fmt.Sprintf("name %s ", name)
		}
		if version != "" {
			message += fmt.Sprintf("version %s ", version)
		}
		message += "has not been found"
		// h.log.Error(ErrGetMeshModels(err))
		// http.Error(rw, ErrGetMeshModels(err).Error(), http.StatusNotFound)
		rw.WriteHeader(http.StatusNotFound)
		// rw.Write([]byte(message))
		if _, err := fmt.Fprintln(rw, message); err != nil {
			h.log.Error(err)
		}
		return
	}

	model := e[0].(*_model.ModelDefinition)
	//This path is used to so that the function can be aware of where the svg file is
	//This is for relative path as we are inside meshery/server/cmd/main.go
	//File is stored in Ui folder so we need to move 2 directories back
	//We do this because we use this in mesheryctl as well
	err = model.ReplaceSVGData("../../")
	if err != nil {
		h.log.Error(err)
	}
	// 2. Convert it to oci
	temp := os.TempDir()
	modelDir := filepath.Join(temp, model.Name)
	versionDir := filepath.Join(modelDir, model.Model.Version, model.Version)

	// Create necessary directories: {modelName}/v1.0.0/1.0.0/{components, relationships}
	dirs := []string{
		versionDir,
		filepath.Join(versionDir, "components"),
		filepath.Join(versionDir, "relationships"),
	}

	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			err = meshkitutils.ErrCreateDir(err, "Error creating temp directory")
			h.log.Error(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	defer func() {
		if err := os.RemoveAll(modelDir); err != nil {
			h.log.Error(err)
		}
	}()

	components := []component.ComponentDefinition{}
	// Components can be nil if hasComponents is false
	if model.Components != nil {
		components = model.Components.([]component.ComponentDefinition)
	}

	relationships := []relationship.RelationshipDefinition{}
	// Relationships can be nil if hasRelationships is false
	if model.Relationships != nil {
		relationships = model.Relationships.([]relationship.RelationshipDefinition)
	}

	model.Relationships = nil
	model.Components = nil
	// Write model.json to {modelname}/v1.0.0/1.0.0/model.json
	err = model.WriteModelDefinition(filepath.Join(versionDir, fmt.Sprintf("model.%s", outputFormat)), outputFormat)
	if err != nil {
		h.log.Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	componentsDir := filepath.Join(versionDir, "components")
	relationshipsDir := filepath.Join(versionDir, "relationships")

	for _, comp := range components {
		_ = comp.ReplaceSVGData("../../")
		comp.Model = model
		_, err := comp.WriteComponentDefinition(componentsDir, outputFormat)
		if err != nil {
			h.log.Error(err)
		}

	}
	for _, rel := range relationships {
		rel.Model = model.ToReference()
		err := rel.WriteRelationshipDefinition(relationshipsDir, outputFormat)
		if err != nil {
			h.log.Error(err)
		}

	}

	// Write components into {modelname}/v1.0.0/1.0.0/components

	// At this point, the data has been written to the directories:
	// {modelname}/v1.0.0/1.0.0/model.json
	// {modelname}/v1.0.0/1.0.0/components/*.json
	// {modelname}/v1.0.0/1.0.0/relationships/*.json

	// Build OCI image for the model from the modelDir
	var tarfileName string
	var byt []byte
	if fileTypes == "oci" {
		img, err := meshkitOci.BuildImage(modelDir)
		if err != nil {
			h.log.Error(err) // TODO: Add appropriate meshkit error
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// Save OCI artifact into a tar file
		tarfileName := filepath.Join(modelDir, "model.tar")
		err = meshkitOci.SaveOCIArtifact(img, tarfileName, model.Name)
		if err != nil {
			h.log.Error(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. Send response
		byt, _ = os.ReadFile(tarfileName)
		rw.Header().Add("Content-Type", "application/x-tar")
		rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.tar\"", model.Name))
		rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(byt)))
	} else {
		var tarData bytes.Buffer
		err := meshkitutils.Compress(modelDir, &tarData)
		if err != nil {
			h.log.Error(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		tarfileName = filepath.Join(modelDir, "model.tar.gz")
		err = os.WriteFile(tarfileName, tarData.Bytes(), 0644)
		if err != nil {
			h.log.Error(err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		byt, _ = os.ReadFile(tarfileName)
		rw.Header().Add("Content-Type", "application/gzip")
		rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.tar.gz\"", model.Name))

	}

	// 3. Send response
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(byt)))
	_, err = rw.Write(byt)
	if err != nil {
		if isClientDisconnect(err) {
			h.log.Debug(ErrEncodeResponse(err))
		} else {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

func RegisterEntity(content []byte, entityType entity.EntityType, h *Handler) error {
	switch entityType {
	case entity.ComponentDefinition:
		var c component.ComponentDefinition
		err := json.Unmarshal(content, &c)
		if err != nil {
			return meshkitutils.ErrUnmarshal(err)
		}
		isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(connection.Connection{
			Kind: c.Model.Registrant.Kind,
		}, &c)
		helpers.HandleError(connection.Connection{
			Kind: c.Model.Registrant.Kind,
		}, &c, err, isModelError, isRegistrantError)
		return nil
	case entity.RelationshipDefinition:
		var r relationship.RelationshipDefinition
		err := json.Unmarshal(content, &r)
		if err != nil {
			return meshkitutils.ErrUnmarshal(err)
		}
		isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(connection.Connection{
			Kind: r.Model.Registrant.Kind,
		}, &r)
		helpers.HandleError(connection.Connection{
			Kind: r.Model.Registrant.Kind,
		}, &r, err, isModelError, isRegistrantError)
		return nil
	}
	return meshkitutils.ErrInvalidSchemaVersion
}

func (h *Handler) DeleteModel(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	modelID := mux.Vars(r)["id"]
	modelUUID, err := uuid.FromString(modelID)
	if err != nil {
		http.Error(rw, ErrInvalidUUID(err).Error(), http.StatusBadRequest)
		return
	}

	err = h.dbHandler.Transaction(func(tx *gorm.DB) error {
		var modelDef _model.ModelDefinition
		if err := tx.First(&modelDef, "id = ?", modelUUID).Error; err != nil {
			return err
		}

		// Delete registry entries for components belonging to this model
		if err := tx.Where("entity IN (?) AND type = ?",
			tx.Model(&component.ComponentDefinition{}).Select("id").Where("model_id = ?", modelUUID),
			entity.ComponentDefinition,
		).Delete(&registry.Registry{}).Error; err != nil {
			return err
		}

		// Delete registry entries for relationships belonging to this model
		if err := tx.Where("entity IN (?) AND type = ?",
			tx.Model(&relationship.RelationshipDefinition{}).Select("id").Where("model_id = ?", modelUUID),
			entity.RelationshipDefinition,
		).Delete(&registry.Registry{}).Error; err != nil {
			return err
		}

		// Delete registry entries for policies belonging to this model
		if err := tx.Where("entity IN (?) AND type = ?",
			tx.Model(&_models.PolicyDefinition{}).Select("id").Where("modelID = ?", modelUUID),
			entity.PolicyDefinition,
		).Delete(&registry.Registry{}).Error; err != nil {
			return err
		}

		// Delete the model's own registry entry
		if err := tx.Where("entity = ? AND type = ?", modelUUID, entity.Model).Delete(&registry.Registry{}).Error; err != nil {
			return err
		}

		// Delete components, relationships, and policies
		if err := tx.Where("model_id = ?", modelUUID).Delete(&component.ComponentDefinition{}).Error; err != nil {
			return err
		}
		if err := tx.Where("model_id = ?", modelUUID).Delete(&relationship.RelationshipDefinition{}).Error; err != nil {
			return err
		}
		if err := tx.Where("modelID = ?", modelUUID).Delete(&_models.PolicyDefinition{}).Error; err != nil {
			return err
		}

		// Delete the model itself
		return tx.Where("id = ?", modelUUID).Delete(&_model.ModelDefinition{}).Error
	})

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(rw, fmt.Sprintf("model with id %s not found", modelID), http.StatusNotFound)
			return
		}
		mesheryErr := models.ErrDBDelete(err, "")
		h.log.Error(mesheryErr)
		http.Error(rw, mesheryErr.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
