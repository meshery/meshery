/**
 * SSE-based events subscription.
 *
 * Replaces the former `subscribeEvents` GraphQL subscription. It opens a native
 * EventSource against the backend SSE endpoint (GET /api/system/events/subscribe)
 * which streams the user's events, each framed as `data: <event-json>\n\n`.
 *
 * Being same-origin, the EventSource automatically carries the `meshery-provider`
 * auth cookie and natively reconnects on transient drops. `onError` is only
 * raised once the browser permanently closes the stream, and that notification
 * is throttled (see RECONNECT_NOTIFY_DELAY_MS) so a persistent failure can't turn
 * caller-driven re-subscription into a request storm. The returned `{ dispose }`
 * mirrors the shape of the old Relay subscription so callers stay unchanged.
 */

const EVENTS_STREAM_URL = '/api/system/events/subscribe';

// When EventSource permanently closes (CLOSED) we notify the caller so it can
// re-subscribe. That notification is delayed so a persistent failure (e.g. an
// expired session where every connection attempt returns 401) cannot drive a
// tight re-subscribe loop that hammers the server; it rate-limits reconnects to
// one attempt per interval, mirroring EventSource's own native backoff.
const RECONNECT_NOTIFY_DELAY_MS = 3000;

export interface EventSubscription {
  dispose: () => void;
}

export type EventCallback = (event: unknown) => void;
export type ErrorCallback = (error: Event) => void;

export function subscribeToEvents(
  onNext: EventCallback,
  onError?: ErrorCallback,
): EventSubscription {
  let disposed = false;
  let reconnectTimer: ReturnType<typeof setTimeout> | undefined;
  const source = new EventSource(EVENTS_STREAM_URL);

  source.onmessage = (message) => {
    try {
      onNext(JSON.parse(message.data));
    } catch (error) {
      console.error('[EventsSubscription] Failed to parse event', error, message.data);
    }
  };

  source.onerror = (error) => {
    // EventSource silently auto-reconnects while readyState is CONNECTING (0);
    // stay quiet and let it handle transient drops itself.
    if (source.readyState !== EventSource.CLOSED) {
      return;
    }
    // CLOSED means the browser gave up (e.g. the server returned a non-2xx
    // status or a non-event-stream body). Surfacing this immediately would let
    // the caller re-subscribe in a tight loop, bypassing EventSource's own
    // backoff and hammering the server on a persistent failure such as an
    // expired session. Delay the notification to rate-limit re-subscription, and
    // coalesce repeat errors onto a single pending timer.
    if (reconnectTimer !== undefined) {
      return;
    }
    reconnectTimer = setTimeout(() => {
      reconnectTimer = undefined;
      if (!disposed) {
        onError?.(error);
      }
    }, RECONNECT_NOTIFY_DELAY_MS);
  };

  return {
    dispose: () => {
      disposed = true;
      if (reconnectTimer !== undefined) {
        clearTimeout(reconnectTimer);
      }
      source.close();
    },
  };
}
