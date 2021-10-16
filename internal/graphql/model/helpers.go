package model

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"sync"

	"github.com/layer5io/meshery/handlers"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/broadcast"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"github.com/spf13/viper"
)

const (
	platform = runtime.GOOS
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

	downloadLocation = path.Join(utils.GetHome(), ".meshery", "charts")
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

// installUsingHelm is for installing helm dependencies. We need this because
// meshery operator and controllers don't have published separate charts but are
// dependencies for meshery's chart.
// We plan to have separate charts for these components once we wish to offer
// users the control over which version of helm want to use
func installUsingHelm(client *mesherykube.Client, delete bool, dependencyName string) error {
	releaseVersion := viper.GetString("BUILD")
	if releaseVersion == "" || releaseVersion == "Not Set" || releaseVersion == "edge-latest" {
		latestReleaseData, err := handlers.CheckLatestVersion("")
		if err != nil {
			releaseVersion = latestReleaseData.Current
		}
	}

	releaseName := fmt.Sprintf("meshery-%s", releaseVersion)

	err := getHelmChart(releaseName)
	if err != nil {
		return err
	}

	dependencyLocation := path.Join(downloadLocation, releaseName, "meshery", "charts", dependencyName)

	err = client.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
		Namespace:       "meshery",
		Delete:          delete,
		CreateNamespace: true,
		LocalPath:       dependencyLocation,
	})
	if err != nil {
		return err
	}

	return nil
}

func getHelmChart(releaseName string) error {
	fmt.Println("Looking for parent chart in ", downloadLocation)
	// see if chart for current release already exist on fs
	_, err := os.Stat(path.Join(downloadLocation, releaseName))
	if os.IsExist(err) {
		return nil
	}

	fmt.Println("Downloading chart: ", releaseName)

	chartURL := fmt.Sprintf("https://meshery.github.io/meshery.io/charts/%s.tgz", releaseName)
	fmt.Println(chartURL)

	err = downloadTar(chartURL, releaseName)
	if err != nil {
		return err
	}

	return nil
}

func downloadTar(url, releaseName string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return err
	}

	// make the directory corresponding to release version
	if err := os.MkdirAll(path.Join(downloadLocation, releaseName), 0750); err != nil {
		return err
	}

	if err := tarxzf(path.Join(downloadLocation, releaseName), resp.Body); err != nil {
		return err
	}

	return nil
}

func tarxzf(location string, stream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(stream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// File traversal is required to store the extracted manifests at the right place
			// #nosec
			if err := os.MkdirAll(path.Join(location, header.Name), 0750); err != nil {
				return err
			}
		case tar.TypeReg:
			// A fallback case because the downloaded tar doesn't have proper header
			// entries for dir
			// #nosec
			if _, err := os.Stat(path.Join(location, path.Dir(header.Name))); os.IsNotExist(err) {
				if err := os.MkdirAll(path.Join(location, path.Dir(header.Name)), 0750); err != nil {
					return err
				}
			}

			outFile, err := os.Create(path.Join(location, header.Name))
			if err != nil {
				return err
			}
			// Trust meshery tar
			// #nosec
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			if err = outFile.Close(); err != nil {
				return err
			}

		default:
			return err
		}
	}

	return nil
}
