package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/layer5io/meshery-operator/pkg/client"
	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshsync/pkg/broker"
	"github.com/layer5io/meshsync/pkg/broker/nats"
	"github.com/layer5io/meshsync/pkg/model"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	operatorYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/manifests/default.yaml"
	operatorCRDS = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/manifests/crd.yaml"
	brokerYaml   = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_broker.yaml"
	meshsyncYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_meshsync.yaml"
)

func (h *Handler) OperatorHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {

	client, err := helpers.NewKubeClient(prefObj.K8SConfig.Config)
	if err != nil {
		fmt.Println(err)
	}
	h.config.KubeClient = client

	enable, err := strconv.ParseBool(req.FormValue("enable"))
	if err != nil {
		fmt.Println(err)
	}

	err = h.initialize(!enable)
	if err != nil {
		fmt.Println(err)
	}

	ch := make(chan *broker.Message)
	go h.subscribeToMeshsync(ch)
	fmt.Println("subscriber started")
	go h.persistEvents(ch, provider)
	fmt.Println("persister started")

	_, err = w.Write([]byte("ok"))
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) initialize(delete bool) error {

	// installOperator
	err := h.applyYaml(delete, operatorYaml)
	handleError(err)
	fmt.Println("applied operator")

	// installCRDS
	err = h.applyYaml(delete, operatorCRDS)
	handleError(err)
	fmt.Println("applied crd")

	// installBroker
	err = h.applyYaml(delete, brokerYaml)
	handleError(err)
	fmt.Println("applied broker")

	// installMeshSync
	err = h.applyYaml(delete, meshsyncYaml)
	handleError(err)
	fmt.Println("applied meshsync")

	return nil
}

func (h *Handler) subscribeToMeshsync(msgch chan *broker.Message) {

	mesheryclient, err := client.New(&h.config.KubeClient.RestConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	broker, err := mesheryclient.CoreV1Alpha1().Brokers("meshery").Get(context.Background(), "broker", metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// subscribing to nats
	natsClient, err := nats.New(broker.Status.Endpoint)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = natsClient.SubscribeWithChannel("meshsync", "meshery", msgch)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (h *Handler) persistEvents(msgch chan *broker.Message, provider models.Provider) {
	// broker.Message
	for {
		select {
		case msg := <-msgch:
			objectJSON, _ := json.Marshal(msg.Object)
			var object model.Object
			err := utils.Unmarshal(string(objectJSON), &object)
			if err != nil {
				fmt.Println(err)
				return
			}
			// persist the object
			err = provider.RecordMeshSyncData(object)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func (h *Handler) applyYaml(delete bool, file string) error {
	contents, err := utils.ReadRemoteFile(file)
	handleError(err)

	err = h.config.KubeClient.ApplyManifest([]byte(contents), mesherykube.ApplyOptions{
		Namespace: "meshery",
		Update:    true,
		Delete:    delete,
	})
	handleError(err)

	return nil
}

func handleError(err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
