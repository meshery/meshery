package sessions

import (
	"strconv"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrNoDriverCode           = "meshery-server-1443"
	ErrUnsupportedSessionCode = "meshery-server-1444"
	ErrTargetNotFoundCode     = "meshery-server-1445"
	ErrInvalidTargetCode      = "meshery-server-1446"
	ErrOpenSessionCode        = "meshery-server-1447"
	ErrSessionLimitCode       = "meshery-server-1448"
	ErrUpgradeCode            = "meshery-server-1449"
	ErrDriverInitCode         = "meshery-server-1450"
	ErrSlowConsumerCode       = "meshery-server-1451"
)

// ErrSlowConsumer is returned when a client stops draining its socket and the
// server's outbound buffer backs up. The session is torn down rather than
// allowed to pin a goroutine and an upstream stream indefinitely.
func ErrSlowConsumer() error {
	return errors.New(
		ErrSlowConsumerCode, errors.None,
		[]string{"Session closed because the client stopped reading"},
		[]string{"Output from the session backed up beyond the server's outbound buffer because the client did not read it, so the session was terminated."},
		[]string{"The browser tab holding the session was suspended or throttled.", "The network path to the client stalled.", "The session produced output faster than the client could consume it."},
		[]string{"Reopen the session.", "If the target emits a very high volume of output, narrow the log stream with a tail limit or a time bound."},
	)
}

// ErrNoDriver is returned when a session is requested against a connection
// whose kind has no registered driver. Emitted with HTTP 400.
func ErrNoDriver(kind string) error {
	return errors.New(
		ErrNoDriverCode, errors.Alert,
		[]string{"Sessions are not supported for this connection kind"},
		[]string{"No session driver is registered for connections of kind '" + kind + "', so terminal and log sessions cannot be opened against its resources."},
		[]string{"The connection's kind does not implement terminal or log sessions.", "A driver for this kind exists but was not registered at startup."},
		[]string{"Open the session against a connection whose kind supports sessions, e.g. a Kubernetes connection.", "If this kind should support sessions, register its driver with the sessions registry."},
	)
}

// ErrUnsupportedSession is returned when a driver serves the connection kind
// but cannot serve the requested session kind for the requested target, e.g. a
// terminal against a resource that has no executable process. Emitted with
// HTTP 400.
func ErrUnsupportedSession(sessionKind, resource, reason string) error {
	ldesc := "A '" + sessionKind + "' session is not available for resource type '" + resource + "'."
	if reason != "" {
		ldesc += " " + reason
	}
	return errors.New(
		ErrUnsupportedSessionCode, errors.Alert,
		[]string{"Session kind is not supported for this target"},
		[]string{ldesc},
		[]string{"The resource type does not support this session kind.", "The resource is not in a state that admits this session kind, e.g. a pod that is not running."},
		[]string{"Request a session kind the target supports; query the capabilities endpoint for the target to discover which kinds are available."},
	)
}

// ErrTargetNotFound is returned when the addressed resource does not exist on
// the remote end. Emitted with HTTP 404.
func ErrTargetNotFound(err error, target string) error {
	return errors.New(
		ErrTargetNotFoundCode, errors.Alert,
		[]string{"Session target not found"},
		[]string{"The resource '" + target + "' could not be found: " + err.Error()},
		[]string{"The resource was deleted or renamed after the client last listed it.", "The namespace or name in the request is misspelled.", "The connection's credentials do not grant visibility of the resource."},
		[]string{"Refresh the resource list and retry against a resource that still exists.", "Verify the namespace and name, and that the connection's credentials can read the resource."},
	)
}

// ErrInvalidTarget is returned when the target reference is malformed or is
// missing a field the driver requires. Emitted with HTTP 400.
func ErrInvalidTarget(reason string) error {
	return errors.New(
		ErrInvalidTargetCode, errors.Alert,
		[]string{"Invalid session target"},
		[]string{reason},
		[]string{"A required query parameter such as `resource` or `name` was omitted.", "A URL template did not get its parameter substituted."},
		[]string{"Supply `resource` and `name` (and `namespace` for namespaced resources) as query parameters on the session request."},
	)
}

// ErrOpenSession wraps a failure to establish or sustain the stream to the
// remote target. It is reported to an already-upgraded client as an error
// control frame, and with HTTP 500 otherwise.
func ErrOpenSession(err error, sessionKind string) error {
	return errors.New(
		ErrOpenSessionCode, errors.Alert,
		[]string{"Failed to open " + sessionKind + " session"},
		[]string{err.Error()},
		[]string{"The remote API server rejected the streaming upgrade.", "The connection's credentials lack the RBAC permission the session requires, e.g. `pods/exec` or `pods/log`.", "The target became unreachable while the session was being established."},
		[]string{"Confirm the cluster is reachable and that the connection's credentials are authorized for the subresource this session uses.", "Check the Meshery server logs for the underlying transport error."},
	)
}

// ErrSessionLimit is returned when a user already holds the maximum number of
// concurrent sessions. Emitted with HTTP 429.
func ErrSessionLimit(limit int) error {
	return errors.New(
		ErrSessionLimitCode, errors.Alert,
		[]string{"Too many concurrent sessions"},
		[]string{"This user already holds the maximum of concurrently open terminal and log sessions."},
		[]string{"Sessions were opened faster than they were closed.", "Browser tabs holding sessions were left open."},
		[]string{"Close an existing terminal or log session and retry.", "The per-user concurrent session limit is " + strconv.Itoa(limit) + "."},
	)
}

// ErrUpgrade is returned when the WebSocket handshake fails. The upgrader has
// already written an HTTP error response by the time this is constructed, so it
// is logged rather than surfaced on the wire.
func ErrUpgrade(err error) error {
	return errors.New(
		ErrUpgradeCode, errors.Alert,
		[]string{"Failed to upgrade connection to a WebSocket"},
		[]string{err.Error()},
		[]string{"The client is not speaking the WebSocket protocol.", "The request's Origin header does not match the server's host and is not in the allowed-origins list.", "An intermediate proxy stripped the upgrade headers."},
		[]string{"Open the session from the Meshery UI, which is served from an allowed origin.", "If Meshery is behind a proxy, configure it to forward `Upgrade` and `Connection` headers."},
	)
}

// ErrDriverInit is returned when a driver cannot be constructed for a
// connection, e.g. its credentials cannot be resolved. Emitted with HTTP 500.
func ErrDriverInit(err error, kind string) error {
	return errors.New(
		ErrDriverInitCode, errors.Alert,
		[]string{"Failed to initialize the session driver"},
		[]string{"A session driver for connection kind '" + kind + "' could not be built: " + err.Error()},
		[]string{"The connection has no credential associated with it.", "The stored credential is malformed or has been revoked.", "The remote provider could not be reached to fetch the credential."},
		[]string{"Re-register the connection so a fresh credential is stored.", "Verify the connection is in a connected state before opening a session against it."},
	)
}
