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
	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	var resources []model.KubernetesResource
	var totalCount int64

	apiVersion := r.URL.Query().Get("apiVersion")
	spec, _ := strconv.ParseBool(r.URL.Query().Get("spec"))
	status, _ := strconv.ParseBool(r.URL.Query().Get("status"))
	isAnnotaion, _ := strconv.ParseBool(r.URL.Query().Get("annotations"))
	isLabels, _ := strconv.ParseBool(r.URL.Query().Get("labels"))
	// kind is an array of strings
	kind := r.URL.Query()["kind"]

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
		Preload("KubernetesResourceMeta").
		Where("kubernetes_resources.cluster_id IN (?)", filter.ClusterIds)

	if len(kind) > 0 {
		// result = result.Where(&model.KubernetesResource{Kind: kind})
		result = result.Where("kubernetes_resources.kind IN (?)", kind)
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

		result = result.
			Joins("JOIN kubernetes_resource_object_meta ON kubernetes_resource_object_meta.id = kubernetes_resources.id").
			Where("kubernetes_resource_object_meta.name LIKE ?", "%"+search+"%")
	}

	result.Count(&totalCount)

	if limit != 0 {
		result = result.Limit(limit)
	}

	if offset != 0 {
		result = result.Offset(offset)
	}

	order = models.SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
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
	if limit == 0 {
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

// swagger:route GET /api/system/meshsync/resources/summary GetMeshSyncResourcesSummary idGetMeshSyncResourcesSummary
// Handle GET request for meshsync discovered resources
//
// ```?clusterId={clusterId}``` clusterId is id of the cluster to get resources for ( multiple supported)
//
//
// responses:
// 200: []meshsyncResourcesSummaryResponseWrapper

func (h *Handler) GetMeshSyncResourcesSummary(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	rw.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(rw)

	clusterIds := r.URL.Query()["clusterId"]
	h.log.Info("Fetching meshsync resources summary", "clusterIds", clusterIds)

	if len(clusterIds) == 0 {
		http.Error(rw, "clusterIds is required", http.StatusBadRequest)
		return
	}

	var kindCounts []struct {
		Kind  string
		Count int64
	}
	var namespaces []string

	err := provider.GetGenericPersister().
		Model(&model.KubernetesResource{}).
		Select("kind, count(*) as count").
		Group("kind").
		Where("kubernetes_resources.cluster_id IN (?)", clusterIds).
		Scan(&kindCounts).Error

	err = provider.GetGenericPersister().
		Model(&model.KubernetesResource{}).
		Joins("JOIN kubernetes_resource_object_meta ON kubernetes_resources.id = kubernetes_resource_object_meta.id").
		Select("distinct namespace").
		Where("kubernetes_resources.cluster_id IN (?)", clusterIds).
		Scan(&namespaces).Error

	if err != nil {
		h.log.Error(ErrFetchMeshSyncResources(err))
		http.Error(rw, ErrFetchMeshSyncResources(err).Error(), http.StatusInternalServerError)
		return
	}

	response := &models.MeshSyncResourcesSummaryAPIResponse{
		Kinds:      kindCounts,
		Namespaces: namespaces,
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
