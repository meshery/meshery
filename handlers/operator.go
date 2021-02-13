package handlers

import (
	"context"
	"encoding/json"
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
	"github.com/sirupsen/logrus"

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

type OperatorStatus string

const (
	// LoadTestError - respresents an error status
	OperatorError OperatorStatus = "error"

	// LoadTestInfo - represents a info status
	OperatorInfo OperatorStatus = "info"

	// LoadTestSuccess - represents a success status
	OperatorSuccess OperatorStatus = "success"
)

type OperatorResponse struct {
	Status  OperatorStatus `json:"status,omitempty"`
	Message string         `json:"message,omitempty"`
	Result  *Status        `json:"result,omitempty"`
}

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
	defer func() {
		_ = req.Body.Close()
	}()

	log := logrus.WithField("file", "operator_handler")

	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Error("Event streaming not supported.")
		http.Error(w, "Event streaming is not supported at the moment.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	notify := req.Context()
	respChan := make(chan *OperatorResponse, 100)
	endChan := make(chan struct{})
	defer close(endChan)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("Recovered from panic: %v.", r)
			}
		}()
		for data := range respChan {
			bd, err := json.Marshal(data)
			if err != nil {
				logrus.Errorf("error: unable to marshal meshery result for shipping: %v", err)
				http.Error(w, "error while invoking meshery Operator", http.StatusInternalServerError)
				return
			}

			log.Debug("received new data on response channel")
			_, _ = fmt.Fprintf(w, "data: %s\n\n", bd)
			if flusher != nil {
				flusher.Flush()
				log.Debugf("Flushed the messages on the wire...")
			}
		}
		endChan <- struct{}{}
		log.Debug("response channel closed")
	}()
	go func() {
		ctx := context.Background()
		h.OperatorHandlerHelper(ctx, w, req, prefObj, user, provider, respChan)
		close(respChan)
	}()
	select {
	case <-notify.Done():
		log.Debugf("received signal to close connection and channels")
		break
	case <-endChan:
		log.Debugf("load test completed")
	}

}

func (h *Handler) OperatorHandlerHelper(ctx context.Context, w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider, respChan chan *OperatorResponse) {
	respChan <- &OperatorResponse{
		Status:  OperatorInfo,
		Message: "Initiating Operator Execution . . . ",
	}
	datach := make(chan *broker.Message)

	client, err := helpers.NewKubeClient(prefObj.K8SConfig.Config)
	if err != nil {
		fmt.Println(err)
		respChan <- &OperatorResponse{
			Status:  OperatorInfo,
			Message: "Error Kube . . . ",
		}
		return
	}
	h.config.KubeClient = client

	enable, err := strconv.ParseBool(req.FormValue("enable"))
	if err != nil {
		fmt.Println(err)
		return
	}

	go h.handleEvents(datach, provider)

	er := h.initialize(!enable, respChan)
	if er != nil {
		status.Error = er.Error()
		respChan <- &OperatorResponse{
			Status:  OperatorError,
			Message: "Error initialize...",
		}
		return
	}
	respChan <- &OperatorResponse{
		Status:  OperatorInfo,
		Message: "Install Operator...",
	}

	er = h.subscribeToBroker(datach, respChan)
	if er != nil {
		status.Error = er.Error()
		respChan <- &OperatorResponse{
			Status:  OperatorError,
			Message: "Sub Error...",
		}
		return
	}
	respChan <- &OperatorResponse{
		Status:  OperatorInfo,
		Message: "Sub Operator...",
	}

	er = h.runMeshSync(!enable, respChan)
	if er != nil {
		status.Error = er.Error()
		respChan <- &OperatorResponse{
			Status:  OperatorError,
			Message: "Syync Error Operator...",
		}
		return
	}
	respChan <- &OperatorResponse{
		Status:  OperatorInfo,
		Message: "Sync Operator...",
	}
}

func (h *Handler) initialize(delete bool, respChan chan *OperatorResponse) error {
	// installOperator
	err := h.applyYaml(delete, operatorYaml)
	if err != nil {
		return err
	}
	status.OperatorInstalled = strconv.FormatBool(!delete)

	respChan <- &OperatorResponse{
		Status:  OperatorInfo,
		Message: "Install Operator...",
	}

	// installBroker
	err = h.applyYaml(delete, brokerYaml)
	if err != nil {
		return err
	}
	status.BrokerInstalled = "true"

	return nil
}

func (h *Handler) subscribeToBroker(datach chan *broker.Message, respChan chan *OperatorResponse) error {
	respChan <- &OperatorResponse{
		Status:  OperatorInfo,
		Message: "SubBroker Started...",
	}
	var broker *operatorv1alpha1.Broker
	mesheryclient, err := client.New(&h.config.KubeClient.RestConfig)
	if err != nil {
		return err
	}
	respChan <- &OperatorResponse{
		Status:  OperatorInfo,
		Message: "MesehryClient Started...",
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
	respChan <- &OperatorResponse{
		Status:  OperatorInfo,
		Message: "Sub Started...",
	}
	return nil
}

func (h *Handler) runMeshSync(delete bool, respChan chan *OperatorResponse) error {
	// installMeshSync
	err := h.applyYaml(delete, meshsyncYaml)
	if err != nil {
		return err
	}
	status.MeshsyncInstalled = "true"
	respChan <- &OperatorResponse{
		Status:  OperatorInfo,
		Message: "Install Operator...",
	}
	return nil
}

func (h *Handler) handleEvents(datach chan *broker.Message, provider models.Provider) {
	status.EventHandlerStarted = "true"
	for {
		select {
		case msg := <-datach:
			objectJSON, _ := utils.Marshal(msg.Object)
			object := model.Object{}
			err := utils.Unmarshal(string(objectJSON), &object)
			if err != nil {
				status.Error += err.Error()
			}

			// persist the object
			err = provider.RecordMeshSyncData(object)
			if err != nil {
				status.Error += err.Error()
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
