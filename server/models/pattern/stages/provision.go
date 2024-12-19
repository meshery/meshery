package stages

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery/server/models/pattern/planner"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/orchestration"

	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

type CompConfigPair struct {
	Component component.ComponentDefinition
	Hosts     []connection.Connection
}

const ProvisionSuffixKey = ".isProvisioned"

func Provision(prov ServiceInfoProvider, act ServiceActionProvider, log logger.Handler) ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			act.Terminate(err)
			return
		}

		processAnnotations(data.Pattern)

		// Create provision plan
		plan, err := planner.CreatePlan(*data.Pattern, prov.IsDelete())
		if err != nil {
			act.Terminate(err)
			return
		}

		// Check feasibility of the generated plan
		if !plan.IsFeasible() {
			act.Terminate(fmt.Errorf("infeasible execution: detected cycle in the plan"))
			return
		}

		errs := []error{}

		// Execute the plan
		_ = plan.Execute(func(name string, component component.ComponentDefinition) bool {
			ccp := CompConfigPair{}

			err := orchestration.EnrichComponentWithMesheryMetadata(&component, data.Pattern.Id.String(), data.Pattern.Version)

			if err != nil {
				fmt.Println("Err while assigning labels", err)
				errs = append(errs, err)
				return false
			}

			// Generate hosts list
			ccp.Hosts = generateHosts(
				data.DeclartionToDefinitionMapping[component.Id],
				act.GetRegistry(),
			)

			ccp.Component = component

			msg, err := act.Provision(ccp)
			if err != nil {
				errs = append(errs, err)
				return false
			}
			data.Lock.Lock()
			// Store that this service was provisioned successfully
			data.Other[fmt.Sprintf("%s%s", name, ProvisionSuffixKey)] = msg
			data.Lock.Unlock()

			return true
		}, log)

		if next != nil {
			next(data, mergeErrors(errs))
		}
	}
}

func processAnnotations(pattern *pattern.PatternFile) {
	components := []*component.ComponentDefinition{}
	for _, component := range pattern.Components {
		if !component.Metadata.IsAnnotation {
			components = append(components, component)
		}
	}
	pattern.Components = components
}

func generateHosts(cd component.ComponentDefinition, reg *meshmodel.RegistryManager) []connection.Connection {
	_connection := reg.GetRegistrant(&cd)
	return []connection.Connection{_connection}
}

func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var errMsg []string
	for _, err := range errs {
		errMsg = append(errMsg, err.Error())
	}

	return fmt.Errorf("%s", strings.Join(errMsg, "\n"))
}
