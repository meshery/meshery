package machines

import (
	"github.com/layer5io/meshkit/errors"
)

var ErrInvalidTransitionCode = "replace_me"

var ErrInvalidTransition = errors.New(ErrInvalidTransitionCode, errors.Alert, []string{}, []string{}, []string{}, []string{})