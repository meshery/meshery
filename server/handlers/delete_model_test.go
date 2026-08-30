package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshkit/database"
	_models "github.com/meshery/meshkit/models/meshmodel/core/v1beta1"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	_model "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta3/component"
	"github.com/stretchr/testify/require"
)

func newDeleteModelFixture(t *testing.T) (*Handler, *database.Handler) {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	require.NoError(t, err, "open database")

	err = db.AutoMigrate(
		&_model.ModelDefinition{},
		&component.ComponentDefinition{},
		&relationship.RelationshipDefinition{},
		&_models.PolicyDefinition{},
		&registry.Registry{},
	)
	require.NoError(t, err, "migrate tables")

	h := &Handler{
		log:       newTestLogger(t),
		dbHandler: &db,
	}

	return h, &db
}

func deleteModelRequest(h *Handler, modelID string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodDelete, "/api/meshmodels/models/"+modelID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": modelID})
	rec := httptest.NewRecorder()

	h.DeleteModel(rec, req, nil, nil, nil)
	return rec
}

func TestDeleteModel_NotFound(t *testing.T) {
	h, _ := newDeleteModelFixture(t)
	unknownID := uuid.Must(uuid.NewV4()).String()

	rec := deleteModelRequest(h, unknownID)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteModel_InvalidUUID(t *testing.T) {
	h, _ := newDeleteModelFixture(t)

	rec := deleteModelRequest(h, "invalid-uuid")

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteModel_WithPolicyDefinition(t *testing.T) {
	h, db := newDeleteModelFixture(t)

	modelUUID := uuid.Must(uuid.NewV4())
	compUUID := uuid.Must(uuid.NewV4())
	relUUID := uuid.Must(uuid.NewV4())
	policyUUID := uuid.Must(uuid.NewV4())

	// 1. Create a model definition
	modelDef := _model.ModelDefinition{
		ID:   modelUUID,
		Name: "test-model",
	}
	require.NoError(t, db.Create(&modelDef).Error)

	// 2. Create associated component definition, relationship definition, policy definition
	compDef := component.ComponentDefinition{
		ID:      compUUID,
		ModelID: &modelUUID,
	}
	require.NoError(t, db.Create(&compDef).Error)

	relDef := relationship.RelationshipDefinition{
		ID:      relUUID,
		ModelId: &modelUUID,
		Kind:    "test-relationship",
	}
	require.NoError(t, db.Create(&relDef).Error)

	policyDef := _models.PolicyDefinition{
		ID:      policyUUID,
		ModelID: modelUUID,
	}
	require.NoError(t, db.Create(&policyDef).Error)

	// 3. Create associated registry entries
	regModel := registry.Registry{
		ID:     uuid.Must(uuid.NewV4()),
		Entity: modelUUID,
		Type:   entity.Model,
	}
	regComp := registry.Registry{
		ID:     uuid.Must(uuid.NewV4()),
		Entity: compUUID,
		Type:   entity.ComponentDefinition,
	}
	regRel := registry.Registry{
		ID:     uuid.Must(uuid.NewV4()),
		Entity: relUUID,
		Type:   entity.RelationshipDefinition,
	}
	regPolicy := registry.Registry{
		ID:     uuid.Must(uuid.NewV4()),
		Entity: policyUUID,
		Type:   entity.PolicyDefinition,
	}
	require.NoError(t, db.Create(&regModel).Error)
	require.NoError(t, db.Create(&regComp).Error)
	require.NoError(t, db.Create(&regRel).Error)
	require.NoError(t, db.Create(&regPolicy).Error)

	// 4. Invoke DeleteModel
	rec := deleteModelRequest(h, modelUUID.String())

	// 5. Verify HTTP StatusNoContent
	require.Equal(t, http.StatusNoContent, rec.Code)

	// 6. Verify model, components, relationships, policies, and registry entries are deleted
	var count int64
	db.Model(&_model.ModelDefinition{}).Where("id = ?", modelUUID).Count(&count)
	require.Equal(t, int64(0), count, "model should be deleted")

	db.Model(&component.ComponentDefinition{}).Where("model_id = ?", modelUUID).Count(&count)
	require.Equal(t, int64(0), count, "components should be deleted")

	db.Model(&relationship.RelationshipDefinition{}).Where("model_id = ?", modelUUID).Count(&count)
	require.Equal(t, int64(0), count, "relationships should be deleted")

	db.Model(&_models.PolicyDefinition{}).Where("modelID = ?", modelUUID).Count(&count)
	require.Equal(t, int64(0), count, "policies should be deleted")

	db.Model(&registry.Registry{}).Where("entity IN (?)", []uuid.UUID{modelUUID, compUUID, relUUID, policyUUID}).Count(&count)
	require.Equal(t, int64(0), count, "registry entries should be deleted")
}
