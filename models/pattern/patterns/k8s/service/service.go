package service

import (
	"encoding/json"
	"fmt"

	patternUtils "github.com/layer5io/meshery/models/pattern/utils"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
	v1 "k8s.io/api/core/v1"
)

func Deploy(kubeClient *meshkube.Client, oamComp v1alpha1.Component, oamConfig v1alpha1.Configuration, isDel bool) error {
	svc, err := createSvcFromComp(oamComp)
	if err != nil {
		return err
	}

	if isDel {
		return patternUtils.DeleteK8sResource(
			kubeClient.DynamicKubeClient,
			"",
			"v1",
			"services",
			svc.Namespace,
			svc.Name,
		)
	}

	return patternUtils.CreateK8sResource(
		kubeClient.DynamicKubeClient,
		"",
		"v1",
		"services",
		svc,
		false,
	)
}

func createSvcFromComp(oamComp v1alpha1.Component) (v1.Service, error) {
	svc := v1.Service{}

	svc.SetName(oamComp.Name)
	svc.SetNamespace(oamComp.Namespace)

	byt, err := json.Marshal(oamComp.Spec.Settings)
	if err != nil {
		return svc, fmt.Errorf("failed to construct service from the settings")
	}

	if err := json.Unmarshal(byt, &svc); err != nil {
		return svc, fmt.Errorf("failed to construct service from the settings")
	}

	return svc, nil
}
