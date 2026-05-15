// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/meshery/meshery/server/helpers"
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
)

// telemetryComponentsQuery selects Service resources from the meshsync DB
// and returns their name + spec + status. Mirrors the SQL used in
// server/internal/graphql/resolver/telemetry.go:getTelemetryComps because
// the resolver function is package-private — duplicating here is preferable
// to widening the resolver's surface for code that the GraphQL teardown PR
// will delete.
const telemetryComponentsQuery = `
SELECT rom.name, rs.attribute, rst.attribute FROM kubernetes_resources kr
	LEFT JOIN kubernetes_resource_object_meta rom on kr.id = rom.id
	INNER JOIN resource_specs rs on kr.id = rs.id
	INNER JOIN resource_statuses rst on kr.id = rst.id
WHERE kr.kind = 'Service' AND kr.cluster_id IN (?);
`

// GetTelemetryComponentsHandler returns the telemetry components
// (Prometheus, Grafana, Jaeger, etc.) discovered in the meshsync DB across
// the supplied cluster IDs. Filters down to the configured
// helpers.TelemetryComps allowlist so non-telemetry Services aren't surfaced.
// Replaces the `fetchTelemetryComponents` GraphQL query.
//
// Query params:
//   - clusterIds : required, repeatable. Cluster IDs to scope the search.
//
// Response shape mirrors the GraphQL TelemetryComp type:
//
//	[ { "name": "...", "spec": "...", "status": "..." }, ... ]
func (h *Handler) GetTelemetryComponentsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	cids := req.URL.Query()["clusterIds"]
	if len(cids) == 0 {
		writeJSONError(w, "at least one clusterIds query parameter is required", http.StatusBadRequest)
		return
	}

	persister := provider.GetGenericPersister()
	if persister == nil {
		writeJSONError(w, "persister unavailable", http.StatusInternalServerError)
		return
	}

	var rows *sql.Rows
	var err error
	rows, err = persister.Raw(telemetryComponentsQuery, cids).Rows()
	if err != nil {
		h.log.Error(err)
		writeJSONError(w, "failed to query telemetry components", http.StatusInternalServerError)
		return
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			h.log.Error(closeErr)
		}
	}()

	components := make([]*model.TelemetryComp, 0)
	for rows.Next() {
		var c model.TelemetryComp
		if err := rows.Scan(&c.Name, &c.Spec, &c.Status); err != nil {
			h.log.Error(err)
			writeJSONError(w, "failed to scan telemetry component row", http.StatusInternalServerError)
			return
		}
		// helpers.TelemetryComps is the allowlist of known telemetry service
		// names (prometheus, grafana, jaeger, etc.). Without this filter we'd
		// return every Service in the cluster, not just the telemetry ones.
		if utils.SliceContains(helpers.TelemetryComps, c.Name) {
			components = append(components, &c)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(components); err != nil {
		h.log.Error(models.ErrMarshal(err, "telemetry components response"))
		return
	}
}
