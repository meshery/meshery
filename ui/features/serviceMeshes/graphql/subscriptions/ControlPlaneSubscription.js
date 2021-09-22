import environment from "@/app/relayEnvironment";
import { graphql, requestSubscription } from "react-relay";

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
  requestSubscription(environment, {
    subscription: controlPlaneSubscription,
    variables: { filter: variables },
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
