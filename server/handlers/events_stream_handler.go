package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
)

// eventStreamKeepAliveInterval bounds how long the SSE connection can sit idle
// before we write a comment line. It lets the browser (and any proxy in front
// of Meshery) notice a dead peer, and keeps intermediaries from closing an
// apparently-silent stream.
const eventStreamKeepAliveInterval = 15 * time.Second

// eventStreamWriteTimeout caps how long a single frame write may block on the
// socket. Without it a client that stops reading (a stalled tab, a half-open
// connection) would wedge this goroutine on the write indefinitely, and - since
// events flow through a per-user buffered channel - back up delivery for that
// user. On timeout the write errors out and the handler unwinds, unsubscribing.
const eventStreamWriteTimeout = 30 * time.Second

// SubscribeEventsHandler streams a user's events over Server-Sent Events. It
// replaces the former `subscribeEvents` GraphQL subscription: both draw from the
// same per-user EventBroadcaster, so any handler that publishes an event (e.g.
// ClientEventHandler) reaches this stream unchanged.
//
// Each message is the raw meshkit events.Event marshalled to JSON, framed as an
// unnamed SSE event (`id: <event.ID>:<seq>\ndata: <json>\n\n` or
// `id: <event.ID>\ndata: <json>\n\n` when event.ID is non-nil, or
// `data: <json>\n\n` otherwise) so the browser's EventSource.onmessage
// receives it and records the last event ID. On reconnect, a `Last-Event-ID`
// header replays buffered events published after that ID (or the full retained
// buffer as best effort if an evicted valid cursor or UUID was provided) before
// streaming live events. Only newly-published (or buffered replay) events are
// streamed; the initial list is still fetched over REST from GET /api/system/events.
//
// The subscription lives for the duration of the request: when the client
// disconnects or the broadcaster closes the channel on buffer overflow, we
// unsubscribe from the broadcaster and the handler returns, prompting
// EventSource to reconnect and recover missed events via Last-Event-ID replay.
func (h *Handler) SubscribeEventsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if user == nil {
		writeJSONError(w, "user unauthorized", http.StatusUnauthorized)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSONError(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Subscribe before writing headers so events published between the request
	// arriving and the stream being ready aren't dropped.
	lastEventID := strings.TrimSpace(req.Header.Get("Last-Event-ID"))
	replay, ch, unsubscribe := h.config.EventBroadcaster.SubscribeFrom(user.ID, lastEventID)
	defer unsubscribe()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// X-Accel-Buffering disables buffering at any nginx hop in front of
	// Meshery so events reach the browser immediately.
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	h.log.Infof("Events SSE subscription started for %s", user.ID)
	defer h.log.Infof("Events SSE subscription stopped for %s", user.ID)

	keepAlive := time.NewTicker(eventStreamKeepAliveInterval)
	defer keepAlive.Stop()

	// responseController lets us bound each frame write with a deadline so a
	// client that stops reading can't pin this goroutine on the write.
	// SetWriteDeadline reports ErrNotSupported for writers that don't implement
	// it (e.g. the test recorder); there we fall back to unbounded writes, which
	// matches the prior behavior.
	responseController := http.NewResponseController(w)

	ctx := req.Context()
	for _, item := range replay {
		if ctx.Err() != nil {
			return
		}
		if !h.writeEventFrame(w, flusher, responseController, item) {
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			return
		case data, open := <-ch:
			// A closed channel yields the zero value without blocking; bail out
			// instead of spinning at 100% CPU on a dead subscription.
			if !open {
				return
			}
			if !h.writeEventFrame(w, flusher, responseController, data) {
				return
			}
		case <-keepAlive.C:
			// SSE comment line: ignored by EventSource, but surfaces a dead
			// connection and keeps proxies from timing the stream out.
			_ = responseController.SetWriteDeadline(time.Now().Add(eventStreamWriteTimeout))
			if _, err := fmt.Fprint(w, ": keepalive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}

// writeEventFrame marshals an event and writes it as an SSE frame. It returns
// false only when writing to w fails, signaling the caller to abort the stream.
func (h *Handler) writeEventFrame(w http.ResponseWriter, flusher http.Flusher, rc *http.ResponseController, data interface{}) bool {
	cursor := ""
	if re, ok := data.(models.ReplayEvent); ok {
		cursor = re.Cursor
		data = re.Data
	}

	var event *events.Event
	switch e := data.(type) {
	case *events.Event:
		if e == nil {
			return true
		}
		event = e
	case events.Event:
		event = &e
	default:
		return true
	}

	payload, err := json.Marshal(event)
	if err != nil {
		h.log.Error(models.ErrMarshal(err, "event"))
		return true
	}

	sseID := cursor
	if sseID == "" && event.ID != uuid.Nil {
		sseID = event.ID.String()
	}

	_ = rc.SetWriteDeadline(time.Now().Add(eventStreamWriteTimeout))
	if sseID != "" {
		_, err = fmt.Fprintf(w, "id: %s\ndata: %s\n\n", sseID, payload)
	} else {
		_, err = fmt.Fprintf(w, "data: %s\n\n", payload)
	}
	if err != nil {
		// Client disconnected or stalled past the write deadline; context
		// cancellation will unsubscribe us on the way out.
		return false
	}
	flusher.Flush()
	return true
}
