package kubernetes

import (
	"github.com/layer5io/meshkit/errors"
)

const (
	ErrConnectActionCode = "replace_me"
)

func ErrConnectAction(err error) error {
	return ErrConnectActionFromString(err.Error())
}

func ErrConnectActionFromString(message string) error {
	return errors.New(ErrConnectActionCode, errors.Alert, []string{"Error in ConnectAction"}, []string{message}, []string{}, []string{})
}
