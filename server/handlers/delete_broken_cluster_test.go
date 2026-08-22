package handlers

import (
	"net/http"
	"testing"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models/connections"
)

func TestDeleteConnectionBrokenCluster(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	tracker := &machines.ConnectionToStateMachineInstanceTracker{}
	h.ConnectionToStateMachineInstanceTracker = tracker

	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name:           "broken-connection",
		Kind:           "kubernetes",
		ConnectionType: "platform",
	})
	if err != nil {
		t.Fatalf("save connection: %v", err)
	}

	// Because we do not inject a valid kubeconfig or reachable cluster into this test,
	// AssignInitialCtx will naturally fail. We want to prove that the HTTP request still
	// returns 200 and removes the connection from the tracker/database.
	
	rec := deleteConnection(t, h, provider, saved.ID)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d for broken cluster delete. body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	if _, err := provider.ConnectionPersister.GetConnection(saved.ID, ""); err == nil {
		t.Error("connection is still readable after a successful delete")
	}
	
	// Tracker should be empty for this connection ID
	if _, ok := tracker.Get(saved.ID); ok {
		t.Error("tracker still contains machine after broken cluster delete")
	}
}
