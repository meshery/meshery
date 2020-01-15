package models

import (
	"net/http"

	"time"

	"github.com/gorilla/sessions"
	"github.com/vmihailenco/taskq"
)

// HandlerInterface defines the methods a Handler should define
type HandlerInterface interface {
	GetProviderType() ProviderType

	AuthMiddleware(http.Handler) http.Handler
	SessionInjectorMiddleware(func(http.ResponseWriter, *http.Request, *sessions.Session, *Preference, *User)) http.Handler

	LoginHandler(w http.ResponseWriter, r *http.Request, fromMiddleWare bool)
	LogoutHandler(w http.ResponseWriter, req *http.Request)
	UserHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, prefObj *Preference, user *User)

	K8SConfigHandler(w http.ResponseWriter, r *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	GetContextsFromK8SConfig(w http.ResponseWriter, req *http.Request)
	KubernetesPingHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	InstalledMeshesHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)

	LoadTestHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	CollectStaticMetrics(config *SubmitMetricsConfig) error
	FetchResultsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)

	MeshAdapterConfigHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	MeshOpsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	GetAllAdaptersHandler(w http.ResponseWriter, req *http.Request)
	EventStreamHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	AdapterPingHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)

	GrafanaConfigHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	GrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	GrafanaQueryHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	GrafanaQueryRangeHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	GrafanaPingHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	SaveSelectedGrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)

	PrometheusConfigHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	GrafanaBoardImportForPrometheusHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	PrometheusQueryHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	PrometheusQueryRangeHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	PrometheusPingHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	PrometheusStaticBoardHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	SaveSelectedPrometheusBoardsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)

	LoadTestPrefencesHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
	AnonymousStatsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)

	SessionSyncHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, prefObj *Preference, user *User)
}

// HandlerConfig holds all the config pieces needed by handler methods
type HandlerConfig struct {
	// SessionName   string
	// RefCookieName string

	// SessionStore sessions.Store

	// SaaSTokenName string
	// SaaSBaseURL   string

	AdapterTracker AdaptersTrackerInterface
	QueryTracker   QueryTrackerInterface

	Queue taskq.Queue

	KubeConfigFolder string

	GrafanaClient         *GrafanaClient
	GrafanaClientForQuery *GrafanaClient

	PrometheusClient         *PrometheusClient
	PrometheusClientForQuery *PrometheusClient

	Provider Provider
}

// SubmitMetricsConfig is used to store config used for submitting metrics
type SubmitMetricsConfig struct {
	TestUUID, ResultID, PromURL string
	StartTime, EndTime          time.Time
	// TokenKey,
	TokenVal string
}
