package handlers

import (
	"testing"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/stretchr/testify/require"
)

func setupTestRegistry(t *testing.T) *meshmodel.RegistryManager {
	log, _ := logger.New("test", logger.Options{})
	dbHandler, err := database.New(database.Options{
		Engine: database.SQLITE,
		Filename: "file::memory:?cache=shared",
		Logger: log,
	})
	require.NoError(t, err)

	regManager, err := meshmodel.NewRegistryManager(&dbHandler)
	require.NoError(t, err)

	return regManager
}

func TestMapModelRelatedData_EmptyPatternFile(t *testing.T) {
	// Passing an empty but non-nil pattern file should be handled safely
	patternFile := &pattern.PatternFile{}
	regManager := setupTestRegistry(t)

	err := mapModelRelatedData(regManager, patternFile)
	// We expect an error because the component is unrecognized or registry is empty, but it should not panic
	_ = err
}

func TestMapModelRelatedData_UnrecognizedComponent_NonNilReg(t *testing.T) {
	patternFile := &pattern.PatternFile{
		Components: []*component.ComponentDefinition{
			{
				Component: component.Component{
					Kind: "UnknownKind",
				},
				Model: &model.ModelDefinition{
					Name: "UnknownModel",
					Model: model.Model{
						Version: "v1.0.0",
					},
				},
			},
		},
	}

	regManager := setupTestRegistry(t)

	err := mapModelRelatedData(regManager, patternFile)
	// We expect an error because the component is unrecognized or registry is empty, but it should not panic
	_ = err
}

func TestMapModelRelatedData_NilMetadataAssignment(t *testing.T) {
	regManager := setupTestRegistry(t)

	// Register a component definition
	def := component.ComponentDefinition{
		Component: component.Component{
			Kind: "TestKind",
			Version: "v1.0.0",
		},
		Model: &model.ModelDefinition{
			Name: "meshery-core",
			Model: model.Model{
				Version: "v1.0.0",
			},
		},
	}
	conn := connection.Connection{}
	_, _, err := regManager.RegisterEntity(conn, &def)
	require.NoError(t, err)

	patternFile := &pattern.PatternFile{
		Components: []*component.ComponentDefinition{
			nil, // Test handling of nil components in the slice
			{
				Component: component.Component{
					Kind: "TestKind",
					Version: "v1.0.0",
				},
				Model: &model.ModelDefinition{
					Name: "meshery-core",
					Model: model.Model{
						Version: "v1.0.0",
					},
				},
			},
		},
	}

	// Should not panic even with nil components in the slice
	err = mapModelRelatedData(regManager, patternFile)
	require.NoError(t, err)
}

func TestMapModelRelatedData_PartialModelDefinition(t *testing.T) {
	regManager := setupTestRegistry(t)

	// Register a component with partial model definitions
	def := component.ComponentDefinition{
		Component: component.Component{
			Kind: "PartialKind",
			Version: "v1.0.0",
		},
		Model: &model.ModelDefinition{
			Name: "meshery-core",
			Model: model.Model{
				Version: "v1.0.0",
			},
			Registrant: connection.Connection{},
		},
	}
	conn := connection.Connection{}
	_, _, err := regManager.RegisterEntity(conn, &def)
	require.NoError(t, err)

	patternFile := &pattern.PatternFile{
		Components: []*component.ComponentDefinition{
			{
				Component: component.Component{
					Kind: "PartialKind",
					Version: "v1.0.0",
				},
				Model: &model.ModelDefinition{
					Name: "meshery-core",
					Model: model.Model{
						Version: "v1.0.0",
					},
				},
			},
		},
	}

	// Should not panic when encountering partial/missing sub-fields
	err = mapModelRelatedData(regManager, patternFile)
	require.NoError(t, err)
}

func TestMapModelRelatedData_NilRegistrantInRegisteredModel(t *testing.T) {
	regManager := setupTestRegistry(t)

	// Register a component whose model definition has a nil/omitted Registrant
	def := component.ComponentDefinition{
		Component: component.Component{
			Kind:    "NilRegistrantKind",
			Version: "v1.0.0",
		},
		Model: &model.ModelDefinition{
			Name: "meshery-core",
			Model: model.Model{
				Version: "v1.0.0",
			},
			// Registrant intentionally omitted to simulate nil/zero-value Registrant
		},
	}

	conn := connection.Connection{}
	_, _, err := regManager.RegisterEntity(conn, &def)
	require.NoError(t, err)

	patternFile := &pattern.PatternFile{
		Components: []*component.ComponentDefinition{
			{
				Component: component.Component{
					Kind:    "NilRegistrantKind",
					Version: "v1.0.0",
				},
				Model: &model.ModelDefinition{
					Name: "meshery-core",
					Model: model.Model{
						Version: "v1.0.0",
					},
					// Registrant again omitted; mapModelRelatedData should
					// handle this without panicking when accessing Status.
				},
			},
		},
	}

	// Should not panic when encountering a registered component with nil Registrant
	err = mapModelRelatedData(regManager, patternFile)
	require.NoError(t, err)
}
