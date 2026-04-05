import { graphql } from 'react-relay';
import { createSubscription } from '../../../lib/subscriptionHelper';

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

export default function subscribeMesheryControllersStatus(
  dataCB: (_data: unknown) => void,
  variables: string[],
) {
  return createSubscription({
    subscription: mesheryControllersStatusSubscription,
    variables: { connectionIDs: variables },
    onNext: dataCB,
    subscriptionName: 'MesheryControllersStatus',
  });
}
