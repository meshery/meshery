package display

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/errors"
)

func GetErrorLog(err error) string {
	var errMsg strings.Builder
	if meshkitErr, ok := err.(*errors.Error); ok {

		errMsg.WriteString(fmt.Sprintf("\n  %s: %s", utils.BoldString("ERROR CODE"), meshkitErr.Code))

		if len(meshkitErr.ShortDescription) != 0 {
			errMsg.WriteString(fmt.Sprintf("\n  %s: %s", utils.BoldString("DESCRIPTION"), strings.Join(meshkitErr.ShortDescription, ". ")))
		}

		if len(meshkitErr.LongDescription) != 0 {
			errMsg.WriteString(fmt.Sprintf("\n  %s: %s", utils.BoldString("ERROR"), strings.Join(meshkitErr.LongDescription, ". ")))
		}

		if len(meshkitErr.ProbableCause) != 0 {
			errMsg.WriteString(fmt.Sprintf("\n  %s: %s", utils.BoldString("PROBABLE CAUSE"), strings.Join(meshkitErr.ProbableCause, ". ")))
		}

		if len(meshkitErr.SuggestedRemediation) != 0 {
			errMsg.WriteString(fmt.Sprintf("\n  %s: %s", utils.BoldString("SUGGESTED REMEDIATION"), strings.Join(meshkitErr.SuggestedRemediation, ". ")))
		}
	}
	return errMsg.String()
}
