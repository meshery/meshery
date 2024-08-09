package model

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrExportModelCode = "mesheryctl-1127"
)

func ErrExportModel(err error, name string) error {
	return errors.New(ErrExportModelCode, errors.Fatal, []string{"Error exporting model"}, []string{fmt.Sprintf("Given model with name: %s could not be exported: %s", name, err)}, []string{"Model may not be present in the registry"}, []string{"Ensure that there are no typos in the model name"})
}
