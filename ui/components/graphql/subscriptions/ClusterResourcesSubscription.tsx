import { graphql, requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export const clusterResourcesSubscription = graphql`
  subscription ClusterResourcesSubscription($k8scontextIDs: [String!], $namespace: String!) {
    clusterResources: subscribeClusterResources(
      k8scontextIDs: $k8scontextIDs
      namespace: $namespace
    ) {
      resources {
        kind
        count
      }
    }
  }
`;

interface ClusterResourcesVariables {
  k8scontextIDs?: string[];
  namespace: string;
}

export default function subscribeClusterResources(
  dataCB: (data: unknown) => void,
  variables: ClusterResourcesVariables,
) {
  const environment = createRelayEnvironment({});
  return requestSubscription(environment, {
    subscription: clusterResourcesSubscription,
    variables: variables,
    onNext: dataCB,
    onError: (error) => console.log(`Cluster Resources Subscription error:`, error),
  });
}
