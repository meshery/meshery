package kubernetes

import (
	"github.com/meshery/meshkit/errors"
)

var (
	ErrResyncK8SResourcesCode = "meshery-server-1372"
)

func ErrResyncK8SResources(err error) error {
	return errors.New(ErrResyncK8SResourcesCode, errors.Critical, []string{"Error resync resources"}, []string{err.Error()}, []string{"Fail to resync resources for the kubernetes machine"}, []string{"Check if machine context is assign to machine", "Check if machine context contains a reference to MesheryControllersHelper"})
}
