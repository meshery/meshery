package handlers

import (
	"fmt"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrHealthCheckFailedCode    = "replace me"
	ErrInvalidK8SConfigCode     = "replace me"
	ErrInvalidAdapterCode       = "replace me"
	ErrNilClientCode            = "replace me"
	ErrPrometheusScanCode       = "replace me"
	ErrGrafanaScanCode          = "replace me"
	ErrRecordPreferencesCode    = "replace me"
	ErrGrafanaConfigCode        = "replace me"
	ErrPrometheusConfigCode     = "replace me"
	ErrGrafanaQueryCode         = "replace me"
	ErrPrometheusQueryCode      = "replace me"
	ErrGrafanaBoardsCode        = "replace me"
	ErrPrometheusBoardsCode     = "replace me"
	ErrStaticBoardsCode         = "replace me"
	ErrRequestBodyCode          = "replace me"
	ErrMarshalCode              = "replace me"
	ErrUnmarshalCode            = "replace me"
	ErrEncodingCode             = "replace me"
	ErrParseBoolCode            = "replace me"
	ErrStreamEventsCode         = "replace me"
	ErrStreamClientCode         = "replace me"
	ErrUnmarshalEventCode       = "replace me"
	ErrPublishSmiResultsCode    = "replace me"
	ErrMarshalEventCode         = "replace me"
	ErrPluginOpenCode           = "replace me"
	ErrPluginLookupCode         = "replace me"
	ErrPluginRunCode            = "replace me"
	ErrParseFormCode            = "replace me"
	ErrQueryGetCode             = "replace me"
	ErrGetResultCode            = "replace me"
	ErrConvertToSpecCode        = "replace me"
	ErrFetchSMIResultsCode      = "replace me"
	ErrFormFileCode             = "replace me"
	ErrDownloadFileCode         = "replace me"
	ErrReadConfigCode           = "replace me"
	ErrLoadConfigCode           = "replace me"
	ErrOpenFileCode             = "replace me"
	ErrKubeVersionCode          = "replace me"
	ErrAddAdapterCode           = "replace me"
	ErrRetrieveDataCode         = "replace me"
	ErrValidAdapterCode         = "replace me"
	ErrOperationIDCode          = "replace me"
	ErrMeshClientCode           = "replace me"
	ErrApplyChangeCode          = "replace me"
	ErrRetrieveMeshDataCode     = "replace me"
	ErrApplicationFailureCode   = "replace me"
	ErrDecodingCode             = "replace me"
	ErrRetrieveUserTokenCode    = "replace me"
	ErrFailToSaveCode           = "replace me"
	ErrFailToDeleteCode         = "replace me"
	ErrFailToLoadExtensionsCode = "replace me"
	ErrConversionCode           = "replace me"
	ErrParseDurationCode        = "replace me"
	ErrLoadTestCode             = "replace me"
	ErrFetchKubernetesCode      = "replace me"
	ErrPanicRecoveryCode        = "replace me"
	ErrBlankNameCode            = "replace me"
	ErrInvalidLTURLCode         = "replace me"
	ErrSaveSessionCode          = "replace me"
)

var (
	ErrInvalidK8SConfig = errors.New(ErrInvalidK8SConfigCode, errors.Alert, []string{"No valid kubernetes config found"}, []string{"Kubernetes config is not initialized with Meshery"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrNilClient        = errors.New(ErrNilClientCode, errors.Alert, []string{"Kubernetes client not initialized"}, []string{"Kubernetes config is not initialized with Meshery"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrPrometheusConfig = errors.New(ErrPrometheusConfigCode, errors.Alert, []string{"Prometheus endpoint not configured"}, []string{"Cannot find valid Prometheus endpoint in user pref"}, []string{}, []string{"Setup your Prometheus Endpoint via the settings dashboard"})
	ErrGrafanaConfig    = errors.New(ErrGrafanaConfigCode, errors.Alert, []string{"Grafana endpoint not configured"}, []string{"Cannot find valid grafana endpoint in user pref"}, []string{}, []string{"Setup your Grafana Endpoint via the settings dashboard"})
	ErrStaticBoards     = errors.New(ErrStaticBoardsCode, errors.Alert, []string{"unable to get static board"}, []string{"unable to get static board"}, []string{}, []string{})
	ErrValidAdapter     = errors.New(ErrValidAdapterCode, errors.Alert, []string{"Unable to find valid Adapter URL"}, []string{"unable to find a valid adapter for the given adapter URL"}, []string{"Given adapter URL is not valid"}, []string{"Please provide a valid Adapter URL"})
	ErrAddAdapter       = errors.New(ErrAddAdapterCode, errors.Alert, []string{"meshLocationURL is empty"}, []string{"meshLocationURL is empty to add an adapter"}, []string{"meshLocationURL cannot be empty to add an adapter"}, []string{"please provide the meshLocationURL"})
	ErrMeshClient       = errors.New(ErrMeshClientCode, errors.Alert, []string{"Error creating a mesh client", "Error pinging the mesh adapter"}, []string{"Unable to create a mesh client", "Unable to ping the mesh adapter"}, []string{"Adapter could not be pinged"}, []string{"Unable to connect to the Mesh adapter using the given config, please try again"})
	ErrParseDuration    = errors.New(ErrParseDurationCode, errors.Alert, []string{"error parsing test duration"}, []string{}, []string{}, []string{"please refer to:  https://docs.meshery.io/guides/mesheryctl#performance-management"})
)

func ErrHealthCheckFailed(err error) error {
	return errors.New(ErrHealthCheckFailedCode, errors.Alert, []string{"Health checks failed"}, []string{err.Error()}, []string{}, []string{})
}

func ErrInvalidAdapter(err error, obj string) error {
	return errors.New(ErrInvalidAdapterCode, errors.Alert, []string{"Invalid adapter ", obj, " specified"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPrometheusScan(err error) error {
	return errors.New(ErrPrometheusScanCode, errors.Alert, []string{"Unable to connect to prometheus"}, []string{err.Error()}, []string{}, []string{"Check if your Prometheus and Grafana Endpoint are correct", "Connect to Prometheus and Grafana from the settings page in the UI"})
}

func ErrGrafanaScan(err error) error {
	return errors.New(ErrGrafanaScanCode, errors.Alert, []string{"Unable to connect to grafana"}, []string{err.Error()}, []string{}, []string{"Check if your Grafana Endpoint is correct", "Connect to Grafana from the settings page in the UI"})
}

func ErrPrometheusQuery(err error) error {
	return errors.New(ErrPrometheusQueryCode, errors.Alert, []string{"Unable to query prometheus"}, []string{err.Error()}, []string{}, []string{"Check if your Prometheus query is correct", "Connect to Prometheus and Grafana from the settings page in the UI"})
}

func ErrGrafanaQuery(err error) error {
	return errors.New(ErrGrafanaQueryCode, errors.Alert, []string{"Unable to query grafana"}, []string{err.Error()}, []string{}, []string{"Check if your Grafana query is correct", "Connect to Grafana from the settings page in the UI"})
}

func ErrGrafanaBoards(err error) error {
	return errors.New(ErrGrafanaBoardsCode, errors.Alert, []string{"unable to get grafana boards"}, []string{err.Error()}, []string{}, []string{"Check if your Grafana endpoint is correct", "Connect to Grafana from the settings page in the UI"})
}

func ErrPrometheusBoards(err error) error {
	return errors.New(ErrGrafanaBoardsCode, errors.Alert, []string{"unable to get Prometheus boards"}, []string{err.Error()}, []string{}, []string{"Check if your Prometheus endpoint is correct", "Connect to Prometheus from the settings page in the UI"})
}

func ErrRecordPreferences(err error) error {
	return errors.New(ErrRecordPreferencesCode, errors.Alert, []string{"unable to save user config data"}, []string{err.Error()}, []string{"User token might be invalid", "db might be corrupted"}, []string{"Relogin to Meshery"})
}

func ErrRequestBody(err error) error {
	return errors.New(ErrRequestBodyCode, errors.Alert, []string{"unable to read the request body"}, []string{err.Error()}, []string{}, []string{})
}

func ErrMarshal(err error, obj string) error {
	return errors.New(ErrMarshalCode, errors.Alert, []string{"Unable to marshal the : ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrUnmarshal(err error, obj string) error {
	return errors.New(ErrUnmarshalCode, errors.Alert, []string{"Unable to unmarshal the : ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrEncoding(err error, obj string) error {
	return errors.New(ErrEncodingCode, errors.Alert, []string{"Error encoding the : ", obj}, []string{err.Error()}, []string{"Unable to encode the : ", obj}, []string{})
}

func ErrParseBool(err error, obj string) error {
	return errors.New(ErrParseBoolCode, errors.Alert, []string{"unable to parse : ", obj}, []string{err.Error()}, []string{"Failed due to invalid value of : ", obj}, []string{"please provide a valid value for : ", obj})
}

func ErrStreamEvents(err error) error {
	return errors.New(ErrStreamEventsCode, errors.Alert, []string{"There was an error connecting to the backend to get events"}, []string{err.Error()}, []string{}, []string{})
}

func ErrStreamClient(err error) error {
	return errors.New(ErrStreamClientCode, errors.Alert, []string{"Event streaming ended"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPublishSmiResults(err error) error {
	return errors.New(ErrPublishSmiResultsCode, errors.Alert, []string{"Error publishing SMI results"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPluginOpen(err error) error {
	return errors.New(ErrPluginOpenCode, errors.Alert, []string{"Error opening the plugin"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPluginLookup(err error) error {
	return errors.New(ErrPluginLookupCode, errors.Alert, []string{"Error performing a plugin lookup"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPluginRun(err error) error {
	return errors.New(ErrPluginRunCode, errors.Alert, []string{"Error running meshery plugin"}, []string{err.Error()}, []string{}, []string{})
}

func ErrParseForm(err error) error {
	return errors.New(ErrParseFormCode, errors.Alert, []string{"unable to parse form"}, []string{err.Error()}, []string{}, []string{})
}

func ErrQueryGet(obj string) error {
	return errors.New(ErrQueryGetCode, errors.Alert, []string{"unable to get: ", obj}, []string{}, []string{}, []string{})
}

func ErrGetResult(err error) error {
	return errors.New(ErrGetResultCode, errors.Alert, []string{"unable to get result"}, []string{err.Error()}, []string{}, []string{})
}

func ErrConvertToSpec(err error) error {
	return errors.New(ErrConvertToSpecCode, errors.Alert, []string{"unable to convert to spec"}, []string{err.Error()}, []string{}, []string{})
}

func ErrFetchSMIResults(err error) error {
	return errors.New(ErrFetchSMIResultsCode, errors.Alert, []string{"unable to fetch SMI results"}, []string{err.Error()}, []string{}, []string{})
}

func ErrFormFile(err error) error {
	return errors.New(ErrFormFileCode, errors.Alert, []string{"error getting k8s file"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDownloadFile(err error, obj string) error {
	return errors.New(ErrDownloadFileCode, errors.Alert, []string{"Error downloading file ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrReadConfig(err error) error {
	return errors.New(ErrReadConfigCode, errors.Alert, []string{"error reading config"}, []string{err.Error()}, []string{}, []string{})
}

func ErrLoadConfig(err error) error {
	return errors.New(ErrLoadConfigCode, errors.Alert, []string{"unable to load kubernetes config"}, []string{err.Error()}, []string{}, []string{})
}

func ErrOpenFile(file string) error {
	return errors.New(ErrOpenFileCode, errors.Alert, []string{"unable to open file: ", file}, []string{}, []string{}, []string{})
}

func ErrKubeVersion(err error) error {
	return errors.New(ErrKubeVersionCode, errors.Alert, []string{"unable to get kubernetes version"}, []string{err.Error()}, []string{}, []string{})
}

func ErrRetrieveData(err error) error {
	return errors.New(ErrRetrieveDataCode, errors.Alert, []string{"Unable to retrieve the requested data"}, []string{err.Error()}, []string{}, []string{})
}

func ErrOperationID(err error) error {
	return errors.New(ErrOperationIDCode, errors.Alert, []string{"Error generating the operation Id"}, []string{err.Error()}, []string{}, []string{})
}

func ErrApplyChange(err error) error {
	return errors.New(ErrApplyChangeCode, errors.Alert, []string{"Error applying the change"}, []string{err.Error()}, []string{}, []string{})
}

func ErrRetrieveMeshData(err error) error {
	return errors.New(ErrRetrieveMeshDataCode, errors.Alert, []string{"Error getting operations for the mesh", "Error getting service mesh name"}, []string{err.Error()}, []string{"unable to retrieve the requested data"}, []string{})
}

func ErrApplicationFailure(err error, obj string) error {
	return errors.New(ErrApplicationFailureCode, errors.Alert, []string{"failed to ", obj, "the application"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDecoding(err error, obj string) error {
	return errors.New(ErrDecodingCode, errors.Alert, []string{"Error decoding the : ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrRetrieveUserToken(err error) error {
	return errors.New(ErrRetrieveUserTokenCode, errors.Alert, []string{"Failed to get the user token"}, []string{err.Error()}, []string{}, []string{})
}
func ErrFailToSave(err error, obj string) error {
	return errors.New(ErrFailToSaveCode, errors.Alert, []string{"Failed to Save: ", obj}, []string{err.Error()}, []string{}, []string{})
}
func ErrFailToDelete(err error, obj string) error {
	return errors.New(ErrFailToDeleteCode, errors.Alert, []string{"Failed to Delete: ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrBlankName(err error) error {
	return errors.New(ErrBlankNameCode, errors.Alert, []string{"Error: name field is blank"}, []string{err.Error()}, []string{}, []string{"Provide a name for the test"})
}

func ErrConversion(err error) error {
	return errors.New(ErrConversionCode, errors.Alert, []string{"unable to convert YAML to JSON"}, []string{err.Error()}, []string{}, []string{})
}

func ErrLoadTest(err error, obj string) error {
	return errors.New(ErrLoadTestCode, errors.Alert, []string{"load test error: ", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrFetchKubernetes(err error) error {
	return errors.New(ErrLoadTestCode, errors.Alert, []string{"unable to ping kubernetes", "unable to scan"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPanicRecovery(r interface{}) error {
	return errors.New(ErrPanicRecoveryCode, errors.Alert, []string{"Recovered from panic"}, []string{fmt.Sprint(r)}, []string{}, []string{})
}

func ErrFailToLoadExtensions(err error) error {
	return errors.New(ErrFailToLoadExtensionsCode, errors.Alert, []string{"Failed to Load Extensions from Package"}, []string{err.Error()}, []string{}, []string{})
}

func ErrInvalidLTURL(url string) error {
	return errors.New(ErrInvalidLTURLCode, errors.Alert, []string{"invalid loadtest url: ", url}, []string{}, []string{}, []string{"please refer to:  https://docs.meshery.io/guides/mesheryctl#performance-management"})
}

func ErrSaveSession(err error) error {
	return errors.New(ErrSaveSessionCode, errors.Alert, []string{"unable to save session"}, []string{err.Error()}, []string{}, []string{})
}
