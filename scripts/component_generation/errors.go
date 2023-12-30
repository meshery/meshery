package main

import (
	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrorFailedRetreivingCode = "1233"
)

func ErrorFailedRetreiving(err error) error {
	return errors.New(
		ErrorFailedRetreivingCode,
		errors.Fatal,
		[]string{"Unable to retrieve data from sheet"},
		[]string{err.Error()}, nil, nil,
	)
}
