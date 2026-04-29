package helpers

import (
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
)

var (
	ErrAutoRegisterCode    = "meshery-server-1219"
	ErrResyncResourcesCode = "meshery-server-1371"
)

func ErrAutoRegister(err error, connType string) error {
	return errors.New(ErrAutoRegisterCode, errors.Alert, []string{}, []string{}, []string{}, []string{})
}

func IsConnectionUpdateErr(err error) bool {
	return errors.GetCode(err) == models.ErrUpdateConnectionStatusCode
}

func ErrResyncResources(err error) error {
	return errors.New(ErrResyncResourcesCode, errors.Critical, []string{"Error resync resources"}, []string{err.Error()}, []string{"Fail to resync resources for the machine"}, []string{"Check if machine supports resource resync"})
}
