package stages

import "github.com/gofrs/uuid"

// ServiceIdentifier takes in a service identity provider and returns a ChainStageFunction
func ServiceIdentifier(prov ServiceInfoProvider) ChainStageFunction {
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
					return // Terminate if ID generation fails
				}

				id = &uid
			}

			// Assign the ID to the service
			svc.ID = id
		}

		if next != nil {
			next(data, nil)
		}
	}
}
