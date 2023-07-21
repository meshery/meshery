package config

import "github.com/layer5io/meshkit/errors"

const (
	ErrInvalidConfigCode = "1000"
)

func ErrInvalidConfig(err error) error {
	return errors.New(ErrInvalidConfigCode, errors.Alert, []string{"Invalid Meshconfig"}, []string{err.Error()}, []string{"MeshConfig File is Invalid"}, []string{})
}
