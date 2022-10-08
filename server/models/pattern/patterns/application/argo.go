package application

import (
	"fmt"

	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/models/pattern/patterns/application/argo/v1alpha1"
	patternUtils "github.com/layer5io/meshery/server/models/pattern/utils"
	"github.com/layer5io/meshkit/utils"
	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const argoManifest = "https://raw.githubusercontent.com/argoproj/argo-rollouts/master/manifests/install.yaml"

type ArgoRollout struct {
	kubeclient *meshkube.Client
	// gRPCClient meshes.MeshServiceClient
}

// Install installs argo rollout in the given namespace
func (ar *ArgoRollout) Install() error {
	manifestContent, err := utils.ReadFileSource(argoManifest)
	if err != nil {
		return err
	}

	return ar.kubeclient.ApplyManifest([]byte(manifestContent), meshkube.ApplyOptions{
		Namespace: "argo-rollouts",
	})
}

// Native will do simple rollout
func (ar *ArgoRollout) Native(opt RolloutEngineGenericOptions) error {
	if opt.Delete {
		// Delete the rollout
		if err := patternUtils.DeleteK8sResource(
			ar.kubeclient.DynamicKubeClient,
			"argoproj.io",
			"v1alpha1",
			"rollouts",
			opt.Namespace,
			opt.Name,
		); err != nil {
			logrus.Error("[Native Rollout]:", err)
			return fmt.Errorf("failed to delete resource with name \"%s\" in namespace \"%s\"", opt.Name, opt.Namespace)
		}

		if !*opt.Advanced.CreateService {
			return nil
		}

		// Delete the service
		if err := patternUtils.DeleteK8sResource(
			ar.kubeclient.DynamicKubeClient,
			"",
			"v1",
			"services",
			opt.Namespace,
			opt.Name,
		); err != nil {
			logrus.Error("[Native Rollout]:", err)
			return fmt.Errorf("failed to delete service with name \"%s\" in namespace \"%s\"", opt.Name, opt.Namespace)
		}

		return nil
	}

	// Create rollout === deployment
	if err := patternUtils.CreateK8sResource(
		ar.kubeclient.DynamicKubeClient,
		"argoproj.io",
		"v1alpha1",
		"rollouts",
		createNativeArgoResource(opt),
		false,
	); err != nil {
		logrus.Error("[Native Rollout]:", err)
		return fmt.Errorf("failed to create resource with name \"%s\" in namespace \"%s\"", opt.Name, opt.Namespace)
	}

	if !*opt.Advanced.CreateService {
		return nil
	}

	// Create service
	if err := patternUtils.CreateK8sResource(
		ar.kubeclient.DynamicKubeClient,
		"",
		"v1",
		"services",
		createNativeService(opt),
		false,
	); err != nil {
		logrus.Error("[Native Rollout]:", err)
		return fmt.Errorf("failed to create service with name \"%s\" in namespace \"%s\"", opt.Name, opt.Namespace)
	}

	return nil
}

// Canary performs canary rollout of the given application
func (ar *ArgoRollout) Canary(opt RolloutEngineCanaryOptions) error {
	if opt.ServiceMesh == "istio" {
		return ar.istioCanary(opt)
	}

	return nil
}

func (ar *ArgoRollout) istioCanary(opt RolloutEngineCanaryOptions) error {
	// Performing canary requires 4 steps:
	// 1. Create a rollout
	// ar.createOrUpdateRollout(opt)
	// 2. Create a service for stable
	// ar.createOrUpdateService(opt.Namespace, &v1.Service{})
	// 3. Create a service for the canary
	// ar.createOrUpdateService(opt.Namespace, &v1.Service{})
	// 4. Create a virtual service for routing
	// invoke remote operation
	return nil
}

func createIstioCanaryArgoResource(opt RolloutEngineCanaryOptions) v1alpha1.Rollout {
	rollout := createNativeArgoResource(opt.RolloutEngineGenericOptions)

	// TODO: Add canary specific information and then return
	return rollout
}

func createNativeArgoResource(opt RolloutEngineGenericOptions) v1alpha1.Rollout {
	replicas := int32(opt.Replicas)
	revisionHistoryLimit := int32(2)
	containers := []v1.Container{}
	weight := int32(100)

	for _, container := range opt.Containers {
		ports := []v1.ContainerPort{}
		for _, port := range container.Ports {
			ports = append(ports, v1.ContainerPort{
				Name:          port.Name,
				ContainerPort: int32(port.ContainerPort),
			})
		}

		containers = append(containers, v1.Container{
			Name:    container.Name,
			Image:   container.Image,
			Ports:   ports,
			Env:     container.Envs,
			Command: container.Commands,
		})
	}

	return v1alpha1.Rollout{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "argoproj.io/v1alpha1",
			Kind:       "Rollout",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      opt.Name,
			Namespace: opt.Namespace,
			Labels: helpers.MergeStringMaps(
				map[string]string{"app": opt.Name},
				opt.Advanced.Labels,
				opt.Metadata.Labels,
			),
			Annotations: helpers.MergeStringMaps(opt.Metadata.Annotations, opt.Metadata.Annotations),
		},
		Spec: v1alpha1.RolloutSpec{
			Replicas:             &replicas,
			RevisionHistoryLimit: &revisionHistoryLimit,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": opt.Name,
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: helpers.MergeStringMaps(
						map[string]string{
							"app": opt.Name,
						},
						opt.Advanced.Labels,
						opt.Metadata.Labels,
					),
				},
				Spec: v1.PodSpec{
					Containers: containers,
				},
			},
			Strategy: v1alpha1.RolloutStrategy{
				Canary: &v1alpha1.CanaryStrategy{
					Steps: []v1alpha1.CanaryStep{{SetWeight: &weight}},
				},
			},
		},
	}
}

func createNativeService(opt RolloutEngineGenericOptions) v1.Service {
	ports := []v1.ServicePort{}

	for _, container := range opt.Containers {
		for _, port := range container.Ports {
			ports = append(ports, v1.ServicePort{
				Port: int32(port.ContainerPort),
				Name: port.Name,
			})
		}
	}

	return v1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: opt.Name,
			Annotations: helpers.MergeStringMaps(map[string]string{
				"meshery-engine": "argo",
			}, opt.Metadata.Annotations),
			Labels:    opt.Metadata.Labels,
			Namespace: opt.Namespace,
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"app": opt.Name,
			},
			Ports: ports,
			Type:  v1.ServiceTypeLoadBalancer, // Should get smarter
		},
	}
}
