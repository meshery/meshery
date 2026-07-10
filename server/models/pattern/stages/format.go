package stages

import (
	"github.com/meshery/meshery/server/models/pattern/core"

	"github.com/sirupsen/logrus"
)

// Format - Format stage performs de-prettification of the pattern file
func Format() ChainStageFunction {
	return func(data *Data, err error, next ChainStageNextFunction) {
		if err != nil {
			next(data, err)
			return
		}
		logrus.Debug("de-prettifying the pattern file")
		for _, component := range data.Pattern.Components {
			component.Configuration = core.Format.DePrettify(component.Configuration, false)
		}

		if next != nil {
			next(data, nil)
		}
	}
}
