package main

import (
	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrorFailedRetreivingSheetCode      = "1233"
	ErrorFailedUpdatingSheetCode        = "1234"
	ErrorFailedRetreivingAHPackagesCode = "1235"
	ErrorIOCode                         = "1236"
	ErrorFailedGeneratingComponentsCode = "1237"
	ErrorFailedWritingComponentsCode    = "1238"
	ErrorFailedConnectingGCPCode        = "1239"
)

func ErrorFailedRetreivingSheet(err error) error {
	return errors.New(
		ErrorFailedRetreivingSheetCode,
		errors.Fatal,
		[]string{"Unable to retrieve data from sheet"},
		[]string{err.Error()}, nil, nil,
	)
}

func ErrorFailedConnectingGCP(err error) error {
	return errors.New(
		ErrorFailedConnectingGCPCode,
		errors.Fatal,
		[]string{"Unable to connect to google api"},
		[]string{err.Error()}, nil, nil,
	)
}

func ErrorFailedUpdatingSheet(err error) error {
	return errors.New(
		ErrorFailedUpdatingSheetCode,
		errors.Fatal,
		[]string{"Unable to update data in sheet"},
		[]string{err.Error()}, nil, nil,
	)
}

func ErrorFailedWritingComponents(err error) error {
	return errors.New(
		ErrorFailedWritingComponentsCode,
		errors.Fatal,
		[]string{"Failed to write components"},
		[]string{err.Error()}, nil, nil,
	)
}

func ErrorFailedGeneratingComponents(err error) error {
	return errors.New(
		ErrorFailedGeneratingComponentsCode,
		errors.Fatal,
		[]string{"Failed to generate components"},
		[]string{err.Error()}, nil, nil,
	)
}

func ErrorFailedRetreivingAHPackages(err error) error {
	return errors.New(
		ErrorFailedRetreivingAHPackagesCode,
		errors.Fatal,
		[]string{"Failed to retreive ArtifactHub packaged"},
		[]string{err.Error()}, nil, nil,
	)
}

func ErrorIOException(err error) error {
	return errors.New(
		ErrorIOCode,
		errors.Fatal,
		[]string{"I/O Error"},
		[]string{err.Error()}, nil, nil,
	)
}
