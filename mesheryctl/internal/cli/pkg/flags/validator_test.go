package mesheryctlflags

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateSemver(t *testing.T) {
	validate := NewValidator()

	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "valid semver with v prefix",
			version: "v1.0.0",
			wantErr: false,
		},
		{
			name:    "valid semver with prerelease",
			version: "v1.0.0-alpha",
			wantErr: false,
		},
		{
			name:    "valid semver with build metadata",
			version: "v1.0.0+20130313144700",
			wantErr: false,
		},
		{
			name:    "valid semver with prerelease and build",
			version: "v1.0.0-beta.1+exp.sha.5114f85",
			wantErr: false,
		},
		{
			name:    "invalid semver without v prefix",
			version: "1.0.0",
			wantErr: true,
		},
		{
			name:    "invalid semver format",
			version: "v1.0",
			wantErr: true,
		},
		{
			name:    "invalid semver with leading zeros",
			version: "v01.0.0",
			wantErr: true,
		},
		{
			name:    "empty string",
			version: "",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			version: "vX.Y.Z",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type testStruct struct {
				Version string `validate:"semver"`
			}
			ts := testStruct{Version: tt.version}
			err := validate.Struct(ts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateBoolean(t *testing.T) {
	validate := NewValidator()

	tests := []struct {
		name    string
		value   bool
		wantErr bool
	}{
		{
			name:    "true value",
			value:   true,
			wantErr: false,
		},
		{
			name:    "false value",
			value:   false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type testStruct struct {
				BoolField bool `validate:"boolean"`
			}
			ts := testStruct{BoolField: tt.value}
			err := validate.Struct(ts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadValidationErrorMessages_Empty(t *testing.T) {
	var emptyErrors validator.ValidationErrors
	err := ReadValidationErrorMessages(emptyErrors)
	assert.Nil(t, err)
}

func TestNewValidator(t *testing.T) {
	validate := NewValidator()
	assert.NotNil(t, validate)

	// Test that custom validations are registered
	type testStruct struct {
		Version   string `validate:"semver"`
		BoolField bool   `validate:"boolean"`
	}

	ts := testStruct{Version: "v1.0.0", BoolField: true}
	err := validate.Struct(ts)
	assert.NoError(t, err)

	// Test json tag name function
	type tagTestStruct struct {
		CustomName string `json:"custom_name" validate:"required"`
	}

	tts := tagTestStruct{CustomName: ""}
	err = validate.Struct(tts)
	assert.Error(t, err)

	validationErrors, ok := err.(validator.ValidationErrors)
	assert.True(t, ok)
	assert.Equal(t, "custom_name", validationErrors[0].Field())
}

func TestReadValidationErrorMessages(t *testing.T) {
	validate := NewValidator()

	tests := []struct {
		name            string
		setupStruct     interface{}
		expectedError   string
		shouldHaveError bool
	}{
		{
			name: "semver validation error",
			setupStruct: struct {
				Version string `json:"version" validate:"semver"`
			}{Version: "1.0.0"},
			expectedError:   "Invalid value for --version '1.0.0': version must be in format vX.X.X",
			shouldHaveError: true,
		},
		{
			name: "oneof validation error",
			setupStruct: struct {
				Format string `json:"format" validate:"oneof=json yaml"`
			}{Format: "xml"},
			expectedError:   "Invalid value for --format 'xml': valid values are json yaml",
			shouldHaveError: true,
		},
		{
			name: "dir validation error",
			setupStruct: struct {
				OutputDir string `json:"output-dir" validate:"dir"`
			}{OutputDir: "/nonexistent/path"},
			expectedError:   "Invalid value for --output-dir '/nonexistent/path': directory does not exist",
			shouldHaveError: true,
		},
		{
			name: "dirpath validation error",
			setupStruct: struct {
				SourcePath string `json:"source-path" validate:"dirpath"`
			}{SourcePath: "/invalid/directory"},
			expectedError:   "Invalid value for --source-path '/invalid/directory': directory does not exist",
			shouldHaveError: true,
		},
		{
			name: "default validation error",
			setupStruct: struct {
				RequiredField string `json:"required-field" validate:"required"`
			}{RequiredField: ""},
			expectedError:   "Invalid value for --required-field ''",
			shouldHaveError: true,
		},
		{
			name: "multiple validation errors",
			setupStruct: struct {
				Version string `json:"version" validate:"semver"`
				Format  string `json:"format" validate:"oneof=json yaml"`
			}{Version: "invalid", Format: "xml"},
			expectedError:   "Invalid value for --version 'invalid': version must be in format vX.X.X, Invalid value for --format 'xml': valid values are json yaml",
			shouldHaveError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.setupStruct)

			if tt.shouldHaveError {
				assert.Error(t, err)
				validationErrors, ok := err.(validator.ValidationErrors)
				assert.True(t, ok)

				resultErr := ReadValidationErrorMessages(validationErrors)
				assert.Error(t, resultErr)
				assert.Equal(t, tt.expectedError, resultErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
