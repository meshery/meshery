package resolver

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrInvalidRequestCode           = "1000"
	ErrNilClientCode                = "1001"
	ErrCreateDataCode               = "1002"
	ErrQueryCode                    = "1003"
	ErrMeshsyncSubscriptionCode     = "1004"
	ErrOperatorSubscriptionCode     = "1005"
	ErrAddonSubscriptionCode        = "1006"
	ErrControlPlaneSubscriptionCode = "1007"
	ErrMesheryClientCode            = "1008"
	ErrSubscribeChannelCode         = "1009"
	ErrPublishBrokerCode            = "1010"
	ErrNoMeshSyncCode               = "1011"
	ErrNoExternalEndpointCode       = "1012"
)

var (
	ErrNilClient      = errors.New(ErrNilClientCode, errors.Alert, []string{"Kubernetes client not initialized"}, []string{"Kubernetes config is not initialized with Meshery"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrInvalidRequest = errors.New(ErrInvalidRequestCode, errors.Alert, []string{"Invalid query, please check syntax"}, []string{"The Graphql query requested is invalid"}, []string{}, []string{"Check the query parameters and syntax of the query being run"})
	ErrNoMeshSync     = errors.New(ErrNoMeshSyncCode, errors.Alert, []string{"MeshSync disabled"}, []string{"MeshSync custom controller is not running in your kubernetes cluster"}, []string{"Meshery Operator is not running in your cluster or is crashing"}, []string{"Enable Meshery Operator from the settings page in the UI", "Check for logs in the meshery-operator pods from inside the application for more information"})
	ErrEmptyHandler   = errors.New(ErrNoMeshSyncCode, errors.Alert, []string{"Database handler not initialized"}, []string{"Meshery Database handler is not accessible to perform operations"}, []string{"Meshery Database is crashed or not reachable"}, []string{"Restart Meshery Server", "Please check if Meshery server is accessible to the Database"})
)

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

func ErrOperatorSubscription(err error) error {
	return errors.New(ErrOperatorSubscriptionCode, errors.Alert, []string{"Operator Subscription failed", err.Error()}, []string{"GraphQL subscription for Operator stopped"}, []string{"Could be a network issue"}, []string{"Check if meshery server is reachable from the browser"})
}

func ErrAddonSubscription(err error) error {
	return errors.New(ErrAddonSubscriptionCode, errors.Alert, []string{"Addons Subscription failed", err.Error()}, []string{"GraphQL subscription for Addons stopped"}, []string{"Could be a network issue"}, []string{"Check if meshery server is reachable from the browser"})
}

func ErrControlPlaneSubscription(err error) error {
	return errors.New(ErrControlPlaneSubscriptionCode, errors.Alert, []string{"Control Plane Subscription failed", err.Error()}, []string{"GraphQL subscription for Control Plane stopped"}, []string{"Could be a network issue"}, []string{"Check if meshery server is reachable from the browser"})
}

func ErrSubscribeChannel(err error) error {
	return errors.New(ErrSubscribeChannelCode, errors.Alert, []string{"Unable to subscribe to channel", err.Error()}, []string{"Unable to create a broker subscription"}, []string{"Could be a network issue", "Meshery Broker could have crashed"}, []string{"Check if Meshery Broker is reachable from Meshery Server", "Check if Meshery Broker is up and running inside the configured cluster"})
}

func ErrPublishBroker(err error) error {
	return errors.New(ErrPublishBrokerCode, errors.Alert, []string{"Unable to publish to broker", err.Error()}, []string{"Unable to create a broker publisher"}, []string{"Could be a network issue", "Meshery Broker could have crashed"}, []string{"Check if Meshery Broker is reachable from Meshery Server", "Check if Meshery Broker is up and running inside the configured cluster"})
}

func ErrMesheryClient(err error) error {
	if err != nil {
		return errors.New(ErrMesheryClientCode, errors.Alert, []string{"Meshery kubernetes client not initialized", err.Error()}, []string{"Kubernetes config is not initialized with Meshery"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	}
	return errors.New(ErrMesheryClientCode, errors.Alert, []string{"Meshery kubernetes client not initialized"}, []string{"Kubernetes config is not initialized with Meshery"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
}

func ErrNoExternalEndpoint(service string) error {
	return errors.New(ErrNoExternalEndpointCode, errors.Alert, []string{"Unable to find the External Endpoint"}, []string{"Unable to find External Endpoint for Service: " + service}, []string{"The service might not be exposed to External IP", "Meshsync cache is outdated"}, []string{"Check if the Service " + service + " is reachable Externally"})
}
