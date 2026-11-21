package utils

import (
	stdErrors "errors"
	"fmt"
	"testing"

	"github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorDescriptionsSeparateCause(t *testing.T) {
	cause := fmt.Errorf("connection refused")

	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"ErrLoadConfig", ErrLoadConfig(cause), "Error processing config: connection refused"},
		{"ErrParseGithubFile", ErrParseGithubFile(cause, "https://github.com/meshery/meshery"), "Failed to parse github file: connection refused"},
		{"ErrRequestResponse", ErrRequestResponse(cause), "Unable to create a response from request: connection refused"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, errors.GetLDescription(tt.err))
		})
	}
}

func TestTokenFileMeshKitErrorsUnwrapAndCode(t *testing.T) {
	nativeErr := stdErrors.New("access denied: 0x5")

	t.Run("ErrAcquireTokenFileLock", func(t *testing.T) {
		err := ErrAcquireTokenFileLock(nativeErr)
		assert.Equal(t, ErrAcquireTokenFileLockCode, errors.GetCode(err))
		assert.True(t, stdErrors.Is(err, nativeErr), "errors.Is must identify the underlying native error")
		var targetMesh *errors.Error
		assert.True(t, stdErrors.As(err, &targetMesh), "errors.As must extract the MeshKit *errors.Error")
		assert.Equal(t, ErrAcquireTokenFileLockCode, targetMesh.Code)
	})

	t.Run("ErrReplaceTokenFile", func(t *testing.T) {
		err := ErrReplaceTokenFile(nativeErr)
		assert.Equal(t, ErrReplaceTokenFileCode, errors.GetCode(err))
		assert.True(t, stdErrors.Is(err, nativeErr), "errors.Is must identify the underlying native error")
		var targetMesh *errors.Error
		assert.True(t, stdErrors.As(err, &targetMesh), "errors.As must extract the MeshKit *errors.Error")
		assert.Equal(t, ErrReplaceTokenFileCode, targetMesh.Code)
	})

	t.Run("ErrVerifyTokenFileACL", func(t *testing.T) {
		err := ErrVerifyTokenFileACL(`C:\Users\test\.meshery\auth.json`, 0x1200a9, 0x001f01ff, nativeErr)
		assert.Equal(t, ErrVerifyTokenFileACLCode, errors.GetCode(err))
		assert.True(t, stdErrors.Is(err, nativeErr), "errors.Is must identify the underlying native error")
		var targetMesh *errors.Error
		assert.True(t, stdErrors.As(err, &targetMesh), "errors.As must extract the MeshKit *errors.Error")
		assert.Equal(t, ErrVerifyTokenFileACLCode, targetMesh.Code)
		assert.Contains(t, err.Error(), "auth.json")
		assert.Contains(t, err.Error(), "0x1200a9")
		assert.Contains(t, err.Error(), "0x1f01ff")
	})
}
