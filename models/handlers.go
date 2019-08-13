package models

import (
	"net/http"

	"time"

	"github.com/gorilla/sessions"
	"github.com/vmihailenco/taskq"
)

type HandlerInterface interface {
	AuthMiddleware(next http.Handler) http.Handler

	LoginHandler(w http.ResponseWriter, r *http.Request)
	LogoutHandler(w http.ResponseWriter, req *http.Request)
	UserHandler(w http.ResponseWriter, r *http.Request)

	K8SConfigHandler(w http.ResponseWriter, r *http.Request)
	GetContextsFromK8SConfig(w http.ResponseWriter, req *http.Request)

	LoadTestHandler(w http.ResponseWriter, req *http.Request)
	CollectStaticMetrics(config *SubmitMetricsConfig) error
	FetchResultsHandler(w http.ResponseWriter, req *http.Request)

	MeshAdapterConfigHandler(w http.ResponseWriter, req *http.Request)
	MeshOpsHandler(w http.ResponseWriter, req *http.Request)
	GetAllAdaptersHandler(w http.ResponseWriter, req *http.Request)
	EventStreamHandler(w http.ResponseWriter, req *http.Request)

	GrafanaConfigHandler(w http.ResponseWriter, req *http.Request)
	GrafanaBoardsHandler(w http.ResponseWriter, req *http.Request)
	GrafanaQueryHandler(w http.ResponseWriter, req *http.Request)
	GrafanaQueryRangeHandler(w http.ResponseWriter, req *http.Request)
	SaveSelectedGrafanaBoardsHandler(w http.ResponseWriter, req *http.Request)

	PrometheusConfigHandler(w http.ResponseWriter, req *http.Request)
	GrafanaBoardImportForPrometheusHandler(w http.ResponseWriter, req *http.Request)
	PrometheusQueryHandler(w http.ResponseWriter, req *http.Request)
	PrometheusQueryRangeHandler(w http.ResponseWriter, req *http.Request)
	PrometheusStaticBoardHandler(w http.ResponseWriter, req *http.Request)
	SaveSelectedPrometheusBoardsHandler(w http.ResponseWriter, req *http.Request)

	SessionSyncHandler(w http.ResponseWriter, req *http.Request)
}

type HandlerConfig struct {
	SessionName   string
	RefCookieName string

	SessionStore sessions.Store

	SaaSTokenName string
	SaaSBaseURL   string

	AdapterTracker AdaptersTrackerInterface
	QueryTracker   QueryTrackerInterface

	Queue taskq.Queue

	SessionPersister SessionPersister

	KubeConfigFolder string
}

type SubmitMetricsConfig struct {
	TestUUID, ResultID, PromURL string
	StartTime, EndTime          time.Time
	TokenKey, TokenVal          string
}
