package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	prometheus "github.com/meshery/meshery/server/models/telemetry/prometheus"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// telemetry_prometheus_handler.go exposes a clean, read-oriented API over a
// Prometheus instance registered as a first-class Meshery connection. It is
// independent of the legacy /api/telemetry/metrics/* handlers and reuses the
// connection + credential plumbing instead of its own config flow. Routes live
// under /api/telemetry/prometheus/{connectionID}/*.

// defaultMetricNamesLimit caps how many metric names are returned by the metrics
// discovery endpoint when the caller does not specify a limit.
const defaultMetricNamesLimit = 2000

// prometheusTelemetryClientForConnection resolves a Prometheus connection and its
// credential into a ready-to-use prometheus.Client. It enforces that the
// connection is of kind "prometheus" and tolerates connections without a
// credential (anonymous Prometheus). The returned int is an HTTP status to use
// when err is non-nil.
func (h *Handler) prometheusTelemetryClientForConnection(token string, connectionID uuid.UUID, p models.Provider) (*prometheus.Client, *connections.Connection, int, error) {
	connection, statusCode, err := p.GetConnectionByID(token, connectionID)
	if err != nil {
		return nil, nil, statusCode, err
	}
	if connection.Kind != "prometheus" {
		return nil, connection, http.StatusBadRequest, ErrInvalidConnectionKind(connection.Kind, "prometheus")
	}

	baseURL, _ := connection.Metadata["url"].(string)
	secret := ""
	if connection.CredentialID != nil {
		cred, sc, cerr := p.GetCredentialByID(token, core.UUIDOrUUIDNil(connection.CredentialID))
		if cerr != nil {
			return nil, connection, sc, cerr
		}
		if cred != nil {
			secret, _ = cred.Secret["secret"].(string)
		}
	}
	return prometheus.New(baseURL, secret, h.log), connection, http.StatusOK, nil
}

// PrometheusTelemetryPingHandler verifies a Prometheus connection is reachable
// and its credential is accepted, emitting a notification event with the outcome.
//
// GET /api/telemetry/prometheus/{connectionID}/ping
func (h *Handler) PrometheusTelemetryPingHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	userID := user.ID
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	eventBuilder := events.NewEvent().
		ActedUpon(connectionID).
		FromOwner(userID).
		FromSystem(*h.SystemID).
		WithCategory("connection").
		WithAction("ping")

	client, connection, statusCode, err := h.prometheusTelemetryClientForConnection(token, connectionID, p)
	if err != nil {
		h.emitPrometheusPingFailure(p, token, userID, eventBuilder, connectionLabel(connection, connectionID), err)
		writeMeshkitError(w, err, statusCode)
		return
	}

	info, err := client.Health(req.Context())
	if err != nil {
		_err := ErrTelemetryPrometheus(err, "ping")
		h.emitPrometheusPingFailure(p, token, userID, eventBuilder, connection.Name, _err)
		writeMeshkitError(w, _err, http.StatusBadGateway)
		return
	}

	description := fmt.Sprintf("Prometheus connection %q is reachable.", connection.Name)
	if info.Version != "" {
		description = fmt.Sprintf("Prometheus connection %q is reachable (Prometheus %s).", connection.Name, info.Version)
	}
	metadata := map[string]interface{}{
		"version": info.Version,
	}
	if connURL, ok := connection.Metadata["url"].(string); ok {
		metadata["url"] = connURL
	}
	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).WithMetadata(metadata).Build()
	_ = p.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)

	h.writeJSON(w, map[string]interface{}{
		"reachable": true,
		"version":   info.Version,
	})
}

// emitPrometheusPingFailure persists and broadcasts an error event describing a
// failed Prometheus connection ping.
func (h *Handler) emitPrometheusPingFailure(p models.Provider, token string, userID core.Uuid, eventBuilder *events.EventBuilder, label string, err error) {
	event := eventBuilder.
		WithSeverity(events.Error).
		WithDescription(fmt.Sprintf("Prometheus connection %q is unreachable.", label)).
		WithMetadata(map[string]interface{}{"error": err}).
		Build()
	_ = p.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)
	h.log.Error(err)
}

// PrometheusTelemetryMetricsHandler lists the instance's metric names, optionally
// filtered by a case-insensitive substring and capped.
//
// GET /api/telemetry/prometheus/{connectionID}/metrics?search=&limit=
func (h *Handler) PrometheusTelemetryMetricsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	client, _, statusCode, err := h.prometheusTelemetryClientForConnection(token, connectionID, p)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}

	q := req.URL.Query()
	limit := defaultMetricNamesLimit
	if l := q.Get("limit"); l != "" {
		if parsed, perr := strconv.Atoi(l); perr == nil && parsed > 0 {
			limit = parsed
		}
	}
	names, err := client.MetricNames(req.Context(), q.Get("search"), limit)
	if err != nil {
		writeMeshkitError(w, ErrTelemetryPrometheus(err, "metric discovery"), http.StatusBadGateway)
		return
	}
	h.writeJSON(w, names)
}

// PrometheusTelemetryLabelsHandler lists the instance's label names, optionally
// constrained to a series selector.
//
// GET /api/telemetry/prometheus/{connectionID}/labels?match=
func (h *Handler) PrometheusTelemetryLabelsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	client, _, statusCode, err := h.prometheusTelemetryClientForConnection(token, connectionID, p)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}
	names, err := client.LabelNames(req.Context(), req.URL.Query().Get("match"))
	if err != nil {
		writeMeshkitError(w, ErrTelemetryPrometheus(err, "label discovery"), http.StatusBadGateway)
		return
	}
	h.writeJSON(w, names)
}

// PrometheusTelemetryLabelValuesHandler lists the values of a single label,
// optionally constrained to a series selector.
//
// GET /api/telemetry/prometheus/{connectionID}/label_values?label=&match=
func (h *Handler) PrometheusTelemetryLabelValuesHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	client, _, statusCode, err := h.prometheusTelemetryClientForConnection(token, connectionID, p)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}
	q := req.URL.Query()
	values, err := client.LabelValues(req.Context(), q.Get("label"), q.Get("match"))
	if err != nil {
		writeMeshkitError(w, ErrTelemetryPrometheus(err, "label value discovery"), http.StatusBadGateway)
		return
	}
	h.writeJSON(w, values)
}

// PrometheusTelemetryMetadataHandler proxies Prometheus' metric metadata response
// back verbatim for the UI to render.
//
// GET /api/telemetry/prometheus/{connectionID}/metadata?metric=
func (h *Handler) PrometheusTelemetryMetadataHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	client, _, statusCode, err := h.prometheusTelemetryClientForConnection(token, connectionID, p)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}
	body, err := client.Metadata(req.Context(), req.URL.Query().Get("metric"))
	if err != nil {
		writeMeshkitError(w, ErrTelemetryPrometheus(err, "metadata"), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

// PrometheusTelemetryQueryHandler proxies an instant Prometheus query and streams
// the response back verbatim for the UI to render.
//
// GET /api/telemetry/prometheus/{connectionID}/query?query=&time=
func (h *Handler) PrometheusTelemetryQueryHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	client, _, statusCode, err := h.prometheusTelemetryClientForConnection(token, connectionID, p)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}
	q := req.URL.Query()
	body, err := client.Query(req.Context(), q.Get("query"), q.Get("time"))
	if err != nil {
		// Panel queries fail one panel at a time and the UI renders each failure
		// inline. We therefore surface them ONLY as the HTTP response and never as
		// broadcast notification events: a single unreachable connection or bad
		// credential would otherwise spawn one notification per panel and bury the
		// user.
		status, merr := prometheusQueryError(err)
		writeMeshkitError(w, merr, status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

// PrometheusTelemetryQueryRangeHandler proxies a Prometheus range query and
// streams the response back verbatim for the UI to render.
//
// GET /api/telemetry/prometheus/{connectionID}/query_range?query=&start=&end=&step=
func (h *Handler) PrometheusTelemetryQueryRangeHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	client, _, statusCode, err := h.prometheusTelemetryClientForConnection(token, connectionID, p)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}
	q := req.URL.Query()
	body, err := client.QueryRange(req.Context(), q.Get("query"), q.Get("start"), q.Get("end"), q.Get("step"))
	if err != nil {
		// See PrometheusTelemetryQueryHandler: per-panel query failures are surfaced
		// inline as the HTTP response, never as broadcast notification events.
		status, merr := prometheusQueryError(err)
		writeMeshkitError(w, merr, status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

// maxPrometheusBatchQueries bounds how many queries a single batch request may
// carry, so a malformed or hostile client can't ask the server to fan out an
// unbounded number of upstream Prometheus calls in one request.
const maxPrometheusBatchQueries = 200

// prometheusBatchQuery is one entry in a batch query_range request body.
type prometheusBatchQuery struct {
	ID    string `json:"id"`
	Query string `json:"query"`
}

// prometheusBatchRequest is the body of a batched query_range request: a shared
// time window plus a set of per-panel queries.
type prometheusBatchRequest struct {
	Start   string                 `json:"start"`
	End     string                 `json:"end"`
	Step    string                 `json:"step"`
	Queries []prometheusBatchQuery `json:"queries"`
}

// prometheusBatchResultItem is one entry in a batch response. It always carries
// the originating id and exactly one of Response (the raw Prometheus query_range
// JSON, unmodified) or Error (a concise, classified failure message).
type prometheusBatchResultItem struct {
	ID       string          `json:"id"`
	Response json.RawMessage `json:"response,omitempty"`
	Error    string          `json:"error,omitempty"`
}

// PrometheusTelemetryQueryRangeBatchHandler runs many Prometheus range queries
// for a dashboard in a single request, fanning them out concurrently. The UI
// sends one request per dashboard instead of one per panel.
//
// A per-query failure is reported inline on that result item (never failing the
// whole batch) and, like the single query_range path, is NOT emitted as a
// notification event: one bad connection would otherwise bury the user under one
// notification per panel. Only a wholly invalid request (bad body, no/too-many
// queries, or an unresolvable connection) returns a non-200 status.
//
// POST /api/telemetry/prometheus/{connectionID}/query_range_batch
func (h *Handler) PrometheusTelemetryQueryRangeBatchHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeMeshkitError(w, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	var request prometheusBatchRequest
	if err := json.Unmarshal(body, &request); err != nil {
		writeMeshkitError(w, models.ErrUnmarshal(err, "prometheus batch query"), http.StatusBadRequest)
		return
	}
	if len(request.Queries) == 0 {
		writeMeshkitError(w, ErrRequestBody(fmt.Errorf("no queries supplied")), http.StatusBadRequest)
		return
	}
	if len(request.Queries) > maxPrometheusBatchQueries {
		writeMeshkitError(w, ErrRequestBody(fmt.Errorf("too many queries: %d (max %d)", len(request.Queries), maxPrometheusBatchQueries)), http.StatusBadRequest)
		return
	}

	client, _, statusCode, err := h.prometheusTelemetryClientForConnection(token, connectionID, p)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}

	queries := make([]prometheus.BatchQuery, len(request.Queries))
	for i, q := range request.Queries {
		queries[i] = prometheus.BatchQuery{ID: q.ID, Query: q.Query}
	}

	results := client.QueryRangeBatch(req.Context(), request.Start, request.End, request.Step, queries)
	items := make([]prometheusBatchResultItem, len(results))
	for i, r := range results {
		item := prometheusBatchResultItem{ID: r.ID}
		if r.Err != nil {
			item.Error = ErrTelemetryPrometheus(r.Err, "datasource query").Error()
		} else {
			item.Response = r.Body
		}
		items[i] = item
	}
	h.writeJSON(w, map[string]interface{}{"results": items})
}

// prometheusQueryError classifies a query failure into a specific, user-facing
// error and the HTTP status to return with it.
func prometheusQueryError(err error) (int, error) {
	switch prometheus.StatusCode(err) {
	case http.StatusUnauthorized, http.StatusForbidden:
		return http.StatusBadGateway, ErrTelemetryPrometheusAuth(err)
	}
	return http.StatusBadGateway, ErrTelemetryPrometheus(err, "datasource query")
}

// PrometheusTelemetryPanelsHandler reads (GET) or replaces (POST) the set of
// saved metric panels for a connection. The selection is persisted inside the
// connection's Metadata, so no dedicated table is required.
//
// GET|POST /api/telemetry/prometheus/{connectionID}/panels
func (h *Handler) PrometheusTelemetryPanelsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	connection, statusCode, err := p.GetConnectionByID(token, connectionID)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}
	if connection.Kind != "prometheus" {
		writeMeshkitError(w, ErrInvalidConnectionKind(connection.Kind, "prometheus"), http.StatusBadRequest)
		return
	}

	if req.Method == http.MethodGet {
		h.writeJSON(w, readMetricPanels(connection))
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeMeshkitError(w, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	panels := []prometheus.MetricPanel{}
	if err := json.Unmarshal(body, &panels); err != nil {
		writeMeshkitError(w, models.ErrUnmarshal(err, "prometheus panels"), http.StatusBadRequest)
		return
	}

	if connection.Metadata == nil {
		connection.Metadata = core.Map{}
	}
	if len(panels) > 0 {
		connection.Metadata[prometheus.MetadataPanelsKey] = panels
	} else {
		delete(connection.Metadata, prometheus.MetadataPanelsKey)
	}

	// Persist via the id-based update (mirrors the connections UI). The
	// kind-based provider.UpdateConnection puts connection.Kind in the remote
	// URL where a UUID is expected, which 500s on remote providers.
	payload := &connections.ConnectionPayload{
		ID:           connection.ID,
		Kind:         connection.Kind,
		SubType:      connection.SubType,
		Type:         connection.ConnectionType,
		Name:         connection.Name,
		Status:       connection.Status,
		MetaData:     connection.Metadata,
		CredentialID: connection.CredentialID,
	}
	if _, err := p.UpdateConnectionById(token, payload, connectionID.String()); err != nil {
		writeMeshkitError(w, ErrUpdateConnection(err), http.StatusInternalServerError)
		return
	}
	h.writeJSON(w, panels)
}

// readMetricPanels extracts the saved metric-panel list from a connection's
// Metadata, normalizing whatever shape it round-tripped through the DB into a
// typed slice.
func readMetricPanels(connection *connections.Connection) []prometheus.MetricPanel {
	panels := []prometheus.MetricPanel{}
	raw, ok := connection.Metadata[prometheus.MetadataPanelsKey]
	if !ok || raw == nil {
		return panels
	}
	// Metadata values come back as generic interfaces after persistence, so
	// re-marshal and decode into the typed slice.
	b, err := json.Marshal(raw)
	if err != nil {
		return panels
	}
	_ = json.Unmarshal(b, &panels)
	return panels
}
