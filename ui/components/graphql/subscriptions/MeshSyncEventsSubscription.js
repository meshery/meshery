import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const meshSyncEventsSubscription = graphql`
 subscription MeshSyncEventsSubscription($k8scontextIDs: [String!]) {
    subscribeMeshSyncEvents(k8scontextIDs: $k8scontextIDs) {
         type
    		 object
    		 contextId
        }
      }
`;

export default function subcribeControllersStatus(dataCB) {
  return requestSubscription(environment, {
    subscription : meshSyncEventsSubscription,
    variables : { k8scontextIDs : [""] },
    onNext : dataCB,
    onError : (error) => console.log(`An error occured:`, error),
  });
}

