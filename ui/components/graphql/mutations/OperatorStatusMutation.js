import { graphql, commitMutation } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const operatorStatusMutation = graphql`
  mutation OperatorStatusMutation($input: OperatorStatusInput) {
    operatorStatus: changeOperatorStatus(input: $input) 
  }
`;

export default function changeOperatorState(onComplete, variables) {
  const vars = { input : { targetStatus : variables.status, contextID : variables.contextID } };

  commitMutation(environment,{
    mutation : operatorStatusMutation,
    variables : vars,
    onCompleted : onComplete,
    onError : error => console.log(`An error occured:`, error),
  });
}
