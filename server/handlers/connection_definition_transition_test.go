package handlers

import (
	"testing"

	"github.com/meshery/meshkit/models/meshmodel/registry"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	hostconn "github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
	"github.com/stretchr/testify/require"
)

// TestConnectionDefinitionTransitionMapRoundTrips guards the `transitionMap`
// persistence path end to end: a connection definition is registered through
// the registry's own RegisterEntity API and then read back via the same
// ConnectionFilter the GetConnectionDefinitions handler uses. The map is
// serialized as JSON on the column (gorm:"type:bytes;serializer:json"); a
// regression to gorm:"-" would silently drop it on write and the UI would
// render "No transitions Available".
func TestConnectionDefinitionTransitionMapRoundTrips(t *testing.T) {
	rm, _ := newTestRegistryManager(t)

	registrant := hostconn.Connection{
		Name:    "test-registrant",
		Kind:    "kubernetes",
		Type:    "platform",
		SubType: "orchestration",
		Status:  hostconn.ConnectionStatusConnected,
	}

	desc := "This will reconnect the cluster and redeploy the operator after maintenance."
	def := connectionv1beta3.ConnectionDefinition{
		Name:           "Kubernetes",
		Kind:           "kubernetes",
		ConnectionType: "platform",
		SubType:        "orchestration",
		Status:         connectionv1beta3.ConnectionStatusDiscovered,
		TransitionMap: map[string][]connectionv1beta3.ConnectionStateTransition{
			"connected": {
				{NextState: connectionv1beta3.Disconnected},
				{NextState: connectionv1beta3.Deleted},
			},
			"disconnected": {
				{NextState: connectionv1beta3.Connected, Description: &desc},
			},
		},
		ModelReference: &model.ModelReference{
			Name:    "kubernetes",
			Version: "v1.25.0",
		},
	}
	id, err := def.GenerateID()
	require.NoError(t, err)
	def.ID = id

	_, _, err = rm.RegisterEntity(registry.RegistrantHostToV1beta3(registrant), &def)
	require.NoError(t, err, "RegisterEntity failed")

	entities, count, _, err := rm.GetEntities(&regv1beta1.ConnectionFilter{Kind: "kubernetes"})
	require.NoError(t, err)
	require.EqualValues(t, 1, count, "expected exactly one connection definition")

	got, ok := entities[0].(*connectionv1beta3.ConnectionDefinition)
	require.True(t, ok, "entity is not a *ConnectionDefinition")

	require.NotEmpty(t, got.TransitionMap, "transitionMap was dropped on the DB round-trip")
	require.Len(t, got.TransitionMap["connected"], 2)
	require.Equal(t,
		connectionv1beta3.Connected,
		got.TransitionMap["disconnected"][0].NextState,
	)
	require.NotNil(t, got.TransitionMap["disconnected"][0].Description)
	require.Equal(t, desc, *got.TransitionMap["disconnected"][0].Description)

	// The model-name filter resolves through the serialized model_reference
	// column (v1beta3 ConnectionDefinition has no model_id FK column). Confirm a
	// definition is findable by its ModelReference.Name and that an unknown model
	// matches nothing.
	byModel, mCount, _, err := rm.GetEntities(&regv1beta1.ConnectionFilter{ModelName: "kubernetes"})
	require.NoError(t, err)
	require.EqualValues(t, 1, mCount, "expected the definition to be found by model name")
	require.Len(t, byModel, 1)

	none, nCount, _, err := rm.GetEntities(&regv1beta1.ConnectionFilter{ModelName: "no-such-model"})
	require.NoError(t, err)
	require.EqualValues(t, 0, nCount, "unknown model name should match no connection definitions")
	require.Empty(t, none)
}
