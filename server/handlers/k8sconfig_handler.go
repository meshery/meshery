// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	mutil "github.com/layer5io/meshery/server/helpers/utils"

	guid "github.com/google/uuid"
	"github.com/layer5io/meshery/server/meshes"
	mcore "github.com/layer5io/meshery/server/models/meshmodel/core"
	meshmodelv1alpha1 "github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"

	// for GKE kube API authentication
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	putils "github.com/layer5io/meshery/server/models/pattern/utils"
	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/events"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// SaveK8sContextResponse - struct used as (json marshaled) response to requests for saving k8s contexts
type SaveK8sContextResponse struct {
	InsertedContexts []models.K8sContext `json:"inserted_contexts"`
	UpdatedContexts  []models.K8sContext `json:"updated_contexts"`
	ErroredContexts  []models.K8sContext `json:"errored_contexts"`
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

func (h *Handler) addK8SConfig(_ *models.User, _ *models.Preference, w http.ResponseWriter, req *http.Request, provider models.Provider) {
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

	// Get meshery instance ID
	mid, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
	if !ok {
		logrus.Error(ErrMesheryInstanceID)
		http.Error(w, ErrMesheryInstanceID.Error(), http.StatusInternalServerError)
		return
	}

	saveK8sContextResponse := SaveK8sContextResponse{
		InsertedContexts: make([]models.K8sContext, 0),
		UpdatedContexts:  make([]models.K8sContext, 0),
		ErroredContexts:  make([]models.K8sContext, 0),
	}

	contexts, respMessage := models.K8sContextsFromKubeconfig(*k8sConfigBytes, mid)
	for _, ctx := range contexts {
		_, err := provider.SaveK8sContext(token, *ctx) // Ignore errors
		if err != nil {
			if err == models.ErrContextAlreadyPersisted {
				saveK8sContextResponse.UpdatedContexts = append(saveK8sContextResponse.UpdatedContexts, *ctx)
			} else {
				saveK8sContextResponse.ErroredContexts = append(saveK8sContextResponse.ErroredContexts, *ctx)
				logrus.Error("failed to persist context")
			}
		} else {
			saveK8sContextResponse.InsertedContexts = append(saveK8sContextResponse.InsertedContexts, *ctx)
		}
	}
	if len(saveK8sContextResponse.InsertedContexts) > 0 || len(saveK8sContextResponse.UpdatedContexts) > 0 {
		h.config.K8scontextChannel.PublishContext()
	}
	if err := json.NewEncoder(w).Encode(saveK8sContextResponse); err != nil {
		logrus.Error(ErrMarshal(err, "kubeconfig"))
		http.Error(w, ErrMarshal(err, "kubeconfig").Error(), http.StatusInternalServerError)
		return
	}

	if respMessage != "" {
		h.EventsBuffer.Publish(&meshes.EventsResponse{
			Component:     "core",
			ComponentName: "kubernetes",
			OperationId:   guid.NewString(),
			EventType:     meshes.EventType_INFO,
			Summary:       "Kubernetes configuration Info",
			Details:       respMessage,
		})
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
func (h *Handler) GetContextsFromK8SConfig(w http.ResponseWriter, req *http.Request) {
	
	k8sConfigBytes, err := readK8sConfigFromBody(req)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get meshery instance ID
	mid, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
	if !ok {
		logrus.Error(ErrMesheryInstanceID)
		http.Error(w, ErrMesheryInstanceID.Error(), http.StatusInternalServerError)
		return
	}

	contexts, _ := models.K8sContextsFromKubeconfig(*k8sConfigBytes, mid)

	err = json.NewEncoder(w).Encode(contexts)
	if err != nil {
		logrus.Error(ErrMarshal(err, "kube-context"))
		http.Error(w, ErrMarshal(err, "kube-context").Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/system/kubernetes/ping?contexts={id} SystemAPI idGetKubernetesPing
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

	ctx := req.URL.Query().Get("context")
	if ctx != "" {
		// Get the context associated with this ID
		k8sContext, err := provider.GetK8sContext(token, ctx)
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
			logrus.Error(ErrMarshal(err, "kube-server-version"))
			http.Error(w, ErrMarshal(err, "kube-server-version").Error(), http.StatusInternalServerError)
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
// 	202:
//  400:
//  500:
func (h *Handler) K8sRegistrationHandler(w http.ResponseWriter, req *http.Request) {
	k8sConfigBytes, err := readK8sConfigFromBody(req)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get meshery instance ID
	mid, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
	if !ok {
		logrus.Error(ErrMesheryInstanceID)
		http.Error(w, ErrMesheryInstanceID.Error(), http.StatusInternalServerError)
		return
	}

	contexts, _ := models.K8sContextsFromKubeconfig(*k8sConfigBytes, mid)
	h.K8sCompRegHelper.UpdateContexts(contexts).RegisterComponents(contexts, []models.K8sRegistrationFunction{RegisterK8sMeshModelComponents}, h.EventsBuffer, h.registryManager, false)
	if _, err = w.Write([]byte(http.StatusText(http.StatusAccepted))); err != nil {
		logrus.Error(ErrWriteResponse)
		logrus.Error(err)
		http.Error(w, ErrWriteResponse.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) LoadContextsAndPersist(token string, prov models.Provider) ([]*models.K8sContext, error) {
	var contexts []*models.K8sContext
	// Get meshery instance ID
	mid, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
	if !ok {
		return contexts, ErrMesheryInstanceID
	}

	// Attempt to get kubeconfig from the filesystem
	if h.config == nil {
		return contexts, ErrInvalidK8SConfig
	}
	data, err := utils.ReadFileSource(fmt.Sprintf("file://%s", filepath.Join(h.config.KubeConfigFolder, "config")))
	if err != nil {
		// Could be an in-cluster deployment
		ctxName := "in-cluster"

		cc, err := models.NewK8sContextFromInClusterConfig(ctxName, mid)
		if err != nil {
			logrus.Warn("failed to generate in cluster context: ", err)
			return contexts, err
		}
		if cc == nil {
			err := fmt.Errorf("nil context generated from in cluster config")
			logrus.Warn(err)
			return contexts, err
		}
		cc.DeploymentType = "in_cluster"
		_, err = prov.SaveK8sContext(token, *cc)
		if err != nil {
			logrus.Warn("failed to save the context for incluster: ", err)
			return contexts, err
		}
		h.config.K8scontextChannel.PublishContext()
		contexts = append(contexts, cc)
		return contexts, nil
	}

	cfg, err := helpers.FlattenMinifyKubeConfig([]byte(data))
	if err != nil {
		return contexts, err
	}

	ctxs, _ := models.K8sContextsFromKubeconfig(cfg, mid)

	// Persist the generated contexts
	for _, ctx := range ctxs {
		ctx.DeploymentType = "out_of_cluster"
		_, err := prov.SaveK8sContext(token, *ctx)
		if err != nil {
			logrus.Warn("failed to save the context: ", err)
			continue
		}
		h.config.K8scontextChannel.PublishContext()
		contexts = append(contexts, ctx)
	}
	return contexts, nil
}

func RegisterK8sMeshModelComponents(_ context.Context, config []byte, ctxID string, reg *meshmodel.RegistryManager, es *events.EventStreamer, ctxName string) (err error) {
	man, err := mcore.GetK8sMeshModelComponents(config)
	if err != nil {
		return ErrCreatingKubernetesComponents(err, ctxID)
	}
	if man == nil {
		return ErrCreatingKubernetesComponents(errors.New("generated components are nil"), ctxID)
	}
	count := 0
	for _, c := range man {
		writeK8sMetadata(&c, reg)
		err = reg.RegisterEntity(meshmodel.Host{
			Hostname:  "kubernetes",
			ContextID: ctxID,
		}, c)
		count++
	}
	es.Publish(&meshes.EventsResponse{
		Component:     "core",
		ComponentName: "kubernetes",
		OperationId:   guid.NewString(),
		EventType:     meshes.EventType_INFO,
		Summary:       fmt.Sprintf("%d Kubernetes components registered from %s", count, ctxName),
		Details:       fmt.Sprintf("%d MeshModel components registered for Kubernetes context \"%s\" (%s)", count, ctxName, ctxID),
	})
	return
}

const k8sMeshModelPath = "../meshmodel/components/kubernetes/model_template.json"

var k8sMeshModelMetadata = make(map[string]interface{})

func writeK8sMetadata(comp *meshmodelv1alpha1.ComponentDefinition, reg *meshmodel.RegistryManager) {
	ent, _ := reg.GetEntities(&meshmodelv1alpha1.ComponentFilter{
		Name:       comp.Kind,
		APIVersion: comp.APIVersion,
	})
	//If component was not available in the registry, then use the generic model level metadata
	if len(ent) == 0 {
		putils.MergeMaps(comp.Metadata, k8sMeshModelMetadata)
		mutil.WriteSVGsOnFileSystem(comp)
	} else {
		existingComp, ok := ent[0].(meshmodelv1alpha1.ComponentDefinition)
		if !ok {
			putils.MergeMaps(comp.Metadata, k8sMeshModelMetadata)
			return
		}
		putils.MergeMaps(comp.Metadata, existingComp.Metadata)
		comp.Model = existingComp.Model
	}
}

// Caches k8sMeshModel metadatas in memory to use at the time of dynamic k8s component generation
func init() {
	f, err := os.Open(filepath.Join(k8sMeshModelPath))
	if err != nil {
		return
	}
	byt, err := io.ReadAll(f)
	if err != nil {
		return
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(byt, &m)
	if err != nil {
		return
	}
	k8sMeshModelMetadata = m
}

func readK8sConfigFromBody(req *http.Request) (*[]byte, error){
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
// func writeDefK8sOnFileSystem(def string, path string) {
// 	err := ioutil.WriteFile(path, []byte(def), 0777)
// 	if err != nil {
// 		fmt.Println("err def: ", err.Error())
// 	}
// }

// func writeSchemaK8sFileSystem(schema string, path string) {
// 	err := ioutil.WriteFile(path, []byte(schema), 0777)
// 	if err != nil {
// 		fmt.Println("err schema: ", err.Error())
// 	}
// }
