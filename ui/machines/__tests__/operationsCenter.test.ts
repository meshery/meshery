/* eslint-disable max-lines */
import { describe, expect, it, vi, beforeEach } from 'vitest';
import { createActor } from 'xstate';

// The operationsCenter machine subscribes to the events SSE stream and bridges
// each event into the redux store + the snackbar notifier. We mock the
// subscription module, the redux store, the rtk-query API and the
// NotificationCenter constants so we can drive events and assert side
// effects.

type EventCallback = (event: unknown) => void;
type ErrorCallback = (err: unknown) => void;

const hoisted = vi.hoisted(() => ({
  subscriptionState: {
    onEvent: undefined as EventCallback | undefined,
    onError: undefined as ErrorCallback | undefined,
    dispose: vi.fn(),
  },
  dispatch: vi.fn(),
  pushEvent: vi.fn((p: unknown) => ({ type: 'events/pushEvent', payload: p })),
}));
const subscriptionState = hoisted.subscriptionState;
const dispatch = hoisted.dispatch;
const pushEvent = hoisted.pushEvent;

vi.mock('lib/eventsSubscription', () => ({
  subscribeToEvents: (next: EventCallback, error: ErrorCallback) => {
    hoisted.subscriptionState.onEvent = next;
    hoisted.subscriptionState.onError = error;
    return { dispose: hoisted.subscriptionState.dispose };
  },
}));

vi.mock('../../store', () => ({
  store: { dispatch: hoisted.dispatch },
}));

vi.mock('@/store/slices/events', () => ({
  pushEvent: (p: unknown) => hoisted.pushEvent(p),
}));

vi.mock('../../rtk-query', () => ({
  api: { util: { invalidateTags: vi.fn((tags) => ({ type: 'api/invalidate', tags })) } },
}));

vi.mock('@/rtk-query/notificationCenter', () => ({
  PROVIDER_TAGS: { EVENT: 'Event' },
}));

vi.mock('@/components/layout/NotificationCenter/constants', () => ({
  SEVERITY_TO_NOTIFICATION_TYPE_MAPPING: {
    informational: 'info',
    error: 'error',
    warning: 'warning',
    success: 'success',
  },
}));

import {
  operationsCenterActor,
  OPERATION_CENTER_EVENTS,
  shouldNotifyEvent,
  extractConnectionStatus,
  extractConnectionId,
  MAX_TRACKED_CONNECTIONS,
} from '../operationsCenter';

const startActor = (notify = vi.fn()) => {
  const actor = createActor(operationsCenterActor, { input: { notify } });
  actor.start();
  return { actor, notify };
};

describe('operationsCenter machine', () => {
  beforeEach(() => {
    dispatch.mockClear();
    pushEvent.mockClear();
    subscriptionState.dispose.mockClear();
    subscriptionState.onEvent = undefined;
    subscriptionState.onError = undefined;
  });

  it('boots into idle after init and spawns the subscription actor', () => {
    const { actor } = startActor();
    const snap = actor.getSnapshot();
    expect(snap.value).toBe('idle');
    expect(subscriptionState.onEvent).toBeTypeOf('function');
    expect(subscriptionState.onError).toBeTypeOf('function');
    actor.stop();
  });

  it('exposes the canonical event constants', () => {
    expect(OPERATION_CENTER_EVENTS).toEqual({
      EVENT_RECEIVED_FROM_SERVER: 'EVENT_RECEIVED_FROM_SERVER',
      ERROR_OCCURRED_IN_SUBSCRIPTION: 'ERROR_OCCURRED_IN_SUBSCRIPTION',
    });
  });

  it('stores received events in redux and triggers the notifier', () => {
    const { actor, notify } = startActor();

    const event = {
      id: 'evt-1',
      description: 'A thing happened',
      severity: 'informational',
    };

    actor.send({
      type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
      data: { event },
    });

    // Redux dispatch fired with pushEvent action
    expect(pushEvent).toHaveBeenCalledWith(event);
    expect(dispatch).toHaveBeenCalledWith({ type: 'events/pushEvent', payload: event });

    // notify called with the mapped event_type
    expect(notify).toHaveBeenCalledWith({
      message: 'A thing happened',
      event_type: 'info',
      id: 'evt-1',
      showInNotificationCenter: true,
    });

    actor.stop();
  });

  it('emits received events back to the parent for consumers', () => {
    const { actor } = startActor();
    const seen: unknown[] = [];
    actor.on('*', (event) => {
      seen.push(event);
    });
    const payload = {
      id: 'evt-2',
      description: 'desc',
      severity: 'error',
    };
    actor.send({
      type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
      data: { event: payload },
    });
    const emitted = seen.find(
      (e: { type?: string }) => e?.type === OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
    );
    expect(emitted).toBeDefined();
    actor.stop();
  });

  it('forwards subscription stream events into the machine via the next callback', () => {
    const { actor } = startActor();
    const event = { id: 'live', description: 'streamed', severity: 'warning' };
    subscriptionState.onEvent?.(event);
    expect(pushEvent).toHaveBeenCalledWith(event);
    actor.stop();
  });

  it('ignores empty stream callbacks but logs', () => {
    const errSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    const { actor } = startActor();
    subscriptionState.onEvent?.(null);
    expect(pushEvent).not.toHaveBeenCalled();
    expect(errSpy).toHaveBeenCalled();
    errSpy.mockRestore();
    actor.stop();
  });

  it('respawns the subscription actor on ERROR_OCCURRED_IN_SUBSCRIPTION', () => {
    const { actor } = startActor();
    // Send error event to trigger respawn
    actor.send({
      type: OPERATION_CENTER_EVENTS.ERROR_OCCURRED_IN_SUBSCRIPTION,
      data: { error: new Error('socket') },
    });

    // After respawn, subscription handlers should still be wired up
    expect(subscriptionState.onEvent).toBeDefined();
    expect(subscriptionState.onError).toBeDefined();

    actor.stop();
  });

  describe('connection lifecycle and notification deduplication (Issue #9392)', () => {
    it('suppresses toasts for generic connection update notifications while storing in redux', () => {
      const { actor, notify } = startActor();
      const genericEvent = {
        id: 'evt-generic-1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'Connection colima-meshery-9392 updated.',
      };

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: genericEvent },
      });

      // Stored in redux for notification center drawer history
      expect(pushEvent).toHaveBeenCalledWith(genericEvent);
      expect(dispatch).toHaveBeenCalledWith({ type: 'events/pushEvent', payload: genericEvent });

      // But notify (toast) must NOT be called for generic connection update
      expect(notify).not.toHaveBeenCalled();
      actor.stop();
    });

    it('suppresses toasts for short generic "Connection updated." notification', () => {
      const { actor, notify } = startActor();
      const genericEvent = {
        id: 'evt-generic-2',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'Connection updated.',
      };

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: genericEvent },
      });

      expect(pushEvent).toHaveBeenCalledWith(genericEvent);
      expect(notify).not.toHaveBeenCalled();
      actor.stop();
    });

    it('toasts legitimate connection lifecycle state transitions', () => {
      const { actor, notify } = startActor();
      const disconnectEvent = {
        id: 'evt-disc-1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'kubernetes connection changed to disconnected',
        metadata: {
          previousStatus: 'connected',
          currentStatus: 'disconnected',
        },
      };

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: disconnectEvent },
      });

      expect(pushEvent).toHaveBeenCalledWith(disconnectEvent);
      expect(notify).toHaveBeenCalledWith({
        message: 'kubernetes connection changed to disconnected',
        event_type: 'info',
        id: 'evt-disc-1',
        showInNotificationCenter: true,
      });
      actor.stop();
    });

    it('suppresses repeated identical lifecycle state transitions on the same connection', () => {
      const { actor, notify } = startActor();
      const disconnectEvent1 = {
        id: 'evt-disc-1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'kubernetes connection changed to disconnected',
        metadata: { currentStatus: 'disconnected' },
      };
      const disconnectEvent2 = {
        id: 'evt-disc-2',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'kubernetes connection changed to disconnected',
        metadata: { currentStatus: 'disconnected' },
      };

      // First transition: should toast
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: disconnectEvent1 },
      });
      expect(notify).toHaveBeenCalledTimes(1);

      // Repeated identical status event on the same connection: stored in redux, but toast suppressed
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: disconnectEvent2 },
      });
      expect(pushEvent).toHaveBeenCalledWith(disconnectEvent2);
      expect(notify).toHaveBeenCalledTimes(1); // Still 1

      actor.stop();
    });

    it('allows transitions between different states on the same connection (disconnect -> connect)', () => {
      const { actor, notify } = startActor();
      const disconnectEvent = {
        id: 'evt-disc-1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'kubernetes connection changed to disconnected',
        metadata: { currentStatus: 'disconnected' },
      };
      const reconnectEvent = {
        id: 'evt-conn-1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'kubernetes connection changed to connected',
        metadata: { currentStatus: 'connected' },
      };

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: disconnectEvent },
      });
      expect(notify).toHaveBeenCalledTimes(1);
      expect(notify).toHaveBeenLastCalledWith(
        expect.objectContaining({
          message: 'kubernetes connection changed to disconnected',
        }),
      );

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: reconnectEvent },
      });
      expect(notify).toHaveBeenCalledTimes(2);
      expect(notify).toHaveBeenLastCalledWith(
        expect.objectContaining({
          message: 'kubernetes connection changed to connected',
        }),
      );

      actor.stop();
    });

    it('tracks state independently across different connections', () => {
      const { actor, notify } = startActor();
      const cluster1Event = {
        id: 'evt-c1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'kubernetes connection changed to connected',
        metadata: { currentStatus: 'connected' },
      };
      const cluster2Event = {
        id: 'evt-c2',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-2',
        description: 'kubernetes connection changed to connected',
        metadata: { currentStatus: 'connected' },
      };

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: cluster1Event },
      });
      expect(notify).toHaveBeenCalledTimes(1);

      // cluster-2 is a different connection; should toast even with same status
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: cluster2Event },
      });
      expect(notify).toHaveBeenCalledTimes(2);

      actor.stop();
    });

    it('always preserves legitimate MeshSync controller notifications', () => {
      const { actor, notify } = startActor();
      const meshsyncEvent = {
        id: 'evt-meshsync-1',
        category: 'controller',
        action: 'update',
        severity: 'informational',
        description: 'MeshSync connected in embedded mode',
      };

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: meshsyncEvent },
      });

      expect(pushEvent).toHaveBeenCalledWith(meshsyncEvent);
      expect(notify).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'MeshSync connected in embedded mode',
        }),
      );

      actor.stop();
    });

    it('always preserves connection errors and warnings', () => {
      const { actor, notify } = startActor();
      const errorEvent = {
        id: 'evt-err-1',
        category: 'connection',
        action: 'update',
        severity: 'error',
        actedUpon: 'cluster-conn-1',
        description: 'Failed to deploy Meshery Operator',
      };

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: errorEvent },
      });

      expect(notify).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'Failed to deploy Meshery Operator',
          event_type: 'error',
        }),
      );

      actor.stop();
    });

    it('suppresses toasts for condition recurrences while preserving resolution toasts', () => {
      const { actor, notify } = startActor();
      const recurrenceEvent = {
        id: 'evt-cond-1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        correlationKey: 'unreachable:cluster-conn-1',
        occurrenceCount: 2,
        description: 'Kubernetes cluster is unreachable',
      };
      const resolutionEvent = {
        id: 'evt-cond-1',
        category: 'connection',
        action: 'update',
        severity: 'success',
        correlationKey: 'unreachable:cluster-conn-1',
        resolvedAt: '2026-11-28T19:00:00Z',
        description: 'kubernetes connection is reachable again and changed to connected',
      };

      // Recurrence: stored in redux, toast suppressed
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: recurrenceEvent },
      });
      expect(pushEvent).toHaveBeenCalledWith(recurrenceEvent);
      expect(notify).not.toHaveBeenCalled();

      // Resolution: state change, should toast
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: resolutionEvent },
      });
      expect(pushEvent).toHaveBeenCalledWith(resolutionEvent);
      expect(notify).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'kubernetes connection is reachable again and changed to connected',
        }),
      );

      actor.stop();
    });

    it('toasts connection deletion and cleans up state from the map', () => {
      const { actor, notify } = startActor();
      const connectEvent = {
        id: 'evt-c1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'kubernetes connection changed to connected',
        metadata: { currentStatus: 'connected' },
      };
      const deleteEvent = {
        id: 'evt-del-1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'kubernetes connection was deleted',
        metadata: { currentStatus: 'deleted' },
      };
      const recreateConnectEvent = {
        id: 'evt-c2',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'kubernetes connection changed to connected',
        metadata: { currentStatus: 'connected' },
      };

      // 1. Connects
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: connectEvent },
      });
      expect(notify).toHaveBeenCalledTimes(1);

      // 2. Deletes -> Should toast deletion and clean up state
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: deleteEvent },
      });
      expect(notify).toHaveBeenCalledTimes(2);
      expect(notify).toHaveBeenLastCalledWith(
        expect.objectContaining({
          message: 'kubernetes connection was deleted',
        }),
      );

      // 3. Re-creating connection with same ID connects cleanly and toasts
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: recreateConnectEvent },
      });
      expect(notify).toHaveBeenCalledTimes(3);
      expect(notify).toHaveBeenLastCalledWith(
        expect.objectContaining({
          message: 'kubernetes connection changed to connected',
        }),
      );

      actor.stop();
    });

    it('handles failure, recovery, and subsequent failure cycles correctly', () => {
      const { actor, notify } = startActor();
      const initialConnect = {
        id: 'evt-conn-1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-1',
        description: 'kubernetes connection changed to connected',
        metadata: { currentStatus: 'connected' },
      };
      const failureEvent = {
        id: 'evt-fail-1',
        category: 'connection',
        action: 'update',
        severity: 'error',
        actedUpon: 'cluster-conn-1',
        description: 'Kubernetes API server unreachable',
        metadata: { currentStatus: 'notfound' },
      };
      const recoveryEvent = {
        id: 'evt-recov-1',
        category: 'connection',
        action: 'update',
        severity: 'success',
        actedUpon: 'cluster-conn-1',
        description: 'kubernetes connection is reachable again and changed to connected',
        metadata: { currentStatus: 'connected' },
        resolvedAt: '2026-11-28T19:00:00Z',
      };
      const subsequentFailure = {
        id: 'evt-fail-2',
        category: 'connection',
        action: 'update',
        severity: 'error',
        actedUpon: 'cluster-conn-1',
        description: 'Connection failed again',
        metadata: { currentStatus: 'notfound' },
      };

      // 1. Initial connect toasts
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: initialConnect },
      });
      expect(notify).toHaveBeenCalledTimes(1);

      // 2. Failure toasts (error severity)
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: failureEvent },
      });
      expect(notify).toHaveBeenCalledTimes(2);

      // 3. Recovery toasts
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: recoveryEvent },
      });
      expect(notify).toHaveBeenCalledTimes(3);
      expect(notify).toHaveBeenLastCalledWith(
        expect.objectContaining({
          message: 'kubernetes connection is reachable again and changed to connected',
        }),
      );

      // 4. Subsequent failure toasts
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: subsequentFailure },
      });
      expect(notify).toHaveBeenCalledTimes(4);

      actor.stop();
    });

    it('safely handles missing connection IDs and metadata without crashing or dropping notifications', () => {
      const { actor, notify } = startActor();
      const nilIdEvent = {
        id: 'evt-nil-id',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: '00000000-0000-0000-0000-000000000000',
        description: 'Something happened with an unknown connection',
      };
      const missingMetadataEvent = {
        id: 'evt-no-meta',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        description: 'General connection update notice',
      };

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: nilIdEvent },
      });
      expect(notify).toHaveBeenCalledTimes(1);

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: missingMetadataEvent },
      });
      expect(notify).toHaveBeenCalledTimes(2);

      actor.stop();
    });

    it('does not suppress condition recurrences for non-connection categories', () => {
      const { actor, notify } = startActor();
      const designRecurrence = {
        id: 'evt-design-1',
        category: 'design',
        action: 'validate',
        severity: 'informational',
        correlationKey: 'validation:design-123',
        occurrenceCount: 3,
        description: 'Design validation in progress',
      };

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: designRecurrence },
      });

      // Because it is NOT connection category, Rule 2 does not suppress it
      expect(notify).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'Design validation in progress',
        }),
      );

      actor.stop();
    });

    it('bounds tracked connection map size to prevent unbounded growth', () => {
      const map = new Map<string, string>();
      for (let i = 0; i < MAX_TRACKED_CONNECTIONS + 5; i++) {
        const event = {
          category: 'connection',
          action: 'update',
          severity: 'informational',
          actedUpon: `conn-${i}`,
          description: `kubernetes connection changed to connected`,
          metadata: { currentStatus: 'connected' },
        };
        shouldNotifyEvent(event, map);
      }
      expect(map.size).toBeLessThanOrEqual(MAX_TRACKED_CONNECTIONS);
      expect(map.has(`conn-${MAX_TRACKED_CONNECTIONS + 4}`)).toBe(true);
      expect(map.has('conn-0')).toBe(false); // Oldest entry pruned
    });

    it('toasts condition resolution even when target status matches already-cached status (Review Issue 1)', () => {
      const { actor, notify } = startActor();
      // 1. Initial connection reaches 'connected'
      const initialConnect = {
        id: 'evt-c-1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-cached',
        description: 'kubernetes connection changed to connected',
        metadata: { currentStatus: 'connected' },
      };
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: initialConnect },
      });
      expect(notify).toHaveBeenCalledTimes(1);

      // 2. Condition resolution arrives with resolvedAt and target status 'connected' (same as cached)
      const resolutionSameStatus = {
        id: 'evt-res-1',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-cached',
        description: 'kubernetes connection is reachable again and changed to connected',
        resolvedAt: '2026-11-28T19:30:00Z',
        metadata: { currentStatus: 'connected' },
      };
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: resolutionSameStatus },
      });
      // MUST NOT be suppressed by lastStatus === targetStatus
      expect(notify).toHaveBeenCalledTimes(2);
      expect(notify).toHaveBeenLastCalledWith(
        expect.objectContaining({
          message: 'kubernetes connection is reachable again and changed to connected',
        }),
      );

      actor.stop();
    });

    it('preserves non-generic connection updates (Review Issue 2)', () => {
      const { actor, notify } = startActor();
      // Informational connection update with non-generic description (e.g. specific custom action)
      const customActionUpdate = {
        id: 'evt-custom-upd',
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-custom',
        description: 'Cluster credentials rotated successfully',
        metadata: { name: 'production-cluster' },
      };

      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: { event: customActionUpdate },
      });
      // Should toast because it is not a generic "Connection updated." message
      expect(notify).toHaveBeenCalledTimes(1);
      expect(notify).toHaveBeenCalledWith(
        expect.objectContaining({
          message: 'Cluster credentials rotated successfully',
        }),
      );

      actor.stop();
    });

    it('validates status extraction against canonical lifecycle states and ignores arbitrary metadata (Review Issue 3)', () => {
      const map = new Map<string, string>();
      // Arbitrary metadata should NOT be extracted as a valid lifecycle status
      const arbitraryEvent = {
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-arb',
        description: 'Custom message without lifecycle keyword',
        metadata: { currentStatus: 'someRandomArbitraryStatus' },
      };
      expect(extractConnectionStatus(arbitraryEvent)).toBeNull();

      // Canonical status in currentStatus or status should be accepted
      const validEvent = {
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: 'cluster-conn-valid',
        description: 'Status update',
        metadata: { currentStatus: 'CONNECTED' },
      };
      expect(extractConnectionStatus(validEvent)).toBe('connected');

      // Nil UUID or malformed actedUpon should be rejected by extractConnectionId
      const nilUuidEvent = {
        category: 'connection',
        action: 'update',
        severity: 'informational',
        actedUpon: '  00000000-0000-0000-0000-000000000000  ',
      };
      expect(extractConnectionId(nilUuidEvent)).toBeNull();
    });

    it('preserves genuine transitions arriving out-of-band or after SSE reconnect (Review Issue 5)', () => {
      const { actor, notify } = startActor();
      // 1. Initial state connected
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: {
          event: {
            id: 'evt-rec-1',
            category: 'connection',
            action: 'update',
            severity: 'informational',
            actedUpon: 'cluster-conn-reconnect',
            description: 'kubernetes connection changed to connected',
            metadata: { currentStatus: 'connected' },
          },
        },
      });
      expect(notify).toHaveBeenCalledTimes(1);

      // 2. Out-of-band state change to disconnected after reconnect
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: {
          event: {
            id: 'evt-rec-2',
            category: 'connection',
            action: 'update',
            severity: 'informational',
            actedUpon: 'cluster-conn-reconnect',
            description: 'kubernetes connection changed to disconnected',
            metadata: { currentStatus: 'disconnected' },
          },
        },
      });
      expect(notify).toHaveBeenCalledTimes(2);

      actor.stop();
    });
    it('resets connectionStates cache on SSE subscription disconnect and reconnects cleanly', () => {
      const { actor, notify } = startActor();
      // 1. Initial state: connection changes to connected
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: {
          event: {
            id: 'evt-recon-1',
            category: 'connection',
            action: 'update',
            severity: 'informational',
            actedUpon: 'conn-missed-cycle',
            description: 'kubernetes connection changed to connected',
            metadata: { currentStatus: 'connected' },
          },
        },
      });
      expect(notify).toHaveBeenCalledTimes(1);

      // 2. SSE subscription drops (e.g. network disconnect, missed an intermediate disconnected event)
      actor.send({
        type: OPERATION_CENTER_EVENTS.ERROR_OCCURRED_IN_SUBSCRIPTION,
        data: { error: new Error('SSE connection lost') },
      });

      // 3. Reconnect occurs: genuine connected event arrives.
      // Because cache was reset on subscription error, the new connected event is NOT suppressed!
      actor.send({
        type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
        data: {
          event: {
            id: 'evt-recon-2',
            category: 'connection',
            action: 'update',
            severity: 'informational',
            actedUpon: 'conn-missed-cycle',
            description: 'kubernetes connection changed to connected',
            metadata: { currentStatus: 'connected' },
          },
        },
      });
      expect(notify).toHaveBeenCalledTimes(2);
      expect(notify).toHaveBeenLastCalledWith(
        expect.objectContaining({
          message: 'kubernetes connection changed to connected',
        }),
      );

      actor.stop();
    });
  });
});
