package handlers

import (
	"github.com/meshery/meshkit/errors"
)

const (
	ErrUpgradingWebSocketCode      = "meshery-server-1400"
	ErrInvalidK8sExecParamsCode    = "meshery-server-1401"
	ErrGettingK8sClientCode        = "meshery-server-1402"
	ErrCreatingSPDYExecutorCode    = "meshery-server-1403"
	ErrStreamingExecCode           = "meshery-server-1404"
	ErrInvalidK8sLogsParamsCode    = "meshery-server-1405"
	ErrGettingPodLogsCode          = "meshery-server-1406"
	ErrWritingToWebSocketCode      = "meshery-server-1407"
	ErrReadingLogsCode             = "meshery-server-1408"
	ErrGettingK8sContextCode       = "meshery-server-1409"
	ErrGeneratingKubeConfigCode    = "meshery-server-1410"
	ErrCreatingK8sClientCode       = "meshery-server-1411"
)

func ErrUpgradingWebSocket(err error) error {
	return errors.New(ErrUpgradingWebSocketCode, errors.Alert, []string{"Failed to upgrade HTTP connection to WebSocket"}, []string{err.Error()}, []string{"Client may not support WebSocket protocol", "CORS configuration issue"}, []string{"Ensure client browser supports WebSocket", "Check WebSocket configuration"})
}

func ErrInvalidK8sExecParams() error {
	return errors.New(ErrInvalidK8sExecParamsCode, errors.Alert, []string{"Invalid parameters for Kubernetes exec"}, []string{"Required parameters (namespace, pod, context) are missing"}, []string{"Missing query parameters in request"}, []string{"Provide namespace, pod, and context parameters"})
}

func ErrGettingK8sClient(err error) error {
	return errors.New(ErrGettingK8sClientCode, errors.Alert, []string{"Failed to get Kubernetes client"}, []string{err.Error()}, []string{"Invalid Kubernetes context", "Connection to cluster failed"}, []string{"Verify Kubernetes context is valid and accessible"})
}

func ErrCreatingSPDYExecutor(err error) error {
	return errors.New(ErrCreatingSPDYExecutorCode, errors.Alert, []string{"Failed to create SPDY executor for pod exec"}, []string{err.Error()}, []string{"Invalid Kubernetes API server configuration", "Network connectivity issue"}, []string{"Check Kubernetes API server accessibility"})
}

func ErrStreamingExec(err error) error {
	return errors.New(ErrStreamingExecCode, errors.Alert, []string{"Error during exec streaming"}, []string{err.Error()}, []string{"Pod terminated", "Network connection lost", "Container not running"}, []string{"Check pod status and network connectivity"})
}

func ErrInvalidK8sLogsParams() error {
	return errors.New(ErrInvalidK8sLogsParamsCode, errors.Alert, []string{"Invalid parameters for Kubernetes logs"}, []string{"Required parameters (namespace, pod, context) are missing"}, []string{"Missing query parameters in request"}, []string{"Provide namespace, pod, and context parameters"})
}

func ErrGettingPodLogs(err error) error {
	return errors.New(ErrGettingPodLogsCode, errors.Alert, []string{"Failed to get pod logs"}, []string{err.Error()}, []string{"Pod not found", "Container not running", "Insufficient permissions"}, []string{"Verify pod exists and is running", "Check RBAC permissions"})
}

func ErrWritingToWebSocket(err error) error {
	return errors.New(ErrWritingToWebSocketCode, errors.Alert, []string{"Failed to write to WebSocket"}, []string{err.Error()}, []string{"WebSocket connection closed", "Network issue"}, []string{"Check WebSocket connection status"})
}

func ErrReadingLogs(err error) error {
	return errors.New(ErrReadingLogsCode, errors.Alert, []string{"Error reading logs"}, []string{err.Error()}, []string{"Log stream interrupted", "Pod terminated"}, []string{"Check pod status"})
}

func ErrGettingK8sContext(err error) error {
	return errors.New(ErrGettingK8sContextCode, errors.Alert, []string{"Failed to get Kubernetes context"}, []string{err.Error()}, []string{"Context not found in database", "Invalid context ID"}, []string{"Verify context ID is correct"})
}

func ErrGeneratingKubeConfig(err error) error {
	return errors.New(ErrGeneratingKubeConfigCode, errors.Alert, []string{"Failed to generate Kubernetes config"}, []string{err.Error()}, []string{"Invalid kubeconfig data", "Missing credentials"}, []string{"Check context configuration and credentials"})
}

func ErrCreatingK8sClient(err error) error {
	return errors.New(ErrCreatingK8sClientCode, errors.Alert, []string{"Failed to create Kubernetes client"}, []string{err.Error()}, []string{"Invalid configuration", "Network connectivity issue"}, []string{"Verify Kubernetes configuration is valid"})
}

