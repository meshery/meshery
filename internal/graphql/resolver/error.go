package resolver

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrResolverInvalidRequestCode           = "2222"
	ErrResolverNilClientCode                = "2223"
	ErrResolverCreateDataCode               = "2224"
	ErrResolverQueryCode                    = "2225"
	ErrResolverMeshsyncSubscriptionCode     = "2226"
	ErrResolverOperatorSubscriptionCode     = "2227"
	ErrResolverAddonSubscriptionCode        = "2228"
	ErrResolverControlPlaneSubscriptionCode = "2229"
	ErrResolverMesheryClientCode            = "2230"
	ErrResolverSubscribeChannelCode         = "2231"
	ErrResolverPublishBrokerCode            = "2232"
	ErrResolverNoMeshSyncCode               = "2233"
	ErrResolverNoExternalEndpointCode       = "2234"
	ErrDataPlaneSubscriptionCode            = "1013"
	ErrBrokerNotConnectedCode               = "2221"
	ErrGettingNamespaceCode                 = "1014"
	ErrFetchingPatternsCode                 = "1015"
	ErrInvalidOAMTypeCode                   = "2178"
	ErrKubectlDescribeCode                  = "2179"
	ErrEmptyCurrentK8sContextCode           = "2180"
)

var (
	ErrNilClient              = errors.New(ErrResolverNilClientCode, errors.Alert, []string{"Kubernetes client not initialized"}, []string{"Kubernetes config is not initialized with Meshery"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrInvalidRequest         = errors.New(ErrResolverInvalidRequestCode, errors.Alert, []string{"Invalid query, please check syntax"}, []string{"The Graphql query requested is invalid"}, []string{}, []string{"Check the query parameters and syntax of the query being run"})
	ErrNoMeshSync             = errors.New(ErrResolverNoMeshSyncCode, errors.Alert, []string{"MeshSync disabled"}, []string{"MeshSync custom controller is not running in your kubernetes cluster"}, []string{"Meshery Operator is not running in your cluster or is crashing"}, []string{"Enable Meshery Operator from the settings page in the UI", "Check for logs in the meshery-operator pods from inside the application for more information"})
	ErrEmptyHandler           = errors.New(ErrResolverNoMeshSyncCode, errors.Alert, []string{"Database handler not initialized"}, []string{"Meshery Database handler is not accessible to perform operations"}, []string{"Meshery Database is crashed or not reachable"}, []string{"Restart Meshery Server", "Please check if Meshery server is accessible to the Database"})
	ErrBrokerNotConnected     = errors.New(ErrBrokerNotConnectedCode, errors.Alert, []string{"Broker not connected"}, []string{"Meshery Broker is not connected to Meshery Server"}, []string{"Meshery Broker is crashed or not reachable"}, []string{"Restart Meshery Server", "Please check if Meshery server has access to the Broker"})
	ErrEmptyCurrentK8sContext = errors.New(ErrEmptyCurrentK8sContextCode, errors.Alert, []string{"Current kubernetes context is empty"}, []string{"Meshery is not able to get the current kubernetes context"}, []string{"Meshery is crashed or not reachable"}, []string{"Restart Meshery Server", "Please check if Meshery server is accessible to the Database"})
)

func ErrUpdateData(err error) error {
	return errors.New(ErrResolverCreateDataCode, errors.Alert, []string{"Error while updating meshsync data", err.Error()}, []string{"Unable to update MeshSync data to the Meshery Database"}, []string{"Meshery Database is crashed or not reachable"}, []string{"Restart Meshery Server", "Please check if Meshery server is accessible to the Database"})
}

func ErrDeleteData(err error) error {
	return errors.New(ErrResolverCreateDataCode, errors.Alert, []string{"Error while deleting meshsync data", err.Error()}, []string{"Unable to read MeshSync data to the Meshery Database"}, []string{"Meshery Database is crashed or not reachable"}, []string{"Restart Meshery Server", "Please check if Meshery server is accessible to the Database"})
}

func ErrMeshsyncSubscription(err error) error {
	return errors.New(ErrResolverMeshsyncSubscriptionCode, errors.Alert, []string{"MeshSync Subscription failed", err.Error()}, []string{"GraphQL subscription for MeshSync stopped"}, []string{"Could be a network issue"}, []string{"Check if meshery server is reachable from the browser"})
}

func ErrOperatorSubscription(err error) error {
	return errors.New(ErrResolverOperatorSubscriptionCode, errors.Alert, []string{"Operator Subscription failed", err.Error()}, []string{"GraphQL subscription for Operator stopped"}, []string{"Could be a network issue"}, []string{"Check if meshery server is reachable from the browser"})
}

func ErrAddonSubscription(err error) error {
	return errors.New(ErrResolverAddonSubscriptionCode, errors.Alert, []string{"Addons Subscription failed", err.Error()}, []string{"GraphQL subscription for Addons stopped"}, []string{"Could be a network issue"}, []string{"Check if meshery server is reachable from the browser"})
}

func ErrPerformanceProfilesSubscription(err error) error {
	return errors.New(ErrResolverAddonSubscriptionCode, errors.Alert, []string{"PerformanceProfiles Subscription failed", err.Error()}, []string{"GraphQL subscription for PerformanceProfiles stopped"}, []string{"Could be a network issue"}, []string{"Confirm that Meshery Server is reachable from your browser."})
}

func ErrPerformanceResultSubscription(err error) error {
	return errors.New(ErrResolverAddonSubscriptionCode, errors.Alert, []string{"PerformanceResult Subscription failed", err.Error()}, []string{"GraphQL subscription for PerformanceResult stopped"}, []string{"Could be a network issue"}, []string{"Confirm that Meshery Server is reachable from your browser."})
}

func ErrControlPlaneSubscription(err error) error {
	return errors.New(ErrResolverControlPlaneSubscriptionCode, errors.Alert, []string{"Control Plane Subscription failed", err.Error()}, []string{"GraphQL subscription for Control Plane stopped"}, []string{"Could be a network issue"}, []string{"Confirm that Meshery Server is reachable from your browser."})
}

func ErrDataPlaneSubscription(err error) error {
	return errors.New(ErrDataPlaneSubscriptionCode, errors.Alert, []string{"Data Plane Subscription failed", err.Error()}, []string{"GraphQL subscription for Data Plane stopped"}, []string{"Could be a network issue"}, []string{"Check if meshery server is reachable from the browser"})
}

func ErrPublishBroker(err error) error {
	return errors.New(ErrResolverPublishBrokerCode, errors.Alert, []string{"Unable to publish to broker", err.Error()}, []string{"Unable to create a broker publisher"}, []string{"Could be a network issue", "Meshery Broker could have crashed"}, []string{"Check if Meshery Broker is reachable from Meshery Server", "Check if Meshery Broker is up and running inside the configured cluster"})
}

func ErrMesheryClient(err error) error {
	if err != nil {
		return errors.New(ErrResolverMesheryClientCode, errors.Alert, []string{"Meshery kubernetes client not initialized", err.Error()}, []string{"Kubernetes config is not initialized with Meshery 4"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	}
	return errors.New(ErrResolverMesheryClientCode, errors.Alert, []string{"Meshery kubernetes client not initialized"}, []string{"Kubernetes config is not initialized with Meshery 5"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
}
func ErrGettingNamespace(err error) error {
	return errors.New(ErrGettingNamespaceCode, errors.Alert, []string{"Cannot get available namespaces"}, []string{err.Error()}, []string{"The table in the database might not exist"}, []string{})
}
func ErrFetchingPatterns(err error) error {
	return errors.New(ErrFetchingPatternsCode, errors.Alert, []string{"Cannot fetch patterns"}, []string{err.Error()}, []string{"There might be something wrong with the Meshery or Meshery Cloud"}, []string{"Try again, if still exist, please post an issue on Meshery repository"})
}

func ErrInvalidOAMType() error {
	return errors.New(
		ErrInvalidOAMTypeCode,
		errors.Alert,
		[]string{"invalid oam type is requested"},
		[]string{"invalid oam type requested, supported types are workload, scope, trait "},
		nil, nil,
	)
}

func ErrKubectlDescribe(err error) error {
	return errors.New(
		ErrKubectlDescribeCode,
		errors.Alert,
		[]string{"failed to find the resource", "invalid resource type", "resource doens't exists"},
		[]string{err.Error(), "invalid kubernetes resource type or couldn't find the specified resource"},
		nil, nil,
	)
}
