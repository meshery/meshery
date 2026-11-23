package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshery/server/console"

	// Registers the Kubernetes console driver. Support for another connection
	// kind arrives as another blank import here, with no change to this file.
	_ "github.com/meshery/meshery/server/console/kubernetes"

	mkerrors "github.com/meshery/meshkit/errors"
)

// copyBufferSize is the chunk size for relaying log output to the client. Log
// lines are small; batching them into a few KiB per frame keeps the frame count
// down without adding perceptible latency to a followed stream.
const copyBufferSize = 8 << 10

// GetConsoleCapabilities handles
// GET /api/integrations/connections/{connectionId}/consoles/capabilities.
//
// It reports which console kinds the addressed resource admits, so the UI can
// decide whether to offer a terminal or a log tail before the user clicks, and
// which containers they may choose between.
func (h *Handler) GetConsoleCapabilities(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	driver, target, ok := h.resolveConsoleDriver(w, req, provider)
	if !ok {
		return
	}

	caps, err := driver.Capabilities(req.Context(), target)
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, consoleErrorStatus(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(caps); err != nil {
		h.log.Error(models.ErrMarshal(err, "console capabilities"))
	}
}

// TerminalConsoleHandler handles
// GET /api/integrations/connections/{connectionId}/consoles/terminal.
//
// It upgrades to a WebSocket and runs an interactive terminal against the
// target for the life of the socket.
func (h *Handler) TerminalConsoleHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	h.serveConsole(w, req, user, provider, console.KindTerminal, func(ctx context.Context, driver console.Driver, target console.TargetRef, socket *console.Socket) error {
		terminal, ok := driver.(console.TerminalDriver)
		if !ok {
			return console.ErrUnsupportedConsole(string(console.KindTerminal), target.Resource, "")
		}
		return terminal.OpenTerminal(ctx, target, console.TerminalOptions{
			Command: req.URL.Query()["command"],
		}, console.TerminalIO{
			Stdin:  socket.Stdin(),
			Stdout: socket.Stdout(),
			Resize: socket.Resize(),
		})
	})
}

// LogConsoleHandler handles
// GET /api/integrations/connections/{connectionId}/consoles/logs.
//
// It upgrades to a WebSocket and relays the target's log output. With
// `follow=true` the socket stays open until the client leaves; otherwise it
// closes once the available history has been sent.
func (h *Handler) LogConsoleHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	opts, err := logOptionsFromQuery(req)
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, http.StatusBadRequest)
		return
	}

	h.serveConsole(w, req, user, provider, console.KindLogs, func(ctx context.Context, driver console.Driver, target console.TargetRef, socket *console.Socket) error {
		streamer, ok := driver.(console.LogStreamer)
		if !ok {
			return console.ErrUnsupportedConsole(string(console.KindLogs), target.Resource, "")
		}
		stream, err := streamer.StreamLogs(ctx, target, opts)
		if err != nil {
			return err
		}
		defer stream.Close()

		// io.CopyBuffer, not io.Copy: the log stream implements WriteTo, which
		// io.Copy would prefer, handing our writer whatever chunk size the
		// stream feels like and bypassing the buffer sizing above.
		_, err = io.CopyBuffer(socket.Stdout(), stream, make([]byte, copyBufferSize))
		switch {
		case err == nil:
			// The stream reached its end. Copy reports that as success.
			return nil
		case errors.Is(err, io.ErrClosedPipe), ctx.Err() != nil:
			// The client left. The socket teardown and the context cancellation
			// race each other, so either may be what surfaces here; neither is
			// a failure worth logging or reporting.
			return nil
		default:
			return console.ErrOpenConsole(err, string(console.KindLogs))
		}
	})
}

// runConsole is the body of a console, run after the socket is upgraded.
type runConsole func(ctx context.Context, driver console.Driver, target console.TargetRef, socket *console.Socket) error

// serveConsole is the shared lifecycle for every console kind.
//
// Everything that can fail with a useful HTTP status — a bad target, an
// unsupported kind, an exhausted console budget — is settled before the
// upgrade. Once the socket is live the only channel back to the client is a
// control frame, and a browser surfaces a failed handshake far better than it
// surfaces a socket that opened and immediately closed.
func (h *Handler) serveConsole(w http.ResponseWriter, req *http.Request, user *models.User, provider models.Provider, kind console.Kind, run runConsole) {
	driver, target, ok := h.resolveConsoleDriver(w, req, provider)
	if !ok {
		return
	}

	caps, err := driver.Capabilities(req.Context(), target)
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, consoleErrorStatus(err))
		return
	}
	if !console.Supports(caps, kind) {
		err := console.ErrUnsupportedConsole(string(kind), target.Resource, caps.Reason)
		h.log.Error(err)
		writeMeshkitError(w, err, http.StatusBadRequest)
		return
	}

	release, err := console.Limit.Acquire(user.ID.String())
	if err != nil {
		h.log.Warn(err)
		writeMeshkitError(w, err, http.StatusTooManyRequests)
		return
	}
	defer release()

	socket, err := console.Upgrade(w, req, h.log)
	if err != nil {
		// The upgrader has already written an HTTP error response.
		h.log.Error(err)
		return
	}

	// From here the socket owns the response. Every exit path below closes it
	// with a code, after a terminal control frame explaining why.
	ctx := socket.Start(req.Context())

	if err := socket.SendControl(console.ControlMessage{
		Type:         console.ControlReady,
		Capabilities: &caps,
	}); err != nil {
		h.log.Debug("console client vanished before it was ready: " + err.Error())
		socket.Close(console.CloseConsoleFailed, "client went away")
		return
	}

	err = run(ctx, driver, target, socket)
	h.finishConsole(socket, kind, err)
}

// finishConsole sends the single terminal control frame that ends every console
// and closes the socket, so a client never has to infer from a bare close code
// why its stream stopped.
func (h *Handler) finishConsole(socket *console.Socket, kind console.Kind, err error) {
	var exitErr *console.ExitError

	switch {
	case err == nil:
		final := console.ControlMessage{Type: console.ControlEOF}
		if kind == console.KindTerminal {
			final = console.ControlMessage{Type: console.ControlExit, ExitCode: 0}
		}
		_ = socket.SendControl(final)
		socket.Close(console.CloseConsoleEnded, "console ended")

	case errors.As(err, &exitErr):
		_ = socket.SendControl(console.ControlMessage{Type: console.ControlExit, ExitCode: exitErr.Code})
		socket.Close(console.CloseConsoleEnded, "command exited")

	default:
		h.log.Error(err)
		_ = socket.SendControl(console.ControlMessage{
			Type:    console.ControlError,
			Code:    mkerrors.GetCode(err),
			Message: mkerrors.GetSDescription(err),
		})
		socket.Close(console.CloseConsoleFailed, "console failed")
	}
}

// resolveConsoleDriver turns the request into a driver bound to the addressed
// connection plus the target within it, writing the HTTP error and returning
// false on any failure.
func (h *Handler) resolveConsoleDriver(w http.ResponseWriter, req *http.Request, provider models.Provider) (console.Driver, console.TargetRef, bool) {
	var target console.TargetRef

	token, _ := req.Context().Value(models.TokenCtxKey).(string)

	connection, ok := h.fetchConnection(w, req, token, provider)
	if !ok {
		return nil, target, false
	}

	target, err := targetFromQuery(req)
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, http.StatusBadRequest)
		return nil, target, false
	}

	driver, err := console.Default.Driver(req.Context(), console.ConnectionContext{
		Connection: connection,
		Provider:   provider,
		Token:      token,
	})
	if err != nil {
		h.log.Error(err)
		writeMeshkitError(w, err, consoleErrorStatus(err))
		return nil, target, false
	}
	return driver, target, true
}

// targetFromQuery reads the target reference out of the query string.
func targetFromQuery(req *http.Request) (console.TargetRef, error) {
	query := req.URL.Query()
	target := console.TargetRef{
		Resource:  strings.TrimSpace(query.Get("resource")),
		Namespace: strings.TrimSpace(query.Get("namespace")),
		Name:      strings.TrimSpace(query.Get("name")),
		Container: strings.TrimSpace(query.Get("container")),
	}
	if target.Resource == "" {
		return target, console.ErrInvalidTarget("the `resource` query parameter is required")
	}
	if target.Name == "" {
		return target, console.ErrInvalidTarget("the `name` query parameter is required")
	}
	return target, nil
}

// logOptionsFromQuery reads the log stream parameters out of the query string.
func logOptionsFromQuery(req *http.Request) (console.LogOptions, error) {
	query := req.URL.Query()
	opts := console.LogOptions{
		Follow:     query.Get("follow") == "true",
		Previous:   query.Get("previous") == "true",
		Timestamps: query.Get("timestamps") == "true",
	}

	tailLines, err := optionalPositiveInt64(query.Get("tailLines"), "tailLines")
	if err != nil {
		return opts, err
	}
	opts.TailLines = tailLines

	sinceSeconds, err := optionalPositiveInt64(query.Get("sinceSeconds"), "sinceSeconds")
	if err != nil {
		return opts, err
	}
	// The Kubernetes API rejects sinceSeconds=0, and a zero-second window is
	// meaningless anyway; treat it as unset.
	if sinceSeconds != nil && *sinceSeconds == 0 {
		sinceSeconds = nil
	}
	opts.SinceSeconds = sinceSeconds

	return opts, nil
}

// optionalPositiveInt64 parses an absent-or-non-negative integer parameter.
func optionalPositiveInt64(raw, name string) (*int64, error) {
	if raw == "" {
		return nil, nil
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value < 0 {
		return nil, console.ErrInvalidTarget("the `" + name + "` query parameter must be a non-negative integer")
	}
	return &value, nil
}

// consoleErrorStatus maps a console error's MeshKit code onto an HTTP status.
// Keeping the mapping in one place is why the consoles package returns coded
// errors instead of bare sentinels.
func consoleErrorStatus(err error) int {
	switch mkerrors.GetCode(err) {
	case console.ErrTargetNotFoundCode:
		return http.StatusNotFound
	case console.ErrInvalidTargetCode, console.ErrNoDriverCode, console.ErrUnsupportedConsoleCode:
		return http.StatusBadRequest
	case console.ErrConsoleLimitCode:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// fetchConnection resolves the {connectionId} path parameter into a connection
// of any kind. Whether that kind can serve consoles is the registry's call, not
// this function's, which is what keeps the handler free of Kubernetes.
func (h *Handler) fetchConnection(w http.ResponseWriter, req *http.Request, token string, provider models.Provider) (*connections.Connection, bool) {
	connectionID, err := uuid.FromString(mux.Vars(req)["connectionId"])
	if err != nil || connectionID == uuid.Nil {
		idErr := ErrEmptyConnectionID()
		h.log.Error(idErr)
		writeMeshkitError(w, idErr, http.StatusBadRequest)
		return nil, false
	}
	connection, statusCode, err := provider.GetConnectionByID(token, connectionID)
	if err != nil {
		h.log.Error(err)
		if statusCode < http.StatusContinue {
			statusCode = http.StatusInternalServerError
		}
		writeMeshkitError(w, err, statusCode)
		return nil, false
	}
	return connection, true
}
