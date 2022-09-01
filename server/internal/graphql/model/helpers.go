package model

import (
	"context"
	"database/sql"
	"strings"
	"sync"

	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/controllers"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/broadcast"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// to be moved elsewhere
const (
	chartRepo = "https://meshery.github.io/meshery.io/charts"
)

var (
	controlPlaneNamespace = map[MeshType][]string{
		MeshTypeIstio:              {"istio-system"},
		MeshTypeLinkerd:            {"linkerd-system"},
		MeshTypeCiliumServiceMesh:  {"kube-system"},
		MeshTypeConsul:             {"consul-system"},
		MeshTypeOctarine:           {"octarine-system"},
		MeshTypeTraefikMesh:        {"traefik-system"},
		MeshTypeOpenServiceMesh:    {"osm-system"},
		MeshTypeKuma:               {"kuma-system"},
		MeshTypeNginxServiceMesh:   {"nginx-system"},
		MeshTypeNetworkServiceMesh: {"nsm-system"},
		MeshTypeCitrixServiceMesh:  {"citrix-system"},
		MeshTypeAppMesh:            {"appmesh-system"},
		//Any namespace added or appended above should also be appended on the AllMesh array
		MeshTypeAllMesh: {"istio-system", "linkerd-system", "consul-system", "octarine-system", "traefik-system", "osm-system", "kuma-system", "nginx-system", "nsm-system", "citrix-system", "appmesh-system"},
	}

	addonPortSelector = map[string]string{
		"grafana":          "service",
		"prometheus":       "http",
		"jaeger-collector": "jaeger-collector-http",
		"kiali":            "http",
		"zipkin":           "http-query",
	}
)
var (
	//TODO: Add the image orgs of other control plane pods. This change is backwards compatible and wont break anything
	controlPlaneImageOrgs = map[MeshType][]string{
		MeshTypeCiliumServiceMesh: {"cilium"},
	}
)

// listernToEvents - scale this function with the number of channels
func ListernToEvents(log logger.Handler,
	handler *database.Handler,
	datach chan *broker.Message,
	meshsyncCh chan struct{},
	broadcast broadcast.Broadcaster,
) {
	var wg sync.WaitGroup
	for msg := range datach {
		wg.Add(1)
		go persistData(*msg, log, handler, meshsyncCh, broadcast, &wg)
	}

	wg.Wait()
}

// persistData - scale this function with the number of events to persist
func persistData(msg broker.Message,
	log logger.Handler,
	handler *database.Handler,
	meshsyncCh chan struct{},
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
		if meshsyncCh != nil {
			meshsyncCh <- struct{}{}
		}
	case broker.SMI:
		log.Info("Received SMI Result")
	}
}

func PersistClusterNames(
	ctx context.Context,
	log logger.Handler,
	handler *database.Handler,
	meshsyncCh chan struct{},
) {
	k8sContexts, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
	if !ok {
		return
	}
	for _, clusterConfig := range k8sContexts {
		clusterName := clusterConfig.Cluster["name"].(string)
		clusterID := clusterConfig.KubernetesServerID.String()
		object := meshsyncmodel.Object{
			Kind: "Cluster",
			ObjectMeta: &meshsyncmodel.ResourceObjectMeta{
				Name:      clusterName,
				ClusterID: clusterID,
			},
			ClusterID: clusterID,
		}

		// persist the object
		log.Info("Incoming object: ", object.ObjectMeta.Name, ", kind: ", object.Kind)
		err := recordMeshSyncData(broker.Add, handler, &object)
		if err != nil {
			log.Error(err)
		}
	}
	if meshsyncCh != nil {
		meshsyncCh <- struct{}{}
	}
}

func applyYaml(client *mesherykube.Client, delete bool, file string) error {
	contents, err := utils.ReadRemoteFile(file)
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

// installs operator
// To be depricated
func installUsingHelm(client *mesherykube.Client, delete bool, adapterTracker models.AdaptersTrackerInterface) error {
	// retrieving meshery's version to apply the appropriate chart
	mesheryReleaseVersion := viper.GetString("BUILD")
	if mesheryReleaseVersion == "" || mesheryReleaseVersion == "Not Set" || mesheryReleaseVersion == "edge-latest" {
		_, latestRelease, err := handlers.CheckLatestVersion(mesheryReleaseVersion)
		// if unable to fetch latest release tag, meshkit helm functions handle
		// this automatically fetch the latest one
		if err != nil {
			logrus.Errorf("Couldn't check release tag: %s. Will use latest version", err)
			mesheryReleaseVersion = ""
		} else {
			mesheryReleaseVersion = latestRelease
		}
	}
	var (
		act   = mesherykube.INSTALL
		chart = "meshery-operator"
	)
	if delete {
		act = mesherykube.UNINSTALL
	}
	// a basic check to see if meshery is installed in cluster
	// this helps decide what chart should be used for installing operator
	if viper.GetString("KUBERNETES_SERVICE_HOST") != "" {
		// act = mesherykube.UPGRADE
		chart = "meshery"
	}

	overrides := SetOverrideValues(delete, adapterTracker)

	err := client.ApplyHelmChart(mesherykube.ApplyHelmChartConfig{
		Namespace:   "meshery",
		ReleaseName: "meshery",
		ChartLocation: mesherykube.HelmChartLocation{
			Repository: chartRepo,
			Chart:      chart,
			Version:    mesheryReleaseVersion,
		},
		// CreateNamespace doesn't have any effect when the action is UNINSTALL
		CreateNamespace: true,
		Action:          act,
		// Setting override values
		OverrideValues: overrides,
	})
	if err != nil {
		return ErrApplyHelmChart(err)
	}

	return nil
}

// SetOverrideValues detects the currently insalled adapters and sets appropriate
// overrides so as to not uninstall them. It also sets override values for
// operator so that it can be enabled or disabled depending on the need

// to be depricated
func SetOverrideValues(delete bool, adapterTracker models.AdaptersTrackerInterface) map[string]interface{} {
	installedAdapters := make([]string, 0)
	adapters := adapterTracker.GetAdapters(context.TODO())

	for _, adapter := range adapters {
		if adapter.Name != "" {
			installedAdapters = append(installedAdapters, strings.Split(adapter.Location, ":")[0])
		}
	}

	overrideValues := map[string]interface{}{
		"fullnameOverride": "meshery-operator",
		"meshery": map[string]interface{}{
			"enabled": false,
		},
		"meshery-istio": map[string]interface{}{
			"enabled": false,
		},
		"meshery-cilium": map[string]interface{}{
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
	}

	for _, adapter := range installedAdapters {
		if _, ok := overrideValues[adapter]; ok {
			overrideValues[adapter] = map[string]interface{}{
				"enabled": true,
			}
		}
	}

	if delete {
		overrideValues["meshery-operator"] = map[string]interface{}{
			"enabled": false,
		}
	}

	return overrideValues
}

// K8sConnectionTracker keeps track of BrokerURLs per kubernetes context
type K8sConnectionTracker struct {
	mx              sync.Mutex
	contextToBroker map[string]string //ContextID -> BrokerURL
}

func NewK8sConnctionTracker() *K8sConnectionTracker {
	return &K8sConnectionTracker{
		contextToBroker: make(map[string]string),
	}
}
func (k *K8sConnectionTracker) Set(id string, url string) {
	k.mx.Lock()
	defer k.mx.Unlock()
	k.contextToBroker[id] = url
}

// Takes a set of endpoints and discard the current endpoint if its not present in the set
func (k *K8sConnectionTracker) ResetEndpoints(available map[string]bool) {
	k.mx.Lock()
	defer k.mx.Unlock()
	c := make(map[string]string)
	for id, url := range k.contextToBroker {
		if available[url] {
			c[id] = url
		}
	}
	k.contextToBroker = c
}
func (k *K8sConnectionTracker) ListBrokerEndpoints() (a []string) {
	k.mx.Lock()
	defer k.mx.Unlock()
	for _, v := range k.contextToBroker {
		a = append(a, v)
	}
	return
}
func (k *K8sConnectionTracker) Get(id string) (url string) {
	k.mx.Lock()
	defer k.mx.Unlock()
	url = k.contextToBroker[id]
	return
}

// Takes the meshkit Logger and logs a comma separated list of currently tracked Broker Endpoints
func (k *K8sConnectionTracker) Log(l logger.Handler) {
	var e = "Connected broker endpoints : "
	k.mx.Lock()
	defer k.mx.Unlock()
	for _, v := range k.contextToBroker {
		e += v + ", "
	}
	l.Info(strings.TrimSuffix(e, ", "))
}

func GetInternalController(controller models.MesheryController) MesheryController {
	switch controller {
	case models.MesheryBroker:
		return MesheryControllerBroker
	case models.MesheryOperator:
		return MesheryControllerOperator
	case models.Meshsync:
		return MesheryControllerMeshsync
	}
	return ""
}

func GetInternalControllerStatus(status controllers.MesheryControllerStatus) MesheryControllerStatus {
	switch status {
	case controllers.Deployed:
		return MesheryControllerStatusDeployed

	case controllers.NotDeployed:
		return MesheryControllerStatusNotdeployed

	case controllers.Deploying:
		return MesheryControllerStatusDeploying

	case controllers.Unknown:
		return MesheryControllerStatusUnkown
	}
	return ""
}

// SelectivelyFetchNamespaces fetches the an array of namespaces from DB based on ClusterIDs (or KubernetesServerIDs) passed as param
func SelectivelyFetchNamespaces(cids []string, provider models.Provider) ([]string, error) {
	namespaces := make([]string, 0)
	var rows *sql.Rows
	var err error
	rows, err = provider.GetGenericPersister().Raw("SELECT DISTINCT rom.name as name FROM objects o LEFT JOIN resource_object_meta rom ON o.id = rom.id WHERE o.kind = 'Namespace' AND o.cluster_id IN ?", cids).Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		namespaces = append(namespaces, name)
	}
	return namespaces, nil
}
