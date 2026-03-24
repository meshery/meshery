package models

import (
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
	schemasEnvironment "github.com/meshery/schemas/models/v1beta1/environment"
	schemasWorkspace "github.com/meshery/schemas/models/v1beta1/workspace"
)

func TestDefaultLocalProvider_InvalidUUID_EnvironmentWorkspaceOperations(t *testing.T) {
	provider := &DefaultLocalProvider{}

	tests := []struct {
		name string
		call func(string) error
	}{
		{
			name: "Given invalid UUID when getting environment by ID then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.GetEnvironmentByID(nil, id, "")
				return err
			},
		},
		{
			name: "Given invalid UUID when deleting environment then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.DeleteEnvironment(nil, id)
				return err
			},
		},
		{
			name: "Given invalid UUID when updating environment then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.UpdateEnvironment(nil, &schemasEnvironment.EnvironmentPayload{}, id)
				return err
			},
		},
		{
			name: "Given invalid UUID when getting workspace by ID then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.GetWorkspaceByID(nil, id, "")
				return err
			},
		},
		{
			name: "Given invalid UUID when deleting workspace then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.DeleteWorkspace(nil, id)
				return err
			},
		},
		{
			name: "Given invalid UUID when updating workspace then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.UpdateWorkspace(nil, &schemasWorkspace.WorkspacePayload{}, id)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call("not-a-uuid")
			if err == nil {
				t.Fatal("expected error for invalid UUID, got nil")
			}

			if meshkiterrors.GetCode(err) != ErrModelInvalidUUIDCode {
				t.Fatalf("expected error code %s, got %s", ErrModelInvalidUUIDCode, meshkiterrors.GetCode(err))
			}
		})
	}
}
