package models

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/meshes"
)

type HandlerInterface interface {
	AuthMiddleware(next http.Handler) http.Handler

	IndexHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	K8SConfigHandler(ctx context.Context, meshClient meshes.MeshClient) func(w http.ResponseWriter, r *http.Request)
	LogoutHandler(w http.ResponseWriter, req *http.Request)
	DashboardHandler(ctx context.Context, meshClient meshes.MeshClient, w http.ResponseWriter, req *http.Request)
	LoadTestHandler(w http.ResponseWriter, req *http.Request)
	MeshOpsHandler(ctx context.Context, meshClient meshes.MeshClient) func(w http.ResponseWriter, req *http.Request)
}

type HandlerConfig struct {
	SessionName   string
	RefCookieName string
	// SessionUserKey
	// SessionSecret string

	SessionStore sessions.Store

	SaaSTokenName string
	SaaSBaseURL   string

	ByPassAuth bool

	// TODO: accept this from the UI. . . so as to dynamically change it
	LoadTestURL string

	FortioURL string
}

//
// sessionUserKey       = "twitterID"
// 	sessionUserName      = "twitterUserName"
// cookieSuffix         = "_referrer"
// saasTokenName        = "meshery_saas"
