import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export const clusterResourcesSubscription = graphql`
  subscription ClusterResourcesSubscription($k8scontextIDs: [String!], $namespaces: [String!]) {
    clusterResources: subscribeClusterResources(k8scontextIDs: $k8scontextIDs, namespaces: $namespaces) {
      resources { 
        kind
        count
      }
    }
  }
`;

export default function subscribeClusterResources(dataCB, variables) {
  return requestSubscription(environment, {
    subscription : clusterResourcesSubscription,
    variables : variables,
    onNext : dataCB,
    onError : (error) => console.log(`Cluster Resources Subscription error:`, error),
  });
}
