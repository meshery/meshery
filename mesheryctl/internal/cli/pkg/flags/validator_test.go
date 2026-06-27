package mesheryctlflags

import (
	"os"
	"path/filepath"
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
			name:    "given valid semver with v prefix when isSemver then returns no error",
			version: "v1.0.0",
			wantErr: false,
		},
		{
			name:    "given valid semver with prerelease when isSemver then returns no error",
			version: "v1.0.0-alpha",
			wantErr: false,
		},
		{
			name:    "given valid semver with build metadata when isSemver then returns no error",
			version: "v1.0.0+20130313144700",
			wantErr: false,
		},
		{
			name:    "given valid semver with prerelease and build when isSemver then returns no error",
			version: "v1.0.0-beta.1+exp.sha.5114f85",
			wantErr: false,
		},
		{
			name:    "given invalid semver without v prefix when isSemver then returns error",
			version: "1.0.0",
			wantErr: true,
		},
		{
			name:    "given invalid semver format when isSemver then returns error",
			version: "v1.0",
			wantErr: true,
		},
		{
			name:    "given invalid semver with leading zeros when isSemver then returns error",
			version: "v01.0.0",
			wantErr: true,
		},
		{
			name:    "given empty string when isSemver then returns error",
			version: "",
			wantErr: true,
		},
		{
			name:    "given invalid characters when isSemver then returns error",
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
			name:    "given true value when isBoolean then returns no error",
			value:   true,
			wantErr: false,
		},
		{
			name:    "given false value when isBoolean then returns no error",
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
			name: "given invalid semver when isSemver then error message contains semver format requirement",
			setupStruct: struct {
				Version string `json:"version" validate:"semver"`
			}{Version: "1.0.0"},
			expectedError:   "Invalid value for --version '1.0.0': version must be in format vX.X.X",
			shouldHaveError: true,
		},
		{
			name: "given invalid oneof value when isOneof then error message contains valid values",
			setupStruct: struct {
				Format string `json:"format" validate:"oneof=json yaml"`
			}{Format: "xml"},
			expectedError:   "Invalid value for --format 'xml': valid values are json yaml",
			shouldHaveError: true,
		},
		{
			name: "given invalid directory when isDir then error message indicates directory does not exist",
			setupStruct: struct {
				OutputDir string `json:"output-dir" validate:"dir"`
			}{OutputDir: "/nonexistent/path"},
			expectedError:   "Invalid value for --output-dir '/nonexistent/path': directory does not exist",
			shouldHaveError: true,
		},
		{
			name: "given invalid directory path when isDirPath then error message indicates directory does not exist",
			setupStruct: struct {
				SourcePath string `json:"source-path" validate:"dirpath"`
			}{SourcePath: "/invalid/directory"},
			expectedError:   "Invalid value for --source-path '/invalid/directory': directory does not exist",
			shouldHaveError: true,
		},
		{
			name: "given empty required field when isRequired then error message indicates field is required",
			setupStruct: struct {
				RequiredField string `json:"required-field" validate:"required"`
			}{RequiredField: ""},
			expectedError:   "Invalid value for --required-field ''",
			shouldHaveError: true,
		},
		{
			name: "given multiple validation errors when multiple fields are invalid then error message contains all errors",
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

	t.Run("given unknown validation tag then falls back to generic message", func(t *testing.T) {
		verrs := getValidationErrors(t, struct {
			Field string `json:"field" validate:"cidr"`
		}{Field: "value"})

		err := fv.ReadValidationErrorMessages(verrs)
		assert.Error(t, err)
		assert.Equal(t, "Invalid value for --field 'value'", err.Error())
	})

	t.Run("given relpath validation error then formats relpath error", func(t *testing.T) {
		verrs := getValidationErrors(t, struct {
			Path string `json:"path" validate:"relpath"`
		}{Path: "/not-a-real-path"})

		err := fv.ReadValidationErrorMessages(verrs)
		assert.Error(t, err)
		assert.Equal(t, "Invalid value for --path '/not-a-real-path': path must be a relative path", err.Error())
	})

	t.Run("given abspath validation error then formats abspath error", func(t *testing.T) {
		verrs := getValidationErrors(t, struct {
			Path string `json:"path" validate:"abspath"`
		}{Path: "://not-a-real-path"})

		err := fv.ReadValidationErrorMessages(verrs)
		assert.Error(t, err)
		assert.Equal(t, "Invalid value for --path '://not-a-real-path': path must be an absolute path", err.Error())
	})

	// Reset custom errors after tests to avoid side effects
	fv.CustomErrors = make(map[string]string)
}

func TestValidateRelativePath(t *testing.T) {
	flagValidator := GetFlagValidator()

	cwd, err := os.Getwd()
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	absolutePath := filepath.Join(cwd, "testdata", "meshery")

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "given empty path when isRelPath then returns no error",
			path:    "",
			wantErr: false,
		},
		{
			name:    "given current relative path when isRelPath then returns no error",
			path:    ".",
			wantErr: false,
		},
		{
			name:    "given simple relative path when isRelPath then returns no error",
			path:    "configs/meshery.yaml",
			wantErr: false,
		},
		{
			name:    "given dot relative path when isRelPath then returns no error",
			path:    "./configs/meshery.yaml",
			wantErr: false,
		},
		{
			name:    "given parent relative path when isRelPath then returns no error",
			path:    "../configs/meshery.yaml",
			wantErr: false,
		},
		{
			name:    "given absolute path when isRelPath then returns error",
			path:    absolutePath,
			wantErr: true,
		},
	}

	// Add a Windows-style absolute path case when running on Windows.
	if vol := filepath.VolumeName(cwd); vol != "" {
		tests = append(tests, struct {
			name    string
			path    string
			wantErr bool
		}{
			name:    "given windows absolute path when isRelPath then returns error",
			path:    vol + `\meshery\config.yaml`,
			wantErr: true,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type testStruct struct {
				Path string `validate:"relpath"`
			}

			err := flagValidator.Validate(testStruct{Path: tt.path})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidateRelativeOrAbsolutePath(t *testing.T) {
	flagValidator := GetFlagValidator()

	cwd, err := os.Getwd()
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	absolutePath := filepath.Join(cwd, "testdata", "meshery")

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "given empty path when isRelOrAbsPath then returns no error",
			path:    "",
			wantErr: false,
		},
		{
			name:    "given current relative path when isRelOrAbsPath then returns no error",
			path:    ".",
			wantErr: false,
		},
		{
			name:    "given simple relative path when isRelOrAbsPath then returns no error",
			path:    "configs/meshery.yaml",
			wantErr: false,
		},
		{
			name:    "given dot relative path when isRelOrAbsPath then returns no error",
			path:    "./configs/meshery.yaml",
			wantErr: false,
		},
		{
			name:    "given parent relative path when isRelOrAbsPath then returns no error",
			path:    "../configs/meshery.yaml",
			wantErr: false,
		},
		{
			name:    "given absolute path when isRelOrAbsPath then returns no error",
			path:    absolutePath,
			wantErr: false,
		},
	}

	// Add a Windows-style absolute path case when running on Windows.
	if vol := filepath.VolumeName(cwd); vol != "" {
		tests = append(tests, struct {
			name    string
			path    string
			wantErr bool
		}{
			name:    "given windows absolute path when isRelOrAbsPath then returns no error",
			path:    vol + `\meshery\config.yaml`,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type testStruct struct {
				Path string `validate:"relabspath"`
			}

			err := flagValidator.Validate(testStruct{Path: tt.path})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestValidateAbsolutePath(t *testing.T) {
	flagValidator := GetFlagValidator()

	cwd, err := os.Getwd()
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	absolutePath := filepath.Join(cwd, "testdata", "meshery")

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "given empty path when isAbsPath then returns no error",
			path:    "",
			wantErr: false,
		},
		{
			name:    "given absolute path when isAbsPath then returns no error",
			path:    absolutePath,
			wantErr: false,
		},
		{
			name:    "given current relative path when isAbsPath then returns error",
			path:    ".",
			wantErr: true,
		},
		{
			name:    "given simple relative path when isAbsPath then returns error",
			path:    "configs/meshery.yaml",
			wantErr: true,
		},
		{
			name:    "given dot relative path when isAbsPath then returns error",
			path:    "./configs/meshery.yaml",
			wantErr: true,
		},
		{
			name:    "given parent relative path when isAbsPath then returns error",
			path:    "../configs/meshery.yaml",
			wantErr: true,
		},
	}

	// Add a Windows-style absolute path case when running on Windows.
	if vol := filepath.VolumeName(cwd); vol != "" {
		tests = append(tests, struct {
			name    string
			path    string
			wantErr bool
		}{
			name:    "given windows absolute path when isAbsPath then returns no error",
			path:    vol + `\meshery\config.yaml`,
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			type testStruct struct {
				Path string `validate:"abspath"`
			}

			err := flagValidator.Validate(testStruct{Path: tt.path})
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
