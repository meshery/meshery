package router

import (
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func registerPerformanceRoutes(gMux *mux.Router, h models.HandlerInterface) {
	gMux.Handle("/api/perf/profile", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.LoadTestHandler), models.ProviderAuth))).
		Methods("GET", "POST")
	gMux.Handle("/api/perf/profile/result", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchAllResultsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/perf/profile/result/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetResultHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/mesh", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetSMPServiceMeshes), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/smi/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchSmiResultsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/smi/results/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchSingleSmiResultHandler), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/user/performance/profiles", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetPerformanceProfilesHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchAllResultsHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetPerformanceProfileHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeletePerformanceProfileHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/user/performance/profiles", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SavePerformanceProfileHandler), models.ProviderAuth))).
		Methods("POST")
	gMux.Handle("/api/user/performance/profiles/{id}/run", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.KubernetesMiddleware(h.LoadTestHandler)), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/performance/profiles/{id}/results", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.FetchResultsHandler), models.ProviderAuth))).
		Methods("GET")

	gMux.Handle("/api/user/schedules", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetSchedulesHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/schedules/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.GetScheduleHandler), models.ProviderAuth))).
		Methods("GET")
	gMux.Handle("/api/user/schedules/{id}", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.DeleteScheduleHandler), models.ProviderAuth))).
		Methods("DELETE")
	gMux.Handle("/api/user/schedules", h.ProviderMiddleware(h.AuthMiddleware(h.SessionInjectorMiddleware(h.SaveScheduleHandler), models.ProviderAuth))).
		Methods("POST")
}
