import { graphql, requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

const meshSyncEventsSubscription = graphql`
  subscription MeshSyncEventsSubscription(
    $connectionIDs: [String!]
    $eventTypes: [MeshSyncEventType!]
  ) {
    meshsyncevents: subscribeMeshSyncEvents(
      connectionIDs: $connectionIDs
      eventTypes: $eventTypes
    ) {
      type
      object
      connectionID
    }
  }
`;

export default function subscribeMeshSyncEvents(dataCB, variables) {
  const environment = createRelayEnvironment({});
  return requestSubscription(environment, {
    subscription: meshSyncEventsSubscription,
    variables: variables,
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
