package patterns

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/layer5io/meshery/models/pattern/patterns/application"
	"github.com/layer5io/meshery/models/pattern/patterns/k8s"
	"github.com/layer5io/meshery/models/pattern/patterns/k8s/service"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	meshkube "github.com/layer5io/meshkit/utils/kubernetes"
)

func ProcessOAM(kubeClient *meshkube.Client, oamComps []string, oamConfig string, isDel bool) (string, error) {
	var comps []v1alpha1.Component
	var config v1alpha1.Configuration

	for _, oamComp := range oamComps {
		var comp v1alpha1.Component
		if err := json.Unmarshal([]byte(oamComp), &comp); err != nil {
			return "", err
		}

		comps = append(comps, comp)
	}

	if err := json.Unmarshal([]byte(oamConfig), &config); err != nil {
		return "", err
	}

	var msgs []string
	var errs []error

	for _, comp := range comps {
		if comp.Spec.Type == "Application" {
			if err := application.Deploy(kubeClient, comp, config, isDel); err != nil {
				errs = append(errs, err)
			}

			if !isDel {
				msgs = append(msgs, "successfully deployed application "+comp.Name)
			} else {
				msgs = append(msgs, "successfully deleted application "+comp.Name)
			}

			continue
		}

		if comp.Spec.Type == "Service.K8s" {
			if err := service.Deploy(kubeClient, comp, config, isDel); err != nil {
				errs = append(errs, err)
			}

			if !isDel {
				msgs = append(msgs, "successfully deployed service "+comp.Name)
			} else {
				msgs = append(msgs, "successfully deleted service "+comp.Name)
			}

			continue
		}

		// Handle all of the k8s component here
		if strings.HasSuffix(strings.ToLower(comp.Spec.Type), ".k8s") {
			if err := k8s.Deploy(kubeClient, comp, config, isDel); err != nil {
				errs = append(errs, err)

				if !isDel {
					msgs = append(msgs, "successfully deployed kubernetes component "+comp.Name)
				} else {
					msgs = append(msgs, "successfully deleted kubernetes component "+comp.Name)
				}

				continue
			}
		}
	}

	return strings.Join(msgs, "\n"), mergeErrors(errs)
}

func mergeErrors(errs []error) error {
	var msgs []string

	for _, err := range errs {
		msgs = append(msgs, err.Error())
	}

	if len(msgs) == 0 {
		return nil
	}

	return fmt.Errorf(strings.Join(msgs, "\n"))
}
