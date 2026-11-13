package models

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/schemas/models/core"
)

// TestRemoteProviderUpdateConnectionStatusIfDisconnectBetweenReadAndPut
// pins the residual window on the remote path (issue #23248). The remote
// provider has no conditional-update API, so UpdateConnectionStatusIf
// reads the stored status and then PUTs the new one, and a disconnect
// landing between the read and the PUT is overwritten by the PUT: after
// it, the stored status is the Discovery's state, not DISCONNECTED, and
// the reporting-path gate re-reads the status the PUT already wrote, so
// it passes too. The fake remote below lands the disconnect inside the
// PUT handler, after the GET already returned NOTFOUND and before the
// PUT body is applied. This test documents that window as it stands:
// update it, do not delete it, when a conditional-update API lands on
// the remote side.
func TestRemoteProviderUpdateConnectionStatusIfDisconnectBetweenReadAndPut(t *testing.T) {
	connectionID := core.Uuid(uuid.Must(uuid.NewV4()))

	var mu sync.Mutex
	stored := connections.Connection{
		ID:     connectionID,
		Name:   "dead-cluster",
		Kind:   "kubernetes",
		Status: connections.NOTFOUND,
	}
	writeStored := func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(stored); err != nil {
			t.Errorf("encode stored connection: %v", err)
		}
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		switch r.Method {
		case http.MethodGet:
			writeStored(w)
		case http.MethodPut:
			// A disconnect lands between the GET and the PUT: the
			// stored status moves to DISCONNECTED before this PUT
			// is applied, exactly the window the local store closes
			// with one UPDATE ... WHERE status IN (...) statement.
			stored.Status = connections.DISCONNECTED
			var payload struct {
				Status connections.ConnectionStatus `json:"status"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			stored.Status = payload.Status
			writeStored(w)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = persistConnectionCapability()

	fromStatuses := []connections.ConnectionStatus{
		connections.DISCOVERED,
		connections.REGISTERED,
		connections.CONNECTED,
		connections.NOTFOUND,
	}
	applied, updated, err := provider.UpdateConnectionStatusIf("test-token", connectionID, fromStatuses, connections.CONNECTED)
	if err != nil {
		t.Fatalf("UpdateConnectionStatusIf returned an error: %v", err)
	}
	if !applied {
		t.Fatalf("applied = false, want true (the read saw NOTFOUND, so the write went ahead)")
	}
	if updated == nil || updated.Status != connections.CONNECTED {
		t.Fatalf("updated status = %v, want CONNECTED", updated)
	}

	// The intervening disconnect did not survive: the PUT overwrote
	// DISCONNECTED with CONNECTED. This is the residual window on the
	// remote path — it needs a conditional-update API on the remote
	// side to close, tracked as a follow-up to issue #23248.
	mu.Lock()
	final := stored.Status
	mu.Unlock()
	if final != connections.CONNECTED {
		t.Fatalf("stored status after PUT = %q, want %q (the test models the disconnect being overwritten)", final, connections.CONNECTED)
	}
}
