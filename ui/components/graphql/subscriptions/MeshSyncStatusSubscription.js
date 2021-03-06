import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const meshSyncStatusSubscription = graphql`
subscription MeshSyncStatusSubscription {
  listenToMeshSyncEvents {
    name
    status
    error {
      code
      description
    }
  }
}
`;

export default function subscribeMeshSyncStatusEvents(dataCB) {
  requestSubscription(environment, {
    subscription: meshSyncStatusSubscription,
    variables: {},
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
  