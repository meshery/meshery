package mesheryctlflags

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

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
	validator *validator.Validate
}

// Based on https://github.com/go-playground/validator/blob/dede3413a22993dd5a091707c6764b6e9724df17/regexes.go#L75 adding `v` prefix
var vSemverRegex = regexp.MustCompile(`^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

func (fv *FlagValidator) Validate(s interface{}) error {
	err := fv.validator.Struct(s)
	if err != nil {
		switch vErr := err.(type) {
		case *validator.InvalidValidationError:
			return utils.ErrFlagsInvalid(err)
		case validator.ValidationErrors:
			return utils.ErrFlagsInvalid(ReadValidationErrorMessages(vErr))
		default:
			return utils.ErrFlagsInvalid(err)
		}
	}
	return nil
}

func validateSemver(fl validator.FieldLevel) bool {
	return vSemverRegex.MatchString(fl.Field().String())
}

// validateBoolean is a custom validation function that checks if a field is a boolean value (true or false)
// This is necessary because the default validator does not have a built-in validation behaving as expected for boolean fields,
// especially when using flags in cobra
func validateBoolean(fl validator.FieldLevel) bool {
	if _, ok := fl.Field().Interface().(bool); ok {
		return true
	}
	return false
}

// ReadValidationErrorMessages reads the validation error and returns a slice of error messages for each validation error encountered
// This is a centralized function to read validation error messages for all the commands in mesheryctl and return user friendly error messages based on the type of validation error encountered
func ReadValidationErrorMessages(err validator.ValidationErrors) error {

	if len(err) == 0 {
		return nil
	}

	errorMessages := make([]string, 0, len(err))
	for _, e := range err {
		switch e.Tag() {
		case "semver":
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v': version must be in format vX.X.X", strings.ToLower(e.Field()), e.Value()))
		case "oneof":
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v': valid values are %s", strings.ToLower(e.Field()), e.Value(), e.Param()))
		case "dir", "dirpath":
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v': directory does not exist", strings.ToLower(e.Field()), e.Value()))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("Invalid value for --%s '%v'", strings.ToLower(e.Field()), e.Value()))
		}
	}

	return errors.New(strings.Join(errorMessages, ", "))
}

func NewFlagValidator() *FlagValidator {
	validate := validator.New()

	// Register the custom validation function with a tag name "semver"
	err := validate.RegisterValidation("semver", validateSemver)
	if err != nil {
		log.Fatalf("Error registering validation: %v", err)
	}

	// Register a custom validation function for boolean values that accepts "true", "false"
	err = validate.RegisterValidation("boolean", validateBoolean)
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

	return &FlagValidator{validator: validate}
}

func InitValidators(cmd *cobra.Command) {
	validate := NewFlagValidator()
	ctx := context.WithValue(context.Background(), FlagValidatorKey, validate)
	cmd.SetContext(ctx)
}
