package application

// import (
// 	"fmt"

// 	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
// 	v1 "k8s.io/api/core/v1"
// )

// type RolloutEngineGenericOptions struct {
// 	Name        string
// 	Namespace   string
// 	ServiceMesh string
// 	Metadata    RolloutEngineGenericOptionsMetadata
// 	Replicas    int
// 	Containers  []RolloutEngineContainer
// 	Delete      bool
// 	Advanced    PatternSettingAdvanced
// }

// type RolloutEngineGenericOptionsMetadata struct {
// 	Labels      map[string]string
// 	Annotations map[string]string
// }

// type RolloutEngineContainer struct {
// 	Name      string
// 	Image     string
// 	Commands  []string
// 	Ports     []RolloutEngineContainerPort
// 	Resources []RolloutEngineContainerResource
// 	Envs      []v1.EnvVar
// }

// type RolloutEngineContainerPort struct {
// 	Name          string
// 	ContainerPort int
// }

// type RolloutEngineContainerResource struct {
// 	v1.ResourceRequirements
// }

// type RolloutEngineCanaryOptions struct {
// 	RolloutEngineGenericOptions
// 	Steps []RolloutEngineCanaryStep
// }

// type RolloutEngineCanaryStep struct {
// 	SetWeight int
// 	Steps     []RolloutEngineCanaryStepPause
// }

// type RolloutEngineCanaryStepPause struct {
// 	Duration int
// 	// Inline Metrics - for the future
// }

// type RolloutEngine interface {
// 	// Canary(RolloutEngineCanaryOptions) error
// 	Install() error
// 	Native(RolloutEngineGenericOptions) error
// }

// func NewRolloutEngine(kubeclient *meshkube.Client, typ string) (RolloutEngine, error) {
// 	if typ == "argo" {
// 		return &ArgoRollout{
// 			kubeclient: kubeclient,
// 		}, nil
// 	}

// 	return nil, fmt.Errorf("invalid rollout engine type: \"%s\"", typ)
// }
