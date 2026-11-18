package stages

const ResourceCapacityResponseKey = "resourceCapacity"

// ResourceCapacity returns the ChainStageFunction that computes the resource capacity
// and scheduling preview for the pattern processing chain.
// Like DryRun, this functions as an informational heads-up preview and does not mutate any cluster resources.
func ResourceCapacity(_ ServiceInfoProvider, act ServiceActionProvider) ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			act.Terminate(err)
			return
		}

		resp, capErr := act.ResourceCapacity(data.Pattern.Components)
		if capErr != nil {
			// Capacity preview errors are non-blocking warnings; log and proceed
			act.Log(capErr.Error())
		} else {
			data.Lock.Lock()
			if data.Other == nil {
				data.Other = make(map[string]interface{})
			}
			data.Other[ResourceCapacityResponseKey] = resp
			data.Lock.Unlock()
		}

		if next != nil {
			next(data, nil)
		}
	}
}
