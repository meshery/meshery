import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function fetchAvailableNamespaces(vars) {
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
