/**
 * SSE-based events subscription.
 *
 * Replaces the former `subscribeEvents` GraphQL subscription. It opens a native
 * EventSource against the backend SSE endpoint (GET /api/system/events/subscribe)
 * which streams the user's events, each framed as `data: <event-json>\n\n`.
 *
 * Being same-origin, the EventSource automatically carries the `meshery-provider`
 * auth cookie and will automatically attempt to reconnect if the connection
 * drops, so no separate retry logic is needed for this stream. The returned
 * `{ dispose }` mirrors the shape of the old Relay subscription so callers stay
 * unchanged.
 */

const EVENTS_STREAM_URL = '/api/system/events/subscribe';

export interface EventSubscription {
  dispose: () => void;
}

export type EventCallback = (event: unknown) => void;
export type ErrorCallback = (error: Event) => void;

export function subscribeToEvents(
  onNext: EventCallback,
  onError?: ErrorCallback,
): EventSubscription {
  const source = new EventSource(EVENTS_STREAM_URL);

  source.onmessage = (message) => {
    try {
      onNext(JSON.parse(message.data));
    } catch (error) {
      console.error('[EventsSubscription] Failed to parse event', error, message.data);
    }
  };

  source.onerror = (error) => {
    // EventSource silently auto-reconnects while readyState is CONNECTING (0).
    // Only surface the error when the browser has permanently closed the stream
    // (CLOSED = 2) so the caller can decide whether to re-subscribe.
    if (source.readyState === EventSource.CLOSED) {
      onError?.(error);
    }
  };

  return {
    dispose: () => source.close(),
  };
}
