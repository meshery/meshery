package registry

import "github.com/layer5io/meshkit/errors"

var (
	ErrGenerateModelCode     = "replace_me"
	ErrGenerateComponentCode = "replace_me"
)

func ErrGenerateModel(err error, modelName string) error {
	return errors.New(ErrGenerateModelCode, errors.Alert, []string{}, []string{}, []string{}, []string{}) 
}

func ErrGenerateComponent(err error, compName string) error {
	return errors.New(ErrGenerateComponentCode, errors.Alert, []string{}, []string{}, []string{}, []string{})
}
