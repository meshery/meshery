package console

import "github.com/meshery/schemas/models/v1beta1/console"

// The console wire protocol.
//
// The frame types are defined by the `v1beta1/console` construct in
// github.com/meshery/schemas and aliased here; see ConsoleControlMessage there
// for the normative description. In summary:
//
//	binary frame  raw payload bytes
//	              client -> server: stdin keystrokes (terminal consoles only)
//	              server -> client: stdout/log bytes
//
//	text frame    a JSON-encoded ControlMessage
//
// Splitting on opcode keeps the hot path allocation-light and, more to the
// point, avoids base64-inflating every keystroke and every log line by a third
// just to smuggle bytes through a JSON string. Terminal output is not valid
// UTF-8 in general — it carries escape sequences and, for a program writing raw
// bytes, arbitrary octets — so a text frame could not carry it faithfully.
//
// The two directions use disjoint sets of control types:
//
//	client -> server   resize
//	server -> client   ready, error, exit, eof
//
// A server that receives an unknown or wrong-direction control type ignores it,
// so that a newer client can send control frames an older server does not know
// about without breaking the console.
//
// Every console ends with exactly one terminal control frame — error, exit, or
// eof — before the socket is closed, so a client never has to infer why a
// stream stopped from the close code alone.

// ControlType identifies a control frame.
type ControlType = console.ConsoleControlType

const (
	// ControlReady is sent once, by the server, after the console has been
	// established with the remote target. It carries the resolved capabilities
	// so the client can render, for example, the container it actually attached
	// to when the request left the container unspecified.
	ControlReady = console.Ready

	// ControlResize is sent by the client whenever its terminal geometry
	// changes. It is meaningful only for terminal console.
	ControlResize = console.Resize

	// ControlError terminates the console after a failure. The console is over
	// once this is sent.
	ControlError = console.Error

	// ControlExit terminates a terminal console after the remote command
	// finished, carrying its exit status.
	ControlExit = console.Exit

	// ControlEOF terminates a non-following log console after the available
	// history has been sent.
	ControlEOF = console.Eof
)

// ControlMessage is the JSON payload of a text frame. Fields not relevant to a
// given Type are omitted.
type ControlMessage = console.ConsoleControlMessage

// Close codes sent when the server tears down a console socket. They sit in the
// private-use range reserved for applications by RFC 6455 §7.4.2, so they never
// collide with the protocol-defined codes gorilla/websocket may send itself.
const (
	// CloseConsoleEnded means the console finished normally: the remote command
	// exited, or a non-following log stream reached the end of its history.
	CloseConsoleEnded = 4000
	// CloseConsoleFailed means the console ended because of an error, which was
	// described in the preceding ControlError frame.
	CloseConsoleFailed = 4001
)
