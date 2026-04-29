package kubernetes

import (
	"github.com/meshery/meshkit/errors"
)

var (
	ErrResyncK8SResourcesCode = "meshery-server-1372"
	ErrConnectActionCode      = "meshery-server-1373"
)

func ErrResyncK8SResources(err error) error {
	return errors.New(ErrResyncK8SResourcesCode, errors.Critical, []string{"Error resync resources"}, []string{err.Error()}, []string{"Fail to resync resources for the kubernetes machine"}, []string{"Check if machine context is assign to machine", "Check if machine context contains a reference to MesheryControllersHelper"})
}

func ErrConnectAction(err error) error {
	return errors.New(ErrConnectActionCode, errors.Critical, []string{"Error connect action"}, []string{err.Error()}, []string{"Fail to perform connect action for the kubernetes machine"}, []string{"Check if token is passed machine from golang context correctly", "Check if there is connection data stored in database for this kubernetes machine"})
}
