package patterns

import (
	"fmt"
	"strings"
	"sync"

	"github.com/layer5io/meshery/server/models"
	_models "github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"

	"github.com/layer5io/meshery/server/models/pattern/patterns/k8s"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/connection"
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

func Process(kconfigs []string, componets []component.ComponentDefinition, isDel bool, patternName string, ec *models.Broadcast, userID string, provider models.Provider, connection connection.Connection, skipCrdAndOperator, upgradeExistingRelease bool) ([]DeploymentMessagePerContext, error) {
	action := "deploy"
	if isDel {
		action = "undeploy"
	}

	depHandler, err := _models.NewDependencyHandler(connection.Kind)
	if err != nil {
		return nil, err
	}

	fmt.Println("\n\n\nTEST DEPHANDLER : ", depHandler)

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
			for _, comp := range componets {
				fmt.Println("TEST INSIDE line 70 : ", comp.Component.Kind)
				if !skipCrdAndOperator && depHandler != nil && comp.Model.Name != (_models.Kubernetes{}).String() {
					fmt.Println("TEST INSIDE line 72 : ")
					deploymentMsg := DeploymentMessagePerComp{
						Kind:       comp.Component.Kind,
						Model:      comp.Model.Name,
						CompName:   comp.DisplayName,
						Success:    true,
						DesignName: patternName,
					}

					// Deploys resources that are required inside cluster for successful deployment of the design.
					result, err := depHandler.HandleDependents(comp, kcli, !isDel, upgradeExistingRelease)
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
					Kind:       comp.Component.Kind,
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
				fmt.Println("TEST INSIDE line 108 after deploying : ", err)
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

	return fmt.Errorf("%s", strings.Join(msgs, "\n"))
}
