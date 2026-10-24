---
title: Messaging Framework
description: Architecture of Meshery events and notifications, with the design goals and current status of CloudEvents alignment.
---

Meshery's messaging framework describes how operations and changes in managed infrastructure are represented as events, recorded, and surfaced to users. It also provides a design direction for correlating related events and exchanging them with external systems.

This page distinguishes **existing Meshery behavior** from **proposals in the Messaging Framework and Notification Center design specification**. The specification was developed over several years; not every proposed interface or event attribute was implemented.

## Events and notifications

Meshery produces events for activities affecting managed infrastructure, workloads, and Meshery itself. Events can report the result of a user-initiated operation or a system-observed condition. They provide the information used to build an operational history and to notify users about relevant changes.

The [Notification Center]({{< ref "guides/infrastructure-management/notification-management/index.md" >}}) presents events in the Meshery UI. It supports severity-based filtering and read/unread status, among other ways of finding and acting on notifications. This is the user-facing view of events, not the event transport itself.

### Current event flow

The current eventing implementation includes the following responsibilities:

1. **Event producers** create Meshery events as operations execute or system conditions are observed.
2. **Persistence** records events through Meshery's provider-backed event persistence path.
3. **Broadcasting** publishes events to connected clients so that the UI can receive updates without polling for every event.
4. **Presentation** displays relevant events in the Notification Center and allows users to manage their read/unread state.

Live events are delivered to the UI using Server-Sent Events (SSE) at `GET /api/system/events/subscribe`. The eventing implementation and its ongoing work on correlation and notification quieting are discussed in [meshery/meshery#20879](https://github.com/meshery/meshery/issues/20879).

These responsibilities should not be confused with the [Meshery Broker]({{< ref "concepts/architecture/broker/index.md" >}}). The Broker uses NATS to stream data among Meshery components. A messaging transport such as NATS does not, by itself, establish that the messages use the CloudEvents format.

## Event information and classification

The design specification aims to give events consistent identity, context, and classification so that they can be understood across Meshery components.

| Concern | Purpose |
| --- | --- |
| Event identity | Distinguish an individual occurrence from other occurrences. |
| Source and context | Identify the component or operation responsible for the event. |
| Category and action | Describe the kind of activity or condition reported. |
| Severity | Express the importance of the event and help determine how it should be presented. |
| Operation correlation | Associate events produced by the same higher-level operation. |
| Event data | Carry the details needed to understand the result or failure. |

Meshery's event model has a severity vocabulary informed by the syslog severity ladder, with an additional `success` value. The Notification Center currently exposes a smaller set of user-facing severity filters, including informational, warning, error, and success. The design's broader classification and notification policy should not be assumed to be implemented end to end.

The specification also distinguishes events that belong in an operational or audit history from events that need to interrupt or notify a user. Its design direction is to use event classification and severity to inform that distinction, rather than requiring every producer to invent its own notification rules.

### Correlation and notification lifecycle

A single user operation may produce several events. The design explores correlating those events through a shared operation identity, including CloudEvents-style `correlationid` and trace-context concepts.

**Operation correlation** is different from **condition correlation**. Multiple events can belong to one operation, while a recurring failure can produce many occurrences of the same underlying condition. The latter needs additional rules for suppressing repetitive notifications, recording recurrence, and recognizing resolution.

Centralized suppression, recurrence roll-up, and resolution semantics remain design and implementation concerns tracked in [meshery/meshery#20879](https://github.com/meshery/meshery/issues/20879). They should not be described as guaranteed framework behavior.

## CloudEvents alignment

[CloudEvents](https://cloudevents.io/) defines a common event envelope for interoperability between services and systems. The Messaging Framework specification proposes aligning Meshery events with CloudEvents so that event producers and consumers can exchange event context in a predictable format.

The design considered:

- A CloudEvents-compatible event envelope with identity, source, type, timestamp, and event data.
- Extension attributes for Meshery-specific context, such as severity and correlation information.
- Centralized event construction so that components use consistent event metadata.
- CloudEvents-based interfaces between components and external event consumers.
- Client-side parsing and an ingress interface for externally produced events.

These are **design proposals**, not a list of currently supported Meshery APIs. In particular, the specification's proposed CloudEvents conversion, adapter RPC migration, and external event-ingestion endpoint must not be treated as shipped functionality. The implementation assessment in [meshery/meshery#20879](https://github.com/meshery/meshery/issues/20879) documents this gap.

### Illustrative CloudEvents envelope

The following example illustrates the intended separation of standard CloudEvents context from Meshery-specific event data. **It is conceptual and is not an example of a response from Meshery's current events API.**

```json
{
  "specversion": "1.0",
  "id": "7a43d8a6-4317-4e4a-a33a-2876cce4a4fd",
  "source": "/meshery/server",
  "type": "io.meshery.operation.failed",
  "time": "2026-01-01T12:00:00Z",
  "subject": "connections/example",
  "correlationid": "operation-example",
  "data": {
    "message": "An operation could not be completed"
  }
}
```

CloudEvents `specversion` is `1.0` in this illustration. The `correlationid` field is an illustrative extension attribute, not a claim that the current Meshery event model emits it. The example's event type and payload are illustrative as well.

The existing [CloudEvents model page]({{< ref "extensions/models/cloudevents/index.md" >}}) describes a Meshery integration catalog entry. It should not be interpreted as confirmation that the eventing framework has completed the CloudEvents migration.

## Design considerations and outstanding work

The original specification establishes useful goals, but several implementation decisions require further work:

- **Stable identity:** Event identity should be available consistently to persistence and live subscribers.
- **Correlation:** Producers need a way to associate events with an operation, while recurring conditions need their own correlation semantics.
- **Notification policy:** Decisions about persistence, broadcasting, severity, and suppression should be consistent across event producers.
- **Lifecycle:** Acknowledging a notification as read is distinct from resolving the underlying condition.
- **Interoperability:** CloudEvents compatibility needs an explicit schema and transport contract, not merely a messaging broker or an integration catalog entry.

The original document also discusses possible changes to adapter communication, UI parsing, and external event ingestion. Those sections are historical design context until their implementation and public interfaces are verified.

## References

- [Messaging Framework and Notification Center Design](https://docs.google.com/document/d/1_zpEWBcC6ngOLen_E_nc-octiUdBbpmdwwIX6qzo70M/edit?usp=sharing) — original design specification and discussion.
- [Issue #6656: Add Meshery as an adopter of CloudEvents](https://github.com/meshery/meshery/issues/6656) — documentation and external integration-card work.
- [Issue #20879: Centralize event correlation and notification quieting](https://github.com/meshery/meshery/issues/20879) — implementation assessment and outstanding design decisions.
- [CloudEvents specification](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md) — standard event context and extension-attribute model.
