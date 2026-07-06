package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	connectionv1beta1 "github.com/meshery/schemas/models/v1beta1/connection"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
)

// connectionDefinitionsFromEntities narrows the polymorphic registry entities
// returned by the registry manager into concrete connection definitions.
func connectionDefinitionsFromEntities(entities []entity.Entity) []*connectionv1beta3.ConnectionDefinition {
	defs := make([]*connectionv1beta3.ConnectionDefinition, 0, len(entities))
	for _, en := range entities {
		if def, ok := en.(*connectionv1beta3.ConnectionDefinition); ok {
			defs = append(defs, def)
		}
	}
	return defs
}

// GetConnectionDefinitions returns a paginated list of connection definitions
// registered in the registry, optionally filtered by model or kind.
//
// swagger:route GET /api/meshmodels/connections GetConnectionDefinitions
// responses:
//
//	200: ConnectionDefinitionPage
func (h *Handler) GetConnectionDefinitions(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	enc := json.NewEncoder(rw)

	page, offset, limit, search, order, sort, _ := getPaginationParams(r)
	q := r.URL.Query()

	entities, count, _, err := h.registryManager.GetEntities(&regv1beta1.ConnectionFilter{
		Name:      search,
		Greedy:    search != "",
		Kind:      q.Get("kind"),
		Type:      q.Get("type"),
		ModelName: q.Get("model"),
		Version:   q.Get("version"),
		Status:    q.Get("status"),
		OrderOn:   order,
		Sort:      sort,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		h.log.Error(ErrQueryGet("connection definitions"))
		writeMeshkitError(rw, ErrQueryGet("connection definitions"), http.StatusInternalServerError)
		return
	}

	pgSize := limit
	if limit == 0 {
		pgSize = int(count)
	}

	response := connectionv1beta3.ConnectionDefinitionPage{
		Page:                  page,
		PageSize:              pgSize,
		TotalCount:            int(count),
		ConnectionDefinitions: connectionDefinitionsFromEntities(entities),
	}

	if err := enc.Encode(response); err != nil {
		h.log.Error(models.ErrEncoding(err, "connection definitions"))
		writeMeshkitError(rw, models.ErrEncoding(err, "connection definitions"), http.StatusInternalServerError)
	}
}

// GetConnectionDefinitionByID returns a single connection definition by its ID.
//
// swagger:route GET /api/meshmodels/connections/{connectionDefinitionId} GetConnectionDefinitionByID
// responses:
//
//	200: ConnectionDefinition
func (h *Handler) GetConnectionDefinitionByID(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")
	connectionDefinitionID := mux.Vars(r)["connectionDefinitionId"]
	if _, err := uuid.FromString(connectionDefinitionID); err != nil {
		h.log.Error(models.ErrInvalidUUID(err))
		writeMeshkitError(rw, models.ErrInvalidUUID(err), http.StatusBadRequest)
		return
	}

	entities, _, _, err := h.registryManager.GetEntities(&regv1beta1.ConnectionFilter{
		Id:    connectionDefinitionID,
		Limit: 1,
	})
	if err != nil {
		h.log.Error(ErrQueryGet("connection definition"))
		writeMeshkitError(rw, ErrQueryGet("connection definition"), http.StatusInternalServerError)
		return
	}
	defs := connectionDefinitionsFromEntities(entities)
	if len(defs) == 0 {
		writeJSONError(rw, fmt.Sprintf("connection definition with id %s not found", connectionDefinitionID), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(rw).Encode(defs[0]); err != nil {
		h.log.Error(models.ErrEncoding(err, "connection definition"))
		writeMeshkitError(rw, models.ErrEncoding(err, "connection definition"), http.StatusInternalServerError)
	}
}

// RegisterConnectionDefinition registers a new connection definition into the
// registry under the model (and registrant) carried in the request body.
//
// swagger:route POST /api/meshmodels/connections RegisterConnectionDefinition
// responses:
//
//	201: ConnectionDefinition
func (h *Handler) RegisterConnectionDefinition(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	obj := "connection definition"

	def := connectionv1beta3.ConnectionDefinition{}
	if err := json.NewDecoder(r.Body).Decode(&def); err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(rw, models.ErrUnmarshal(err, obj), http.StatusBadRequest)
		return
	}

	if def.ModelReference == nil {
		err := ErrFailToSave(fmt.Errorf("a modelReference (with registrant) is required to register a connection definition"), obj)
		h.log.Error(err)
		writeMeshkitError(rw, err, http.StatusBadRequest)
		return
	}

	registrant := connectionv1beta1.Connection{Kind: def.ModelReference.Registrant.Kind}
	if _, _, err := h.registryManager.RegisterEntity(registry.RegistrantHostToV1beta3(registrant), &def); err != nil {
		_err := ErrFailToSave(err, obj)
		h.log.Error(_err)
		writeMeshkitError(rw, _err, http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(def); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
	}
}

// UpdateConnectionDefinition updates an existing connection definition.
//
// swagger:route PUT /api/meshmodels/connections/{connectionDefinitionId} UpdateConnectionDefinition
// responses:
//
//	200: ConnectionDefinition
func (h *Handler) UpdateConnectionDefinition(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	obj := "connection definition"
	connectionDefinitionID := mux.Vars(r)["connectionDefinitionId"]
	id, err := uuid.FromString(connectionDefinitionID)
	if err != nil {
		h.log.Error(models.ErrInvalidUUID(err))
		writeMeshkitError(rw, models.ErrInvalidUUID(err), http.StatusBadRequest)
		return
	}

	def := connectionv1beta3.ConnectionDefinition{}
	if err := json.NewDecoder(r.Body).Decode(&def); err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(rw, models.ErrUnmarshal(err, obj), http.StatusBadRequest)
		return
	}
	def.ID = id

	// Persist only the connection definition row; associations (e.g. model) are
	// managed through registration, not this update.
	if err := h.dbHandler.Omit(clause.Associations).Save(&def).Error; err != nil {
		_err := ErrFailToSave(err, obj)
		h.log.Error(_err)
		writeMeshkitError(rw, _err, http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(rw).Encode(def); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
	}
}

// DeleteConnectionDefinition removes a connection definition and its registry
// entry. Mirrors the registry-cleanup approach used by DeleteModel.
//
// swagger:route DELETE /api/meshmodels/connections/{connectionDefinitionId} DeleteConnectionDefinition
// responses:
//
//	204:
func (h *Handler) DeleteConnectionDefinition(rw http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	connectionDefinitionID := mux.Vars(r)["connectionDefinitionId"]
	id, err := uuid.FromString(connectionDefinitionID)
	if err != nil {
		h.log.Error(models.ErrInvalidUUID(err))
		writeMeshkitError(rw, models.ErrInvalidUUID(err), http.StatusBadRequest)
		return
	}

	err = h.dbHandler.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("entity = ? AND type = ?", id, entity.ConnectionDefinition).
			Delete(&registry.Registry{}).Error; err != nil {
			return err
		}
		return tx.Delete(&connectionv1beta3.ConnectionDefinition{}, "id = ?", id).Error
	})
	if err != nil {
		_err := ErrFailToDelete(err, "connection definition")
		h.log.Error(_err)
		writeMeshkitError(rw, _err, http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
