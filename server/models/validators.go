package models

import (
	"net/url"
	"slices"
	"time"

	SMP "github.com/layer5io/service-mesh-performance/spec"
)

// SMPPerformanceTestConfigValidator performs validations on the given PerformanceTestConfig object
func SMPPerformanceTestConfigValidator(perfTest *SMP.PerformanceTestConfig) error {
	if perfTest.Name == "" {
		return ErrField
	}
	if _, err := time.ParseDuration(perfTest.Duration); err != nil {
		return ErrParsingTest
	}

	if len(perfTest.Clients) < 1 {
		return ErrTestClient
	}

	validGenerators := []string{Wrk2LG.Name(), FortioLG.Name()}
	for _, testClient := range perfTest.Clients {
		if testClient.Protocol.String() == "" {
			return ErrProtocol
		}
		isValidGenerator := slices.Contains(validGenerators, testClient.LoadGenerator)
		if !isValidGenerator {
			return ErrLoadgenerator
		}
		if len(testClient.EndpointUrls) < 1 {
			return ErrTestEndpoint
		}
		for _, URL := range testClient.EndpointUrls {
			if _, err := url.Parse(URL); err != nil {
				return ErrValidURL
			}
		}
	}
	return nil
}
