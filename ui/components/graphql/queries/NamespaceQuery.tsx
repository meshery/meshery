import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

interface NamespaceQueryVariables {
  k8sClusterIDs?: string[];
}

export default function fetchAvailableNamespaces(vars: NamespaceQueryVariables) {
  const environment = createRelayEnvironment({});

  const NamespaceQueryNode = graphql`
    query NamespaceQuery($k8sClusterIDs: [String!]) {
      namespaces: getAvailableNamespaces(k8sClusterIDs: $k8sClusterIDs) {
        namespace
      }
    }
  `;

  return fetchQuery(environment, NamespaceQueryNode, vars);
}
