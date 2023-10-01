package handlers

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrInvalidK8SConfigNilCode          = "1014"
	ErrNilClientCode                    = "1015"
	ErrPrometheusScanCode               = "1016"
	ErrGrafanaScanCode                  = "1017"
	ErrRecordPreferencesCode            = "1018"
	ErrGrafanaConfigCode                = "1019"
	ErrPrometheusConfigCode             = "1020"
	ErrGrafanaQueryCode                 = "1021"
	ErrPrometheusQueryCode              = "1022"
	ErrPrometheusBoardsCode             = "1024"
	ErrStaticBoardsCode                 = "1025"
	ErrRequestBodyCode                  = "1026"
	ErrParseBoolCode                    = "1030"
	ErrStreamEventsCode                 = "1031"
	ErrStreamClientCode                 = "1032"
	ErrPublishSmiResultsCode            = "1033"
	ErrPluginOpenCode                   = "1034"
	ErrPluginLookupCode                 = "1035"
	ErrPluginRunCode                    = "1036"
	ErrParseFormCode                    = "1037"
	ErrQueryGetCode                     = "1038"
	ErrGetResultCode                    = "1039"
	ErrConvertToSpecCode                = "1040"
	ErrFetchSMIResultsCode              = "1041"
	ErrFormFileCode                     = "1042"
	ErrReadConfigCode                   = "1043"
	ErrLoadConfigCode                   = "1044"
	ErrOpenFileCode                     = "1045"
	ErrKubeVersionCode                  = "1046"
	ErrAddAdapterCode                   = "1047"
	ErrRetrieveDataCode                 = "1048"
	ErrValidAdapterCode                 = "1049"
	ErrOperationIDCode                  = "1050"
	ErrMeshClientCode                   = "1051"
	ErrApplyChangeCode                  = "1052"
	ErrRetrieveMeshDataCode             = "1053"
	ErrApplicationFailureCode           = "1054"
	ErrDecodingCode                     = "1055"
	ErrRetrieveUserTokenCode            = "1056"
	ErrFailToSaveCode                   = "1057"
	ErrFailToDeleteCode                 = "1058"
	ErrWriteResponseCode                = "1059"
	ErrTestConfigsCode                  = "1060"
	ErrInvalidGenValueCode              = "1061"
	ErrFailToLoadExtensionsCode         = "1062"
	ErrConversionCode                   = "1063"
	ErrParseDurationCode                = "1064"
	ErrLoadTestCode                     = "1065"
	ErrFetchKubernetesCode              = "1066"
	ErrPanicRecoveryCode                = "1067"
	ErrBlankNameCode                    = "1068"
	ErrInvalidLTURLCode                 = "1069"
	ErrVersionCompareCode               = "1070"
	ErrSaveSessionCode                  = "1071"
	ErrKubeClientCode                   = "1072"
	ErrWorkloadDefinitionCode           = "1073"
	ErrTraitDefinitionCode              = "1074"
	ErrScopeDefinitionCode              = "1075"
	ErrPatternFileCode                  = "1076"
	ErrExecutionPlanCode                = "1077"
	ErrInvalidPatternCode               = "1078"
	ErrCompConfigPairsCode              = "1079"
	ErrCreateDirCode                    = "1080"
	ErrInvalidRequestObjectCode         = "1081"
	ErrChangeK8sContextCode             = "1082"
	ErrSavingUserPreferenceCode         = "1083"
	ErrGetFilterCode                    = "1084"
	ErrSaveFilterCode                   = "1085"
	ErrDecodeFilterCode                 = "1086"
	ErrEncodeFilterCode                 = "1087"
	ErrImportFilterCode                 = "1088"
	ErrFetchFilterCode                  = "1089"
	ErrDeleteFilterCode                 = "1090"
	ErrSavePatternCode                  = "1091"
	ErrSaveApplicationCode              = "1092"
	ErrGetPatternCode                   = "1093"
	ErrDeletePatternCode                = "1094"
	ErrFetchPatternCode                 = "1095"
	ErrImportPatternCode                = "1096"
	ErrEncodePatternCode                = "1097"
	ErrDecodePatternCode                = "1098"
	ErrParsePatternCode                 = "1099"
	ErrConvertPatternCode               = "1100"
	ErrInvalidKubeConfigCode            = "1102"
	ErrInvalidKubeHandlerCode           = "1103"
	ErrInvalidKubeContextCode           = "1104"
	ErrCreatingKubernetesComponentsCode = "1105"
	ErrValidateCode                     = "1106"
	ErrApplicationContentCode           = "1107"
	ErrRemoteApplicationURL             = "1108"
	ErrClonePatternCode                 = "1109"
	ErrCloneFilterCode                  = "1110"
	ErrGenerateComponentsCode           = "1111"
	ErrPublishCatalogPatternCode        = "1112"
	ErrPublishCatalogFilterCode         = "1113"
	ErrGetMeshModelsCode                = "1114"
	ErrGetUserDetailsCode               = "1115"
	ErrResolvingRelationship            = "1116"
	ErrGetLatestVersionCode             = "1117"
	ErrCreateFileCode                   = "1118"
	ErrLoadCertificateCode              = "1119"
	ErrCleanupCertificateCode           = "1120"
	ErrDownlaodWASMFileCode             = "1121"
	ErrFetchProfileCode                 = "1122"
	ErrPerformanceTestCode              = "1123"
	ErrFetchApplicationCode             = "1124"
	ErrDeleteApplicationCode            = "1125"
	ErrGetEventsCode                    = "1126"
	ErrUpdateEventCode                  = "1127"
	ErrDeleteEventCode                  = "1128"
	ErrUnsupportedEventStatusCode       = "1129"
	ErrInvalidRepoURLCode               = "1537"
)

var (
	ErrInvalidK8SConfigNil = errors.New(ErrInvalidK8SConfigNilCode, errors.Alert, []string{"No valid kubernetes config found. Make sure to pass contextIDs in query parameters."}, []string{"Kubernetes config is not initialized with Meshery"}, []string{"Kubernetes config is not accessible to meshery or not valid"}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrNilClient           = errors.New(ErrNilClientCode, errors.Alert, []string{"Kubernetes client not initialized"}, []string{"Kubernetes config is not initialized with Meshery"}, []string{"Kubernetes config is not accessible to meshery or not valid"}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrPrometheusConfig    = errors.New(ErrPrometheusConfigCode, errors.Alert, []string{"Prometheus endpoint not configured"}, []string{"Cannot find valid Prometheus endpoint in user pref"}, []string{"Prometheus endpoint might not be reachable from meshery"}, []string{"Setup your Prometheus Endpoint via the settings dashboard"})
	ErrGrafanaConfig       = errors.New(ErrGrafanaConfigCode, errors.Alert, []string{"Grafana endpoint not configured"}, []string{"Cannot find valid grafana endpoint in user pref"}, []string{"Grafana endpoint might not be reachable from meshery"}, []string{"Setup your Grafana Endpoint via the settings dashboard"})
	ErrStaticBoards        = errors.New(ErrStaticBoardsCode, errors.Alert, []string{"unable to get static board"}, []string{"unable to get static board"}, []string{"No boards could be available in grafana"}, []string{})
	ErrValidAdapter        = errors.New(ErrValidAdapterCode, errors.Alert, []string{"Unable to find valid Adapter URL"}, []string{"unable to find a valid adapter for the given adapter URL"}, []string{"Given adapter URL is not valid"}, []string{"Please provide a valid Adapter URL"})
	ErrAddAdapter          = errors.New(ErrAddAdapterCode, errors.Alert, []string{"meshLocationURL is empty"}, []string{"meshLocationURL is empty to add an adapter"}, []string{"meshLocationURL cannot be empty to add an adapter"}, []string{"please provide the meshLocationURL"})
	ErrMeshClient          = errors.New(ErrMeshClientCode, errors.Alert, []string{"Error creating a mesh client", "Error pinging the mesh adapter"}, []string{"Unable to create a mesh client", "Unable to ping the mesh adapter"}, []string{"Adapter could not be pinged"}, []string{"Unable to connect to the Mesh adapter using the given config, please try again"})
	ErrWriteResponse       = errors.New(ErrWriteResponseCode, errors.Alert, []string{"Error writing response"}, []string{}, []string{}, []string{})
	ErrTestConfigs         = errors.New(ErrTestConfigsCode, errors.Alert, []string{"Error fetching test configs"}, []string{}, []string{}, []string{})
	ErrInvalidGenValue     = errors.New(ErrInvalidGenValueCode, errors.Alert, []string{"Invalid value for gen"}, []string{}, []string{}, []string{"please provide a valid value for gen (load generator)"})
	ErrParseDuration       = errors.New(ErrParseDurationCode, errors.Alert, []string{"error parsing test duration"}, []string{}, []string{"The format of the duration passed could be incorrect"}, []string{"please refer to:  https://docs.meshery.io/guides/mesheryctl#performance-management"})
	ErrPerformanceTest     = errors.New(ErrPerformanceTestCode, errors.Alert, []string{"Load test error"}, []string{}, []string{"Load test endpoint could be not reachable"}, []string{"Make sure load test endpoint is reachable"})
)

func ErrGenerateComponents(err error) error {
	return errors.New(ErrGenerateComponentsCode, errors.Alert, []string{"failed to generate components for the given payload"}, []string{err.Error()}, []string{}, []string{"Make sure the payload is valid"})
}

func ErrValidate(err error) error {
	return errors.New(ErrValidateCode, errors.Alert, []string{"failed to validate the given value against the schema"}, []string{err.Error()}, []string{"unable to validate the value against given schema", "either value or schema might not be a valid cue expression"}, []string{"Make sure that the schema and value provided are valid cue values", "Make sure both schema and value are sent", "Make sure appropriate value types are sent"})
}

func ErrCreatingKubernetesComponents(err error, ctxID string) error {
	return errors.New(ErrCreatingKubernetesComponentsCode, errors.Alert, []string{"failed to register/create kubernetes components for contextID " + ctxID}, []string{err.Error()}, []string{"component generation was canceled due to deletion or reload of K8s context", "Invalid kubeconfig", "Filters passed incorrectly in config", "Could not fetch API resources from Kubernetes server"}, []string{"If there is the log \"Starting to register ...\" for the same contextID after this error means that for some reason the context was reloaded which caused this run to abort. In that case, this error can be ignored.", "Make sure that the configuration filters passed are in accordance with output from /openapi/v2"})
}

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

func ErrPrometheusBoards(err error) error {
	return errors.New(ErrPrometheusBoardsCode, errors.Alert, []string{"unable to get Prometheus boards"}, []string{err.Error()}, []string{"Prometheus endpoint might not be reachable from meshery", "Prometheus endpoint is incorrect"}, []string{"Check if your Prometheus endpoint is correct", "Connect to Prometheus from the settings page in the UI"})
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
	return errors.New(ErrPatternFileCode, errors.Alert, []string{"Failed to Parse design File", err.Error()}, []string{err.Error()}, []string{"Trait Definition is invalid or unable to process"}, []string{"Check Trait Definition"})
}

func ErrInvalidPattern(err error) error {
	return errors.New(ErrInvalidPatternCode, errors.Alert, []string{"Invalid design, execution is infeasible", err.Error()}, []string{err.Error()}, []string{"Trait Definition is invalid or unable to process"}, []string{"Check Trait Definition"})
}

func ErrExecutionPlan(err error) error {
	return errors.New(ErrExecutionPlanCode, errors.Alert, []string{"Failed to Create Execution Plan", err.Error()}, []string{err.Error()}, []string{"Trait Definition is invalid or unable to process"}, []string{"Check Trait Definition"})
}

func ErrCompConfigPairs(err error) error {
	return errors.New(ErrCompConfigPairsCode, errors.Alert, []string{"unable to Create Comp Config.", err.Error()}, []string{err.Error()}, []string{}, []string{})
}

func ErrRequestBody(err error) error {
	return errors.New(ErrRequestBodyCode, errors.Alert, []string{"unable to read the request body"}, []string{err.Error()}, []string{"Request body is empty or faulty"}, []string{"Check if the request is sent with proper values"})
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
	return errors.New(ErrQueryGetCode, errors.Alert, []string{"unable to get: ", obj}, []string{}, []string{"Query parameter is not a part of the request", "Query parameter is malformed."}, []string{"Make sure to pass the query paramater in the request", "Make sure to pass valid query parameter"})
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

func ErrApplicationSourceContent(err error, obj string) error {
	return errors.New(ErrApplicationContentCode, errors.Alert, []string{"failed to ", obj, "the application content"}, []string{err.Error()}, []string{"Remote provider might be not reachable", "Remote provider doesn't support this capability"}, []string{"Ensure you have required permissions or retry after sometime."})
}

func ErrDownloadWASMFile(err error, obj string) error {
	return errors.New(ErrDownlaodWASMFileCode, errors.Alert, []string{"failed to ", obj, "the WASM file"}, []string{err.Error()}, []string{"Ensure that DB is not corrupted", "Ensure Remote Provider is working properly", "Ensure Meshery Server is working properly and connected to remote provider"}, []string{"Try restarting Meshery server"})
}

func ErrDecoding(err error, obj string) error {
	return errors.New(ErrDecodingCode, errors.Alert, []string{"Error decoding the : ", obj}, []string{err.Error()}, []string{"Object is not a valid json object"}, []string{"Make sure if the object passed is a valid json"})
}

func ErrRetrieveUserToken(err error) error {
	return errors.New(ErrRetrieveUserTokenCode, errors.Alert, []string{"Failed to get the user token"}, []string{err.Error()}, []string{"User token could be expired"}, []string{"Re-initiate login"})
}

func ErrFailToSave(err error, obj string) error {
	return errors.New(ErrFailToSaveCode, errors.Alert, []string{"Failed to Save: ", obj}, []string{err.Error()}, []string{"Provider Database could be down or not reachable"}, []string{"Make sure provider is up and reachable"})
}
func ErrFailToDelete(err error, obj string) error {
	return errors.New(ErrFailToDeleteCode, errors.Alert, []string{"Failed to Delete: ", obj}, []string{err.Error()}, []string{"Provider Database could be down or not reachable"}, []string{"Make sure provider is up and reachable"})
}

func ErrBlankName(err error) error {
	return errors.New(ErrBlankNameCode, errors.Alert, []string{"Error: name field is blank"}, []string{err.Error()}, []string{"Load test name empty or not valid"}, []string{"Provide a name for the test"})
}

func ErrConversion(err error) error {
	return errors.New(ErrConversionCode, errors.Alert, []string{"unable to convert YAML to JSON"}, []string{err.Error()}, []string{"Yaml provided is not valid"}, []string{"Make sure the yaml is valid and has the right parameters"})
}

func ErrLoadTest(err error, obj string) error {
	return errors.New(ErrLoadTestCode, errors.Alert, []string{"Load test error: ", obj}, []string{err.Error()}, []string{"Load test endpoint could be not reachable"}, []string{"Make sure load test endpoint is reachable"})
}

func ErrFetchKubernetes(err error) error {
	return errors.New(ErrFetchKubernetesCode, errors.Alert, []string{"unable to ping kubernetes", "unable to scan"}, []string{err.Error()}, []string{"Kubernetes might not be reachable from meshery"}, []string{"Make sure meshery has connectivity to kubernetes"})
}

func ErrPanicRecovery(r interface{}) error {
	return errors.New(ErrPanicRecoveryCode, errors.Alert, []string{"Recovered from panic"}, []string{fmt.Sprint(r)}, []string{"Meshery crashes"}, []string{"Restart Meshery"})
}

func ErrFailToLoadExtensions(err error) error {
	return errors.New(ErrFailToLoadExtensionsCode, errors.Alert, []string{"Failed to Load Extensions from Package"}, []string{err.Error()}, []string{"Plugin is not available in the location", "plugin does not match with meshery version"}, []string{"Make sure the plugin is compatible with Meshery server"})
}

func ErrInvalidLTURL(url string) error {
	return errors.New(ErrInvalidLTURLCode, errors.Alert, []string{"invalid loadtest url: ", url}, []string{}, []string{"URL for load test could be invalid"}, []string{"please refer to: https://docs.meshery.io/tasks/performance-management"})
}

func ErrVersionCompare(err error) error {
	return errors.New(ErrVersionCompareCode, errors.Alert, []string{"failed to compare latest and current version of Meshery"}, []string{err.Error()}, []string{}, []string{})
}

func ErrGetLatestVersion(err error) error {
	return errors.New(ErrGetLatestVersionCode, errors.Alert, []string{"failed to get latest version of Meshery"}, []string{err.Error()}, []string{}, []string{})
}

func ErrSaveSession(err error) error {
	return errors.New(ErrSaveSessionCode, errors.Alert, []string{"unable to save session"}, []string{err.Error()}, []string{"User session could be expired"}, []string{"Re-initiate login"})
}

func ErrCreateDir(err error, obj string) error {
	return errors.New(ErrCreateDirCode, errors.Alert, []string{"Error creating directory ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrInvalidRequestObject(fields ...string) error {
	return errors.New(ErrInvalidRequestObjectCode, errors.Alert, []string{"Error invalid request object:"}, []string{strings.Join(fields, " ")}, []string{""}, []string{""})
}

func ErrChangeK8sContext(err error) error {
	return errors.New(ErrChangeK8sContextCode, errors.Alert, []string{"Error changing context"}, []string{err.Error()}, []string{"Context Name might be invalid or not present in the uploaded kubeconfig"}, []string{"Check the context name, if the context name is correct and is present in the kubeconfig then try uploading the kubeconfig again"})
}

func ErrInvalidKubeConfig(err error, content string) error {
	return errors.New(ErrInvalidKubeConfigCode, errors.Alert, []string{"Invalid Kube Config ", content}, []string{err.Error()}, []string{"Meshery handler failed to find a valid kubernetes config for the deployment"}, []string{"Try uploading a new kubeconfig and also ensure that meshery can reach kubernetes API server"})
}

func ErrInvalidKubeHandler(err error, content string) error {
	return errors.New(ErrInvalidKubeHandlerCode, errors.Alert, []string{"Invalid Kube Handler", content}, []string{err.Error()}, []string{"Meshery handler failed to find a valid kubernetes handler for the deployment"}, []string{"Try uploading a new kubeconfig and also ensure that meshery can reach kubernetes API server"})
}

func ErrInvalidKubeContext(err error, content string) error {
	return errors.New(ErrInvalidKubeContextCode, errors.Alert, []string{"Invalid Kube Context", content}, []string{err.Error()}, []string{"Meshery handler failed to find a valid kubernetes context for the deployment"}, []string{"Try uploading a new kubeconfig and also ensure that meshery can reach kubernetes API server"})
}

func ErrSavingUserPreference(err error) error {
	return errors.New(ErrSavingUserPreferenceCode, errors.Alert, []string{"Error saving user preference."}, []string{err.Error()}, []string{"Invalid data passed", "Unable to connect with provider"}, []string{"Pass valid values for preferences", "Make sure provider supports saving user preferences", "Make sure you're connected with provider", "Make sure extension provides these preferences"})
}

func ErrGetFilter(err error) error {
	return errors.New(ErrGetFilterCode, errors.Alert, []string{"Error failed to get filter"}, []string{err.Error()}, []string{"Cannot get the filter with the given Filter ID"}, []string{"Check if the given Filter ID is correct"})
}

func ErrSaveFilter(err error) error {
	return errors.New(ErrSaveFilterCode, errors.Alert, []string{"Error failed to save filter"}, []string{err.Error()}, []string{"Cannot save the Filter due to wrong path or URL", "Filter is corrupted."}, []string{"Check if the given path or URL of the filter is correct", "Try uplaoding a different filter"})
}

func ErrDecodeFilter(err error) error {
	return errors.New(ErrDecodeFilterCode, errors.Alert, []string{"Error failed to decode filters data into go slice"}, []string{err.Error()}, []string{}, []string{})
}

func ErrEncodeFilter(err error) error {
	return errors.New(ErrEncodeFilterCode, errors.Alert, []string{"Error failed to encode filter"}, []string{err.Error()}, []string{}, []string{})
}

func ErrImportFilter(err error) error {
	return errors.New(ErrImportFilterCode, errors.Alert, []string{"Error failed to import filter"}, []string{err.Error()}, []string{"Cannot save the Filter due to wrong path or URL"}, []string{"Check if the given path or URL of the Filter is correct"})
}

func ErrFetchFilter(err error) error {
	return errors.New(ErrFetchFilterCode, errors.Alert, []string{"Error failed to fetch filter"}, []string{err.Error()}, []string{"Failed to retrieve the list of all the Filters"}, []string{})
}

func ErrDeleteFilter(err error) error {
	return errors.New(ErrDeleteFilterCode, errors.Alert, []string{"Error failed to delete filter"}, []string{err.Error()}, []string{"Failed to delete Filter with the given ID"}, []string{"Check if the Filter ID is correct"})
}

func ErrSavePattern(err error) error {
	return errors.New(ErrSavePatternCode, errors.Alert, []string{"Error failed to save design"}, []string{err.Error()}, []string{"Cannot save the design due to an invalid path or URL"}, []string{"Confirm the correct path / URL to the design"})
}

func ErrSaveApplication(err error) error {
	return errors.New(ErrSaveApplicationCode, errors.Alert, []string{"Error failed to save application"}, []string{err.Error()}, []string{"Cannot save the Application due to wrong path or URL"}, []string{"Check if the given path or URL of the Application is correct"})
}

func ErrFetchApplication(err error) error {
	return errors.New(ErrFetchApplicationCode, errors.Alert, []string{"Error failed to fetch applications"}, []string{err.Error()}, []string{"Remote provider might be not reachable.", "Token might have expired."}, []string{"Refresh your browser"})
}

func ErrDeleteApplication(err error) error {
	return errors.New(ErrDeleteApplicationCode, errors.Alert, []string{"Error failed to delete application"}, []string{err.Error()}, []string{"Application might already have been deleted", "You might not have enough permissions to perform the operation."}, []string{"Check the owner of the application."})
}

func ErrGetPattern(err error) error {
	return errors.New(ErrGetPatternCode, errors.Alert, []string{"Error failed to get design"}, []string{err.Error()}, []string{"Cannot get the design with the given design ID"}, []string{"Check if the given design ID is correct"})
}

func ErrDeletePattern(err error) error {
	return errors.New(ErrDeletePatternCode, errors.Alert, []string{"Error failed to delete design"}, []string{err.Error()}, []string{"Failed to delete design with the given ID"}, []string{"Check if the design ID is correct"})
}

func ErrFetchPattern(err error) error {
	return errors.New(ErrFetchPatternCode, errors.Alert, []string{"Error failed to fetch design"}, []string{err.Error()}, []string{"Failed to retrieve the list of all the designs"}, []string{})
}

func ErrFetchProfile(err error) error {
	return errors.New(ErrFetchProfileCode, errors.Alert, []string{"Error failed to fetch profile"}, []string{err.Error()}, []string{"Invalid profile ID"}, []string{"Check if the profile ID is correct"})
}

func ErrImportPattern(err error) error {
	return errors.New(ErrImportPatternCode, errors.Alert, []string{"Error failed to import design"}, []string{err.Error()}, []string{"Cannot save the design due to wrong path or URL"}, []string{"Check if the given path or URL of the design is correct"})
}

func ErrEncodePattern(err error) error {
	return errors.New(ErrEncodePatternCode, errors.Alert, []string{"Error failed to encode design"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDecodePattern(err error) error {
	return errors.New(ErrDecodePatternCode, errors.Alert, []string{"Error failed to decode design data into go slice"}, []string{err.Error()}, []string{}, []string{})
}

func ErrParsePattern(err error) error {
	return errors.New(ErrParsePatternCode, errors.Alert, []string{"Error failed to parse pattern file"}, []string{err.Error()}, []string{}, []string{})
}

func ErrConvertPattern(err error) error {
	return errors.New(ErrConvertPatternCode, errors.Alert, []string{"Error failed to convert design file to Cytoscape object"}, []string{err.Error()}, []string{}, []string{})
}

func ErrRemoteApplication(err error) error {
	return errors.New(ErrRemoteApplicationURL, errors.Alert, []string{"Error failed to persist remote application"}, []string{err.Error()}, []string{}, []string{})
}

func ErrClonePattern(err error) error {
	return errors.New(ErrClonePatternCode, errors.Alert, []string{"Error failed to clone design"}, []string{err.Error()}, []string{"Failed to clone design with the given ID"}, []string{"Check if the design ID is correct and the design is published"})
}

func ErrCloneFilter(err error) error {
	return errors.New(ErrCloneFilterCode, errors.Alert, []string{"Error failed to clone filter"}, []string{err.Error()}, []string{"Failed to clone Filter with the given ID"}, []string{"Check if the Filter ID is correct and the Filter is published"})
}

func ErrPublishCatalogPattern(err error) error {
	return errors.New(ErrPublishCatalogPatternCode, errors.Alert, []string{"Error failed to publish catalog design"}, []string{err.Error()}, []string{"Failed to publish catalog design"}, []string{"Check if the design ID is correct and you are admin"})
}

func ErrPublishCatalogFilter(err error) error {
	return errors.New(ErrPublishCatalogFilterCode, errors.Alert, []string{"Error failed to publish catalog filter"}, []string{err.Error()}, []string{"Failed to publish catalog filter"}, []string{"Check if the filter ID is correct and you are admin"})
}

func ErrGetMeshModels(err error) error {
	return errors.New(ErrGetMeshModelsCode, errors.Alert, []string{"could not get meshmodel entitities"}, []string{err.Error()}, []string{"Meshmodel entity could not be converted into valid json", "data in the registry was inconsistent"}, []string{"make sure correct and consistent data is present inside the registry", "drop the Meshmodel tables and restart Meshery server"})
}

func ErrGetUserDetails(err error) error {
	return errors.New(ErrGetUserDetailsCode, errors.Alert, []string{"could not get user details"}, []string{err.Error()}, []string{"User details could not be fetched from provider", "Your provider may not be reachable", "No user exists for the provided token"}, []string{"Make sure provider is reachable", "Make sure you are logged in", "Make sure you are using a valid token"})
}

func ErrResolvingRegoRelationship(err error) error {
	return errors.New(ErrResolvingRelationship, errors.Alert, []string{"could not resolve rego relationship"}, []string{err.Error()}, []string{"The rego evaluation engine failed to resolve policies", "Design-File/Application-File is in incorrect format", "The policy query is invalid", "The evaluation engine response is unexpected for the code written"}, []string{"Make sure the design-file/application-file is a valid yaml", "Make sure you're proving correct rego query", "Make sure the server is evaluating the query correctly, add some logs"})
}

func ErrCreateFile(err error, obj string) error {
	return errors.New(ErrCreateFileCode, errors.Alert, []string{"Could not create file", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrLoadCertificate(err error) error {
	return errors.New(ErrLoadCertificateCode, errors.Alert, []string{"Could not load certificates associated with performance profile"}, []string{err.Error()}, []string{"Remote provider might be not reachable"}, []string{"try running performance profile test without using certificates, update the profile without certificates"})
}

func ErrCleanupCertificate(err error, obj string) error {
	return errors.New(ErrCleanupCertificateCode, errors.Alert, []string{"Could not delete certificates from ", obj}, []string{err.Error()}, []string{"might be due to insufficient permissions", "file was deleted manually"}, []string{"please delete the file if present, path: ", obj})
}

func ErrGetEvents(err error) error {
	return errors.New(ErrGetEventsCode, errors.Alert, []string{"Could not retrieve events"}, []string{err.Error()}, []string{"Request contains unknown query variables.", "Database is not reachable or corrupt."}, []string{"Check the request URL and try again."})
}

func ErrUpdateEvent(err error, id string) error {
	return errors.New(ErrUpdateEventCode, errors.Alert, []string{fmt.Sprintf("Could not update event status for %s", id)}, []string{err.Error()}, []string{"Provided event status not supported", "Event has been deleted or does not exist", "Database is corrupt."}, []string{"Verify event filter settings", "Reset database."})
}

func ErrDeleteEvent(err error, id string) error {
	return errors.New(ErrDeleteEventCode, errors.Alert, []string{fmt.Sprintf("Could not delete event %s", id)}, []string{err.Error()}, []string{"Event might have been deleted and doesn't exist", "Database is corrupt."}, []string{"Verify event filter settings", "Reset database."})
}

func ErrUnsupportedEventStatus(err error, status string) error {
	return errors.New(ErrUnsupportedEventStatusCode, errors.Alert, []string{fmt.Sprintf("Event status %s not supported.", status)}, []string{err.Error()}, []string{"Unsupported event status for the current Meshery Server."}, []string{"Try upgrading Meshery to latest version.", "Use one of the supported event statuses."})
}

func ErrInvalidRepoURL(err error) error {
	return errors.New(ErrInvalidRepoURLCode, errors.Alert, []string{"cannot get url from the request body"}, []string{err.Error()}, []string{"request body is empty/do not contain repository URL"}, []string{"provide a valid request body"})
}
