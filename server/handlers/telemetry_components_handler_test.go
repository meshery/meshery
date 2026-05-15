package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTelemetryComponentsHandler_MissingClusterIds(t *testing.T) {
	h := newSSETestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/telemetry/components", nil)
	w := httptest.NewRecorder()

	h.GetTelemetryComponentsHandler(w, req, nil, nil, nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("missing clusterIds should 400, got %d", w.Code)
	}
}

func TestGetTelemetryComponentsHandler_AcceptsRepeatedClusterIds(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/telemetry/components?clusterIds=a&clusterIds=b", nil)
	if got, want := len(req.URL.Query()["clusterIds"]), 2; got != want {
		t.Fatalf("query parse: got %d clusterIds, want %d", got, want)
	}
}
