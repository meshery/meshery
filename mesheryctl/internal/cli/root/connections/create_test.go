package connections

import (
	"encoding/json"
	"testing"
)

// TestConnectionIDForContext exercises the connection-id resolution that backs
// the `connection create` -> `connection view`/`connection delete` handoff. The
// id must be resolved by an EXACT context-name match; a multi-context kubeconfig
// registers every context in one response, so an arbitrary fallback could return
// a different context's id than the one the user selected (the "arbitrary
// conn-id" review fix). When the server does not echo the requested context the
// function must return "" rather than guess.
func TestConnectionIDForContext(t *testing.T) {
	const (
		minikubeID = "aaaaaaaa-1111-1111-1111-111111111111"
		otherID    = "bbbbbbbb-2222-2222-2222-222222222222"
	)

	tests := []struct {
		name     string
		response saveK8sContextResponse
		cname    string
		want     string
	}{
		{
			name: "exact match in registeredContexts returns its id (positive path)",
			response: saveK8sContextResponse{
				RegisteredContexts: []registeredK8sContext{
					{Name: "minikube", ConnectionID: minikubeID},
				},
			},
			cname: "minikube",
			want:  minikubeID,
		},
		{
			name: "exact match in connectedContexts returns its id",
			response: saveK8sContextResponse{
				ConnectedContexts: []registeredK8sContext{
					{Name: "minikube", ConnectionID: minikubeID},
				},
			},
			cname: "minikube",
			want:  minikubeID,
		},
		{
			name: "multi-context response returns the requested context, not the first",
			response: saveK8sContextResponse{
				RegisteredContexts: []registeredK8sContext{
					{Name: "prod-cluster", ConnectionID: otherID},
					{Name: "minikube", ConnectionID: minikubeID},
				},
			},
			cname: "minikube",
			want:  minikubeID,
		},
		{
			name: "requested context not echoed returns empty, never an arbitrary id",
			response: saveK8sContextResponse{
				RegisteredContexts: []registeredK8sContext{
					{Name: "prod-cluster", ConnectionID: otherID},
					{Name: "staging", ConnectionID: "cccccccc-3333-3333-3333-333333333333"},
				},
			},
			cname: "minikube",
			want:  "",
		},
		{
			name: "name matches but id empty returns empty",
			response: saveK8sContextResponse{
				RegisteredContexts: []registeredK8sContext{
					{Name: "minikube", ConnectionID: ""},
				},
			},
			cname: "minikube",
			want:  "",
		},
		{
			name:     "empty response returns empty",
			response: saveK8sContextResponse{},
			cname:    "minikube",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := connectionIDForContext(tt.response, tt.cname); got != tt.want {
				t.Fatalf("connectionIDForContext(%q) = %q, want %q", tt.cname, got, tt.want)
			}
		})
	}
}

// TestSaveK8sContextResponseWireContract pins the camelCase wire tags the CLI
// parses off the server's /api/system/kubernetes response. If a tag drifts to
// snake_case the id resolution silently returns "" and the create->view/delete
// handoff regresses without a compile error, so assert the contract directly.
func TestSaveK8sContextResponseWireContract(t *testing.T) {
	const body = `{
		"connectedContexts": [
			{"name": "minikube", "connectionId": "aaaaaaaa-1111-1111-1111-111111111111"}
		],
		"registeredContexts": [
			{"name": "prod-cluster", "connectionId": "bbbbbbbb-2222-2222-2222-222222222222"}
		]
	}`

	var response saveK8sContextResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatalf("failed to unmarshal server response: %v", err)
	}

	if got, want := connectionIDForContext(response, "minikube"), "aaaaaaaa-1111-1111-1111-111111111111"; got != want {
		t.Fatalf("connectionIDForContext(minikube) = %q, want %q", got, want)
	}
	if got, want := connectionIDForContext(response, "prod-cluster"), "bbbbbbbb-2222-2222-2222-222222222222"; got != want {
		t.Fatalf("connectionIDForContext(prod-cluster) = %q, want %q", got, want)
	}
}
