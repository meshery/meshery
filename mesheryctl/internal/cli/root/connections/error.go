package connections

import (
	"fmt"
	"strings"

	"github.com/meshery/meshkit/errors"
)

var (
	ErrConnectionTypeCode = "replace_me"
)

func errInvalidConnectionType(connectionType string) error {
	return errors.New(ErrConnectionTypeCode, errors.Alert, []string{fmt.Sprintf("Invalid connection type %s provided", connectionType)}, []string{fmt.Sprintf("the provided connection type is not yet supported")}, []string{}, []string{fmt.Sprintf("Please provide a valid connection type: %s", strings.Join(supportedConnectionTypes, ", "))})
}
