package stages

import (
	stderrors "errors"
	"fmt"

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

// enrichComponent stamps a component with the design it came from. It is a
// variable so that a test can make it fail: the outcome of every component has
// to reach the summary whichever way it turns out, and that is not something
// the deployment paths can be asked to demonstrate on their own.
var enrichComponent = orchestration.EnrichComponentWithMesheryMetadata

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
			act.Terminate(planner.ErrCyclicDependency(data.Pattern.Name, prov.IsDelete()))
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
			err := enrichComponent(v1beta3Comp, data.Pattern.ID.String(), string(data.Pattern.Version))
			patternutils.ApplyV1beta3MetadataChanges(v1beta3Comp, &component)

			if err != nil {
				log.Error(err)
				data.Lock.Lock()
				data.Other[fmt.Sprintf("%s%s", name, ProvisionSuffixKey)] = dispatchFailureMessage(component, err, data.Pattern.Name, prov.IsDelete())
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

			// The outcome is recorded before anything is decided from it, so
			// that a component never disappears from the summary - not even
			// when it could not be dispatched at all, and not when a
			// fulfillment path comes back with nothing to say about it.
			//
			// The second case is what keeps the decision below sound: an empty
			// outcome is indistinguishable from a successful one, so it would
			// let the components that depend on this one deploy on a success
			// nothing ever reported.
			if len(msg) == 0 {
				if err != nil {
					msg = dispatchFailureMessage(component, err, data.Pattern.Name, prov.IsDelete())
				} else {
					msg = noOutcomeMessage(component, data.Pattern.Name, prov.IsDelete())
				}
			}

			data.Lock.Lock()
			// Store the outcome reported for this service
			data.Other[fmt.Sprintf("%s%s", name, ProvisionSuffixKey)] = msg
			data.Lock.Unlock()

			if err != nil {
				data.Lock.Lock()
				errs = append(errs, err)
				data.Lock.Unlock()
				return false
			}

			return provisionSucceeded(msg)
		}, func(name string, component component.ComponentDefinition, failedDependency string) {
			data.Lock.Lock()
			data.Other[fmt.Sprintf("%s%s", name, ProvisionSuffixKey)] = withheldMessage(component, displayNames[failedDependency], data.Pattern.Name, prov.IsDelete())
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

// provisionSucceeded reports whether the component itself was applied
// successfully everywhere it was acted on.
//
// Neither fulfillment path returns an error when the component itself fails to
// apply: the failure is carried on the per-component message. The messages are
// therefore what decides whether the components that depend on this one may
// proceed.
//
// Outcomes reported for a component's prerequisites are not part of that
// decision. Installing them is fail-forward and the component is applied
// regardless, so counting them here would withhold the dependents of a
// component that did deploy.
func provisionSucceeded(msgs []patterns.DeploymentMessagePerContext) bool {
	for _, msg := range msgs {
		for _, summary := range msg.Summary {
			if summary.IsPrerequisite {
				continue
			}

			if !summary.Success {
				return false
			}
		}
	}

	return true
}

// action names what the Provision stage is doing, so that a message reads
// correctly whichever direction the design is being applied in.
func action(isDelete bool) string {
	if isDelete {
		return "undeploy"
	}

	return "deploy"
}

// withheldMessage reports a component that was never dispatched because a
// component it declared a dependency on did not succeed.
func withheldMessage(comp component.ComponentDefinition, failedDependency, designName string, isDelete bool) []patterns.DeploymentMessagePerContext {
	dependency := failedDependency
	if dependency == "" {
		dependency = "a component it depends on"
	}

	return []patterns.DeploymentMessagePerContext{
		{
			Summary: []patterns.DeploymentMessagePerComp{
				{
					Kind:       comp.Component.Kind,
					Model:      modelName(comp),
					CompName:   comp.DisplayName,
					DesignName: designName,
					Success:    false,
					Message:    fmt.Sprintf("Withheld %s: %s did not %s successfully", comp.DisplayName, dependency, action(isDelete)),
					Error:      ErrDependencyNotSatisfied(comp.DisplayName, dependency, isDelete),
				},
			},
		},
	}
}

// dispatchFailureMessage reports a component that could not be handed to either
// fulfillment path at all, so that it is still named in the summary rather than
// silently missing from it.
func dispatchFailureMessage(comp component.ComponentDefinition, err error, designName string, isDelete bool) []patterns.DeploymentMessagePerContext {
	return []patterns.DeploymentMessagePerContext{
		{
			Summary: []patterns.DeploymentMessagePerComp{
				{
					Kind:       comp.Component.Kind,
					Model:      modelName(comp),
					CompName:   comp.DisplayName,
					DesignName: designName,
					Success:    false,
					Message:    fmt.Sprintf("Could not %s %s", action(isDelete), comp.DisplayName),
					Error:      err,
				},
			},
		},
	}
}

// noOutcomeMessage reports a component that a fulfillment path returned nothing
// at all about - neither an outcome nor a reason there is none.
//
// No fulfillment path does that today, and the guard is here because the
// alternative is silent rather than loud: whatever reads the summary cannot
// tell an empty outcome from a successful one, so a path that ever stopped
// reporting would deploy this component's dependents on a success no one
// claimed.
func noOutcomeMessage(comp component.ComponentDefinition, designName string, isDelete bool) []patterns.DeploymentMessagePerContext {
	return []patterns.DeploymentMessagePerContext{
		{
			Summary: []patterns.DeploymentMessagePerComp{
				{
					Kind:       comp.Component.Kind,
					Model:      modelName(comp),
					CompName:   comp.DisplayName,
					DesignName: designName,
					Success:    false,
					Message:    fmt.Sprintf("Could not confirm that %s was %sed: nothing was reported for it", comp.DisplayName, action(isDelete)),
				},
			},
		},
	}
}

func modelName(comp component.ComponentDefinition) string {
	if comp.Model == nil {
		return ""
	}

	return comp.Model.Name
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

// mergeErrors combines the errors collected over a stage into one.
//
// It joins rather than reformats, so that a structured MeshKit error keeps its
// code and stays recognizable to errors.GetCode as it leaves the chain, rather
// than being flattened into an opaque message. Whether that code survives any
// further is up to whoever handles the terminated chain: the deploy handler
// re-wraps what it is given, so the code it reports is its own.
func mergeErrors(errs []error) error {
	return stderrors.Join(errs...)
}
