package models

import (
	"reflect"
	"testing"
)

// The shapes below are the ones read-only production inspection found on
// layer5io/meshery-cloud issue #5918: ~36.9k Kubernetes credentials shaped
// {auth, cluster}, ~154 shaped {secret: "<string>"}, and 8 double-nested
// {credentialName, secret: {...}} rows written by Meshery UI's credential form.
// The canonical shape is what meshery/schemas declares and what Cloud is moving
// to. All four have to keep resolving.

func TestCredentialPayload(t *testing.T) {
	tests := []struct {
		name   string
		secret map[string]interface{}
		want   map[string]interface{}
	}{
		{
			name:   "canonical prometheus payload is the secret map itself",
			secret: map[string]interface{}{"prometheusURL": "https://prom.example"},
			want:   map[string]interface{}{"prometheusURL": "https://prom.example"},
		},
		{
			name: "canonical kubernetes payload is the secret map itself",
			secret: map[string]interface{}{
				"clusterName":      "cluster-a",
				"clusterServerURL": "https://k8s.example",
				"auth":             map[string]interface{}{"clusterToken": "tok"},
			},
			want: map[string]interface{}{
				"clusterName":      "cluster-a",
				"clusterServerURL": "https://k8s.example",
				"auth":             map[string]interface{}{"clusterToken": "tok"},
			},
		},
		{
			name: "stored kubernetes shape is the secret map itself",
			secret: map[string]interface{}{
				"auth":    map[string]interface{}{"clusterToken": "tok"},
				"cluster": map[string]interface{}{"server": "https://k8s.example"},
			},
			want: map[string]interface{}{
				"auth":    map[string]interface{}{"clusterToken": "tok"},
				"cluster": map[string]interface{}{"server": "https://k8s.example"},
			},
		},
		{
			name: "legacy double-nested payload is unwrapped",
			secret: map[string]interface{}{
				"credentialName": "kube-cred",
				"secret": map[string]interface{}{
					"auth":    map[string]interface{}{"clusterToken": "tok"},
					"cluster": map[string]interface{}{"server": "https://k8s.example"},
				},
			},
			want: map[string]interface{}{
				"auth":    map[string]interface{}{"clusterToken": "tok"},
				"cluster": map[string]interface{}{"server": "https://k8s.example"},
			},
		},
		{
			name: "legacy double-nested payload keyed by name is unwrapped",
			secret: map[string]interface{}{
				"name":   "kube-cred",
				"secret": map[string]interface{}{"auth": map[string]interface{}{"clusterToken": "tok"}},
			},
			want: map[string]interface{}{"auth": map[string]interface{}{"clusterToken": "tok"}},
		},
		{
			name:   "legacy string shape carries no object payload",
			secret: map[string]interface{}{"secret": "super-secret"},
			want:   nil,
		},
		{
			name: "a payload carrying its own secret field is not unwrapped",
			secret: map[string]interface{}{
				"grafanaURL": "https://grafana.example",
				"secret":     map[string]interface{}{"nested": "value"},
			},
			want: map[string]interface{}{
				"grafanaURL": "https://grafana.example",
				"secret":     map[string]interface{}{"nested": "value"},
			},
		},
		{
			name:   "nil secret",
			secret: nil,
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CredentialPayload(tt.secret)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("CredentialPayload() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestCredentialAuthSecret(t *testing.T) {
	tests := []struct {
		name   string
		secret map[string]interface{}
		want   string
	}{
		{
			name:   "legacy string shape",
			secret: map[string]interface{}{"secret": "super-secret"},
			want:   "super-secret",
		},
		{
			name:   "legacy string shape alongside a credential name",
			secret: map[string]interface{}{"credentialName": "prom-token", "secret": "super-secret"},
			want:   "super-secret",
		},
		{
			name: "canonical grafana payload",
			secret: map[string]interface{}{
				"grafanaURL":    "https://grafana.example",
				"grafanaAPIKey": "canonical-key",
			},
			want: "canonical-key",
		},
		{
			name: "legacy double-nested grafana payload",
			secret: map[string]interface{}{
				"credentialName": "grafana-cred",
				"secret": map[string]interface{}{
					"grafanaURL":    "https://grafana.example",
					"grafanaAPIKey": "nested-key",
				},
			},
			want: "nested-key",
		},
		{
			name:   "canonical prometheus payload is anonymous",
			secret: map[string]interface{}{"prometheusURL": "https://prom.example"},
			want:   "",
		},
		{
			name: "kubernetes payload carries no string auth material",
			secret: map[string]interface{}{
				"auth":    map[string]interface{}{"clusterToken": "tok"},
				"cluster": map[string]interface{}{"server": "https://k8s.example"},
			},
			want: "",
		},
		{
			name: "canonical field wins over a sibling legacy secret string",
			secret: map[string]interface{}{
				"grafanaAPIKey": "canonical-key",
				"secret":        "legacy-key",
			},
			want: "canonical-key",
		},
		{
			name:   "nil secret",
			secret: nil,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CredentialAuthSecret(tt.secret); got != tt.want {
				t.Fatalf("CredentialAuthSecret() = %q, want %q", got, tt.want)
			}
		})
	}
}
