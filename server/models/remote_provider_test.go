package models

import (
	"net/http/httptest"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
	schemasEnvironment "github.com/meshery/schemas/models/v1beta1/environment"
	schemasWorkspace "github.com/meshery/schemas/models/v1beta1/workspace"
)

func TestRemoteProvider_InvalidUUID_EnvironmentWorkspaceOperations(t *testing.T) {
	provider := &RemoteProvider{
		ProviderProperties: ProviderProperties{
			ProviderName: "test-provider",
			Capabilities: Capabilities{
				{Feature: PersistEnvironments, Endpoint: "/environments"},
				{Feature: PersistWorkspaces, Endpoint: "/workspaces"},
			},
		},
		RemoteProviderURL: "http://example.com",
	}

	req := httptest.NewRequest("GET", "/", nil)
	tests := []struct {
		name string
		call func(string) error
	}{
		{
			name: "Given invalid UUID when getting environment by ID then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.GetEnvironmentByID(req, id, "")
				return err
			},
		},
		{
			name: "Given invalid UUID when deleting environment then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.DeleteEnvironment(req, id)
				return err
			},
		},
		{
			name: "Given invalid UUID when updating environment then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.UpdateEnvironment(req, &schemasEnvironment.EnvironmentPayload{}, id)
				return err
			},
		},
		{
			name: "Given invalid UUID when getting workspace by ID then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.GetWorkspaceByID(req, id, "")
				return err
			},
		},
		{
			name: "Given invalid UUID when deleting workspace then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.DeleteWorkspace(req, id)
				return err
			},
		},
		{
			name: "Given invalid UUID when updating workspace then returns invalid UUID error code",
			call: func(id string) error {
				_, err := provider.UpdateWorkspace(req, &schemasWorkspace.WorkspacePayload{}, id)
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

			if meshkiterrors.GetCode(err) != ErrInvalidUUIDCode {
				t.Fatalf("expected error code %s, got %s", ErrInvalidUUIDCode, meshkiterrors.GetCode(err))
			}
		})
	}
}
