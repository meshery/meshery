// Package handlers - Error codes for MEshery handlers
package handlers

import (
	"github.com/layer5io/meshkit/errors"
	"github.com/layer5io/meshkit/utils"
)

const (
	ErrInvalidK8SConfigCode  = "1000"
	ErrNilClientCode         = "1001"
	ErrPrometheusScanCode    = "1002"
	ErrGrafanaScanCode       = "1003"
	ErrRecordPreferencesCode = "1004"
	ErrGrafanaConfigCode     = "1005"
	ErrPrometheusConfigCode  = "1006"
	ErrGrafanaQueryCode      = "1007"
	ErrPrometheusQueryCode   = "1008"
	ErrGrafanaBoardsCode     = "1009"
	ErrPrometheusBoardsCode  = "1010"
	ErrStaticBoardsCode      = "1011"
	ErrRequestBodyCode       = "1012"
)

var (
	ErrInvalidK8SConfig = errors.New(ErrInvalidK8SConfigCode, errors.Alert, []string{"No valid kubernetes config found"}, []string{"Kubernetes config is not initialized with Meshery"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrNilClient        = errors.New(ErrNilClientCode, errors.Alert, []string{"Kubernetes client not initialized"}, []string{"Kubernetes config is not initialized with Meshery"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrPrometheusConfig = errors.New(ErrGrafanaConfigCode, errors.Alert, []string{"Prometheus endpoint not configured."}, []string{"Cannot find valid Prometheus endpoint in user pref"}, []string{}, []string{"Setup your Prometheus Endpoint via the settings dashboard."})
	ErrGrafanaConfig    = errors.New(ErrGrafanaConfigCode, errors.Alert, []string{"Grafana endpoint not configured."}, []string{"Cannot find valid grafana endpoint in user pref"}, []string{}, []string{"Setup your Grafana Endpoint via the settings dashboard."})
	ErrStaticBoards     = errors.New(ErrStaticBoardsCode, errors.Alert, []string{"unable to get static board"}, []string{"unable to get static board"}, []string{}, []string{})
)

func ErrPrometheusScan(err error) error {
	return errors.New(ErrPrometheusScanCode, errors.Alert, []string{"Unable to connect to prometheus."}, []string{"Failed connecting to Prometheus Client through Grafna Proxy", err.Error()}, []string{}, []string{"Check if your Prometheus and Grafana Endpoint are correct.", "Connect to Prometheus and Grafana from the settings page in the UI"})
}

func ErrGrafanaScan(err error) error {
	return errors.New(ErrGrafanaScanCode, errors.Alert, []string{"Unable to connect to grafana."}, []string{"Failed connecting to Grafna Client", err.Error()}, []string{}, []string{"Check if your Grafana Endpoint is correct.", "Connect to Grafana from the settings page in the UI"})
}

func ErrPrometheusQuery(err error) error {
	return errors.New(ErrPrometheusQueryCode, errors.Alert, []string{"Unable to query prometheus. " + err.Error()}, []string{"Failed querying Prometheus Client through Grafna Proxy", err.Error()}, []string{}, []string{"Check if your Prometheus query is correct.", "Connect to Prometheus and Grafana from the settings page in the UI"})
}

func ErrGrafanaQuery(err error) error {
	return errors.New(ErrGrafanaQueryCode, errors.Alert, []string{"Unable to query grafana." + err.Error()}, []string{"Failed querying Grafana Client", err.Error()}, []string{}, []string{"Check if your Grafana query is correct.", "Connect to Grafana from the settings page in the UI"})
}

func ErrGrafanaBoards(err error) error {
	return errors.New(ErrGrafanaBoardsCode, errors.Alert, []string{"unable to get grafana boards." + err.Error()}, []string{"Failed querying boards in Grafana Client. " + err.Error()}, []string{}, []string{"Check if your Grafana endpoint is correct.", "Connect to Grafana from the settings page in the UI"})
}

func ErrPrometheusBoards(err error) error {
	return errors.New(ErrGrafanaBoardsCode, errors.Alert, []string{"unable to get Prometheus boards." + err.Error()}, []string{"Failed querying boards for Prometheus Client. " + err.Error()}, []string{}, []string{"Check if your Prometheus endpoint is correct.", "Connect to Prometheus from the settings page in the UI"})
}

func ErrRecordPreferences(err error) error {
	return errors.New(ErrRecordPreferencesCode, errors.Alert, []string{"unable to save user config data:", err.Error()}, []string{"Failed to record User Preferences for ", err.Error()}, []string{"User token is invalid"}, []string{"Relogin to Meshery"})
}

func ErrRequestBody(err error) error {
	return errors.New(ErrRequestBodyCode, errors.Alert, []string{"unable to read the request body.", err.Error()}, []string{"unable to read the request body" + err.Error()}, []string{}, []string{})
}

func ErrMarshal(err error) error {
	return utils.ErrMarshal(err)
}

func ErrUnmarshal(err error) error {
	return utils.ErrUnmarshal(err)
}
