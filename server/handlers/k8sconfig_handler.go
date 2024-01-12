// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/layer5io/meshery/server/machines"
	mhelpers "github.com/layer5io/meshery/server/machines/helpers"
	"github.com/layer5io/meshery/server/machines/kubernetes"

	"github.com/layer5io/meshery/server/models/connections"
	mcore "github.com/layer5io/meshery/server/models/meshmodel/core"

	// for GKE kube API authentication
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/layer5io/meshkit/models/events"

	"github.com/layer5io/meshkit/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SaveK8sContextResponse - struct used as (json marshaled) response to requests for saving k8s contexts
type SaveK8sContextResponse struct {
	RegisteredContexts []models.K8sContext `json:"registered_contexts"`
	ConnectedContexts  []models.K8sContext `json:"connected_contexts"`
	IgnoredContexts    []models.K8sContext `json:"ignored_contexts"`
	ErroredContexts    []models.K8sContext `json:"errored_contexts"`
}

// K8SConfigHandler is used for persisting kubernetes config and context info
func (h *Handler) K8SConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodPost && req.Method != http.MethodDelete {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	if req.Method == http.MethodPost {
		h.addK8SConfig(user, prefObj, w, req, provider)
		return
	}
	if req.Method == http.MethodDelete {
		h.deleteK8SConfig(user, prefObj, w, req, provider)
		return
	}
}

// swagger:route POST /api/system/kubernetes SystemAPI idPostK8SConfig
// Handle POST request for Kubernetes Config
//
// Used to add kubernetes config to System
// responses:
// 	200: k8sConfigRespWrapper

// The function is called only when user uploads a kube config.
// Connections which have state as "registered" are the only new ones, hence the GraphQL K8sContext subscription only sends an update to UI if any connection has registered state.
// A registered connection might have been regsitered previously and is not required for K8sContext Subscription to notify, but this case is not considered here.
func (h *Handler) addK8SConfig(user *models.User, _ *models.Preference, w http.ResponseWriter, req *http.Request, provider models.Provider) {
	userID := uuid.FromStringOrNil(user.ID)

	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		err := ErrRetrieveUserToken(fmt.Errorf("failed to retrieve user token"))
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	k8sConfigBytes, err := readK8sConfigFromBody(req)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Flatten kubeconfig. If that fails, go ahead with non-flattened config file
	flattenedK8sConfig, err := helpers.FlattenMinifyKubeConfig(*k8sConfigBytes)
	if err == nil {
		k8sConfigBytes = &flattenedK8sConfig
	}

	saveK8sContextResponse := SaveK8sContextResponse{
		RegisteredContexts: make([]models.K8sContext, 0),
		ConnectedContexts:  make([]models.K8sContext, 0),
		IgnoredContexts:    make([]models.K8sContext, 0),
		ErroredContexts:    make([]models.K8sContext, 0),
	}

	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("create").
		WithDescription("Kubernetes config uploaded.").WithSeverity(events.Informational)
	eventMetadata := map[string]interface{}{}
	contexts := models.K8sContextsFromKubeconfig(provider, user.ID, h.config.EventBroadcaster, *k8sConfigBytes, h.SystemID, eventMetadata)
	len := len(contexts)

	smInstanceTracker := h.ConnectionToStateMachineInstanceTracker
	for idx, ctx := range contexts {
		metadata := map[string]interface{}{}
		metadata["context"] = models.RedactCredentialsForContext(ctx)
		metadata["description"] = fmt.Sprintf("Connection established with context \"%s\" at %s", ctx.Name, ctx.Server)

		connection, err := provider.SaveK8sContext(token, *ctx)
		if err != nil {
			saveK8sContextResponse.ErroredContexts = append(saveK8sContextResponse.ErroredContexts, *ctx)
			metadata["description"] = fmt.Sprintf("Unable to establish connection with context \"%s\" at %s", ctx.Name, ctx.Server)
			metadata["error"] = err
		} else {
			ctx.ConnectionID = connection.ID.String()
			eventBuilder.ActedUpon(connection.ID)
			status := connection.Status
			machineCtx := &kubernetes.MachineCtx{
				K8sContext:         *ctx,
				MesheryCtrlsHelper: h.MesheryCtrlsHelper,
				K8sCompRegHelper:   h.K8sCompRegHelper,
				OperatorTracker:    h.config.OperatorTracker,
				K8scontextChannel:  h.config.K8scontextChannel,
				EventBroadcaster:   h.config.EventBroadcaster,
				RegistryManager:    h.registryManager,
			}

			if status == connections.CONNECTED {
				saveK8sContextResponse.ConnectedContexts = append(saveK8sContextResponse.ConnectedContexts, *ctx)
				metadata["description"] = fmt.Sprintf("Connection already exists with Kubernetes context \"%s\" at %s", ctx.Name, ctx.Server)
			} else if status == connections.IGNORED {
				saveK8sContextResponse.IgnoredContexts = append(saveK8sContextResponse.IgnoredContexts, *ctx)
				metadata["description"] = fmt.Sprintf("Kubernetes context \"%s\" is set to ignored state.", ctx.Name)
			} else if status == connections.DISCOVERED {
				saveK8sContextResponse.RegisteredContexts = append(saveK8sContextResponse.RegisteredContexts, *ctx)
				metadata["description"] = fmt.Sprintf("Connection registered with kubernetes context \"%s\" at %s.", ctx.Name, ctx.Server)
			}

			inst, err := mhelpers.InitializeMachineWithContext(
				machineCtx,
				req.Context(),
				connection.ID,
				userID,
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

			go func(inst *machines.StateMachine) {
				event, err := inst.SendEvent(req.Context(), machines.EventType(mhelpers.StatusToEvent(status)), nil)
				if err != nil {
					_ = provider.PersistEvent(event)
					go h.config.EventBroadcaster.Publish(userID, event)
				}
			}(inst)
		}

		eventMetadata[ctx.Name] = metadata

		if idx == len-1 {
			h.config.K8scontextChannel.PublishContext()
		}
	}

	event := eventBuilder.WithMetadata(eventMetadata).Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)

	if err := json.NewEncoder(w).Encode(saveK8sContextResponse); err != nil {
		logrus.Error(models.ErrMarshal(err, "kubeconfig"))
		http.Error(w, models.ErrMarshal(err, "kubeconfig").Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route DELETE /api/system/kubernetes SystemAPI idDeleteK8SConfig
// Handle DELETE request for Kubernetes Config
//
// Used to delete kubernetes config to System
// responses:
// 	200:

func (h *Handler) deleteK8SConfig(_ *models.User, _ *models.Preference, w http.ResponseWriter, _ *http.Request, _ models.Provider) {
	// prefObj.K8SConfig = nil
	// err := provider.RecordPreferences(req, user.UserID, prefObj)
	// if err != nil {
	// 	logrus.Error(ErrRecordPreferences(err))
	// 	http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
	// 	return
	// }

	ctxID := "0" //To be replaced with actual context ID after multi context support
	go core.DeleteK8sWorkloads(ctxID)
	_, _ = w.Write([]byte("{}"))
}

// swagger:route POST /api/system/kubernetes/contexts SystemAPI idPostK8SContexts
// Handle POST requests for Kubernetes Context list
//
// Returns the context list for a given k8s config
// responses:
// 	200: k8sContextsRespWrapper

// GetContextsFromK8SConfig returns the context list for a given k8s config
func (h *Handler) GetContextsFromK8SConfig(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {

	k8sConfigBytes, err := readK8sConfigFromBody(req)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userUUID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().FromUser(userUUID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("discovered").
		WithDescription("Kubernetes config uploaded.").WithSeverity(events.Informational)

	eventMetadata := map[string]interface{}{}

	contexts := models.K8sContextsFromKubeconfig(provider, user.ID, h.config.EventBroadcaster, *k8sConfigBytes, h.SystemID, eventMetadata)

	event := eventBuilder.WithMetadata(eventMetadata).Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userUUID, event)

	err = json.NewEncoder(w).Encode(contexts)
	if err != nil {
		logrus.Error(models.ErrMarshal(err, "kube-context"))
		http.Error(w, models.ErrMarshal(err, "kube-context").Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/system/kubernetes/ping?connection_id={id} SystemAPI idGetKubernetesPing
// Handle GET request for Kubernetes ping
//
// Fetches server version to simulate ping
// responses:
// 	200:

// KubernetesPingHandler - fetches server version to simulate ping
func (h *Handler) KubernetesPingHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to get the token for the user")
		return
	}

	connectionID := req.URL.Query().Get("connection_id")
	if connectionID != "" {
		// Get the context associated with this ID
		k8sContext, err := provider.GetK8sContext(token, connectionID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "failed to get kubernetes context for the given ID")
			return
		}

		// Create handler for the context
		kubeclient, err := k8sContext.GenerateKubeHandler()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "failed to get kubernetes config for the user")
			return
		}
		version, err := kubeclient.KubeClient.ServerVersion()
		if err != nil {
			logrus.Error(ErrKubeVersion(err))
			http.Error(w, ErrKubeVersion(err).Error(), http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(map[string]string{
			"server_version": version.String(),
		}); err != nil {
			err = errors.Wrap(err, "unable to marshal the payload")
			logrus.Error(models.ErrMarshal(err, "kube-server-version"))
			http.Error(w, models.ErrMarshal(err, "kube-server-version").Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Empty contextID. Pass the context ID(in query parameter \"context\") of the kuberenetes to be pinged", http.StatusBadRequest)
}

// swagger:route POST /api/system/kubernetes/register SystemAPI idPostK8SRegistration
// Handle registration request for Kubernetes components
//
// Used to register Kubernetes components to Meshery from a kubeconfig file
// responses:
//
//		202:
//	 400:
//	 500:
func (h *Handler) K8sRegistrationHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	k8sConfigBytes, err := readK8sConfigFromBody(req)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contexts := models.K8sContextsFromKubeconfig(provider, user.ID, h.config.EventBroadcaster, *k8sConfigBytes, h.SystemID, map[string]interface{}{}) // here we are not concerned for the events becuase inside the middleware the contexts would have been verified.
	h.K8sCompRegHelper.UpdateContexts(contexts).RegisterComponents(contexts, []models.K8sRegistrationFunction{mcore.RegisterK8sMeshModelComponents}, h.registryManager, h.config.EventBroadcaster, provider, user.ID, false)
	if _, err = w.Write([]byte(http.StatusText(http.StatusAccepted))); err != nil {
		logrus.Error(ErrWriteResponse)
		logrus.Error(err)
		http.Error(w, ErrWriteResponse.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) DiscoverK8SContextFromKubeConfig(userID string, token string, prov models.Provider) ([]*models.K8sContext, error) {
	var contexts []*models.K8sContext
	// userUUID := uuid.FromStringOrNil(userID)

	// Get meshery instance ID
	mid, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
	if !ok {
		return contexts, models.ErrMesheryInstanceID
	}

	// Attempt to get kubeconfig from the filesystem
	if h.config == nil {
		return contexts, ErrInvalidK8SConfigNil
	}
	kubeconfigSource := fmt.Sprintf("file://%s", filepath.Join(h.config.KubeConfigFolder, "config"))
	data, err := utils.ReadFileSource(kubeconfigSource)

	eventMetadata := map[string]interface{}{}
	metadata := map[string]interface{}{}
	if err != nil {
		// Could be an in-cluster deployment
		ctxName := "in-cluster"

		cc, err := models.NewK8sContextFromInClusterConfig(ctxName, mid)
		if err != nil {
			metadata["description"] = "Failed to import in-cluster kubeconfig."
			metadata["error"] = err
			logrus.Warn("failed to generate in cluster context: ", err)
			return contexts, err
		}
		if cc == nil {
			metadata["description"] = "No contexts detected in the in-cluster kubeconfig."
			err := fmt.Errorf("nil context generated from in cluster config")
			logrus.Warn(err)
			return contexts, err
		}
		cc.DeploymentType = "in_cluster"
		conn, err := prov.SaveK8sContext(token, *cc)
		if err != nil {
			metadata["description"] = fmt.Sprintf("Unable to establish connection with context \"%s\" at %s", cc.Name, cc.Server)
			metadata["error"] = err
			logrus.Warn("failed to save the context for incluster: ", err)
			return contexts, err
		}
		h.log.Debug(conn)

		cc.ConnectionID = conn.ID.String()
		contexts = append(contexts, cc)
		metadata["context"] = models.RedactCredentialsForContext(cc)
		eventMetadata["in-cluster"] = metadata
		return contexts, nil
	}

	cfg, err := helpers.FlattenMinifyKubeConfig([]byte(data))
	if err != nil {
		return contexts, err
	}

	ctxs := models.K8sContextsFromKubeconfig(prov, userID, h.config.EventBroadcaster, cfg, mid, eventMetadata)

	// Do not persist the generated contexts
	// consolidate this func and addK8sConfig. In this we explicitly updated status as well as this func perfomr greeedy upload so while consolidating make sure to handle the case.
	for _, ctx := range ctxs {
		metadata := map[string]interface{}{}
		metadata["context"] = models.RedactCredentialsForContext(ctx)
		metadata["description"] = fmt.Sprintf("K8S context \"%s\" discovered with cluster at %s", ctx.Name, ctx.Server)
		metadata["description"] = fmt.Sprintf("Connection established with context \"%s\" at %s", ctx.Name, ctx.Server)
		ctx.DeploymentType = "out_of_cluster"
		conn, err := prov.SaveK8sContext(token, *ctx)
		if err != nil {
			logrus.Warn("failed to save the context: ", err)
			metadata["description"] = fmt.Sprintf("Unable to establish connection with context \"%s\" at %s", ctx.Name, ctx.Server)
			metadata["error"] = err
			continue
		}
		ctx.ConnectionID = conn.ID.String()
		h.log.Debug(conn)

		contexts = append(contexts, ctx)
	}
	return contexts, nil
}

func readK8sConfigFromBody(req *http.Request) (*[]byte, error) {
	_ = req.ParseMultipartForm(1 << 20)

	k8sfile, _, err := req.FormFile("k8sfile")
	if err != nil {
		return nil, ErrFormFile(err)
	}
	defer func() {
		_ = k8sfile.Close()
	}()

	k8sConfigBytes, err := io.ReadAll(k8sfile)
	if err != nil {
		return nil, ErrReadConfig(err)
	}
	return &k8sConfigBytes, nil
}
