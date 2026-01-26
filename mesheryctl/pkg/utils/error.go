package utils

import (
	"fmt"
	"strconv"

	"github.com/meshery/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshery.io/project/contributing/contributing-error
// https://github.com/meshery/meshkit/blob/master/errors/errors.go
var (
	ErrFailRequestCode             = "mesheryctl-1090"
	ErrInvalidTokenCode            = "mesheryctl-1091"
	ErrFailReqStatusCode           = "mesheryctl-1092"
	ErrAttachAuthTokenCode         = "mesheryctl-1093"
	ErrUnmarshalCode               = "mesheryctl-1094"
	ErrFileReadCode                = "mesheryctl-1095"
	ErrCreatingRequestCode         = "mesheryctl-1096"
	ErrMarshalCode                 = "mesheryctl-1097"
	ErrReadResponseBodyCode        = "mesheryctl-1098"
	ErrParsingUrlCode              = "mesheryctl-1099"
	ErrNotFoundCode                = "mesheryctl-1100"
	ErrUnauthenticatedCode         = "mesheryctl-1101"
	ErrInvalidFileCode             = "mesheryctl-1102"
	ErrInvalidNameOrIDCode         = "mesheryctl-1103"
	ErrInvalidAPIResponseCode      = "mesheryctl-1104"
	ErrReadConfigFileCode          = "mesheryctl-1105"
	ErrMarshalIndentCode           = "mesheryctl-1106"
	ErrLoadConfigCode              = "mesheryctl-1107"
	ErrResponseStatusBodyCode      = "mesheryctl-1108"
	ErrResponseStatusCode          = "mesheryctl-1109"
	ErrJSONToYAMLCode              = "mesheryctl-1110"
	ErrOutFormatFlagCode           = "mesheryctl-1111"
	ErrParseGithubFileCode         = "mesheryctl-1112"
	ErrReadTokenCode               = "mesheryctl-1113"
	ErrRequestResponseCode         = "mesheryctl-1114"
	ErrBadRequestCode              = "mesheryctl-1117"
	ErrInvalidArgumentCode         = "mesheryctl-1118"
	ErrGeneratingIconsCode         = "mesheryctl-1119"
	ErrClearLineCode               = "mesheryctl-1120"
	ErrGeneratesModelCode          = "mesheryctl-1132"
	ErrUpdateComponentsCode        = "mesheryctl-1134"
	ErrMissingCommandsCode         = "mesheryctl-1137"
	ErrKubernetesConnectivityCode  = "mesheryctl-1138"
	ErrKubernetesQueryCode         = "mesheryctl-1139"
	ErrCreateManifestsFolderCode   = "mesheryctl-1141"
	ErrDownloadFileCode            = "mesheryctl-1142"
	ErrNoManifestFilesFoundCode    = "mesheryctl-1143"
	ErrWalkManifestsCode           = "mesheryctl-1144"
	ErrGetChannelVersionCode       = "mesheryctl-1145"
	ErrInvalidModelCode            = "mesheryctl-1150"
	ErrInvalidUUIDCode             = "mesheryctl-1152"
	ErrFetchEnvironmentsCode       = "mesheryctl-1153"
	ErrTableRenderCode             = "mesheryctl-1154"
	ErrFlagsInvalidCode            = "mesheryctl-1155"
	ErrMesheryServerNotRunningCode = "mesheryctl-1156"
	ErrHandlePaginationCode        = "mesheryctl-1172"
	ErrCreateFileCode              = "mesheryctl-1123"
	ErrRetrieveHomeDirCode         = "mesheryctl-1124"
	ErrReadFromBodyCode            = "mesheryctl-1125"
	ErrMarkFlagRequireCode         = "mesheryctl-1126"
	ErrGetKubernetesContextsCode   = "mesheryctl-1165"
	ErrSetKubernetesContextCode    = "mesheryctl-1166"
	ErrReadInputCode               = "mesheryctl-1193"
	ErrUploadFileWithParamsCode    = "mesheryctl-1185"
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
	case "import":
		return formatError(msg, cmdModelImport)
	default:
		return formatError(msg, cmdModel)
	}
}

func ComponentSubError(msg string, cmd string) string {
	switch cmd {
	case "list":
		return formatError(msg, cmdComponentList)
	case "view":
		return formatError(msg, cmdComponentView)
	case "search":
		return formatError(msg, cmdComponentSearch)
	default:
		return formatError(msg, cmdComponent)
	}
}

func EnvironmentSubError(msg string, cmd string) string {
	switch cmd {
	case "create":
		return formatError(msg, cmdEnvironmentCreate)
	case "delete":
		return formatError(msg, cmdEnvironmentDelete)
	case "list":
		return formatError(msg, cmdEnvironmentList)
	case "view":
		return formatError(msg, cmdEnvironmentView)
	default:
		return formatError(msg, cmdEnvironment)
	}
}

func WorkspaceSubError(msg string, cmd string) string {
	switch cmd {
	case "list":
		return formatError(msg, cmdExpWorkspaceList)
	case "create":
		return formatError(msg, cmdExpWorkspaceCreate)
	default:
		return formatError(msg, cmdExpWorkspace)
	}
}

func RegistryError(msg string, cmd string) string {
	switch cmd {
	case "publish":
		return formatError(msg, cmdRegistryPublish)
	case "generate":
		return formatError(msg, cmdRegistryGenerate)
	case "update":
		return formatError(msg, cmdRegistryUpdate)
	default:
		return formatError(msg, cmdRegistry)
	}
}

func RelationshipsError(msg string, cmd string) string {
	switch cmd {
	case "view":
		return formatError(msg, cmdRelationshipView)
	case "generate":
		return formatError(msg, cmdRelationshipGenerateDocs)
	case "search":
		return formatError(msg, cmdRelationshipSearch)
	case "list":
		return formatError(msg, cmdRelationshipList)
	default:
		return formatError(msg, cmdRelationships)
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

// DesignError returns a formatted error message with a link to 'design' command usage page in addition to the error message
func DesignError(msg string) string {
	return formatError(msg, cmdDesign)
}

// DesignViewError returns a formatted error message with a link to the 'design view' command usage page in addition to the error message
func DesignViewError(msg string) string {
	return formatError(msg, cmdDesignView)
}

// formatError returns a formatted error message with a link to the meshery command URL
func formatError(msg string, cmd cmdType) string {
	switch cmd {
	case cmdRoot:
		return formatUsageDetails(msg, rootUsageURL)
	case cmdPerf:
		return formatUsageDetails(msg, perfUsageURL)
	case cmdMesh:
		return formatUsageDetails(msg, meshUsageURL)
	case cmdSystem:
		return formatUsageDetails(msg, systemUsageURL)
	case cmdSystemStop:
		return formatUsageDetails(msg, systemStopURL)
	case cmdSystemUpdate:
		return formatUsageDetails(msg, systemUpdateURL)
	case cmdSystemReset:
		return formatUsageDetails(msg, systemResetURL)
	case cmdSystemStatus:
		return formatUsageDetails(msg, systemStatusURL)
	case cmdSystemRestart:
		return formatUsageDetails(msg, systemRestartURL)
	case cmdExp:
		return formatUsageDetails(msg, expUsageURL)
	case cmdFilter:
		return formatUsageDetails(msg, filterUsageURL)
	case cmdFilterImport:
		return formatUsageDetails(msg, filterImportURL)
	case cmdFilterDelete:
		return formatUsageDetails(msg, filterDeleteURL)
	case cmdFilterList:
		return formatUsageDetails(msg, filterListURL)
	case cmdFilterView:
		return formatUsageDetails(msg, filterViewURL)
	case cmdDesign:
		return formatUsageDetails(msg, designUsageURL)
	case cmdDesignView:
		return formatUsageDetails(msg, designViewURL)
	case cmdDesignExport:
		return formatUsageDetails(msg, designExportURL)
	case cmdContextDelete:
		return formatUsageDetails(msg, contextDeleteURL)
	case cmdContextCreate:
		return formatUsageDetails(msg, contextCreateURL)
	case cmdContextView:
		return formatUsageDetails(msg, contextViewURL)
	case cmdContext:
		return formatUsageDetails(msg, contextUsageURL)
	case cmdChannelSwitch:
		return formatUsageDetails(msg, channelSwitchURL)
	case cmdChannelView:
		return formatUsageDetails(msg, channelViewURL)
	case cmdChannelSet:
		return formatUsageDetails(msg, channelSetURL)
	case cmdChannel:
		return formatUsageDetails(msg, channelUsageURL)
	case cmdProviderView:
		return formatUsageDetails(msg, providerViewURL)
	case cmdProviderList:
		return formatUsageDetails(msg, providerListURL)
	case cmdProviderSet:
		return formatUsageDetails(msg, providerSetURL)
	case cmdProviderSwitch:
		return formatUsageDetails(msg, providerSwitchURL)
	case cmdProviderReset:
		return formatUsageDetails(msg, providerResetURL)
	case cmdProvider:
		return formatUsageDetails(msg, providerUsageURL)
	case cmdToken:
		return formatUsageDetails(msg, tokenUsageURL)
	case cmdModel:
		return formatUsageDetails(msg, modelUsageURL)
	case cmdModelList:
		return formatUsageDetails(msg, modelListURL)
	case cmdModelImport:
		return formatUsageDetails(msg, modelImportURl)
	case cmdModelView:
		return formatUsageDetails(msg, modelViewURL)
	case cmdRegistry:
		return formatUsageDetails(msg, registryUsageURL)
	case cmdRegistryPublish:
		return formatUsageDetails(msg, registryPublishURL)
	case cmdRegistryGenerate:
		return formatUsageDetails(msg, registryGenerateURL)
	case cmdRegistryUpdate:
		return formatUsageDetails(msg, registryUpdateURL)
	case cmdEnvironment:
		return formatUsageDetails(msg, environmentUsageURL)
	case cmdEnvironmentCreate:
		return formatUsageDetails(msg, environmentCreateURL)
	case cmdEnvironmentDelete:
		return formatUsageDetails(msg, environmentDeleteURL)
	case cmdEnvironmentList:
		return formatUsageDetails(msg, environmentListURL)
	case cmdEnvironmentView:
		return formatUsageDetails(msg, environmentViewURL)
	case cmdExpWorkspace:
		return formatUsageDetails(msg, workspaceUsageURL)
	case cmdExpWorkspaceCreate:
		return formatUsageDetails(msg, workspaceCreateURL)
	case cmdExpWorkspaceList:
		return formatUsageDetails(msg, workspaceListURL)
	case cmdRelationshipView:
		return formatUsageDetails(msg, relationshipViewURL)
	case cmdRelationships:
		return formatUsageDetails(msg, relationshipUsageURL)
	case cmdRelationshipGenerateDocs:
		return formatUsageDetails(msg, cmdRelationshipGenerateDocsURL)
	case cmdComponent:
		return formatUsageDetails(msg, componentUsageURL)
	case cmdComponentList:
		return formatUsageDetails(msg, componentListURL)
	case cmdComponentSearch:
		return formatUsageDetails(msg, componentSearchURL)
	case cmdComponentView:
		return formatUsageDetails(msg, componentViewURL)
	case cmdConnection:
		return formatUsageDetails(msg, connectionUsageURL)
	case cmdConnectionDelete:
		return formatUsageDetails(msg, connectionDeleteURL)
	case cmdConnectionList:
		return formatUsageDetails(msg, connectionListURL)
	case cmdExpRelationship:
		return formatUsageDetails(msg, expRelationshipUsageURL)
	case cmdExpRelationshipGenerate:
		return formatUsageDetails(msg, expRelationshipGenerateURL)
	case cmdExpRelationshipView:
		return formatUsageDetails(msg, expRelationshipViewURL)
	case cmdExpRelationshipList:
		return formatUsageDetails(msg, expRelationshipListURL)
	default:
		return fmt.Sprintf("%s\n", msg)
	}
}

func formatUsageDetails(msg string, docURL string) string {
	return fmt.Sprintf("%s\nSee %s for usage details\n", msg, docURL)
}

func ErrFailRequest(err error) error {
	return errors.New(ErrFailRequestCode, errors.Alert,
		[]string{"Failed to make a request"},
		[]string{err.Error()},
		[]string{"Meshery server is not reachable."},
		[]string{"Please ensure that the Meshery server is running and accessible. You can try running 'mesheryctl system restart'."})
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
		[]string{"Authentication token not found: " + err.Error() + "\nLog in with `mesheryctl system login`"},
		[]string{"The user is not logged in to generate a token."},
		[]string{"Log in with `mesheryctl system login` or supply a valid user token using the --token (or -t) flag."})
}

func ErrCreateManifestsFolder(err error) error {
	return errors.New(ErrCreateManifestsFolderCode, errors.Alert, []string{"Error creating manifest folder"}, []string{err.Error()}, []string{"system error in creating manifest folder"}, []string{"Make sure manifest folder (.meshery/manifests) is created properly"})
}

func ErrFailReqStatus(statusCode int, obj string) error {
	return errors.New(ErrFailReqStatusCode, errors.Alert,
		[]string{"Failed response server error"},
		[]string{"Response Status Code " + strconv.Itoa(statusCode) + ". Server emitted an error: " + obj},
		[]string{"Invalid API call"},
		[]string{"Check your network connection and the status of Meshery Server via `mesheryctl system status`."})
}

func ErrGenerateModel(err error, modelName string) error {
	return errors.New(ErrGeneratesModelCode, errors.Alert, []string{fmt.Sprintf("error generating model: %s", modelName)}, []string{fmt.Sprintf("Error generating model: %s\n %s", modelName, err.Error())}, []string{"Registrant used for the model is not supported", "Verify the model's source URL.", "Failed to create a local directory in the filesystem for this model."}, []string{"Ensure that each kind of registrant used is a supported kind.", "Ensure correct model source URL is provided and properly formatted.", "Ensure sufficient permissions to allow creation of model directory."})
}

func ErrDownloadFile(err error, obj string) error {
	return errors.New(ErrDownloadFileCode, errors.Alert, []string{"Error downloading file ", obj}, []string{err.Error()}, []string{"Failed to download docker-compose or manifest file due to system/config/network issues"}, []string{"Make sure docker-compose or manifest file is downloaded"})
}

func ErrWalkManifests(err error) error {
	return errors.New(
		ErrWalkManifestsCode,
		errors.Alert,
		[]string{"Error walking through manifests"},
		[]string{err.Error()},
		[]string{"Unable to traverse git repository or manifests due to filesystem or permission issues."},
		[]string{"Ensure the repository and manifests directory are accessible and have proper permissions."},
	)
}

func ErrNoManifestFilesFound(path string) error {
	return errors.New(
		ErrNoManifestFilesFoundCode,
		errors.Alert,
		[]string{"No manifest files found in the specified path"},
		[]string{fmt.Sprintf("No manifest files present in path: %s", path)},
		[]string{"The provided directory may be empty, incorrect, or manifests were not properly downloaded."},
		[]string{"Verify the specified path contains valid manifest files."},
	)
}

func ErrGetChannelVersion(err error) error {
	return errors.New(
		ErrGetChannelVersionCode,
		errors.Alert,
		[]string{"Unable to retrieve release channel and version information."},
		[]string{err.Error()},
		[]string{"Failed to determine version from context or GitHub releases, possibly due to network or configuration issues."},
		[]string{"Check your network connection and context configuration; ensure GitHub is accessible."},
	)
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

func ErrBadRequest(err error) error {
	return errors.New(ErrBadRequestCode, errors.Alert,
		[]string{"Failed to delete the connection"},
		[]string{err.Error()},
		[]string{"Error occurred while deleting the connection"},
		[]string{"Check your network connection and the status of Meshery Server via `mesheryctl system status`."})
}

func ErrInvalidArgument(err error) error {
	return errors.New(ErrInvalidArgumentCode, errors.Alert, []string{"Invalid Argument"}, []string{err.Error()}, []string{"Invalid Argument"}, []string{"Please check the arguments passed"})
}

func ErrGeneratingIcons(err error, path string) error {
	return errors.New(
		ErrGeneratingIconsCode,
		errors.Alert,
		[]string{"error generating icons at ", path},
		[]string{err.Error()},
		[]string{"Model SVG data is missing", "Model name formatting issue"},
		[]string{"Ensure model SVG data is provided in model definition", "Ensure model name formatting is correct"},
	)
}

func ErrClearLine(err error) error {
	return errors.New(ErrClearLineCode, errors.Alert,
		[]string{"Failed to clear terminal"},
		[]string{err.Error()},
		[]string{"Error occurred while attempting to clear the command-line interface"},
		[]string{"Check if the required clear commands ('clear' or 'cls') are available in the system's PATH"})
}

func ErrInvalidUUID(err error) error {
	return errors.New(ErrInvalidUUIDCode, errors.Alert,
		[]string{"Invalid ID format"},
		[]string{err.Error()},
		[]string{"ID is not a valid UUID format", "ID contains invalid characters"},
		[]string{"Ensure the ID is a valid UUID format", "Check the ID parameter for typos or formatting issues"})
}

func ErrFetchEnvironments(err error) error {
	return errors.New(ErrFetchEnvironmentsCode, errors.Alert,
		[]string{"Failed to fetch environments"},
		[]string{err.Error()},
		[]string{"Network connection issues", "Meshery server is not reachable", "Invalid API endpoint", "Authentication issues"},
		[]string{"Check your network connection", "Verify Meshery server status via `mesheryctl system status`", "Ensure you have proper authentication credentials"})
}

func ErrUpdateComponent(err error, modelName, compName string) error {
	return errors.New(ErrUpdateComponentsCode, errors.Alert, []string{fmt.Sprintf("error updating component %s of model %s ", compName, modelName)}, []string{err.Error()}, []string{"Component does not exist", "Component definition is corrupted"}, []string{"Ensure existence of component, check for typo in component name", "Regenerate corrupted component"})
}

func ErrInvalidModel() error {
	return errors.New(ErrInvalidModelCode, errors.Alert, []string{"No valid component or relationship found in the model provided"}, []string{"No valid component or relationship found in the Model provided. A Model can be only imported if it contains at least one valid Component or Relationship."}, []string{"Provided components or relationships might have incorrect format", "Folder structure might be incorrect"}, []string{"Know about Meshery Models and Importing instructions here: https://docs.meshery.io/guides/configuration-management/importing-models"})
}

func ErrMissingCommands(err error) error {
	return errors.New(ErrMissingCommandsCode, errors.Alert, []string{"Missing required commands"}, []string{err.Error()}, []string{"Required commands are not installed or not in PATH"}, []string{"Install the missing commands and ensure they are in PATH"})
}

func ErrKubernetesConnectivity(err error) error {
	return errors.New(ErrKubernetesConnectivityCode, errors.Alert, []string{"Failed to connect to Kubernetes API server"}, []string{err.Error()}, []string{"Kubernetes API server is not reachable"}, []string{"Ensure your Kubernetes cluster is running and accessible"})
}

func ErrKubernetesQuery(err error) error {
	return errors.New(ErrKubernetesQueryCode, errors.Alert, []string{"Failed to query Kubernetes API"}, []string{err.Error()}, []string{"Kubernetes API query failed"}, []string{"Ensure your Kubernetes cluster is running and accessible"})
}

func ErrTableRender(err error) error {
	return errors.New(ErrTableRenderCode, errors.Alert,
		[]string{"Failed to display output in table format"},
		[]string{err.Error()},
		[]string{"Table rendering issue"},
		[]string{"Ensure the data being rendered is valid and properly structured."})
}

func ErrFlagsInvalid(msg string) error {
	return errors.New(ErrFlagsInvalidCode, errors.Alert,
		[]string{"Invalid flag provided"},
		[]string{msg},
		[]string{"The flag provided is invalid."},
		[]string{"Provide a valid flag"})
}

func ErrMesheryServerNotRunning(platform string) error {
	return errors.New(ErrMesheryServerNotRunningCode, errors.Alert,
		[]string{"Meshery Server is not running"},
		[]string{fmt.Sprintf("Meshery Server is not available on platform: %s", platform)},
		[]string{"Meshery Server is not running or is unreachable", "Docker or Kubernetes environment is not ready", "Network connectivity issues"},
		[]string{"Start Meshery Server with `mesheryctl system start`", "Verify system readiness with `mesheryctl system check --preflight`", "Check your network connection and firewall settings"})
}

func ErrHandlePagination(err error) error {
	return errors.New(ErrHandlePaginationCode, errors.Alert,
		[]string{"Unable to display paginated results"},
		[]string{err.Error()},
		[]string{"Interactive pagination requires keyboard input support"},
		[]string{
			"Ensure you are running in an interactive terminal",
			"If running in a non-interactive environment, use '--page' flag to skip pagination",
		})
}

func ErrCreateFile(filepath string, err error) error {
	return errors.New(ErrCreateFileCode, errors.Alert,
		[]string{"Error creating file"},
		[]string{fmt.Sprintf("Failed to create the file at path: %s", filepath), err.Error()},
		[]string{"Insufficient disk page, filepath could be invalid."},
		[]string{"Verify that the file path is valid, and ensure there is sufficient disk space available."})
}

func ErrRetrieveHomeDir(err error) error {
	return errors.New(ErrRetrieveHomeDirCode, errors.Alert,
		[]string{"Error retrieving user home/root directory"},
		[]string{"Failed to retrieve the home/root directory,", err.Error()},
		[]string{"Operating system environment issue or insufficient permissions."},
		[]string{"Ensure that the operating system environment is set up correctly and run the application with elevated privileges."})
}

func ErrMarkFlagRequire(flagName string, err error) error {
	return errors.New(ErrMarkFlagRequireCode, errors.Alert,
		[]string{fmt.Sprintf("Failed to mark the flag '%s' as required", flagName)},
		[]string{err.Error()},
		[]string{"The flag may not exist or there was some error while specifying the flag."},
		[]string{"Please ensure that the required flag '%s' is correctly specified and set before running the command."})
}

func ErrReadFromBody(err error) error {
	return errors.New(ErrReadFromBodyCode, errors.Alert,
		[]string{"Unable to read data from the response body"},
		[]string{err.Error()},
		[]string{"The data for the pattern (design) file might be corrupted."},
		[]string{"Please ensure that your network connection is stable. If the issue continues, check the server response or data format for potential problems."})
}

func ErrGetKubernetesContexts(err error) error {
	return errors.New(
		ErrGetKubernetesContextsCode,
		errors.Fatal,
		[]string{"Unable to get kubernetes contexts"},
		[]string{err.Error()},
		[]string{"No kubernetes contexts found"},
		[]string{"Ensure you have at least one valid context in your meshconfig file."})
}

func ErrSetKubernetesContext(err error) error {
	return errors.New(
		ErrSetKubernetesContextCode,
		errors.Fatal,
		[]string{"Unable to set kubernetes context"},
		[]string{err.Error()},
		[]string{"The specified Kubernetes context does not exist. "},
		[]string{"Verify that the Kubernetes context provided is valid and try again."})
}

func ErrReadInput(err error) error {
	return errors.New(
		ErrReadInputCode,
		errors.Fatal,
		[]string{"Unable to read the input"},
		[]string{err.Error()},
		[]string{"The provided input was invalid or could not be read"},
		[]string{"Validate the input and try again."})
}

func ErrUploadFileWithParams(err error, fileName string) error {
	return errors.New(
		ErrUploadFileWithParamsCode,
		errors.Alert,
		[]string{fmt.Sprintf("Failed to upload file: %s", fileName)},
		[]string{err.Error()},
		[]string{"File upload failed due to network issues or server errors"},
		[]string{"Check your network connection and ensure the server is reachable."},
	)
}
