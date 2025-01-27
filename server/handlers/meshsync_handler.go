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
	"gorm.io/gorm"
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

		var spec interface{}
		if resource.Spec != nil {
			spec = JsonParse(&resource.Spec.Attribute, true, map[string]interface{}{})
		}

		labels := map[string]string{}
		annotations := map[string]string{}

		if resource.KubernetesResourceMeta.Labels != nil {
			for _, label := range resource.KubernetesResourceMeta.Labels {
				labels[label.Key] = label.Value
			}
		}

		if resource.KubernetesResourceMeta.Annotations != nil {
			for _, annotation := range resource.KubernetesResourceMeta.Annotations {
				annotations[annotation.Key] = annotation.Value
			}
		}

		metadata := map[string]interface{}{
			"labels":      labels,
			"annotations": annotations,
			"name":        resource.KubernetesResourceMeta.Name,
		}

		// only set namespace if it is not empty otherwise evaluator creates empty namespace
		if resource.KubernetesResourceMeta.Namespace != "" {
			metadata["namespace"] = resource.KubernetesResourceMeta.Namespace
		}

		componentDef.Configuration = map[string]interface{}{
			"metadata": metadata,
			"spec":     spec,
			"data":     JsonParse(&resource.Data, true, map[string]interface{}{}),
		}

		if componentDef.Metadata.AdditionalProperties == nil {
			componentDef.Metadata.AdditionalProperties = make(map[string]interface{})
		}

		componentDef.Metadata.AdditionalProperties["instanceDetails"] = resource
		if stripSchema {
			componentDef.Component.Schema = ""
		}

		components = append(components, &componentDef)
	}

	var emptyUUID uuid.UUID

	return pattern.PatternFile{
		Name:          "ClusterSnapshot",
		Id:            emptyUUID,
		SchemaVersion: "designs.meshery.io/v1beta1",
		Version:       "v1",
		Components:    components,
		Relationships: []*relationship.RelationshipDefinition{},
	}
}

// scopes down the query based on the namespaces
func filterByNamespaces(query *gorm.DB, namespaces []string) *gorm.DB {
	if len(namespaces) > 0 {
		return query.Where("kubernetes_resource_object_meta.namespace IN (?) or ( kubernetes_resources.kind = 'Namespace' AND kubernetes_resource_object_meta.name IN (?) )", namespaces, namespaces)

	}
	return query
}

func searchResources(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		return query.Where("kubernetes_resource_object_meta.name LIKE ?", "%"+search+"%")
	}
	return query
}

func filterByKinds(query *gorm.DB, kinds []string) *gorm.DB {
	if len(kinds) > 0 {
		return query.Where("kubernetes_resources.kind IN (?)", kinds)
	}
	return query
}

func filterByClusters(query *gorm.DB, clusterIDs []string) *gorm.DB {
	if len(clusterIDs) > 0 {
		return query.Where("kubernetes_resources.cluster_id IN (?)", clusterIDs)
	}
	return query
}

func filterByPatternIds(query *gorm.DB, patternIDs []string) *gorm.DB {
	if len(patternIDs) > 0 {
		return query.Where("kubernetes_resources.pattern_resource IN (?)", patternIDs)
	}
	return query
}

// filterByKey filters a GORM query to only include resources with specific key-value pairs in their metadata.
//
// This function performs the following steps:
// 1. Joins the `kubernetes_resources` table with `kubernetes_resource_object_meta` to access metadata for resources.
// 2. Joins the `kubernetes_resource_object_meta` table with `kubernetes_key_values` to access key-value pairs.
// 3. Filters the results based on:
//   - The `kind` of the key-value pair (e.g., "label", "annotation").
//   - A set of key-value pairs provided in the `keyValues` slice, formatted as "key=value".
//
// The query ensures only matching resources are returned.
//
// Parameters:
// - query: The initial GORM query to apply the filters to.
// - kind: The kind of key-value pair to filter by (e.g., "label").
// - keyValues: A slice of strings in the format "key=value" used to match key-value pairs.
//
// Returns:
// - A filtered GORM query that includes only resources with matching key-value pairs.
func filterByKey(query *gorm.DB, kind string, keyValues []string) *gorm.DB {
	if len(keyValues) == 0 {
		return query
	}
	return query.
		// Joins("join kubernetes_resource_object_meta on kubernetes_resources.id = kubernetes_resource_object_meta.id").
		Joins("join kubernetes_key_values on kubernetes_resource_object_meta.id = kubernetes_key_values.id").
		// We need to use the `||` operator to concatenate the key and value, so we can compare it to the key-value pairs. and not cross match with other key value pairs
		Where("kubernetes_key_values.kind = ? AND  (kubernetes_key_values.key || '=' || kubernetes_key_values.value) IN (?)", kind, keyValues)
}

func selectDistinctKeyValues(db *gorm.DB, kind string) *gorm.DB {

	return db.
		Select("DISTINCT kubernetes_key_values.key, kubernetes_key_values.value").
		Joins("join kubernetes_resource_object_meta on kubernetes_resources.id = kubernetes_resource_object_meta.id").
		Joins("join kubernetes_key_values on kubernetes_resource_object_meta.id = kubernetes_key_values.id").
		Where("kubernetes_key_values.kind = ?", kind)
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
	patternIds := r.URL.Query()["patternId"]

	namespaces := r.URL.Query()["namespace"] // namespace is an array of strings to scope the resources
	labels := r.URL.Query()["label"]         // label is an array of strings
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

	query := provider.GetGenericPersister().Model(&model.KubernetesResource{}).
		Joins("JOIN kubernetes_resource_object_meta ON kubernetes_resource_object_meta.id = kubernetes_resources.id").
		Preload("KubernetesResourceMeta").
		Where("kubernetes_resources.cluster_id IN (?)", filter.ClusterIds)

	query = filterByNamespaces(query, namespaces)
	query = searchResources(query, search)
	query = filterByKinds(query, kind)
	query = filterByPatternIds(query, patternIds)
	query = filterByKey(query, model.KindLabel, labels)

	if apiVersion != "" {
		query = query.Where(&model.KubernetesResource{APIVersion: apiVersion})
	}

	if isLabels {
		query = query.Preload("KubernetesResourceMeta.Labels", "kind = ?", model.KindLabel)
	}
	if isAnnotaion {
		query = query.Preload("KubernetesResourceMeta.Annotations", "kind = ?", model.KindAnnotation)
	}

	if spec {
		query = query.Preload("Spec")
	}

	if status {
		query = query.Preload("Status")
	}

	query.Count(&totalCount)

	if limit != 0 {
		query = query.Limit(limit)
	}

	if offset != 0 {
		query = query.Offset(offset)
	}

	order = models.SanitizeOrderInput(order, []string{"created_at", "updated_at", "name"})
	if order != "" {
		if sort == "desc" {
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: order}, Desc: true})
		} else {
			query = query.Order(order)
		}
	}

	// prrint the query
	h.log.Info("Resources query", query.Statement.SQL.String())

	err := query.Find(&resources).Error
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
		// evalResponse, error := h.Rego.RegoPolicyHandler(rawDesign, RelationshipPolicyPackageName)
		evalResponse, error := h.EvaluateDesign(pattern.EvaluationRequest{
			Design: rawDesign,
		})

		if error != nil {
			design = rawDesign
			h.log.Error(fmt.Errorf("Error evaluating design: %v", error))
		} else {
			design = evalResponse.Design // use the evaluated design
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
	namespaceScope := r.URL.Query()["namespace"]
	patternIds := r.URL.Query()["patternId"]
	h.log.Info("Fetching meshsync resources summary", "clusterIds", clusterIds)

	if len(clusterIds) == 0 {
		http.Error(rw, "clusterIds is required", http.StatusBadRequest)
		return
	}

	var kindCounts []struct {
		Kind  string
		Model string
		Count int64
	}
	var namespaces []string

	// TODO: simplify into one query if possible
	kindsQuery := provider.GetGenericPersister().
		Model(&model.KubernetesResource{}).
		Joins("JOIN kubernetes_resource_object_meta ON kubernetes_resources.id = kubernetes_resource_object_meta.id").
		Select("kind, model, count(*) as count").
		Group("kind, model").
		Having("model IS NOT NULL")

	kindsQuery = filterByClusters(kindsQuery, clusterIds)
	kindsQuery = filterByNamespaces(kindsQuery, namespaceScope)
	kindsQuery = filterByPatternIds(kindsQuery, patternIds)

	err1 := kindsQuery.Scan(&kindCounts).Error

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

	var labels []model.KubernetesKeyValue

	labelsQuery := selectDistinctKeyValues(provider.GetGenericPersister().Model(&model.KubernetesResource{}), "label")
	labelsQuery = filterByClusters(labelsQuery, clusterIds)
	labelsQuery = filterByNamespaces(labelsQuery, namespaceScope)
	labelsQuery = filterByPatternIds(labelsQuery, patternIds)

	err := labelsQuery.Scan(&labels).Error

	if err != nil {
		h.log.Error(ErrFetchMeshSyncResources(err))
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
		Labels:     labels,
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
