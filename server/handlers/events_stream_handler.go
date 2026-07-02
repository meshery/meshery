package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
)

// eventStreamKeepAliveInterval bounds how long the SSE connection can sit idle
// before we write a comment line. It lets the browser (and any proxy in front
// of Meshery) notice a dead peer, and keeps intermediaries from closing an
// apparently-silent stream.
const eventStreamKeepAliveInterval = 15 * time.Second

// SubscribeEventsHandler streams a user's events over Server-Sent Events. It
// replaces the former `subscribeEvents` GraphQL subscription: both draw from the
// same per-user EventBroadcaster, so any handler that publishes an event (e.g.
// ClientEventHandler) reaches this stream unchanged.
//
// Each message is the raw meshkit events.Event marshalled to JSON, framed as an
// unnamed SSE event (`data: <json>\n\n`) so the browser's EventSource.onmessage
// receives it. Only newly-published events are streamed; the initial list is
// still fetched over REST from GET /api/system/events.
//
// The subscription lives for the duration of the request: when the client
// disconnects the request context is cancelled, we unsubscribe from the
// broadcaster, and the goroutine returns.
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
	ch, unsubscribe := h.config.EventBroadcaster.Subscribe(user.ID)
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

	ctx := req.Context()
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
			var event *events.Event
			switch e := data.(type) {
			case *events.Event:
				event = e
			case events.Event:
				event = &e
			default:
				continue
			}
			payload, err := json.Marshal(event)
			if err != nil {
				h.log.Error(models.ErrMarshal(err, "event"))
				continue
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", payload); err != nil {
				// Client disconnected mid-write; context cancellation will
				// unsubscribe us on the way out.
				return
			}
			flusher.Flush()
		case <-keepAlive.C:
			// SSE comment line: ignored by EventSource, but surfaces a dead
			// connection and keeps proxies from timing the stream out.
			if _, err := fmt.Fprint(w, ": keepalive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
