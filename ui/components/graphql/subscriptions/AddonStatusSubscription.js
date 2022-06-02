import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const addonStatusSubscription = graphql`
  subscription AddonStatusSubscription($filter: ServiceMeshFilter) {
    addonsState: listenToAddonState(filter: $filter) {
      name
      owner
    }
  }
`;

export default function subscribeAddonStatusEvents(dataCB, variables) { // todo: changes in variable
  requestSubscription(environment, {
    subscription : addonStatusSubscription,
    variables : { filter : variables },
    onNext : dataCB,
    onError : (error) => console.log(`An error occured:`, error),
  });
}
