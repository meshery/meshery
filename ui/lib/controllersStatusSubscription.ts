/**
 * SSE-based Meshery controllers-status subscription.
 *
 * Replaces the former `subscribeMesheryControllersStatus` GraphQL subscription.
 * It opens a native EventSource against the backend SSE endpoint
 * (GET /api/system/controllers/status/subscribe?connectionIds=<id>&connectionIds=<id>)
 * which streams the controller (operator / MeshSync / broker) status for the
 * requested connections. Each message is the FULL status array framed as
 * `data: <json-array>\n\n`, so the consumer just replaces its state — no
 * client-side delta merge is needed (that was a source of bugs in the old path).
 *
 * Being same-origin, the EventSource automatically carries the `meshery-provider`
 * auth cookie and natively reconnects on transient drops. `onError` is only
 * raised once the browser permanently closes the stream, and that notification
 * is throttled (see RECONNECT_NOTIFY_DELAY_MS) so a persistent failure can't turn
 * caller-driven re-subscription into a request storm. The returned `{ dispose }`
 * mirrors the shape of the old Relay subscription so callers stay unchanged.
 */

const CONTROLLERS_STATUS_STREAM_URL = '/api/system/controllers/status/subscribe';

// Mirrors eventsSubscription.ts: delay surfacing a permanent close so a caller
// that re-subscribes on error can't bypass EventSource's own backoff and hammer
// the server on a persistent failure (e.g. an expired session).
const RECONNECT_NOTIFY_DELAY_MS = 3000;

export interface ControllerStatusItem {
  connectionId: string;
  controller: string; // OPERATOR | MESHSYNC | BROKER
  status: string;
  version: string;
}

export interface ControllersStatusSubscription {
  dispose: () => void;
}

export type ControllersStatusCallback = (items: ControllerStatusItem[]) => void;
export type ErrorCallback = (error: Event) => void;

function buildStreamUrl(connectionIds: string[]): string {
  const params = new URLSearchParams();
  connectionIds.forEach((id) => id && params.append('connectionIds', id));
  const query = params.toString();
  return query ? `${CONTROLLERS_STATUS_STREAM_URL}?${query}` : CONTROLLERS_STATUS_STREAM_URL;
}

export function subscribeToControllersStatus(
  connectionIds: string[],
  onNext: ControllersStatusCallback,
  onError?: ErrorCallback,
): ControllersStatusSubscription {
  let disposed = false;
  let reconnectTimer: ReturnType<typeof setTimeout> | undefined;
  const source = new EventSource(buildStreamUrl(connectionIds));

  source.onmessage = (message) => {
    try {
      const items = JSON.parse(message.data);
      onNext(Array.isArray(items) ? items : []);
    } catch (error) {
      console.error('[ControllersStatusSubscription] Failed to parse frame', error, message.data);
    }
  };

  source.onerror = (error) => {
    // EventSource silently auto-reconnects while readyState is CONNECTING (0);
    // stay quiet and let it handle transient drops itself.
    if (source.readyState !== EventSource.CLOSED) {
      return;
    }
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
