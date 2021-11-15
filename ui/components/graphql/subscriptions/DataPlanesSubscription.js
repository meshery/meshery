import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export const dataPlaneSubscription = graphql`
  subscription DataPlanesSubscription($filter: ServiceMeshFilter) {
    dataPlanesState: listenToDataPlaneState(filter: $filter) {
      name
      proxies {
        controlPlaneMemberName
        containerName
        image
        status {
          containerStatusName
          image
          state
          lastState
          ready
          restartCount
          started
          imageID
          containerID
        }
        ports {
          name
          containerPort
          protocol
        }
        resources
      }
    }
  }
`;

export default function subscribeDataPlaneEvents(dataCB, variables) {
  return requestSubscription(environment, {
    subscription : dataPlaneSubscription,
    variables : {
      filter : variables,
    },
    onNext : dataCB,
    onError : (error) => console.log(`An error occured:`, error),
  });
}
