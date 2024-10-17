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
	ErrInvalidK8SConfigNilCode             = "meshery-server-1014"
	ErrNilClientCode                       = "meshery-server-1015"
	ErrRecordPreferencesCode               = "meshery-server-1016"
	ErrGrafanaConfigCode                   = "meshery-server-1017"
	ErrPrometheusConfigCode                = "meshery-server-1018"
	ErrGrafanaQueryCode                    = "meshery-server-1019"
	ErrPrometheusQueryCode                 = "meshery-server-1020"
	ErrPrometheusBoardsCode                = "meshery-server-1021"
	ErrStaticBoardsCode                    = "meshery-server-1022"
	ErrRequestBodyCode                     = "meshery-server-1023"
	ErrParseBoolCode                       = "meshery-server-1024"
	ErrStreamEventsCode                    = "meshery-server-1025"
	ErrStreamClientCode                    = "meshery-server-1026"
	ErrPublishSmiResultsCode               = "meshery-server-1027"
	ErrPluginOpenCode                      = "meshery-server-1028"
	ErrPluginLookupCode                    = "meshery-server-1029"
	ErrPluginRunCode                       = "meshery-server-1030"
	ErrParseFormCode                       = "meshery-server-1031"
	ErrQueryGetCode                        = "meshery-server-1032"
	ErrGetResultCode                       = "meshery-server-1033"
	ErrConvertToSpecCode                   = "meshery-server-1034"
	ErrFetchSMIResultsCode                 = "meshery-server-1035"
	ErrFormFileCode                        = "meshery-server-1036"
	ErrReadConfigCode                      = "meshery-server-1037"
	ErrLoadConfigCode                      = "meshery-server-1038"
	ErrOpenFileCode                        = "meshery-server-1039"
	ErrKubeVersionCode                     = "meshery-server-1040"
	ErrAddAdapterCode                      = "meshery-server-1041"
	ErrRetrieveDataCode                    = "meshery-server-1042"
	ErrValidAdapterCode                    = "meshery-server-1043"
	ErrOperationIDCode                     = "meshery-server-1044"
	ErrMeshClientCode                      = "meshery-server-1045"
	ErrApplyChangeCode                     = "meshery-server-1046"
	ErrRetrieveMeshDataCode                = "meshery-server-1047"
	ErrApplicationFailureCode              = "meshery-server-1048"
	ErrDecodingCode                        = "meshery-server-1049"
	ErrRetrieveUserTokenCode               = "meshery-server-1050"
	ErrFailToSaveCode                      = "meshery-server-1051"
	ErrFailToDeleteCode                    = "meshery-server-1052"
	ErrWriteResponseCode                   = "meshery-server-1053"
	ErrTestConfigsCode                     = "meshery-server-1054"
	ErrInvalidGenValueCode                 = "meshery-server-1055"
	ErrFailToLoadExtensionsCode            = "meshery-server-1056"
	ErrConversionCode                      = "meshery-server-1057"
	ErrParseDurationCode                   = "meshery-server-1058"
	ErrLoadTestCode                        = "meshery-server-1059"
	ErrFetchKubernetesCode                 = "meshery-server-1060"
	ErrPanicRecoveryCode                   = "meshery-server-1061"
	ErrBlankNameCode                       = "meshery-server-1062"
	ErrInvalidLTURLCode                    = "meshery-server-1063"
	ErrVersionCompareCode                  = "meshery-server-1064"
	ErrSaveSessionCode                     = "meshery-server-1065"
	ErrKubeClientCode                      = "meshery-server-1066"
	ErrWorkloadDefinitionCode              = "meshery-server-1067"
	ErrTraitDefinitionCode                 = "meshery-server-1068"
	ErrScopeDefinitionCode                 = "meshery-server-1069"
	ErrPatternFileCode                     = "meshery-server-1070"
	ErrExecutionPlanCode                   = "meshery-server-1071"
	ErrInvalidPatternCode                  = "meshery-server-1072"
	ErrPatternDeployCode                   = "meshery-server-1073"
	ErrCreateDirCode                       = "meshery-server-1074"
	ErrInvalidRequestObjectCode            = "meshery-server-1075"
	ErrChangeK8sContextCode                = "meshery-server-1076"
	ErrSavingUserPreferenceCode            = "meshery-server-1077"
	ErrGetFilterCode                       = "meshery-server-1078"
	ErrSaveFilterCode                      = "meshery-server-1079"
	ErrDecodeFilterCode                    = "meshery-server-1080"
	ErrEncodeFilterCode                    = "meshery-server-1081"
	ErrImportFilterCode                    = "meshery-server-1082"
	ErrFetchFilterCode                     = "meshery-server-1083"
	ErrDeleteFilterCode                    = "meshery-server-1084"
	ErrSavePatternCode                     = "meshery-server-1085"
	ErrSaveApplicationCode                 = "meshery-server-1086"
	ErrGetPatternCode                      = "meshery-server-1087"
	ErrDeletePatternCode                   = "meshery-server-1088"
	ErrFetchPatternCode                    = "meshery-server-1089"
	ErrImportPatternCode                   = "meshery-server-1090"
	ErrEncodePatternCode                   = "meshery-server-1091"
	ErrDecodePatternCode                   = "meshery-server-1092"
	ErrParsePatternCode                    = "meshery-server-1093"
	ErrConvertPatternCode                  = "meshery-server-1094"
	ErrInvalidKubeConfigCode               = "meshery-server-1095"
	ErrInvalidKubeHandlerCode              = "meshery-server-1096"
	ErrInvalidKubeContextCode              = "meshery-server-1097"
	ErrValidateCode                        = "meshery-server-1098"
	ErrApplicationContentCode              = "meshery-server-1099"
	ErrRemoteApplicationURLCode            = "meshery-server-1100"
	ErrClonePatternCode                    = "meshery-server-1101"
	ErrCloneFilterCode                     = "meshery-server-1102"
	ErrGenerateComponentsCode              = "meshery-server-1103"
	ErrPublishCatalogPatternCode           = "meshery-server-1104"
	ErrPublishCatalogFilterCode            = "meshery-server-1105"
	ErrGetMeshModelsCode                   = "meshery-server-1106"
	ErrGetUserDetailsCode                  = "meshery-server-1107"
	ErrResolvingRelationshipCode           = "meshery-server-1108"
	ErrGetLatestVersionCode                = "meshery-server-1109"
	ErrCreateFileCode                      = "meshery-server-1110"
	ErrLoadCertificateCode                 = "meshery-server-1111"
	ErrCleanupCertificateCode              = "meshery-server-1112"
	ErrDownlaodWASMFileCode                = "meshery-server-1113"
	ErrFetchProfileCode                    = "meshery-server-1114"
	ErrPerformanceTestCode                 = "meshery-server-1115"
	ErrFetchApplicationCode                = "meshery-server-1116"
	ErrDeleteApplicationCode               = "meshery-server-1117"
	ErrGetEventsCode                       = "meshery-server-1118"
	ErrUpdateEventCode                     = "meshery-server-1119"
	ErrDeleteEventCode                     = "meshery-server-1120"
	ErrUnsupportedEventStatusCode          = "meshery-server-1121"
	ErrBulkUpdateEventCode                 = "meshery-server-1122"
	ErrBulkDeleteEventCode                 = "meshery-server-1123"
	ErrFetchMeshSyncResourcesCode          = "meshery-server-1124"
	ErrDesignSourceContentCode             = "meshery-server-1125"
	ErrGetConnectionsCode                  = "meshery-server-1126"
	ErrWritingIntoFileCode                 = "meshery-server-1127"
	ErrBuildOCIImgCode                     = "meshery-server-1128"
	ErrSaveOCIArtifactCode                 = "meshery-server-1129"
	ErrIOReaderCode                        = "meshery-server-1130"
	ErrUnCompressOCIArtifactCode           = "meshery-server-1131"
	ErrWaklingLocalDirectoryCode           = "meshery-server-1132"
	ErrConvertingK8sManifestToDesignCode   = "meshery-server-1133"
	ErrConvertingDockerComposeToDesignCode = "meshery-server-1134"
	ErrConvertingHelmChartToDesignCode     = "meshery-server-1136"
	ErrInvalidUUIDCode                     = "meshery-server-1137"
	ErrPersistEventToRemoteProviderCode    = "meshery-server-1320"
	ErrEventStreamingNotSupportedCode      = "meshery-server-1324"
	ErrGenerateClusterContextCode          = "meshery-server-1325"
	ErrNilClusterContextCode               = "meshery-server-1326"
	ErrFailToSaveContextCode               = "meshery-server-1327"
	ErrParsingCallBackUrlCode              = "meshery-server-1328"
	ErrReadSessionPersistorCode            = "meshery-server-1329"
	ErrFailToGetK8SContextCode             = "meshery-server-1330"
	ErrFailToLoadK8sContextCode            = "meshery-server-1331"
	ErrNoTarInsideOCiCode                  = "meshery-server-1365"
	ErrEmptyOCIImageCode                   = "meshery-server-1360"
	ErrGetComponentDefinitionCode          = "meshery-server-1362"
	ErrGetCapabilitiesCode                 = "meshery-server-1363"
	ErrExportPatternInFormatCode           = "meshery-server-1364"
)

var (
	ErrInvalidK8SConfigNil        = errors.New(ErrInvalidK8SConfigNilCode, errors.Alert, []string{"No valid Kubernetes config found. Make sure to pass contextIDs in query parameters."}, []string{"Kubernetes config is not initialized with Meshery"}, []string{"Kubernetes config is not accessible to Meshery or not valid"}, []string{"Upload your Kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrNilClient                  = errors.New(ErrNilClientCode, errors.Alert, []string{"Kubernetes client not initialized"}, []string{"Kubernetes config is not initialized with Meshery"}, []string{"Kubernetes config is not accessible to Meshery or not valid"}, []string{"Upload your Kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrPrometheusConfig           = errors.New(ErrPrometheusConfigCode, errors.Alert, []string{"Prometheus endpoint not configured"}, []string{"Cannot find valid Prometheus endpoint in user pref"}, []string{"Prometheus endpoint might not be reachable from Meshery"}, []string{"Setup your Prometheus Endpoint via the settings dashboard"})
	ErrGrafanaConfig              = errors.New(ErrGrafanaConfigCode, errors.Alert, []string{"Grafana endpoint not configured"}, []string{"Cannot find valid grafana endpoint in user pref"}, []string{"Grafana endpoint might not be reachable from Meshery"}, []string{"Setup your Grafana Endpoint via the settings dashboard"})
	ErrStaticBoards               = errors.New(ErrStaticBoardsCode, errors.Alert, []string{"unable to get static board"}, []string{"unable to get static board"}, []string{"No boards could be available in grafana"}, []string{})
	ErrValidAdapter               = errors.New(ErrValidAdapterCode, errors.Alert, []string{"Unable to find valid Adapter URL"}, []string{"unable to find a valid adapter for the given adapter URL"}, []string{"Given adapter URL is not valid"}, []string{"Please provide a valid Adapter URL"})
	ErrAddAdapter                 = errors.New(ErrAddAdapterCode, errors.Alert, []string{"meshLocationURL is empty"}, []string{"meshLocationURL is empty to add an adapter"}, []string{"meshLocationURL cannot be empty to add an adapter"}, []string{"please provide the meshLocationURL"})
	ErrMeshClient                 = errors.New(ErrMeshClientCode, errors.Alert, []string{"Error creating a mesh client", "Error pinging the mesh adapter"}, []string{"Unable to create a mesh client", "Unable to ping the mesh adapter"}, []string{"Adapter could not be pinged"}, []string{"Unable to connect to the Mesh adapter using the given config, please try again"})
	ErrTestConfigs                = errors.New(ErrTestConfigsCode, errors.Alert, []string{"Error fetching test configs"}, []string{}, []string{}, []string{})
	ErrInvalidGenValue            = errors.New(ErrInvalidGenValueCode, errors.Alert, []string{"Invalid value for gen"}, []string{}, []string{}, []string{"please provide a valid value for gen (load generator)"})
	ErrParseDuration              = errors.New(ErrParseDurationCode, errors.Alert, []string{"error parsing test duration"}, []string{}, []string{"The format of the duration passed could be incorrect"}, []string{"please refer to:  https://docs.meshery.io/guides/mesheryctl#performance-management"})
	ErrPerformanceTest            = errors.New(ErrPerformanceTestCode, errors.Alert, []string{"Load test error"}, []string{}, []string{"Load test endpoint could be not reachable"}, []string{"Make sure load test endpoint is reachable"})
	ErrEventStreamingNotSupported = errors.New(ErrEventStreamingNotSupportedCode, errors.Alert, []string{"Live events stream not supported."}, []string{"Our server cannot provide live events streaming at the moment. This might be due to a technical issue with our system."}, []string{}, []string{})
	ErrReadSessionPersistor       = errors.New(ErrReadSessionPersistorCode, errors.Alert, []string{"Unable to read session from the session persister, starting with a new one"}, []string{}, []string{}, []string{})
	ErrFailToGetK8SContext        = errors.New(ErrFailToGetK8SContextCode, errors.Alert, []string{"Failed to get Kubernetes context"}, []string{}, []string{}, []string{})
)

func ErrGenerateClusterContext(err error) error {
	return errors.New(ErrGenerateClusterContextCode, errors.Alert, []string{"Failed to generate in cluster context."}, []string{err.Error()}, []string{}, []string{})
}
func ErrNilClusterContext(err error) error {
	return errors.New(ErrNilClusterContextCode, errors.Alert, []string{"Nil context generated from in cluster config"}, []string{err.Error()}, []string{}, []string{})
}
func ErrWriteResponse(err error) error {
	return errors.New(ErrWriteResponseCode, errors.Alert, []string{"Error writing response"}, []string{err.Error()}, []string{}, []string{})
}
func ErrFailToSaveContext(err error) error {
	return errors.New(ErrFailToSaveContextCode, errors.Alert, []string{"Failed to save the context"}, []string{err.Error()}, []string{}, []string{})
}
func ErrGenerateComponents(err error) error {
	return errors.New(ErrGenerateComponentsCode, errors.Alert, []string{"failed to generate components for the given payload"}, []string{err.Error()}, []string{}, []string{"Make sure the payload is valid"})
}

func ErrValidate(err error) error {
	return errors.New(ErrValidateCode, errors.Alert, []string{"failed to validate the given value against the schema"}, []string{err.Error()}, []string{"unable to validate the value against given schema", "either value or schema might not be a valid cue expression"}, []string{"Make sure that the schema and value provided are valid cue values", "Make sure both schema and value are sent", "Make sure appropriate value types are sent"})
}
func ErrParsingCallBackUrl(err error) error {
	return errors.New(ErrParsingCallBackUrlCode, errors.Alert, []string{"Failed to parse the callback URL"}, []string{err.Error()}, []string{"callback URL might be empty"}, []string{"Make sure the callback URL is not empty"})
}
func ErrFailToLoadK8sContext(err error) error {
	return errors.New(ErrFailToLoadK8sContextCode, errors.Alert, []string{"Failed to load Kubernetes context"}, []string{err.Error()}, []string{}, []string{})
}
func ErrPrometheusQuery(err error) error {
	return errors.New(ErrPrometheusQueryCode, errors.Alert, []string{"Unable to query prometheus"}, []string{err.Error()}, []string{"Prometheus query did not get executed from Meshery", "Prometheus query is invalid"}, []string{"Check if your Prometheus query is correct", "Connect to Prometheus and Grafana from the settings page in the UI"})
}

func ErrGrafanaQuery(err error) error {
	return errors.New(ErrGrafanaQueryCode, errors.Alert, []string{"Unable to query Grafana"}, []string{err.Error()}, []string{"Grafana query did not get executed from Meshery", "Grafana query is invalid"}, []string{"Check if your Grafana query is correct", "Connect to Grafana from the settings page in the UI"})
}

func ErrPrometheusBoards(err error) error {
	return errors.New(ErrPrometheusBoardsCode, errors.Alert, []string{"unable to get Prometheus boards"}, []string{err.Error()}, []string{"Prometheus endpoint might not be reachable from Meshery", "Prometheus endpoint is incorrect"}, []string{"Check if your Prometheus endpoint is correct", "Connect to Prometheus from the settings page in the UI"})
}

func ErrRecordPreferences(err error) error {
	return errors.New(ErrRecordPreferencesCode, errors.Alert, []string{"unable to save user config data"}, []string{err.Error()}, []string{"User token might be invalid", "db might be corrupted"}, []string{"Relogin to Meshery"})
}

func ErrKubeClient(err error) error {
	return errors.New(ErrKubeClientCode, errors.Alert, []string{"Failed to Create Kube Client", err.Error()}, []string{err.Error()}, []string{"Check Kubernetes"}, []string{"Check your kubeconfig if valid", "Ensure Meshery is able to reach the Kubernetes cluster"})
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

func ErrPatternDeploy(err error, patternName string) error {
	return errors.New(ErrPatternDeployCode, errors.Alert, []string{"Unable to deploy the selected design \"%s\"", patternName}, []string{err.Error()}, []string{"Connection Error: There was an error connecting to the selected target platform (i.e. Kubernetes cluster(s)).", "This connection might not be assigned to the selected environment."}, []string{"Verify that the Kubernetes connection status is 'Connected' or try uploading a new kubeconfig.", "Assign the current Kubernetes connection to the selected environment."})
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
	return errors.New(ErrPublishSmiResultsCode, errors.Alert, []string{"Error publishing SMI results"}, []string{err.Error()}, []string{"Meshery Cloud is not functional or reachable"}, []string{"Make sure Meshery cloud is up and reachable"})
}

func ErrPluginOpen(err error) error {
	return errors.New(ErrPluginOpenCode, errors.Alert, []string{"Error opening the plugin"}, []string{err.Error()}, []string{"Plugin is not available in the location", "plugin does not match with Meshery version"}, []string{"Make sure the plugin is compatible with Meshery server"})
}

func ErrPluginLookup(err error) error {
	return errors.New(ErrPluginLookupCode, errors.Alert, []string{"Error performing a plugin lookup"}, []string{err.Error()}, []string{"Plugin is not available in the location"}, []string{"Make sure the plugin is compatible with Meshery server"})
}

func ErrPluginRun(err error) error {
	return errors.New(ErrPluginRunCode, errors.Alert, []string{"Error running Meshery plugin"}, []string{err.Error()}, []string{"plugin does not match with Meshery version"}, []string{"Make sure the plugin is compatible with Meshery server"})
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
	return errors.New(ErrLoadConfigCode, errors.Alert, []string{"unable to load Kubernetes config"}, []string{err.Error()}, []string{"The kubeconfig file is empty or not valid"}, []string{"Make sure to upload the correct kubeconfig file"})
}

func ErrOpenFile(file string) error {
	return errors.New(ErrOpenFileCode, errors.Alert, []string{"unable to open file: ", file}, []string{}, []string{"The file does not exist in the location"}, []string{"Make sure to upload the correct file"})
}

func ErrKubeVersion(err error) error {
	return errors.New(ErrKubeVersionCode, errors.Alert, []string{"unable to get Kubernetes version"}, []string{err.Error()}, []string{"Kubernetes might not be reachable from Meshery"}, []string{"Make sure Meshery has connectivity to Kubernetes"})
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
	return errors.New(ErrRetrieveMeshDataCode, errors.Alert, []string{"Error getting operations for the mesh", "Error getting component name"}, []string{err.Error()}, []string{"unable to retrieve the requested data"}, []string{"Make sure adapter is reachable and running"})
}

func ErrApplicationFailure(err error, obj string) error {
	return errors.New(ErrApplicationFailureCode, errors.Alert, []string{"failed to ", obj, "the application"}, []string{err.Error()}, []string{"uploaded application source content might be converted", "incorrect source type selected"}, []string{"Select the correct source type", "Make sure the uploaded application source content is valid"})
}

func ErrApplicationSourceContent(err error, obj string) error {
	return errors.New(ErrApplicationContentCode, errors.Alert, []string{"failed to ", obj, "the application content"}, []string{err.Error()}, []string{"Remote provider might be not reachable", "Remote provider doesn't support this capability"}, []string{"Ensure you have required permissions or retry after sometime."})
}

func ErrDesignSourceContent(err error, obj string) error {
	return errors.New(ErrDesignSourceContentCode, errors.Alert, []string{"failed to ", obj, "the design content"}, []string{err.Error()}, []string{"Remote provider might be not reachable", "Remote provider doesn't support this capability"}, []string{"Ensure you have required permissions or retry after sometime."})
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
	return errors.New(ErrFetchKubernetesCode, errors.Alert, []string{"unable to ping Kubernetes", "unable to scan"}, []string{err.Error()}, []string{"Kubernetes might not be reachable from Meshery"}, []string{"Make sure Meshery has connectivity to Kubernetes"})
}

func ErrPanicRecovery(r interface{}) error {
	return errors.New(ErrPanicRecoveryCode, errors.Alert, []string{"Recovered from panic"}, []string{fmt.Sprint(r)}, []string{"Meshery crashes"}, []string{"Restart Meshery"})
}

func ErrFailToLoadExtensions(err error) error {
	return errors.New(ErrFailToLoadExtensionsCode, errors.Alert, []string{"Failed to Load Extensions from Package"}, []string{err.Error()}, []string{"Plugin is not available in the location", "plugin does not match with Meshery version"}, []string{"Make sure the plugin is compatible with Meshery server"})
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
	return errors.New(ErrCreateDirCode, errors.Alert, []string{"Error creating directory ", obj}, []string{err.Error()}, []string{"Insufficient permission", "Insufficient storage"}, []string{"check if sufficient permissions are available to create dir", "check if sufficient storage is available to create dir"})
}

func ErrInvalidRequestObject(fields ...string) error {
	return errors.New(ErrInvalidRequestObjectCode, errors.Alert, []string{"Error invalid request object:"}, []string{strings.Join(fields, " ")}, []string{""}, []string{""})
}

func ErrChangeK8sContext(err error) error {
	return errors.New(ErrChangeK8sContextCode, errors.Alert, []string{"Error changing context"}, []string{err.Error()}, []string{"Context Name might be invalid or not present in the uploaded kubeconfig"}, []string{"Check the context name, if the context name is correct and is present in the kubeconfig then try uploading the kubeconfig again"})
}

func ErrInvalidKubeConfig(err error, content string) error {
	return errors.New(ErrInvalidKubeConfigCode, errors.Alert, []string{"Invalid Kube Config ", content}, []string{err.Error()}, []string{"Meshery handler failed to find a valid Kubernetes config for the deployment"}, []string{"Try uploading a new kubeconfig and also ensure that Meshery can reach Kubernetes API server"})
}

func ErrInvalidKubeHandler(err error, userId string) error {
	return errors.New(ErrInvalidKubeHandlerCode, errors.Alert, []string{"Kubernetes cluster is unavailable for ", userId}, []string{err.Error()}, []string{"There might be a network disruption or the Meshery server does not have valid credentials."}, []string{"Try uploading a new kubeconfig.", "Check the network connection and Kubernetes cluster status.", "Verify that the Meshery server has valid and updated credentials to access the Kubernetes cluster."})
}

func ErrInvalidKubeContext(err error, content string) error {
	return errors.New(ErrInvalidKubeContextCode, errors.Alert, []string{"Invalid Kube Context", content}, []string{err.Error()}, []string{"Meshery handler failed to find a valid Kubernetes context for the deployment"}, []string{"Try uploading a new kubeconfig and also ensure that Meshery can reach Kubernetes API server"})
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
	return errors.New(ErrImportPatternCode, errors.Alert, []string{"Error failed to import design"}, []string{err.Error()}, []string{"Cannot save the design due to wrong path or URL"}, []string{"Check if the provided path or URL of the design is correct. If you are providing a URL, it should be a direct URL to a downloadable file. For example, if the file is stored on GitHub, the URL should be 'https://raw.githubusercontent.com/path-to-file'."})
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
	return errors.New(ErrRemoteApplicationURLCode, errors.Alert, []string{"Error failed to persist remote application"}, []string{err.Error()}, []string{}, []string{})
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

func ErrGetComponentDefinition(err error) error {
	return errors.New(ErrGetComponentDefinitionCode, errors.Alert, []string{"component definition not found"}, []string{err.Error()}, []string{"Component might not have been registered", "The component might not be supported by default, in the version of Meshery you are currently using."}, []string{"Ensure component definition is valid JSON/YAML.", "Import the model and try again."})
}

func ErrGetUserDetails(err error) error {
	return errors.New(ErrGetUserDetailsCode, errors.Alert, []string{"could not get user details"}, []string{err.Error()}, []string{"User details could not be fetched from provider", "Your provider may not be reachable", "No user exists for the provided token"}, []string{"Make sure provider is reachable", "Make sure you are logged in", "Make sure you are using a valid token"})
}

func ErrResolvingRegoRelationship(err error) error {
	return errors.New(ErrResolvingRelationshipCode, errors.Alert, []string{"could not resolve rego relationship"}, []string{err.Error()}, []string{"The rego evaluation engine failed to resolve policies", "Design-File/Application-File is in incorrect format", "The policy query is invalid", "The evaluation engine response is unexpected for the code written"}, []string{"Make sure the design-file/application-file is a valid yaml", "Make sure you're proving correct rego query", "Make sure the server is evaluating the query correctly, add some logs"})
}

func ErrCreateFile(err error, obj string) error {
	return errors.New(ErrCreateFileCode, errors.Alert, []string{"Could not create file", obj}, []string{err.Error()}, []string{"Insufficient permission", "Insufficient storage"}, []string{"check if sufficient permissions are available to create file", "check if sufficient storage is available to create file"})
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

func ErrBulkUpdateEvent(err error) error {
	return errors.New(ErrBulkUpdateEventCode, errors.Alert, []string{"Could not update status for one or more events."}, []string{err.Error()}, []string{"Event has been deleted or does not exist.", "The requested event status is invalid.", "Meshery Database is corrupt."}, []string{"Verify that the event still exists.", "Verify that the requested event status is supported.", "Visit Settings and reset the Meshery database."})
}

func ErrDeleteEvent(err error, id string) error {
	return errors.New(ErrDeleteEventCode, errors.Alert, []string{fmt.Sprintf("Could not delete event %s", id)}, []string{err.Error()}, []string{"Event might have been deleted and doesn't exist", "Database is corrupt."}, []string{"Verify event filter settings", "Reset database."})
}

func ErrBulkDeleteEvent(err error) error {
	return errors.New(ErrBulkDeleteEventCode, errors.Alert, []string{"Could not delete one or more events."}, []string{err.Error()}, []string{"Event has been deleted or does not exist.", "Meshery Database is corrupt."}, []string{"Confirm that the status you are using is valid and a supported event status. Refer to Meshery Docs for a list of event statuses. Check for availability of a new version of Meshery Server. Try upgrading to the latest version.", "Visit Settings and reset the Meshery database."})
}

func ErrUnsupportedEventStatus(err error, status string) error {
	return errors.New(ErrUnsupportedEventStatusCode, errors.Alert, []string{fmt.Sprintf("Event status '%s' is not a supported status.", status)}, []string{err.Error()}, []string{"Unsupported event status for your current version of Meshery Server."}, []string{"Confirm that the status you are using is valid and a supported event status. Refer to Meshery Docs for a list of event statuses.", "Check for availability of a new version of Meshery Server. Try upgrading to the latest version."})
}

// ErrFetchMeshSyncResources
func ErrFetchMeshSyncResources(err error) error {
	return errors.New(ErrFetchMeshSyncResourcesCode, errors.Alert, []string{"Error fetching MeshSync resources", "DB might be corrupted"}, []string{err.Error()}, []string{"MeshSync might not be reachable from Meshery"}, []string{"Make sure Meshery has connectivity to MeshSync", "Try restarting Meshery server"})
}

func ErrGetConnections(err error) error {
	return errors.New(ErrGetConnectionsCode, errors.Alert, []string{"Failed to retrieve connections"}, []string{err.Error()}, []string{"Unable to retrieve the connections"}, []string{"Check if the cluster is connected and healthy, you can check it from k8s switcher in header"})
}

func ErrWritingIntoFile(err error, obj string) error {
	return errors.New(ErrWritingIntoFileCode, errors.Alert, []string{fmt.Sprintf("failed to write into file %s", obj)}, []string{err.Error()}, []string{"Insufficient permissions to write into file", "file might be corrupted"}, []string{"check if sufficient permissions are givent to the file", "check if the file is corrupted"})
}

func ErrBuildOCIImg(err error) error {
	return errors.New(ErrBuildOCIImgCode, errors.Alert, []string{"Failed to build OCI image"}, []string{err.Error()}, []string{"unable to read source directory", "source directory is corrupted"}, []string{"check if the source directory is valid and has sufficient permissions", "check if the source directory is not corrupted"})
}

func ErrSaveOCIArtifact(err error) error {
	return errors.New(ErrSaveOCIArtifactCode, errors.Alert, []string{"Failed to persist OCI artifact"}, []string{err.Error()}, []string{"unable to read source directory", "source directory is corrupted", "unable to persist in requested location", "OCI img may be corrupted"}, []string{"check if the source directory is valid and has sufficient permissions", "check if the source directory is not corrupted", "check if sufficient permissions are available to write in requested location", "check if the OCI img is not corrupted"})
}

func ErrIOReader(err error) error {
	return errors.New(ErrIOReaderCode, errors.Alert, []string{"Failed to read from io.Reader"}, []string{err.Error()}, []string{"unable to read from io.Reader"}, []string{"check if the io.Reader is valid"})
}

func ErrUnCompressOCIArtifact(err error) error {
	return errors.New(ErrUnCompressOCIArtifactCode, errors.Alert, []string{"Failed to uncompress OCI artifact"}, []string{err.Error()}, []string{"unable to uncompress OCI artifact", "OCI artifact may be corrupted"}, []string{"check if the OCI artifact is valid and not corrupted"})
}

func ErrWaklingLocalDirectory(err error) error {
	return errors.New(ErrWaklingLocalDirectoryCode, errors.Alert, []string{"Failed to walk local directory"}, []string{err.Error()}, []string{"unable to walk local directory", "local directory may be corrupted"}, []string{"check if the local directory is valid and not corrupted"})
}

func ErrConvertingK8sManifestToDesign(err error) error {
	return errors.New(ErrConvertingK8sManifestToDesignCode, errors.Alert, []string{"Failed to convert k8s manifest to design"}, []string{err.Error()}, []string{"unable to convert k8s manifest to design", "k8s manifest may be corrupted", "incorrect source type selected"}, []string{"check if the k8s manifest is valid and not corrupted", "check if the source type selected is Kubernetes Manifest"})
}

func ErrConvertingDockerComposeToDesign(err error) error {
	return errors.New(ErrConvertingDockerComposeToDesignCode, errors.Alert, []string{"Failed to convert docker compose to design"}, []string{err.Error()}, []string{"unable to convert docker compose to design", "docker compose may be corrupted", "incorrect source type selected"}, []string{"check if the docker compose is valid and not corrupted", "check if the source type selected is Docker Compose"})
}

func ErrConvertingHelmChartToDesign(err error) error {
	return errors.New(ErrConvertingHelmChartToDesignCode, errors.Alert, []string{"Failed to convert helm chart to design"}, []string{err.Error()}, []string{"unable to convert helm chart to design", "helm chart may be corrupted", "incorrect source type selected"}, []string{"check if the helm chart is valid and not corrupted", "check if the source type selected is Helm Chart"})
}

func ErrInvalidUUID(err error) error {
	return errors.New(ErrInvalidUUIDCode, errors.Alert, []string{"invalid or empty uuid"}, []string{err.Error()}, []string{"provided id is not a valid uuid"}, []string{"provide a valid uuid"})
}

func ErrPersistEventToRemoteProvider(err error) error {
	return errors.New(ErrPersistEventToRemoteProviderCode, errors.Alert, []string{"failed to persist event to remote provider"}, []string{err.Error()}, []string{"token is expired/revoked", "Remote Provider is not reachable or unavailable"}, []string{"Try re-authenticating with the remote provider", "Verify remote provider for its reachability or availability."})
}
func ErrNoTarInsideOCi() error {
	return errors.New(ErrNoTarInsideOCiCode, errors.Alert, []string{"No tar file found inside OCI image"}, []string{"Unable to locate the compressed file(.tar.gz) inside the OCI image."}, []string{"The OCI image does not contain a ziped file."}, []string{"Verify that the OCI image contains a ziped file."})
}
func ErrEmptyOCIImage(err error) error {
	return errors.New(ErrEmptyOCIImageCode, errors.Alert, []string{}, []string{}, []string{}, []string{})
}

func ErrGetCapabilities(err error, userId string) error {
	return errors.New(ErrGetCapabilitiesCode, errors.Alert, []string{fmt.Sprintf("failed to get capabilities for the user with id: \"%s\"", userId)}, []string{err.Error()}, []string{"Remote provider server may be down or not accepting requests."}, []string{"Make sure remote provider server is healthy and accepting requests."})
}

func ErrExportPatternInFormat(err error, format, designName string) error {
	return errors.New(ErrExportPatternInFormatCode, errors.Alert, []string{fmt.Sprintf("Failed to export design file \"%s\" as \"%s\"", designName, format)}, []string{err.Error()}, []string{fmt.Sprintf("Current version of Meshery does not support exporting in \"%s\" format", format)}, []string{"Export design in one of the supported format."})
}
