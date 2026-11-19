package main

import (
	"errors"
	"strings"
	"testing"

	meshkitErrors "github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrListenAndServe(t *testing.T) {
	origErr := errors.New("bind: address already in use")
	err := ErrListenAndServe(origErr)

	assert.NotNil(t, err)
	assert.Equal(t, ErrListenAndServeCode, meshkitErrors.GetCode(err))
	assert.Equal(t, meshkitErrors.Severity(meshkitErrors.Fatal), meshkitErrors.GetSeverity(err))
	assert.Equal(t, "ListenAndServe Error", meshkitErrors.GetSDescription(err))
	assert.Contains(t, meshkitErrors.GetLDescription(err), "bind: address already in use")

	probableCauses := meshkitErrors.GetCause(err)
	assert.NotEmpty(t, probableCauses, "probable causes must not be empty (#23956)")
	assert.True(t, strings.Contains(probableCauses, "port is already in use") || strings.Contains(probableCauses, "configured port"), "expected probable causes to mention port in use")

	remediations := meshkitErrors.GetRemedy(err)
	assert.NotEmpty(t, remediations, "remediations must not be empty (#23956)")
	assert.True(t, strings.Contains(remediations, "PORT") || strings.Contains(remediations, "conflicting process"), "expected remediations to guide resolving port conflict")
}
