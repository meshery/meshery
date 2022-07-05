package model

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrInvalidRequestCode                       = "1000"
	ErrNilClientCode                            = "1001"
	ErrCreateDataCode                           = "1002"
	ErrQueryCode                                = "1003"
	ErrMeshsyncSubscriptionCode                 = "1004"
	ErrOperatorSubscriptionCode                 = "1005"
	ErrAddonSubscriptionCode                    = "1006"
	ErrControlPlaneSubscriptionCode             = "1007"
	ErrMesheryClientCode                        = "1008"
	ErrSubscribeChannelCode                     = "1009"
	ErrPublishBrokerCode                        = "1010"
	ErrNoMeshSyncCode                           = "1011"
	ErrNoExternalEndpointCode                   = "1012"
	ErrApplyHelmChartCode                       = "2155"
	ErrMesheryControllersStatusSubscriptionCode = "2239"
	ErrMeshSyncEventsSubscriptionCode           = "2240"
)

var (
	ErrEmptyHandler = errors.New(ErrNoMeshSyncCode, errors.Alert, []string{"Database handler not initialized"}, []string{"Meshery Database handler is not accessible to perform operations"}, []string{"Meshery Database is crashed or not reachable"}, []string{"Restart Meshery Server", "Please check if Meshery server is accessible to the Database"})
)

func ErrMeshSyncEventsSubscription(err error) error {
	return errors.New(ErrMeshSyncEventsSubscriptionCode, errors.Alert, []string{"Could not create MeshSync events subcription", err.Error()}, []string{"Meshery controller handlers may not be available"}, []string{""}, []string{"Make sure the controllers are deployed and their handlers are configured"})
}

func ErrMesheryControllersStatusSubscription(err error) error {
	return errors.New(ErrMesheryControllersStatusSubscriptionCode, errors.Alert, []string{"Could not create meshery controllers status subcription", err.Error()}, []string{"Meshery controller handlers may not be available"}, []string{""}, []string{"Make sure the controllers are deployed and their handlers are configured"})
}

func ErrCreateData(err error) error {
	return errors.New(ErrCreateDataCode, errors.Alert, []string{"Error while writing meshsync data", err.Error()}, []string{"Unable to write MeshSync data to the Meshery Database"}, []string{"Meshery Database is crashed or not reachable"}, []string{"Restart Meshery Server", "Please check if Meshery server is accessible to the Database"})
}

func ErrUpdateData(err error) error {
	return errors.New(ErrCreateDataCode, errors.Alert, []string{"Error while updating meshsync data", err.Error()}, []string{"Unable to update MeshSync data to the Meshery Database"}, []string{"Meshery Database is crashed or not reachable"}, []string{"Restart Meshery Server", "Please check if Meshery server is accessible to the Database"})
}

func ErrDeleteData(err error) error {
	return errors.New(ErrCreateDataCode, errors.Alert, []string{"Error while deleting meshsync data", err.Error()}, []string{"Unable to read MeshSync data to the Meshery Database"}, []string{"Meshery Database is crashed or not reachable"}, []string{"Restart Meshery Server", "Please check if Meshery server is accessible to the Database"})
}

func ErrQuery(err error) error {
	return errors.New(ErrQueryCode, errors.Alert, []string{"Error while querying data", err.Error()}, []string{"Invalid Query performed in Meshery Database"}, []string{}, []string{})
}

func ErrMeshsyncSubscription(err error) error {
	return errors.New(ErrMeshsyncSubscriptionCode, errors.Alert, []string{"MeshSync Subscription failed", err.Error()}, []string{"GraphQL subscription for MeshSync stopped"}, []string{"Could be a network issue"}, []string{"Check if meshery server is reachable from the browser"})
}

func ErrSubscribeChannel(err error) error {
	return errors.New(ErrSubscribeChannelCode, errors.Alert, []string{"Unable to subscribe to channel", err.Error()}, []string{"Unable to create a broker subscription"}, []string{"Could be a network issue", "Meshery Broker could have crashed"}, []string{"Check if Meshery Broker is reachable from Meshery Server", "Check if Meshery Broker is up and running inside the configured cluster"})
}

func ErrPublishBroker(err error) error {
	return errors.New(ErrPublishBrokerCode, errors.Alert, []string{"Unable to publish to broker", err.Error()}, []string{"Unable to create a broker publisher"}, []string{"Could be a network issue", "Meshery Broker could have crashed"}, []string{"Check if Meshery Broker is reachable from Meshery Server", "Check if Meshery Broker is up and running inside the configured cluster"})
}

func ErrMesheryClient(err error) error {
	if err != nil {
		return errors.New(ErrMesheryClientCode, errors.Alert, []string{"Meshery kubernetes client not initialized", err.Error()}, []string{"Kubernetes config is not initialized with Meshery 1"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	}
	return errors.New(ErrMesheryClientCode, errors.Alert, []string{"Meshery kubernetes client not initialized"}, []string{"Kubernetes config is not initialized with Meshery 2"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
}

// ErrApplyHelmChart is the error which occurs while applying helm chart
func ErrApplyHelmChart(err error) error {
	return errors.New(ErrApplyHelmChartCode, errors.Alert, []string{"Error occurred while applying Helm Chart"}, []string{err.Error()}, []string{"Kubernetes cluster might not be connected", "Leftover resources from partial install"}, []string{"Try reinstalling", "Try reconnecting your kubernetes cluster", "Clean up artifacts from preinstalled helm release manually"})
}
