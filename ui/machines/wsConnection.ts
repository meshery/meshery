import { setup, fromCallback, assign } from 'xstate';
import { subscriptionClient } from '../lib/relayEnvironment';

/**
 * XState machine managing the graphql-ws WebSocket connection lifecycle.
 * Tracks connection state (connecting, connected, reconnecting, disconnected)
 * and emits events for UI consumption (connection status indicators, error toasts).
 */

export const WS_CONNECTION_EVENTS = {
  CONNECTED: 'WS_CONNECTED',
  DISCONNECTED: 'WS_DISCONNECTED',
  ERROR: 'WS_ERROR',
  RECONNECTING: 'WS_RECONNECTING',
} as const;

const wsMonitorActor = fromCallback(({ sendBack }) => {
  if (!subscriptionClient || typeof window === 'undefined') {
    return;
  }

  const unsubConnected = subscriptionClient.on('connected', () => {
    sendBack({ type: WS_CONNECTION_EVENTS.CONNECTED });
  });

  const unsubClosed = subscriptionClient.on('closed', (event) => {
    sendBack({ type: WS_CONNECTION_EVENTS.DISCONNECTED, data: { event } });
  });

  const unsubError = subscriptionClient.on('error', (error) => {
    sendBack({ type: WS_CONNECTION_EVENTS.ERROR, data: { error } });
  });

  return () => {
    unsubConnected();
    unsubClosed();
    unsubError();
  };
});

export const wsConnectionMachine = setup({
  types: {
    context: {} as {
      retryCount: number;
      lastError: unknown | null;
      lastConnectedAt: number | null;
    },
  },
  actors: {
    wsMonitor: wsMonitorActor,
  },
  actions: {
    incrementRetry: assign({
      retryCount: ({ context }) => context.retryCount + 1,
    }),
    resetRetry: assign({
      retryCount: 0,
      lastConnectedAt: () => Date.now(),
    }),
    recordError: assign({
      lastError: (_, params: { error: unknown }) => params.error,
    }),
  },
}).createMachine({
  id: 'wsConnection',
  initial: 'connecting',
  context: {
    retryCount: 0,
    lastError: null,
    lastConnectedAt: null,
  },

  invoke: {
    src: 'wsMonitor',
    id: 'wsMonitor',
  },

  states: {
    connecting: {
      on: {
        [WS_CONNECTION_EVENTS.CONNECTED]: {
          target: 'connected',
          actions: 'resetRetry',
        },
        [WS_CONNECTION_EVENTS.ERROR]: {
          target: 'disconnected',
          actions: 'incrementRetry',
        },
      },
    },
    connected: {
      on: {
        [WS_CONNECTION_EVENTS.DISCONNECTED]: {
          target: 'reconnecting',
        },
        [WS_CONNECTION_EVENTS.ERROR]: {
          target: 'reconnecting',
        },
      },
    },
    reconnecting: {
      entry: 'incrementRetry',
      on: {
        [WS_CONNECTION_EVENTS.CONNECTED]: {
          target: 'connected',
          actions: 'resetRetry',
        },
        [WS_CONNECTION_EVENTS.ERROR]: {
          target: 'disconnected',
        },
      },
    },
    disconnected: {
      on: {
        [WS_CONNECTION_EVENTS.CONNECTED]: {
          target: 'connected',
          actions: 'resetRetry',
        },
      },
    },
  },
});
