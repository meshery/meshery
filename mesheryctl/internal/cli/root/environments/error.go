package environments

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

var ErrCreateEnvironmentCode = "mesheryctl-1195"

func errCreateEnvironment(name, orgID string) error {
	return errors.New(ErrCreateEnvironmentCode, errors.Alert,
		[]string{fmt.Sprintf("Failed to create environment %s in organization %s", name, orgID)},
		[]string{"There was an error while creating the environment"},
		[]string{"Ensure Meshery server is running and reachable", "Check if the organization ID is correct"},
		[]string{"Start Meshery server and try again", "Provide a valid organization ID"},
	)
}
