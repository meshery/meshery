import { graphql, requestSubscription } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const addonStatusSubscription = graphql`
  subscription AddonStatusSubscription($selector: MeshType) {
    addonEvent: listenToAddonState(selector: $selector) {
      type
      status
      config {
        serviceName
        endpoint
      }
    }
  }
`;

export default function subscribeAddonStatusEvents(dataCB, variables) {
  requestSubscription(environment, {
    subscription: addonStatusSubscription,
    variables: {
      selector: variables.serviceMesh
    },
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
