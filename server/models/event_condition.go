package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
)

// Producer-side API for the meshery/schemas event correlation contract
// (docs/event-correlation-contract.md in meshery/schemas): where an event came
// from (`source`), which operation it belongs to (`operationId`, W3C
// `traceId`/`parentId`) and which ongoing condition it reports
// (`correlationKey`). Every field is optional; an event that sets none of them
// is persisted and sent exactly as before.

// EventDispositionHeader is the createEvent response header that says what the
// server did with the event: `created`, `replayed`, `incremented` or
// `resolved`. A server that predates condition correlation omits it.
const EventDispositionHeader = "Event-Disposition"

// EventDisposition is what persisting an event did, as createEvent reports it
// in EventDispositionHeader. Empty means a remote provider did not say.
type EventDisposition string

const (
	// EventCreated: a new event was recorded.
	EventCreated EventDisposition = "created"
	// EventReplayed: the event was already recorded (its idempotencyKey was
	// already used, or the same event, by ID, was persisted before); nothing
	// changed.
	EventReplayed EventDisposition = "replayed"
	// EventIncremented: an occurrence was counted on the open condition.
	EventIncremented EventDisposition = "incremented"
	// EventResolved: the open condition was closed.
	EventResolved EventDisposition = "resolved"
)

// PersistedEvent is the event as recorded and what recording it did. For an
// incremented or resolved condition, Event is the condition's stored event,
// whose ID is not the ID of the event that was persisted, so it is the one to
// broadcast or return to a client.
//
// Event is nil only when a remote provider recorded nothing new (incremented,
// resolved or replayed) and did not answer with the stored event: there is
// then no stored event to broadcast or return, and nothing new to report.
type PersistedEvent struct {
	Event       *events.Event
	Disposition EventDisposition
}

// EventSourceServer is the `source` of every event Meshery Server emits. It is
// the CloudEvents-compatible name of the emitting surface; `systemId` names the
// instance.
const EventSourceServer = "meshery-server"

// ApplyClientEventFields copies the optional correlation, provenance and trace
// fields of an event a client sent (the schemas EventPayload) onto the event
// the server records for it, so a condition a client reports correlates as it
// would on the remote provider. The client's operationId groups the event with
// the rest of its operation when it sent one. A client that names no source
// leaves it unset. Check the result with events.Event.Validate.
func ApplyClientEventFields(event *events.Event, sent events.Event) {
	event.CorrelationKey = sent.CorrelationKey
	event.OccurrenceCount = sent.OccurrenceCount
	event.FirstSeenAt = sent.FirstSeenAt
	event.LastSeenAt = sent.LastSeenAt
	event.ResolvedAt = sent.ResolvedAt
	event.Type = sent.Type
	event.SpecVersion = sent.SpecVersion
	event.TraceID = sent.TraceID
	event.ParentID = sent.ParentID
	event.Source = sent.Source
	if !sent.OperationID.IsNil() {
		event.OperationID = sent.OperationID
	}
}

// NewServerEvent starts an event Meshery Server emits on behalf of the
// operation ctx belongs to. It carries source EventSourceServer and, when ctx
// carries a valid OpenTelemetry span, that span's trace context: the router's
// otelmux middleware starts one per request, continuing the W3C `traceparent`
// the client sent, so the event joins the trace of the operation the client
// started. Nothing mints a trace context at persist time.
//
// A nil ctx is allowed and sets no trace context. Callers chain the usual
// builder setters, plus WithOperationID to group the events of one operation
// and WithCorrelationKey (with ConditionKey) to report a condition.
func NewServerEvent(ctx context.Context) *events.EventBuilder {
	builder := events.NewEvent().WithSource(EventSourceServer)
	if ctx != nil {
		builder = builder.WithTraceContextFrom(ctx)
	}
	return builder
}

// ConditionKey returns the correlationKey of a condition Meshery Server
// reports: `meshery-server/<component>/<subject>/<condition>`, for example
// `meshery-server/k8s-connection/<connectionId>/unreachable`. The key names the
// condition, never one occurrence of it, and is qualified by the emitting
// source and component, so it is never an error code alone.
//
// component is the part of the server that watches the condition
// (`k8s-connection`), subject identifies what the condition is about (a
// connection id) and condition says what is wrong (`unreachable`). Each part
// must be non-empty; a slash inside one is replaced so the key keeps exactly
// four segments.
func ConditionKey(component, subject, condition string) string {
	return strings.Join([]string{
		EventSourceServer,
		conditionKeySegment(component),
		conditionKeySegment(subject),
		conditionKeySegment(condition),
	}, "/")
}

func conditionKeySegment(segment string) string {
	segment = strings.TrimSpace(segment)
	if segment == "" {
		return "unknown"
	}
	return strings.ReplaceAll(segment, "/", "-")
}

// PrepareEventForPersistence is applied to every event Meshery Server
// persists, locally or with a remote provider, before it is written or sent.
// It gives an event without an ID (one constructed without events.NewEvent,
// which mints one) a fresh one,
// so no event is stored, sent to a remote provider or broadcast with the zero
// ID, and it bounds the metadata with GuardEventMetadataSize. It never sets
// `source`; withServerSource does that for the events the server itself emits.
//
// The event is a copy and its metadata map is never mutated: the same map is
// often shared with the copy published to the EventBroadcaster.
func PrepareEventForPersistence(event events.Event, log logger.Handler) events.Event {
	if event.ID.IsNil() {
		id, err := uuid.NewV4()
		if err != nil {
			// Persisting fails on its own: the database insert mints an ID the
			// same way and returns the same error.
			if log != nil {
				log.Warn(fmt.Errorf("could not assign an ID to event for category %q action %q: %w", event.Category, event.Action, err))
			}
		} else {
			event.ID = id
		}
	}
	return GuardEventMetadataSize(event, log)
}

// withServerSource sets source EventSourceServer on an event Meshery Server
// emitted without naming its source. Provider.PersistEvent applies it, as every
// event persisted through it was emitted by this server; PersistClientEvent,
// which records an event a client sent, does not.
func withServerSource(event events.Event) events.Event {
	if event.Source == nil || *event.Source == "" {
		source := EventSourceServer
		event.Source = &source
	}
	return event
}

// correlatesEvent reports whether a create of event takes part in condition
// correlation: it names a condition and identifies the user it belongs to.
// Like the contract's createEvent, an event that identifies no user (a system
// event without an owner) is recorded without being correlated.
func correlatesEvent(event events.Event) bool {
	if event.CorrelationKey == nil || *event.CorrelationKey == "" {
		return false
	}
	return event.Owner != nil && !event.Owner.IsNil()
}
