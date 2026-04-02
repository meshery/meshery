import { createClient } from 'graphql-ws';
import { Environment, Network, Observable, RecordSource, Store } from 'relay-runtime';
import { promisifiedDataFetch } from './data-fetch';

function fetchQuery(operation, variables) {
  return promisifiedDataFetch('/api/system/graphql/query', {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query: operation.text,
      variables,
    }),
  });
}

export let subscriptionClient;

if (typeof window !== 'undefined') {
  const isWss = window.location.protocol === 'https:';
  const wsProtocol = isWss ? 'wss://' : 'ws://';
  subscriptionClient = createClient({
    url: wsProtocol + window.location.host + '/api/system/graphql/query',
    retryAttempts: Infinity,
    shouldRetry: () => true,
    retryWait: (retries) => {
      const baseDelay = 1000;
      const maxDelay = 30000;
      const delay = Math.min(baseDelay * 2 ** retries, maxDelay);
      return new Promise((resolve) => setTimeout(resolve, delay));
    },
    keepAlive: 30_000,
    on: {
      closed: (event) => {
        if (event && typeof event !== 'number') {
          console.warn('[GraphQL WS] Connection closed:', event);
        }
      },
      error: (error) => {
        console.error('[GraphQL WS] Connection error:', error);
      },
      connected: () => {
        console.info('[GraphQL WS] Connected');
      },
    },
  });
}

function fetchOrSubscribe(operation, variables) {
  return Observable.create((sink) => {
    if (!operation.text) {
      return sink.error(new Error('Operation text cannot be empty'));
    }
    return subscriptionClient.subscribe(
      {
        operationName: operation.name,
        query: operation.text,
        variables,
      },
      sink,
    );
  });
}

export const serializeRelayEnvironment = (environment) => {
  return environment.getStore().getSource().toJSON();
};

// Singleton environment for client-side use.
// Avoids creating a new Environment + RecordSource + Store per query,
// which defeats Relay's normalized cache.
let clientEnvironment: Environment | null = null;

export const createRelayEnvironment = (records = {}) => {
  // Server-side: always create a fresh environment per request
  if (typeof window === 'undefined') {
    return new Environment({
      store: new Store(new RecordSource(records)),
      network: Network.create(fetchQuery, fetchOrSubscribe),
    });
  }

  // Client-side: reuse singleton
  if (!clientEnvironment) {
    clientEnvironment = new Environment({
      store: new Store(new RecordSource(records)),
      network: Network.create(fetchQuery, fetchOrSubscribe),
    });
  }

  return clientEnvironment;
};
