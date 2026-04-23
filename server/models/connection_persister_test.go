package models

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/environments"
	"github.com/meshery/schemas/models/core"
)

func TestToSchemaEnvironments(t *testing.T) {
	environmentID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate environment ID: %v", err)
	}

	organizationID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate organization ID: %v", err)
	}

	ownerID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate owner ID: %v", err)
	}

	now := time.Now().UTC().Round(time.Second)
	schemaVersion := core.VersionString("v1beta1")
	metadata := core.Map{"source": "test"}
	environmentsFetched := []*environments.EnvironmentData{
		nil,
		{
			ID:             environmentID,
			SchemaVersion:  schemaVersion,
			Name:           "Prod",
			Description:    "Production environment",
			OrganizationID: organizationID,
			Owner:          &ownerID,
			CreatedAt:      now,
			Metadata:       metadata,
			UpdatedAt:      now,
		},
	}

	converted := toSchemaEnvironments(environmentsFetched)
	if len(converted) != 1 {
		t.Fatalf("expected 1 converted environment, got %d", len(converted))
	}

	got := converted[0]
	if got.ID != environmentID {
		t.Fatalf("expected environment ID %s, got %s", environmentID, got.ID)
	}
	if got.SchemaVersion != schemaVersion {
		t.Fatalf("expected schema version %q, got %q", schemaVersion, got.SchemaVersion)
	}
	if got.Name != "Prod" {
		t.Fatalf("expected name %q, got %q", "Prod", got.Name)
	}
	if got.Description != "Production environment" {
		t.Fatalf("expected description %q, got %q", "Production environment", got.Description)
	}
	if got.OrganizationID != organizationID {
		t.Fatalf("expected organization ID %s, got %s", organizationID, got.OrganizationID)
	}
	if got.Owner == nil || *got.Owner != ownerID {
		t.Fatalf("expected owner ID %s, got %v", ownerID, got.Owner)
	}
	if got.CreatedAt != now {
		t.Fatalf("expected createdAt %v, got %v", now, got.CreatedAt)
	}
	if got.UpdatedAt != now {
		t.Fatalf("expected updatedAt %v, got %v", now, got.UpdatedAt)
	}
	if got.Metadata["source"] != metadata["source"] {
		t.Fatalf("expected metadata source %q, got %v", metadata["source"], got.Metadata["source"])
	}
}

func TestToSchemaEnvironmentsReturnsNilForEmptyInput(t *testing.T) {
	if converted := toSchemaEnvironments(nil); converted != nil {
		t.Fatalf("expected nil for nil input, got %#v", converted)
	}

	if converted := toSchemaEnvironments([]*environments.EnvironmentData{nil}); converted != nil {
		t.Fatalf("expected nil when all environments are nil, got %#v", converted)
	}
}
