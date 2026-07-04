package handlers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
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
	"github.com/meshery/schemas/models/v1beta1/connection"
	_model "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta3/component"

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
		Page:       page,
		PageSize:   int(pgSize),
		TotalCount: count,
		Models:     models.FindDuplicateModels(modelDefs),
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
		Page:       page,
		PageSize:   int(pgSize),
		TotalCount: count,
		Models:     models.FindDuplicateModels(modelDefs),
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
		Page:       page,
		PageSize:   int(pgSize),
		TotalCount: count,
		Models:     models.FindDuplicateModels(modelDefs),
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
		Page:       page,
		PageSize:   int(pgSize),
		TotalCount: count,
		Models:     models.FindDuplicateModels(modelDefs),
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
		TotalCount: count,
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
		TotalCount: count,
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
		TotalCount: count,
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
		TotalCount: count,
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
		TotalCount: count,
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
		TotalCount: count,
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
		TotalCount: count,
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
		TotalCount: count,
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
		TotalCount: count,
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
		TotalCount: count,
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
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	var c component.ComponentDefinition
	switch cc.EntityType {
	case entity.ComponentDefinition:
		var isModelError bool
		var isRegistranError bool
		err = json.Unmarshal(cc.Entity, &c)
		if err != nil {
			h.log.Error(models.ErrUnmarshal(err, "component definition"))
			writeMeshkitError(rw, models.ErrUnmarshal(err, "component definition"), http.StatusBadRequest)
			return
		}
		utils.WriteSVGsOnFileSystem(&c)
		isRegistranError, isModelError, err = h.registryManager.RegisterEntity(cc.Connection, &c)
		helpers.HandleError(registry.RegistrantHostToV1beta1(cc.Connection), &c, err, isModelError, isRegistranError)
	}
	err = helpers.WriteLogsToFiles()
	if err != nil {
		// WriteLogsToFiles is an internal flush of registry-attempt
		// state to REGISTRY_LOG_FILE — the failure is server-side
		// (filesystem permissions, disk full, marshal error), so
		// surface a 500 with structured remediation instead of the
		// previous raw 400.
		wrappedErr := ErrWriteRegistryLogs(err)
		h.log.Error(wrappedErr)
		writeMeshkitError(rw, wrappedErr, http.StatusInternalServerError)
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
		writeMeshkitError(rw, ErrGetMeshModels(err), http.StatusInternalServerError)
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
		TotalCount:  count,
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
		writeMeshkitError(rw, ErrRetrieveUserToken(err), http.StatusInternalServerError)
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
		writeMeshkitError(rw, ErrRequestBody(err), http.StatusBadRequest)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(userID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory(entityType).WithAction("update")
	err = h.registryManager.UpdateEntityStatus(updateData.ID, updateData.Status, entityType)
	if err != nil {
		wrappedErr := ErrUpdateEntityStatus(err)
		eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Failed to update '%s' status to %s", updateData.DisplayName, updateData.Status)).WithMetadata(map[string]interface{}{
			"error": wrappedErr,
		})
		_event := eventBuilder.Build()
		_ = provider.PersistEvent(*_event, token)
		go h.config.EventBroadcaster.Publish(userID, _event)
		writeMeshkitError(rw, wrappedErr, http.StatusInternalServerError)
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
		writeMeshkitError(rw, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}

	//Here the codes handles to decode and store the data from the payload
	var importRequest schemav1beta1.ImportRequest

	err = json.NewDecoder(r.Body).Decode(&importRequest)
	if err != nil {
		h.log.Info("Error in unmarshalling request body")
		h.sendErrorEvent(userID, provider, "Error in unmarshalling request body", err, token)
		writeMeshkitError(rw, models.ErrUnmarshal(err, "import request"), http.StatusBadRequest)
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
			h.log.Error(models.ErrMarshal(err, "model file"))
			writeMeshkitError(rw, models.ErrMarshal(err, "model file"), http.StatusInternalServerError)
			return
		}
		base64String := string(base64Data)
		// Remove double quotes
		base64String = base64String[1 : len(base64String)-1]

		decodedBytes, err := base64.StdEncoding.DecodeString(base64String)
		if err != nil {
			h.log.Error(fmt.Errorf("invalid base64 data: %w", err))
			writeMeshkitError(rw, ErrInvalidBase64Data(err), http.StatusBadRequest)
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
		writeMeshkitError(rw, ErrGetMeshModels(err), http.StatusInternalServerError)
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
		writeJSONError(rw, message, http.StatusNotFound)
		return
	}

	model, ok := e[0].(*_model.ModelDefinition)
	if !ok {
		writeJSONError(rw, "invalid model definition type", http.StatusInternalServerError)
		return
	}
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
			writeMeshkitError(rw, ErrExportModel(err, "temp directory creation"), http.StatusInternalServerError)
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
		writeMeshkitError(rw, ErrExportModel(err, "model definition write"), http.StatusInternalServerError)
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
			h.log.Error(err)
			writeMeshkitError(rw, ErrExportModel(err, "OCI image build"), http.StatusInternalServerError)
			return
		}

		// Save OCI artifact into a tar file
		tarfileName := filepath.Join(modelDir, "model.tar")
		err = meshkitOci.SaveOCIArtifact(img, tarfileName, model.Name)
		if err != nil {
			h.log.Error(err)
			writeMeshkitError(rw, ErrExportModel(err, "OCI artifact save"), http.StatusInternalServerError)
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
			writeMeshkitError(rw, ErrExportModel(err, "tar.gz compress"), http.StatusInternalServerError)
			return
		}
		tarfileName = filepath.Join(modelDir, "model.tar.gz")
		err = os.WriteFile(tarfileName, tarData.Bytes(), 0644)
		if err != nil {
			h.log.Error(err)
			writeMeshkitError(rw, ErrExportModel(err, "tar.gz write"), http.StatusInternalServerError)
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

func (h *Handler) PushModel(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	var req struct {
		ModelID      string `json:"modelId"`
		ModelName    string `json:"modelName"`
		Version      string `json:"version"`
		Registry     string `json:"registry"`
		Repository   string `json:"repository"`
		Tag          string `json:"tag"`
		CredentialID string `json:"credentialId"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(rw, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.ModelID == "" && req.ModelName == "" {
		writeJSONError(rw, "modelId or modelName is required", http.StatusBadRequest)
		return
	}

	username, password, err := resolveCredentials(r, &req.CredentialID, req.Username, req.Password, provider)
	// Clear inline credentials from the request struct to minimise
	// accidental exposure lifetime. This is not secure memory erasure;
	// strings in Go are immutable and may linger in memory until GC.
	req.Username = ""
	req.Password = ""
	if err != nil {
		writeMeshkitError(rw, ErrPushModel(err, "credential resolution"), http.StatusBadRequest)
		return
	}

	if err := validateOCIRegistryDestination(req.Registry, req.Repository); err != nil {
		writeJSONError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Tag == "" {
		req.Tag = "latest"
	}

	modelFilter := &regv1beta1.ModelFilter{
		Id:            req.ModelID,
		Name:          req.ModelName,
		Components:    true,
		Relationships: true,
		Greedy:        true,
		Version:       req.Version,
	}
	e, _, _, err := h.registryManager.GetEntities(modelFilter)
	if err != nil {
		h.log.Error(ErrGetMeshModels(err))
		writeMeshkitError(rw, ErrGetMeshModels(err), http.StatusInternalServerError)
		return
	}
	if len(e) == 0 {
		writeJSONError(rw, "model not found", http.StatusNotFound)
		return
	}

	model, ok := e[0].(*_model.ModelDefinition)
	if !ok {
		writeJSONError(rw, "invalid model definition type", http.StatusInternalServerError)
		return
	}

	modelDir, err := os.MkdirTemp("", "model-push-")
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(rw, ErrPushModel(err, "temp directory creation"), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := os.RemoveAll(modelDir); err != nil {
			h.log.Error(err)
		}
	}()

	versionDir := filepath.Join(modelDir, model.Model.Version, model.Version)
	for _, dir := range []string{versionDir, filepath.Join(versionDir, "components"), filepath.Join(versionDir, "relationships")} {
		if err := os.MkdirAll(dir, 0700); err != nil {
			h.log.Error(err)
			writeMeshkitError(rw, ErrPushModel(err, "temp directory creation"), http.StatusInternalServerError)
			return
		}
	}

	var components []component.ComponentDefinition
	if model.Components != nil {
		var ok bool
		components, ok = model.Components.([]component.ComponentDefinition)
		if !ok {
			err := fmt.Errorf("unexpected type for Components: %T", model.Components)
			h.log.Error(err)
			writeMeshkitError(rw, ErrPushModel(err, "type assertion"), http.StatusInternalServerError)
			return
		}
	}
	var relationships []relationship.RelationshipDefinition
	if model.Relationships != nil {
		var ok bool
		relationships, ok = model.Relationships.([]relationship.RelationshipDefinition)
		if !ok {
			err := fmt.Errorf("unexpected type for Relationships: %T", model.Relationships)
			h.log.Error(err)
			writeMeshkitError(rw, ErrPushModel(err, "type assertion"), http.StatusInternalServerError)
			return
		}
	}

	_ = model.ReplaceSVGData("../../")
	model.Relationships = nil
	model.Components = nil

	if err := model.WriteModelDefinition(filepath.Join(versionDir, "model.json"), "json"); err != nil {
		h.log.Error(err)
		writeMeshkitError(rw, ErrPushModel(err, "model definition write"), http.StatusInternalServerError)
		return
	}
	for _, comp := range components {
		_ = comp.ReplaceSVGData("../../")
		comp.Model = model
		if _, err := comp.WriteComponentDefinition(filepath.Join(versionDir, "components"), "json"); err != nil {
			h.log.Error(err)
			writeMeshkitError(rw, ErrPushModel(err, "component definition write"), http.StatusInternalServerError)
			return
		}
	}
	for _, rel := range relationships {
		rel.Model = model.ToReference()
		if err := rel.WriteRelationshipDefinition(filepath.Join(versionDir, "relationships"), "json"); err != nil {
			h.log.Error(err)
			writeMeshkitError(rw, ErrPushModel(err, "relationship definition write"), http.StatusInternalServerError)
			return
		}
	}

	if err := meshkitOci.PushToOCIRegistry(modelDir, req.Registry, req.Repository, req.Tag, username, password); err != nil {
		h.log.Error(err)
		writeMeshkitError(rw, ErrPushModel(err, "OCI push"), http.StatusInternalServerError)
		return
	}

	writeJSONMessage(rw, map[string]string{"message": "Model pushed successfully"}, http.StatusOK)
}

func (h *Handler) PullModel(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	var req struct {
		Registry     string `json:"registry"`
		Repository   string `json:"repository"`
		Tag          string `json:"tag"`
		CredentialID string `json:"credentialId"`
		Username     string `json:"username"`
		Password     string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(rw, "Invalid request body", http.StatusBadRequest)
		return
	}

	username, password, err := resolveCredentials(r, &req.CredentialID, req.Username, req.Password, provider)
	req.Username = ""
	req.Password = ""
	if err != nil {
		writeMeshkitError(rw, ErrPullModel(err, "credential resolution"), http.StatusBadRequest)
		return
	}

	if err := validateOCIRegistryDestination(req.Registry, req.Repository); err != nil {
		writeJSONError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Tag == "" {
		req.Tag = "latest"
	}

	pullDir, err := os.MkdirTemp("", "model-pull-")
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(rw, ErrPullModel(err, "pull directory creation"), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := os.RemoveAll(pullDir); err != nil {
			h.log.Error(err)
		}
	}()

	if err := meshkitOci.PullFromOCIRegistry(pullDir, req.Registry, req.Repository, req.Tag, username, password); err != nil {
		h.log.Error(err)
		writeMeshkitError(rw, ErrPullModel(err, "OCI pull"), http.StatusInternalServerError)
		return
	}

	extractDir, err := os.MkdirTemp("", "model-extract-")
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(rw, ErrPullModel(err, "extract directory creation"), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := os.RemoveAll(extractDir); err != nil {
			h.log.Error(err)
		}
	}()

	if err := extractFromOCIStore(pullDir, req.Tag, extractDir); err != nil {
		h.log.Error(err)
		writeMeshkitError(rw, ErrPullModel(err, "artifact extraction"), http.StatusInternalServerError)
		return
	}

	regErrorStore := models.NewRegistrationFailureLogHandler()
	registrationHelper := registration.NewRegistrationHelper(utils.UI, h.registryManager, regErrorStore)
	dir := registration.NewDir(extractDir)
	registrationHelper.Register(dir)

	if len(registrationHelper.PkgUnits) == 0 {
		h.log.Error(ErrPullModel(fmt.Errorf("registration produced no valid package units"), "registration"))
		writeMeshkitError(rw, ErrPullModel(fmt.Errorf("registration produced no valid package units"), "registration"), http.StatusInternalServerError)
		return
	}

	writeJSONMessage(rw, map[string]string{"message": "Model pulled and registered successfully"}, http.StatusOK)
}

// extractFromOCIStore reads an ORAS file store (created by PullFromOCIRegistry),
// finds its first layer blob, decompresses and untars it into destDir.
// It validates digests before reading files, streams layer content instead of
// reading the entire blob into memory, rejects symlinks/hardlinks, enforces
// resource limits, and prevents path traversal.
func extractFromOCIStore(pullDir, tag, destDir string) error {
	const (
		maxTotalBytes = 500 << 20 // 500 MiB total decompressed
		maxFileSize   = 100 << 20 // 100 MiB per file
		maxFileCount  = 10000
	)

	indexFile := filepath.Join(pullDir, "index.json")
	data, err := os.ReadFile(indexFile)
	if err != nil {
		return fmt.Errorf("cannot read index.json: %w", err)
	}

	var index struct {
		Manifests []struct {
			Digest      string            `json:"digest"`
			Annotations map[string]string `json:"annotations,omitempty"`
		} `json:"manifests"`
	}
	if err := json.Unmarshal(data, &index); err != nil {
		return fmt.Errorf("cannot parse index.json: %w", err)
	}
	if len(index.Manifests) == 0 {
		return fmt.Errorf("no manifests found in pulled artifact")
	}

	manifestDigest := index.Manifests[0].Digest
	for _, m := range index.Manifests {
		if ref, ok := m.Annotations["org.opencontainers.image.ref.name"]; ok && ref == tag {
			manifestDigest = m.Digest
			break
		}
	}

	if err := validateOCIDigest(manifestDigest); err != nil {
		return fmt.Errorf("invalid manifest digest: %w", err)
	}

	manifestFile := filepath.Join(pullDir, "blobs", "sha256", strings.TrimPrefix(manifestDigest, "sha256:"))
	manifestData, err := os.ReadFile(manifestFile)
	if err != nil {
		return fmt.Errorf("cannot read manifest blob %s: %w", manifestDigest, err)
	}

	var manifest struct {
		Layers []struct {
			Digest string `json:"digest"`
		} `json:"layers"`
	}
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return fmt.Errorf("cannot parse manifest: %w", err)
	}
	if len(manifest.Layers) == 0 {
		return fmt.Errorf("no layers found in pulled artifact")
	}

	layerDigest := manifest.Layers[0].Digest
	if err := validateOCIDigest(layerDigest); err != nil {
		return fmt.Errorf("invalid layer digest: %w", err)
	}

	layerFile := filepath.Join(pullDir, "blobs", "sha256", strings.TrimPrefix(layerDigest, "sha256:"))
	f, err := os.Open(layerFile)
	if err != nil {
		return fmt.Errorf("cannot open layer blob %s: %w", layerDigest, err)
	}
	defer func() { _ = f.Close() }()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("cannot decompress layer: %w", err)
	}
	defer func() { _ = gzr.Close() }()

	var totalExtracted int64
	var fileCount int
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading tar stream: %w", err)
		}

		if header.Typeflag == tar.TypeSymlink || header.Typeflag == tar.TypeLink {
			return fmt.Errorf("rejected symlink/hardlink entry: %s", header.Name)
		}

		target := filepath.Join(destDir, filepath.Clean(header.Name))
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("path traversal rejected: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			mode := os.FileMode(header.Mode) & os.ModePerm
			if mode > 0777 {
				mode = 0755
			}
			if err := os.MkdirAll(target, mode); err != nil {
				return fmt.Errorf("cannot create dir %s: %w", target, err)
			}
		case tar.TypeReg:
			if header.Size > maxFileSize {
				return fmt.Errorf("file %s exceeds max size %d", header.Name, maxFileSize)
			}
			if totalExtracted+header.Size > maxTotalBytes {
				return fmt.Errorf("total extracted size would exceed limit %d", maxTotalBytes)
			}
			if fileCount >= maxFileCount {
				return fmt.Errorf("extracted file count exceeds limit %d", maxFileCount)
			}
			fileCount++

			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return fmt.Errorf("cannot create parent dir for %s: %w", target, err)
			}
			// Restrict created file permissions; ignore header mode bits
			// that grant execute, setuid, setgid, or group/world write.
			mode := os.FileMode(header.Mode) & 0644
			fw, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
			if err != nil {
				return fmt.Errorf("cannot create file %s: %w", target, err)
			}
			written, err := io.Copy(fw, io.LimitReader(tr, header.Size))
			if cerr := fw.Close(); cerr != nil && err == nil {
				err = fmt.Errorf("cannot close file %s: %w", target, cerr)
			}
			if err != nil {
				return fmt.Errorf("cannot write file %s: %w", target, err)
			}
			if written != header.Size {
				return fmt.Errorf("short write for %s: wrote %d, expected %d", target, written, header.Size)
			}
			totalExtracted += written
		}
	}

	return nil
}

// validateOCIDigest checks that a digest string has the form "sha256:<hex>".
func validateOCIDigest(digest string) error {
	if !strings.HasPrefix(digest, "sha256:") {
		return fmt.Errorf("unsupported digest algorithm: %s", digest)
	}
	hexPart := strings.TrimPrefix(digest, "sha256:")
	if len(hexPart) != 64 {
		return fmt.Errorf("invalid digest length: %s", digest)
	}
	for _, c := range hexPart {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return fmt.Errorf("invalid hex character in digest: %s", digest)
		}
	}
	return nil
}

// resolveCredentials retrieves OCI registry credentials either from a stored
// Credential (by ID) or from inline username/password. If credentialId is set
// it takes precedence over inline values.
func resolveCredentials(r *http.Request, credentialID *string, inlineUsername, inlinePassword string, provider models.Provider) (username, password string, err error) {
	if credentialID != nil && *credentialID != "" {
		token, ok := r.Context().Value(models.TokenCtxKey).(string)
		if !ok {
			return "", "", fmt.Errorf("failed to retrieve auth token from context")
		}
		cid, parseErr := uuid.FromString(*credentialID)
		if parseErr != nil {
			return "", "", fmt.Errorf("invalid credentialId: %w", parseErr)
		}
		cred, _, cerr := provider.GetCredentialByID(token, cid)
		if cerr != nil {
			return "", "", fmt.Errorf("failed to retrieve credential: %w", cerr)
		}
		if cred == nil {
			return "", "", fmt.Errorf("credential not found")
		}
		u, _ := cred.Secret["username"].(string)
		p, _ := cred.Secret["password"].(string)
		if u == "" {
			return "", "", fmt.Errorf("credential %q has no username", *credentialID)
		}
		return u, p, nil
	}
	return inlineUsername, inlinePassword, nil
}

// validateOCIRegistryDestination validates the registry address and repository
// for OCI operations. It rejects private/loopback/link-local destinations to
// prevent SSRF attacks.
func validateOCIRegistryDestination(registry, repository string) error {
	if registry == "" {
		return fmt.Errorf("registry is required")
	}
	if repository == "" {
		return fmt.Errorf("repository is required")
	}

	// Strip optional scheme for parsing; only https/http are potentially
	// acceptable, but we strip them to focus on the host.
	raw := registry
	if strings.Contains(raw, "://") {
		u, perr := url.Parse(raw)
		if perr != nil {
			return fmt.Errorf("invalid registry URL: %w", perr)
		}
		if u.Scheme != "" && u.Scheme != "https" && u.Scheme != "http" {
			return fmt.Errorf("unsupported registry scheme: %s", u.Scheme)
		}
		raw = u.Host
	}

	// Extract host (and optional port)
	host := raw
	if h, _, err := net.SplitHostPort(raw); err == nil {
		host = h
	}

	if host == "" {
		return fmt.Errorf("registry host is empty")
	}

	// Reject path separators in host (must be pure host:port)
	if strings.ContainsAny(raw, "/\\") {
		return fmt.Errorf("registry must be a host:port, not a path")
	}

	// Reject repository path traversal
	if strings.Contains(repository, "..") {
		return fmt.Errorf("repository must not contain path traversal")
	}

	// Reject known-bad host patterns including unspecified addresses.
	lower := strings.ToLower(host)
	if lower == "localhost" || lower == "127.0.0.1" || lower == "::1" {
		return fmt.Errorf("registry must not be a loopback address")
	}
	if lower == "0.0.0.0" || lower == "::" {
		return fmt.Errorf("registry must not be an unspecified address")
	}

	// Resolve hostname and check IP ranges.
	// NOTE: DNS-rebinding TOCTOU — the resolution below happens at
	// SSRF-check time, not at connection time. The meshkit OCI transport
	// (PushToOCIRegistry / PullFromOCIRegistry) does not expose a dial
	// control function, so the resolved addresses cannot be re-checked
	// on the actual TCP connection. This is an upstream limitation of
	// github.com/meshery/meshkit v1.0.19. The correct fix would require
	// meshkit to accept a net.Dialer or dial-context function.
	ips, err := net.LookupHost(host)
	if err != nil {
		return fmt.Errorf("registry host %q could not be resolved: %w", host, err)
	}
	for _, ipStr := range ips {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			continue
		}
		if ip.IsLoopback() {
			return fmt.Errorf("registry resolved to loopback address: %s", ipStr)
		}
		if ip.IsUnspecified() {
			return fmt.Errorf("registry resolved to unspecified address: %s", ipStr)
		}
		if ip.IsPrivate() {
			return fmt.Errorf("registry resolved to private address: %s", ipStr)
		}
		if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return fmt.Errorf("registry resolved to link-local address: %s", ipStr)
		}
	}

	return nil
}

func RegisterEntity(content []byte, entityType entity.EntityType, h *Handler) error {
	switch entityType {
	case entity.ComponentDefinition:
		var c component.ComponentDefinition
		err := json.Unmarshal(content, &c)
		if err != nil {
			return meshkitutils.ErrUnmarshal(err)
		}
		isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(registry.RegistrantHostToV1beta3(connection.Connection{
			Kind: c.Model.Registrant.Kind,
		}), &c)
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
		isRegistrantError, isModelError, err := h.registryManager.RegisterEntity(registry.RegistrantHostToV1beta3(connection.Connection{
			Kind: r.Model.Registrant.Kind,
		}), &r)
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
		h.log.Error(models.ErrInvalidUUID(err))
		writeMeshkitError(rw, models.ErrInvalidUUID(err), http.StatusBadRequest)
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
			writeJSONError(rw, fmt.Sprintf("model with id %s not found", modelID), http.StatusNotFound)
			return
		}
		mesheryErr := models.ErrDBDelete(err, "")
		h.log.Error(mesheryErr)
		writeMeshkitError(rw, mesheryErr, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
