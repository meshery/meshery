import { graphql, commitMutation } from "react-relay";
import environment from "../../../lib/relayEnvironment";

const operatorStatusMutation = graphql`
  mutation OperatorStatusMutation($targetStatus: Status!) {
    operatorStatus: changeOperatorStatus(targetStatus: $targetStatus) 
  }
`;

export default function changeOperatorState(onComplete, variables) {
  const vars = {
    targetStatus:variables.status
  };

  commitMutation(environment,{
    mutation: operatorStatusMutation,
    variables: vars,
    onCompleted: onComplete,
    onError: error => console.log(`An error occured:`, error),
  });
}
