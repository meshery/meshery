import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const addonEventsSubscription = graphql`
  subscription AddonEventsSubscription($selector: MeshType) {
    listenToAddonEvents(selector: $selector) {
      type
      status
    }
  }
`;

export default function subscribeOperatorEvents(updater, variables) {
  requestSubscription(environment, {
    subscription: addonEventsSubscription,
    variables: {
      selector: variables.serviceMesh
    },
    updater: updater,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
