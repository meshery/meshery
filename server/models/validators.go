package models

import (
	"fmt"
	"net/url"
	"slices"
	"time"

	"github.com/gofrs/uuid"
	SMP "github.com/layer5io/service-mesh-performance/spec"
	"github.com/meshery/schemas/models/v1beta1/environment"
	"github.com/meshery/schemas/models/v1beta1/workspace"
)

func parseUUIDWithField(id, field string) (uuid.UUID, error) {
	parsedID, err := uuid.FromString(id)
	if err != nil {
		return uuid.Nil, ErrInvalidUUID(err, field)
	}

	return parsedID, nil
}

func validateOptionalUUIDWithField(id, field string) error {
	if id == "" {
		return nil
	}

	_, err := parseUUIDWithField(id, field)

	return err
}

func validateRequiredUUIDWithField(id, field string) error {
	if id == "" {
		return ErrInvalidUUID(fmt.Errorf("value is empty"), field)
	}

	_, err := parseUUIDWithField(id, field)

	return err
}

func validateRequiredUUIDValue(id uuid.UUID, field string) error {
	if id == uuid.Nil {
		return ErrInvalidUUID(fmt.Errorf("value is empty"), field)
	}

	return nil
}

func validateEnvironmentPayload(payload *environment.EnvironmentPayload) error {
	return validateRequiredUUIDWithField(payload.OrgId, "organization ID")
}

func validateWorkspacePayload(payload *workspace.WorkspacePayload) error {
	return validateRequiredUUIDValue(payload.OrganizationID, "organization ID")
}

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
