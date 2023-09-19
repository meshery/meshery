import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

/**
 *
 * @param {{k8scontextID: string}} variable
 * @returns
 */
export default function deployMeshSync(variable) {
  const environment = createRelayEnvironment({});

  const DeployMeshSyncQueryNode = graphql`
    query DeployMeshSyncQuery($k8scontextID: String!) {
      deployMeshsync(k8scontextID: $k8scontextID)
    }
  `;

  return fetchQuery(environment, DeployMeshSyncQueryNode, variable);
}
