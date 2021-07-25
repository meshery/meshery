package system

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrHealthCheckFailedCode = "1000"
	ErrInvalidAdapterCode    = "1001"
	ErrDownloadFileCode      = "1002"
)

func ErrHealthCheckFailed(err error) error {
	return errors.New(ErrHealthCheckFailedCode, errors.Alert, []string{"Health checks failed"}, []string{err.Error()}, []string{}, []string{})
}

func ErrInvalidAdapter(err error, obj string) error {
	return errors.New(ErrInvalidAdapterCode, errors.Alert, []string{"Invalid adapter ", obj, " specified"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDownloadFile(err error, obj string) error {
	return errors.New(ErrDownloadFileCode, errors.Alert, []string{"Error downloading file ", obj}, []string{err.Error()}, []string{}, []string{})
}
