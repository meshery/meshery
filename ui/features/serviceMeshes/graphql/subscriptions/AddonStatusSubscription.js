import environment from "@/app/relayEnvironment";
import { graphql, requestSubscription } from "react-relay";

const addonStatusSubscription = graphql`
  subscription AddonStatusSubscription($selector: MeshType) {
    addonsState: listenToAddonState(selector: $selector) {
      name
      owner
      endpoint
    }
  }
`;

export default function subscribeAddonStatusEvents(dataCB, variables) {
  requestSubscription(environment, {
    subscription: addonStatusSubscription,
    variables: { selector: variables.serviceMesh },
    onNext: dataCB,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
