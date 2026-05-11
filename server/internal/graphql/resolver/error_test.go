package resolver

import (
	"errors"
	"testing"

	meshkitErrors "github.com/meshery/meshkit/errors"
)

func TestErrAdapterOperationUsesMeshKitError(t *testing.T) {
	err := ErrAdapterOperation("Deploy", errors.New("adapter deployment failed"))

	if got := meshkitErrors.GetCode(err); got != ErrAdapterOperationCode {
		t.Fatalf("expected MeshKit error code %q, got %q", ErrAdapterOperationCode, got)
	}
}
