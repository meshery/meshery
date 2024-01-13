package helpers

import (
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/errors"
)

var (
	ErrAutoRegisterCode = "1556"
)

func ErrAutoRegister(err error, connType string) error {
	return errors.New(ErrAutoRegisterCode, errors.Alert, []string{}, []string{}, []string{}, []string{})
}

func IsConnectionUpdateErr(err error) bool {
	return errors.GetCode(err) == models.ErrUpdateConnectionStatusCode
}
