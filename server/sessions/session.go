// Package sessions provides long-lived, interactive streaming sessions
// (interactive terminals and log tails) against resources reachable through a
// Meshery connection.
//
// The package is deliberately transport- and provider-agnostic so that support
// for a new resource universe is additive:
//
//   - A [Driver] knows how to open sessions against the resources of exactly
//     one connection kind ("kubernetes" today; a container runtime, a managed
//     database, or a remote host tomorrow).
//   - A [Registry] maps a connection's Kind onto the driver that serves it, so
//     teaching Meshery a new kind means registering a [Factory] rather than
//     touching the HTTP layer.
//   - The WebSocket transport (websocket.go) speaks the driver-independent
//     frame protocol in protocol.go and never imports Kubernetes.
//
// Sessions are always addressed through a connection, never through the
// deprecated k8scontext: the connection is the credential-bearing construct,
// and it is what makes the abstraction portable to non-Kubernetes targets.
//
// The wire types below are aliases for the generated models of the
// `v1beta1/session` construct in github.com/meshery/schemas, which is the single
// source of truth for the contract this package serves. They are aliased rather
// than re-declared so that a schema change is a compile error here rather than a
// silent divergence, and so a caller may pass a schemas value straight through.
package sessions

import (
	"context"
	"io"

	"github.com/meshery/schemas/models/v1beta1/session"
)

// Kind enumerates the session types a driver may support for a given target.
type Kind = session.SessionKind

const (
	// KindTerminal is a bidirectional, TTY-backed interactive shell.
	KindTerminal = session.Terminal
	// KindLogs is a unidirectional, server-to-client stream of log output.
	KindLogs = session.Logs
)

// ValidKind reports whether k is a session kind this package understands. It is
// a function rather than a method because Kind aliases a generated type.
func ValidKind(k Kind) bool {
	return k == KindTerminal || k == KindLogs
}

// TargetRef addresses the resource a session attaches to, within the universe
// of a single connection.
//
// The field names read Kubernetes-ish because that is the first driver, but
// they are deliberately generic and each driver interprets them for itself:
// Namespace is any parent scope, Container is any sub-target within Name. A
// driver ignores the fields it has no use for, and rejects a target it cannot
// address via [Driver.Capabilities].
type TargetRef = session.SessionTarget

// Capabilities reports what a driver can do with one specific target. It is
// resolved against live state rather than a static table, because whether a
// target admits a terminal generally depends on its current status.
type Capabilities = session.SessionCapabilities

// Supports reports whether the target admits a session of the given kind. It is
// a function rather than a method because Capabilities aliases a generated type.
func Supports(c Capabilities, k Kind) bool {
	switch k {
	case KindTerminal:
		return c.Terminal
	case KindLogs:
		return c.Logs
	default:
		return false
	}
}

// TerminalOptions parameterizes an interactive terminal session.
type TerminalOptions struct {
	// Command is the argv to execute. When empty the driver picks a sensible
	// interactive shell for the target.
	Command []string
}

// LogOptions parameterizes a log session. Pointer fields are unset when nil,
// which lets a driver distinguish "no tail limit" from "tail zero lines".
type LogOptions struct {
	// Follow keeps the stream open and appends new output as it is produced.
	Follow bool
	// Previous requests the logs of the target's prior instance, if any.
	Previous bool
	// Timestamps prefixes each line with an RFC3339 timestamp.
	Timestamps bool
	// TailLines caps how much history is replayed before following.
	TailLines *int64
	// SinceSeconds bounds history to output produced within the last N seconds.
	SinceSeconds *int64
}

// TerminalSize is a terminal window geometry, in character cells.
type TerminalSize struct {
	Width  uint16 `json:"cols"`
	Height uint16 `json:"rows"`
}

// TerminalIO carries the streams a driver wires an interactive session onto.
//
// Stderr is nil for TTY sessions, where the remote end multiplexes both output
// streams onto Stdout; a driver must tolerate a nil Stderr.
type TerminalIO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	// Resize yields a new geometry each time the client's window changes. It is
	// closed when the session ends.
	Resize <-chan TerminalSize
}

// ExitError reports that the remote command terminated with a non-zero status.
// It is an expected end to a session rather than a transport failure, so the
// handler surfaces it as an exit control frame, not an error frame.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "session exited with a non-zero status"
}

func (e *ExitError) Unwrap() error { return e.Err }

// Driver opens sessions against the resources of a single connection.
// Implementations are built per connection by a [Factory] and serve a single
// session, so they need not be safe for concurrent use.
type Driver interface {
	// Capabilities resolves what the driver can do with target, consulting live
	// state as needed. It returns an ErrTargetNotFound-coded error when target
	// does not exist.
	Capabilities(ctx context.Context, target TargetRef) (Capabilities, error)
}

// TerminalDriver is implemented by drivers that can open interactive terminals.
// A driver that cannot simply does not implement it, and the handler rejects
// the request before upgrading the socket.
type TerminalDriver interface {
	Driver

	// OpenTerminal runs an interactive session, copying between stream and the
	// remote target until ctx is cancelled or the remote command exits. It
	// returns *[ExitError] for a non-zero remote exit status.
	OpenTerminal(ctx context.Context, target TargetRef, opts TerminalOptions, stream TerminalIO) error
}

// LogStreamer is implemented by drivers that can tail a target's log output.
type LogStreamer interface {
	Driver

	// StreamLogs opens the target's log stream. The caller closes the returned
	// reader. With LogOptions.Follow the reader blocks until ctx is cancelled;
	// otherwise it reports io.EOF once history is exhausted.
	StreamLogs(ctx context.Context, target TargetRef, opts LogOptions) (io.ReadCloser, error)
}
