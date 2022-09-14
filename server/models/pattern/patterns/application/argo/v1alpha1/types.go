package v1alpha1

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Rollout is a specification for a Rollout resource
type Rollout struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   RolloutSpec   `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	Status RolloutStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// RolloutSpec is the spec for a Rollout resource
type RolloutSpec struct {
	TemplateResolvedFromRef bool `json:"-"`
	SelectorResolvedFromRef bool `json:"-"`
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`
	// Label selector for pods. Existing ReplicaSets whose pods are
	// selected by this will be the ones affected by this rollout.
	// It must match the pod template's labels.
	// +optional
	Selector *metav1.LabelSelector `json:"selector" protobuf:"bytes,2,opt,name=selector"`
	// Template describes the pods that will be created.
	// +optional
	Template corev1.PodTemplateSpec `json:"template" protobuf:"bytes,3,opt,name=template"`
	// WorkloadRef holds a references to a workload that provides Pod template
	// +optional
	WorkloadRef *ObjectRef `json:"workloadRef,omitempty" protobuf:"bytes,10,opt,name=workloadRef"`
	// Minimum number of seconds for which a newly created pod should be ready
	// without any of its container crashing, for it to be considered available.
	// Defaults to 0 (pod will be considered available as soon as it is ready)
	// +optional
	MinReadySeconds int32 `json:"minReadySeconds,omitempty" protobuf:"varint,4,opt,name=minReadySeconds"`
	// The deployment strategy to use to replace existing pods with new ones.
	// +optional
	Strategy RolloutStrategy `json:"strategy" protobuf:"bytes,5,opt,name=strategy"`
	// The number of old ReplicaSets to retain. If unspecified, will retain 10 old ReplicaSets
	RevisionHistoryLimit *int32 `json:"revisionHistoryLimit,omitempty" protobuf:"varint,6,opt,name=revisionHistoryLimit"`
	// Paused pauses the rollout at its current step.
	Paused bool `json:"paused,omitempty" protobuf:"varint,7,opt,name=paused"`
	// ProgressDeadlineSeconds The maximum time in seconds for a rollout to
	// make progress before it is considered to be failed. Argo Rollouts will
	// continue to process failed rollouts and a condition with a
	// ProgressDeadlineExceeded reason will be surfaced in the rollout status.
	// Note that progress will not be estimated during the time a rollout is paused.
	// Defaults to 600s.
	ProgressDeadlineSeconds *int32 `json:"progressDeadlineSeconds,omitempty" protobuf:"varint,8,opt,name=progressDeadlineSeconds"`
	// RestartAt indicates when all the pods of a Rollout should be restarted
	RestartAt *metav1.Time `json:"restartAt,omitempty" protobuf:"bytes,9,opt,name=restartAt"`
}

func (s *RolloutSpec) SetResolvedSelector(selector *metav1.LabelSelector) {
	s.SelectorResolvedFromRef = true
	s.Selector = selector
}

func (s *RolloutSpec) SetResolvedTemplate(template corev1.PodTemplateSpec) {
	s.TemplateResolvedFromRef = true
	s.Template = template
}

func (s *RolloutSpec) MarshalJSON() ([]byte, error) {
	type Alias RolloutSpec

	if s.TemplateResolvedFromRef || s.SelectorResolvedFromRef {
		obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&struct {
			Alias `json:",inline"`
		}{
			Alias: (Alias)(*s),
		})
		if err != nil {
			return nil, err
		}
		if s.TemplateResolvedFromRef {
			unstructured.RemoveNestedField(obj, "template")
		}
		if s.SelectorResolvedFromRef {
			unstructured.RemoveNestedField(obj, "selector")
		}

		return json.Marshal(obj)
	}
	return json.Marshal(&struct{ *Alias }{
		Alias: (*Alias)(s),
	})
}

// ObjectRef holds a references to the Kubernetes object
type ObjectRef struct {
	// API Version of the referent
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,1,opt,name=apiVersion"`
	// Kind of the referent
	Kind string `json:"kind,omitempty" protobuf:"bytes,2,opt,name=kind"`
	// Name of the referent
	Name string `json:"name,omitempty" protobuf:"bytes,3,opt,name=name"`
}

const (
	// DefaultRolloutUniqueLabelKey is the default key of the selector that is added
	// to existing ReplicaSets (and label key that is added to its pods) to prevent the existing ReplicaSets
	// to select new pods (and old pods being select by new ReplicaSet).
	DefaultRolloutUniqueLabelKey string = "rollouts-pod-template-hash"
	// DefaultReplicaSetScaleDownDeadlineAnnotationKey is the default key attached to an old stable ReplicaSet after
	// the rollout transitioned to a new version. It contains the time when the controller can scale down the RS.
	DefaultReplicaSetScaleDownDeadlineAnnotationKey = "scale-down-deadline"
	// ManagedByRolloutKey is the key used to indicate which rollout(s) manage a resource but doesn't own it.
	ManagedByRolloutsKey = "argo-rollouts.argoproj.io/managed-by-rollouts"
	// DefaultReplicaSetRestartAnnotationKey indicates that the ReplicaSet with this annotation was restarted at the
	// time listed in the value
	DefaultReplicaSetRestartAnnotationKey = "argo-rollouts.argoproj.io/restarted-after"
	// LabelKeyControllerInstanceID is the label the controller uses for the rollout, experiment, analysis segregation
	// between controllers. Controllers will only operate on objects with the same instanceID as the controller.
	LabelKeyControllerInstanceID = "argo-rollouts.argoproj.io/controller-instance-id"
)

// RolloutStrategy defines strategy to apply during next rollout
type RolloutStrategy struct {
	// +optional
	BlueGreen *BlueGreenStrategy `json:"blueGreen,omitempty" protobuf:"bytes,1,opt,name=blueGreen"`
	// +optional
	Canary *CanaryStrategy `json:"canary,omitempty" protobuf:"bytes,2,opt,name=canary"`
}

// BlueGreenStrategy defines parameters for Blue Green deployment
type BlueGreenStrategy struct {
	// Name of the service that the rollout modifies as the active service.
	ActiveService string `json:"activeService" protobuf:"bytes,1,opt,name=activeService"`
	// Name of the service that the rollout modifies as the preview service.
	// +optional
	PreviewService string `json:"previewService,omitempty" protobuf:"bytes,2,opt,name=previewService"`
	// PreviewReplicaCount is the number of replicas to run for the preview stack before the
	// switchover. Once the rollout is resumed the desired replicaset will be full scaled up before the switch occurs
	// +optional
	PreviewReplicaCount *int32 `json:"previewReplicaCount,omitempty" protobuf:"varint,3,opt,name=previewReplicaCount"`
	// AutoPromotionEnabled indicates if the rollout should automatically promote the new ReplicaSet
	// to the active service or enter a paused state. If not specified, the default value is true.
	// +optional
	AutoPromotionEnabled *bool `json:"autoPromotionEnabled,omitempty" protobuf:"varint,4,opt,name=autoPromotionEnabled"`
	// AutoPromotionSeconds is a duration in seconds in which to delay auto-promotion (default: 0).
	// The countdown begins after the preview ReplicaSet have reached full availability.
	// This option is ignored if autoPromotionEnabled is set to false.
	// +optional
	AutoPromotionSeconds int32 `json:"autoPromotionSeconds,omitempty" protobuf:"varint,5,opt,name=autoPromotionSeconds"`
	// MaxUnavailable The maximum number of pods that can be unavailable during a restart operation.
	// Defaults to 25% of total replicas.
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty" protobuf:"bytes,6,opt,name=maxUnavailable"`
	// ScaleDownDelaySeconds adds a delay before scaling down the previous replicaset.
	// If omitted, the Rollout waits 30 seconds before scaling down the previous ReplicaSet.
	// A minimum of 30 seconds is recommended to ensure IP table propagation across the nodes in
	// a cluster. See https://github.com/argoproj/argo-rollouts/issues/19#issuecomment-476329960 for
	// more information
	// +optional
	ScaleDownDelaySeconds *int32 `json:"scaleDownDelaySeconds,omitempty" protobuf:"varint,7,opt,name=scaleDownDelaySeconds"`
	// ScaleDownDelayRevisionLimit limits the number of old RS that can run at one time before getting scaled down
	// +optional
	ScaleDownDelayRevisionLimit *int32 `json:"scaleDownDelayRevisionLimit,omitempty" protobuf:"varint,8,opt,name=scaleDownDelayRevisionLimit"`
	// PrePromotionAnalysis configuration to run analysis before a selector switch
	PrePromotionAnalysis *RolloutAnalysis `json:"prePromotionAnalysis,omitempty" protobuf:"bytes,9,opt,name=prePromotionAnalysis"`
	// AntiAffinity enables anti-affinity rules for Blue Green deployment
	// +optional
	AntiAffinity *AntiAffinity `json:"antiAffinity,omitempty" protobuf:"bytes,10,opt,name=antiAffinity"`
	// PostPromotionAnalysis configuration to run analysis after a selector switch
	PostPromotionAnalysis *RolloutAnalysis `json:"postPromotionAnalysis,omitempty" protobuf:"bytes,11,opt,name=postPromotionAnalysis"`
	// PreviewMetadata specify labels and annotations which will be attached to the preview pods for
	// the duration which they act as a preview pod, and will be removed after
	PreviewMetadata *PodTemplateMetadata `json:"previewMetadata,omitempty" protobuf:"bytes,12,opt,name=previewMetadata"`
	// ActiveMetadata specify labels and annotations which will be attached to the active pods for
	// the duration which they act as a active pod, and will be removed after
	ActiveMetadata *PodTemplateMetadata `json:"activeMetadata,omitempty" protobuf:"bytes,13,opt,name=activeMetadata"`
}

// AntiAffinity defines which inter-pod scheduling rule to use for anti-affinity injection
type AntiAffinity struct {
	// +optional
	PreferredDuringSchedulingIgnoredDuringExecution *PreferredDuringSchedulingIgnoredDuringExecution `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,opt,name=preferredDuringSchedulingIgnoredDuringExecution"`
	// +optional
	RequiredDuringSchedulingIgnoredDuringExecution *RequiredDuringSchedulingIgnoredDuringExecution `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,opt,name=requiredDuringSchedulingIgnoredDuringExecution"`
}

// PreferredDuringSchedulingIgnoredDuringExecution defines the weight of the anti-affinity injection
type PreferredDuringSchedulingIgnoredDuringExecution struct {
	// Weight associated with matching the corresponding podAffinityTerm, in the range 1-100.
	Weight int32 `json:"weight" protobuf:"varint,1,opt,name=weight"`
}

// RequiredDuringSchedulingIgnoredDuringExecution defines inter-pod scheduling rule to be RequiredDuringSchedulingIgnoredDuringExecution
type RequiredDuringSchedulingIgnoredDuringExecution struct{}

// CanaryStrategy defines parameters for a Replica Based Canary
type CanaryStrategy struct {
	// CanaryService holds the name of a service which selects pods with canary version and don't select any pods with stable version.
	// +optional
	CanaryService string `json:"canaryService,omitempty" protobuf:"bytes,1,opt,name=canaryService"`
	// StableService holds the name of a service which selects pods with stable version and don't select any pods with canary version.
	// +optional
	StableService string `json:"stableService,omitempty" protobuf:"bytes,2,opt,name=stableService"`
	// Steps define the order of phases to execute the canary deployment
	// +optional
	Steps []CanaryStep `json:"steps,omitempty" protobuf:"bytes,3,rep,name=steps"`
	// TrafficRouting hosts all the supported service meshes supported to enable more fine-grained traffic routing
	TrafficRouting *RolloutTrafficRouting `json:"trafficRouting,omitempty" protobuf:"bytes,4,opt,name=trafficRouting"`

	// MaxUnavailable The maximum number of pods that can be unavailable during the update.
	// Value can be an absolute number (ex: 5) or a percentage of total pods at the start of update (ex: 10%).
	// Absolute number is calculated from percentage by rounding down.
	// This can not be 0 if MaxSurge is 0.
	// By default, a fixed value of 25% is used.
	// Example: when this is set to 30%, the old RC can be scaled down by 30%
	// immediately when the rolling update starts. Once new pods are ready, old RC
	// can be scaled down further, followed by scaling up the new RC, ensuring
	// that at least 70% of original number of pods are available at all times
	// during the update.
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty" protobuf:"bytes,5,opt,name=maxUnavailable"`

	// MaxSurge The maximum number of pods that can be scheduled above the original number of
	// pods.
	// Value can be an absolute number (ex: 5) or a percentage of total pods at
	// the start of the update (ex: 10%). This can not be 0 if MaxUnavailable is 0.
	// Absolute number is calculated from percentage by rounding up.
	// By default, a value of 25% is used.
	// Example: when this is set to 30%, the new RC can be scaled up by 30%
	// immediately when the rolling update starts. Once old pods have been killed,
	// new RC can be scaled up further, ensuring that total number of pods running
	// at any time during the update is at most 130% of original pods.
	// +optional
	MaxSurge *intstr.IntOrString `json:"maxSurge,omitempty" protobuf:"bytes,6,opt,name=maxSurge"`
	// Analysis runs a separate analysisRun while all the steps execute. This is intended to be a continuous validation of the new ReplicaSet
	Analysis *RolloutAnalysisBackground `json:"analysis,omitempty" protobuf:"bytes,7,opt,name=analysis"`
	// AntiAffinity enables anti-affinity rules for Canary deployment
	// +optional
	AntiAffinity *AntiAffinity `json:"antiAffinity,omitempty" protobuf:"bytes,8,opt,name=antiAffinity"`
	// CanaryMetadata specify labels and annotations which will be attached to the canary pods for
	// the duration which they act as a canary, and will be removed after
	CanaryMetadata *PodTemplateMetadata `json:"canaryMetadata,omitempty" protobuf:"bytes,9,opt,name=canaryMetadata"`
	// StableMetadata specify labels and annotations which will be attached to the stable pods for
	// the duration which they act as a canary, and will be removed after
	StableMetadata *PodTemplateMetadata `json:"stableMetadata,omitempty" protobuf:"bytes,10,opt,name=stableMetadata"`

	// ScaleDownDelaySeconds adds a delay before scaling down the previous ReplicaSet when the
	// canary strategy is used with traffic routing (default 30 seconds). A delay in scaling down
	// the previous ReplicaSet is needed after switching the stable service selector to point to
	// the new ReplicaSet, in order to give time for traffic providers to re-target the new pods.
	// This value is ignored with basic, replica-weighted canary without traffic routing.
	// +optional
	ScaleDownDelaySeconds *int32 `json:"scaleDownDelaySeconds,omitempty" protobuf:"varint,11,opt,name=scaleDownDelaySeconds"`
	// ScaleDownDelayRevisionLimit limits the number of old RS that can run at one time before getting scaled down
	// +optional
	ScaleDownDelayRevisionLimit *int32 `json:"scaleDownDelayRevisionLimit,omitempty" protobuf:"varint,12,opt,name=scaleDownDelayRevisionLimit"`
}

// ALBTrafficRouting configuration for ALB ingress controller to control traffic routing
type ALBTrafficRouting struct {
	// Ingress refers to the name of an `Ingress` resource in the same namespace as the `Rollout`
	Ingress string `json:"ingress" protobuf:"bytes,1,opt,name=ingress"`
	// ServicePort refers to the port that the Ingress action should route traffic to
	ServicePort int32 `json:"servicePort" protobuf:"varint,2,opt,name=servicePort"`
	// RootService references the service in the ingress to the controller should add the action to
	RootService string `json:"rootService,omitempty" protobuf:"bytes,3,opt,name=rootService"`
	// AnnotationPrefix has to match the configured annotation prefix on the alb ingress controller
	// +optional
	AnnotationPrefix string `json:"annotationPrefix,omitempty" protobuf:"bytes,4,opt,name=annotationPrefix"`
}

// RolloutTrafficRouting hosts all the different configuration for supported service meshes to enable more fine-grained traffic routing
type RolloutTrafficRouting struct {
	// Istio holds Istio specific configuration to route traffic
	Istio *IstioTrafficRouting `json:"istio,omitempty" protobuf:"bytes,1,opt,name=istio"`
	// Nginx holds Nginx Ingress specific configuration to route traffic
	Nginx *NginxTrafficRouting `json:"nginx,omitempty" protobuf:"bytes,2,opt,name=nginx"`
	// Nginx holds ALB Ingress specific configuration to route traffic
	ALB *ALBTrafficRouting `json:"alb,omitempty" protobuf:"bytes,3,opt,name=alb"`
	// SMI holds TrafficSplit specific configuration to route traffic
	SMI *SMITrafficRouting `json:"smi,omitempty" protobuf:"bytes,4,opt,name=smi"`
	// Ambassador holds specific configuration to use Ambassador to route traffic
	Ambassador *AmbassadorTrafficRouting `json:"ambassador,omitempty" protobuf:"bytes,5,opt,name=ambassador"`
}

// AmbassadorTrafficRouting defines the configuration required to use Ambassador as traffic
// router
type AmbassadorTrafficRouting struct {
	// Mappings refer to the name of the Ambassador Mappings used to route traffic to the
	// service
	Mappings []string `json:"mappings" protobuf:"bytes,1,rep,name=mappings"`
}

// SMITrafficRouting configuration for TrafficSplit Custom Resource to control traffic routing
type SMITrafficRouting struct {
	// RootService holds the name of that clients use to communicate.
	// +optional
	RootService string `json:"rootService,omitempty" protobuf:"bytes,1,opt,name=rootService"`
	// TrafficSplitName holds the name of the TrafficSplit.
	// +optional
	TrafficSplitName string `json:"trafficSplitName,omitempty" protobuf:"bytes,2,opt,name=trafficSplitName"`
}

// NginxTrafficRouting configuration for Nginx ingress controller to control traffic routing
type NginxTrafficRouting struct {
	// AnnotationPrefix has to match the configured annotation prefix on the nginx ingress controller
	// +optional
	AnnotationPrefix string `json:"annotationPrefix,omitempty" protobuf:"bytes,1,opt,name=annotationPrefix"`
	// StableIngress refers to the name of an `Ingress` resource in the same namespace as the `Rollout`
	StableIngress string `json:"stableIngress" protobuf:"bytes,2,opt,name=stableIngress"`
	// +optional
	AdditionalIngressAnnotations map[string]string `json:"additionalIngressAnnotations,omitempty" protobuf:"bytes,3,rep,name=additionalIngressAnnotations"`
}

// IstioTrafficRouting configuration for Istio service mesh to enable fine grain configuration
type IstioTrafficRouting struct {
	// VirtualService references an Istio VirtualService to modify to shape traffic
	VirtualService IstioVirtualService `json:"virtualService" protobuf:"bytes,1,opt,name=virtualService"`
	// DestinationRule references an Istio DestinationRule to modify to shape traffic
	DestinationRule *IstioDestinationRule `json:"destinationRule,omitempty" protobuf:"bytes,2,opt,name=destinationRule"`
}

// IstioVirtualService holds information on the virtual service the rollout needs to modify
type IstioVirtualService struct {
	// Name holds the name of the VirtualService
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Routes list of routes within VirtualService to edit
	Routes []string `json:"routes" protobuf:"bytes,2,rep,name=routes"`
}

// IstioDestinationRule is a reference to an Istio DestinationRule to modify and shape traffic
type IstioDestinationRule struct {
	// Name holds the name of the DestinationRule
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// CanarySubsetName is the subset name to modify labels with canary ReplicaSet pod template hash value
	CanarySubsetName string `json:"canarySubsetName" protobuf:"bytes,2,opt,name=canarySubsetName"`
	// StableSubsetName is the subset name to modify labels with stable ReplicaSet pod template hash value
	StableSubsetName string `json:"stableSubsetName" protobuf:"bytes,3,opt,name=stableSubsetName"`
}

// RolloutExperimentStep defines a template that is used to create a experiment for a step
type RolloutExperimentStep struct {
	// Templates what templates that should be added to the experiment. Should be non-nil
	// +patchMergeKey=name
	// +patchStrategy=merge
	Templates []RolloutExperimentTemplate `json:"templates" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,1,rep,name=templates"`
	// Duration is a duration string (e.g. 30s, 5m, 1h) that the experiment should run for
	// +optional
	Duration DurationString `json:"duration,omitempty" protobuf:"bytes,2,opt,name=duration,casttype=DurationString"`
	// Analyzes reference which analysis templates to run with the experiment
	// +patchMergeKey=name
	// +patchStrategy=merge
	Analyzes []RolloutExperimentStepAnalysisTemplateRef `json:"analyses,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,3,rep,name=analyses"` //nolint
}

type RolloutExperimentStepAnalysisTemplateRef struct {
	// Name is a name for this analysis template invocation
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// TemplateName reference of the AnalysisTemplate name used by the Experiment to create the run
	TemplateName string `json:"templateName" protobuf:"bytes,2,opt,name=templateName"`
	// Whether to look for the templateName at cluster scope or namespace scope
	// +optional
	ClusterScope bool `json:"clusterScope,omitempty" protobuf:"varint,3,opt,name=clusterScope"`
	// Args the arguments that will be added to the AnalysisRuns
	// +patchMergeKey=name
	// +patchStrategy=merge
	Args []AnalysisRunArgument `json:"args,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,4,rep,name=args"`
	// RequiredForCompletion blocks the Experiment from completing until the analysis has completed
	RequiredForCompletion bool `json:"requiredForCompletion,omitempty" protobuf:"varint,5,opt,name=requiredForCompletion"`
}

// RolloutExperimentTemplate defines the template used to create experiments for the Rollout's experiment canary step
type RolloutExperimentTemplate struct {
	// Name description of template that passed to the template
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// SpecRef indicates where the rollout should get the RS template from
	SpecRef ReplicaSetSpecRef `json:"specRef" protobuf:"bytes,2,opt,name=specRef,casttype=ReplicaSetSpecRef"`
	// Replicas replica count for the template
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,3,opt,name=replicas"`
	// Metadata sets labels and annotations to use for the RS created from the template
	// +optional
	Metadata PodTemplateMetadata `json:"metadata,omitempty" protobuf:"bytes,4,opt,name=metadata"`
	// Selector overrides the selector to be used for the template's ReplicaSet. If omitted, will
	// use the same selector as the Rollout
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty" protobuf:"bytes,5,opt,name=selector"`
}

// PodTemplateMetadata extra labels to add to the template
type PodTemplateMetadata struct {
	// Labels Additional labels to add to the experiment
	// +optional
	Labels map[string]string `json:"labels,omitempty" protobuf:"bytes,1,rep,name=labels"`
	// Annotations additional annotations to add to the experiment
	// +optional
	Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,2,rep,name=annotations"`
}

// ReplicaSetSpecRef defines which RS that the experiment's template will use.
type ReplicaSetSpecRef string

const (
	// CanarySpecRef indicates the RS template should be pulled from the newRS's template
	CanarySpecRef ReplicaSetSpecRef = "canary"
	// StableSpecRef indicates the RS template should be pulled from the stableRS's template
	StableSpecRef ReplicaSetSpecRef = "stable"
)

// CanaryStep defines a step of a canary deployment.
type CanaryStep struct {
	// SetWeight sets what percentage of the newRS should receive
	SetWeight *int32 `json:"setWeight,omitempty" protobuf:"varint,1,opt,name=setWeight"`
	// Pause freezes the rollout by setting spec.Paused to true.
	// A Rollout will resume when spec.Paused is reset to false.
	// +optional
	Pause *RolloutPause `json:"pause,omitempty" protobuf:"bytes,2,opt,name=pause"`
	// Experiment defines the experiment object that should be created
	Experiment *RolloutExperimentStep `json:"experiment,omitempty" protobuf:"bytes,3,opt,name=experiment"`
	// Analysis defines the AnalysisRun that will run for a step
	Analysis *RolloutAnalysis `json:"analysis,omitempty" protobuf:"bytes,4,opt,name=analysis"`
	// SetCanaryScale defines how to scale the newRS without changing traffic weight
	// +optional
	SetCanaryScale *SetCanaryScale `json:"setCanaryScale,omitempty" protobuf:"bytes,5,opt,name=setCanaryScale"`
}

// CanaryStepString returns a string representation of a canary step
func CanaryStepString(c CanaryStep) string {
	if c.SetWeight != nil {
		return fmt.Sprintf("setWeight: %d", *c.SetWeight)
	}
	if c.Pause != nil {
		str := "pause"
		if c.Pause.Duration != nil {
			str = fmt.Sprintf("%s: %s", str, c.Pause.Duration.String())
		}
		return str
	}
	if c.Experiment != nil {
		return "experiment"
	}
	if c.Analysis != nil {
		return "analysis"
	}
	if c.SetCanaryScale != nil {
		if c.SetCanaryScale.Weight != nil {
			return fmt.Sprintf("setCanaryScale{weight: %d}", c.SetCanaryScale.Weight)
		} else if c.SetCanaryScale.MatchTrafficWeight {
			return "setCanaryScale{matchTrafficWeight: true}"
		} else if c.SetCanaryScale.Replicas != nil {
			return fmt.Sprintf("setCanaryScale{replicas: %d}", *c.SetCanaryScale.Replicas)
		}
	}
	return "invalid"
}

// SetCanaryScale defines how to scale the newRS without changing traffic weight
type SetCanaryScale struct {
	// Weight sets the percentage of replicas the newRS should have
	// +optional
	Weight *int32 `json:"weight,omitempty" protobuf:"varint,1,opt,name=weight"`
	// Replicas sets the number of replicas the newRS should have
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,2,opt,name=replicas"`
	// MatchTrafficWeight cancels out previously set Replicas or Weight, effectively activating SetWeight
	// +optional
	MatchTrafficWeight bool `json:"matchTrafficWeight,omitempty" protobuf:"varint,3,opt,name=matchTrafficWeight"`
}

// RolloutAnalysisBackground defines a template that is used to create a background analysisRun
type RolloutAnalysisBackground struct {
	RolloutAnalysis `json:",inline" protobuf:"bytes,1,opt,name=rolloutAnalysis"`
	// StartingStep indicates which step the background analysis should start on
	// If not listed, controller defaults to 0
	StartingStep *int32 `json:"startingStep,omitempty" protobuf:"varint,2,opt,name=startingStep"`
}

// RolloutAnalysis defines a template that is used to create a analysisRun
type RolloutAnalysis struct {
	//Templates reference to a list of analysis templates to combine for an AnalysisRun
	Templates []RolloutAnalysisTemplate `json:"templates,omitempty" protobuf:"bytes,1,rep,name=templates"`
	// Args the arguments that will be added to the AnalysisRuns
	// +patchMergeKey=name
	// +patchStrategy=merge
	Args []AnalysisRunArgument `json:"args,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=args"`
}

type RolloutAnalysisTemplate struct {
	//TemplateName name of template to use in AnalysisRun
	// +optional
	TemplateName string `json:"templateName" protobuf:"bytes,1,opt,name=templateName"`
	// Whether to look for the templateName at cluster scope or namespace scope
	// +optional
	ClusterScope bool `json:"clusterScope,omitempty" protobuf:"varint,2,opt,name=clusterScope"`
}

// AnalysisRunArgument argument to add to analysisRun
type AnalysisRunArgument struct {
	// Name argument name
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Value a hardcoded value for the argument. This field is a one of field with valueFrom
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
	// ValueFrom A reference to where the value is stored. This field is a one of field with valueFrom
	ValueFrom *ArgumentValueFrom `json:"valueFrom,omitempty" protobuf:"bytes,3,opt,name=valueFrom"`
}

// ArgumentValueFrom defines references to fields within resources to grab for the value (i.e. Pod Template Hash)
type ArgumentValueFrom struct {
	// PodTemplateHashValue gets the value from one of the children ReplicaSet's Pod Template Hash
	PodTemplateHashValue *ValueFromPodTemplateHash `json:"podTemplateHashValue,omitempty" protobuf:"bytes,1,opt,name=podTemplateHashValue,casttype=ValueFromPodTemplateHash"`
	//FieldRef
	FieldRef *FieldRef `json:"fieldRef,omitempty" protobuf:"bytes,2,opt,name=fieldRef"`
}

type FieldRef struct {
	// Required: Path of the field to select in the specified API version
	FieldPath string `json:"fieldPath" protobuf:"bytes,1,opt,name=fieldPath"`
}

// ValueFromPodTemplateHash indicates which ReplicaSet pod template pod hash to use
type ValueFromPodTemplateHash string

const (
	// Stable tells the Rollout to get the pod template hash from the stable ReplicaSet
	Stable ValueFromPodTemplateHash = "Stable"
	// Latest tells the Rollout to get the pod template hash from the latest ReplicaSet
	Latest ValueFromPodTemplateHash = "Latest"
)

const (
	// RolloutTypeLabel indicates how the rollout created the analysisRun
	RolloutTypeLabel = "rollout-type"
	// RolloutTypeStepLabel indicates that the analysisRun was created as a canary step
	RolloutTypeStepLabel = "Step"
	// RolloutTypeBackgroundRunLabel indicates that the analysisRun was created in Background to an execution
	RolloutTypeBackgroundRunLabel = "Background"
	// RolloutTypePrePromotionLabel indicates that the analysisRun was created before the active service promotion
	RolloutTypePrePromotionLabel = "PrePromotion"
	// RolloutTypePostPromotionLabel indicates that the analysisRun was created after the active service promotion
	RolloutTypePostPromotionLabel = "PostPromotion"
	// RolloutCanaryStepIndexLabel indicates which step created this analysisRun
	RolloutCanaryStepIndexLabel = "step-index"
)

// RolloutPause defines a pause stage for a rollout
type RolloutPause struct {
	// Duration the amount of time to wait before moving to the next step.
	// +optional
	Duration *intstr.IntOrString `json:"duration,omitempty" protobuf:"bytes,1,opt,name=duration"`
}

// DurationSeconds converts the pause duration to seconds
// If Duration is nil 0 is returned
// if Duration values is string and does not contain a valid unit -1 is returned
func (p RolloutPause) DurationSeconds() int32 {
	if p.Duration != nil {
		if p.Duration.Type == intstr.String {
			s, err := strconv.ParseInt(p.Duration.StrVal, 10, 32)
			if err != nil {
				d, err := time.ParseDuration(p.Duration.StrVal)
				if err != nil {
					return -1
				}
				return int32(d.Seconds())
			}
			// special case where no unit was specified
			return int32(s)
		}
		return p.Duration.IntVal
	}
	return 0
}

// DurationFromInt creates duration in seconds from int value
func DurationFromInt(i int) *intstr.IntOrString {
	d := intstr.FromInt(i)
	return &d
}

// DurationFromString creates duration from string
// value must be a string representation of an int with optional time unit (see time.ParseDuration)
func DurationFromString(s string) *intstr.IntOrString {
	d := intstr.FromString(s)
	return &d
}

// PauseReason reasons that the rollout can pause
type PauseReason string

const (
	// PauseReasonInconclusiveAnalysis pauses rollout when rollout has an inconclusive analysis run
	PauseReasonInconclusiveAnalysis PauseReason = "InconclusiveAnalysisRun"
	// PauseReasonInconclusiveExperiment pauses rollout when rollout has an inconclusive experiment
	PauseReasonInconclusiveExperiment PauseReason = "InconclusiveExperiment"
	// PauseReasonCanaryPauseStep pause rollout for canary pause step
	PauseReasonCanaryPauseStep PauseReason = "CanaryPauseStep"
	// PauseReasonBlueGreenPause pause rollout before promoting rollout
	PauseReasonBlueGreenPause PauseReason = "BlueGreenPause"
)

// PauseCondition the reason for a pause and when it started
type PauseCondition struct {
	Reason    PauseReason `json:"reason" protobuf:"bytes,1,opt,name=reason,casttype=PauseReason"`
	StartTime metav1.Time `json:"startTime" protobuf:"bytes,2,opt,name=startTime"`
}

// RolloutStatus is the status for a Rollout resource
type RolloutStatus struct {
	// Abort cancel the current rollout progression
	Abort bool `json:"abort,omitempty" protobuf:"varint,1,opt,name=abort"`
	// PauseConditions indicates why the rollout is currently paused
	PauseConditions []PauseCondition `json:"pauseConditions,omitempty" protobuf:"bytes,2,rep,name=pauseConditions"`
	// ControllerPause indicates the controller has paused the rollout. It is set to true when
	// the controller adds a pause condition. This field helps to discern the scenario where a
	// rollout was resumed after being paused by the controller (e.g. via the plugin). In that
	// situation, the pauseConditions would have been cleared , but controllerPause would still be
	// set to true.
	ControllerPause bool `json:"controllerPause,omitempty" protobuf:"varint,3,opt,name=controllerPause"`
	// AbortedAt indicates the controller reconciled an aborted rollout. The controller uses this to understand if
	// the controller needs to do some specific work when a Rollout is aborted. For example, the reconcileAbort is used
	// to indicate if the Rollout should enter an aborted state when the latest AnalysisRun is a failure, or the controller
	// has already put the Rollout into an aborted and should create a new AnalysisRun.
	AbortedAt *metav1.Time `json:"abortedAt,omitempty" protobuf:"bytes,4,opt,name=abortedAt"`
	// CurrentPodHash the hash of the current pod template
	// +optional
	CurrentPodHash string `json:"currentPodHash,omitempty" protobuf:"bytes,5,opt,name=currentPodHash"`
	// CurrentStepHash the hash of the current list of steps for the current strategy. This is used to detect when the
	// list of current steps change
	// +optional
	CurrentStepHash string `json:"currentStepHash,omitempty" protobuf:"bytes,6,opt,name=currentStepHash"`
	// Total number of non-terminated pods targeted by this rollout (their labels match the selector).
	// +optional
	Replicas int32 `json:"replicas,omitempty" protobuf:"varint,7,opt,name=replicas"`
	// Total number of non-terminated pods targeted by this rollout that have the desired template spec.
	// +optional
	UpdatedReplicas int32 `json:"updatedReplicas,omitempty" protobuf:"varint,8,opt,name=updatedReplicas"`
	// Total number of ready pods targeted by this rollout.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty" protobuf:"varint,9,opt,name=readyReplicas"`
	// Total number of available pods (ready for at least minReadySeconds) targeted by this rollout.
	// +optional
	AvailableReplicas int32 `json:"availableReplicas,omitempty" protobuf:"varint,10,opt,name=availableReplicas"`
	// CurrentStepIndex defines the current step of the rollout is on. If the current step index is null, the
	// controller will execute the rollout.
	// +optional
	CurrentStepIndex *int32 `json:"currentStepIndex,omitempty" protobuf:"varint,11,opt,name=currentStepIndex"`
	// Count of hash collisions for the Rollout. The Rollout controller uses this
	// field as a collision avoidance mechanism when it needs to create the name for the
	// newest ReplicaSet.
	// +optional
	CollisionCount *int32 `json:"collisionCount,omitempty" protobuf:"varint,12,opt,name=collisionCount"`
	// The generation observed by the rollout controller by taking a hash of the spec.
	// +optional
	ObservedGeneration string `json:"observedGeneration,omitempty" protobuf:"bytes,13,opt,name=observedGeneration"`
	// Conditions a list of conditions a rollout can have.
	// +optional
	Conditions []RolloutCondition `json:"conditions,omitempty" protobuf:"bytes,14,rep,name=conditions"`
	// Canary describes the state of the canary rollout
	// +optional
	Canary CanaryStatus `json:"canary,omitempty" protobuf:"bytes,15,opt,name=canary"`
	// BlueGreen describes the state of the bluegreen rollout
	// +optional
	BlueGreen BlueGreenStatus `json:"blueGreen,omitempty" protobuf:"bytes,16,opt,name=blueGreen"`
	// HPAReplicas the number of non-terminated replicas that are receiving active traffic
	// +optional
	HPAReplicas int32 `json:"HPAReplicas,omitempty" protobuf:"varint,17,opt,name=HPAReplicas"`
	// Selector that identifies the pods that are receiving active traffic
	// +optional
	Selector string `json:"selector,omitempty" protobuf:"bytes,18,opt,name=selector"`
	// StableRS indicates the replicaset that has successfully rolled out
	// +optional
	StableRS string `json:"stableRS,omitempty" protobuf:"bytes,19,opt,name=stableRS"`
	// RestartedAt indicates last time a Rollout was restarted
	RestartedAt *metav1.Time `json:"restartedAt,omitempty" protobuf:"bytes,20,opt,name=restartedAt"`
	// PromoteFull indicates if the rollout should perform a full promotion, skipping analysis and pauses.
	PromoteFull bool `json:"promoteFull,omitempty" protobuf:"varint,21,opt,name=promoteFull"`
}

// BlueGreenStatus status fields that only pertain to the blueGreen rollout
type BlueGreenStatus struct {
	// PreviewSelector indicates which replicas set the preview service is serving traffic to
	// +optional
	PreviewSelector string `json:"previewSelector,omitempty" protobuf:"bytes,1,opt,name=previewSelector"`
	// ActiveSelector indicates which replicas set the active service is serving traffic to
	// +optional
	ActiveSelector string `json:"activeSelector,omitempty" protobuf:"bytes,2,opt,name=activeSelector"`
	// ScaleUpPreviewCheckPoint indicates that the Replicaset receiving traffic from the preview service is ready to be scaled up after the rollout is unpaused
	// +optional
	ScaleUpPreviewCheckPoint bool `json:"scaleUpPreviewCheckPoint,omitempty" protobuf:"varint,3,opt,name=scaleUpPreviewCheckPoint"`
	// PrePromotionAnalysisRunStatus indicates the status of the current prepromotion analysis run
	PrePromotionAnalysisRunStatus *RolloutAnalysisRunStatus `json:"prePromotionAnalysisRunStatus,omitempty" protobuf:"bytes,4,opt,name=prePromotionAnalysisRunStatus"`
	// PostPromotionAnalysisRunStatus indicates the status of the current post promotion analysis run
	PostPromotionAnalysisRunStatus *RolloutAnalysisRunStatus `json:"postPromotionAnalysisRunStatus,omitempty" protobuf:"bytes,5,opt,name=postPromotionAnalysisRunStatus"`
}

// CanaryStatus status fields that only pertain to the canary rollout
type CanaryStatus struct {
	// CurrentStepAnalysisRunStatus indicates the status of the current step analysis run
	CurrentStepAnalysisRunStatus *RolloutAnalysisRunStatus `json:"currentStepAnalysisRunStatus,omitempty" protobuf:"bytes,1,opt,name=currentStepAnalysisRunStatus"`
	// CurrentBackgroundAnalysisRunStatus indicates the status of the current background analysis run
	CurrentBackgroundAnalysisRunStatus *RolloutAnalysisRunStatus `json:"currentBackgroundAnalysisRunStatus,omitempty" protobuf:"bytes,2,opt,name=currentBackgroundAnalysisRunStatus"`
	// CurrentExperiment indicates the running experiment
	CurrentExperiment string `json:"currentExperiment,omitempty" protobuf:"bytes,3,opt,name=currentExperiment"`
}

type RolloutAnalysisRunStatus struct {
	Name    string        `json:"name" protobuf:"bytes,1,opt,name=name"`
	Status  AnalysisPhase `json:"status" protobuf:"bytes,2,opt,name=status,casttype=AnalysisPhase"`
	Message string        `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`
}

// RolloutConditionType defines the conditions of Rollout
type RolloutConditionType string

// These are valid conditions of a rollout.
const (
	// InvalidSpec means the rollout has an invalid spec and will not progress until
	// the spec is fixed.
	InvalidSpec RolloutConditionType = "InvalidSpec"
	// RolloutAvailable means the rollout is available, ie. the active service is pointing at a
	// replicaset with the required replicas up and running for at least minReadySeconds.
	RolloutAvailable RolloutConditionType = "Available"
	// RolloutProgressing means the rollout is progressing. Progress for a rollout is
	// considered when a new replica set is created or adopted, when pods scale
	// up or old pods scale down, or when the services are updated. Progress is not estimated
	// for paused rollouts.
	RolloutProgressing RolloutConditionType = "Progressing"
	// RolloutReplicaFailure ReplicaFailure is added in a deployment when one of its pods
	// fails to be created or deleted.
	RolloutReplicaFailure RolloutConditionType = "ReplicaFailure"
	// RolloutPaused means that rollout is in a paused state. It is still progressing at this point.
	RolloutPaused RolloutConditionType = "Paused"
	// RolloutCompleted means that rollout is in a completed state. It is still progressing at this point.
	RolloutCompleted RolloutConditionType = "Completed"
)

// RolloutCondition describes the state of a rollout at a certain point.
type RolloutCondition struct {
	// Type of deployment condition.
	Type RolloutConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=RolloutConditionType"`
	// Phase of the condition, one of True, False, Unknown.
	Status corev1.ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	// The last time this condition was updated.
	LastUpdateTime metav1.Time `json:"lastUpdateTime" protobuf:"bytes,3,opt,name=lastUpdateTime"`
	// Last time the condition transitioned from one status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	Reason string `json:"reason" protobuf:"bytes,5,opt,name=reason"`
	// A human readable message indicating details about the transition.
	Message string `json:"message" protobuf:"bytes,6,opt,name=message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RolloutList is a list of Rollout resources
type RolloutList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Items []Rollout `json:"items" protobuf:"bytes,2,rep,name=items"`
}
