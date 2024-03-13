package k8s

import (
	"encoding/json"
	"fmt"
	meshkitutils "github.com/layer5io/meshkit/utils/kubernetes"

	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshkit/errors"
	kubeerror "k8s.io/apimachinery/pkg/api/errors"
)

const (
	ErrDryRunCode = "meshery-server-1318"
)

func isErrKubeStatusErr(err error) bool {
	switch err.(type) {
	case *kubeerror.StatusError:
		return true
	default:
		return false
	}
}

func formatKubeStatusErrToMeshkitErr(status *[]byte, componentName string) error {
	var shortDescription, longDescription []string

	var kubeErr kubeerror.StatusError

	err := json.Unmarshal(*status, &kubeErr)
	kubeStatus := kubeErr.ErrStatus
	if err != nil {
		return err
	}

	st := string(kubeStatus.Status)

	if kubeStatus.Details != nil {

		sd := kubeStatus.Details.Kind

		sd = fmt.Sprintf("%s %s", sd, st)

		sd = fmt.Sprintf("%s \"%s\"", sd, kubeStatus.Details.Name)

		for _, cause := range kubeStatus.Details.Causes {
			var longDes, field string
			longDes = utils.FormatK8sMessage(cause.Message)
			longDescription = append(longDescription, longDes)
			sd = fmt.Sprintf("%s %s", sd, cause.Type)
			field = componentName + "." + utils.GetComponentFieldPathFromK8sFieldPath(cause.Field)
			shortDescription = append(shortDescription, sd, field)
		}
		if len(kubeStatus.Details.Causes) == 0 {
			longDescription = append(longDescription, kubeStatus.Message)
		}
	}
	return errors.New(meshkitutils.ErrApplyManifestCode, errors.Alert, shortDescription, longDescription, []string{}, []string{})
}

func ErrDryRun(err error, obj string) error {
	return errors.New(ErrDryRunCode, errors.Alert, []string{"error performing a dry run on the design"}, []string{err.Error()}, []string{obj}, []string{})
}
