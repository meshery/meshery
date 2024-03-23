package connections

import "github.com/layer5io/meshkit/errors"

const (
	// Error name for the error when the argument are not passed correctly
	ErrInvalidArgumentCode = "mesheryctl-1118"
)

func ErrInvalidArgument(err error) error {
	return errors.New(ErrInvalidArgumentCode, errors.Alert, []string{"Invalid Argument"}, []string{err.Error()}, []string{"Invalid Argument"}, []string{"Please check the arguments passed"})
}
