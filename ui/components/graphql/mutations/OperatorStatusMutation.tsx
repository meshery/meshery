import { graphql, commitMutation } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

interface OperatorStatusVariables {
  status: string;
  contextID: string;
}

export default function changeOperatorState(
  onComplete: (response: unknown) => void,
  variables: OperatorStatusVariables,
) {
  const environment = createRelayEnvironment({});
  const vars = { input: { targetStatus: variables.status, contextID: variables.contextID } };

  const operatorStatusMutation = graphql`
    mutation OperatorStatusMutation($input: OperatorStatusInput) {
      operatorStatus: changeOperatorStatus(input: $input)
    }
  `;

  commitMutation(environment, {
    mutation: operatorStatusMutation,
    variables: vars,
    onCompleted: onComplete,
    onError: (error) => console.log(`An error occured:`, error),
  });
}
