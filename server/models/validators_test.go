package models

import (
	"testing"

	SMP "github.com/layer5io/service-mesh-performance/spec"
)

func TestSMPPerformanceTestConfigValidatorRejectsRelativeEndpoints(t *testing.T) {
	perfTest := &SMP.PerformanceTestConfig{
		Name:     "invalid-endpoint-test",
		Duration: "30s",
		Clients: []*SMP.PerformanceTestConfig_Client{
			{
				LoadGenerator: FortioLG.Name(),
				Protocol:      SMP.PerformanceTestConfig_Client_PROTOCOL_HTTP,
				EndpointUrls:  []string{"not-a-url"},
			},
		},
	}

	err := SMPPerformanceTestConfigValidator(perfTest)
	if err != ErrValidURL {
		t.Fatalf("expected ErrValidURL, got %v", err)
	}
}

func TestSMPPerformanceTestConfigValidatorAcceptsAbsoluteEndpoints(t *testing.T) {
	perfTest := &SMP.PerformanceTestConfig{
		Name:     "valid-endpoint-test",
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
	}

	if err := SMPPerformanceTestConfigValidator(perfTest); err != nil {
		t.Fatalf("expected valid absolute endpoints to pass validation, got %v", err)
	}
}
