package patterns

import (
	"fmt"
	"strings"
	"sync"

	"github.com/meshery/meshery/server/models"
	_models "github.com/meshery/meshkit/models/meshmodel/core/v1beta1"

	"github.com/meshery/meshery/server/models/pattern/patterns/k8s"
	patternutils "github.com/meshery/meshery/server/models/pattern/utils"
	"github.com/meshery/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta2/component"
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

	// IsPrerequisite marks an outcome reported for the prerequisites a
	// component needs rather than for the component itself.
	//
	// Installing prerequisites is deliberately fail-forward - what a component
	// needs may already be present in the cluster - and the component is
	// applied either way. A failed install therefore carries Success: false
	// like any other entry, and this flag is the only thing telling the two
	// apart: whatever decides whether the components that depend on this one
	// may proceed has to skip the entries carrying it, or a component that did
	// deploy ends up withholding its dependents.
	IsPrerequisite bool
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
				if !skipCrdAndOperator && depHandler != nil && comp.Model.Name != (_models.Kubernetes{}).String() {
					deploymentMsg := DeploymentMessagePerComp{
						Kind:           comp.Component.Kind,
						Model:          comp.Model.Name,
						CompName:       comp.DisplayName,
						Success:        true,
						DesignName:     patternName,
						IsPrerequisite: true,
					}

					// Deploys resources that are required inside cluster for successful deployment of the design.
					// meshkit's DependencyHandler.HandleDependents takes
					// *v1beta3/component.ComponentDefinition (the registry-
					// canonical casing); PatternFile.Components pins
					// v1beta2/component, so we bridge just at the meshkit
					// boundary rather than holding v1beta3 everywhere.
					v1beta3Comp := patternutils.ComponentV1beta2ToV1beta3(&comp)
					result, err := depHandler.HandleDependents(*v1beta3Comp, kcli, !isDel, upgradeExistingRelease)
					// If dependencies were not resolved fail forward, there can be case that dependency already exist in the cluster.
					//
					// The failure stays on this component's own summary entry and is
					// deliberately not returned from Process: the error Process returns
					// means the design could not be dispatched at all, and a component
					// that goes on to apply successfully must not be treated as failed.
					deploymentMsg.Message = result
					if err != nil {
						deploymentMsg.Success = false
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
