package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// registerPatternRoutes registers routes for patterns/designs, filters,
// mesh models, schemas, and content sharing.
func registerPatternRoutes(gMux *mux.Router, h models.HandlerInterface) {
	// Patterns / Designs
	gMux.Handle("/api/pattern/deploy", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.PatternFileHandler)), models.ProviderAuth))).
		Methods("POST", "DELETE")
	gMux.Handle("/api/pattern", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PatternFileRequestHandler), models.ProviderAuth))).
		Methods("POST", "GET")
	gMux.Handle("/api/pattern/import", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DesignFileImportHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/pattern/{sourcetype}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DesignFileRequestHandlerWithSourceType), models.ProviderAuth))).
		Methods("POST", "PUT")
	gMux.Handle("/api/pattern/types", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryDesignTypesHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/pattern/catalog", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetCatalogMesheryPatternsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/pattern/catalog/unpublish", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UnPublishCatalogPatternHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/pattern/catalog/publish", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PublishCatalogPatternHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/pattern/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryPatternHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/pattern/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMesheryPatternHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/pattern/clone/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.CloneMesheryPatternHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/pattern/download/{id}/{sourcetype}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryPatternSourceHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/pattern/download/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DownloadMesheryPatternHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/patterns/delete", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMultiMesheryPatternsHandler), models.ProviderAuth))).
		Methods("POST")

	// Filters
	gMux.Handle("/api/filter/deploy", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.FilterFileHandler)), models.ProviderAuth))).
		Methods("POST", "DELETE")
	gMux.Handle("/api/filter", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FilterFileRequestHandler), models.ProviderAuth))).
		Methods("POST", "GET")
	gMux.Handle("/api/filter/catalog", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetCatalogMesheryFiltersHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/filter/catalog/publish", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.PublishCatalogFilterHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/filter/catalog/unpublish", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UnPublishCatalogFilterHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/filter/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryFilterHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/filter/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMesheryFilterHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/filter/download/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMesheryFilterFileHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/filter/clone/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.CloneMesheryFilterHandler), models.ProviderAuth))).
		Methods("POST")

	// Content sharing
	gMux.Handle("/api/content/design/share", h.ProviderMiddleware((h.AuthMiddleware(h.SessionInjectorMiddleware(h.ShareDesignHandler), models.ProviderAuth)))).
		Methods("POST")
	gMux.Handle("/api/content/filter/share", h.ProviderMiddleware((h.AuthMiddleware(h.SessionInjectorMiddleware(h.ShareFilterHandler), models.ProviderAuth)))).
		Methods("POST")

	// Schemas
	gMux.Handle("/api/schema/resource/{resourceName}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.HandleResourceSchemas), models.NoAuth))).
		Methods("GET")

	// Mesh Models
	gMux.Handle("/api/meshmodels/validate", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.ValidationHandler), models.NoAuth))).Methods("POST")
	gMux.Handle("/api/meshmodels/components", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.RegisterMeshmodelComponents), models.NoAuth))).Methods("POST")
	gMux.Handle("/api/meshmodel/components/register", h.ProviderMiddleware((http.HandlerFunc(h.RegisterMeshmodelComponents)))).Methods("POST")
	gMux.Handle("/api/meshmodels/{entityType}/status", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.UpdateEntityStatus), models.NoAuth))).Methods("POST")
	gMux.Handle("/api/meshmodels/components", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelComponents), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/categories", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelCategories), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/models", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelModels), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/models/{model}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelModelsByName), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/models/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteModel), models.ProviderAuth))).Methods("DELETE")
	gMux.Handle("/api/meshmodels/registrants", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelRegistrants), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/register", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.RegisterMeshmodels), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/meshmodels/categories/{category}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelCategoriesByName), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/categories/{category}/models", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelModelsByCategories), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/categories/{category}/models/{model}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelModelsByCategoriesByModel), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/categories/{category}/models/{model}/components", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentByModelByCategory), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/categories/{category}/models/{model}/components/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentsByNameByModelByCategory), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/categories/{category}/components", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentByCategory), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/categories/{category}/components/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentsByNameByCategory), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/components/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelComponentsByName), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/generate", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.MeshModelGenerationHandler), models.NoAuth))).Methods("POST")
	gMux.Handle("/api/meshmodels/relationships", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelRelationships), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/models/{model}/components", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentByModel), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/models/{model}/components/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelComponentsByNameByModel), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/models/{model}/relationships", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelRelationships), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/models/{model}/relationships/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetMeshmodelRelationshipByName), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/relationships", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.RegisterMeshmodelRelationships), models.NoAuth))).Methods("POST")
	gMux.Handle("/api/meshmodels/relationships/evaluate", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.EvaluateRelationshipPolicy), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/meshmodels/models/{model}/policies", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelPolicies), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/models/{model}/policies/{name}", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.GetAllMeshmodelPoliciesByName), models.NoAuth))).Methods("GET")
	gMux.Handle("/api/meshmodels/export", h.ProviderMiddleware(h.AuthMiddleware(http.HandlerFunc(h.ExportModel), models.NoAuth))).Methods("GET")

	// MeshSync resources (grouped under patterns for proximity to mesh data)
	gMux.Handle("/api/system/meshsync/resources", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMeshSyncResources), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/meshsync/resources/summary", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMeshSyncResourcesSummary), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/system/meshsync/resources/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteMeshSyncResource), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/system/meshsync/resources/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetMeshSyncResourceByID), models.ProviderAuth))).
		Methods("GET")
}
