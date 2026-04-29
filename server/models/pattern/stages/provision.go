package stages

import (
	"fmt"
	"strings"

	"github.com/meshery/meshery/server/models/pattern/planner"
	patternutils "github.com/meshery/meshery/server/models/pattern/utils"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/orchestration"

	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
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

			// meshkit's orchestration.EnrichComponentWithMesheryMetadata operates
			// on *v1beta3/component.ComponentDefinition; bridge to that type,
			// then copy mutations back onto the v1beta2 value held by the
			// pattern so the plan sees the enriched Configuration/Metadata.
			v1beta3Comp := patternutils.ComponentV1beta2ToV1beta3(&component)
			err := orchestration.EnrichComponentWithMesheryMetadata(v1beta3Comp, data.Pattern.ID.String(), string(data.Pattern.Version))
			patternutils.ApplyV1beta3MetadataChanges(v1beta3Comp, &component)

			if err != nil {
				fmt.Println("Err while assigning labels", err)
				errs = append(errs, err)
				return false
			}

			// Generate hosts list
			ccp.Hosts = generateHosts(
				data.DeclartionToDefinitionMapping[component.ID],
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
	// registry.GetRegistrant expects an entity.Entity; v1beta3/component is
	// the canonical-casing version that implements that interface. Bridge
	// just for the registry lookup — the enclosing pattern continues to
	// hold the component in its native v1beta2 representation.
	v1beta3Cd := patternutils.ComponentV1beta2ToV1beta3(&cd)
	_connection := reg.GetRegistrant(v1beta3Cd)
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
