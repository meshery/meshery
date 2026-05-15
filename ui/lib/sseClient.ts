/**
 * sseClient.ts
 *
 * Typed wrapper around the browser-native EventSource for the SSE migration.
 * EventSource handles reconnection natively (the cadence is implementation-
 * defined and can be tuned by the server via the SSE `retry:` field), so this
 * module deliberately does NOT layer extra retry on top — adding our own
 * retry would either double-retry or fight the built-in semantics, both of
 * which we have been burned by in the WS client.
 *
 * The wrapper integrates with lib/wsStatus.ts so that the existing
 * connection-status XState machine continues to receive `connected | closed |
 * error` notifications during the migration period.
 */
import {
  _notifySubscriptionClose,
  _notifySubscriptionError,
  _notifySubscriptionOpen,
} from './wsStatus';

export interface SSESubscribeOptions {
  /** Server-relative path, e.g. '/api/system/kubernetes/contexts/stream'. */
  path: string;
  /**
   * Optional query parameters. Arrays expand to repeated keys
   * (`?id=a&id=b`). undefined values are dropped.
   */
  params?: Record<string, string | string[] | undefined>;
  /** Called per SSE message with the JSON-parsed payload. */
  onMessage: (data: unknown) => void;
  /**
   * Optional error handler. If omitted, a default handler logs to the
   * console with the subscription name — callers opt into surfacing errors
   * because most SSE error events are spurious reconnect cycles the browser
   * is already healing.
   */
  onError?: (err: Event) => void;
  /** SSE named event to listen for; defaults to the unnamed 'message'. */
  eventName?: string;
  /** Used in diagnostic log lines; defaults to 'Unknown'. */
  subscriptionName?: string;
}

export interface SSESubscription {
  /** Tear down the EventSource. Safe to call more than once. */
  dispose: () => void;
  /**
   * Close the current EventSource and open a new one to the same path with
   * the supplied params. Preserves onMessage / onError / eventName /
   * subscriptionName. Useful when a filter the user controls (e.g. a context
   * picker) changes and we want a single logical subscription to follow it.
   */
  rebind: (newParams: Record<string, string | string[] | undefined>) => void;
}

/**
 * Build a URL with query parameters. Arrays expand to repeated keys,
 * undefined values are skipped, everything else is coerced to string and
 * URL-encoded by URLSearchParams.
 */
function buildUrl(
  path: string,
  params: Record<string, string | string[] | undefined> | undefined,
): string {
  if (!params) return path;
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined) continue;
    if (Array.isArray(value)) {
      for (const v of value) {
        search.append(key, v);
      }
    } else {
      search.append(key, value);
    }
  }
  const qs = search.toString();
  if (!qs) return path;
  return `${path}${path.includes('?') ? '&' : '?'}${qs}`;
}

/**
 * Open an SSE subscription.
 *
 * Lifecycle:
 *   1. Create EventSource with withCredentials so the meshery session cookie
 *      rides along.
 *   2. On open, mark this subscription as "active" with wsStatus so the
 *      connection indicator lights up.
 *   3. On message of the requested eventName, JSON-parse the payload and
 *      hand it to onMessage. Parse failures are logged but not thrown —
 *      a malformed event should not kill the stream.
 *   4. On error, hand the event to onError (or the default logger) and let
 *      the browser's built-in reconnect logic recover.
 *   5. On dispose, close the EventSource and decrement the wsStatus count.
 */
export function sseSubscribe(opts: SSESubscribeOptions): SSESubscription {
  const { path, onMessage, onError, eventName = 'message', subscriptionName = 'Unknown' } = opts;

  // Mutable closure-captured state so rebind() can swap the EventSource
  // without forcing callers to re-register handlers.
  let currentParams = opts.params;
  let source: EventSource | null = null;
  let disposed = false;
  // Track whether THIS subscription is currently counted as open in wsStatus.
  // We only flip the count on the actual open <-> closed transitions to keep
  // the aggregate connection indicator honest across reconnects.
  let countedOpen = false;

  const errorHandler = onError ?? defaultOnError(subscriptionName);

  function open(): void {
    const url = buildUrl(path, currentParams);
    const es = new EventSource(url, { withCredentials: true });
    source = es;

    es.addEventListener('open', () => {
      if (disposed) return;
      if (!countedOpen) {
        countedOpen = true;
        _notifySubscriptionOpen();
      }
    });

    es.addEventListener(eventName, (raw: Event) => {
      // A message can arrive between dispose() and the browser actually
      // closing the EventSource — drop it to keep dispose semantically
      // clean (no callbacks fire after dispose returns).
      if (disposed) return;
      // EventSource only fires MessageEvent for named events, but the DOM
      // typings declare a plain Event here so we narrow defensively.
      const msg = raw as MessageEvent<string>;
      try {
        const parsed: unknown = JSON.parse(msg.data);
        onMessage(parsed);
      } catch (err) {
        // A malformed payload from the server is a server bug, not a
        // transport problem — log it but keep the stream open so the next
        // (presumably well-formed) event still gets through.
        console.error(`[SSE: ${subscriptionName}] failed to parse message:`, err);
      }
    });

    es.addEventListener('error', (err: Event) => {
      if (disposed) return;
      // EventSource onerror fires when the connection drops; the browser
      // will then attempt to reconnect, firing `open` again on success. We
      // must mark the subscription closed here so the upcoming reopen
      // notifies wsStatus on the 0->1 edge — otherwise countedOpen stays
      // true and the status machine never observes that we reconnected.
      if (countedOpen) {
        countedOpen = false;
        _notifySubscriptionClose();
      }
      _notifySubscriptionError(err);
      errorHandler(err);
    });
  }

  function close(): void {
    if (source) {
      source.close();
      source = null;
    }
    if (countedOpen) {
      countedOpen = false;
      _notifySubscriptionClose();
    }
  }

  open();

  return {
    dispose(): void {
      if (disposed) return;
      disposed = true;
      close();
    },
    rebind(newParams: Record<string, string | string[] | undefined>): void {
      if (disposed) return;
      currentParams = newParams;
      close();
      // Re-open under the new params, preserving handlers via closure.
      // disposed is still false so open() will install fresh listeners on
      // the new EventSource.
      open();
    },
  };
}

/**
 * Default error handler used when callers don't supply one. Matches the
 * "log-only, don't throw" convention from lib/subscriptionHelper.ts so UI
 * surfaces aren't spammed with reconnect-cycle errors that the browser is
 * already healing.
 */
function defaultOnError(subscriptionName: string): (err: Event) => void {
  return (err) => {
    console.error(`[SSE: ${subscriptionName}] Error:`, err);
  };
}
