package handlers

import (
	"net/http"
)

//func (h *Handler) GraphqlSystemHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
//	queryEndpoint := "/api/system/graphql/query"
//	playgroundEndpoint := "/api/system/graphql/playground"

//	if req.URL.Path == queryEndpoint {
//		provider.GetGraphqlHandler().ServeHTTP(w, req)
//	} else if req.URL.Path == playgroundEndpoint {
//		provider.GetGraphqlPlayground().ServeHTTP(w, req)
//	} else {
//		http.Error(w, "Invalid endpoint", http.StatusInternalServerError)
//	}
//}

func (h *Handler) GetGraphQLHandler() http.Handler {
	return h.config.GraphQLHandler
}

func (h *Handler) GetGraphQLPlaygroundHandler() http.Handler {
	return h.config.GraphQLPlaygroundHandler
}
