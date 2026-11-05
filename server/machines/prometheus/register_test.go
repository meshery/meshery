package prometheus

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/schemas/models/core"
)

func registerTestContext(t *testing.T) context.Context {
	t.Helper()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("user uuid: %v", err)
	}
	sysID := core.Uuid(uuid.Nil)
	ctx := context.Background()
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: core.Uuid(userID)})
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	return ctx
}

// TestRegisterAction_ExecuteReturnsNoOpOnHealthyEndpoint is the wizard success
// path: probe Prometheus, then settle. Returning Exit used to be misread as a
// follow-up event with no edge from REGISTERED.
func TestRegisterAction_ExecuteReturnsNoOpOnHealthyEndpoint(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/status/buildinfo" {
			t.Errorf("unexpected path %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]string{"version": "2.53.0"},
		})
	}))
	t.Cleanup(srv.Close)

	next, event, err := (&RegisterAction{}).Execute(registerTestContext(t), nil, connections.ConnectionPayload{
		MetaData: map[string]interface{}{"url": srv.URL},
		CredentialSecret: map[string]interface{}{
			"secret": "",
		},
	})
	if err != nil {
		t.Fatalf("expected healthy register to succeed, got %v", err)
	}
	if event != nil {
		t.Fatalf("expected no error event on success, got severity %q", event.Severity)
	}
	if next != machines.NoOp {
		t.Fatalf("expected NoOp so the machine stays REGISTERED, got %q", next)
	}
}
