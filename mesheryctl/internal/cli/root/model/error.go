package model

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrExportModelCode                      = "mesheryctl-1127"
	ErrTemplateFileEmptyImportModelCode     = "mesheryctl-1128"
	ErrReadTemplateFileEmptyImportModelCode = "mesheryctl-1129"
	ErrRegisterModelCode                    = "mesheryctl-1130"
	ErrReadFilePathModelCode                = "mesheryctl-1131"
	ErrPathErrorModelCode                   = "mesheryctl-1133"
)

func formatErrorWithReference(cmd string) string {
	baseURL := "https://docs.meshery.io/reference/mesheryctl/model"
	switch cmd {
	case "import":
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL+"/import")
	default:
		return fmt.Sprintf("\nSee %s for usage details\n", baseURL)
	}
}

func ErrExportModel(err error, name string) error {
	return errors.New(ErrExportModelCode, errors.Fatal,
		[]string{"Error exporting model"},
		[]string{fmt.Sprintf("Given model with name: %s could not be exported: %s", name, err)},
		[]string{"Model may not be present in the registry"},
		[]string{"Ensure that there are no typos in the model name"})
}

func ErrTemplateFileEmptyImportModel() error {
	return errors.New(ErrTemplateFileEmptyImportModelCode, errors.Fatal,
		[]string{"Error template file is empty on model import"},
		[]string{fmt.Sprintf("Template file is required to import model from an URL")},
		[]string{"Import a model from a URL required a template file"},
		[]string{"Ensure to provide a template file to import a model from an URL"})
}

func ErrReadTemplateFileEmptyImportModel(err error) error {
	return errors.New(ErrReadTemplateFileEmptyImportModelCode, errors.Fatal,
		[]string{"Error reading template file provided on model import"},
		[]string{fmt.Sprintf("Could not read the provided template file: %v", err)},
		[]string{"An issue happened while reading the template file provided"}, []string{"Not able to read the provided template file"})
}

func ErrRegisterModel(err error) error {
	return errors.New(ErrRegisterModelCode, errors.Fatal,
		[]string{"Error to register model"},
		[]string{fmt.Sprintf("Error while trying to register model: %v", err)},
		[]string{"An error happens while trying to register a model"},
		[]string{"An error happens while trying to register a model"})
}

func ErrPathErrorModel(err error, cmd string) error {
	return errors.New(ErrPathErrorModelCode, errors.Alert,
		[]string{"Error check if provided path is a directory"},
		[]string{fmt.Sprintf("Error while validating provided path: %v", err)},
		[]string{"An error happens while heck if provided path is a directory"},
		[]string{"Provide a valid folder path", formatErrorWithReference(cmd)})
}

func ErrReadFilePathModel(err error, cmd string) error {
	return errors.New(ErrReadFilePathModelCode, errors.Alert,
		[]string{"Unable to read file"},
		[]string{err.Error()},
		[]string{"The provided file is not present or has an invalid path"},
		[]string{"Provide a valid file path", formatErrorWithReference(cmd)})
}
