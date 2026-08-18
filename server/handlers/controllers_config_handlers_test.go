package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// The four controllers-configuration endpoints are the only way the layered
// document is written and read, so this file covers them at the wire: what
// they reject and with which status code, what shape the layered GET returns,
// and what an empty document does at each layer. The cluster-propagation half
// of the same feature is covered by
// server/models/controllers_config_apply_test.go; here the state-machine
// tracker is deliberately nil so an apply is a no-op and every assertion is
// about storage, layering and validation - which is exactly the part that must
// hold with no cluster in reach.

// controllersConfigFakeProvider is a connection store standing in for a
// provider. It hands out deep copies on read, so a handler that mutates the
// metadata map it was given cannot silently mutate the "persisted" state - the
// same separation a real provider has across the wire, and the thing that
// makes the inherit round-trip assertions meaningful.
type controllersConfigFakeProvider struct {
	*models.DefaultLocalProvider
	stored map[core.Uuid]*connections.Connection
	events []events.Event
}

func newControllersConfigFakeProvider() *controllersConfigFakeProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &controllersConfigFakeProvider{
		DefaultLocalProvider: base,
		stored:               map[core.Uuid]*connections.Connection{},
	}
}

func (p *controllersConfigFakeProvider) GetConnectionByID(_ string, connectionID core.Uuid) (*connections.Connection, int, error) {
	conn, ok := p.stored[connectionID]
	if !ok {
		return nil, http.StatusNotFound, connections.ErrControllersConfigInvalid("connection not found")
	}
	return copyConnection(conn), http.StatusOK, nil
}

func (p *controllersConfigFakeProvider) UpdateConnectionById(_ string, payload *connections.ConnectionPayload, connID string) (*connections.Connection, error) {
	id, err := uuid.FromString(connID)
	if err != nil {
		return nil, err
	}
	conn, ok := p.stored[id]
	if !ok {
		return nil, connections.ErrControllersConfigInvalid("connection not found")
	}
	updated := copyConnection(conn)
	updated.Metadata = copyMetadata(payload.MetaData)
	p.stored[id] = updated
	return copyConnection(updated), nil
}

// PersistEvent is recorded rather than written: the handlers emit events on
// every write path and a DefaultLocalProvider with no events persister would
// panic.
func (p *controllersConfigFakeProvider) PersistEvent(event events.Event, _ string) error {
	p.events = append(p.events, event)
	return nil
}

func copyConnection(conn *connections.Connection) *connections.Connection {
	dup := *conn
	dup.Metadata = copyMetadata(conn.Metadata)
	return &dup
}

func copyMetadata(metadata core.Map) core.Map {
	if metadata == nil {
		return nil
	}
	encoded, err := json.Marshal(metadata)
	if err != nil {
		return core.Map{}
	}
	dup := core.Map{}
	if err := json.Unmarshal(encoded, &dup); err != nil {
		return core.Map{}
	}
	return dup
}

func newControllersConfigTestHandler(t *testing.T) (*Handler, *controllersConfigFakeProvider) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}
	// database.Handler embeds a *sync.Mutex that the settings accessors take;
	// composite-literal construction has to supply it or the first Lock()
	// nil-derefs.
	dbHandler := &database.Handler{DB: db, Mutex: &sync.Mutex{}}
	if err := dbHandler.AutoMigrate(&models.SystemSetting{}); err != nil {
		t.Fatalf("failed to migrate system_settings: %v", err)
	}

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config: &models.HandlerConfig{
			EventBroadcaster: models.NewBroadcaster("test"),
		},
		log:       newTestLogger(t),
		dbHandler: dbHandler,
		SystemID:  &systemID,
		// The built-in default is embedded; pinning the server default to
		// operator keeps the seeded connection's materialized mode and the
		// resolved mode in step, so no test here accidentally exercises the
		// mode-change teardown path (which needs a live state machine).
		MeshsyncDefaultDeploymentMode: connections.MeshsyncDeploymentModeOperator,
		// Deliberately nil: an apply to a cluster is then a no-op, so these
		// tests assert only storage, layering and validation.
		ConnectionToStateMachineInstanceTracker: nil,
	}
	return h, newControllersConfigFakeProvider()
}

// seedKubernetesConnection registers a Kubernetes connection whose materialized
// deployment mode already matches the server default, i.e. a connection at rest.
func seedKubernetesConnection(p *controllersConfigFakeProvider, override *controllersconfig.MesheryControllersConfig) core.Uuid {
	id := uuid.Must(uuid.NewV4())
	metadata := core.Map{
		connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeOperator),
	}
	if override != nil {
		_ = connections.SetControllersConfigToMetadata(metadata, override)
	}
	p.stored[id] = &connections.Connection{
		ID:             id,
		Name:           "test-cluster",
		Kind:           "kubernetes",
		ConnectionType: "platform",
		SubType:        "orchestrator",
		Metadata:       metadata,
	}
	return id
}

func testUser() *models.User {
	return &models.User{ID: uuid.Must(uuid.NewV4())}
}

func withToken(req *http.Request) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
}

// putDefaults drives PUT /api/system/controllers/config.
func putDefaults(h *Handler, p models.Provider, body string) *httptest.ResponseRecorder {
	req := withToken(httptest.NewRequest(http.MethodPut, "/api/system/controllers/config", bytes.NewBufferString(body)))
	rec := httptest.NewRecorder()
	h.UpdateControllersDefaultConfig(rec, req, nil, testUser(), p)
	return rec
}

// getDefaults drives GET /api/system/controllers/config.
func getDefaults(h *Handler, p models.Provider) *httptest.ResponseRecorder {
	req := withToken(httptest.NewRequest(http.MethodGet, "/api/system/controllers/config", nil))
	rec := httptest.NewRecorder()
	h.GetControllersDefaultConfig(rec, req, nil, testUser(), p)
	return rec
}

// putConnectionConfig drives
// PUT /api/integrations/connections/{connectionId}/controllers/config through
// mux so the {connectionId} path variable resolves the way it does in the
// router.
func putConnectionConfig(h *Handler, p models.Provider, connectionID string, body string) *httptest.ResponseRecorder {
	req := withToken(httptest.NewRequest(http.MethodPut,
		"/api/integrations/connections/"+connectionID+"/controllers/config", bytes.NewBufferString(body)))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connectionID})
	rec := httptest.NewRecorder()
	h.UpdateConnectionControllersConfig(rec, req, nil, testUser(), p)
	return rec
}

// getConnectionConfig drives
// GET /api/integrations/connections/{connectionId}/controllers/config.
func getConnectionConfig(h *Handler, p models.Provider, connectionID string) *httptest.ResponseRecorder {
	req := withToken(httptest.NewRequest(http.MethodGet,
		"/api/integrations/connections/"+connectionID+"/controllers/config", nil))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connectionID})
	rec := httptest.NewRecorder()
	h.GetConnectionControllersConfig(rec, req, nil, testUser(), p)
	return rec
}

func decodeControllersConfig(t *testing.T, rec *httptest.ResponseRecorder) *controllersconfig.MesheryControllersConfig {
	t.Helper()
	cfg := &controllersconfig.MesheryControllersConfig{}
	if err := json.Unmarshal(rec.Body.Bytes(), cfg); err != nil {
		t.Fatalf("failed to decode controllers config response %q: %v", rec.Body.String(), err)
	}
	return cfg
}

func decodeLayered(t *testing.T, rec *httptest.ResponseRecorder) *controllersconfig.ConnectionControllersConfig {
	t.Helper()
	layered := &controllersconfig.ConnectionControllersConfig{}
	if err := json.Unmarshal(rec.Body.Bytes(), layered); err != nil {
		t.Fatalf("failed to decode layered response %q: %v", rec.Body.String(), err)
	}
	return layered
}

// A server that has never had defaults set must answer with an empty document
// carrying the schema-version discriminator, not with a 404 or a bare `null`:
// clients key the whole editor off "every field absent means built-in".
func TestGetControllersDefaultConfig_UnsetReturnsEmptyStampedDocument(t *testing.T) {
	h, p := newControllersConfigTestHandler(t)

	rec := getDefaults(h, p)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body.String())
	}
	cfg := decodeControllersConfig(t, rec)
	if string(cfg.SchemaVersion) != connections.ControllersConfigSchemaVersion {
		t.Fatalf("schemaVersion = %q, want %q", cfg.SchemaVersion, connections.ControllersConfigSchemaVersion)
	}
	if cfg.Operator != nil || cfg.Meshsync != nil || cfg.Broker != nil {
		t.Fatalf("expected an empty defaults document, got %s", rec.Body.String())
	}
}

// The write endpoint must persist and echo the stored document, and a
// subsequent read must return the same thing: the editor reloads from the GET,
// so an echo that is not what was stored would silently lose the save.
func TestUpdateControllersDefaultConfig_PersistsAndReadsBack(t *testing.T) {
	h, p := newControllersConfigTestHandler(t)

	body := `{"meshsync":{"replicas":3,"debugLogging":true},"broker":{"service":{"type":"NodePort"}}}`
	rec := putDefaults(h, p, body)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body.String())
	}

	for name, got := range map[string]*controllersconfig.MesheryControllersConfig{
		"PUT echo": decodeControllersConfig(t, rec),
		"GET read": decodeControllersConfig(t, getDefaults(h, p)),
	} {
		if got.Meshsync == nil || got.Meshsync.Replicas == nil || *got.Meshsync.Replicas != 3 {
			t.Fatalf("%s: meshsync.replicas not persisted: %+v", name, got.Meshsync)
		}
		if got.Meshsync.DebugLogging == nil || !*got.Meshsync.DebugLogging {
			t.Fatalf("%s: meshsync.debugLogging not persisted: %+v", name, got.Meshsync)
		}
		if got.Broker == nil || got.Broker.Service == nil || got.Broker.Service.Type == nil ||
			*got.Broker.Service.Type != controllersconfig.NodePort {
			t.Fatalf("%s: broker.service.type not persisted: %+v", name, got.Broker)
		}
		if string(got.SchemaVersion) != connections.ControllersConfigSchemaVersion {
			t.Fatalf("%s: schemaVersion = %q, want %q", name, got.SchemaVersion, connections.ControllersConfigSchemaVersion)
		}
	}
}

// Returning every field to Inherit at the server-wide layer means posting an
// empty document, and that must clear the stored defaults rather than leave the
// last-saved values in place. Proven from the store, not just from the
// response, so a handler that echoes an empty document while still holding the
// old one cannot pass.
func TestUpdateControllersDefaultConfig_EmptyDocumentClearsDefaults(t *testing.T) {
	h, p := newControllersConfigTestHandler(t)

	if rec := putDefaults(h, p, `{"meshsync":{"replicas":3}}`); rec.Code != http.StatusOK {
		t.Fatalf("seed status = %d, want 200; body: %s", rec.Code, rec.Body.String())
	}
	stored, err := models.GetControllersConfigDefaults(h.dbHandler)
	if err != nil || stored == nil || stored.Meshsync == nil {
		t.Fatalf("precondition: defaults were not stored (err=%v, stored=%+v)", err, stored)
	}

	rec := putDefaults(h, p, `{}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body.String())
	}

	stored, err = models.GetControllersConfigDefaults(h.dbHandler)
	if err != nil {
		t.Fatalf("reading defaults back: %v", err)
	}
	if stored != nil {
		t.Fatalf("expected the stored defaults to be cleared, got %+v", stored)
	}
	if cfg := decodeControllersConfig(t, rec); cfg.Meshsync != nil || cfg.Broker != nil || cfg.Operator != nil {
		t.Fatalf("expected an empty document in the response, got %s", rec.Body.String())
	}
}

// Every guardrail the schema cannot express is enforced by the server on both
// write endpoints, and rejection is a 400 that names the offending field -
// never a 500 and never a silent accept. Both endpoints share
// readControllersConfigPayload, so both are driven here: a future endpoint that
// forgets to call it fails this table.
func TestControllersConfigWriteEndpoints_RejectInvalidDocuments(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		wantCode int
		wantMsg  string
	}{
		{
			name:     "meshsync.replicas below range",
			body:     `{"meshsync":{"replicas":0}}`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "meshsync.replicas must be between 1 and 10",
		},
		{
			name:     "meshsync.replicas above range",
			body:     `{"meshsync":{"replicas":11}}`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "meshsync.replicas must be between 1 and 10",
		},
		{
			name:     "broker.replicas above range",
			body:     `{"broker":{"replicas":11}}`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "broker.replicas must be between 1 and 10",
		},
		{
			name:     "watchList sets both whitelist and blacklist",
			body:     `{"meshsync":{"watchList":{"whitelist":[{"resource":"pods.v1."}],"blacklist":["deployments.v1.apps"]}}}`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "whitelist and blacklist are mutually exclusive",
		},
		{
			name:     "watchList whitelist entry without a resource",
			body:     `{"meshsync":{"watchList":{"whitelist":[{"events":["ADDED"]}]}}}`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "whitelist entries must set resource",
		},
		{
			name:     "loadBalancerClass without service.type LoadBalancer",
			body:     `{"broker":{"service":{"type":"ClusterIP","loadBalancerClass":"service.k8s.aws/nlb"}}}`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "broker.service.loadBalancerClass is only valid when broker.service.type is LoadBalancer",
		},
		{
			name:     "loadBalancerSourceRanges without service.type LoadBalancer",
			body:     `{"broker":{"service":{"type":"NodePort","loadBalancerSourceRanges":["10.0.0.0/8"]}}}`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "broker.service.loadBalancerSourceRanges is only valid when broker.service.type is LoadBalancer",
		},
		{
			name:     "unknown broker.service.type",
			body:     `{"broker":{"service":{"type":"Headless"}}}`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "broker.service.type must be one of ClusterIP, NodePort, LoadBalancer",
		},
		{
			name:     "unknown operator.deploymentMode",
			body:     `{"operator":{"deploymentMode":"sidecar"}}`,
			wantCode: http.StatusBadRequest,
			// The quoted mode names are JSON-escaped in the error body, so
			// match up to them.
			wantMsg: "operator.deploymentMode must be either",
		},
		{
			name:     "malformed JSON body",
			body:     `{"meshsync":`,
			wantCode: http.StatusBadRequest,
			wantMsg:  "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("server defaults", func(t *testing.T) {
				h, p := newControllersConfigTestHandler(t)
				rec := putDefaults(h, p, tc.body)
				assertRejected(t, rec, tc.wantCode, tc.wantMsg)
				stored, err := models.GetControllersConfigDefaults(h.dbHandler)
				if err != nil {
					t.Fatalf("reading defaults back: %v", err)
				}
				if stored != nil {
					t.Fatalf("a rejected document must not be persisted, got %+v", stored)
				}
			})

			t.Run("connection override", func(t *testing.T) {
				h, p := newControllersConfigTestHandler(t)
				connectionID := seedKubernetesConnection(p, nil)
				rec := putConnectionConfig(h, p, connectionID.String(), tc.body)
				assertRejected(t, rec, tc.wantCode, tc.wantMsg)
				if _, exists := p.stored[connectionID].Metadata[connections.ControllersConfigMetadataKey]; exists {
					t.Fatalf("a rejected document must not be written to the connection's metadata")
				}
			})
		})
	}
}

func assertRejected(t *testing.T, rec *httptest.ResponseRecorder, wantCode int, wantMsg string) {
	t.Helper()
	if rec.Code != wantCode {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, wantCode, rec.Body.String())
	}
	if wantMsg != "" && !strings.Contains(rec.Body.String(), wantMsg) {
		t.Fatalf("expected the response to name the offending field with %q, got: %s", wantMsg, rec.Body.String())
	}
}

// The per-connection GET is the editor's entire data source: it has to expose
// each precedence layer separately (so a field can be shown as inherited or
// overridden) plus the resolved effective document. All three keys are part of
// the contract even when a layer is empty.
func TestGetConnectionControllersConfig_LayeredResponseShape(t *testing.T) {
	h, p := newControllersConfigTestHandler(t)

	if rec := putDefaults(h, p, `{"meshsync":{"version":"v1.0.2"}}`); rec.Code != http.StatusOK {
		t.Fatalf("seeding defaults: status %d, body %s", rec.Code, rec.Body.String())
	}
	replicas := 4
	connectionID := seedKubernetesConnection(p, &controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{Replicas: &replicas},
	})

	rec := getConnectionConfig(h, p, connectionID.String())
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body.String())
	}

	// Assert on the raw JSON keys too: the layer names are the wire contract
	// the UI keys its inherit/override indicators off.
	for _, key := range []string{`"override"`, `"default"`, `"effective"`} {
		if !strings.Contains(rec.Body.String(), key) {
			t.Fatalf("layered response is missing the %s layer: %s", key, rec.Body.String())
		}
	}

	layered := decodeLayered(t, rec)
	if layered.Override == nil || layered.Override.Meshsync == nil || layered.Override.Meshsync.Replicas == nil ||
		*layered.Override.Meshsync.Replicas != 4 {
		t.Fatalf("override layer does not carry the connection's own meshsync.replicas: %s", rec.Body.String())
	}
	if layered.Override.Meshsync.Version != nil {
		t.Fatalf("override layer must carry only what the connection set, not the inherited meshsync.version: %s", rec.Body.String())
	}
	if layered.Default == nil || layered.Default.Meshsync == nil || layered.Default.Meshsync.Version == nil ||
		*layered.Default.Meshsync.Version != "v1.0.2" {
		t.Fatalf("default layer does not carry the server-wide meshsync.version: %s", rec.Body.String())
	}
	// The effective document reports the deployment mode the connection
	// actually runs, which is what clients gate the inert fields on.
	if layered.Effective.Operator == nil || layered.Effective.Operator.DeploymentMode == nil ||
		*layered.Effective.Operator.DeploymentMode != controllersconfig.Operator {
		t.Fatalf("effective document does not report the resolved operator.deploymentMode: %s", rec.Body.String())
	}
}

// Precedence, end to end through the endpoints: the per-connection override
// wins over the server-wide default, which wins over the built-in default.
// Each of the three layers contributes at least one field to `effective`, so a
// regression in any single link of the chain fails here.
func TestGetConnectionControllersConfig_PrecedenceOverrideBeatsDefaultBeatsBuiltIn(t *testing.T) {
	h, p := newControllersConfigTestHandler(t)

	// Server-wide default: meshsync.replicas=3 (to be overridden) and
	// broker.replicas=4 (nothing overrides it, so it must survive).
	if rec := putDefaults(h, p, `{"meshsync":{"replicas":3},"broker":{"replicas":4}}`); rec.Code != http.StatusOK {
		t.Fatalf("seeding defaults: status %d, body %s", rec.Code, rec.Body.String())
	}
	overrideReplicas := 5
	connectionID := seedKubernetesConnection(p, &controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{Replicas: &overrideReplicas},
	})

	layered := decodeLayered(t, getConnectionConfig(h, p, connectionID.String()))

	if layered.Effective.Meshsync == nil || layered.Effective.Meshsync.Replicas == nil ||
		*layered.Effective.Meshsync.Replicas != 5 {
		t.Fatalf("override must beat the server-wide default for meshsync.replicas, got %+v", layered.Effective.Meshsync)
	}
	if layered.Effective.Broker == nil || layered.Effective.Broker.Replicas == nil ||
		*layered.Effective.Broker.Replicas != 4 {
		t.Fatalf("server-wide default must supply broker.replicas when the connection does not, got %+v", layered.Effective.Broker)
	}
	// Nothing sets these at either editable layer, so the built-in default is
	// the only thing that can be answering.
	if layered.Effective.Meshsync.RedactSecrets == nil || *layered.Effective.Meshsync.RedactSecrets {
		t.Fatalf("built-in default must supply meshsync.redactSecrets=false, got %+v", layered.Effective.Meshsync.RedactSecrets)
	}
	if layered.Effective.Broker.Service == nil || layered.Effective.Broker.Service.Type == nil ||
		*layered.Effective.Broker.Service.Type != controllersconfig.ClusterIP {
		t.Fatalf("built-in default must supply broker.service.type=ClusterIP, got %+v", layered.Effective.Broker.Service)
	}
}

// Inherit round-trip: a field set on the connection and then returned to
// Inherit must LEAVE the stored override document, not linger at its last
// value, and the lower layer must apply again. The write endpoint therefore
// has to replace the override wholesale rather than merge into it - a merging
// write would make Inherit unreachable for any field ever set.
func TestUpdateConnectionControllersConfig_InheritRoundTripRemovesField(t *testing.T) {
	h, p := newControllersConfigTestHandler(t)

	if rec := putDefaults(h, p, `{"meshsync":{"replicas":3}}`); rec.Code != http.StatusOK {
		t.Fatalf("seeding defaults: status %d, body %s", rec.Code, rec.Body.String())
	}
	connectionID := seedKubernetesConnection(p, nil)

	// Override meshsync.replicas on the connection.
	if rec := putConnectionConfig(h, p, connectionID.String(), `{"meshsync":{"replicas":7,"debugLogging":true}}`); rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body.String())
	}
	stored := storedOverride(t, p, connectionID)
	if stored == nil || stored.Meshsync == nil || stored.Meshsync.Replicas == nil || *stored.Meshsync.Replicas != 7 {
		t.Fatalf("precondition: meshsync.replicas was not stored as an override, got %+v", stored)
	}

	// Return meshsync.replicas to Inherit while keeping the sibling field set,
	// so this cannot pass by simply clearing the whole override.
	rec := putConnectionConfig(h, p, connectionID.String(), `{"meshsync":{"debugLogging":true}}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body.String())
	}

	stored = storedOverride(t, p, connectionID)
	if stored == nil || stored.Meshsync == nil {
		t.Fatalf("the sibling override must survive, got %+v", stored)
	}
	if stored.Meshsync.Replicas != nil {
		t.Fatalf("meshsync.replicas must leave the stored override document, still present as %d", *stored.Meshsync.Replicas)
	}
	if stored.Meshsync.DebugLogging == nil || !*stored.Meshsync.DebugLogging {
		t.Fatalf("meshsync.debugLogging must remain overridden, got %+v", stored.Meshsync.DebugLogging)
	}

	// The lower layer applies again for the field that was released.
	layered := decodeLayered(t, rec)
	if layered.Effective.Meshsync == nil || layered.Effective.Meshsync.Replicas == nil ||
		*layered.Effective.Meshsync.Replicas != 3 {
		t.Fatalf("the server-wide default must apply again once the override is released, got %+v", layered.Effective.Meshsync)
	}
	if layered.Override == nil || layered.Override.Meshsync == nil || layered.Override.Meshsync.Replicas != nil {
		t.Fatalf("the response's override layer must no longer carry meshsync.replicas: %s", rec.Body.String())
	}
}

// Clearing every field at the connection layer removes the override key from
// the connection's metadata entirely - it does not persist an empty husk that
// later readers would have to special-case - and the server-wide default
// applies again.
func TestUpdateConnectionControllersConfig_EmptyDocumentClearsOverride(t *testing.T) {
	h, p := newControllersConfigTestHandler(t)

	if rec := putDefaults(h, p, `{"meshsync":{"replicas":3}}`); rec.Code != http.StatusOK {
		t.Fatalf("seeding defaults: status %d, body %s", rec.Code, rec.Body.String())
	}
	replicas := 9
	connectionID := seedKubernetesConnection(p, &controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{Replicas: &replicas},
	})

	rec := putConnectionConfig(h, p, connectionID.String(), `{}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rec.Code, rec.Body.String())
	}

	if _, exists := p.stored[connectionID].Metadata[connections.ControllersConfigMetadataKey]; exists {
		t.Fatalf("the override key must be removed from metadata, still present: %+v",
			p.stored[connectionID].Metadata[connections.ControllersConfigMetadataKey])
	}
	layered := decodeLayered(t, rec)
	if layered.Override != nil {
		t.Fatalf("the response must report no override layer, got %+v", layered.Override)
	}
	if layered.Effective.Meshsync == nil || layered.Effective.Meshsync.Replicas == nil ||
		*layered.Effective.Meshsync.Replicas != 3 {
		t.Fatalf("the server-wide default must apply again, got %+v", layered.Effective.Meshsync)
	}
}

// Both per-connection endpoints resolve {connectionId} through the same helper,
// which must reject a malformed id and a non-Kubernetes connection with 400
// rather than reaching the store with uuid.Nil or writing a controllers
// configuration onto, say, a Grafana connection.
func TestConnectionControllersConfigEndpoints_RejectUnusableConnections(t *testing.T) {
	nonKubernetes := func(p *controllersConfigFakeProvider) core.Uuid {
		id := uuid.Must(uuid.NewV4())
		p.stored[id] = &connections.Connection{
			ID: id, Name: "grafana", Kind: "grafana", ConnectionType: "observability", SubType: "metrics",
			Metadata: core.Map{},
		}
		return id
	}

	cases := []struct {
		name    string
		id      func(p *controllersConfigFakeProvider) string
		wantMsg string
	}{
		{
			name:    "malformed connection id",
			id:      func(*controllersConfigFakeProvider) string { return "not-a-uuid" },
			wantMsg: "",
		},
		{
			name:    "non-kubernetes connection",
			id:      func(p *controllersConfigFakeProvider) string { return nonKubernetes(p).String() },
			wantMsg: "Kubernetes connections only",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("GET", func(t *testing.T) {
				h, p := newControllersConfigTestHandler(t)
				rec := getConnectionConfig(h, p, tc.id(p))
				assertRejected(t, rec, http.StatusBadRequest, tc.wantMsg)
				if strings.Contains(rec.Body.String(), uuid.Nil.String()) {
					t.Fatalf("response must not echo the zero UUID: %s", rec.Body.String())
				}
			})
			t.Run("PUT", func(t *testing.T) {
				h, p := newControllersConfigTestHandler(t)
				rec := putConnectionConfig(h, p, tc.id(p), `{"meshsync":{"replicas":2}}`)
				assertRejected(t, rec, http.StatusBadRequest, tc.wantMsg)
			})
		})
	}
}

// storedOverride reads the override back out of the provider's persisted
// connection, which is the only place the round-trip can honestly be observed.
func storedOverride(t *testing.T, p *controllersConfigFakeProvider, connectionID core.Uuid) *controllersconfig.MesheryControllersConfig {
	t.Helper()
	conn, ok := p.stored[connectionID]
	if !ok {
		t.Fatalf("connection %s is not in the store", connectionID)
	}
	cfg, err := connections.ControllersConfigFromMetadata(conn.Metadata)
	if err != nil {
		t.Fatalf("reading the stored override: %v", err)
	}
	return cfg
}

// TestReconcileInheritedDeploymentMode_SkipsConnectionWithUnparseableOverride
// pins the fault reported on #21146: an override that merely fails to parse used
// to be treated as no override at all. Resolution then produced the server-wide
// default, that differed from the materialized mode, and the connection was
// reconciled - tearing down and redeploying MeshSync into a mode the user never
// chose, on the strength of corrupt data. Unreadable is unknown intent, not
// absent intent, so the connection must be left exactly as it is.
func TestReconcileInheritedDeploymentMode_SkipsConnectionWithUnparseableOverride(t *testing.T) {
	h, provider := newControllersConfigTestHandler(t)

	// A connection materialized as embedded, carrying an override that cannot
	// be decoded, while the server-wide default says operator. Treating the
	// override as absent would resolve to operator and reconcile.
	id := uuid.Must(uuid.NewV4())
	connection := &connections.Connection{
		ID:   id,
		Kind: "kubernetes",
		Metadata: core.Map{
			connections.ControllersConfigMetadataKey:      "{not json",
			connections.MeshsyncDeploymentModeMetadataKey: string(connections.MeshsyncDeploymentModeEmbedded),
		},
	}
	provider.stored[id] = copyConnection(connection)

	operatorMode := controllersconfig.MesheryOperatorConfigDeploymentMode(connections.MeshsyncDeploymentModeOperator)
	serverDefaults := &controllersconfig.MesheryControllersConfig{
		Operator: &controllersconfig.MesheryOperatorConfig{DeploymentMode: &operatorMode},
	}

	eventBuilder := events.NewEvent().FromSystem(*h.SystemID).WithCategory("connection").WithAction("update")
	h.reconcileInheritedDeploymentMode(connection, serverDefaults, "token", uuid.Must(uuid.NewV4()), provider, eventBuilder)

	persisted := connections.MeshsyncDeploymentModeFromMetadata(provider.stored[id].Metadata)
	if persisted != connections.MeshsyncDeploymentModeEmbedded {
		t.Fatalf("an unparseable override must not persist a new materialization; stored mode is now %q", persisted)
	}
	if mode := connections.MeshsyncDeploymentModeFromMetadata(connection.Metadata); mode != connections.MeshsyncDeploymentModeEmbedded {
		t.Fatalf("the in-memory materialized mode must be left untouched, got %q", mode)
	}
	if raw := connection.Metadata[connections.ControllersConfigMetadataKey]; raw != "{not json" {
		t.Fatalf("the unreadable override must be left intact for deliberate correction, got %#v", raw)
	}
}
