// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	mhelpers "github.com/meshery/meshery/server/machines/helpers"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/meshsync/pkg/model"
	"github.com/meshery/schemas/models/core"
	"github.com/spf13/viper"
)

// isTransientProviderError returns true if the error indicates the remote
// provider is temporarily unreachable (network error, timeout, 5xx) as
// opposed to a genuine authentication failure (missing token, invalid token).
// This distinction prevents destroying a user's valid session when the
// remote provider has a transient outage.
func isTransientProviderError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "Could not reach remote provider") ||
		strings.Contains(errMsg, models.ErrUnreachableRemoteProviderCode)
}

const providerQParamName = "provider"

// backgroundDiscoveryContext keeps request values while allowing discovery to
// outlive the HTTP response, with a deadline for context-aware operations.
func backgroundDiscoveryContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(ctx), 2*time.Minute)
}

// legacyProviderAliasWarnOnce ensures the "None"->"Local" deprecation warning
// fires at most once per process lifetime, regardless of how many requests
// arrive with the legacy alias.
var legacyProviderAliasWarnOnce sync.Once

// resolveProviderName returns the provider name for a request based on, in
// order: the provider cookie, the provider header, the ?provider= query
// param, and finally the enforced default. Returns "" only in
// multi-provider mode when the client supplied no hint. The returned name is
// passed through models.NormalizeProviderName so the legacy "None" alias
// resolves to "Local"; the second return value indicates whether the raw
// input was a casing of the legacy alias so the caller can log a one-shot
// deprecation warning.
//
// When PROVIDER is unset, cookie then header then ?provider= select the
// provider. When PROVIDER is set, those client hints are ignored so a stale
// Local cookie cannot override a pinned remote.
func resolveProviderName(req *http.Request, cookieName, enforcedProvider string) (string, bool) {
	// When PROVIDER is set, it is the only provider this process will use.
	// Cookie, header, and ?provider= cannot override it; those were the
	// paths that let a client keep using Local on a pinned-remote deployment.
	if enforced := models.NormalizeProviderName(enforcedProvider); enforced != "" {
		usedLegacyAlias := strings.EqualFold(enforcedProvider, models.LocalProviderLegacyAlias)
		return enforced, usedLegacyAlias
	}
	var name string
	if ck, err := req.Cookie(cookieName); err == nil && ck.Value != "" {
		name = ck.Value
	} else if hdr := req.Header.Get(cookieName); hdr != "" {
		name = hdr
	} else if q := req.URL.Query().Get(providerQParamName); q != "" {
		name = q
	}
	usedLegacyAlias := strings.EqualFold(name, models.LocalProviderLegacyAlias)
	return models.NormalizeProviderName(name), usedLegacyAlias
}

// ProviderMiddleware is a middleware to validate if a provider is set
func (h *Handler) ProviderMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		var provider models.Provider
		providerName, usedLegacyAlias := resolveProviderName(req, h.config.ProviderCookieName, h.Provider)
		if usedLegacyAlias {
			legacyProviderAliasWarnOnce.Do(func() {
				h.log.Warnf("DEPRECATION: the local provider name %q is the legacy alias for %q and will be removed in a future Meshery release. Update meshery-provider cookies, the PROVIDER environment variable, and ~/.meshery/config.yaml to use %q. This warning is logged once per server process.", models.LocalProviderLegacyAlias, models.LocalProviderName, models.LocalProviderName)
			})
		}
		if providerName != "" {
			if providerKey, ok := models.ResolveProviderKey(providerName, h.config.Providers); ok {
				provider = h.config.Providers[providerKey]
			}
			if provider == nil && h.Provider != "" && providerName == models.NormalizeProviderName(h.Provider) {
				h.log.Errorf("enforced provider %q is not registered in h.config.Providers; ProviderUIHandler will degrade to serving the provider-selection UI instead of auto-login. Register %q in PROVIDERS or unset PROVIDER on this deployment.", h.Provider, h.Provider)
			}
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
			h.log.Error(ErrParsingCallBackUrl(err))
		} else {
			ctx = context.WithValue(ctx, models.MesheryServerURL, fmt.Sprintf("%s://%s", _url.Scheme, _url.Host))
		}
		req1 := req.WithContext(ctx)
		next.ServeHTTP(w, req1)
	}
	return http.HandlerFunc(fn)
}

// NoCacheMiddleware sets response headers that stop clients and proxies from caching the response.
func (h *Handler) NoCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Set headers to disable caching
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		// Call the next handler
		next.ServeHTTP(w, req)
	})
}

// AuthMiddleware is a middleware to validate if a user is authenticated
func (h *Handler) AuthMiddleware(next http.Handler, auth models.AuthenticationMechanism) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		refURLB64 := GetRefURL(req)

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
			// Propagate existing request parameters, if present.
			queryParams := req.URL.Query()
			if refURLB64 != "" {
				queryParams["ref"] = []string{refURLB64}
			}
			providerI := req.Context().Value(models.ProviderCtxKey)
			// logrus.Debugf("models.ProviderCtxKey %s", models.ProviderCtxKey)
			provider, ok := providerI.(models.Provider)
			if !ok {
				http.Redirect(w, req, fmt.Sprintf("/provider?%s", queryParams.Encode()), http.StatusFound)
				return
			}

			if providerH != "" {
				enforcedKey, ok := models.ResolveProviderKey(providerH, h.config.Providers)
				if ok {
					gotKey, gotOK := models.ResolveProviderKey(provider.Name(), h.config.Providers)
					if !gotOK || gotKey != enforcedKey {
						// Log the mismatch: without it this 401 is indistinguishable
						// from an ordinary auth failure, and the enforced-provider
						// case is exactly the one an operator needs to recognise.
						h.log.Infof("[AUTH_FLOW] step=AuthMiddleware action=provider_mismatch path=%s enforced=%q got=%q", req.URL.Path, enforcedKey, provider.Name())
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
				}
			}
			// logrus.Debugf("provider %s", provider)
			isValid, err := h.validateAuth(provider, req)
			if !isValid {
				h.log.Info(fmt.Sprintf("[AUTH_FLOW] step=AuthMiddleware action=auth_failed path=%s error=%v", req.URL.Path, err))
				if !errors.Is(err, models.ErrEmptySession) && provider.GetProviderType() == models.RemoteProviderType {
					provider.HandleUnAuthenticated(w, req)
					return
				}

				// Local Provider
				h.LoginHandler(w, req, provider, true)
				return
			}
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// RegistrationAuthMiddleware requires provider authentication for component
// registration only when a provider is actually in play: a provider is
// enforced through PROVIDER, or the request carries a selector for a remote
// provider. Adapters register against a local deployment with no provider
// selector, and in local mode there is no session to validate, so forcing
// ProviderAuth there would redirect every selector-free registration to
// /provider and break adapter registration. Requests that select Local
// (including the legacy None alias) pass through for the same reason: local
// mode has no credentials to check. The traversal itself is still blocked by
// the single-path-component validation in RegisterMeshmodelComponents, which
// runs regardless of auth.
func (h *Handler) RegistrationAuthMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		name, _ := resolveProviderName(r, h.config.ProviderCookieName, h.Provider)
		if h.Provider != "" || (name != "" && name != models.LocalProviderName) {
			h.AuthMiddleware(next, models.ProviderAuth).ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// validateAuth reports whether the request carries a valid provider session, returning the
// session error when it does not.
func (h *Handler) validateAuth(provider models.Provider, req *http.Request) (bool, error) {
	err := provider.GetSession(req)
	if err == nil {
		// logrus.Debugf("session: %v", sess)
		return true, nil
	}
	// logrus.Errorf("session invalid, error: %v", err)
	return false, err
}

// KubernetesMiddleware is a middleware that is responsible for handling kubernetes related stuff such as
// setting contexts, component generation etc.
func (h *Handler) KubernetesMiddleware(next func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)) func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider) {
	return func(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
		ctx := req.Context()
		ctx, err := KubernetesMiddleware(ctx, h, provider, user, req.URL.Query()["contexts"])
		if err != nil {
			writeMeshkitError(w, err, http.StatusInternalServerError)
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
		// err := provider.GetSession(req)
		// if err != nil {
		// 	err1 := provider.Logout(w, req)
		// 	if err1 != nil {
		// 		logrus.Errorf("Error performing logout: %v", err1.Error())
		// 		provider.HandleUnAuthenticated(w, req)
		// 		return
		// 	}
		// 	logrus.Errorf("Error: unable to get session: %v", err)
		// 	http.Error(w, "unable to get session", http.StatusUnauthorized)
		// 	return
		// }

		user, err := provider.GetUserDetails(req)
		if err != nil {
			h.log.Error(ErrGetUserDetails(err))
			// INTENTIONAL log/wire divergence. We log ErrGetUserDetails so the
			// operational trail captures *what* failed (the get-user call), but
			// we surface ErrTransientProvider on the wire so the client can
			// distinguish between a transient failure (e.g., Cloud unreachable)
			// and a genuine auth failure. Conflating the two would either flood
			// logs with misleading transient classifications or hide the
			// auth-vs-network distinction from clients. Don't "fix" this to
			// match without reading PR #18919.
			//
			// Behavioral consequence: on a transient provider error we must NOT
			// destroy the user's session by logging them out — that would cause
			// a redirect loop when Cloud recovers. A missing/invalid token
			// cookie still falls through to the genuine auth-failure path below.
			if isTransientProviderError(err) {
				writeMeshkitError(w, ErrTransientProvider(err), http.StatusServiceUnavailable)
				return
			}
			// Genuine auth failure — logout and redirect to login
			err1 := provider.Logout(w, req)
			if err1 != nil {
				h.log.Error(models.ErrLogout(err1))
				provider.HandleUnAuthenticated(w, req)
				return
			}
			writeMeshkitError(w, ErrGetUserDetails(err), http.StatusUnauthorized)
			return
		}
		prefObj, err := provider.ReadFromPersister(user.ID.String())
		if err != nil {
			// log underlying error from persister along with high-level context
			h.log.Warn(fmt.Errorf("%w: userID=%s: %v", ErrReadSessionPersistor, user.ID.String(), err))
			prefObj = models.NewDefaultPreference()
		} else if prefObj == nil {
			// persister unexpectedly returned a nil preference without error
			h.log.Warn(fmt.Errorf("%w: persister returned nil preference without error for userID=%s", ErrReadSessionPersistor, user.ID.String()))
			prefObj = models.NewDefaultPreference()
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

// GraphqlMiddleware adapts next to the handler-chain signature and forwards the request to it, without injecting anything
func (h *Handler) GraphqlMiddleware(next http.Handler) func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider) {
	return func(w http.ResponseWriter, req *http.Request, pref *models.Preference, user *models.User, prov models.Provider) {
		next.ServeHTTP(w, req)
	}
}

// KubernetesMiddleware loads connected Kubernetes contexts for the authenticated user and
// initiates background discovery on associated state machine instances. It guards against
// nil user, provider, or token inputs to prevent runtime panics.
func KubernetesMiddleware(ctx context.Context, h *Handler, provider models.Provider, user *models.User, k8sContextIDs []string) (context.Context, error) {
	if user == nil {
		err := fmt.Errorf("user is required for KubernetesMiddleware")
		if h != nil && h.log != nil {
			h.log.Error(err)
		}
		return nil, err
	}
	if provider == nil {
		err := fmt.Errorf("provider is required for KubernetesMiddleware")
		if h != nil && h.log != nil {
			h.log.Error(err)
		}
		return nil, err
	}

	token, ok := ctx.Value(models.TokenCtxKey).(string)
	if !ok {
		err := ErrRetrieveUserToken(fmt.Errorf("failed to retrieve user token"))
		if h != nil && h.log != nil {
			h.log.Error(err)
		}
		return nil, err
	}
	userUUID := user.ID
	smInstanceTracker := h.ConnectionToStateMachineInstanceTracker
	connectedK8sContexts, err := provider.LoadAllK8sContext(token)

	k8sContextPassedByUser := []models.K8sContext{}
	k8sContextsFromKubeConfig := []*models.K8sContext{}

	// A failed read is not the same as having no contexts, so log its cause
	// distinctly. Both still fall back to kubeconfig discovery: this middleware
	// fronts session bootstrap routes such as /api/system/sync, and failing them
	// on a transient read error would take the UI down with it.
	if err != nil {
		h.log.Error(ErrGetK8sContexts(err))
	} else if len(connectedK8sContexts) == 0 {
		h.log.Warn(ErrFailToGetK8SContext)
	}
	if err != nil || len(connectedK8sContexts) == 0 {
		k8sContextsFromKubeConfig, err = h.DiscoverK8SContextFromKubeConfig(user.ID.String(), token, provider)
		if err != nil {
			h.log.Warn(ErrFailToLoadK8sContext(err))
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

	var operatorTracker *models.OperatorTracker
	var k8scontextChannel *models.K8scontextChan
	var eventBroadcaster *models.Broadcast
	if h.config != nil {
		operatorTracker = h.config.OperatorTracker
		k8scontextChannel = h.config.K8scontextChannel
		eventBroadcaster = h.config.EventBroadcaster
	}

	for _, k8sContext := range k8sContextsFromKubeConfig {
		if k8sContext == nil {
			continue
		}
		machineCtx := &kubernetes.MachineCtx{
			K8sContext:         *k8sContext,
			MesheryCtrlsHelper: h.MesheryCtrlsHelper,
			K8sCompRegHelper:   h.K8sCompRegHelper,
			OperatorTracker:    operatorTracker,
			K8scontextChannel:  k8scontextChannel,
			EventBroadcaster:   eventBroadcaster,
			RegistryManager:    h.registryManager,
		}
		connectionUUID := uuid.FromStringOrNil(k8sContext.ConnectionID)

		inst, err := mhelpers.InitializeMachineWithContext(
			machineCtx,
			ctx,
			connectionUUID,
			userUUID,
			smInstanceTracker,
			h.log,
			provider,
			machines.DefaultState,
			"kubernetes",
			kubernetes.AssignInitialCtx,
		)
		if err != nil && h.log != nil {
			h.log.Error(err)
		}

		// Skip a machine that failed to initialize rather than nil-dereferencing
		// on ResetState/SendEvent. The connection stays stuck until the tracker
		// entry is removed; see mhelpers.HasMachineContext for why.
		if !mhelpers.HasMachineContext(inst) {
			continue
		}
		inst.ResetState()
		go func(inst *machines.StateMachine) {
			// Discovery outlives the request but must not wait indefinitely on Kubernetes.
			discoveryCtx, cancel := backgroundDiscoveryContext(ctx)
			defer cancel()
			event, err := inst.SendEvent(discoveryCtx, machines.Discovery, nil)
			h.reportBackgroundDiscovery(provider, token, userUUID, event, err)
		}(inst)
	}
	return ctx, nil
}

// reportBackgroundDiscovery persists and broadcasts what a background
// Discovery of a connection has to tell the user: its failure, or the
// resolution of the connection's unreachable condition when it became
// reachable again. A repeat of a failure that left the status unchanged is a
// Recurrence: it is counted on the open condition without notifying anyone. A
// routine successful Discovery reports nothing.
func (h *Handler) reportBackgroundDiscovery(provider models.Provider, token string, userUUID core.Uuid, event *events.Event, err error) {
	if provider == nil || event == nil || (err == nil && event.ResolvedAt == nil && !machines.IsRecurrence(event, err)) {
		return
	}
	h.reportMachineOutcome(provider, token, userUUID, event, err)
}

// reportMachineOutcome reports what a connection's StateMachine.SendEvent
// returned: nothing for no event, a Recurrence counted on its open condition
// without a broadcast, and anything else persisted and broadcast. A Recurrence
// is sent only to a provider that correlates conditions
// (Provider.CorrelatesConditions); any other would record it as a new event,
// so it is dropped, as before condition correlation.
func (h *Handler) reportMachineOutcome(provider models.Provider, token string, userID core.Uuid, event *events.Event, err error) {
	if provider == nil || event == nil {
		return
	}
	if machines.IsRecurrence(event, err) {
		if !provider.CorrelatesConditions() {
			return
		}
		if _, perr := provider.PersistEvent(*event, token); perr != nil {
			if h.log != nil {
				h.log.Error(perr)
			}
		}
		return
	}
	if err == nil && event.ResolvedAt != nil && event.CorrelationKey != nil {
		// This is a server-generated recovery for shared connection state.
		// Client-submitted resolutions still go through the owner-scoped API.
		if localProvider, ok := provider.(*models.DefaultLocalProvider); ok {
			resolved, resolveErr := localProvider.ResolveConnectionCondition(*event)
			if resolveErr != nil {
				if h.log != nil {
					h.log.Error(resolveErr)
				}
				return
			}
			if h.config != nil && h.config.EventBroadcaster != nil {
				for i := range resolved {
					h.config.EventBroadcaster.Publish(*resolved[i].Owner, &resolved[i])
				}
			}
			return
		}
	}
	h.persistAndPublish(provider, token, userID, event)
}

// persistAndPublish persists event and broadcasts it as stored: an occurrence
// counted on, or the resolution of, an open condition is broadcast as the
// condition's stored event, so a client updates that notification in place
// instead of showing one whose ID was never recorded. When persisting fails the
// event is still broadcast, as sent.
//
// Every server site that persists an event and broadcasts it goes through this
// helper (or persistAndPublishEvent, outside a Handler; a system event goes
// through models.PersistAndPublishSystemEvent), never publishing its own copy
// after PersistEvent: that copy carries no stored ID.
func (h *Handler) persistAndPublish(provider models.Provider, token string, userID core.Uuid, event *events.Event) {
	if h == nil || provider == nil || event == nil {
		return
	}
	var broadcaster *models.Broadcast
	if h.config != nil {
		broadcaster = h.config.EventBroadcaster
	}
	persistAndPublishEvent(h.log, provider, token, broadcaster, userID, event)
}

// persistAndPublishEvent is persistAndPublish for code that has no Handler.
//
// A nil event reports nothing: SendEvent returns none for a background
// Discovery whose failure it suppressed.
func persistAndPublishEvent(log logger.Handler, provider models.Provider, token string, broadcaster *models.Broadcast, userID core.Uuid, event *events.Event) {
	if provider == nil || event == nil {
		return
	}
	// Prepare the event here, not only on the provider's copy, so the event
	// broadcast when persisting fails has its own ID too: events that share
	// the zero ID overwrite one another in the UI's store.
	prepared := models.PrepareEventForPersistence(*event, log)
	event = &prepared
	persisted, err := provider.PersistEvent(*event, token)
	if err != nil {
		if log != nil {
			log.Error(err)
		}
	} else {
		if persisted.Event == nil {
			// Nothing new was recorded and the stored event is unknown.
			return
		}
		event = persisted.Event
	}
	if broadcaster != nil {
		broadcaster.Publish(userID, event)
	}
}

// K8sFSMMiddleware wraps an HTTP handler function with Kubernetes finite state machine
// background discovery, safely passing context and identity down the middleware chain.
func (h *Handler) K8sFSMMiddleware(next func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider)) func(http.ResponseWriter, *http.Request, *models.Preference, *models.User, models.Provider) {
	return func(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
		K8sFSMMiddleware(req.Context(), h, provider, user)
		next(w, req, prefObj, user, provider)
	}
}

type dataHandlerToClusterID struct {
	mdh       models.MeshsyncDataHandler
	clusterID string
}

// K8sFSMMiddleware re-runs background discovery across connected Kubernetes contexts and registers meshsync handlers.
// It defensively extracts cluster contexts and guards against nil user or provider values to prevent runtime panics.
func K8sFSMMiddleware(ctx context.Context, h *Handler, provider models.Provider, user *models.User) {
	if h == nil {
		return
	}

	token, _ := ctx.Value(models.TokenCtxKey).(string)
	smInstanceTracker := h.ConnectionToStateMachineInstanceTracker

	var connectedK8sContexts []*models.K8sContext
	if val := ctx.Value(models.AllKubeClusterKey); val != nil {
		if contexts, ok := val.([]*models.K8sContext); ok {
			connectedK8sContexts = contexts
		}
	}

	var userUUID uuid.UUID
	if user != nil {
		userUUID = user.ID
	}

	var operatorTracker *models.OperatorTracker
	var k8scontextChannel *models.K8scontextChan
	var eventBroadcaster *models.Broadcast
	if h.config != nil {
		operatorTracker = h.config.OperatorTracker
		k8scontextChannel = h.config.K8scontextChannel
		eventBroadcaster = h.config.EventBroadcaster
	}

	dataHandlers := []*dataHandlerToClusterID{}
	clusterIDs := []string{}
	for _, k8sContext := range connectedK8sContexts {
		if k8sContext == nil {
			continue
		}
		machineCtx := &kubernetes.MachineCtx{
			K8sContext:         *k8sContext,
			MesheryCtrlsHelper: h.MesheryCtrlsHelper,
			K8sCompRegHelper:   h.K8sCompRegHelper,
			OperatorTracker:    operatorTracker,
			K8scontextChannel:  k8scontextChannel,
			EventBroadcaster:   eventBroadcaster,
			RegistryManager:    h.registryManager,
		}
		connectionUUID := uuid.FromStringOrNil(k8sContext.ConnectionID)

		inst, err := mhelpers.InitializeMachineWithContext(
			machineCtx,
			ctx,
			connectionUUID,
			userUUID,
			smInstanceTracker,
			h.log,
			provider,
			machines.DefaultState,
			"kubernetes",
			kubernetes.AssignInitialCtx,
		)
		if err != nil && h.log != nil {
			h.log.Error(err)
		}

		// Skip a machine that failed to initialize rather than nil-dereferencing
		// on ResetState/SendEvent, or type-asserting its nil Context on the Cast
		// below. See mhelpers.HasMachineContext for why.
		if !mhelpers.HasMachineContext(inst) {
			continue
		}
		inst.ResetState()
		go func(inst *machines.StateMachine) {
			discoveryCtx, cancel := backgroundDiscoveryContext(ctx)
			defer cancel()
			event, err := inst.SendEvent(discoveryCtx, machines.Discovery, nil)
			h.reportBackgroundDiscovery(provider, token, userUUID, event, err)
		}(inst)
		kubernesMachineCtx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
		if err != nil {
			if h.log != nil {
				h.log.Error(err)
			}
			continue
		}
		if kubernesMachineCtx != nil && kubernesMachineCtx.MesheryCtrlsHelper != nil {
			mdh := kubernesMachineCtx.MesheryCtrlsHelper.GetMeshSyncDataHandlersForEachContext()
			if mdh != nil {
				if k8sContext.KubernetesServerID == nil {
					continue
				}
				dataHandlers = append(dataHandlers, &dataHandlerToClusterID{
					mdh:       *mdh,
					clusterID: k8sContext.KubernetesServerID.String(),
				})
				clusterIDs = append(clusterIDs, k8sContext.KubernetesServerID.String())
			}
		}
	}

	if provider == nil || len(clusterIDs) == 0 {
		return
	}
	genericPersister := provider.GetGenericPersister()
	if genericPersister == nil {
		return
	}

	var resources []model.KubernetesResource

	err := genericPersister.Model(&model.KubernetesResource{}).
		Preload("KubernetesResourceMeta").
		Joins("JOIN kubernetes_resource_object_meta ON kubernetes_resource_object_meta.id = kubernetes_resources.id").
		Where("kubernetes_resources.cluster_id IN (?)", clusterIDs).Where(&model.KubernetesResource{Kind: "Service"}).Where("lower(kubernetes_resource_object_meta.name) LIKE ? OR lower(kubernetes_resource_object_meta.name) LIKE ?", "%grafana%", "%prometheus%").Find(&resources).Error

	if err != nil {
		if h.log != nil {
			h.log.Error(ErrFetchMeshSyncResources(err))
		}
		return
	}

	regQueue := models.GetMeshSyncRegistrationQueue()

	for _, resource := range resources {
		for _, dh := range dataHandlers {
			if dh.clusterID == resource.ClusterID {
				go regQueue.Send(models.MeshSyncRegistrationData{MeshsyncDataHandler: dh.mdh, Obj: resource})
			}
		}
	}
}
