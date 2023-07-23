package config

import "github.com/layer5io/meshkit/errors"

const (
	ErrInvalidMeshConfigCode = "replace_me"
)

func ErrInvalidMeshConfig(err error) error {
	return errors.New(ErrInvalidMeshConfigCode, errors.Alert, []string{"Invalid Meshconfig"},
		[]string{err.Error()}, []string{"Loading the Invalid MeshConfig data caused error"}, []string{"Make sure that your 'config.yaml' file in your '.meshery' is valid"})
}
