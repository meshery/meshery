// Package models - Error codes for Meshery models
package models

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrGrafanaClientCode        = "3000"
	ErrPageSizeCode             = "replace_me"
	ErrPageNumberCode           = "replace_me"
	ErrResultIDCode             = "replace_me"
	ErrPerfIDCode               = "replace_me"
	ErrMarshalCode              = "replace_me"
	ErrUnmarshalCode            = "replace_me"
	ErrGenerateUUIDCode         = "replace_me"
	ErrLocalProviderSupportCode = "replace_me"
)

var (
	ErrResultID             = errors.New(ErrResultIDCode, errors.Alert, []string{"Given resultID is not valid"}, []string{"Given resultID is nil"}, []string{}, []string{})
	ErrLocalProviderSupport = errors.New(ErrLocalProviderSupportCode, errors.Alert, []string{"Method not supported by local provider"}, []string{}, []string{}, []string{})
)

func ErrGrafanaClient(err error) error {
	return errors.New(ErrGrafanaClientCode, errors.Alert, []string{"Unable to initialize Grafana Client"}, []string{"Unable to initializes client for interacting with an instance of Grafana server", err.Error()}, []string{"Invalid Grafana Endpoint or API-Key"}, []string{"Update your Grafana URL and API-Key from the settings page in the UI"})
}

func ErrPageSize(err error) error {
	return errors.New(ErrPageSizeCode, errors.Alert, []string{"Unable to prase the Page Size"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPageNumber(err error) error {
	return errors.New(ErrPageNumberCode, errors.Alert, []string{"Unable to prase the Page Numer"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPerfID(err error) error {
	return errors.New(ErrPerfIDCode, errors.Alert, []string{"Invalid peformance profile ID"}, []string{err.Error()}, []string{}, []string{})
}

func ErrMarshal(err error, obj string) error {
	return errors.New(ErrMarshalCode, errors.Alert, []string{"Unable to marshal the object", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrUnmarshal(err error, obj string) error {
	return errors.New(ErrUnmarshalCode, errors.Alert, []string{"Unable to unmarshal the object", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrGenerateUUID(err error) error {
	return errors.New(ErrGenerateUUIDCode, errors.Alert, []string{"Unable to generate a new UUID"}, []string{err.Error()}, []string{}, []string{})
}
