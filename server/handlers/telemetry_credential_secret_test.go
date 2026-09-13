package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
)

// Meshery has to read every credential `secret` shape that exists in production,
// because Layer5 Cloud is moving to the canonical shape meshery/schemas declares
// while the legacy rows stay exactly as they were written. These tests pin the
// telemetry read path - the one that turns a stored credential into the
// Authorization header sent to Grafana/Prometheus - against all of them.
//
// The shape catalogue and the resolution rules live in
// server/models/credential_secret.go.

// newTelemetryCredentialFixture wires a handler over an in-memory database with
// a single connection of the given kind pointing at fakeURL and backed by a
// credential holding secret.
func newTelemetryCredentialFixture(
	t *testing.T,
	kind, fakeURL string,
	secret map[string]interface{},
) (*Handler, models.Provider, uuid.UUID) {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, models.Credential{}); err != nil {
		t.Fatalf("migrate tables: %v", err)
	}

	provider := &models.DefaultLocalProvider{
		ConnectionPersister: &models.ConnectionPersister{DB: &db},
		GenericPersister:    &db,
	}

	credential, err := provider.SaveUserCredential("", &models.Credential{
		Name:   "telemetry-cred",
		Type:   kind,
		Secret: secret,
	})
	if err != nil {
		t.Fatalf("save credential: %v", err)
	}

	connectionID := uuid.Must(uuid.NewV4())
	credentialID := credential.ID
	connection := &connections.Connection{
		ID:             connectionID,
		Name:           kind + "-connection",
		Kind:           kind,
		ConnectionType: "telemetry",
		SubType:        "metrics",
		Status:         connections.CONNECTED,
		CredentialID:   &credentialID,
		Metadata:       core.Map{"url": fakeURL},
	}
	if _, err := provider.ConnectionPersister.SaveConnection(connection); err != nil {
		t.Fatalf("save connection: %v", err)
	}

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	return h, provider, connectionID
}

// newAuthCapturingServer stands in for a Grafana/Prometheus instance and records
// the Authorization header of the last request it served.
func newAuthCapturingServer(t *testing.T) (*httptest.Server, *string) {
	t.Helper()

	var got string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":"11.0.0","database":"ok","status":"success"}`))
	}))
	t.Cleanup(server.Close)
	return server, &got
}

// TestGrafanaClientForConnectionCredentialShapes proves the Grafana telemetry
// read path resolves the API key out of every persisted credential shape. Before
// the tolerant read it reached straight into `secret["secret"]`, so the canonical
// shape (`secret` object IS the payload) and the double-nested shape Meshery UI
// writes both produced an unauthenticated client.
func TestGrafanaClientForConnectionCredentialShapes(t *testing.T) {
	tests := []struct {
		name    string
		secret  map[string]interface{}
		wantHdr string
	}{
		{
			name:    "canonical shape",
			secret:  map[string]interface{}{"grafanaURL": "https://grafana.example", "grafanaAPIKey": "canonical-key"},
			wantHdr: "Bearer canonical-key",
		},
		{
			name: "legacy double-nested shape",
			secret: map[string]interface{}{
				"credentialName": "grafana-cred",
				"secret":         map[string]interface{}{"grafanaURL": "https://grafana.example", "grafanaAPIKey": "nested-key"},
			},
			wantHdr: "Bearer nested-key",
		},
		{
			name:    "legacy string shape",
			secret:  map[string]interface{}{"secret": "legacy-key"},
			wantHdr: "Bearer legacy-key",
		},
		{
			name:    "anonymous credential",
			secret:  map[string]interface{}{"grafanaURL": "https://grafana.example"},
			wantHdr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, gotHdr := newAuthCapturingServer(t)
			h, provider, connectionID := newTelemetryCredentialFixture(t, "grafana", server.URL, tt.secret)

			client, _, _, err := h.grafanaClientForConnection("", connectionID, provider)
			if err != nil {
				t.Fatalf("grafanaClientForConnection: %v", err)
			}
			if _, err := client.Health(context.Background()); err != nil {
				t.Fatalf("health: %v", err)
			}
			if *gotHdr != tt.wantHdr {
				t.Fatalf("Authorization = %q, want %q", *gotHdr, tt.wantHdr)
			}
		})
	}
}

// TestPrometheusClientForConnectionCredentialShapes is the Prometheus half of
// TestGrafanaClientForConnectionCredentialShapes. The canonical Prometheus
// credential form carries no auth field, so a canonical credential is correctly
// anonymous; the legacy shapes must keep authenticating.
func TestPrometheusClientForConnectionCredentialShapes(t *testing.T) {
	tests := []struct {
		name    string
		secret  map[string]interface{}
		wantHdr string
	}{
		{
			name:    "legacy string shape",
			secret:  map[string]interface{}{"secret": "legacy-key"},
			wantHdr: "Bearer legacy-key",
		},
		{
			name:    "legacy string shape with basic auth",
			secret:  map[string]interface{}{"credentialName": "prom-cred", "secret": "user:pass"},
			wantHdr: "Basic dXNlcjpwYXNz",
		},
		{
			name:    "canonical shape is anonymous",
			secret:  map[string]interface{}{"prometheusURL": "https://prom.example"},
			wantHdr: "",
		},
		{
			name: "legacy double-nested shape is anonymous",
			secret: map[string]interface{}{
				"credentialName": "prom-cred",
				"secret":         map[string]interface{}{"prometheusURL": "https://prom.example"},
			},
			wantHdr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, gotHdr := newAuthCapturingServer(t)
			h, provider, connectionID := newTelemetryCredentialFixture(t, "prometheus", server.URL, tt.secret)

			client, _, _, err := h.prometheusTelemetryClientForConnection("", connectionID, provider)
			if err != nil {
				t.Fatalf("prometheusTelemetryClientForConnection: %v", err)
			}
			if _, err := client.Health(context.Background()); err != nil {
				t.Fatalf("health: %v", err)
			}
			if *gotHdr != tt.wantHdr {
				t.Fatalf("Authorization = %q, want %q", *gotHdr, tt.wantHdr)
			}
		})
	}
}
