// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// keepConnectionAlive runs on a ticker for as long as the port forward is up,
// so every call has to return its connection to the pool. If the response body
// is left undrained or unclosed, each tick opens a new connection instead.
func TestKeepConnectionAliveReusesTheConnection(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	var mu sync.Mutex
	opened := 0

	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "meshery")
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			mu.Lock()
			opened++
			mu.Unlock()
		}
	}
	server.Start()
	defer server.Close()

	const ticks = 5
	for i := 0; i < ticks; i++ {
		keepConnectionAlive(server.URL)
	}

	mu.Lock()
	defer mu.Unlock()
	if opened != 1 {
		t.Fatalf("expected %d calls to share a single connection, but %d were opened", ticks, opened)
	}
}

// A failed request must not be reported as a success.
func TestKeepConnectionAliveDoesNotLogSuccessOnError(t *testing.T) {
	b := utils.SetupMeshkitLoggerTesting(t, true)

	// Closed immediately so the address is not listening.
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	unreachableURL := server.URL
	server.Close()

	keepConnectionAlive(unreachableURL)

	output := b.String()
	if !strings.Contains(output, "connection request failed") {
		t.Errorf("expected the failure to be logged, got: %q", output)
	}
	if strings.Contains(output, "connection request success") {
		t.Errorf("a failed request was reported as a success: %q", output)
	}
}
