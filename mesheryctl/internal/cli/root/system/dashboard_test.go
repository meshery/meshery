package system

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestKeepConnectionAlive_Success(t *testing.T) {
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	requestReceived := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	keepConnectionAlive(server.URL)

	if !requestReceived {
		t.Fatalf("expected request to be received by server, but it was not")
	}
}

func TestKeepConnectionAlive_Failure(t *testing.T) {
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	// Point to a closed server URL that immediately refuses connection
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := server.URL
	server.Close() // close immediately so connection fails

	// keepConnectionAlive should log error and return without panicking
	keepConnectionAlive(url)
}
