/**
 * wsStatus.ts
 *
 * Compatibility shim for the in-flight WebSocket-to-SSE migration. The
 * existing XState machine in machines/wsConnection.ts was wired to the
 * graphql-ws subscriptionClient on/off contract (`connected | closed | error`)
 * so it could light up a connection status indicator. As routes migrate from
 * graphql-ws to SSE we still want that indicator to reflect "do we have a
 * live server connection at all?".
 *
 * This module exposes the same `on(event, cb) => unsubscribe` / `off(event, cb)`
 * surface, but the underlying signal is the number of open SSE subscriptions
 * managed by sseClient.ts. When the first subscription opens we emit
 * `connected`; when the last one closes we emit `closed`; any subscription's
 * onerror fires `error`. This keeps consumers of the XState machine
 * (connection toasts, header indicators) working unchanged across the
 * migration without forcing every UI surface to subscribe to a different API.
 */
export type WsStatusEvent = 'connected' | 'closed' | 'error';
export type WsStatusListener<E extends WsStatusEvent = WsStatusEvent> = (
  payload?: E extends 'error' ? unknown : Event | undefined,
) => void;

type ListenerMap = {
  [E in WsStatusEvent]: Set<WsStatusListener<E>>;
};

const listeners: ListenerMap = {
  connected: new Set(),
  closed: new Set(),
  error: new Set(),
};

// activeCount is the number of currently-open SSE subscriptions. We treat
// the aggregate as "connected" whenever it is > 0 and "closed" when it
// returns to 0. We don't try to coalesce concurrent transitions: if a caller
// opens two subscriptions in the same tick the second open is a no-op for
// status purposes, which matches the graphql-ws behaviour where there is a
// single underlying socket.
let activeCount = 0;

function emit<E extends WsStatusEvent>(
  event: E,
  payload?: E extends 'error' ? unknown : Event | undefined,
): void {
  // Snapshot the listener set before iterating: a handler is allowed to
  // call off() on itself, which would mutate the live Set mid-iteration.
  const snapshot = Array.from(listeners[event]) as WsStatusListener<E>[];
  for (const cb of snapshot) {
    try {
      cb(payload);
    } catch (err) {
      // A buggy listener shouldn't poison the others. Surface via console
      // so it shows up in dev tools but keep going.
      console.error('[wsStatus] listener for', event, 'threw:', err);
    }
  }
}

/**
 * Subscribe to a wsStatus event. Returns an unsubscribe function so callers
 * can use it inline (e.g. inside xstate fromCallback).
 */
export function on<E extends WsStatusEvent>(event: E, cb: WsStatusListener<E>): () => void {
  listeners[event].add(cb as WsStatusListener);
  return () => off(event, cb);
}

/**
 * Remove a previously-registered listener.
 */
export function off<E extends WsStatusEvent>(event: E, cb: WsStatusListener<E>): void {
  listeners[event].delete(cb as WsStatusListener);
}

// --- internal hooks for sseClient -------------------------------------------
// These are not exported in the public type surface — they are accessed from
// sseClient.ts via the module-level binding because the two modules are
// tightly coupled by design.

/**
 * @internal Called by sseClient when a new EventSource transitions to OPEN.
 * Fires `connected` only on the 0 -> 1 edge.
 */
export function _notifySubscriptionOpen(): void {
  activeCount++;
  if (activeCount === 1) {
    emit('connected');
  }
}

/**
 * @internal Called by sseClient when an EventSource is disposed or naturally
 * closes. Fires `closed` only on the 1 -> 0 edge.
 */
export function _notifySubscriptionClose(): void {
  if (activeCount === 0) return;
  activeCount--;
  if (activeCount === 0) {
    emit('closed');
  }
}

/**
 * @internal Called by sseClient on EventSource onerror. Always fires —
 * the XState machine decides whether to transition to reconnecting.
 */
export function _notifySubscriptionError(err: Event): void {
  emit('error', err);
}

/**
 * @internal Test-only escape hatch to reset internal counters. Production
 * callers should never need this; exported with an underscore prefix to
 * keep it out of typical autocomplete suggestions.
 */
export function _resetForTests(): void {
  activeCount = 0;
  listeners.connected.clear();
  listeners.closed.clear();
  listeners.error.clear();
}

/**
 * @internal Read-only view of the current count, for diagnostics and tests.
 */
export function _activeCountForTests(): number {
  return activeCount;
}
