package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	grafana "github.com/meshery/meshery/server/models/telemetry/grafana"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// telemetry_grafana_handler.go exposes a clean, read-oriented API over a Grafana
// instance registered as a first-class Meshery connection. It is independent of
// the legacy /api/telemetry/metrics/grafana/* handlers and reuses the connection
// + credential plumbing instead of its own config flow. Routes live under
// /api/telemetry/grafana/{connectionID}/*.

// grafanaClientForConnection resolves a Grafana connection and its credential
// into a ready-to-use grafana.Client. It enforces that the connection is of
// kind "grafana" and tolerates connections without a credential (anonymous
// Grafana). The returned int is an HTTP status to use when err is non-nil.
func (h *Handler) grafanaClientForConnection(token string, connectionID uuid.UUID, p models.Provider) (*grafana.Client, *connections.Connection, int, error) {
	connection, statusCode, err := p.GetConnectionByID(token, connectionID)
	if err != nil {
		return nil, nil, statusCode, err
	}
	if connection.Kind != "grafana" {
		return nil, connection, http.StatusBadRequest, ErrInvalidConnectionKind(connection.Kind, "grafana")
	}

	baseURL, _ := connection.Metadata["url"].(string)
	secret := ""
	if connection.CredentialID != nil {
		cred, sc, cerr := p.GetCredentialByID(token, core.UUIDOrUUIDNil(connection.CredentialID))
		if cerr != nil {
			return nil, connection, sc, cerr
		}
		secret, _ = cred.Secret["secret"].(string)
	}
	return grafana.New(baseURL, secret, h.log), connection, http.StatusOK, nil
}

// writeJSON serializes v as the JSON response body with a 200 status.
func (h *Handler) writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// The body may already be partially written, so a fresh error response
		// would corrupt the stream — log only.
		h.log.Error(models.ErrMarshal(err, "telemetry grafana payload"))
	}
}

// GrafanaTelemetryPingHandler verifies a Grafana connection is reachable and its
// credential is accepted, emitting a notification event with the outcome.
//
// GET /api/telemetry/grafana/{connectionID}/ping
func (h *Handler) GrafanaTelemetryPingHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	userID := user.ID
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	eventBuilder := events.NewEvent().
		ActedUpon(connectionID).
		FromUser(userID).
		FromSystem(*h.SystemID).
		WithCategory("connection").
		WithAction("ping")

	client, connection, statusCode, err := h.grafanaClientForConnection(token, connectionID, p)
	if err != nil {
		h.emitGrafanaPingFailure(p, token, userID, eventBuilder, connectionLabel(connection, connectionID), err)
		writeMeshkitError(w, err, statusCode)
		return
	}

	info, err := client.Health(req.Context())
	if err != nil {
		_err := ErrTelemetryGrafana(err, "ping")
		h.emitGrafanaPingFailure(p, token, userID, eventBuilder, connection.Name, _err)
		writeMeshkitError(w, _err, http.StatusBadGateway)
		return
	}

	description := fmt.Sprintf("Grafana connection %q is reachable.", connection.Name)
	if info.Version != "" {
		description = fmt.Sprintf("Grafana connection %q is reachable (Grafana %s).", connection.Name, info.Version)
	}
	metadata := map[string]interface{}{
		"version":  info.Version,
		"database": info.Database,
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
		"database":  info.Database,
	})
}

// emitGrafanaPingFailure persists and broadcasts an error event describing a
// failed Grafana connection ping.
func (h *Handler) emitGrafanaPingFailure(p models.Provider, token string, userID core.Uuid, eventBuilder *events.EventBuilder, label string, err error) {
	event := eventBuilder.
		WithSeverity(events.Error).
		WithDescription(fmt.Sprintf("Grafana connection %q is unreachable.", label)).
		WithMetadata(map[string]interface{}{"error": err}).
		Build()
	_ = p.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)
	h.log.Error(err)
}

// connectionLabel returns a human-friendly label for a connection, falling back
// to its ID when the connection could not be loaded.
func connectionLabel(connection *connections.Connection, connectionID uuid.UUID) string {
	if connection != nil && connection.Name != "" {
		return connection.Name
	}
	return connectionID.String()
}

// GrafanaTelemetryBoardsHandler searches the Grafana instance's dashboards.
//
// GET /api/telemetry/grafana/{connectionID}/boards?search=
func (h *Handler) GrafanaTelemetryBoardsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	client, _, statusCode, err := h.grafanaClientForConnection(token, connectionID, p)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}
	boards, err := client.SearchBoards(req.Context(), req.URL.Query().Get("search"))
	if err != nil {
		writeMeshkitError(w, ErrTelemetryGrafana(err, "board search"), http.StatusBadGateway)
		return
	}
	h.writeJSON(w, boards)
}

// GrafanaTelemetryBoardHandler returns a single dashboard's panels.
//
// GET /api/telemetry/grafana/{connectionID}/boards/{uid}
func (h *Handler) GrafanaTelemetryBoardHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	client, _, statusCode, err := h.grafanaClientForConnection(token, connectionID, p)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}
	board, err := client.GetBoard(req.Context(), mux.Vars(req)["uid"])
	if err != nil {
		writeMeshkitError(w, ErrTelemetryGrafana(err, "board fetch"), http.StatusBadGateway)
		return
	}
	h.writeJSON(w, board)
}

// GrafanaTelemetryDatasourcesHandler lists the Grafana instance's datasources.
//
// GET /api/telemetry/grafana/{connectionID}/datasources
func (h *Handler) GrafanaTelemetryDatasourcesHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	client, _, statusCode, err := h.grafanaClientForConnection(token, connectionID, p)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}
	ds, err := client.ListDatasources(req.Context())
	if err != nil {
		writeMeshkitError(w, ErrTelemetryGrafana(err, "datasource listing"), http.StatusBadGateway)
		return
	}
	h.writeJSON(w, ds)
}

// GrafanaTelemetryQueryRangeHandler proxies a Prometheus-style range query
// through Grafana's datasource proxy and streams the response back verbatim for
// the UI to render.
//
// GET /api/telemetry/grafana/{connectionID}/query_range?ds=&query=&start=&end=&step=
func (h *Handler) GrafanaTelemetryQueryRangeHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	client, _, statusCode, err := h.grafanaClientForConnection(token, connectionID, p)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}

	q := req.URL.Query()
	dsUID := q.Get("ds")
	params := url.Values{}
	params.Set("query", q.Get("query"))
	params.Set("start", q.Get("start"))
	params.Set("end", q.Get("end"))
	params.Set("step", q.Get("step"))

	body, err := client.QueryRange(req.Context(), dsUID, params)
	if err != nil {
		writeMeshkitError(w, ErrTelemetryGrafana(err, "datasource query"), http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

// GrafanaTelemetryPinnedBoardsHandler reads (GET) or replaces (POST) the set of
// boards a user has "added" to a connection. The selection is persisted inside
// the connection's Metadata, so no dedicated table is required.
//
// GET|POST /api/telemetry/grafana/{connectionID}/pinned
func (h *Handler) GrafanaTelemetryPinnedBoardsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	connection, statusCode, err := p.GetConnectionByID(token, connectionID)
	if err != nil {
		writeMeshkitError(w, err, statusCode)
		return
	}
	if connection.Kind != "grafana" {
		writeMeshkitError(w, ErrInvalidConnectionKind(connection.Kind, "grafana"), http.StatusBadRequest)
		return
	}

	if req.Method == http.MethodGet {
		h.writeJSON(w, readPinnedBoards(connection))
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeMeshkitError(w, ErrRequestBody(err), http.StatusBadRequest)
		return
	}
	boards := []grafana.PinnedBoard{}
	if err := json.Unmarshal(body, &boards); err != nil {
		writeMeshkitError(w, models.ErrUnmarshal(err, "pinned boards"), http.StatusBadRequest)
		return
	}

	if connection.Metadata == nil {
		connection.Metadata = core.Map{}
	}
	if len(boards) > 0 {
		connection.Metadata[grafana.MetadataPinnedBoardsKey] = boards
	} else {
		delete(connection.Metadata, grafana.MetadataPinnedBoardsKey)
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
	h.writeJSON(w, boards)
}

// readPinnedBoards extracts the pinned-board list from a connection's Metadata,
// normalizing whatever shape it round-tripped through the DB into a typed slice.
func readPinnedBoards(connection *connections.Connection) []grafana.PinnedBoard {
	boards := []grafana.PinnedBoard{}
	raw, ok := connection.Metadata[grafana.MetadataPinnedBoardsKey]
	if !ok || raw == nil {
		return boards
	}
	// Metadata values come back as generic interfaces after persistence, so
	// re-marshal and decode into the typed slice.
	b, err := json.Marshal(raw)
	if err != nil {
		return boards
	}
	_ = json.Unmarshal(b, &boards)
	return boards
}
