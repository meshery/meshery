package system

import (
	"os"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitutils "github.com/meshery/meshkit/utils"
)

// TestResolveDockerEndpointReturnsDialableHost verifies that the endpoint
// system start waits on is a host and port TcpCheck can dial, not a URL with a
// scheme, which no TCP dial can reach.
func TestResolveDockerEndpointReturnsDialableHost(t *testing.T) {
	tests := []struct {
		name         string
		endpoint     string
		silent       bool
		answer       string
		wantAddress  string
		wantEndpoint string
	}{
		{
			name:         "localhost endpoint",
			endpoint:     "http://localhost:9081",
			silent:       true,
			wantAddress:  "localhost",
			wantEndpoint: "http://localhost:9081",
		},
		{
			name:         "remote endpoint kept when the change to localhost is declined",
			endpoint:     "http://meshery.example.com:9081",
			answer:       "n\n",
			wantAddress:  "meshery.example.com",
			wantEndpoint: "http://meshery.example.com:9081",
		},
		{
			name:         "remote endpoint stored without a scheme",
			endpoint:     "meshery.example.com:9081",
			answer:       "n\n",
			wantAddress:  "meshery.example.com",
			wantEndpoint: "meshery.example.com:9081",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupContextTestEnv(t)
			t.Cleanup(BreakupFunc)
			prevSilent := utils.SilentFlag
			utils.SilentFlag = tt.silent
			t.Cleanup(func() { utils.SilentFlag = prevSilent })

			if tt.answer != "" {
				r, w, err := os.Pipe()
				if err != nil {
					t.Fatal(err)
				}
				_, _ = w.WriteString(tt.answer)
				_ = w.Close()
				stdin := os.Stdin
				os.Stdin = r
				t.Cleanup(func() {
					os.Stdin = stdin
					_ = r.Close()
				})
			}

			currCtx, err := mctlCfg.GetCurrentContext()
			if err != nil {
				t.Fatal(err)
			}
			currCtx.SetEndpoint(tt.endpoint)

			endpoint, err := resolveDockerEndpoint(mctlCfg, currCtx, strings.Split(tt.endpoint, ":"))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if endpoint.Address != tt.wantAddress || endpoint.Port != 9081 {
				t.Errorf("expected endpoint %s:9081, got %s:%d", tt.wantAddress, endpoint.Address, endpoint.Port)
			}
			if !meshkitutils.TcpCheck(&endpoint, &meshkitutils.MockOptions{DesiredEndpoint: tt.wantAddress + ":9081"}) {
				t.Errorf("expected TcpCheck to dial %s:9081, got %s:%d", tt.wantAddress, endpoint.Address, endpoint.Port)
			}
			if got := currCtx.GetEndpoint(); got != tt.wantEndpoint {
				t.Errorf("expected the context endpoint %s, got %s", tt.wantEndpoint, got)
			}
		})
	}
}
