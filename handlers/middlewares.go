//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"net/http"

	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/sirupsen/logrus"
)

// ProviderMiddleware is a middleware to validate if a provider is set
func (h *Handler) ProviderMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		var providerName string
		var provider models.Provider
		ck, err := req.Cookie(h.config.ProviderCookieName)
		if err == nil {
			providerName = ck.Value
		} else {
			providerName = req.Header.Get(h.config.ProviderCookieName)
		}
		if providerName != "" {
			provider = h.config.Providers[providerName]
		}
		if provider == nil {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		ctx := context.WithValue(req.Context(), models.ProviderCtxKey, provider) // nolint
		req1 := req.WithContext(ctx)
		next.ServeHTTP(w, req1)
	}
	return http.HandlerFunc(fn)
}

// AuthMiddleware is a middleware to validate if a user is authenticated
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		// logrus.Debugf("models.ProviderCtxKey %s", models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		// logrus.Debugf("provider %s", provider)
		isValid := h.validateAuth(provider, req)
		// logrus.Debugf("validate auth: %t", isValid)
		if !isValid {
			// if h.GetProviderType() == models.RemoteProviderType {
			// 	http.Redirect(w, req, "/user/login", http.StatusFound)
			// } else { // Local Provider
			// 	h.LoginHandler(w, req)
			// }
			// return
			if provider.GetProviderType() == models.RemoteProviderType {
				provider.Logout(w, req)
				return
			}
			// Local Provider
			h.LoginHandler(w, req, provider, true)
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func (h *Handler) validateAuth(provider models.Provider, req *http.Request) bool {
	if err := provider.GetSession(req); err == nil {
		// logrus.Debugf("session: %v", sess)
		return true
	}
	// logrus.Errorf("session invalid, error: %v", err)
	return false
}

// SessionInjectorMiddleware - is a middleware which injects user and session object
func (h *Handler) SessionInjectorMiddleware(next func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		providerI := req.Context().Value(models.ProviderCtxKey)
		provider, ok := providerI.(models.Provider)
		if !ok {
			http.Redirect(w, req, "/provider", http.StatusFound)
			return
		}
		// ensuring session is intact before running load test
		err := provider.GetSession(req)
		if err != nil {
			provider.Logout(w, req)
			logrus.Errorf("Error: unable to get session: %v", err)
			http.Error(w, "unable to get session", http.StatusUnauthorized)
			return
		}

		user, _ := provider.GetUserDetails(req)
		prefObj, err := provider.ReadFromPersister(user.UserID)
		if err != nil {
			logrus.Warn("unable to read session from the session persister, starting with a new one")
		}

		token := provider.UpdateToken(w, req)

		ctx := context.WithValue(req.Context(), models.TokenCtxKey, token)
		ctx = context.WithValue(ctx, models.PerfObjCtxKey, prefObj)
		ctx = context.WithValue(ctx, models.UserCtxKey, user)
		ctx = context.WithValue(ctx, models.BrokerURLCtxKey, h.config.BrokerEndpointURL) // nolint

		k8scontext, err := h.GetCurrentContext(token, provider)
		if err != nil {
			logrus.Warn("failed to find kubernetes context")

			// Set some defaults in the context so that the casting doesn't fails
			ctx = context.WithValue(ctx, models.KubeContextKey, nil)
			ctx = context.WithValue(ctx, models.KubeHanderKey, nil)
			ctx = context.WithValue(ctx, models.KubeConfigKey, nil)
		} else {
			cfg, err := k8scontext.GenerateKubeConfig()
			if err != nil {
				logrus.Warn("failed to load kube config for the user: ", err)
			}

			// Create mesherykube handler
			client, err := kubernetes.New(cfg)
			if err != nil {
				logrus.Warn("failed to create kubeconfig handler for the user")
				// http.Error(w, "failed to create kubeconfig handler for the user", http.StatusInternalServerError)
				// return
			}

			ctx = context.WithValue(ctx, models.KubeContextKey, k8scontext)
			ctx = context.WithValue(ctx, models.KubeHanderKey, client)
			ctx = context.WithValue(ctx, models.KubeConfigKey, cfg)
		}

		// Identify custom contexts, if provided
		k8sContextIDs := req.URL.Query()["contexts"]
		k8scontexts := []models.K8sContext{}

		if len(k8sContextIDs) == 1 && k8sContextIDs[0] == "all" {
			contexts, err := provider.LoadAllK8sContext(token)
			if err != nil {
				logrus.Warn("failed to load all k8scontext")
			}

			for _, c := range contexts {
				if c != nil {
					k8scontexts = append(k8scontexts, *c)
				}
			}
		} else {
			for _, kctxID := range k8sContextIDs {
				kctx, err := provider.GetK8sContext(token, kctxID)
				if err != nil {
					logrus.Warn("invalid context ID found")
					continue
				}

				k8scontexts = append(k8scontexts, kctx)
			}
		}

		ctx = context.WithValue(ctx, models.KubeClustersKey, k8scontexts)

		req1 := req.WithContext(ctx)

		next(w, req1, prefObj, user, provider)
	})
}

// GraphqlSessionInjectorMiddleware - is a middleware which injects user and session object
func (h *Handler) GraphqlMiddleware(next http.Handler) func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider) {
	return func(w http.ResponseWriter, req *http.Request, pref *models.Preference, user *models.User, prov models.Provider) {
		next.ServeHTTP(w, req)
	}
}
