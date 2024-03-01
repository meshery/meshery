import { graphql, requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

const mesheryControllersStatusSubscription = graphql`
  subscription MesheryControllersStatusSubscription($connectionIDs: [String!]) {
    subscribeMesheryControllersStatus(connectionIDs: $connectionIDs) {
      connectionID
      controller
      status
      version
    }
  }
`;

export default function subscribeMesheryControllersStatus(dataCB, variables) {
  const environment = createRelayEnvironment({});
  return requestSubscription(environment, {
    subscription: mesheryControllersStatusSubscription,
    variables: { connectionIDs: variables },
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
