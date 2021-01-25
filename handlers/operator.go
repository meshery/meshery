package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	operatorv1alpha1 "github.com/layer5io/meshery-operator/api/v1alpha1"
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
	namespace = "meshery"

	operatorYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/manifests/default.yaml"
	brokerYaml   = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_broker.yaml"
	meshsyncYaml = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_meshsync.yaml"
)

var (
	status = Status{
		OperatorInstalled:   "false",
		BrokerInstalled:     "false",
		MeshsyncInstalled:   "false",
		SubscriptionStarted: "false",
		EventHandlerStarted: "false",
		Error:               "none",
	}
)

type Status struct {
	OperatorInstalled   string `json:"operator-installed,omitempty"`
	BrokerInstalled     string `json:"broker-installed,omitempty"`
	MeshsyncInstalled   string `json:"meshsync-installed,omitempty"`
	SubscriptionStarted string `json:"subscription-started,omitempty"`
	EventHandlerStarted string `json:"event-handler-started,omitempty"`
	Error               string `json:"error,omitempty"`
}

func (h *Handler) MeshSyncDataHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	objects, err := provider.ReadMeshSyncData()
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := utils.Marshal(objects)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = w.Write([]byte(data))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *Handler) OperatorStatusHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	s, err := utils.Marshal(status)
	if err != nil {
		return
	}

	_, err = w.Write([]byte(s))
	if err != nil {
		return
	}
}

func (h *Handler) OperatorHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	datach := make(chan *broker.Message)

	client, err := helpers.NewKubeClient(prefObj.K8SConfig.Config)
	if err != nil {
		fmt.Println(err)
		return
	}
	h.config.KubeClient = client

	enable, err := strconv.ParseBool(req.FormValue("enable"))
	if err != nil {
		fmt.Println(err)
		return
	}

	go h.handleEvents(datach, provider)

	go func(e bool, d chan *broker.Message) {
		er := h.initialize(e)
		if er != nil {
			status.Error = er.Error()
			return
		}

		er = h.subscribeToBroker(d)
		if er != nil {
			status.Error = er.Error()
			return
		}

		er = h.runMeshSync(e)
		if er != nil {
			status.Error = er.Error()
			return
		}
	}(!enable, datach)

	_, err = w.Write([]byte(`{"response": "ok"}`))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (h *Handler) initialize(delete bool) error {
	// installOperator
	err := h.applyYaml(delete, operatorYaml)
	if err != nil {
		return err
	}
	status.OperatorInstalled = "true"

	// installBroker
	err = h.applyYaml(delete, brokerYaml)
	if err != nil {
		return err
	}
	status.BrokerInstalled = "true"

	return nil
}

func (h *Handler) subscribeToBroker(datach chan *broker.Message) error {
	var broker *operatorv1alpha1.Broker
	mesheryclient, err := client.New(&h.config.KubeClient.RestConfig)
	if err != nil {
		return err
	}

	for {
		broker, err = mesheryclient.CoreV1Alpha1().Brokers(namespace).Get(context.Background(), "meshery-broker", metav1.GetOptions{})
		if err == nil && broker.Status.Endpoint.External != "" {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// subscribing to nats
	natsClient, err := nats.New(broker.Status.Endpoint.External)
	if err != nil {
		return err
	}

	err = natsClient.SubscribeWithChannel("meshsync", namespace, datach)
	if err != nil {
		return err
	}

	status.SubscriptionStarted = "true"
	return nil
}

func (h *Handler) runMeshSync(delete bool) error {
	// installMeshSync
	err := h.applyYaml(delete, meshsyncYaml)
	if err != nil {
		return err
	}
	status.MeshsyncInstalled = "true"
	return nil
}

func (h *Handler) handleEvents(datach chan *broker.Message, provider models.Provider) {
	status.EventHandlerStarted = "true"
	for {
		select {
		case msg := <-datach:
			objectJSON, _ := utils.Marshal(msg.Object)
			var object model.Object
			err := utils.Unmarshal(string(objectJSON), &object)
			if err != nil {
				status.Error = err.Error()
			}
			// persist the object
			err = provider.RecordMeshSyncData(object)
			if err != nil {
				status.Error = err.Error()
			}
		}
	}
}

func (h *Handler) applyYaml(delete bool, file string) error {
	contents, err := utils.ReadRemoteFile(file)
	if err != nil {
		return err
	}

	err = h.config.KubeClient.ApplyManifest([]byte(contents), mesherykube.ApplyOptions{
		Namespace: namespace,
		Update:    true,
		Delete:    delete,
	})
	if err != nil {
		return err
	}

	return nil
}
