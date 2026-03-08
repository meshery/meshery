package mesheryctlflags

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateSemver(t *testing.T) {
	flagValidator := GetFlagValidator()

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
			err := flagValidator.Validate(ts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateBoolean(t *testing.T) {
	flagValidator := GetFlagValidator()

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
			err := flagValidator.Validate(ts)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewValidator(t *testing.T) {
	flagValidator := GetFlagValidator()
	assert.NotNil(t, flagValidator)

	// Test that custom validations are registered
	type testStruct struct {
		Version   string `validate:"semver"`
		BoolField bool   `validate:"boolean"`
	}

	ts := testStruct{Version: "v1.0.0", BoolField: true}
	err := flagValidator.Validate(ts)
	assert.NoError(t, err)

	// Test json tag name function
	type tagTestStruct struct {
		CustomName string `json:"custom_name" validate:"required"`
	}

	tts := tagTestStruct{CustomName: ""}
	err = flagValidator.Validate(tts)
	assert.Error(t, err)
}

func TestReadValidationErrorMessages(t *testing.T) {
	flagValidator := GetFlagValidator()

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
			err := flagValidator.Validate(tt.setupStruct)

			if tt.shouldHaveError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReadValidationErrorMessages_Direct(t *testing.T) {

	fv := GetFlagValidator()

	getValidationErrors := func(t *testing.T, v interface{}) validator.ValidationErrors {
		t.Helper()
		err := fv.Validator.Struct(v)
		assert.Error(t, err)

		verrs, ok := err.(validator.ValidationErrors)
		assert.True(t, ok)
		return verrs
	}

	t.Run("given no custom error exists then falls back to default message", func(t *testing.T) {

		verrs := getValidationErrors(t, struct {
			RequiredField string `json:"required-field" validate:"required"`
		}{RequiredField: ""})

		err := fv.ReadValidationErrorMessages(verrs)
		assert.Error(t, err)
		assert.Equal(t, "Invalid value for --required-field ''", err.Error())
	})

	t.Run("given empty validation errors then returns nil", func(t *testing.T) {

		var empty validator.ValidationErrors
		err := fv.ReadValidationErrorMessages(empty)
		assert.NoError(t, err)
	})

	t.Run("given semver validation error then formats semver error", func(t *testing.T) {

		verrs := getValidationErrors(t, struct {
			Version string `json:"version" validate:"semver"`
		}{Version: "1.2.3"})

		err := fv.ReadValidationErrorMessages(verrs)
		assert.Error(t, err)
		assert.Equal(t, "Invalid value for --version '1.2.3': version must be in format vX.X.X", err.Error())
	})

	t.Run("given oneof validation error then formats oneof error", func(t *testing.T) {

		verrs := getValidationErrors(t, struct {
			Format string `json:"format" validate:"oneof=json yaml"`
		}{Format: "xml"})

		err := fv.ReadValidationErrorMessages(verrs)
		assert.Error(t, err)
		assert.Equal(t, "Invalid value for --format 'xml': valid values are json yaml", err.Error())
	})

	t.Run("given dir validation error then formats dir error", func(t *testing.T) {

		verrs := getValidationErrors(t, struct {
			OutputDir string `json:"output-dir" validate:"dir"`
		}{OutputDir: "/this/path/should/not/exist"})

		err := fv.ReadValidationErrorMessages(verrs)
		assert.Error(t, err)
		assert.Equal(t, "Invalid value for --output-dir '/this/path/should/not/exist': directory does not exist", err.Error())
	})

	t.Run("given custom error exists then uses custom error", func(t *testing.T) {

		fv.CustomErrors["required"] = "required-field must be provided"

		verrs := getValidationErrors(t, struct {
			RequiredField string `json:"required-field" validate:"required"`
		}{RequiredField: ""})

		err := fv.ReadValidationErrorMessages(verrs)
		assert.Error(t, err)
		assert.Equal(t, "required-field must be provided", err.Error())
	})

	t.Run("given multiple validation errors then joins messages with comma", func(t *testing.T) {

		verrs := getValidationErrors(t, struct {
			Version string `json:"version" validate:"semver"`
			Format  string `json:"format" validate:"oneof=json yaml"`
		}{
			Version: "invalid",
			Format:  "xml",
		})

		err := fv.ReadValidationErrorMessages(verrs)
		assert.Error(t, err)
		assert.Equal(
			t,
			"Invalid value for --version 'invalid': version must be in format vX.X.X, Invalid value for --format 'xml': valid values are json yaml",
			err.Error(),
		)
	})

	// Reset custom errors after tests to avoid side effects
	fv.CustomErrors = make(map[string]string)
}
