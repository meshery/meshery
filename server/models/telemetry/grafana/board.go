// Package grafana provides a small, self-contained client and data model for
// reading dashboards ("boards") and proxying datasource queries from a Grafana
// instance that has been registered as a first-class Meshery connection.
//
// It is intentionally independent of the legacy github.com/meshery/meshery/server/models
// Grafana support (which is built on the abandoned grafana-tools/sdk). The types
// here are the wire contract between the Meshery telemetry handlers and the UI.
package grafana

import "encoding/json"

// BoardSummary is a lightweight dashboard descriptor returned by Grafana's
// search API. It carries just enough to list and pick a board.
type BoardSummary struct {
	UID   string   `json:"uid"`
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Tags  []string `json:"tags,omitempty"`
}

// Board is a dashboard with the subset of its definition needed to render its
// panels inside Meshery.
type Board struct {
	UID          string        `json:"uid"`
	Title        string        `json:"title"`
	Tags         []string      `json:"tags,omitempty"`
	Panels       []Panel       `json:"panels"`
	TemplateVars []TemplateVar `json:"templateVars,omitempty"`
}

// GridPos mirrors Grafana's panel layout coordinates (24-column grid).
type GridPos struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// Panel is a single visualization within a board.
type Panel struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	GridPos GridPos  `json:"gridPos"`
	Unit    string   `json:"unit,omitempty"`
	Targets []Target `json:"targets"`
}

// Target is a single query backing a panel.
type Target struct {
	RefID         string `json:"refId"`
	DatasourceUID string `json:"datasourceUid,omitempty"`
	Expr          string `json:"expr"`
	LegendFormat  string `json:"legendFormat,omitempty"`
}

// TemplateVar is a board-level template/dashboard variable. Current holds the
// variable's currently-selected value(s), which the UI uses to resolve `$var`
// references inside panel queries.
type TemplateVar struct {
	Name       string   `json:"name"`
	Label      string   `json:"label,omitempty"`
	Query      string   `json:"query,omitempty"`
	Type       string   `json:"type,omitempty"`
	Current    []string `json:"current,omitempty"`
	Multi      bool     `json:"multi,omitempty"`
	IncludeAll bool     `json:"includeAll,omitempty"`
}

// HealthInfo is Grafana's GET /api/health response: reachability plus the
// instance's reported version and database status.
type HealthInfo struct {
	Commit   string `json:"commit,omitempty"`
	Database string `json:"database,omitempty"`
	Version  string `json:"version,omitempty"`
}

// Datasource is a Grafana datasource the panels may query against.
type Datasource struct {
	UID       string `json:"uid"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsDefault bool   `json:"isDefault"`
}

// PinnedBoard is a board a user has "added" to a connection. Pinned boards are
// persisted inside the connection's Metadata (no dedicated table).
type PinnedBoard struct {
	UID   string `json:"uid"`
	Title string `json:"title"`
}

// MetadataPinnedBoardsKey is the connection.Metadata key under which the list
// of PinnedBoard is stored.
const MetadataPinnedBoardsKey = "telemetryPinnedBoards"

// --- raw Grafana API shapes (internal decoding targets) ---

// rawSearchHit is one entry of Grafana's GET /api/search response.
type rawSearchHit struct {
	UID   string   `json:"uid"`
	Title string   `json:"title"`
	URL   string   `json:"url"`
	Type  string   `json:"type"`
	Tags  []string `json:"tags"`
}

// rawDashboardResponse is Grafana's GET /api/dashboards/uid/:uid response.
type rawDashboardResponse struct {
	Dashboard rawDashboard `json:"dashboard"`
}

type rawDashboard struct {
	UID        string        `json:"uid"`
	Title      string        `json:"title"`
	Tags       []string      `json:"tags"`
	Panels     []rawPanel    `json:"panels"`
	Templating rawTemplating `json:"templating"`
}

type rawTemplating struct {
	List []rawTemplateVar `json:"list"`
}

type rawTemplateVar struct {
	Name       string          `json:"name"`
	Label      string          `json:"label"`
	Type       string          `json:"type"`
	Query      json.RawMessage `json:"query"`
	Multi      bool            `json:"multi"`
	IncludeAll bool            `json:"includeAll"`
	Current    rawCurrent      `json:"current"`
}

type rawCurrent struct {
	Value json.RawMessage `json:"value"`
}

type rawPanel struct {
	ID          int             `json:"id"`
	Title       string          `json:"title"`
	Type        string          `json:"type"`
	GridPos     GridPos         `json:"gridPos"`
	Datasource  json.RawMessage `json:"datasource"`
	FieldConfig rawFieldConfig  `json:"fieldConfig"`
	Targets     []rawTarget     `json:"targets"`
	// Panels holds a collapsed row's child panels. Grafana nests them here when
	// the row is collapsed; an expanded row's children are top-level siblings.
	Panels []rawPanel `json:"panels"`
}

type rawFieldConfig struct {
	Defaults struct {
		Unit string `json:"unit"`
	} `json:"defaults"`
}

type rawTarget struct {
	RefID        string          `json:"refId"`
	Expr         string          `json:"expr"`
	LegendFormat string          `json:"legendFormat"`
	Datasource   json.RawMessage `json:"datasource"`
}
