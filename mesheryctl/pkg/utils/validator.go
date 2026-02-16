package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// Based on https://github.com/go-playground/validator/blob/dede3413a22993dd5a091707c6764b6e9724df17/regexes.go#L75 + v prefix
var vSemverRegex = regexp.MustCompile(`^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

func IsUUID(value string) bool {
	if _, err := uuid.Parse(value); err != nil {
		return false
	}
	return true
}

func ValidateSemver(fl validator.FieldLevel) bool {
	return vSemverRegex.MatchString(fl.Field().String())
}

func ValidationErrorMessage(structField string, err validator.ValidationErrors) string {

	if len(err) == 0 {
		return ""
	}

	for _, e := range err {
		switch e.Tag() {
		case "semver":
			return fmt.Sprintf("Invalid value '%v' for flag --%s: version must be in format vX.X.X", e.Value(), strings.ToLower(e.Field()))
		case "oneof":
			return fmt.Sprintf("Invalid value '%v' for flag --%s: valid values are %s", e.Value(), strings.ToLower(e.Field()), e.Param())
		case "dir":
			return fmt.Sprintf("Invalid value '%v' for flag --%s: directory does not exist", e.Value(), strings.ToLower(e.Field()))
		default:
			return fmt.Sprintf("Invalid value '%v' for flag --%s", e.Value(), strings.ToLower(e.Field()))
		}
	}
	return ""
}
