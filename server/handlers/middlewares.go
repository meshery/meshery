// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/server/models"
	"github.com/sirupsen/logrus"
)

const providerQParamName = "provider"

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
			// allow provider to be set using query parameter
			// this is OK since provider information is not sensitive
			if providerName == "" {
				providerName = req.URL.Query().Get(providerQParamName)
			}
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
				provider.HandleUnAuthenticated(w, req)
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

// MesheryControllersMiddleware is a middleware that is responsible for handling meshery controllers(operator, meshsync and broker) related stuff such as
// getting status, reconciling their deployments etc.
func (h *Handler) MesheryControllersMiddleware(next func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)) func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider) {
	return func(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
		ctx := req.Context()
		mk8sContexts, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
		if !ok || len(mk8sContexts) == 0 {
			h.log.Error(ErrInvalidK8SConfig)
			// this should not block the request
			next(w, req, prefObj, user, provider)
			return
		}

		// 1. get the status of controller deployments for each cluster and make sure that all the contexts have meshery controllers deployed
		ctrlHlpr := h.MesheryCtrlsHelper.UpdateCtxControllerHandlers(mk8sContexts).UpdateOperatorsStatusMap().DeployUndeployedOperators()
		ctx = context.WithValue(ctx, models.MesheryControllerHandlersKey, h.MesheryCtrlsHelper.GetControllerHandlersForEachContext())

		// 2. make sure that the data from meshsync for all the clusters are persisted properly
		ctrlHlpr.UpdateMeshsynDataHandlers()
		ctx = context.WithValue(ctx, models.MeshSyncDataHandlersKey, h.MesheryCtrlsHelper.GetMeshSyncDataHandlersForEachContext())

		req1 := req.WithContext(ctx)
		next(w, req1, prefObj, user, provider)
	}
}

// KubernetesMiddleware is a middleware that is responsible for handling kubernetes related stuff such as
// setting contexts, component generation etc.
func (h *Handler) KubernetesMiddleware(next func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)) func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider) {
	return func(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
		ctx := req.Context()
		token, ok := ctx.Value(models.TokenCtxKey).(string)
		if !ok {
			err := ErrRetrieveUserToken(fmt.Errorf("failed to retrieve user token"))
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		contexts, err := provider.LoadAllK8sContext(token)
		if err != nil || len(contexts) == 0 { //Try to load the contexts when there are no contexts available
			logrus.Warn("failed to get kubernetes contexts")
			// only the contexts that are successfully pinged will be persisted
			contexts, err = h.LoadContextsAndPersist(token, provider)
			if err != nil {
				logrus.Warn("failed to load kubernetes contexts: ", err.Error())
			}
		}

		// register kubernetes components
		h.K8sCompRegHelper.UpdateContexts(contexts).RegisterComponents(contexts, RegisterK8sComponents, h.EventsBuffer)

		// Identify custom contexts, if provided
		k8sContextIDs := req.URL.Query()["contexts"]
		k8scontexts := []models.K8sContext{}    //The contexts passed by the user
		allk8scontexts := []models.K8sContext{} //All contexts to track all the connected clusters

		if len(k8sContextIDs) == 0 { //This is for backwards compabitibility with clients. This will work fine for single cluster.
			//For multi cluster, it is expected of clients to explicitly pass the k8scontextID.
			//So for now, randomly one of the contexts from available ones will be pushed to the array to stop anything from breaking in case of no contexts received(with single cluster, the behavior would be as expected).
			if len(contexts) > 0 && contexts[0] != nil {
				k8scontexts = append(k8scontexts, *contexts[0])
			}
		} else if len(k8sContextIDs) == 1 && k8sContextIDs[0] == "all" {
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
		for _, k8scontext := range contexts {
			if k8scontext != nil {
				allk8scontexts = append(allk8scontexts, *k8scontext)
			}
		}

		ctx = context.WithValue(ctx, models.KubeClustersKey, k8scontexts)
		ctx = context.WithValue(ctx, models.AllKubeClusterKey, allk8scontexts)
		req1 := req.WithContext(ctx)
		next(w, req1, prefObj, user, provider)
	}
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
