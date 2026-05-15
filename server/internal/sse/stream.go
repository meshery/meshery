// Package sse provides a reusable Server-Sent Events writer for streaming
// JSON-serialized values to HTTP clients. It is the canonical SSE plumbing
// for the SSE migration: handlers should depend on Stream rather than rolling
// their own copy of the header set + flusher loop.
package sse

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/meshery/meshkit/logger"
)

// HeartbeatInterval is how often Stream emits an SSE comment line so that
// idle reverse proxies (nginx, traefik, ALB) don't drop the connection.
// 25 s is below the common 30 s default idle timeout.
//
// It is a var rather than a const so tests can shrink it; production code
// should not mutate it.
var HeartbeatInterval = 25 * time.Second

// ErrFlusherUnavailable is returned by Stream when the supplied
// http.ResponseWriter doesn't implement http.Flusher. Callers should detect
// this BEFORE writing SSE headers if they want to fall back to a regular
// JSON error envelope; once Stream has set the text/event-stream headers,
// any further error must be surfaced via the stream itself.
var ErrFlusherUnavailable = errors.New("sse: response writer does not support flushing")

// Event is the optional envelope for publishing on a Stream channel when the
// caller wants to set the SSE event name (matching the eventName parameter on
// sseSubscribe in the UI). Callers that don't need named events can keep
// publishing plain JSON-marshalable values directly on the channel — Stream
// emits them as the default (unnamed) `message` event.
type Event struct {
	// Name is the SSE event name. Empty string emits a default `message`
	// event (matching the prior behavior).
	Name string
	// Data is the JSON-marshalable payload.
	Data any
}

// Stream writes JSON-serialized values from events to w as SSE messages
// until ctx is cancelled or events closes. It emits a `: keepalive\n\n`
// comment every HeartbeatInterval so idle reverse proxies don't drop the
// connection. Sets Content-Type, Cache-Control, Connection, and
// X-Accel-Buffering headers before the first write.
//
// Values received on events are written as default (unnamed) message frames
// unless the caller wraps them in an Event{Name, Data} envelope, in which
// case the SSE `event:` line is emitted before the `data:` line — matching
// the eventName parameter on the UI's sseSubscribe.
//
// Returns nil on clean shutdown (ctx cancelled or events closed), and a
// non-nil error if the writer doesn't support flushing or if a write fails
// (which almost always means the client disconnected).
func Stream(ctx context.Context, w http.ResponseWriter, events <-chan any, log logger.Handler) error {
	flusher, ok := w.(http.Flusher)
	if !ok {
		if log != nil {
			log.Error(ErrFlusherUnavailable)
		}
		return ErrFlusherUnavailable
	}

	// Headers must be set before the first Write/Flush. Cache-Control: no-cache
	// keeps intermediate caches from buffering; Connection: keep-alive is a
	// hint to legacy HTTP/1.1 proxies. X-Accel-Buffering: no disables nginx's
	// response buffering specifically — without it, nginx will hold our
	// chunks in 8 KB increments and break the live stream.
	h := w.Header()
	h.Set("Content-Type", "text/event-stream")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
	h.Set("X-Accel-Buffering", "no")

	// Flush headers so the client knows the stream is open even before the
	// first real event lands. Useful for clients that want to render a
	// "connected" indicator independent of payload arrival.
	flusher.Flush()

	heartbeat := time.NewTicker(HeartbeatInterval)
	defer heartbeat.Stop()

	// Reused across loop iterations so high-frequency streams don't
	// allocate a fresh buffer per event. Reset() keeps the underlying
	// backing array; the buffer grows toward the largest frame we've seen.
	var compact bytes.Buffer

	for {
		select {
		case <-ctx.Done():
			if log != nil {
				log.Debug("sse: context done, closing stream")
			}
			return nil

		case ev, ok := <-events:
			if !ok {
				if log != nil {
					log.Debug("sse: events channel closed, stopping stream")
				}
				return nil
			}

			// Unwrap the optional event-name envelope. A bare payload uses
			// the default (unnamed) `message` event. Both value and pointer
			// forms are handled so callers can publish &Event{...} without
			// silently degrading to an unnamed payload containing the
			// serialized envelope.
			name := ""
			data := ev
			switch env := ev.(type) {
			case Event:
				name = env.Name
				data = env.Data
			case *Event:
				if env != nil {
					name = env.Name
					data = env.Data
				}
			}

			payload, err := json.Marshal(data)
			if err != nil {
				// A marshal failure is a programmer error in the publisher,
				// not a transport problem. Log it and skip the event rather
				// than tearing down the whole stream — the next event has a
				// reasonable chance of being well-formed.
				if log != nil {
					log.Error(fmt.Errorf("sse: marshal event: %w", err))
				}
				continue
			}

			// SSE treats a bare LF as a field terminator, so a payload that
			// contains literal newlines (e.g. from json.RawMessage that's
			// already pretty-printed) would arrive at the client as a
			// truncated or split frame. Compact the JSON to a single line
			// so the wire format stays valid regardless of how the publisher
			// produced it. Reset() preserves the buffer's backing array
			// across iterations to keep alloc pressure low.
			compact.Reset()
			if err := json.Compact(&compact, payload); err != nil {
				// Should not happen for output we just produced via
				// json.Marshal, but if it does, treat it like a marshal
				// failure — log and skip rather than corrupt the stream.
				if log != nil {
					log.Error(fmt.Errorf("sse: compact event: %w", err))
				}
				continue
			}

			// Write directly to w to avoid materializing the whole frame as
			// a string first. Frame layout: optional "event: NAME\n",
			// then "data: <compact-json>\n\n".
			if name != "" {
				if _, err := fmt.Fprintf(w, "event: %s\n", name); err != nil {
					if log != nil {
						log.Error(fmt.Errorf("sse: write event name: %w", err))
					}
					return err
				}
			}
			if _, err := fmt.Fprintf(w, "data: %s\n\n", compact.Bytes()); err != nil {
				// A write error here is almost always a broken pipe (client
				// disconnected). Return so the caller can clean up — there
				// is no point spinning publishing into a dead socket.
				if log != nil {
					log.Error(fmt.Errorf("sse: write event: %w", err))
				}
				return err
			}
			flusher.Flush()

		case <-heartbeat.C:
			// SSE comment lines start with ":" and are ignored by the
			// EventSource spec on the client side, so they make a safe
			// keepalive that doesn't surface as a message event.
			if _, err := fmt.Fprint(w, ": keepalive\n\n"); err != nil {
				if log != nil {
					log.Error(fmt.Errorf("sse: write heartbeat: %w", err))
				}
				return err
			}
			flusher.Flush()
		}
	}
}
