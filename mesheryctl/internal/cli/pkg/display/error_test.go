package display

import (
	"errors"
	"fmt"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitErrors "github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorListPagination(t *testing.T) {
	tests := []struct {
		name        string
		inputError  error
		currentPage int
		expectAuth  bool // true if we expect the original error to be returned (auth error)
	}{
		{
			name:        "String-based authentication error",
			inputError:  errors.New("authentication failed"),
			currentPage: 0,
			expectAuth:  true,
		},
		{
			name:        "String-based token error",
			inputError:  errors.New("invalid token provided"),
			currentPage: 1,
			expectAuth:  true,
		},
		{
			name:        "String-based unauthorized error",
			inputError:  errors.New("unauthorized access"),
			currentPage: 2,
			expectAuth:  true,
		},
		{
			name:        "String-based login error",
			inputError:  errors.New("please login to continue"),
			currentPage: 3,
			expectAuth:  true,
		},
		{
			name:        "Mixed case authentication error",
			inputError:  errors.New("Authentication Failed"),
			currentPage: 0,
			expectAuth:  true,
		},
		{
			name: "Meshkit unauthenticated error",
			inputError: meshkitErrors.New(utils.ErrUnauthenticatedCode, meshkitErrors.Alert,
				[]string{"Unauthenticated User"},
				[]string{"Access to this resource is unauthorized."},
				[]string{"You haven't logged in to Meshery."},
				[]string{"To proceed, log in using `mesheryctl system login`."}),
			currentPage: 0,
			expectAuth:  true,
		},
		{
			name: "Meshkit invalid token error",
			inputError: meshkitErrors.New(utils.ErrInvalidTokenCode, meshkitErrors.Alert,
				[]string{"Invalid authentication Token"},
				[]string{"The authentication token has expired or is invalid."},
				[]string{"The token in auth.json has expired or is invalid."},
				[]string{"Provide a valid user token by logging in with `mesheryctl system login`."}),
			currentPage: 1,
			expectAuth:  true,
		},
		{
			name: "Meshkit attach auth token error",
			inputError: meshkitErrors.New(utils.ErrAttachAuthTokenCode, meshkitErrors.Alert,
				[]string{"Authentication token Not Found"},
				[]string{"Authentication token not found"},
				[]string{"The user is not logged in to generate a token."},
				[]string{"Log in with `mesheryctl system login` or supply a valid user token using the --token (or -t) flag."}),
			currentPage: 2,
			expectAuth:  true,
		},
		{
			name: "Non-authentication meshkit error",
			inputError: meshkitErrors.New(utils.ErrFailRequestCode, meshkitErrors.Alert,
				[]string{"Failed to make a request"},
				[]string{"Network error occurred"},
				[]string{"Meshery server is not reachable."},
				[]string{"Ensure your Kubernetes cluster is running and your network connection is active."}),
			currentPage: 1,
			expectAuth:  false,
		},
		{
			name:        "Generic network error",
			inputError:  errors.New("connection refused"),
			currentPage: 2,
			expectAuth:  false,
		},
		{
			name:        "Generic server error",
			inputError:  errors.New("internal server error"),
			currentPage: 3,
			expectAuth:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ErrorListPagination(tt.inputError, tt.currentPage)

			if tt.expectAuth {
				// For authentication errors, the original error should be returned unchanged
				assert.Equal(t, tt.inputError, result, "Authentication errors should be returned as-is")
			} else {
				// For non-authentication errors, should be wrapped with pagination context
				assert.NotEqual(t, tt.inputError, result, "Non-authentication errors should be wrapped")
				
				// Check that it's a meshkit error with the correct code
				meshkitErr, ok := result.(*meshkitErrors.Error)
				assert.True(t, ok, "Result should be a meshkit error")
				assert.Equal(t, ErrListPaginationCode, meshkitErr.Code, "Should have pagination error code")
				
				// Check that the error message contains pagination context
				assert.Contains(t, result.Error(), fmt.Sprintf("failed to fetch data for page %d", tt.currentPage),
					"Error should contain pagination context")
			}
		})
	}
}

func TestErrorListPagination_EdgeCases(t *testing.T) {
	t.Run("Empty error message", func(t *testing.T) {
		emptyErr := errors.New("")
		result := ErrorListPagination(emptyErr, 0)
		
		// Empty error should be wrapped since it doesn't contain auth keywords
		assert.NotEqual(t, emptyErr, result)
		meshkitErr, ok := result.(*meshkitErrors.Error)
		assert.True(t, ok)
		assert.Equal(t, ErrListPaginationCode, meshkitErr.Code)
	})

	t.Run("Negative page number", func(t *testing.T) {
		genericErr := errors.New("some error")
		result := ErrorListPagination(genericErr, -1)
		
	     _, ok := result.(*meshkitErrors.Error)
		assert.True(t, ok)
		assert.Contains(t, result.Error(), "failed to fetch data for page -1")
	})

	t.Run("Case insensitive auth detection", func(t *testing.T) {
		testCases := []string{
			"AUTHENTICATION failed",
			"Token expired",
			"UNAUTHORIZED request",
			"Please LOGIN",
		}
		
		for _, errMsg := range testCases {
			result := ErrorListPagination(errors.New(errMsg), 0)
			assert.Equal(t, errors.New(errMsg).Error(), result.Error(), 
				"Case insensitive auth detection should work for: %s", errMsg)
		}
	})
}
