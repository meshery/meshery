package model

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"fortio.org/fortio/log"
	"github.com/layer5io/meshery-operator/api/v1alpha1"
	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/broadcast"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	controllerConfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	//platform = runtime.GOOS
	chartRepo = "https://meshery.github.io/meshery.io/charts"
	brokerResourceName   = "brokers"
	brokerInstanceName   = "meshery-broker"
	meshsyncResourceName = "meshsyncs"
	meshsyncInstanceName = "meshery-meshsync"
	brokderCRDName  = "brokers.meshery.layer5.io"
	meshsyncCRDName = "meshsyncs.meshery.layer5.io"
)

var (
	controlPlaneNamespace = map[MeshType]string{
		MeshTypeIstio:              "istio-system",
		MeshTypeLinkerd:            "linkerd-system",
		MeshTypeConsul:             "consul-system",
		MeshTypeOctarine:           "octarine-system",
		MeshTypeTraefikMesh:        "traefik-system",
		MeshTypeOpenServiceMesh:    "osm-system",
		MeshTypeKuma:               "kuma-system",
		MeshTypeNginxServiceMesh:   "nginx-system",
		MeshTypeNetworkServiceMesh: "nsm-system",
		MeshTypeCitrixServiceMesh:  "ctrix-system",
		MeshTypeAppMesh:            "appmesh-system",
	}

	addonPortSelector = map[string]string{
		"grafana":          "service",
		"prometheus":       "http",
		"jaeger-collector": "jaeger-collector-http",
		"kiali":            "http",
		"zipkin":           "http-query",
	}

	//downloadLocation = path.Join(utils.GetHome(), ".meshery", "manifests")
)

// listernToEvents - scale this function with the number of channels
func ListernToEvents(log logger.Handler,
	handler *database.Handler,
	datach chan *broker.Message,
	meshsyncCh chan struct{},
	operatorSyncChannel chan bool,
	controlPlaneSyncChannel chan struct{},
	meshsyncLivenessChannel chan struct{},
	broadcast broadcast.Broadcaster,
) {
	var wg sync.WaitGroup
	for msg := range datach {
		wg.Add(1)
		go persistData(*msg, log, handler, meshsyncCh, operatorSyncChannel, controlPlaneSyncChannel, broadcast, &wg)
	}

	wg.Wait()
}

// persistData - scale this function with the number of events to persist
func persistData(msg broker.Message,
	log logger.Handler,
	handler *database.Handler,
	meshsyncCh chan struct{},
	operatorSyncChannel chan bool,
	controlPlaneSyncChannel chan struct{},
	broadcaster broadcast.Broadcaster,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	objectJSON, _ := utils.Marshal(msg.Object)
	switch msg.ObjectType {
	case broker.MeshSync:
		object := meshsyncmodel.Object{}
		err := utils.Unmarshal(string(objectJSON), &object)
		if err != nil {
			log.Error(err)
			return
		}

		// persist the object
		log.Info("Incoming object: ", object.ObjectMeta.Name, ", kind: ", object.Kind)
		if object.ObjectMeta.Name == "meshery-operator" || object.ObjectMeta.Name == "meshery-broker" || object.ObjectMeta.Name == "meshery-meshsync" {
			// operatorSyncChannel <- false
			broadcaster.Submit(broadcast.BroadcastMessage{
				Source: broadcast.OperatorSyncChannel,
				Data:   false,
				Type:   "health",
			})
		}
		err = recordMeshSyncData(msg.EventType, handler, &object)
		if err != nil {
			log.Error(err)
			return
		}
		meshsyncCh <- struct{}{}
	case broker.SMI:
		log.Info("Received SMI Result")
	}
}

func applyYaml(client *mesherykube.Client, delete bool, file string) error {
	contents, err := utils.ReadLocalFile(file)
	if err != nil {
		return err
	}

	err = client.ApplyManifest([]byte(contents), mesherykube.ApplyOptions{
		Namespace: Namespace,
		Update:    true,
		Delete:    delete,
	})
	if err != nil {
		return err
	}

	return nil
}

// installUsingHelm is for installing helm dependencies. We need this because
// meshery operator and controllers don't have published separate charts but
// exist as subcharts for meshery's chart.
// We plan to have separate charts for these components once we wish to offer
// users the control over which version of helm want to use
func installUsingHelm(client *mesherykube.Client, delete bool) error {
	releaseVersion := viper.GetString("BUILD")
	if releaseVersion == "" || releaseVersion == "Not Set" || releaseVersion == "edge-latest" {
		latestReleaseData, err := handlers.CheckLatestVersion("")
		if err != nil {
			releaseVersion = latestReleaseData.Current
		}
	}

	var (
		act = mesherykube.INSTALL
		chart = "meshery-operator"
		overrideValues map[string]interface{}
		err error
	)

	// if meshery server is in cluster, we upgrade the release
	if viper.GetString("KUBERNETES_SERVICE_HOST") != "" {
		act = mesherykube.UPGRADE
		chart = "meshery"
	} else if delete {
		act = mesherykube.UNINSTALL
	}

	if delete {
		// Delete the CR instances for brokers and meshsyncs
		// this needs to be executed before deleting the helm release, or the CR instances cannot be found for some reason
		if err := DeleteCR(brokerResourceName, brokerInstanceName, client); err != nil {
			return errors.Wrap(err, "cannot delete CR "+brokerInstanceName)
		}
		if err := DeleteCR(meshsyncResourceName, meshsyncInstanceName, client); err != nil {
			return errors.Wrap(err, "cannot delete CR "+meshsyncInstanceName)
		}
	}

	overrideValues, err = SetOverrideValues(delete)
	if err != nil {
		return err
	}

	log.Infof("Installing operator using chart: %s", chart)

	err = client.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
		Namespace: "meshery",
		ChartLocation: mesherykube.HelmChartLocation{
			Repository: chartRepo,
			Chart: chart,
			Version: releaseVersion,
		},
		CreateNamespace: true,
		Action:    act,
		OverrideValues: overrideValues,
	})
	if err != nil {
		return err
	}

	if delete {
		// Delete the CRDs for brokers and meshsyncs
		// These need to be deleted after deleting the chart
		if err = DeleteCRD(brokderCRDName); err != nil {
			return errors.Wrap(err, "cannot delete CRD "+brokderCRDName)
		}

		if err = DeleteCRD(meshsyncCRDName); err != nil {
			return errors.Wrap(err, "cannot delete CRD "+meshsyncCRDName)
		}
	}

	return nil
}

// SetOverrideValues detects the currently insalled adapters and sets uninstalls
// enable flag for operator to true or false
func SetOverrideValues(delete bool) (map[string]interface{}, error) {
	adapters := make([]string, 0)
	installedAdapters := helpers.NewAdaptersTracker(viper.GetStringSlice("ADAPTER_URLS")).GetAdapters(context.TODO())

	for _, adapter := range installedAdapters {
		if adapter.Name != "" {
			adapters = append(adapters, strings.Split(adapter.Location, ":")[0])
		}
	}

	overrideValues := map[string]interface{} {
		// TODO: update operator chart and remove meshery
		"meshery": map[string]interface{}{
			"enabled": false,
		},
		"meshery-istio": map[string]interface{}{
			"enabled": false,
		},
		"meshery-linkerd": map[string]interface{}{
			"enabled": false,
		},
		"meshery-consul": map[string]interface{}{
			"enabled": false,
		},
		"meshery-kuma": map[string]interface{}{
			"enabled": false,
		},
		"meshery-osm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nsm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nginx-sm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-traefik-mesh": map[string]interface{}{
			"enabled": false,
		},
		"meshery-cpx": map[string]interface{}{
			"enabled": false,
		},
		"meshery-app-mesh": map[string]interface{}{
			"enabled": false,
		},
		"meshery-operator": map[string]interface{}{
			"enabled": true,
		},
		"meshery-meshsync": map[string]interface{}{
			"enabled": true,
		},
		"meshery-broker": map[string]interface{}{
			"enabled": true,
		},
	}

	for _, adapter := range adapters {
		if _, ok := overrideValues[adapter]; ok {
			fmt.Println(adapter)
			overrideValues[adapter] = map[string]interface{} {
				"enabled": true,
			}
		}
	}

	if delete {
		overrideValues["meshery-operator"] = map[string]interface{} {
			"enabled": false,
		}
		overrideValues["meshery-meshsync"] = map[string]interface{} {
			"enabled": false,
		}
		overrideValues["meshery-broker"] = map[string]interface{} {
			"enabled": false,
		}
	}

	return overrideValues, nil
}

// DeleteCRs delete the specified CR instance in the clusters
func DeleteCR(resourceName, instanceName string, client *mesherykube.Client) error {
	return client.DynamicKubeClient.Resource(schema.GroupVersionResource{
		Group:    v1alpha1.GroupVersion.Group,
		Version:  v1alpha1.GroupVersion.Version,
		Resource: resourceName,
	}).Namespace("meshery").Delete(context.TODO(), instanceName, metav1.DeleteOptions{})
}

// DeleteCRs delete the specified CRD in the clusters
func DeleteCRD(name string) error {
	cfg := controllerConfig.GetConfigOrDie()
	client, err := apiextension.NewForConfig(cfg)
	if err != nil {
		return errors.Wrap(err, "cannot invoke delete CRDs")
	}
	return client.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), name, metav1.DeleteOptions{})
}
