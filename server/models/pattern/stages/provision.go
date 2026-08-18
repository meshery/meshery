package stages

import (
	"fmt"
	"strings"

	"github.com/meshery/meshery/server/models/pattern/patterns"
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

		annotations := processAnnotations(data.Pattern)

		// Create provision plan
		plan, err := planner.CreatePlan(*data.Pattern, annotations, prov.IsDelete())
		if err != nil {
			act.Terminate(err)
			return
		}

		// Check feasibility of the generated plan
		feasible, err := plan.IsFeasible()
		if err != nil {
			act.Terminate(err)
			return
		}

		if !feasible {
			act.Terminate(planner.ErrCyclicDependency(data.Pattern.Name))
			return
		}

		errs := []error{}

		// The plan is keyed by component id; carry the display names so that a
		// component withheld by a failed dependency can name that dependency.
		displayNames := map[string]string{}
		for _, component := range data.Pattern.Components {
			displayNames[component.ID.String()] = component.DisplayName
		}

		// Execute the plan
		err = plan.Execute(func(name string, component component.ComponentDefinition) bool {
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
				data.Lock.Lock()
				errs = append(errs, err)
				data.Lock.Unlock()
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
				data.Lock.Lock()
				errs = append(errs, err)
				data.Lock.Unlock()
				return false
			}
			data.Lock.Lock()
			// Store the outcome reported for this service
			data.Other[fmt.Sprintf("%s%s", name, ProvisionSuffixKey)] = msg
			data.Lock.Unlock()

			return provisionSucceeded(msg)
		}, func(name string, component component.ComponentDefinition, failedDependency string) {
			data.Lock.Lock()
			data.Other[fmt.Sprintf("%s%s", name, ProvisionSuffixKey)] = withheldMessage(component, displayNames[failedDependency], data.Pattern.Name)
			data.Lock.Unlock()
		}, log)
		if err != nil {
			act.Terminate(err)
			return
		}

		if next != nil {
			next(data, mergeErrors(errs))
		}
	}
}

// provisionSucceeded reports whether every context that acted on a component
// reported success.
//
// Neither fulfillment path returns an error when the component itself fails to
// apply: the failure is carried on the per-component message. The messages are
// therefore what decides whether the components that depend on this one may
// proceed.
func provisionSucceeded(msgs []patterns.DeploymentMessagePerContext) bool {
	for _, msg := range msgs {
		for _, summary := range msg.Summary {
			if !summary.Success {
				return false
			}
		}
	}

	return true
}

// withheldMessage reports a component that was never dispatched because a
// component it declared a dependency on did not succeed.
func withheldMessage(comp component.ComponentDefinition, failedDependency, designName string) []patterns.DeploymentMessagePerContext {
	dependency := failedDependency
	if dependency == "" {
		dependency = "a component it depends on"
	}

	model := ""
	if comp.Model != nil {
		model = comp.Model.Name
	}

	return []patterns.DeploymentMessagePerContext{
		{
			Summary: []patterns.DeploymentMessagePerComp{
				{
					Kind:       comp.Component.Kind,
					Model:      model,
					CompName:   comp.DisplayName,
					DesignName: designName,
					Success:    false,
					Message:    fmt.Sprintf("Withheld %s: %s did not deploy successfully", comp.DisplayName, dependency),
					Error:      ErrDependencyNotSatisfied(comp.DisplayName, dependency),
				},
			},
		},
	}
}

// processAnnotations sets the non-semantic components of a design aside so that
// they are never deployed, and returns them.
func processAnnotations(pattern *pattern.PatternFile) []*component.ComponentDefinition {
	components := []*component.ComponentDefinition{}
	annotations := []*component.ComponentDefinition{}

	for _, component := range pattern.Components {
		if component.Metadata.IsAnnotation {
			annotations = append(annotations, component)
			continue
		}

		components = append(components, component)
	}

	pattern.Components = components

	return annotations
}

func generateHosts(cd component.ComponentDefinition, reg *meshmodel.RegistryManager) []connection.Connection {
	// registry.GetRegistrant expects an entity.Entity; v1beta3/component is
	// the canonical-casing version that implements that interface. Bridge
	// just for the registry lookup — the enclosing pattern continues to
	// hold the component in its native v1beta2 representation.
	v1beta3Cd := patternutils.ComponentV1beta2ToV1beta3(&cd)
	_connection := reg.GetRegistrant(v1beta3Cd)
	return []connection.Connection{meshmodel.RegistrantHostToV1beta1(_connection)}
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
