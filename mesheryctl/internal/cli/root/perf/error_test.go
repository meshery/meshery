package perf

import (
	"fmt"
	"testing"

	"github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrPerformanceProfileResultSeparatesCause(t *testing.T) {
	err := ErrPerformanceProfileResult(fmt.Errorf("profile not found"))

	assert.Equal(t, "Unable to display performance Profile: profile not found", errors.GetLDescription(err))
}
