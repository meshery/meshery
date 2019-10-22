package models

import (
	"net/http"

	"time"

	"github.com/gorilla/sessions"
	"github.com/vmihailenco/taskq"
)

// HandlerInterface defines the methods a Handler should define
type HandlerInterface interface {
	AuthMiddleware(http.Handler) http.Handler
	SessionInjectorMiddleware(func(http.ResponseWriter, *http.Request, *sessions.Session, *User)) http.Handler

	LoginHandler(w http.ResponseWriter, r *http.Request)
	LogoutHandler(w http.ResponseWriter, req *http.Request)
	UserHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, user *User)

	K8SConfigHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, user *User)
	GetContextsFromK8SConfig(w http.ResponseWriter, req *http.Request)
	KubernetesPingHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	InstalledMeshesHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)

	LoadTestHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	CollectStaticMetrics(config *SubmitMetricsConfig) error
	FetchResultsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)

	MeshAdapterConfigHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	MeshOpsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	GetAllAdaptersHandler(w http.ResponseWriter, req *http.Request)
	EventStreamHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	AdapterPingHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)

	GrafanaConfigHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	GrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	GrafanaQueryHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	GrafanaQueryRangeHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	SaveSelectedGrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)

	PrometheusConfigHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	GrafanaBoardImportForPrometheusHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	PrometheusQueryHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	PrometheusQueryRangeHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	PrometheusStaticBoardHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
	SaveSelectedPrometheusBoardsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)

	SessionSyncHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *User)
}

// HandlerConfig holds all the config pieces needed by handler methods
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

// SubmitMetricsConfig is used to store config used for submitting metrics
type SubmitMetricsConfig struct {
	TestUUID, ResultID, PromURL string
	StartTime, EndTime          time.Time
	TokenKey, TokenVal          string
}
