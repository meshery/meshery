package system

import (
	"fmt"
	"testing"

	"github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrHealthCheckFailedSeparatesCause(t *testing.T) {
	err := ErrHealthCheckFailed(fmt.Errorf("valid context is not available in meshconfig"))

	assert.Equal(t, "Failed to initialize healthchecker: valid context is not available in meshconfig", errors.GetLDescription(err))
}
