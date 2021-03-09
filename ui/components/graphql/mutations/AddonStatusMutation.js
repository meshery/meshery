import { graphql, commitMutation } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const addonStatusMutation = graphql`
  mutation AddonStatusMutation($selector: MeshType, $targetStatus: Status!) {
    addonstate: changeAddonStatus(selector: $selector, targetStatus: $targetStatus)
  }
`;

export default function changeAddonStatus(onComplete, variables) {
  const vars = {
    selector: variables.serviceMesh,
    targetStatus:variables.status
  };

  commitMutation(environment,{
    mutation: addonStatusMutation,
    variables: vars,
    onCompleted: onComplete,
    onError: error => console.log(`An error occured:`, error),
  });
}
