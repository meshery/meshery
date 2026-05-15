// Package sse provides Server-Sent Events infrastructure for Meshery.
//
// The Event/write helpers here are the source of truth for the SSE wire
// format inside this server.  Once meshery/schemas publishes an AsyncAPI
// document for these channels (see docs/asyncapi/controller-status-sse.yaml)
// and generates Go stubs, this file should be updated to import those
// generated types and remove the hand-rolled Event struct.
//
// TODO(schemas-canonical): replace Event struct with import from
//   github.com/meshery/schemas/models/asyncapi/sse once generated.
package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// heartbeatInterval is how often a keep-alive comment is sent on an idle
// SSE connection so that proxies and load-balancers do not close it.
const heartbeatInterval = 30 * time.Second

// Event is the SSE envelope.  Every message pushed to the client is wrapped
// in this type.  The Name field maps to the SSE "event:" line; Data is
// JSON-serialised and emitted as the "data:" line.
type Event struct {
	Name string `json:"name"`
	Data any    `json:"data"`
}

// SetSSEHeaders writes the standard headers required by the SSE protocol.
// It must be called before any data is written to w.
func SetSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
}

// Write serialises e and pushes it to the client.  It flushes immediately so
// that the browser/client receives the event without buffering.
// Returns an error if the response writer does not support flushing, or if
// JSON serialisation fails.
func Write(w http.ResponseWriter, e Event) error {
	// Guard against SSE protocol injection: the event name is written directly
	// onto the wire as the "event:" field.  A name containing CR or LF would
	// terminate that field early and allow an attacker to inject arbitrary SSE
	// frames into the stream.
	if strings.ContainsAny(e.Name, "\r\n") {
		return fmt.Errorf("sse: event name %q contains illegal CR/LF characters", e.Name)
	}

	payload, err := json.Marshal(e.Data)
	if err != nil {
		return fmt.Errorf("sse: marshal event %q data: %w", e.Name, err)
	}

	if _, err := fmt.Fprintf(w, "event: %s\ndata: %s\n\n", e.Name, payload); err != nil {
		return fmt.Errorf("sse: write event %q: %w", e.Name, err)
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// WriteHeartbeat sends an SSE comment (": heartbeat") that is invisible to
// application-level listeners but keeps the TCP connection alive.
func WriteHeartbeat(w http.ResponseWriter) error {
	if _, err := fmt.Fprintf(w, ": heartbeat\n\n"); err != nil {
		return fmt.Errorf("sse: write heartbeat: %w", err)
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// HeartbeatInterval returns the recommended interval between heartbeats.
// Callers should send a heartbeat at least this often on idle connections.
func HeartbeatInterval() time.Duration {
	return heartbeatInterval
}
