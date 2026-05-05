import { graphql } from 'react-relay';
import { createSubscription } from '../../../lib/subscriptionHelper';

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
  dataCB: (_data: unknown) => void,
  variables: ClusterResourcesVariables,
) {
  return createSubscription({
    subscription: clusterResourcesSubscription,
    variables,
    onNext: dataCB,
    subscriptionName: 'ClusterResources',
  });
}
