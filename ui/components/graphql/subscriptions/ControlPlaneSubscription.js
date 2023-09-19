import { graphql, requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export const controlPlaneSubscription = graphql`
  subscription ControlPlaneSubscription($filter: ServiceMeshFilter) {
    controlPlanesState: listenToControlPlaneState(filter: $filter) {
      name
      members {
        name
        version
        component
        namespace
      }
    }
  }
`;

export default function subscribeControlPlaneEvents(dataCB, variables) {
  const environment = createRelayEnvironment({});
  return requestSubscription(environment, {
    subscription: controlPlaneSubscription,
    variables: { filter: variables },
    onNext: dataCB,
    onError: (error) => console.log(`ControlPlane Subscription error:`, error),
  });
}
