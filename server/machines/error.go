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
	ErrConnectionNotFoundCode     = "meshery-server-1444"
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

func ErrConnectionNotFound(id string) error {
	return errors.New(ErrConnectionNotFoundCode, errors.Alert, []string{fmt.Sprintf("No connection found with id \"%s\" while updating its status.", id)}, []string{"The provider returned no connection for the given id without reporting an error."}, []string{"The connection may have been deleted, or the persisted record is missing."}, []string{"Reinitialise the connection, or verify the connection id exists in the provider."})
}
