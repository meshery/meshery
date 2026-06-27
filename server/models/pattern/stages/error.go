package stages

import "github.com/meshery/meshkit/errors"

const (
	ErrResolveReferenceCode = "meshery-server-1361"
	ErrYAMLUnmarshalCode    = "meshery-server-1436"
)

func ErrResolveReference(err error) error {
	return errors.New(ErrResolveReferenceCode, errors.Alert, []string{}, []string{err.Error()}, []string{}, []string{})
}

func ErrYAMLUnmarshal(err error) error {
	return errors.New(ErrYAMLUnmarshalCode, errors.Alert, []string{"failed to parse YAML configuration"}, []string{err.Error()}, []string{"malformed YAML in component configuration"}, []string{"verify the YAML syntax in your design configuration"})
}
