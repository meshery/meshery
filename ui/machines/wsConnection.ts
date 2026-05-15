import { setup, fromCallback, assign } from 'xstate';
import * as wsStatus from '../lib/wsStatus';

/**
 * XState machine managing the live server connection lifecycle for the
 * SSE migration. Historically wired to the graphql-ws subscriptionClient;
 * now wired to lib/wsStatus, which aggregates SSE subscription state under
 * the same `connected | closed | error` event surface. Consumers of the
 * machine (status indicators, reconnect toasts) see no shape change.
 */

export const WS_CONNECTION_EVENTS = {
  CONNECTED: 'WS_CONNECTED',
  DISCONNECTED: 'WS_DISCONNECTED',
  ERROR: 'WS_ERROR',
  RECONNECTING: 'WS_RECONNECTING',
} as const;

const wsMonitorActor = fromCallback(({ sendBack }) => {
  // SSR guard: wsStatus events fire from EventSource callbacks which only
  // exist in the browser. Returning a no-op cleanup keeps the actor happy
  // when the machine is hydrated server-side.
  if (typeof window === 'undefined') {
    return;
  }

  const unsubConnected = wsStatus.on('connected', () => {
    sendBack({ type: WS_CONNECTION_EVENTS.CONNECTED });
  });

  const unsubClosed = wsStatus.on('closed', (event) => {
    sendBack({ type: WS_CONNECTION_EVENTS.DISCONNECTED, data: { event } });
  });

  const unsubError = wsStatus.on('error', (error) => {
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
