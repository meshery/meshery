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
		call func() error
	}{
		{
			name: "Given invalid organization UUID when getting environments then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetEnvironments("", "", "", "", "", "", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid UUID when getting environment by ID then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetEnvironmentByID(nil, "not-a-uuid", "")
				return err
			},
		},
		{
			name: "Given invalid organization UUID when getting environment by ID then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetEnvironmentByID(nil, "123e4567-e89b-12d3-a456-426614174000", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid UUID when deleting environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.DeleteEnvironment(nil, "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid organization UUID when saving environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.SaveEnvironment(nil, &schemasEnvironment.EnvironmentPayload{OrgId: "not-a-uuid"}, "", false)
				return err
			},
		},
		{
			name: "Given invalid UUID when updating environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.UpdateEnvironment(nil, &schemasEnvironment.EnvironmentPayload{}, "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid organization UUID when updating environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.UpdateEnvironment(nil, &schemasEnvironment.EnvironmentPayload{OrgId: "not-a-uuid"}, "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid environment UUID when adding connection to environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddConnectionToEnvironment(nil, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid connection UUID when adding connection to environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddConnectionToEnvironment(nil, "123e4567-e89b-12d3-a456-426614174000", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid environment UUID when removing connection from environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.RemoveConnectionFromEnvironment(nil, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid environment UUID when getting connections of environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetConnectionsOfEnvironment(nil, "not-a-uuid", "", "", "", "", "")
				return err
			},
		},
		{
			name: "Given invalid organization UUID when getting workspaces then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetWorkspaces("", "", "", "", "", "", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid UUID when getting workspace by ID then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetWorkspaceByID(nil, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given missing organization UUID when getting workspace by ID then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetWorkspaceByID(nil, "123e4567-e89b-12d3-a456-426614174000", "")
				return err
			},
		},
		{
			name: "Given invalid organization UUID when getting workspace by ID then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetWorkspaceByID(nil, "123e4567-e89b-12d3-a456-426614174000", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid UUID when deleting workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.DeleteWorkspace(nil, "not-a-uuid")
				return err
			},
		},
		{
			name: "Given missing organization UUID when saving workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.SaveWorkspace(nil, &schemasWorkspace.WorkspacePayload{}, "", false)
				return err
			},
		},
		{
			name: "Given invalid UUID when updating workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.UpdateWorkspace(nil, &schemasWorkspace.WorkspacePayload{}, "not-a-uuid")
				return err
			},
		},
		{
			name: "Given missing organization UUID when updating workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.UpdateWorkspace(nil, &schemasWorkspace.WorkspacePayload{}, "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid workspace UUID when getting environments of workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetEnvironmentsOfWorkspace(nil, "not-a-uuid", "", "", "", "", "")
				return err
			},
		},
		{
			name: "Given invalid workspace UUID when adding environment to workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddEnvironmentToWorkspace(nil, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid environment UUID when adding environment to workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddEnvironmentToWorkspace(nil, "123e4567-e89b-12d3-a456-426614174000", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid workspace UUID when removing environment from workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.RemoveEnvironmentFromWorkspace(nil, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid workspace UUID when adding design to workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddDesignToWorkspace(nil, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid design UUID when adding design to workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddDesignToWorkspace(nil, "123e4567-e89b-12d3-a456-426614174000", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid workspace UUID when getting designs of workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetDesignsOfWorkspace(nil, "not-a-uuid", "", "", "", "", "", nil)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			if err == nil {
				t.Fatal("expected error for invalid UUID, got nil")
			}

			if meshkiterrors.GetCode(err) != ErrModelInvalidUUIDCode {
				t.Fatalf("expected error code %s, got %s", ErrModelInvalidUUIDCode, meshkiterrors.GetCode(err))
			}
		})
	}
}
