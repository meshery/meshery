package connections

import (
	"fmt"
	"strings"

	"github.com/meshery/meshkit/errors"
)

var (
	ErrConnectionTypeCode = "mesheryctl-1163"

	invalidOuptputFormatMsg = "output-format choice is invalid, use [json|yaml]"
)

func errInvalidConnectionType(connectionType string) error {
	return errors.New(ErrConnectionTypeCode, errors.Alert, []string{fmt.Sprintf("Invalid connection type %s provided", connectionType)}, []string{fmt.Sprintf("the provided connection type is not yet supported. Supported type(s) are: %s", strings.Join(supportedConnectionTypes, ", "))}, []string{}, []string{fmt.Sprintf("Please provide a valid connection type: %s", strings.Join(supportedConnectionTypes, ", "))})
}
