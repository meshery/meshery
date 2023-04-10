import { graphql, commitMutation } from "react-relay";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";

export default function changeAdaptorState(onComplete, variables) {
  const environment = createRelayEnvironment({});
  const vars = { input : { targetStatus : variables.status, targetPort : variables.targetPort, adapter : variables.adapter } };

  const adaptorStatusMutation = graphql`
  mutation AdaptorStatusMutation($input: AdaptorStatusInput) {
    adaptorStatus: changeAdaptorStatus(input: $input) 
  }
`;

  commitMutation(environment,{
    mutation : adaptorStatusMutation,
    variables : vars,
    onCompleted : onComplete,
    onError : error => console.log(`An error occured:`, error),
  });
}