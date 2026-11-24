package models

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// The propagation layer is exercised against a stand-in API server rather than
// client-go's fake clientset: meshkit's kubernetes Client holds a concrete
// *kubernetes.Clientset, so a fake typed clientset cannot be substituted into
// it. Driving both the typed and the dynamic client through one httptest
// server also puts the assertions on the wire content of each server-side
// apply, which is what the settings contract is written in terms of.
const (
	ccMeshSyncCRPath = "/apis/meshery.io/v1alpha1/namespaces/meshery/meshsyncs/meshery-meshsync"
	ccBrokerCRPath   = "/apis/meshery.io/v1alpha1/namespaces/meshery/brokers/meshery-broker"
	ccDeploymentPath = "/apis/apps/v1/namespaces/meshery/deployments/meshery-meshsync"
)

// ccStoredObject is one object the stand-in API server holds, together with
// the payload the meshery-server field manager applied last. Keeping the
// previous payload is what lets the server reproduce the withdrawal half of
// server-side apply: a field this manager owned and no longer sets is removed.
type ccStoredObject struct {
	object      map[string]interface{}
	prevApplied map[string]interface{}
}

// ccFakeCluster is a minimal Kubernetes API server that serves GET and
// server-side-apply PATCH for the three objects the controllers config
// propagates to, and records every applied payload.
type ccFakeCluster struct {
	t       *testing.T
	mu      sync.Mutex
	objects map[string]*ccStoredObject
	applied map[string][]map[string]interface{}
	client  *mesherykube.Client
}

// newCCFakeCluster starts a stand-in API server holding exactly the objects
// named in present. Anything omitted is served as a 404, which is how an
// operator that has not deployed yet - and an embedded-mode cluster that never
// will - looks to the propagation layer.
func newCCFakeCluster(t *testing.T, present ...string) *ccFakeCluster {
	t.Helper()

	cluster := &ccFakeCluster{
		t:       t,
		objects: map[string]*ccStoredObject{},
		applied: map[string][]map[string]interface{}{},
	}
	for _, path := range present {
		switch path {
		case ccMeshSyncCRPath:
			cluster.objects[path] = &ccStoredObject{object: ccMeshSyncCR()}
		case ccBrokerCRPath:
			cluster.objects[path] = &ccStoredObject{object: ccBrokerCR()}
		case ccDeploymentPath:
			cluster.objects[path] = &ccStoredObject{object: ccMeshSyncDeployment()}
		default:
			t.Fatalf("unknown object path %q", path)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(cluster.serve))
	t.Cleanup(server.Close)

	config := &rest.Config{Host: server.URL}
	typed, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("typed client: %v", err)
	}
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		t.Fatalf("dynamic client: %v", err)
	}
	cluster.client = &mesherykube.Client{RestConfig: *config, KubeClient: typed, DynamicKubeClient: dyn}

	return cluster
}

func (c *ccFakeCluster) serve(w http.ResponseWriter, r *http.Request) {
	c.mu.Lock()
	defer c.mu.Unlock()

	stored, ok := c.objects[r.URL.Path]
	if !ok {
		ccWriteNotFound(w, r.URL.Path)
		return
	}

	switch r.Method {
	case http.MethodGet:
		ccWriteJSON(w, stored.object)
	case http.MethodPatch:
		if got := r.Header.Get("Content-Type"); got != string(types.ApplyPatchType) {
			c.t.Errorf("%s: patch content type = %q, want server-side apply (%q)", r.URL.Path, got, types.ApplyPatchType)
		}
		if got := r.URL.Query().Get("fieldManager"); got != controllersConfigFieldManager {
			c.t.Errorf("%s: fieldManager = %q, want %q", r.URL.Path, got, controllersConfigFieldManager)
		}
		if got := r.URL.Query().Get("force"); got != "true" {
			c.t.Errorf("%s: force = %q, want true", r.URL.Path, got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			c.t.Errorf("%s: read patch body: %v", r.URL.Path, err)
			return
		}
		payload := map[string]interface{}{}
		if err := json.Unmarshal(body, &payload); err != nil {
			c.t.Errorf("%s: patch body is not JSON: %v", r.URL.Path, err)
			return
		}
		c.applied[r.URL.Path] = append(c.applied[r.URL.Path], payload)
		stored.object = ccServerSideApply(stored.object, stored.prevApplied, payload)
		stored.prevApplied = payload
		ccWriteJSON(w, stored.object)
	default:
		ccWriteNotFound(w, r.URL.Path)
	}
}

// lastApplied returns the most recent payload applied to path, failing the
// test when nothing was applied there.
func (c *ccFakeCluster) lastApplied(path string) map[string]interface{} {
	c.t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()
	payloads := c.applied[path]
	if len(payloads) == 0 {
		c.t.Fatalf("nothing was applied to %s", path)
	}
	return payloads[len(payloads)-1]
}

func (c *ccFakeCluster) applyCount(path string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.applied[path])
}

// stored returns the object as the cluster now holds it, after every apply.
func (c *ccFakeCluster) stored(path string) map[string]interface{} {
	c.t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()
	object, ok := c.objects[path]
	if !ok {
		c.t.Fatalf("%s is not present on the cluster", path)
	}
	return object.object
}

func (c *ccFakeCluster) apply(merged *controllersconfig.MesheryControllersConfig) (*ControllersConfigApplyResult, error) {
	return ApplyControllersConfigToCluster(context.Background(), nil, c.client, merged)
}

// ccServerSideApply reproduces the half of server-side apply this code depends
// on: fields the applying manager owned and no longer sets are withdrawn, and
// the rest of the applied document is merged over what is stored.
func ccServerSideApply(stored, prevApplied, applied map[string]interface{}) map[string]interface{} {
	merged := ccDeepCopy(stored)
	ccWithdraw(merged, prevApplied, applied)
	ccDeepMerge(merged, applied)
	return merged
}

func ccWithdraw(into, prevApplied, applied map[string]interface{}) {
	for key, previous := range prevApplied {
		current, stillApplied := applied[key]
		if !stillApplied {
			delete(into, key)
			continue
		}
		previousMap, previousIsMap := previous.(map[string]interface{})
		currentMap, currentIsMap := current.(map[string]interface{})
		intoMap, intoIsMap := into[key].(map[string]interface{})
		if previousIsMap && currentIsMap && intoIsMap {
			ccWithdraw(intoMap, previousMap, currentMap)
		}
	}
}

func ccDeepMerge(into, from map[string]interface{}) {
	for key, value := range from {
		valueMap, valueIsMap := value.(map[string]interface{})
		intoMap, intoIsMap := into[key].(map[string]interface{})
		if valueIsMap && intoIsMap {
			ccDeepMerge(intoMap, valueMap)
			continue
		}
		into[key] = value
	}
}

func ccDeepCopy(object map[string]interface{}) map[string]interface{} {
	encoded, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	copied := map[string]interface{}{}
	if err := json.Unmarshal(encoded, &copied); err != nil {
		panic(err)
	}
	return copied
}

func ccWriteJSON(w http.ResponseWriter, object map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(object)
}

func ccWriteNotFound(w http.ResponseWriter, path string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"kind":       "Status",
		"apiVersion": "v1",
		"status":     "Failure",
		"reason":     "NotFound",
		"code":       http.StatusNotFound,
		"message":    path + " not found",
	})
}

func ccMeshSyncCR() map[string]interface{} {
	return map[string]interface{}{
		"apiVersion": "meshery.io/v1alpha1",
		"kind":       "MeshSync",
		"metadata": map[string]interface{}{
			"name":      meshSyncCRName,
			"namespace": controllersNamespace,
		},
		"spec": map[string]interface{}{},
	}
}

func ccBrokerCR() map[string]interface{} {
	return map[string]interface{}{
		"apiVersion": "meshery.io/v1alpha1",
		"kind":       "Broker",
		"metadata": map[string]interface{}{
			"name":      brokerCRName,
			"namespace": controllersNamespace,
		},
		"spec": map[string]interface{}{},
	}
}

func ccMeshSyncDeployment() map[string]interface{} {
	return map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name":      meshSyncDeploymentName,
			"namespace": controllersNamespace,
		},
		"spec": map[string]interface{}{
			"selector": map[string]interface{}{
				"matchLabels": map[string]interface{}{"app": meshSyncDeploymentName},
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{"app": meshSyncDeploymentName},
				},
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  meshSyncContainerName,
							"image": "meshery/meshsync:stable-latest",
						},
					},
				},
			},
		},
	}
}

func ccStr(value string) *string { return &value }

func ccInt(value int) *int { return &value }

func ccBool(value bool) *bool { return &value }

func ccServiceType(value controllersconfig.MesheryBrokerServiceConfigType) *controllersconfig.MesheryBrokerServiceConfigType {
	return &value
}

// ccAt walks a decoded JSON document, failing the test when the path is absent.
func ccAt(t *testing.T, document map[string]interface{}, path ...string) interface{} {
	t.Helper()
	var cursor interface{} = document
	for index, key := range path {
		asMap, ok := cursor.(map[string]interface{})
		if !ok {
			t.Fatalf("%v: %q is not an object", path[:index], path[index-1])
		}
		cursor, ok = asMap[key]
		if !ok {
			t.Fatalf("%v is absent from the applied document", path[:index+1])
		}
	}
	return cursor
}

// ccAbsent asserts a path is not present, which is how the propagation layer
// withdraws a field: the applied document simply stops describing it.
func ccAbsent(t *testing.T, document map[string]interface{}, path ...string) {
	t.Helper()
	var cursor interface{} = document
	for _, key := range path {
		asMap, ok := cursor.(map[string]interface{})
		if !ok {
			return
		}
		cursor, ok = asMap[key]
		if !ok {
			return
		}
	}
	t.Fatalf("%v is still present as %#v; it should have been withdrawn", path, cursor)
}

// ccDecodeJSONString decodes a JSON document that the propagation layer
// embedded as a string, which is how the MeshSync CR carries its watch-list.
func ccDecodeJSONString(t *testing.T, raw interface{}) interface{} {
	t.Helper()
	asString, ok := raw.(string)
	if !ok {
		t.Fatalf("expected an embedded JSON string, got %#v", raw)
	}
	var decoded interface{}
	if err := json.Unmarshal([]byte(asString), &decoded); err != nil {
		t.Fatalf("embedded JSON %q: %v", asString, err)
	}
	return decoded
}

// ccMeshSyncContainer returns the single container the Deployment overlay
// describes.
func ccMeshSyncContainer(t *testing.T, payload map[string]interface{}) map[string]interface{} {
	t.Helper()
	containers, ok := ccAt(t, payload, "spec", "template", "spec", "containers").([]interface{})
	if !ok || len(containers) != 1 {
		t.Fatalf("expected exactly one container in the overlay, got %#v", containers)
	}
	container, ok := containers[0].(map[string]interface{})
	if !ok {
		t.Fatalf("container is not an object: %#v", containers[0])
	}
	if name := container["name"]; name != meshSyncContainerName {
		t.Fatalf("overlay targets container %v, want %q", name, meshSyncContainerName)
	}
	return container
}

// ccContainerEnv reduces the overlay's env list to a name/value map.
func ccContainerEnv(t *testing.T, container map[string]interface{}) map[string]string {
	t.Helper()
	env := map[string]string{}
	entries, ok := container["env"].([]interface{})
	if !ok {
		return env
	}
	for _, raw := range entries {
		entry, ok := raw.(map[string]interface{})
		if !ok {
			t.Fatalf("env entry is not an object: %#v", raw)
		}
		name, _ := entry["name"].(string)
		value, _ := entry["value"].(string)
		env[name] = value
	}
	return env
}

func ccContainerArgs(t *testing.T, container map[string]interface{}) []string {
	t.Helper()
	args := []string{}
	entries, ok := container["args"].([]interface{})
	if !ok {
		return args
	}
	for _, raw := range entries {
		arg, ok := raw.(string)
		if !ok {
			t.Fatalf("arg is not a string: %#v", raw)
		}
		args = append(args, arg)
	}
	return args
}

// TestApplyControllersConfigMeshSyncCRSettings covers meshsync.version,
// meshsync.replicas and meshsync.watchList: the three settings the doc's
// MeshSync table routes to the MeshSync custom resource.
func TestApplyControllersConfigMeshSyncCRSettings(t *testing.T) {
	t.Run("version, replicas and a whitelist", func(t *testing.T) {
		cluster := newCCFakeCluster(t, ccMeshSyncCRPath, ccBrokerCRPath, ccDeploymentPath)

		result, err := cluster.apply(&controllersconfig.MesheryControllersConfig{
			Meshsync: &controllersconfig.MeshSyncConfig{
				Version:  ccStr("v1.0.2"),
				Replicas: ccInt(3),
				WatchList: &controllersconfig.MeshSyncWatchList{
					Whitelist: []controllersconfig.MeshSyncWatchedResource{
						{Resource: "pods.v1.", Events: []controllersconfig.MeshSyncWatchedResourceEvents{controllersconfig.ADDED, controllersconfig.MODIFIED}},
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("apply: %v", err)
		}
		if !result.MeshSyncCRPatched {
			t.Fatalf("result does not report the MeshSync custom resource as patched: %+v", result)
		}

		payload := cluster.lastApplied(ccMeshSyncCRPath)
		if got := ccAt(t, payload, "spec", "version"); got != "v1.0.2" {
			t.Errorf("meshsync.version -> spec.version = %#v, want %q", got, "v1.0.2")
		}
		if got := ccAt(t, payload, "spec", "size"); got != float64(3) {
			t.Errorf("meshsync.replicas -> spec.size = %#v, want 3", got)
		}

		whitelist := ccDecodeJSONString(t, ccAt(t, payload, "spec", "watch-list", "data", "whitelist"))
		want := []interface{}{
			map[string]interface{}{"Resource": "pods.v1.", "Events": []interface{}{"ADDED", "MODIFIED"}},
		}
		if !reflect.DeepEqual(whitelist, want) {
			t.Errorf("meshsync.watchList -> spec.watch-list.data.whitelist = %#v, want %#v", whitelist, want)
		}
		if got := ccAt(t, payload, "spec", "watch-list", "data", "blacklist"); got != "" {
			t.Errorf("the unused watch-list key is %#v, want it claimed as empty so a chart-set counterpart cannot linger", got)
		}
	})

	t.Run("a blacklist claims the whitelist key as empty", func(t *testing.T) {
		cluster := newCCFakeCluster(t, ccMeshSyncCRPath, ccBrokerCRPath, ccDeploymentPath)

		if _, err := cluster.apply(&controllersconfig.MesheryControllersConfig{
			Meshsync: &controllersconfig.MeshSyncConfig{
				WatchList: &controllersconfig.MeshSyncWatchList{Blacklist: []string{"pods.v1.", "events.v1."}},
			},
		}); err != nil {
			t.Fatalf("apply: %v", err)
		}

		payload := cluster.lastApplied(ccMeshSyncCRPath)
		blacklist := ccDecodeJSONString(t, ccAt(t, payload, "spec", "watch-list", "data", "blacklist"))
		want := []interface{}{"pods.v1.", "events.v1."}
		if !reflect.DeepEqual(blacklist, want) {
			t.Errorf("meshsync.watchList -> spec.watch-list.data.blacklist = %#v, want %#v", blacklist, want)
		}
		if got := ccAt(t, payload, "spec", "watch-list", "data", "whitelist"); got != "" {
			t.Errorf("the unused watch-list key is %#v, want it claimed as empty", got)
		}
	})
}

// TestApplyControllersConfigBrokerCRSettings covers every row of the doc's
// Meshery Broker table: broker.version, broker.replicas and the five
// broker.service.* settings.
func TestApplyControllersConfigBrokerCRSettings(t *testing.T) {
	cluster := newCCFakeCluster(t, ccMeshSyncCRPath, ccBrokerCRPath, ccDeploymentPath)

	result, err := cluster.apply(&controllersconfig.MesheryControllersConfig{
		Broker: &controllersconfig.MesheryBrokerConfig{
			Version:  ccStr("2.10.18"),
			Replicas: ccInt(2),
			Service: &controllersconfig.MesheryBrokerServiceConfig{
				Type:                     ccServiceType(controllersconfig.LoadBalancer),
				Annotations:              map[string]string{"service.beta.kubernetes.io/aws-load-balancer-internal": "true"},
				LoadBalancerClass:        ccStr("service.k8s.aws/nlb"),
				LoadBalancerSourceRanges: []string{"10.0.0.0/8", "192.168.0.0/16"},
				ExternalEndpointOverride: ccStr("broker.example.com:4222"),
			},
		},
	})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if !result.BrokerCRPatched {
		t.Fatalf("result does not report the Broker custom resource as patched: %+v", result)
	}

	payload := cluster.lastApplied(ccBrokerCRPath)
	if got := ccAt(t, payload, "spec", "version"); got != "2.10.18" {
		t.Errorf("broker.version -> spec.version = %#v, want %q", got, "2.10.18")
	}
	if got := ccAt(t, payload, "spec", "size"); got != float64(2) {
		t.Errorf("broker.replicas -> spec.size = %#v, want 2", got)
	}
	if got := ccAt(t, payload, "spec", "service", "type"); got != "LoadBalancer" {
		t.Errorf("broker.service.type -> spec.service.type = %#v, want LoadBalancer", got)
	}
	wantAnnotations := map[string]interface{}{"service.beta.kubernetes.io/aws-load-balancer-internal": "true"}
	if got := ccAt(t, payload, "spec", "service", "annotations"); !reflect.DeepEqual(got, wantAnnotations) {
		t.Errorf("broker.service.annotations -> spec.service.annotations = %#v, want %#v", got, wantAnnotations)
	}
	if got := ccAt(t, payload, "spec", "service", "loadBalancerClass"); got != "service.k8s.aws/nlb" {
		t.Errorf("broker.service.loadBalancerClass -> spec.service.loadBalancerClass = %#v", got)
	}
	wantRanges := []interface{}{"10.0.0.0/8", "192.168.0.0/16"}
	if got := ccAt(t, payload, "spec", "service", "loadBalancerSourceRanges"); !reflect.DeepEqual(got, wantRanges) {
		t.Errorf("broker.service.loadBalancerSourceRanges -> spec.service.loadBalancerSourceRanges = %#v, want %#v", got, wantRanges)
	}
	if got := ccAt(t, payload, "spec", "service", "externalEndpointOverride"); got != "broker.example.com:4222" {
		t.Errorf("broker.service.externalEndpointOverride -> spec.service.externalEndpointOverride = %#v", got)
	}
}

// TestApplyControllersConfigDeploymentOverlaySettings covers the five MeshSync
// rows the doc routes to the MeshSync Deployment rather than to its custom
// resource: three env knobs and two output filters.
func TestApplyControllersConfigDeploymentOverlaySettings(t *testing.T) {
	cluster := newCCFakeCluster(t, ccMeshSyncCRPath, ccBrokerCRPath, ccDeploymentPath)

	result, err := cluster.apply(&controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{
			RedactSecrets:      ccBool(true),
			BrokerContentDedup: ccBool(true),
			DebugLogging:       ccBool(false),
			OutputNamespaces:   []string{"meshery", "kube-system"},
			OutputResources:    []string{"pod", "deployment"},
		},
	})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if !result.DeploymentOverlayApplied {
		t.Fatalf("result does not report the Deployment overlay as applied: %+v", result)
	}

	container := ccMeshSyncContainer(t, cluster.lastApplied(ccDeploymentPath))

	env := ccContainerEnv(t, container)
	wantEnv := map[string]string{
		envRedactSecrets:      "true",
		envBrokerContentDedup: "true",
		envDebug:              "false",
	}
	if !reflect.DeepEqual(env, wantEnv) {
		t.Errorf("meshsync.redactSecrets/brokerContentDedup/debugLogging -> env = %#v, want %#v", env, wantEnv)
	}

	args := ccContainerArgs(t, container)
	wantArgs := []string{"--outputNamespaces=meshery,kube-system", "--outputResources=pod,deployment"}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Errorf("meshsync.outputNamespaces/outputResources -> args = %#v, want %#v", args, wantArgs)
	}
}

// TestApplyControllersConfigWithdrawsClearedFields covers the withdrawal
// contract: clearing a field at every layer stops the applied document from
// describing it, which under server-side apply removes it from the cluster
// object instead of leaving the last-applied value behind.
func TestApplyControllersConfigWithdrawsClearedFields(t *testing.T) {
	cluster := newCCFakeCluster(t, ccMeshSyncCRPath, ccBrokerCRPath, ccDeploymentPath)

	if _, err := cluster.apply(&controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{
			Version:            ccStr("v1.0.2"),
			Replicas:           ccInt(3),
			WatchList:          &controllersconfig.MeshSyncWatchList{Blacklist: []string{"pods.v1."}},
			RedactSecrets:      ccBool(true),
			BrokerContentDedup: ccBool(true),
			DebugLogging:       ccBool(true),
			OutputNamespaces:   []string{"meshery"},
			OutputResources:    []string{"pod"},
		},
		Broker: &controllersconfig.MesheryBrokerConfig{
			Version:  ccStr("2.10.18"),
			Replicas: ccInt(2),
			Service: &controllersconfig.MesheryBrokerServiceConfig{
				Type:                     ccServiceType(controllersconfig.LoadBalancer),
				LoadBalancerClass:        ccStr("service.k8s.aws/nlb"),
				LoadBalancerSourceRanges: []string{"10.0.0.0/8"},
				ExternalEndpointOverride: ccStr("broker.example.com:4222"),
				Annotations:              map[string]string{"a": "b"},
			},
		},
	}); err != nil {
		t.Fatalf("first apply: %v", err)
	}

	// Every layer now clears every field. A nil document is the same case.
	if _, err := cluster.apply(nil); err != nil {
		t.Fatalf("withdrawing apply: %v", err)
	}

	meshsyncPayload := cluster.lastApplied(ccMeshSyncCRPath)
	for _, path := range [][]string{{"spec", "version"}, {"spec", "size"}, {"spec", "watch-list"}} {
		ccAbsent(t, meshsyncPayload, path...)
	}

	brokerPayload := cluster.lastApplied(ccBrokerCRPath)
	for _, path := range [][]string{{"spec", "version"}, {"spec", "size"}, {"spec", "service"}} {
		ccAbsent(t, brokerPayload, path...)
	}

	container := ccMeshSyncContainer(t, cluster.lastApplied(ccDeploymentPath))
	if env := ccContainerEnv(t, container); len(env) != 0 {
		t.Errorf("cleared env knobs are still applied: %#v", env)
	}
	if args := ccContainerArgs(t, container); len(args) != 0 {
		t.Errorf("cleared output filters are still applied: %#v", args)
	}

	// The fields left the cluster objects, not just the applied documents.
	ccAbsent(t, cluster.stored(ccMeshSyncCRPath), "spec", "version")
	ccAbsent(t, cluster.stored(ccMeshSyncCRPath), "spec", "watch-list")
	ccAbsent(t, cluster.stored(ccBrokerCRPath), "spec", "service")
}

// TestApplyControllersConfigRestartsMeshSyncOnlyOnWatchListChange covers the
// restart contract: MeshSync reads its watch-list at startup only, so a
// watch-list change must roll its pods and any other change must not.
func TestApplyControllersConfigRestartsMeshSyncOnlyOnWatchListChange(t *testing.T) {
	cluster := newCCFakeCluster(t, ccMeshSyncCRPath, ccBrokerCRPath, ccDeploymentPath)

	watchList := &controllersconfig.MeshSyncWatchList{Blacklist: []string{"pods.v1."}}

	result, err := cluster.apply(&controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{WatchList: watchList},
	})
	if err != nil {
		t.Fatalf("apply setting a watch-list: %v", err)
	}
	if !result.MeshSyncRestarted {
		t.Fatalf("setting a watch-list did not restart MeshSync: %+v", result)
	}
	firstRestartValue, _ := ccAt(t, cluster.lastApplied(ccDeploymentPath), "spec", "template", "metadata", "annotations", meshSyncRestartAnnotation).(string)
	if firstRestartValue == "" {
		t.Fatalf("no restart annotation was stamped on the MeshSync pod template")
	}

	// Re-applying the same watch-list changes nothing, so it must not roll pods.
	result, err = cluster.apply(&controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{WatchList: watchList},
	})
	if err != nil {
		t.Fatalf("re-applying the same watch-list: %v", err)
	}
	if result.MeshSyncRestarted {
		t.Errorf("re-applying an unchanged watch-list restarted MeshSync: %+v", result)
	}
	if got, _ := ccAt(t, cluster.lastApplied(ccDeploymentPath), "spec", "template", "metadata", "annotations", meshSyncRestartAnnotation).(string); got != firstRestartValue {
		t.Errorf("restart annotation = %q, want the previous value %q carried forward so the apply itself does not roll pods", got, firstRestartValue)
	}

	// A change to any other setting must not roll pods either.
	result, err = cluster.apply(&controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{
			WatchList:     watchList,
			Replicas:      ccInt(4),
			RedactSecrets: ccBool(true),
		},
	})
	if err != nil {
		t.Fatalf("applying a non-watch-list change: %v", err)
	}
	if result.MeshSyncRestarted {
		t.Errorf("a replicas/redactSecrets change restarted MeshSync: %+v", result)
	}
	if got, _ := ccAt(t, cluster.lastApplied(ccDeploymentPath), "spec", "template", "metadata", "annotations", meshSyncRestartAnnotation).(string); got != firstRestartValue {
		t.Errorf("restart annotation = %q, want the previous value %q carried forward", got, firstRestartValue)
	}

	// Withdrawing the watch-list changes what MeshSync watches, so it rolls.
	result, err = cluster.apply(&controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{Replicas: ccInt(4)},
	})
	if err != nil {
		t.Fatalf("withdrawing the watch-list: %v", err)
	}
	if !result.MeshSyncRestarted {
		t.Errorf("withdrawing the watch-list did not restart MeshSync: %+v", result)
	}
}

// TestApplyControllersConfigReportsAbsentTargetsAsSkipped covers the skip
// contract: an operator that has not deployed yet, and an embedded-mode
// cluster that never will, are reported in Skipped rather than raised as
// errors, because the configuration re-applies when the connection reconnects.
func TestApplyControllersConfigReportsAbsentTargetsAsSkipped(t *testing.T) {
	t.Run("nothing deployed", func(t *testing.T) {
		cluster := newCCFakeCluster(t)

		result, err := cluster.apply(&controllersconfig.MesheryControllersConfig{
			Meshsync: &controllersconfig.MeshSyncConfig{Replicas: ccInt(3), RedactSecrets: ccBool(true)},
			Broker:   &controllersconfig.MesheryBrokerConfig{Replicas: ccInt(2)},
		})
		if err != nil {
			t.Fatalf("an absent operator must not be an error, got: %v", err)
		}
		if result.MeshSyncCRPatched || result.BrokerCRPatched || result.DeploymentOverlayApplied || result.MeshSyncRestarted {
			t.Errorf("nothing is deployed, yet the result reports writes: %+v", result)
		}
		if len(result.Skipped) != 3 {
			t.Fatalf("Skipped = %#v, want one entry per absent target", result.Skipped)
		}
		for _, want := range []string{"MeshSync custom resource not present", "Broker custom resource not present", "MeshSync deployment not present"} {
			found := false
			for _, skipped := range result.Skipped {
				if len(skipped) >= len(want) && skipped[:len(want)] == want {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Skipped = %#v, missing an entry starting %q", result.Skipped, want)
			}
		}
	})

	t.Run("custom resources absent, deployment present", func(t *testing.T) {
		cluster := newCCFakeCluster(t, ccDeploymentPath)

		result, err := cluster.apply(&controllersconfig.MesheryControllersConfig{
			Meshsync: &controllersconfig.MeshSyncConfig{DebugLogging: ccBool(true)},
		})
		if err != nil {
			t.Fatalf("apply: %v", err)
		}
		if !result.DeploymentOverlayApplied {
			t.Errorf("the present Deployment was not written: %+v", result)
		}
		if result.MeshSyncCRPatched || result.BrokerCRPatched {
			t.Errorf("absent custom resources were reported as patched: %+v", result)
		}
		if len(result.Skipped) != 2 {
			t.Errorf("Skipped = %#v, want the two absent custom resources", result.Skipped)
		}
		if cluster.applyCount(ccDeploymentPath) != 1 {
			t.Errorf("Deployment applied %d times, want 1", cluster.applyCount(ccDeploymentPath))
		}
	})
}

// TestApplyControllersConfigRequiresAKubernetesClient covers the one case that
// is an error rather than a skip: there is no client to reach the cluster
// with, so nothing can be said about what is deployed.
func TestApplyControllersConfigRequiresAKubernetesClient(t *testing.T) {
	result, err := ApplyControllersConfigToCluster(context.Background(), nil, nil, &controllersconfig.MesheryControllersConfig{})
	if err == nil {
		t.Fatalf("a nil kubernetes client must be an error, got result %+v", result)
	}
	if result == nil {
		t.Fatalf("a result must be returned alongside the error so callers can report what was reached")
	}
}
