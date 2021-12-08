import { graphql, fetchQuery } from 'react-relay';
import environment from '../../../lib/relayEnvironment';

export default function connectToNATS() {

  const query = graphql`
    query DeployNatsQuery {
      connectToNats 
    }
  `;

  return fetchQuery(environment, query);
}
