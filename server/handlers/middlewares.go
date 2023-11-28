// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines/kubernetes"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/machines"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const providerQParamName = "provider"

type ConnectionToStateMachineInstanceTracker struct {
	ConnectToInstanceMap map[uuid.UUID]*machines.StateMachine
	mx                   sync.RWMutex
}

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

// SHOULD GO AWAY AS THIS LOGIC WILL BE HANDLED IN THE CONNECTED STATE
func (h *Handler) MesheryControllersMiddleware(next func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)) func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider) {
	return func(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
		ctx := req.Context()
		// ctx, err := MesheryControllersMiddleware(ctx, h)
		// if err != nil {
		// 	h.log.Error(err)
		// 	next(w, req, prefObj, user, provider)
		// 	return
		// }
		req1 := req.WithContext(ctx)
		next(w, req1, prefObj, user, provider)
	}
}

// KubernetesMiddleware is a middleware that is responsible for handling kubernetes related stuff such as
// setting contexts, component generation etc.
func (h *Handler) KubernetesMiddleware(next func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)) func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider) {
	return func(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
		ctx := req.Context()
		ctx, err := KubernetesMiddleware(ctx, h, provider, user, req.URL.Query()["contexts"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
		ctx = context.WithValue(ctx, models.HandlerKey, h)
		ctx = context.WithValue(ctx, models.SystemIDKey, h.SystemID)
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

func KubernetesMiddleware(ctx context.Context, h *Handler, provider models.Provider, user *models.User, k8sContextIDs []string) (context.Context, error) {
	token, ok := ctx.Value(models.TokenCtxKey).(string)
	if !ok {
		err := ErrRetrieveUserToken(fmt.Errorf("failed to retrieve user token"))
		logrus.Error(err)
		return nil, err
	}
	userUUID := uuid.FromStringOrNil(user.ID)
	smInstanceTracker := h.ConnectionToStateMachineInstanceTracker
	connectedK8sContexts, err := provider.LoadAllK8sContext(token)

	k8sContextPassedByUser := []models.K8sContext{}
	k8sContextsFromKubeConfig := []*models.K8sContext{}

	if err != nil || len(connectedK8sContexts) == 0 {
		logrus.Warn("failed to get kubernetes contexts")
		k8sContextsFromKubeConfig, err = h.DiscoverK8SContextFromKubeConfig(user.ID, token, provider)
		if err != nil {
			logrus.Warn("failed to load kubernetes contexts: ", err.Error())
		}
	}

	if len(k8sContextIDs) == 1 && k8sContextIDs[0] == "all" {
		for _, c := range connectedK8sContexts {
			if c != nil {
				k8sContextPassedByUser = append(k8sContextPassedByUser, *c)
			}
		}
	} else {
		for _, kctxID := range k8sContextIDs {
			for _, c := range connectedK8sContexts {
				if c != nil && c.ID == kctxID {
					k8sContextPassedByUser = append(k8sContextPassedByUser, *c)
				}
			}
		}
	}

	ctx = context.WithValue(ctx, models.KubeClustersKey, k8sContextPassedByUser)
	ctx = context.WithValue(ctx, models.AllKubeClusterKey, connectedK8sContexts)

	for _, k8sContext := range connectedK8sContexts {
		machineCtx := &kubernetes.MachineCtx{
			K8sContext:         *k8sContext,
			MesheryCtrlsHelper: h.MesheryCtrlsHelper,
			K8sCompRegHelper:   h.K8sCompRegHelper,
			OperatorTracker:    h.config.OperatorTracker,
			Provider:           provider,
			K8scontextChannel:  h.config.K8scontextChannel,
			EventBroadcaster:   h.config.EventBroadcaster,
			RegistryManager:    h.registryManager,
		}
		connectionUUID := uuid.FromStringOrNil(k8sContext.ConnectionID)
		smInstanceTracker.mx.Lock()
		inst, ok := smInstanceTracker.ConnectToInstanceMap[connectionUUID]
		if !ok {
			inst, err = InitializeMachineWithContext(
				machineCtx,
				ctx,
				connectionUUID,
				smInstanceTracker,
				h.log,
				provider,
				machines.DefaultState,
				"kubernetes",
				kubernetes.AssignInitialCtx,
			)
			if err != nil {
				h.log.Error(err)
			}
		}
		inst.ResetState()
		go func(inst *machines.StateMachine) {
			event, err := inst.SendEvent(ctx, machines.Discovery, nil)
			if err != nil {
				_ = provider.PersistEvent(event)
				go h.config.EventBroadcaster.Publish(userUUID, event)
			}
		}(inst)
		smInstanceTracker.mx.Unlock()
	}

	for _, k8sContext := range k8sContextsFromKubeConfig {
		machineCtx := &kubernetes.MachineCtx{
			K8sContext:         *k8sContext,
			MesheryCtrlsHelper: h.MesheryCtrlsHelper,
			K8sCompRegHelper:   h.K8sCompRegHelper,
			OperatorTracker:    h.config.OperatorTracker,
			Provider:           provider,
			K8scontextChannel:  h.config.K8scontextChannel,
			EventBroadcaster:   h.config.EventBroadcaster,
			RegistryManager:    h.registryManager,
		}
		connectionUUID := uuid.FromStringOrNil(k8sContext.ConnectionID)
		smInstanceTracker.mx.Lock()
		inst, ok := smInstanceTracker.ConnectToInstanceMap[connectionUUID]
		if !ok {
			inst, err = InitializeMachineWithContext(
				machineCtx,
				ctx,
				connectionUUID,
				smInstanceTracker,
				h.log,
				provider,
				machines.DefaultState,
				"kubernetes",
				kubernetes.AssignInitialCtx,
			)
			if err != nil {
				h.log.Error(err)
			}
		}
		inst.ResetState()
		go func(inst *machines.StateMachine) {
			event, err := inst.SendEvent(ctx, machines.Discovery, nil)
			if err != nil {
				_ = provider.PersistEvent(event)
				go h.config.EventBroadcaster.Publish(userUUID, event)
			}
		}(inst)
		smInstanceTracker.mx.Unlock()
	}
	return ctx, nil
}
