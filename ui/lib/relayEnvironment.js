import {
  Environment,
  Network,
  RecordSource,
  Store,
} from 'relay-runtime';

import { SubscriptionClient } from 'subscriptions-transport-ws'
  
function fetchQuery(
  operation,
  variables,
) {
  return fetch('http://localhost:9081/api/system/graphql/query', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      query: operation.text,
      variables,
    }),
  }).then(response => {
    return response.json();
  });
}

const setupSubscription = (config, variables, cacheConfig, observer) => {
  const query = config.text

  const subscriptionClient = new SubscriptionClient('wss://localhost:9081/api/system/graphql/query', {reconnect: true})
  subscriptionClient.subscribe({query, variables}, (error, result) => {
    observer.onNext({data: result})
  })
}
  
const environment = new Environment({
  network: Network.create(fetchQuery, setupSubscription),
  store: new Store(new RecordSource()),
});
  
export default environment;