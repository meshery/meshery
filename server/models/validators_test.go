package models

import (
	"testing"

	SMP "github.com/layer5io/service-mesh-performance/spec"
)

func TestSMPPerformanceTestConfigValidatorRejectsInvalidEndpoints(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		endpoint string
		wantErr  bool
	}{
		{
			name:     "rejects bare text",
			endpoint: "not-a-url",
			wantErr:  true,
		},
		{
			name:     "rejects relative path",
			endpoint: "/relative/path",
			wantErr:  true,
		},
		{
			name:     "rejects unsupported scheme",
			endpoint: "ftp://example.com",
			wantErr:  true,
		},
		{
			name:     "accepts valid https endpoint",
			endpoint: "https://example.com/api",
			wantErr:  false,
		},
		{
			name:     "accepts valid http endpoint",
			endpoint: "http://example.com/api",
			wantErr:  false,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := SMPPerformanceTestConfigValidator(validPerformanceTestConfig(tc.endpoint))
			if tc.wantErr {
				if err != ErrValidURL {
					t.Fatalf("expected ErrValidURL for %q, got %v", tc.endpoint, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("expected valid config for %q, got %v", tc.endpoint, err)
			}
		})
	}
}

func validPerformanceTestConfig(endpoint string) *SMP.PerformanceTestConfig {
	return &SMP.PerformanceTestConfig{
		Name:     "test",
		Duration: "1s",
		Clients: []*SMP.PerformanceTestConfig_Client{
			{
				Protocol:      SMP.PerformanceTestConfig_Client_PROTOCOL_HTTP,
				LoadGenerator: Wrk2LG.Name(),
				EndpointUrls:  []string{endpoint},
			},
		},
	}
}
