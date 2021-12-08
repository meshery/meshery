import { graphql, fetchQuery } from 'react-relay';
import environment from '../../../lib/relayEnvironment';

export default function deployMeshSync() {

  const query = graphql`
    query DeployMeshSyncQuery {
    deployMeshsync 
    }
  `;

  return fetchQuery(environment, query);
}
