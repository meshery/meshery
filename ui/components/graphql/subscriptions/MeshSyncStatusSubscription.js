import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const meshSyncStatusSubscription = graphql`
subscription MeshSyncStatusSubscription($k8scontextID: String!) {
  listenToMeshSyncEvents(k8scontextID: $k8scontextID) {
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
    subscription : meshSyncStatusSubscription,
    variables : {},
    onNext : dataCB,
    onError : (error) => console.log(`An error occured:`, error),
  });
}
