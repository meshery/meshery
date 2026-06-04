package machines

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrInvalidTransitionCode      = "meshery-server-1214"
	ErrInvalidTransitionEventCode = "meshery-server-1215"
	ErrInititalizeK8sMachineCode  = "meshery-server-1216"
	ErrAssetMachineCtxCode        = "meshery-server-1217"
	ErrInvalidTypeCode            = "meshery-server-1218"
	ErrInitLoggerCode             = "meshery-server-1432"
)

func ErrInvalidTransition(from, to StateType) error {
	return errors.New(ErrInvalidTransitionCode, errors.Alert, []string{fmt.Sprintf("transition restricted from \"%s\" to \"%s\"", from, to)}, []string{}, []string{}, []string{})
}

func ErrInvalidTransitionEvent(from StateType, event EventType) error {
	return errors.New(ErrInvalidTransitionEventCode, errors.Alert, []string{fmt.Sprintf("unsupported transition event received \"%s\" for the state \"%s\"", event, from)}, []string{}, []string{}, []string{})
}

func ErrInititalizeK8sMachine(err error) error {
	return errors.New(ErrInititalizeK8sMachineCode, errors.Alert, []string{"Provided connection id is invalid"}, []string{err.Error()}, []string{"Provided ID is not a valid uuid."}, []string{"Hard delete and reinitialise the connection process."})
}

func ErrAssertMachineCtx(err error) error {
	return errors.New(ErrAssetMachineCtxCode, errors.Alert, []string{"type assertion of context to *kubernetes.MachineCtx failed"}, []string{err.Error()}, []string{"The machine context has become invalid."}, []string{"Retry the opereation or try archiving and reinitializing the connection"})
}

func ErrInvalidType(err error) error {
	return errors.New(ErrInvalidTypeCode, errors.Alert, []string{"Provided connection id is invalid"}, []string{err.Error()}, []string{"Provided ID is not a valid uuid."}, []string{"Hard delete and reinitialise the connection process."})
}

func ErrInitLogger(err error) error {
	return errors.New(ErrInitLoggerCode, errors.Alert, []string{"Failed to initialize Prometheus logger"}, []string{err.Error()}, []string{"Invalid logger configuration or system resources failure"}, []string{"Verify the logger formatting options and system resources"})
}
