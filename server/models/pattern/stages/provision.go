package stages

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/planner"
	meshmodelv1alpha1 "github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
)

type CompConfigPair struct {
	Component     v1alpha1.Component
	Configuration v1alpha1.Configuration
	Hosts         map[meshmodel.Host]bool
}

const ProvisionSuffixKey = ".isProvisioned"

func Provision(prov ServiceInfoProvider, act ServiceActionProvider) ChainStageFunction {
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

		config, err := data.Pattern.GenerateApplicationConfiguration()
		if err != nil {
			act.Terminate(fmt.Errorf("failed to generate application configuration: %s", err))
			return
		}

		errs := []error{}

		// Execute the plan
		_ = plan.Execute(func(name string, svc core.Service) bool {
			ccp := CompConfigPair{}

			// Create application component
			comp, err := data.Pattern.GetApplicationComponent(name)
			if err != nil {
				return false
			}

			// Generate hosts list
			ccp.Hosts = generateHosts(
				data.PatternSvcWorkloadCapabilities[name],
				data.PatternSvcTraitCapabilities[name],
				act.GetRegistry(),
			)
			// Get annotations for the component, if any
			comp.ObjectMeta.Annotations = helpers.MergeStringMaps(
				v1alpha1.GetAnnotationsForWorkload(data.PatternSvcWorkloadCapabilities[name]),
				comp.ObjectMeta.Annotations,
			)
			if core.Format { //deprettify the component before deploying
				comp.Spec.Settings = core.Format.DePrettify(comp.Spec.Settings, false)
			}
			ccp.Component = comp
			// Add configuration only if traits are applied to the component
			if len(svc.Traits) > 0 {
				ccp.Configuration = config
			}

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
		})

		if next != nil {
			next(data, mergeErrors(errs))
		}
	}
}

func processAnnotations(pattern *core.Pattern) {
	for name, svc := range pattern.Services {
		if svc.IsAnnotation {
			// this particular block is present so that designs with previous filters don't break 
			// also UI is dependent but not exactly sure how?
			delete(pattern.Services, name)
		}
		data, ok := svc.Traits["meshmap"]
		if ok {
			metadata, ok2 := data.(map[string]interface{})
			if ok2 {
				compMetadata, ok3 := metadata["meshmodel-metadata"].(map[string]interface{})
				if ok3 {
					isAnnotation, ok4 := compMetadata["isAnnotation"].(bool)
					if ok4 && isAnnotation {
						delete(pattern.Services, name)
					}
				}
			}
		}
	}
}

func generateHosts(wc meshmodelv1alpha1.ComponentDefinition, _ []core.TraitCapability, reg *meshmodel.RegistryManager) map[meshmodel.Host]bool {
	res := map[meshmodel.Host]bool{}
	host := reg.GetRegistrant(wc)
	res[host] = true
	// for _, tc := range tcs {
	// 	res[tc.Host] = true
	// }

	return res
}

func mergeErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}

	var errMsg []string
	for _, err := range errs {
		errMsg = append(errMsg, err.Error())
	}

	return fmt.Errorf(strings.Join(errMsg, "\n"))
}
