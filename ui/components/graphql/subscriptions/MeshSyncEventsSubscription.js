import { graphql, requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

const meshSyncEventsSubscription = graphql`
 subscription MeshSyncEventsSubscription($k8scontextIDs: [String!], $eventTypes: [MeshSyncEventType!]) {
    meshsyncevents: subscribeMeshSyncEvents(k8scontextIDs: $k8scontextIDs, eventTypes: $eventTypes) {
         type
    		 object
    		 contextId
        }
      }
`;

export default function subscribeMeshSyncEvents(dataCB, variables) {
  const environment = createRelayEnvironment({});
  return requestSubscription(environment, {
    subscription : meshSyncEventsSubscription,
    variables : variables,
    onNext : dataCB,
    onError : (error) => console.log(`An error occured:`, error),
  });
}
