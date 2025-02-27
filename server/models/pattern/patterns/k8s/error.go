package k8s

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrDryRunCode = "meshery-server-1318"
)

func ErrDryRun(err error, obj string) error {
	return errors.New(ErrDryRunCode, errors.Alert, []string{"error performing a dry run on the design"}, []string{err.Error()}, []string{obj}, []string{})
}
