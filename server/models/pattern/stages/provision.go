package stages

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"github.com/layer5io/meshery/server/models/pattern/planner"
	"github.com/layer5io/meshery/server/models/pattern/resource/selector"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
)

type CompConfigPair struct {
	Component     v1alpha1.Component
	Configuration v1alpha1.Configuration
	Hosts         map[string]bool
}

const ProvisionSuffixKey = ".isProvisioned"

func Provision(prov ServiceInfoProvider, act ServiceActionProvider) ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			act.Terminate(err)
			return
		}

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
			)

			comp.SetLabels(map[string]string{
				"resource.pattern.meshery.io/id": svc.ID.String(),
			})

			// Get annotations for the component, if any
			comp.Annotations = helpers.MergeStringMaps(
				selector.GetAnnotationsForWorkload(data.PatternSvcWorkloadCapabilities[name]),
				comp.Annotations,
			)

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

func generateHosts(wc core.WorkloadCapability, tcs []core.TraitCapability) map[string]bool {
	res := map[string]bool{}

	res[wc.Host] = true
	for _, tc := range tcs {
		res[tc.Host] = true
	}

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
