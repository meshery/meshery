import { graphql, commitMutation } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function changeOperatorState(onComplete, variables) {
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
