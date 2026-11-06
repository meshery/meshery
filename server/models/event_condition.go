package models

import (
	"context"
	"strings"

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

// EventSourceServer is the `source` of every event Meshery Server emits. It is
// the CloudEvents-compatible name of the emitting surface; `systemId` names the
// instance.
const EventSourceServer = "meshery-server"

// EventSourceUnknownClient is the `source` of an event a client sent to
// POST /api/events without saying which surface it is. The server records the
// event on the client's behalf, so it is never attributed to EventSourceServer.
const EventSourceUnknownClient = "unknown-client"

// ApplyClientEventFields copies the optional correlation, provenance and trace
// fields of an event a client sent (the schemas EventPayload) onto the event
// the server records for it, so a condition a client reports correlates as it
// would on the remote provider. The client's operationId groups the event with
// the rest of its operation when it sent one. Check the result with
// events.Event.Validate.
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
	if !sent.OperationID.IsNil() {
		event.OperationID = sent.OperationID
	}
	if sent.Source != nil && *sent.Source != "" {
		event.Source = sent.Source
	} else {
		source := EventSourceUnknownClient
		event.Source = &source
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
// It sets source EventSourceServer when the producer did not say (every event
// that reaches a persister was emitted by this server) and bounds the metadata
// with GuardEventMetadataSize.
//
// The event is a copy and its metadata map is never mutated: the same map is
// often shared with the copy published to the EventBroadcaster.
func PrepareEventForPersistence(event events.Event, log logger.Handler) events.Event {
	if event.Source == nil || *event.Source == "" {
		source := EventSourceServer
		event.Source = &source
	}
	return GuardEventMetadataSize(event, log)
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
