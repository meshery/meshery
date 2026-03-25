package models

import (
	"net/http/httptest"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
	schemasEnvironment "github.com/meshery/schemas/models/v1beta1/environment"
	schemasWorkspace "github.com/meshery/schemas/models/v1beta1/workspace"
	"github.com/stretchr/testify/assert"
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
				_, err := provider.GetEnvironmentByID(req, "not-a-uuid", "")
				return err
			},
		},
		{
			name: "Given invalid organization UUID when getting environment by ID then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetEnvironmentByID(req, "123e4567-e89b-12d3-a456-426614174000", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid UUID when deleting environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.DeleteEnvironment(req, "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid organization UUID when saving environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.SaveEnvironment(req, &schemasEnvironment.EnvironmentPayload{OrgId: "not-a-uuid"}, "", true)
				return err
			},
		},
		{
			name: "Given invalid UUID when updating environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.UpdateEnvironment(req, &schemasEnvironment.EnvironmentPayload{}, "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid organization UUID when updating environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.UpdateEnvironment(req, &schemasEnvironment.EnvironmentPayload{OrgId: "not-a-uuid"}, "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid environment UUID when adding connection to environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddConnectionToEnvironment(req, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid connection UUID when adding connection to environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddConnectionToEnvironment(req, "123e4567-e89b-12d3-a456-426614174000", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid environment UUID when removing connection from environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.RemoveConnectionFromEnvironment(req, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid environment UUID when getting connections of environment then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetConnectionsOfEnvironment(req, "not-a-uuid", "", "", "", "", "")
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
				_, err := provider.GetWorkspaceByID(req, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given missing organization UUID when getting workspace by ID then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetWorkspaceByID(req, "123e4567-e89b-12d3-a456-426614174000", "")
				return err
			},
		},
		{
			name: "Given invalid organization UUID when getting workspace by ID then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetWorkspaceByID(req, "123e4567-e89b-12d3-a456-426614174000", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid UUID when deleting workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.DeleteWorkspace(req, "not-a-uuid")
				return err
			},
		},
		{
			name: "Given missing organization UUID when saving workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.SaveWorkspace(req, &schemasWorkspace.WorkspacePayload{}, "", true)
				return err
			},
		},
		{
			name: "Given invalid UUID when updating workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.UpdateWorkspace(req, &schemasWorkspace.WorkspacePayload{}, "not-a-uuid")
				return err
			},
		},
		{
			name: "Given missing organization UUID when updating workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.UpdateWorkspace(req, &schemasWorkspace.WorkspacePayload{}, "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid workspace UUID when adding environment to workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddEnvironmentToWorkspace(req, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid environment UUID when adding environment to workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddEnvironmentToWorkspace(req, "123e4567-e89b-12d3-a456-426614174000", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid workspace UUID when removing environment from workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.RemoveEnvironmentFromWorkspace(req, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid workspace UUID when getting environments of workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetEnvironmentsOfWorkspace(req, "not-a-uuid", "", "", "", "", "")
				return err
			},
		},
		{
			name: "Given invalid workspace UUID when adding design to workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddDesignToWorkspace(req, "not-a-uuid", "123e4567-e89b-12d3-a456-426614174000")
				return err
			},
		},
		{
			name: "Given invalid design UUID when adding design to workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.AddDesignToWorkspace(req, "123e4567-e89b-12d3-a456-426614174000", "not-a-uuid")
				return err
			},
		},
		{
			name: "Given invalid workspace UUID when getting designs of workspace then returns invalid UUID error code",
			call: func() error {
				_, err := provider.GetDesignsOfWorkspace(req, "not-a-uuid", "", "", "", "", "", nil)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call()
			assert.Error(t, err)
			assert.Equal(t, ErrModelInvalidUUIDCode, meshkiterrors.GetCode(err))
		})
	}
}
