// Package models - Error codes for Meshery models
package models

import (
	"fmt"
	"time"

	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
const (
	ErrGrafanaClientCode                  = "1450"
	ErrPageSizeCode                       = "1451"
	ErrPageNumberCode                     = "1452"
	ErrResultIDCode                       = "1453"
	ErrPerfIDCode                         = "1454"
	ErrMarshalCode                        = "1455"
	ErrUnmarshalCode                      = "1456"
	ErrGenerateUUIDCode                   = "1457"
	ErrLocalProviderSupportCode           = "1458"
	ErrGrafanaOrgCode                     = "1459"
	ErrGrafanaBoardsCode                  = "1460"
	ErrGrafanaDashboardCode               = "1461"
	ErrGrafanaDataSourceCode              = "1462"
	ErrNilQueryCode                       = "1463"
	ErrGrafanaDataCode                    = "1464"
	ErrApplicationFileNameCode            = "1465"
	ErrFilterFileNameCode                 = "1466"
	ErrPatternFileNameCode                = "1467"
	ErrMakeDirCode                        = "1468"
	ErrFolderStatCode                     = "1469"
	ErrUserIDCode                         = "1470"
	ErrDBConnectionCode                   = "1471"
	ErrNilConfigDataCode                  = "1472"
	ErrDBOpenCode                         = "1473"
	ErrDBRLockCode                        = "1474"
	ErrDBLockCode                         = "1475"
	ErrDBReadCode                         = "1476"
	ErrDBDeleteCode                       = "1477"
	ErrCopyCode                           = "1478"
	ErrDBPutCode                          = "1479"
	ErrPrometheusGetNodesCode             = "1480"
	ErrPrometheusLabelSeriesCode          = "1481"
	ErrPrometheusQueryRangeCode           = "1482"
	ErrPrometheusStaticBoardCode          = "1483"
	ErrTokenRefreshCode                   = "1484"
	ErrGetTokenCode                       = "1485"
	ErrDataReadCode                       = "1486"
	ErrTokenDecodeCode                    = "1487"
	ErrNilJWKsCode                        = "1488"
	ErrNilKeysCode                        = "1489"
	ErrTokenExpiredCode                   = "1490"
	ErrTokenClaimsCode                    = "1491"
	ErrTokenClientCheckCode               = "1492"
	ErrTokenPraseCode                     = "1493"
	ErrJWKsKeysCode                       = "1494"
	ErrDecodeBase64Code                   = "1495"
	ErrMarshalPKIXCode                    = "1496"
	ErrEncodingPEMCode                    = "1497"
	ErrPraseUnverifiedCode                = "1498"
	ErrEncodingCode                       = "1499"
	ErrFetchCode                          = "1500"
	ErrPostCode                           = "1501"
	ErrDeleteCode                         = "1502"
	ErrInvalidCapabilityCode              = "1503"
	ErrResultDataCode                     = "1504"
	ErrUnableToPersistsResultCode         = "1505"
	ErrValidURLCode                       = "1506"
	ErrTestEndpointCode                   = "1507"
	ErrLoadgeneratorCode                  = "1508"
	ErrProtocolCode                       = "1509"
	ErrTestClientCode                     = "1510"
	ErrParsingTestCode                    = "1511"
	ErrFieldCode                          = "1512"
	ErrFetchDataCode                      = "1513"
	ErrIndexOutOfRangeCode                = "1514"
	ErrSessionCopyCode                    = "1515"
	ErrSavingSeededComponentsCode         = "1516"
	ErrGettingSeededComponentsCode        = "1517"
	ErrDownloadingSeededComponentsCode    = "1518"
	ErrContextIDCode                      = "1519"
	ErrMesheryInstanceIDCode              = "1520"
	ErrMesheryNotInClusterCode            = "1521"
	ErrBrokerNotFoundCode                 = "1522"
	ErrCreateOperatorDeploymentConfigCode = "1523"
	ErrRequestMeshsyncStoreCode           = "1524"
	ErrBrokerSubscriptionCode             = "1525"
	ErrContextAlreadyPersistedCode        = "1526"
	ErrGetPackageCode                     = "1527"
	ErrTokenRevokeCode                    = "1528"
	ErrTokenIntrospectCode                = "1529"
	ErrShareDesignCode                    = "1530"
	ErrUnreachableRemoteProviderCode      = "1531"
	ErrShareFilterCode                    = "1532"
	ErrPersistEventCode                   = "1533"
	ErrUnreachableKubeAPICode             = "1534"
	ErrFlushMeshSyncDataCode              = "1535"
	ErrUpdateConnectionStatusCode         = "1540"
	ErrResultNotFoundCode                 = "1546"
	ErrPersistCredentialCode              = "1547"
	ErrPersistConnectionCode              = "1548"
	ErrPrometheusScanCode                 = "1549"
	ErrGrafanaScanCode                    = "1550"
)

var (
	ErrResultID                = errors.New(ErrResultIDCode, errors.Alert, []string{"Given resultID is not valid"}, []string{"Given resultID is nil"}, []string{}, []string{})
	ErrLocalProviderSupport    = errors.New(ErrLocalProviderSupportCode, errors.Alert, []string{"Method not supported by local provider"}, []string{}, []string{}, []string{})
	ErrNilQuery                = errors.New(ErrNilQueryCode, errors.Alert, []string{"Query data passed is nil"}, []string{}, []string{}, []string{})
	ErrApplicationFileName     = errors.New(ErrApplicationFileNameCode, errors.Alert, []string{"Invalid Applicationfile"}, []string{"Name field is either not present or is not valid"}, []string{}, []string{})
	ErrFilterFileName          = errors.New(ErrFilterFileNameCode, errors.Alert, []string{"Invalid Filterfile"}, []string{"Name field is either not present or is not valid"}, []string{}, []string{})
	ErrPatternFileName         = errors.New(ErrPatternFileNameCode, errors.Alert, []string{"Invalid Patternfile"}, []string{"Name field is either not present or is not valid"}, []string{}, []string{})
	ErrUserID                  = errors.New(ErrUserIDCode, errors.Alert, []string{"User ID is empty"}, []string{}, []string{}, []string{})
	ErrDBConnection            = errors.New(ErrDBConnectionCode, errors.Alert, []string{"Connection to DataBase does not exist"}, []string{}, []string{}, []string{})
	ErrNilConfigData           = errors.New(ErrNilConfigDataCode, errors.Alert, []string{"Given config data is nil"}, []string{}, []string{}, []string{})
	ErrNilJWKs                 = errors.New(ErrNilJWKsCode, errors.Alert, []string{"Invalid JWks"}, []string{"Value of JWKs is nil"}, []string{}, []string{})
	ErrNilKeys                 = errors.New(ErrNilKeysCode, errors.Alert, []string{"Key not found"}, []string{"JWK not found for the given KeyID"}, []string{}, []string{})
	ErrTokenExpired            = errors.New(ErrTokenExpiredCode, errors.Alert, []string{"Token has expired"}, []string{"Token is invalid, it has expired"}, []string{}, []string{})
	ErrTokenClaims             = errors.New(ErrTokenClaimsCode, errors.Alert, []string{"Error occurred while prasing claims"}, []string{}, []string{}, []string{})
	ErrValidURL                = errors.New(ErrValidURLCode, errors.Alert, []string{"Enter valid URLs"}, []string{}, []string{}, []string{})
	ErrTestEndpoint            = errors.New(ErrTestEndpointCode, errors.Alert, []string{"minimum one test endpoint needs to be specified"}, []string{}, []string{}, []string{})
	ErrLoadgenerator           = errors.New(ErrLoadgeneratorCode, errors.Alert, []string{"specify valid Loadgenerator"}, []string{}, []string{}, []string{})
	ErrProtocol                = errors.New(ErrProtocolCode, errors.Alert, []string{"specify the Protocol for all clients"}, []string{}, []string{}, []string{})
	ErrTestClient              = errors.New(ErrTestClientCode, errors.Alert, []string{"minimum one test client needs to be specified"}, []string{}, []string{}, []string{})
	ErrParsingTest             = errors.New(ErrParsingTestCode, errors.Alert, []string{"error parsing test duration, please refer to: https://docs.meshery.io/guides/mesheryctl#performance-management"}, []string{}, []string{}, []string{})
	ErrField                   = errors.New(ErrFieldCode, errors.Alert, []string{"Error: name field is blank"}, []string{}, []string{}, []string{})
	ErrIndexOutOfRange         = errors.New(ErrIndexOutOfRangeCode, errors.Alert, []string{"Error: index out of range"}, []string{}, []string{}, []string{})
	ErrContextID               = errors.New(ErrContextIDCode, errors.Alert, []string{"Error: Context ID is empty"}, []string{}, []string{}, []string{})
	ErrMesheryInstanceID       = errors.New(ErrMesheryInstanceIDCode, errors.Alert, []string{"Error: Meshery Instance ID is empty or is invalid"}, []string{}, []string{}, []string{})
	ErrMesheryNotInCluster     = errors.New(ErrMesheryNotInClusterCode, errors.Alert, []string{"Error: Meshery is not running inside a cluster"}, []string{}, []string{}, []string{})
	ErrContextAlreadyPersisted = errors.New(ErrContextAlreadyPersistedCode, errors.Alert, []string{"kubernetes context already persisted with provider"}, []string{"kubernetes context already persisted with provider"}, []string{}, []string{})
)

func ErrGetPackage(err error) error {
	return errors.New(ErrGetPackageCode, errors.Alert, []string{"Could not get the package"}, []string{"", err.Error()}, []string{""}, []string{"Make sure the configurations are correct"})
}

func ErrBrokerSubscription(err error) error {
	return errors.New(ErrBrokerSubscriptionCode, errors.Alert, []string{"Could not subscribe to the broker subject"}, []string{"", err.Error()}, []string{""}, []string{"Make sure meshery broker is healthy"})
}

func ErrRequestMeshsyncStore(err error) error {
	return errors.New(ErrRequestMeshsyncStoreCode, errors.Alert, []string{"Meshsync store request could not be issued"}, []string{"", err.Error()}, []string{""}, []string{"Make sure meshery broker is healthy"})
}

func ErrCreateOperatorDeploymentConfig(err error) error {
	return errors.New(ErrCreateOperatorDeploymentConfigCode, errors.Alert, []string{"Operator deployment configuration could not be created."}, []string{"", err.Error()}, []string{""}, []string{""})
}

func ErrBrokerNotFound(err error) error {
	return errors.New(ErrBrokerNotFoundCode, errors.Alert, []string{"Meshery broker not found"}, []string{"Unable to find meshery broker in the cluster", err.Error()}, []string{"Invalid Grafana Endpoint or API-Key"}, []string{"Update your Grafana URL and API-Key from the settings page in the UI"})
}

func ErrGrafanaClient(err error) error {
	return errors.New(ErrGrafanaClientCode, errors.Alert, []string{"Unable to initialize Grafana Client"}, []string{"Unable to initializes client for interacting with an instance of Grafana server", err.Error()}, []string{"Invalid Grafana Endpoint or API-Key"}, []string{"Update your Grafana URL and API-Key from the settings page in the UI"})
}

func ErrPageSize(err error) error {
	return errors.New(ErrPageSizeCode, errors.Alert, []string{"Unable to prase the Page Size"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPageNumber(err error) error {
	return errors.New(ErrPageNumberCode, errors.Alert, []string{"Unable to prase the Page Numer"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPerfID(err error) error {
	return errors.New(ErrPerfIDCode, errors.Alert, []string{"Invalid peformance profile ID"}, []string{err.Error()}, []string{}, []string{})
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

func ErrFetch(err error, obj string, statusCode int) error {
	return errors.New(ErrFetchCode, errors.Alert, []string{"Unable to fetch data from the Provider", obj}, []string{"Status Code: " + fmt.Sprint(statusCode), err.Error()}, []string{}, []string{})
}

func ErrPost(err error, obj string, statusCode int) error {
	return errors.New(ErrPostCode, errors.Alert, []string{"Unable to post data to the Provider", obj}, []string{"Status Code: " + fmt.Sprint(statusCode), err.Error()}, []string{}, []string{})
}

func ErrDelete(err error, obj string, statusCode int) error {
	return errors.New(ErrDeleteCode, errors.Alert, []string{"Unable to de-register Meshery Server from Remote Provider", obj}, []string{"Status Code: " + fmt.Sprint(statusCode) + " ", err.Error()}, []string{"Network connectivity to Remote Provider may not be available. Session might have expired; token could be invalid."}, []string{"Verify that the Remote Provider is available. Ensure that you have an active session / valid token."})
}

func ErrDecodeBase64(err error, obj string) error {
	return errors.New(ErrDecodeBase64Code, errors.Alert, []string{"Error occurred while decoding base65 string", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrMarshalPKIX(err error) error {
	return errors.New(ErrMarshalPKIXCode, errors.Alert, []string{"Error occurred while marshaling PKIX"}, []string{err.Error()}, []string{}, []string{})
}

func ErrEncodingPEM(err error) error {
	return errors.New(ErrEncodingPEMCode, errors.Alert, []string{"Error occurred while encoding jwk to pem"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPraseUnverified(err error) error {
	return errors.New(ErrPraseUnverifiedCode, errors.Alert, []string{"Error occurred while prasing tokens (unverified)"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDataRead(err error, r string) error {
	return errors.New(ErrDataReadCode, errors.Alert, []string{"Error occurred while reading from the Reader", r}, []string{err.Error()}, []string{}, []string{})
}

func ErrResultData() error {
	return errors.New(ErrResultDataCode, errors.Alert, []string{"given result data is nil"}, []string{}, []string{}, []string{})
}

func ErrUnableToPersistsResult(err error) error {
	return errors.New(ErrUnableToPersistsResultCode, errors.Alert, []string{"unable to persists the result data"}, []string{err.Error()}, []string{}, []string{})
}

func ErrGenerateUUID(err error) error {
	return errors.New(ErrGenerateUUIDCode, errors.Alert, []string{"Unable to generate a new UUID"}, []string{err.Error()}, []string{}, []string{})
}

func ErrGrafanaOrg(err error) error {
	return errors.New(ErrGrafanaOrgCode, errors.Alert, []string{"Failed to get Org data from Grafana"}, []string{err.Error()}, []string{"Invalid URL", "Invalid API-Key"}, []string{})
}

func ErrGrafanaBoards(err error) error {
	return errors.New(ErrGrafanaBoardsCode, errors.Alert, []string{"Unable to get Grafana Boards"}, []string{err.Error()}, []string{"Grafana endpoint might not be reachable from meshery", "Grafana endpoint is incorrect"}, []string{"Check if your Grafana endpoint is correct", "Connect to Grafana from the settings page in the UI"})
}

func ErrGrafanaDashboard(err error, UID string) error {
	return errors.New(ErrGrafanaDashboardCode, errors.Alert, []string{"Error getting grafana dashboard from UID", UID}, []string{err.Error()}, []string{}, []string{})
}

func ErrGrafanaDataSource(err error, ds string) error {
	return errors.New(ErrGrafanaDataSourceCode, errors.Alert, []string{"Error getting Grafana Board's Datasource", ds}, []string{err.Error()}, []string{}, []string{})
}

func ErrGrafanaData(err error, apiEndpoint string) error {
	return errors.New(ErrGrafanaDataCode, errors.Alert, []string{"Error getting data from Grafana API", apiEndpoint}, []string{err.Error()}, []string{}, []string{})
}

func ErrMakeDir(err error, dir string) error {
	return errors.New(ErrMakeDirCode, errors.Alert, []string{"Unable to create directory/folder", dir}, []string{err.Error()}, []string{}, []string{})
}

func ErrFolderStat(err error, dir string) error {
	return errors.New(ErrFolderStatCode, errors.Alert, []string{"Unable to find (os.stat) the folder", dir}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBOpen(err error) error {
	return errors.New(ErrDBOpenCode, errors.Alert, []string{"Unable to open the database"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBRLock(err error) error {
	return errors.New(ErrDBRLockCode, errors.Alert, []string{"Unable to obtain read lock from bitcask store"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBLock(err error) error {
	return errors.New(ErrDBLockCode, errors.Alert, []string{"Unable to obtain write lock from bitcask store"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBRead(err error) error {
	return errors.New(ErrDBReadCode, errors.Alert, []string{"Unable to read data from bitcast store"}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBPut(err error) error {
	return errors.New(ErrDBPutCode, errors.Alert, []string{"Unable to Persist config data."}, []string{err.Error()}, []string{}, []string{})
}

func ErrDBDelete(err error, user string) error {
	return errors.New(ErrDBDeleteCode, errors.Alert, []string{"Unable to delete config data for the user", user}, []string{err.Error()}, []string{}, []string{})
}

func ErrCopy(err error, obj string) error {
	return errors.New(ErrCopyCode, errors.Alert, []string{"Error occurred while copying", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrPrometheusGetNodes(err error) error {
	return errors.New(ErrPrometheusGetNodesCode, errors.Alert, []string{"Prometheus Client unable to get all nodes"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPrometheusStaticBoard(err error) error {
	return errors.New(ErrPrometheusStaticBoardCode, errors.Alert, []string{"Unbale to get Static Boards"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPrometheusLabelSeries(err error) error {
	return errors.New(ErrPrometheusLabelSeriesCode, errors.Alert, []string{"Unable to get the label set series"}, []string{err.Error()}, []string{}, []string{})
}

func ErrPrometheusQueryRange(err error, query string, startTime, endTime time.Time, step time.Duration) error {
	return errors.New(ErrPrometheusQueryRangeCode, errors.Alert, []string{"Unable to fetch data for the query", fmt.Sprintf("Query: %s, with start: %v, end: %v, step: %v", query, startTime, endTime, step)}, []string{err.Error()}, []string{}, []string{})
}

func ErrTokenRefresh(err error) error {
	return errors.New(ErrTokenRefreshCode, errors.Alert, []string{"Error occurred while Refresing the token"}, []string{err.Error()}, []string{}, []string{})
}

func ErrTokenRevoke(err error) error {
	return errors.New(ErrTokenRevokeCode, errors.Alert, []string{"Error occurred while revoking the token"}, []string{err.Error()}, []string{"Unable to revoke token. Token appears to be a malformed base64 token."}, []string{"Try logging out (again) in order to fully close your session (and revoke the session token)."})
}

func ErrTokenIntrospect(err error) error {
	return errors.New(ErrTokenIntrospectCode, errors.Alert, []string{"token introspection failed"}, []string{err.Error()}, []string{"Invalid session token. Token has revoked."}, []string{"Login again to establish a new session with valid token."})
}

func ErrGetToken(err error) error {
	return errors.New(ErrGetTokenCode, errors.Alert, []string{"Error occurred while getting token from the Browser Cookie"}, []string{err.Error()}, []string{}, []string{})
}

func ErrTokenDecode(err error) error {
	return errors.New(ErrTokenDecodeCode, errors.Alert, []string{"Error occurred while Decoding Token Data"}, []string{err.Error()}, []string{}, []string{})
}

func ErrTokenClientCheck(err error) error {
	return errors.New(ErrTokenClientCheckCode, errors.Alert, []string{"Error occurred while performing token check HTTP request"}, []string{err.Error()}, []string{}, []string{})
}

func ErrTokenPrase(err error) error {
	return errors.New(ErrTokenPraseCode, errors.Alert, []string{"Error occurred while Prasing and validating the token"}, []string{err.Error()}, []string{}, []string{})
}

func ErrJWKsKeys(err error) error {
	return errors.New(ErrJWKsKeysCode, errors.Alert, []string{"Unable to fetch JWKs keys from the remote provider"}, []string{err.Error()}, []string{}, []string{})
}

func ErrInvalidCapability(capability string, provider string) error {
	return errors.New(ErrInvalidCapabilityCode, errors.Alert, []string{"Capablity is not supported by your Provider", capability}, []string{"You dont have access to the capability", "Provider: " + provider, "Capability: " + capability}, []string{"Not logged in to the vaild remote Provider"}, []string{"Connect to the vaild remote Provider", "Ask the Provider Adim for access"})
}

func ErrFetchData(err error) error {
	return errors.New(ErrFetchDataCode, errors.Alert, []string{"unable to fetch result data"}, []string{err.Error()}, []string{}, []string{})
}

func ErrSessionCopy(err error) error {
	return errors.New(ErrSessionCopyCode, errors.Alert, []string{"Error: session copy error"}, []string{err.Error()}, []string{}, []string{})
}

func ErrGettingSeededComponents(err error, content string) error {
	return errors.New(ErrGettingSeededComponentsCode, errors.Alert, []string{"Error while getting ", content, " from sample content"}, []string{err.Error()}, []string{"Sample content does not exist.\nContent file format not supported.\nUser doesn't have permission to read sample content.\nContent file corrupt."}, []string{"Try restarting Meshery.\nTry fetching content again."})
}

func ErrSavingSeededComponents(err error, content string) error {
	return errors.New(ErrSavingSeededComponentsCode, errors.Alert, []string{"Error while saving sample ", content}, []string{err.Error()}, []string{"User doesn't have permission to save content.\nDatabase unreachable.\nDatabase locked or corrupt.\nContent unsupported."}, []string{"Retry fetching content\nRetry sigining in\nLogin with correct user account\nRetry after deleting '~/.meshery/config'."})
}

func ErrDownloadingSeededComponents(err error, content string) error {
	return errors.New(ErrDownloadingSeededComponentsCode, errors.Alert, []string{"Could not download seed content for" + content}, []string{err.Error()}, []string{"The content is not present at the specified url endpoint", "HTTP requests failed"}, []string{"Make sure the content is available at the endpoints", "Make sure that Github is reachable and the http requests are not failing"})
}

func ErrShareDesign(err error) error {
	return errors.New(ErrShareDesignCode, errors.Alert, []string{"cannot make design public"}, []string{err.Error()}, []string{"email address provided might not be valid", "insufficient permission"}, []string{"Ensure that you are the owner of the design you are sharing", "Try again later", "Try using an alternate email address"})
}

func ErrShareFilter(err error) error {
	return errors.New(ErrShareFilterCode, errors.Alert, []string{"Cannot make filter public"}, []string{err.Error()}, []string{"Email address provided might not be valid", "Verify that you have sufficient permission to share the filter. You should be the owner of the filter"}, []string{"Verify the spelling of the email address. Try using an alternate email address.", "Ensure that you are the owner of the filter you are sharing or have sharing permission assigned.", ""})
}

func ErrUnreachableRemoteProvider(err error) error {
	return errors.New(ErrUnreachableRemoteProviderCode, errors.Alert, []string{"Could not reach remote provider"}, []string{"", err.Error()}, []string{"Remote provider server may be down or not accepting requests."}, []string{"Make sure remote provider server is healthy and accepting requests."})
}

func ErrPersistEvent(err error) error {
	return errors.New(ErrPersistEventCode, errors.Alert, []string{"Could not persist event"}, []string{err.Error()}, []string{"Database could be down or not reachable", "Meshery Database handler is not accessible to perform operations"}, []string{"Restart Meshery Server or Perform Hard Reset"})
}

func ErrUnreachableKubeAPI(err error, server string) error {
	return errors.New(ErrUnreachableKubeAPICode, errors.Alert, []string{fmt.Sprintf("Error communicating with KubeAPI at %s.", server)}, []string{err.Error()}, []string{"The Kubernetes API server is not reachable.", "Credentials are invalid."}, []string{"Verify network connectivity and Kubernetes API responsiveness between Meshery Server and your cluster.", "Ensure client credential is not expired and is properly formed.", "Remove the cluster credential and enable 'insecure-skip-tls-verify'."})
}

func ErrFlushMeshSyncData(err error, contextName, server string) error {
	return errors.New(ErrFlushMeshSyncDataCode, errors.Alert, []string{"Unable to flush MeshSync data for context %s at %s "}, []string{err.Error()}, []string{"Meshery Database handler is not accessible to perform operations"}, []string{"Restart Meshery Server or Perform Hard Reset"})
}

func ErrUpdateConnectionStatus(err error, statusCode int) error {
	return errors.New(ErrUpdateConnectionStatusCode, errors.Alert, []string{"Unable to update connection status"}, []string{err.Error()}, []string{"Connection was already deleted", "User might not have necessary privileges"}, []string{"Try refresing, you might be seeing stale data on the dashboard", "Check if the user has necessary privileges"})
}

func ErrResultNotFound(err error) error {
	return errors.New(ErrResultNotFoundCode, errors.Alert, []string{err.Error()}, []string{"The record in the database does not exist."}, []string{"The record might have been deleted."}, []string{""})
}

func ErrPersistCredential(err error) error {
	return errors.New(ErrPersistCredentialCode, errors.Alert, []string{"unable to persist credential details"}, []string{err.Error()}, []string{"The credential object is not valid"}, []string{"Ensure all the required fields are provided"})
}

func ErrPersistConnection(err error) error {
	return errors.New(ErrPersistConnectionCode, errors.Alert, []string{"unable to persist connection details"}, []string{err.Error()}, []string{"The connection object is not valid"}, []string{"Ensure all the required fields are provided"})
}

func ErrGrafanaScan(err error) error {
	return errors.New(ErrGrafanaScanCode, errors.Alert, []string{"Unable to connect to grafana"}, []string{err.Error()}, []string{"Grafana endpoint might not be reachable from meshery", "Grafana endpoint is incorrect"}, []string{"Check if your Grafana Endpoint is correct", "Connect to Grafana from the settings page in the UI"})
}

func ErrPrometheusScan(err error) error {
	return errors.New(ErrPrometheusScanCode, errors.Alert, []string{"Unable to connect to prometheus"}, []string{err.Error()}, []string{"Prometheus endpoint might not be reachable from meshery", "Prometheus endpoint is incorrect"}, []string{"Check if your Prometheus and Grafana Endpoint are correct", "Connect to Prometheus and Grafana from the settings page in the UI"})
}
