package environments

import "github.com/layer5io/meshkit/errors"

var (
	ErrEnvironmentCreateCode = "mesheryctl-1151"
)

func ErrEnvironmentCreate(err error) error {
	return errors.New(ErrEnvironmentCreateCode, errors.Alert, []string{"error creating environment"}, []string{err.Error()}, []string{"Invalid organization id provided"}, []string{"Ensure that the environment creation parameters are correct"})
}
