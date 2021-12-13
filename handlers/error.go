package handlers

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrInvalidK8SConfigCode     = "2000"
	ErrNilClientCode            = "2001"
	ErrPrometheusScanCode       = "2002"
	ErrGrafanaScanCode          = "2003"
	ErrRecordPreferencesCode    = "2004"
	ErrGrafanaConfigCode        = "2005"
	ErrPrometheusConfigCode     = "2006"
	ErrGrafanaQueryCode         = "2007"
	ErrPrometheusQueryCode      = "2008"
	ErrGrafanaBoardsCode        = "2009"
	ErrPrometheusBoardsCode     = "2010"
	ErrStaticBoardsCode         = "2011"
	ErrRequestBodyCode          = "2012"
	ErrMarshalCode              = "2013"
	ErrUnmarshalCode            = "2014"
	ErrEncodingCode             = "2015"
	ErrParseBoolCode            = "2016"
	ErrStreamEventsCode         = "2017"
	ErrStreamClientCode         = "2018"
	ErrUnmarshalEventCode       = "2019"
	ErrPublishSmiResultsCode    = "2020"
	ErrMarshalEventCode         = "2021"
	ErrPluginOpenCode           = "2022"
	ErrPluginLookupCode         = "2023"
	ErrPluginRunCode            = "2024"
	ErrParseFormCode            = "2025"
	ErrQueryGetCode             = "2026"
	ErrGetResultCode            = "2027"
	ErrConvertToSpecCode        = "2028"
	ErrFetchSMIResultsCode      = "2029"
	ErrFormFileCode             = "2030"
	ErrReadConfigCode           = "2031"
	ErrLoadConfigCode           = "2032"
	ErrOpenFileCode             = "2033"
	ErrKubeVersionCode          = "2034"
	ErrAddAdapterCode           = "2035"
	ErrRetrieveDataCode         = "2036"
	ErrValidAdapterCode         = "2037"
	ErrOperationIDCode          = "2038"
	ErrMeshClientCode           = "2039"
	ErrApplyChangeCode          = "2040"
	ErrRetrieveMeshDataCode     = "2041"
	ErrApplicationFailureCode   = "2042"
	ErrDecodingCode             = "2043"
	ErrRetrieveUserTokenCode    = "2044"
	ErrFailToSaveCode           = "2045"
	ErrFailToDeleteCode         = "2046"
	ErrWriteResponseCode        = "2052"
	ErrTestConfigsCode          = "2053"
	ErrInvalidGenValueCode      = "2054"
	ErrFailToLoadExtensionsCode = "2047"
	ErrConversionCode           = "2048"
	ErrParseDurationCode        = "2049"
	ErrLoadTestCode             = "2050"
	ErrFetchKubernetesCode      = "2051"
	ErrPanicRecoveryCode        = "2052"
	ErrBlankNameCode            = "2053"
	ErrInvalidLTURLCode         = "2053"
	ErrDataSendCode             = "2137"
	ErrVersionCompareCode       = "2138"
	ErrSaveSessionCode          = "2136"
	ErrKubeClientCode           = "2139"
	ErrWorkloadDefinitionCode   = "2140"
	ErrTraitDefinitionCode      = "2141"
	ErrScopeDefinitionCode      = "2142"
	ErrPatternFileCode          = "2143"
	ErrExecutionPlanCode        = "2144"
	ErrInvalidPatternCode       = "2145"
	ErrCompConfigPairsCode      = "2146"
	ErrCreateDirCode            = "2150"
	ErrInvalidRequestObjectCode = "2151"
	ErrChangeK8sContextCode     = "2152"
	ErrSavingUserPreferenceCode = "some_code"
)

var (
	ErrInvalidK8SConfig = errors.New(ErrInvalidK8SConfigCode, errors.Alert, []string{"No valid kubernetes config found"}, []string{"Kubernetes config is not initialized with Meshery"}, []string{"Kubernetes config is not accessible to meshery or not valid"}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrNilClient        = errors.New(ErrNilClientCode, errors.Alert, []string{"Kubernetes client not initialized"}, []string{"Kubernetes config is not initialized with Meshery"}, []string{"Kubernetes config is not accessible to meshery or not valid"}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrPrometheusConfig = errors.New(ErrPrometheusConfigCode, errors.Alert, []string{"Prometheus endpoint not configured"}, []string{"Cannot find valid Prometheus endpoint in user pref"}, []string{"Prometheus endpoint might not be reachable from meshery"}, []string{"Setup your Prometheus Endpoint via the settings dashboard"})
	ErrGrafanaConfig    = errors.New(ErrGrafanaConfigCode, errors.Alert, []string{"Grafana endpoint not configured"}, []string{"Cannot find valid grafana endpoint in user pref"}, []string{"Grafana endpoint might not be reachable from meshery"}, []string{"Setup your Grafana Endpoint via the settings dashboard"})
	ErrStaticBoards     = errors.New(ErrStaticBoardsCode, errors.Alert, []string{"unable to get static board"}, []string{"unable to get static board"}, []string{"No boards could be available in grafana"}, []string{})
	ErrValidAdapter     = errors.New(ErrValidAdapterCode, errors.Alert, []string{"Unable to find valid Adapter URL"}, []string{"unable to find a valid adapter for the given adapter URL"}, []string{"Given adapter URL is not valid"}, []string{"Please provide a valid Adapter URL"})
	ErrAddAdapter       = errors.New(ErrAddAdapterCode, errors.Alert, []string{"meshLocationURL is empty"}, []string{"meshLocationURL is empty to add an adapter"}, []string{"meshLocationURL cannot be empty to add an adapter"}, []string{"please provide the meshLocationURL"})
	ErrMeshClient       = errors.New(ErrMeshClientCode, errors.Alert, []string{"Error creating a mesh client", "Error pinging the mesh adapter"}, []string{"Unable to create a mesh client", "Unable to ping the mesh adapter"}, []string{"Adapter could not be pinged"}, []string{"Unable to connect to the Mesh adapter using the given config, please try again"})
	ErrWriteResponse    = errors.New(ErrWriteResponseCode, errors.Alert, []string{"Error writing response"}, []string{}, []string{}, []string{})
	ErrTestConfigs      = errors.New(ErrTestConfigsCode, errors.Alert, []string{"Error fetching test configs"}, []string{}, []string{}, []string{})
	ErrInvalidGenValue  = errors.New(ErrInvalidGenValueCode, errors.Alert, []string{"Invalid value for gen"}, []string{}, []string{}, []string{"please provide a valid value for gen (load generator)"})
	ErrParseDuration    = errors.New(ErrParseDurationCode, errors.Alert, []string{"error parsing test duration"}, []string{}, []string{"The format of the duration passed could be incorrect"}, []string{"please refer to:  https://docs.meshery.io/guides/mesheryctl#performance-management"})
)

func ErrPrometheusScan(err error) error {
	return errors.New(ErrPrometheusScanCode, errors.Alert, []string{"Unable to connect to prometheus"}, []string{err.Error()}, []string{"Prometheus endpoint might not be reachable from meshery", "Prometheus endpoint is incorrect"}, []string{"Check if your Prometheus and Grafana Endpoint are correct", "Connect to Prometheus and Grafana from the settings page in the UI"})
}

func ErrGrafanaScan(err error) error {
	return errors.New(ErrGrafanaScanCode, errors.Alert, []string{"Unable to connect to grafana"}, []string{err.Error()}, []string{"Grafana endpoint might not be reachable from meshery", "Grafana endpoint is incorrect"}, []string{"Check if your Grafana Endpoint is correct", "Connect to Grafana from the settings page in the UI"})
}

func ErrPrometheusQuery(err error) error {
	return errors.New(ErrPrometheusQueryCode, errors.Alert, []string{"Unable to query prometheus"}, []string{err.Error()}, []string{"Prometheus query did not get executed from meshery", "Prometheus query is invalid"}, []string{"Check if your Prometheus query is correct", "Connect to Prometheus and Grafana from the settings page in the UI"})
}

func ErrGrafanaQuery(err error) error {
	return errors.New(ErrGrafanaQueryCode, errors.Alert, []string{"Unable to query grafana"}, []string{err.Error()}, []string{"Grafana query did not get executed from meshery", "Grafana query is invalid"}, []string{"Check if your Grafana query is correct", "Connect to Grafana from the settings page in the UI"})
}

func ErrGrafanaBoards(err error) error {
	return errors.New(ErrGrafanaBoardsCode, errors.Alert, []string{"unable to get grafana boards"}, []string{err.Error()}, []string{"Grafana endpoint might not be reachable from meshery", "Grafana endpoint is incorrect"}, []string{"Check if your Grafana endpoint is correct", "Connect to Grafana from the settings page in the UI"})
}

func ErrPrometheusBoards(err error) error {
	return errors.New(ErrGrafanaBoardsCode, errors.Alert, []string{"unable to get Prometheus boards"}, []string{err.Error()}, []string{"Prometheus endpoint might not be reachable from meshery", "Prometheus endpoint is incorrect"}, []string{"Check if your Prometheus endpoint is correct", "Connect to Prometheus from the settings page in the UI"})
}

func ErrRecordPreferences(err error) error {
	return errors.New(ErrRecordPreferencesCode, errors.Alert, []string{"unable to save user config data"}, []string{err.Error()}, []string{"User token might be invalid", "db might be corrupted"}, []string{"Relogin to Meshery"})
}

func ErrKubeClient(err error) error {
	return errors.New(ErrKubeClientCode, errors.Alert, []string{"Failed to Create Kube Client", err.Error()}, []string{err.Error()}, []string{"Check Kubernetes"}, []string{"Check your kubeconfig if valid", "Ensure meshery is able to reach the kubernetes cluster"})
}

func ErrWorkloadDefinition(err error) error {
	return errors.New(ErrWorkloadDefinitionCode, errors.Alert, []string{"Failed to load Workload Definition", err.Error()}, []string{err.Error()}, []string{"Workload Definition is invalid or unable to process"}, []string{"Check Workload Definition"})
}

func ErrTraitDefinition(err error) error {
	return errors.New(ErrTraitDefinitionCode, errors.Alert, []string{"Failed to Encode Trait Definition", err.Error()}, []string{err.Error()}, []string{"Trait Definition is invalid or unable to process"}, []string{"Check Trait Definition"})
}

func ErrScopeDefinition(err error) error {
	return errors.New(ErrScopeDefinitionCode, errors.Alert, []string{"Failed to Encode Scope Definition", err.Error()}, []string{err.Error()}, []string{"Trait Definition is invalid or unable to process"}, []string{"Check Trait Definition"})
}

func ErrPatternFile(err error) error {
	return errors.New(ErrPatternFileCode, errors.Alert, []string{"Failed to Parse Pattern File", err.Error()}, []string{err.Error()}, []string{"Trait Definition is invalid or unable to process"}, []string{"Check Trait Definition"})
}

func ErrInvalidPattern(err error) error {
	return errors.New(ErrInvalidPatternCode, errors.Alert, []string{"Invalid Pattern, execution is infeasible", err.Error()}, []string{err.Error()}, []string{"Trait Definition is invalid or unable to process"}, []string{"Check Trait Definition"})
}

func ErrExecutionPlan(err error) error {
	return errors.New(ErrExecutionPlanCode, errors.Alert, []string{"Failed to Create Execution Plan", err.Error()}, []string{err.Error()}, []string{"Trait Definition is invalid or unable to process"}, []string{"Check Trait Definition"})
}

func ErrCompConfigPairs(err error) error {
	return errors.New(ErrRequestBodyCode, errors.Alert, []string{"unable to Create Comp Config.", err.Error()}, []string{err.Error()}, []string{}, []string{})
}

func ErrRequestBody(err error) error {
	return errors.New(ErrRequestBodyCode, errors.Alert, []string{"unable to read the request body"}, []string{err.Error()}, []string{"Request body is empty or faulty"}, []string{"Check if the request is sent with proper values"})
}

func ErrMarshal(err error, obj string) error {
	return errors.New(ErrMarshalCode, errors.Alert, []string{"Unable to marshal the : ", obj}, []string{err.Error()}, []string{"Object is not a valid json object"}, []string{"Make sure if the object passed has json tags"})
}

func ErrUnmarshal(err error, obj string) error {
	return errors.New(ErrUnmarshalCode, errors.Alert, []string{"Unable to unmarshal the : ", obj}, []string{err.Error()}, []string{"Object is not a valid json object"}, []string{"Make sure if the object passed is a valid json"})
}

func ErrEncoding(err error, obj string) error {
	return errors.New(ErrEncodingCode, errors.Alert, []string{"Error encoding the : ", obj}, []string{err.Error()}, []string{"Object is not a valid json object"}, []string{"Make sure if the object passed is a valid json"})
}

func ErrParseBool(err error, obj string) error {
	return errors.New(ErrParseBoolCode, errors.Alert, []string{"unable to parse : ", obj}, []string{err.Error()}, []string{"Failed due to invalid value of : ", obj}, []string{"please provide a valid value for : ", obj})
}

func ErrStreamEvents(err error) error {
	return errors.New(ErrStreamEventsCode, errors.Alert, []string{"There was an error connecting to the backend to get events"}, []string{err.Error()}, []string{"Websocket is blocked in the network", "Meshery UI is not able to reach the Meshery server"}, []string{"Ensure Meshery UI is able to reach the Meshery server"})
}

func ErrStreamClient(err error) error {
	return errors.New(ErrStreamClientCode, errors.Alert, []string{"Event streaming ended"}, []string{err.Error()}, []string{"Websocket is blocked in the network", "Meshery UI is not able to reach the Meshery server"}, []string{"Ensure Meshery UI is able to reach the Meshery server"})
}

func ErrPublishSmiResults(err error) error {
	return errors.New(ErrPublishSmiResultsCode, errors.Alert, []string{"Error publishing SMI results"}, []string{err.Error()}, []string{"Meshery Cloud is not functional or reachable"}, []string{"Make sure meshery cloud is up and reachable"})
}

func ErrPluginOpen(err error) error {
	return errors.New(ErrPluginOpenCode, errors.Alert, []string{"Error opening the plugin"}, []string{err.Error()}, []string{"Plugin is not available in the location", "plugin does not match with meshery version"}, []string{"Make sure the plugin is compatible with Meshery server"})
}

func ErrPluginLookup(err error) error {
	return errors.New(ErrPluginLookupCode, errors.Alert, []string{"Error performing a plugin lookup"}, []string{err.Error()}, []string{"Plugin is not available in the location"}, []string{"Make sure the plugin is compatible with Meshery server"})
}

func ErrPluginRun(err error) error {
	return errors.New(ErrPluginRunCode, errors.Alert, []string{"Error running meshery plugin"}, []string{err.Error()}, []string{"plugin does not match with meshery version"}, []string{"Make sure the plugin is compatible with Meshery server"})
}

func ErrParseForm(err error) error {
	return errors.New(ErrParseFormCode, errors.Alert, []string{"unable to parse form"}, []string{err.Error()}, []string{"The data provided could be invalid"}, []string{"Make sure to enter valid parameters in the form"})
}

func ErrQueryGet(obj string) error {
	return errors.New(ErrQueryGetCode, errors.Alert, []string{"unable to get: ", obj}, []string{}, []string{"Query parameter is not a part of the request"}, []string{"Make sure to pass the query paramater in the request"})
}

func ErrGetResult(err error) error {
	return errors.New(ErrGetResultCode, errors.Alert, []string{"unable to get result"}, []string{err.Error()}, []string{"Result Identifier provided is not valid", "Result did not persist in the database"}, []string{"Make sure to provide the correct identifier for the result"})
}

func ErrConvertToSpec(err error) error {
	return errors.New(ErrConvertToSpecCode, errors.Alert, []string{"unable to convert to spec"}, []string{err.Error()}, []string{"The performance spec format is invalid"}, []string{"Make sure to provide the correct spec"})
}

func ErrFetchSMIResults(err error) error {
	return errors.New(ErrFetchSMIResultsCode, errors.Alert, []string{"unable to fetch SMI results"}, []string{err.Error()}, []string{"SMI results did not get persisted", "Result identifier is invalid"}, []string{"Make sure to provide the correct identifier for the result"})
}

func ErrFormFile(err error) error {
	return errors.New(ErrFormFileCode, errors.Alert, []string{"error getting k8s file"}, []string{err.Error()}, []string{"The kubeconfig file does not exist in the location"}, []string{"Make sure to upload the correct kubeconfig file"})
}

func ErrReadConfig(err error) error {
	return errors.New(ErrReadConfigCode, errors.Alert, []string{"error reading config"}, []string{err.Error()}, []string{"The kubeconfig file is empty or not valid"}, []string{"Make sure to upload the correct kubeconfig file"})
}

func ErrLoadConfig(err error) error {
	return errors.New(ErrLoadConfigCode, errors.Alert, []string{"unable to load kubernetes config"}, []string{err.Error()}, []string{"The kubeconfig file is empty or not valid"}, []string{"Make sure to upload the correct kubeconfig file"})
}

func ErrOpenFile(file string) error {
	return errors.New(ErrOpenFileCode, errors.Alert, []string{"unable to open file: ", file}, []string{}, []string{"The file does not exist in the location"}, []string{"Make sure to upload the correct file"})
}

func ErrKubeVersion(err error) error {
	return errors.New(ErrKubeVersionCode, errors.Alert, []string{"unable to get kubernetes version"}, []string{err.Error()}, []string{"Kubernetes might not be reachable from meshery"}, []string{"Make sure meshery has connectivity to kubernetes"})
}

func ErrRetrieveData(err error) error {
	return errors.New(ErrRetrieveDataCode, errors.Alert, []string{"Unable to retrieve the requested data"}, []string{err.Error()}, []string{"Adapter operation invalid"}, []string{"Make sure adapter is reachable and running"})
}

func ErrOperationID(err error) error {
	return errors.New(ErrOperationIDCode, errors.Alert, []string{"Error generating the operation Id"}, []string{err.Error()}, []string{"Adapter operation invalid"}, []string{"Make sure adapter is reachable and running"})
}

func ErrApplyChange(err error) error {
	return errors.New(ErrApplyChangeCode, errors.Alert, []string{"Error applying the change"}, []string{err.Error()}, []string{"Adapter operation invalid"}, []string{"Make sure adapter is reachable and running"})
}

func ErrRetrieveMeshData(err error) error {
	return errors.New(ErrRetrieveMeshDataCode, errors.Alert, []string{"Error getting operations for the mesh", "Error getting service mesh name"}, []string{err.Error()}, []string{"unable to retrieve the requested data"}, []string{"Make sure adapter is reachable and running"})
}

func ErrApplicationFailure(err error, obj string) error {
	return errors.New(ErrApplicationFailureCode, errors.Alert, []string{"failed to ", obj, "the application"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDecoding(err error, obj string) error {
	return errors.New(ErrDecodingCode, errors.Alert, []string{"Error decoding the : ", obj}, []string{err.Error()}, []string{"Object is not a valid json object"}, []string{"Make sure if the object passed is a valid json"})
}

func ErrRetrieveUserToken(err error) error {
	return errors.New(ErrRetrieveUserTokenCode, errors.Alert, []string{"Failed to get the user token"}, []string{err.Error()}, []string{"User token could be expired"}, []string{"Re-initiate login"})
}

func ErrFailToSave(err error, obj string) error {
	return errors.New(ErrFailToSaveCode, errors.Alert, []string{"Failed to Save: ", obj}, []string{err.Error()}, []string{"Meshery Database could be down or not reachable"}, []string{"Restart Meshery instance and make sure database is up and reachable"})
}
func ErrFailToDelete(err error, obj string) error {
	return errors.New(ErrFailToDeleteCode, errors.Alert, []string{"Failed to Delete: ", obj}, []string{err.Error()}, []string{"Meshery Database could be down or not reachable"}, []string{"Restart Meshery instance and make sure database is up and reachable"})
}

func ErrBlankName(err error) error {
	return errors.New(ErrBlankNameCode, errors.Alert, []string{"Error: name field is blank"}, []string{err.Error()}, []string{"Load test name empty or not valid"}, []string{"Provide a name for the test"})
}

func ErrConversion(err error) error {
	return errors.New(ErrConversionCode, errors.Alert, []string{"unable to convert YAML to JSON"}, []string{err.Error()}, []string{"Yaml provided is not valid"}, []string{"Make sure the yaml is valid and has the right parameters"})
}

func ErrLoadTest(err error, obj string) error {
	return errors.New(ErrLoadTestCode, errors.Alert, []string{"load test error: ", obj}, []string{err.Error()}, []string{"Load test endpoint could be not reachable"}, []string{"Make sure load test endpoint is reachable"})
}

func ErrFetchKubernetes(err error) error {
	return errors.New(ErrLoadTestCode, errors.Alert, []string{"unable to ping kubernetes", "unable to scan"}, []string{err.Error()}, []string{"Kubernetes might not be reachable from meshery"}, []string{"Make sure meshery has connectivity to kubernetes"})
}

func ErrPanicRecovery(r interface{}) error {
	return errors.New(ErrPanicRecoveryCode, errors.Alert, []string{"Recovered from panic"}, []string{fmt.Sprint(r)}, []string{"Meshery crashes"}, []string{"Restart Meshery"})
}

func ErrFailToLoadExtensions(err error) error {
	return errors.New(ErrFailToLoadExtensionsCode, errors.Alert, []string{"Failed to Load Extensions from Package"}, []string{err.Error()}, []string{"Plugin is not available in the location", "plugin does not match with meshery version"}, []string{"Make sure the plugin is compatible with Meshery server"})
}

func ErrInvalidLTURL(url string) error {
	return errors.New(ErrInvalidLTURLCode, errors.Alert, []string{"invalid loadtest url: ", url}, []string{}, []string{"URL for load test could be invalid"}, []string{"please refer to:  https://docs.meshery.io/guides/mesheryctl#performance-management"})
}

func ErrVersionCompare(err error) error {
	return errors.New(ErrVersionCompareCode, errors.Alert, []string{"failed to compare latest and current version of Meshery"}, []string{err.Error()}, []string{}, []string{})
}

func ErrSaveSession(err error) error {
	return errors.New(ErrSaveSessionCode, errors.Alert, []string{"unable to save session"}, []string{err.Error()}, []string{"User session could be expired"}, []string{"Re-initiate login"})
}

func ErrCreateDir(err error, obj string) error {
	return errors.New(ErrCreateDirCode, errors.Alert, []string{"Error creating directory ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrInvalidRequestObject(fields ...string) error {
	return errors.New(ErrCreateDirCode, errors.Alert, append([]string{"Error invalid request object:"}, fields...), []string{}, []string{}, []string{})
}

func ErrChangeK8sContext(err error) error {
	return errors.New(ErrCreateDirCode, errors.Alert, []string{"Error changing context"}, []string{err.Error()}, []string{"Context Name might be invalid or not present in the uploaded kubeconfig"}, []string{"Check the context name, if the context name is correct and is present in the kubeconfig then try uploading the kubeconfig again"})
}

func ErrSavingUserPreference(err error) error {
	return errors.New(ErrSavingUserPreferenceCode, errors.Alert, []string{"Error saving user preference."}, []string{err.Error()}, []string{"Invalid data passed", "Unable to connect with provider"}, []string{"Pass valid values for preferences", "Make sure provider supports saving user preferences", "Make sure you're connected with provider", "Make sure extension provides these preferences"})
}
