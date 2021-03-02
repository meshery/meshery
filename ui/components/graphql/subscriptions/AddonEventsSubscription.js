import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const addonEventsSubscription = graphql`
  subscription AddonEventsSubscription($selector: MeshType) {
    addonEvent: listenToAddonState(selector: $selector) {
      type
      status
      config {
        serviceName
      }
    }
  }
`;

export default function subscribeOperatorEvents(dataCB, variables) {
  requestSubscription(environment, {
    subscription: addonEventsSubscription,
    variables: {
      selector: variables.serviceMesh
    },
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
