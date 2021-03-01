import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const meshScanSubscription = graphql`
  subscription meshScanSubscription($filter: ControlPlaneFilter) {
    listenToControlPlaneEvents(filter: $filter) {
        name
        version
        members {
          component
          status
          namespace
        }
      }
  }
`;

export default function subscribeMeshScanEvents(dataCB, variables) {
  requestSubscription(environment, {
    subscription: meshScanSubscription,
    variables: {
      filter: variables
    },
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
