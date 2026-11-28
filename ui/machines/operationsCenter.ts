import { SEVERITY_TO_NOTIFICATION_TYPE_MAPPING } from '@/components/layout/NotificationCenter/constants';
import { subscribeToEvents } from 'lib/eventsSubscription';
import { emit, fromCallback, setup, spawnChild, stopChild } from 'xstate';
import { store } from '../store';
import { pushEvent } from '@/store/slices/events';
import { api as mesheryApi } from '../rtk-query';
import { PROVIDER_TAGS } from '@/rtk-query/notificationCenter';

export interface OperationsCenterMeshkitEvent {
  id?: string;
  category?: string;
  action?: string;
  severity?: string;
  description?: string;
  actedUpon?: string;
  correlationKey?: string;
  resolvedAt?: string | null;
  occurrenceCount?: number;
  isRecurrence?: boolean;
  metadata?: Record<string, unknown>;
}

export interface OperationsCenterNotification {
  message: string;
  event_type: string;
  id?: string;
  showInNotificationCenter?: boolean;
}

export type OperationsCenterNotifyCallback = (notification: OperationsCenterNotification) => void;

/**
 * Supported canonical connection lifecycle states per Meshery connection schemas & state machine
 * (schemas/models/v1beta3/connection and server/models/connections/connections.go).
 */
export const CANONICAL_CONNECTION_STATES = new Set([
  'connected',
  'disconnected',
  'discovered',
  'registered',
  'ignored',
  'maintenance',
  'notfound',
  'deleted',
]);

/**
 * Generic connection update events emitted by server REST endpoints (e.g. UpdateConnectionById)
 * report "Connection <name> updated." or "Connection updated." without specific lifecycle state.
 * These are emitted alongside NotifySmOfConnectionStatusChange on status changes and direct UI mutation toasts.
 */
const GENERIC_CONNECTION_UPDATE_PATTERN = /^Connection(?:\s+.*)?\s+updated\.?$/i;

/**
 * Lifecycle state transition description pattern emitted by the state machine
 * e.g. "kubernetes connection changed to disconnected", "kubernetes connection was deleted",
 * "kubernetes connection is reachable again and changed to connected"
 */
const LIFECYCLE_STATUS_PATTERN =
  /(?:connection.*(?:changed to|status updated to|is reachable again and changed to)\s+(\w+)|connection\s+was\s+(\w+))/i;

/** Bound on tracked connections to prevent unbounded map growth in long-running UI sessions */
export const MAX_TRACKED_CONNECTIONS = 200;

const NIL_UUID = '00000000-0000-0000-0000-000000000000';

export function extractConnectionStatus(event: OperationsCenterMeshkitEvent): string | null {
  if (!event) {
    return null;
  }

  // Helper to validate and return canonical connection status
  const validateStatus = (raw: unknown): string | null => {
    if (typeof raw !== 'string') {
      return null;
    }
    const normalized = raw.trim().toLowerCase();
    return CANONICAL_CONNECTION_STATES.has(normalized) ? normalized : null;
  };

  if (event.metadata && typeof event.metadata === 'object') {
    const fromCurrent = validateStatus(event.metadata.currentStatus);
    if (fromCurrent) {
      return fromCurrent;
    }
    const fromStatus = validateStatus(event.metadata.status);
    if (fromStatus) {
      return fromStatus;
    }
  }

  if (typeof event.description === 'string') {
    const match = event.description.match(LIFECYCLE_STATUS_PATTERN);
    if (match) {
      const candidate = (match[1] || match[2] || '').trim().toLowerCase();
      if (CANONICAL_CONNECTION_STATES.has(candidate)) {
        return candidate;
      }
    }
  }

  return null;
}

export function extractConnectionId(event: OperationsCenterMeshkitEvent): string | null {
  if (!event) {
    return null;
  }

  const cleanId = (val: unknown): string | null => {
    if (typeof val !== 'string') {
      return null;
    }
    const trimmed = val.trim();
    if (!trimmed || trimmed === NIL_UUID) {
      return null;
    }
    return trimmed;
  };

  const idFromActedUpon = cleanId(event.actedUpon);
  if (idFromActedUpon) {
    return idFromActedUpon;
  }

  if (event.metadata && typeof event.metadata === 'object') {
    const idFromCamel = cleanId(event.metadata.connectionId);
    if (idFromCamel) {
      return idFromCamel;
    }
    const idFromSnake = cleanId(event.metadata.connection_id);
    if (idFromSnake) {
      return idFromSnake;
    }
  }

  return null;
}

function trackConnectionStatus(
  connectionStates: Map<string, string>,
  connectionId: string,
  status: string,
): void {
  // When a connection is deleted, remove it from the map so invalid state is not retained
  // and memory is released.
  if (status === 'deleted') {
    connectionStates.delete(connectionId);
    return;
  }
  // Enforce capacity to prevent unbounded growth in very large / long-running sessions
  if (connectionStates.size >= MAX_TRACKED_CONNECTIONS && !connectionStates.has(connectionId)) {
    const oldestKey = connectionStates.keys().next().value;
    if (oldestKey) {
      connectionStates.delete(oldestKey);
    }
  }
  connectionStates.set(connectionId, status);
}

export function shouldNotifyEvent(
  event: OperationsCenterMeshkitEvent | undefined | null,
  connectionStates: Map<string, string>,
): boolean {
  if (!event) {
    return false;
  }

  const severity = String(event.severity || '').toLowerCase();
  const isConnectionCategory = event.category?.toLowerCase() === 'connection';
  const connectionId = isConnectionCategory ? extractConnectionId(event) : null;
  const targetStatus = isConnectionCategory ? extractConnectionStatus(event) : null;

  // Rule 1: Always toast errors, warnings, and critical issues regardless of category
  if (severity === 'error' || severity === 'warning' || severity === 'critical') {
    // If this error/warning event also transitions connection status (e.g. to notfound/disconnected),
    // update the tracked state so subsequent recovery events are recognized as valid state changes.
    if (connectionId && targetStatus) {
      trackConnectionStatus(connectionStates, connectionId, targetStatus);
    }
    return true;
  }

  // Condition resolutions (event.resolvedAt) represent state recovery on an issue or condition.
  // They must always toast even if targetStatus matches the currently recorded state in connectionStates.
  if (isConnectionCategory && event.resolvedAt) {
    if (connectionId && targetStatus) {
      trackConnectionStatus(connectionStates, connectionId, targetStatus);
    }
    return true;
  }

  // Rule 2: Open condition recurrences on connections (repeat occurrences without resolution)
  // should not produce repeated toast popups. Scoped strictly to the connection category.
  if (
    isConnectionCategory &&
    event.correlationKey &&
    !event.resolvedAt &&
    ((event.occurrenceCount !== undefined && event.occurrenceCount > 1) || event.isRecurrence)
  ) {
    return false;
  }

  if (isConnectionCategory) {
    const description = typeof event.description === 'string' ? event.description.trim() : '';

    // Rule 3: Suppress generic connection update notifications ("Connection updated." / "Connection <name> updated.")
    // when action is update and severity is informational.
    // In server/handlers/connections_handlers.go:688, UpdateConnectionById unconditionally emits
    // "Connection %s updated." with Informational severity for all update calls.
    // In the UI, direct mutation hooks (e.g. useUpdateConnectionMutation) already display user-facing
    // toasts on mutation completion, and state-machine transitions emit specific lifecycle toasts
    // ("kubernetes connection changed to ..."). Emitting toasts for generic updates creates redundant
    // banner storms across both user-initiated and system reconciliations.
    // These events remain fully preserved in Redux store and Notification Center history.
    if (
      event.action === 'update' &&
      severity === 'informational' &&
      GENERIC_CONNECTION_UPDATE_PATTERN.test(description)
    ) {
      return false;
    }

    // Rule 4: Deduplicate repeated connection lifecycle state transitions for the same connection
    if (connectionId && targetStatus) {
      const lastStatus = connectionStates.get(connectionId);
      if (lastStatus === targetStatus) {
        // Redundant repeat of the exact same state for this connection
        return false;
      }
      // Informative state transition: record new state and toast
      trackConnectionStatus(connectionStates, connectionId, targetStatus);
      return true;
    }
  }

  // All other events (MeshSync controller notifications, component registrations, non-connection events, etc.) are notified
  return true;
}

export const OPERATION_CENTER_EVENTS = {
  EVENT_RECEIVED_FROM_SERVER: 'EVENT_RECEIVED_FROM_SERVER',
  ERROR_OCCURRED_IN_SUBSCRIPTION: 'ERROR_OCCURRED_IN_SUBSCRIPTION',
};
const events = {
  eventReceivedFromServer: (event) => ({
    type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
    data: {
      event,
    },
  }),
  errorOccurredInSubscription: (error) => ({
    type: 'ERROR_OCCURRED_IN_SUBSCRIPTION',
    data: {
      error,
    },
  }),
};

const subscriptionActor = fromCallback(({ sendBack }) => {
  const subscription = subscribeToEvents(
    (event) => {
      try {
        if (!event) {
          console.error('Invalid event received', event);
          return;
        }

        // The SSE stream delivers the raw meshkit Event JSON (camelCase), the
        // same shape the REST /api/system/events list endpoint returns, so it
        // is consumed as-is by the downstream UI.
        sendBack(events.eventReceivedFromServer(event));
      } catch (error) {
        console.error('[operationsCenter] An error occurred in processing event', error);
      }
    },
    (error) => {
      console.error('[operationsCenter] An error occurred in subscription to events', error);
      sendBack(events.errorOccurredInSubscription(error));
    },
  );

  return () => {
    subscription.dispose();
  };
});

export const operationsCenterActor = setup({
  types: {
    context: {} as {
      notify: OperationsCenterNotifyCallback;
      connectionStates: Map<string, string>;
    },
    input: {} as {
      notify: OperationsCenterNotifyCallback;
    },
  },
  actions: {
    resetConnectionStates: ({ context }) => {
      context.connectionStates.clear();
    },
    spawnSubscriptionActor: spawnChild(subscriptionActor, {
      id: 'subscriptionActor',
    }),
    storeInRedux: ({ event }) => {
      store.dispatch(pushEvent(event.data.event));
    },
    invalidateRtk: () => {
      store.dispatch(mesheryApi.util.invalidateTags([PROVIDER_TAGS.EVENT]));
    },
    notifyUI: ({ context, event }) => {
      const validatedEvent = event.data.event;
      if (!context.connectionStates) {
        context.connectionStates = new Map<string, string>();
      }
      if (!shouldNotifyEvent(validatedEvent, context.connectionStates)) {
        return;
      }
      context.notify({
        message: validatedEvent.description,
        event_type: SEVERITY_TO_NOTIFICATION_TYPE_MAPPING[validatedEvent.severity],
        id: validatedEvent.id,
        showInNotificationCenter: true,
      });
    },
    emitBack: emit(({ event }) => event),
  },
}).createMachine({
  id: 'operationsCenter',
  initial: 'init',
  context: ({ input }) => ({
    notify: input.notify,
    connectionStates: new Map<string, string>(),
  }),

  states: {
    init: {
      entry: 'spawnSubscriptionActor',
      always: 'idle',
    },
    idle: {
      on: {
        [OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER]: {
          actions: ['storeInRedux', 'notifyUI', 'emitBack'],
        },
        [OPERATION_CENTER_EVENTS.ERROR_OCCURRED_IN_SUBSCRIPTION]: {
          actions: [
            'resetConnectionStates',
            stopChild('subscriptionActor'),
            'spawnSubscriptionActor',
          ],
        },
      },
    },
  },
});
