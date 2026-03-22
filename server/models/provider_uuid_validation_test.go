package models

import (
	"net/http/httptest"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
	schemasEnvironment "github.com/meshery/schemas/models/v1beta1/environment"
	schemasWorkspace "github.com/meshery/schemas/models/v1beta1/workspace"
)

func TestDefaultLocalProvider_InvalidUUIDForEnvironmentWorkspaceOperations(t *testing.T) {
	provider := &DefaultLocalProvider{}

	tests := []struct {
		name string
		call func(string) error
	}{
		{
			name: "get environment by id",
			call: func(id string) error {
				_, err := provider.GetEnvironmentByID(nil, id, "")
				return err
			},
		},
		{
			name: "delete environment",
			call: func(id string) error {
				_, err := provider.DeleteEnvironment(nil, id)
				return err
			},
		},
		{
			name: "update environment",
			call: func(id string) error {
				_, err := provider.UpdateEnvironment(nil, &schemasEnvironment.EnvironmentPayload{}, id)
				return err
			},
		},
		{
			name: "get workspace by id",
			call: func(id string) error {
				_, err := provider.GetWorkspaceByID(nil, id, "")
				return err
			},
		},
		{
			name: "delete workspace",
			call: func(id string) error {
				_, err := provider.DeleteWorkspace(nil, id)
				return err
			},
		},
		{
			name: "update workspace",
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

			if meshkiterrors.GetCode(err) != ErrProviderInvalidUUIDCode {
				t.Fatalf("expected error code %s, got %s", ErrProviderInvalidUUIDCode, meshkiterrors.GetCode(err))
			}
		})
	}
}

func TestRemoteProvider_InvalidUUIDForEnvironmentWorkspaceOperations(t *testing.T) {
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
			name: "get environment by id",
			call: func(id string) error {
				_, err := provider.GetEnvironmentByID(req, id, "")
				return err
			},
		},
		{
			name: "delete environment",
			call: func(id string) error {
				_, err := provider.DeleteEnvironment(req, id)
				return err
			},
		},
		{
			name: "update environment",
			call: func(id string) error {
				_, err := provider.UpdateEnvironment(req, &schemasEnvironment.EnvironmentPayload{}, id)
				return err
			},
		},
		{
			name: "get workspace by id",
			call: func(id string) error {
				_, err := provider.GetWorkspaceByID(req, id, "")
				return err
			},
		},
		{
			name: "delete workspace",
			call: func(id string) error {
				_, err := provider.DeleteWorkspace(req, id)
				return err
			},
		},
		{
			name: "update workspace",
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

			if meshkiterrors.GetCode(err) != ErrProviderInvalidUUIDCode {
				t.Fatalf("expected error code %s, got %s", ErrProviderInvalidUUIDCode, meshkiterrors.GetCode(err))
			}
		})
	}
}
