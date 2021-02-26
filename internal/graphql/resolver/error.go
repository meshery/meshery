package resolver

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	errCode               = "test_code"
	ErrInvalidRequestCode = "test_code"
)

var (
	ErrInvalidRequest = errors.NewDefault(ErrInvalidRequestCode, "Invalid query, please check syntax")
)
