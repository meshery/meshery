// Package prometheus provides a small, self-contained client and data model for
// querying a Prometheus instance that has been registered as a first-class
// Meshery connection.
//
// It is intentionally independent of the legacy github.com/meshery/meshery/server/models
// Prometheus support. Unlike the sibling grafana package, Prometheus IS the
// datasource, so there is no datasource ("ds") concept here. The types here are
// the wire contract between the Meshery telemetry handlers and the UI.
package prometheus

import "encoding/json"

// HealthInfo is the reachability + version summary derived from Prometheus'
// GET /api/v1/status/buildinfo response. A reachable instance whose body can't be
// parsed still yields a (versionless) HealthInfo rather than an error.
type HealthInfo struct {
	Version string `json:"version,omitempty"`
}

// BatchQuery is a single resolved query within a batch. ID is an opaque,
// caller-chosen key echoed back on the matching result; Query is the PromQL
// expression. There is no datasource reference — Prometheus is the datasource.
type BatchQuery struct {
	ID    string
	Query string
}

// BatchResult carries the outcome of one BatchQuery. ID matches the input query.
// Exactly one of Body (the raw Prometheus query_range JSON, unmodified) or Err
// (a per-query failure) is set.
type BatchResult struct {
	ID   string
	Body json.RawMessage
	Err  error
}

// GridPos mirrors a panel's layout coordinates (24-column grid), matching the
// grafana package's GridPos so the UI can share layout logic.
type GridPos struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// MetricPanel is a user-saved panel: a PromQL expression plus how to render it.
// Saved panels are persisted inside the connection's Metadata (no dedicated
// table), the same way grafana pinned boards are.
type MetricPanel struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Expr    string  `json:"expr"`
	Type    string  `json:"type"` // timeseries|stat|gauge|bar
	Unit    string  `json:"unit,omitempty"`
	GridPos GridPos `json:"gridPos"`
}

// MetadataPanelsKey is the connection.Metadata key under which the list of
// MetricPanel is stored.
const MetadataPanelsKey = "telemetryPrometheusPanels"
