package patterns

import (
	"fmt"
	"strings"
	"sync"

	"github.com/layer5io/meshery/server/models"
		model "github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"

	"github.com/layer5io/meshery/server/models/pattern/patterns/k8s"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/v1beta1"

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

func Process(kconfigs []string, componets []v1beta1.ComponentDefinition, isDel bool, patternName string, ec *models.Broadcast, userID string, provider models.Provider, hostname model.IHost, skipCrdAndOperator, upgradeExistingRelease bool) ([]DeploymentMessagePerContext, error) {
	var comps []v1beta1.ComponentDefinition
	action := "deploy"
	if isDel {
		action = "undeploy"
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

				if !skipCrdAndOperator && hostname != nil && comp.Model.Name != (model.Kubernetes{}).String() {
					deploymentMsg := DeploymentMessagePerComp{
						Kind:       comp.Kind,
						Model:      comp.Model.Name,
						CompName:   comp.DisplayName,
						Success:    true,
						DesignName: patternName,
					}

					// Deploys resources that are required inside cluster for successful deployment of the design.
					result, err := hostname.HandleDependents(comp, kcli, !isDel, upgradeExistingRelease)
					// If dependencies were not resolved fail forward, there can be case that dependency already exist in the cluster.
					deploymentMsg.Message = result
					if err != nil {
						deploymentMsg.Success = false
						errs = append(errs, err)
						deploymentMsg.Error = err
					}
					msgsPerComp = append(msgsPerComp, deploymentMsg)
				}
				//All other components will be handled directly by Kubernetes
				//TODO: Add a Mapper utility function which carries the logic for X hosts can handle Y components under Z circumstances.

				_msg := DeploymentMessagePerComp{
					Kind:       comp.Kind,
					Model:      comp.Model.Name,
					CompName:   comp.DisplayName,
					Success:    true,
					DesignName: patternName,
					Message:    fmt.Sprintf("%sed %s/%s", action, patternName, comp.DisplayName),
				}

				if err := k8s.Deploy(kcli, comp, isDel); err != nil {
					_msg.Message = fmt.Sprintf("Error %sing %s/%s", action, patternName, comp.DisplayName)
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
