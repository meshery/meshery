package registry

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

var (
	ErrGenerateModelCode     = "1190"
	ErrGenerateComponentCode = "1191"
	ErrUpdateModelCode       = "1192"
	ErrUpdateComponentCode   = "1193"
	ErrUpdateRegistryCode    = "1194"
)

func ErrUpdateRegistry(err error, path string) error {
	return errors.New(ErrUpdateRegistryCode, errors.Alert, []string{"error updating registry at ", path}, []string{err.Error()}, []string{"Provided spreadsheet ID is incorrect", "Provided credentials are incorrect"}, []string{"Ensure correct spreadsheet ID is provided", "Ensure correct credentials are provided"})
}

func ErrGenerateModel(err error, modelName string) error {
	return errors.New(ErrGenerateModelCode, errors.Alert, []string{fmt.Sprintf("error generating model %s", modelName)}, []string{err.Error()}, []string{"Registrant used for the model is not available", "Failed to create model directory"}, []string{"Check network connevtivity and try again.", "Ensure sufficient permissions to allow creation of model directory"})
}

func ErrGenerateComponent(err error, modelName, compName string) error {
	return errors.New(ErrGenerateComponentCode, errors.Alert, []string{"error generating comp %s of model %s", compName, modelName}, []string{err.Error()}, []string{}, []string{})
}

func ErrUpdateModel(err error, modelName string) error {
	return errors.New(ErrUpdateModelCode, errors.Alert, []string{"error updating model ", modelName}, []string{err.Error()}, []string{"Model does not exist"}, []string{"Ensure existence of model, check for typo in model name"})
}

func ErrUpdateComponent(err error, modelName, compName string) error {
	return errors.New(ErrUpdateComponentCode, errors.Alert, []string{fmt.Sprintf("error updating component %s of model %s ", compName, modelName)}, []string{err.Error()}, []string{"Component does not exist", "Component definition is corrupted"}, []string{"Ensure existence of component, check for typo in component name", "Regenerate corrupted component"})
}
