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
	ErrFailRequestCode        = "1044"
	ErrFailReqStatusCode      = "1045"
	ErrAttachAuthTokenCode    = "1046"
	ErrReadConfigFileCode     = "replace_me"
	ErrMarshalIndentCode      = "replace_me"
	ErrProcessingConfigCode   = "replace_me"
	ErrResponseStatusBodyCode = "replace_me"
	ErrResponseStatusCode     = "replace_me"
	ErrJSONToYAMLCode         = "replace_me"
	ErrOutFormatFlagCode      = "replace_me"
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

// AppError returns a formatted error message with a link to 'app' command usage page in addition to the error message
func AppError(msg string) string {
	return formatError(msg, cmdApp)
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
	case cmdApp:
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
	default:
		return fmt.Sprintf("%s\n", msg)
	}
}

func ErrAttachAuthToken(err error) error {
	return errors.New(ErrAttachAuthTokenCode, errors.Alert, []string{"Authentication token not found"},
		[]string{"Authentication token not found", err.Error()}, []string{"User not logged in to generate token"}, []string{ "Login with `mesheryctl system login` or supply a valid user token using the --token (or -t) flag."})
}

func ErrFailRequest(err error) error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{"Failed to make a request"},
		[]string{"Failed to make request" + err.Error()}, []string{"Probable network connection problem"}, []string{"Check your network connection and check the status of meshery server via 'mesheryctl system status'"})
}

func ErrUnauthenticated() error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{"Authentication token is invalid"},
		[]string{"Authentication token is invalid. Please supply a valid user token. "}, []string{"User not logged in correctly to generate token"}, []string{"Please login with `mesheryctl system login`"})
}

func InvalidToken() error {
	return errors.New(ErrFailRequestCode, errors.Alert, []string{},
		[]string{"Authentication token is expired/invalid."}, []string{"Token expired"}, []string{" Please login with `mesheryctl system login` to generate a new token"})
}

func ErrFailReqStatus(statusCode int) error {
	return errors.New(ErrFailReqStatusCode, errors.Alert, []string{"Failed response Server Error"},
		[]string{"Response Status Code " + strconv.Itoa(statusCode) + ", possible Server Error"}, []string{"Bad api call"}, []string{"Check your network connection and check the status of meshery server via 'mesheryctl system status'"})
}

func ErrProcessingConfig(err error) error {
	return errors.New(ErrProcessingConfigCode, errors.Alert, []string{"Error processing config"},
		[]string{"Error processing config", err.Error()}, []string{"Couldn't load mesh config"}, []string{"Ensure mesheryctl has the right configurations"})
}

func ErrMarshalIndent(err error) error {
	return errors.New(ErrMarshalIndentCode, errors.Alert,
		[]string{"Error indenting json body"}, []string{err.Error()}, []string{"Can't marshal indent body filters due to invalid content"}, []string{"Please check the data structure you are providing for Marshalling"})
}

func ErrResponseStatusBody(statusCode int, body string) error {
	return errors.New(ErrResponseStatusBodyCode, errors.Alert, []string{"Wrong status code"},
		[]string{"Server returned with status code: " + fmt.Sprint(statusCode) + "\n" + "Response: " + body}, []string{"Error occured while generating a response body"}, []string{"Check your network connection and check the status of meshery server via 'mesheryctl system status'"})
}

func ErrResponseStatus(statusCode int) error {
	return errors.New(ErrResponseStatusCode, errors.Alert, []string{"Wrong status code"},
		[]string{"Server returned with status code: " + fmt.Sprint(statusCode)}, []string{"Error occured while generating a respons"}, []string{"Check your network connection and check the status of meshery server via 'mesheryctl system status'"})
}

func ErrJSONToYAML(err error) error {
	return errors.New(ErrJSONToYAMLCode, errors.Alert, []string{"Can't convert JSON to yaml"},
		[]string{"Failed to convert json to yaml", "" + err.Error()}, []string{"The Content provided to Marshall is invalid"}, []string{"Please check the json structure you are providing for yaml conversion"})
}

func ErrOutFormatFlag() error {
	return errors.New(ErrOutFormatFlagCode, errors.Alert, []string{"Output-format choice invalid"},
		[]string{"Output-format choice invalid, use [json|yaml]"}, []string{"Invalid json or yaml content"}, []string{"Please check the json or yaml structure"})
}

func ErrReadConfigFile(err error) error {
	return errors.New(ErrReadConfigFileCode, errors.Alert, []string{"Unable to read config file"},
		[]string{"Can't read config file from Path", "" + err.Error()}, []string{"Provided file is not present or invalid path"}, []string{"Please provide a valid file path with a valid config file"})
}
