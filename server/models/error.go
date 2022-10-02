// Package models - Error codes for Meshery models
package models

import (
	"fmt"
	"time"

	"github.com/layer5io/meshkit/errors"
)

const (
	ErrGrafanaClientCode                  = "2073"
	ErrPageSizeCode                       = "2074"
	ErrPageNumberCode                     = "2075"
	ErrResultIDCode                       = "2076"
	ErrPerfIDCode                         = "2077"
	ErrMarshalCode                        = "2078"
	ErrUnmarshalCode                      = "2079"
	ErrGenerateUUIDCode                   = "2080"
	ErrLocalProviderSupportCode           = "2081"
	ErrGrafanaOrgCode                     = "2082"
	ErrGrafanaBoardsCode                  = "2083"
	ErrGrafanaDashboardCode               = "2084"
	ErrGrafanaDataSourceCode              = "2085"
	ErrNilQueryCode                       = "2086"
	ErrGrafanaDataCode                    = "2087"
	ErrApplicationFileNameCode            = "2088"
	ErrFilterFileNameCode                 = "2089"
	ErrPatternFileNameCode                = "2090"
	ErrMakeDirCode                        = "2091"
	ErrFolderStatCode                     = "2092"
	ErrUserIDCode                         = "2093"
	ErrDBConnectionCode                   = "2094"
	ErrNilConfigDataCode                  = "2095"
	ErrDBOpenCode                         = "2096"
	ErrDBRLockCode                        = "2097"
	ErrDBLockCode                         = "2098"
	ErrDBReadCode                         = "2099"
	ErrDBDeleteCode                       = "2100"
	ErrCopyCode                           = "2101"
	ErrDBPutCode                          = "2102"
	ErrPrometheusGetNodesCode             = "2103"
	ErrPrometheusLabelSeriesCode          = "2104"
	ErrPrometheusQueryRangeCode           = "2105"
	ErrPrometheusStaticBoardCode          = "2106"
	ErrTokenRefreshCode                   = "2107"
	ErrGetTokenCode                       = "2108"
	ErrDataReadCode                       = "2109"
	ErrTokenDecodeCode                    = "2110"
	ErrNilJWKsCode                        = "2111"
	ErrNilKeysCode                        = "2112"
	ErrTokenExpiredCode                   = "2113"
	ErrTokenClaimsCode                    = "2114"
	ErrTokenClientCheckCode               = "2115"
	ErrTokenPraseCode                     = "2116"
	ErrJWKsKeysCode                       = "2117"
	ErrDecodeBase64Code                   = "2118"
	ErrMarshalPKIXCode                    = "2119"
	ErrEncodingPEMCode                    = "2120"
	ErrPraseUnverifiedCode                = "2121"
	ErrEncodingCode                       = "2122"
	ErrFetchCode                          = "2123"
	ErrPostCode                           = "2124"
	ErrDeleteCode                         = "2125"
	ErrInvalidCapabilityCode              = "2126"
	ErrResultDataCode                     = "2127"
	ErrUnableToPersistsResultCode         = "2128"
	ErrValidURLCode                       = "2129"
	ErrTestEndpointCode                   = "2130"
	ErrLoadgeneratorCode                  = "2131"
	ErrProtocolCode                       = "2132"
	ErrTestClientCode                     = "2133"
	ErrParsingTestCode                    = "2134"
	ErrFieldCode                          = "2135"
	ErrFetchDataCode                      = "2147"
	ErrIndexOutOfRangeCode                = "2148"
	ErrSessionCopyCode                    = "2149"
	ErrSavingSeededComponentsCode         = "2215"
	ErrGettingSeededComponentsCode        = "2216"
	ErrDownloadingSeededComponentsCode    = "2217"
	ErrContextIDCode                      = "2218"
	ErrMesheryInstanceIDCode              = "2219"
	ErrMesheryNotInClusterCode            = "2220"
	ErrBrokerNotFoundCode                 = "2235"
	ErrCreateOperatorDeploymentConfigCode = "2236"
	ErrRequestMeshsyncStoreCode           = "2237"
	ErrBrokerSubscriptionCode             = "2238"
	ErrContextAlreadyPersistedCode        = "2241"
	ErrGetPackageCode                     = "2252"
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
	return errors.New(ErrMarshalCode, errors.Alert, []string{"Unable to marshal the object", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrUnmarshal(err error, obj string) error {
	return errors.New(ErrUnmarshalCode, errors.Alert, []string{"Unable to unmarshal the object", obj}, []string{err.Error()}, []string{}, []string{})
}

func ErrEncoding(err error, obj string) error {
	return errors.New(ErrEncodingCode, errors.Alert, []string{"Error encoding the object", obj}, []string{err.Error()}, []string{"Unable to encode the : ", obj}, []string{})
}

func ErrFetch(err error, obj string, statusCode int) error {
	return errors.New(ErrFetchCode, errors.Alert, []string{"Unable to fetch data from the Provider", obj}, []string{"Status Code: " + fmt.Sprint(statusCode), err.Error()}, []string{}, []string{})
}

func ErrPost(err error, obj string, statusCode int) error {
	return errors.New(ErrPostCode, errors.Alert, []string{"Unable to post data to the Provider", obj}, []string{"Status Code: " + fmt.Sprint(statusCode), err.Error()}, []string{}, []string{})
}

func ErrDelete(err error, obj string, statusCode int) error {
	return errors.New(ErrDeleteCode, errors.Alert, []string{"Unable to delete data from the Provider", obj}, []string{"Status Code: " + fmt.Sprint(statusCode), err.Error()}, []string{}, []string{})
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
	return errors.New(ErrDataReadCode, errors.Alert, []string{"Eeror occurred while reading from the Reader", r}, []string{err.Error()}, []string{}, []string{})
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
	return errors.New(ErrGrafanaBoardsCode, errors.Alert, []string{"Unable to get Grafana Boards"}, []string{err.Error()}, []string{}, []string{})
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
