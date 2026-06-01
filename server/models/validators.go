package models

import (
	"net/url"
	"slices"
	"time"

	"github.com/meshery/meshery/server/models/performance"
)

// SMPPerformanceTestConfigValidator performs validations on the given PerformanceTestConfig object
func SMPPerformanceTestConfigValidator(perfTest *performance.PerformanceTestConfig) error {
	if perfTest.Name == "" {
		return ErrField
	}
	if _, err := time.ParseDuration(perfTest.Duration); err != nil {
		return ErrParsingTest
	}

	// Support both Clients and ClientsNested for backward compatibility
	clients := len(perfTest.Clients)
	nestedClients := len(perfTest.ClientsNested)
	
	if clients < 1 && nestedClients < 1 {
		return ErrTestClient
	}

	// fortio is the only supported load generator. The removed "wrk2" is
	// retained purely as a backward-compatibility alias so performance
	// profiles saved before its removal still validate; such profiles
	// transparently run on fortio (see executeLoadTest). An empty value is
	// treated as the default (fortio). Any other value is rejected.
	validGenerators := []string{FortioLG.Name(), "wrk2", ""}
	
	// Validate Clients (new format)
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
		for _, rawURL := range testClient.EndpointUrls {
			parsedURL, err := url.ParseRequestURI(rawURL)
			if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
				return ErrValidURL
			}
		}
	}
	
	// Validate ClientsNested (old format for tests)
	for _, testClient := range perfTest.ClientsNested {
		if testClient.Protocol == 0 {
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
