// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/layer5io/meshery/server/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
		ctx := context.WithValue(req.Context(), models.ProviderCtxKey, provider) // nolint

		// Incase Meshery is configured for deployments scenario: Istio, Azure Kubernetes Service etc
		// then we can expect a MESHERY_SERVER_CALLBACK_URL in env var
		callbackURL := viper.GetString("MESHERY_SERVER_CALLBACK_URL")
		if callbackURL == "" {
			// if MESHERY_SERVER_CALLBACK_URL is not set then we can assume standard CALLBACK_URL
			callbackURL = "http://" + req.Host + "/api/user/token" // Hard coding the path because this is what meshery expects
		}
		ctx = context.WithValue(ctx, models.MesheryServerCallbackURL, callbackURL)
		_url, err := url.Parse(callbackURL)
		if err != nil {
			logrus.Errorf("Error parsing callback url: %v", err)
		} else {
			ctx = context.WithValue(ctx, models.MesheryServerURL, fmt.Sprintf("%s://%s", _url.Scheme, _url.Host))
		}
		req1 := req.WithContext(ctx)
		next.ServeHTTP(w, req1)
	}
	return http.HandlerFunc(fn)
}

// AuthMiddleware is a middleware to validate if a user is authenticated
func (h *Handler) AuthMiddleware(next http.Handler, auth models.AuthenticationMechanism) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		providerH := h.Provider
		if auth == models.NoAuth && providerH != "" {
			auth = models.ProviderAuth //If a provider is enforced then use provider authentication even in case of NoAuth
		}
		switch auth {
		// case models.NoAuth:
		// 	if providerH != "" {
		// 		w.WriteHeader(http.StatusUnauthorized)
		// 		return
		// 	}
		case models.ProviderAuth:
			providerI := req.Context().Value(models.ProviderCtxKey)
			// logrus.Debugf("models.ProviderCtxKey %s", models.ProviderCtxKey)
			provider, ok := providerI.(models.Provider)
			if !ok {
				http.Redirect(w, req, "/provider", http.StatusFound)
				return
			}
			if providerH != "" && providerH != provider.Name() {
				w.WriteHeader(http.StatusUnauthorized)
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
			h.log.Error(ErrInvalidK8SConfigNil)
			// this should not block the request
			next(w, req, prefObj, user, provider)
			return
		}

		// 1. get the status of controller deployments for each cluster and make sure that all the contexts have meshery controllers deployed
		ctrlHlpr := h.MesheryCtrlsHelper.UpdateCtxControllerHandlers(mk8sContexts).UpdateOperatorsStatusMap(h.config.OperatorTracker).DeployUndeployedOperators(h.config.OperatorTracker)
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
		h.K8sCompRegHelper.UpdateContexts(contexts).RegisterComponents(contexts, []models.K8sRegistrationFunction{RegisterK8sMeshModelComponents}, h.EventsBuffer, h.registryManager, true)
		go h.config.MeshModelSummaryChannel.Publish()

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
		// ensuring session is intact
		err := provider.GetSession(req)
		if err != nil {
			err1 := provider.Logout(w, req)
			if err1 != nil {
				logrus.Errorf("Error performing logout: %v", err1.Error())
				provider.HandleUnAuthenticated(w, req)
				return
			}
			logrus.Errorf("Error: unable to get session: %v", err)
			http.Error(w, "unable to get session", http.StatusUnauthorized)
			return
		}

		user, err := provider.GetUserDetails(req)
		// if user details are not available,
		// then logout current user session and redirect to login page
		if err != nil {
			err1 := provider.Logout(w, req)
			if err1 != nil {
				logrus.Errorf("Error performing logout: %v", err1.Error())
				provider.HandleUnAuthenticated(w, req)
				return
			}
			h.log.Error(ErrGetUserDetails(err))
			http.Error(w, "unable to get user details", http.StatusUnauthorized)
			return
		}
		prefObj, err := provider.ReadFromPersister(user.UserID)
		if err != nil {
			logrus.Warn("unable to read session from the session persister, starting with a new one")
		}

		token := provider.UpdateToken(w, req)
		ctx := context.WithValue(req.Context(), models.TokenCtxKey, token)
		ctx = context.WithValue(ctx, models.PerfObjCtxKey, prefObj)
		ctx = context.WithValue(ctx, models.UserCtxKey, user)
		ctx = context.WithValue(ctx, models.RegistryManagerKey, h.registryManager)

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
