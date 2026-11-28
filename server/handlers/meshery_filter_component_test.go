package handlers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta3/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// registerWASMFilterComponent registers the WASMFilter component definition
// shipped in the meshery-core model, so the lookup in generateFilterComponent
// is tested against the real definition rather than a hand-written copy.
func registerWASMFilterComponent(t *testing.T, rm *registry.RegistryManager) {
	t.Helper()

	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)
	path := filepath.Join(filepath.Dir(filename), "..", "..", "models", "meshery-core", "0.7.2", "v1.0.0", "components", "WASMFilter.json")
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	var comp component.ComponentDefinition
	require.NoError(t, json.Unmarshal(content, &comp))
	id, err := comp.GenerateID()
	require.NoError(t, err)
	comp.ID = id

	conn := connection.Connection{
		Name:   "meshery",
		Kind:   "meshery",
		Type:   "registry",
		Status: connection.ConnectionStatusConnected,
	}
	_, _, err = rm.RegisterEntity(registry.RegistrantHostToV1beta3(conn), &comp)
	require.NoError(t, err, "registering the WASMFilter component")
}

func newFilterComponentTestHandler(t *testing.T, withWASMFilter bool) *Handler {
	t.Helper()

	rm, _ := newTestRegistryDB(t)
	if withWASMFilter {
		registerWASMFilterComponent(t, rm)
	}
	log, err := logger.New("meshery_filter_component_test", logger.Options{})
	require.NoError(t, err)
	return &Handler{registryManager: rm, log: log}
}

// TestGenerateFilterComponentCarriesConfig guards the filter configuration a
// save persists: it is stored in filterResource under configuration.config, and
// is lost if the WASMFilter component cannot be found in the registry.
func TestGenerateFilterComponentCarriesConfig(t *testing.T) {
	h := newFilterComponentTestHandler(t, true)
	const config = "match: route-a\nratelimit: 100\n"

	filterResource, err := h.generateFilterComponent(config)
	require.NoError(t, err)
	require.NotEmpty(t, filterResource, "the configuration must not be dropped")

	var stored struct {
		Component struct {
			Kind string `json:"kind"`
		} `json:"component"`
		Configuration map[string]interface{} `json:"configuration"`
	}
	require.NoError(t, json.Unmarshal([]byte(filterResource), &stored))
	assert.Equal(t, "WASMFilter", stored.Component.Kind)
	assert.Equal(t, config, stored.Configuration["config"])
}

// TestGenerateFilterComponentWithoutWASMFilter pins what happens when the
// registry has no WASMFilter component: a configuration must not be silently
// discarded, while a save without one has nothing to lose.
func TestGenerateFilterComponentWithoutWASMFilter(t *testing.T) {
	h := newFilterComponentTestHandler(t, false)

	_, err := h.generateFilterComponent("match: route-a\n")
	assert.Error(t, err, "a configuration that cannot be stored must be reported")

	filterResource, err := h.generateFilterComponent("")
	require.NoError(t, err)
	assert.Empty(t, filterResource)
}
