package models

import (
	"errors"
	"testing"

	SMP "github.com/layer5io/service-mesh-performance/spec"
)

func TestSMPPerformanceTestConfigValidator(t *testing.T) {
	cases := []struct {
		name      string
		config    *SMP.PerformanceTestConfig
		wantError error
	}{
		{
			name: "given relative or invalid endpoints, when validating config, then error is returned",
			config: &SMP.PerformanceTestConfig{
				Name:     "invalid endpoints",
				Duration: "30s",
				Clients: []*SMP.PerformanceTestConfig_Client{
					{
						LoadGenerator: FortioLG.Name(),
						Protocol:      SMP.PerformanceTestConfig_Client_PROTOCOL_HTTP,
						EndpointUrls:  []string{"not-a-url", "/health", "health"},
					},
				},
			},
			wantError: ErrValidURL,
		},
		{
			name: "given valid absolute endpoints, when validating config, then no error is returned",
			config: &SMP.PerformanceTestConfig{
				Name:     "valid endpoints",
				Duration: "30s",
				Clients: []*SMP.PerformanceTestConfig_Client{
					{
						LoadGenerator: FortioLG.Name(),
						Protocol:      SMP.PerformanceTestConfig_Client_PROTOCOL_HTTP,
						EndpointUrls:  []string{"https://meshery.io/api/health"},
					},
					{
						LoadGenerator: Wrk2LG.Name(),
						Protocol:      SMP.PerformanceTestConfig_Client_PROTOCOL_TCP,
						EndpointUrls:  []string{"tcp://127.0.0.1:8080"},
					},
				},
			},
			wantError: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := SMPPerformanceTestConfigValidator(tc.config)
			if tc.wantError != nil {
				if !errors.Is(err, tc.wantError) {
					t.Errorf("expected error %v, got %v", tc.wantError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
