package patterns

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/patterns/k8s"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/kubernetes"
)

type DeploymentMessagePerComp struct {
	Kind       string
	Model      string
	CompName   string
	Success    bool
	DesignName string
	Message    string
	Metadata   map[string]interface{}
	Error      error
}

type DeploymentMessagePerContext struct {
	Summary    []DeploymentMessagePerComp
	SystemName string
	Location   string
}

func Process(kconfigs []string, oamComps []string, oamConfig string, isDel bool, patternName string, ec *models.Broadcast, userID string, provider models.Provider, hostname v1beta1.IHost, skipCrdAndOperator, upgradeExistingRelease bool) ([]DeploymentMessagePerContext, error) {
	var comps []v1beta1.Component
	var config v1alpha1.Configuration
	action := "deploy"
	if isDel {
		action = "undeploy"
	}
	for _, oamComp := range oamComps {
		var comp v1beta1.Component
		if err := json.Unmarshal([]byte(oamComp), &comp); err != nil {
			return nil, err
		}

		comps = append(comps, comp)
	}

	if err := json.Unmarshal([]byte(oamConfig), &config); err != nil {
		return nil, err
	}

	msgs := make([]DeploymentMessagePerContext, 0)
	var msgsMx sync.Mutex

	var errs []error
	var kclis []*kubernetes.Client
	for _, config := range kconfigs {
		cli, err := kubernetes.New([]byte(config))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		kclis = append(kclis, cli)
	}

	var wg sync.WaitGroup
	for _, kcli := range kclis {
		wg.Add(1)
		go func(kcli *kubernetes.Client) {
			defer wg.Done()

			msgsPerComp := make([]DeploymentMessagePerComp, 0)
			for _, comp := range comps {

				if !skipCrdAndOperator && hostname != nil && comp.Spec.Model != (v1beta1.Kubernetes{}).String() {
					deploymentMsg := DeploymentMessagePerComp{
						Kind: comp.Kind,
						Model: comp.Spec.Model,
						CompName: comp.Name,
						Success: true,
						DesignName: patternName,
					}

					// if !isDel {
					// 	description = fmt.Sprintf("Detected dependency for %s/%s, deploying dependent model %s.", patternName, comp.Name, comp.Spec.Model)
					// } else {
					// 	description = fmt.Sprintf("Detected dependency for %s/%s, undeploying dependent model %s.", patternName, comp.Name, comp.Spec.Model)
					// }

					fmt.Println("host: ", hostname, " comp: ", comp.Name, "type: ", comp.Spec.Type, comp)
					// Deploys resources that are required inside cluster for successful deployment of the design.
					result, err := hostname.HandleDependents(comp, kcli, !isDel, upgradeExistingRelease)
					// If dependencies were not resolved fail forward, there can be case that dependency already exist in the cluster.
					deploymentMsg.Message = result
					if err != nil {
						deploymentMsg.Success = false
						errs = append(errs, err)
						deploymentMsg.Error = err
					}
				}
				//All other components will be handled directly by Kubernetes
				//TODO: Add a Mapper utility function which carries the logic for X hosts can handle Y components under Z circumstances.

				_msg := DeploymentMessagePerComp{
					Kind:       comp.Spec.Type,
					Model:      comp.Spec.Model,
					CompName:   comp.Name,
					Success:    true,
					DesignName: patternName,
					Message:    fmt.Sprintf("%sed %s/%s", action, patternName, comp.Name),
				}

				if err := k8s.Deploy(kcli, comp, config, isDel); err != nil {
					_msg.Message = fmt.Sprintf("Error %sing %s/%s", action, patternName, comp.Name)
					_msg.Error = err
					_msg.Success = false

				}
				msgsPerComp = append(msgsPerComp, _msg)
			}

			msgsMx.Lock()
			msgs = append(msgs, DeploymentMessagePerContext{
				Summary:    msgsPerComp,
				SystemName: kcli.RestConfig.ServerName,
				Location:   kcli.RestConfig.Host,
			})
			msgsMx.Unlock()

		}(kcli)
	}
	wg.Wait()
	return msgs, mergeErrors(errs)
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
