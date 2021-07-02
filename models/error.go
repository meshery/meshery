// Package models - Error codes for Meshery models
package models

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrGrafanaClientCode = "3000"
)

func ErrGrafanaClient(err error) error {
	return errors.New(ErrGrafanaClientCode, errors.Alert, []string{"Unable to initialize Grafana Client"}, []string{"Unable to initializes client for interacting with an instance of Grafana server", err.Error()}, []string{"Invalid Grafana Endpoint or API-Key"}, []string{"Update your Grafana URL and API-Key from the settings page in the UI"})
}
