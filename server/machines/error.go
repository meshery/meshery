package machines

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrAutoRegisterCode           = "replace_me"
)
func ErrAutoRegister(err error, connType string) error {
	return errors.New(ErrAutoRegisterCode, errors.Alert, []string{}, []string{err.Error()}, []string{""}, []string{""})
}
