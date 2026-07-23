import { Environment, Network, RecordSource, Store } from 'relay-runtime';
import { promisifiedDataFetch } from './data-fetch';

// Meshery Server no longer exposes a `type Subscription`, so this environment is
// request/response only — there is no graphql-ws client and no subscribe handler.
// Real-time surfaces use Server-Sent Events (see lib/eventsSubscription.ts and
// lib/controllersStatusSubscription.ts); everything else is REST.
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
      network: Network.create(fetchQuery),
    });
  }

  // Client-side: reuse singleton
  if (!clientEnvironment) {
    clientEnvironment = new Environment({
      store: new Store(new RecordSource(records)),
      network: Network.create(fetchQuery),
    });
  }

  return clientEnvironment;
};
