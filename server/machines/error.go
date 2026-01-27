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
	ErrMissingUserContextCode     = "meshery-server-1219"
	ErrMissingSystemIDContextCode = "meshery-server-1220"
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

func ErrMissingUserContext() error {
	return errors.New(ErrMissingUserContextCode, errors.Critical, []string{"User context is missing or invalid"}, []string{"Failed to extract user information from context"}, []string{"The request context does not contain valid user information.", "This can occur when a background operation runs after the HTTP request has completed.", "The user session may have expired."}, []string{"Ensure the operation is initiated through a properly authenticated request.", "Check if the user session is still valid.", "Retry the operation after re-authenticating."})
}

func ErrMissingSystemIDContext() error {
	return errors.New(ErrMissingSystemIDContextCode, errors.Critical, []string{"System ID context is missing or invalid"}, []string{"Failed to extract system ID from context"}, []string{"The request context does not contain a valid system ID.", "This can occur during server initialization or when context propagation fails."}, []string{"Ensure the server is properly initialized.", "Check the middleware chain is correctly configured.", "Retry the operation."})
}
