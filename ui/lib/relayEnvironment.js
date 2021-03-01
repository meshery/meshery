import { Environment, Network, Observable, RecordSource, Store } from "relay-runtime";
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

  const subscriptionClient = new SubscriptionClient("ws://localhost:9081/api/system/graphql/query", {
    reconnect: true,
  });

  const subscribeObservable = subscriptionClient.request({ query, variables }, (error, result) => {
    if (error) {
      console.error(error);
      return
    }

    observer.onNext({ data: result });
  });

  return Observable.from(subscribeObservable);
}

const environment = new Environment({
  network: Network.create(fetchQuery, setupSubscription),
  store: new Store(new RecordSource()),
});

export default environment;
