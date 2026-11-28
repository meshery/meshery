package prometheus

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
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
)

func registerTestContext() context.Context {
	userID := uuid.Must(uuid.NewV4())
	sysID := uuid.Must(uuid.NewV4())
	ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: userID})
	return context.WithValue(ctx, models.SystemIDKey, &sysID)
}

func TestRegisterAction_ResolvesAllCredentialShapes(t *testing.T) {
	basicHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("prom:secret"))

	cases := []struct {
		name       string
		credSecret map[string]interface{}
		wantAuth   string
	}{
		{
			name:       "canonical anonymous prometheusURL",
			credSecret: map[string]interface{}{"prometheusURL": "http://prometheus:9090"},
			wantAuth:   "",
		},
		{
			name: "legacy nested anonymous prometheusURL",
			credSecret: map[string]interface{}{
				"credentialName": "my-prom",
				"secret":         map[string]interface{}{"prometheusURL": "http://prometheus:9090"},
			},
			wantAuth: "",
		},
		{
			name: "legacy nested with string secret (basic auth)",
			credSecret: map[string]interface{}{
				"name":   "my-prom",
				"secret": map[string]interface{}{"secret": "prom:secret"},
			},
			wantAuth: basicHeader,
		},
		{
			name:       "legacy flat string secret (bearer token)",
			credSecret: map[string]interface{}{"name": "my-prom", "secret": "prom_token"},
			wantAuth:   "Bearer prom_token",
		},
		{
			name:       "anonymous nil credential secret",
			credSecret: nil,
			wantAuth:   "",
		},
		{
			name:       "malformed non-string secret falls back to anonymous",
			credSecret: map[string]interface{}{"secret": 123},
			wantAuth:   "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotAuth string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/status/buildinfo" {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				gotAuth = r.Header.Get("Authorization")
				if gotAuth != tc.wantAuth {
					w.WriteHeader(http.StatusUnauthorized)
					_, _ = w.Write([]byte(`{"status":"error"}`))
					return
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"status":"success","data":{"version":"2.50.0"}}`))
			}))
			defer srv.Close()

			payload := connections.ConnectionPayload{
				Kind:             "prometheus",
				Name:             "test-prom",
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
		_, _ = w.Write([]byte(`{"status":"error"}`))
	}))
	defer srv.Close()

	payload := connections.ConnectionPayload{
		Kind:             "prometheus",
		Name:             "test-prom",
		MetaData:         core.Map{"url": srv.URL},
		CredentialSecret: map[string]interface{}{"secret": "bad-token"},
	}

	action := &RegisterAction{}
	eventType, event, err := action.Execute(registerTestContext(), nil, payload)
	if err == nil {
		t.Fatal("expected error for 401 health response, got nil")
	}
	if code := errors.GetCode(err); code != models.ErrPrometheusScanCode {
		t.Errorf("error code = %q, want %q", code, models.ErrPrometheusScanCode)
	}
	if eventType != machines.NoOp {
		t.Errorf("eventType = %q, want %q", eventType, machines.NoOp)
	}
	if event == nil {
		t.Error("expected error event to be returned")
	}
}
