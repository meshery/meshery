import { Environment, Network, RecordSource, Store } from "relay-runtime";
import { promisifiedDataFetch } from "./data-fetch";
import { SubscriptionClient } from "subscriptions-transport-ws";

function fetchQuery(operation, variables) {
  return promisifiedDataFetch("/api/system/graphql/query", {
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    method: "POST",
    body: JSON.stringify({
      query: operation.text,
      variables,
    }),
  });
}

function setupSubscription(config, variables, cacheConfig, observer) {
  const query = config.text;

  const subscriptionClient = new SubscriptionClient("/api/system/graphql/query", {
    reconnect: true,
  });
  subscriptionClient.subscribe({ query, variables }, (error, result) => {
    observer.onNext({ data: result });
  });
}

const environment = new Environment({
  network: Network.create(fetchQuery, setupSubscription),
  store: new Store(new RecordSource()),
});

export default environment;
