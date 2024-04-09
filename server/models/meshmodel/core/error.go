package core

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrCreatingKubernetesComponentsCode = "meshery-server-1313"
	ErrRegisterEntityCode               = "meshery-server-1319"
)

func ErrCreatingKubernetesComponents(err error, ctxID string) error {
	return errors.New(ErrCreatingKubernetesComponentsCode, errors.Alert, []string{"failed to register/create kubernetes components for contextID " + ctxID}, []string{err.Error()}, []string{"component generation was canceled due to deletion or reload of K8s context", "Invalid kubeconfig", "Filters passed incorrectly in config", "Could not fetch API resources from Kubernetes server"}, []string{"If there is the log \"Starting to register ...\" for the same contextID after this error means that for some reason the context was reloaded which caused this run to abort. In that case, this error can be ignored.", "Make sure that the configuration filters passed are in accordance with output from /openapi/v2"})
}

func ErrRegisterEntity(err error, name, entity string) error {
	return errors.New(ErrRegisterEntityCode, errors.Alert, []string{fmt.Sprintf("failed to register %s %s", name, entity)}, []string{err.Error()}, []string{fmt.Sprintf("%s definition violates the definition schema", entity), fmt.Sprintf("%s might be missing model details", entity)}, []string{fmt.Sprintf("ensure the %s definition follows the correct schema", entity), fmt.Sprintf("ensure %s definition belongs to correct model", entity)})
}
