package models

import (
	"errors"
	"testing"

	SMP "github.com/service-mesh-performance/service-mesh-performance/spec"
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
			name: "given a legacy wrk2 generator, when validating config, then it is tolerated and no error is returned",
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
						// "wrk2" was removed as a load generator; existing
						// profiles carrying it must still validate (they run
						// on fortio). Regression guard for graceful fallback.
						LoadGenerator: "wrk2",
						Protocol:      SMP.PerformanceTestConfig_Client_PROTOCOL_TCP,
						EndpointUrls:  []string{"tcp://127.0.0.1:8080"},
					},
				},
			},
			wantError: nil,
		},
		{
			name: "given an unrecognized generator, when validating config, then ErrLoadgenerator is returned",
			config: &SMP.PerformanceTestConfig{
				Name:     "bogus generator",
				Duration: "30s",
				Clients: []*SMP.PerformanceTestConfig_Client{
					{
						LoadGenerator: "bogus",
						Protocol:      SMP.PerformanceTestConfig_Client_PROTOCOL_HTTP,
						EndpointUrls:  []string{"https://meshery.io/api/health"},
					},
				},
			},
			wantError: ErrLoadgenerator,
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
