package stages

import "github.com/layer5io/meshkit/errors"

const (
	ErrResolveReferenceCode = ""
)

func ErrResolveReference(err error) error {
	return errors.New(ErrResolveReferenceCode, errors.Alert, []string{}, []string{err.Error()}, []string{}, []string{})
}