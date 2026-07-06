package models

import (
	"context"
	"encoding/json"
	goerrors "errors"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/meshery/meshkit/logger"
	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

const (
	// controllersNamespace is where the Meshery Operator deploys the MeshSync
	// and Broker workloads and their custom resources.
	controllersNamespace = "meshery"
	// meshSyncCRName is the MeshSync custom resource the operator chart
	// installs and the MeshSync agent reads its watch-list from.
	meshSyncCRName = "meshery-meshsync"
	// brokerCRName is the Broker custom resource the operator chart installs.
	brokerCRName = "meshery-broker"
	// meshSyncDeploymentName is the Deployment the operator reconciles for
	// the MeshSync CR (same name as the CR).
	meshSyncDeploymentName = "meshery-meshsync"
	// meshSyncContainerName is the MeshSync container within that Deployment.
	meshSyncContainerName = "meshsync"
	// controllersConfigFieldManager is the server-side-apply field manager
	// Meshery Server uses for its Deployment overlay. The operator applies
	// the Deployment under its own field manager and never sets these env
	// names, args, or the restart annotation, so ownership stays disjoint.
	controllersConfigFieldManager = "meshery-server"
	// meshSyncRestartAnnotation triggers a rolling restart of MeshSync pods
	// when its value changes (same mechanism as kubectl rollout restart).
	meshSyncRestartAnnotation = "meshery.io/restarted-at"

	envRedactSecrets      = "MESHSYNC_REDACT_SECRETS"
	envBrokerContentDedup = "MESHSYNC_BROKER_CONTENT_DEDUP"
	envDebug              = "DEBUG"
)

var (
	meshSyncGVR = schema.GroupVersionResource{Group: "meshery.io", Version: "v1alpha1", Resource: "meshsyncs"}
	brokerGVR   = schema.GroupVersionResource{Group: "meshery.io", Version: "v1alpha1", Resource: "brokers"}
)

// ControllersConfigApplyResult reports which propagation targets a
// configuration apply reached on a cluster.
type ControllersConfigApplyResult struct {
	MeshSyncCRPatched        bool     `json:"meshSyncCRPatched"`
	BrokerCRPatched          bool     `json:"brokerCRPatched"`
	DeploymentOverlayApplied bool     `json:"deploymentOverlayApplied"`
	MeshSyncRestarted        bool     `json:"meshSyncRestarted"`
	Skipped                  []string `json:"skipped,omitempty"`
}

// ApplyControllersConfigToCluster propagates the merged (explicitly-set)
// controllers configuration to a cluster running in operator deployment
// mode. Every target is managed with server-side apply under the
// meshery-server field manager, so the applied configuration always
// describes the complete set of fields Meshery Server owns: setting a field
// takes ownership, and unsetting it at every layer withdraws it on the next
// apply, reverting to the operator's or chart's own value.
//
//   - MeshSync CR: spec.version, spec.size, spec.watch-list
//   - Broker CR: spec.version, spec.size, spec.service
//   - MeshSync Deployment: env (MESHSYNC_REDACT_SECRETS,
//     MESHSYNC_BROKER_CONTENT_DEDUP, DEBUG) and args (--outputNamespaces,
//     --outputResources)
//
// A nil merged document is applied as an empty one: that is the withdrawal
// case (no layer sets anything anymore), not a no-op. MeshSync reads its
// watch-list at startup only, so a watch-list change also triggers a rolling
// restart of the MeshSync Deployment. Absent CRs or Deployment (operator not
// deployed yet, or embedded-mode cluster) are skipped and reported, not
// treated as errors: configuration re-applies when the connection
// reconnects.
func ApplyControllersConfigToCluster(
	ctx context.Context,
	log logger.Handler,
	kubeClient *mesherykube.Client,
	merged *controllersconfig.MesheryControllersConfig,
) (*ControllersConfigApplyResult, error) {
	result := &ControllersConfigApplyResult{}
	if kubeClient == nil {
		return result, ErrApplyControllersConfig(k8serrors.NewBadRequest("no kubernetes client available for the connection"))
	}
	if merged == nil {
		merged = &controllersconfig.MesheryControllersConfig{}
	}

	var applyErrs []error

	watchListChanged, err := applyMeshSyncCR(ctx, kubeClient, merged.Meshsync, result)
	if err != nil {
		applyErrs = append(applyErrs, err)
	}

	if err := applyBrokerCR(ctx, kubeClient, merged.Broker, result); err != nil {
		applyErrs = append(applyErrs, err)
	}

	if err := applyMeshSyncDeploymentOverlay(ctx, kubeClient, merged.Meshsync, watchListChanged, result); err != nil {
		applyErrs = append(applyErrs, err)
	}

	if len(applyErrs) > 0 {
		return result, ErrApplyControllersConfig(goerrors.Join(applyErrs...))
	}
	if log != nil {
		log.Debugf("controllers config applied: %+v", result)
	}
	return result, nil
}

// applyMeshSyncCR server-side-applies the explicitly-set MeshSync fields
// onto the MeshSync custom resource under the meshery-server field manager.
// The applied configuration always describes the complete owned set, so
// fields cleared at every layer are withdrawn and revert to the operator's
// or chart's own values. It returns whether the watch-list changed as a
// result of the apply (which requires a MeshSync restart to take effect,
// since MeshSync reads its CR at startup only).
func applyMeshSyncCR(
	ctx context.Context,
	kubeClient *mesherykube.Client,
	cfg *controllersconfig.MeshSyncConfig,
	result *ControllersConfigApplyResult,
) (bool, error) {
	crClient := kubeClient.DynamicKubeClient.Resource(meshSyncGVR).Namespace(controllersNamespace)
	before, err := crClient.Get(ctx, meshSyncCRName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) || isNoKindMatch(err) {
			result.Skipped = append(result.Skipped, "MeshSync custom resource not present; watch-list, version, and replica settings apply when the Meshery Operator is deployed")
			return false, nil
		}
		return false, err
	}

	spec := map[string]interface{}{}
	if cfg != nil {
		if cfg.Version != nil {
			spec["version"] = *cfg.Version
		}
		if cfg.Replicas != nil {
			spec["size"] = *cfg.Replicas
		}
		if cfg.WatchList != nil {
			// Claim both keys whenever a watch-list is managed: the unused
			// key is set to the empty string (which MeshSync ignores) so a
			// pre-existing chart-set counterpart cannot linger alongside the
			// managed key and trip MeshSync's whitelist-XOR-blacklist
			// validation. Withdrawing the watch-list releases both keys.
			data := map[string]interface{}{"whitelist": "", "blacklist": ""}
			if len(cfg.WatchList.Whitelist) > 0 {
				encoded, err := json.Marshal(watchListWhitelistWireEntries(cfg.WatchList.Whitelist))
				if err != nil {
					return false, err
				}
				data["whitelist"] = string(encoded)
			}
			if len(cfg.WatchList.Blacklist) > 0 {
				encoded, err := json.Marshal(cfg.WatchList.Blacklist)
				if err != nil {
					return false, err
				}
				data["blacklist"] = string(encoded)
			}
			spec["watch-list"] = map[string]interface{}{"data": data}
		}
	}

	applyConfig := map[string]interface{}{
		"apiVersion": meshSyncGVR.Group + "/" + meshSyncGVR.Version,
		"kind":       "MeshSync",
		"metadata": map[string]interface{}{
			"name":      meshSyncCRName,
			"namespace": controllersNamespace,
		},
		"spec": spec,
	}
	payload, err := json.Marshal(applyConfig)
	if err != nil {
		return false, err
	}
	force := true
	after, err := crClient.Patch(ctx, meshSyncCRName, types.ApplyPatchType, payload, metav1.PatchOptions{FieldManager: controllersConfigFieldManager, Force: &force})
	if err != nil {
		return false, err
	}
	result.MeshSyncCRPatched = true
	return !watchListDataEqual(before.Object, after.Object), nil
}

// watchListWhitelistWireEntries converts schema whitelist entries into the
// wire shape MeshSync's config parser expects: {"Resource": ..., "Events":
// [...]} with Go-style field names (meshsync unmarshals into an untagged
// struct).
func watchListWhitelistWireEntries(entries []controllersconfig.MeshSyncWatchedResource) []map[string]interface{} {
	wire := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		events := make([]string, 0, len(entry.Events))
		for _, ev := range entry.Events {
			events = append(events, string(ev))
		}
		wire = append(wire, map[string]interface{}{
			"Resource": entry.Resource,
			"Events":   events,
		})
	}
	return wire
}

// watchListDataEqual reports whether two MeshSync CR objects carry the same
// effective spec.watch-list.data (whitelist/blacklist), comparing the JSON
// values regardless of formatting. Absent and empty entries are equivalent.
func watchListDataEqual(a, b map[string]interface{}) bool {
	dataOf := func(obj map[string]interface{}) map[string]interface{} {
		if spec, ok := obj["spec"].(map[string]interface{}); ok {
			if wl, ok := spec["watch-list"].(map[string]interface{}); ok {
				if data, ok := wl["data"].(map[string]interface{}); ok {
					return data
				}
			}
		}
		return map[string]interface{}{}
	}
	aData, bData := dataOf(a), dataOf(b)
	for _, key := range []string{"whitelist", "blacklist"} {
		aString, _ := aData[key].(string)
		bString, _ := bData[key].(string)
		if aString == "" && bString == "" {
			continue
		}
		if !jsonStringsEquivalent(aString, bString) {
			return false
		}
	}
	return true
}

// jsonStringsEquivalent reports whether two JSON documents encoded as strings
// carry the same value regardless of formatting.
func jsonStringsEquivalent(a, b string) bool {
	var av, bv interface{}
	if err := json.Unmarshal([]byte(a), &av); err != nil {
		return a == b
	}
	if err := json.Unmarshal([]byte(b), &bv); err != nil {
		return a == b
	}
	return reflect.DeepEqual(av, bv)
}

// applyBrokerCR server-side-applies the explicitly-set Broker fields onto
// the Broker custom resource under the meshery-server field manager, with
// the same complete-owned-set withdrawal semantics as applyMeshSyncCR.
// Broker service changes reconcile in place; version and size changes roll
// the NATS statefulset under operator control.
func applyBrokerCR(
	ctx context.Context,
	kubeClient *mesherykube.Client,
	cfg *controllersconfig.MesheryBrokerConfig,
	result *ControllersConfigApplyResult,
) error {
	crClient := kubeClient.DynamicKubeClient.Resource(brokerGVR).Namespace(controllersNamespace)
	if _, err := crClient.Get(ctx, brokerCRName, metav1.GetOptions{}); err != nil {
		if k8serrors.IsNotFound(err) || isNoKindMatch(err) {
			result.Skipped = append(result.Skipped, "Broker custom resource not present; broker settings apply when the Meshery Operator is deployed")
			return nil
		}
		return err
	}

	spec := map[string]interface{}{}
	if cfg != nil {
		if cfg.Version != nil {
			spec["version"] = *cfg.Version
		}
		if cfg.Replicas != nil {
			spec["size"] = *cfg.Replicas
		}
		if svc := cfg.Service; svc != nil {
			service := map[string]interface{}{}
			if svc.Type != nil {
				service["type"] = string(*svc.Type)
			}
			if svc.Annotations != nil {
				service["annotations"] = svc.Annotations
			}
			if svc.LoadBalancerClass != nil {
				service["loadBalancerClass"] = *svc.LoadBalancerClass
			}
			if svc.LoadBalancerSourceRanges != nil {
				service["loadBalancerSourceRanges"] = svc.LoadBalancerSourceRanges
			}
			if svc.ExternalEndpointOverride != nil {
				service["externalEndpointOverride"] = *svc.ExternalEndpointOverride
			}
			if len(service) > 0 {
				spec["service"] = service
			}
		}
	}

	applyConfig := map[string]interface{}{
		"apiVersion": brokerGVR.Group + "/" + brokerGVR.Version,
		"kind":       "Broker",
		"metadata": map[string]interface{}{
			"name":      brokerCRName,
			"namespace": controllersNamespace,
		},
		"spec": spec,
	}
	payload, err := json.Marshal(applyConfig)
	if err != nil {
		return err
	}
	force := true
	if _, err := crClient.Patch(ctx, brokerCRName, types.ApplyPatchType, payload, metav1.PatchOptions{FieldManager: controllersConfigFieldManager, Force: &force}); err != nil {
		return err
	}
	result.BrokerCRPatched = true
	return nil
}

// applyMeshSyncDeploymentOverlay server-side-applies Meshery Server's
// env/args overlay onto the MeshSync Deployment. The applied configuration
// always describes the complete set of fields this field manager owns, so
// clearing a knob at every layer withdraws its env entry or argument on the
// next apply. When restartMeshSync is true (watch-list changed), the pod
// template restart annotation is refreshed; otherwise any previously-applied
// annotation value is carried forward unchanged so the apply itself does not
// roll pods.
func applyMeshSyncDeploymentOverlay(
	ctx context.Context,
	kubeClient *mesherykube.Client,
	cfg *controllersconfig.MeshSyncConfig,
	restartMeshSync bool,
	result *ControllersConfigApplyResult,
) error {
	deployment, err := kubeClient.KubeClient.AppsV1().Deployments(controllersNamespace).Get(ctx, meshSyncDeploymentName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			result.Skipped = append(result.Skipped, "MeshSync deployment not present; env and output-filter settings apply when the Meshery Operator has deployed MeshSync")
			return nil
		}
		return err
	}

	env := []map[string]interface{}{}
	args := []string{}
	if cfg != nil {
		if cfg.RedactSecrets != nil {
			env = append(env, map[string]interface{}{"name": envRedactSecrets, "value": strconv.FormatBool(*cfg.RedactSecrets)})
		}
		if cfg.BrokerContentDedup != nil {
			env = append(env, map[string]interface{}{"name": envBrokerContentDedup, "value": strconv.FormatBool(*cfg.BrokerContentDedup)})
		}
		if cfg.DebugLogging != nil {
			env = append(env, map[string]interface{}{"name": envDebug, "value": strconv.FormatBool(*cfg.DebugLogging)})
		}
		if len(cfg.OutputNamespaces) > 0 {
			args = append(args, "--outputNamespaces="+strings.Join(cfg.OutputNamespaces, ","))
		}
		if len(cfg.OutputResources) > 0 {
			args = append(args, "--outputResources="+strings.Join(cfg.OutputResources, ","))
		}
	}

	container := map[string]interface{}{"name": meshSyncContainerName}
	if len(env) > 0 {
		container["env"] = env
	}
	if len(args) > 0 {
		container["args"] = args
	}

	templateMeta := map[string]interface{}{}
	previousRestartValue := deployment.Spec.Template.Annotations[meshSyncRestartAnnotation]
	switch {
	case restartMeshSync:
		templateMeta["annotations"] = map[string]interface{}{meshSyncRestartAnnotation: time.Now().UTC().Format(time.RFC3339)}
		result.MeshSyncRestarted = true
	case previousRestartValue != "":
		// Carry the previously-applied annotation forward: dropping it from
		// this manager's applied set would remove it from the pod template,
		// which itself rolls the pods.
		templateMeta["annotations"] = map[string]interface{}{meshSyncRestartAnnotation: previousRestartValue}
	}

	template := map[string]interface{}{
		"spec": map[string]interface{}{
			"containers": []interface{}{container},
		},
	}
	if len(templateMeta) > 0 {
		template["metadata"] = templateMeta
	}

	applyConfig := map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name":      meshSyncDeploymentName,
			"namespace": controllersNamespace,
		},
		"spec": map[string]interface{}{
			"template": template,
		},
	}

	payload, err := json.Marshal(applyConfig)
	if err != nil {
		return err
	}
	force := true
	if _, err := kubeClient.KubeClient.AppsV1().Deployments(controllersNamespace).Patch(ctx, meshSyncDeploymentName, types.ApplyPatchType, payload, metav1.PatchOptions{FieldManager: controllersConfigFieldManager, Force: &force}); err != nil {
		return err
	}
	result.DeploymentOverlayApplied = true
	return nil
}

// isNoKindMatch reports whether the error indicates the CRD itself is not
// installed on the cluster (embedded-mode clusters never install the
// meshery.io CRDs).
func isNoKindMatch(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no matches for kind") ||
		strings.Contains(err.Error(), "could not find the requested resource") ||
		strings.Contains(err.Error(), "the server could not find the requested resource")
}
