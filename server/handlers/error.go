package handlers

import (
	"fmt"
	"strings"

	"github.com/meshery/meshkit/errors"
)

// isClientDisconnect checks if an error is caused by the client disconnecting
// (broken pipe or connection reset). These are expected during normal operation
// when clients close connections before the server finishes writing.
func isClientDisconnect(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer")
}

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
	ErrWalkingLocalDirectoryCode           = "meshery-server-1132"
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
	ErrFileTypeCode                        = "meshery-server-1366"
	ErrCreatingOPAInstanceCode             = "meshery-server-1367"
	ErrEncodeResponseCode                  = "meshery-server-1374"
	ErrTransientProviderCode               = "meshery-server-1378"
	ErrServeSchemaCode                     = "meshery-server-1381"
	ErrInvalidFileRequestCode              = "meshery-server-1382"
	ErrReadFileContentCode                 = "meshery-server-1383"
	ErrExtensionEndpointNotRegisteredCode  = "meshery-server-1384"
	ErrUserNotFoundCode                    = "meshery-server-1385"
	ErrFetchTokenCode                      = "meshery-server-1386"
	ErrHandlerShareDesignCode              = "meshery-server-1387"
	ErrHandlerShareFilterCode              = "meshery-server-1388"
	ErrGetUserCredentialCode               = "meshery-server-1389"
	ErrSaveUserCredentialCode              = "meshery-server-1390"
	ErrUpdateUserCredentialCode            = "meshery-server-1391"
	ErrDeleteUserCredentialCode            = "meshery-server-1392"
	ErrEncodeUserCredentialCode            = "meshery-server-1393"
	ErrUnknownConnectionKindCode           = "meshery-server-1394"
	ErrGetK8sContextsCode                  = "meshery-server-1395"
	ErrEncodeK8sContextsCode               = "meshery-server-1396"
	ErrCreateDatabaseArchiveDirCode        = "meshery-server-1397"
	ErrOpenDatabaseFileCode                = "meshery-server-1398"
	ErrCreateDatabaseArchiveFileCode       = "meshery-server-1399"
	ErrCopyDatabaseFileCode                = "meshery-server-1400"
	ErrObtainDatabaseHandlerCode           = "meshery-server-1401"
	ErrAccessDatabaseTablesCode            = "meshery-server-1402"
	ErrDropDatabaseTableCode               = "meshery-server-1403"
	ErrMigrateDatabaseTablesCode           = "meshery-server-1404"
	ErrFetchResultsCode                    = "meshery-server-1405"
	ErrMissingResultIDCode                 = "meshery-server-1406"
	ErrHandlerGenerateUUIDCode             = "meshery-server-1407"
	ErrMethodNotAllowedCode                = "meshery-server-1408"
	ErrMissingRouteVariableCode            = "meshery-server-1409"
	ErrRetrieveEventTypesCode              = "meshery-server-1410"
	ErrDeprecatedAPICode                   = "meshery-server-1411"
	ErrInvalidConnectionKindCode           = "meshery-server-1412"
	ErrUpdateConnectionCode                = "meshery-server-1413"
	ErrExportModelCode                     = "meshery-server-1414"
	ErrInvalidContextsConfigCode           = "meshery-server-1415"
	ErrEmptyConnectionIDCode               = "meshery-server-1416"
	ErrPolicyEvalTimeoutCode               = "meshery-server-1417"
	ErrPolicyEvalCode                      = "meshery-server-1418"
	ErrInvalidBase64DataCode               = "meshery-server-1419"
	ErrInvalidImportRequestCode            = "meshery-server-1422"
	ErrConvertToDesignCode                 = "meshery-server-1423"
	ErrCompressArtifactCode                = "meshery-server-1424"
	ErrWriteRegistryLogsCode               = "meshery-server-1425"
	ErrUpdateEntityStatusCode              = "meshery-server-1426"
	ErrExtensionProxyCode                  = "meshery-server-1427"
	ErrInitializeMachineCode               = "meshery-server-1428"
	ErrSendMachineEventCode                = "meshery-server-1429"
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

// ErrTransientProvider wraps a transient failure talking to the remote provider
// (e.g. Meshery Cloud). Callers should map this to HTTP 503 so clients can
// distinguish it from an auth failure and retry rather than logging the user out.
func ErrTransientProvider(err error) error {
	return errors.New(
		ErrTransientProviderCode,
		errors.Alert,
		[]string{"Remote provider temporarily unavailable"},
		[]string{fmt.Sprintf("Meshery Cloud (the remote provider) did not respond: %v", err)},
		[]string{"Network connectivity issue, Cloud outage, or firewall blocking egress."},
		[]string{"Retry the request. If this persists, check https://status.meshery.io and verify Meshery Server can reach Meshery Cloud."},
	)
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
	return errors.New(ErrPatternDeployCode, errors.Alert, []string{fmt.Sprintf("Unable to deploy the selected design \"%s\"", patternName)}, []string{err.Error()}, []string{"Connection Error: There was an error connecting to the selected target platform (i.e. Kubernetes cluster(s)).", "This connection might not be assigned to the selected environment."}, []string{"Verify that the Kubernetes connection status is 'Connected' or try uploading a new kubeconfig.", "Assign the current Kubernetes connection to the selected environment."})
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
	return errors.New(ErrFormFileCode, errors.Alert, []string{"error getting kubeconfig file"}, []string{err.Error()}, []string{"The kubeconfig file does not exist in the location."}, []string{"Make sure to upload the correct kubeconfig file."})
}

func ErrReadConfig(err error) error {
	return errors.New(ErrReadConfigCode, errors.Alert, []string{"error reading config."}, []string{err.Error()}, []string{"The kubeconfig file is empty or not valid."}, []string{"Make sure to upload the correct kubeconfig file."})
}

func ErrLoadConfig(err error) error {
	return errors.New(ErrLoadConfigCode, errors.Alert, []string{"unable to load Kubernetes config"}, []string{err.Error()}, []string{"The kubeconfig file is empty or not valid"}, []string{"Make sure to upload the correct kubeconfig file"})
}

func ErrOpenFile(file string) error {
	return errors.New(ErrOpenFileCode, errors.Alert, []string{"Unable to open file: ", file}, []string{}, []string{"The file does not exist in the location."}, []string{"Make sure to upload the correct file."})
}

func ErrKubeVersion(err error) error {
	return errors.New(ErrKubeVersionCode, errors.Alert, []string{"Unable to get Kubernetes version."}, []string{err.Error()}, []string{"Kubernetes might not be reachable from Meshery."}, []string{"Make sure Meshery has connectivity to Kubernetes."})
}

func ErrRetrieveData(err error) error {
	return errors.New(ErrRetrieveDataCode, errors.Alert, []string{"Unable to retrieve the requested data"}, []string{err.Error()}, []string{"Adapter operation invalid"}, []string{"Make sure adapter is running and reachable by Meshery Server."})
}

func ErrOperationID(err error) error {
	return errors.New(ErrOperationIDCode, errors.Alert, []string{"Error generating the operation Id"}, []string{err.Error()}, []string{"Adapter operation invalid"}, []string{"Make sure adapter is reachable and running"})
}

func ErrApplyChange(err error) error {
	return errors.New(ErrApplyChangeCode, errors.Alert, []string{"Error applying the change"}, []string{err.Error()}, []string{"Adapter operation invalid"}, []string{"Make sure adapter is reachable and running"})
}

func ErrRetrieveMeshData(err error) error {
	return errors.New(ErrRetrieveMeshDataCode, errors.Alert, []string{"Error getting operations for the adapter.", "Error getting component name."}, []string{err.Error()}, []string{"Unable to retrieve the requested data."}, []string{"Make sure adapter is reachable and running."})
}

func ErrApplicationFailure(err error, obj string) error {
	return errors.New(ErrApplicationFailureCode, errors.Alert, []string{"failed to ", obj, "the application"}, []string{err.Error()}, []string{"uploaded application source content might be converted", "incorrect source type selected"}, []string{"Select the correct source type", "Make sure the uploaded application source content is valid"})
}

func ErrApplicationSourceContent(err error, obj string) error {
	return errors.New(ErrApplicationContentCode, errors.Alert, []string{"failed to ", obj, "the application content"}, []string{err.Error()}, []string{"Remote provider might be not reachable", "Remote provider doesn't support this capability"}, []string{"Ensure you have required permissions or retry after sometime."})
}

func ErrDesignSourceContent(err error, obj string) error {
	return errors.New(ErrDesignSourceContentCode, errors.Alert, []string{"failed to ", obj, "the design content"}, []string{err.Error()}, []string{"Remote provider might be not reachable.", "Remote provider doesn't support this capability."}, []string{"Ensure you have required permissions or retry after sometime."})
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
	return errors.New(ErrBlankNameCode, errors.Alert, []string{"Error: name field is blank."}, []string{err.Error()}, []string{"Load test name empty or not valid."}, []string{"Provide a name for the test."})
}

func ErrConversion(err error) error {
	return errors.New(ErrConversionCode, errors.Alert, []string{"unable to convert YAML to JSON"}, []string{err.Error()}, []string{"Yaml provided is not valid"}, []string{"Make sure the yaml is valid and has the right parameters"})
}

func ErrLoadTest(err error, obj string) error {
	return errors.New(ErrLoadTestCode, errors.Alert, []string{"Load test error: ", obj}, []string{err.Error()}, []string{"Load test endpoint might be not reachable."}, []string{"Make sure load test endpoint is reachable by Meshery Server."})
}

func ErrFetchKubernetes(err error) error {
	return errors.New(ErrFetchKubernetesCode, errors.Alert, []string{"Unable to ping Kubernetes.", "Unable to run a discovery scan for Kubernetes cluster."}, []string{err.Error()}, []string{"Kubernetes might not be reachable from Meshery."}, []string{"Make sure Meshery has connectivity to Kubernetes."})
}

func ErrPanicRecovery(r interface{}) error {
	return errors.New(ErrPanicRecoveryCode, errors.Alert, []string{"Recovered from panic."}, []string{fmt.Sprint(r)}, []string{"Meshery crashed."}, []string{"Restart Meshery."})
}

func ErrFailToLoadExtensions(err error) error {
	return errors.New(ErrFailToLoadExtensionsCode, errors.Alert, []string{"Failed to Load Extensions from Package."}, []string{err.Error()}, []string{"Plugin is not available in the location.", "Plugin does not match with Meshery version."}, []string{"Make sure the plugin is compatible with Meshery server."})
}

func ErrInvalidLTURL(url string) error {
	return errors.New(ErrInvalidLTURLCode, errors.Alert, []string{"invalid load test url: ", url}, []string{}, []string{"URL for load test could be invalid."}, []string{"Please refer to: https://docs.meshery.io/tasks/performance-management"})
}

func ErrVersionCompare(err error) error {
	return errors.New(ErrVersionCompareCode, errors.Alert, []string{"failed to compare latest and current version of Meshery."}, []string{err.Error()}, []string{}, []string{})
}

func ErrGetLatestVersion(err error) error {
	return errors.New(ErrGetLatestVersionCode, errors.Alert, []string{"failed to get latest version of Meshery."}, []string{err.Error()}, []string{}, []string{})
}

func ErrSaveSession(err error) error {
	return errors.New(ErrSaveSessionCode, errors.Alert, []string{"unable to save session."}, []string{err.Error()}, []string{"User session could be expired."}, []string{"Re-initiate login."})
}

func ErrCreateDir(err error, obj string) error {
	return errors.New(ErrCreateDirCode, errors.Alert, []string{"Error creating directory ", obj}, []string{err.Error()}, []string{"Insufficient permission", "Insufficient storage"}, []string{"check if sufficient permissions are available to create dir", "check if sufficient storage is available to create dir"})
}

func ErrInvalidRequestObject(fields ...string) error {
	return errors.New(ErrInvalidRequestObjectCode, errors.Alert, []string{"Error invalid request object:"}, []string{strings.Join(fields, " ")}, []string{""}, []string{""})
}

func ErrChangeK8sContext(err error) error {
	return errors.New(ErrChangeK8sContextCode, errors.Alert, []string{"Error changing context"}, []string{err.Error()}, []string{"Context Name might be invalid or not present in the uploaded kubeconfig"}, []string{"Check the context name, if the context name is correct and is present in the kubeconfig, then try uploading the kubeconfig again."})
}

func ErrInvalidKubeConfig(err error, content string) error {
	return errors.New(ErrInvalidKubeConfigCode, errors.Alert, []string{"Invalid Kube Config ", content}, []string{err.Error()}, []string{"Meshery handler failed to find a valid Kubernetes config for the deployment"}, []string{"Try uploading a new kubeconfig and also ensure that Meshery can reach Kubernetes API server"})
}

func ErrInvalidKubeHandler(err error, userId string) error {
	return errors.New(ErrInvalidKubeHandlerCode, errors.Alert, []string{"Kubernetes cluster is unavailable for ", userId}, []string{err.Error()}, []string{"There might be a network disruption or the Meshery server does not have valid credentials."}, []string{"Try uploading a new kubeconfig.", "Check the network connection and Kubernetes cluster status.", "Verify that the Meshery server has valid and updated credentials to access the Kubernetes cluster."})
}

func ErrInvalidKubeContext(err error, content string) error {
	return errors.New(ErrInvalidKubeContextCode, errors.Alert, []string{"Invalid Kube Context", content}, []string{err.Error()}, []string{"Meshery handler failed to find a valid Kubernetes context for the deployment."}, []string{"Try uploading a new kubeconfig and also ensure that Meshery can reach Kubernetes API server."})
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
	return errors.New(ErrFetchFilterCode, errors.Alert, []string{"Error failed to fetch filter."}, []string{err.Error()}, []string{"Failed to retrieve the list of all the Filters"}, []string{})
}

func ErrDeleteFilter(err error) error {
	return errors.New(ErrDeleteFilterCode, errors.Alert, []string{"Error failed to delete filter."}, []string{err.Error()}, []string{"Failed to delete Filter with the given ID"}, []string{"Check if the Filter ID is correct"})
}

func ErrSavePattern(err error) error {
	return errors.New(ErrSavePatternCode, errors.Alert, []string{"Error failed to save design."}, []string{err.Error()}, []string{"Cannot save the design due to an invalid path or URL"}, []string{"Verify that you have an active user session. Try logging and in again.", "Confirm that you have sufficient permissions to save the design.", "Try reducing the size of the design file by removing the number of images, using alternative image formats or removing other non-critical components from the design. See https://docs.layer5.io/kanvas/advanced/performance/."})
}

func ErrSaveApplication(err error) error {
	return errors.New(ErrSaveApplicationCode, errors.Alert, []string{"Error failed to save application."}, []string{err.Error()}, []string{"Cannot save the Application due to wrong path or URL."}, []string{"Check if the given path or URL of the Application is correct."})
}

func ErrFetchApplication(err error) error {
	return errors.New(ErrFetchApplicationCode, errors.Alert, []string{"Error failed to fetch applications"}, []string{err.Error()}, []string{"Remote provider might be not reachable.", "Token might have expired."}, []string{"Refresh your browser"})
}

func ErrDeleteApplication(err error) error {
	return errors.New(ErrDeleteApplicationCode, errors.Alert, []string{"Error failed to delete application"}, []string{err.Error()}, []string{"Application might already have been deleted", "You might not have enough permissions to perform the operation."}, []string{"Check the owner of the application."})
}

func ErrGetPattern(err error) error {
	return errors.New(
		ErrGetPatternCode,
		errors.Alert,
		[]string{"Unable to open this design"},
		[]string{fmt.Sprintf("The server could not return the requested design. Underlying error: %v", err)},
		[]string{
			"The design ID in the URL is malformed or does not exist.",
			"The design has been deleted by its owner.",
			"Your account does not have permission to view this design — it may belong to another organization or be set to private.",
			"Your session has expired or the remote provider is currently unreachable.",
		},
		[]string{
			"Verify the design link — confirm the full design ID is intact, with no missing or extra characters.",
			"Open My Designs and confirm the design still exists; if it was shared with you, ask the owner to re-share or grant access.",
			"Sign out and sign back in to refresh your session, then retry.",
		},
	)
}

func ErrDeletePattern(err error) error {
	return errors.New(
		ErrDeletePatternCode,
		errors.Alert,
		[]string{"Unable to delete this design"},
		[]string{fmt.Sprintf("The server could not delete the requested design. Underlying error: %v", err)},
		[]string{
			"The design ID is malformed or no longer exists — it may already have been deleted.",
			"Your account does not have permission to delete this design — only authorized users (such as the owner) can delete it.",
			"Your session has expired or the remote provider is currently unreachable.",
		},
		[]string{
			"Refresh the designs list to confirm the design still exists before retrying.",
			"If the design was shared with you, ask its owner to delete it.",
			"Sign out and sign back in to refresh your session, then retry.",
		},
	)
}

func ErrFetchPattern(err error) error {
	return errors.New(
		ErrFetchPatternCode,
		errors.Alert,
		[]string{"Unable to load your designs"},
		[]string{fmt.Sprintf("The server could not retrieve the list of designs. Underlying error: %v", err)},
		[]string{
			"The remote provider is currently unreachable or returned an error.",
			"Your session has expired or your authentication token is no longer valid.",
			"Network connectivity between Meshery Server and the remote provider has been interrupted.",
		},
		[]string{
			"Check your network connection and reload the page.",
			"Sign out and sign back in to refresh your session, then retry.",
			"Confirm that the configured remote provider is online.",
		},
	)
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

// ErrConvertPattern wraps failures in the in-place schema migration that
// converts a stored design from an older schema version (e.g. v1alpha2)
// to the current v1beta3 PatternFile shape. Emitted with HTTP 500
// because the failure originates server-side during a download/view of
// an already-persisted design, not from caller input.
func ErrConvertPattern(err error) error {
	return errors.New(ErrConvertPatternCode, errors.Alert, []string{"Failed to migrate design to current schema version"}, []string{err.Error()}, []string{"The persisted design is in an older schema version (e.g. v1alpha2) whose ConvertFrom path rejected one of its fields.", "A field present in the older schema has no clean mapping to the current v1beta3 PatternFile shape."}, []string{"Inspect server logs for the underlying conversion error. If the design is irrecoverable, re-import its source content via the design import endpoint, which re-runs the full conversion pipeline."})
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

// ErrRetrieveEventTypes wraps failures from provider.GetEventTypes — the
// endpoint that returns the catalogue of event categories + actions the
// UI uses to populate the notification filter drop-downs.
func ErrRetrieveEventTypes(err error) error {
	return errors.New(ErrRetrieveEventTypesCode, errors.Alert, []string{"Could not retrieve event categories and actions"}, []string{err.Error()}, []string{"Database is not reachable or corrupt.", "Remote provider event-type endpoint is unavailable."}, []string{"Check database connectivity and provider availability, then retry."})
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
func ErrFileType(obj string) error {
	return errors.New(ErrFileTypeCode, errors.Alert, []string{"Error: ", obj, " is not a valid file type"}, []string{fmt.Sprintf("file is not the expected type %splease verify the file type", obj)}, []string{"The provided file type is wrong"}, []string{"Verify the file type"})
}
func ErrIOReader(err error) error {
	return errors.New(ErrIOReaderCode, errors.Alert, []string{"Failed to read from io.Reader"}, []string{err.Error()}, []string{"unable to read from io.Reader"}, []string{"check if the io.Reader is valid"})
}

func ErrUnCompressOCIArtifact(err error) error {
	return errors.New(ErrUnCompressOCIArtifactCode, errors.Alert, []string{"Failed to uncompress OCI artifact"}, []string{err.Error()}, []string{"unable to uncompress OCI artifact", "OCI artifact may be corrupted"}, []string{"check if the OCI artifact is valid and not corrupted"})
}

func ErrWalkingLocalDirectory(err error, path string) error {
	return errors.New(ErrWalkingLocalDirectoryCode, errors.Alert, []string{"Failed to walk local directory: ", path}, []string{err.Error()}, []string{"unable to walk local directory at ", path, "local directory may be corrupted or inaccessible"}, []string{"check if the local directory is valid and has correct permissions"})
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
func ErrEncodeResponse(err error) error {
	return errors.New(ErrEncodeResponseCode, errors.Alert, []string{"failed to encode and write response"}, []string{err.Error()}, []string{"Client may have disconnected before the response was fully written", "Response data could not be serialized to JSON"}, []string{"If the issue persists, check server logs for details"})
}
func ErrCreatingOPAInstance(err error) error {
	return errors.New(ErrCreatingOPAInstanceCode, errors.Alert, []string{"Error creating OPA Instance."}, []string{err.Error()}, []string{"Unable to create OPA instance, policies will not be evaluated."}, []string{"Ensure relationships are registered"})
}
func ErrServeSchema(err error) error {
	return errors.New(ErrServeSchemaCode, errors.Alert, []string{"Failed to serve the requested schema"}, []string{err.Error()}, []string{"Requested resource's schema could not be found or read"}, []string{"Ensure the resource name is spelled correctly and that the schema is bundled with the server"})
}
func ErrInvalidFileRequest(err error) error {
	return errors.New(ErrInvalidFileRequestCode, errors.Alert, []string{"Invalid file request"}, []string{err.Error()}, []string{"The provided file query parameter could not be decoded"}, []string{"Ensure the file parameter is a properly URL-encoded path"})
}
func ErrReadFileContent(err error, file string) error {
	return errors.New(ErrReadFileContentCode, errors.Alert, []string{"Failed to read file content", file}, []string{err.Error()}, []string{"The file could not be opened or streamed to the response"}, []string{"Verify the file exists and the server has permission to read it"})
}
func ErrExtensionEndpointNotRegistered(endpoint string) error {
	return errors.New(ErrExtensionEndpointNotRegisteredCode, errors.Alert, []string{"No extension is registered for endpoint: ", endpoint}, []string{}, []string{"Requested extension is not loaded into this Meshery server"}, []string{"Install the extension or check that its provider registered the route on startup"})
}

func ErrUserNotFound(userID string) error {
	return errors.New(ErrUserNotFoundCode, errors.Alert, []string{"User not found"}, []string{fmt.Sprintf("No user exists with id: %s", userID)}, []string{"The user may have been deleted or the id is incorrect"}, []string{"Verify the user id is correct and that the user has not been removed"})
}

func ErrFetchToken(err error) error {
	return errors.New(ErrFetchTokenCode, errors.Alert, []string{"Failed to obtain authentication token from request context"}, []string{err.Error()}, []string{"The request may not be authenticated or the token was not propagated into the context"}, []string{"Ensure the request is authenticated and that upstream middleware injects the token into context"})
}

func ErrShareDesign(err error) error {
	return errors.New(ErrHandlerShareDesignCode, errors.Alert, []string{"Failed to share design"}, []string{err.Error()}, []string{"The remote provider rejected the share request", "Network connectivity issue", "Invalid payload"}, []string{"Verify the share request payload and target recipients, then retry"})
}

func ErrShareFilter(err error) error {
	return errors.New(ErrHandlerShareFilterCode, errors.Alert, []string{"Failed to share filter"}, []string{err.Error()}, []string{"The remote provider rejected the share request", "Network connectivity issue", "Invalid payload"}, []string{"Verify the share request payload and target recipients, then retry"})
}

func ErrGetUserCredential(err error) error {
	return errors.New(ErrGetUserCredentialCode, errors.Alert, []string{"Failed to get user credential"}, []string{err.Error()}, []string{"Credential may not exist, caller may lack access, or the provider is unreachable"}, []string{"Verify the credential id and that the current user has access to it"})
}

func ErrSaveUserCredential(err error) error {
	return errors.New(ErrSaveUserCredentialCode, errors.Alert, []string{"Failed to save user credential"}, []string{err.Error()}, []string{"Invalid credential payload or the provider rejected the request"}, []string{"Verify the credential payload matches the expected schema and retry"})
}

func ErrUpdateUserCredential(err error) error {
	return errors.New(ErrUpdateUserCredentialCode, errors.Alert, []string{"Failed to update user credential"}, []string{err.Error()}, []string{"Credential may not exist, caller may lack access, or the payload is invalid"}, []string{"Verify the credential id, caller permissions, and payload, then retry"})
}

func ErrDeleteUserCredential(err error) error {
	return errors.New(ErrDeleteUserCredentialCode, errors.Alert, []string{"Failed to delete user credential"}, []string{err.Error()}, []string{"Credential may not exist or the caller may lack permission"}, []string{"Verify the credential id and caller permissions, then retry"})
}

func ErrEncodeUserCredential(err error) error {
	return errors.New(ErrEncodeUserCredentialCode, errors.Alert, []string{"Failed to encode user credential response"}, []string{err.Error()}, []string{"Credential payload could not be serialized to JSON"}, []string{"Check server logs for serialization details"})
}

func ErrUnknownConnectionKind(kind string) error {
	return errors.New(ErrUnknownConnectionKindCode, errors.Alert, []string{"Unable to register resource as connection"}, []string{fmt.Sprintf("No matching connection definition found in the registry for kind: %s", kind)}, []string{"The specified connection kind is not registered with Meshery"}, []string{"Verify the connection kind is spelled correctly and that its definition is installed"})
}

func ErrGetK8sContexts(err error) error {
	return errors.New(ErrGetK8sContextsCode, errors.Alert, []string{"Failed to get Kubernetes contexts"}, []string{err.Error()}, []string{"The remote provider may be unreachable or the stored context data is corrupted"}, []string{"Verify the remote provider is reachable and retry; if the issue persists, check server logs"})
}

func ErrEncodeK8sContexts(err error) error {
	return errors.New(ErrEncodeK8sContextsCode, errors.Alert, []string{"Failed to encode Kubernetes context response"}, []string{err.Error()}, []string{"Response data could not be serialized to JSON"}, []string{"Check server logs for serialization details"})
}

func ErrCreateDatabaseArchiveDir(err error) error {
	return errors.New(ErrCreateDatabaseArchiveDirCode, errors.Alert, []string{"Directory could not be created due to a non-existent path"}, []string{err.Error()}, []string{"Meshery's configuration directory does not exist or is not writable"}, []string{"Verify that ~/.meshery/config exists and that the server process has write permission"})
}

func ErrOpenDatabaseFile(err error) error {
	return errors.New(ErrOpenDatabaseFileCode, errors.Alert, []string{"The database does not exist or you don't have enough permission to access it"}, []string{err.Error()}, []string{"mesherydb.sql is missing or not readable by the server process"}, []string{"Verify that mesherydb.sql exists under ~/.meshery/config and that the server process can read it"})
}

func ErrCreateDatabaseArchiveFile(err error) error {
	return errors.New(ErrCreateDatabaseArchiveFileCode, errors.Alert, []string{"Destination file can not be created"}, []string{err.Error()}, []string{"Archive target path is not writable or disk is full"}, []string{"Ensure ~/.meshery/config/.archive is writable and has free space"})
}

func ErrCopyDatabaseFile(err error) error {
	return errors.New(ErrCopyDatabaseFileCode, errors.Alert, []string{"Can not copy file from source to destination"}, []string{err.Error()}, []string{"I/O error while archiving the current database, or disk is full"}, []string{"Check disk space and filesystem permissions for ~/.meshery/config"})
}

func ErrObtainDatabaseHandler() error {
	return errors.New(ErrObtainDatabaseHandlerCode, errors.Alert, []string{"Failed to obtain database handler"}, []string{"Provider returned a nil generic persister"}, []string{"The configured provider does not expose a persister, or initialization has not completed"}, []string{"Verify the active provider supports system database reset and that the server is fully initialized"})
}

func ErrAccessDatabaseTables(err error) error {
	return errors.New(ErrAccessDatabaseTablesCode, errors.Alert, []string{"Can not access database tables"}, []string{err.Error()}, []string{"The database file is locked, corrupt, or inaccessible"}, []string{"Check that no other process is holding a lock on mesherydb and that the file is not corrupt"})
}

func ErrDropDatabaseTable(err error) error {
	return errors.New(ErrDropDatabaseTableCode, errors.Alert, []string{"Cannot drop table from database"}, []string{err.Error()}, []string{"Migrator could not drop an existing table — schema or connection issue"}, []string{"Check server logs for details and verify database integrity"})
}

func ErrMigrateDatabaseTables(err error) error {
	return errors.New(ErrMigrateDatabaseTablesCode, errors.Alert, []string{"Can not migrate tables to database"}, []string{err.Error()}, []string{"Auto-migrate or registry manager setup failed during database reset"}, []string{"Check server logs for migration details and verify the database is accessible"})
}

func ErrFetchResults(err error) error {
	return errors.New(ErrFetchResultsCode, errors.Alert, []string{"Error while getting load test results"}, []string{err.Error()}, []string{"The remote provider is unreachable, the profile id is invalid, or the results store is unavailable"}, []string{"Verify connectivity to the remote provider and retry; confirm the profile id is correct"})
}

func ErrMissingResultID() error {
	return errors.New(ErrMissingResultIDCode, errors.Alert, []string{"Missing result id in request"}, []string{"No result id was supplied in the URL path"}, []string{"The client did not include a result identifier"}, []string{"Provide the result id in the request URL, for example /api/user/performance/results/{id}"})
}

func ErrGenerateUUID(err error) error {
	return errors.New(ErrHandlerGenerateUUIDCode, errors.Alert, []string{"Failed to generate a unique identifier"}, []string{err.Error()}, []string{"Secure random source is unavailable on this host"}, []string{"Retry the request; if the error persists, verify the host's /dev/urandom or equivalent entropy source is accessible"})
}

func ErrMethodNotAllowed(method string) error {
	return errors.New(ErrMethodNotAllowedCode, errors.Alert, []string{"HTTP method not allowed for this endpoint"}, []string{fmt.Sprintf("Received %s, but the endpoint only accepts the methods registered on its route", method)}, []string{"The client used an HTTP verb the route does not support"}, []string{"Use one of the supported HTTP methods for this endpoint"})
}

// ErrMissingRouteVariable wraps requests where a required path/route variable
// was not supplied by the caller. Used by handlers that rely on mux.Vars(r)
// entries the router is expected to provide.
func ErrMissingRouteVariable(name string, allowed ...string) error {
	cause := fmt.Sprintf("Required route variable %q was missing from the request", name)
	remedy := fmt.Sprintf("Include %q in the request path", name)
	if len(allowed) > 0 {
		remedy = fmt.Sprintf("Include %q (one of %s) in the request path", name, strings.Join(allowed, ", "))
	}
	return errors.New(ErrMissingRouteVariableCode, errors.Alert, []string{fmt.Sprintf("Missing required route variable %q", name)}, []string{cause}, []string{"The client called the endpoint without supplying a required path segment."}, []string{remedy})
}

// ErrDeprecatedAPI is returned when a client invokes an endpoint that has
// been retired. Emitted with HTTP 410 Gone. The `replacement` argument
// should name the new endpoint so the UI can surface a clear upgrade path.
func ErrDeprecatedAPI(replacement string) error {
	return errors.New(ErrDeprecatedAPICode, errors.Alert, []string{"This API is deprecated"}, []string{fmt.Sprintf("The requested endpoint has been removed in favor of %s.", replacement)}, []string{"Client is using an older Meshery UI or mesheryctl that still targets the retired endpoint."}, []string{fmt.Sprintf("Use %s instead, and update the client to the latest Meshery release.", replacement)})
}

// ErrInvalidConnectionKind is returned when an endpoint scoped to a specific
// connection kind (e.g. prometheus, grafana) is called with a connection
// whose Kind does not match. Emitted with HTTP 400.
func ErrInvalidConnectionKind(actual, expected string) error {
	return errors.New(ErrInvalidConnectionKindCode, errors.Alert, []string{fmt.Sprintf("Connection is not of kind %q", expected)}, []string{fmt.Sprintf("The referenced connection is of kind %q but the endpoint expects %q.", actual, expected)}, []string{"The connectionID in the URL references a connection created for a different integration."}, []string{fmt.Sprintf("Pass the ID of a %s connection, or use the endpoint that matches the %s connection kind.", expected, actual)})
}

// ErrUpdateConnection wraps failures persisting connection metadata changes
// (e.g. board selections on prometheus/grafana connections). Emitted with
// HTTP 500.
func ErrUpdateConnection(err error) error {
	return errors.New(ErrUpdateConnectionCode, errors.Alert, []string{"Could not update the connection"}, []string{err.Error()}, []string{"Remote provider is unreachable.", "Connection has been deleted since it was loaded.", "Persisted metadata is corrupt."}, []string{"Verify provider connectivity and that the connection still exists, then retry."})
}

// ErrExportModel wraps failures in the ExportModel pipeline — building the
// OCI image, writing the model definition, saving the tar/gzip archive, or
// creating the scratch directories. The origin string identifies which
// sub-step failed so the UI can surface a specific remediation.
func ErrExportModel(err error, stage string) error {
	return errors.New(ErrExportModelCode, errors.Alert, []string{fmt.Sprintf("Failed to export model during %s", stage)}, []string{err.Error()}, []string{"Temp directory is not writable.", "Model definition references a missing dependency.", "OCI/tar tooling produced an invalid artifact."}, []string{"Retry the export; if it persists, inspect server logs and disk availability under the temp directory."})
}

// ErrInvalidContextsConfig wraps failures to parse the `contexts` form field
// of the kubeconfig upload endpoint. Emitted with HTTP 400.
func ErrInvalidContextsConfig(err error) error {
	return errors.New(ErrInvalidContextsConfigCode, errors.Alert, []string{"Invalid contexts configuration"}, []string{err.Error()}, []string{"The `contexts` form field is not a valid JSON map from context ID to {meshsync_deployment_mode: string}."}, []string{"Send a well-formed JSON object in the `contexts` form field, e.g. {\"<contextID>\": {\"meshsync_deployment_mode\": \"operator\"}}."})
}

// ErrEmptyConnectionID is returned from endpoints that require a connection
// ID query parameter when none is supplied. Emitted with HTTP 400.
func ErrEmptyConnectionID() error {
	return errors.New(ErrEmptyConnectionIDCode, errors.Alert, []string{"Empty connection ID"}, []string{"No connection ID was supplied in the canonical `connectionId` query parameter (the legacy `connection_id` spelling is also accepted during the Phase 2 deprecation window)."}, []string{"The client did not pass `connectionId` in the query string.", "A URL template did not get its parameter substituted."}, []string{"Pass the connection ID of the Kubernetes context to be pinged, e.g. /api/system/kubernetes/ping?connectionId=<id>."})
}

// ErrPolicyEvalTimeout wraps the `errEvalTimeout` sentinel from
// policy_relationship_handler.go with structured MeshKit metadata when the
// handler surfaces it over the wire. Emitted with HTTP 504 Gateway Timeout.
func ErrPolicyEvalTimeout(timeout fmt.Stringer) error {
	return errors.New(ErrPolicyEvalTimeoutCode, errors.Alert, []string{"Relationship policy evaluation timed out"}, []string{fmt.Sprintf("The relationship policy evaluator did not return a response within %s.", timeout)}, []string{"The design is large enough that the policy engine needs more than the configured timeout.", "One of the Rego/Go policies under evaluation is in a slow path or infinite loop."}, []string{"Increase the POLICY_EVAL_TIMEOUT setting (default 3m) if the design genuinely needs more time, or retry with a smaller design."})
}

// ErrPolicyEval wraps a generic policy-evaluation failure returned by the
// relationship evaluator (OPA/Rego or the Go engine). Emitted with HTTP 500.
func ErrPolicyEval(err error) error {
	return errors.New(ErrPolicyEvalCode, errors.Alert, []string{"Relationship policy evaluation failed"}, []string{err.Error()}, []string{"A registered relationship policy returned an error during evaluation.", "A cycle or invalid declaration in the design triggered an engine panic recovered as an error."}, []string{"Inspect server logs for the underlying policy error. If the design has recently been edited, revert the most recent change and retry."})
}

// ErrInvalidBase64Data wraps a base64 decoding failure on a request payload
// (e.g. a model file uploaded as a base64-encoded blob in the import body).
// Emitted with HTTP 400 because the client supplied malformed input.
func ErrInvalidBase64Data(err error) error {
	return errors.New(ErrInvalidBase64DataCode, errors.Alert, []string{"Invalid base64 data"}, []string{err.Error()}, []string{"The supplied payload was not a valid base64-encoded string.", "The payload may have been corrupted in transit or encoded with the wrong alphabet (URL-safe vs. standard)."}, []string{"Verify the client is base64-encoding the file with the standard alphabet before sending it.", "If the payload was hand-edited, re-encode it from the source bytes and retry."})
}

// ErrInvalidImportRequest wraps oneOf-invariant violations on the design
// import endpoint — e.g. the request body had both a File and URL variant
// set, or neither. Emitted with HTTP 400 because the caller needs to
// correct the request shape, not the server to recover.
func ErrInvalidImportRequest(err error) error {
	return errors.New(ErrInvalidImportRequestCode, errors.Alert, []string{"Invalid design import request"}, []string{err.Error()}, []string{"The request body did not match exactly one variant of the import oneOf — the File variant requires `file` and `file_name`, the URL variant requires `url`.", "Both variants were provided, or neither was."}, []string{"Send a request body with exactly one variant set: either {\"file\": <bytes>, \"file_name\": \"design.yml\"} or {\"url\": \"https://...\"}."})
}

// ErrConvertToDesign wraps failures in the conversion pipeline that
// turns a source file (Helm chart, Kubernetes manifest, Docker Compose,
// Kustomize, or a Meshery design) into a v1beta3 design. These failures
// are often rooted in malformed or unsupported input — a corrupt
// archive, an unrecognized file extension, or a manifest the registry
// could not map onto known component definitions — but the same
// pipeline is also re-run server-side during download/view of a
// non-design pattern, where the same wrap surfaces. The HTTP status
// returned for this error is therefore determined by the calling
// handler and request flow: 400 for the import endpoint where the
// input is the request body, 500 for download/view where the input is
// already-persisted SourceContent.
func ErrConvertToDesign(err error) error {
	return errors.New(ErrConvertToDesignCode, errors.Alert, []string{"Failed to convert uploaded file to a design"}, []string{err.Error()}, []string{"The uploaded file extension is not one of the supported import formats (.yml, .yaml, .json, .tar, .tar.gz, .tgz, .zip).", "The file is corrupt or its content does not parse as the type implied by its extension.", "The Kubernetes manifest references a kind/apiVersion that does not match any registered component definition."}, []string{"Verify the file is one of the supported formats and content-types, and that it parses cleanly outside Meshery.", "If the source is a Kubernetes manifest, ensure each kind it references has a corresponding model registered in the Meshery registry."})
}

// ErrCompressArtifact wraps tar-writer failures encountered while
// packaging a Meshery design and its companion artifacthub-pkg.yml into
// an OCI artifact for download. Emitted with HTTP 500 because the
// failure originates server-side (in-memory buffer / archive writer)
// rather than from caller input.
func ErrCompressArtifact(err error) error {
	return errors.New(ErrCompressArtifactCode, errors.Alert, []string{"Failed to compress design as OCI artifact"}, []string{err.Error()}, []string{"The server-side tar writer hit an I/O error while packaging the design YAML and the artifacthub-pkg metadata.", "Available memory or disk for the in-memory archive buffer was exhausted."}, []string{"Retry the export. If the failure persists, inspect server logs for the underlying writer error and ensure the host has sufficient resources."})
}

// ErrWriteRegistryLogs wraps failures of the helpers.WriteLogsToFiles
// post-registration step that flushes per-host registration attempts
// (component / model / relationship / policy / registry) to the file
// pointed to by the REGISTRY_LOG_FILE setting. Emitted with HTTP 500
// because the error is internal to the registry-log subsystem (file
// permissions, disk full, marshal failure) rather than caller input.
func ErrWriteRegistryLogs(err error) error {
	return errors.New(ErrWriteRegistryLogsCode, errors.Alert, []string{"Failed to flush registry logs"}, []string{err.Error()}, []string{"The path configured in REGISTRY_LOG_FILE is not writable by the Meshery process (missing directory, missing permissions, or the disk is full).", "An in-memory log entry could not be marshaled into the on-disk format."}, []string{"Verify REGISTRY_LOG_FILE points to a writable path with sufficient free disk and that the Meshery process owns or has write permission to its parent directory."})
}

// ErrUpdateEntityStatus wraps registry failures when toggling the
// `status` (enabled/disabled/etc.) of a model or component definition
// through the entity-status endpoint. Emitted with HTTP 500 because
// the failure originates in the registry persistence layer, not in
// caller input — input validation runs upstream of this site.
func ErrUpdateEntityStatus(err error) error {
	return errors.New(ErrUpdateEntityStatusCode, errors.Alert, []string{"Failed to update entity status"}, []string{err.Error()}, []string{"The entity ID does not exist in the registry.", "The registry's persistence layer rejected the status update — typically a database connection or transaction failure."}, []string{"Verify the entity ID exists by listing the entities of the same type. If it does, retry the request and inspect server logs for the underlying registry error."})
}

// ErrExtensionProxy wraps failures of provider.ExtensionProxy when a
// `/api/extensions/...` passthrough request cannot be served. Caller
// chooses the HTTP status: 502 Bad Gateway for upstream/remote-provider
// failures (the common case), 501 Not Implemented when the active
// provider is the local provider, which has no extensions backend.
func ErrExtensionProxy(err error) error {
	return errors.New(ErrExtensionProxyCode, errors.Alert, []string{"Extension proxy request failed"}, []string{err.Error()}, []string{"The remote provider could not be reached (network failure, DNS resolution, or TLS handshake failure).", "The remote provider returned a non-2xx response that the proxy could not translate.", "The active provider is the local provider, which has no extensions backend."}, []string{"Verify the remote provider is reachable from this Meshery instance and that the user's session token has not expired. Retry the request after re-authenticating if the failure persists.", "If running against the local provider, switch to a remote provider (Layer5 Cloud) that exposes the extensions surface."})
}

// ErrInitializeMachine wraps failures of the connection state-machine
// initializer (machines/helpers.InitializeMachineWithContext). Emitted
// with HTTP 500 because the failure is internal to the state-machine
// bootstrap (unknown machine type, persistence failure, or initial
// state assignment failure) rather than caller input.
func ErrInitializeMachine(err error) error {
	return errors.New(ErrInitializeMachineCode, errors.Alert, []string{"Failed to initialize connection state machine"}, []string{err.Error()}, []string{"The connection kind has no registered machine type in the state-machine registry.", "The state machine's persister could not write the initial state to the database."}, []string{"Verify the connection kind is one of the supported types. If it is, retry the request and inspect server logs for the underlying machine-bootstrap error."})
}

// ErrSendMachineEvent wraps failures of *StateMachine.SendEvent, which
// drives a connection through its registered transitions (e.g.
// REGISTERED → DISCOVERED → CONNECTED). Emitted with HTTP 500 because
// the event-driven transition failed inside the state machine, not in
// caller input.
func ErrSendMachineEvent(err error) error {
	return errors.New(ErrSendMachineEventCode, errors.Alert, []string{"Failed to advance connection state machine"}, []string{err.Error()}, []string{"The requested event is not valid from the connection's current state.", "A side-effect action attached to the transition (e.g. provisioning, discovery) returned an error."}, []string{"Inspect the connection's current status before retrying. If the failure originates from a side-effect action, address the underlying cause (e.g. cluster reachability, credential validity) and retry."})
}
