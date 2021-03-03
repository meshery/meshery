import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const controlPlaneSubscription = graphql`
  subscription ControlPlaneSubscription($filter: ControlPlaneFilter) {
    listenToControlPlaneState(filter: $filter) {
        name
        members {
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
    variables: {
      filter: variables
    },
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
