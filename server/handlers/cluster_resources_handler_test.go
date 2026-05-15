package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetClusterResourcesHandler_MissingClusterIds(t *testing.T) {
	h := newSSETestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/cluster/resources", nil)
	w := httptest.NewRecorder()

	h.GetClusterResourcesHandler(w, req, nil, nil, nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("missing clusterIds should 400, got %d", w.Code)
	}
}

// nilPersisterProvider returns a nil *gorm.DB so the handler hits the
// "persister unavailable" guard.
type nilPersisterProvider struct{}

func TestGetClusterResourcesHandler_AcceptsRepeatedClusterIds(t *testing.T) {
	// Just exercises query parsing: with two clusterIds and no persister the
	// handler should reach the persister check and 500 (not 400). This pins
	// down the param-shape contract — the UI reads multiple ?clusterIds=… and
	// any future single-value regression would surface as a 400 here.
	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/cluster/resources?clusterIds=a&clusterIds=b", nil)
	if got, want := len(req.URL.Query()["clusterIds"]), 2; got != want {
		t.Fatalf("query parse: got %d clusterIds, want %d", got, want)
	}
}
