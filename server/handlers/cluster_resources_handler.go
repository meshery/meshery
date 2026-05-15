// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
)

// clusterResourcesQuery counts the kubernetes_resources rows grouped by kind,
// scoped to the supplied cluster IDs and (optionally) namespace. Mirrors the
// SQL used in server/internal/graphql/resolver/kubernetes.go:getClusterResources
// because the resolver function is package-private — duplicating here is
// preferable to widening the resolver's surface for code that the GraphQL
// teardown PR will delete.
//
// Three legs unioned together:
//  1. Cluster-scoped resources (no namespace, kind != "Namespace")
//  2. Namespaced resources matching the supplied namespace
//  3. Namespace resources themselves (kind == "Namespace")
const clusterResourcesQuery = `
	SELECT count(kind) as count, kind FROM kubernetes_resources kr LEFT JOIN kubernetes_resource_object_meta rom on kr.id = rom.id
		WHERE kr.kind <> 'Namespace' AND rom.namespace = '' AND kr.cluster_id IN (?) GROUP BY kind
			UNION
	SELECT count(kind) as count, kind FROM kubernetes_resources kr LEFT JOIN kubernetes_resource_object_meta rom on kr.id = rom.id
		WHERE rom.namespace IN (?) AND kr.cluster_id IN (?) GROUP BY kind
			UNION
	SELECT count(kind) as count, kind FROM kubernetes_resources kr
		WHERE kr.kind = 'Namespace' AND kr.cluster_id IN (?) GROUP BY kind`

// GetClusterResourcesHandler returns a per-kind count of kubernetes resources
// observed across the supplied cluster IDs (and, optionally, namespace).
// Replaces the `subscribeClusterResources` GraphQL subscription with a plain
// REST poll endpoint — push semantics weren't load-bearing for this surface
// (the resolver re-ran the same SQL on every tick).
//
// Query params:
//   - clusterIds : required, repeatable. Cluster IDs to scope the count.
//   - namespace  : optional. When set, includes namespaced resources matching
//     this namespace in the count; otherwise namespaced resources are not
//     included (matching the resolver's three-leg UNION semantics).
//
// Response shape mirrors the GraphQL ClusterResources type:
//
//	{ "resources": [ { "kind": "Pod", "count": 12 }, ... ] }
func (h *Handler) GetClusterResourcesHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	cids := req.URL.Query()["clusterIds"]
	if len(cids) == 0 {
		writeJSONError(w, "at least one clusterIds query parameter is required", http.StatusBadRequest)
		return
	}
	namespace := req.URL.Query().Get("namespace")

	persister := provider.GetGenericPersister()
	if persister == nil {
		writeJSONError(w, "persister unavailable", http.StatusInternalServerError)
		return
	}

	var rows *sql.Rows
	var err error
	rows, err = persister.Raw(clusterResourcesQuery, cids, namespace, cids, cids).Rows()
	if err != nil {
		h.log.Error(err)
		writeJSONError(w, "failed to query cluster resources", http.StatusInternalServerError)
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			h.log.Error(closeErr)
		}
	}()

	resources := make([]*model.Resource, 0)
	for rows.Next() {
		var r model.Resource
		if err := rows.Scan(&r.Count, &r.Kind); err != nil {
			h.log.Error(err)
			writeJSONError(w, "failed to scan cluster resource row", http.StatusInternalServerError)
			return
		}
		resources = append(resources, &r)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&model.ClusterResources{Resources: resources}); err != nil {
		h.log.Error(models.ErrMarshal(err, "cluster resources response"))
		return
	}
}
