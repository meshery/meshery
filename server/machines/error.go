package machines

import (
	"github.com/layer5io/meshkit/errors"
)

var ErrInvalidTransitionCode = "replace_me"
var ErrInititalizeK8sMachineCode = "replace_me"
var ErrAssetMachineCtxCode = "replace_me"

var ErrInvalidTransition = errors.New(ErrInvalidTransitionCode, errors.Alert, []string{}, []string{}, []string{}, []string{})


func ErrInititalizeK8sMachine(err error) error {
	return errors.New(ErrInititalizeK8sMachineCode, errors.Alert, []string{"Provided connection id is invalid"}, []string{err.Error()}, []string{"Provided ID is not a valid uuid."}, []string{"Hard delete and reinitialise the connection process."})
} 

func ErrAssertMachineCtx(err error) error {
	return errors.New(ErrAssetMachineCtxCode, errors.Alert, []string{"type assertion of context to *kubernetes.MachineCtx failed"}, []string{err.Error()}, []string{"The machine context has become invalid."}, []string{"Retry the opereation or try archiving and reinitializing the connection"})
} 