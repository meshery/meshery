package grafana

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/schemas/models/core"
)

func registerTestContext() context.Context {
	userID := uuid.Must(uuid.NewV4())
	sysID := uuid.Must(uuid.NewV4())
	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: userID})
	return context.WithValue(ctx, models.SystemIDKey, &sysID)
}

func TestRegisterAction_ResolvesAllCredentialShapes(t *testing.T) {
	basicHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:secret"))

	cases := []struct {
		name       string
		credSecret map[string]interface{}
		wantAuth   string
	}{
		{
			name:       "canonical grafanaAPIKey",
			credSecret: map[string]interface{}{"grafanaURL": "http://grafana:3000", "grafanaAPIKey": "glsa_canonical"},
			wantAuth:   "Bearer glsa_canonical",
		},
		{
			name: "legacy nested with grafanaAPIKey",
			credSecret: map[string]interface{}{
				"credentialName": "my-grafana",
				"secret":         map[string]interface{}{"grafanaAPIKey": "glsa_nested"},
			},
			wantAuth: "Bearer glsa_nested",
		},
		{
			name: "legacy nested with string secret (basic auth)",
			credSecret: map[string]interface{}{
				"name":   "my-grafana",
				"secret": map[string]interface{}{"secret": "admin:secret"},
			},
			wantAuth: basicHeader,
		},
		{
			name:       "legacy flat string secret",
			credSecret: map[string]interface{}{"name": "my-grafana", "secret": "glsa_flat"},
			wantAuth:   "Bearer glsa_flat",
		},
		{
			name:       "anonymous nil credential secret",
			credSecret: nil,
			wantAuth:   "",
		},
		{
			name:       "anonymous empty credential secret",
			credSecret: map[string]interface{}{},
			wantAuth:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotAuth string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/health" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				gotAuth = r.Header.Get("Authorization")
				if gotAuth != tc.wantAuth {
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = w.Write([]byte(`{"message":"unauthorized"}`))
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"database":"ok","version":"10.4.2"}`))
			}))
			defer srv.Close()

			payload := connections.ConnectionPayload{
				Kind:             "grafana",
				Name:             "test-grafana",
				MetaData:         core.Map{"url": srv.URL},
				CredentialSecret: tc.credSecret,
			}

			action := &RegisterAction{}
			eventType, event, err := action.Execute(registerTestContext(), nil, payload)
			if err != nil {
				t.Fatalf("Execute failed: %v (event=%+v)", err, event)
			}
			if eventType != machines.Exit {
				t.Errorf("eventType = %q, want %q", eventType, machines.Exit)
			}
			if gotAuth != tc.wantAuth {
				t.Errorf("Authorization header = %q, want %q", gotAuth, tc.wantAuth)
			}
		})
	}
}

func TestRegisterAction_RejectsUnauthorizedCredential(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"invalid api key"}`))
	}))
	defer srv.Close()

	payload := connections.ConnectionPayload{
		Kind:             "grafana",
		Name:             "test-grafana",
		MetaData:         core.Map{"url": srv.URL},
		CredentialSecret: map[string]interface{}{"grafanaAPIKey": "bad-token"},
	}

	action := &RegisterAction{}
	eventType, event, err := action.Execute(registerTestContext(), nil, payload)
	if err == nil {
		t.Fatal("expected error for 401 health response, got nil")
	}
	if eventType != machines.NoOp {
		t.Errorf("eventType = %q, want %q", eventType, machines.NoOp)
	}
	if event == nil {
		t.Error("expected error event to be returned")
	}
}
