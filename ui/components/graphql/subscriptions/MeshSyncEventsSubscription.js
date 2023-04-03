import { graphql, requestSubscription } from "react-relay";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";

const meshSyncEventsSubscription = graphql`
 subscription MeshSyncEventsSubscription($k8scontextIDs: [String!]) {
    meshsyncevents: subscribeMeshSyncEvents(k8scontextIDs: $k8scontextIDs) {
         type
    		 object
    		 contextId
        }
      }
`;

export default function subscribeMeshSyncEvents(dataCB) {
  const environment = createRelayEnvironment({});
  return requestSubscription(environment, {
    subscription : meshSyncEventsSubscription,
    variables : { k8scontextIDs : [""] },
    onNext : dataCB,
    onError : (error) => console.log(`An error occured:`, error),
  });
}

