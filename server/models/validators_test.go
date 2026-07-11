package models

import (
	"errors"
	"testing"

	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

func TestPerformanceTestConfigValidator(t *testing.T) {
	cases := []struct {
		name      string
		config    *perfprofile.PerformanceTestConfig
		wantError error
	}{
		{
			name: "given relative or invalid endpoints, when validating config, then error is returned",
			config: &perfprofile.PerformanceTestConfig{
				Name:     "invalid endpoints",
				Duration: "30s",
				Clients: []perfprofile.PerformanceTestClient{
					{
						LoadGenerator: FortioLG.Name(),
						Protocol:      "http",
						EndpointUrls:  []string{"not-a-url", "/health", "health"},
					},
				},
			},
			wantError: ErrValidURL,
		},
		{
			name: "given a legacy wrk2 generator, when validating config, then it is tolerated and no error is returned",
			config: &perfprofile.PerformanceTestConfig{
				Name:     "valid endpoints",
				Duration: "30s",
				Clients: []perfprofile.PerformanceTestClient{
					{
						LoadGenerator: FortioLG.Name(),
						Protocol:      "http",
						EndpointUrls:  []string{"https://meshery.io/api/health"},
					},
					{
						// "wrk2" was removed as a load generator; existing
						// profiles carrying it must still validate (they run
						// on fortio). Regression guard for graceful fallback.
						LoadGenerator: "wrk2",
						Protocol:      "tcp",
						EndpointUrls:  []string{"tcp://127.0.0.1:8080"},
					},
				},
			},
			wantError: nil,
		},
		{
			name: "given an unrecognized generator, when validating config, then ErrLoadgenerator is returned",
			config: &perfprofile.PerformanceTestConfig{
				Name:     "bogus generator",
				Duration: "30s",
				Clients: []perfprofile.PerformanceTestClient{
					{
						LoadGenerator: "bogus",
						Protocol:      "http",
						EndpointUrls:  []string{"https://meshery.io/api/health"},
					},
				},
			},
			wantError: ErrLoadgenerator,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := PerformanceTestConfigValidator(tc.config)
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
