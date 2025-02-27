package model

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrExportModelCode            = "mesheryctl-1127"
	ErrTemplateFileNotPresentCode = "mesheryctl-1131"
)

func ErrExportModel(err error, name string) error {
	return errors.New(ErrExportModelCode, errors.Fatal, []string{"Error exporting model"}, []string{fmt.Sprintf("Given model with name: %s could not be exported: %s", name, err)}, []string{"Model may not be present in the registry"}, []string{"Ensure that there are no typos in the model name"})
}
func ErrTemplateFileNotPresent() error {
	return errors.New(ErrTemplateFileNotPresentCode, errors.Fatal, []string{"error no template file provided"}, []string{"no template file is provided while using url for importing a model "}, []string{"template file not present"}, []string{"ensure that the template file is present in the given path"})
}
