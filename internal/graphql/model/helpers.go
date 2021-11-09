package model

import (
	"context"
	"sync"

	"github.com/layer5io/meshery-operator/api/v1alpha1"
	"github.com/layer5io/meshery/handlers"
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

//	releaseName := fmt.Sprintf("meshery-%s", releaseVersion)

//	err := getHelmChart(releaseName)
//	if err != nil {
//		return err
//	}

// downloadedChartLocation := path.Join(downloadLocation, releaseName, "meshery")

//	// install CRDs
//	err = applyYaml(client, delete, path.Join(downloadedChartLocation, "crds", "crds.yaml"))
//	if err != nil {
//		return err
//	}

// dependencyLocation := path.Join(downloadedChartLocation, "charts", dependencyName)
	var act mesherykube.HelmChartAction
	act = mesherykube.INSTALL
	if delete {
		act = mesherykube.UNINSTALL
		// Delete the CR instances for brokers and meshsyncs
		// this needs to be executed before deleting the helm release, or the CR instances cannot be found for some reason
		if err := DeleteCR(brokerResourceName, brokerInstanceName, client); err != nil {
			return errors.Wrap(err, "cannot delete CR "+brokerInstanceName)
		}
		if err := DeleteCR(meshsyncResourceName, meshsyncInstanceName, client); err != nil {
			return errors.Wrap(err, "cannot delete CR "+meshsyncInstanceName)
		}
	}

	err := client.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
		Namespace: "meshery",
		ChartLocation: mesherykube.HelmChartLocation{
			Repository: chartRepo,
			Chart: "meshery-operator",
			Version: releaseVersion,
		},
		CreateNamespace: true,
		Action:    act,
		OverrideValues: map[string]interface{}{
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
		},
//		LocalPath: dependencyLocation,
	})
	if err != nil {
		return err
	}

	if delete {
		// Delete the CRDs for brokers and meshsyncs
		if err = DeleteCRD(brokderCRDName); err != nil {
			return errors.Wrap(err, "cannot delete CRD "+brokderCRDName)
		}

		if err = DeleteCRD(meshsyncCRDName); err != nil {
			return errors.Wrap(err, "cannot delete CRD "+meshsyncCRDName)
		}
	}

	return nil
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

//func getHelmChart(releaseName string) error {
//	logrus.Println("Looking for parent chart in ", downloadLocation)
//	// see if chart for current release already exist on fs
//	_, err := os.Stat(path.Join(downloadLocation, releaseName))
//	if err == nil {
//		return nil
//	}
//
//	logrus.Println("Downloading chart: ", releaseName)
//
//	chartURL := fmt.Sprintf("https://meshery.github.io/meshery.io/charts/%s.tgz", releaseName)
//
//	err = downloadTar(chartURL, releaseName)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func downloadTar(url, releaseName string) error {
//	resp, err := http.Get(url)
//	if err != nil {
//		return err
//	}
//
//	if resp.StatusCode != http.StatusOK {
//		_ = resp.Body.Close()
//		return err
//	}
//
//	// make the directory corresponding to release version
//	if err := os.MkdirAll(path.Join(downloadLocation, releaseName), 0750); err != nil {
//		return err
//	}
//
//	err = tarxzf(path.Join(downloadLocation, releaseName), resp.Body)
//	return err
//}
//
//func tarxzf(location string, stream io.Reader) error {
//	uncompressedStream, err := gzip.NewReader(stream)
//	if err != nil {
//		return err
//	}
//
//	tarReader := tar.NewReader(uncompressedStream)
//
//	for {
//		header, err := tarReader.Next()
//
//		if err == io.EOF {
//			break
//		}
//
//		if err != nil {
//			return err
//		}
//
//		switch header.Typeflag {
//		case tar.TypeDir:
//			// File traversal is required to store the extracted manifests at the right place
//			// #nosec
//			if err := os.MkdirAll(path.Join(location, header.Name), 0750); err != nil {
//				return err
//			}
//		case tar.TypeReg:
//			// A fallback case because the downloaded tar doesn't have proper header
//			// entries for dir
//			// #nosec
//			if _, err := os.Stat(path.Join(location, path.Dir(header.Name))); err != nil {
//				if err := os.MkdirAll(path.Join(location, path.Dir(header.Name)), 0750); err != nil {
//					return err
//				}
//			}
//
//			outFile, err := os.Create(path.Join(location, header.Name))
//			if err != nil {
//				return err
//			}
//			if _, err := io.CopyN(outFile, tarReader, header.Size); err != nil {
//				return err
//			}
//			if err = outFile.Close(); err != nil {
//				return err
//			}
//
//		default:
//			return err
//		}
//	}
//
//	return nil
//}
