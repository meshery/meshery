import environment from "app/relayEnvironment";
import { graphql, requestSubscription } from "react-relay";

const meshSyncStatusSubscription = graphql`
  subscription MeshSyncStatusSubscription {
    listenToMeshSyncEvents {
      name
      status
      version
      error {
        code
        description
      }
    }
  }
`;

export default function subscribeMeshSyncStatusEvents(dataCB) {
  return requestSubscription(environment, {
    subscription: meshSyncStatusSubscription,
    variables: {},
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
