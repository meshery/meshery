// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
)

// GetMeshesAddonsHandler returns the addons (e.g. zipkin, grafana, prometheus)
// installed for a given service mesh / cluster. It replaces the
// `getAvailableAddons` GraphQL query.
//
// Query params:
//   - meshType  : optional. One of model.MeshType values (e.g. "ISTIO").
//                 When absent/empty, all mesh types are queried.
//   - clusterId : optional, repeatable. Filters by Kubernetes cluster IDs.
func (h *Handler) GetMeshesAddonsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	q := req.URL.Query()

	// Build the selector list. The legacy resolver nil-derefs filter.Type when
	// Type is unset (see addons.go:17); guarding the nil here is the bug-fix
	// half of this migration.
	var selectors []model.MeshType
	meshType := q.Get("meshType")
	if meshType == "" {
		selectors = append(selectors, model.AllMeshType...)
	} else {
		mt := model.MeshType(meshType)
		if !mt.IsValid() {
			writeJSONError(w, "invalid meshType", http.StatusBadRequest)
			return
		}
		if mt == model.MeshTypeAllMesh {
			selectors = append(selectors, model.AllMeshType...)
		} else {
			selectors = append(selectors, mt)
		}
	}

	cids := q["clusterId"]

	addons, err := model.GetAddonsState(req.Context(), selectors, provider, cids)
	if err != nil {
		h.log.Error(err)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if addons == nil {
		addons = []*model.AddonList{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(addons); err != nil {
		h.log.Error(models.ErrMarshal(err, "addons response"))
		return
	}
}

// GetKubernetesNamespacesHandler returns the distinct namespace names
// observed across the given cluster IDs. It replaces the
// `getAvailableNamespaces` GraphQL query.
//
// Query params:
//   - clusterIds : optional, repeatable. Filters by Kubernetes cluster IDs.
//     If empty the underlying provider behavior applies (which may return
//     an empty list when no clusters are scoped — same as the resolver).
//
// Response shape:
//
//	{ "namespaces": [{ "namespace": "default" }, ...] }
//
// matching what the UI sseSubscribe/RTK Query consumers expect.
func (h *Handler) GetKubernetesNamespacesHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	cids := req.URL.Query()["clusterIds"]
	if len(cids) == 1 {
		// Accept a comma-separated single value as well, so callers can
		// pass either ?clusterIds=a&clusterIds=b or ?clusterIds=a,b.
		for _, part := range splitAndTrim(cids[0], ',') {
			if part != "" && part != cids[0] {
				cids = splitAndTrim(cids[0], ',')
				break
			}
		}
	}

	names, err := model.SelectivelyFetchNamespaces(cids, provider)
	if err != nil {
		h.log.Error(err)
		writeJSONError(w, "failed to fetch namespaces", http.StatusInternalServerError)
		return
	}

	type namespaceEntry struct {
		Namespace string `json:"namespace"`
	}
	out := struct {
		Namespaces []namespaceEntry `json:"namespaces"`
	}{Namespaces: make([]namespaceEntry, 0, len(names))}
	for _, n := range names {
		out.Namespaces = append(out.Namespaces, namespaceEntry{Namespace: n})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		h.log.Error(models.ErrMarshal(err, "namespaces response"))
		return
	}
}

// splitAndTrim splits s by sep and trims whitespace around each piece,
// dropping empty pieces. Local helper to avoid pulling in strings for one
// call site.
func splitAndTrim(s string, sep rune) []string {
	out := make([]string, 0)
	start := 0
	for i, r := range s {
		if r == sep {
			seg := s[start:i]
			for len(seg) > 0 && (seg[0] == ' ' || seg[0] == '\t') {
				seg = seg[1:]
			}
			for len(seg) > 0 && (seg[len(seg)-1] == ' ' || seg[len(seg)-1] == '\t') {
				seg = seg[:len(seg)-1]
			}
			if seg != "" {
				out = append(out, seg)
			}
			start = i + 1
		}
	}
	seg := s[start:]
	for len(seg) > 0 && (seg[0] == ' ' || seg[0] == '\t') {
		seg = seg[1:]
	}
	for len(seg) > 0 && (seg[len(seg)-1] == ' ' || seg[len(seg)-1] == '\t') {
		seg = seg[:len(seg)-1]
	}
	if seg != "" {
		out = append(out, seg)
	}
	return out
}

// GetMeshesControlPlanesHandler returns the control-plane state for the
// specified service-mesh type and cluster IDs. It replaces the
// `getControlPlanes` GraphQL query.
//
// Query params:
//   - type      : optional. One of model.MeshType values. Absent => all.
//   - clusterId : optional, repeatable. Filters by Kubernetes cluster IDs.
func (h *Handler) GetMeshesControlPlanesHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	q := req.URL.Query()

	var selectors []model.MeshType
	meshType := q.Get("type")
	if meshType == "" {
		selectors = append(selectors, model.AllMeshType...)
	} else {
		mt := model.MeshType(meshType)
		if !mt.IsValid() {
			writeJSONError(w, "invalid type", http.StatusBadRequest)
			return
		}
		if mt == model.MeshTypeAllMesh {
			selectors = append(selectors, model.AllMeshType...)
		} else {
			selectors = append(selectors, mt)
		}
	}

	cids := q["clusterId"]

	planes, err := model.GetControlPlaneState(req.Context(), selectors, provider, cids)
	if err != nil {
		h.log.Error(err)
		writeJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if planes == nil {
		planes = []*model.ControlPlane{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(planes); err != nil {
		h.log.Error(models.ErrMarshal(err, "control planes response"))
		return
	}
}

// GetMeshesDataPlanesHandler returns the data-plane (sidecar) state for the
// specified service-mesh type and cluster IDs. It replaces the
// `getDataPlanes` GraphQL query.
//
// Query params:
//   - type      : optional. One of model.MeshType values. Absent => all.
//   - clusterId : optional, repeatable. Filters by Kubernetes cluster IDs.
func (h *Handler) GetMeshesDataPlanesHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	q := req.URL.Query()

	var selectors []model.MeshType
	meshType := q.Get("type")
	if meshType == "" {
		selectors = append(selectors, model.AllMeshType...)
	} else {
		mt := model.MeshType(meshType)
		if !mt.IsValid() {
			writeJSONError(w, "invalid type", http.StatusBadRequest)
			return
		}
		if mt == model.MeshTypeAllMesh {
			selectors = append(selectors, model.AllMeshType...)
		} else {
			selectors = append(selectors, mt)
		}
	}

	cids := q["clusterId"]

	planes, err := model.GetDataPlaneState(req.Context(), selectors, provider, cids)
	if err != nil {
		h.log.Error(err)
		writeJSONError(w, "failed to fetch data planes", http.StatusInternalServerError)
		return
	}

	if planes == nil {
		planes = []*model.DataPlane{}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(planes); err != nil {
		h.log.Error(models.ErrMarshal(err, "data planes response"))
		return
	}
}
