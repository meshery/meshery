import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export const clusterInfoSubscription = graphql`
  subscription ClusterInfoSubscription($k8scontextIDs: [String!]) {
    subscribeClusterInfo(k8scontextIDs: $k8scontextIDs) {
      resources { 
        kind
        number
      }
    }
  }
`;

export default function subscribeClusterInfo(dataCB, variables) {
  return requestSubscription(environment, {
    subscription : clusterInfoSubscription,
    variables : variables,
    onNext : dataCB,
    onError : (error) => console.log(`Cluster Info Subscription error:`, error),
  });
}
