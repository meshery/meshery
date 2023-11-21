import { graphql, requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

const mesheryControllersStatusSubscription = graphql`
  subscription MesheryControllersStatusSubscription($k8scontextIDs: [String!]) {
    subscribeMesheryControllersStatus(k8scontextIDs: $k8scontextIDs) {
      contextId
      controller
      status
    }
  }
`;

export default function subscribeMesheryControllersStatus(dataCB, variables) {
  const environment = createRelayEnvironment({});
  return requestSubscription(environment, {
    subscription: mesheryControllersStatusSubscription,
    variables: { k8scontextIDs: variables },
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
