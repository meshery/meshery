package relationships

import "fmt"

var (
	viewUsageMsg           = "\n\nUsage: mesheryctl exp relationship view [model-name]\nRun 'mesheryctl exp relationship view --help' to see detailed help message"
	errNoModelNameProvided = fmt.Errorf("[model-name] isn't specified%s", viewUsageMsg)
	errTooManyArgs         = fmt.Errorf("too many arguments, only [model-name] is expected%s", viewUsageMsg)
)
