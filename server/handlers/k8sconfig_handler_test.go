package handlers

import (
	"errors"
	"testing"
)

// TestK8sEventMetadataHasError guards the decision that raises a Kubernetes
// connection receipt event to Error severity. A receipt whose per-context
// metadata records any failure must be reported as Error so it persists in the
// notification center and stays retrievable under the Error severity filter
// (issue #20725); a receipt describing only successful connections must remain
// Informational.
func TestK8sEventMetadataHasError(t *testing.T) {
	tests := []struct {
		name          string
		eventMetadata map[string]interface{}
		want          bool
	}{
		{
			name:          "nil metadata reports no error",
			eventMetadata: nil,
			want:          false,
		},
		{
			name:          "empty metadata reports no error",
			eventMetadata: map[string]interface{}{},
			want:          false,
		},
		{
			name: "present-but-nil error value reports no error",
			eventMetadata: map[string]interface{}{
				"prod": map[string]interface{}{
					"error": nil,
				},
			},
			want: false,
		},
		{
			name: "only successful contexts report no error",
			eventMetadata: map[string]interface{}{
				"prod": map[string]interface{}{
					"description": "Connection registered with kubernetes context \"prod\".",
				},
				"staging": map[string]interface{}{
					"description": "Connection already exists with Kubernetes context \"staging\".",
				},
			},
			want: false,
		},
		{
			name: "a single failed context reports an error",
			eventMetadata: map[string]interface{}{
				"unreachable": map[string]interface{}{
					"description": "Unable to establish connection with context \"unreachable\".",
					"error":       errors.New("api server unreachable"),
				},
			},
			want: true,
		},
		{
			name: "a failure mixed with successes reports an error",
			eventMetadata: map[string]interface{}{
				"prod": map[string]interface{}{
					"description": "Connection registered with kubernetes context \"prod\".",
				},
				"unreachable": map[string]interface{}{
					"error": errors.New("api server unreachable"),
				},
			},
			want: true,
		},
		{
			name: "non-map entries are ignored",
			eventMetadata: map[string]interface{}{
				"weird": "not a metadata map",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := k8sEventMetadataHasError(tt.eventMetadata); got != tt.want {
				t.Errorf("k8sEventMetadataHasError() = %v, want %v", got, tt.want)
			}
		})
	}
}
