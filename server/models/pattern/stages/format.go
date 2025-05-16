package stages

import (
	"fmt"

	"github.com/layer5io/meshery/server/models/pattern/core"
)

// Format - Format stage performs de-prettification of the pattern file
func Format() ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			next(data, err)
			return
		}
		fmt.Println("Deprettifying the pattern file")
		for _, component := range data.Pattern.Components {
			component.Configuration = core.Format.DePrettify(component.Configuration, false)
		}

		if next != nil {
			next(data, nil)
		}
	}
}
