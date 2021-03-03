package resolver

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	errCode               = "test_code"
	ErrInvalidRequestCode = "test_code"
	ErrNilClientCode      = "test_code"
	ErrCreateDataCode     = "test_code"
	ErrQueryCode          = "test_code"
)

var (
	ErrNilClient      = errors.NewDefault(ErrNilClientCode, "Kubernetes client not initialized")
	ErrInvalidRequest = errors.NewDefault(ErrInvalidRequestCode, "Invalid query, please check syntax")
)

func ErrCreateData(err error) error {
	return errors.NewDefault(ErrCreateDataCode, "Error while writing meshsync data", err.Error())
}

func ErrUpdateData(err error) error {
	return errors.NewDefault(ErrCreateDataCode, "Error while updating meshsync data", err.Error())
}

func ErrDeleteData(err error) error {
	return errors.NewDefault(ErrCreateDataCode, "Error while deleting meshsync data", err.Error())
}

func ErrQuery(err error) error {
	return errors.NewDefault(ErrQueryCode, "Error while querying data", err.Error())
}
