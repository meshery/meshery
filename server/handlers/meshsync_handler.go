package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshsync/pkg/model"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"gorm.io/gorm/clause"
)

// MapToStruct converts a map[string]interface{} to a specified struct type.
// The targetStruct should be a pointer to the desired struct.
func MapToStruct(data map[string]interface{}, targetStruct interface{}) error {
	// Validate input
	if data == nil {
		return fmt.Errorf("data is nil")
	}
	if targetStruct == nil {
		return fmt.Errorf("targetStruct is nil")
	}

	// Marshal the map into JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal map: %w", err)
	}

	// Unmarshal the JSON into the target struct
	if err := json.Unmarshal(jsonData, targetStruct); err != nil {
		return fmt.Errorf("failed to unmarshal JSON into struct: %w", err)
	}

	return nil
}

// JsonParse safely parses a JSON string into an interface{} with a fallback default value.
func JsonParse(input *string, allowEmpty bool, defaultValue interface{}) interface{} {
	if input == nil {
		return defaultValue
	}
	var result interface{}
	if err := json.Unmarshal([]byte(*input), &result); err != nil {
		if allowEmpty {
			return defaultValue
		}
		return nil
	}
	return result
}

// ConvertToPatternFile converts a list of Kubernetes resources to a PatternFile.
func ConvertToPatternFile(resources []model.KubernetesResource, stripSchema bool) pattern.PatternFile {
	components := []*component.ComponentDefinition{}

	for _, resource := range resources {
		var componentDef component.ComponentDefinition
		err := MapToStruct(resource.ComponentMetadata, &componentDef)

		if err != nil {
			continue
		}

		componentDef.Id = uuid.FromStringOrNil(resource.KubernetesResourceMeta.UID)
		componentDef.DisplayName = resource.KubernetesResourceMeta.Name
		componentDef.Configuration = map[string]interface{}{
			"metadata": resource.KubernetesResourceMeta,
			"spec":     JsonParse(&resource.Spec.Attribute, true, map[string]interface{}{}),
			"data":     JsonParse(&resource.Data, true, map[string]interface{}{}),
		}
		// componentDef.Metadata.InstanceDetails = resource

		if stripSchema {
			componentDef.Component.Schema = ""
		}

		components = append(components, &componentDef)
	}

	var emptyUUID uuid.UUID

	return pattern.PatternFile{
		Name:          "EmptyDesignName",
		Id:            emptyUUID,
		SchemaVersion: "designs.meshery.io/v1beta1",
		Version:       "v1",
		Components:    components,
		Relationships: []*relationship.RelationshipDefinition{},
	}
}

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
// ```?asDesign={bool}``` asDesign is a boolean value. If true then the response is returned as a design and resources are omitted
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
	asDesign, _ := strconv.ParseBool(r.URL.Query().Get("asDesign"))
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

	var design pattern.PatternFile

	if asDesign {
		rawDesign := ConvertToPatternFile(resources, true) // strip schema
		resources = []model.KubernetesResource{}           // clear resources to save memory
		evalResponse, error := h.Rego.RegoPolicyHandler(rawDesign, RelationshipPolicyPackageName)
		if error != nil {
			design = rawDesign
			h.log.Error(fmt.Errorf("Error evaluating design: %v", error))
		} else {
			// if there is error in evaluation, return the raw design (without any relationships)
			design = evalResponse.Design
		}

	}

	response := &models.MeshSyncResourcesAPIResponse{
		Page:       page,
		PageSize:   pgSize,
		TotalCount: totalCount,
		Resources:  resources,
		Design:     design,
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

	// TODO: simplify into one query if possible
	err1 := provider.GetGenericPersister().
		Model(&model.KubernetesResource{}).
		Select("kind, count(*) as count").
		Group("kind").
		Where("kubernetes_resources.cluster_id IN (?)", clusterIds).
		Scan(&kindCounts).Error

	if err1 != nil {
		h.log.Error(ErrFetchMeshSyncResources(err1))
	}

	err2 := provider.GetGenericPersister().
		Model(&model.KubernetesResource{}).
		Joins("JOIN kubernetes_resource_object_meta ON kubernetes_resources.id = kubernetes_resource_object_meta.id").
		Select("distinct namespace").
		Where("kubernetes_resources.cluster_id IN (?)", clusterIds).
		Scan(&namespaces).Error

	if err2 != nil {
		h.log.Error(ErrFetchMeshSyncResources(err2))
	}

	// only return error if both queries failed
	if err1 != nil && err2 != nil {
		combinedErr := fmt.Errorf("Error fetching meshsync resources summary: %v, %v", err1, err2)
		http.Error(rw, ErrFetchMeshSyncResources(combinedErr).Error(), http.StatusInternalServerError)
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
