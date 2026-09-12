package connections

import (
	"github.com/meshery/meshkit/errors"
)

const (
	ErrControllersConfigInvalidCode  = "meshery-server-1437"
	ErrControllersConfigMetadataCode = "meshery-server-1438"
)

// ErrControllersConfigInvalid is returned when a controllers configuration
// document violates a guardrail (replica range, watch-list mutual exclusion,
// broker service coherence, or deployment-mode enum).
func ErrControllersConfigInvalid(reason string) error {
	return errors.New(
		ErrControllersConfigInvalidCode,
		errors.Alert,
		[]string{"Invalid controllers configuration."},
		[]string{reason},
		[]string{"The submitted Meshery Operator / MeshSync / Broker configuration violates a validation rule."},
		[]string{"Correct the highlighted field and resubmit. Replica counts must be 1-10, a watch list sets at most one of whitelist or blacklist, and load-balancer settings require the LoadBalancer service type."},
	)
}

// ErrControllersConfigMetadata is returned when the controllers configuration
// stored on a connection's metadata cannot be encoded or decoded.
func ErrControllersConfigMetadata(err error) error {
	return errors.New(
		ErrControllersConfigMetadataCode,
		errors.Alert,
		[]string{"Unable to read or write the connection's controllers configuration."},
		[]string{err.Error()},
		[]string{"The controllers_config entry on the connection's metadata is not valid JSON or does not match the expected schema."},
		[]string{"Clear the connection's controllers configuration override and reapply it from the Connections page."},
	)
}
