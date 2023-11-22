package core

import "github.com/layer5io/meshkit/errors"

const (
	ErrCreatingKubernetesComponentsCode = "1545"
)

func ErrCreatingKubernetesComponents(err error, ctxID string) error {
	return errors.New(ErrCreatingKubernetesComponentsCode, errors.Alert, []string{"failed to register/create kubernetes components for contextID " + ctxID}, []string{err.Error()}, []string{"component generation was canceled due to deletion or reload of K8s context", "Invalid kubeconfig", "Filters passed incorrectly in config", "Could not fetch API resources from Kubernetes server"}, []string{"If there is the log \"Starting to register ...\" for the same contextID after this error means that for some reason the context was reloaded which caused this run to abort. In that case, this error can be ignored.", "Make sure that the configuration filters passed are in accordance with output from /openapi/v2"})
}
