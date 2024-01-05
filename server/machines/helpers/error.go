package helpers

import "github.com/layer5io/meshkit/errors"
var (
	ErrAutoRegisterCode = "replace_me"
)

func ErrAutoRegister(err error, connType string) error {
	return errors.New(ErrAutoRegisterCode, errors.Alert, []string{}, []string{}, []string{}, []string{})
}