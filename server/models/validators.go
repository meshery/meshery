package models

import (
	"net/url"
	"slices"
	"time"

	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

// PerformanceTestConfigValidator performs validations on the given PerformanceTestConfig object
func PerformanceTestConfigValidator(perfTest *perfprofile.PerformanceTestConfig) error {
	if perfTest == nil {
		return ErrField
	}
	if perfTest.Name == "" {
		return ErrField
	}
	if _, err := time.ParseDuration(perfTest.Duration); err != nil {
		return ErrParsingTest
	}

	if len(perfTest.Clients) < 1 {
		return ErrTestClient
	}

	// fortio is the only supported load generator. The removed "wrk2" is
	// retained purely as a backward-compatibility alias so performance
	// profiles saved before its removal still validate; such profiles
	// transparently run on fortio (see executeLoadTest). An empty value is
	// treated as the default (fortio). Any other value is rejected.
	validGenerators := []string{FortioLG.Name(), "wrk2", ""}
	for _, testClient := range perfTest.Clients {
		if testClient.Protocol == "" {
			return ErrProtocol
		}
		isValidGenerator := slices.Contains(validGenerators, testClient.LoadGenerator)
		if !isValidGenerator {
			return ErrLoadgenerator
		}
		if len(testClient.EndpointUrls) < 1 {
			return ErrTestEndpoint
		}
		for _, rawURL := range testClient.EndpointUrls {
			parsedURL, err := url.ParseRequestURI(rawURL)
			if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
				return ErrValidURL
			}
		}
	}
	return nil
}
