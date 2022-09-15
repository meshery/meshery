package application

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type serviceMesh string

const (
	istio serviceMesh = "istio"
)

type rolloutEngine string

const (
	argo rolloutEngine = "argo"
)

type PatternSetting struct {
	Replicas   int                      `json:"replicas,omitempty"`
	Mesh       serviceMesh              `json:"mesh,omitempty"`
	Containers []RolloutEngineContainer `json:"containers,omitempty"`
	Advanced   PatternSettingAdvanced   `json:"advanced,omitempty"`
}

type PatternSettingAdvanced struct {
	CreateService *bool             `json:"create_service,omitempty"`
	Labels        map[string]string `json:"labels,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty"`
}

type PatternConfiguration struct {
	RolloutStrategy *PatternConfigurationRolloutStrategy `json:"rollout_strategy,omitempty"`
}

type PatternConfigurationRolloutStrategy struct {
}

// Deploy will deploy the application
func Deploy(
	kubeClient *meshkube.Client,
	oamComp v1alpha1.Component,
	oamConfig v1alpha1.Configuration,
	isDel bool,
) error {
	if oamComp.Spec.Type == "Application" {
		settings, err := getApplicationPatternSettings(oamComp)
		if err != nil {
			return err
		}
		config, err := getApplicationPatternConfiguration(oamComp.Name, oamConfig)
		if err != nil {
			return nil
		}

		mesh := settings.Mesh
		if mesh == "" {
			mesh = detectServiceMesh(kubeClient)
		}

		var engineName rolloutEngine
		detectedEngine := detectRolloutEngine(kubeClient)
		if detectedEngine == "" {
			engineName = argo
		} else {
			engineName = detectedEngine
		}

		engine, err := NewRolloutEngine(kubeClient, string(engineName))
		if err != nil {
			return err
		}

		// If no engine was detected then proceed to install the default one
		if detectedEngine == "" {
			if err := engine.Install(); err != nil {
				return err
			}
		}

		setupDefaults(&settings)

		if config.RolloutStrategy == nil {
			return engine.Native(RolloutEngineGenericOptions{
				Name:        oamComp.Name,
				Namespace:   oamComp.Namespace,
				ServiceMesh: string(mesh),
				Replicas:    settings.Replicas,
				Containers:  settings.Containers,
				Metadata: RolloutEngineGenericOptionsMetadata{
					Labels:      oamComp.Labels,
					Annotations: oamComp.Annotations,
				},
				Delete:   isDel,
				Advanced: settings.Advanced,
			})
		}

		return fmt.Errorf("strategy not supported")
	}

	return fmt.Errorf("%s is not an application pattern", oamComp.Spec.Type)
}

func getApplicationPatternSettings(oamComp v1alpha1.Component) (PatternSetting, error) {
	var settings PatternSetting

	jsonByt, err := json.Marshal(oamComp.Spec.Settings)
	if err != nil {
		return PatternSetting{}, err
	}

	if err := json.Unmarshal(jsonByt, &settings); err != nil {
		return settings, err
	}

	return settings, nil
}

func getApplicationPatternConfiguration(compName string, oamConfig v1alpha1.Configuration) (PatternConfiguration, error) {
	var config PatternConfiguration

	for _, cfg := range oamConfig.Spec.Components {
		if cfg.ComponentName == compName {
			// Find the "RolloutStrategy" trait
			for _, trait := range cfg.Traits {
				if trait.Name == "RolloutStrategy" {
					jsonByt, err := json.Marshal(trait.Properties)
					if err != nil {
						return config, err
					}

					if err := json.Unmarshal(jsonByt, &config); err != nil {
						return config, err
					}

					break
				}
			}
		}
	}

	return config, nil
}

// detectServiceMesh will detect available service mesh in the k8s cluster
func detectServiceMesh(kubeClient *meshkube.Client) serviceMesh {
	return ""
}

// detectRolloutEngine will detect available rollout engine in the k8s cluster
func detectRolloutEngine(kubeClient *meshkube.Client) rolloutEngine {
	if detectArgoRollout(kubeClient) {
		return argo
	}

	return ""
}

// detectArgoRollout returns true if argo rollout is present in the cluster
func detectArgoRollout(kubeClient *meshkube.Client) bool {
	deployments := []string{"argo-rollouts"}
	nsList := findNamespaces(kubeClient)

	for _, deployment := range deployments {
		for _, ns := range nsList {
			if _, err := kubeClient.
				DynamicKubeClient.
				Resource(schema.GroupVersionResource{
					Group:    "apps",
					Version:  "v1",
					Resource: "deployments",
				}).
				Namespace(ns).
				Get(context.TODO(), deployment, v1.GetOptions{}); err == nil {
				logrus.Debug("found argo rollout in the cluster")
				return true
			}

			logrus.Debugf("%s not found in namespace %s", deployment, ns)
		}
	}

	return false
}

func findNamespaces(kubeClient *meshkube.Client) []string {
	list, err := kubeClient.KubeClient.CoreV1().Namespaces().List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return []string{}
	}

	res := []string{}

	for _, ns := range list.Items {
		res = append(res, ns.Name)
	}

	return res
}

func setupDefaults(s *PatternSetting) {
	// If nothing is specified then set the value to true
	if s.Advanced.CreateService == nil {
		csvc := true
		s.Advanced.CreateService = &csvc
	}
}
