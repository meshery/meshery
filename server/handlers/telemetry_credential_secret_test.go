package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/events"
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

func newTelemetryDBFixture(t *testing.T) (*Handler, *models.DefaultLocalProvider) {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if err := db.AutoMigrate(connections.Connection{}, models.Credential{}, events.Event{}); err != nil {
		t.Fatalf("migrate tables: %v", err)
	}

	provider := &models.DefaultLocalProvider{
		ConnectionPersister: &models.ConnectionPersister{DB: &db},
		EventsPersister:     &models.EventsPersister{DB: &db},
		GenericPersister:    &db,
	}

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	return h, provider
}

// TestTelemetryClientForConnection_AnonymousZeroCredentialID proves that
// anonymous Prometheus and Grafana connections work end-to-end when stored
// without a credential. DefaultLocalProvider.SaveConnection persists a nil
// CredentialID as &uuid.Nil for compatibility; telemetry client resolution must
// treat both nil and uuid.Nil as anonymous rather than looking up the nil UUID
// in the credentials table and failing with 404.
func TestTelemetryClientForConnection_AnonymousZeroCredentialID(t *testing.T) {
	kinds := []string{"prometheus", "grafana"}
	modes := []string{
		"saved via DefaultLocalProvider.SaveConnection",
		"explicit nil CredentialID",
		"explicit zero-UUID CredentialID",
	}

	for _, kind := range kinds {
		for _, mode := range modes {
			t.Run(kind+"/"+mode, func(t *testing.T) {
				server, gotHdr := newAuthCapturingServer(t)
				h, provider := newTelemetryDBFixture(t)

				var connID uuid.UUID
				switch mode {
				case "saved via DefaultLocalProvider.SaveConnection":
					saved, err := provider.SaveConnection(&connections.ConnectionPayload{
						Name:     kind + "-anonymous",
						Kind:     kind,
						Type:     "telemetry",
						SubType:  "metrics",
						Status:   connections.CONNECTED,
						MetaData: map[string]interface{}{"url": server.URL},
					}, "", false)
					if err != nil {
						t.Fatalf("SaveConnection: %v", err)
					}
					if saved.CredentialID == nil || !saved.CredentialID.IsNil() {
						t.Fatalf("expected DefaultLocalProvider.SaveConnection to store zero UUID sentinel, got %v", saved.CredentialID)
					}
					connID = saved.ID
				case "explicit nil CredentialID":
					connID = uuid.Must(uuid.NewV4())
					_, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
						ID:             connID,
						Name:           kind + "-nil-cred",
						Kind:           kind,
						ConnectionType: "telemetry",
						SubType:        "metrics",
						Status:         connections.CONNECTED,
						CredentialID:   nil,
						Metadata:       core.Map{"url": server.URL},
					})
					if err != nil {
						t.Fatalf("ConnectionPersister.SaveConnection: %v", err)
					}
				case "explicit zero-UUID CredentialID":
					connID = uuid.Must(uuid.NewV4())
					zeroID := core.Uuid(uuid.Nil)
					_, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
						ID:             connID,
						Name:           kind + "-zero-cred",
						Kind:           kind,
						ConnectionType: "telemetry",
						SubType:        "metrics",
						Status:         connections.CONNECTED,
						CredentialID:   &zeroID,
						Metadata:       core.Map{"url": server.URL},
					})
					if err != nil {
						t.Fatalf("ConnectionPersister.SaveConnection: %v", err)
					}
				}

				user := &models.User{ID: uuid.Must(uuid.NewV4())}
				req := httptest.NewRequest(http.MethodGet, "/api/telemetry/"+kind+"/"+connID.String()+"/ping", nil)
				req = mux.SetURLVars(req, map[string]string{"connectionID": connID.String()})
				rec := httptest.NewRecorder()

				if kind == "prometheus" {
					h.PrometheusTelemetryPingHandler(rec, req, nil, user, provider)
				} else {
					h.GrafanaTelemetryPingHandler(rec, req, nil, user, provider)
				}

				if rec.Code != http.StatusOK {
					t.Fatalf("ping status = %d, want 200 (body=%q)", rec.Code, rec.Body.String())
				}
				var payload struct {
					Reachable bool `json:"reachable"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
					t.Fatalf("decode ping response %q: %v", rec.Body.String(), err)
				}
				if !payload.Reachable {
					t.Fatalf("expected reachable=true, got body=%q", rec.Body.String())
				}
				if *gotHdr != "" {
					t.Fatalf("Authorization = %q, want empty header for anonymous connection", *gotHdr)
				}
			})
		}
	}
}

// TestTelemetryClientForConnection_NonZeroMissingCredentialStillReturns404
// ensures that skipping credential lookup for uuid.Nil does not mask a real
// non-zero CredentialID that references a deleted or missing credential.
func TestTelemetryClientForConnection_NonZeroMissingCredentialStillReturns404(t *testing.T) {
	for _, kind := range []string{"prometheus", "grafana"} {
		t.Run(kind, func(t *testing.T) {
			h, provider := newTelemetryDBFixture(t)
			connID := uuid.Must(uuid.NewV4())
			missingCredID := core.Uuid(uuid.Must(uuid.NewV4()))

			_, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
				ID:             connID,
				Name:           kind + "-missing-cred",
				Kind:           kind,
				ConnectionType: "telemetry",
				SubType:        "metrics",
				Status:         connections.CONNECTED,
				CredentialID:   &missingCredID,
				Metadata:       core.Map{"url": "http://127.0.0.1:9999"},
			})
			if err != nil {
				t.Fatalf("SaveConnection: %v", err)
			}

			var (
				statusCode int
				clientErr  error
			)
			if kind == "prometheus" {
				_, _, statusCode, clientErr = h.prometheusTelemetryClientForConnection("", connID, provider)
			} else {
				_, _, statusCode, clientErr = h.grafanaClientForConnection("", connID, provider)
			}

			if clientErr == nil {
				t.Fatal("expected error for missing non-zero credential ID, got nil")
			}
			if statusCode != http.StatusNotFound {
				t.Fatalf("statusCode = %d, want 404 (err=%v)", statusCode, clientErr)
			}
		})
	}
}
