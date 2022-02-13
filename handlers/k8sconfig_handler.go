//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	// for GKE kube API authentication
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/gofrs/uuid"

	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshery/models/pattern/core"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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

func (h *Handler) addK8SConfig(user *models.User, prefObj *models.Preference, w http.ResponseWriter, req *http.Request, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		err := ErrRetrieveUserToken(fmt.Errorf("failed to retrieve user token"))
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_ = req.ParseMultipartForm(1 << 20)

	inClusterConfig := req.FormValue("inClusterConfig")
	logrus.Debugf("inClusterConfig: %s", inClusterConfig)

	k8sfile, _, err := req.FormFile("k8sfile")
	if err != nil {
		logrus.Error(ErrFormFile(err))
		http.Error(w, ErrFormFile(err).Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		_ = k8sfile.Close()
	}()

	k8sConfigBytes, err := ioutil.ReadAll(k8sfile)
	if err != nil {
		logrus.Error(ErrReadConfig(err))
		http.Error(w, ErrReadConfig(err).Error(), http.StatusBadRequest)
		return
	}

	// Get meshery instance ID
	mid, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
	if !ok {
		logrus.Error(ErrMesheryInstanceID)
		http.Error(w, ErrMesheryInstanceID.Error(), http.StatusInternalServerError)
		return
	}

	contexts := models.K8sContextsFromKubeconfig(k8sConfigBytes, mid)

	for _, ctx := range contexts {
		_, err := provider.SaveK8sContext(token, ctx) // Ignore errors
		if err != nil {
			logrus.Error("failed to persist context")
		}
	}

	if err := json.NewEncoder(w).Encode(contexts); err != nil {
		logrus.Error(ErrMarshal(err, "kubeconfig"))
		http.Error(w, ErrMarshal(err, "kubeconfig").Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route DELETE /api/system/kubernetes SystemAPI idDeleteK8SConfig
// Handle DELETE request for Kubernetes Config
//
// Used to delete kubernetes config to System
// responses:
// 	200:

func (h *Handler) deleteK8SConfig(user *models.User, prefObj *models.Preference, w http.ResponseWriter, req *http.Request, provider models.Provider) {
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
	_ = req.ParseMultipartForm(1 << 20)
	var k8sConfigBytes []byte

	k8sfile, _, err := req.FormFile("k8sfile")
	if err != nil {
		logrus.Error(ErrFormFile(err))
		http.Error(w, ErrFormFile(err).Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		_ = k8sfile.Close()
	}()
	k8sConfigBytes, err = io.ReadAll(k8sfile)
	if err != nil {
		logrus.Error(ErrReadConfig(err))
		http.Error(w, ErrReadConfig(err).Error(), http.StatusBadRequest)
		return
	}

	// Get meshery instance ID
	mid, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
	if !ok {
		logrus.Error(ErrMesheryInstanceID)
		http.Error(w, ErrMesheryInstanceID.Error(), http.StatusInternalServerError)
		return
	}

	contexts := models.K8sContextsFromKubeconfig(k8sConfigBytes, mid)

	err = json.NewEncoder(w).Encode(contexts)
	if err != nil {
		logrus.Error(ErrMarshal(err, "kube-context"))
		http.Error(w, ErrMarshal(err, "kube-context").Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/system/kubernetes/ping SystemAPI idGetKubernetesPing
// Handle GET request for Kubernetes ping
//
// Fetches server version to simulate ping
// responses:
// 	200:

// KubernetesPingHandler - fetches server version to simulate ping
func (h *Handler) KubernetesPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to get the token for the user")
		return
	}

	kubeclient, ok := req.Context().Value(models.KubeHanderKey).(*mesherykube.Client)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to get kube client for the user")
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
		kubeclient, err = k8sContext.GenerateKubeHandler()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "failed to get kubernetes config for the user")
			return
		}
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
		return
	}
}

func (h *Handler) GetCurrentContext(token string, prov models.Provider) (*models.K8sContext, error) {
	// Try to get current context
	cc, err := prov.GetCurrentContext(token)
	if err != nil {
		// No current context implies that this is the first time meshery is loading up
		// or it could mean that meshery does not have access to the kubernetes cluster

		// Get meshery instance ID
		mid, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
		if !ok {
			return nil, ErrMesheryInstanceID
		}

		// Attempt to get kubeconfig from the filesystem
		data, err := utils.ReadFileSource(fmt.Sprintf("file://%s", filepath.Join(h.config.KubeConfigFolder, "config")))
		if err != nil {
			// Could be an in-cluster deployment
			ctxName := "in-cluster"

			cc, err := models.NewK8sContextFromInClusterConfig(ctxName, mid)
			if err != nil {
				logrus.Warn("failed to generate in cluster context: ", err)
				return nil, err
			}

			if _, err := prov.SaveK8sContext(token, *cc); err == nil {
				_, _ = prov.SetCurrentContext(token, cc.ID) // Ignore the error
			}

			// Generate Kube Config
			if !viper.GetBool("SKIP_COMP_GEN") {
				ctxID := cc.ID
				cfg, err := cc.GenerateKubeConfig()
				if err != nil {
					return cc, nil
				}

				go func(l logger.Handler, config []byte, ctx string) {
					err := registerK8sComponents(h.log, config, ctxID)
					if err != nil {
						logrus.Error(err)
					}
				}(h.log, cfg, ctxID)
			}

			return cc, nil
		}

		cfg, err := helpers.FlattenMinifyKubeConfig([]byte(data))
		if err != nil {
			return nil, err
		}

		ctxs := models.K8sContextsFromKubeconfig(cfg, mid)

		// Persist the generated contexts
		for _, ctx := range ctxs {
			_, err := prov.SaveK8sContext(token, ctx)
			if err != nil {
				logrus.Warn("failed to save the context: ", err)
				continue
			}

			if ctx.IsCurrentContext {
				_, _ = prov.SetCurrentContext(token, ctx.ID) // Ignore the error

				if !viper.GetBool("SKIP_COMP_GEN") {
					ctxID := cc.ID
					cfg, err := cc.GenerateKubeConfig()
					if err != nil {
						return &ctx, nil
					}

					go func(l logger.Handler, config []byte, ctx string) {
						err := registerK8sComponents(h.log, config, ctxID)
						if err != nil {
							logrus.Error(err)
						}
					}(h.log, cfg, ctxID)
				}

				return &ctx, nil
			}
		}
	}

	return &cc, nil
}
func registerK8sComponents(l logger.Handler, config []byte, ctx string) error {
	l.Info("Starting to register k8s native components")
	man, err := core.GetK8Components(config, ctx)
	if err != nil {
		return err
	}
	if man == nil {
		l.Error(errors.New("Could not get k8s components"))
		return err
	}
	for i, def := range man.Definitions {
		var ord core.WorkloadCapability
		ord.Metadata = make(map[string]string)
		ord.Metadata["io.meshery.ctxid"] = ctx
		ord.Metadata["adapter.meshery.io/name"] = "kubernetes"
		ord.Host = "<none-local>"
		ord.OAMRefSchema = man.Schemas[i]

		var definition v1alpha1.WorkloadDefinition
		err := json.Unmarshal([]byte(def), &definition)
		if err != nil {
			return err
		}
		ord.OAMDefinition = definition
		content, err := json.Marshal(ord)
		if err != nil {
			return err
		}
		err = core.RegisterWorkload(content)
		if err != nil {
			return err
		}
	}
	l.Info("Registration of k8s native components completed")
	return nil
}
