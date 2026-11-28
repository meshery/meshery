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

// Backoff between reopens, on top of RECONNECT_NOTIFY_DELAY_MS: 1s, 2s, 4s, ...
// capped at 30s, so a long outage doesn't have every open tab reconnecting
// every few seconds. It resets once the stream delivers a frame again.
const REOPEN_BASE_DELAY_MS = 1000;
const REOPEN_MAX_DELAY_MS = 30_000;

/**
 * Like subscribeToControllersStatus, but reopens the stream after the browser
 * permanently closes it (e.g. a proxy answering 502 while Meshery Server
 * restarts). Without this, controller status freezes until a page reload.
 * Reopens back off (see REOPEN_BASE_DELAY_MS), and the server writes a full
 * snapshot on connect, so the reopened stream restores current state at once.
 */
export function watchControllersStatus(
  connectionIds: string[],
  onNext: ControllersStatusCallback,
): ControllersStatusSubscription {
  let current: ControllersStatusSubscription;
  let attempt = 0;
  let reopenTimer: ReturnType<typeof setTimeout> | undefined;

  const open = () => {
    current = subscribeToControllersStatus(
      connectionIds,
      (items) => {
        attempt = 0;
        onNext(items);
      },
      () => {
        current.dispose();
        const delay = Math.min(REOPEN_MAX_DELAY_MS, REOPEN_BASE_DELAY_MS * 2 ** attempt);
        attempt += 1;
        reopenTimer = setTimeout(open, delay);
      },
    );
  };
  open();

  return {
    // onError is never raised after dispose, so clearing a pending reopen and
    // disposing the current stream stops any further reopening.
    dispose: () => {
      clearTimeout(reopenTimer);
      current.dispose();
    },
  };
}
