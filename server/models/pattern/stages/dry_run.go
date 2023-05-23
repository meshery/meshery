package stages

import (
	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
)

// import "github.com/layer5io/meshery/server/models/pattern/patterns"

const DryRunResponseKey = "dryRunResponse"

var mesheryDefinedAPIVersions = map[string]bool{
	"core.oam.dev/v1alpha1": true,
}

// There are two types of errors here:
// 1. Error while performing the Dry Run (when the DryRun request could not be sent)
// 2. Errors in Dry Run (when the Dry Run request was performed successfully but there are errors in the Object sent for DryRun)
// We are not considering #2 as an error. #2 is treated as the Response of the Dry Run.
// DryRun stage does not terminate the whole chain when an error of type #2 is encountered, it rather
// stores that information in the `Other` placeholder for future use.
// This is in contrast with the Validation stage where the Validation errors terminate the chain.
func DryRun(_ ServiceInfoProvider, act ServiceActionProvider) ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			act.Terminate(err)
			return
		}
		var comps []v1alpha1.Component
		for name, svc := range data.Pattern.Services {
			if mesheryDefinedAPIVersions[svc.APIVersion] {
				continue
			}
			comp, err := data.Pattern.GetApplicationComponent(name)
			if err != nil {
				continue
			}
			comp.ObjectMeta.Annotations = helpers.MergeStringMaps(
				v1alpha1.GetAnnotationsForWorkload(data.PatternSvcWorkloadCapabilities[name]),
				comp.ObjectMeta.Annotations,
			)
			comps = append(comps, comp)
		}
		resp, err := act.DryRun(comps)
		if err != nil {
			act.Terminate(err)
			return
		}
		data.Lock.Lock()
		if data.Other == nil {
			data.Other = make(map[string]interface{})
		}
		data.Other[DryRunResponseKey] = resp
		data.Lock.Unlock()
		if next != nil {
			next(data, err)
		}
	}
}
