import { graphql, commitMutation } from "react-relay";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";

export default function changeAdapterState(onComplete, variables) {
  const environment = createRelayEnvironment({});
  const vars = { input : { targetStatus : variables.status, targetPort : variables.targetPort, adapter : variables.adapter } };

  const adapterStatusMutation = graphql`
  mutation AdapterStatusMutation($input: AdapterStatusInput) {
    adapterStatus: changeAdapterStatus(input: $input) 
  }
`;

  commitMutation(environment,{
    mutation : adapterStatusMutation,
    variables : vars,
    onCompleted : onComplete,
    onError : error => console.log(`An error occured:`, error),
  });
}