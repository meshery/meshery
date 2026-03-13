package components

const (
	searchUsageMsg      = "Usage: mesheryctl component search [query-text]\nRun 'mesheryctl component search --help' to see detailed help message"
	viewUsageMsg        = "Usage: mesheryctl component view [component-name | component-id]\nRun 'mesheryctl component view --help' to see detailed help message"
	errInvalidArg       = "Only one argument must be provided and needs to be enclosed by double quotes if it contains spaces (eg. \"model name\", modelName)"
	errNoComponentFound = "No component(s) found with the name or id: "
)
