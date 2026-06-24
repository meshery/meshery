package mesheryctlflags

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	// FlagValidatorKey is the key used to store the flag validator in the command context
	FlagValidatorKey pkg.ContextKey = "flags-validator"
)

type FlagValidator struct {
	Validator *validator.Validate
	// CustomErrors is a map of custom error messages for specific validation tags.
	// The key is the validation tag and the value is the error message to be returned when that validation fails
	// It allows us to return more user friendly error messages for specific validation errors instead of the default error messages returned by the validator package.
	CustomErrors map[string]string
}

var flagValidator *FlagValidator
var once sync.Once

// Based on https://github.com/go-playground/validator/blob/dede3413a22993dd5a091707c6764b6e9724df17/regexes.go#L75 adding `v` prefix
var vSemverRegex = regexp.MustCompile(`^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

func (fv *FlagValidator) Validate(s interface{}) error {
	err := fv.Validator.Struct(s)
	if err != nil {
		switch vErr := err.(type) {
		case *validator.InvalidValidationError:
			return utils.ErrFlagsInvalid(err)
		case validator.ValidationErrors:
			return utils.ErrFlagsInvalid(fv.ReadValidationErrorMessages(vErr))
		default:
			return utils.ErrFlagsInvalid(err)
		}
	}
	return nil
}

func isSemver(fl validator.FieldLevel) bool {
	return vSemverRegex.MatchString(fl.Field().String())
}

func isRelativePath(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if path == "" {
		return true // Omit empty check should be handled by 'omitempty' tag
	}
	return !filepath.IsAbs(path)
}

func isAbsolutePath(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if path == "" {
		return true // Omit empty check should be handled by 'omitempty' tag
	}
	return filepath.IsAbs(path)
}

func isRelativeOrAbsolutePath(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	if path == "" {
		return true // Omit empty check should be handled by 'omitempty' tag
	}
	return isRelativePath(fl) || isAbsolutePath(fl)
}

// isBoolean is a custom validation function that checks if a field is a boolean value (true or false)
// This is necessary because the default validator does not have a built-in validation behaving as expected for boolean fields,
// especially when using flags in cobra
func isBoolean(fl validator.FieldLevel) bool {
	if _, ok := fl.Field().Interface().(bool); ok {
		return true
	}
	return false
}

// ReadValidationErrorMessages reads the validation error and returns a slice of error messages for each validation error encountered
// This is a centralized function to read validation error messages for all the commands in mesheryctl and return user friendly error messages based on the type of validation error encountered
func (fv *FlagValidator) ReadValidationErrorMessages(err validator.ValidationErrors) error {
	if len(err) == 0 {
		return nil
	}

	errorMessages := make([]string, 0, len(err))
	for _, e := range err {
		switch e.Tag() {
		case "semver":
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v': version must be in format vX.X.X", e.Field(), e.Value()))
		case "oneof":
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v': valid values are %s", strings.ToLower(e.Field()), e.Value(), e.Param()))
		case "uuid":
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v': must be a valid UUID", strings.ToLower(e.Field()), e.Value()))
		case "dir", "dirpath":
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v': directory does not exist", e.Field(), e.Value()))
		case "relpath":
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v': path must be a relative path", e.Field(), e.Value()))
		case "abspath":
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v': path must be an absolute path", e.Field(), e.Value()))
		case "relabspath":
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v': path must be a relative or absolute path", e.Field(), e.Value()))
		default:
			if customErr, exists := fv.CustomErrors[e.Tag()]; exists {
				errorMessages = append(errorMessages, customErr)
			} else {
				errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v'", e.Field(), e.Value()))
			}
		}
	}

	return errors.New(strings.Join(errorMessages, ", "))
}

func GetFlagValidator() *FlagValidator {

	once.Do(func() {
		validate := validator.New()

		// Register the custom validation function with a tag name "semver"
		err := validate.RegisterValidation("semver", isSemver)
		if err != nil {
			log.Fatalf("Error registering validation: %v", err)
		}

		// Register a custom validation function for boolean values that accepts "true", "false" with tag name "boolean"
		err = validate.RegisterValidation("boolean", isBoolean)
		if err != nil {
			log.Fatalf("Error registering validation: %v", err)
		}

		// Register a custom validation function for relative file paths with a tag name "relpath"
		err = validate.RegisterValidation("relpath", isRelativePath)
		if err != nil {
			log.Fatalf("Error registering validation: %v", err)
		}

		// Register a custom validation function for absolute file paths with a tag name "abspath"
		err = validate.RegisterValidation("abspath", isAbsolutePath)
		if err != nil {
			log.Fatalf("Error registering validation: %v", err)
		}

		// Register a custom validation function for relative or absolute file paths with a tag name "relabspath"
		err = validate.RegisterValidation("relabspath", isRelativeOrAbsolutePath)
		if err != nil {
			log.Fatalf("Error registering validation: %v", err)
		}

		// Register a custom function to use the json tags as field names in errors
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		flagValidator = &FlagValidator{Validator: validate, CustomErrors: make(map[string]string)}
	})

	return flagValidator
}

func InitValidators(cmd *cobra.Command) {
	flagValidator := GetFlagValidator()
	ctx := context.WithValue(context.Background(), FlagValidatorKey, flagValidator)
	cmd.SetContext(ctx)
}

func ValidateCmdFlags[T any](cmd *cobra.Command, cmdFlags *T) error {
	flagValidator, ok := cmd.Context().Value(FlagValidatorKey).(*FlagValidator)
	if !ok || flagValidator == nil {
		return utils.ErrCommandContextMissing(string(FlagValidatorKey))
	}

	return flagValidator.Validate(cmdFlags)
}
