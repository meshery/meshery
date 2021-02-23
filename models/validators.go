package models

import (
	"net/url"
	"time"

	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/pkg/errors"
)

// SMPPerformanceTestConfigValidator performs validations on the given PerformanceTestConfig object
func SMPPerformanceTestConfigValidator(perfTest *SMP.PerformanceTestConfig) error {
	if perfTest.Name == "" {
		return errors.Errorf("Error: name field is blank")
	}
	if _, err := time.ParseDuration(perfTest.Duration); err != nil {
		return errors.Wrapf(err, "error parsing test duration, please refer to: https://docs.meshery.io/guides/mesheryctl#performance-management")
	}

	if len(perfTest.Clients) < 1 {
		return errors.Errorf("minimum one test client needs to be specified")
	}

	for _, testClient := range perfTest.Clients {
		if testClient.Protocol.String() == "" {
			return errors.Errorf("specify the Protocol for all clients")
		}
		if !(testClient.LoadGenerator == Wrk2LG.Name() || testClient.LoadGenerator == FortioLG.Name()) {
			return errors.Errorf("specify valid Loadgenerator")
		}
		if len(testClient.EndpointUrls) < 1 {
			return errors.Errorf("minimum one test endpoint needs to be specified")
		}
		for _, URL := range testClient.EndpointUrls {
			if _, err := url.Parse(URL); err != nil {
				return errors.Wrapf(err, "Enter valid URLs")
			}
		}
	}
	return nil
}
