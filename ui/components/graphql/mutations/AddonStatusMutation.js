import { graphql, commitMutation } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const addonStatusMutation = graphql`
  mutation AddonStatusMutation($input: AddonStatusInput) {
    addonstate: changeAddonStatus(input: $input)
  }
`;

export default function changeAddonStatus(onComplete, variables) {
  const vars = {
    input: {
      selector: variables.serviceMesh,
      targetStatus:variables.status
    }
  };

  commitMutation(environment,{
    mutation: addonStatusMutation,
    variables: vars,
    onCompleted: onComplete,
    onError: error => console.log(`An error occured:`, error),
  });
}
