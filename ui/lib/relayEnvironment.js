import { Environment, Network, Observable, RecordSource, Store } from "relay-runtime";
import { promisifiedDataFetch } from "./data-fetch";
import { SubscriptionClient } from "subscriptions-transport-ws";

function fetchQuery(operation, variables) {
  return promisifiedDataFetch("/api/system/graphql/query", {
    headers : {
      "Content-Type" : "application/json",
    },
    credentials : "include",
    method : "POST",
    body : JSON.stringify({
      query : operation.text,
      variables,
    }),
  });
}

export let subscriptionClient;

if (typeof window !== "undefined"){
  subscriptionClient = new SubscriptionClient("ws://"+window.location.host+"/api/system/graphql/query", {
    reconnect : true,
  });

}

function setupSubscription(config, variables, cacheConfig, observer) {
  const query = config.text;


  const subscribeObservable = subscriptionClient.request({ query, variables }, (error, result) => {
    if (error) {
      console.error(error);
      return
    }

    observer.onNext({ data : result });
  });

  return Observable.from(subscribeObservable);
}

const environment = new Environment({
  network : Network.create(fetchQuery, setupSubscription),
  store : new Store(new RecordSource()),
});

export default environment;
