package components

const (
<<<<<<< HEAD
<<<<<<< HEAD
	errViewCmdMsg       = "Usage: mesheryctl component view [component-name | component-id]\nRun 'mesheryctl component view --help' to see detailed help message"
=======
	searchUsageMsg      = "Usage: mesheryctl component search [query-text]\nRun 'mesheryctl component search --help' to see detailed help message"
	viewUsageMsg        = "Usage: mesheryctl component view [component-name | component-id]\nRun 'mesheryctl component view --help' to see detailed help message"
	errInvalidArg       = "Only one argument must be provided and needs to be enclosed by double quotes if it contains spaces (eg. \"model name\", modelName)"
>>>>>>> fab0f34435c (resolved conflits)
	errNoComponentFound = "No component(s) found with the name or id: "
=======
	errViewCmdMsg      = "Usage: mesheryctl component view [component-name]\nRun 'mesheryctl component view --help' to see detailed help message"
	errNoArg           = "no component name specified as an argument"
	errMultiArg        = "too many arguments specified, name with spaces must be enclosed in double quotes"
	searchUsageMessage = "Usage: mesheryctl component search [component-name]\nRun 'mesheryctl component search --help' to see detailed help message"
>>>>>>> ac983b60c5c (added page flag)
)
