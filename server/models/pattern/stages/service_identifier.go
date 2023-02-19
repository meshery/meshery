package stages

import (
	"fmt"

	"github.com/gofrs/uuid"
)

const UpdateSuffixKey = ".isUpdate"

// ServiceIdentifier takes in a service identity provider and returns a ChainStageFunction
func ServiceIdentifierAndMutator(prov ServiceInfoProvider, act ServiceActionProvider) ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		// Find the ID of the resources
		for svcName, svc := range data.Pattern.Services {
			id, err := prov.GetMesheryPatternResource(
				svcName,
				svc.Namespace,
				svc.Type,
				"workload",
			)

			if err != nil || id == nil {
				// Don't terminate - assign new ID instead
				uid, err := uuid.NewV4()
				if err != nil {
					act.Terminate(err)
					return // Terminate if ID generation fails
				}

				svc.ID = &uid
				continue
			}

			// Assign the ID to the service
			svc.ID = id
			data.Lock.Lock()
			data.Other[fmt.Sprintf("%s%s", svcName, UpdateSuffixKey)] = true
			data.Lock.Unlock()
		}
		act.Mutate(data.Pattern)
		if next != nil {
			next(data, nil)
		}
	}
}
