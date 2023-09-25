package utils

import (
	"fmt"
	"strconv"

	"github.com/layer5io/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
var (
	ErrFailRequestCode        = "1163"
	ErrInvalidTokenCode       = "1164"
	ErrFailReqStatusCode      = "1165"
	ErrAttachAuthTokenCode    = "1166"
	ErrUnmarshalCode          = "1167"
	ErrFileReadCode           = "1168"
	ErrCreatingRequestCode    = "1169"
	ErrMarshalCode            = "1170"
	ErrReadResponseBodyCode   = "1171"
	ErrParsingUrlCode         = "1172"
	ErrNotFoundCode           = "1173"
	ErrUnauthenticatedCode    = "1174"
	ErrInvalidFileCode        = "1175"
	ErrInvalidNameOrIDCode    = "1176"
	ErrInvalidAPIResponseCode = "1177"
	ErrReadConfigFileCode     = "1178"
	ErrMarshalIndentCode      = "1179"
	ErrLoadConfigCode         = "1180"
	ErrResponseStatusBodyCode = "1181"
	ErrResponseStatusCode     = "1182"
	ErrJSONToYAMLCode         = "1183"
	ErrOutFormatFlagCode      = "1184"
	ErrInvalidAPICallCode     = "1185"
	ErrParseGithubFileCode    = "1186"
	ErrReadTokenCode          = "1187"
	ErrRequestResponseCode    = "1188"
)

// RootError returns a formatted error message with a link to 'root' command usage page at
// in addition to the error message
func RootError(msg string) string {
	return formatError(msg, cmdRoot)
}

// PerfError returns a formatted error message with a link to 'perf' command usage page at
// in addition to the error message
func PerfError(msg string) string {
	return formatError(msg, cmdPerf)
}

// SystemError returns a formatted error message with a link to 'system' command usage page
// in addition to the error message
func SystemError(msg string) string {
	return formatError(msg, cmdSystem)
}

// SystemTokenError returns a formatted error message with a link to 'token' command usage page
// in addition to the error message
func SystemTokenError(msg string) string {
	return formatError(msg, cmdToken)
}

func SystemLifeCycleError(msg string, cmd string) string {
	switch cmd {
	case "stop":
		return formatError(msg, cmdSystemStop)
	case "update":
		return formatError(msg, cmdSystemUpdate)
	case "reset":
		return formatError(msg, cmdSystemReset)
	case "status":
		return formatError(msg, cmdSystemStatus)
	case "restart":
		return formatError(msg, cmdSystemRestart)
	default:
		return formatError(msg, cmdSystem)
	}
}

// SystemContextSubError returns a formatted error message with a link to `context` command usage page
// in addition to the error message
func SystemContextSubError(msg string, cmd string) string {
	switch cmd {
	case "delete":
		return formatError(msg, cmdContextDelete)
	case "create":
		return formatError(msg, cmdContextCreate)
	case "view":
		return formatError(msg, cmdContextView)
	default:
		return formatError(msg, cmdContext)
	}
}

// SystemChannelSubError returns a formatted error message with a link to `channel` command usage page
// in addition to the error message
func SystemChannelSubError(msg string, cmd string) string {
	switch cmd {
	case "switch":
		return formatError(msg, cmdChannelSwitch)
	case "view":
		return formatError(msg, cmdChannelView)
	case "set":
		return formatError(msg, cmdChannelSet)
	default:
		return formatError(msg, cmdChannel)
	}
}

// SystemProviderSubError returns a formatted error message with a link to `provider` command usage page
// in addition to the error message
func SystemProviderSubError(msg string, cmd string) string {
	switch cmd {
	case "switch":
		return formatError(msg, cmdProviderSwitch)
	case "view":
		return formatError(msg, cmdProviderView)
	case "set":
		return formatError(msg, cmdProviderSet)
	case "list":
		return formatError(msg, cmdProviderList)
	case "reset":
		return formatError(msg, cmdProviderReset)
	default:
		return formatError(msg, cmdProvider)
	}
}

// SystemProviderSubError returns a formatted error message with a link to `provider` command usage page
// in addition to the error message
func SystemModelSubError(msg string, cmd string) string {
	switch cmd {
	case "list":
		return formatError(msg, cmdModelList)
	case "view":
		return formatError(msg, cmdModelView)
	default:
		return formatError(msg, cmdModel)
	}
}

// MeshError returns a formatted error message with a link to 'mesh' command usage page in addition to the error message
func MeshError(msg string) string {
	return formatError(msg, cmdMesh)
}

// ExpError returns a formatted error message with a link to 'exp' command usage page in addition to the error message
func ExpError(msg string) string {
	return formatError(msg, cmdExp)
}

// FilterError returns a formatted error message with a link to 'filter' command usage page in addition to the error message
func FilterError(msg string) string {
	return formatError(msg, cmdFilter)
}

// FilterImportError returns a formatted error message with a link to 'filter import' command usage page in addition to the error message
func FilterImportError(msg string) string {
	return formatError(msg, cmdFilterImport)
}

// FilterDeleteError returns a formatted error message with a link to 'filter delete' command usage page in addition to the error message
func FilterDeleteError(msg string) string {
	return formatError(msg, cmdFilterDelete)
}

// FilterListError returns a formatted error message with a link to 'filter list' command usage page in addition to the error message
func FilterListError(msg string) string {
	return formatError(msg, cmdFilterList)
}

// FilterViewError returns a formatted error message with a link to 'filter view' command usage page in addition to the error message
func FilterViewError(msg string) string {
	return formatError(msg, cmdFilterView)
}

// PatternError returns a formatted error message with a link to 'pattern' command usage page in addition to the error message
func PatternError(msg string) string {
	return formatError(msg, cmdPattern)
}

// PatternViewError returns a formatted error message with a link to the 'pattern view' commad usage page in addition to the error message
func PatternViewError(msg string) string {
	return formatError(msg, cmdPatternView)
}

// AppError returns a formatted error message with a link to 'app' command usage page in addition to the error message
func AppError(msg string) string {
	return formatError(msg, cmdApp)
}

// AppError returns a formatted error message with a link to 'app view' command usage page in addition to the error message
func AppViewError(msg string) string {
	return formatError(msg, cmdAppView)
}

// formatError returns a formatted error message with a link to the meshery command URL
func formatError(msg string, cmd cmdType) string {
	switch cmd {
	case cmdRoot:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, rootUsageURL)
	case cmdPerf:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, perfUsageURL)
	case cmdMesh:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, meshUsageURL)
	case cmdSystem:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, systemUsageURL)
	case cmdSystemStop:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, systemStopURL)
	case cmdSystemUpdate:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, systemUpdateURL)
	case cmdSystemReset:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, systemResetURL)
	case cmdSystemStatus:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, systemStatusURL)
	case cmdSystemRestart:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, systemRestartURL)
	case cmdExp:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, expUsageURL)
	case cmdFilter:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, filterUsageURL)
	case cmdFilterImport:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, filterImportURL)
	case cmdFilterDelete:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, filterDeleteURL)
	case cmdFilterList:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, filterListURL)
	case cmdFilterView:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, filterViewURL)
	case cmdPattern:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, patternUsageURL)
	case cmdPatternView:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, patternViewURL)
	case cmdApp:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, appViewURL)
	case cmdAppView:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, appUsageURL)
	case cmdContextDelete:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, contextDeleteURL)
	case cmdContextCreate:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, contextCreateURL)
	case cmdContextView:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, contextViewURL)
	case cmdContext:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, contextUsageURL)
	case cmdChannelSwitch:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, channelSwitchURL)
	case cmdChannelView:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, channelViewURL)
	case cmdChannelSet:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, channelSetURL)
	case cmdChannel:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, channelUsageURL)
	case cmdProviderView:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, providerViewURL)
	case cmdProviderList:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, providerListURL)
	case cmdProviderSet:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, providerSetURL)
	case cmdProviderSwitch:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, providerSwitchURL)
	case cmdProviderReset:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, providerResetURL)
	case cmdProvider:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, providerUsageURL)
	case cmdToken:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, tokenUsageURL)
	case cmdModel:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, modelUsageURL)
	case cmdModelList:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, modelListURL)
	case cmdModelView:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, modelViewURL)
	default:
		return fmt.Sprintf("%s\n", msg)
	}
}

func ErrFailRequest(err error) error {
	return errors.New(ErrFailRequestCode, errors.Alert,
		[]string{"Failed to make a request"},
		[]string{err.Error()},
		[]string{"Meshery server is not reachable."},
		[]string{"Ensure your Kubernetes cluster is running and your network connection is active. You can also try running 'mesheryctl system restart'."})
}

func ErrUnauthenticated() error {
	return errors.New(ErrUnauthenticatedCode, errors.Alert, []string{"Unauthenticated User"},
		[]string{"Access to this resource is unauthorized."},
		[]string{"You haven't logged in to Meshery."},
		[]string{"To proceed, log in using `mesheryctl system login`."})
}

func ErrInvalidToken() error {
	return errors.New(ErrInvalidTokenCode, errors.Alert,
		[]string{"Invalid authentication Token"},
		[]string{"The authentication token has expired or is invalid."},
		[]string{"The token in auth.json has expired or is invalid."},
		[]string{"Provide a valid user token by logging in with `mesheryctl system login`."})
}

func ErrUnmarshal(err error) error {
	return errors.New(ErrUnmarshalCode, errors.Alert,
		[]string{"Error unmarshalling response"},
		[]string{"Unable to process JSON response from server.\n" + err.Error()},
		[]string{"The JSON format from the response body is not valid."},
		[]string{"Ensure a valid JSON is provided for processing."})
}

func ErrFileRead(err error) error {
	return errors.New(ErrFileReadCode, errors.Alert,
		[]string{"File read error"},
		[]string{err.Error()},
		[]string{"The provided file is not present or has an invalid path."},
		[]string{"To proceed, provide a valid file path with a valid file."})
}

func ErrCreatingRequest(err error) error {
	return errors.New(ErrCreatingRequestCode, errors.Fatal,
		[]string{"Error occurred while making an HTTP request."},
		[]string{err.Error()},
		[]string{"Meshery is not running or there is a network issue."},
		[]string{"Check your network connection and verify the status of the Meshery server using `mesheryctl system status`."})
}

func ErrMarshal(err error) error {
	return errors.New(ErrMarshalCode, errors.Fatal,
		[]string{"Error while marshalling the content"},
		[]string{err.Error()},
		[]string{"The content provided for marshalling is invalid."},
		[]string{"Check the data structure you are providing for marshalling."})
}

func ErrReadResponseBody(err error) error {
	return errors.New(ErrReadResponseBodyCode, errors.Alert,
		[]string{"Failed to read response body from request"},
		[]string{err.Error()},
		[]string{"There might be connection failure with Meshery Server"},
		[]string{"Check the status via `mesheryctl system status`"})
}

func ErrParsingUrl(err error) error {
	return errors.New(ErrParsingUrlCode, errors.Fatal,
		[]string{"Error parsing the URL"},
		[]string{err.Error()},
		[]string{"The provided URL does not exist or the relative path is incorrect."},
		[]string{"Double-check the correctness of the URL you have inputted."})
}

func ErrNotFound(err error) error {
	return errors.New(ErrNotFoundCode, errors.Fatal,
		[]string{"Item Not Found"},
		[]string{err.Error()},
		[]string{"The item you are searching for is not present."},
		[]string{"Check whether the item is present."})
}

func ErrInvalidFile(err error) error {
	return errors.New(ErrInvalidFileCode, errors.Fatal,
		[]string{"Invalid File"},
		[]string{err.Error()},
		[]string{"File does not meet the criteria."},
		[]string{"Check the file's validity."})
}

func ErrInvalidNameOrID(err error) error {
	return errors.New(ErrInvalidNameOrIDCode, errors.Fatal,
		[]string{"Invalid Name or ID Provided"},
		[]string{err.Error()},
		[]string{"The provided Name or ID is either not present or invalid."},
		[]string{"Provide a valid Name or ID using the `list` command with the appropriate subcommand."})
}

func ErrAttachAuthToken(err error) error {
	return errors.New(ErrAttachAuthTokenCode, errors.Alert,
		[]string{"Authentication token Not Found"},
		[]string{"Authentication token not found: " + err.Error()},
		[]string{"The user is not logged in to generate a token."},
		[]string{"Log in with `mesheryctl system login` or supply a valid user token using the --token (or -t) flag."})
}

func ErrFailReqStatus(statusCode int) error {
	return errors.New(ErrFailReqStatusCode, errors.Alert,
		[]string{"Failed response server error"},
		[]string{"Response Status Code " + strconv.Itoa(statusCode) + ", possibly Server error"},
		[]string{"Invalid API call"},
		[]string{"Check your network connection and the status of Meshery Server via `mesheryctl system status`."})
}

func ErrMarshalIndent(err error) error {
	return errors.New(ErrMarshalIndentCode, errors.Alert,
		[]string{"Error indenting JSON body"},
		[]string{err.Error()},
		[]string{"Unable to format the JSON body filters due to invalid content"},
		[]string{"Check the data structure provided for proper formatting."})
}

func ErrResponseStatusBody(statusCode int, body string) error {
	return errors.New(ErrResponseStatusBodyCode, errors.Alert,
		[]string{"Incorrect status code"},
		[]string{"Server returned with status code: " + fmt.Sprint(statusCode) + "\n" + "Response: " + body},
		[]string{"Error occurred while generating a response body"},
		[]string{"Check your network connection and the status of Meshery Server via `mesheryctl system status`."})
}

func ErrResponseStatus(statusCode int) error {
	return errors.New(ErrResponseStatusCode, errors.Alert,
		[]string{"Incorrect status code"},
		[]string{"Server returned with status code: " + fmt.Sprint(statusCode)},
		[]string{"Error occurred while generating a response"},
		[]string{"Check your network connection and the status of Meshery Server via `mesheryctl system status`."})
}

func ErrJSONToYAML(err error) error {
	return errors.New(ErrJSONToYAMLCode, errors.Alert,
		[]string{"Failed to convert JSON to YAML"},
		[]string{"Error occurred while converting JSON to YAML: " + err.Error()},
		[]string{"The content provided for marshalling is invalid."},
		[]string{"Check the JSON structure you are providing for YAML conversion."})
}

func ErrOutFormatFlag() error {
	return errors.New(ErrOutFormatFlagCode, errors.Alert,
		[]string{"Invalid output format choice"},
		[]string{"Output format choice is invalid, use [json|yaml]"},
		[]string{"Invalid JSON or YAML content"},
		[]string{"Check the JSON or YAML structure."})
}

func ErrReadConfigFile(err error) error {
	return errors.New(ErrReadConfigFileCode, errors.Alert,
		[]string{"Unable to read meshconfig file"},
		[]string{"Unable to read the meshconfig file from the specified path: " + err.Error()},
		[]string{"The provided file is not present or has an invalid path"},
		[]string{"Provide a valid file path with a valid meshconfig file."})
}

func ErrInvalidAPIResponse(err error) error {
	return errors.New(ErrInvalidAPIResponseCode, errors.Fatal,
		[]string{"Invalid API response encountered"},
		[]string{"Invalid API response encountered: " + err.Error()},
		[]string{"Error occurred while generating a response body"},
		[]string{"Check your network connection and the status of Meshery Server via `mesheryctl system status`."})
}

func ErrLoadConfig(err error) error {
	return errors.New(ErrLoadConfigCode, errors.Alert,
		[]string{"Error processing config"},
		[]string{"Error processing config:" + err.Error()},
		[]string{"Unable to load meshconfig due to wrong configurations"},
		[]string{"Ensure your `config.yaml` file in your `.meshery` is valid, or run `mesheryctl system config`."})
}

func ErrParseGithubFile(err error, URL string) error {
	return errors.New(ErrParseGithubFileCode, errors.Alert,
		[]string{"Failed to parse github file"},
		[]string{"Failed to parse github file" + err.Error()},
		[]string{"Unable to retrieve file from URL: %s", URL},
		[]string{"Ensure you have a github url in file path"})
}

func ErrReadToken(err error) error {
	return errors.New(ErrReadTokenCode, errors.Alert,
		[]string{"Could not read token"},
		[]string{err.Error()},
		[]string{"Token file is invalid"},
		[]string{"Provide a valid user token by logging in with `mesheryctl system login`."})
}

func ErrRequestResponse(err error) error {
	return errors.New(ErrRequestResponseCode, errors.Alert,
		[]string{"Failed to handle request"},
		[]string{"Unable to create a response from request" + err.Error()},
		[]string{"Error occurred while generating a response"},
		[]string{"Check your network connection and the status of Meshery Server via `mesheryctl system status`."})
}
