package models

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/meshmodel/core/v1beta1"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistrantImportSummariesDeduplicatesByConnectionID(t *testing.T) {
	artifactHubID := uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))
	githubID := uuid.Must(uuid.FromString("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"))

	hosts := []v1beta1.MeshModelHostsWithEntitySummary{
		{
			Connection: connection.Connection{ID: artifactHubID, Kind: "artifacthub", Name: "Artifact Hub", Type: "registry"},
			Summary:    v1beta1.EntitySummary{Models: 150, Components: 1326, Relationships: 26},
		},
		{
			Connection: connection.Connection{ID: artifactHubID, Kind: "artifacthub", Name: "Artifact Hub", Type: "registry"},
			Summary:    v1beta1.EntitySummary{Models: 2, Components: 10},
		},
		{
			Connection: connection.Connection{ID: githubID, Kind: "github", Name: "github"},
			Summary:    v1beta1.EntitySummary{Models: 129, Components: 1441, Relationships: 559},
		},
	}

	got := registrantImportSummaries(hosts)
	require.Len(t, got, 2)

	artifactHubSummary := got[0].summary
	assert.Equal(t, v1beta1.EntitySummary{Models: 152, Components: 1336, Relationships: 26}, artifactHubSummary)

	message := registrySuccessMessage(got[0].displayName, artifactHubSummary)
	want := "For registrant artifacthub imported 152 models, 1336 components, 26 relationships."
	assert.Equal(t, want, message)
}

func TestRegistrantImportSummariesKeepsDistinctRegistrantNamesWithSameKind(t *testing.T) {
	firstID := uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))
	secondID := uuid.Must(uuid.FromString("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"))

	hosts := []v1beta1.MeshModelHostsWithEntitySummary{
		{
			Connection: connection.Connection{
				ID:   firstID,
				Kind: "artifacthub",
				Name: "Artifact Hub Catalog",
				Type: "registry",
			},
			Summary: v1beta1.EntitySummary{Models: 150},
		},
		{
			Connection: connection.Connection{
				ID:   secondID,
				Kind: "artifacthub",
				Name: "Artifact Hub Extension",
				Type: "registry",
			},
			Summary: v1beta1.EntitySummary{Models: 2},
		},
	}

	got := registrantImportSummaries(hosts)
	require.Len(t, got, 2)
	assert.Equal(t, v1beta1.EntitySummary{Models: 150}, got[0].summary)
	assert.Equal(t, v1beta1.EntitySummary{Models: 2}, got[1].summary)
	assert.NotEqual(t, got[0].displayName, got[1].displayName)
}

func TestRegistrantImportSummariesDeduplicatesSameNameAndKindWithDifferentGeneratedIDs(t *testing.T) {
	firstID := uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))
	secondID := uuid.Must(uuid.FromString("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"))

	hosts := []v1beta1.MeshModelHostsWithEntitySummary{
		{
			Connection: connection.Connection{ID: firstID, Kind: "artifacthub", Name: "Artifact Hub", Type: "registry"},
			Summary:    v1beta1.EntitySummary{Models: 150},
		},
		{
			Connection: connection.Connection{ID: secondID, Kind: "artifacthub", Name: "Artifact Hub", Type: "registry"},
			Summary:    v1beta1.EntitySummary{Models: 2},
		},
	}

	got := registrantImportSummaries(hosts)
	require.Len(t, got, 1)
	assert.Equal(t, v1beta1.EntitySummary{Models: 152}, got[0].summary)
}
