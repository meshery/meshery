package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshsync/pkg/model"
	"gorm.io/gorm/clause"
)

// swagger:route GET /api/system/meshsync/resources GetMeshSyncResources idGetMeshSyncResources
// Handle GET request for meshsync discovered resources
//
// ```?apiVersion={apiVersion}``` If apiVersion is unspecified then all models are returned
//
// ```?search={componentname}``` If search is non empty then a greedy search is performed
//
// ```?order={field}``` orders on the passed field
//
// ```?sort={[asc/desc]}``` Default behavior is asc
//
// ```?page={page-number}``` Default page number is 1
//
// ```?pagesize={pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
//
// ```?annotation={annotaion}``` annotation is a boolean value. If true then annotations are returned
//
// ```?labels={labels}``` labels is a boolean value. If true then labels are returned
//
// ```?spec={spec}``` spec is a boolean value. If true then spec is returned
//
// ```?status={status}``` status is a boolean value. If true then status is returned
//
// ```?clusterId={[clusterId]}``` clusterId is array of string values. Required.
//
// responses:
// 200: []meshsyncResourcesResponseWrapper

func (h *Handler) GetMeshSyncResources(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	rw.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(rw)

	var resources []model.KubernetesResource
	var totalCount int64
	limitstr := r.URL.Query().Get("pagesize")
	var limit int
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit <= 0 {
			limit = defaultPageSize
		}
	}

	pagestr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pagestr)

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")
	search := r.URL.Query().Get("search")
	apiVersion := r.URL.Query().Get("apiVersion")
	spec, _ := strconv.ParseBool(r.URL.Query().Get("spec"))
	status, _ := strconv.ParseBool(r.URL.Query().Get("status"))
	isAnnotaion, _ := strconv.ParseBool(r.URL.Query().Get("annotations"))
	isLabels, _ := strconv.ParseBool(r.URL.Query().Get("labels"))
	kind := r.URL.Query().Get("kind")

	filter := struct {
		ClusterIds []string `json:"clusterIds"`
	}{}

	clusterIds := r.URL.Query().Get("clusterIds")
	if clusterIds != "" {
		err := json.Unmarshal([]byte(clusterIds), &filter.ClusterIds)
		if err != nil {
			h.log.Error(ErrFetchMeshSyncResources(err))
			http.Error(rw, ErrFetchMeshSyncResources(err).Error(), http.StatusInternalServerError)
			return
		}
	} else {
		filter.ClusterIds = []string{}
	}

	result := provider.GetGenericPersister().Model(&model.KubernetesResource{}).
		Where("kubernetes_resources.cluster_id IN (?)", filter.ClusterIds).
		Joins("JOIN kubernetes_resource_object_meta ON kubernetes_resource_object_meta.id = kubernetes_resources.id")

	if kind != "" {
		result = result.Where(&model.KubernetesResource{Kind: kind})
	}

	if apiVersion != "" {
		result = result.Where(&model.KubernetesResource{APIVersion: apiVersion})
	}

	if isLabels {
		result = result.Preload("KubernetesResourceMeta.Labels", "kind = ?", model.KindLabel)
	}
	if isAnnotaion {
		result = result.Preload("KubernetesResourceMeta.Annotations", "kind = ?", model.KindAnnotation)
	}

	if spec {
		result = result.Preload("Spec")
	}

	if status {
		result = result.Preload("Status")
	}

	if search != "" {
		result = result.Where("kubernetes_resource_object_meta.name LIKE ?", "%"+search+"%")
	}

	result.Count(&totalCount)

	if limit != 0 {
		result = result.Limit(limit)
	}

	if offset != 0 {
		result = result.Offset(offset)
	}

	if order != "" {
		if sort == "desc" {
			result = result.Order(clause.OrderByColumn{Column: clause.Column{Name: order}, Desc: true})
		} else {
			result = result.Order(order)
		}
	}

	err := result.Find(&resources).Error
	if err != nil {
		h.log.Error(ErrFetchMeshSyncResources(err))
		http.Error(rw, ErrFetchMeshSyncResources(err).Error(), http.StatusInternalServerError)
		return
	}

	var pgSize int
	if limitstr == "all" {
		pgSize = len(resources)
	} else {
		pgSize = limit
	}

	response := &models.MeshSyncResourcesAPIResponse{
		Page:       page,
		PageSize:   pgSize,
		TotalCount: totalCount,
		Resources:  resources,
	}

	if err := enc.Encode(response); err != nil {
		h.log.Error(ErrFetchMeshSyncResources(err))
		http.Error(rw, ErrFetchMeshSyncResources(err).Error(), http.StatusInternalServerError)
	}
}

// swagger:route GET /api/system/meshsync/resources/kinds GetMeshSyncResourcesKinds idGetMeshSyncResourcesKinds
// Handle GET request for meshsync discovered resources kinds
//
// ```?clusterId={[clusterId]}``` clusterId is array of string values. Required.
//
// ```?page = {page-number}``` Default page number is 1
//
// ```?pagesize = {pagesize}``` Default pagesize is 25. To return all results: ```pagesize=all```
//
// ```?search = {componentname}``` If search is non empty then a greedy search is performed
//
// ```?order = {field}``` orders on the passed field
//
// responses:
// 200: []meshsyncResourcesKindsResponseWrapper

func (h *Handler) GetMeshSyncResourcesKinds(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	rw.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(rw)

	page, offset, limit,
		search, order, sort, _ := getPaginationParams(r)

	filter := struct {
		ClusterIds []string `json:"clusterIds"`
	}{}

	var kinds []string
	var totalCount int64

	clusterIds := r.URL.Query().Get("clusterIds")
	if clusterIds != "" {
		err := json.Unmarshal([]byte(clusterIds), &filter.ClusterIds)
		if err != nil {
			h.log.Error(ErrFetchMeshSyncResources(err))
			http.Error(rw, ErrFetchMeshSyncResources(err).Error(), http.StatusInternalServerError)
			return
		}
	} else {
		filter.ClusterIds = []string{}
	}

	result := provider.GetGenericPersister().Model(&model.KubernetesResource{}).Distinct("kind").
		Where("kubernetes_resources.cluster_id IN (?)", filter.ClusterIds)

	if search != "" {
		result = result.Where("kubernetes_resources.kind LIKE ?", "%"+search+"%")
	}

	if limit != 0 {
		result = result.Limit(limit)
	}

	if offset != 0 {
		result = result.Offset(offset)
	}

	if order != "" {
		if sort == "desc" {
			result = result.Order(clause.OrderByColumn{Column: clause.Column{Name: order}, Desc: true})
		} else {
			result = result.Order(order)
		}
	}

	err := result.Pluck("kinds", &kinds).Error
	if err != nil {
		h.log.Error(ErrFetchMeshSyncResources(err))
		http.Error(rw, ErrFetchMeshSyncResources(err).Error(), http.StatusInternalServerError)
		return
	}

	response := &models.MeshSyncResourcesKindsAPIResponse{
		Kinds:      kinds,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   limit,
	}

	if err := enc.Encode(response); err != nil {
		h.log.Error(ErrFetchMeshSyncResources(err))
		http.Error(rw, ErrFetchMeshSyncResources(err).Error(), http.StatusInternalServerError)
	}
}
func (h *Handler) DeleteMeshSyncResource(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	resourceID := mux.Vars(r)["id"]
	db := provider.GetGenericPersister()
	err := db.Model(&model.KubernetesResource{}).Delete(&model.KubernetesResource{ID: resourceID}).Error
	if err != nil {
		h.log.Error(models.ErrDelete(err, "meshsync data", http.StatusInternalServerError))
	}
}
